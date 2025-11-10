// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func dataSourceNsxtVpcSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcSubnetRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
		},
	}
}

func dataSourceNsxtVpcSubnetRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC Subnet data source requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)

	_, err := policyDataSourceResourceReadWithValidation(d, connector, getSessionContext(d, m), "VpcSubnet", nil, false)
	if err != nil {
		return err
	}

	return nil
}
