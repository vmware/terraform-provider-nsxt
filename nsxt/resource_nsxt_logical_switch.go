/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"log"
	"net/http"
	"time"
)

var logicalSwitchReplicationModeValues = []string{"MTEP", "SOURCE", ""}

// TODO: consider splitting this resource to overlay_ls and vlan_ls
func resourceNsxtLogicalSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtLogicalSwitchCreate,
		Read:   resourceNsxtLogicalSwitchRead,
		Update: resourceNsxtLogicalSwitchUpdate,
		Delete: resourceNsxtLogicalSwitchDelete,

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
			"tag":             getTagsSchema(),
			"address_binding": getAddressBindingsSchema(),
			"admin_state":     getAdminStateSchema(),
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
			"switching_profile_id": getSwitchingProfileIdsSchema(),
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

func resourceNsxtLogicalSwitchCreateRollback(nsxClient *api.APIClient, id string) {
	log.Printf("[ERROR] Rollback switch %s creation due to unrealized state", id)

	localVarOptionals := make(map[string]interface{})
	_, err := nsxClient.LogicalSwitchingApi.DeleteLogicalSwitch(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		log.Printf("[ERROR] Rollback failed!")
	}
}

func resourceNsxtLogicalSwitchCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	address_bindings := getAddressBindingsFromSchema(d)
	admin_state := d.Get("admin_state").(string)
	ip_pool_id := d.Get("ip_pool_id").(string)
	mac_pool_id := d.Get("mac_pool_id").(string)
	replication_mode := d.Get("replication_mode").(string)
	switching_profile_id := getSwitchingProfileIdsFromSchema(d)
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
		SwitchingProfileIds: switching_profile_id,
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
		stateConf := &resource.StateChangeConf{
			Pending: []string{"in_progress", "pending", "partial_success"},
			Target:  []string{"success"},
			Refresh: func() (interface{}, string, error) {
				state, resp, err := nsxClient.LogicalSwitchingApi.GetLogicalSwitchState(nsxClient.Context, logical_switch.Id)
				if err != nil {
					return nil, "", fmt.Errorf("Error while querying realization state: %v", err)
				}

				if resp.StatusCode != http.StatusOK {
					return nil, "", fmt.Errorf("Unexpected return status %d", resp.StatusCode)
				}

				if state.FailureCode != 0 {
					return nil, "", fmt.Errorf("Error in switch realization: %s", state.FailureMessage)
				}

				log.Printf("[DEBUG] Realization state: %s", state.State)
				return logical_switch, state.State, nil
			},
			Timeout:    d.Timeout(schema.TimeoutCreate),
			MinTimeout: 1 * time.Second,
			Delay:      1 * time.Second,
		}
		_, err = stateConf.WaitForState()
		if err != nil {
			resourceNsxtLogicalSwitchCreateRollback(nsxClient, logical_switch.Id)
			return err
		}
	}

	d.SetId(logical_switch.Id)

	return resourceNsxtLogicalSwitchRead(d, m)
}

func resourceNsxtLogicalSwitchRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical switch id")
	}

	logical_switch, resp, err := nsxClient.LogicalSwitchingApi.GetLogicalSwitch(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LogicalSwitch %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during LogicalSwitch read: %v", err)
	}

	d.Set("revision", logical_switch.Revision)
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

func resourceNsxtLogicalSwitchUpdate(d *schema.ResourceData, m interface{}) error {
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
	switching_profile_id := getSwitchingProfileIdsFromSchema(d)
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
		SwitchingProfileIds: switching_profile_id,
		TransportZoneId:     transport_zone_id,
		Vlan:                vlan,
		Vni:                 vni,
		Revision:            revision,
	}

	logical_switch, resp, err := nsxClient.LogicalSwitchingApi.UpdateLogicalSwitch(nsxClient.Context, id, logical_switch)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LogicalSwitch update: %v", err)
	}

	return resourceNsxtLogicalSwitchRead(d, m)
}

func resourceNsxtLogicalSwitchDelete(d *schema.ResourceData, m interface{}) error {
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
		log.Printf("[DEBUG] LogicalSwitch %s not found", id)
		d.SetId("")
	}
	return nil
}
