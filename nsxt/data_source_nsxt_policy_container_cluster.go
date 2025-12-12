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

	enforcement_points "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func dataSourceNsxtPolicyContainerCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyContainerClusterRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
		},
	}
}

func listClusterControlPlanes(siteID, epID string, connector client.Connector, sessionContext utl.SessionContext) ([]model.ClusterControlPlane, error) {
	client := enforcement_points.NewClusterControlPlanesClient(sessionContext, connector)
	if client == nil {
		return nil, policyResourceNotSupportedError()
	}

	var results []model.ClusterControlPlane
	boolFalse := false
	var cursor *string
	total := 0

	for {
		policies, err := client.List(siteID, epID, cursor, nil, nil, nil, &boolFalse, nil)
		if err != nil {
			return results, err
		}
		results = append(results, policies.Results...)
		if total == 0 && policies.ResultCount != nil {
			// first response
			total = int(*policies.ResultCount)
		}

		cursor = policies.Cursor
		if len(results) >= total {
			return results, nil
		}
	}
}

func dataSourceNsxtPolicyContainerClusterRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	client := enforcement_points.NewClusterControlPlanesClient(sessionContext, connector)

	// As Project resource type paths reside under project and not under /infra or /global_infra or such, and since
	// this data source fetches extra attributes, e.g site_info and tier0_gateway_paths, it's simpler to implement using .List()
	// instead of using search API.

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	var obj model.ClusterControlPlane

	if objID != "" {
		// Get by id
		objGet, err := client.Get(defaultSite, defaultEnforcementPoint, objID)
		if err != nil {
			return handleDataSourceReadError(d, "ClusterControlPlane", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining ClusterControlPlane ID or name during read")
	} else {
		// Get by full name/prefix
		objList, err := listClusterControlPlanes(defaultSite, defaultEnforcementPoint, connector, sessionContext)
		if err != nil {
			return handleListError("ClusterControlPlane", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.ClusterControlPlane
		var prefixMatch []model.ClusterControlPlane
		for _, objInList := range objList {
			if strings.HasPrefix(*objInList.DisplayName, objName) {
				prefixMatch = append(prefixMatch, objInList)
			}
			if *objInList.DisplayName == objName {
				perfectMatch = append(perfectMatch, objInList)
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple ClusterControlPlanes with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple ClusterControlPlanes with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("ClusterControlPlane with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)

	return nil
}
