/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtPolicyIPSecVpnLocalEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIPSecVpnLocalEndpointRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"service_path": getPolicyPathSchema(false, false, "Policy path for IPSec VPN service"),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
		},
	}
}

func dataSourceNsxtPolicyIPSecVpnLocalEndpointRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	servicePath := d.Get("service_path").(string)
	query := make(map[string]string)
	if len(servicePath) > 0 {
		query["parent_path"] = servicePath
	}
	_, err := policyDataSourceResourceReadWithValidation(d, connector, isPolicyGlobalManager(m), "IPSecVpnLocalEndpoint", query, false)
	if err != nil {
		return err
	}

	return nil
}
