/* Copyright © 2024 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/fabric"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

func dataSourceNsxtDiscoverNode() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtDiscoverNodeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "External id of the discovered node, ex. a mo-ref from VC",
				Optional:    true,
				Computed:    true,
			},
			"ip_address": {
				Type:         schema.TypeString,
				Description:  "IP Address of the the discovered node.",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateCidrOrIPOrRange(),
			},
		},
	}
}

func dataSourceNsxtDiscoverNodeRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	discoverNodeClient := fabric.NewDiscoveredNodesClient(connector)

	objID := d.Get("id").(string)
	ipAddress := d.Get("ip_address").(string)

	var obj model.DiscoveredNode
	if objID != "" {
		// Get by ID
		objGet, err := discoverNodeClient.Get(objID)
		if isNotFoundError(err) {
			return fmt.Errorf("Discover Node %s was not found", objID)
		}
		if err != nil {
			return fmt.Errorf("Error while reading Discover Node %s: %v", objID, err)
		}
		obj = objGet
	} else if ipAddress == "" {
		return fmt.Errorf("Error obtaining Discover Node external ID or IP address during read")
	} else {
		// Get by IP address
		objList, err := discoverNodeClient.List(nil, nil, nil, nil, nil, nil, &ipAddress, nil, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("Error while reading Discover Node: %v", err)
		}
		if len(objList.Results) == 0 {
			return fmt.Errorf("Discover Node with IP %s was not found", ipAddress)
		}
		obj = objList.Results[0]
	}

	d.SetId(*obj.ExternalId)
	if len(obj.IpAddresses) > 0 {
		d.Set("ip_address", obj.IpAddresses[0])
	}
	return nil
}
