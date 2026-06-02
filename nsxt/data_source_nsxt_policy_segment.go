// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func dataSourceNsxtPolicySegment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicySegmentRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaWithSpec(utl.SessionContextSpec{IsRequired: false, IsComputed: false, IsVpc: false, AllowDefaultProject: false, FromGlobal: true}),
			"transport_zone_path": {
				Type:        schema.TypeString,
				Description: "Policy path to the transport zone",
				Computed:    true,
			},
			"connectivity_path": {
				Type:        schema.TypeString,
				Description: "Policy path to the connecting Tier-0 or Tier-1 gateway",
				Computed:    true,
			},
			"vlan_ids": {
				Type:        schema.TypeList,
				Description: "VLAN IDs for a VLAN-backed segment",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"subnet": {
				Type:        schema.TypeList,
				Description: "Subnet configuration",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:        schema.TypeString,
							Description: "Gateway IP address in CIDR format",
							Computed:    true,
						},
						"network": {
							Type:        schema.TypeString,
							Description: "Network CIDR for this subnet",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNsxtPolicySegmentRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	obj, err := policyDataSourceResourceRead(d, connector, getSessionContext(d, m), "Segment", nil)
	if err != nil {
		return err
	}

	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToGolang(obj, model.SegmentBindingType())
	if len(errors) > 0 {
		return errors[0]
	}
	segment := dataValue.(model.Segment)

	d.Set("transport_zone_path", segment.TransportZonePath)
	d.Set("connectivity_path", segment.ConnectivityPath)
	d.Set("vlan_ids", segment.VlanIds)

	var subnets []interface{}
	for _, s := range segment.Subnets {
		entry := make(map[string]interface{})
		entry["cidr"] = ""
		entry["network"] = ""
		if s.GatewayAddress != nil {
			entry["cidr"] = *s.GatewayAddress
		}
		if s.Network != nil {
			entry["network"] = *s.Network
		}
		subnets = append(subnets, entry)
	}
	d.Set("subnet", subnets)

	return nil
}
