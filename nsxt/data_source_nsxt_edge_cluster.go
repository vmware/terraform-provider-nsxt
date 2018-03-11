/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
	"strings"
)

func dataSourceNsxtEdgeCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNsxtEdgeClusterRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Unique ID of this resource",
				Optional:    true,
				Computed:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource",
				Optional:    true,
				Computed:    true,
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
				Computed:    true,
			},
			"deployment_type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The deployment type of edge cluster members (UNKNOWN/VIRTUAL_MACHINE|PHYSICAL_MACHINE)",
				Optional:    true,
				Computed:    true,
			},
			"member_node_type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Type of transport nodes",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtEdgeClusterRead(d *schema.ResourceData, m interface{}) error {
	// Read an edge cluster by name or id
	nsxClient := m.(*api.APIClient)
	objID := d.Get("id").(string)
	objName := d.Get("display_name").(string)
	var obj manager.EdgeCluster
	if objID != "" {
		// Get by id
		objGet, resp, err := nsxClient.NetworkTransportApi.ReadEdgeCluster(nsxClient.Context, objID)

		if err != nil {
			return fmt.Errorf("Error while reading edge cluster %s: %v", objID, err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Edge cluster %s was not found", objID)
		}
		obj = objGet
	} else if objName != "" {
		// Get by name prefix
		// TODO use 2nd parameter localVarOptionals for paging
		objList, _, err := nsxClient.NetworkTransportApi.ListEdgeClusters(nsxClient.Context, nil)
		if err != nil {
			return fmt.Errorf("Error while reading edge clusters: %v", err)
		}
		// go over the list to find the correct one
		// TODO: prefer full match
		found := false
		for _, objInList := range objList.Results {
			if strings.HasPrefix(objInList.DisplayName, objName) {
				if found {
					return fmt.Errorf("Found multiple edge clusters with name '%s'", objName)
				}
				obj = objInList
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Edge cluster '%s' was not found", objName)
		}
	} else {
		return fmt.Errorf("Error obtaining edge cluster ID or name during read")
	}

	d.SetId(obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("deployment_type", obj.DeploymentType)
	d.Set("member_node_type", obj.MemberNodeType)

	return nil
}
