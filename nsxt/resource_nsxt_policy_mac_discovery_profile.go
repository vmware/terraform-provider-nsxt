/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var macDiscoveryProfileMacLimitPolicyValues = []string{
	model.MacDiscoveryProfile_MAC_LIMIT_POLICY_ALLOW,
	model.MacDiscoveryProfile_MAC_LIMIT_POLICY_DROP,
}

func resourceNsxtPolicyMacDiscoveryProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyMacDiscoveryProfileCreate,
		Read:   resourceNsxtPolicyMacDiscoveryProfileRead,
		Update: resourceNsxtPolicyMacDiscoveryProfileUpdate,
		Delete: resourceNsxtPolicyMacDiscoveryProfileDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtPolicyPathResourceImporter,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"context":      getContextSchema(false, false),
			"mac_change_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"mac_learning_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"mac_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 4096),
				Default:      4096,
			},
			"mac_limit_policy": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(macDiscoveryProfileMacLimitPolicyValues, false),
				Optional:     true,
				Default:      "ALLOW",
			},
			"remote_overlay_mac_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2048, 8192),
				Default:      2048,
			},
			"unknown_unicast_flooding_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceNsxtPolicyMacDiscoveryProfileExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	var err error
	client := infra.NewMacDiscoveryProfilesClient(sessionContext, connector)
	_, err = client.Get(id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyMacDiscoveryProfileCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID2(d, m, resourceNsxtPolicyMacDiscoveryProfileExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	macChangeEnabled := d.Get("mac_change_enabled").(bool)
	macLearningEnabled := d.Get("mac_learning_enabled").(bool)
	macLimit := int64(d.Get("mac_limit").(int))
	macLimitPolicy := d.Get("mac_limit_policy").(string)
	remoteOverlayMacLimit := int64(d.Get("remote_overlay_mac_limit").(int))
	unknownUnicastFloodingEnabled := d.Get("unknown_unicast_flooding_enabled").(bool)

	obj := model.MacDiscoveryProfile{
		DisplayName:                   &displayName,
		Description:                   &description,
		Tags:                          tags,
		MacChangeEnabled:              &macChangeEnabled,
		MacLearningEnabled:            &macLearningEnabled,
		MacLimit:                      &macLimit,
		MacLimitPolicy:                &macLimitPolicy,
		RemoteOverlayMacLimit:         &remoteOverlayMacLimit,
		UnknownUnicastFloodingEnabled: &unknownUnicastFloodingEnabled,
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating MacDiscoveryProfile with ID %s", id)
	boolFalse := false
	client := infra.NewMacDiscoveryProfilesClient(getSessionContext(d, m), connector)
	err = client.Patch(id, obj, &boolFalse)
	if err != nil {
		return handleCreateError("MacDiscoveryProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyMacDiscoveryProfileRead(d, m)
}

func resourceNsxtPolicyMacDiscoveryProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining MacDiscoveryProfile ID")
	}

	client := infra.NewMacDiscoveryProfilesClient(getSessionContext(d, m), connector)
	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "MacDiscoveryProfile", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	d.Set("mac_change_enabled", obj.MacChangeEnabled)
	d.Set("mac_learning_enabled", obj.MacLearningEnabled)
	d.Set("mac_limit", obj.MacLimit)
	d.Set("mac_limit_policy", obj.MacLimitPolicy)
	d.Set("remote_overlay_mac_limit", obj.RemoteOverlayMacLimit)
	d.Set("unknown_unicast_flooding_enabled", obj.UnknownUnicastFloodingEnabled)

	return nil
}

func resourceNsxtPolicyMacDiscoveryProfileUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining MacDiscoveryProfile ID")
	}

	// Read the rest of the configured parameters
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)

	macChangeEnabled := d.Get("mac_change_enabled").(bool)
	macLearningEnabled := d.Get("mac_learning_enabled").(bool)
	macLimit := int64(d.Get("mac_limit").(int))
	macLimitPolicy := d.Get("mac_limit_policy").(string)
	remoteOverlayMacLimit := int64(d.Get("remote_overlay_mac_limit").(int))
	unknownUnicastFloodingEnabled := d.Get("unknown_unicast_flooding_enabled").(bool)
	revision := int64(d.Get("revision").(int))

	obj := model.MacDiscoveryProfile{
		DisplayName:                   &displayName,
		Description:                   &description,
		Tags:                          tags,
		MacChangeEnabled:              &macChangeEnabled,
		MacLearningEnabled:            &macLearningEnabled,
		MacLimit:                      &macLimit,
		MacLimitPolicy:                &macLimitPolicy,
		RemoteOverlayMacLimit:         &remoteOverlayMacLimit,
		UnknownUnicastFloodingEnabled: &unknownUnicastFloodingEnabled,
		Revision:                      &revision,
	}

	// Update the resource using PATCH
	boolFalse := false
	client := infra.NewMacDiscoveryProfilesClient(getSessionContext(d, m), connector)
	_, err := client.Update(id, obj, &boolFalse)
	if err != nil {
		return handleUpdateError("MacDiscoveryProfile", id, err)
	}

	return resourceNsxtPolicyMacDiscoveryProfileRead(d, m)
}

func resourceNsxtPolicyMacDiscoveryProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining MacDiscoveryProfile ID")
	}

	connector := getPolicyConnector(m)
	boolFalse := false
	client := infra.NewMacDiscoveryProfilesClient(getSessionContext(d, m), connector)
	err := client.Delete(id, &boolFalse)

	if err != nil {
		return handleDeleteError("MacDiscoveryProfile", id, err)
	}

	return nil
}
