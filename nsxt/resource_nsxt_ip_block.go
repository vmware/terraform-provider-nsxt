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

func resourceNsxtIpBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtIpBlockCreate,
		Read:   resourceNsxtIpBlockRead,
		Update: resourceNsxtIpBlockUpdate,
		Delete: resourceNsxtIpBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			"cidr": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Represents network address and the prefix length which will be associated with a layer-2 broadcast domain",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceNsxtIpBlockCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	cidr := d.Get("cidr").(string)
	ipBlock := manager.IpBlock{
		Description: description,
		DisplayName: displayName,
		Tags:        tags,
		Cidr:        cidr,
	}
	// Create the IP Block
	ipBlock, resp, err := nsxClient.PoolManagementApi.CreateIpBlock(nsxClient.Context, ipBlock)

	if err != nil {
		return fmt.Errorf("Error during IpBlock create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned during IpBlock create: %v", resp.StatusCode)
	}
	d.SetId(ipBlock.Id)

	return resourceNsxtIpBlockRead(d, m)
}

func resourceNsxtIpBlockRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	ipBlock, resp, err := nsxClient.PoolManagementApi.ReadIpBlock(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] IpBlock %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during IpBlock read: %v", err)
	}

	d.Set("revision", ipBlock.Revision)
	d.Set("description", ipBlock.Description)
	d.Set("display_name", ipBlock.DisplayName)
	setTagsInSchema(d, ipBlock.Tags)
	d.Set("cidr", ipBlock.Cidr)

	return nil
}

func resourceNsxtIpBlockUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	cidr := d.Get("cidr").(string)
	tags := getTagsFromSchema(d)
	revision := int64(d.Get("revision").(int))
	ipBlock := manager.IpBlock{
		DisplayName: displayName,
		Description: description,
		Cidr:        cidr,
		Tags:        tags,
		Revision:    revision,
	}
	// Update the IP block
	ipBlock, resp, err := nsxClient.PoolManagementApi.UpdateIpBlock(nsxClient.Context, id, ipBlock)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during IpBlock update: %v", err)
	}

	return resourceNsxtIpBlockRead(d, m)
}

func resourceNsxtIpBlockDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}
	resp, err := nsxClient.PoolManagementApi.DeleteIpBlock(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during IpBlock delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] IpBlock %s not found", id)
		d.SetId("")
	}
	return nil
}
