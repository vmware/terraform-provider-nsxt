// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func dataSourceNsxtVpcServiceEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcServiceEndpointRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
			"service_endpoint_ip": {
				Type:        schema.TypeString,
				Description: "IP address of the VM providing the service",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtVpcServiceEndpointRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.2.0") {
		return fmt.Errorf("VPC Service Endpoint data source requires NSX version 9.2.0 or higher")
	}
	connector := getPolicyConnector(m)

	query := make(map[string]string)
	if ip, ok := d.GetOk("service_endpoint_ip"); ok {
		query["service_endpoint_ip"] = ip.(string)
	}

	_, err := policyDataSourceResourceReadWithValidation(d, connector, getSessionContext(d, m), "VpcServiceEndpoint", query, false)
	if err != nil {
		return err
	}

	return nil
}
