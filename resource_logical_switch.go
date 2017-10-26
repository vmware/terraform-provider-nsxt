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
			"Revision": GetRevisionSchema(),
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
			"Tags":            GetTagsSchema(),
			"AddressBindings": GetAddressBindingsSchema(),
			"AdminState": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Represents Desired state of the Logical Switch",
				Optional:    true,
			},
			"IpPoolId": &schema.Schema{
				Type:        schema.TypeString,
				Description: "IP pool id that associated with a LogicalSwitch",
				Optional:    true,
			},
			"MacPoolId": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Mac pool id that associated with a LogicalSwitch",
				Optional:    true,
			},
			"ReplicationMode": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Replication mode of the Logical Switch",
				Required:    true,
			},
			"SwitchingProfileIds": GetSwitchingProfileIdsSchema(),
			"TransportZoneId": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Id of the TransportZone to which this LogicalSwitch is associated",
				Required:    true,
				ForceNew:    true,
			},
			"Vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"Vni": &schema.Schema{
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

	description := d.Get("Description").(string)
	display_name := d.Get("DisplayName").(string)
	tags := GetTagsFromSchema(d)
	address_bindings := GetAddressBindingsFromSchema(d)
	admin_state := d.Get("AdminState").(string)
	ip_pool_id := d.Get("IpPoolId").(string)
	mac_pool_id := d.Get("MacPoolId").(string)
	replication_mode := d.Get("ReplicationMode").(string)
	switching_profile_ids := GetSwitchingProfileIdsFromSchema(d)
	transport_zone_id := d.Get("TransportZoneId").(string)
	vlan := int64(d.Get("Vlan").(int))
	vni := int32(d.Get("Vni").(int))
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
		return fmt.Errorf("Error during LogicalSwitch read: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalSwitch not found")
		d.SetId("")
		return nil
	}
	d.Set("Revision", logical_switch.Revision)
	d.Set("SystemOwned", logical_switch.SystemOwned)
	d.Set("Description", logical_switch.Description)
	d.Set("DisplayName", logical_switch.DisplayName)
	SetTagsInSchema(d, logical_switch.Tags)
	SetAddressBindingsInSchema(d, logical_switch.AddressBindings)
	d.Set("AdminState", logical_switch.AdminState)
	d.Set("IpPoolId", logical_switch.IpPoolId)
	d.Set("MacPoolId", logical_switch.MacPoolId)
	d.Set("ReplicationMode", logical_switch.ReplicationMode)
	SetSwitchingProfileIdsInSchema(d, logical_switch.SwitchingProfileIds)
	d.Set("TransportZoneId", logical_switch.TransportZoneId)
	d.Set("Vlan", logical_switch.Vlan)
	d.Set("Vni", logical_switch.Vni)

	return nil
}

func resourceLogicalSwitchUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*nsxt.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	description := d.Get("Description").(string)
	display_name := d.Get("DisplayName").(string)
	tags := GetTagsFromSchema(d)
	address_bindings := GetAddressBindingsFromSchema(d)
	admin_state := d.Get("AdminState").(string)
	ip_pool_id := d.Get("IpPoolId").(string)
	mac_pool_id := d.Get("MacPoolId").(string)
	replication_mode := d.Get("ReplicationMode").(string)
	switching_profile_ids := GetSwitchingProfileIdsFromSchema(d)
	transport_zone_id := d.Get("TransportZoneId").(string)
	vlan := int64(d.Get("Vlan").(int))
	vni := int32(d.Get("Vni").(int))
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
