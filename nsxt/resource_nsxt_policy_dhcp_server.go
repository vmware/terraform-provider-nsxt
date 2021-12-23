/* Copyright © 2019 VMware, Inc. All Rights Reserved.
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

func resourceNsxtPolicyDhcpServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyDhcpServerCreate,
		Read:   resourceNsxtPolicyDhcpServerRead,
		Update: resourceNsxtPolicyDhcpServerUpdate,
		Delete: resourceNsxtPolicyDhcpServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":            getNsxIDSchema(),
			"path":              getPathSchema(),
			"display_name":      getDisplayNameSchema(),
			"description":       getDescriptionSchema(),
			"revision":          getRevisionSchema(),
			"tag":               getTagsSchema(),
			"edge_cluster_path": getPolicyPathSchema(false, false, "Edge Cluster path"),
			"lease_time": {
				Type:         schema.TypeInt,
				Description:  "IP Address lease time in seconds",
				Optional:     true,
				Default:      86400,
				ValidateFunc: validation.IntAtLeast(60),
			},
			"preferred_edge_paths": {
				Type:        schema.TypeList,
				Description: "The first edge node is assigned as active edge, and second one as standby edge",
				Optional:    true,
				MaxItems:    2,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validatePolicyPath(),
				},
			},
			"server_addresses": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "DHCP server address in CIDR format",
				MaxItems:    2,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateIPCidr(),
				},
			},
		},
	}
}

func resourceNsxtPolicyDhcpServerExists(id string, connector *client.RestConnector, isGlobalManager bool) (bool, error) {

	var err error
	if isGlobalManager {
		client := gm_infra.NewDhcpServerConfigsClient(connector)
		_, err = client.Get(id)
	} else {
		client := infra.NewDhcpServerConfigsClient(connector)
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

func resourceNsxtPolicyDhcpServerSchemaToModel(d *schema.ResourceData) model.DhcpServerConfig {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	edgeClusterPath := d.Get("edge_cluster_path").(string)
	leaseTime := int64(d.Get("lease_time").(int))
	preferredEdgePaths := interface2StringList(d.Get("preferred_edge_paths").([]interface{}))
	serverAddresses := interface2StringList(d.Get("server_addresses").([]interface{}))

	obj := model.DhcpServerConfig{
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		LeaseTime:   &leaseTime,
	}

	if edgeClusterPath != "" {
		obj.EdgeClusterPath = &edgeClusterPath
	}
	if len(preferredEdgePaths) > 0 {
		obj.PreferredEdgePaths = preferredEdgePaths
	}
	if len(serverAddresses) > 0 {
		obj.ServerAddresses = serverAddresses
		// TODO: remove setting of server_address once this deprecated property is removed from NSX API
		obj.ServerAddress = &serverAddresses[0]
	}

	return obj
}

func resourceNsxtPolicyDhcpServerCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicyDhcpServerExists)
	if err != nil {
		return err
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating DhcpServer with ID %s", id)
	if isPolicyGlobalManager(m) {
		obj := resourceNsxtPolicyDhcpServerSchemaToModel(d)
		gmObj, err1 := convertModelBindingType(obj, model.DhcpServerConfigBindingType(), gm_model.DhcpServerConfigBindingType())
		if err1 != nil {
			return err1
		}

		client := gm_infra.NewDhcpServerConfigsClient(connector)
		err = client.Patch(id, gmObj.(gm_model.DhcpServerConfig))
	} else {
		client := infra.NewDhcpServerConfigsClient(connector)
		err = client.Patch(id, resourceNsxtPolicyDhcpServerSchemaToModel(d))
	}
	if err != nil {
		return handleCreateError("DhcpServer", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyDhcpServerRead(d, m)
}

func resourceNsxtPolicyDhcpServerRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpServer ID")
	}

	var obj model.DhcpServerConfig
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewDhcpServerConfigsClient(connector)
		gmObj, err := client.Get(id)
		if err != nil {
			return handleReadError(d, "DhcpServer", id, err)
		}
		rawObj, err := convertModelBindingType(gmObj, gm_model.DhcpServerConfigBindingType(), model.DhcpServerConfigBindingType())
		if err != nil {
			return err
		}
		obj = rawObj.(model.DhcpServerConfig)
	} else {
		var err error
		client := infra.NewDhcpServerConfigsClient(connector)
		obj, err = client.Get(id)
		if err != nil {
			return handleReadError(d, "DhcpServer", id, err)
		}
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)

	d.Set("edge_cluster_path", obj.EdgeClusterPath)
	d.Set("lease_time", obj.LeaseTime)
	d.Set("preferred_edge_paths", obj.PreferredEdgePaths)
	d.Set("server_addresses", obj.ServerAddresses)

	return nil
}

func resourceNsxtPolicyDhcpServerUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDhcpServerConfigsClient(connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpServer ID")
	}

	// Update the resource using PATCH
	var err error
	if isPolicyGlobalManager(m) {
		obj := resourceNsxtPolicyDhcpServerSchemaToModel(d)
		gmObj, err1 := convertModelBindingType(obj, model.DhcpServerConfigBindingType(), gm_model.DhcpServerConfigBindingType())
		if err1 != nil {
			return err1
		}

		client := gm_infra.NewDhcpServerConfigsClient(connector)
		err = client.Patch(id, gmObj.(gm_model.DhcpServerConfig))
	} else {
		client := infra.NewDhcpServerConfigsClient(connector)
		err = client.Patch(id, resourceNsxtPolicyDhcpServerSchemaToModel(d))
	}
	if err != nil {
		return handleUpdateError("DhcpServer", id, err)
	}

	return resourceNsxtPolicyDhcpServerRead(d, m)
}

func resourceNsxtPolicyDhcpServerDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining DhcpServer ID")
	}

	var err error
	connector := getPolicyConnector(m)
	if isPolicyGlobalManager(m) {
		client := gm_infra.NewDhcpServerConfigsClient(connector)
		err = client.Delete(id)
	} else {
		client := infra.NewDhcpServerConfigsClient(connector)
		err = client.Delete(id)
	}

	if err != nil {
		return handleDeleteError("DhcpServer", id, err)
	}

	return nil
}
