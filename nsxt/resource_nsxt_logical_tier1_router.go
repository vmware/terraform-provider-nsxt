package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"log"
	"net/http"
)

// TODO: add advanced config
func resourceLogicalTier1Router() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogicalTier1RouterCreate,
		Read:   resourceLogicalTier1RouterRead,
		Update: resourceLogicalTier1RouterUpdate,
		Delete: resourceLogicalTier1RouterDelete,

		Schema: map[string]*schema.Schema{
			"revision":     getRevisionSchema(),
			"system_owned": getSystemOwnedSchema(),
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
				Type:        schema.TypeString,
				Description: "This failover mode determines, whether the preferred service router instance for given logical router will preempt the peer. Note - It can be specified if and only if logical router is ACTIVE_STANDBY and NON_PREEMPTIVE mode is supported only for a Tier1 logical router. For Tier0 ACTIVE_STANDBY logical router, failover mode is always PREEMPTIVE, i.e. once the preferred node comes up after a failure, it will preempt the peer causing failover from current active to preferred node. For ACTIVE_ACTIVE logical routers, this field must not be populated",
				Optional:    true,
			},
			"firewall_sections": getResourceReferencesSchema(false, true),
			"high_availability_mode": &schema.Schema{
				Type:        schema.TypeString,
				Description: "High availability mode",
				Optional:    true,
			},
			"edge_cluster_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Edge Cluster Id",
				Optional:    true,
			},
		},
	}
}

func resourceLogicalTier1RouterCreateRollback(nsxClient *api.APIClient, id string) {
	log.Printf("[ERROR] Rollback router %d creation", id)

	localVarOptionals := make(map[string]interface{})
	_, err := nsxClient.LogicalRoutingAndServicesApi.DeleteLogicalRouter(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		log.Printf("[ERROR] Rollback failed!")
	}
}

func resourceLogicalTier1RouterCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	failover_mode := d.Get("failover_mode").(string)
	firewall_sections := getResourceReferencesFromSchema(d, "firewall_sections")
	high_availability_mode := d.Get("high_availability_mode").(string)
	router_type := "TIER1"
	edge_cluster_id := d.Get("edge_cluster_id").(string)
	logical_router := manager.LogicalRouter{
		Description:          description,
		DisplayName:          display_name,
		Tags:                 tags,
		FailoverMode:         failover_mode,
		FirewallSections:     firewall_sections,
		HighAvailabilityMode: high_availability_mode,
		RouterType:           router_type,
		EdgeClusterId:        edge_cluster_id,
	}

	logical_router, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateLogicalRouter(nsxClient.Context, logical_router)

	if err != nil {
		return fmt.Errorf("Error during LogicalTier1Router create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned: %s", resp.StatusCode)
	}

	// Add advertisement config
	// TODO - make advertisement parameters configurable
	adv_config := manager.AdvertisementConfig{
		Enabled:                     true,
		AdvertiseNsxConnectedRoutes: true,
		AdvertiseStaticRoutes:       true,
		AdvertiseNatRoutes:          true,
	}

	_, resp, err1 := nsxClient.LogicalRoutingAndServicesApi.UpdateAdvertisementConfig(nsxClient.Context, logical_router.Id, adv_config)

	if err1 != nil {
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
		return fmt.Errorf("Error obtaining logical object id")
	}

	logical_router, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouter(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalTier1Router not found")
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during LogicalTier1Router read: %v", err)
	}

	d.Set("revision", logical_router.Revision)
	d.Set("system_owned", logical_router.SystemOwned)
	d.Set("description", logical_router.Description)
	d.Set("display_name", logical_router.DisplayName)
	setTagsInSchema(d, logical_router.Tags)
	d.Set("advanced_config", logical_router.AdvancedConfig)
	d.Set("edge_cluster_id", logical_router.EdgeClusterId)
	d.Set("failover_mode", logical_router.FailoverMode)
	setResourceReferencesInSchema(d, logical_router.FirewallSections, "firewall_sections")
	d.Set("firewall_sections", logical_router.FirewallSections)
	d.Set("high_availability_mode", logical_router.HighAvailabilityMode)
	d.Set("preferred_edge_cluster_member_index", logical_router.PreferredEdgeClusterMemberIndex)
	d.Set("router_type", logical_router.RouterType)
	d.Set("resource_type", logical_router.ResourceType)

	return nil
}

func resourceLogicalTier1RouterUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := getTagsFromSchema(d)
	failover_mode := d.Get("failover_mode").(string)
	firewall_sections := getResourceReferencesFromSchema(d, "firewall_sections")
	high_availability_mode := d.Get("high_availability_mode").(string)
	router_type := "TIER1"
	edge_cluster_id := d.Get("edge_cluster_id").(string)
	logical_router := manager.LogicalRouter{
		Revision:             revision,
		Description:          description,
		DisplayName:          display_name,
		Tags:                 tags,
		FailoverMode:         failover_mode,
		FirewallSections:     firewall_sections,
		HighAvailabilityMode: high_availability_mode,
		RouterType:           router_type,
		EdgeClusterId:        edge_cluster_id,
	}
	logical_router, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateLogicalRouter(nsxClient.Context, id, logical_router)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LogicalTier1Router update error: %v, resp %+v", err, resp)
	}

	return resourceLogicalTier1RouterRead(d, m)
}

func resourceLogicalTier1RouterDelete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteLogicalRouter(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during LogicalTier1Router delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalTier1Router not found")
		d.SetId("")
	}
	return nil
}
