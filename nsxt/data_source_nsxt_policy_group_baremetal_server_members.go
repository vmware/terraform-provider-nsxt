// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	groupsAPI "github.com/vmware/terraform-provider-nsxt/api/infra/domains/groups"
)

var cliGroupBareMetalServerMembersClient = groupsAPI.NewBmsMembersClient

func dataSourceNsxtPolicyGroupBareMetalServerMembers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyGroupBareMetalServerMembersRead,

		Schema: map[string]*schema.Schema{
			"id":     getDataSourceIDSchema(),
			"domain": getDomainNameSchema(),
			"group_id": {
				Type:        schema.TypeString,
				Description: "ID of the group to read effective bare metal server members",
				Required:    true,
			},
			"enforcement_point_path": {
				Type:        schema.TypeString,
				Description: "Enforcement point path",
				Optional:    true,
			},
			"items": {
				Type:        schema.TypeList,
				Description: "Bare metal server members of the group",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_id":    {Type: schema.TypeString, Computed: true},
						"display_name":   {Type: schema.TypeString, Computed: true},
						"source_id":      {Type: schema.TypeString, Computed: true},
						"cpu_cores":      {Type: schema.TypeInt, Computed: true},
						"os_name":        {Type: schema.TypeString, Computed: true},
						"os_version":     {Type: schema.TypeString, Computed: true},
						"last_sync_time": {Type: schema.TypeInt, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceNsxtPolicyGroupBareMetalServerMembersRead(d *schema.ResourceData, m interface{}) error {
	if err := validateBMSVersionRequirement(); err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	client := cliGroupBareMetalServerMembersClient(getSessionContext(d, m), connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	domainID := d.Get("domain").(string)
	if domainID == "" {
		domainID = "default"
	}
	groupID := d.Get("group_id").(string)

	enforcementPointPath := d.Get("enforcement_point_path").(string)
	var enforcementPointPathPtr *string
	if enforcementPointPath != "" {
		enforcementPointPathPtr = &enforcementPointPath
	}

	obj, err := client.List(domainID, groupID, nil, enforcementPointPathPtr, nil, nil, nil, nil, nil)
	if err != nil {
		return handleReadError(d, "Group Bare Metal Server Members", groupID, err)
	}

	d.SetId(groupID)

	var items []map[string]interface{}
	for _, server := range obj.Results {
		elem := make(map[string]interface{})
		elem["external_id"] = server.ExternalId
		elem["display_name"] = server.DisplayName
		elem["source_id"] = server.SourceId
		elem["last_sync_time"] = server.LastSyncTime

		if server.CpuCores != nil {
			elem["cpu_cores"] = *server.CpuCores
		}

		if server.OsInfo != nil {
			if server.OsInfo.OsName != nil {
				elem["os_name"] = *server.OsInfo.OsName
			}
			if server.OsInfo.OsVersion != nil {
				elem["os_version"] = *server.OsInfo.OsVersion
			}
		}

		items = append(items, elem)
	}
	d.Set("items", items)

	return nil
}
