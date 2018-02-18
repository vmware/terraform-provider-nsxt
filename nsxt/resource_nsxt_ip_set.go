/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

func resourceIpSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpSetCreate,
		Read:   resourceIpSetRead,
		Update: resourceIpSetUpdate,
		Delete: resourceIpSetDelete,

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
			"tag": getTagsSchema(),
			"ip_addresses": &schema.Schema{
				Type:        schema.TypeSet,
				Description: "Set of IP addresses",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCidrOrIPOrRange(),
				},
				Optional: true,
			},
		},
	}
}

func resourceIpSetCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	ip_addresses := getStringListFromSchemaSet(d, "ip_addresses")
	ip_set := manager.IpSet{
		Description: description,
		DisplayName: display_name,
		Tags:        tags,
		IpAddresses: ip_addresses,
	}

	ip_set, resp, err := nsxClient.GroupingObjectsApi.CreateIPSet(nsxClient.Context, ip_set)

	if err != nil {
		return fmt.Errorf("Error during IpSet create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during IpSet create: %v", resp.StatusCode)
	}
	d.SetId(ip_set.Id)

	return resourceIpSetRead(d, m)
}

func resourceIpSetRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	ip_set, resp, err := nsxClient.GroupingObjectsApi.ReadIPSet(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("IpSet %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during IpSet read: %v", err)
	}

	d.Set("revision", ip_set.Revision)
	d.Set("description", ip_set.Description)
	d.Set("display_name", ip_set.DisplayName)
	setTagsInSchema(d, ip_set.Tags)
	d.Set("ip_addresses", ip_set.IpAddresses)

	return nil
}

func resourceIpSetUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	ip_addresses := interface2StringList(d.Get("ip_addresses").(*schema.Set).List())
	ip_set := manager.IpSet{
		Revision:    revision,
		Description: description,
		DisplayName: display_name,
		Tags:        tags,
		IpAddresses: ip_addresses,
	}

	ip_set, resp, err := nsxClient.GroupingObjectsApi.UpdateIPSet(nsxClient.Context, id, ip_set)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during IpSet update: %v", err)
	}

	return resourceIpSetRead(d, m)
}

func resourceIpSetDelete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.GroupingObjectsApi.DeleteIPSet(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during IpSet delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("IpSet %s not found", id)
		d.SetId("")
	}
	return nil
}
