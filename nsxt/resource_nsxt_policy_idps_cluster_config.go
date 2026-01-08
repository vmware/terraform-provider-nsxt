// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func resourceNsxtPolicyIdpsClusterConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIdpsClusterConfigCreate,
		Read:   resourceNsxtPolicyIdpsClusterConfigRead,
		Update: resourceNsxtPolicyIdpsClusterConfigUpdate,
		Delete: resourceNsxtPolicyIdpsClusterConfigDelete,
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
			"context":      getContextSchema(false, false, false),
			"ids_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable or disable IDS for the cluster",
				Required:    true,
			},
			"cluster": {
				Type:        schema.TypeList,
				Description: "Cluster reference configuration",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_id": {
							Type:        schema.TypeString,
							Description: "Cluster target ID (e.g., domain-c123)",
							Required:    true,
							ForceNew:    true,
						},
						"target_type": {
							Type:        schema.TypeString,
							Description: "Target type (e.g., VC_Cluster)",
							Required:    true,
							ForceNew:    true,
						},
					},
				},
			},
		},
	}
}

func resourceNsxtPolicyIdpsClusterConfigExists(sessionContext utl.SessionContext, id string, connector client.Connector) (bool, error) {
	client := intrusion_services.NewClusterConfigsClient(connector)
	if client == nil {
		return false, policyResourceNotSupportedError()
	}

	_, err := client.Get(id, nil)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving IDPS Cluster Config", err)
}

func buildIdsClusterConfig(d *schema.ResourceData, id string) model.IdsClusterConfig {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	idsEnabled := d.Get("ids_enabled").(bool)

	clusterList := d.Get("cluster").([]interface{})
	clusterData := clusterList[0].(map[string]interface{})
	targetID := clusterData["target_id"].(string)
	targetType := clusterData["target_type"].(string)

	// Build PolicyResourceReference for cluster
	cluster := model.PolicyResourceReference{
		TargetId:   &targetID,
		TargetType: &targetType,
	}

	idsClusterConfig := model.IdsClusterConfig{
		Id:          &id,
		DisplayName: &displayName,
		Description: &description,
		Tags:        tags,
		IdsEnabled:  &idsEnabled,
		Cluster:     &cluster,
	}

	return idsClusterConfig
}

func resourceNsxtPolicyIdpsClusterConfigCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := intrusion_services.NewClusterConfigsClient(connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	// Use the cluster target ID as the resource ID for NSX-T IDPS cluster config
	clusterList := d.Get("cluster").([]interface{})
	clusterData := clusterList[0].(map[string]interface{})
	targetID := clusterData["target_id"].(string)

	// The ID for IDPS cluster config should be the cluster target ID
	id := targetID

	idsClusterConfig := buildIdsClusterConfig(d, id)

	log.Printf("[INFO] Creating IDPS Cluster Config with ID %s for cluster %s (ids_enabled: %v)", id, *idsClusterConfig.Cluster.TargetId, *idsClusterConfig.IdsEnabled)
	err := client.Patch(id, idsClusterConfig)
	if err != nil {
		return handleCreateError("IdsClusterConfig", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyIdpsClusterConfigRead(d, m)
}

func resourceNsxtPolicyIdpsClusterConfigRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IDPS Cluster Config ID")
	}

	client := intrusion_services.NewClusterConfigsClient(connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	obj, err := client.Get(id, nil)
	if err != nil {
		return handleReadError(d, "IdsClusterConfig", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("ids_enabled", obj.IdsEnabled)

	// Set cluster information
	if obj.Cluster != nil {
		clusterList := make([]map[string]interface{}, 0, 1)
		clusterMap := make(map[string]interface{})
		clusterMap["target_id"] = obj.Cluster.TargetId
		clusterMap["target_type"] = obj.Cluster.TargetType
		clusterList = append(clusterList, clusterMap)
		d.Set("cluster", clusterList)
	}

	return nil
}

func resourceNsxtPolicyIdpsClusterConfigUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IDPS Cluster Config ID")
	}

	client := intrusion_services.NewClusterConfigsClient(connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	idsClusterConfig := buildIdsClusterConfig(d, id)

	log.Printf("[INFO] Updating IDPS Cluster Config with ID %s for cluster %s (ids_enabled: %v)", id, *idsClusterConfig.Cluster.TargetId, *idsClusterConfig.IdsEnabled)
	err := client.Patch(id, idsClusterConfig)
	if err != nil {
		return handleUpdateError("IdsClusterConfig", id, err)
	}

	return resourceNsxtPolicyIdpsClusterConfigRead(d, m)
}

func resourceNsxtPolicyIdpsClusterConfigDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IDPS Cluster Config ID")
	}

	client := intrusion_services.NewClusterConfigsClient(connector)
	if client == nil {
		return policyResourceNotSupportedError()
	}

	// IDPS Cluster Config cannot be truly deleted, only disabled
	// This is a configuration resource that persists in NSX
	// The delete operation disables IDPS on the cluster via PATCH
	log.Printf("[INFO] Disabling IDPS Cluster Config with ID %s (DELETE not supported by NSX API)", id)

	// Disable IDPS on the cluster by setting ids_enabled to false
	idsClusterConfig := buildIdsClusterConfig(d, id)
	idsClusterConfig.IdsEnabled = &[]bool{false}[0] // Set to false to disable

	err := client.Patch(id, idsClusterConfig)
	if err != nil {
		return handleDeleteError("IdsClusterConfig", id, err)
	}

	return nil
}
