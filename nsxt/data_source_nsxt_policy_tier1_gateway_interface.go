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

func dataSourceNsxtPolicyTier1GatewayInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyTier1GatewayInterfaceRead,
		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"t1_gateway_path": {
				Type:        schema.TypeString,
				Description: "The name of the Tier1 gateway where the interface is linked",
				Required:    true,
			},
			"path": getPathSchema(),
			"segment_path": {
				Type:        schema.TypeString,
				Description: "Policy path for segment to be connected with the Gateway.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyTier1GatewayInterfaceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()
	t1Gw := d.Get("t1_gateway_path").(string)
	query := make(map[string]string)
	query["parent_path"] = t1Gw + "/locale-services/*"
	obj, err := policyDataSourceResourceRead(d, connector, getSessionContext(d, m), "Tier1Interface", query)
	if err != nil {
		return err
	}

	dataValue, errors := converter.ConvertToGolang(obj, model.Tier1InterfaceBindingType())
	if len(errors) > 0 {
		return errors[0]
	}

	currGwInterface := dataValue.(model.Tier1Interface)
	if currGwInterface.Path != nil {
		err := d.Set("path", *currGwInterface.Path)
		if err != nil {
			return fmt.Errorf("Error while setting interface path : %v", err)
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
