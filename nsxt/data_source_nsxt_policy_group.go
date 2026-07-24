// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyGroupRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"domain":       getDomainNameSchema(),
			"context":      getContextSchema(false, false, false),
		},
	}
}

func dataSourceNsxtPolicyGroupRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	objID := d.Get("id").(string)

	if _, ok := cacheAwareDataSourceReadByID[model.Group](d, m, connector, objID, resourceTypeGroup, model.GroupBindingType()); ok {
		return nil
	}

	domain := d.Get("domain").(string)
	query := make(map[string]string)
	query["parent_path"] = "*/" + domain
	_, err := policyDataSourceResourceRead(d, connector, getSessionContext(d, m), "Group", query)
	if err != nil {
		return err
	}
	return nil
}
