/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ipsec_vpn_services "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/locale_services/ipsec_vpn_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dataSourceNsxtPolicyIPSecVpnLocalEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtPolicyIPSecVpnLocalEndpointRead,

		Schema: map[string]*schema.Schema{
			"id":           getDataSourceIDSchema(),
			"display_name": getDataSourceDisplayNameSchema(),
			"description":  getDataSourceDescriptionSchema(),
			"path":         getPathSchema(),
			"tier0_id": {
				Type:        schema.TypeString,
				Description: "Policy path referencing Local endpoint.",
				Optional:    true,
				Default:     "vmc",
			},
			"locale_service": {
				Type:        schema.TypeString,
				Description: "Local_service",
				Optional:    true,
				Default:     "default",
			},
			"service_id": {
				Type:        schema.TypeString,
				Description: "Policy path referencing Local endpoint.",
				Optional:    true,
				Default:     "default",
			},
			"local_address": {
				Type:        schema.TypeString,
				Description: "Local IP Address",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyIPSecVpnLocalEndpointRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	Tier0ID := d.Get("tier0_id").(string)
	LocaleService := d.Get("locale_service").(string)
	ServiceID := d.Get("service_id").(string)

	objID := d.Get("id").(string)
	log.Println(objID)

	objName := d.Get("display_name").(string)
	client := ipsec_vpn_services.NewDefaultLocalEndpointsClient(connector)
	var obj model.IPSecVpnLocalEndpoint
	if objID != "" {
		// Get by id
		objGet, err := client.Get(Tier0ID, LocaleService, ServiceID, objID)
		if isNotFoundError(err) {
			return fmt.Errorf("IPSecVpnLocalEndpoint with ID %s was not found", objID)
		}

		if err != nil {
			return fmt.Errorf("Error while reading IPSecVpnLocalEndpoint %s: %v", objID, err)
		}
		obj = objGet
	} else if objName == "" {
		return fmt.Errorf("Error obtaining IPSecVpnLocalEndpoint ID or name during read")
	} else {
		// Get by full name/prefix
		includeMarkForDeleteObjectsParam := false
		objList, err := client.List(Tier0ID, LocaleService, ServiceID, nil, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("Error while reading <!RESOURCES!>: %v", err)
		}
		// go over the list to find the correct one (prefer a perfect match. If not - prefix match)
		var perfectMatch []model.IPSecVpnLocalEndpoint
		var prefixMatch []model.IPSecVpnLocalEndpoint
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
				return fmt.Errorf("Found multiple <!RESOURCES!> with name '%s'", objName)
			}
			obj = perfectMatch[0]
		} else if len(prefixMatch) > 0 {
			if len(prefixMatch) > 1 {
				return fmt.Errorf("Found multiple <!RESOURCES!> with name starting with '%s'", objName)
			}
			obj = prefixMatch[0]
		} else {
			return fmt.Errorf("IPSecVpnLocalEndpoint with name '%s' was not found", objName)
		}
	}

	d.SetId(*obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	d.Set("local_address", obj.LocalAddress)
	return nil
}
