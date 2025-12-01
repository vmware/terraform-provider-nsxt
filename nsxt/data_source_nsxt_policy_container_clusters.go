// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtPolicyContainerClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyContainerClustersRead,

		Schema: map[string]*schema.Schema{
			"items": {
				Type:        schema.TypeMap,
				Description: "Mapping of service policy path by display name",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNsxtPolicyContainerClustersRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	groupsMap := make(map[string]string)
	results, err := listClusterControlPlanes(defaultSite, defaultEnforcementPoint, connector)
	if err != nil {
		return err
	}
	for _, r := range results {
		groupsMap[*r.DisplayName] = *r.Path
	}

	d.Set("items", groupsMap)
	d.SetId(newUUID())
	return nil
}
