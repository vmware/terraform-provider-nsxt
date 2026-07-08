// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyLbService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyLbServiceRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
		},
	}
}

func dataSourceNsxtPolicyLbServiceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	objID := d.Get("id").(string)
	displayName := d.Get("display_name").(string)
	lookupKey := objID
	if lookupKey == "" {
		lookupKey = displayName
	}

	if lookupKey != "" && IsCacheEnabled() {
		val, err := gcache.readCache(lookupKey, "LBService", d, m, connector)
		if err == nil {
			converter := bindings.NewTypeConverter()
			goVal, convErrs := converter.ConvertToGolang(val.(*data.StructValue), model.LBServiceBindingType())
			if len(convErrs) == 0 {
				obj, ok := goVal.(model.LBService)
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

	_, err := policyDataSourceResourceRead(d, connector, getSessionContext(d, m), "LBService", nil)
	if err != nil {
		return err
	}

	return nil
}
