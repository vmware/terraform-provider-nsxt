// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	groupsAPI "github.com/vmware/terraform-provider-nsxt/api/infra/domains/groups"
)

var cliGroupBareMetalServerInterfaceMembersClient = groupsAPI.NewBmsiMembersClient

func dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembersRead,

		Schema: map[string]*schema.Schema{
			"id":     getDataSourceIDSchema(),
			"domain": getDomainNameSchema(),
			"group_id": {
				Type:        schema.TypeString,
				Description: "ID of the group to read effective bare metal server interface members",
				Required:    true,
			},
			"enforcement_point_path": {
				Type:        schema.TypeString,
				Description: "Enforcement point path",
				Optional:    true,
			},
			"items": {
				Type:        schema.TypeList,
				Description: "Bare metal server interface members of the group",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_id":       {Type: schema.TypeString, Computed: true},
						"display_name":      {Type: schema.TypeString, Computed: true},
						"bms_external_id":   {Type: schema.TypeString, Computed: true},
						"source_id":         {Type: schema.TypeString, Computed: true},
						"mac_address":       {Type: schema.TypeString, Computed: true},
						"is_mgmt_interface": {Type: schema.TypeBool, Computed: true},
						"state":             {Type: schema.TypeString, Computed: true},
						"ip_addresses":      {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"last_sync_time":    {Type: schema.TypeInt, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembersRead(d *schema.ResourceData, m interface{}) error {
	if err := validateBMSVersionRequirement(); err != nil {
		return err
	}

	connector := getPolicyConnector(m)
	client := cliGroupBareMetalServerInterfaceMembersClient(getSessionContext(d, m), connector)
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
		return handleReadError(d, "Group Bare Metal Server Interface Members", groupID, err)
	}

	d.SetId(groupID)

	var items []map[string]interface{}
	for _, iface := range obj.Results {
		elem := make(map[string]interface{})
		elem["external_id"] = iface.ExternalId
		elem["display_name"] = iface.DisplayName
		elem["bms_external_id"] = iface.BmsExternalId
		elem["source_id"] = iface.SourceId
		elem["mac_address"] = iface.MacAddress
		if iface.State != nil {
			elem["state"] = *iface.State
		} else {
			elem["state"] = "" // Set empty string when state is not available
		}
		elem["last_sync_time"] = iface.LastSyncTime

		if iface.IsMgmtInterface != nil {
			elem["is_mgmt_interface"] = *iface.IsMgmtInterface
		}

		if iface.IpAddresses != nil {
			elem["ip_addresses"] = iface.IpAddresses
		}

		items = append(items, elem)
	}
	d.Set("items", items)

	return nil
}
