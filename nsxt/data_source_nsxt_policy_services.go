/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
)

func dataSourceNsxtPolicyServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyServicesRead,

		Schema: map[string]*schema.Schema{
			"context": getContextSchema(false, false, false),
			"items": {
				Type:        schema.TypeList,
				Description: "List of services",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNsxtPolicyServicesRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewServicesClient(getSessionContext(d, m), connector)

	var servicesList []interface{}
	results, err := client.List(nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return err
	}
	for _, r := range results.Results {
		serviceMap := make(map[string]interface{})
		serviceMap["id"] = r.Id
		serviceMap["display_name"] = r.DisplayName
		serviceMap["path"] = r.Path
		servicesList = append(servicesList, serviceMap)
	}

	d.SetId(newUUID())
	d.Set("items", servicesList)
	return nil
}
