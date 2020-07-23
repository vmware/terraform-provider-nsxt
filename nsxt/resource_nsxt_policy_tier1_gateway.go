/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	gm_tier_1s "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/tier_1s"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
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
			"locale_service":    getPolicyLocaleServiceSchema(true),
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
			"intersite_config":         getGatewayIntersiteConfigSchema(),
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

func listGlobalManagerTier1GatewayLocaleServices(connector *client.RestConnector, gwID string, cursor *string) (model.LocaleServicesListResult, error) {
	client := gm_tier_1s.NewDefaultLocaleServicesClient(connector)
	markForDelete := false
	listResponse, err := client.List(gwID, cursor, &markForDelete, nil, nil, nil, nil)
	if err != nil {
		return model.LocaleServicesListResult{}, err
	}

	convertedResult, conversionErr := convertModelBindingType(listResponse, gm_model.LocaleServicesListResultBindingType(), model.LocaleServicesListResultBindingType())
	if conversionErr != nil {
		return model.LocaleServicesListResult{}, conversionErr
	}

	return convertedResult.(model.LocaleServicesListResult), nil
}

func listLocalManagerTier1GatewayLocaleServices(connector *client.RestConnector, gwID string, cursor *string) (model.LocaleServicesListResult, error) {
	client := tier_1s.NewDefaultLocaleServicesClient(connector)
	markForDelete := false
	return client.List(gwID, cursor, &markForDelete, nil, nil, nil, nil)
}

func listPolicyTier1GatewayLocaleServices(connector *client.RestConnector, gwID string, isGlobalManager bool) ([]model.LocaleServices, error) {

	if isGlobalManager {
		return listPolicyGatewayLocaleServices(connector, gwID, listGlobalManagerTier1GatewayLocaleServices)
	}

	return listPolicyGatewayLocaleServices(connector, gwID, listLocalManagerTier1GatewayLocaleServices)
}

func getPolicyTier1GatewayLocaleServiceEntry(gwID string, connector *client.RestConnector) (*model.LocaleServices, error) {
	// Get the locale services of this Tier1 for the edge-cluster id
	client := tier_1s.NewDefaultLocaleServicesClient(connector)
	obj, err := client.Get(gwID, defaultPolicyLocaleServiceID)
	if err == nil {
		return &obj, nil
	}

	// No locale-service with the default ID
	// List all the locale services
	objList, errList := listPolicyTier1GatewayLocaleServices(connector, gwID, false)
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

func resourceNsxtPolicyTier1GatewayReadEdgeCluster(d *schema.ResourceData, connector *client.RestConnector) error {
	// Get the locale services of this Tier1 for the edge-cluster id
	obj, err := getPolicyTier1GatewayLocaleServiceEntry(d.Id(), connector)
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

func resourceNsxtPolicyTier1GatewayExists(id string, connector *client.RestConnector, isGlobalManager bool) bool {
	var err error
	if isGlobalManager {
		client := gm_infra.NewDefaultTier1sClient(connector)
		_, err = client.Get(id)
	} else {
		client := infra.NewDefaultTier1sClient(connector)
		_, err = client.Get(id)
	}

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

func initSingleTier1GatewayLocaleService(d *schema.ResourceData, connector *client.RestConnector) (*data.StructValue, error) {

	edgeClusterPath := d.Get("edge_cluster_path").(string)
	var serviceStruct *model.LocaleServices
	var err error
	if len(d.Id()) > 0 {
		// This is an update flow - fetch existing locale service to reuse if needed
		serviceStruct, err = getPolicyTier1GatewayLocaleServiceEntry(d.Id(), connector)
		if err != nil {
			return nil, err
		}
	}

	if serviceStruct == nil {
		lsType := "LocaleServices"
		serviceStruct = &model.LocaleServices{
			Id:           &defaultPolicyLocaleServiceID,
			ResourceType: &lsType,
		}
	}
	if len(edgeClusterPath) > 0 {
		serviceStruct.EdgeClusterPath = &edgeClusterPath
	} else {
		serviceStruct.EdgeClusterPath = nil
	}

	log.Printf("[DEBUG] Using Locale Service with ID %s and Edge Cluster %v", *serviceStruct.Id, serviceStruct.EdgeClusterPath)
	return initChildLocaleService(serviceStruct, false)
}

func policyTier1GatewayResourceToInfraStruct(d *schema.ResourceData, connector *client.RestConnector, id string, isGlobalManager bool) (model.Infra, error) {
	var infraChildren, gwChildren []*data.StructValue
	var infraStruct model.Infra
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

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
	revision := int64(d.Get("revision").(int))

	t1Type := "Tier1"
	obj := model.Tier1{
		Id:                      &id,
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
		ResourceType:            &t1Type,
	}

	if dhcpPath != "" {
		dhcpPaths := []string{dhcpPath}
		obj.DhcpConfigPaths = dhcpPaths
	} else {
		obj.DhcpConfigPaths = []string{}
	}
	if len(d.Id()) > 0 {
		// This is update flow
		obj.Revision = &revision
	}

	resourceNsxtPolicyTier1GatewaySetVersionDependentAttrs(d, &obj)

	if isGlobalManager {
		intersiteConfig := getPolicyGatewayIntersiteConfigFromSchema(d)
		obj.IntersiteConfig = intersiteConfig
	}

	// set edge cluster for local manager if needed
	if d.HasChange("edge_cluster_path") && !isGlobalManager {
		dataValue, err := initSingleTier1GatewayLocaleService(d, connector)
		if err != nil {
			return infraStruct, err
		}

		gwChildren = append(gwChildren, dataValue)
	}

	if isGlobalManager {
		localeServices, err := initGatewayLocaleServices(d)
		if err != nil {
			return infraStruct, err
		}

		if len(localeServices) > 0 {
			gwChildren = append(gwChildren, localeServices...)
		}
	}

	obj.Children = gwChildren
	childTier1 := model.ChildTier1{
		Tier1:        &obj,
		ResourceType: "ChildTier1",
	}

	dataValue, errors := converter.ConvertToVapi(childTier1, model.ChildTier1BindingType())
	if errors != nil {
		return infraStruct, fmt.Errorf("Error converting Tier1 Child: %v", errors[0])
	}
	infraChildren = append(infraChildren, dataValue.(*data.StructValue))

	infraType := "Infra"
	infraStruct = model.Infra{
		Children:     infraChildren,
		ResourceType: &infraType,
	}

	return infraStruct, nil
}

func resourceNsxtPolicyTier1GatewayCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, m, resourceNsxtPolicyTier1GatewayExists)
	if err != nil {
		return err
	}

	obj, err := policyTier1GatewayResourceToInfraStruct(d, connector, id, isPolicyGlobalManager(m))
	if err != nil {
		return err
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Using H-API to create Tier1 with ID %s", id)
	err = policyInfraPatch(obj, isPolicyGlobalManager(m), connector, false)
	if err != nil {
		return handleCreateError("Tier1", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyTier1GatewayRead(d, m)
}

func resourceNsxtPolicyTier1GatewayRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	var obj model.Tier1
	var err error

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Tier1 id")
	}

	isGlobalManager := isPolicyGlobalManager(m)
	if isGlobalManager {
		client := gm_infra.NewDefaultTier1sClient(connector)
		gmObj, getErr := client.Get(id)
		if getErr != nil {
			return handleReadError(d, "Tier0", id, getErr)
		}

		convertedObj, convErr := convertModelBindingType(gmObj, model.Tier1BindingType(), model.Tier1BindingType())
		if convErr != nil {
			return convErr
		}
		obj = convertedObj.(model.Tier1)
	} else {
		client := infra.NewDefaultTier1sClient(connector)
		obj, err = client.Get(id)
	}

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

	// Get the edge cluster Id or locale services
	localeServices, err := listPolicyTier1GatewayLocaleServices(connector, id, isGlobalManager)
	if err != nil {
		return handleReadError(d, "Locale Service for T1", id, err)
	}
	var services []map[string]interface{}
	if len(localeServices) > 0 {

		for _, service := range localeServices {
			if isGlobalManager {
				cfgMap := make(map[string]interface{})
				cfgMap["path"] = service.Path
				cfgMap["edge_cluster_path"] = service.EdgeClusterPath
				cfgMap["preferred_edge_paths"] = service.PreferredEdgePaths
				cfgMap["revision"] = service.Revision
				services = append(services, cfgMap)

			} else {
				if service.EdgeClusterPath != nil {
					d.Set("edge_cluster_path", service.EdgeClusterPath)
				}
			}
		}

	} else {
		d.Set("edge_cluster_path", "")
	}

	if isGlobalManager {
		d.Set("locale_service", services)
	}

	err = setAdvRulesInSchema(d, obj.RouteAdvertisementRules)
	if err != nil {
		return fmt.Errorf("Error during Tier1 advertisement rules set in schema: %v", err)
	}

	err = setIpv6ProfilePathsInSchema(d, obj.Ipv6ProfilePaths)
	if err != nil {
		return fmt.Errorf("Failed to get Tier1 %s ipv6 profiles: %v", *obj.Id, err)
	}

	if isGlobalManager {
		err = setPolicyGatewayIntersiteConfigInSchema(d, obj.IntersiteConfig)
		if err != nil {
			return fmt.Errorf("Failed to get Tier1 %s interset config: %v", *obj.Id, err)
		}
	}

	return nil
}

func resourceNsxtPolicyTier1GatewayUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Tier1 id")
	}

	obj, err := policyTier1GatewayResourceToInfraStruct(d, connector, id, isPolicyGlobalManager(m))
	if err != nil {
		return err
	}

	log.Printf("[INFO] Using H-API to update Tier1 with ID %s", id)
	err = policyInfraPatch(obj, isPolicyGlobalManager(m), connector, true)
	if err != nil {
		return handleUpdateError("Tier1", id, err)
	}

	return resourceNsxtPolicyTier1GatewayRead(d, m)
}

func resourceNsxtPolicyTier1GatewayDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Tier1 id")
	}

	var infraChildren []*data.StructValue
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)
	boolTrue := true

	t1Type := "Tier1"
	t1obj := model.Tier1{
		Id:           &id,
		ResourceType: &t1Type,
	}
	childT1 := model.ChildTier1{
		MarkedForDelete: &boolTrue,
		Tier1:           &t1obj, ResourceType: "ChildTier1",
	}
	dataValue, errors := converter.ConvertToVapi(childT1, model.ChildTier1BindingType())
	if errors != nil {
		return fmt.Errorf("Error converting Child Tier1: %v", errors[0])
	}
	infraChildren = append(infraChildren, dataValue.(*data.StructValue))
	infraType := "Infra"
	obj := model.Infra{
		Children:     infraChildren,
		ResourceType: &infraType,
	}

	log.Printf("[DEBUG] Using H-API to delete Tier1 with ID %s", id)
	err := policyInfraPatch(obj, isPolicyGlobalManager(m), getPolicyConnector(m), false)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Success deleting Tier1 with ID %s", id)

	return nil
}
