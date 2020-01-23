/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
)

var conditionKeyValues = []string{
	model.Condition_KEY_TAG,
	model.Condition_KEY_COMPUTERNAME,
	model.Condition_KEY_OSNAME,
	model.Condition_KEY_NAME}
var conditionMemberTypeValues = []string{
	model.Condition_MEMBER_TYPE_IPSET,
	model.Condition_MEMBER_TYPE_LOGICALPORT,
	model.Condition_MEMBER_TYPE_LOGICALSWITCH,
	model.Condition_MEMBER_TYPE_SEGMENT,
	model.Condition_MEMBER_TYPE_SEGMENTPORT,
	model.Condition_MEMBER_TYPE_VIRTUALMACHINE,
}
var conditionOperatorValues = []string{
	model.Condition_OPERATOR_CONTAINS,
	model.Condition_OPERATOR_ENDSWITH,
	model.Condition_OPERATOR_EQUALS,
	model.Condition_OPERATOR_NOTEQUALS,
	model.Condition_OPERATOR_STARTSWITH,
}
var conjunctionOperatorValues = []string{
	model.ConjunctionOperator_CONJUNCTION_OPERATOR_OR,
	model.ConjunctionOperator_CONJUNCTION_OPERATOR_AND,
}

func resourceNsxtPolicyGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyGroupCreate,
		Read:   resourceNsxtPolicyGroupRead,
		Update: resourceNsxtPolicyGroupUpdate,
		Delete: resourceNsxtPolicyGroupDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtDomainResourceImporter,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"domain":       getDomainNameSchema(),
			"criteria": {
				Type:        schema.TypeList,
				Description: "Criteria to determine Group membership",
				Elem:        getCriteriaSetSchema(),
				Optional:    true,
			},
			"conjunction": {
				Type:        schema.TypeList,
				Description: "A conjunction applied to 2 sets of criteria.",
				Elem:        getConjunctionSchema(),
				Optional:    true,
			},
		},
	}
}

func getIPAddressExpressionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip_addresses": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of; single IP addresses, IP address ranges or Subnets. Cannot mix IPv4 and IPv6 in a single list",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCidrOrIPOrRange(),
				},
			},
		},
	}
}

func getConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The resource key attribute to apply the condition to.",
				ValidateFunc: validation.StringInSlice(conditionKeyValues, false),
			},
			"member_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The NSX member to apply the condition to. Can be one of; IPSet, LogicalPort, LogicalSwitch, Segment, SegmentPort or VirtualMachine",
				ValidateFunc: validation.StringInSlice(conditionMemberTypeValues, false),
			},
			"operator": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The operator to use for the condition. Can be one of; CONTAINS, ENDSWITH, EQUALS, NOTEQUALS or STARTSWITH",
				ValidateFunc: validation.StringInSlice(conditionOperatorValues, false),
			},
			"value": {
				Type:        schema.TypeString,
				Description: "The value to check for in the condition",
				Required:    true,
			},
		},
	}
}

func getConjunctionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The conjunction operator; either OR or AND",
				ValidateFunc: validation.StringInSlice(conjunctionOperatorValues, false),
			},
		},
	}
}

func getCriteriaSetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ipaddress_expression": {
				Type:        schema.TypeList,
				Description: "An IP Address expression specifying IP Address members in the Group",
				Elem:        getIPAddressExpressionSchema(),
				Optional:    true,
				MaxItems:    1,
			},
			"condition": {
				Type:        schema.TypeList,
				Description: "A Condition querying resources for membership in the Group",
				Elem:        getConditionSchema(),
				Optional:    true,
				MaxItems:    5,
			},
		},
	}
}

func resourceNsxtPolicyGroupExistsInDomain(id string, domain string, connector *client.RestConnector) bool {
	client := domains.NewDefaultGroupsClient(connector)

	_, err := client.Get(domain, id)
	if err == nil {
		return true
	}

	if isNotFoundError(err) {
		return false
	}

	logAPIError("Error retrieving Group", err)
	return false

}

func resourceNsxtPolicyGroupExistsInDomainPartial(domain string) func(id string, connector *client.RestConnector) bool {
	return func(id string, connector *client.RestConnector) bool {
		return resourceNsxtPolicyGroupExistsInDomain(id, domain, connector)
	}
}

func validateNestedGroupConditions(conditions []interface{}) (string, error) {
	memberType := ""
	for _, cond := range conditions {
		condMap := cond.(map[string]interface{})
		condMemberType := condMap["member_type"].(string)
		if memberType != "" && condMemberType != memberType {
			return "", fmt.Errorf("Nested conditions must all use the same member_type, but found '%v' with '%v'", condMemberType, memberType)
		}
		if condMemberType != model.Condition_MEMBER_TYPE_VIRTUALMACHINE && condMap["key"] != model.Condition_KEY_TAG {
			return "", fmt.Errorf("Only Tag can be used for the key of '%v'", condMemberType)
		}
		memberType = condMemberType
	}
	return memberType, nil
}

type criteriaMeta struct {
	ExpressionType string
	MemberType     string
	IsNested       bool
	criteriaBlocks []interface{}
}

func validateGroupCriteriaSets(criteriaSets []interface{}) ([]criteriaMeta, error) {
	var criteria []criteriaMeta
	for _, criteriaBlock := range criteriaSets {
		seenExp := ""
		criteriaMap := criteriaBlock.(map[string]interface{})
		for expName, expVal := range criteriaMap {
			memberType := ""
			expValList := expVal.([]interface{})
			if len(expValList) > 0 {
				if seenExp != "" {
					return nil, fmt.Errorf("Criteria blocks are homogeneous, but found '%v' with '%v'", expName, seenExp)
				}
				if expName == "condition" {
					mType, err := validateNestedGroupConditions(expValList)
					if err != nil {
						return nil, err
					}
					memberType = mType
				} else if expName == "ipaddress_expression" {
					memberType = ""
				} else {
					return nil, fmt.Errorf("Unknown criteria: %v", expName)
				}
				criteriaType := criteriaMeta{
					MemberType:     memberType,
					ExpressionType: expName,
					IsNested:       len(expValList) > 1,
					criteriaBlocks: expValList}
				criteria = append(criteria, criteriaType)
				seenExp = expName
			}
		}
	}

	return criteria, nil
}

func validateGroupConjunctions(conjunctions []interface{}, criteriaMeta []criteriaMeta) error {
	for index, conjunctionIFace := range conjunctions {
		conjunction := conjunctionIFace.(map[string]interface{})
		if conjunction["operator"] == model.ConjunctionOperator_CONJUNCTION_OPERATOR_AND {
			metaA := criteriaMeta[index]
			metaB := criteriaMeta[index+1]
			if metaA.ExpressionType != metaB.ExpressionType {
				return fmt.Errorf("AND conjunctions must use the same types of criteria expressions, but got %v and %v",
					metaA.ExpressionType, metaB.ExpressionType)
			}
			if metaA.MemberType != metaB.MemberType {
				return fmt.Errorf("AND conjunctions with conditions must have the same member types, but got %v and %v",
					metaA.MemberType, metaB.MemberType)
			}
		}
	}
	return nil
}

func buildGroupConditionData(condition interface{}) (*data.StructValue, error) {
	conditionMap := condition.(map[string]interface{})
	conditionModel := model.Condition{
		Key:          conditionMap["key"].(string),
		MemberType:   conditionMap["member_type"].(string),
		Operator:     conditionMap["operator"].(string),
		Value:        conditionMap["value"].(string),
		ResourceType: model.Condition__TYPE_IDENTIFIER,
	}
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)
	dataValue, errors := converter.ConvertToVapi(conditionModel, model.ConditionBindingType())
	if errors != nil {
		return nil, errors[0]
	}
	return dataValue.(*data.StructValue), nil
}

func buildGroupConjunctionData(conjunction string) (*data.StructValue, error) {
	conjunctionStruct := model.ConjunctionOperator{
		ConjunctionOperator: conjunction,
		ResourceType:        model.ConjunctionOperator__TYPE_IDENTIFIER,
	}
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)
	dataValue, errors := converter.ConvertToVapi(conjunctionStruct, model.ConjunctionOperatorBindingType())
	if errors != nil {
		return nil, errors[0]
	}
	return dataValue.(*data.StructValue), nil
}

func buildGroupIPAddressData(ipaddr interface{}) (*data.StructValue, error) {
	ipaddrMap := ipaddr.(map[string]interface{})
	var ipList []string
	for _, ip := range ipaddrMap["ip_addresses"].([]interface{}) {
		ipList = append(ipList, ip.(string))
	}
	ipaddrStruct := model.IPAddressExpression{
		IpAddresses:  ipList,
		ResourceType: model.IPAddressExpression__TYPE_IDENTIFIER,
	}
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)
	dataValue, errors := converter.ConvertToVapi(ipaddrStruct, model.IPAddressExpressionBindingType())
	if errors != nil {
		return nil, errors[0]
	}
	return dataValue.(*data.StructValue), nil
}

func buildNestedGroupExpressionData(expressions []*data.StructValue) (*data.StructValue, error) {
	var completeExpressions []*data.StructValue
	for i := 0; i < len(expressions)-1; i++ {
		completeExpressions = append(completeExpressions, expressions[i])
		// nested conditions have implicit AND
		conjData, err := buildGroupConjunctionData(model.ConjunctionOperator_CONJUNCTION_OPERATOR_AND)
		if err != nil {
			return nil, err
		}
		completeExpressions = append(completeExpressions, conjData)
	}
	completeExpressions = append(completeExpressions, expressions[len(expressions)-1])

	nestedStruct := model.NestedExpression{
		Expressions:  completeExpressions,
		ResourceType: model.NestedExpression__TYPE_IDENTIFIER,
	}
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)
	dataValue, errors := converter.ConvertToVapi(nestedStruct, model.NestedExpressionBindingType())
	if errors != nil {
		return nil, errors[0]
	}
	return dataValue.(*data.StructValue), nil
}

func buildGroupExpressionDataFromType(expressionType string, datum interface{}) (*data.StructValue, error) {
	if expressionType == "condition" {
		data, err := buildGroupConditionData(datum)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else if expressionType == "ipaddress_expression" {
		data, err := buildGroupIPAddressData(datum)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, fmt.Errorf("Unknown expression type: %v", expressionType)
}

func buildGroupExpressionData(criteriaMeta []criteriaMeta, conjunctions []interface{}) ([]*data.StructValue, error) {
	var expressionData []*data.StructValue
	for index, meta := range criteriaMeta {
		if meta.IsNested {
			var nestedExpressionData []*data.StructValue
			for _, nestedMeta := range meta.criteriaBlocks {
				nestedData, err := buildGroupExpressionDataFromType(meta.ExpressionType, nestedMeta)
				if err != nil {
					return nil, err
				}
				nestedExpressionData = append(nestedExpressionData, nestedData)
			}
			nested, err := buildNestedGroupExpressionData(nestedExpressionData)
			if err != nil {
				return nil, err
			}
			expressionData = append(expressionData, nested)
		} else {
			data, err := buildGroupExpressionDataFromType(meta.ExpressionType, meta.criteriaBlocks[0])
			if err != nil {
				return nil, err
			}
			expressionData = append(expressionData, data)
		}
		if index < len(conjunctions) {
			conjMap := conjunctions[index].(map[string]interface{})
			data, err := buildGroupConjunctionData(conjMap["operator"].(string))
			if err != nil {
				return nil, err
			}
			expressionData = append(expressionData, data)
		}
	}
	return expressionData, nil
}

func groupConditionDataToMap(expData *data.StructValue) (map[string]interface{}, error) {
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)
	condData, errors := converter.ConvertToGolang(expData, model.ConditionBindingType())
	if len(errors) > 0 {
		return nil, errors[0]
	}
	condStruct := condData.(model.Condition)
	var condMap = make(map[string]interface{})
	condMap["key"] = condStruct.Key
	condMap["member_type"] = condStruct.MemberType
	condMap["operator"] = condStruct.Operator
	condMap["value"] = condStruct.Value
	return condMap, nil
}

func fromGroupExpressionData(expressions []*data.StructValue) ([]interface{}, []interface{}, error) {
	var parsedConjunctions []interface{}
	var parsedCriteria []interface{}
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

	for _, expression := range expressions {
		expData, errs := converter.ConvertToGolang(expression, model.ExpressionBindingType())
		if len(errs) > 0 {
			return nil, nil, errs[0]
		}
		expStruct := expData.(model.Expression)

		if expStruct.ResourceType == model.Expression_RESOURCE_TYPE_CONJUNCTIONOPERATOR {
			conjData, errors := converter.ConvertToGolang(expression, model.ConjunctionOperatorBindingType())
			if len(errors) > 0 {
				return nil, nil, errors[0]
			}
			conjStruct := conjData.(model.ConjunctionOperator)
			var conjMap = make(map[string]string)
			conjMap["operator"] = conjStruct.ConjunctionOperator
			parsedConjunctions = append(parsedConjunctions, conjMap)
		} else if expStruct.ResourceType == model.IPAddressExpression__TYPE_IDENTIFIER {
			ipData, errors := converter.ConvertToGolang(expression, model.IPAddressExpressionBindingType())
			if len(errors) > 0 {
				return nil, nil, errors[0]
			}
			ipStruct := ipData.(model.IPAddressExpression)
			var addrMap = make(map[string][]string)
			addrMap["ip_addresses"] = ipStruct.IpAddresses
			var ipMap = make(map[string]interface{})
			ipMap["ipaddress_expression"] = addrMap
			parsedCriteria = append(parsedCriteria, ipMap)
		} else if expStruct.ResourceType == model.Condition__TYPE_IDENTIFIER {
			condMap, err := groupConditionDataToMap(expression)
			if err != nil {
				return nil, nil, err
			}
			var criteriaMap = make(map[string]interface{})
			criteriaMap["condition"] = condMap
			parsedCriteria = append(parsedCriteria, criteriaMap)
		} else if expStruct.ResourceType == model.NestedExpression__TYPE_IDENTIFIER {
			nestedData, errors := converter.ConvertToGolang(expression, model.NestedExpressionBindingType())
			if len(errors) > 0 {
				return nil, nil, errors[0]
			}
			var criteriaList []interface{}
			nestedExp := nestedData.(model.NestedExpression)
			for _, nestedExpression := range nestedExp.Expressions {
				expData, errors := converter.ConvertToGolang(nestedExpression, model.ExpressionBindingType())
				if len(errors) > 0 {
					return nil, nil, errors[0]
				}
				eStruct := expData.(model.Expression)
				if eStruct.ResourceType == model.Condition__TYPE_IDENTIFIER {
					condMap, err := groupConditionDataToMap(nestedExpression)
					if err != nil {
						return nil, nil, err
					}
					criteriaList = append(criteriaList, condMap)
				}
			}
			criteriaMap := make(map[string]interface{})
			criteriaMap["condition"] = criteriaList
			parsedCriteria = append(parsedCriteria, criteriaMap)
		} else {
			return nil, nil, fmt.Errorf("Unsupported criteria type: %v", expStruct.ResourceType)
		}
	}
	return parsedCriteria, parsedConjunctions, nil
}

func validateGroupCriteriaAndConjunctions(criteriaSets []interface{}, conjunctions []interface{}) ([]criteriaMeta, error) {
	if len(criteriaSets)+len(conjunctions) == 0 {
		return nil, nil
	}
	if (len(criteriaSets)+len(conjunctions))%2 == 0 {
		if len(conjunctions) < len(criteriaSets)-1 {
			return nil, fmt.Errorf("Missing conjunction for criteria")
		}
		return nil, fmt.Errorf("Missing criteria for last conjunction")
	}
	criteriaMeta, err := validateGroupCriteriaSets(criteriaSets)
	if err != nil {
		return nil, err
	}

	err = validateGroupConjunctions(conjunctions, criteriaMeta)
	if err != nil {
		return nil, err
	}
	return criteriaMeta, nil
}

func resourceNsxtPolicyGroupCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := domains.NewDefaultGroupsClient(connector)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, connector, resourceNsxtPolicyGroupExistsInDomainPartial(d.Get("domain").(string)))
	if err != nil {
		return err
	}

	criteriaSets := d.Get("criteria").([]interface{})
	conjunctions := d.Get("conjunction").([]interface{})

	criteriaMeta, err := validateGroupCriteriaAndConjunctions(criteriaSets, conjunctions)
	if err != nil {
		return err
	}

	expressionData, err := buildGroupExpressionData(criteriaMeta, conjunctions)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.Group{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Expression:  expressionData,
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating Group with ID %s", id)
	err = client.Patch(d.Get("domain").(string), id, obj)
	if err != nil {
		return handleCreateError("Group", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyGroupRead(d, m)
}

func resourceNsxtPolicyGroupRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := domains.NewDefaultGroupsClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Group ID")
	}

	domainName := d.Get("domain").(string)
	obj, err := client.Get(domainName, id)
	if err != nil {
		return handleReadError(d, "Group", id, err)
	}

	criteria, conditions, err := fromGroupExpressionData(obj.Expression)
	if err != nil {
		return err
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("domain", getDomainFromResourcePath(*obj.Path))
	d.Set("revision", obj.Revision)
	d.Set("criteria", criteria)
	d.Set("conjunction", conditions)

	return nil
}

func resourceNsxtPolicyGroupUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := domains.NewDefaultGroupsClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Group ID")
	}

	criteriaSets := d.Get("criteria").([]interface{})
	conjunctions := d.Get("conjunction").([]interface{})

	criteriaMeta, err := validateGroupCriteriaAndConjunctions(criteriaSets, conjunctions)
	if err != nil {
		return err
	}

	expressionData, err := buildGroupExpressionData(criteriaMeta, conjunctions)
	if err != nil {
		return err
	}

	// Read the rest of the configured parameters
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.Group{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Expression:  expressionData,
	}

	// Update the resource using PATCH
	err = client.Patch(d.Get("domain").(string), id, obj)
	if err != nil {
		return handleUpdateError("Group", id, err)
	}

	return resourceNsxtPolicyGroupRead(d, m)
}

func resourceNsxtPolicyGroupDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Group ID")
	}

	connector := getPolicyConnector(m)
	client := domains.NewDefaultGroupsClient(connector)

	forceDelete := true
	failIfSubtreeExists := false
	err := client.Delete(d.Get("domain").(string), id, &failIfSubtreeExists, &forceDelete)
	if err != nil {
		return handleDeleteError("Group", id, err)
	}

	return nil
}
