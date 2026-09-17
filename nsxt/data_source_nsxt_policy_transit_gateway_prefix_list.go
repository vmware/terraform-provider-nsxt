// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtPolicyTransitGatewayPrefixList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyTransitGatewayPrefixListRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"parent_path":  getPolicyPathSchema(true, false, "Policy path of the parent Transit Gateway"),
		},
	}
}

func dataSourceNsxtPolicyTransitGatewayPrefixListRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// PrefixList is shared across Tier-0/Tier-1 and Transit Gateway prefix
	// lists, so parent_path is required to scope the search to this
	// Transit Gateway.
	parentPath := d.Get("parent_path").(string)
	query := make(map[string]string)
	query["parent_path"] = escapeSpecialCharacters(parentPath) + "*"

	_, err := policyDataSourceResourceReadWithValidation(d, connector, getParentContext(d, m, parentPath), "PrefixList", query, false)
	if err != nil {
		return err
	}

	return nil
}
