// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	intrusionservicepolicies "github.com/vmware/terraform-provider-nsxt/api/infra/domains/intrusion_service_policies"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliIntrusionServicePolicyRulesClient = intrusionservicepolicies.NewRulesClient

func resourceNsxtPolicyIntrusionServicePolicyRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIntrusionServicePolicyRuleCreate,
		Read:   resourceNsxtPolicyIntrusionServicePolicyRuleRead,
		Update: resourceNsxtPolicyIntrusionServicePolicyRuleUpdate,
		Delete: resourceNsxtPolicyIntrusionServicePolicyRuleDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtIntrusionServicePolicyRuleImporter,
		},
		Schema: getSecurityPolicyAndGatewayRuleSchema(false, true, false, true),
	}
}

func intrusionServicePolicyRuleSchemaToModel(d *schema.ResourceData, id string) model.IdsRule {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	action := d.Get("action").(string)
	logged := d.Get("logged").(bool)
	tag := d.Get("log_label").(string)
	disabled := d.Get("disabled").(bool)
	sourcesExcluded := d.Get("sources_excluded").(bool)
	destinationsExcluded := d.Get("destinations_excluded").(bool)
	ipProtocol := d.Get("ip_version").(string)
	direction := d.Get("direction").(string)
	notes := d.Get("notes").(string)
	seq := d.Get("sequence_number").(int)
	sequenceNumber := int64(seq)
	tagStructs := getPolicyTagsFromSet(d.Get("tag").(*schema.Set))

	resourceType := "IdsRule"
	rule := model.IdsRule{
		ResourceType:         &resourceType,
		Id:                   &id,
		DisplayName:          &displayName,
		Notes:                &notes,
		Description:          &description,
		Action:               &action,
		Logged:               &logged,
		Tag:                  &tag,
		Tags:                 tagStructs,
		Disabled:             &disabled,
		SourcesExcluded:      &sourcesExcluded,
		DestinationsExcluded: &destinationsExcluded,
		IpProtocol:           &ipProtocol,
		Direction:            &direction,
		SourceGroups:         getPathListFromSchema(d, "source_groups"),
		DestinationGroups:    getPathListFromSchema(d, "destination_groups"),
		Services:             getPathListFromSchema(d, "services"),
		Scope:                getPathListFromSchema(d, "scope"),
		SequenceNumber:       &sequenceNumber,
		IdsProfiles:          getPathListFromSchema(d, "ids_profiles"),
	}

	return rule
}

func intrusionServicePolicyRuleModelToSchema(d *schema.ResourceData, rule model.IdsRule) {
	d.Set("display_name", rule.DisplayName)
	d.Set("description", rule.Description)
	d.Set("path", rule.Path)
	d.Set("notes", rule.Notes)
	d.Set("logged", rule.Logged)
	d.Set("log_label", rule.Tag)
	d.Set("action", rule.Action)
	d.Set("destinations_excluded", rule.DestinationsExcluded)
	d.Set("sources_excluded", rule.SourcesExcluded)
	d.Set("ip_version", rule.IpProtocol)
	d.Set("direction", rule.Direction)
	d.Set("disabled", rule.Disabled)
	d.Set("revision", rule.Revision)
	setPathListInSchema(d, "source_groups", rule.SourceGroups)
	setPathListInSchema(d, "destination_groups", rule.DestinationGroups)
	setPathListInSchema(d, "services", rule.Services)
	setPathListInSchema(d, "scope", rule.Scope)
	setPathListInSchema(d, "ids_profiles", rule.IdsProfiles)
	d.Set("sequence_number", rule.SequenceNumber)
	d.Set("nsx_id", rule.Id)

	setPolicyTagsInSchema(d, rule.Tags)
}

func resourceNsxtPolicyIntrusionServicePolicyRuleCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	policyPath := d.Get("policy_path").(string)
	projectID := getProjectIDFromResourcePath(policyPath)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyIntrusionServicePolicyRuleExistsPartial(d, m, policyPath))
	if err != nil {
		return err
	}

	if err := setSecurityOrGatewayPolicyRuleContext(d, projectID); err != nil {
		return handleCreateError("IntrusionServicePolicyRule", fmt.Sprintf("%s/%s", policyPath, id), err)
	}

	log.Printf("[INFO] Creating Intrusion Service Policy Rule with ID %s under policy %s", id, policyPath)
	client := cliIntrusionServicePolicyRulesClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	rule := intrusionServicePolicyRuleSchemaToModel(d, id)
	err = client.Patch(domain, policyID, id, rule)
	if err != nil {
		return handleCreateError("Intrusion Service Policy Rule", fmt.Sprintf("%s/%s", policyPath, id), err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyIntrusionServicePolicyRuleRead(d, m)
}

func resourceNsxtPolicyIntrusionServicePolicyRuleExistsPartial(d *schema.ResourceData, m interface{}, policyPath string) func(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	parentContext := getParentContext(d, m, policyPath)
	return func(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
		return resourceNsxtPolicyIntrusionServicePolicyRuleExists(parentContext, id, policyPath, connector)
	}
}

func resourceNsxtPolicyIntrusionServicePolicyRuleExists(sessionContext utl.SessionContext, id string, policyPath string, connector client.Connector) (bool, error) {
	client := cliIntrusionServicePolicyRulesClient(sessionContext, connector)
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

	return false, logAPIError("Error retrieving Intrusion Service Policy Rule", err)
}

func resourceNsxtPolicyIntrusionServicePolicyRuleRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Intrusion Service Policy Rule ID")
	}

	policyPath := d.Get("policy_path").(string)
	projectID := getProjectIDFromResourcePath(policyPath)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	if err := setSecurityOrGatewayPolicyRuleContext(d, projectID); err != nil {
		return handleReadError(d, "IntrusionServicePolicyRule", fmt.Sprintf("%s/%s", policyPath, id), err)
	}

	client := cliIntrusionServicePolicyRulesClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}
	rule, err := client.Get(domain, policyID, id)
	if err != nil {
		return handleReadError(d, "Intrusion Service Policy Rule", fmt.Sprintf("%s/%s", policyPath, id), err)
	}

	intrusionServicePolicyRuleModelToSchema(d, rule)
	return nil
}

func resourceNsxtPolicyIntrusionServicePolicyRuleUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Intrusion Service Policy Rule ID")
	}

	policyPath := d.Get("policy_path").(string)
	log.Printf("[INFO] Updating Intrusion Service Policy Rule with ID %s under policy %s", id, policyPath)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	client := cliIntrusionServicePolicyRulesClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	rule := intrusionServicePolicyRuleSchemaToModel(d, id)
	err := client.Patch(domain, policyID, id, rule)
	if err != nil {
		return handleUpdateError("Intrusion Service Policy Rule", fmt.Sprintf("%s/%s", policyPath, id), err)
	}

	return resourceNsxtPolicyIntrusionServicePolicyRuleRead(d, m)
}

func resourceNsxtPolicyIntrusionServicePolicyRuleDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Intrusion Service Policy Rule ID")
	}

	connector := getPolicyConnector(m)

	policyPath := d.Get("policy_path").(string)
	log.Printf("[INFO] Deleting Intrusion Service Policy Rule with ID %s under policy %s", id, policyPath)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	client := cliIntrusionServicePolicyRulesClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}
	err := client.Delete(domain, policyID, id)
	if err != nil {
		return handleDeleteError("Intrusion Service Policy Rule", id, err)
	}
	return nil
}

func nsxtIntrusionServicePolicyRuleImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	rd, err := nsxtPolicyPathResourceImporterHelper(d, m)
	if err != nil {
		return rd, err
	}
	ruleIdx := strings.Index(importID, "rule")
	if ruleIdx <= 0 {
		return nil, fmt.Errorf("invalid path of Intrusion Service Policy Rule to import")
	}
	d.Set("policy_path", importID[:ruleIdx-1])
	return []*schema.ResourceData{d}, nil
}
