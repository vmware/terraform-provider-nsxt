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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/locale_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
)

var haModeValues = []string{
	model.Tier0_HA_MODE_ACTIVE,
	model.Tier0_HA_MODE_STANDBY}

var nsxtPolicyTier0GatewayBgpGracefulRestartModes = []string{
	model.BgpGracefulRestartConfig_MODE_DISABLE,
	model.BgpGracefulRestartConfig_MODE_GR_AND_HELPER,
	model.BgpGracefulRestartConfig_MODE_HELPER_ONLY,
}

var policyVRFRouteValues = []string{
	model.VrfRouteTargets_ADDRESS_FAMILY_EVPN,
}

var policyBGPGracefulRestartTimerDefault = 180
var policyBGPGracefulRestartStaleRouteTimerDefault = 600
var policyBGPLocalAsNumDefault = "65000"

func resourceNsxtPolicyTier0Gateway() *schema.Resource {

	return &schema.Resource{
		Create: resourceNsxtPolicyTier0GatewayCreate,
		Read:   resourceNsxtPolicyTier0GatewayRead,
		Update: resourceNsxtPolicyTier0GatewayUpdate,
		Delete: resourceNsxtPolicyTier0GatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":        getNsxIDSchema(),
			"path":          getPathSchema(),
			"display_name":  getDisplayNameSchema(),
			"description":   getDescriptionSchema(),
			"revision":      getRevisionSchema(),
			"tag":           getTagsSchema(),
			"failover_mode": getFailoverModeSchema(failOverModeDefaultPolicyT0Value),
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
			"force_whitelisting": {
				Type:        schema.TypeBool,
				Description: "Force whitelisting",
				Default:     false,
				Optional:    true,
			},
			"ha_mode": {
				Type:         schema.TypeString,
				Description:  "High-availability Mode for Tier-0",
				ValidateFunc: validation.StringInSlice(haModeValues, false),
				Optional:     true,
				Default:      model.Tier0_HA_MODE_ACTIVE,
				ForceNew:     true,
			},
			"internal_transit_subnets": {
				Type:        schema.TypeList,
				Description: "Internal transit subnets in CIDR format",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCidr(),
				},
				Optional: true,
				MaxItems: 1,
				Computed: true,
				ForceNew: true,
			},
			"transit_subnets": {
				Type:        schema.TypeList,
				Description: "Transit subnets in CIDR format",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCidr(),
				},
				Optional: true,
				Computed: true,
				ForceNew: true, // Modification of transit subnet not allowed after Tier-0 deployment
			},
			"ipv6_ndra_profile_path": getIPv6NDRAPathSchema(),
			"ipv6_dad_profile_path":  getIPv6DadPathSchema(),
			"edge_cluster_path":      getPolicyEdgeClusterPathSchema(),
			"bgp_config":             getPolicyBGPConfigSchema(),
			"vrf_config":             getPolicyVRFConfigSchema(),
			"dhcp_config_path":       getPolicyPathSchema(false, false, "Policy path to DHCP server or relay configuration to use for this Tier0"),
		},
	}
}

func getPolicyBGPConfigSchema() *schema.Schema {
	return &schema.Schema{
		// NOTE: setting bpg_config requires a edge_cluster_path
		Type:        schema.TypeList,
		Description: "BGP routing configuration",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tag":      getTagsSchema(),
				"revision": getRevisionSchema(),
				"path":     getPathSchema(),
				"ecmp": {
					Type:        schema.TypeBool,
					Description: "Flag to enable ECMP",
					Optional:    true,
					Default:     true,
				},
				"enabled": {
					Type:        schema.TypeBool,
					Description: "Flag to enable BGP configuration",
					Optional:    true,
					Default:     true,
				},
				"inter_sr_ibgp": {
					Type:        schema.TypeBool,
					Description: "Enable inter SR IBGP configuration",
					Optional:    true,
					Default:     true,
				},
				"local_as_num": {
					Type: schema.TypeString,
					Description: "	BGP AS number in ASPLAIN/ASDOT Format",
					Optional:     true,
					Default:      policyBGPLocalAsNumDefault, //NOTE: empty string disables
					ValidateFunc: validate4ByteASNPlain,
				},
				"multipath_relax": {
					Type:        schema.TypeBool,
					Description: "Flag to enable BGP multipath relax option",
					Optional:    true,
					Default:     true,
				},
				"route_aggregation": {
					Type:        schema.TypeList,
					Description: "List of routes to be aggregated",
					Optional:    true,
					MaxItems:    1000,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"prefix": {
								Type:         schema.TypeString,
								Description:  "CIDR of aggregate address",
								Optional:     true,
								ValidateFunc: validateCidr(),
							},
							"summary_only": {
								Type:        schema.TypeBool,
								Description: "Send only summarized route",
								Optional:    true,
								Default:     true,
							},
						},
					},
				},
				"graceful_restart_mode": {
					// BgpGracefulRestartConfig.mode
					Type:         schema.TypeString,
					Description:  "BGP Graceful Restart Configuration Mode",
					ValidateFunc: validation.StringInSlice(nsxtPolicyTier0GatewayBgpGracefulRestartModes, false),
					Optional:     true,
					Default:      model.BgpGracefulRestartConfig_MODE_HELPER_ONLY,
				},
				"graceful_restart_timer": {
					// BgpGracefulRestartConfig.timer.restart_timer
					Type:         schema.TypeInt,
					Description:  "BGP Graceful Restart Timer",
					Optional:     true,
					Default:      policyBGPGracefulRestartTimerDefault,
					ValidateFunc: validation.IntBetween(1, 3600),
				},
				"graceful_restart_stale_route_timer": {
					// BgpGracefulRestartConfig.timer.stale_route_timer
					Type:         schema.TypeInt,
					Description:  "BGP Stale Route Timer",
					Optional:     true,
					Default:      policyBGPGracefulRestartStaleRouteTimerDefault,
					ValidateFunc: validation.IntBetween(1, 3600),
				},
			},
		},
	}
}

func getVRFRouteSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validateIPorASNPair,
	}
}

func getVRFRouteElemSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validateASNPair,
	}
}

func getPolicyVRFConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "VRF configuration",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tag":          getTagsSchema(),
				"path":         getPathSchema(),
				"gateway_path": getPolicyPathSchema(true, true, "Default tier0 path"),
				"evpn_transit_vni": {
					Type:        schema.TypeInt,
					Description: "L3 VNI associated with the VRF for overlay traffic. VNI must be unique and belong to configured VNI pool",
					Optional:    true,
				},
				"route_distinguisher": getVRFRouteSchema(),
				"route_target": {
					Type:        schema.TypeList,
					Description: "Route targets",
					Optional:    true,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"address_family": {
								Type:         schema.TypeString,
								Optional:     true,
								Default:      policyVRFRouteValues[0],
								ValidateFunc: validation.StringInSlice(policyVRFRouteValues, false),
							},
							"auto_mode": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     true,
								Description: "When set to false, targets should be configured",
							},
							"import_targets": {
								Type:     schema.TypeList,
								Optional: true,
								Elem:     getVRFRouteElemSchema(),
							},
							"export_targets": {
								Type:     schema.TypeList,
								Optional: true,
								Elem:     getVRFRouteElemSchema(),
							},
						},
					},
				},
			},
		},
	}
}

func resourceNsxtPolicyTier0GatewayListEdgeClusterLocaleServiceEntries(connector *client.RestConnector, t0ID string) ([]model.LocaleServices, error) {
	client := tier_0s.NewDefaultLocaleServicesClient(connector)
	var results []model.LocaleServices
	var cursor *string
	total := 0

	for {
		includeMarkForDeleteObjectsParam := false
		searchResponse, err := client.List(t0ID, cursor, &includeMarkForDeleteObjectsParam, nil, nil, nil, nil)
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

func resourceNsxtPolicyTier0GatewayGetLocaleServiceEntry(gwID string, connector *client.RestConnector) (*model.LocaleServices, error) {
	// Get the locale services of this Tier0 for the edge-cluster id
	client := tier_0s.NewDefaultLocaleServicesClient(connector)
	obj, err := client.Get(gwID, defaultPolicyLocaleServiceID)
	if err == nil {
		return &obj, nil
	}

	// No locale-service with the default ID
	// List all the locale services
	objList, errList := resourceNsxtPolicyTier0GatewayListEdgeClusterLocaleServiceEntries(connector, gwID)
	if errList != nil {
		if isNotFoundError(errList) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while reading Tier0 %v locale-services: %v", gwID, errList)
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

	return nil, nil
}

func resourceNsxtPolicyTier0GatewayReadBGPConfig(d *schema.ResourceData, connector *client.RestConnector, localeService model.LocaleServices) error {
	var bgpConfigs []map[string]interface{}
	client := locale_services.NewDefaultBgpClient(connector)

	t0Id := d.Id()
	bgpConfig, err := client.Get(t0Id, *localeService.Id)
	if err != nil {
		if isNotFoundError(err) {
			return d.Set("bgp_config", bgpConfigs)
		}
		return err
	}

	cfgMap := make(map[string]interface{})
	cfgMap["revision"] = int(*bgpConfig.Revision)
	cfgMap["path"] = bgpConfig.Path
	cfgMap["ecmp"] = bgpConfig.Ecmp
	cfgMap["enabled"] = bgpConfig.Enabled
	cfgMap["inter_sr_ibgp"] = bgpConfig.InterSrIbgp
	cfgMap["local_as_num"] = bgpConfig.LocalAsNum
	cfgMap["multipath_relax"] = bgpConfig.MultipathRelax

	if bgpConfig.GracefulRestartConfig != nil {
		cfgMap["graceful_restart_mode"] = bgpConfig.GracefulRestartConfig.Mode
		if bgpConfig.GracefulRestartConfig.Timer != nil {
			cfgMap["graceful_restart_timer"] = int(*bgpConfig.GracefulRestartConfig.Timer.RestartTimer)
			cfgMap["graceful_restart_stale_route_timer"] = int(*bgpConfig.GracefulRestartConfig.Timer.StaleRouteTimer)
		}
	} else {
		// Assign defaults
		cfgMap["graceful_restart_mode"] = model.BgpGracefulRestartConfig_MODE_HELPER_ONLY
		cfgMap["graceful_restart_timer"] = policyBGPGracefulRestartTimerDefault
		cfgMap["graceful_restart_stale_route_timer"] = policyBGPGracefulRestartStaleRouteTimerDefault
	}

	var tagList []map[string]string
	for _, tag := range bgpConfig.Tags {
		elem := make(map[string]string)
		elem["scope"] = *tag.Scope
		elem["tag"] = *tag.Tag
		tagList = append(tagList, elem)
	}
	cfgMap["tag"] = tagList

	var aggregationList []map[string]interface{}
	for _, agg := range bgpConfig.RouteAggregations {
		elem := make(map[string]interface{})
		elem["prefix"] = agg.Prefix
		elem["summary_only"] = *agg.SummaryOnly
		aggregationList = append(aggregationList, elem)
	}
	cfgMap["route_aggregation"] = aggregationList

	bgpConfigs = append(bgpConfigs, cfgMap)

	return d.Set("bgp_config", bgpConfigs)
}

func getPolicyVRFConfigFromSchema(d *schema.ResourceData) *model.Tier0VrfConfig {

	if nsxVersionLower("3.0.0") {
		// VRF Lite is supported from 3.0.0 onwards
		return nil
	}

	vrfConfigs := d.Get("vrf_config").([]interface{})
	if len(vrfConfigs) == 0 {
		return nil
	}

	vrfConfig := vrfConfigs[0].(map[string]interface{})
	vni := int64(vrfConfig["evpn_transit_vni"].(int))
	gwPath := vrfConfig["gateway_path"].(string)
	routeDist := vrfConfig["route_distinguisher"].(string)

	var tagStructs []model.Tag
	if vrfConfig["tag"] != nil {
		vrfTags := vrfConfig["tag"].(*schema.Set).List()
		for _, tag := range vrfTags {
			data := tag.(map[string]interface{})
			tagScope := data["scope"].(string)
			tagTag := data["tag"].(string)
			elem := model.Tag{
				Scope: &tagScope,
				Tag:   &tagTag}
			tagStructs = append(tagStructs, elem)
		}
	}

	config := model.Tier0VrfConfig{
		Tier0Path: &gwPath,
		Tags:      tagStructs,
	}

	if len(routeDist) > 0 {
		config.RouteDistinguisher = &routeDist
	}

	if vni > 0 {
		config.EvpnTransitVni = &vni
	}

	routeTargets := vrfConfig["route_target"].([]interface{})
	if len(routeTargets) > 0 {
		routeTarget := routeTargets[0].(map[string]interface{})
		addressFamily := routeTarget["address_family"].(string)
		exportTargets := interface2StringList(routeTarget["export_targets"].([]interface{}))
		importTargets := interface2StringList(routeTarget["import_targets"].([]interface{}))
		// Only one is supported for now
		targets := model.VrfRouteTargets{
			AddressFamily:      &addressFamily,
			ExportRouteTargets: exportTargets,
			ImportRouteTargets: importTargets,
		}

		config.RouteTargets = []model.VrfRouteTargets{targets}
	}

	if vni > 0 {
		config.EvpnTransitVni = &vni
	}

	return &config
}

func setPolicyVRFConfigInSchema(d *schema.ResourceData, config *model.Tier0VrfConfig) error {
	if config == nil {
		return nil
	}

	var vrfConfigs []map[string]interface{}
	elem := make(map[string]interface{})
	elem["gateway_path"] = config.Tier0Path
	elem["route_distinguisher"] = config.RouteDistinguisher
	if config.RouteTargets != nil {
		routeTarget := make(map[string]interface{})
		routeTarget["address_family"] = config.RouteTargets[0].AddressFamily
		routeTarget["auto_mode"] = false
		if len(config.RouteTargets[0].ImportRouteTargets) > 0 || len(config.RouteTargets[0].ExportRouteTargets) > 0 {
			routeTarget["import_targets"] = config.RouteTargets[0].ImportRouteTargets
			routeTarget["export_targets"] = config.RouteTargets[0].ExportRouteTargets
		}
		var routeTargets []map[string]interface{}
		routeTargets = append(routeTargets, routeTarget)
		elem["route_target"] = routeTargets
	}

	var tagList []map[string]string
	for _, tag := range config.Tags {
		tagElem := make(map[string]string)
		tagElem["scope"] = *tag.Scope
		tagElem["tag"] = *tag.Tag
		tagList = append(tagList, tagElem)
	}
	elem["tag"] = tagList

	vrfConfigs = append(vrfConfigs, elem)

	return d.Set("vrf_config", vrfConfigs)
}

func resourceNsxtPolicyTier0GatewayExists(id string, connector *client.RestConnector) bool {
	client := infra.NewDefaultTier0sClient(connector)

	_, err := client.Get(id)
	if err == nil {
		return true
	}

	if isNotFoundError(err) {
		return false
	}

	logAPIError("Error retrieving Tier0", err)
	return false
}

func resourceNsxtPolicyTier0GatewayBGPConfigSchemaToStruct(cfg interface{}, isVrf bool, gwID string) model.BgpRoutingConfig {
	cfgMap := cfg.(map[string]interface{})
	revision := int64(cfgMap["revision"].(int))
	ecmp := cfgMap["ecmp"].(bool)
	enabled := cfgMap["enabled"].(bool)
	interSrIbgp := cfgMap["inter_sr_ibgp"].(bool)
	localAsNum := cfgMap["local_as_num"].(string)
	multipathRelax := cfgMap["multipath_relax"].(bool)
	restartMode := cfgMap["graceful_restart_mode"].(string)
	restartTimer := int64(cfgMap["graceful_restart_timer"].(int))
	staleTimer := int64(cfgMap["graceful_restart_stale_route_timer"].(int))

	var tagStructs []model.Tag
	if cfgMap["tag"] != nil {
		cfgTags := cfgMap["tag"].(*schema.Set).List()
		for _, tag := range cfgTags {
			data := tag.(map[string]interface{})
			tagScope := data["scope"].(string)
			tagTag := data["tag"].(string)
			elem := model.Tag{
				Scope: &tagScope,
				Tag:   &tagTag}

			tagStructs = append(tagStructs, elem)
		}
	}

	var aggregationStructs []model.RouteAggregationEntry
	routeAggregations := cfgMap["route_aggregation"].([]interface{})
	if len(routeAggregations) > 0 {
		for _, agg := range routeAggregations {
			data := agg.(map[string]interface{})
			prefix := data["prefix"].(string)
			summary := data["summary_only"].(bool)
			elem := model.RouteAggregationEntry{
				Prefix:      &prefix,
				SummaryOnly: &summary,
			}

			aggregationStructs = append(aggregationStructs, elem)
		}
	}

	restartTimerStruct := model.BgpGracefulRestartTimer{
		RestartTimer:    &restartTimer,
		StaleRouteTimer: &staleTimer,
	}

	restartConfigStruct := model.BgpGracefulRestartConfig{
		Mode:  &restartMode,
		Timer: &restartTimerStruct,
	}

	id := "bgp"
	bgpcType := "BgpRoutingConfig"
	routeStruct := model.BgpRoutingConfig{
		Ecmp:              &ecmp,
		Enabled:           &enabled,
		RouteAggregations: aggregationStructs,
		ResourceType:      &bgpcType,
		Tags:              tagStructs,
		Id:                &id,
		Revision:          &revision,
	}

	if !isVrf {
		// backend complains if the below config appears on VRF gateway.
		// We print a warning if property differs from default
		if interSrIbgp == false {
			log.Printf("[WARNING] BGP setting inter_sr_ibgp is not applicable for VRF gateway %s, and will be ignored", gwID)
		}
		if localAsNum != policyBGPLocalAsNumDefault {
			log.Printf("[WARNING] BGP setting local_as_num is not applicable for VRF gateway %s, and will be ignored", gwID)
		}
		if multipathRelax == false {
			log.Printf("[WARNING] BGP setting multipath_relax is not applicable for VRF gateway %s, and will be ignored", gwID)
		}
		if (restartMode != model.BgpGracefulRestartConfig_MODE_HELPER_ONLY) || (restartTimer != int64(policyBGPGracefulRestartStaleRouteTimerDefault)) || (staleTimer != int64(policyBGPGracefulRestartStaleRouteTimerDefault)) {
			log.Printf("[WARNING] BGP graceful restart settings are not applicable for VRF gateway %s, and will be ignored", gwID)
		}
		routeStruct.InterSrIbgp = &interSrIbgp
		routeStruct.LocalAsNum = &localAsNum
		routeStruct.MultipathRelax = &multipathRelax
		routeStruct.GracefulRestartConfig = &restartConfigStruct
	}

	return routeStruct
}

func resourceNsxtPolicyTier0GatewayResourceToInfraStruct(d *schema.ResourceData, connector *client.RestConnector, isCreate bool, id string) (model.Infra, error) {
	var infraChildren, gwChildren, lsChildren []*data.StructValue
	var infraStruct model.Infra
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	failoverMode := d.Get("failover_mode").(string)
	defaultRuleLogging := d.Get("default_rule_logging").(bool)
	disableFirewall := !d.Get("enable_firewall").(bool)
	forceWhitelisting := d.Get("force_whitelisting").(bool)
	haMode := d.Get("ha_mode").(string)
	revision := int64(d.Get("revision").(int))
	internalSubnets := interfaceListToStringList(d.Get("internal_transit_subnets").([]interface{}))
	transitSubnets := interfaceListToStringList(d.Get("transit_subnets").([]interface{}))
	ipv6ProfilePaths := getIpv6ProfilePathsFromSchema(d)
	vrfConfig := getPolicyVRFConfigFromSchema(d)
	dhcpPath := d.Get("dhcp_config_path").(string)

	t0Type := "Tier0"
	t0Struct := model.Tier0{
		DisplayName:            &displayName,
		Description:            &description,
		Tags:                   tags,
		FailoverMode:           &failoverMode,
		DefaultRuleLogging:     &defaultRuleLogging,
		DisableFirewall:        &disableFirewall,
		ForceWhitelisting:      &forceWhitelisting,
		Ipv6ProfilePaths:       ipv6ProfilePaths,
		HaMode:                 &haMode,
		InternalTransitSubnets: internalSubnets,
		TransitSubnets:         transitSubnets,
		ResourceType:           &t0Type,
		Id:                     &id,
		VrfConfig:              vrfConfig,
	}

	if !isCreate {
		t0Struct.Revision = &revision
	}
	if dhcpPath != "" {
		dhcpPaths := []string{dhcpPath}
		t0Struct.DhcpConfigPaths = dhcpPaths
	} else {
		t0Struct.DhcpConfigPaths = []string{}
	}

	bgpConfig := d.Get("bgp_config").([]interface{})
	if len(bgpConfig) > 0 {
		routingConfigStruct := resourceNsxtPolicyTier0GatewayBGPConfigSchemaToStruct(bgpConfig[0], vrfConfig != nil, id)
		childConfig := model.ChildBgpRoutingConfig{
			ResourceType:     "ChildBgpRoutingConfig",
			BgpRoutingConfig: &routingConfigStruct,
		}
		dataValue, errors := converter.ConvertToVapi(childConfig, model.ChildBgpRoutingConfigBindingType())
		if errors != nil {
			return infraStruct, fmt.Errorf("Error converting child BGP Routing Configuration: %v", errors[0])
		}
		lsChildren = append(lsChildren, dataValue.(*data.StructValue))
	}

	if len(lsChildren) > 0 || d.Get("edge_cluster_path") != "" {
		if d.Get("edge_cluster_path") == "" {
			bgpMap := bgpConfig[0].(map[string]interface{})
			if bgpMap["enabled"].(bool) {
				// BGP requires edge cluster
				// TODO: validate at plan time once multi-attribute validation is supported
				return infraStruct, fmt.Errorf("A valid edge_cluster_path is required when BGP is enabled")
			}
		}

		var serviceStruct *model.LocaleServices
		var err error
		if !isCreate {
			// BGP config and edge cluster share the same locale service
			serviceStruct, err = resourceNsxtPolicyTier0GatewayGetLocaleServiceEntry(d.Id(), connector)
			if err != nil {
				return infraStruct, err
			}
		}

		lsType := "LocaleServices"
		if serviceStruct == nil {
			// Locale Service required for edge cluster path and/or BGP config
			serviceStruct = &model.LocaleServices{
				Id: &defaultPolicyLocaleServiceID,
			}
		}

		serviceStruct.ResourceType = &lsType

		edgeClusterPath := d.Get("edge_cluster_path").(string)
		serviceStruct.EdgeClusterPath = &edgeClusterPath

		log.Printf("[DEBUG] Using Locale Service with ID %s and Edge Cluster %v", *serviceStruct.Id, serviceStruct.EdgeClusterPath)

		if len(lsChildren) > 0 {
			serviceStruct.Children = lsChildren
		}

		childService := model.ChildLocaleServices{
			ResourceType:   "ChildLocaleServices",
			LocaleServices: serviceStruct,
		}
		dataValue, errors := converter.ConvertToVapi(childService, model.ChildLocaleServicesBindingType())
		if errors != nil {
			return infraStruct, fmt.Errorf("Error converting child Locale Service: %v", errors[0])
		}
		gwChildren = append(gwChildren, dataValue.(*data.StructValue))
	}

	t0Struct.Children = gwChildren
	childTier0 := model.ChildTier0{
		Tier0:        &t0Struct,
		ResourceType: "ChildTier0",
	}
	dataValue, errors := converter.ConvertToVapi(childTier0, model.ChildTier0BindingType())
	if errors != nil {
		return infraStruct, fmt.Errorf("Error converting Tier0 Child: %v", errors[0])
	}
	infraChildren = append(infraChildren, dataValue.(*data.StructValue))

	infraType := "Infra"
	infraStruct = model.Infra{
		Children:     infraChildren,
		ResourceType: &infraType,
	}

	return infraStruct, nil
}

func resourceNsxtPolicyTier0GatewayCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	infraClient := nsx_policy.NewDefaultInfraClient(connector)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, connector, resourceNsxtPolicyTier0GatewayExists)
	if err != nil {
		return err
	}

	infraModel, err := resourceNsxtPolicyTier0GatewayResourceToInfraStruct(d, connector, true, id)
	if err != nil {
		return handleCreateError("Tier0", id, err)
	}

	boolFalse := false
	log.Printf("[INFO] H-API Creating Tier0 with ID %s", id)
	err = infraClient.Patch(infraModel, &boolFalse)
	if err != nil {
		return handleCreateError("Tier0", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyTier0GatewayRead(d, m)
}

func resourceNsxtPolicyTier0GatewayRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultTier0sClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Tier0 ID")
	}

	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "Tier0", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("failover_mode", obj.FailoverMode)
	d.Set("default_rule_logging", obj.DefaultRuleLogging)
	d.Set("enable_firewall", !(*obj.DisableFirewall))
	d.Set("ha_mode", obj.HaMode)
	d.Set("force_whitelisting", obj.ForceWhitelisting)
	d.Set("internal_transit_subnets", obj.InternalTransitSubnets)
	d.Set("transit_subnets", obj.TransitSubnets)
	d.Set("revision", obj.Revision)
	setPolicyVRFConfigInSchema(d, obj.VrfConfig)

	dhcpPaths := obj.DhcpConfigPaths
	if len(dhcpPaths) > 0 {
		d.Set("dhcp_config_path", dhcpPaths[0])
	}
	// Get the edge cluster Id
	localeService, err := resourceNsxtPolicyTier0GatewayGetLocaleServiceEntry(d.Id(), connector)
	if err != nil {
		return handleReadError(d, "Locale Service for T0", id, err)
	}
	if localeService != nil {
		d.Set("edge_cluster_path", localeService.EdgeClusterPath)
		err = resourceNsxtPolicyTier0GatewayReadBGPConfig(d, connector, *localeService)
		if err != nil {
			return handleReadError(d, "BGP Configuration for T0", id, err)
		}
	} else {
		// set empty bgp_config to keep empty plan
		d.Set("bgp_config", make([]map[string]interface{}, 0))
	}

	err = setIpv6ProfilePathsInSchema(d, obj.Ipv6ProfilePaths)
	if err != nil {
		return fmt.Errorf("Failed to get Tier0 %s ipv6 profiles: %v", *obj.Id, err)
	}

	return nil
}

func resourceNsxtPolicyTier0GatewayUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	infraClient := nsx_policy.NewDefaultInfraClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Tier0 ID")
	}

	infraModel, err := resourceNsxtPolicyTier0GatewayResourceToInfraStruct(d, connector, false, id)
	if err != nil {
		return err
	}

	boolFalse := false
	log.Printf("[INFO] H-API Updating Tier0 with ID %s", id)
	err = infraClient.Patch(infraModel, &boolFalse)
	if err != nil {
		return handleUpdateError("Tier0", id, err)
	}

	return resourceNsxtPolicyTier0GatewayRead(d, m)
}

func resourceNsxtPolicyTier0GatewayDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	var infraChildren []*data.StructValue
	infraClient := nsx_policy.NewDefaultInfraClient(connector)
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)
	boolTrue := true
	boolFalse := false

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Tier0 ID")
	}

	t0Type := "Tier0"
	t0obj := model.Tier0{
		Id:           &id,
		ResourceType: &t0Type,
	}

	childT0 := model.ChildTier0{
		MarkedForDelete: &boolTrue,
		Tier0:           &t0obj,
		ResourceType:    "ChildTier0",
	}
	dataValue, errors := converter.ConvertToVapi(childT0, model.ChildTier0BindingType())
	if errors != nil {
		return fmt.Errorf("Error converting Child Tier0: %v", errors[0])
	}
	infraChildren = append(infraChildren, dataValue.(*data.StructValue))

	infraType := "Infra"
	infraModel := model.Infra{
		Children:     infraChildren,
		ResourceType: &infraType,
	}

	log.Printf("[DEBUG] HAPI deleting Tier0 with ID %s", id)
	err := infraClient.Patch(infraModel, &boolFalse)
	if err != nil {
		return handleDeleteError("Tier0", id, err)
	}

	return nil

}
