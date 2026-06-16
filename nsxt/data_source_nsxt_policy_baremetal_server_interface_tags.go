// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func dataSourceNsxtPolicyBareMetalServerInterfaceTags() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyBareMetalServerInterfaceTagsRead,

		Schema: map[string]*schema.Schema{
			"id":          getDataSourceIDSchema(),
			"external_id": getBMSExternalIDDataSourceSchema("External ID of the bare metal server interface", true),
			"tag":         getTagsSchema(),
		},
	}
}

func dataSourceNsxtPolicyBareMetalServerInterfaceTagsRead(d *schema.ResourceData, m interface{}) error {
	if err := validateBMSVersionRequirement(); err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	ctx := utl.SessionContext{ClientType: utl.Local}

	interfaceId := d.Get("external_id").(string)
	if interfaceId == "" {
		return fmt.Errorf("external_id is required")
	}

	// Find the BMS interface using Policy API search to get its tags
	bmsi, err := findBareMetalServerInterfaceByExternalID(connector, ctx, interfaceId)
	if err != nil {
		return fmt.Errorf("Failed to find Bare Metal Server Interface with external_id %s: %v", interfaceId, err)
	}

	d.SetId(interfaceId)
	d.Set("external_id", bmsi.ExternalId)
	setPolicyTagsInSchema(d, bmsi.Tags)

	return nil
}
