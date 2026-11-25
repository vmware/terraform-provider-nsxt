// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func dataSourceNsxtVpcLBMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtVpcLBMonitorRead,

		Schema: map[string]*schema.Schema{
			"id": getDataSourceIDSchema(),
			"type": {
				Type:         schema.TypeString,
				Description:  "Load Balancer Monitor Type",
				Optional:     true,
				Default:      "ANY",
				ValidateFunc: validation.StringInSlice(lbMonitorTypeValues, false),
			},
			"display_name": getDataSourceExtendedDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"context":      getContextSchemaExtended(true, false, true, true),
		},
	}
}

func dataSourceNsxtVpcLBMonitorRead(d *schema.ResourceData, m interface{}) error {
	if !util.NsxVersionHigherOrEqual("9.0.0") {
		return fmt.Errorf("VPC LB Monitor data source requires NSX version 9.0.0 or higher")
	}
	connector := getPolicyConnector(m)
	sessionContext := getSessionContext(d, m)
	parents := getVpcParentsFromContext(sessionContext)
	c := cliVpcLbMonitorProfilesClient(sessionContext, connector)

	objID := d.Get("id").(string)
	objTypeValue, typeSet := d.GetOk("type")
	objType := objTypeValue.(string)
	objName := d.Get("display_name").(string)
	var result *model.LBMonitorProfile
	if objID != "" {
		// Get by id
		objGet, err := c.Get(parents[0], parents[1], parents[2], objID)

		if err != nil {
			return handleDataSourceReadError(d, "VPC LBMonitor", objID, err)
		}
		result, err = policyLbMonitorConvert(objGet, objType)
		if err != nil {
			return fmt.Errorf("Error while converting VPC LBMonitor %s: %v", objID, err)
		}
		if result == nil {
			return fmt.Errorf("VPC LBMonitor with ID '%s' and type %s was not found", objID, objType)
		}
	} else if objName == "" && !typeSet {
		return fmt.Errorf("Error obtaining VPC LBMonitor ID or name or type during read")
	} else {
		// Get by full name/prefix
		includeMarkForDeleteObjectsParam := false
		var perfectMatch []model.LBMonitorProfile
		var prefixMatch []model.LBMonitorProfile
		var cursor *string
		for {
			objList, err := c.List(parents[0], parents[1], parents[2], cursor, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
			if err != nil {
				return handleListError("VPC LBMonitor", err)
			}
			// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
			for _, objInList := range objList.Results {
				obj, err := policyLbMonitorConvert(objInList, objType)
				if err != nil {
					return fmt.Errorf("Error while converting VPC LBMonitor for name %q and type %s: %w", objName, objType, err)
				}
				if obj == nil {
					continue
				}
				if objName != "" && obj.DisplayName != nil && strings.HasPrefix(*obj.DisplayName, objName) {
					prefixMatch = append(prefixMatch, *obj)
				}
				if obj.DisplayName != nil && *obj.DisplayName == objName {
					perfectMatch = append(perfectMatch, *obj)
				}
				if objName == "" && typeSet {
					// match only by type
					perfectMatch = append(perfectMatch, *obj)
				}
			}
			cursor = objList.Cursor
			if cursor == nil || *cursor == "" {
				break
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple VPC LBMonitors with name '%s'", objName)
			}
			result = &perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple VPC LBMonitors with name starting with '%s'", objName)
			}
			result = &prefixMatch[0]
		} else {
			return fmt.Errorf("VPC LBMonitor with name '%s' and type %s was not found", objName, objType)
		}
	}

	d.SetId(*result.Id)
	d.Set("display_name", result.DisplayName)
	d.Set("type", lbMonitorTypeMap[result.ResourceType])
	d.Set("description", result.Description)
	d.Set("path", result.Path)
	return nil
}
