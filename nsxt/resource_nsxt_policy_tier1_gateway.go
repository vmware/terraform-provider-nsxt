/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_1s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
)

var advertismentTypeValues = []string{
	model.RouteAdvertisementRule_ROUTE_ADVERTISEMENT_TYPES_STATIC_ROUTES,
	model.RouteAdvertisementRule_ROUTE_ADVERTISEMENT_TYPES_CONNECTED,
	model.RouteAdvertisementRule_ROUTE_ADVERTISEMENT_TYPES_NAT,
	model.RouteAdvertisementRule_ROUTE_ADVERTISEMENT_TYPES_LB_VIP,
	model.RouteAdvertisementRule_ROUTE_ADVERTISEMENT_TYPES_LB_SNAT,
	model.RouteAdvertisementRule_ROUTE_ADVERTISEMENT_TYPES_DNS_FORWARDER_IP,
	model.RouteAdvertisementRule_ROUTE_ADVERTISEMENT_TYPES_IPSEC_LOCAL_ENDPOINT}
var advertismentRuleActionValues = []string{
	model.RouteAdvertisementRule_ACTION_PERMIT,
	model.RouteAdvertisementRule_ACTION_DENY}
var advertismentRuleOperatorValues = []string{
	model.RouteAdvertisementRule_PREFIX_OPERATOR_GE,
	model.RouteAdvertisementRule_PREFIX_OPERATOR_EQ}

var poolAllocationValues = []string{
	model.Tier1_POOL_ALLOCATION_ROUTING,
	model.Tier1_POOL_ALLOCATION_LB_SMALL,
	model.Tier1_POOL_ALLOCATION_LB_MEDIUM,
	model.Tier1_POOL_ALLOCATION_LB_LARGE,
	model.Tier1_POOL_ALLOCATION_LB_XLARGE,
}

func resourceNsxtPolicyTier1Gateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyTier1GatewayCreate,
		Read:   resourceNsxtPolicyTier1GatewayRead,
		Update: resourceNsxtPolicyTier1GatewayUpdate,
		Delete: resourceNsxtPolicyTier1GatewayDelete,
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
			"edge_cluster_path": getPolicyEdgeClusterPathSchema(),
			"failover_mode":     getFailoverModeSchema(failOverModeDefaultValue),
			"default_rule_logging": {
				Type:        schema.TypeBool,
				Description: "Default rule logging",
				Default:     false,
				Optional:    true,
			},
			"enable_firewall": {
				Type:        schema.TypeBool,
				Description: "Enable edge firewall",
				Default:     true,
				Optional:    true,
			},
			"enable_standby_relocation": {
				Type:        schema.TypeBool,
				Description: "Enable standby relocation",
				Default:     false,
				Optional:    true,
			},
			"force_whitelisting": {
				Type:        schema.TypeBool,
				Description: "Force whitelisting",
				Default:     false,
				Optional:    true,
			},
			"tier0_path": {
				Type:         schema.TypeString,
				Description:  "The path of the connected Tier0",
				Optional:     true,
				ValidateFunc: validatePolicyPath(),
			},
			"route_advertisement_types": {
				Type:        schema.TypeSet,
				Description: "Enable different types of route advertisements",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(advertismentTypeValues, false),
				},
				Optional: true,
			},
			"route_advertisement_rule": getAdvRulesSchema(),
			"ipv6_ndra_profile_path":   getIPv6NDRAPathSchema(),
			"ipv6_dad_profile_path":    getIPv6DadPathSchema(),
			"dhcp_config_path":         getPolicyPathSchema(false, false, "Policy path to DHCP server or relay configuration to use for this Tier1"),
			"pool_allocation": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Description:  "Edge node allocation at different sizes for routing and load balancer service to meet performance and scalability requirements",
				Optional:     true,
				Default:      model.Tier1_POOL_ALLOCATION_ROUTING,
				ValidateFunc: validation.StringInSlice(poolAllocationValues, false),
			},
			"ingress_qos_profile_path": getPolicyPathSchema(false, false, "Policy path to gateway QoS profile in ingress direction"),
			"egress_qos_profile_path":  getPolicyPathSchema(false, false, "Policy path to gateway QoS profile in egress direction"),
		},
	}
}

func getAdvRulesSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of route advertisement rules",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "Name of this rule",
					Required:    true,
				},
				"subnets": {
					Type:        schema.TypeSet,
					Description: "List of network CIDRs to be routed",
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validateCidr(),
					},
					Required: true,
				},
				"action": {
					Type:         schema.TypeString,
					Description:  "Action to advertise filtered routes to the connected Tier0 gateway",
					Default:      model.RouteAdvertisementRule_ACTION_PERMIT,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(advertismentRuleActionValues, false),
				},
				"prefix_operator": {
					Type:         schema.TypeString,
					Description:  "Prefix operator to apply on networks",
					Default:      model.RouteAdvertisementRule_PREFIX_OPERATOR_GE,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(advertismentRuleOperatorValues, false),
				},
				"route_advertisement_types": {
					Type:        schema.TypeSet,
					Description: "Enable different types of route advertisements",
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.StringInSlice(advertismentTypeValues, false),
					},
					Optional: true,
				},
			},
		},
	}
}

func resourceNsxtPolicyTier1GatewayListLocaleServiceEntries(connector *client.RestConnector, t1ID string) ([]model.LocaleServices, error) {
	client := tier_1s.NewDefaultLocaleServicesClient(connector)
	var results []model.LocaleServices
	var cursor *string
	total := 0

	for {
		includeMarkForDeleteObjectsParam := false
		searchResponse, err := client.List(t1ID, cursor, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
		if err != nil {
			return results, err
		}
		results = append(results, searchResponse.Results...)
		if total == 0 {
			// first response
			total = int(*searchResponse.ResultCount)
		}
		cursor = searchResponse.Cursor
		if len(results) >= total {
			return results, nil
		}
	}
}

func resourceNsxtPolicyTier1GatewayCreateEdgeCluster(d *schema.ResourceData, connector *client.RestConnector) error {
	// Create a Tier1 locale service with the edge-cluster ID
	client := tier_1s.NewDefaultLocaleServicesClient(connector)
	edgeClusterPath := d.Get("edge_cluster_path").(string)
	objID := d.Id()
	// The default ID of the locale service will be the Tier1 ID
	obj := model.LocaleServices{
		EdgeClusterPath: &edgeClusterPath,
	}

	err := client.Patch(objID, defaultPolicyLocaleServiceID, obj)
	if err != nil {
		return handleCreateError("Tier1 locale service", objID, err)
	}
	return nil
}

func resourceNsxtPolicyTier1GatewayGetLocaleServiceEntry(gwID string, connector *client.RestConnector) (*model.LocaleServices, error) {
	// Get the locale services of this Tier1 for the edge-cluster id
	client := tier_1s.NewDefaultLocaleServicesClient(connector)
	obj, err := client.Get(gwID, defaultPolicyLocaleServiceID)
	if err == nil {
		return &obj, nil
	}

	// No locale-service with the default ID
	// List all the locale services
	objList, errList := resourceNsxtPolicyTier1GatewayListLocaleServiceEntries(connector, gwID)
	if errList != nil {
		return nil, fmt.Errorf("Error while reading Tier1 %v locale-services: %v", gwID, err)
	}
	for _, objInList := range objList {
		// Find the one with the edge cluster path
		if objInList.EdgeClusterPath != nil {
			return &objInList, nil
		}
	}
	// No locale service with edge cluster path found.
	// Return any of the locale services (To avoid creating a new one)
	for _, objInList := range objList {
		return &objInList, nil
	}

	// No locale service with edge cluster path found
	return nil, nil
}

func resourceNsxtPolicyTier1GatewayUpdateEdgeCluster(d *schema.ResourceData, connector *client.RestConnector) error {
	// Create or update a Tier1 locale service with the edge-cluster ID
	// The ID if the locale service should be searches as in case of imported Tier1 it is unknown
	client := tier_1s.NewDefaultLocaleServicesClient(connector)
	edgeClusterPath := d.Get("edge_cluster_path").(string)
	if edgeClusterPath == "" {
		return resourceNsxtPolicyTier1GatewayDeleteEdgeCluster(d, connector)
	}
	obj := model.LocaleServices{
		EdgeClusterPath: &edgeClusterPath,
	}
	objID := d.Id()

	err := client.Patch(objID, defaultPolicyLocaleServiceID, obj)
	if err != nil {
		return handleUpdateError("Tier1 locale service", defaultPolicyLocaleServiceID, err)
	}
	return nil
}

func resourceNsxtPolicyTier1GatewayReadEdgeCluster(d *schema.ResourceData, connector *client.RestConnector) error {
	// Get the locale services of this Tier1 for the edge-cluster id
	obj, err := resourceNsxtPolicyTier1GatewayGetLocaleServiceEntry(d.Id(), connector)
	if err != nil || obj == nil {
		// No locale-service found
		return nil
	}
	if obj.EdgeClusterPath != nil {
		d.Set("edge_cluster_path", obj.EdgeClusterPath)
		return nil
	}
	// No edge cluster found
	return nil
}

func resourceNsxtPolicyTier1GatewayDeleteEdgeCluster(d *schema.ResourceData, connector *client.RestConnector) error {
	// Find and delete the locale service of this Tier1 for the edge-cluster id
	client := tier_1s.NewDefaultLocaleServicesClient(connector)
	objID := d.Id()

	_, err := client.Get(objID, defaultPolicyLocaleServiceID)
	if err == nil {
		err = client.Delete(objID, defaultPolicyLocaleServiceID)
		if err != nil {
			logAPIError("Error During Tier1 locale-services deletion", err)
			return fmt.Errorf("Failed to delete Tier1 %s locale-services", objID)
		}
	}
	return nil
}

func resourceNsxtPolicyTier1GatewayExists(id string, connector *client.RestConnector) bool {
	client := infra.NewDefaultTier1sClient(connector)

	_, err := client.Get(id)
	if err == nil {
		return true
	}

	if isNotFoundError(err) {
		return false
	}

	logAPIError("Error retrieving Tier1", err)
	return false
}

func setAdvRulesInSchema(d *schema.ResourceData, rules []model.RouteAdvertisementRule) error {
	var rulesList []map[string]interface{}
	for _, rule := range rules {
		elem := make(map[string]interface{})
		elem["name"] = rule.Name
		elem["action"] = rule.Action
		elem["subnets"] = rule.Subnets
		elem["route_advertisement_types"] = rule.RouteAdvertisementTypes
		elem["prefix_operator"] = *rule.PrefixOperator
		rulesList = append(rulesList, elem)
	}
	err := d.Set("route_advertisement_rule", rulesList)
	return err
}

func getAdvRulesFromSchema(d *schema.ResourceData) []model.RouteAdvertisementRule {
	rules := d.Get("route_advertisement_rule").([]interface{})
	var ruleList []model.RouteAdvertisementRule
	for _, rule := range rules {
		data := rule.(map[string]interface{})
		prefix := data["prefix_operator"].(string)
		name := data["name"].(string)
		action := data["action"].(string)
		elem := model.RouteAdvertisementRule{
			Name:                    &name,
			Action:                  &action,
			Subnets:                 interface2StringList(data["subnets"].(*schema.Set).List()),
			RouteAdvertisementTypes: interface2StringList(data["route_advertisement_types"].(*schema.Set).List()),
			PrefixOperator:          &prefix,
		}
		ruleList = append(ruleList, elem)
	}
	return ruleList
}

func resourceNsxtPolicyTier1GatewaySetQos(d *schema.ResourceData, obj *model.Tier1) {
	ingressQosProfile := d.Get("ingress_qos_profile_path").(string)
	egressQosProfile := d.Get("egress_qos_profile_path").(string)

	if ingressQosProfile != "" || egressQosProfile != "" {
		qosConfig := model.GatewayQosProfileConfig{
			IngressQosProfilePath: &ingressQosProfile,
			EgressQosProfilePath:  &egressQosProfile,
		}
		obj.QosProfile = &qosConfig
	}

}

func resourceNsxtPolicyTier1GatewaySetVersionDependentAttrs(d *schema.ResourceData, obj *model.Tier1) {
	if nsxVersionLower("3.0.0") {
		return
	}

	resourceNsxtPolicyTier1GatewaySetQos(d, obj)
	poolAllocation := d.Get("pool_allocation").(string)
	if poolAllocation != "" {
		obj.PoolAllocation = &poolAllocation
	}

}

func resourceNsxtPolicyTier1GatewayCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultTier1sClient(connector)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, connector, resourceNsxtPolicyTier1GatewayExists)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	failoverMode := d.Get("failover_mode").(string)
	defaultRuleLogging := d.Get("default_rule_logging").(bool)
	disableFirewall := !d.Get("enable_firewall").(bool)
	enableStandbyRelocation := d.Get("enable_standby_relocation").(bool)
	forceWhitelisting := d.Get("force_whitelisting").(bool)
	tier0Path := d.Get("tier0_path").(string)
	routeAdvertisementTypes := getStringListFromSchemaSet(d, "route_advertisement_types")
	routeAdvertisementRules := getAdvRulesFromSchema(d)
	ipv6ProfilePaths := getIpv6ProfilePathsFromSchema(d)
	dhcpPath := d.Get("dhcp_config_path").(string)

	obj := model.Tier1{
		DisplayName:             &displayName,
		Description:             &description,
		Tags:                    tags,
		FailoverMode:            &failoverMode,
		DefaultRuleLogging:      &defaultRuleLogging,
		DisableFirewall:         &disableFirewall,
		EnableStandbyRelocation: &enableStandbyRelocation,
		ForceWhitelisting:       &forceWhitelisting,
		Tier0Path:               &tier0Path,
		RouteAdvertisementTypes: routeAdvertisementTypes,
		RouteAdvertisementRules: routeAdvertisementRules,
		Ipv6ProfilePaths:        ipv6ProfilePaths,
	}

	if dhcpPath != "" {
		dhcpPaths := []string{dhcpPath}
		obj.DhcpConfigPaths = dhcpPaths
	} else {
		obj.DhcpConfigPaths = []string{}
	}

	resourceNsxtPolicyTier1GatewaySetVersionDependentAttrs(d, &obj)

	// Create the resource using PATCH
	log.Printf("[INFO] Creating tier1 with ID %s", id)
	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("Tier1", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)
	// Add the edge cluster
	if d.Get("edge_cluster_path") != "" {
		err = resourceNsxtPolicyTier1GatewayCreateEdgeCluster(d, connector)
		if err != nil {
			log.Printf("[INFO] Rolling back Tier1 creation")
			client.Delete(id)
			return err
		}
	}

	return resourceNsxtPolicyTier1GatewayRead(d, m)
}

func resourceNsxtPolicyTier1GatewayRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultTier1sClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Tier1 id")
	}

	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "Tier1", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("failover_mode", obj.FailoverMode)
	d.Set("default_rule_logging", obj.DefaultRuleLogging)
	d.Set("enable_firewall", !(*obj.DisableFirewall))
	d.Set("enable_standby_relocation", obj.EnableStandbyRelocation)
	d.Set("force_whitelisting", obj.ForceWhitelisting)
	if obj.Tier0Path != nil {
		d.Set("tier0_path", *obj.Tier0Path)
	}
	d.Set("route_advertisement_types", obj.RouteAdvertisementTypes)
	d.Set("revision", obj.Revision)
	if obj.PoolAllocation == nil {
		// This will happen with NSX version < 3.0.0
		d.Set("pool_allocation", model.Tier1_POOL_ALLOCATION_ROUTING)
	} else {
		d.Set("pool_allocation", obj.PoolAllocation)
	}
	dhcpPaths := obj.DhcpConfigPaths

	if len(dhcpPaths) > 0 {
		d.Set("dhcp_config_path", dhcpPaths[0])
	}

	if obj.QosProfile != nil {
		d.Set("ingressQosProfile", obj.QosProfile.IngressQosProfilePath)
		d.Set("egressQosProfile", obj.QosProfile.EgressQosProfilePath)
	}

	// Get the edge cluster Id
	err = resourceNsxtPolicyTier1GatewayReadEdgeCluster(d, connector)
	if err != nil {
		return fmt.Errorf("Failed to get Tier1 %s locale-services: %v", *obj.Id, err)
	}

	err = setAdvRulesInSchema(d, obj.RouteAdvertisementRules)
	if err != nil {
		return fmt.Errorf("Error during Tier1 advertisement rules set in schema: %v", err)
	}

	err = setIpv6ProfilePathsInSchema(d, obj.Ipv6ProfilePaths)
	if err != nil {
		return fmt.Errorf("Failed to get Tier1 %s ipv6 profiles: %v", *obj.Id, err)
	}

	return nil
}

func resourceNsxtPolicyTier1GatewayUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultTier1sClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Tier1 id")
	}

	// Read the rest of the configured parameters
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)
	failoverMode := d.Get("failover_mode").(string)
	defaultRuleLogging := d.Get("default_rule_logging").(bool)
	disableFirewall := !d.Get("enable_firewall").(bool)
	enableStandbyRelocation := d.Get("enable_standby_relocation").(bool)
	forceWhitelisting := d.Get("force_whitelisting").(bool)
	tier0Path := d.Get("tier0_path").(string)
	routeAdvertisementTypes := getStringListFromSchemaSet(d, "route_advertisement_types")
	routeAdvertisementRules := getAdvRulesFromSchema(d)
	ipv6ProfilePaths := getIpv6ProfilePathsFromSchema(d)
	dhcpPath := d.Get("dhcp_config_path").(string)
	revision := int64(d.Get("revision").(int))

	obj := model.Tier1{
		DisplayName:             &displayName,
		Description:             &description,
		Tags:                    tags,
		FailoverMode:            &failoverMode,
		DefaultRuleLogging:      &defaultRuleLogging,
		DisableFirewall:         &disableFirewall,
		EnableStandbyRelocation: &enableStandbyRelocation,
		ForceWhitelisting:       &forceWhitelisting,
		Tier0Path:               &tier0Path,
		RouteAdvertisementTypes: routeAdvertisementTypes,
		RouteAdvertisementRules: routeAdvertisementRules,
		Ipv6ProfilePaths:        ipv6ProfilePaths,
		Revision:                &revision,
	}

	if dhcpPath != "" {
		dhcpPaths := []string{dhcpPath}
		obj.DhcpConfigPaths = dhcpPaths
	} else {
		obj.DhcpConfigPaths = []string{}
	}

	resourceNsxtPolicyTier1GatewaySetVersionDependentAttrs(d, &obj)

	// Update the resource using PUT
	_, err := client.Update(id, obj)
	if err != nil {
		return handleUpdateError("Tier1", id, err)
	}

	if d.HasChange("edge_cluster_path") {
		// Update edge cluster
		err = resourceNsxtPolicyTier1GatewayUpdateEdgeCluster(d, connector)
		if err != nil {
			return err
		}
	}

	return resourceNsxtPolicyTier1GatewayRead(d, m)
}

func resourceNsxtPolicyTier1GatewayDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Tier1 id")
	}

	connector := getPolicyConnector(m)
	err := resourceNsxtPolicyTier1GatewayDeleteEdgeCluster(d, connector)
	if err != nil {
		err = handleDeleteError("Tier1 locale service", id, err)
		if err != nil {
			return err
		}
	}

	client := infra.NewDefaultTier1sClient(connector)
	err = client.Delete(id)
	if err != nil {
		return handleDeleteError("Tier1", id, err)
	}

	return nil
}
