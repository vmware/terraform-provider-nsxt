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
	"time"
)

var logicalSwitchReplicationModeValues = []string{"MTEP", "SOURCE"}

// TODO: consider splitting this resource to overlay_ls and vlan_ls
func resourceLogicalSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogicalSwitchCreate,
		Read:   resourceLogicalSwitchRead,
		Update: resourceLogicalSwitchUpdate,
		Delete: resourceLogicalSwitchDelete,

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
			"tags":             getTagsSchema(),
			"address_bindings": getAddressBindingsSchema(),
			"admin_state": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Represents Desired state of the Logical Switch",
				Required:    true,
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
				Type:         schema.TypeString,
				Description:  "Replication mode of the Logical Switch",
				Optional:     true,
				Default:      "MTEP",
				ValidateFunc: validation.StringInSlice(logicalSwitchReplicationModeValues, false),
			},
			"switching_profile_ids": getSwitchingProfileIdsSchema(),
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
			"verify_realization": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Wait for realization to complete",
				Default:     true,
				Optional:    true,
			},
		},
	}
}

func resourceLogicalSwitchCreateRollback(nsxClient *api.APIClient, id string) {
	log.Printf("[ERROR] Rollback switch %d creation due to unrealized state", id)

	localVarOptionals := make(map[string]interface{})
	_, err := nsxClient.LogicalSwitchingApi.DeleteLogicalSwitch(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		log.Printf("[ERROR] Rollback failed!")
	}
}

func resourceLogicalSwitchCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	address_bindings := getAddressBindingsFromSchema(d)
	admin_state := d.Get("admin_state").(string)
	ip_pool_id := d.Get("ip_pool_id").(string)
	mac_pool_id := d.Get("mac_pool_id").(string)
	replication_mode := d.Get("replication_mode").(string)
	switching_profile_ids := getSwitchingProfileIdsFromSchema(d)
	transport_zone_id := d.Get("transport_zone_id").(string)
	vlan := int64(d.Get("vlan").(int))
	vni := int32(d.Get("vni").(int))

	verify_realization := d.Get("verify_realization").(bool)

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
		return fmt.Errorf("Unexpected status returned during LogicalSwitch create: %v", resp.StatusCode)
	}

	if verify_realization {

		for {
			time.Sleep(time.Second)

			state, resp, err := nsxClient.LogicalSwitchingApi.GetLogicalSwitchState(nsxClient.Context, logical_switch.Id)

			if err != nil {
				resourceLogicalSwitchCreateRollback(nsxClient, logical_switch.Id)
				return fmt.Errorf("Error while querying realization state: %v", err)
			}

			if resp.StatusCode != http.StatusOK {
				resourceLogicalSwitchCreateRollback(nsxClient, logical_switch.Id)
				return fmt.Errorf("Unexpected return status %d", resp.StatusCode)
			}

			if state.FailureCode != 0 {
				resourceLogicalSwitchCreateRollback(nsxClient, logical_switch.Id)
				return fmt.Errorf("Error in switch realization: %s", state.FailureMessage)
			}

			log.Printf("[DEBUG] Realization state: %s", state.State)

			if state.State == "success" {
				break
			}

		}
	}

	d.SetId(logical_switch.Id)

	return resourceLogicalSwitchRead(d, m)
}

func resourceLogicalSwitchRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical switch id")
	}

	logical_switch, resp, err := nsxClient.LogicalSwitchingApi.GetLogicalSwitch(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalSwitch %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during LogicalSwitch read: %v", err)
	}

	d.Set("revision", logical_switch.Revision)
	d.Set("system_owned", logical_switch.SystemOwned)
	d.Set("description", logical_switch.Description)
	d.Set("display_name", logical_switch.DisplayName)
	setTagsInSchema(d, logical_switch.Tags)
	setAddressBindingsInSchema(d, logical_switch.AddressBindings)
	d.Set("admin_state", logical_switch.AdminState)
	d.Set("ip_pool_id", logical_switch.IpPoolId)
	d.Set("mac_pool_id", logical_switch.MacPoolId)
	d.Set("replication_mode", logical_switch.ReplicationMode)
	setSwitchingProfileIdsInSchema(d, nsxClient, logical_switch.SwitchingProfileIds)
	d.Set("transport_zone_id", logical_switch.TransportZoneId)
	d.Set("vlan", logical_switch.Vlan)
	d.Set("vni", logical_switch.Vni)

	return nil
}

func resourceLogicalSwitchUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical switch id")
	}

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	address_bindings := getAddressBindingsFromSchema(d)
	admin_state := d.Get("admin_state").(string)
	ip_pool_id := d.Get("ip_pool_id").(string)
	mac_pool_id := d.Get("mac_pool_id").(string)
	replication_mode := d.Get("replication_mode").(string)
	switching_profile_ids := getSwitchingProfileIdsFromSchema(d)
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
	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LogicalSwitch update: %v", err)
	}

	return resourceLogicalSwitchRead(d, m)
}

func resourceLogicalSwitchDelete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical switch id")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.LogicalSwitchingApi.DeleteLogicalSwitch(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during LogicalSwitch delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalSwitch %s not found", id)
		d.SetId("")
	}
	return nil
}
