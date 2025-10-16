// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var vpcNatTypes = []string{
	model.PolicyNat_NAT_TYPE_INTERNAL,
	model.PolicyNat_NAT_TYPE_USER,
	model.PolicyNat_NAT_TYPE_DEFAULT,
	model.PolicyNat_NAT_TYPE_NAT64,
}

func dataSourceNsxtVpcNat() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcNatRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"nat_type": {
				Type:         schema.TypeString,
				Description:  "Nat Type",
				Required:     true,
				ValidateFunc: validation.StringInSlice(vpcNatTypes, false),
			},
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
		},
	}
}

func dataSourceNsxtVpcNatRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC NAT data source requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)

	natType := d.Get("nat_type").(string)
	query := make(map[string]string)
	query["nat_type"] = natType

	_, err := policyDataSourceResourceReadWithValidation(d, connector, commonSessionContext, "PolicyNat", query, false)
	if err != nil {
		return err
	}

	return nil
}
