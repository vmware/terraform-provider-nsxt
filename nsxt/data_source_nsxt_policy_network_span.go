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

func dataSourceNsxtPolicyNetworkSpan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyNetworkSpanRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func dataSourceNsxtPolicyNetworkSpanRead(d *schema.ResourceData, m interface{}) error {
	// Using deprecated API because GetOk is not behaving as expected when is_default = "false".
	// It does not return true for a key that's explicitly set to false.
	value, defaultOK := d.GetOkExists("is_default")
	_, dpOk := d.GetOk("display_name")
	_, idOk := d.GetOk("id")

	if defaultOK && !dpOk && !idOk {
		query := make(map[string]string)
		query["is_default"] = strconv.FormatBool(value.(bool))
		_, err := policyDataSourceReadWithCustomField(d, getPolicyConnector(m), getSessionContext(d, m), "NetworkSpan", query)
		return err
	}
	obj, err := policyDataSourceResourceRead(d, getPolicyConnector(m), getSessionContext(d, m), "NetworkSpan", nil)
	if err != nil {
		return err
	}
	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToGolang(obj, model.NetworkSpanBindingType())
	if len(errors) > 0 {
		return errors[0]
	}
	netSpan := dataValue.(model.NetworkSpan)
	d.Set("is_default", netSpan.IsDefault)
	return nil
}
