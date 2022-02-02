/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
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
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
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

func resourceNsxtPolicyMacDiscoveryProfileExists(id string, connector *client.RestConnector, isGlobalManager bool) (bool, error) {
	var err error
	if isGlobalManager {
		client := gm_infra.NewMacDiscoveryProfilesClient(connector)
		_, err = client.Get(id)
	} else {
		client := infra.NewMacDiscoveryProfilesClient(connector)
		_, err = client.Get(id)
	}
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
	id, err := getOrGenerateID(d, m, resourceNsxtPolicyMacDiscoveryProfileExists)
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
	if isPolicyGlobalManager(m) {
		gmObj, convErr := convertModelBindingType(obj, model.MacDiscoveryProfileBindingType(), gm_model.MacDiscoveryProfileBindingType())
		if convErr != nil {
			return convErr
		}
		client := gm_infra.NewMacDiscoveryProfilesClient(connector)
		err = client.Patch(id, gmObj.(gm_model.MacDiscoveryProfile), &boolFalse)
	} else {
		client := infra.NewMacDiscoveryProfilesClient(connector)
		err = client.Patch(id, obj, &boolFalse)
	}
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

	var obj model.MacDiscoveryProfile
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewMacDiscoveryProfilesClient(connector)
		gmObj, err := client.Get(id)
		if err != nil {
			return handleReadError(d, "MacDiscoveryProfile", id, err)
		}

		lmObj, err := convertModelBindingType(gmObj, gm_model.MacDiscoveryProfileBindingType(), model.MacDiscoveryProfileBindingType())
		if err != nil {
			return err
		}
		obj = lmObj.(model.MacDiscoveryProfile)
	} else {
		client := infra.NewMacDiscoveryProfilesClient(connector)
		var err error
		obj, err = client.Get(id)
		if err != nil {
			return handleReadError(d, "MacDiscoveryProfile", id, err)
		}
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
	var err error
	boolFalse := false
	if isPolicyGlobalManager(m) {
		gmObj, convErr := convertModelBindingType(obj, model.MacDiscoveryProfileBindingType(), gm_model.MacDiscoveryProfileBindingType())
		if convErr != nil {
			return convErr
		}
		client := gm_infra.NewMacDiscoveryProfilesClient(connector)
		_, err = client.Update(id, gmObj.(gm_model.MacDiscoveryProfile), &boolFalse)
	} else {
		client := infra.NewMacDiscoveryProfilesClient(connector)
		_, err = client.Update(id, obj, &boolFalse)
	}
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
	var err error
	boolFalse := false
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewMacDiscoveryProfilesClient(connector)
		err = client.Delete(id, &boolFalse)
	} else {
		client := infra.NewMacDiscoveryProfilesClient(connector)
		err = client.Delete(id, &boolFalse)
	}

	if err != nil {
		return handleDeleteError("MacDiscoveryProfile", id, err)
	}

	return nil
}
