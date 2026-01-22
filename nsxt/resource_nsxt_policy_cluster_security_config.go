// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/security"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func resourceNsxtPolicyClusterSecurityConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyClusterSecurityConfigCreate,
		Read:   resourceNsxtPolicyClusterSecurityConfigRead,
		Update: resourceNsxtPolicyClusterSecurityConfigUpdate,
		Delete: resourceNsxtPolicyClusterSecurityConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Cluster Security Configuration resource for NSX 9.1+. Manages DFW (Distributed Firewall) at the cluster level. Note: This resource updates existing cluster configuration created automatically by NSX when a cluster is added. It cannot create or delete cluster configurations.",

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Cluster external ID (e.g., uuid:domain-c20). This is the cluster identifier from compute collection.",
				Required:    true,
				ForceNew:    true,
			},
			"dfw_enabled": {
				Type:         schema.TypeBool,
				Description:  "Enable or disable Distributed Firewall (DFW) on the cluster. DFW must be enabled before IDPS can be configured.",
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"dfw_enabled"},
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "Display name of the cluster security configuration",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the cluster security configuration",
				Computed:    true,
			},
			"path": {
				Type:        schema.TypeString,
				Description: "NSX path of the cluster security configuration",
				Computed:    true,
			},
			"revision": {
				Type:        schema.TypeInt,
				Description: "Revision number of the cluster security configuration",
				Computed:    true,
			},
		},
	}
}

func getClusterSecurityConfigFromSchema(d *schema.ResourceData) model.ClusterSecurityConfiguration {
	featureName := "DFW"
	featureList := []model.ClusterSecurityFeature{}

	// Always include DFW feature in the update
	// Read dfw_enabled value (either from config or from state)
	dfwEnabled := d.Get("dfw_enabled").(bool)
	featureList = append(featureList, model.ClusterSecurityFeature{
		Feature: &featureName,
		Enabled: &dfwEnabled,
	})

	return model.ClusterSecurityConfiguration{
		Features: featureList,
	}
}

func setClusterSecurityConfigInSchema(d *schema.ResourceData, config model.ClusterSecurityConfiguration, isDataSource bool) error {
	// Extract DFW enabled status from the features list
	dfwEnabled := false
	if config.Features != nil {
		for _, feature := range config.Features {
			if feature.Feature != nil && *feature.Feature == "DFW" && feature.Enabled != nil {
				dfwEnabled = *feature.Enabled
				break
			}
		}
	}
	d.Set("dfw_enabled", dfwEnabled)

	if config.DisplayName != nil {
		d.Set("display_name", *config.DisplayName)
	}

	if config.Description != nil {
		d.Set("description", *config.Description)
	}

	if config.Path != nil {
		d.Set("path", *config.Path)
	}

	// Only set revision for resources, not data sources
	if !isDataSource && config.Revision != nil {
		d.Set("revision", int(*config.Revision))
	}

	return nil
}

func resourceNsxtPolicyClusterSecurityConfigCreate(d *schema.ResourceData, m interface{}) error {
	// Version check: This resource requires NSX 9.1.0 or higher
	if !util.NsxVersionHigherOrEqual("9.1.0") {
		return fmt.Errorf("Cluster Security Config resource requires NSX version 9.1.0 or higher (current version: %s)", util.NsxVersion)
	}

	// Cluster security config is auto-created by NSX when cluster is added
	// This "create" is actually an update operation
	clusterID := d.Get("cluster_id").(string)
	d.SetId(clusterID)

	log.Printf("[INFO] Creating/Updating Cluster Security Config for cluster %s", clusterID)

	return resourceNsxtPolicyClusterSecurityConfigUpdate(d, m)
}

func resourceNsxtPolicyClusterSecurityConfigRead(d *schema.ResourceData, m interface{}) error {
	// Version check: This resource requires NSX 9.1.0 or higher
	if !util.NsxVersionHigherOrEqual("9.1.0") {
		return fmt.Errorf("Cluster Security Config resource requires NSX version 9.1.0 or higher (current version: %s)", util.NsxVersion)
	}

	connector := getPolicyConnector(m)
	client := security.NewClusterConfigsClient(connector)

	clusterID := d.Id()
	if clusterID == "" {
		return fmt.Errorf("Error obtaining Cluster Security Config ID")
	}

	config, err := client.Get(clusterID, nil)
	if err != nil {
		return handleReadError(d, "Cluster Security Config", clusterID, err)
	}

	d.Set("cluster_id", clusterID)
	return setClusterSecurityConfigInSchema(d, config, false)
}

func resourceNsxtPolicyClusterSecurityConfigUpdate(d *schema.ResourceData, m interface{}) error {
	// Version check: This resource requires NSX 9.1.0 or higher
	if !util.NsxVersionHigherOrEqual("9.1.0") {
		return fmt.Errorf("Cluster Security Config resource requires NSX version 9.1.0 or higher (current version: %s)", util.NsxVersion)
	}

	connector := getPolicyConnector(m)
	client := security.NewClusterConfigsClient(connector)

	clusterID := d.Id()
	if clusterID == "" {
		return fmt.Errorf("Error obtaining Cluster Security Config ID")
	}

	config := getClusterSecurityConfigFromSchema(d)

	log.Printf("[INFO] Updating Cluster Security Config for cluster %s", clusterID)

	err := client.Patch(clusterID, config)
	if err != nil {
		return handleUpdateError("Cluster Security Config", clusterID, err)
	}

	return resourceNsxtPolicyClusterSecurityConfigRead(d, m)
}

func resourceNsxtPolicyClusterSecurityConfigDelete(d *schema.ResourceData, m interface{}) error {
	// Version check: This resource requires NSX 9.1.0 or higher
	if !util.NsxVersionHigherOrEqual("9.1.0") {
		return fmt.Errorf("Cluster Security Config resource requires NSX version 9.1.0 or higher (current version: %s)", util.NsxVersion)
	}

	connector := getPolicyConnector(m)
	client := security.NewClusterConfigsClient(connector)

	clusterID := d.Id()
	if clusterID == "" {
		return fmt.Errorf("Error obtaining Cluster Security Config ID")
	}

	// Cluster security config cannot be deleted, only features can be disabled
	// On delete, we disable all features
	log.Printf("[INFO] Disabling all security features for cluster %s (DELETE operation)", clusterID)

	// Get current config
	currentConfig, err := client.Get(clusterID, nil)
	if err != nil {
		// If config doesn't exist, that's fine for delete
		if isNotFoundError(err) {
			return nil
		}
		log.Printf("[WARN] Could not read cluster config during delete: %v", err)
		return nil
	}

	// Disable all features
	disabledFeatures := make([]model.ClusterSecurityFeature, 0, len(currentConfig.Features))
	for _, feature := range currentConfig.Features {
		disabled := false
		disabledFeatures = append(disabledFeatures, model.ClusterSecurityFeature{
			Feature: feature.Feature,
			Enabled: &disabled,
		})
	}

	config := model.ClusterSecurityConfiguration{
		Features: disabledFeatures,
	}

	err = client.Patch(clusterID, config)
	if err != nil {
		return handleDeleteError("Cluster Security Config", clusterID, err)
	}

	return nil
}
