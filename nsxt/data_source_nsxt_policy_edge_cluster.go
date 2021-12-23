/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyEdgeCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyEdgeClusterRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"site_path": {
				Type:         schema.TypeString,
				Description:  "Path of the site this Edge cluster belongs to",
				Optional:     true,
				ValidateFunc: validatePolicyPath(),
			},
		},
	}
}

func dataSourceNsxtPolicyEdgeClusterRead(d *schema.ResourceData, m interface{}) error {
	// Read an edge cluster by name or id
	objSitePath := d.Get("site_path").(string)

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)

	if !isPolicyGlobalManager(m) && objSitePath != "" {
		return globalManagerOnlyError()
	}
	if isPolicyGlobalManager(m) {
		if objSitePath == "" {
			return attributeRequiredGlobalManagerError("site_path", "nsxt_policy_edge_cluster")
		}

		query := make(map[string]string)
		globalPolicyEnforcementPointPath := getGlobalPolicyEnforcementPointPath(m, &objSitePath)
		query["parent_path"] = globalPolicyEnforcementPointPath
		_, err := policyDataSourceResourceReadWithValidation(d, getPolicyConnector(m), true, "PolicyEdgeCluster", query, false)
		if err != nil {
			return err
		}
		return nil
	}

	// Local manager
	connector := getPolicyConnector(m)
	client := enforcement_points.NewEdgeClustersClient(connector)
	var obj model.PolicyEdgeCluster
	if objID != "" {
		// Get by id
		objGet, err := client.Get(defaultSite, getPolicyEnforcementPoint(m), objID)

		if err != nil {
			return handleDataSourceReadError(d, "Edge Cluster", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining edge cluster ID or name during read")
	} else {
		// Get by full name/prefix
		includeMarkForDeleteObjectsParam := false
		objList, err := client.List(defaultSite, getPolicyEnforcementPoint(m), nil, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
		if err != nil {
			return handleListError("Edge Cluster", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.PolicyEdgeCluster
		var prefixMatch []model.PolicyEdgeCluster
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
				return fmt.Errorf("Found multiple edge clusters with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple edge clusters with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("edge cluster '%s' was not found", objName)
		}
	}
	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	return nil
}
