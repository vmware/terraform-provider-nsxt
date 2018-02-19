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

var algTypeValues = []string{"ORACLE_TNS", "FTP", "SUN_RPC_TCP", "SUN_RPC_UDP", "MS_RPC_TCP", "MS_RPC_UDP", "NBNS_BROADCAST", "NBDG_BROADCAST", "TFTP"}

func resourceNsxtAlgorithmTypeNsService() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtAlgorithmTypeNsServiceCreate,
		Read:   resourceNsxtAlgorithmTypeNsServiceRead,
		Update: resourceNsxtAlgorithmTypeNsServiceUpdate,
		Delete: resourceNsxtAlgorithmTypeNsServiceDelete,

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
			"destination_ports": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "A single destination port",
				Required:     true,
				ValidateFunc: validateSinglePort(),
			},
			"source_ports": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "Set of source ports or ranges",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validatePortRange(),
				},
				Optional: true,
			},
			"algorithm": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Algorithm",
				Required:     true,
				ValidateFunc: validation.StringInSlice(algTypeValues, false),
			},
		},
	}
}

func resourceNsxtAlgorithmTypeNsServiceCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	default_service := d.Get("default_service").(bool)
	alg := d.Get("algorithm").(string)
	source_ports := getStringListFromSchemaSet(d, "source_ports")
	destination_ports := make([]string, 0, 1)
	destination_ports = append(destination_ports, d.Get("destination_ports").(string))

	ns_service := manager.AlgTypeNsService{
		NsService: manager.NsService{
			Description:    description,
			DisplayName:    display_name,
			Tags:           tags,
			DefaultService: default_service,
		},
		NsserviceElement: manager.AlgTypeNsServiceEntry{
			ResourceType:     "ALGTypeNSService",
			Alg:              alg,
			DestinationPorts: destination_ports,
			SourcePorts:      source_ports,
		},
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.CreateAlgTypeNSService(nsxClient.Context, ns_service)

	if err != nil {
		return fmt.Errorf("Error during NsService create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during NsService create: %v", resp.StatusCode)
	}
	d.SetId(ns_service.Id)
	return resourceNsxtAlgorithmTypeNsServiceRead(d, m)
}

func resourceNsxtAlgorithmTypeNsServiceRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.ReadAlgTypeNSService(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] NsService %s not found", id)
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
	d.Set("algorithm", nsservice_element.Alg)
	d.Set("destination_ports", nsservice_element.DestinationPorts)
	d.Set("source_ports", nsservice_element.SourcePorts)

	return nil
}

func resourceNsxtAlgorithmTypeNsServiceUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining ns service id")
	}

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	default_service := d.Get("default_service").(bool)
	alg := d.Get("algorithm").(string)
	source_ports := getStringListFromSchemaSet(d, "source_ports")
	destination_ports := make([]string, 0, 1)
	destination_ports = append(destination_ports, d.Get("destination_ports").(string))
	revision := int64(d.Get("revision").(int))

	ns_service := manager.AlgTypeNsService{
		NsService: manager.NsService{
			Description:    description,
			DisplayName:    display_name,
			Tags:           tags,
			DefaultService: default_service,
			Revision:       revision,
		},
		NsserviceElement: manager.AlgTypeNsServiceEntry{
			ResourceType:     "ALGTypeNSService",
			Alg:              alg,
			DestinationPorts: destination_ports,
			SourcePorts:      source_ports,
		},
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.UpdateAlgTypeNSService(nsxClient.Context, id, ns_service)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during NsService update: %v %v", err, resp)
	}

	return resourceNsxtAlgorithmTypeNsServiceRead(d, m)
}

func resourceNsxtAlgorithmTypeNsServiceDelete(d *schema.ResourceData, m interface{}) error {
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
