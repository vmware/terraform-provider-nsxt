/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/loadbalancer"
	"log"
	"net/http"
)

func resourceNsxtLbFastUdpApplicationProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtLbFastUdpApplicationProfileCreate,
		Read:   resourceNsxtLbFastUdpApplicationProfileRead,
		Update: resourceNsxtLbFastUdpApplicationProfileUpdate,
		Delete: resourceNsxtLbFastUdpApplicationProfileDelete,
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
			"idle_timeout": &schema.Schema{
				Type:         schema.TypeInt,
				Description:  "Timeout in seconds to specify how long an idle UDP connection in ESTABLISHED state should be kept for this application before cleaning up",
				Optional:     true,
				Default:      300,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"ha_flow_mirroring": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "A boolean flag which reflects whether flow mirroring is enabled, and all the flows to the bounded virtual server are mirrored to the standby node",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceNsxtLbFastUdpApplicationProfileCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	haFlowMirroringEnabled := d.Get("ha_flow_mirroring").(bool)
	idleTimeout := int64(d.Get("idle_timeout").(int))
	lbFastUdpProfile := loadbalancer.LbFastUdpProfile{
		Description:          description,
		DisplayName:          displayName,
		Tags:                 tags,
		IdleTimeout:          idleTimeout,
		FlowMirroringEnabled: haFlowMirroringEnabled,
	}

	lbFastUdpProfile, resp, err := nsxClient.ServicesApi.CreateLoadBalancerFastUdpProfile(nsxClient.Context, lbFastUdpProfile)

	if err != nil {
		return fmt.Errorf("Error during LbFastUdpProfile create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status returned during LbFastUdpProfile create: %v", resp.StatusCode)
	}
	d.SetId(lbFastUdpProfile.Id)

	return resourceNsxtLbFastUdpApplicationProfileRead(d, m)
}

func resourceNsxtLbFastUdpApplicationProfileRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	lbFastUdpProfile, resp, err := nsxClient.ServicesApi.ReadLoadBalancerFastUdpProfile(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LbFastUdpProfile read: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LbFastUdpProfile %s not found", id)
		d.SetId("")
		return nil
	}
	d.Set("revision", lbFastUdpProfile.Revision)
	d.Set("description", lbFastUdpProfile.Description)
	d.Set("display_name", lbFastUdpProfile.DisplayName)
	setTagsInSchema(d, lbFastUdpProfile.Tags)
	d.Set("ha_flow_mirroring", lbFastUdpProfile.FlowMirroringEnabled)
	d.Set("idle_timeout", lbFastUdpProfile.IdleTimeout)

	return nil
}

func resourceNsxtLbFastUdpApplicationProfileUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	revision := int32(d.Get("revision").(int))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	haFlowMirroringEnabled := d.Get("ha_flow_mirroring").(bool)
	idleTimeout := int64(d.Get("idle_timeout").(int))
	lbFastUdpProfile := loadbalancer.LbFastUdpProfile{
		Revision:             revision,
		Description:          description,
		DisplayName:          displayName,
		Tags:                 tags,
		IdleTimeout:          idleTimeout,
		FlowMirroringEnabled: haFlowMirroringEnabled,
	}

	lbFastUdpProfile, resp, err := nsxClient.ServicesApi.UpdateLoadBalancerFastUdpProfile(nsxClient.Context, id, lbFastUdpProfile)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LbFastUdpProfile update: %v", err)
	}

	return resourceNsxtLbFastUdpApplicationProfileRead(d, m)
}

func resourceNsxtLbFastUdpApplicationProfileDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	resp, err := nsxClient.ServicesApi.DeleteLoadBalancerApplicationProfile(nsxClient.Context, id)
	if err != nil {
		return fmt.Errorf("Error during LbFastUdpProfile delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LbFastUdpProfile %s not found", id)
		d.SetId("")
	}
	return nil
}
