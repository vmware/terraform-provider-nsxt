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

func dataSourceNsxtVPC() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVPCRead,
		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, false, true),
			"short_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNsxtVPCRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC data source requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)
	objID := d.Get("id").(string)
	displayName := d.Get("display_name").(string)
	lookupKey := objID
	if lookupKey == "" {
		lookupKey = displayName
	}
	if lookupKey != "" && IsCacheEnabled() {
		val, err := gcache.readCache(lookupKey, "Vpc", d, m, connector)
		if err == nil {
			converter := bindings.NewTypeConverter()
			goVal, convErrs := converter.ConvertToGolang(val.(*data.StructValue), model.VpcBindingType())
			if len(convErrs) == 0 {
				obj, ok := goVal.(model.Vpc)
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
					d.Set("short_id", obj.ShortId)
					return nil
				}
			}
		}
	}

	obj, err := policyDataSourceResourceRead(d, connector, getSessionContext(d, m), "Vpc", nil)
	if err != nil {
		return err
	}

	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToGolang(obj, model.VpcBindingType())
	if len(errors) > 0 {
		return errors[0]
	}
	vpc := dataValue.(model.Vpc)

	d.Set("short_id", vpc.ShortId)

	return nil
}
