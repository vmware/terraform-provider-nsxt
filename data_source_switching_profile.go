package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func dataSourceSwitchingProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSwitchingProfileRead,

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
			"resource_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSwitchingProfileRead(d *schema.ResourceData, m interface{}) error {
	// Read a switching profile by name or id
	nsxClient := m.(*nsxt.APIClient)
	obj_id := d.Get("Id").(string)
	obj_name := d.Get("DisplayName").(string)
	var obj manager.BaseSwitchingProfile
	if obj_id != "" {
		// Get by id
		obj_get, resp, err := nsxClient.LogicalSwitchingApi.GetSwitchingProfile(nsxClient.Context, obj_id)

		if err != nil {
			return fmt.Errorf("Error while reading switching profile %s: %v\n", obj_id, err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("switching profile %s was not found\n", obj_id)
		}
		obj = obj_get
	} else if obj_name != "" {
		// Get by name
		// TODO use 2nd parameter localVarOptionals for paging
		localVarOptionals := make(map[string]interface{})
		localVarOptionals["includeSystemOwned"] = true
		obj_list, _, err := nsxClient.LogicalSwitchingApi.ListSwitchingProfiles(nsxClient.Context, localVarOptionals)
		if err != nil {
			return fmt.Errorf("Error while reading switching profiles: %v\n", err)
		}
		// go over the list to find the correct one
		found := false
		for _, obj_in_list := range obj_list.Results {
			if obj_in_list.DisplayName == obj_name {
				if found == true {
					return fmt.Errorf("Found multiple switching profiles with name '%s'\n", obj_name)
				}
				obj = obj_in_list
				found = true
			}
		}
		if found == false {
			return fmt.Errorf("Switching profile '%s' was not found\n", obj_name)
		}
	} else {
		return fmt.Errorf("Error obtaining switching profile ID or name during read")
	}

	d.SetId(obj.Id)
	d.Set("DisplayName", obj.DisplayName)
	d.Set("resource_type", obj.ResourceType)
	return nil
}
