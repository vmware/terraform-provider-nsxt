// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyEdgeHighAvailabilityProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyEdgeHighAvailabilityProfileRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"site_path": {
				Type:         schema.TypeString,
				Description:  "Path of the site this Transport Node Collection belongs to",
				Optional:     true,
				ValidateFunc: validatePolicyPath(),
			},
			"unique_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique identifier assigned by the system",
			},
		},
	}
}

func dataSourceNsxtPolicyEdgeHighAvailabilityProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := enforcement_points.NewEdgeClusterHighAvailabilityProfilesClient(connector)

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	objSitePath := d.Get("site_path").(string)
	objSiteID := defaultSite
	// For local manager, if site path is not provided, use default site
	if objSitePath != "" {
		objSiteID = getPolicyIDFromPath(objSitePath)
	}
	var obj model.PolicyEdgeHighAvailabilityProfile
	if objID != "" {
		// Get by id
		objGet, err := client.Get(objSiteID, getPolicyEnforcementPoint(m), objID)
		if err != nil {
			return handleDataSourceReadError(d, "EdgeHighAvailabilityProfile", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining EdgeHighAvailabilityProfile ID or name during read")
	} else {
		// Get by full name/prefix
		objList, err := client.List(objSiteID, getPolicyEnforcementPoint(m), nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return handleListError("EdgeHighAvailabilityProfile", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.PolicyEdgeHighAvailabilityProfile
		var prefixMatch []model.PolicyEdgeHighAvailabilityProfile
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
				return fmt.Errorf("Found multiple EdgeHighAvailabilityProfile with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple EdgeHighAvailabilityProfiles with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("EdgeHighAvailabilityProfile with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	d.Set("unique_id", obj.UniqueId)

	return nil
}
