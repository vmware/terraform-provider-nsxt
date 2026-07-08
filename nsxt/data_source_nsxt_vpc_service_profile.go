// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"strconv"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtVpcServiceProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcServiceProfileRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, false, true),
			"is_default": {
				Type:         schema.TypeBool,
				Optional:     true,
				ExactlyOneOf: []string{"id", "display_name", "is_default"},
			},
		},
	}
}

func dataSourceNsxtVpcServiceProfileRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC Service Profile data source requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)
	// Using deprecated API because GetOk is not behaving as expected when is_default = "false".
	// It does not return true for a key that's explicitly set to false.
	value, defaultOK := d.GetOkExists("is_default")
	if defaultOK {
		query := make(map[string]string)
		query["is_default"] = strconv.FormatBool(value.(bool))
		_, err := policyDataSourceReadWithCustomField(d, connector, getSessionContext(d, m), "VpcServiceProfile", query)
		return err
	}

	objID := d.Get("id").(string)
	displayName := d.Get("display_name").(string)
	lookupKey := objID
	if lookupKey == "" {
		lookupKey = displayName
	}
	if lookupKey != "" && IsCacheEnabled() {
		val, err := gcache.readCache(lookupKey, "VpcServiceProfile", d, m, connector)
		if err == nil {
			converter := bindings.NewTypeConverter()
			goVal, convErrs := converter.ConvertToGolang(val.(*data.StructValue), model.VpcServiceProfileBindingType())
			if len(convErrs) == 0 {
				obj, ok := goVal.(model.VpcServiceProfile)
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
					d.Set("is_default", obj.IsDefault)
					return nil
				}
			}
		}
	}

	obj, err := policyDataSourceResourceRead(d, connector, getSessionContext(d, m), "VpcServiceProfile", nil)
	if err != nil {
		return err
	}
	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToGolang(obj, model.VpcServiceProfileBindingType())
	if len(errors) > 0 {
		return errors[0]
	}
	vpcSvcProfile := dataValue.(model.VpcServiceProfile)
	d.Set("is_default", vpcSvcProfile.IsDefault)
	return nil
}
