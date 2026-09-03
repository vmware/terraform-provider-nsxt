// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	connector := getPolicyConnector(m)
	objID := d.Get("id").(string)

	if _, ok := cacheAwareDataSourceReadByID[model.VpcSubnet](d, m, connector, objID, resourceTypeVpcSubnet, model.VpcSubnetBindingType()); ok {
		return nil
	}

	_, err := policyDataSourceResourceReadWithValidation(d, connector, getSessionContext(d, m), "VpcSubnet", nil, false)
	if err != nil {
		return err
	}

	return nil
}
