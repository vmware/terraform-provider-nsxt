// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
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
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	client := infra.NewTagsClient(sessionContext, connector)
	if client == nil {
		return nil, policyResourceNotSupportedError()
	}
	var cursor *string
	var tagsListResults []model.TagInfo
	for {
		tagsList, err := client.List(cursor, nil, nil, nil, nil, &scope, nil, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		if tagsList.Results != nil {
			tagsListResults = append(tagsListResults, tagsList.Results...)
		}

		cursor = tagsList.Cursor
		if cursor == nil {
			return tagsListResults, nil
		}
	}
}
