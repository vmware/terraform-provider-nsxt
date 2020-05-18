/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"strings"
)

var policyTransportZoneTransportTypes = [](string){
	model.PolicyTransportZone_TZ_TYPE_OVERLAY_STANDARD,
	model.PolicyTransportZone_TZ_TYPE_OVERLAY_ENS,
	model.PolicyTransportZone_TZ_TYPE_VLAN_BACKED,
	model.PolicyTransportZone_TZ_TYPE_UNKNOWN,
}

func dataSourceNsxtPolicyTransportZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyTransportZoneRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"is_default": {
				Type:        schema.TypeBool,
				Description: "Indicates whether the transport zone is default",
				Optional:    true,
				Computed:    true,
			},
			"transport_type": {
				Type:         schema.TypeString,
				Description:  "Type of Transport Zone",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(policyTransportZoneTransportTypes, false),
			},
		},
	}
}

func dataSourceNsxtPolicyTransportZoneRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := enforcement_points.NewDefaultTransportZonesClient(connector)

	// TODO: support non-default site and enforcement point possibly as a triple; site/point/tz_id
	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	defaultVal, isDefaultSet := d.GetOkExists("is_default")
	isDefault := isDefaultSet && defaultVal.(bool)
	transportType := d.Get("transport_type").(string)
	var obj model.PolicyTransportZone
	if objID != "" {
		// Get by id
		objGet, err := client.Get(defaultSite, getPolicyEnforcementPoint(m), objID)

		if err != nil {
			return handleDataSourceReadError(d, "TransportZone", objID, err)
		}
		obj = objGet
	} else if objName == "" && !(isDefault && transportType != "") {
		return fmt.Errorf("Please specify id, display_name or is_default and transport_type in order to identify Transport Zone")
	} else {
		// Get by full name/prefix
		includeMarkForDeleteObjectsParam := false
		objList, err := client.List(defaultSite, getPolicyEnforcementPoint(m), nil, &includeMarkForDeleteObjectsParam, nil, nil, &includeMarkForDeleteObjectsParam, nil)
		if err != nil {
			return handleListError("TransportZone", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.PolicyTransportZone
		var prefixMatch []model.PolicyTransportZone
		for _, objInList := range objList.Results {
			if transportType != "" && transportType != *objInList.TzType {
				// no match for transport type
				continue
			}

			if isDefault && *objInList.IsDefault {
				// user is looking for default TZ
				perfectMatch = append(perfectMatch, objInList)
				break
			}

			if strings.HasPrefix(*objInList.DisplayName, objName) {
				prefixMatch = append(prefixMatch, objInList)
			}
			if *objInList.DisplayName == objName {
				perfectMatch = append(perfectMatch, objInList)
			}
		}
		if len(perfectMatch) > 0 {
			if len(perfectMatch) > 1 {
				return fmt.Errorf("Found multiple TransportZones with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple TransportZones with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("TransportZone '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	d.Set("is_default", obj.IsDefault)
	d.Set("transport_type", obj.TzType)
	return nil
}
