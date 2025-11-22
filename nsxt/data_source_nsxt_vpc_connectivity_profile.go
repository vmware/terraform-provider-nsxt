// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtVpcConnectivityProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcConnectivityProfileRead,

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

func dataSourceNsxtVpcConnectivityProfileRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC Connectivity Profile data source requires NSX version 9.0.0 or higher")
	}
	// Using deprecated API because GetOk is not behaving as expected when is_default = "false".
	// It does not return true for a key that's explicitly set to false.
	value, defaultOK := d.GetOkExists("is_default")
	if defaultOK {
		query := make(map[string]string)
		query["is_default"] = strconv.FormatBool(value.(bool))
		_, err := policyDataSourceReadWithCustomField(d, getPolicyConnector(m), getSessionContext(d, m), "VpcConnectivityProfile", query)
		return err
	}
	obj, err := policyDataSourceResourceRead(d, getPolicyConnector(m), getSessionContext(d, m), "VpcConnectivityProfile", nil)
	if err != nil {
		return err
	}

	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToGolang(obj, model.VpcConnectivityProfileBindingType())
	if len(errors) > 0 {
		return errors[0]
	}

	vpcConProfile := dataValue.(model.VpcConnectivityProfile)
	d.Set("is_default", vpcConProfile.IsDefault)
	return nil
}
