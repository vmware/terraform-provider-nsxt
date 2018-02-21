/* Copyright © 2018 VMware, Inc. All Rights Reserved.
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

func resourceNsxtIpSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtIpSetCreate,
		Read:   resourceNsxtIpSetRead,
		Update: resourceNsxtIpSetUpdate,
		Delete: resourceNsxtIpSetDelete,

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

func resourceNsxtIpSetCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	ipAddresses := getStringListFromSchemaSet(d, "ip_addresses")
	ipSet := manager.IpSet{
		Description: description,
		DisplayName: displayName,
		Tags:        tags,
		IpAddresses: ipAddresses,
	}

	ipSet, resp, err := nsxClient.GroupingObjectsApi.CreateIPSet(nsxClient.Context, ipSet)

	if err != nil {
		return fmt.Errorf("Error during IpSet create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during IpSet create: %v", resp.StatusCode)
	}
	d.SetId(ipSet.Id)

	return resourceNsxtIpSetRead(d, m)
}

func resourceNsxtIpSetRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	ipSet, resp, err := nsxClient.GroupingObjectsApi.ReadIPSet(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] IpSet %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during IpSet read: %v", err)
	}

	d.Set("revision", ipSet.Revision)
	d.Set("description", ipSet.Description)
	d.Set("display_name", ipSet.DisplayName)
	setTagsInSchema(d, ipSet.Tags)
	d.Set("ip_addresses", ipSet.IpAddresses)

	return nil
}

func resourceNsxtIpSetUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	ipAddresses := interface2StringList(d.Get("ip_addresses").(*schema.Set).List())
	ipSet := manager.IpSet{
		Revision:    revision,
		Description: description,
		DisplayName: displayName,
		Tags:        tags,
		IpAddresses: ipAddresses,
	}

	ipSet, resp, err := nsxClient.GroupingObjectsApi.UpdateIPSet(nsxClient.Context, id, ipSet)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during IpSet update: %v", err)
	}

	return resourceNsxtIpSetRead(d, m)
}

func resourceNsxtIpSetDelete(d *schema.ResourceData, m interface{}) error {
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
		log.Printf("[DEBUG] IpSet %s not found", id)
		d.SetId("")
	}
	return nil
}
