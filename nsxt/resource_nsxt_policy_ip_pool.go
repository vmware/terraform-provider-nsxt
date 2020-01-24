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

func resourceNsxtPolicyIPPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIPPoolCreate,
		Read:   resourceNsxtPolicyIPPoolRead,
		Update: resourceNsxtPolicyIPPoolUpdate,
		Delete: resourceNsxtPolicyIPPoolDelete,
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
		},
	}
}

func resourceNsxtPolicyIPPoolExists(id string, connector *client.RestConnector) bool {
	client := infra.NewDefaultIpPoolsClient(connector)

	_, err := client.Get(id)
	if err == nil {
		return true
	}

	if isNotFoundError(err) {
		return false
	}

	logAPIError("Error retrieving IP Pool", err)
	return false
}

func resourceNsxtPolicyIPPoolRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultIpPoolsClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IP Pool ID")
	}

	pool, err := client.Get(id)
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			log.Printf("[DEBUG] IP Pool %s not found", id)
			return nil
		}
		return handleReadError(d, "IP Pool", id, err)
	}

	d.Set("display_name", pool.DisplayName)
	d.Set("description", pool.Description)
	setPolicyTagsInSchema(d, pool.Tags)
	d.Set("nsx_id", pool.Id)
	d.Set("path", pool.Path)
	d.Set("revision", pool.Revision)

	return nil
}

func resourceNsxtPolicyIPPoolCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultIpPoolsClient(connector)

	id, err := getOrGenerateID(d, connector, resourceNsxtPolicyIPPoolExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.IpAddressPool{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Id:          &id,
	}

	log.Printf("[INFO] Creating IP Pool with ID %s", id)
	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("IP Pool", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyIPPoolRead(d, m)
}

func resourceNsxtPolicyIPPoolUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultIpPoolsClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IP Pool ID")
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)

	obj := model.IpAddressPool{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		Id:          &id,
	}

	log.Printf("[INFO] Updating IP Pool with ID %s", id)
	err := client.Patch(id, obj)
	if err != nil {
		return handleUpdateError("IP Pool", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyIPPoolRead(d, m)
}

func resourceNsxtPolicyIPPoolDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultIpPoolsClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IP Pool ID")
	}

	log.Printf("[INFO] Deleting IP Pool with ID %s", id)
	err := client.Delete(id)
	if err != nil {
		return handleDeleteError("IP Pool", id, err)
	}

	return nil
}
