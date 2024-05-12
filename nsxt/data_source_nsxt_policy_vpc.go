/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyVPC() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyVPCRead,
		Schema: map[string]*schema.Schema{
			"id":                   getDataSourceIDSchema(),
			"display_name":         getDataSourceDisplayNameSchema(),
			"description":          getDataSourceDescriptionSchema(),
			"path":                 getPathSchema(),
			"context":              getContextSchema(true, false),
			"site_info":            getSiteInfoSchema(),
			"default_gateway_path": getPathSchema(),
			"short_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNsxtPolicyVPCRead(d *schema.ResourceData, m interface{}) error {
	obj, err := policyDataSourceResourceRead(d, getPolicyConnector(m), getSessionContext(d, m), "Vpc", nil)
	if err != nil {
		return err
	}

	converter := bindings.NewTypeConverter()
	dataValue, errors := converter.ConvertToGolang(obj, model.VpcBindingType())
	if len(errors) > 0 {
		return errors[0]
	}
	vpc := dataValue.(model.Vpc)

	d.Set("default_gateway_path", vpc.DefaultGatewayPath)
	d.Set("short_id", vpc.ShortId)

	var siteInfosList []map[string]interface{}
	for _, item := range vpc.SiteInfos {
		data := make(map[string]interface{})
		data["edge_cluster_paths"] = item.EdgeClusterPaths
		data["site_path"] = item.SitePath
		siteInfosList = append(siteInfosList, data)
	}
	d.Set("site_info", siteInfosList)

	return nil
}
