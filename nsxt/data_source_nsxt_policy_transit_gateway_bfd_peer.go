// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtPolicyTransitGatewayBfdPeer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyTransitGatewayBfdPeerRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"parent_path":  getPolicyPathSchema(true, false, "Policy path of the parent Transit Gateway"),
		},
	}
}

func dataSourceNsxtPolicyTransitGatewayBfdPeerRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	parentPath := d.Get("parent_path").(string)
	query := make(map[string]string)
	query["parent_path"] = escapeSpecialCharacters(parentPath) + "*"

	_, err := policyDataSourceResourceReadWithValidation(d, connector, getParentContext(d, m, parentPath), "TransitGatewayBfdPeer", query, false)
	if err != nil {
		return err
	}

	return nil
}
