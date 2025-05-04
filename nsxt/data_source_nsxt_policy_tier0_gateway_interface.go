// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyTier0GatewayInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyTier0GatewayInterfaceRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"t0_gateway_path": {
				Type:        schema.TypeString,
				Description: "The name of the Tier0 gateway where the interface is linked",
				Required:    true,
			},
			"path": getPathSchema(),
			"edge_cluster_path": {
				Type:        schema.TypeString,
				Description: "The path of the edge cluster connected to the Tier0 gateway linked to this interface",
				Optional:    true,
				Computed:    true,
			},
			"segment_path": {
				Type:        schema.TypeString,
				Description: "Policy path for segment to be connected with the Gateway.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyTier0GatewayInterfaceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()
	t0Gw := d.Get("t0_gateway_path").(string)
	query := make(map[string]string)
	query["parent_path"] = t0Gw + "/locale-services/*"
	obj, err := policyDataSourceResourceRead(d, connector, getSessionContext(d, m), "Tier0Interface", query)
	if err != nil {
		return err
	}
	dataValue, errors := converter.ConvertToGolang(obj, model.Tier0InterfaceBindingType())
	if len(errors) > 0 {
		return errors[0]
	}

	currGwInterface := dataValue.(model.Tier0Interface)
	if currGwInterface.Path != nil {
		err := d.Set("path", *currGwInterface.Path)
		if err != nil {
			return fmt.Errorf("Error while setting interface path : %v", err)
		}
	}
	if isPolicyGlobalManager(m) {
		d.Set("edge_cluster_path", "")
	} else if currGwInterface.EdgePath != nil {
		err = d.Set("edge_cluster_path", *currGwInterface.EdgePath)
		if err != nil {
			return fmt.Errorf("Error while setting the interface edge cluster path : %v", err)
		}
	}
	if currGwInterface.Description != nil {
		err = d.Set("description", *currGwInterface.Description)
		if err != nil {
			return fmt.Errorf("Error while setting the interface description : %v", err)
		}
	}
	if currGwInterface.SegmentPath != nil {
		err = d.Set("segment_path", *currGwInterface.SegmentPath)
		if err != nil {
			return fmt.Errorf("Error while setting the segment connected to the interface : %v", err)
		}
	}
	d.SetId(newUUID())
	return nil
}
