// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func dataSourceNsxtVpcIpAddressAllocation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcIpAddressAllocationRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"allocation_ips": {
				Type:        schema.TypeString,
				Description: "Allocation IPs - single IP or cidr",
				Required:    true,
			},
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
		},
	}
}

func dataSourceNsxtVpcIpAddressAllocationRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC IP Address Allocation data source requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)

	ips := d.Get("allocation_ips").(string)
	query := make(map[string]string)
	query["allocation_ips"] = ips

	_, err := policyDataSourceResourceReadWithValidation(d, connector, commonSessionContext, "VpcIpAddressAllocation", query, false)
	if err != nil {
		return err
	}
	// validate that IP address indeed matches

	return nil
}
