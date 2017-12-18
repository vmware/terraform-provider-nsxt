package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
	"net/http"
)

// TODO: add advanced config
func resourceLogicalRouter() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogicalRouterCreate,
		Read:   resourceLogicalRouterRead,
		Update: resourceLogicalRouterUpdate,
		Delete: resourceLogicalRouterDelete,

		Schema: map[string]*schema.Schema{
			"revision":     GetRevisionSchema(),
			"system_owned": GetSystemOwnedSchema(),
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
			"tags": GetTagsSchema(),
			"edge_cluster_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Used for tier0 routers",
				Optional:    true,
			},
			"failover_mode": &schema.Schema{
				Type:        schema.TypeString,
				Description: "This failover mode determines, whether the preferred service router instance for given logical router will preempt the peer. Note - It can be specified if and only if logical router is ACTIVE_STANDBY and NON_PREEMPTIVE mode is supported only for a Tier1 logical router. For Tier0 ACTIVE_STANDBY logical router, failover mode is always PREEMPTIVE, i.e. once the preferred node comes up after a failure, it will preempt the peer causing failover from current active to preferred node. For ACTIVE_ACTIVE logical routers, this field must not be populated",
				Optional:    true,
			},
			"firewall_sections": GetResourceReferencesSchema(false, true),
			"high_availability_mode": &schema.Schema{
				Type:        schema.TypeString,
				Description: "High availability mode",
				Optional:    true,
			},
			"preferred_edge_cluster_member_index": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Used for tier0 routers only",
				Optional:    true,
			},
			"router_type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Type of Logical Router",
				Optional:    true,
				ForceNew:    true,
				Default:     "TIER1",
			},
		},
	}
}

func resourceLogicalRouterCreate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := GetTagsFromSchema(d)
	edge_cluster_id := d.Get("edge_cluster_id").(string)
	failover_mode := d.Get("failover_mode").(string)
	firewall_sections := GetResourceReferencesFromSchema(d, "firewall_sections")
	high_availability_mode := d.Get("high_availability_mode").(string)
	preferred_edge_cluster_member_index := int64(d.Get("preferred_edge_cluster_member_index").(int))
	router_type := d.Get("router_type").(string)
	logical_router := manager.LogicalRouter{
		Description:                     description,
		DisplayName:                     display_name,
		Tags:                            tags,
		EdgeClusterId:                   edge_cluster_id,
		FailoverMode:                    failover_mode,
		FirewallSections:                firewall_sections,
		HighAvailabilityMode:            high_availability_mode,
		PreferredEdgeClusterMemberIndex: preferred_edge_cluster_member_index,
		RouterType:                      router_type,
	}

	logical_router, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateLogicalRouter(nsxClient.Context, logical_router)

	if err != nil {
		return fmt.Errorf("Error during LogicalRouter create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned: %s", resp.StatusCode)
	}
	d.SetId(logical_router.Id)

	return resourceLogicalRouterRead(d, m)
}

func resourceLogicalRouterRead(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	logical_router, resp, err := nsxClient.LogicalRoutingAndServicesApi.ReadLogicalRouter(nsxClient.Context, id)
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalRouter not found")
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error during LogicalRouter read: %v", err)
	}

	d.Set("Revision", logical_router.Revision)
	d.Set("SystemOwned", logical_router.SystemOwned)
	d.Set("Description", logical_router.Description)
	d.Set("DisplayName", logical_router.DisplayName)
	SetTagsInSchema(d, logical_router.Tags)
	d.Set("AdvancedConfig", logical_router.AdvancedConfig)
	d.Set("EdgeClusterId", logical_router.EdgeClusterId)
	d.Set("FailoverMode", logical_router.FailoverMode)
	SetResourceReferencesInSchema(d, logical_router.FirewallSections, "firewall_sections")
	d.Set("FirewallSections", logical_router.FirewallSections)
	d.Set("HighAvailabilityMode", logical_router.HighAvailabilityMode)
	d.Set("PreferredEdgeClusterMemberIndex", logical_router.PreferredEdgeClusterMemberIndex)
	d.Set("RouterType", logical_router.RouterType)

	return nil
}

func resourceLogicalRouterUpdate(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	revision := int64(d.Get("revision").(int))
	description := d.Get("description").(string)
	display_name := d.Get("display_name").(string)
	tags := GetTagsFromSchema(d)
	edge_cluster_id := d.Get("edge_cluster_id").(string)
	failover_mode := d.Get("failover_mode").(string)
	firewall_sections := GetResourceReferencesFromSchema(d, "firewall_sections")
	high_availability_mode := d.Get("high_availability_mode").(string)
	preferred_edge_cluster_member_index := int64(d.Get("preferred_edge_cluster_member_index").(int))
	router_type := d.Get("router_type").(string)
	logical_router := manager.LogicalRouter{
		Revision:                        revision,
		Description:                     description,
		DisplayName:                     display_name,
		Tags:                            tags,
		EdgeClusterId:                   edge_cluster_id,
		FailoverMode:                    failover_mode,
		FirewallSections:                firewall_sections,
		HighAvailabilityMode:            high_availability_mode,
		PreferredEdgeClusterMemberIndex: preferred_edge_cluster_member_index,
		RouterType:                      router_type,
	}

	logical_router, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateLogicalRouter(nsxClient.Context, id, logical_router)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LogicalRouter update: %v", err)
	}

	return resourceLogicalRouterRead(d, m)
}

func resourceLogicalRouterDelete(d *schema.ResourceData, m interface{}) error {

	nsxClient := m.(*api.APIClient)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining logical object id")
	}

	localVarOptionals := make(map[string]interface{})
	resp, err := nsxClient.LogicalRoutingAndServicesApi.DeleteLogicalRouter(nsxClient.Context, id, localVarOptionals)
	if err != nil {
		return fmt.Errorf("Error during LogicalRouter delete: %v", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("LogicalRouter not found")
		d.SetId("")
	}
	return nil
}
