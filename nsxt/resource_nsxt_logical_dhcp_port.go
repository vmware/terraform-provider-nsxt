/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"log"
	"net/http"
)

var dhcpType = "DHCP_SERVICE"

func resourceNsxtLogicalDhcpPort() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtLogicalDhcpPortCreate,
		Read:   resourceNsxtLogicalDhcpPortRead,
		Update: resourceNsxtLogicalDhcpPortUpdate,
		Delete: resourceNsxtLogicalDhcpPortDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"logical_switch_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Id of the Logical switch that this port belongs to",
				Required:    true,
				ForceNew:    true, // Cannot change the logical switch of a logical port
			},
			"dhcp_server_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Id of the Logical DHCP server this port belongs to",
				Required:    true,
			},
			"admin_state": getAdminStateSchema(),
			//TODO(asarfaty): switching profiles for DHCP ports are currently ignored by NSX
			//"switching_profile_id": getSwitchingProfileIdsSchema(),
			"tag": getTagsSchema(),
		},
	}
}

func resourceNsxtLogicalDhcpPortCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	name := d.Get("display_name").(string)
	description := d.Get("description").(string)
	lsID := d.Get("logical_switch_id").(string)
	adminState := d.Get("admin_state").(string)
	//profilesList := getSwitchingProfileIdsFromSchema(d)
	tagList := getTagsFromSchema(d)
	dhcpServerId := d.Get("dhcp_server_id").(string)
	attachment := manager.LogicalPortAttachment{
		AttachmentType: dhcpType,
		Id:             dhcpServerId,
	}
	lp := manager.LogicalPort{
		DisplayName:     name,
		Description:     description,
		LogicalSwitchId: lsID,
		AdminState:      adminState,
		//SwitchingProfileIds: profilesList,
		Tags:       tagList,
		Attachment: &attachment,
	}

	lp, resp, err := nsxClient.LogicalSwitchingApi.CreateLogicalPort(nsxClient.Context, lp)

	if err != nil {
		return fmt.Errorf("Error while creating logical DHCP port %s: %v", name, err)
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during Logical DHCP port create: %v", resp.StatusCode)
	}

	resourceID := lp.Id
	d.SetId(resourceID)

	return resourceNsxtLogicalDhcpPortRead(d, m)
}

func resourceNsxtLogicalDhcpPortRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical DHCP port ID from state during read")
	}
	LogicalDhcpPort, resp, err := nsxClient.LogicalSwitchingApi.GetLogicalPort(nsxClient.Context, id)

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		log.Printf("[DEBUG] Logical DHCP port %s not found", id)
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error while reading logical DHCP port %s: %v", id, err)
	}
	if LogicalDhcpPort.Attachment == nil || LogicalDhcpPort.Attachment.AttachmentType != dhcpType {
		return fmt.Errorf("Error reading DHCP port %s: This not a DHCP port", id)
	}

	d.Set("revision", LogicalDhcpPort.Revision)
	d.Set("display_name", LogicalDhcpPort.DisplayName)
	d.Set("description", LogicalDhcpPort.Description)
	d.Set("logical_switch_id", LogicalDhcpPort.LogicalSwitchId)
	d.Set("admin_state", LogicalDhcpPort.AdminState)
	d.Set("dhcp_server_id", LogicalDhcpPort.Attachment.Id)
	//err = setSwitchingProfileIdsInSchema(d, nsxClient, LogicalDhcpPort.SwitchingProfileIds)
	//if err != nil {
	//	return fmt.Errorf("Error during logical DHCP port switching profiles set in schema: %v", err)
	//}
	setTagsInSchema(d, LogicalDhcpPort.Tags)

	return nil
}

func resourceNsxtLogicalDhcpPortUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	name := d.Get("display_name").(string)
	description := d.Get("description").(string)
	adminState := d.Get("admin_state").(string)
	lsID := d.Get("logical_switch_id").(string)
	//profilesList := getSwitchingProfileIdsFromSchema(d)
	tagList := getTagsFromSchema(d)
	revision := int64(d.Get("revision").(int))
	dhcpServerId := d.Get("dhcp_server_id").(string)
	attachment := manager.LogicalPortAttachment{
		AttachmentType: "DHCP_SERVICE",
		Id:             dhcpServerId,
	}
	lp := manager.LogicalPort{
		Revision:        revision,
		DisplayName:     name,
		Description:     description,
		LogicalSwitchId: lsID,
		AdminState:      adminState,
		//SwitchingProfileIds: profilesList,
		Tags:       tagList,
		Attachment: &attachment,
	}

	lp, resp, err := nsxClient.LogicalSwitchingApi.UpdateLogicalPort(nsxClient.Context, id, lp)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error while updating logical DHCP port %s: %v", id, err)
	}
	return resourceNsxtLogicalDhcpPortRead(d, m)
}

func resourceNsxtLogicalDhcpPortDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	lpID := d.Id()
	if lpID == "" {
		return fmt.Errorf("Error obtaining logical DHCP port ID from state during delete")
	}
	localVarOptionals := make(map[string]interface{})
	localVarOptionals["detach"] = true
	resp, err := nsxClient.LogicalSwitchingApi.DeleteLogicalPort(nsxClient.Context, lpID, localVarOptionals)

	if err != nil {
		return fmt.Errorf("Error while deleting logical DHCP port %s: %v", lpID, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] Logical DHCP port %s was not found", lpID)
		d.SetId("")
	}

	return nil
}
