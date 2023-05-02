/* Copyright © 2019 VMware, Inc. All Rights Reserved.
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

func resourceNsxtPolicyDhcpRelayConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyDhcpRelayConfigCreate,
		Read:   resourceNsxtPolicyDhcpRelayConfigRead,
		Update: resourceNsxtPolicyDhcpRelayConfigUpdate,
		Delete: resourceNsxtPolicyDhcpRelayConfigDelete,
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
			"server_addresses": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateSingleIP(),
				},
			},
		},
	}
}

func resourceNsxtPolicyDhcpRelayConfigExists(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
	var err error
	if isGlobalManager {
		client := gm_infra.NewDhcpRelayConfigsClient(connector)
		_, err = client.Get(id)
	} else {
		client := infra.NewDhcpRelayConfigsClient(connector)
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

func resourceNsxtPolicyDhcpRelayConfigCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicyDhcpRelayConfigExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	serverAddresses := getStringListFromSchemaList(d, "server_addresses")

	obj := model.DhcpRelayConfig{
		DisplayName:     &displayName,
		Description:     &description,
		Tags:            tags,
		ServerAddresses: serverAddresses,
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating DhcpRelayConfig with ID %s", id)
	if isPolicyGlobalManager(m) {
		gmObj, err1 := convertModelBindingType(obj, model.DhcpRelayConfigBindingType(), gm_model.DhcpRelayConfigBindingType())
		if err1 != nil {
			return err1
		}

		client := gm_infra.NewDhcpRelayConfigsClient(connector)
		err = client.Patch(id, gmObj.(gm_model.DhcpRelayConfig))
	} else {
		client := infra.NewDhcpRelayConfigsClient(connector)
		err = client.Patch(id, obj)
	}
	if err != nil {
		return handleCreateError("DhcpRelayConfig", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyDhcpRelayConfigRead(d, m)
}

func resourceNsxtPolicyDhcpRelayConfigRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDhcpRelayConfigsClient(connector)

	if client == nil {
		return policyResourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpRelayConfig ID")
	}
	var obj model.DhcpRelayConfig
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewDhcpRelayConfigsClient(connector)
		gmObj, err1 := client.Get(id)
		if err1 != nil {
			return handleReadError(d, "DhcpRelayConfig", id, err1)
		}
		rawObj, err := convertModelBindingType(gmObj, gm_model.DhcpRelayConfigBindingType(), model.DhcpRelayConfigBindingType())
		if err != nil {
			return err
		}
		obj = rawObj.(model.DhcpRelayConfig)
	} else {
		var err error
		client := infra.NewDhcpRelayConfigsClient(connector)
		obj, err = client.Get(id)
		if err != nil {
			return handleReadError(d, "DhcpRelayConfig", id, err)
		}
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	d.Set("server_addresses", obj.ServerAddresses)

	return nil
}

func resourceNsxtPolicyDhcpRelayConfigUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDhcpRelayConfigsClient(connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpRelayConfig ID")
	}

	// Read the rest of the configured parameters
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)
	revision := int64(d.Get("revision").(int))

	serverAddresses := getStringListFromSchemaList(d, "server_addresses")

	obj := model.DhcpRelayConfig{
		DisplayName:     &displayName,
		Description:     &description,
		Tags:            tags,
		ServerAddresses: serverAddresses,
		Revision:        &revision,
	}

	var err error
	if isPolicyGlobalManager(m) {
		gmObj, err1 := convertModelBindingType(obj, model.DhcpRelayConfigBindingType(), gm_model.DhcpRelayConfigBindingType())
		if err1 != nil {
			return err1
		}

		client := gm_infra.NewDhcpRelayConfigsClient(connector)
		_, err = client.Update(id, gmObj.(gm_model.DhcpRelayConfig))
	} else {
		client := infra.NewDhcpRelayConfigsClient(connector)
		_, err = client.Update(id, obj)
	}
	if err != nil {
		return handleUpdateError("DhcpRelayConfig", id, err)
	}

	return resourceNsxtPolicyDhcpRelayConfigRead(d, m)
}

func resourceNsxtPolicyDhcpRelayConfigDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpRelayConfig ID")
	}

	connector := getPolicyConnector(m)
	client := infra.NewDhcpRelayConfigsClient(connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	var err error
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewDhcpRelayConfigsClient(connector)
		err = client.Delete(id)
	} else {
		client := infra.NewDhcpRelayConfigsClient(connector)
		err = client.Delete(id)
	}

	if err != nil {
		return handleDeleteError("DhcpRelayConfig", id, err)
	}

	return nil
}
