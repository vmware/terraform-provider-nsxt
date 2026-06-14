// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtPolicyRouteControllerInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyRouteControllerInterfaceRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"parent_path":  getPolicyPathSchema(false, false, "Policy path of the parent Route Controller"),
			"path":         getPathSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDescriptionSchema(),
		},
	}
}

func dataSourceNsxtPolicyRouteControllerInterfaceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	parentPath := d.Get("parent_path").(string)
	query := make(map[string]string)
	if len(parentPath) > 0 {
		if _, err := parseStandardPolicyPathVerifySize(parentPath, 1, routeControllerInterfacePathExample); err != nil {
			return fmt.Errorf("invalid parent_path: %w", err)
		}
		query["parent_path"] = fmt.Sprintf("%s*", parentPath)
	}

	_, err := policyDataSourceResourceReadWithValidation(d, connector, getSessionContext(d, m), "RouteControllerInterface", query, false)
	if err != nil {
		return err
	}

	return nil
}
