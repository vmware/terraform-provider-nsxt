// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtTags() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtTagsRead,
		Schema: map[string]*schema.Schema{
			"scope": getRequiredStringSchema("The scope of the tags to retrieve."),
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
	scope := d.Get("scope").(string)
	tagsList, err := listTags(connector, scope)
	if err != nil {
		return err
	}

	var tags []string
	for _, tag := range tagsList {
		tags = append(tags, *tag.Tag)
	}
	d.Set("items", tags)

	d.SetId(newUUID())

	return nil
}

func listTags(connector client.Connector, scope string) ([]model.TagInfo, error) {
	client := infra.NewTagsClient(connector)
	if client == nil {
		return nil, policyResourceNotSupportedError()
	}
	var cursor *string
	var tagsListResults []model.TagInfo
	total := 0
	for {
		tagsList, err := client.List(cursor, nil, nil, nil, nil, &scope, nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		tagsListResults = append(tagsListResults, tagsList.Results...)
		if total == 0 && tagsList.ResultCount != nil {
			total = int(*tagsList.ResultCount)
		}

		cursor = tagsList.Cursor
		if len(tagsListResults) >= total {
			return tagsListResults, nil
		}
	}

}
