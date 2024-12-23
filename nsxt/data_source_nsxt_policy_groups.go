/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/terraform-provider-nsxt/api/infra/domains"
)

func dataSourceNsxtPolicyGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyGroupsRead,

		Schema: map[string]*schema.Schema{
			"context": getContextSchema(false, false, false),
			"domain":  getDomainNameSchema(),
			"items": {
				Type:        schema.TypeList,
				Description: "Mapping of service UUID by display name",
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

func dataSourceNsxtPolicyGroupsRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	domainName := d.Get("domain").(string)

	client := domains.NewGroupsClient(getSessionContext(d, m), connector)

	var groupsList []interface{}
	results, err := client.List(domainName, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return err
	}
	for _, r := range results.Results {
		groupMap := make(map[string]interface{})
		groupMap["id"] = r.Id
		groupMap["display_name"] = r.DisplayName
		groupMap["path"] = r.Path
		groupsList = append(groupsList, groupMap)
	}

	d.SetId(newUUID())
	d.Set("items", groupsList)
	return nil
}
