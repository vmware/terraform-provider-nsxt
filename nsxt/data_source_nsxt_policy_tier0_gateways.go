// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//	data "nsxt_policy_tier0_gateways" "all" {
//		context {
//			project_id = ""
//		}
//	}
//
// list the tier0 gateways in map ID:displayname
func dataSourceNsxtPolicyTier0Gateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyTier0GatewaysRead,

		Schema: map[string]*schema.Schema{
			"items": {
				Type:        schema.TypeMap,
				Description: "Mapping of Tier0 instance ID by display name",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNsxtPolicyTier0GatewaysRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	resultMap := make(map[string]string)
	err := policyDataSourceCreateMap(connector, getSessionContext(d, m), "Tier0", resultMap, nil)
	if err != nil {
		return fmt.Errorf("error in listing the Tier0 gateways items : %v", err)
	}
	d.SetId(newUUID())
	d.Set("items", resultMap)
	return nil
}
