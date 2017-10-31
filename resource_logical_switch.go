package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func resourceLogicalSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogicalSwitchCreate,
		Read:   resourceLogicalSwitchRead,
		Update: resourceLogicalSwitchUpdate,
		Delete: resourceLogicalSwitchDelete,

		Schema: map[string]*schema.Schema{
			"revision": GetRevisionSchema(),
			"system_owned": GetSystemOwnedSchema(),
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
			"tags":             GetTagsSchema(),
			"address_bindings": GetAddressBindingsSchema(),
			"admin_state": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Represents Desired state of the Logical Switch",
				Optional:    true,
			},
			"ip_pool_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "IP pool id that associated with a LogicalSwitch",
				Optional:    true,
			},
			"mac_pool_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Mac pool id that associated with a LogicalSwitch",
				Optional:    true,
			},
			"replication_mode": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Replication mode of the Logical Switch",
				Required:    true,
			},
			"switching_profile_ids": GetSwitchingProfileIdsSchema(),
			"transport_zone_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Id of the TransportZone to which this LogicalSwitch is associated",
				Required:    true,
				ForceNew:    true,
			},
			"vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vni": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "VNI for this LogicalSwitch",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceLogicalSwitchCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*nsxt.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := GetTagsFromSchema(d)
	address_bindings := GetAddressBindingsFromSchema(d)
	admin_state := d.Get("admin_state").(string)
	ip_pool_id := d.Get("ip_pool_id").(string)
	mac_pool_id := d.Get("mac_pool_id").(string)
	replication_mode := d.Get("replication_mode").(string)
	switching_profile_ids := GetSwitchingProfileIdsFromSchema(d)
	transport_zone_id := d.Get("transport_zone_id").(string)
	vlan := int64(d.Get("vlan").(int))
	vni := int32(d.Get("vni").(int))
	logical_switch := manager.LogicalSwitch{
		Description:         description,
		DisplayName:         display_name,
		Tags:                tags,
		AddressBindings:     address_bindings,
		AdminState:          admin_state,
		IpPoolId:            ip_pool_id,
		MacPoolId:           mac_pool_id,
		ReplicationMode:     replication_mode,
		SwitchingProfileIds: switching_profile_ids,
		TransportZoneId:     transport_zone_id,
		Vlan:                vlan,
		Vni:                 vni,
	}

	logical_switch, resp, err := nsxClient.LogicalSwitchingApi.CreateLogicalSwitch(nsxClient.Context, logical_switch)

	if err != nil {
		return fmt.Errorf("Error during LogicalSwitch create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Unexpected status returned")
		return nil
	}
	d.SetId(logical_switch.Id)

	return resourceLogicalSwitchRead(d, m)
}

func resourceLogicalSwitchRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*nsxt.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	logical_switch, resp, err := nsxClient.LogicalSwitchingApi.GetLogicalSwitch(nsxClient.Context, id)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			fmt.Printf("LogicalSwitch not found")
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error during LogicalSwitch read: %v", err)
	}

	d.Set("revision", logical_switch.Revision)
	d.Set("system_owned", logical_switch.SystemOwned)
	d.Set("description", logical_switch.Description)
	d.Set("display_name", logical_switch.DisplayName)
	SetTagsInSchema(d, logical_switch.Tags)
	SetAddressBindingsInSchema(d, logical_switch.AddressBindings)
	d.Set("admin_state", logical_switch.AdminState)
	d.Set("ip_pool_id", logical_switch.IpPoolId)
	d.Set("mac_pool_id", logical_switch.MacPoolId)
	d.Set("replication_mode", logical_switch.ReplicationMode)
	SetSwitchingProfileIdsInSchema(d, nsxClient, logical_switch.SwitchingProfileIds)
	d.Set("transport_zone_id", logical_switch.TransportZoneId)
	d.Set("vlan", logical_switch.Vlan)
	d.Set("vni", logical_switch.Vni)

	return nil
}

func resourceLogicalSwitchUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*nsxt.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := GetTagsFromSchema(d)
	address_bindings := GetAddressBindingsFromSchema(d)
	admin_state := d.Get("admin_state").(string)
	ip_pool_id := d.Get("ip_pool_id").(string)
	mac_pool_id := d.Get("mac_pool_id").(string)
	replication_mode := d.Get("replication_mode").(string)
	switching_profile_ids := GetSwitchingProfileIdsFromSchema(d)
	transport_zone_id := d.Get("transport_zone_id").(string)
	vlan := int64(d.Get("vlan").(int))
	vni := int32(d.Get("vni").(int))
	revision := int64(d.Get("revision").(int))
	logical_switch := manager.LogicalSwitch{
		Description:         description,
		DisplayName:         display_name,
		Tags:                tags,
		AddressBindings:     address_bindings,
		AdminState:          admin_state,
		IpPoolId:            ip_pool_id,
		MacPoolId:           mac_pool_id,
		ReplicationMode:     replication_mode,
		SwitchingProfileIds: switching_profile_ids,
		TransportZoneId:     transport_zone_id,
		Vlan:                vlan,
		Vni:                 vni,
		Revision:            revision,
	}

	logical_switch, resp, err := nsxClient.LogicalSwitchingApi.UpdateLogicalSwitch(nsxClient.Context, id, logical_switch)

	if err != nil {
		return fmt.Errorf("Error during LogicalSwitch update: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalSwitch not found")
		d.SetId("")
		return nil
	}
	return resourceLogicalSwitchRead(d, m)
}

func resourceLogicalSwitchDelete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*nsxt.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.LogicalSwitchingApi.DeleteLogicalSwitch(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during LogicalSwitch delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalSwitch not found")
		d.SetId("")
	}
	return nil
}
