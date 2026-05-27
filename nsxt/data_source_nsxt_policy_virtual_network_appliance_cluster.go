// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func dataSourceNsxtPolicyVirtualNetworkApplianceCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyVirtualNetworkApplianceClusterRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"site_path": {
				Type:         schema.TypeString,
				Description:  "Path to the site this Virtual Network Appliance Cluster belongs to",
				Optional:     true,
				Default:      defaultInfraSitePath,
				ValidateFunc: validatePolicyPath(),
			},
			"enforcement_point": {
				Type:        schema.TypeString,
				Description: "ID of the enforcement point this Virtual Network Appliance Cluster belongs to",
				Optional:    true,
				Default:     "default",
			},
			"appliance_form_factor": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Form factor for virtual network appliances in this cluster",
			},
			"service_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service type of the cluster",
			},
		},
	}
}

func listVirtualNetworkApplianceClusters(siteID, epID string, connector client.Connector, sessionContext utl.SessionContext) ([]model.VirtualNetworkApplianceCluster, error) {
	client := cliVNAClustersClient(sessionContext, connector)
	if client == nil {
		return nil, policyResourceNotSupportedError()
	}

	var results []model.VirtualNetworkApplianceCluster
	boolFalse := false
	var cursor *string
	total := 0

	for {
		listResult, err := client.List(siteID, epID, cursor, &boolFalse, nil, nil, nil, nil)
		if err != nil {
			return results, err
		}
		results = append(results, listResult.Results...)
		if total == 0 && listResult.ResultCount != nil {
			total = int(*listResult.ResultCount)
		}
		cursor = listResult.Cursor
		if len(results) >= total {
			return results, nil
		}
	}
}

func dataSourceNsxtPolicyVirtualNetworkApplianceClusterRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)

	sitePath := d.Get("site_path").(string)
	siteID := getResourceIDFromResourcePath(sitePath, "sites")
	if siteID == "" {
		return fmt.Errorf("error obtaining Site ID from site path %s", sitePath)
	}
	epID := d.Get("enforcement_point").(string)
	if epID == "" {
		epID = getPolicyEnforcementPoint(m)
	}

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)

	var obj model.VirtualNetworkApplianceCluster

	client := cliVNAClustersClient(sessionContext, connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	if objID != "" {
		objGet, err := client.Get(siteID, epID, objID)
		if err != nil {
			return handleDataSourceReadError(d, "VirtualNetworkApplianceCluster", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("error obtaining VirtualNetworkApplianceCluster ID or name during read")
	} else {
		objList, err := listVirtualNetworkApplianceClusters(siteID, epID, connector, sessionContext)
		if err != nil {
			return handleListError("VirtualNetworkApplianceCluster", err)
		}

		var perfectMatch []model.VirtualNetworkApplianceCluster
		var prefixMatch []model.VirtualNetworkApplianceCluster
		for _, objInList := range objList {
			if objInList.DisplayName == nil {
				continue
			}
			if strings.HasPrefix(*objInList.DisplayName, objName) {
				prefixMatch = append(prefixMatch, objInList)
			}
			if *objInList.DisplayName == objName {
				perfectMatch = append(perfectMatch, objInList)
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("found multiple VirtualNetworkApplianceClusters with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("found multiple VirtualNetworkApplianceClusters with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("VirtualNetworkApplianceCluster with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	d.Set("appliance_form_factor", obj.ApplianceFormFactor)
	d.Set("service_type", obj.ServiceType)

	return nil
}
