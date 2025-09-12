// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtPolicyShare() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyShareRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDescriptionSchema(),
		},
	}
}

func dataSourceNsxtPolicyShareRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	_, err := policyDataSourceResourceRead(d, connector, getSessionContext(d, m), "Share", nil)
	if isNotFoundError(err) {
		return fmt.Errorf("Share with name '%s' was not found", d.Get("display_name").(string))
	} else if err != nil {
		return fmt.Errorf("encountered an error while searching Share with name '%s', error is %v", d.Get("display_name").(string), err)
	}
	return nil
}
