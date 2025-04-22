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
			"t1_gateway_name": {
				Type:        schema.TypeString,
				Description: "The name of the Tier1 gateway where the interface is linked",
				Required:    true,
			},
			"path": getPathSchema(),
			"segment_path": {
				Type:        schema.TypeString,
				Description: "The path of the edge cluster connected to the Tier1 gateway linked to this interface",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyTier1GatewayInterfaceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()
	interfaceName := d.Get("display_name").(string)
	interfaceId := d.Get("id").(string)
	t1Gw := d.Get("t1_gateway_name").(string)

	// Get the T1 gateway ID as in some case the gw name and the id are different
	t1GwObj := GenericResourceData{
		DisplayName: t1Gw,
		Id:          "",
	}
	gwObjList, err := t1GwObj.policyGenericDataSourceResourceRead(connector, getSessionContext(d, m), "Tier1", nil)
	if err != nil {
		return err
	}
	t1GwList := []model.Tier1{}
	converter = bindings.NewTypeConverter()
	for _, obj := range gwObjList {
		dataValue, errors := converter.ConvertToGolang(obj, model.Tier1BindingType())
		if len(errors) > 0 {
			return errors[0]
		}
		curGw := dataValue.(model.Tier1)
		if *curGw.DisplayName == t1Gw {
			t1GwList = append(t1GwList, dataValue.(model.Tier1))
		}
	}

	// Get the interface
	interfaceObj := GenericResourceData{
		DisplayName: interfaceName,
		Id:          interfaceId,
	}
	interfaceObjList, err := interfaceObj.policyGenericDataSourceResourceRead(connector, getSessionContext(d, m), "Tier1Interface", nil)
	if err != nil {
		return err
	}

	isOp := false
	for _, obj := range interfaceObjList {
		dataValue, errors := converter.ConvertToGolang(obj, model.Tier1InterfaceBindingType())
		if len(errors) > 0 {
			return errors[0]
		}
		currInt := dataValue.(model.Tier1Interface)
		isT0, gwID, _, _ := parseGatewayInterfacePolicyPath(*currInt.Path)

		if !isT0 && gwID == *t1GwList[0].Id {
			isOp = true
			if currInt.Path != nil {
				err := d.Set("path", *currInt.Path)
				if err != nil {
					return fmt.Errorf("Error while setting interface path : %v", err)
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
		return fmt.Errorf("The T1 gateway %s doesn't have a linked interface with name %s.", t1Gw, interfaceName)
	}
	return nil
}
