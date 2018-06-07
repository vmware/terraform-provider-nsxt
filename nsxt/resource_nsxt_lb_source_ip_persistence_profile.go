/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/loadbalancer"
	"log"
	"net/http"
)

func resourceNsxtLbSourceIpPersistenceProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtLbSourceIpPersistenceProfileCreate,
		Read:   resourceNsxtLbSourceIpPersistenceProfileRead,
		Update: resourceNsxtLbSourceIpPersistenceProfileUpdate,
		Delete: resourceNsxtLbSourceIpPersistenceProfileDelete,
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
			"persistence_shared": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "A boolean flag which reflects whether the cookie persistence is private or shared",
				Optional:    true,
			},
			"ha_persistence_mirroring": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "A boolean flag which reflects whether persistence entries will be synchronized to the HA peer",
				Optional:    true,
				Default:     false,
			},
			"purge_when_full": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "A boolean flag which reflects whether entries will be purged when the persistence table is full",
				Optional:    true,
				Default:     true,
			},
			"timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Persistence expiration time in seconds, counted from the time all the connections are completed",
				Optional:    true,
				Default:     300,
			},
		},
	}
}

func resourceNsxtLbSourceIpPersistenceProfileCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	persistenceShared := d.Get("persistence_shared").(bool)
	haPersistenceMirroring := d.Get("ha_persistence_mirroring").(bool)
	purge_flag := d.Get("purge_when_full").(bool)
	purge := "FULL"
	if purge_flag == false {
		purge = "NO_PURGE"
	}
	timeout := int64(d.Get("timeout").(int))
	lbSourceIpPersistenceProfile := loadbalancer.LbSourceIpPersistenceProfile{
		Description:                   description,
		DisplayName:                   displayName,
		Tags:                          tags,
		PersistenceShared:             persistenceShared,
		HaPersistenceMirroringEnabled: haPersistenceMirroring,
		Purge:   purge,
		Timeout: timeout,
	}

	lbSourceIpPersistenceProfile, resp, err := nsxClient.ServicesApi.CreateLoadBalancerSourceIpPersistenceProfile(nsxClient.Context, lbSourceIpPersistenceProfile)

	if err != nil {
		return fmt.Errorf("Error during LbSourceIpPersistenceProfile create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status returned during LbSourceIpPersistenceProfile create: %v", resp.StatusCode)
	}
	d.SetId(lbSourceIpPersistenceProfile.Id)

	return resourceNsxtLbSourceIpPersistenceProfileRead(d, m)
}

func resourceNsxtLbSourceIpPersistenceProfileRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	lbSourceIpPersistenceProfile, resp, err := nsxClient.ServicesApi.ReadLoadBalancerSourceIpPersistenceProfile(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LbSourceIpPersistenceProfile read: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LbSourceIpPersistenceProfile %s not found", id)
		d.SetId("")
		return nil
	}
	d.Set("revision", lbSourceIpPersistenceProfile.Revision)
	d.Set("description", lbSourceIpPersistenceProfile.Description)
	d.Set("display_name", lbSourceIpPersistenceProfile.DisplayName)
	setTagsInSchema(d, lbSourceIpPersistenceProfile.Tags)
	d.Set("persistence_shared", lbSourceIpPersistenceProfile.PersistenceShared)
	d.Set("ha_persistence_mirroring", lbSourceIpPersistenceProfile.HaPersistenceMirroringEnabled)
	if lbSourceIpPersistenceProfile.Purge == "FULL" {
		d.Set("purge_when_full", true)
	} else {
		d.Set("purge_when_full", false)
	}
	d.Set("timeout", lbSourceIpPersistenceProfile.Timeout)

	return nil
}

func resourceNsxtLbSourceIpPersistenceProfileUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	revision := int32(d.Get("revision").(int))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	persistenceShared := d.Get("persistence_shared").(bool)
	haPersistenceMirroring := d.Get("ha_persistence_mirroring").(bool)
	purge_flag := d.Get("purge_when_full").(bool)
	purge := "FULL"
	if purge_flag == false {
		purge = "NO_PURGE"
	}
	timeout := int64(d.Get("timeout").(int))
	lbSourceIpPersistenceProfile := loadbalancer.LbSourceIpPersistenceProfile{
		Revision:                      revision,
		Description:                   description,
		DisplayName:                   displayName,
		Tags:                          tags,
		PersistenceShared:             persistenceShared,
		HaPersistenceMirroringEnabled: haPersistenceMirroring,
		Purge:   purge,
		Timeout: timeout,
	}

	lbSourceIpPersistenceProfile, resp, err := nsxClient.ServicesApi.UpdateLoadBalancerSourceIpPersistenceProfile(nsxClient.Context, id, lbSourceIpPersistenceProfile)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LbSourceIpPersistenceProfile update: %v", err)
	}

	return resourceNsxtLbSourceIpPersistenceProfileRead(d, m)
}

func resourceNsxtLbSourceIpPersistenceProfileDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	resp, err := nsxClient.ServicesApi.DeleteLoadBalancerPersistenceProfile(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LbSourceIpPersistenceProfile delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LbSourceIpPersistenceProfile %s not found", id)
		d.SetId("")
	}
	return nil
}
