/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
)

func resourceNsxtPolicyIPBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIPBlockCreate,
		Read:   resourceNsxtPolicyIPBlockRead,
		Update: resourceNsxtPolicyIPBlockUpdate,
		Delete: resourceNsxtPolicyIPBlockDelete,
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
			"cidr": {
				Type:         schema.TypeString,
				Description:  "Network address and the prefix length which will be associated with a layer-2 broadcast domain",
				Required:     true,
				ValidateFunc: validateCidr(),
			},
		},
	}
}

func resourceNsxtPolicyIPBlockExists(id string, connector *client.RestConnector) bool {
	client := infra.NewDefaultIpBlocksClient(connector)

	_, err := client.Get(id)
	if err == nil {
		return true
	}

	if isNotFoundError(err) {
		return false
	}

	logAPIError("Error retrieving IP Block", err)
	return false
}

func resourceNsxtPolicyIPBlockRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultIpBlocksClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IP Block ID")
	}

	block, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "IP Block", id, err)
	}

	d.Set("display_name", block.DisplayName)
	d.Set("description", block.Description)
	setPolicyTagsInSchema(d, block.Tags)
	d.Set("nsx_id", block.Id)
	d.Set("path", block.Path)
	d.Set("revision", block.Revision)
	d.Set("cidr", block.Cidr)

	return nil
}

func resourceNsxtPolicyIPBlockCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultIpBlocksClient(connector)

	id, err := getOrGenerateID(d, connector, resourceNsxtPolicyIPBlockExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	cidr := d.Get("cidr").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.IpAddressBlock{
		DisplayName: &displayName,
		Description: &description,
		Cidr:        &cidr,
		Tags:        tags,
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating IP Block with ID %s", id)
	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("IP Block", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyIPBlockRead(d, m)
}

func resourceNsxtPolicyIPBlockUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultIpBlocksClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IP Block ID")
	}

	// Read the rest of the configured parameters
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	cidr := d.Get("cidr").(string)
	revision := int64(d.Get("revision").(int))
	tags := getPolicyTagsFromSchema(d)

	obj := model.IpAddressBlock{
		Id:          &id,
		DisplayName: &displayName,
		Description: &description,
		Cidr:        &cidr,
		Tags:        tags,
		Revision:    &revision,
	}

	_, err := client.Update(id, obj)
	if err != nil {
		return handleUpdateError("IP Block", id, err)
	}
	return resourceNsxtPolicyIPBlockRead(d, m)

}

func resourceNsxtPolicyIPBlockDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IP Block ID")
	}

	connector := getPolicyConnector(m)
	client := infra.NewDefaultIpBlocksClient(connector)
	err := client.Delete(id)
	if err != nil {
		return handleDeleteError("IP Block", id, err)
	}

	return nil

}
