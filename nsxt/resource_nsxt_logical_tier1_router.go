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

var failOverModeValues = []string{"PREEMPTIVE", "NON_PREEMPTIVE"}
var HighAvailabilityValues = []string{"ACTIVE_ACTIVE", "ACTIVE_STANDBY"}

// TODO: add advanced config
func resourceLogicalTier1Router() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogicalTier1RouterCreate,
		Read:   resourceLogicalTier1RouterRead,
		Update: resourceLogicalTier1RouterUpdate,
		Delete: resourceLogicalTier1RouterDelete,

		Schema: map[string]*schema.Schema{
			"revision": getRevisionSchema(),
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Description of this resource",
				Optional:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Defaults to ID if not set",
				Optional:    true,
			},
			"tags": getTagsSchema(),
			"failover_mode": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "This failover mode determines, whether the preferred service router instance for given logical router will preempt the peer. Note - It can be specified if and only if logical router is ACTIVE_STANDBY and NON_PREEMPTIVE mode is supported only for a Tier1 logical router. For Tier0 ACTIVE_STANDBY logical router, failover mode is always PREEMPTIVE, i.e. once the preferred node comes up after a failure, it will preempt the peer causing failover from current active to preferred node. For ACTIVE_ACTIVE logical routers, this field must not be populated",
				Default:      "PREEMPTIVE",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(failOverModeValues, false),
			},
			"firewall_sections": getResourceReferencesSchema(false, true, []string{}),
			"high_availability_mode": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "High availability mode",
				Default:      "ACTIVE_STANDBY",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(HighAvailabilityValues, false),
			},
			"edge_cluster_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Edge Cluster Id",
				Optional:    true,
			},
			"enable_router_advertisement": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Enable router advertisement",
				Default:     false,
				Optional:    true,
			},
			"advertise_connected_routes": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Enable connected NSX routes advertisement",
				Default:     false,
				Optional:    true,
			},
			"advertise_static_routes": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Enable static routes advertisement",
				Default:     false,
				Optional:    true,
			},
			"advertise_nat_routes": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Enable NAT routes advertisement",
				Default:     false,
				Optional:    true,
			},
			"advertise_config_revision": getRevisionSchema(),
		},
	}
}

func resourceLogicalTier1RouterCreateRollback(nsxClient *api.APIClient, id string) {
	log.Printf("[ERROR] Rollback router %s creation", id)

	localVarOptionals := make(map[string]interface{})
	_, err := nsxClient.LogicalRoutingAndServicesApi.DeleteLogicalRouter(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		log.Printf("[ERROR] Rollback failed!")
	}
}

func resourceLogicalTier1RouterReadAdv(d *schema.ResourceData, nsxClient *api.APIClient, id string) error {
	adv_config, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadAdvertisementConfig(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("LogicalTier1Router %s Advertisement config not found", id)
	}
	if err != nil {
		return fmt.Errorf("Error during LogicalTier1Router Advertisement config read: %v", err)
	}
	d.Set("enable_router_advertisement", adv_config.Enabled)
	d.Set("advertise_connected_routes", adv_config.AdvertiseNsxConnectedRoutes)
	d.Set("advertise_static_routes", adv_config.AdvertiseStaticRoutes)
	d.Set("advertise_nat_routes", adv_config.AdvertiseNatRoutes)
	d.Set("advertise_config_revision", adv_config.Revision)
	return nil
}

func resourceLogicalTier1RouterCreateAdv(d *schema.ResourceData, nsxClient *api.APIClient, id string) error {
	enable_router_advertisement := d.Get("enable_router_advertisement").(bool)
	if enable_router_advertisement {
		adv_connected := d.Get("advertise_connected_routes").(bool)
		adv_static := d.Get("advertise_static_routes").(bool)
		adv_nat := d.Get("advertise_nat_routes").(bool)
		adv_config := manager.AdvertisementConfig{
			Enabled:                     true,
			AdvertiseNsxConnectedRoutes: adv_connected,
			AdvertiseStaticRoutes:       adv_static,
			AdvertiseNatRoutes:          adv_nat,
		}
		_, _, err := nsxClient.LogicalRoutingAndServicesApi.UpdateAdvertisementConfig(nsxClient.Context, id, adv_config)
		return err
	}
	return nil
}

func resourceLogicalTier1RouterUpdateAdv(d *schema.ResourceData, nsxClient *api.APIClient, id string) error {
	enable_router_advertisement := d.Get("enable_router_advertisement").(bool)
	adv_connected := d.Get("advertise_connected_routes").(bool)
	adv_static := d.Get("advertise_static_routes").(bool)
	adv_nat := d.Get("advertise_nat_routes").(bool)
	adv_revision := int64(d.Get("advertise_config_revision").(int))
	adv_config := manager.AdvertisementConfig{
		Enabled:                     enable_router_advertisement,
		AdvertiseNsxConnectedRoutes: adv_connected,
		AdvertiseStaticRoutes:       adv_static,
		AdvertiseNatRoutes:          adv_nat,
		Revision:                    adv_revision,
	}
	_, _, err := nsxClient.LogicalRoutingAndServicesApi.UpdateAdvertisementConfig(nsxClient.Context, id, adv_config)
	return err
}

func resourceLogicalTier1RouterCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	failover_mode := d.Get("failover_mode").(string)
	high_availability_mode := d.Get("high_availability_mode").(string)
	router_type := "TIER1"
	edge_cluster_id := d.Get("edge_cluster_id").(string)
	logical_router := manager.LogicalRouter{
		Description:          description,
		DisplayName:          display_name,
		Tags:                 tags,
		FailoverMode:         failover_mode,
		HighAvailabilityMode: high_availability_mode,
		RouterType:           router_type,
		EdgeClusterId:        edge_cluster_id,
	}

	logical_router, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateLogicalRouter(nsxClient.Context, logical_router)

	if err != nil {
		return fmt.Errorf("Error during LogicalTier1Router create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned: %d", resp.StatusCode)
	}

	// Add advertisement config
	err = resourceLogicalTier1RouterCreateAdv(d, nsxClient, logical_router.Id)
	if err != nil {
		resourceLogicalTier1RouterCreateRollback(nsxClient, logical_router.Id)
		return fmt.Errorf("Error while setting config advertisement state: %v", err)
	}

	d.SetId(logical_router.Id)

	return resourceLogicalTier1RouterRead(d, m)
}

func resourceLogicalTier1RouterRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical tier1 router id")
	}

	logical_router, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouter(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalTier1Router %s not found", id)
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during LogicalTier1Router read: %v", err)
	}

	d.Set("revision", logical_router.Revision)
	d.Set("description", logical_router.Description)
	d.Set("display_name", logical_router.DisplayName)
	setTagsInSchema(d, logical_router.Tags)
	d.Set("advanced_config", logical_router.AdvancedConfig)
	d.Set("edge_cluster_id", logical_router.EdgeClusterId)
	d.Set("failover_mode", logical_router.FailoverMode)
	setResourceReferencesInSchema(d, logical_router.FirewallSections, "firewall_sections")
	d.Set("high_availability_mode", logical_router.HighAvailabilityMode)
	d.Set("preferred_edge_cluster_member_index", logical_router.PreferredEdgeClusterMemberIndex)
	d.Set("router_type", logical_router.RouterType)

	err = resourceLogicalTier1RouterReadAdv(d, nsxClient, id)
	if err != nil {
		return fmt.Errorf("Error during LogicalTier1Router read: %v", err)
	}

	return nil
}

func resourceLogicalTier1RouterUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical tier1 router id")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	failover_mode := d.Get("failover_mode").(string)
	high_availability_mode := d.Get("high_availability_mode").(string)
	router_type := "TIER1"
	edge_cluster_id := d.Get("edge_cluster_id").(string)
	logical_router := manager.LogicalRouter{
		Revision:             revision,
		Description:          description,
		DisplayName:          display_name,
		Tags:                 tags,
		FailoverMode:         failover_mode,
		HighAvailabilityMode: high_availability_mode,
		RouterType:           router_type,
		EdgeClusterId:        edge_cluster_id,
	}
	logical_router, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateLogicalRouter(nsxClient.Context, id, logical_router)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LogicalTier1Router update error: %v, resp %+v", err, resp)
	}

	// Update advertisement config
	err = resourceLogicalTier1RouterUpdateAdv(d, nsxClient, id)
	if err != nil {
		return fmt.Errorf("Error while setting config advertisement state: %v", err)
	}

	return resourceLogicalTier1RouterRead(d, m)
}

func resourceLogicalTier1RouterDelete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical tier1 router id")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteLogicalRouter(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during LogicalTier1Router delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalTier1Router %s not found", id)
		d.SetId("")
	}
	return nil
}
