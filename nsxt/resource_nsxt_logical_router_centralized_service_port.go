/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"log"
	"net/http"
)

func resourceNsxtLogicalRouterCentralizedServicePort() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtLogicalRouterCentralizedServicePortCreate,
		Read:   resourceNsxtLogicalRouterCentralizedServicePortRead,
		Update: resourceNsxtLogicalRouterCentralizedServicePortUpdate,
		Delete: resourceNsxtLogicalRouterCentralizedServicePortDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"tag": getTagsSchema(),
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
			"ip_address": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Logical router port subnet (ipAddress / prefix length)",
				Required:     true,
				ValidateFunc: validatePortAddress(),
			},
			"urpf_mode": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Unicast Reverse Path Forwarding mode",
				Optional:     true,
				Default:      "STRICT",
				ValidateFunc: validation.StringInSlice(logicalRouterPortUrpfModeValues, false),
			},
		},
	}
}

func resourceNsxtLogicalRouterCentralizedServicePortCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	logicalRouterID := d.Get("logical_router_id").(string)
	linkedLogicalSwitchPortID := d.Get("linked_logical_switch_port_id").(string)
	subnets := getIPSubnetsFromCidr(d.Get("ip_address").(string))
	urpfMode := d.Get("urpf_mode").(string)
	LogicalRouterCentralizedServicePort := manager.LogicalRouterCentralizedServicePort{
		Description:               description,
		DisplayName:               displayName,
		Tags:                      tags,
		LogicalRouterId:           logicalRouterID,
		LinkedLogicalSwitchPortId: makeResourceReference("LogicalPort", linkedLogicalSwitchPortID),
		Subnets:                   subnets,
		UrpfMode:                  urpfMode,
	}

	LogicalRouterCentralizedServicePort, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateLogicalRouterCentralizedServicePort(nsxClient.Context, LogicalRouterCentralizedServicePort)

	if err != nil {
		return fmt.Errorf("Error during LogicalRouterCentralizedServicePort create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during LogicalRouterCentralizedServicePort create: %v", resp.StatusCode)
	}
	d.SetId(LogicalRouterCentralizedServicePort.Id)

	return resourceNsxtLogicalRouterCentralizedServicePortRead(d, m)
}

func resourceNsxtLogicalRouterCentralizedServicePortRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router centralized port id while reading")
	}

	LogicalRouterCentralizedServicePort, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouterCentralizedServicePort(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LogicalRouterCentralizedServicePort read: %v", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LogicalRouterCentralizedServicePort %s not found", id)
		d.SetId("")
		return nil
	}

	d.Set("revision", LogicalRouterCentralizedServicePort.Revision)
	d.Set("description", LogicalRouterCentralizedServicePort.Description)
	d.Set("display_name", LogicalRouterCentralizedServicePort.DisplayName)
	setTagsInSchema(d, LogicalRouterCentralizedServicePort.Tags)
	d.Set("logical_router_id", LogicalRouterCentralizedServicePort.LogicalRouterId)
	d.Set("linked_logical_switch_port_id", LogicalRouterCentralizedServicePort.LinkedLogicalSwitchPortId.TargetId)
	setIPSubnetsInSchema(d, LogicalRouterCentralizedServicePort.Subnets)
	d.Set("urpf_mode", LogicalRouterCentralizedServicePort.UrpfMode)

	return nil
}

func resourceNsxtLogicalRouterCentralizedServicePortUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router centralized port id while updating")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	logicalRouterID := d.Get("logical_router_id").(string)
	linkedLogicalSwitchPortID := d.Get("linked_logical_switch_port_id").(string)
	subnets := getIPSubnetsFromCidr(d.Get("ip_address").(string))
	urpfMode := d.Get("urpf_mode").(string)
	LogicalRouterCentralizedServicePort := manager.LogicalRouterCentralizedServicePort{
		Revision:                  revision,
		Description:               description,
		DisplayName:               displayName,
		Tags:                      tags,
		LogicalRouterId:           logicalRouterID,
		LinkedLogicalSwitchPortId: makeResourceReference("LogicalPort", linkedLogicalSwitchPortID),
		Subnets:                   subnets,
		UrpfMode:                  urpfMode,
		ResourceType:              "LogicalRouterCentralizedServicePort",
	}

	LogicalRouterCentralizedServicePort, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateLogicalRouterCentralizedServicePort(nsxClient.Context, id, LogicalRouterCentralizedServicePort)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LogicalRouterCentralizedServicePort update: %v", err)
	}

	return resourceNsxtLogicalRouterCentralizedServicePortRead(d, m)
}

func resourceNsxtLogicalRouterCentralizedServicePortDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router centralized port id while deleting")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteLogicalRouterPort(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during LogicalRouterCentralizedServicePort delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LogicalRouterCentralizedServicePort %s not found", id)
		d.SetId("")
	}

	return nil
}
