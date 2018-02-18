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

func resourceLogicalRouterLinkPortOnTier1() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogicalRouterLinkPortOnTier1Create,
		Read:   resourceLogicalRouterLinkPortOnTier1Read,
		Update: resourceLogicalRouterLinkPortOnTier1Update,
		Delete: resourceLogicalRouterLinkPortOnTier1Delete,

		Schema: map[string]*schema.Schema{
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
			},
			"tags": getTagsSchema(),
			"logical_router_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Identifier for logical router on which this port is created",
				Required:    true,
				ForceNew:    true,
			},
			"linked_logical_router_port_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Identifier for port on logical router to connect to",
				Required:    true,
				ForceNew:    true,
			},
			"service_binding": getResourceReferencesSchema(false, false, []string{"LogicalService"}),
		},
	}
}

func resourceLogicalRouterLinkPortOnTier1Create(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	logical_router_id := d.Get("logical_router_id").(string)
	linked_logical_router_port_id := d.Get("linked_logical_router_port_id").(string)
	service_binding := getServiceBindingsFromSchema(d, "service_binding")
	logical_router_link_port := manager.LogicalRouterLinkPortOnTier1{
		Description:               description,
		DisplayName:               display_name,
		Tags:                      tags,
		LogicalRouterId:           logical_router_id,
		LinkedLogicalRouterPortId: makeResourceReference("LogicalPort", linked_logical_router_port_id),
		ServiceBindings:           service_binding,
	}

	logical_router_link_port, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateLogicalRouterLinkPortOnTier1(nsxClient.Context, logical_router_link_port)

	if err != nil {
		return fmt.Errorf("Error during LogicalRouterLinkPortOnTier1 create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during LogicalRouterLinkPortOnTier1 create: %v", resp.StatusCode)
	}
	d.SetId(logical_router_link_port.Id)

	return resourceLogicalRouterLinkPortOnTier1Read(d, m)
}

func resourceLogicalRouterLinkPortOnTier1Read(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router link port on tier1 id")
	}

	logical_router_link_port, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouterLinkPortOnTier1(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalRouterLinkPortOnTier1 %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during LogicalRouterLinkPortOnTier1 read: %v", err)
	}

	d.Set("revision", logical_router_link_port.Revision)
	d.Set("description", logical_router_link_port.Description)
	d.Set("display_name", logical_router_link_port.DisplayName)
	setTagsInSchema(d, logical_router_link_port.Tags)
	d.Set("logical_router_id", logical_router_link_port.LogicalRouterId)
	d.Set("linked_logical_router_port_id", logical_router_link_port.LinkedLogicalRouterPortId)
	setServiceBindingsInSchema(d, logical_router_link_port.ServiceBindings, "service_binding")

	return nil
}

func resourceLogicalRouterLinkPortOnTier1Update(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router link port on tier1 id")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	logical_router_id := d.Get("logical_router_id").(string)
	linked_logical_router_port_id := d.Get("linked_logical_router_port_id").(string)
	service_binding := getServiceBindingsFromSchema(d, "service_binding")
	logical_router_link_port := manager.LogicalRouterLinkPortOnTier1{
		Revision:                  revision,
		Description:               description,
		DisplayName:               display_name,
		Tags:                      tags,
		LogicalRouterId:           logical_router_id,
		LinkedLogicalRouterPortId: makeResourceReference("LogicalPort", linked_logical_router_port_id),
		ServiceBindings:           service_binding,
		ResourceType:              "LogicalRouterLinkPortOnTIER1",
	}

	logical_router_link_port, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateLogicalRouterLinkPortOnTier1(nsxClient.Context, id, logical_router_link_port)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LogicalRouterLinkPortOnTier1 %v update: %v (%+v)", id, err, resp)
	}

	return resourceLogicalRouterLinkPortOnTier1Read(d, m)
}

func resourceLogicalRouterLinkPortOnTier1Delete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router link port on tier1 id")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteLogicalRouterPort(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during LogicalRouterLinkPortOnTier1 delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalRouterLinkPortOnTier1 %s not found", id)
		d.SetId("")
	}

	return nil
}
