// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyVirtualNetworkAppliance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyVirtualNetworkApplianceRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"cluster_path": {
				Type:         schema.TypeString,
				Description:  "Policy path of the parent Virtual Network Appliance Cluster",
				Required:     true,
				ValidateFunc: validatePolicyPath(),
			},
			"hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Hostname or FQDN of the Virtual Network Appliance VM",
			},
			"failure_domain_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy path of the failure domain",
			},
		},
	}
}

func dataSourceNsxtPolicyVirtualNetworkApplianceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	clusterPath := d.Get("cluster_path").(string)
	siteID, epID, clusterID, err := getVNAClusterPathComponents(clusterPath)
	if err != nil {
		return err
	}

	vnaClient := cliVNACRUDClient(sessionContext, connector)
	if vnaClient == nil {
		return policyResourceNotSupportedError()
	}

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)

	var obj model.VirtualNetworkAppliance

	if objID != "" {
		objGet, err := vnaClient.Get(siteID, epID, clusterID, objID)
		if err != nil {
			return handleDataSourceReadError(d, "VirtualNetworkAppliance", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("error obtaining VirtualNetworkAppliance ID or display_name during read")
	} else {
		listResult, err := vnaClient.List(siteID, epID, clusterID)
		if err != nil {
			return handleListError("VirtualNetworkAppliance", err)
		}

		var perfectMatch []model.VirtualNetworkAppliance
		var prefixMatch []model.VirtualNetworkAppliance
		for _, item := range listResult.Results {
			if item.DisplayName == nil {
				continue
			}
			if strings.HasPrefix(*item.DisplayName, objName) {
				prefixMatch = append(prefixMatch, item)
			}
			if *item.DisplayName == objName {
				perfectMatch = append(perfectMatch, item)
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("found multiple VirtualNetworkAppliances with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("found multiple VirtualNetworkAppliances with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("VirtualNetworkAppliance with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	d.Set("hostname", obj.Hostname)
	d.Set("failure_domain_path", obj.FailureDomainPath)

	return nil
}
