// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyIPPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIPPoolRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchema(false, false, false),
			"realized_id": {
				Type:        schema.TypeString,
				Description: "The ID of the realized resource",
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyIPPoolRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	objID := d.Get("id").(string)

	if obj, ok := cacheAwareDataSourceReadByID[model.IpAddressPool](d, m, connector, objID, resourceTypeIpAddressPool, model.IpAddressPoolBindingType()); ok {
		d.Set("realized_id", obj.RealizationId)
		return nil
	}

	obj, err := policyDataSourceResourceRead(d, connector, getSessionContext(d, m), "IpAddressPool", nil)
	if err != nil {
		return err
	}

	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToGolang(obj, model.IpAddressPoolBindingType())
	if len(errors) > 0 {
		return errors[0]
	}
	poolObj := dataValue.(model.IpAddressPool)
	d.Set("realized_id", poolObj.RealizationId)
	return nil
}
