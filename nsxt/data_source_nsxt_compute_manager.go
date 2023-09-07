/* Copyright © 2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/fabric"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

func dataSourceNsxtComputeManager() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtComputeManagerRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"server": {
				Type:        schema.TypeString,
				Description: "IP address or hostname of compute manager",
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtComputeManagerRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := fabric.NewComputeManagersClient(connector)

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)

	var obj model.ComputeManager

	if objID != "" {
		// Get by id
		objGet, err := client.Get(objID)
		if err != nil {
			return fmt.Errorf("failed to read ComputeManager %s: %v", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("error obtaining Compute Manager ID or name during read")
	} else {
		// Get by full name/prefix
		objList, err := client.List(nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to read Compute Managers: %v", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.ComputeManager
		var prefixMatch []model.ComputeManager
		for _, objInList := range objList.Results {
			if strings.HasPrefix(*objInList.DisplayName, objName) {
				prefixMatch = append(prefixMatch, objInList)
			}
			if *objInList.DisplayName == objName {
				perfectMatch = append(perfectMatch, objInList)
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("found multiple Compute Managers with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("found multiple Compute Managers with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("Compute Manager with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("server", obj.Server)
	return nil
}
