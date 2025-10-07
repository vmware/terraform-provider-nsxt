// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func dataSourceNsxtVpcConnectivityProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcConnectivityProfileRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, false, true),
		},
	}
}

func dataSourceNsxtVpcConnectivityProfileRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC Connectivity Profile data source requires NSX version 9.0.0 or higher")
	}
	_, err := policyDataSourceResourceRead(d, getPolicyConnector(m), getSessionContext(d, m), "VpcConnectivityProfile", nil)
	if err != nil {
		return err
	}
	return nil
}
