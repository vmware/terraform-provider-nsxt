// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

// list the tier0 gateways in map ID:displayname
func dataSourceNsxtPolicyTier1Gateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyTier1GatewaysRead,

		Schema: map[string]*schema.Schema{
			"items": {
				Type:        schema.TypeMap,
				Description: "Mapping of Tier1 instance ID by display name",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"context": getContextSchemaWithSpec(utl.SessionContextSpec{IsRequired: false, IsComputed: false, IsVpc: false, AllowDefaultProject: false, FromGlobal: true}),
		},
	}
}

func dataSourceNsxtPolicyTier1GatewaysRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	resultMap := make(map[string]string)
	err := policyDataSourceCreateMap(connector, getSessionContext(d, m), "Tier1", resultMap, nil)
	if err != nil {
		return fmt.Errorf("error in listing the Tier1 gateways items : %v", err)
	}
	d.SetId(newUUID())
	d.Set("items", resultMap)
	return nil
}
