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

func dataSourceNsxtTransportZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtTransportZoneRead,

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
			"host_switch_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Name of the host switch on all transport nodes in this transport zone that will be used to run NSX network traffic",
				Optional:    true,
				Computed:    true,
			},
			"transport_type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The transport type of this transport zone (OVERLAY or VLAN)",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtTransportZoneRead(d *schema.ResourceData, m interface{}) error {
	// Read a transport zone by name or id
	nsxClient := m.(*api.APIClient)
	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	var obj manager.TransportZone
	if objID != "" {
		// Get by id
		objGet, resp, err := nsxClient.NetworkTransportApi.GetTransportZone(nsxClient.Context, objID)

		if err != nil {
			return fmt.Errorf("Error while reading transport zone %s: %v\n", objID, err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Transport zone %s was not found\n", objID)
		}
		obj = objGet
	} else if objName != "" {
		// Get by name prefix
		// TODO use 2nd parameter localVarOptionals for paging
		objList, _, err := nsxClient.NetworkTransportApi.ListTransportZones(nsxClient.Context, nil)
		if err != nil {
			return fmt.Errorf("Error while reading transport zones: %v\n", err)
		}
		// go over the list to find the correct one
		// TODO: prefer full match
		found := false
		for _, objInList := range objList.Results {
			if strings.HasPrefix(objInList.DisplayName, objName) {
				if found == true {
					return fmt.Errorf("Found multiple transport zones with name '%s'\n", objName)
				}
				obj = objInList
				found = true
			}
		}
		if found == false {
			return fmt.Errorf("Transport zone '%s' was not found\n", objName)
		}
	} else {
		return fmt.Errorf("Error obtaining transport zone ID or name during read")
	}

	d.SetId(obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("host_switch_name", obj.HostSwitchName)
	d.Set("transport_type", obj.TransportType)

	return nil
}
