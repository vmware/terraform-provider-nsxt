package main

import (
	"fmt"
	"strings"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func dataSourceTransportZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTransportZoneRead,

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
			"host_switch_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"transport_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceTransportZoneRead(d *schema.ResourceData, m interface{}) error {
	// Read a transport zone by name or id
	nsxClient := m.(*nsxt.APIClient)
	obj_id := d.Get("id").(string)
	obj_name := d.Get("display_name").(string)
	var obj manager.TransportZone
	if obj_id != "" {
		// Get by id
		obj_get, resp, err := nsxClient.NetworkTransportApi.GetTransportZone(nsxClient.Context, obj_id)

		if err != nil {
			return fmt.Errorf("Error while reading transport zone %s: %v\n", obj_id, err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Transport zone %s was not found\n", obj_id)
		}
		obj = obj_get
	} else if obj_name != "" {
		// Get by name prefix
		// TODO use 2nd parameter localVarOptionals for paging
		obj_list, _, err := nsxClient.NetworkTransportApi.ListTransportZones(nsxClient.Context, nil)
		if err != nil {
			return fmt.Errorf("Error while reading transport zones: %v\n", err)
		}
		// go over the list to find the correct one
		// TODO: prefer full match
		found := false
		for _, obj_in_list := range obj_list.Results {
			if strings.HasPrefix(obj_in_list.DisplayName, obj_name) {
				if found == true {
					return fmt.Errorf("Found multiple transport zones with name '%s'\n", obj_name)
				}
				obj = obj_in_list
				found = true
			}
		}
		if found == false {
			return fmt.Errorf("Transport zone '%s' was not found\n", obj_name)
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
