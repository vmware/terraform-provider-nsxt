/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"strings"
)

func dataSourceNsxtPolicyService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyServiceRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
		},
	}
}

func dataSourceNsxtPolicyServiceReadAllServices(connector *client.RestConnector) ([]model.Service, error) {
	var results []model.Service
	client := infra.NewDefaultServicesClient(connector)
	boolFalse := false
	var cursor *string
	total := 0

	for {
		services, err := client.List(cursor, nil, &boolFalse, nil, nil, &boolFalse, nil)
		if err != nil {
			return results, err
		}
		results = append(results, services.Results...)
		if total == 0 && services.ResultCount != nil {
			// first response
			total = int(*services.ResultCount)
		}
		cursor = services.Cursor
		if len(results) >= total {
			return results, nil
		}
	}
}

func dataSourceNsxtPolicyServiceRead(d *schema.ResourceData, m interface{}) error {
	// Read a service by name or id
	connector := getPolicyConnector(m)
	client := infra.NewDefaultServicesClient(connector)

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	var obj model.Service
	if objID != "" {
		// Get by id
		objGet, err := client.Get(objID)

		if err != nil {
			return handleDataSourceReadError(d, "Service", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining service ID or name during read")
	} else {
		// Get by full name/prefix
		objList, err := dataSourceNsxtPolicyServiceReadAllServices(connector)
		if err != nil {
			return handleListError("Service", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.Service
		var prefixMatch []model.Service
		for _, objInList := range objList {
			if strings.HasPrefix(*objInList.DisplayName, objName) {
				prefixMatch = append(prefixMatch, objInList)
			}
			if *objInList.DisplayName == objName {
				perfectMatch = append(perfectMatch, objInList)
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple services with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple services with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("Service '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	return nil
}
