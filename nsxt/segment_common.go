package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/segments"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"log"
	"time"
)

var connectivityValues = []string{
	model.SegmentAdvancedConfig_CONNECTIVITY_ON,
	model.SegmentAdvancedConfig_CONNECTIVITY_OFF,
}

func getPolicySegmentDhcpV4ConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"server_address": {
				Type:        schema.TypeString,
				Description: "IP address of the DHCP server in CIDR format",
				Optional:    true,
				// TODO: validate IPv4 only
				ValidateFunc: validateIPCidr(),
			},
			"dns_servers": {
				Type:        schema.TypeList,
				Description: "IP addresses of DNS servers for subnet",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateSingleIP(),
				},
				Optional: true,
			},
			"lease_time": {
				Type:         schema.TypeInt,
				Description:  "DHCP lease time in seconds",
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(60),
			},
			"dhcp_option_121":     getDhcpOptions121Schema(),
			"dhcp_generic_option": getDhcpGenericOptionsSchema(),
		},
	}
}

func getPolicySegmentDhcpV6ConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"server_address": {
				Type:        schema.TypeString,
				Description: "IP address of the DHCP server in CIDR format",
				Optional:    true,
				// TODO: validate IPv6 only
				ValidateFunc: validateIPCidr(),
			},
			"dns_servers": {
				Type:        schema.TypeList,
				Description: "IP addresses of DNS servers for subnet",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateSingleIP(),
				},
				Optional: true,
			},
			"lease_time": {
				Type:         schema.TypeInt,
				Description:  "DHCP lease time in seconds",
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(60),
			},
			"domain_names": {
				Type:        schema.TypeList,
				Description: "Domain names for subnet",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"excluded_range": getAllocationRangeListSchema(false, "Excluded addresses to define dynamic ip allocation ranges"),
			"preferred_time": {
				Type:         schema.TypeInt,
				Description:  "The time interval in seconds, in which the prefix is advertised as preferred",
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(60),
			},
			"sntp_servers": {
				Type:        schema.TypeList,
				Description: "IPv6 address of SNTP servers for subnet",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateSingleIP(),
				},
				Optional: true,
			},
			// TODO: add options
		},
	}
}

func getPolicySegmentSubnetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"dhcp_v4_config": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          getPolicySegmentDhcpV4ConfigSchema(),
				MaxItems:      1,
				ConflictsWith: []string{"subnet.dhcp_v6_config"},
			},
			"dhcp_v6_config": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          getPolicySegmentDhcpV6ConfigSchema(),
				MaxItems:      1,
				ConflictsWith: []string{"subnet.dhcp_v4_config"},
			},
			"dhcp_ranges": {
				Type:        schema.TypeList,
				Description: "DHCP address ranges for dynamic IP allocation",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCidrOrIPOrRange(),
				},
				Optional: true,
			},
			"cidr": {
				Type:         schema.TypeString,
				Description:  "Gateway IP address in CIDR format",
				Required:     true,
				ValidateFunc: validateIPCidr(),
			},
			"network": {
				Type:        schema.TypeString,
				Description: "Network CIDR for subnet",
				Computed:    true,
			},
		},
	}
}

func getPolicySegmentL2ExtensionConfigurationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"l2vpn_paths": {
				Type:        schema.TypeList,
				Description: "Policy paths of associated L2 VPN sessions",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validatePolicyPath(),
				},
				Optional: true,
			},
			"tunnel_id": {
				Type:         schema.TypeInt,
				Description:  "Tunnel ID",
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 4093),
			},
		},
	}
}

func getPolicySegmentAdvancedConfigurationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"address_pool_path": {
				Type:         schema.TypeString,
				Description:  "Policy path to IP address pool",
				Optional:     true,
				ValidateFunc: validatePolicyPath(),
			},
			"connectivity": {
				Type:         schema.TypeString,
				Description:  "Connectivity configuration to manually connect (ON) or disconnect (OFF)",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(connectivityValues, false),
			},
			"hybrid": {
				Type:        schema.TypeBool,
				Description: "Flag to identify a hybrid logical switch",
				Optional:    true,
				Default:     false,
			},
			"local_egress": {
				Type:        schema.TypeBool,
				Description: "Flag to enable local egress",
				Optional:    true,
				Default:     false,
			},
			"uplink_teaming_policy": {
				Type:        schema.TypeString,
				Description: "The name of the switching uplink teaming policy for the bridge endpoint",
				Optional:    true,
			},
		},
	}
}

func getPolicyCommonSegmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"nsx_id":       getNsxIDSchema(),
		"path":         getPathSchema(),
		"display_name": getDisplayNameSchema(),
		"description":  getDescriptionSchema(),
		"revision":     getRevisionSchema(),
		"tag":          getTagsSchema(),
		"advanced_config": {
			Type:        schema.TypeList,
			Description: "Advanced segment configuration",
			Elem:        getPolicySegmentAdvancedConfigurationSchema(),
			Optional:    true,
			MaxItems:    1,
		},
		"connectivity_path": {
			Type:         schema.TypeString,
			Description:  "Policy path to the connecting Tier-0 or Tier-1",
			Optional:     true,
			ValidateFunc: validatePolicyPath(),
		},
		"domain_name": {
			Type:        schema.TypeString,
			Description: "DNS domain names",
			Optional:    true,
		},
		"l2_extension": {
			Type:        schema.TypeList,
			Description: "Configuration for extending Segment through L2 VPN",
			Elem:        getPolicySegmentL2ExtensionConfigurationSchema(),
			Optional:    true,
			MaxItems:    1,
		},
		"overlay_id": {
			Type:        schema.TypeInt,
			Description: "Overlay connectivity ID for this Segment",
			Optional:    true,
			Computed:    true,
		},
		"subnet": {
			Type:        schema.TypeList,
			Description: "Subnet configuration with at most 1 IPv4 CIDR and multiple IPv6 CIDRs",
			Elem:        getPolicySegmentSubnetSchema(),
			Optional:    true,
		},
		"dhcp_config_path": getPolicyPathSchema(false, false, "Policy path to DHCP server or relay configuration to use for subnets configured on this segment"),
		"transport_zone_path": {
			Type:         schema.TypeString,
			Description:  "Policy path to the transport zone",
			Required:     true,
			ValidateFunc: validatePolicyPath(),
		},
		"vlan_ids": {
			Type:        schema.TypeList,
			Description: "VLAN IDs for VLAN backed Segment",
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateVLANId,
			},
			Required: true,
		},
	}
}

func getPolicyDhcpOptions121(opts []interface{}) model.DhcpOption121 {
	var opt121Struct model.DhcpOption121
	var routes []model.ClasslessStaticRoute
	for _, opt121 := range opts {
		data := opt121.(map[string]interface{})
		network := data["network"].(string)
		nextHop := data["next_hop"].(string)
		elem := model.ClasslessStaticRoute{
			Network: &network,
			NextHop: &nextHop,
		}
		routes = append(routes, elem)
	}
	if len(routes) > 0 {
		opt121Struct = model.DhcpOption121{
			StaticRoutes: routes,
		}
	}
	return opt121Struct
}

func getPolicyDhcpOptions121FromStruct(opt *model.DhcpOption121) []map[string]interface{} {
	var dhcpOpt121 []map[string]interface{}
	for _, route := range opt.StaticRoutes {
		elem := make(map[string]interface{})
		elem["network"] = route.Network
		elem["next_hop"] = route.NextHop
		dhcpOpt121 = append(dhcpOpt121, elem)
	}
	return dhcpOpt121
}

func getPolicyDhcpGenericOptions(opts []interface{}) []model.GenericDhcpOption {
	var options []model.GenericDhcpOption
	for _, opt := range opts {
		data := opt.(map[string]interface{})
		code := int64(data["code"].(int))
		elem := model.GenericDhcpOption{
			Code:   &code,
			Values: interface2StringList(data["values"].([]interface{})),
		}
		options = append(options, elem)
	}
	return options
}

func getPolicyDhcpGenericOptionsFromStruct(opts []model.GenericDhcpOption) []map[string]interface{} {
	var dhcpOptions []map[string]interface{}
	for _, opt := range opts {
		elem := make(map[string]interface{})
		elem["code"] = opt.Code
		elem["values"] = opt.Values
		dhcpOptions = append(dhcpOptions, elem)
	}
	return dhcpOptions
}

func getSegmentSubnetDhcpConfigFromSchema(schemaConfig map[string]interface{}) (*data.StructValue, error) {
	if nsxVersionLower("3.0.0") {
		return nil, nil
	}

	dhcpV4Config := schemaConfig["dhcp_v4_config"].([]interface{})
	dhcpV6Config := schemaConfig["dhcp_v6_config"].([]interface{})

	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

	if len(dhcpV4Config) > 0 {
		dhcpConfig := dhcpV4Config[0].(map[string]interface{})
		serverAddress := dhcpConfig["server_address"].(string)
		dnsServers := dhcpConfig["dns_servers"].([]interface{})
		leaseTime := int64(dhcpConfig["lease_time"].(int))

		config := model.SegmentDhcpV4Config{
			ResourceType: model.SegmentDhcpConfig_RESOURCE_TYPE_SEGMENTDHCPV4CONFIG,
			DnsServers:   interface2StringList(dnsServers),
		}

		if len(serverAddress) > 0 {
			config.ServerAddress = &serverAddress
		}

		dhcpOpts := model.DhcpV4Options{}

		dhcp121Opts := dhcpConfig["dhcp_option_121"].([]interface{})
		if len(dhcp121Opts) > 0 {
			dhcp121OptStruct := getPolicyDhcpOptions121(dhcp121Opts)
			dhcpOpts.Option121 = &dhcp121OptStruct
		}

		otherDhcpOpts := dhcpConfig["dhcp_generic_option"].([]interface{})
		if len(otherDhcpOpts) > 0 {
			otherOptStructs := getPolicyDhcpGenericOptions(otherDhcpOpts)
			dhcpOpts.Others = otherOptStructs
		}

		if len(dhcp121Opts)+len(otherDhcpOpts) > 0 {
			config.Options = &dhcpOpts
		}

		if leaseTime > 0 {
			config.LeaseTime = &leaseTime
		}

		dataValue, errs := converter.ConvertToVapi(config, model.SegmentDhcpV4ConfigBindingType())
		if errs != nil {
			return nil, errs[0]
		}

		return dataValue.(*data.StructValue), nil
	}

	if len(dhcpV6Config) > 0 {
		dhcpConfig := dhcpV6Config[0].(map[string]interface{})
		serverAddress := dhcpConfig["server_address"].(string)
		dnsServers := dhcpConfig["dns_servers"].([]interface{})
		sntpServers := dhcpConfig["sntp_servers"].([]interface{})
		domainNames := dhcpConfig["domain_names"].([]interface{})
		excludedRanges := dhcpConfig["excluded_range"].([]interface{})
		leaseTime := int64(dhcpConfig["lease_time"].(int))
		preferredTime := int64(dhcpConfig["preferred_time"].(int))

		config := model.SegmentDhcpV6Config{
			ResourceType: model.SegmentDhcpConfig_RESOURCE_TYPE_SEGMENTDHCPV6CONFIG,
			DnsServers:   interface2StringList(dnsServers),
		}

		if len(serverAddress) > 0 {
			config.ServerAddress = &serverAddress
		}

		if len(domainNames) > 0 {
			config.DomainNames = interface2StringList(domainNames)
		}

		if len(sntpServers) > 0 {
			config.SntpServers = interface2StringList(sntpServers)
		}

		if len(excludedRanges) > 0 {
			config.ExcludedRanges = interface2StringList(excludedRanges)
		}

		if leaseTime > 0 {
			config.LeaseTime = &leaseTime
		}

		if preferredTime > 0 {
			config.PreferredTime = &preferredTime
		}

		dataValue, errs := converter.ConvertToVapi(config, model.SegmentDhcpV6ConfigBindingType())
		if errs != nil {
			return nil, errs[0]
		}

		return dataValue.(*data.StructValue), nil
	}

	return nil, nil

}

func setSegmentSubnetDhcpConfigInSchema(schemaConfig map[string]interface{}, subnetConfig model.SegmentSubnet) error {
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

	resultConfig := make(map[string]interface{})

	if subnetConfig.DhcpConfig == nil {
		return nil
	}

	obj, errs := converter.ConvertToGolang(subnetConfig.DhcpConfig, model.SegmentDhcpConfigBindingType())
	if errs != nil {
		return errs[0]
	}

	resourceType := obj.(model.SegmentDhcpConfig).ResourceType
	if resourceType == model.SegmentDhcpConfig_RESOURCE_TYPE_SEGMENTDHCPV4CONFIG {
		obj, errs := converter.ConvertToGolang(subnetConfig.DhcpConfig, model.SegmentDhcpV4ConfigBindingType())
		if errs != nil {
			return errs[0]
		}
		dhcpV4Config := obj.(model.SegmentDhcpV4Config)
		resultConfig["server_address"] = dhcpV4Config.ServerAddress
		resultConfig["lease_time"] = dhcpV4Config.LeaseTime
		resultConfig["dns_servers"] = dhcpV4Config.DnsServers

		if dhcpV4Config.Options != nil {
			if dhcpV4Config.Options.Option121 != nil {
				opts := getPolicyDhcpOptions121FromStruct(dhcpV4Config.Options.Option121)
				resultConfig["dhcp_option_121"] = opts
			}
			if len(dhcpV4Config.Options.Others) > 0 {
				opts := getPolicyDhcpGenericOptionsFromStruct(dhcpV4Config.Options.Others)
				resultConfig["dhcp_generic_option"] = opts
			}
		}
		schemaConfig["dhcp_v4_config"] = resultConfig
		return nil
	}

	if resourceType == model.SegmentDhcpConfig_RESOURCE_TYPE_SEGMENTDHCPV6CONFIG {
		obj, errs := converter.ConvertToGolang(subnetConfig.DhcpConfig, model.SegmentDhcpV6ConfigBindingType())
		if errs != nil {
			return errs[0]
		}

		dhcpV6Config := obj.(model.SegmentDhcpV6Config)
		resultConfig["server_address"] = dhcpV6Config.ServerAddress
		resultConfig["lease_time"] = dhcpV6Config.LeaseTime
		resultConfig["preferred_time"] = dhcpV6Config.PreferredTime
		resultConfig["dns_servers"] = dhcpV6Config.DnsServers
		resultConfig["sntp_servers"] = dhcpV6Config.SntpServers
		resultConfig["domain_names"] = dhcpV6Config.DomainNames
		resultConfig["excluded_range"] = dhcpV6Config.ExcludedRanges

		schemaConfig["dhcp_v6_config"] = resultConfig
		return nil
	}

	return fmt.Errorf("Unrecognized DHCP Config Resource Type %s", resourceType)

}

func policySegmentResourceToStruct(d *schema.ResourceData, isVlan bool) (model.Segment, error) {
	// Read the rest of the configured parameters
	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)
	domainName := d.Get("domain_name").(string)
	tzPath := d.Get("transport_zone_path").(string)
	dhcpConfigPath := d.Get("dhcp_config_path").(string)
	revision := int64(d.Get("revision").(int))

	obj := model.Segment{
		DisplayName: &displayName,
		Tags:        tags,
		Revision:    &revision,
	}

	if description != "" {
		obj.Description = &description
	}
	if domainName != "" {
		obj.DomainName = &domainName
	}
	if tzPath != "" {
		obj.TransportZonePath = &tzPath
	}
	if dhcpConfigPath != "" && nsxVersionHigherOrEqual("3.0.0") {
		obj.DhcpConfigPath = &dhcpConfigPath
	}

	var vlanIds []string
	var subnets []interface{}
	var subnetStructs []model.SegmentSubnet
	if isVlan {
		// VLAN specific fields
		for _, vlanID := range d.Get("vlan_ids").([]interface{}) {
			vlanIds = append(vlanIds, vlanID.(string))
		}
		obj.VlanIds = vlanIds
	} else {
		// overlay specific fields
		connectivityPath := d.Get("connectivity_path").(string)
		overlayID, exists := d.GetOkExists("overlay_id")
		if exists {
			overlayID64 := int64(overlayID.(int))
			obj.OverlayId = &overlayID64
		}
		if connectivityPath != "" {
			obj.ConnectivityPath = &connectivityPath
		}
	}
	subnets = d.Get("subnet").([]interface{})
	if len(subnets) > 0 {
		for _, subnet := range subnets {
			subnetMap := subnet.(map[string]interface{})
			dhcpRanges := subnetMap["dhcp_ranges"].([]interface{})
			var dhcpRangeList []string
			if len(dhcpRanges) > 0 {
				for _, dhcpRange := range dhcpRanges {
					dhcpRangeList = append(dhcpRangeList, dhcpRange.(string))
				}
			}
			gwAddr := subnetMap["cidr"].(string)
			network := subnetMap["network"].(string)
			subnetStruct := model.SegmentSubnet{
				DhcpRanges:     dhcpRangeList,
				GatewayAddress: &gwAddr,
				Network:        &network,
			}
			config, err := getSegmentSubnetDhcpConfigFromSchema(subnetMap)
			if err != nil {
				return obj, err
			}

			subnetStruct.DhcpConfig = config

			subnetStructs = append(subnetStructs, subnetStruct)
		}
	}
	obj.Subnets = subnetStructs

	advConfig := d.Get("advanced_config").([]interface{})
	if len(advConfig) > 0 {
		advConfigMap := advConfig[0].(map[string]interface{})
		connectivity := advConfigMap["connectivity"].(string)
		hybrid := advConfigMap["hybrid"].(bool)
		egress := advConfigMap["local_egress"].(bool)
		var poolPaths []string
		if advConfigMap["cidr"] != nil {
			poolPaths = append(poolPaths, advConfigMap["cidr"].(string))
		}
		advConfigStruct := model.SegmentAdvancedConfig{
			AddressPoolPaths: poolPaths,
			Hybrid:           &hybrid,
			LocalEgress:      &egress,
		}

		if connectivity != "" {
			advConfigStruct.Connectivity = &connectivity
		}

		if nsxVersionHigherOrEqual("3.0.0") {
			teamingPolicy := advConfigMap["uplink_teaming_policy"].(string)
			if teamingPolicy != "" {
				advConfigStruct.UplinkTeamingPolicyName = &teamingPolicy
			}
		}
		obj.AdvancedConfig = &advConfigStruct
	}

	l2Ext := d.Get("l2_extension").([]interface{})
	if len(l2Ext) > 0 {
		l2ExtMap := l2Ext[0].(map[string]interface{})
		vpnPaths := interfaceListToStringList(l2ExtMap["l2vpn_paths"].([]interface{}))
		tunnelID := int64(l2ExtMap["tunnel_id"].(int))
		l2Struct := model.L2Extension{
			L2vpnPaths: vpnPaths,
			TunnelId:   &tunnelID,
		}
		obj.L2Extension = &l2Struct
	}

	return obj, nil
}

func resourceNsxtPolicySegmentExists(id string, connector *client.RestConnector) bool {
	client := infra.NewDefaultSegmentsClient(connector)

	_, err := client.Get(id)
	if err == nil {
		return true
	}

	if isNotFoundError(err) {
		return false
	}

	logAPIError("Error retrieving Segment", err)
	return false
}

func nsxtPolicySegmentRead(d *schema.ResourceData, m interface{}, isVlan bool) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultSegmentsClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Segment ID")
	}

	obj, err := client.Get(id)
	if err != nil {
		return handleReadError(d, "Segment", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("connectivity_path", obj.ConnectivityPath)
	d.Set("dhcp_config_path", obj.DhcpConfigPath)
	d.Set("domain_name", obj.DomainName)
	d.Set("transport_zone_path", obj.TransportZonePath)

	if isVlan {
		d.Set("vlan_ids", obj.VlanIds)
	} else {
		if obj.OverlayId != nil {
			d.Set("overlay_id", int(*obj.OverlayId))
		} else {
			d.Set("overlay_id", "")
		}
	}

	if obj.AdvancedConfig != nil {
		advConfig := make(map[string]interface{})
		poolPaths := obj.AdvancedConfig.AddressPoolPaths
		if len(poolPaths) > 0 {
			advConfig["address_pool_path"] = poolPaths[0]
		}
		advConfig["connectivity"] = obj.AdvancedConfig.Connectivity
		advConfig["hybrid"] = obj.AdvancedConfig.Hybrid
		advConfig["local_egress"] = obj.AdvancedConfig.LocalEgress
		if obj.AdvancedConfig.UplinkTeamingPolicyName != nil {
			advConfig["uplink_teaming_policy"] = *obj.AdvancedConfig.UplinkTeamingPolicyName
		}
		// This is a list with 1 element
		var advConfigList []map[string]interface{}
		advConfigList = append(advConfigList, advConfig)
		d.Set("advanced_config", advConfigList)
	}

	if obj.L2Extension != nil {
		l2Ext := make(map[string]interface{})
		l2Ext["l2vpn_paths"] = obj.L2Extension.L2vpnPaths
		l2Ext["tunnel_id"] = obj.L2Extension.TunnelId
		// This is a list with 1 element
		var l2ExtList []map[string]interface{}
		l2ExtList = append(l2ExtList, l2Ext)
		d.Set("l2_extension", l2ExtList)
	}

	var subnetSegments []interface{}
	for _, subnetSeg := range obj.Subnets {
		seg := make(map[string]interface{})
		seg["dhcp_ranges"] = subnetSeg.DhcpRanges
		seg["cidr"] = subnetSeg.GatewayAddress
		seg["network"] = subnetSeg.Network
		setSegmentSubnetDhcpConfigInSchema(seg, subnetSeg)
		subnetSegments = append(subnetSegments, seg)
	}

	d.Set("subnet", subnetSegments)

	return nil
}

func nsxtPolicySegmentCreate(d *schema.ResourceData, m interface{}, isVlan bool) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultSegmentsClient(connector)

	// Initialize resource Id and verify this ID is not yet used
	id, err := getOrGenerateID(d, connector, resourceNsxtPolicySegmentExists)
	if err != nil {
		return err
	}

	// Create the resource using PATCH
	log.Printf("[INFO] Creating Segment with ID %s", id)
	obj, err := policySegmentResourceToStruct(d, isVlan)
	if err != nil {
		return err
	}
	err = client.Patch(id, obj)
	if err != nil {
		return handleCreateError("Segment", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return nsxtPolicySegmentRead(d, m, isVlan)
}

func nsxtPolicySegmentUpdate(d *schema.ResourceData, m interface{}, isVlan bool) error {
	connector := getPolicyConnector(m)
	client := infra.NewDefaultSegmentsClient(connector)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Segment ID")
	}

	obj, err := policySegmentResourceToStruct(d, isVlan)
	if err != nil {
		return err
	}
	_, err = client.Update(id, obj)
	if err != nil {
		return handleUpdateError("Segment", id, err)
	}

	return nsxtPolicySegmentRead(d, m, isVlan)
}

func nsxtPolicySegmentDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining Segment ID")
	}

	connector := getPolicyConnector(m)
	client := infra.NewDefaultSegmentsClient(connector)
	portsClient := segments.NewDefaultPortsClient(connector)

	// During buld destroy, VMs might be destroyed before segments, but
	// VIF release is not yet propagated to NSX. NSX will reply with
	// InvalidRequest on attempted delete if ports are present on the
	// segment. The code below waits till possible ports are deleted.
	pendingStates := []string{"pending"}
	targetStates := []string{"ok", "error"}
	stateConf := &resource.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {
			ports, err := portsClient.List(id, nil, nil, nil, nil, nil, nil)
			if err != nil {
				return ports, "error", logAPIError("Error listing segment ports", err)
			}

			log.Printf("[DEBUG] Current number of ports on segment %s is %d", id, len(ports.Results))

			if len(ports.Results) > 0 {
				return ports, "pending", nil
			}
			return ports, "ok", nil

		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		MinTimeout: 1 * time.Second,
		Delay:      1 * time.Second,
	}
	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Failed to get port information for segment %s: %v", id, err)
	}
	err = client.Delete(id)
	if err != nil {
		return handleDeleteError("Segment", id, err)
	}

	return nil
}
