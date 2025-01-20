// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	domain := d.Get("domain").(string)
	query := make(map[string]string)
	query["parent_path"] = "*/" + domain
	_, err := policyDataSourceResourceRead(d, getPolicyConnector(m), getSessionContext(d, m), "Group", query)
	if err != nil {
		return err
	}
	return nil
}
