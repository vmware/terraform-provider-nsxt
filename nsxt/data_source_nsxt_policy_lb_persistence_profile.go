/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"strings"
)

var lbPersistenceTypeValues = []string{"SOURCE_IP", "COOKIE", "GENERIC", "ANY"}
var lbPersistenceTypeMap = map[string]string{
	model.LBPersistenceProfile_RESOURCE_TYPE_LBSOURCEIPPERSISTENCEPROFILE: "SOURCE_IP",
	model.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE:   "COOKIE",
	model.LBPersistenceProfile_RESOURCE_TYPE_LBGENERICPERSISTENCEPROFILE:  "GENERIC",
}

func dataSourceNsxtPolicyLbPersistenceProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyLbPersistenceProfileRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"type": {
				Type:         schema.TypeString,
				Description:  "Load Balancer Persistence Type",
				Optional:     true,
				Default:      "ANY",
				ValidateFunc: validation.StringInSlice(lbPersistenceTypeValues, false),
			},

			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
		},
	}
}

func dataSourceNsxtPolicyLbPersistenceProfileTypeMatches(profile model.PolicyLbPersistenceProfile, profileType string) bool {
	if profileType == "ANY" {
		return true
	}
	if lbPersistenceTypeMap[profile.ResourceType] == profileType {
		return true
	}
	return false
}

func dataSourceNsxtPolicyLbPersistenceProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultLbPersistenceProfilesClient(connector)
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	objTypeValue, typeSet := d.GetOkExists("type")
	objType := objTypeValue.(string)

	var obj model.PolicyLbPersistenceProfile
	if objID != "" {
		// Get by id
		objGet, err := client.Get(objID)

		if err != nil {
			return handleDataSourceReadError(d, "LbPersistenceProfile", objID, err)
		}
		profile, errs := converter.ConvertToGolang(objGet, model.PolicyLbPersistenceProfileBindingType())
		if errs != nil {
			return errs[0]
		}
		obj = profile.(model.PolicyLbPersistenceProfile)
	} else if objName == "" && !typeSet {
		return fmt.Errorf("Error obtaining LbPersistenceProfile name or type during read")
	} else {
		// Get by full name/prefix
		includeMarkForDeleteObjectsParam := false
		objList, err := client.List(nil, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
		if err != nil {
			return handleListError("LbPersistenceProfile", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.PolicyLbPersistenceProfile
		var prefixMatch []model.PolicyLbPersistenceProfile
		for _, objInList := range objList.Results {
			profile, errs := converter.ConvertToGolang(objInList, model.PolicyLbPersistenceProfileBindingType())
			if errs != nil {
				return errs[0]
			}
			lbProfile := profile.(model.PolicyLbPersistenceProfile)

			if objName != "" && strings.HasPrefix(*lbProfile.DisplayName, objName) && dataSourceNsxtPolicyLbPersistenceProfileTypeMatches(lbProfile, objType) {
				prefixMatch = append(prefixMatch, lbProfile)
			}
			if *lbProfile.DisplayName == objName {
				perfectMatch = append(perfectMatch, lbProfile)
			}
			if objName == "" && typeSet && dataSourceNsxtPolicyLbPersistenceProfileTypeMatches(lbProfile, objType) {
				// match only by type
				perfectMatch = append(perfectMatch, lbProfile)
			}

		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple PolicyLbPersistenceProfiles with name '%s' and type '%s'", objName, objType)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple PolicyLbPersistenceProfiles with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("PolicyLbPersistenceProfile with name '%s' and type '%s' was not found", objName, objType)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("type", lbPersistenceTypeMap[obj.ResourceType])
	d.Set("path", obj.Path)
	return nil
}
