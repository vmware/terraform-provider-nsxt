package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt"
	global_policy "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
	"strings"
)

var nsxtPolicyTier0GatewayRedistributionRuleTypes = []string{
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER0_STATIC,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER0_CONNECTED,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER0_EXTERNAL_INTERFACE,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER0_SEGMENT,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER0_ROUTER_LINK,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER0_SERVICE_INTERFACE,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER0_LOOPBACK_INTERFACE,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER0_DNS_FORWARDER_IP,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER0_IPSEC_LOCAL_IP,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER0_NAT,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER0_EVPN_TEP_IP,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER1_NAT,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER1_STATIC,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER1_LB_VIP,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER1_LB_SNAT,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER1_DNS_FORWARDER_IP,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER1_CONNECTED,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER1_SERVICE_INTERFACE,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER1_SEGMENT,
	model.Tier0RouteRedistributionRule_ROUTE_REDISTRIBUTION_TYPES_TIER1_IPSEC_LOCAL_ENDPOINT,
}

func getFailoverModeSchema(defaultValue string) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "Failover mode",
		Default:      defaultValue,
		Optional:     true,
		ValidateFunc: validation.StringInSlice(policyFailOverModeValues, false),
	}
}

func getIPv6NDRAPathSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "The path of an IPv6 NDRA profile",
		Optional:    true,
		Computed:    true,
	}
}

func getIPv6DadPathSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "The path of an IPv6 DAD profile",
		Optional:    true,
		Computed:    true,
	}
}

func getPolicyEdgeClusterPathSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  "The path of the edge cluster connected to this gateway",
		Optional:     true,
		ValidateFunc: validatePolicyPath(),
		Computed:     true,
	}
}

func getPolicyLocaleServiceSchema(isTier1 bool) *schema.Schema {
	nodeConficts := []string{}
	if isTier1 {
		// for Tier1, enable_standby_relocation can not be enabled if
		// preferred nodes are specified
		nodeConficts = append(nodeConficts, "enable_standby_relocation")
	}

	elemSchema := map[string]*schema.Schema{
		"edge_cluster_path": {
			Type:         schema.TypeString,
			Description:  "The path of the edge cluster connected to this gateway",
			Required:     true,
			ValidateFunc: validatePolicyPath(),
		},
		"preferred_edge_paths": {
			Type:          schema.TypeSet,
			Description:   "Paths of specific edge nodes",
			Optional:      true,
			Elem:          getElemPolicyPathSchema(),
			ConflictsWith: nodeConficts,
		},
		"redistribution_config": getRedistributionConfigSchema(),
		"path":                  getPathSchema(),
		"revision":              getRevisionSchema(),
	}
	if isTier1 {
		delete(elemSchema, "redistribution_config")
	}

	result := &schema.Schema{
		Type:          schema.TypeSet,
		Optional:      true,
		Description:   "Locale Service for the gateway",
		ConflictsWith: []string{"edge_cluster_path"},
		Elem: &schema.Resource{
			Schema: elemSchema,
		},
	}

	return result
}

func getIpv6ProfilePathsFromSchema(d *schema.ResourceData) []string {
	var profiles []string
	if d.Get("ipv6_ndra_profile_path") != "" {
		profiles = append(profiles, d.Get("ipv6_ndra_profile_path").(string))
	}
	if d.Get("ipv6_dad_profile_path") != "" {
		profiles = append(profiles, d.Get("ipv6_dad_profile_path").(string))
	}
	return profiles
}

func setIpv6ProfilePathsInSchema(d *schema.ResourceData, paths []string) error {
	for _, path := range paths {
		if strings.HasPrefix(path, "/infra/ipv6-ndra-profiles") {
			d.Set("ipv6_ndra_profile_path", path)
		}
		if strings.HasPrefix(path, "/infra/ipv6-dad-profiles") {
			d.Set("ipv6_dad_profile_path", path)
		}
	}
	return nil
}

func getGatewayInterfaceSubnetsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of IP addresses and network prefixes for this interface",
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateIPCidr(),
		},
		Required: true,
	}
}

func getMtuSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Description:  "Maximum transmission unit specifies the size of the largest packet that a network protocol can transmit",
		ValidateFunc: validation.IntAtLeast(64),
	}
}

var gatewayInterfaceUrpfModeValues = []string{
	model.Tier0Interface_URPF_MODE_NONE,
	model.Tier0Interface_URPF_MODE_STRICT,
}

func getGatewayInterfaceUrpfModeSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "Unicast Reverse Path Forwarding mode",
		ValidateFunc: validation.StringInSlice(gatewayInterfaceUrpfModeValues, false),
		Default:      model.Tier0Interface_URPF_MODE_STRICT,
	}
}

func getGatewayIntersiteConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: "Locale Service for the gateway",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"transit_subnet": {
					Type:         schema.TypeString,
					Description:  "IPv4 subnet for inter-site transit segment connecting service routers across sites for stretched gateway. For IPv6 link local subnet is auto configured",
					Computed:     true,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validateCidr(),
				},
				"primary_site_path": {
					Type:         schema.TypeString,
					Description:  "Primary egress site for gateway",
					Optional:     true,
					ValidateFunc: validatePolicyPath(),
				},
				"fallback_site_paths": {
					Type:        schema.TypeSet,
					Description: "Fallback sites to be used as new primary site on current primary site failure",
					Optional:    true,
					Elem:        getElemPolicyPathSchema(),
				},
			},
		},
	}
}

func listPolicyGatewayLocaleServices(connector *client.RestConnector, gwID string, listLocaleServicesFunc func(*client.RestConnector, string, *string) (model.LocaleServicesListResult, error)) ([]model.LocaleServices, error) {
	var results []model.LocaleServices
	var cursor *string
	var count int64
	total := int64(0)

	for {
		listResponse, err := listLocaleServicesFunc(connector, gwID, cursor)
		if err != nil {
			return results, err
		}
		cursor = listResponse.Cursor
		count = *listResponse.ResultCount
		results = append(results, listResponse.Results...)
		if total == 0 {
			// first response
			total = count
		}
		if int64(len(results)) >= total {
			return results, nil
		}
	}
}

func getGlobalPolicyGatewayLocaleServiceIDWithSite(localeServices []model.LocaleServices, sitePath string, gatewayID string) (string, error) {
	if len(localeServices) == 0 {
		return "", fmt.Errorf("Edge cluster is mandatory on gateway %s in order to create interfaces", gatewayID)
	}
	for _, localeService := range localeServices {
		if localeService.EdgeClusterPath != nil {
			if strings.HasPrefix(*localeService.EdgeClusterPath, sitePath) {
				localeServiceID := *localeService.Id
				return localeServiceID, nil
			}
		}
	}
	return "", fmt.Errorf("Edge cluster is mandatory on gateway %s in order to create interfaces", gatewayID)
}

func initChildLocaleService(serviceStruct *model.LocaleServices, markForDelete bool) (*data.StructValue, error) {
	childService := model.ChildLocaleServices{
		ResourceType:    "ChildLocaleServices",
		LocaleServices:  serviceStruct,
		MarkedForDelete: &markForDelete,
	}

	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)
	dataValue, err := converter.ConvertToVapi(childService, model.ChildLocaleServicesBindingType())

	if err != nil {
		return nil, err[0]
	}

	return dataValue.(*data.StructValue), nil
}

func initGatewayLocaleServices(d *schema.ResourceData) ([]*data.StructValue, error) {
	var localeServices []*data.StructValue

	oldServices, newServices := d.GetChange("locale_service")

	existingServices := make(map[string]bool)
	if len(d.Id()) > 0 {
		for _, obj := range oldServices.(*schema.Set).List() {
			cfg := obj.(map[string]interface{})
			serviceID := getPolicyIDFromPath(cfg["path"].(string))
			existingServices[serviceID] = true
		}
	}
	lsType := "LocaleServices"
	for _, service := range newServices.(*schema.Set).List() {
		cfg := service.(map[string]interface{})
		edgeClusterPath := cfg["edge_cluster_path"].(string)
		edgeNodes := interface2StringList(cfg["preferred_edge_paths"].(*schema.Set).List())
		path := cfg["path"].(string)
		revision := int64(cfg["revision"].(int))

		var serviceID string
		if path != "" {
			serviceID = getPolicyIDFromPath(path)
		} else {
			serviceID = newUUID()
			log.Printf("[DEBUG] Preparing to create locale service %s for gateway %s", serviceID, d.Id())
		}
		serviceStruct := model.LocaleServices{
			Id:                 &serviceID,
			ResourceType:       &lsType,
			EdgeClusterPath:    &edgeClusterPath,
			PreferredEdgePaths: edgeNodes,
		}

		redistributionConfigs := cfg["redistribution_config"]
		if redistributionConfigs != nil {
			setLocaleServiceRedistributionConfig(redistributionConfigs.([]interface{}), &serviceStruct)
		}

		if _, ok := existingServices[serviceID]; ok {
			// if this is an update for existing locale service,
			// we need revision
			serviceStruct.Revision = &revision
		}
		/*
			bgpConfig := cfg["bgp_config"].([]interface{})
			if len(bgpConfig) > 0 {
				var lsChildren []*data.StructValue
				// For Global Manager BGP is defined under locale services
				routingConfigStruct := resourceNsxtPolicyTier0GatewayBGPConfigSchemaToStruct(bgpConfig[0], false, d.Id())
				structValue, err := initPolicyTier0ChildBgpConfig(&routingConfigStruct)
				if err != nil {
					return localeServices, err
				}
				lsChildren = append(lsChildren, structValue)
				serviceStruct.Children = lsChildren
			}
		*/
		dataValue, err := initChildLocaleService(&serviceStruct, false)
		if err != nil {
			return localeServices, err
		}

		localeServices = append(localeServices, dataValue)
		existingServices[serviceID] = false
	}
	// Add instruction to delete services that are no longer present in intent
	for id, shouldDelete := range existingServices {
		if shouldDelete {
			serviceStruct := model.LocaleServices{
				Id:           &id,
				ResourceType: &lsType,
			}
			dataValue, err := initChildLocaleService(&serviceStruct, true)
			if err != nil {
				return localeServices, err
			}
			localeServices = append(localeServices, dataValue)
			log.Printf("[DEBUG] Preparing to delete locale service %s for gateway %s", id, d.Id())
		}
	}

	return localeServices, nil
}

func getPolicyGatewayIntersiteConfigFromSchema(d *schema.ResourceData) *model.IntersiteGatewayConfig {
	cfg, isSet := d.GetOk("intersite_config")
	if !isSet {
		return nil
	}

	configs := cfg.([]interface{})
	for _, elem := range configs {
		data := elem.(map[string]interface{})

		subnet := data["transit_subnet"].(string)
		primarySitePath := data["primary_site_path"].(string)
		fallbackSites := interface2StringList(data["fallback_site_paths"].(*schema.Set).List())
		intersiteConfig := model.IntersiteGatewayConfig{
			PrimarySitePath: &primarySitePath,
			FallbackSites:   fallbackSites,
		}
		if len(subnet) > 0 {
			intersiteConfig.IntersiteTransitSubnet = &subnet
		}

		return &intersiteConfig
	}

	return nil
}

func setPolicyGatewayIntersiteConfigInSchema(d *schema.ResourceData, config *model.IntersiteGatewayConfig) error {
	if config == nil {
		return nil
	}

	var result []map[string]interface{}
	elem := make(map[string]interface{})

	if config.IntersiteTransitSubnet != nil {
		elem["transit_subnet"] = config.IntersiteTransitSubnet
	}

	if config.PrimarySitePath != nil {
		elem["primary_site_path"] = config.PrimarySitePath
	}

	elem["fallback_site_paths"] = config.FallbackSites
	result = append(result, elem)

	return d.Set("intersite_config", result)
}

func policyInfraPatch(obj model.Infra, isGlobalManager bool, connector *client.RestConnector, enforceRevision bool) error {
	if isGlobalManager {
		infraClient := global_policy.NewDefaultGlobalInfraClient(connector)
		gmObj, err := convertModelBindingType(obj, model.InfraBindingType(), gm_model.InfraBindingType())
		if err != nil {
			return err
		}

		return infraClient.Patch(gmObj.(gm_model.Infra), &enforceRevision)
	}

	infraClient := nsx_policy.NewDefaultInfraClient(connector)
	return infraClient.Patch(obj, &enforceRevision)
}

func getRedistributionConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Route Redistribution configuration",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:        schema.TypeBool,
					Description: "Flag to enable route redistribution for BGP",
					Optional:    true,
					Default:     true,
				},
				"rule": {
					Type:        schema.TypeList,
					Description: "List of routes to be aggregated",
					Optional:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:        schema.TypeString,
								Description: "Rule name",
								Optional:    true,
							},
							"route_map_path": {
								Type:        schema.TypeString,
								Description: "Route map to be associated with the redistribution rule",
								Optional:    true,
							},
							"types": {
								Type:        schema.TypeSet,
								Description: "List of redistribution types",
								Optional:    true,
								Elem: &schema.Schema{
									Type:         schema.TypeString,
									ValidateFunc: validation.StringInSlice(nsxtPolicyTier0GatewayRedistributionRuleTypes, false),
								},
							},
						},
					},
				},
			},
		},
	}
}

func setLocaleServiceRedistributionConfig(redistributionConfigs []interface{}, serviceStruct *model.LocaleServices) {
	var rules []model.Tier0RouteRedistributionRule
	if len(redistributionConfigs) == 0 {
		return
	}

	redistributionConfig := redistributionConfigs[0].(map[string]interface{})
	bgp := redistributionConfig["enabled"].(bool)
	rulesConfig := redistributionConfig["rule"].([]interface{})

	for _, ruleConfig := range rulesConfig {
		data := ruleConfig.(map[string]interface{})
		name := data["name"].(string)
		routeMapPath := data["route_map_path"].(string)
		types := data["types"].(*schema.Set).List()

		rule := model.Tier0RouteRedistributionRule{
			RouteRedistributionTypes: interface2StringList(types),
		}

		if len(name) > 0 {
			rule.Name = &name
		}

		if len(routeMapPath) > 0 {
			rule.RouteMapPath = &routeMapPath
		}

		rules = append(rules, rule)
	}

	redistributionStruct := model.Tier0RouteRedistributionConfig{
		BgpEnabled:          &bgp,
		RedistributionRules: rules,
	}

	serviceStruct.RouteRedistributionConfig = &redistributionStruct
}

func getLocaleServiceRedistributionConfig(serviceStruct *model.LocaleServices) []map[string]interface{} {
	var redistributionConfigs []map[string]interface{}
	config := serviceStruct.RouteRedistributionConfig
	if config == nil {
		return redistributionConfigs
	}

	var rules []map[string]interface{}
	elem := make(map[string]interface{})
	elem["enabled"] = config.BgpEnabled
	for _, ruleConfig := range config.RedistributionRules {
		rule := make(map[string]interface{})
		rule["name"] = ruleConfig.Name
		rule["route_map_path"] = ruleConfig.RouteMapPath
		rule["types"] = ruleConfig.RouteRedistributionTypes
		rules = append(rules, rule)
	}

	elem["rule"] = rules
	redistributionConfigs = append(redistributionConfigs, elem)
	return redistributionConfigs
}

func findTier0LocaleServiceForSite(connector *client.RestConnector, gwID string, sitePath string) (string, error) {
	localeServices, err := listPolicyTier0GatewayLocaleServices(connector, gwID, true)
	if err != nil {
		return "", err
	}
	return getGlobalPolicyGatewayLocaleServiceIDWithSite(localeServices, sitePath, gwID)
}
