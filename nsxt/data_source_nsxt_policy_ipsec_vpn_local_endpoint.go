/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyIPSecVpnLocalEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIPSecVpnLocalEndpointRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"service_path": getPolicyPathSchema(false, false, "Policy path for IPSec VPN service"),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"local_address": {
				Type:        schema.TypeString,
				Description: "Local IPv4 IP address",
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyIPSecVpnLocalEndpointRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	servicePath := d.Get("service_path").(string)
	query := make(map[string]string)
	if len(servicePath) > 0 {
		// In newer NSX versions, NSX removes locale service from the parent path when search API is concerned
		objID := d.Get("id").(string)
		objName := d.Get("display_name").(string)
		client, err := newLocalEndpointClient(servicePath)
		if err != nil {
			return err
		}
		if objID != "" {
			obj, err := client.Get(connector, objID)
			if err != nil {
				return fmt.Errorf("Failed to locate Local Endpoint %s/%s: %v", servicePath, objID, err)
			}
			d.SetId(*obj.Id)
			d.Set("display_name", obj.DisplayName)
			d.Set("description", obj.Description)
			d.Set("path", obj.Path)
			d.Set("local_address", obj.LocalAddress)
			return nil
		}

		objList, err := client.List(connector)
		if err != nil {
			return fmt.Errorf("Failed to list local endpoints: %v", err)
		}

		for _, obj := range objList {
			if *obj.DisplayName == objName {
				d.SetId(*obj.Id)
				d.Set("display_name", obj.DisplayName)
				d.Set("description", obj.Description)
				d.Set("path", obj.Path)
				d.Set("local_address", obj.LocalAddress)
				return nil
			}
		}
		return fmt.Errorf("Failed to locate Local Endpoint under %s named %s", servicePath, objName)
	}
	objInt, err := policyDataSourceResourceReadWithValidation(d, connector, isPolicyGlobalManager(m), "IPSecVpnLocalEndpoint", query, false)
	if err != nil {
		return err
	}

	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToGolang(objInt, model.IPSecVpnLocalEndpointBindingType())
	if len(errors) > 0 {
		return fmt.Errorf("Failed to convert type for Local Endpoint: %v", errors[0])
	}
	obj := dataValue.(model.IPSecVpnLocalEndpoint)
	d.Set("local_address", obj.LocalAddress)
	return nil
}
