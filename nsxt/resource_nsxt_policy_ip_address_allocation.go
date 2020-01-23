/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/ip_pools"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
	"strings"
)

func resourceNsxtPolicyIPAddressAllocation() *schema.Resource {
	displayNameSchema := getDisplayNameSchema()
	displayNameSchema.ForceNew = true
	descriptionSchema := getDescriptionSchema()
	descriptionSchema.ForceNew = true
	tagSchema := getTagsSchema()
	tagSchema.ForceNew = true

	return &schema.Resource{
		Create: resourceNsxtPolicyIPAddressAllocationCreate,
		Read:   resourceNsxtPolicyIPAddressAllocationRead,
		Delete: resourceNsxtPolicyIPAddressAllocationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyIPAddressAllocationImport,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": displayNameSchema,
			"description":  descriptionSchema,
			"revision":     getRevisionSchema(),
			"tag":          tagSchema,
			"pool_path":    getPolicyPathSchema(true, true, "The path of the IP Pool for this allocation"),
			"allocation_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The IP allocated. If unspecified any free IP will be allocated.",
				ValidateFunc: validateSingleIP(),
				Computed:     true,
				ForceNew:     true,
			},
		},
	}
}

func resourceNsxtPolicyIPAddressParsePoolIDFromPath(poolPath string, connector *client.RestConnector) (string, error) {
	poolID := getPolicyIDFromPath(poolPath)
	if !resourceNsxtPolicyIPPoolExists(poolID, connector) {
		return "", fmt.Errorf("IP Pool specified in path '%s' does not exist", poolPath)
	}
	return poolID, nil
}

func resourceNsxtPolicyIPAddressAllocationExists(poolID string, allocationID string, connector *client.RestConnector) bool {
	client := ip_pools.NewDefaultIpAllocationsClient(connector)

	_, err := client.Get(poolID, allocationID)
	if err == nil {
		return true
	}

	if isNotFoundError(err) {
		return false
	}

	logAPIError("Error retrieving resource", err)

	return false
}

func resourceNsxtPolicyIPAddressAllocationCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := ip_pools.NewDefaultIpAllocationsClient(connector)

	if client == nil {
		return policyResourceNotSupportedError()
	}

	poolID, err := resourceNsxtPolicyIPAddressParsePoolIDFromPath(d.Get("pool_path").(string), connector)

	if err != nil {
		return err
	}

	id := d.Get("nsx_id").(string)
	if id == "" {
		uuid, _ := uuid.NewRandom()
		id = uuid.String()
	}

	if resourceNsxtPolicyIPAddressAllocationExists(poolID, id, connector) {
		return fmt.Errorf("Resource with ID %s already exists", id)
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	allocationIP := d.Get("allocation_ip").(string)

	obj := model.IpAddressAllocation{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
	}

	if allocationIP != "" {
		obj.AllocationIp = &allocationIP
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating IPAddressAllocation with ID %s", id)
	err = client.Patch(poolID, id, obj)
	if err != nil {
		return handleCreateError("IPAddressAllocation", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyIPAddressAllocationRead(d, m)
}

func resourceNsxtPolicyIPAddressAllocationRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := ip_pools.NewDefaultIpAllocationsClient(connector)

	if client == nil {
		return policyResourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPAddressAllocation ID")
	}

	poolID, err := resourceNsxtPolicyIPAddressParsePoolIDFromPath(d.Get("pool_path").(string), connector)

	if err != nil {
		return err
	}

	obj, err := client.Get(poolID, id)
	if err != nil {
		return handleReadError(d, "IPAddressAllocation", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("pool_path", obj.ParentPath)

	d.Set("allocation_ip", obj.AllocationIp)

	if d.Get("allocation_ip").(string) == "" {
		log.Printf("[DEBUG] Waiting for realization of IP Address for IP Allocation with ID %s", id)

		stateConf := nsxtPolicyWaitForRealizationStateConf(connector, d, d.Get("path").(string))
		entity, err := stateConf.WaitForState()
		if err != nil {
			return err
		}
		realizedResource := entity.(model.GenericPolicyRealizedResource)
		for _, attr := range realizedResource.ExtendedAttributes {
			if *attr.Key == "allocation_ip" {
				d.Set("allocation_ip", attr.Values[0])
				return nil
			}
		}
		return fmt.Errorf("Failed to get realized IP for path %s", d.Get("path"))
	}

	return nil
}

func resourceNsxtPolicyIPAddressAllocationDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := ip_pools.NewDefaultIpAllocationsClient(connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPAddressAllocation ID")
	}

	poolID, err := resourceNsxtPolicyIPAddressParsePoolIDFromPath(d.Get("pool_path").(string), connector)

	if err != nil {
		return err
	}

	err = client.Delete(poolID, id)
	if err != nil {
		return handleDeleteError("IPAddressAllocation", id, err)
	}

	return nil
}

func resourceNsxtPolicyIPAddressAllocationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	s := strings.Split(importID, "/")
	if len(s) != 2 {
		return nil, fmt.Errorf("Please provide <ip-pool-id>/<allocation-id> as an input")
	}

	poolID := s[0]
	connector := getPolicyConnector(m)
	client := infra.NewDefaultIpPoolsClient(connector)

	pool, err := client.Get(poolID)
	if err != nil {
		return nil, err
	}
	d.Set("pool_path", pool.Path)

	d.SetId(s[1])

	return []*schema.ResourceData{d}, nil
}
