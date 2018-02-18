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
				Description: "This field could show deployment_type of members. It would return UNKNOWN if there is no members, and return VIRTUAL_MACHINE|PHYSICAL_MACHINE if all edge members are VIRTUAL_MACHINE|PHYSICAL_MACHINE.",
				Optional:    true,
				Computed:    true,
			},
			"member_node_type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Edge cluster is homogenous collection of transport nodes. Hence all transport nodes of the cluster must be of same type. This field shows the type of transport nodes.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtEdgeClusterRead(d *schema.ResourceData, m interface{}) error {
	// Read an edge cluster by name or id
	nsxClient := m.(*api.APIClient)
	obj_id := d.Get("id").(string)
	obj_name := d.Get("display_name").(string)
	var obj manager.EdgeCluster
	if obj_id != "" {
		// Get by id
		obj_get, resp, err := nsxClient.NetworkTransportApi.ReadEdgeCluster(nsxClient.Context, obj_id)

		if err != nil {
			return fmt.Errorf("Error while reading edge cluster %s: %v\n", obj_id, err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Edge cluster %s was not found\n", obj_id)
		}
		obj = obj_get
	} else if obj_name != "" {
		// Get by name prefix
		// TODO use 2nd parameter localVarOptionals for paging
		obj_list, _, err := nsxClient.NetworkTransportApi.ListEdgeClusters(nsxClient.Context, nil)
		if err != nil {
			return fmt.Errorf("Error while reading edge clusters: %v\n", err)
		}
		// go over the list to find the correct one
		// TODO: prefer full match
		found := false
		for _, obj_in_list := range obj_list.Results {
			if strings.HasPrefix(obj_in_list.DisplayName, obj_name) {
				if found == true {
					return fmt.Errorf("Found multiple edge clusters with name '%s'\n", obj_name)
				}
				obj = obj_in_list
				found = true
			}
		}
		if found == false {
			return fmt.Errorf("Edge cluster '%s' was not found\n", obj_name)
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
