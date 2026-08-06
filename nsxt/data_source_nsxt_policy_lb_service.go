// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	if _, ok := cacheAwareDataSourceReadByID[model.LBService](d, m, connector, objID, resourceTypeLBService, model.LBServiceBindingType()); ok {
		return nil
	}

	_, err := policyDataSourceResourceRead(d, connector, getSessionContext(d, m), "LBService", nil)
	if err != nil {
		return err
	}

	return nil
}
