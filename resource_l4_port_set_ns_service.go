package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func resourceL4PortSetNsService() *schema.Resource {
	return &schema.Resource{
		Create: resourceL4PortSetNsServiceCreate,
		Read:   resourceL4PortSetNsServiceRead,
		Update: resourceL4PortSetNsServiceUpdate,
		Delete: resourceL4PortSetNsServiceDelete,

		Schema: map[string]*schema.Schema{
			"Revision": GetRevisionSchema(),
			"CreateTime": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Timestamp of resource creation",
				Optional:    true,
				Computed:    true,
			},
			"CreateUser": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the user who created this resource",
				Optional:    true,
				Computed:    true,
			},
			"LastModifiedTime": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Timestamp of last modification",
				Optional:    true,
				Computed:    true,
			},
			"LastModifiedUser": &schema.Schema{
				Type:        schema.TypeString,
				Description: "ID of the user who last modified this resource",
				Optional:    true,
				Computed:    true,
			},
			"SystemOwned": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Indicates system owned resource",
				Optional:    true,
				Computed:    true,
			},
			"Description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"DisplayName": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Defaults to ID if not set",
				Optional:    true,
			},
			"Tags": GetTagsSchema(),
			"DefaultService": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "The default NSServices are created in the system by default. These NSServices can't be modified/deleted",
				Optional:    true,
			},
			"DestinationPorts": &schema.Schema{
				Type:        schema.TypeList,
				Description: "Set of destination ports",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"SourcePorts": &schema.Schema{
				Type:        schema.TypeList,
				Description: "Set of source ports",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"L4Protocol": &schema.Schema{
				Type:        schema.TypeString,
				Description: "L4 Protocol",
				Optional:    true,
			},
		},
	}
}

func resourceL4PortSetNsServiceCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*nsxt.APIClient)

	description := d.Get("Description").(string)
	display_name := d.Get("DisplayName").(string)
	tags := GetTagsFromSchema(d)
	default_service := d.Get("DefaultService").(bool)
	l4_protocol := d.Get("L4Protocol").(string)
	source_ports := Interface2StringList(d.Get("SourcePorts").([]interface{}))
	destination_ports := Interface2StringList(d.Get("DestinationPorts").([]interface{}))

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
		fmt.Printf("Unexpected status returned")
		return nil
	}
	d.SetId(ns_service.Id)
	return resourceL4PortSetNsServiceRead(d, m)
}

func resourceL4PortSetNsServiceRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*nsxt.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	ns_service, resp, err := nsxClient.GroupingObjectsApi.ReadL4PortSetNSService(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during NsService read: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("NsService not found")
		d.SetId("")
		return nil
	}

	nsservice_element := ns_service.NsserviceElement

	d.Set("Revision", ns_service.Revision)
	d.Set("CreateTime", ns_service.CreateTime)
	d.Set("CreateUser", ns_service.CreateUser)
	d.Set("LastModifiedTime", ns_service.LastModifiedTime)
	d.Set("LastModifiedUser", ns_service.LastModifiedUser)
	d.Set("SystemOwned", ns_service.SystemOwned)
	d.Set("Description", ns_service.Description)
	d.Set("DisplayName", ns_service.DisplayName)
	SetTagsInSchema(d, ns_service.Tags)
	d.Set("DefaultService", ns_service.DefaultService)
	d.Set("ResourceType", nsservice_element.ResourceType)
	//d.Set("DestinationPorts", flattenStringList(nsservice_element.DestinationPorts))
	//d.Set("SourcePorts", flattenStringList(nsservice_element.SourcePorts))
	d.Set("DestinationPorts", nsservice_element.DestinationPorts)
	d.Set("SourcePorts", nsservice_element.SourcePorts)

	return nil
}

func resourceL4PortSetNsServiceUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*nsxt.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	description := d.Get("Description").(string)
	display_name := d.Get("DisplayName").(string)
	tags := GetTagsFromSchema(d)
	default_service := d.Get("DefaultService").(bool)
	l4_protocol := d.Get("L4Protocol").(string)
	source_ports := Interface2StringList(d.Get("SourcePorts").([]interface{}))
	destination_ports := Interface2StringList(d.Get("DestinationPorts").([]interface{}))

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

	//ns_service, resp, err := nsxClient.GroupingObjectsApi.UpdateNSService(nsxClient.Context, id, ns_service)
	ns_service, resp, err := nsxClient.GroupingObjectsApi.UpdateL4PortSetNSService(nsxClient.Context, id, ns_service)

	if err != nil {
		return fmt.Errorf("Error during NsService update: %v %s", err, resp.Body)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("NsService not found")
		d.SetId("")
		return nil
	}
	return resourceL4PortSetNsServiceRead(d, m)
}

func resourceL4PortSetNsServiceDelete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*nsxt.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.GroupingObjectsApi.DeleteNSService(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during NsService delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("NsService not found")
		d.SetId("")
	}
	return nil
}
