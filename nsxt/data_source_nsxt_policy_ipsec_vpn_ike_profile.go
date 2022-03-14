/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyIPSecVpnIkeProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIPSecVpnIkeProfileRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
		},
	}
}

func dataSourceNsxtPolicyIPSecVpnIkeProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	if isPolicyGlobalManager(m) {
		_, err := policyDataSourceResourceRead(d, connector, false, "IPSecVpnIkeProfile", nil)
		if err != nil {
			return err
		}

		return nil
	}

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	log.Println("########################################################")
	log.Println(objName)
	log.Println("########################################################")
	client := infra.NewDefaultIpsecVpnIkeProfilesClient(connector)

	var obj model.IPSecVpnIkeProfile
	if objID != "" {
		// Get by id
		objGet, err := client.Get(objID)
		if isNotFoundError(err) {
			return fmt.Errorf("IPSecVpnIkeProfile with ID %s was not found", objID)
		}

		if err != nil {
			return fmt.Errorf("Error while reading IPSecVpnIkeProfile %s: %v", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining IPSecVpnIkeProfile ID or name during read")
	} else {
		// Get by full name/prefix
		includeMarkForDeleteObjectsParam := false
		objList, err := client.List(nil, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("Error while reading IPSecVpnIkeProfiles: %v", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.IPSecVpnIkeProfile
		var prefixMatch []model.IPSecVpnIkeProfile
		for _, objInList := range objList.Results {
			if strings.HasPrefix(*objInList.DisplayName, objName) {
				prefixMatch = append(prefixMatch, objInList)
				log.Println("HAS PREFIX?########################################################")
			}
			if *objInList.DisplayName == objName {
				perfectMatch = append(perfectMatch, objInList)
				log.Println("PERFECT MATCH?########################################################")
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple IPSecVpnIkeProfiles with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple IPSecVpnIkeProfiles with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("IPSecVpnIkeProfile with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	return nil
}
