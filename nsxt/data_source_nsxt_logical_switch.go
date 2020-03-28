/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/go-vmware-nsxt/manager"
)

func dataSourceNsxtLogicalSwitch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtLogicalSwitchRead,

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
			"description": {
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtLogicalSwitchRead(d *schema.ResourceData, m interface{}) error {
	// Read a logical tier1 router by name or id
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return dataSourceNotSupportedError()
	}

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	var obj manager.LogicalSwitch
	if objID != "" {
		// Get by id
		objGet, resp, err := nsxClient.LogicalSwitchingApi.GetLogicalSwitch(nsxClient.Context, objID)

		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Logical tier1 router %s was not found", objID)
		}
		if err != nil {
			return fmt.Errorf("Error while reading logical switch %s: %v", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining logical switch ID or name during read")
	} else {
		// Get by full name/prefix
		objList, _, err := nsxClient.LogicalSwitchingApi.ListLogicalSwitches(nsxClient.Context, nil)
		if err != nil {
			return fmt.Errorf("Error while reading logical switches: %v", err)
		}

		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []manager.LogicalSwitch
		var prefixMatch []manager.LogicalSwitch
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
				return fmt.Errorf("Found multiple logical switches with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple logical switches with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("Logical switch with name '%s' was not found", objName)
		}
	}

	d.SetId(obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("transport_zone_id", obj.TransportZoneId)

	return nil
}
