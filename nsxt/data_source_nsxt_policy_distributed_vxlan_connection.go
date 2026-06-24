// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtPolicyDistributedVxlanConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyDistributedVxlanConnectionRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"path":         getPathSchema(),
			"display_name": getOptionalDisplayNameSchema(true),
			"description":  getDescriptionSchema(),
		},
	}
}

func dataSourceNsxtPolicyDistributedVxlanConnectionRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	_, err := policyDataSourceResourceReadWithValidation(d, connector, getSessionContext(d, m), "DistributedVxlanConnection", nil, false)
	if err != nil {
		return err
	}

	return nil
}
