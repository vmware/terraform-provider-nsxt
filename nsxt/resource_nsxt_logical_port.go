/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func resourceLogicalPort() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogicalPortCreate,
		Read:   resourceLogicalPortRead,
		Update: resourceLogicalPortUpdate,
		Delete: resourceLogicalPortDelete,

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"logical_switch_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // Cannot change the logical switch of a logical port
			},
			"admin_state":          getAdminStateSchema(),
			"switching_profile_id": getSwitchingProfileIdsSchema(),
			"tags":                 getTagsSchema(),
		},
	}
}

func resourceLogicalPortCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)

	name := d.Get("display_name").(string)
	description := d.Get("description").(string)
	ls_id := d.Get("logical_switch_id").(string)
	admin_state := d.Get("admin_state").(string)
	profilesList := getSwitchingProfileIdsFromSchema(d)
	tagList := getTagsFromSchema(d)

	lp := manager.LogicalPort{
		DisplayName:         name,
		Description:         description,
		LogicalSwitchId:     ls_id,
		AdminState:          admin_state,
		SwitchingProfileIds: profilesList,
		Tags:                tagList}

	lp, resp, err := nsxClient.LogicalSwitchingApi.CreateLogicalPort(nsxClient.Context, lp)

	if err != nil {
		return fmt.Errorf("Error while creating logical port %s: %v\n", lp.DisplayName, err)
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during Logical port create: %v", resp.StatusCode)
	}

	resource_id := lp.Id
	d.SetId(resource_id)

	return resourceLogicalPortRead(d, m)
}

func resourceLogicalPortRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical port ID from state during read")
	}
	logical_port, resp, err := nsxClient.LogicalSwitchingApi.GetLogicalPort(nsxClient.Context, id)

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		fmt.Printf("Logical port %s not found", id)
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error while reading logical port %s: %v\n", id, err)
	}

	d.Set("revision", logical_port.Revision)
	d.Set("display_name", logical_port.DisplayName)
	d.Set("description", logical_port.Description)
	d.Set("logical_switch_id", logical_port.LogicalSwitchId)
	d.Set("admin_state", logical_port.AdminState)
	setSwitchingProfileIdsInSchema(d, nsxClient, logical_port.SwitchingProfileIds)
	setTagsInSchema(d, logical_port.Tags)

	return nil
}

func resourceLogicalPortUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)

	id := d.Id()
	name := d.Get("display_name").(string)
	description := d.Get("description").(string)
	admin_state := d.Get("admin_state").(string)
	profilesList := getSwitchingProfileIdsFromSchema(d)
	tagList := getTagsFromSchema(d)
	revision := int64(d.Get("revision").(int))

	// Some of the port attributes (attachment) are not exposed to terraform.
	// If we try to update port based on terraform attributes only, apply will fail
	// due to missing info.
	// We don't expose attachment to terraform, since it will become out of sync
	// once attachment info is updated/remove outside the scope of port management.

	lp, resp, err := nsxClient.LogicalSwitchingApi.GetLogicalPort(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Logical port %s was not found\n", id)
	}
	if err != nil {
		return fmt.Errorf("Error while reading logical port %s: %v\n", id, err)
	}

	lp.DisplayName = name
	lp.Description = description
	lp.AdminState = admin_state
	lp.SwitchingProfileIds = profilesList
	lp.Tags = tagList
	lp.Revision = revision

	lp, resp, err = nsxClient.LogicalSwitchingApi.UpdateLogicalPort(nsxClient.Context, id, lp)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error while updating logical port %s: %v\n", id, err)
	}
	return resourceLogicalPortRead(d, m)
}

func resourceLogicalPortDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)

	lp_id := d.Id()
	if lp_id == "" {
		return fmt.Errorf("Error obtaining logical port ID from state during delete")
	}
	//TODO: add optional detach param
	localVarOptionals := make(map[string]interface{})

	resp, err := nsxClient.LogicalSwitchingApi.DeleteLogicalPort(nsxClient.Context, lp_id, localVarOptionals)

	if err != nil {
		return fmt.Errorf("Error while deleting logical port %s: %v\n", lp_id, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("Logical port %s was not found\n", lp_id)
		d.SetId("")
	}

	return nil
}
