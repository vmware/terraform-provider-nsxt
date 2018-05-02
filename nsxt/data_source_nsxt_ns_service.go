/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
	"strings"
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
			return fmt.Errorf("Error while reading NS service %s: %v", objID, err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("NS service %s was not found", objID)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining NS service ID or name during read")
	} else {
		// Get by full name/prefix
		// TODO use 2nd parameter localVarOptionals for paging
		objList, _, err := nsxClient.GroupingObjectsApi.ListNSServices(nsxClient.Context, nil)
		if err != nil {
			return fmt.Errorf("Error while reading NS services: %v", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []manager.NsService
		var prefixMatch []manager.NsService
		for _, objInList := range objList.Results {
			if strings.HasPrefix(objInList.DisplayName, objName) {
				prefixMatch = append(prefixMatch, objInList)
			}
			if objInList.DisplayName == objName {
				perfectMatch = append(perfectMatch, objInList)
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple NS services with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple NS services with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("NS service '%s' was not found", objName)
		}
	}

	d.SetId(obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)

	return nil
}
