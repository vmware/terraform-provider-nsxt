// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func dataSourceNsxtVpcLBServerSslProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcLBServerSslProfileRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
		},
	}
}

func dataSourceNsxtVpcLBServerSslProfileRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Server SSL Profile data source requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbServerSslProfilesClient(sessionContext, connector)

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	var obj model.LBServerSslProfile
	if objID != "" {
		// Get by id
		objGet, err := c.Get(parents[0], parents[1], parents[2], objID)
		if err != nil {
			return handleDataSourceReadError(d, "VPC LBServerSslProfile", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining VPC LBServerSslProfile ID or name during read")
	} else {
		// Get by full name/prefix
		includeMarkForDeleteObjectsParam := false
		objList, err := c.List(parents[0], parents[1], parents[2], nil, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
		if err != nil {
			return handleListError("VPC LBServerSslProfile", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.LBServerSslProfile
		var prefixMatch []model.LBServerSslProfile
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
				return fmt.Errorf("Found multiple VPC LBServerSslProfiles with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple VPC LBServerSslProfiles with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("VPC LBServerSslProfile with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	return nil
}
