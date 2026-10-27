// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func dataSourceNsxtVpcLBPersistenceProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcLBPersistenceProfileRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"type": {
				Type:         schema.TypeString,
				Description:  "Load Balancer Persistence Type",
				Optional:     true,
				Default:      "ANY",
				ValidateFunc: validation.StringInSlice(lbPersistenceTypeValues, false),
			},
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
		},
	}
}

func dataSourceNsxtVpcLBPersistenceProfileRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Persistence Profile data source requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbPersistenceProfilesClient(sessionContext, connector)
	converter := bindings.NewTypeConverter()

	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	objTypeValue, typeSet := d.GetOk("type")
	objType := objTypeValue.(string)

	var obj model.LBPersistenceProfile
	if objID != "" {
		// Get by id
		objGet, err := c.Get(parents[0], parents[1], parents[2], objID)
		if err != nil {
			return handleDataSourceReadError(d, "VPC LBPersistenceProfile", objID, err)
		}
		profile, errs := converter.ConvertToGolang(objGet, model.LBPersistenceProfileBindingType())
		if errs != nil {
			return errs[0]
		}
		obj = profile.(model.LBPersistenceProfile)
		if !dataSourceNsxtPolicyLbPersistenceProfileTypeMatches(obj, objType) {
			return fmt.Errorf("VPC LBPersistenceProfile with ID '%s' and type '%s' was not found", objID, objType)
		}
	} else if objName == "" && !typeSet {
		return fmt.Errorf("Error obtaining VPC LBPersistenceProfile name or type during read")
	} else {
		// Get by full name/prefix
		includeMarkForDeleteObjectsParam := false
		var perfectMatch, prefixMatch []model.LBPersistenceProfile
		var cursor *string
		for {
			objList, err := c.List(parents[0], parents[1], parents[2], cursor, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
			if err != nil {
				return handleListError("VPC LBPersistenceProfile", err)
			}
			// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
			for _, objInList := range objList.Results {
				profile, errs := converter.ConvertToGolang(objInList, model.LBPersistenceProfileBindingType())
				if errs != nil {
					return errs[0]
				}
				lbProfile := profile.(model.LBPersistenceProfile)
				if !dataSourceNsxtPolicyLbPersistenceProfileTypeMatches(lbProfile, objType) {
					continue
				}

				if objName != "" && lbProfile.DisplayName != nil && strings.HasPrefix(*lbProfile.DisplayName, objName) && dataSourceNsxtPolicyLbPersistenceProfileTypeMatches(lbProfile, objType) {
					prefixMatch = append(prefixMatch, lbProfile)
				}
				if objName != "" && lbProfile.DisplayName != nil && *lbProfile.DisplayName == objName && dataSourceNsxtPolicyLbPersistenceProfileTypeMatches(lbProfile, objType) {
					perfectMatch = append(perfectMatch, lbProfile)
				}
				if objName == "" && typeSet && dataSourceNsxtPolicyLbPersistenceProfileTypeMatches(lbProfile, objType) {
					// match only by type
					perfectMatch = append(perfectMatch, lbProfile)
				}
			}
			cursor = objList.Cursor
			if cursor == nil || *cursor == "" {
				break
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple VPC LBPersistenceProfiles with name '%s' and type '%s'", objName, objType)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple VPC LBPersistenceProfiles with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("VPC LBPersistenceProfile with name '%s' and type '%s' was not found", objName, objType)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("type", lbPersistenceTypeMap[obj.ResourceType])
	d.Set("path", obj.Path)
	return nil
}
