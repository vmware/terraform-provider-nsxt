// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func dataSourceNsxtVpcEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcEndpointRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
			"vpc_service_endpoint": {
				Type:        schema.TypeString,
				Description: "Policy path to the VPC service endpoint being consumed",
				Computed:    true,
			},
			"ip_allocation_path": {
				Type:        schema.TypeString,
				Description: "Policy path to the VPC IP address allocation",
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtVpcEndpointRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.2.0") {
		return fmt.Errorf("VPC Endpoint data source requires NSX version 9.2.0 or higher")
	}
	connector := getPolicyConnector(m)

	_, err := policyDataSourceResourceReadWithValidation(d, connector, getSessionContext(d, m), "VpcEndpoint", nil, false)
	if err != nil {
		return err
	}

	return nil
}
