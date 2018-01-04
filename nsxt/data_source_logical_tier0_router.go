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

func dataSourceLogicalTier0Router() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLogicalTier0RouterRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"router_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"edge_cluster_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"high_availability_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceLogicalTier0RouterRead(d *schema.ResourceData, m interface{}) error {
	// Read a logical tier0 router by name or id
	nsxClient := m.(*api.APIClient)
	obj_id := d.Get("id").(string)
	obj_name := d.Get("display_name").(string)
	var obj manager.LogicalRouter
	if obj_id != "" {
		// Get by id
		obj_get, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouter(nsxClient.Context, obj_id)

		if err != nil {
			return fmt.Errorf("Error while reading logical tier0 router %s: %v\n", obj_id, err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("logical tier0 router %s was not found\n", obj_id)
		}
		obj = obj_get
	} else if obj_name != "" {
		// Get by name prefix
		// TODO use 2nd parameter localVarOptionals for paging
		obj_list, _, err := nsxClient.LogicalRoutingAndServicesApi.ListLogicalRouters(nsxClient.Context, nil)
		if err != nil {
			return fmt.Errorf("Error while reading logical tier0 routers: %v\n", err)
		}
		// go over the list to find the correct one
		// TODO: prefer full match
		found := false
		for _, obj_in_list := range obj_list.Results {
			if strings.HasPrefix(obj_in_list.DisplayName, obj_name) {
				if found == true {
					return fmt.Errorf("Found multiple logical tier0 routers with name '%s'\n", obj_name)
				}
				obj = obj_in_list
				found = true
			}
		}
		if found == false {
			return fmt.Errorf("logical tier0 router '%s' was not found\n", obj_name)
		}
	} else {
		return fmt.Errorf("Error obtaining logical tier0 router ID or name during read")
	}

	d.SetId(obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("router_type", obj.RouterType)
	d.Set("edge_cluster_id", obj.EdgeClusterId)
	d.Set("high_availability_mode", obj.HighAvailabilityMode)

	return nil
}
