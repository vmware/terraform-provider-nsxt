/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func resourceNsxtPolicySpoofGuardProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicySpoofGuardProfileCreate,
		Read:   resourceNsxtPolicySpoofGuardProfileRead,
		Update: resourceNsxtPolicySpoofGuardProfileUpdate,
		Delete: resourceNsxtPolicySpoofGuardProfileDelete,
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
			"address_binding_allowlist": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceNsxtPolicySpoofGuardProfileExists(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
	var err error
	if isGlobalManager {
		client := gm_infra.NewSpoofguardProfilesClient(connector)
		_, err = client.Get(id)
	} else {
		client := infra.NewSpoofguardProfilesClient(connector)
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

func resourceNsxtPolicySpoofGuardProfilePatch(d *schema.ResourceData, m interface{}, id string) error {
	connector := getPolicyConnector(m)

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	addressBindingAllowlist := d.Get("address_binding_allowlist").(bool)

	obj := model.SpoofGuardProfile{
		DisplayName:             &displayName,
		Description:             &description,
		Tags:                    tags,
		AddressBindingAllowlist: &addressBindingAllowlist,
	}

	log.Printf("[INFO] Patching SpoofGuardProfile with ID %s", id)
	if isPolicyGlobalManager(m) {
		gmObj, convErr := convertModelBindingType(obj, model.SpoofGuardProfileBindingType(), gm_model.SpoofGuardProfileBindingType())
		if convErr != nil {
			return convErr
		}
		client := gm_infra.NewSpoofguardProfilesClient(connector)
		return client.Patch(id, gmObj.(gm_model.SpoofGuardProfile), nil)
	}

	client := infra.NewSpoofguardProfilesClient(connector)
	return client.Patch(id, obj, nil)
}

func resourceNsxtPolicySpoofGuardProfileCreate(d *schema.ResourceData, m interface{}) error {

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicySpoofGuardProfileExists)
	if err != nil {
		return err
	}

	err = resourceNsxtPolicySpoofGuardProfilePatch(d, m, id)
	if err != nil {
		return handleCreateError("SpoofGuardProfile", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicySpoofGuardProfileRead(d, m)
}

func resourceNsxtPolicySpoofGuardProfileRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining SpoofGuardProfile ID")
	}

	var obj model.SpoofGuardProfile
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewSpoofguardProfilesClient(connector)
		gmObj, err := client.Get(id)
		if err != nil {
			return handleReadError(d, "SpoofGuardProfile", id, err)
		}

		lmObj, err := convertModelBindingType(gmObj, gm_model.SpoofGuardProfileBindingType(), model.SpoofGuardProfileBindingType())
		if err != nil {
			return err
		}
		obj = lmObj.(model.SpoofGuardProfile)
	} else {
		client := infra.NewSpoofguardProfilesClient(connector)
		var err error
		obj, err = client.Get(id)
		if err != nil {
			return handleReadError(d, "SpoofGuardProfile", id, err)
		}
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	d.Set("address_binding_allowlist", obj.AddressBindingAllowlist)

	return nil
}

func resourceNsxtPolicySpoofGuardProfileUpdate(d *schema.ResourceData, m interface{}) error {

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining SpoofGuardProfile ID")
	}

	err := resourceNsxtPolicySpoofGuardProfilePatch(d, m, id)
	if err != nil {
		return handleUpdateError("SpoofGuardProfile", id, err)
	}

	return resourceNsxtPolicySpoofGuardProfileRead(d, m)
}

func resourceNsxtPolicySpoofGuardProfileDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining SpoofGuardProfile ID")
	}

	connector := getPolicyConnector(m)
	var err error
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewSpoofguardProfilesClient(connector)
		err = client.Delete(id, nil)
	} else {
		client := infra.NewSpoofguardProfilesClient(connector)
		err = client.Delete(id, nil)
	}

	if err != nil {
		return handleDeleteError("SpoofGuardProfile", id, err)
	}

	return nil
}
