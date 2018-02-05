/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func resourceLogicalRouterDownLinkPort() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogicalRouterDownLinkPortCreate,
		Read:   resourceLogicalRouterDownLinkPortRead,
		Update: resourceLogicalRouterDownLinkPortUpdate,
		Delete: resourceLogicalRouterDownLinkPortDelete,

		Schema: map[string]*schema.Schema{
			"revision":     getRevisionSchema(),
			"system_owned": getSystemOwnedSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Defaults to ID if not set",
				Optional:    true,
			},
			"tags": getTagsSchema(),
			"logical_router_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Identifier for logical router on which this port is created",
				Required:    true,
				ForceNew:    true,
			},
			"linked_logical_switch_port_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Identifier for port on logical switch to connect to",
				Required:    true,
				ForceNew:    true,
			},
			"subnets": getIpSubnetsSchema(true, false),
			"mac_address": &schema.Schema{
				Type:        schema.TypeString,
				Description: "MAC address",
				Optional:    true,
			},
			"urpf_mode": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Unicast Reverse Path Forwarding mode",
				Optional:    true,
			},
			"service_bindings": getResourceReferencesSchema(false, false, []string{"LogicalService"}),
		},
	}
}

func resourceLogicalRouterDownLinkPortCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	logical_router_id := d.Get("logical_router_id").(string)
	mac_address := d.Get("mac_address").(string)
	linked_logical_switch_port_id := d.Get("linked_logical_switch_port_id").(string)
	subnets := getIpSubnetsFromSchema(d)
	urpf_mode := d.Get("urpf_mode").(string)
	service_bindings := getServiceBindingsFromSchema(d, "service_bindings")
	logical_router_down_link_port := manager.LogicalRouterDownLinkPort{
		Description:               description,
		DisplayName:               display_name,
		Tags:                      tags,
		LogicalRouterId:           logical_router_id,
		MacAddress:                mac_address,
		LinkedLogicalSwitchPortId: makeResourceReference("LogicalPort", linked_logical_switch_port_id),
		Subnets:                   subnets,
		UrpfMode:                  urpf_mode,
		ServiceBindings:           service_bindings,
	}

	logical_router_down_link_port, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateLogicalRouterDownLinkPort(nsxClient.Context, logical_router_down_link_port)

	if err != nil {
		return fmt.Errorf("Error during LogicalRouterDownLinkPort create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during LogicalRouterDownLinkPort create: %v", resp.StatusCode)
	}
	d.SetId(logical_router_down_link_port.Id)

	return resourceLogicalRouterDownLinkPortRead(d, m)
}

func resourceLogicalRouterDownLinkPortRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router downlink port id while reading")
	}

	logical_router_down_link_port, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouterDownLinkPort(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalRouterDownLinkPort %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during LogicalRouterDownLinkPort read: %v", err)
	}

	d.Set("Revision", logical_router_down_link_port.Revision)
	d.Set("SystemOwned", logical_router_down_link_port.SystemOwned)
	d.Set("Description", logical_router_down_link_port.Description)
	d.Set("DisplayName", logical_router_down_link_port.DisplayName)
	setTagsInSchema(d, logical_router_down_link_port.Tags)
	d.Set("LogicalRouterId", logical_router_down_link_port.LogicalRouterId)
	d.Set("MacAddress", logical_router_down_link_port.MacAddress)
	d.Set("LinkedLogicalSwitchPortId", logical_router_down_link_port.LinkedLogicalSwitchPortId.TargetId)
	setIpSubnetsInSchema(d, logical_router_down_link_port.Subnets)
	d.Set("UrpfMode", logical_router_down_link_port.UrpfMode)
	d.Set("ResourceType", logical_router_down_link_port.ResourceType)
	setServiceBindingsInSchema(d, logical_router_down_link_port.ServiceBindings, "service_bindings")

	return nil
}

func resourceLogicalRouterDownLinkPortUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router downlink port id while updating")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	logical_router_id := d.Get("logical_router_id").(string)
	linked_logical_switch_port_id := d.Get("linked_logical_switch_port_id").(string)
	subnets := getIpSubnetsFromSchema(d)
	mac_address := d.Get("mac_address").(string)
	urpf_mode := d.Get("urpf_mode").(string)
	service_bindings := getServiceBindingsFromSchema(d, "service_bindings")
	logical_router_down_link_port := manager.LogicalRouterDownLinkPort{
		Revision:                  revision,
		Description:               description,
		DisplayName:               display_name,
		Tags:                      tags,
		LogicalRouterId:           logical_router_id,
		MacAddress:                mac_address,
		LinkedLogicalSwitchPortId: makeResourceReference("LogicalPort", linked_logical_switch_port_id),
		Subnets:                   subnets,
		UrpfMode:                  urpf_mode,
		ServiceBindings:           service_bindings,
		ResourceType:              "LogicalRouterDownLinkPort",
	}

	logical_router_down_link_port, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateLogicalRouterDownLinkPort(nsxClient.Context, id, logical_router_down_link_port)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LogicalRouterDownLinkPort update: %v", err)
	}

	return resourceLogicalRouterDownLinkPortRead(d, m)
}

func resourceLogicalRouterDownLinkPortDelete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router downlink port id while deleting")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteLogicalRouterPort(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during LogicalRouterDownLinkPort delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalRouterDownLinkPort %s not found", id)
		d.SetId("")
	}

	return nil
}
