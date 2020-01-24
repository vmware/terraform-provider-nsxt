/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func dataSourceNsxtSwitchingProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtSwitchingProfileRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Unique ID of this resource",
				Optional:    true,
				Computed:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The display name of this resource",
				Optional:    true,
				Computed:    true,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Description: "The resource type representing the specific type of this profile",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtSwitchingProfileRead(d *schema.ResourceData, m interface{}) error {
	// Read a switching profile by name or id
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return dataSourceNotSupportedError()
	}

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	var obj manager.BaseSwitchingProfile
	if objID != "" {
		// Get by id
		objGet, resp, err := nsxClient.LogicalSwitchingApi.GetSwitchingProfile(nsxClient.Context, objID)

		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("switching profile %s was not found", objID)
		}
		if err != nil {
			return fmt.Errorf("Error while reading switching profile %s: %v", objID, err)
		}
		obj = objGet
	} else if objName != "" {
		// Get by full name
		// TODO use localVarOptionals for paging
		localVarOptionals := make(map[string]interface{})
		localVarOptionals["includeSystemOwned"] = true
		objList, _, err := nsxClient.LogicalSwitchingApi.ListSwitchingProfiles(nsxClient.Context, localVarOptionals)
		if err != nil {
			return fmt.Errorf("Error while reading switching profiles: %v", err)
		}
		// go over the list to find the correct one
		found := false
		for _, objInList := range objList.Results {
			if objInList.DisplayName == objName {
				if found {
					return fmt.Errorf("Found multiple switching profiles with name '%s'", objName)
				}
				obj = objInList
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Switching profile with name '%s' was not found", objName)
		}
	} else {
		return fmt.Errorf("Error obtaining switching profile ID or name during read")
	}

	d.SetId(obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("resource_type", obj.ResourceType)
	d.Set("description", obj.Description)
	return nil
}
