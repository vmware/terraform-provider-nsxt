/* Copyright © 2017 VMware, Inc. All Rights Reserved.
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

func resourceNsxtIpProtocolNsService() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtIpProtocolNsServiceCreate,
		Read:   resourceNsxtIpProtocolNsServiceRead,
		Update: resourceNsxtIpProtocolNsServiceUpdate,
		Delete: resourceNsxtIpProtocolNsServiceDelete,

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
			"default_service": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "The default NSServices are created in the system by default. These NSServices can't be modified/deleted",
				Computed:    true,
			},
			"protocol": &schema.Schema{
				Type:         schema.TypeInt,
				Description:  "Ip protocol number",
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 255),
			},
		},
	}
}

func resourceNsxtIpProtocolNsServiceCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	defaultService := d.Get("default_service").(bool)
	protocol := int64(d.Get("protocol").(int))

	nsService := manager.IpProtocolNsService{
		NsService: manager.NsService{
			Description:    description,
			DisplayName:    displayName,
			Tags:           tags,
			DefaultService: defaultService,
		},
		NsserviceElement: manager.IpProtocolNsServiceEntry{
			ResourceType:   "IPProtocolNSService",
			ProtocolNumber: protocol,
		},
	}

	nsService, resp, err := nsxClient.GroupingObjectsApi.CreateIpProtocolNSService(nsxClient.Context, nsService)

	if err != nil {
		return fmt.Errorf("Error during NsService create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during NsService create: %v", resp.StatusCode)
	}
	d.SetId(nsService.Id)
	return resourceNsxtIpProtocolNsServiceRead(d, m)
}

func resourceNsxtIpProtocolNsServiceRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	nsService, resp, err := nsxClient.GroupingObjectsApi.ReadIpProtocolNSService(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] NsService %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during NsService read: %v", err)
	}

	nsserviceElement := nsService.NsserviceElement

	d.Set("revision", nsService.Revision)
	d.Set("description", nsService.Description)
	d.Set("display_name", nsService.DisplayName)
	setTagsInSchema(d, nsService.Tags)
	d.Set("default_service", nsService.DefaultService)
	d.Set("protocol", nsserviceElement.ProtocolNumber)

	return nil
}

func resourceNsxtIpProtocolNsServiceUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	defaultService := d.Get("default_service").(bool)
	revision := int64(d.Get("revision").(int))
	protocol := int64(d.Get("protocol").(int))

	nsService := manager.IpProtocolNsService{
		NsService: manager.NsService{
			Description:    description,
			DisplayName:    displayName,
			Tags:           tags,
			DefaultService: defaultService,
			Revision:       revision,
		},
		NsserviceElement: manager.IpProtocolNsServiceEntry{
			ResourceType:   "IPProtocolNSService",
			ProtocolNumber: protocol,
		},
	}

	nsService, resp, err := nsxClient.GroupingObjectsApi.UpdateIpProtocolNSService(nsxClient.Context, id, nsService)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during NsService update: %v %v", err, resp)
	}

	return resourceNsxtIpProtocolNsServiceRead(d, m)
}

func resourceNsxtIpProtocolNsServiceDelete(d *schema.ResourceData, m interface{}) error {
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
		log.Printf("[DEBUG] NsService %s not found", id)
		d.SetId("")
	}
	return nil
}
