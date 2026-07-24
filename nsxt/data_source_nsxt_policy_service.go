// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyServiceRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchema(false, false, false),
		},
	}
}

func dataSourceNsxtPolicyServiceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	objID := d.Get("id").(string)

	if _, ok := cacheAwareDataSourceReadByID[model.Service](d, m, connector, objID, resourceTypeService, model.ServiceBindingType()); ok {
		return nil
	}

	_, err := policyDataSourceResourceRead(d, connector, getSessionContext(d, m), "Service", nil)
	if err != nil {
		return err
	}
	return nil
}
