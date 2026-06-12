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

var cliBareMetalServerInterfaceTagsClient = infraAPI.NewBareMetalServerInterfaceTagsClient

func resourceNsxtPolicyBareMetalServerInterfaceTagsImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	if isSpaceString(importID) {
		return nil, ErrEmptyImportID
	}
	if !isValidResourceID(importID) {
		return nil, fmt.Errorf("invalid import ID %q: expected bare metal server interface external_id", importID)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceNsxtPolicyBareMetalServerInterfaceTags() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyBareMetalServerInterfaceTagsCreate,
		Read:   resourceNsxtPolicyBareMetalServerInterfaceTagsRead,
		Update: resourceNsxtPolicyBareMetalServerInterfaceTagsUpdate,
		Delete: resourceNsxtPolicyBareMetalServerInterfaceTagsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyBareMetalServerInterfaceTagsImporter,
		},
		Description: "This resource provides a way to configure tags on Bare Metal Server Interfaces discovered and managed by NSX-T. Tags applied to bare metal server interfaces can be used in dynamic group membership criteria and security policies.",

		Schema: map[string]*schema.Schema{
			"external_id": getBMSExternalIDSchema("External ID of the bare metal server interface"),
			"tag":         getTagsSchema(),
		},
	}
}

// findBareMetalServerInterfaceByExternalID and convertSearchResultToBareMetalServerInterfaceList
// are defined in data_source_nsxt_policy_baremetal_server_interface.go and
// data_source_nsxt_policy_baremetal_server_interfaces.go respectively

func updateNsxtPolicyBareMetalServerInterfaceTags(ctx utl.SessionContext, connector client.Connector, externalID string, tags []model.Tag) error {
	c := cliBareMetalServerInterfaceTagsClient(ctx, connector)
	if c == nil {
		return policyResourceNotSupportedError()
	}
	payload := model.BareMetalServerInterfaceTagList{
		BmsInterfaceExternalId: &externalID,
		Tags:                   tags,
	}
	_, err := c.Create(payload)
	return err
}

func resourceNsxtPolicyBareMetalServerInterfaceTagsCreate(d *schema.ResourceData, m interface{}) error {
	// Enhanced version validation with detailed error message
	if err := validateBMSVersionRequirement(); err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	ctx := utl.SessionContext{ClientType: utl.Local}
	externalID := d.Get("external_id").(string)

	if _, err := findBareMetalServerInterfaceByExternalID(connector, ctx, externalID); err != nil {
		return handleCreateError("BareMetalServerInterfaceTags", externalID, err)
	}

	tags := getPolicyTagsFromSchema(d)
	if tags == nil {
		tags = make([]model.Tag, 0)
	}
	if err := updateNsxtPolicyBareMetalServerInterfaceTags(ctx, connector, externalID, tags); err != nil {
		return handleUpdateError("BareMetalServerInterfaceTags", externalID, err)
	}
	d.SetId(externalID)
	return resourceNsxtPolicyBareMetalServerInterfaceTagsRead(d, m)
}

func resourceNsxtPolicyBareMetalServerInterfaceTagsRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	ctx := utl.SessionContext{ClientType: utl.Local}
	externalID := d.Id()
	if externalID == "" {
		return fmt.Errorf("error obtaining bare metal server interface external_id")
	}

	bmsi, err := findBareMetalServerInterfaceByExternalID(connector, ctx, externalID)
	if err != nil {
		return handleReadError(d, "BareMetalServerInterfaceTags", externalID, err)
	}

	setPolicyTagsInSchema(d, bmsi.Tags)
	if d.Get("external_id") == "" {
		d.Set("external_id", bmsi.ExternalId)
	}
	return nil
}

func resourceNsxtPolicyBareMetalServerInterfaceTagsUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceNsxtPolicyBareMetalServerInterfaceTagsCreate(d, m)
}

func resourceNsxtPolicyBareMetalServerInterfaceTagsDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	ctx := utl.SessionContext{ClientType: utl.Local}
	externalID := d.Get("external_id").(string)

	// Check if interface exists before attempting delete
	if _, err := findBareMetalServerInterfaceByExternalID(connector, ctx, externalID); err != nil {
		return handleDeleteError("BareMetalServerInterfaceTags", externalID, err)
	}

	if err := updateNsxtPolicyBareMetalServerInterfaceTags(ctx, connector, externalID, []model.Tag{}); err != nil {
		return handleDeleteError("BareMetalServerInterfaceTags", externalID, err)
	}
	return nil
}
