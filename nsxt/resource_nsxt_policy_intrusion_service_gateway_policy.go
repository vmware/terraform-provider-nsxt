// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliIntrusionServiceGatewayPoliciesClient = domains.NewIntrusionServiceGatewayPoliciesClient

func resourceNsxtPolicyIntrusionServiceGatewayPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIntrusionServiceGatewayPolicyCreate,
		Read:   resourceNsxtPolicyIntrusionServiceGatewayPolicyRead,
		Update: resourceNsxtPolicyIntrusionServiceGatewayPolicyUpdate,
		Delete: resourceNsxtPolicyIntrusionServiceGatewayPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtDomainResourceImporter,
		},
		Schema: getIntrusionServiceGatewayPolicySchema(true),
	}
}

func getIntrusionServiceGatewayPolicySchema(withRule bool) map[string]*schema.Schema {
	result := map[string]*schema.Schema{
		"nsx_id":       getNsxIDSchema(),
		"path":         getPathSchema(),
		"display_name": getDisplayNameSchema(),
		"description":  getDescriptionSchema(),
		"revision":     getRevisionSchema(),
		"tag":          getTagsSchema(),
		"domain":       getDomainNameSchema(),
		"category": {
			Type:         schema.TypeString,
			Description:  "Category",
			ValidateFunc: validation.StringInSlice(gatewayPolicyCategoryValues, false),
			Required:     true,
			ForceNew:     true,
		},
		"comments": {
			Type:        schema.TypeString,
			Description: "Comments for security policy lock/unlock",
			Optional:    true,
		},
		"locked": {
			Type:        schema.TypeBool,
			Description: "Indicates whether a security policy should be locked",
			Optional:    true,
			Default:     false,
		},
		"sequence_number": {
			Type:        schema.TypeInt,
			Description: "This field is used to resolve conflicts between security policies across domains",
			Optional:    true,
			Default:     0,
		},
		"stateful": {
			Type:        schema.TypeBool,
			Description: "When it is stateful, the state of the network connects are tracked and a stateful packet inspection is performed",
			Optional:    true,
			Default:     true,
		},
		"tcp_strict": {
			Type:        schema.TypeBool,
			Description: "Ensures that a 3 way TCP handshake is done before the data packets are sent",
			Optional:    true,
			Computed:    true,
		},
	}

	if withRule {
		result["rule"] = getIntrusionServiceGatewayPolicyRulesSchema()
	}

	return result
}

func getIntrusionServiceGatewayPolicyRulesSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of Intrusion Service rules in the policy",
		Optional:    true,
		MaxItems:    1000,
		Elem: &schema.Resource{
			Schema: getIntrusionServiceGatewayPolicyRuleSchema(true),
		},
	}
}

// getIntrusionServiceGatewayPolicyRuleSchema returns the rule schema for both embedded
// (embedded=true) and standalone (embedded=false) use cases for IDPS gateway policies.
// When embedded is true, nsx_id is Computed-only and sequence_number is Optional+Computed
// (rules are managed inline with the policy).
// When embedded is false, nsx_id is Optional+Computed and sequence_number is Required,
// with policy_path added for the standalone rule resource.
func getIntrusionServiceGatewayPolicyRuleSchema(embedded bool) map[string]*schema.Schema {
	ruleSchema := map[string]*schema.Schema{
		"nsx_id":       getFlexNsxIDSchema(embedded),
		"display_name": getDisplayNameSchema(),
		"description":  getDescriptionSchema(),
		"path":         getPathSchema(),
		"revision":     getRevisionSchema(),
		"destination_groups": {
			Type:        schema.TypeSet,
			Description: "List of destination groups",
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validatePolicySourceDestinationGroups(),
			},
			Optional: true,
		},
		"destinations_excluded": {
			Type:        schema.TypeBool,
			Description: "Negation of destination groups",
			Optional:    true,
			Default:     false,
		},
		"direction": {
			Type:         schema.TypeString,
			Description:  "Traffic direction",
			Optional:     true,
			ValidateFunc: validation.StringInSlice(securityPolicyDirectionValues, false),
			Default:      model.Rule_DIRECTION_IN_OUT,
		},
		"disabled": {
			Type:        schema.TypeBool,
			Description: "Flag to disable the rule",
			Optional:    true,
			Default:     false,
		},
		"ip_version": {
			Type:         schema.TypeString,
			Description:  "IP version",
			Optional:     true,
			ValidateFunc: validation.StringInSlice(securityPolicyIPProtocolValues, false),
			Default:      model.Rule_IP_PROTOCOL_IPV4_IPV6,
		},
		"logged": {
			Type:        schema.TypeBool,
			Description: "Flag to enable packet logging",
			Optional:    true,
			Default:     false,
		},
		"notes": {
			Type:        schema.TypeString,
			Description: "Text for additional notes on changes",
			Optional:    true,
		},
		"scope": {
			Type:        schema.TypeSet,
			Description: "List of policy paths where the rule is applied",
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validatePolicyPath(),
			},
			Required: true,
		},
		"services": {
			Type:        schema.TypeSet,
			Description: "List of services to match",
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validatePolicyPath(),
			},
			Optional: true,
		},
		"source_groups": {
			Type:        schema.TypeSet,
			Description: "List of source groups",
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validatePolicySourceDestinationGroups(),
			},
			Optional: true,
		},
		"sources_excluded": {
			Type:        schema.TypeBool,
			Description: "Negation of source groups",
			Optional:    true,
			Default:     false,
		},
		"tag": getTagsSchema(),
		"log_label": {
			Type:        schema.TypeString,
			Description: "Additional information (string) which will be propagated to the rule syslog",
			Optional:    true,
		},
		"action": {
			Type:         schema.TypeString,
			Description:  "Action",
			Optional:     true,
			ValidateFunc: validation.StringInSlice(policyIntrusionServiceRuleActionValues, false),
			Default:      model.IdsRule_ACTION_DETECT,
		},
		"ids_profiles": getIdsProfilesSchema(),
	}

	if embedded {
		ruleSchema["sequence_number"] = &schema.Schema{
			Type:        schema.TypeInt,
			Description: "Sequence number of this rule",
			Optional:    true,
			Computed:    true,
		}
	} else {
		ruleSchema["sequence_number"] = &schema.Schema{
			Type:        schema.TypeInt,
			Description: "Sequence number of this rule",
			Required:    true,
		}
		ruleSchema["policy_path"] = getPolicyPathSchema(true, true, "Intrusion Service Gateway Policy path")
	}

	return ruleSchema
}

func resourceNsxtPolicyIntrusionServiceGatewayPolicyExistsInDomain(sessionContext utl.SessionContext, id string, domainName string, connector client.Connector) (bool, error) {
	client := cliIntrusionServiceGatewayPoliciesClient(sessionContext, connector)
	if client == nil {
		return false, policyResourceNotSupportedError()
	}
	_, err := client.Get(domainName, id)

	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving Intrusion Service Gateway Policy", err)
}

func resourceNsxtPolicyIntrusionServiceGatewayPolicyExistsPartial(sessionContext utl.SessionContext, domainName string) func(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	return func(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
		return resourceNsxtPolicyIntrusionServiceGatewayPolicyExistsInDomain(sessionContext, id, domainName, connector)
	}
}

func createChildDomainWithIntrusionServiceGatewayPolicy(domain string, policyID string, policy model.IdsGatewayPolicy) (*data.StructValue, error) {
	converter := bindings.NewTypeConverter()

	childPolicy := model.ChildIdsGatewayPolicy{
		Id:               &policyID,
		ResourceType:     "ChildIdsGatewayPolicy",
		IdsGatewayPolicy: &policy,
	}

	dataValue, errors := converter.ConvertToVapi(childPolicy, model.ChildIdsGatewayPolicyBindingType())
	if len(errors) > 0 {
		return nil, errors[0]
	}

	var domainChildren []*data.StructValue
	domainChildren = append(domainChildren, dataValue.(*data.StructValue))

	targetType := "Domain"
	childDomain := model.ChildResourceReference{
		Id:           &domain,
		ResourceType: "ChildResourceReference",
		TargetType:   &targetType,
		Children:     domainChildren,
	}

	dataValue, errors = converter.ConvertToVapi(childDomain, model.ChildResourceReferenceBindingType())
	if len(errors) > 0 {
		return nil, errors[0]
	}
	return dataValue.(*data.StructValue), nil
}

func intrusionServiceGatewayPolicyInfraPatch(context utl.SessionContext, policy model.IdsGatewayPolicy, domain string, m interface{}) error {
	childDomain, err := createChildDomainWithIntrusionServiceGatewayPolicy(domain, *policy.Id, policy)
	if err != nil {
		return fmt.Errorf("Failed to create H-API for Intrusion Service Gateway Policy: %s", err)
	}

	var infraChildren []*data.StructValue
	infraChildren = append(infraChildren, childDomain)

	infraType := "Infra"
	infraObj := model.Infra{
		Children:     infraChildren,
		ResourceType: &infraType,
	}

	return policyInfraPatch(context, infraObj, getPolicyConnector(m), false)
}

func updateIntrusionServiceGatewayPolicy(id string, d *schema.ResourceData, m interface{}) error {

	domain := d.Get("domain").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	comments := d.Get("comments").(string)
	locked := d.Get("locked").(bool)
	sequenceNumber := int64(d.Get("sequence_number").(int))
	stateful := d.Get("stateful").(bool)
	category := d.Get("category").(string)
	resourceType := "IdsGatewayPolicy"

	obj := model.IdsGatewayPolicy{
		Id:             &id,
		DisplayName:    &displayName,
		Description:    &description,
		Tags:           tags,
		Comments:       &comments,
		Locked:         &locked,
		SequenceNumber: &sequenceNumber,
		Stateful:       &stateful,
		Category:       &category,
		ResourceType:   &resourceType,
	}

	_, isSet := d.GetOkExists("tcp_strict")
	if isSet {
		tcpStrict := d.Get("tcp_strict").(bool)
		obj.TcpStrict = &tcpStrict
	}

	if len(d.Id()) > 0 {
		revision := int64(d.Get("revision").(int))
		obj.Revision = &revision
	}

	var childRules []*data.StructValue
	if d.HasChange("rule") {
		oldRules, _ := d.GetChange("rule")
		rules := getPolicyIdsRulesFromSchema(d)

		existingRules := make(map[string]bool)
		for _, rule := range rules {
			ruleID := newUUID()
			if rule.Id != nil {
				ruleID = *rule.Id
				existingRules[ruleID] = true
			} else {
				rule.Id = &ruleID
			}

			childRule, err := createPolicyChildIdsRule(ruleID, rule, false)
			if err != nil {
				return err
			}
			log.Printf("[DEBUG]: Adding child Intrusion Service gateway rule with id %s", ruleID)
			childRules = append(childRules, childRule)
		}

		for _, oldRule := range oldRules.([]interface{}) {
			oldRuleMap := oldRule.(map[string]interface{})
			oldRuleID := oldRuleMap["nsx_id"].(string)
			if _, exists := existingRules[oldRuleID]; !exists {
				resourceType := "IdsRule"
				rule := model.IdsRule{
					Id:           &oldRuleID,
					ResourceType: &resourceType,
				}

				childRule, err := createPolicyChildIdsRule(oldRuleID, rule, true)
				if err != nil {
					return err
				}
				log.Printf("[DEBUG]: Deleting child Intrusion Service gateway rule with id %s", oldRuleID)
				childRules = append(childRules, childRule)

			}
		}
	}

	log.Printf("[DEBUG]: Updating Intrusion Service Gateway policy %s with %d child rules", id, len(childRules))
	if len(childRules) > 0 {
		obj.Children = childRules
	}

	return intrusionServiceGatewayPolicyInfraPatch(getSessionContext(d, m), obj, domain, m)
}

func resourceNsxtPolicyIntrusionServiceGatewayPolicyCreate(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyIntrusionServiceGatewayPolicyGeneralCreate(d, m, true)
}

func resourceNsxtPolicyIntrusionServiceGatewayPolicyGeneralCreate(d *schema.ResourceData, m interface{}, withRule bool) error {
	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyIntrusionServiceGatewayPolicyExistsPartial(getSessionContext(d, m), d.Get("domain").(string)))
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Intrusion Service Gateway Policy with ID %s", id)
	err = updateIntrusionServiceGatewayPolicy(id, d, m)

	if err != nil {
		return handleCreateError("Intrusion Service Gateway Policy", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyIntrusionServiceGatewayPolicyGeneralRead(d, m, withRule)
}

func resourceNsxtPolicyIntrusionServiceGatewayPolicyRead(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyIntrusionServiceGatewayPolicyGeneralRead(d, m, true)
}

func resourceNsxtPolicyIntrusionServiceGatewayPolicyGeneralRead(d *schema.ResourceData, m interface{}, withRule bool) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	domainName := d.Get("domain").(string)
	if id == "" {
		return fmt.Errorf("Error obtaining Intrusion Service Gateway Policy id")
	}
	client := cliIntrusionServiceGatewayPoliciesClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}
	obj, err := client.Get(domainName, id)
	if err != nil {
		return handleReadError(d, "Intrusion Service Gateway Policy", id, err)
	}
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("domain", getDomainFromResourcePath(*obj.Path))
	d.Set("category", obj.Category)
	d.Set("comments", obj.Comments)
	d.Set("locked", obj.Locked)
	d.Set("sequence_number", obj.SequenceNumber)
	d.Set("stateful", obj.Stateful)
	if obj.TcpStrict != nil {
		d.Set("tcp_strict", *obj.TcpStrict)
	}
	d.Set("revision", obj.Revision)
	if withRule {
		return setPolicyIdsRulesInSchema(d, obj.Rules)
	}
	return nil
}

func resourceNsxtPolicyIntrusionServiceGatewayPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyIntrusionServiceGatewayPolicyGeneralUpdate(d, m, true)
}

func resourceNsxtPolicyIntrusionServiceGatewayPolicyGeneralUpdate(d *schema.ResourceData, m interface{}, withRule bool) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Intrusion Service Gateway Policy id")
	}

	log.Printf("[INFO] Updating Intrusion Service Gateway Policy with ID %s", id)
	err := updateIntrusionServiceGatewayPolicy(id, d, m)

	if err != nil {
		return handleUpdateError("Intrusion Service Gateway Policy", id, err)
	}

	return resourceNsxtPolicyIntrusionServiceGatewayPolicyGeneralRead(d, m, withRule)
}

func resourceNsxtPolicyIntrusionServiceGatewayPolicyDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Intrusion Service Gateway Policy id")
	}

	connector := getPolicyConnector(m)

	client := cliIntrusionServiceGatewayPoliciesClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}
	err := client.Delete(d.Get("domain").(string), id)

	if err != nil {
		return handleDeleteError("Intrusion Service Gateway Policy", id, err)
	}

	return nil
}
