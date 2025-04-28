// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	// t0interface "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s"
)

func dataSourceNsxtPolicyTier0GatewayInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyTier0GatewayInterfaceRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"t0_gateway_name": {
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
				Description: "The path of the edge cluster connected to the Tier0 gateway linked to this interface",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyTier0GatewayInterfaceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()
	interfaceName := d.Get("display_name").(string)
	interfaceId := d.Get("id").(string)
	t0Gw := d.Get("t0_gateway_name").(string)

	// Get the T0 gateway ID as in some case the gw name and the id are different
	t0GwObj := GenericResourceData{
		DisplayName: t0Gw,
		Id:          "",
	}
	gwObjList, err := t0GwObj.policyGenericDataSourceResourceRead(connector, getSessionContext(d, m), "Tier0", nil)
	if err != nil {
		return err
	}
	t0GwList := []model.Tier0{}
	for _, obj := range gwObjList {
		dataValue, errors := converter.ConvertToGolang(obj, model.Tier0BindingType())
		if len(errors) > 0 {
			return errors[0]
		}
		curGw := dataValue.(model.Tier0)
		if *curGw.DisplayName == t0Gw {
			t0GwList = append(t0GwList, dataValue.(model.Tier0))
		}
	}

	// Get the interface
	interfaceObj := GenericResourceData{
		DisplayName: interfaceName,
		Id:          interfaceId,
	}
	interfaceObjList, err := interfaceObj.policyGenericDataSourceResourceRead(connector, getSessionContext(d, m), "Tier0Interface", nil)
	if err != nil {
		return err
	}

	isOp := false
	for _, obj := range interfaceObjList {
		dataValue, errors := converter.ConvertToGolang(obj, model.Tier0InterfaceBindingType())
		if len(errors) > 0 {
			return errors[0]
		}
		currInt := dataValue.(model.Tier0Interface)
		isT0, gwID, _, _ := parseGatewayInterfacePolicyPath(*currInt.Path)

		if isT0 && gwID == *t0GwList[0].Id {
			isOp = true
			if currInt.Path != nil {
				err := d.Set("path", *currInt.Path)
				if err != nil {
					return fmt.Errorf("Error while setting interface path : %v", err)
				}
			}
			if isPolicyGlobalManager(m) {
				d.Set("edge_cluster_path", "")
			} else if currInt.EdgePath != nil {
				err = d.Set("edge_cluster_path", *currInt.EdgePath)
				if err != nil {
					return fmt.Errorf("Error while setting the interface edge cluster path : %v", err)
				}
			}
			if currInt.Description != nil {
				err = d.Set("description", *currInt.Description)
				if err != nil {
					return fmt.Errorf("Error while setting the interface description : %v", err)
				}
			}
			if currInt.SegmentPath != nil {
				err = d.Set("segment_path", *currInt.SegmentPath)
				if err != nil {
					return fmt.Errorf("Error while setting the segment connected to the interface : %v", err)
				}
			}
			d.SetId(newUUID())
			break
		}
	}
	if !isOp {
		return fmt.Errorf("The T0 gateway %s doesn't have a linked interface with name %s.", t0Gw, interfaceName)
	}
	return nil

}
