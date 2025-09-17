// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
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
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func dataSourceNsxtVpcServiceProfileRead(d *schema.ResourceData, m interface{}) error {
	// Using deprecated API because GetOk is not behaving as expected when is_default = "false".
	// It does not return true for a key that's explicitly set to false.
	value, defaultOK := d.GetOkExists("is_default")
	_, dpOk := d.GetOk("display_name")
	_, idOk := d.GetOk("id")

	if defaultOK && !dpOk && !idOk {
		query := make(map[string]string)
		query["is_default"] = strconv.FormatBool(value.(bool))
		_, err := policyDataSourceReadWithCustomField(d, getPolicyConnector(m), getSessionContext(d, m), "VpcServiceProfile", query)
		return err
	}
	obj, err := policyDataSourceResourceRead(d, getPolicyConnector(m), getSessionContext(d, m), "VpcServiceProfile", nil)
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
