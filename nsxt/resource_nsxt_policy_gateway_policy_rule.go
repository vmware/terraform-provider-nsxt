// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gatewaypolicies "github.com/vmware/terraform-provider-nsxt/api/infra/domains/gateway_policies"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
)

var cliGatewayPoliciesRulesClient = gatewaypolicies.NewRulesClient

func resourceNsxtPolicyGatewayRulePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyGatewayPolicyRuleCreate,
		Read:   resourceNsxtPolicyGatewayPolicyRuleRead,
		Update: resourceNsxtPolicyGatewayPolicyRuleUpdate,
		Delete: resourceNsxtPolicyGatewayPolicyRuleDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtGatewayPolicyRuleImporter,
		},

		Schema: getSecurityPolicyAndGatewayRuleSchema(false, false, false, true),
	}
}

func resourceNsxtPolicyGatewayPolicyRuleCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	ruleName := d.Get("display_name").(string)
	policyPath := d.Get("policy_path").(string)
	projectID := getProjectIDFromResourcePath(policyPath)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyGatewayPolicyRuleExistsPartial(d, m, policyPath))
	if err != nil {
		return err
	}

	d.Set("nsx_id", id)
	if err := setSecurityOrGatewayPolicyRuleContext(d, projectID); err != nil {
		return handleCreateError("GatewayPolicyRule", fmt.Sprintf("%s/%s", policyPath, ruleName), err)
	}

	log.Printf("[INFO] Creating Gateway Policy Rule with ID %s under policy %s", id, policyPath)
	client := cliGatewayPoliciesRulesClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	rule := securityAndGatewayPolicyRuleSchemaToModel(d, id)
	err = client.Patch(domain, policyID, id, rule)
	if err != nil {
		return handleCreateError("GatewayPolicyRule", fmt.Sprintf("%s/%s", policyPath, ruleName), err)
	}

	d.SetId(id)

	return resourceNsxtPolicyGatewayPolicyRuleRead(d, m)
}

func resourceNsxtPolicyGatewayPolicyRuleExistsPartial(d *schema.ResourceData, m interface{}, policyPath string) func(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parentContext := getParentContext(d, m, policyPath)
	return func(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
		return resourceNsxtPolicyGatewayPolicyRuleExists(parentContext, id, policyPath, connector)
	}
}

func resourceNsxtPolicyGatewayPolicyRuleExists(sessionContext utl.SessionContext, id string, policyPath string, connector client.Connector) (bool, error) {
	client := cliGatewayPoliciesRulesClient(sessionContext, connector)
	if client == nil {
		return false, policyResourceNotSupportedError()
	}

	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)
	_, err := client.Get(domain, policyID, id)

	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving Gateway Policy Rule", err)
}

func resourceNsxtPolicyGatewayPolicyRuleRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Gateway Policy Rule ID")
	}

	policyPath := d.Get("policy_path").(string)
	projectID := getProjectIDFromResourcePath(policyPath)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	if err := setSecurityOrGatewayPolicyRuleContext(d, projectID); err != nil {
		return handleReadError(d, "GatewayPolicyRule", fmt.Sprintf("%s/%s", policyPath, id), err)
	}

	client := cliGatewayPoliciesRulesClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}
	rule, err := client.Get(domain, policyID, id)
	if err != nil {
		return handleReadError(d, "GatewayPolicyRule", fmt.Sprintf("%s/%s", policyPath, id), err)
	}

	securityAndGatewayPolicyRuleModelToSchema(d, rule)
	return nil
}

func resourceNsxtPolicyGatewayPolicyRuleUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Gateway Policy Rule ID")
	}

	policyPath := d.Get("policy_path").(string)
	log.Printf("[INFO] Updating Gateway Policy Rule with ID %s under policy %s", id, policyPath)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	client := cliGatewayPoliciesRulesClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	rule := securityAndGatewayPolicyRuleSchemaToModel(d, id)
	err := client.Patch(domain, policyID, id, rule)
	if err != nil {
		return handleUpdateError("GatewayPolicyRule", fmt.Sprintf("%s/%s", policyPath, id), err)
	}

	return resourceNsxtPolicyGatewayPolicyRuleRead(d, m)
}

func resourceNsxtPolicyGatewayPolicyRuleDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Get("nsx_id").(string)
	if id == "" {
		return fmt.Errorf("Error obtaining Gateway Policy Rule ID")
	}

	connector := getPolicyConnector(m)

	policyPath := d.Get("policy_path").(string)
	log.Printf("[INFO] Deleting Gateway Policy Rule with ID %s under policy %s", id, policyPath)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	client := cliGatewayPoliciesRulesClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}
	return client.Delete(domain, policyID, id)
}

func nsxtGatewayPolicyRuleImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	rd, err := nsxtPolicyPathResourceImporterHelper(d, m)
	if err != nil {
		return rd, err
	}
	ruleIdx := strings.Index(importID, "rule")
	if ruleIdx <= 0 {
		return nil, fmt.Errorf("invalid path of Security Policy Rule to import")
	}
	d.Set("policy_path", importID[:ruleIdx-1])
	return []*schema.ResourceData{d}, nil
}
