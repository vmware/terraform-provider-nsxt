/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyTier0Gateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyTier0GatewayRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"edge_cluster_path": {
				Type:        schema.TypeString,
				Description: "The path of the edge cluster connected to this Tier0 gateway",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyTier0GatewayRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	if isPolicyGlobalManager(m) {
		_, err := policyDataSourceResourceRead(d, connector, "Tier0", nil)
		if err != nil {
			return err
		}

		// Single edgle cluster is not informative for global manager
		d.Set("edge_cluster_path", "")
		return nil
	}

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	client := infra.NewDefaultTier0sClient(connector)
	var obj model.Tier0
	if objID != "" {
		// Get by id
		objGet, err := client.Get(objID)
		if isNotFoundError(err) {
			return fmt.Errorf("Tier0 with ID %s was not found", objID)
		}

		if err != nil {
			return fmt.Errorf("Error while reading Tier0 %s: %v", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining Tier0 ID or name during read")
	} else {
		// Get by full name/prefix
		includeMarkForDeleteObjectsParam := false
		objList, err := client.List(nil, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("Error while reading Tier0s: %v", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.Tier0
		var prefixMatch []model.Tier0
		for _, objInList := range objList.Results {
			if strings.HasPrefix(*objInList.DisplayName, objName) {
				prefixMatch = append(prefixMatch, objInList)
			}
			if *objInList.DisplayName == objName {
				perfectMatch = append(perfectMatch, objInList)
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple Tier0s with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple Tier0s with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("Tier0 with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)

	localeServices, err := listPolicyTier0GatewayLocaleServices(connector, *obj.Id, false)
	if err != nil {
		return fmt.Errorf("Failed to read locale services for '%s'", objName)
	}
	for _, service := range localeServices {
		if service.EdgeClusterPath != nil {
			d.Set("edge_cluster_path", service.EdgeClusterPath)
			break
		}
	}
	return nil
}
