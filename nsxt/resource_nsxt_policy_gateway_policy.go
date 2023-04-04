/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gm_domains "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/domains"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func resourceNsxtPolicyGatewayPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyGatewayPolicyCreate,
		Read:   resourceNsxtPolicyGatewayPolicyRead,
		Update: resourceNsxtPolicyGatewayPolicyUpdate,
		Delete: resourceNsxtPolicyGatewayPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtDomainResourceImporter,
		},

		Schema: getPolicyGatewayPolicySchema(),
	}
}

func getGatewayPolicyInDomain(id string, domainName string, connector client.Connector, isGlobalManager bool) (model.GatewayPolicy, error) {
	if isGlobalManager {
		client := gm_domains.NewGatewayPoliciesClient(connector)
		gmObj, err := client.Get(domainName, id)
		if err != nil {
			return model.GatewayPolicy{}, err
		}
		rawObj, convErr := convertModelBindingType(gmObj, gm_model.GatewayPolicyBindingType(), model.GatewayPolicyBindingType())
		if convErr != nil {
			return model.GatewayPolicy{}, convErr
		}
		return rawObj.(model.GatewayPolicy), nil
	}
	client := domains.NewGatewayPoliciesClient(connector)
	return client.Get(domainName, id)

}

func resourceNsxtPolicyGatewayPolicyExistsInDomain(id string, domainName string, connector client.Connector, isGlobalManager bool) (bool, error) {
	_, err := getGatewayPolicyInDomain(id, domainName, connector, isGlobalManager)

	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving Gateway Policy", err)
}

func resourceNsxtPolicyGatewayPolicyExistsPartial(domainName string) func(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
	return func(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
		return resourceNsxtPolicyGatewayPolicyExistsInDomain(id, domainName, connector, isGlobalManager)
	}
}

func getUpdatedRuleChildren(d *schema.ResourceData) ([]*data.StructValue, error) {
	var policyChildren []*data.StructValue

	if !d.HasChange("rule") {
		return nil, nil
	}

	oldRules, newRules := d.GetChange("rule")
	rules := getPolicyRulesFromSchema(d)
	newRulesCount := len(newRules.([]interface{}))
	oldRulesCount := len(oldRules.([]interface{}))
	for ruleNo := 0; ruleNo < newRulesCount; ruleNo++ {
		ruleIndicator := fmt.Sprintf("rule.%d", ruleNo)
		if d.HasChange(ruleIndicator) {
			rule := rules[ruleNo]
			// New or updated rule
			ruleID := newUUID()
			if rule.Id != nil {
				ruleID = *rule.Id
				log.Printf("[DEBUG]: Updating child rule with id %s", *rule.Id)
			} else {
				log.Printf("[DEBUG]: Adding child rule with id %s", ruleID)
				rule.Id = &ruleID
			}

			childRule, err := createPolicyChildRule(ruleID, rule, false)
			if err != nil {
				return policyChildren, err
			}
			policyChildren = append(policyChildren, childRule)
		}
	}

	resourceType := "Rule"
	for ruleNo := newRulesCount; ruleNo < oldRulesCount; ruleNo++ {
		// delete
		ruleIndicator := fmt.Sprintf("rule.%d", ruleNo)
		oldRule, _ := d.GetChange(ruleIndicator)
		oldRuleMap := oldRule.(map[string]interface{})
		oldRuleID := oldRuleMap["nsx_id"].(string)
		rule := model.Rule{
			Id:           &oldRuleID,
			ResourceType: &resourceType,
		}

		childRule, err := createPolicyChildRule(oldRuleID, rule, true)
		if err != nil {
			return policyChildren, err
		}
		log.Printf("[DEBUG]: Deleting child rule with id %s", oldRuleID)
		policyChildren = append(policyChildren, childRule)

	}

	return policyChildren, nil

}

func policyGatewayPolicyBuildAndPatch(d *schema.ResourceData, m interface{}, connector client.Connector, isGlobalManager bool, id string) error {

	domain := d.Get("domain").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	category := d.Get("category").(string)
	comments := d.Get("comments").(string)
	locked := d.Get("locked").(bool)
	sequenceNumber := int64(d.Get("sequence_number").(int))
	stateful := d.Get("stateful").(bool)
	revision := int64(d.Get("revision").(int))
	objType := "GatewayPolicy"

	obj := model.GatewayPolicy{
		DisplayName:    &displayName,
		Description:    &description,
		Tags:           tags,
		Category:       &category,
		Comments:       &comments,
		Locked:         &locked,
		SequenceNumber: &sequenceNumber,
		Stateful:       &stateful,
		ResourceType:   &objType,
		Id:             &id,
	}
	_, isSet := d.GetOkExists("tcp_strict")
	if isSet {
		tcpStrict := d.Get("tcp_strict").(bool)
		obj.TcpStrict = &tcpStrict
	}

	if len(d.Id()) > 0 {
		// This is update flow
		obj.Revision = &revision
	}

	policyChildren, err := getUpdatedRuleChildren(d)
	if err != nil {
		return err
	}
	if len(policyChildren) > 0 {
		obj.Children = policyChildren
	}

	return gatewayPolicyInfraPatch(obj, domain, m)
}

func resourceNsxtPolicyGatewayPolicyCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicyGatewayPolicyExistsPartial(d.Get("domain").(string)))
	if err != nil {
		return err
	}

	err = policyGatewayPolicyBuildAndPatch(d, m, connector, isPolicyGlobalManager(m), id)
	if err != nil {
		return handleCreateError("Gateway Policy", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyGatewayPolicyRead(d, m)
}

func resourceNsxtPolicyGatewayPolicyRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Gateway Policy ID")
	}

	obj, err := getGatewayPolicyInDomain(id, d.Get("domain").(string), connector, isPolicyGlobalManager(m))
	if err != nil {
		return handleReadError(d, "Gateway Policy", id, err)
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
		// tcp_strict is dependant on stateful and maybe nil
		d.Set("tcp_strict", *obj.TcpStrict)
	}
	d.Set("revision", obj.Revision)
	return setPolicyRulesInSchema(d, obj.Rules)
}

func resourceNsxtPolicyGatewayPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Gateway Policy ID")
	}

	err := policyGatewayPolicyBuildAndPatch(d, m, connector, isPolicyGlobalManager(m), id)
	if err != nil {
		return handleUpdateError("Gateway Policy", id, err)
	}

	return resourceNsxtPolicyGatewayPolicyRead(d, m)
}

func resourceNsxtPolicyGatewayPolicyDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Gateway Policy ID")
	}

	connector := getPolicyConnector(m)
	var err error
	if isPolicyGlobalManager(m) {
		client := gm_domains.NewGatewayPoliciesClient(connector)
		err = client.Delete(d.Get("domain").(string), id)
	} else {
		client := domains.NewGatewayPoliciesClient(connector)
		err = client.Delete(d.Get("domain").(string), id)
	}
	if err != nil {
		return handleDeleteError("Gateway Policy", id, err)
	}

	return nil
}
