/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"log"
	"net/http"
)

var highAvailabilityValues = []string{"ACTIVE_ACTIVE", "ACTIVE_STANDBY"}

// TODO: add advanced config
func resourceNsxtLogicalTier0Router() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtLogicalTier0RouterCreate,
		Read:   resourceNsxtLogicalTier0RouterRead,
		Update: resourceNsxtLogicalTier0RouterUpdate,
		Delete: resourceNsxtLogicalTier0RouterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource. Defaults to ID if not set",
				Optional:    true,
				Computed:    true,
			},
			"tag": getTagsSchema(),
			"high_availability_mode": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "High availability mode",
				Default:      "ACTIVE_STANDBY",
				Optional:     true,
				ForceNew:     true, // Cannot change the HA mode of a router
				ValidateFunc: validation.StringInSlice(highAvailabilityValues, false),
			},
			"firewall_sections": getResourceReferencesSchema(false, true, []string{}, "List of Firewall sections related to the logical router"),
			"edge_cluster_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Edge Cluster Id",
				Optional:    true,
			},
		},
	}
}

func resourceNsxtLogicalTier0RouterCreate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	highAvailabilityMode := d.Get("high_availability_mode").(string)
	routerType := "TIER0"
	edgeClusterID := d.Get("edge_cluster_id").(string)
	logicalRouter := manager.LogicalRouter{
		Description:          description,
		DisplayName:          displayName,
		Tags:                 tags,
		RouterType:           routerType,
		EdgeClusterId:        edgeClusterID,
		HighAvailabilityMode: highAvailabilityMode,
	}
	logicalRouter, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateLogicalRouter(nsxClient.Context, logicalRouter)

	if err != nil {
		return fmt.Errorf("Error during LogicalTier0Router create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned: %d", resp.StatusCode)
	}

	d.SetId(logicalRouter.Id)

	return resourceNsxtLogicalTier0RouterRead(d, m)
}

func resourceNsxtLogicalTier0RouterRead(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical Tier-1 router id")
	}

	logicalRouter, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouter(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LogicalTier0Router %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during LogicalTier0Router read: %v", err)
	}

	d.Set("revision", logicalRouter.Revision)
	d.Set("description", logicalRouter.Description)
	d.Set("display_name", logicalRouter.DisplayName)
	setTagsInSchema(d, logicalRouter.Tags)
	d.Set("advanced_config", logicalRouter.AdvancedConfig)
	d.Set("edge_cluster_id", logicalRouter.EdgeClusterId)
	d.Set("high_availability_mode", logicalRouter.HighAvailabilityMode)
	err = setResourceReferencesInSchema(d, logicalRouter.FirewallSections, "firewall_sections")
	if err != nil {
		return fmt.Errorf("Error during LogicalTier0Router firewall sections set in schema: %v", err)
	}
	d.Set("preferred_edge_cluster_member_index", logicalRouter.PreferredEdgeClusterMemberIndex)
	d.Set("router_type", logicalRouter.RouterType)

	return nil
}

func resourceNsxtLogicalTier0RouterUpdate(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical tier1 router id")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	highAvailabilityMode := d.Get("high_availability_mode").(string)
	routerType := "TIER0"
	edgeClusterID := d.Get("edge_cluster_id").(string)
	logicalRouter := manager.LogicalRouter{
		Revision:             revision,
		Description:          description,
		DisplayName:          displayName,
		Tags:                 tags,
		RouterType:           routerType,
		EdgeClusterId:        edgeClusterID,
		HighAvailabilityMode: highAvailabilityMode,
	}
	logicalRouter, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateLogicalRouter(nsxClient.Context, id, logicalRouter)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LogicalTier0Router update error: %v, resp %+v", err, resp)
	}

	return resourceNsxtLogicalTier0RouterRead(d, m)
}

func resourceNsxtLogicalTier0RouterDelete(d *schema.ResourceData, m interface{}) error {
	nsxClient := m.(*api.APIClient)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical tier1 router id")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteLogicalRouter(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during LogicalTier0Router delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] LogicalTier0Router %s not found", id)
		d.SetId("")
	}
	return nil
}
