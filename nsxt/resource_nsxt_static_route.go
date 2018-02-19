/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/common"
	"github.com/vmware/go-vmware-nsxt/manager"
	"log"
	"net/http"
)

func resourceNsxtStaticRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtStaticRouteCreate,
		Read:   resourceNsxtStaticRouteRead,
		Update: resourceNsxtStaticRouteUpdate,
		Delete: resourceNsxtStaticRouteDelete,

		Schema: map[string]*schema.Schema{
			"logical_router_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Logical router id",
				Required:    true,
			},
			"network": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "CIDR",
				Required:     true,
				ValidateFunc: validation.CIDRNetwork(0, 32),
			},
			"next_hop": getNextHopsSchema(),
			"revision": getRevisionSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"tag": getTagsSchema(),
		},
	}
}

func getNextHopsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"administrative_distance": &schema.Schema{
					Type:        schema.TypeInt,
					Description: "Administrative Distance for the next hop IP",
					Optional:    true,
				},
				"bfd_enabled": &schema.Schema{
					Type:        schema.TypeBool,
					Description: "Status of bfd for this next hop where bfd_enabled = true indicate bfd is enabled for this next hop and bfd_enabled = false indicate bfd peer is disabled or not configured for this next hop.",
					Computed:    true,
				},
				"blackhole_action": &schema.Schema{
					Type:        schema.TypeString,
					Description: "Action to be taken on matching packets for NULL routes",
					Computed:    true,
				},
				"ip_address": &schema.Schema{
					Type:         schema.TypeString,
					Description:  "Next Hop IP",
					Optional:     true,
					ValidateFunc: validateSingleIP(),
				},
				"logical_router_port_id": &schema.Schema{
					Type:        schema.TypeString,
					Description: "Logical router port id",
					Optional:    true,
				},
			},
		},
	}
}

func getNextHopsFromSchema(d *schema.ResourceData) []manager.StaticRouteNextHop {
	hops := d.Get("next_hop").([]interface{})
	var nextHopsList []manager.StaticRouteNextHop
	for _, hop := range hops {
		data := hop.(map[string]interface{})
		administrative_distance := int64(data["administrative_distance"].(int))
		bfd_enabled := data["bfd_enabled"].(bool)
		blackhole_action := data["blackhole_action"].(string)
		ip_address := data["ip_address"].(string)
		logical_router_port_id := &common.ResourceReference{
			TargetType: "LogicalPort",
			TargetId:   data["logical_router_port_id"].(string),
		}
		elem := manager.StaticRouteNextHop{
			AdministrativeDistance: administrative_distance,
			BfdEnabled:             bfd_enabled,
			BlackholeAction:        blackhole_action,
			IpAddress:              ip_address,
			LogicalRouterPortId:    logical_router_port_id,
		}
		nextHopsList = append(nextHopsList, elem)
	}
	return nextHopsList
}

func setNextHopsInSchema(d *schema.ResourceData, next_hops []manager.StaticRouteNextHop) {
	var nextHopsList []map[string]interface{}
	for _, static_route_next_hop := range next_hops {
		elem := make(map[string]interface{})
		elem["administrative_distance"] = static_route_next_hop.AdministrativeDistance
		elem["bfd_enabled"] = static_route_next_hop.BfdEnabled
		elem["blackhole_action"] = static_route_next_hop.BlackholeAction
		elem["ip_address"] = static_route_next_hop.IpAddress
		elem["logical_router_port_id"] = static_route_next_hop.LogicalRouterPortId
		nextHopsList = append(nextHopsList, elem)
	}
	d.Set("next_hop", nextHopsList)
}

func resourceNsxtStaticRouteCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	logical_router_id := d.Get("logical_router_id").(string)
	if logical_router_id == "" {
		return fmt.Errorf("Error obtaining logical router id during static route creation")
	}

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	network := d.Get("network").(string)
	next_hops := getNextHopsFromSchema(d)
	static_route := manager.StaticRoute{
		Description:     description,
		DisplayName:     display_name,
		Tags:            tags,
		LogicalRouterId: logical_router_id,
		Network:         network,
		NextHops:        next_hops,
	}

	static_route, resp, err := nsxClient.LogicalRoutingAndServicesApi.AddStaticRoute(nsxClient.Context, logical_router_id, static_route)

	if err != nil {
		return fmt.Errorf("Error during StaticRoute create on router %s: %v", logical_router_id, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status returned during StaticRoute create on router %s: %v", logical_router_id, resp.StatusCode)
	}
	d.SetId(static_route.Id)

	return resourceNsxtStaticRouteRead(d, m)
}

func resourceNsxtStaticRouteRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	logical_router_id := d.Get("logical_router_id").(string)
	if logical_router_id == "" {
		return fmt.Errorf("Error obtaining logical router id during static route read")
	}

	static_route, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadStaticRoute(nsxClient.Context, logical_router_id, id)
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] StaticRoute %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during StaticRoute read: %v", err)
	}

	d.Set("revision", static_route.Revision)
	d.Set("description", static_route.Description)
	d.Set("display_name", static_route.DisplayName)
	setTagsInSchema(d, static_route.Tags)
	d.Set("logical_router_id", static_route.LogicalRouterId)
	d.Set("network", static_route.Network)
	setNextHopsInSchema(d, static_route.NextHops)

	return nil
}

func resourceNsxtStaticRouteUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	logical_router_id := d.Get("logical_router_id").(string)
	if logical_router_id == "" {
		return fmt.Errorf("Error obtaining logical router id during static route update")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	network := d.Get("network").(string)
	next_hops := getNextHopsFromSchema(d)
	static_route := manager.StaticRoute{
		Revision:        revision,
		Description:     description,
		DisplayName:     display_name,
		Tags:            tags,
		LogicalRouterId: logical_router_id,
		Network:         network,
		NextHops:        next_hops,
	}

	static_route, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateStaticRoute(nsxClient.Context, logical_router_id, id, static_route)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during StaticRoute update: %v", err)
	}

	return resourceNsxtStaticRouteRead(d, m)
}

func resourceNsxtStaticRouteDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	logical_router_id := d.Get("logical_router_id").(string)
	if logical_router_id == "" {
		return fmt.Errorf("Error obtaining logical router id during static route deletion")
	}

	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteStaticRoute(nsxClient.Context, logical_router_id, id)
	if err != nil {
		return fmt.Errorf("Error during StaticRoute delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] StaticRoute %s for router %s not found", id, logical_router_id)
		d.SetId("")
	}
	return nil
}
