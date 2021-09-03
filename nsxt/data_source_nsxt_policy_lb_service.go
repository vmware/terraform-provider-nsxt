/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyLbService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyLbServiceRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
		},
	}
}

func dataSourceNsxtPolicyLbServiceRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultLbServicesClient(connector)
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)

	var obj model.LBService
	if objID != "" {
		// Get by id
		objGet, err := client.Get(objID)

		if err != nil {
			return handleDataSourceReadError(d, "LbService", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining LbService name or type during read")
	} else {
		// Get by full name/prefix
		includeMarkForDeleteObjectsParam := false
		objList, err := client.List(nil, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
		if err != nil {
			return handleListError("LbService", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.LBService
		var prefixMatch []model.LBService
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
				return fmt.Errorf("Found multiple PolicyLbServices with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple PolicyLbServices with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("PolicyLbService with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	return nil
}
