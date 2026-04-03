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

	intrusionservicegatewaypolicies "github.com/vmware/terraform-provider-nsxt/api/infra/domains/intrusion_service_gateway_policies"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliIntrusionServiceGatewayPolicyRulesClient = intrusionservicegatewaypolicies.NewRulesClient

func resourceNsxtPolicyIntrusionServiceGatewayPolicyRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleCreate,
		Read:   resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead,
		Update: resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleUpdate,
		Delete: resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtIntrusionServiceGatewayPolicyRuleImporter,
		},
		Schema: getIntrusionServiceGatewayPolicyRuleSchema(false),
	}
}

func intrusionServiceGatewayPolicyRuleSchemaToModel(d *schema.ResourceData, id string) model.IdsRule {
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

func intrusionServiceGatewayPolicyRuleModelToSchema(d *schema.ResourceData, rule model.IdsRule) {
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

func resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	policyPath := d.Get("policy_path").(string)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleExistsPartial(d, m, policyPath))
	if err != nil {
		return err
	}

	d.Set("nsx_id", id)

	log.Printf("[INFO] Creating Intrusion Service Gateway Policy Rule with ID %s under policy %s", id, policyPath)
	client := cliIntrusionServiceGatewayPolicyRulesClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	rule := intrusionServiceGatewayPolicyRuleSchemaToModel(d, id)
	err = client.Patch(domain, policyID, id, rule)
	if err != nil {
		return handleCreateError("Intrusion Service Gateway Policy Rule", fmt.Sprintf("%s/%s", policyPath, id), err)
	}

	d.SetId(id)

	return resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, m)
}

func resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleExistsPartial(d *schema.ResourceData, m interface{}, policyPath string) func(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	return func(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
		return resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleExists(sessionContext, id, policyPath, connector)
	}
}

func resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleExists(sessionContext utl.SessionContext, id string, policyPath string, connector client.Connector) (bool, error) {
	client := cliIntrusionServiceGatewayPolicyRulesClient(sessionContext, connector)
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

	return false, logAPIError("Error retrieving Intrusion Service Gateway Policy Rule", err)
}

func resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Intrusion Service Gateway Policy Rule ID")
	}

	policyPath := d.Get("policy_path").(string)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	client := cliIntrusionServiceGatewayPolicyRulesClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}
	rule, err := client.Get(domain, policyID, id)
	if err != nil {
		return handleReadError(d, "Intrusion Service Gateway Policy Rule", fmt.Sprintf("%s/%s", policyPath, id), err)
	}

	intrusionServiceGatewayPolicyRuleModelToSchema(d, rule)
	return nil
}

func resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Intrusion Service Gateway Policy Rule ID")
	}

	policyPath := d.Get("policy_path").(string)
	log.Printf("[INFO] Updating Intrusion Service Gateway Policy Rule with ID %s under policy %s", id, policyPath)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	client := cliIntrusionServiceGatewayPolicyRulesClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	rule := intrusionServiceGatewayPolicyRuleSchemaToModel(d, id)
	err := client.Patch(domain, policyID, id, rule)
	if err != nil {
		return handleUpdateError("Intrusion Service Gateway Policy Rule", fmt.Sprintf("%s/%s", policyPath, id), err)
	}

	return resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, m)
}

func resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Intrusion Service Gateway Policy Rule ID")
	}

	connector := getPolicyConnector(m)

	policyPath := d.Get("policy_path").(string)
	log.Printf("[INFO] Deleting Intrusion Service Gateway Policy Rule with ID %s under policy %s", id, policyPath)
	domain := getDomainFromResourcePath(policyPath)
	policyID := getPolicyIDFromPath(policyPath)

	client := cliIntrusionServiceGatewayPolicyRulesClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}
	err := client.Delete(domain, policyID, id)

	if err != nil {
		return handleDeleteError("Intrusion Service Gateway Policy Rule", id, err)
	}
	return nil
}

func nsxtIntrusionServiceGatewayPolicyRuleImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	rd, err := nsxtPolicyPathResourceImporterHelper(d, m)
	if err != nil {
		return rd, err
	}
	ruleIdx := strings.Index(importID, "rule")
	if ruleIdx <= 0 {
		return nil, fmt.Errorf("invalid path of Intrusion Service Gateway Policy Rule to import")
	}
	d.Set("policy_path", importID[:ruleIdx-1])
	return []*schema.ResourceData{d}, nil
}
