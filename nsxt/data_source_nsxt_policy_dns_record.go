// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNsxtPolicyDnsRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyDnsRecordRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, false, true),
			"zone_path":    getPolicyPathSchema(false, false, "Policy path of the parent DNS Zone"),
		},
	}
}

func dataSourceNsxtPolicyDnsRecordRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	zonePath := d.Get("zone_path").(string)
	query := make(map[string]string)
	if len(zonePath) > 0 {
		query["zone_path"] = escapeSpecialCharacters(zonePath)
	}
	_, err := policyDataSourceResourceReadWithValidation(d, connector, getSessionContext(d, m), "DnsRecord", query, false)
	if err != nil {
		return err
	}

	return nil
}
