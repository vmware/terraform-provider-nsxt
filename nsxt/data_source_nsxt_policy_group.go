/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/domains"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyGroupRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"domain":       getDomainNameSchema(),
		},
	}
}

func dataSourceNsxtPolicyGroupRead(d *schema.ResourceData, m interface{}) error {
	if isPolicyGlobalManager(m) {
		domain := d.Get("domain").(string)
		query := make(map[string]string)
		query["parent_path"] = "*/" + domain
		_, err := policyDataSourceResourceRead(d, getPolicyConnector(m), "Group", query)
		if err != nil {
			return err
		}
		return nil
	}

	connector := getPolicyConnector(m)
	client := domains.NewDefaultGroupsClient(connector)
	domain := d.Get("domain").(string)
	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	var obj model.Group
	if objID != "" {
		// Get by id
		objGet, err := client.Get(domain, objID)

		if err != nil {
			return handleDataSourceReadError(d, "Group", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining Group ID or name during read")
	} else {
		// Get by full name/prefix
		objList, err := client.List(domain, nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return handleListError("Group", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.Group
		var prefixMatch []model.Group
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
				return fmt.Errorf("Found multiple Groups with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple Groups with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("Group with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	return nil
}
