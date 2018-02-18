/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

var ipProtocolValues = []string{"TCP", "UDP"}

func resourceL4PortSetNsService() *schema.Resource {
	return &schema.Resource{
		Create: resourceL4PortSetNsServiceCreate,
		Read:   resourceL4PortSetNsServiceRead,
		Update: resourceL4PortSetNsServiceUpdate,
		Delete: resourceL4PortSetNsServiceDelete,

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
			"tag": getTagsSchema(),
			"default_service": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "The default NSServices are created in the system by default. These NSServices can't be modified/deleted",
				Computed:    true,
			},
			"destination_ports": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "Set of destination ports",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validatePortRange(),
				},
				Optional: true,
			},
			"source_ports": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "Set of source ports",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validatePortRange(),
				},
				Optional: true,
			},
			"protocol": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "L4 Protocol",
				Required:     true,
				ValidateFunc: validation.StringInSlice(ipProtocolValues, false),
			},
		},
	}
}

func resourceL4PortSetNsServiceCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	default_service := d.Get("default_service").(bool)
	l4_protocol := d.Get("protocol").(string)
	source_ports := getStringListFromSchemaSet(d, "source_ports")
	destination_ports := getStringListFromSchemaSet(d, "destination_ports")

	ns_service := manager.L4PortSetNsService{
		NsService: manager.NsService{
			Description:    description,
			DisplayName:    display_name,
			Tags:           tags,
			DefaultService: default_service,
		},
		NsserviceElement: manager.L4PortSetNsServiceEntry{
			ResourceType:     "L4PortSetNSService",
			L4Protocol:       l4_protocol,
			DestinationPorts: destination_ports,
			SourcePorts:      source_ports,
		},
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.CreateL4PortSetNSService(nsxClient.Context, ns_service)

	if err != nil {
		return fmt.Errorf("Error during NsService create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during NsService create: %v", resp.StatusCode)
	}
	d.SetId(ns_service.Id)
	return resourceL4PortSetNsServiceRead(d, m)
}

func resourceL4PortSetNsServiceRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.ReadL4PortSetNSService(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("NsService %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during NsService read: %v", err)
	}

	nsservice_element := ns_service.NsserviceElement

	d.Set("revision", ns_service.Revision)
	d.Set("description", ns_service.Description)
	d.Set("display_name", ns_service.DisplayName)
	setTagsInSchema(d, ns_service.Tags)
	d.Set("default_service", ns_service.DefaultService)
	d.Set("protocol", nsservice_element.L4Protocol)
	d.Set("destination_ports", nsservice_element.DestinationPorts)
	d.Set("source_ports", nsservice_element.SourcePorts)

	return nil
}

func resourceL4PortSetNsServiceUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	default_service := d.Get("default_service").(bool)
	l4_protocol := d.Get("protocol").(string)
	source_ports := getStringListFromSchemaSet(d, "source_ports")
	destination_ports := getStringListFromSchemaSet(d, "destination_ports")
	revision := int64(d.Get("revision").(int))

	ns_service := manager.L4PortSetNsService{
		NsService: manager.NsService{
			Description:    description,
			DisplayName:    display_name,
			Tags:           tags,
			DefaultService: default_service,
			Revision:       revision,
		},
		NsserviceElement: manager.L4PortSetNsServiceEntry{
			ResourceType:     "L4PortSetNSService",
			L4Protocol:       l4_protocol,
			DestinationPorts: destination_ports,
			SourcePorts:      source_ports,
		},
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.UpdateL4PortSetNSService(nsxClient.Context, id, ns_service)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during NsService update: %v %v", err, resp)
	}

	return resourceL4PortSetNsServiceRead(d, m)
}

func resourceL4PortSetNsServiceDelete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.GroupingObjectsApi.DeleteNSService(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during NsService delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("NsService %s not found", id)
		d.SetId("")
	}
	return nil
}
