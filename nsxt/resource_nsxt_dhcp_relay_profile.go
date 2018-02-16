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

func resourceDhcpRelayProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceDhcpRelayProfileCreate,
		Read:   resourceDhcpRelayProfileRead,
		Update: resourceDhcpRelayProfileUpdate,
		Delete: resourceDhcpRelayProfileDelete,

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
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
			"server_addresses": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "Set of dhcp relay server addresses",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateSingleIP(),
				},
				Required: true,
			},
		},
	}
}

func resourceDhcpRelayProfileCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	server_addresses := getStringListFromSchemaSet(d, "server_addresses")
	dhcp_relay_profile := manager.DhcpRelayProfile{
		Description:     description,
		DisplayName:     display_name,
		Tags:            tags,
		ServerAddresses: server_addresses,
	}

	dhcp_relay_profile, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateDhcpRelayProfile(nsxClient.Context, dhcp_relay_profile)

	if err != nil {
		return fmt.Errorf("Error during DhcpRelayProfile create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during DhcpRelayProfile create: %v", resp.StatusCode)
	}
	d.SetId(dhcp_relay_profile.Id)

	return resourceDhcpRelayProfileRead(d, m)
}

func resourceDhcpRelayProfileRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining dhcp relay profile id")
	}

	dhcp_relay_profile, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadDhcpRelayProfile(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("DhcpRelayProfile %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during DhcpRelayProfile read: %v", err)
	}

	d.Set("revision", dhcp_relay_profile.Revision)
	d.Set("description", dhcp_relay_profile.Description)
	d.Set("display_name", dhcp_relay_profile.DisplayName)
	setTagsInSchema(d, dhcp_relay_profile.Tags)
	d.Set("server_addresses", dhcp_relay_profile.ServerAddresses)

	return nil
}

func resourceDhcpRelayProfileUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining dhcp relay profile id")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	server_addresses := interface2StringList(d.Get("server_addresses").(*schema.Set).List())
	dhcp_relay_profile := manager.DhcpRelayProfile{
		Revision:        revision,
		Description:     description,
		DisplayName:     display_name,
		Tags:            tags,
		ServerAddresses: server_addresses,
	}

	dhcp_relay_profile, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateDhcpRelayProfile(nsxClient.Context, id, dhcp_relay_profile)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during DhcpRelayProfile update: %v", err)
	}

	return resourceDhcpRelayProfileRead(d, m)
}

func resourceDhcpRelayProfileDelete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining dhcp relay profile id")
	}

	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteDhcpRelayProfile(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during DhcpRelayProfile delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("DhcpRelayProfile %s not found", id)
		d.SetId("")
	}
	return nil
}
