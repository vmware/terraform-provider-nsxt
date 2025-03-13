// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
)

func dataSourceNsxtTags() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtTagsRead,
		Schema: map[string]*schema.Schema{
			"scope": getDataSourceStringSchema("The scope of the tags to retrieve. If this is not set it will retrieve all the tags."),
			"items": {
				Type:        schema.TypeList,
				Description: "List of tags based on the scope.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNsxtTagsRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewTagsClient(connector)
	tagsList, err := client.List(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return err
	}

	scopeFilter, isSet := d.GetOk("scope")

	var tags []string
	for _, tag := range tagsList.Results {
		if isSet == false || *tag.Scope == scopeFilter {
			tags = append(tags, *tag.Tag)
		}
	}
	d.Set("items", tags)

	d.SetId(newUUID())

	return nil
}
