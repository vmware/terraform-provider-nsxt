/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func dataSourceNsxtNsService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtNsServiceRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Unique ID of this resource",
				Optional:    true,
				Computed:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource",
				Optional:    true,
				Computed:    true,
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtNsServiceRead(d *schema.ResourceData, m interface{}) error {
	// Read NS Service by name or id
	nsxClient := m.(*api.APIClient)
	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	var obj manager.NsService
	if objID != "" {
		// Get by id
		objGet, resp, err := nsxClient.GroupingObjectsApi.ReadNSService(nsxClient.Context, objID)

		if err != nil {
			return fmt.Errorf("Error while reading ns service %s: %v", objID, err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("ns service %s was not found", objID)
		}
		obj = objGet
	} else if objName != "" {
		// Get by name
		// TODO use 2nd parameter localVarOptionals for paging
		objList, _, err := nsxClient.GroupingObjectsApi.ListNSServices(nsxClient.Context, nil)
		if err != nil {
			return fmt.Errorf("Error while reading ns services: %v", err)
		}
		// go over the list to find the correct one
		found := false
		for _, objInList := range objList.Results {
			if objInList.DisplayName == objName {
				if found {
					return fmt.Errorf("Found multiple ns services with name '%s'", objName)
				}
				obj = objInList
				found = true
			}
		}
		if !found {
			return fmt.Errorf("ns service '%s' was not found out of %d services", objName, len(objList.Results))
		}
	} else {
		return fmt.Errorf("Error obtaining ns service ID or name during read")
	}

	d.SetId(obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)

	return nil
}
