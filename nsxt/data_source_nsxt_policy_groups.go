/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"

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
				Type:        schema.TypeMap,
				Description: "Mapping of service UUID by display name",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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

	groupsMap := make(map[string]interface{})
	results, err := client.List(domainName, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return err
	}
	for _, r := range results.Results {
		if _, ok := groupsMap[*r.DisplayName]; ok {
			return fmt.Errorf("found duplicate policy group %s", *r.DisplayName)
		}
		groupMap := make(map[string]interface{})
		groupMap["id"] = r.Id
		groupMap["path"] = r.Path
		groupsMap[*r.DisplayName] = groupMap
	}

	d.SetId(newUUID())
	d.Set("items", groupsMap)
	return nil
}
