/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points/edge_clusters"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"strings"
)

func dataSourceNsxtPolicyEdgeNode() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyEdgeNodeRead,

		Schema: map[string]*schema.Schema{
			"edge_cluster_path": getPolicyPathSchema(true, false, "Edge cluster Path"),
			"member_index": {
				Type:         schema.TypeInt,
				Description:  "Index of this node within edge cluster",
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
		},
	}
}

func dataSourceNsxtPolicyEdgeNodeRead(d *schema.ResourceData, m interface{}) error {
	// Read an edge node by name or id
	connector := getPolicyConnector(m)
	client := edge_clusters.NewDefaultEdgeNodesClient(connector)

	edgeClusterPath := d.Get("edge_cluster_path").(string)
	edgeClusterID := getPolicyIDFromPath(edgeClusterPath)
	objID := d.Get("id").(string)
	name, nameSet := d.GetOkExists("display_name")
	objName := name.(string)
	memberIndex, memberIndexSet := d.GetOkExists("member_index")
	objMemberIndex := int64(memberIndex.(int))
	var obj model.PolicyEdgeNode
	if objID != "" {
		// Get by id
		objGet, err := client.Get(defaultSite, getPolicyEnforcementPoint(m), edgeClusterID, objID)

		if err != nil {
			return handleDataSourceReadError(d, "Edge Node", objID, err)
		}
		obj = objGet
	} else {
		// Get by full name/prefix
		includeMarkForDeleteObjectsParam := false
		objList, err := client.List(defaultSite, getPolicyEnforcementPoint(m), edgeClusterID, nil, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
		if err != nil {
			return handleListError("Edge Node", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.PolicyEdgeNode
		var prefixMatch []model.PolicyEdgeNode
		for _, objInList := range objList.Results {
			indexMatch := true
			if memberIndexSet && objMemberIndex != *objInList.MemberIndex {
				indexMatch = false
			}
			if nameSet && strings.HasPrefix(*objInList.DisplayName, objName) && indexMatch {
				prefixMatch = append(prefixMatch, objInList)
			}
			if nameSet && *objInList.DisplayName == objName && indexMatch {
				perfectMatch = append(perfectMatch, objInList)
			}
			if !nameSet && indexMatch {
				perfectMatch = append(perfectMatch, objInList)
			}
		}

		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple edge nodes with name '%s' and index %d", objName, memberIndex)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple edge nodes with name starting with '%s' and index %d", objName, memberIndex)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("edge node '%s' was not found and %d", objName, memberIndex)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	return nil
}
