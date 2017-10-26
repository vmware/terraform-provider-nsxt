package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func dataSourceTransportZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTransportZoneRead,

		Schema: map[string]*schema.Schema{
			"Id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"DisplayName": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"Description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"HostSwitchName": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"TransportType": &schema.Schema{
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
	obj_id := d.Get("Id").(string)
	obj_name := d.Get("DisplayName").(string)
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
		// Get by name
		// TODO use 2nd parameter localVarOptionals for paging
		obj_list, _, err := nsxClient.NetworkTransportApi.ListTransportZones(nsxClient.Context, nil)
		if err != nil {
			return fmt.Errorf("Error while reading transport zones: %v\n", err)
		}
		// go over the list to find the correct one
		found := false
		for _, obj_in_list := range obj_list.Results {
			if obj_in_list.DisplayName == obj_name {
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
	d.Set("DisplayName", obj.DisplayName)
	d.Set("Description", obj.Description)
	d.Set("HostSwitchName", obj.HostSwitchName)
	d.Set("TransportType", obj.TransportType)

	return nil
}
