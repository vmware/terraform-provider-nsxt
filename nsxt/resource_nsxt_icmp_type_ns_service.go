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

var icmpProtocolValues = []string{"ICMPv4", "ICMPv6"}

func resourceIcmpTypeNsService() *schema.Resource {
	return &schema.Resource{
		Create: resourceIcmpTypeNsServiceCreate,
		Read:   resourceIcmpTypeNsServiceRead,
		Update: resourceIcmpTypeNsServiceUpdate,
		Delete: resourceIcmpTypeNsServiceDelete,

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
			"default_service": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "The default NSServices are created in the system by default. These NSServices can't be modified/deleted",
				Computed:    true,
			},
			"icmp_code": &schema.Schema{
				Type:         schema.TypeInt,
				Description:  "ICMP message code",
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 255),
			},
			"icmp_type": &schema.Schema{
				Type:         schema.TypeInt,
				Description:  "ICMP message type",
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 255),
			},
			"protocol": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Version of ICMP protocol (ICMPv4/ICMPv6)",
				Required:     true,
				ValidateFunc: validation.StringInSlice(icmpProtocolValues, false),
			},
		},
	}
}

func resourceIcmpTypeNsServiceCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	default_service := d.Get("default_service").(bool)
	icmp_code := int64(d.Get("icmp_code").(int))
	icmp_type := int64(d.Get("icmp_type").(int))
	protocol := d.Get("protocol").(string)

	ns_service := manager.IcmpTypeNsService{
		NsService: manager.NsService{
			Description:    description,
			DisplayName:    display_name,
			Tags:           tags,
			DefaultService: default_service,
		},
		NsserviceElement: manager.IcmpTypeNsServiceEntry{
			ResourceType: "ICMPTypeNSService",
			IcmpCode:     icmp_code,
			IcmpType:     icmp_type,
			Protocol:     protocol,
		},
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.CreateIcmpTypeNSService(nsxClient.Context, ns_service)

	if err != nil {
		return fmt.Errorf("Error during NsService create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during NsService create: %v", resp.StatusCode)
	}
	d.SetId(ns_service.Id)
	return resourceIcmpTypeNsServiceRead(d, m)
}

func resourceIcmpTypeNsServiceRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.ReadIcmpTypeNSService(nsxClient.Context, id)
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
	d.Set("icmp_type", nsservice_element.IcmpType)
	d.Set("icmp_code", nsservice_element.IcmpCode)
	d.Set("protocol", nsservice_element.Protocol)

	return nil
}

func resourceIcmpTypeNsServiceUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	default_service := d.Get("default_service").(bool)
	icmp_code := int64(d.Get("icmp_code").(int))
	icmp_type := int64(d.Get("icmp_type").(int))
	protocol := d.Get("protocol").(string)
	revision := int64(d.Get("revision").(int))

	ns_service := manager.IcmpTypeNsService{
		NsService: manager.NsService{
			Description:    description,
			DisplayName:    display_name,
			Tags:           tags,
			DefaultService: default_service,
			Revision:       revision,
		},
		NsserviceElement: manager.IcmpTypeNsServiceEntry{
			ResourceType: "ICMPTypeNSService",
			IcmpCode:     icmp_code,
			IcmpType:     icmp_type,
			Protocol:     protocol,
		},
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.UpdateIcmpTypeNSService(nsxClient.Context, id, ns_service)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during NsService update: %v %v", err, resp)
	}

	return resourceIcmpTypeNsServiceRead(d, m)
}

func resourceIcmpTypeNsServiceDelete(d *schema.ResourceData, m interface{}) error {

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
