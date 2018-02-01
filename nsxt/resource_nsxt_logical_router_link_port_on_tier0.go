package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func resourceLogicalRouterLinkPortOnTier0() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogicalRouterLinkPortOnTier0Create,
		Read:   resourceLogicalRouterLinkPortOnTier0Read,
		Update: resourceLogicalRouterLinkPortOnTier0Update,
		Delete: resourceLogicalRouterLinkPortOnTier0Delete,

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
			"linked_logical_router_port_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Identifier for port on logical router to connect to",
				Computed:    true,
			},
			"service_bindings": getResourceReferencesSchema(false, false, []string{"LogicalService"}),
		},
	}
}

func resourceLogicalRouterLinkPortOnTier0Create(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	logical_router_id := d.Get("logical_router_id").(string)
	linked_logical_router_port_id := d.Get("linked_logical_router_port_id").(string)
	service_bindings := getServiceBindingsFromSchema(d, "service_bindings")
	logical_router_link_port := manager.LogicalRouterLinkPortOnTier0{
		Description:               description,
		DisplayName:               display_name,
		Tags:                      tags,
		LogicalRouterId:           logical_router_id,
		LinkedLogicalRouterPortId: linked_logical_router_port_id,
		ServiceBindings:           service_bindings,
	}
	logical_router_link_port, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateLogicalRouterLinkPortOnTier0(nsxClient.Context, logical_router_link_port)

	if err != nil {
		return fmt.Errorf("Error during LogicalRouterLinkPortOnTier0 create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during LogicalRouterLinkPortOnTier0 create: %v", resp.StatusCode)
	}
	d.SetId(logical_router_link_port.Id)

	return resourceLogicalRouterLinkPortOnTier0Read(d, m)
}

func resourceLogicalRouterLinkPortOnTier0Read(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router link port on tier0 id")
	}

	logical_router_link_port, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouterLinkPortOnTier0(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalRouterLinkPortOnTier0 %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during LogicalRouterLinkPortOnTier0 read: %v", err)
	}

	d.Set("revision", logical_router_link_port.Revision)
	d.Set("system_owned", logical_router_link_port.SystemOwned)
	d.Set("description", logical_router_link_port.Description)
	d.Set("display_name", logical_router_link_port.DisplayName)
	setTagsInSchema(d, logical_router_link_port.Tags)
	d.Set("logical_router_id", logical_router_link_port.LogicalRouterId)
	d.Set("linked_logical_router_port_id", logical_router_link_port.LinkedLogicalRouterPortId)
	setServiceBindingsInSchema(d, logical_router_link_port.ServiceBindings, "service_bindings")
	d.Set("resource_type", logical_router_link_port.ResourceType)

	return nil
}

func resourceLogicalRouterLinkPortOnTier0Update(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router link port on tier0 id")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	logical_router_id := d.Get("logical_router_id").(string)
	linked_logical_router_port_id := d.Get("linked_logical_router_port_id").(string)
	service_bindings := getServiceBindingsFromSchema(d, "service_bindings")
	logical_router_link_port := manager.LogicalRouterLinkPortOnTier0{
		Revision:                  revision,
		Description:               description,
		DisplayName:               display_name,
		Tags:                      tags,
		LogicalRouterId:           logical_router_id,
		LinkedLogicalRouterPortId: linked_logical_router_port_id,
		ServiceBindings:           service_bindings,
		ResourceType:              "LogicalRouterLinkPortOnTIER0",
	}

	logical_router_link_port, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateLogicalRouterLinkPortOnTier0(nsxClient.Context, id, logical_router_link_port)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LogicalRouterLinkPortOnTier0 update: %v", err)
	}

	return resourceLogicalRouterLinkPortOnTier0Read(d, m)
}

func resourceLogicalRouterLinkPortOnTier0Delete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router link port on tier0 id")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteLogicalRouterPort(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during LogicalRouterLinkPortOnTier0 delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalRouterLinkPortOnTier0 %s not found", id)
		d.SetId("")
	}

	return nil
}
