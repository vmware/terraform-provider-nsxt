// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	infraAPI "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliBareMetalServerTagsClient = infraAPI.NewBareMetalServerTagsClient

func resourceNsxtPolicyBareMetalServerTagsImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	if isSpaceString(importID) {
		return nil, ErrEmptyImportID
	}
	if !isValidBMSExternalID(importID) {
		return nil, fmt.Errorf("invalid import ID %q: expected bare metal server external_id in UUID format", importID)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceNsxtPolicyBareMetalServerTags() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyBareMetalServerTagsCreate,
		Read:   resourceNsxtPolicyBareMetalServerTagsRead,
		Update: resourceNsxtPolicyBareMetalServerTagsUpdate,
		Delete: resourceNsxtPolicyBareMetalServerTagsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyBareMetalServerTagsImporter,
		},
		Description: "This resource provides a way to configure tags on Bare Metal Servers discovered and managed by NSX-T. Tags applied to bare metal servers can be used in dynamic group membership criteria and security policies.",

		Schema: map[string]*schema.Schema{
			"external_id": getBMSExternalIDSchema("External ID of the bare metal server"),
			"tag":         getTagsSchema(),
		},
	}
}

// findBareMetalServerByExternalID and convertSearchResultToBareMetalServerList
// are defined in data_source_nsxt_policy_baremetal_server.go and
// data_source_nsxt_policy_baremetal_servers.go respectively

func updateNsxtPolicyBareMetalServerTags(ctx utl.SessionContext, connector client.Connector, externalID string, tags []model.Tag) error {
	c := cliBareMetalServerTagsClient(ctx, connector)
	if c == nil {
		return policyResourceNotSupportedError()
	}
	payload := model.BareMetalServerTagList{
		BmsExternalId: &externalID,
		Tags:          tags,
	}
	_, err := c.Create(payload)
	return err
}

func resourceNsxtPolicyBareMetalServerTagsCreate(d *schema.ResourceData, m interface{}) error {
	// Enhanced version validation with detailed error message
	if err := validateBMSVersionRequirement(); err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	ctx := utl.SessionContext{ClientType: utl.Local}
	externalID := d.Get("external_id").(string)

	// Verify server exists before applying tags
	if _, err := findBareMetalServerByExternalID(connector, ctx, externalID); err != nil {
		return handleCreateError("BareMetalServerTags", externalID, err)
	}

	// Validate tag combinations before applying
	tags := getPolicyTagsFromSchema(d)
	if tags == nil {
		tags = make([]model.Tag, 0)
	}

	// Basic tag validation - delegate business logic validation to NSX manager

	if err := updateNsxtPolicyBareMetalServerTags(ctx, connector, externalID, tags); err != nil {
		return handleUpdateError("BareMetalServerTags", externalID, err)
	}
	d.SetId(externalID)
	return resourceNsxtPolicyBareMetalServerTagsRead(d, m)
}

func resourceNsxtPolicyBareMetalServerTagsRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	ctx := utl.SessionContext{ClientType: utl.Local}
	externalID := d.Id()
	if externalID == "" {
		return fmt.Errorf("error obtaining bare metal server external_id")
	}

	bms, err := findBareMetalServerByExternalID(connector, ctx, externalID)
	if err != nil {
		return handleReadError(d, "BareMetalServerTags", externalID, err)
	}

	setPolicyTagsInSchema(d, bms.Tags)
	if d.Get("external_id") == "" {
		d.Set("external_id", bms.ExternalId)
	}
	return nil
}

func resourceNsxtPolicyBareMetalServerTagsUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyBareMetalServerTagsCreate(d, m)
}

func resourceNsxtPolicyBareMetalServerTagsDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	ctx := utl.SessionContext{ClientType: utl.Local}
	externalID := d.Get("external_id").(string)

	// Check if server exists before attempting delete
	if _, err := findBareMetalServerByExternalID(connector, ctx, externalID); err != nil {
		return handleDeleteError("BareMetalServerTags", externalID, err)
	}

	if err := updateNsxtPolicyBareMetalServerTags(ctx, connector, externalID, []model.Tag{}); err != nil {
		return handleDeleteError("BareMetalServerTags", externalID, err)
	}
	return nil
}
