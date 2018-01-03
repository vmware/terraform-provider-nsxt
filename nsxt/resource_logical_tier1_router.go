package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/manager"
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
		},
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
	logical_router := manager.LogicalRouter{
		Description:                     description,
		DisplayName:                     display_name,
		Tags:                            tags,
		FailoverMode:                    failover_mode,
		FirewallSections:                firewall_sections,
		HighAvailabilityMode:            high_availability_mode,
		RouterType:                      router_type,
	}

	logical_router, resp, err := nsxClient.LogicalRoutingAndServicesApi.CreateLogicalRouter(nsxClient.Context, logical_router)

	if err != nil {
		return fmt.Errorf("Error during LogicalTier1Router create: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Unexpected status returned: %s", resp.StatusCode)
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

	d.Set("Revision", logical_router.Revision)
	d.Set("SystemOwned", logical_router.SystemOwned)
	d.Set("Description", logical_router.Description)
	d.Set("DisplayName", logical_router.DisplayName)
	setTagsInSchema(d, logical_router.Tags)
	d.Set("AdvancedConfig", logical_router.AdvancedConfig)
	d.Set("EdgeClusterId", logical_router.EdgeClusterId)
	d.Set("FailoverMode", logical_router.FailoverMode)
	setResourceReferencesInSchema(d, logical_router.FirewallSections, "firewall_sections")
	d.Set("FirewallSections", logical_router.FirewallSections)
	d.Set("HighAvailabilityMode", logical_router.HighAvailabilityMode)
	d.Set("PreferredEdgeClusterMemberIndex", logical_router.PreferredEdgeClusterMemberIndex)
	d.Set("RouterType", logical_router.RouterType)

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
	logical_router := manager.LogicalRouter{
		Revision:                        revision,
		Description:                     description,
		DisplayName:                     display_name,
		Tags:                            tags,
		FailoverMode:                    failover_mode,
		FirewallSections:                firewall_sections,
		HighAvailabilityMode:            high_availability_mode,
	}

	logical_router, resp, err := nsxClient.LogicalRoutingAndServicesApi.UpdateLogicalRouter(nsxClient.Context, id, logical_router)

	if err != nil || resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Error during LogicalTier1Router update: %v", err)
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
