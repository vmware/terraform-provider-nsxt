// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtPolicyTransitGatewayBgpNeighbor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyTransitGatewayBgpNeighborRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"parent_path":  getPolicyPathSchema(true, false, "Policy path of the parent Transit Gateway"),
		},
	}
}

func dataSourceNsxtPolicyTransitGatewayBgpNeighborRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// BgpNeighborConfig is shared across Tier-0/Tier-1 and Transit Gateway BGP
	// neighbors, so parent_path is required to scope the search to this
	// Transit Gateway.
	parentPath := d.Get("parent_path").(string)
	query := make(map[string]string)
	query["parent_path"] = escapeSpecialCharacters(parentPath) + "*"

	_, err := policyDataSourceResourceReadWithValidation(d, connector, getParentContext(d, m, parentPath), "BgpNeighborConfig", query, false)
	if err != nil {
		return err
	}

	return nil
}
