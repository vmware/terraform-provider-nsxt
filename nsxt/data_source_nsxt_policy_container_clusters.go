// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points"
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

	client := enforcement_points.NewClusterControlPlanesClient(connector)

	boolFalse := false
	var cursor *string
	total := 0
	groupsMap := make(map[string]string)

	for {
		results, err := client.List(defaultSite, defaultEnforcementPoint, cursor, nil, nil, nil, &boolFalse, nil)
		if err != nil {
			return err
		}
		if total == 0 && results.ResultCount != nil {
			// first response
			total = int(*results.ResultCount)
		}
		for _, r := range results.Results {
			groupsMap[*r.DisplayName] = *r.Path
		}
		cursor = results.Cursor
		if len(groupsMap) >= total {
			break
		}
	}
	d.Set("items", groupsMap)
	d.SetId(newUUID())
	return nil
}
