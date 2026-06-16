// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func dataSourceNsxtPolicyBareMetalServerTags() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyBareMetalServerTagsRead,

		Schema: map[string]*schema.Schema{
			"id":          getDataSourceIDSchema(),
			"external_id": getBMSExternalIDDataSourceSchema("External ID of the bare metal server", true),
			"tag":         getTagsSchema(),
		},
	}
}

func dataSourceNsxtPolicyBareMetalServerTagsRead(d *schema.ResourceData, m interface{}) error {
	if err := validateBMSVersionRequirement(); err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	ctx := utl.SessionContext{ClientType: utl.Local}

	serverId := d.Get("external_id").(string)
	if serverId == "" {
		return fmt.Errorf("external_id is required")
	}

	// Find the BMS server using Policy API search to get its tags
	bms, err := findBareMetalServerByExternalID(connector, ctx, serverId)
	if err != nil {
		return fmt.Errorf("Failed to find Bare Metal Server with external_id %s: %v", serverId, err)
	}

	d.SetId(serverId)
	d.Set("external_id", bms.ExternalId)
	setPolicyTagsInSchema(d, bms.Tags)

	return nil
}
