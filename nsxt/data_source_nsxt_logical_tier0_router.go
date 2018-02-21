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

func dataSourceNsxtLogicalTier0Router() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtLogicalTier0RouterRead,

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
			"edge_cluster_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The ID of the edge cluster connected to this router",
				Optional:    true,
				Computed:    true,
			},
			"high_availability_mode": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The High availability mode of this router",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtLogicalTier0RouterRead(d *schema.ResourceData, m interface{}) error {
	// Read a logical tier0 router by name or id
	nsxClient := m.(*api.APIClient)
	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	var obj manager.LogicalRouter
	if objID != "" {
		// Get by id
		objGet, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouter(nsxClient.Context, objID)

		if err != nil {
			return fmt.Errorf("Error while reading logical tier0 router %s: %v", objID, err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("logical tier0 router %s was not found", objID)
		}
		obj = objGet
	} else if objName != "" {
		// Get by name prefix
		// TODO use 2nd parameter localVarOptionals for paging
		objList, _, err := nsxClient.LogicalRoutingAndServicesApi.ListLogicalRouters(nsxClient.Context, nil)
		if err != nil {
			return fmt.Errorf("Error while reading logical tier0 routers: %v", err)
		}
		// go over the list to find the correct one
		// TODO: prefer full match
		found := false
		for _, objInList := range objList.Results {
			if strings.HasPrefix(objInList.DisplayName, objName) {
				if found {
					return fmt.Errorf("Found multiple logical tier0 routers with name '%s'", objName)
				}
				obj = objInList
				found = true
			}
		}
		if !found {
			return fmt.Errorf("logical tier0 router '%s' was not found", objName)
		}
	} else {
		return fmt.Errorf("Error obtaining logical tier0 router ID or name during read")
	}

	// TODO(asarfaty): Make sure this is a TIER0 router

	d.SetId(obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("edge_cluster_id", obj.EdgeClusterId)
	d.Set("high_availability_mode", obj.HighAvailabilityMode)

	return nil
}
