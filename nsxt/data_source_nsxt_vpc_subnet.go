// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtVpcSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcSubnetRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
		},
	}
}

func dataSourceNsxtVpcSubnetRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC Subnet data source requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)
	objID := d.Get("id").(string)
	displayName := d.Get("display_name").(string)
	lookupKey := objID
	if lookupKey == "" {
		lookupKey = displayName
	}

	// Try serving from cache when the caller specified an ID or display name.
	// Data sources don't participate in refresh phase detection the same way resources do,
	// so we call cache directly here while reusing the same cache storage and conversion.
	if lookupKey != "" && IsCacheEnabled() {
		val, err := gcache.readCache(lookupKey, "VpcSubnet", d, m, connector)
		if err == nil {
			converter := bindings.NewTypeConverter()
			goVal, convErrs := converter.ConvertToGolang(val.(*data.StructValue), model.VpcSubnetBindingType())
			if len(convErrs) == 0 {
				obj, ok := goVal.(model.VpcSubnet)
				if ok {
					id := lookupKey
					if obj.Id != nil {
						id = *obj.Id
					}
					d.SetId(id)
					d.Set("id", id)
					d.Set("display_name", obj.DisplayName)
					d.Set("description", obj.Description)
					d.Set("path", obj.Path)
					return nil
				}
			}
		}
	}

	_, err := policyDataSourceResourceReadWithValidation(d, connector, getSessionContext(d, m), "VpcSubnet", nil, false)
	if err != nil {
		return err
	}

	return nil
}
