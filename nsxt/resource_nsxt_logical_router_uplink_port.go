package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/go-vmware-nsxt/manager"
	"log"
	"net/http"
)

const LogicalRouterUpLinkPortResourceType = "LogicalRouterUpLinkPort"

func resourceNsxtLogicalRouterUpLinkPort() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtLogicalRouterUpLinkPortCreate,
		Read:   resourceNsxtLogicalRouterUpLinkPortRead,
		Update: resourceNsxtLogicalRouterUpLinkPortUpdate,
		Delete: resourceNsxtLogicalRouterUpLinkPortDelete,
		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": {
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"tag": getTagsSchema(),
			"logical_router_id": {
				Type:        schema.TypeString,
				Description: "Identifier for logical router on which this port is created",
				Required:    true,
				ForceNew:    true,
			},
			"linked_logical_switch_port_id": {
				Type:        schema.TypeString,
				Description: "Identifier for the logical switch port to connect to",
				Required:    true,
			},
			"edge_cluster_member_index": {
				Type:        schema.TypeList,
				Description: "Member index of the edge node on the cluster (start at 1)",
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"subnets": {
				Type:        schema.TypeList,
				Description: "Logical router port subnets",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_addresses": {
							Type:        schema.TypeList,
							Description: "IP Addresses",
							Required:    true,
							MinItems:    1,
							MaxItems:    2,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"prefix_length": {
							Type:         schema.TypeInt,
							Description:  "Subnet Prefix Length",
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 128),
						},
					},
				},
			},
			"urpf_mode": {
				Type:         schema.TypeString,
				Description:  "Unicast Reverse Path Forwarding mode",
				Optional:     true,
				Default:      "STRICT",
				ValidateFunc: validation.StringInSlice(logicalRouterPortUrpfModeValues, false),
			},
			"mac_address": {
				Type:        schema.TypeString,
				Description: "MAC address",
				Computed:    true,
			},
		},
	}
}

func resourceNsxtLogicalRouterUpLinkPortCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tag := getTagsFromSchema(d)
	logicalRouterID := d.Get("logical_router_id").(string)
	linkedLogicalSwitchPortID := d.Get("linked_logical_switch_port_id").(string)
	edgeClusterMemberIndex := getInt64ListFromSchemaList(d, "edge_cluster_member_index")
	subnets := getIPSubnetsFromSchema(d)
	urpfMode := d.Get("urpf_mode").(string)
	macAddress := d.Get("mac_address").(string)
	logicalRouterUpLinkPort := manager.LogicalRouterUpLinkPort{
		Description:               description,
		DisplayName:               displayName,
		Tags:                      tag,
		LogicalRouterId:           logicalRouterID,
		LinkedLogicalSwitchPortId: makeResourceReference("LogicalPort", linkedLogicalSwitchPortID),
		EdgeClusterMemberIndex:    edgeClusterMemberIndex,
		Subnets:                   subnets,
		UrpfMode:                  urpfMode,
		MacAddress:                macAddress,
	}
	logicalRouterUpLinkPort, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateLogicalRouterUpLinkPort(nsxClient.Context, logicalRouterUpLinkPort)

	if err != nil {
		return fmt.Errorf("Error during LogicalRouterUpLinkPort create: %v", err)
	}

	if resp != nil && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during LogicalRouterUpLinkPort create: %v", resp.StatusCode)
	}

	d.SetId(logicalRouterUpLinkPort.Id)

	return resourceNsxtLogicalRouterUpLinkPortRead(d, m)
}

func resourceNsxtLogicalRouterUpLinkPortRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router uplink port id")
	}

	logicalRouterUpLinkPort, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouterUpLinkPort(nsxClient.Context, id)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LogicalRouterUpLinkPort %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during LogicalRouterUpLinkPort read: %v", err)
	}

	d.Set("revision", logicalRouterUpLinkPort.Revision)
	d.Set("description", logicalRouterUpLinkPort.Description)
	d.Set("display_name", logicalRouterUpLinkPort.DisplayName)
	setTagsInSchema(d, logicalRouterUpLinkPort.Tags)
	d.Set("logical_router_id", logicalRouterUpLinkPort.LogicalRouterId)
	d.Set("linked_logical_switch_port_id", logicalRouterUpLinkPort.LinkedLogicalSwitchPortId.TargetId)
	d.Set("edge_cluster_member_index", logicalRouterUpLinkPort.EdgeClusterMemberIndex)
	d.Set("subnets", logicalRouterUpLinkPort.Subnets)
	d.Set("urpf_mode", logicalRouterUpLinkPort.UrpfMode)
	d.Set("mac_address", logicalRouterUpLinkPort.MacAddress)
	return nil
}

func resourceNsxtLogicalRouterUpLinkPortUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router uplink port id while updating")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	logicalRouterID := d.Get("logical_router_id").(string)
	linkedLogicalSwitchPortID := d.Get("linked_logical_switch_port_id").(string)
	edgeClusterMemberIndex := getInt64ListFromSchemaList(d, "edge_cluster_member_index")
	subnets := getIPSubnetsFromSchema(d)
	urpfMode := d.Get("urpf_mode").(string)
	macAddress := d.Get("mac_address").(string)
	logicalRouterUpLinkPort := manager.LogicalRouterUpLinkPort{
		Revision:                  revision,
		Description:               description,
		DisplayName:               displayName,
		Tags:                      tags,
		LogicalRouterId:           logicalRouterID,
		LinkedLogicalSwitchPortId: makeResourceReference("LogicalPort", linkedLogicalSwitchPortID),
		EdgeClusterMemberIndex:    edgeClusterMemberIndex,
		Subnets:                   subnets,
		UrpfMode:                  urpfMode,
		MacAddress:                macAddress,
		ResourceType:              LogicalRouterUpLinkPortResourceType,
	}

	logicalRouterUpLinkPort, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateLogicalRouterUpLinkPort(nsxClient.Context, id, logicalRouterUpLinkPort)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LogicalRouterUpLinkPort update: %v", err)
	}

	return resourceNsxtLogicalRouterUpLinkPortRead(d, m)
}

func resourceNsxtLogicalRouterUpLinkPortDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(nsxtClients).NsxtClient
	if nsxClient == nil {
		return resourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical router uplink port id while deleting")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteLogicalRouterPort(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during LogicalRouterUpLinkPort delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LogicalRouterUpLinkPort %s not found", id)
		d.SetId("")
	}

	return nil
}
