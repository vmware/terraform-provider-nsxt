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

func resourceIgmpTypeNsService() *schema.Resource {
	return &schema.Resource{
		Create: resourceIgmpTypeNsServiceCreate,
		Read:   resourceIgmpTypeNsServiceRead,
		Update: resourceIgmpTypeNsServiceUpdate,
		Delete: resourceIgmpTypeNsServiceDelete,

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
		},
	}
}

func resourceIgmpTypeNsServiceCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	default_service := d.Get("default_service").(bool)

	ns_service := manager.IgmpTypeNsService{
		NsService: manager.NsService{
			Description:    description,
			DisplayName:    display_name,
			Tags:           tags,
			DefaultService: default_service,
		},
		NsserviceElement: manager.IgmpTypeNsServiceEntry{
			ResourceType: "IGMPTypeNSService",
		},
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.CreateIgmpTypeNSService(nsxClient.Context, ns_service)

	if err != nil {
		return fmt.Errorf("Error during NsService create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during NsService create: %v", resp.StatusCode)
	}
	d.SetId(ns_service.Id)
	return resourceIgmpTypeNsServiceRead(d, m)
}

func resourceIgmpTypeNsServiceRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.ReadIgmpTypeNSService(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("NsService %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during NsService read: %v", err)
	}

	d.Set("revision", ns_service.Revision)
	d.Set("description", ns_service.Description)
	d.Set("display_name", ns_service.DisplayName)
	setTagsInSchema(d, ns_service.Tags)
	d.Set("default_service", ns_service.DefaultService)

	return nil
}

func resourceIgmpTypeNsServiceUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	default_service := d.Get("default_service").(bool)
	revision := int64(d.Get("revision").(int))

	ns_service := manager.IgmpTypeNsService{
		NsService: manager.NsService{
			Description:    description,
			DisplayName:    display_name,
			Tags:           tags,
			DefaultService: default_service,
			Revision:       revision,
		},
		NsserviceElement: manager.IgmpTypeNsServiceEntry{
			ResourceType: "IGMPTypeNSService",
		},
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.UpdateIgmpTypeNSService(nsxClient.Context, id, ns_service)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during NsService update: %v %v", err, resp)
	}

	return resourceIgmpTypeNsServiceRead(d, m)
}

func resourceIgmpTypeNsServiceDelete(d *schema.ResourceData, m interface{}) error {

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
