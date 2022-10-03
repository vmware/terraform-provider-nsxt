/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"

	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/locale_services/ipsec_vpn_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

const policyBasedIPSecVpnSession string = "PolicyBased"
const routeBasedIPSecVpnSession string = "RouteBased"

var IPSecVpnSessionResourceType = []string{
	policyBasedIPSecVpnSession,
	routeBasedIPSecVpnSession,
}

var IPSecVpnSessionAuthenticationMode = []string{
	model.IPSecVpnSession_AUTHENTICATION_MODE_PSK,
	model.IPSecVpnSession_AUTHENTICATION_MODE_CERTIFICATE,
}

var IPSecVpnSessionConnectionInitiationMode = []string{
	model.IPSecVpnSession_CONNECTION_INITIATION_MODE_INITIATOR,
	model.IPSecVpnSession_CONNECTION_INITIATION_MODE_RESPOND_ONLY,
	model.IPSecVpnSession_CONNECTION_INITIATION_MODE_ON_DEMAND,
}
var IPSecVpnSessionComplianceSuite = []string{
	model.IPSecVpnSession_COMPLIANCE_SUITE_CNSA,
	model.IPSecVpnSession_COMPLIANCE_SUITE_SUITE_B_GCM_128,
	model.IPSecVpnSession_COMPLIANCE_SUITE_SUITE_B_GCM_256,
	model.IPSecVpnSession_COMPLIANCE_SUITE_PRIME,
	model.IPSecVpnSession_COMPLIANCE_SUITE_FOUNDATION,
	model.IPSecVpnSession_COMPLIANCE_SUITE_FIPS,
	model.IPSecVpnSession_COMPLIANCE_SUITE_NONE,
}

var IPSecRulesActionValues = []string{
	model.IPSecVpnRule_ACTION_PROTECT,
	model.IPSecVpnRule_ACTION_BYPASS,
}

func resourceNsxtPolicyIPSecVpnSession() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIPSecVpnSessionCreate,
		Read:   resourceNsxtPolicyIPSecVpnSessionRead,
		Update: resourceNsxtPolicyIPSecVpnSessionUpdate,
		Delete: resourceNsxtPolicyIPSecVpnSessionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nsx_id":              getNsxIDSchema(),
			"path":                getPathSchema(),
			"display_name":        getDisplayNameSchema(),
			"description":         getDescriptionSchema(),
			"revision":            getRevisionSchema(),
			"tag":                 getTagsSchema(),
			"tunnel_profile_path": getPolicyPathSchema(true, false, "Policy path referencing Tunnel profile to be used."),
			"local_endpoint_path": getPolicyPathSchema(true, false, "Policy path referencing Local endpoint."),
			"ike_profile_path":    getPolicyPathSchema(false, false, "Policy path referencing Ike profile."),
			"dpd_profile_path":    getPolicyPathSchema(false, false, "Policy path referencing dpd profile."),
			"vpn_type": {
				Type:         schema.TypeString,
				Description:  "A Policy Based VPN requires to define protect rules that match local and peer subnets. IPSec security associations is negotiated for each pair of local and peer subnet. A Route Based VPN is more flexible, more powerful and recommended over policy based VPN. IP Tunnel port is created and all traffic routed via tunnel port is protected. Routes can be configured statically or can be learned through BGP. A route based VPN is must for establishing redundant VPN session to remote site.",
				ValidateFunc: validation.StringInSlice(IPSecVpnSessionResourceType, false),
				Required:     true,
			},
			"compliance_suite": {
				Type:         schema.TypeString,
				Description:  "Compliance suite.",
				ValidateFunc: validation.StringInSlice(IPSecVpnSessionComplianceSuite, false),
				Optional:     true,
				Default:      model.IPSecVpnSession_COMPLIANCE_SUITE_NONE,
			},
			"connection_initiation_mode": {
				Type:         schema.TypeString,
				Description:  "Connection initiation mode used by local endpoint to establish ike connection with peer site. INITIATOR - In this mode local endpoint initiates tunnel setup and will also respond to incoming tunnel setup requests from peer gateway. RESPOND_ONLY - In this mode, local endpoint shall only respond to incoming tunnel setup requests. It shall not initiate the tunnel setup. ON_DEMAND - In this mode local endpoint will initiate tunnel creation once first packet matching the policy rule is received and will also respond to incoming initiation request.",
				ValidateFunc: validation.StringInSlice(IPSecVpnSessionConnectionInitiationMode, false),
				Optional:     true,
				Default:      model.IPSecVpnSession_CONNECTION_INITIATION_MODE_INITIATOR,
			},
			"authentication_mode": {
				Type:         schema.TypeString,
				Description:  "Peer authentication mode. PSK - In this mode a secret key shared between local and peer sites is to be used for authentication. The secret key can be a string with a maximum length of 128 characters. CERTIFICATE - In this mode a certificate defined at the global level is to be used for authentication.",
				ValidateFunc: validation.StringInSlice(IPSecVpnSessionAuthenticationMode, false),
				Optional:     true,
				Default:      model.IPSecVpnSession_AUTHENTICATION_MODE_PSK,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Enable/Disable IPSec VPN session.",
				Optional:    true,
				Default:     true,
			},
			"psk": {
				Type:        schema.TypeString,
				Description: "IPSec Pre-shared key. Maximum length of this field is 128 characters.",
				Optional:    true,
				Sensitive:   true,
			},
			"peer_id": {
				Type:         schema.TypeString,
				Description:  "Peer ID to uniquely identify the peer site. The peer ID is the public IP address of the remote device terminating the VPN tunnel. When NAT is configured for the peer, enter the private IP address of the peer.",
				Required:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"peer_address": {
				Type:         schema.TypeString,
				Description:  "Public IPV4 address of the remote device terminating the VPN connection.",
				Required:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"service_path": getPolicyPathSchema(true, true, "Policy path for IPSec VPN service"),
			"ip_addresses": {
				Type:        schema.TypeList,
				Description: "IP Tunnel interface (commonly referred as VTI) ip addresses.",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateSingleIP(),
				},
			},
			"rule": getIPSecVPNRulesSchema(),
			"prefix_length": {
				Type:         schema.TypeInt,
				Description:  "Subnet Prefix Length.",
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 255),
			},
			"direction": {
				Type:         schema.TypeString,
				Description:  "The traffic direction apply to the MSS clamping",
				Optional:     true,
				ValidateFunc: validation.StringInSlice(L2VpnSessionTCPSegmentClampingDirection, false),
			},
			"max_segment_size": {
				Type:         schema.TypeInt,
				Description:  "Maximum amount of data the host will accept in a Tcp segment.",
				Optional:     true,
				ValidateFunc: validation.IntBetween(108, 8860),
			},
		},
	}
}

func getIPSecVPNSessionFromSchema(d *schema.ResourceData) (*data.StructValue, error) {
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

	psk := d.Get("psk").(string)
	peerID := d.Get("peer_id").(string)
	peerAddress := d.Get("peer_address").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	ikeProfilePath := d.Get("ike_profile_path").(string)
	resourceType := d.Get("vpn_type").(string)
	localEndpointPath := d.Get("local_endpoint_path").(string)
	dpdProfilePath := d.Get("dpd_profile_path").(string)
	tunnelProfilePath := d.Get("tunnel_profile_path").(string)
	connectionInitiationMode := d.Get("connection_initiation_mode").(string)
	authenticationMode := d.Get("authentication_mode").(string)
	complianceSuite := d.Get("compliance_suite").(string)
	prefixLengh := int64(d.Get("prefix_length").(int))
	enabled := d.Get("enabled").(bool)
	direction := d.Get("direction").(string)
	mss := int64(d.Get("max_segment_size").(int))
	tags := getPolicyTagsFromSchema(d)

	if resourceType == routeBasedIPSecVpnSession {
		tunnelInterface := interfaceListToStringList(d.Get("ip_addresses").([]interface{}))
		var ipSubnets []model.TunnelInterfaceIPSubnet
		ipSubnet := model.TunnelInterfaceIPSubnet{
			IpAddresses:  tunnelInterface,
			PrefixLength: &prefixLengh,
		}
		ipSubnets = append(ipSubnets, ipSubnet)
		var vtiList []model.IPSecVpnTunnelInterface

		vti := model.IPSecVpnTunnelInterface{
			IpSubnets:   ipSubnets,
			DisplayName: &displayName,
		}

		vtiList = append(vtiList, vti)
		rsType := model.IPSecVpnSession_RESOURCE_TYPE_ROUTEBASEDIPSECVPNSESSION
		routeObj := model.RouteBasedIPSecVpnSession{
			DisplayName:              &displayName,
			Description:              &description,
			ConnectionInitiationMode: &connectionInitiationMode,
			ComplianceSuite:          &complianceSuite,
			AuthenticationMode:       &authenticationMode,
			ResourceType:             rsType,
			Enabled:                  &enabled,
			TunnelInterfaces:         vtiList,
			PeerAddress:              &peerAddress,
			PeerId:                   &peerID,
			Psk:                      &psk,
			Tags:                     tags,
		}
		if nsxVersionHigherOrEqual("3.2.0") {
			if direction != "" {
				tcpMSSClamping := model.TcpMaximumSegmentSizeClamping{
					Direction:      &direction,
					MaxSegmentSize: &mss,
				}
				routeObj.TcpMssClamping = &tcpMSSClamping
			}
		}
		if len(ikeProfilePath) > 0 {
			routeObj.IkeProfilePath = &ikeProfilePath
		}
		if len(localEndpointPath) > 0 {
			routeObj.LocalEndpointPath = &localEndpointPath
		}
		if len(tunnelProfilePath) > 0 {
			routeObj.TunnelProfilePath = &tunnelProfilePath
		}
		if len(dpdProfilePath) > 0 {
			routeObj.DpdProfilePath = &dpdProfilePath
		}

		dataValue, err := converter.ConvertToVapi(routeObj, model.RouteBasedIPSecVpnSessionBindingType())
		if err != nil {
			return nil, err[0]
		}
		return dataValue.(*data.StructValue), nil
	}

	if resourceType == policyBasedIPSecVpnSession {
		rsType := model.IPSecVpnSession_RESOURCE_TYPE_POLICYBASEDIPSECVPNSESSION
		ipSecVpnRules := getIPSecVPNRulesFromSchema(d)
		policyObj := model.PolicyBasedIPSecVpnSession{
			DisplayName:              &displayName,
			Description:              &description,
			ConnectionInitiationMode: &connectionInitiationMode,
			ComplianceSuite:          &complianceSuite,
			AuthenticationMode:       &authenticationMode,
			ResourceType:             rsType,
			Enabled:                  &enabled,
			Rules:                    ipSecVpnRules,
			PeerAddress:              &peerAddress,
			PeerId:                   &peerID,
			Psk:                      &psk,
			Tags:                     tags,
		}
		if nsxVersionHigherOrEqual("3.2.0") {
			if direction != "" {
				tcpMSSClamping := model.TcpMaximumSegmentSizeClamping{
					Direction:      &direction,
					MaxSegmentSize: &mss,
				}
				policyObj.TcpMssClamping = &tcpMSSClamping
			}
		}
		if len(ikeProfilePath) > 0 {
			policyObj.IkeProfilePath = &ikeProfilePath
		}
		if len(localEndpointPath) > 0 {
			policyObj.LocalEndpointPath = &localEndpointPath
		}
		if len(tunnelProfilePath) > 0 {
			policyObj.TunnelProfilePath = &tunnelProfilePath
		}
		if len(dpdProfilePath) > 0 {
			policyObj.DpdProfilePath = &dpdProfilePath
		}
		dataValue, err := converter.ConvertToVapi(policyObj, model.PolicyBasedIPSecVpnSessionBindingType())
		if err != nil {
			return nil, err[0]
		}
		return dataValue.(*data.StructValue), nil

	}

	return nil, fmt.Errorf("Failed to get IPSecVPNSession from Schema")
}

func getIPSecVPNRulesSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "For policy-based IPsec VPNs, a security policy specifies as its action the VPN tunnel to be used for transit traffic that meets the policy’s match criteria.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"nsx_id": getNsxIDSchema(),
				"sources": {
					Type:        schema.TypeSet,
					Description: "List of local subnets. Specifying no value is interpreted as 0.0.0.0/0.",
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validateCidr(),
					},
					Required: true,
				},
				"destinations": {
					Type:        schema.TypeSet,
					Description: "List of remote subnets used in policy-based L3Vpn.",
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validateCidr(),
					},
					Required: true,
				},
				"action": {
					Type:         schema.TypeString,
					Description:  "PROTECT - Protect rules are defined per policy based IPSec VPN session. BYPASS - Bypass rules are defined per IPSec VPN service and affects all policy based IPSec VPN sessions. Bypass rules are prioritized over protect rules.",
					Default:      model.IPSecVpnRule_ACTION_PROTECT,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(IPSecRulesActionValues, false),
				},
			},
		},
	}
}

func parseIPSecVPNServicePolicyPath(path string) (bool, string, string, string, error) {
	segs := strings.Split(path, "/")
	// Path should be like /infra/tier-1s/aaa/locale-services/default/ipsec-vpn-services/ccc
	segCount := len(segs)
	if (segCount < 8) || (segs[segCount-2] != "ipsec-vpn-services") {
		// error - this is not a segment path
		return false, "", "", "", fmt.Errorf("Invalid IPSec VPN service path %s", path)
	}

	serviceID := segs[segCount-1]
	localeServiceID := segs[5]
	gwPath := strings.Join(segs[:4], "/")

	isT0, gwID := parseGatewayPolicyPath(gwPath)
	return isT0, gwID, localeServiceID, serviceID, nil
}

func setRuleInSchema(d *schema.ResourceData, bypassRules []model.IPSecVpnRule) {
	var ruleList []map[string]interface{}
	for _, rule := range bypassRules {
		elem := make(map[string]interface{})
		var srcList []string
		for _, src := range rule.Sources {
			srcList = append(srcList, *src.Subnet)
		}
		var destList []string
		for _, dest := range rule.Destinations {
			destList = append(destList, *dest.Subnet)
		}
		nsxID := rule.Id
		elem["nsx_id"] = nsxID
		elem["sources"] = srcList
		elem["destinations"] = destList
		elem["action"] = rule.Action
		ruleList = append(ruleList, elem)
	}
	err := d.Set("rule", ruleList)
	if err != nil {
		log.Printf("[WARNING] Failed to set bypass_rules in schema: %v", err)
	}
}

func getIPSecVPNRulesFromSchema(d *schema.ResourceData) []model.IPSecVpnRule {
	rules := d.Get("rule")
	if rules != nil {
		rules := rules.([]interface{})
		var ruleList []model.IPSecVpnRule
		for _, rule := range rules {
			data := rule.(map[string]interface{})
			action := data["action"].(string)
			sourceRanges := interface2StringList(data["sources"].(*schema.Set).List())
			destinationRanges := interface2StringList(data["destinations"].(*schema.Set).List())

			/// Source Subnets

			SourceIPSecVpnSubnetList := make([]model.IPSecVpnSubnet, 0)
			if len(sourceRanges) > 0 {
				for _, element := range sourceRanges {
					subnet := element
					IPSecVpnSubnet := model.IPSecVpnSubnet{
						Subnet: &subnet,
					}
					SourceIPSecVpnSubnetList = append(SourceIPSecVpnSubnetList, IPSecVpnSubnet)
				}
			}

			/// Destination Subnets
			DestinationIPSecVpnSubnetList := make([]model.IPSecVpnSubnet, 0)
			if len(destinationRanges) > 0 {
				for _, element := range destinationRanges {
					subnet := element
					IPSecVpnSubnet := model.IPSecVpnSubnet{
						Subnet: &subnet,
					}
					DestinationIPSecVpnSubnetList = append(DestinationIPSecVpnSubnetList, IPSecVpnSubnet)
				}
			}
			ruleID := data["nsx_id"].(string)
			if ruleID == "" {
				ruleID = newUUID()
			}
			elem := model.IPSecVpnRule{
				Action:       &action,
				Sources:      SourceIPSecVpnSubnetList,
				Destinations: DestinationIPSecVpnSubnetList,
				UniqueId:     &ruleID,
				Id:           &ruleID,
			}
			ruleList = append(ruleList, elem)
		}
		return ruleList
	}
	return nil
}

func resourceNsxtPolicyIPSecVpnSessionCreate(d *schema.ResourceData, m interface{}) error {

	servicePath := d.Get("service_path").(string)
	isT0, gwID, localeServiceID, serviceID, err := parseIPSecVPNServicePolicyPath(servicePath)
	if err != nil {
		return err
	}
	if !isT0 {
		return fmt.Errorf("VPN is supported only on Tier-0 with ACTIVE-STANDBY HA mode")
	}
	id := d.Get("nsx_id").(string)
	if id == "" {
		id = newUUID()
	}

	connector := getPolicyConnector(m)

	_, err1 := resourceNsxtPolicyIPSecVpnSessionExists(gwID, localeServiceID, serviceID, id, connector)
	if err1 == nil {
		return fmt.Errorf("IPSecVpnSession with nsx_id '%s' already exists under IPSecVpnService '%s' of localeService '%s'", id, serviceID, localeServiceID)
	} else if !isNotFoundError(err1) {
		return err1
	}

	obj, err2 := getIPSecVPNSessionFromSchema(d)
	if err2 != nil {
		return err2
	}

	client := ipsec_vpn_services.NewSessionsClient(connector)

	err3 := client.Patch(gwID, localeServiceID, serviceID, id, obj)

	if err3 != nil {
		return handleCreateError("IPSecVpnSession", id, err3)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyIPSecVpnSessionRead(d, m)
}

func resourceNsxtPolicyIPSecVpnSessionExists(tier0ID string, localeServiceID string, serviceID string, sessionID string, connector *client.RestConnector) (bool, error) {
	client := ipsec_vpn_services.NewSessionsClient(connector)
	_, err := client.Get(tier0ID, localeServiceID, serviceID, sessionID)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, err
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyIPSecVpnSessionRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

	id := d.Id()
	resourceType := d.Get("vpn_type").(string)
	if id == "" {
		return fmt.Errorf("Error obtaining IPSecVpnSession ID")
	}

	servicePath := d.Get("service_path").(string)
	isT0, gwID, localeServiceID, serviceID, err := parseIPSecVPNServicePolicyPath(servicePath)
	if err != nil {
		return err
	}
	if !isT0 {
		return fmt.Errorf("VPN is supported only on Tier-0 with ACTIVE-STANDBY HA mode")
	}

	client := ipsec_vpn_services.NewSessionsClient(connector)

	obj, err1 := client.Get(gwID, localeServiceID, serviceID, id)

	if err1 != nil {
		if isNotFoundError(err1) {
			d.SetId("")
			return nil
		}
		return handleReadError(d, "VPN Session", id, err1)
	}

	if resourceType == routeBasedIPSecVpnSession {
		interfaceVpn, errs := converter.ConvertToGolang(obj, model.RouteBasedIPSecVpnSessionBindingType())
		if len(errs) > 0 {
			return fmt.Errorf("Error converting VPN Session %s", errs[0])
		}
		blockVPN := interfaceVpn.(model.RouteBasedIPSecVpnSession)

		d.Set("display_name", blockVPN.DisplayName)
		d.Set("description", blockVPN.Description)
		setPolicyTagsInSchema(d, blockVPN.Tags)
		d.Set("nsx_id", blockVPN.Id)
		d.Set("path", blockVPN.Path)
		d.Set("revision", blockVPN.Revision)
		d.Set("vpn_type", routeBasedIPSecVpnSession)
		d.Set("authentication_mode", blockVPN.AuthenticationMode)
		d.Set("compliance_suite", blockVPN.ComplianceSuite)
		d.Set("connection_initiation_mode", blockVPN.ConnectionInitiationMode)
		d.Set("enabled", blockVPN.Enabled)
		d.Set("ike_profile_path", blockVPN.IkeProfilePath)
		d.Set("local_endpoint_path", blockVPN.LocalEndpointPath)
		d.Set("dpd_profile_path", blockVPN.DpdProfilePath)
		d.Set("tunnel_profile_path", blockVPN.TunnelProfilePath)
		d.Set("peer_address", blockVPN.PeerAddress)
		d.Set("peer_id", blockVPN.PeerId)
		if nsxVersionHigherOrEqual("3.2.0") {
			if blockVPN.TcpMssClamping != nil {
				direction := blockVPN.TcpMssClamping.Direction
				mss := blockVPN.TcpMssClamping.MaxSegmentSize
				d.Set("direction", direction)
				d.Set("max_segment_size", mss)
			}
		}
		var subnets []string
		var prefixLength int64
		if blockVPN.TunnelInterfaces != nil && len(blockVPN.TunnelInterfaces) > 0 {
			for _, tunnelInterface := range blockVPN.TunnelInterfaces {
				ipSubnets := tunnelInterface.IpSubnets
				for _, ipSubnet := range ipSubnets {
					ipAddresses := ipSubnet.IpAddresses
					if len(ipAddresses) > 0 {
						subnets = append(subnets, ipAddresses...)
					}
					if prefixLength == 0 {
						prefixLength = *ipSubnet.PrefixLength
					}
				}
			}
		}
		d.Set("prefix_length", prefixLength)
		if len(subnets) > 0 {
			d.Set("ip_addresses", subnets)
		}

	} else {
		// Policy based ipsec vpn session
		interfaceVpn, errs := converter.ConvertToGolang(obj, model.PolicyBasedIPSecVpnSessionBindingType())
		if len(errs) > 0 {
			return fmt.Errorf("Error converting VPN Session %s", errs[0])
		}
		blockVPN := interfaceVpn.(model.PolicyBasedIPSecVpnSession)

		d.Set("display_name", blockVPN.DisplayName)
		d.Set("description", blockVPN.Description)
		setPolicyTagsInSchema(d, blockVPN.Tags)
		d.Set("nsx_id", blockVPN.Id)
		d.Set("path", blockVPN.Path)
		d.Set("revision", blockVPN.Revision)
		d.Set("vpn_type", policyBasedIPSecVpnSession)
		d.Set("authentication_mode", blockVPN.AuthenticationMode)
		d.Set("compliance_suite", blockVPN.ComplianceSuite)
		d.Set("connection_initiation_mode", blockVPN.ConnectionInitiationMode)
		d.Set("enabled", blockVPN.Enabled)
		d.Set("ike_profile_path", blockVPN.IkeProfilePath)
		d.Set("local_endpoint_path", blockVPN.LocalEndpointPath)
		d.Set("dpd_profile_path", blockVPN.DpdProfilePath)
		d.Set("tunnel_profile_path", blockVPN.TunnelProfilePath)
		if blockVPN.Rules != nil {
			setRuleInSchema(d, blockVPN.Rules)
		}
		if nsxVersionHigherOrEqual("3.2.0") {
			if blockVPN.TcpMssClamping != nil {
				direction := blockVPN.TcpMssClamping.Direction
				mss := blockVPN.TcpMssClamping.MaxSegmentSize
				d.Set("direction", direction)
				d.Set("max_segment_size", mss)
			}
		}
	}
	return nil
}

func resourceNsxtPolicyIPSecVpnSessionUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPSecVpnSession ID")
	}

	servicePath := d.Get("service_path").(string)
	isT0, gwID, localeServiceID, serviceID, err := parseIPSecVPNServicePolicyPath(servicePath)
	if err != nil {
		return err
	}
	if !isT0 {
		return fmt.Errorf("VPN is supported only on Tier-0 with ACTIVE-STANDBY HA mode")
	}
	obj, err := getIPSecVPNSessionFromSchema(d)
	if err != nil {
		return err
	}

	client := ipsec_vpn_services.NewSessionsClient(connector)

	err = client.Patch(gwID, localeServiceID, serviceID, id, obj)

	if err != nil {
		return handleUpdateError("IPSecVpnSession", id, err)
	}
	d.SetId(id)
	d.Set("nsx_id", id)
	return resourceNsxtPolicyIPSecVpnSessionRead(d, m)

}

func resourceNsxtPolicyIPSecVpnSessionDelete(d *schema.ResourceData, m interface{}) error {

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPSecVpnSession ID")
	}
	servicePath := d.Get("service_path").(string)
	isT0, gwID, localeServiceID, serviceID, err := parseIPSecVPNServicePolicyPath(servicePath)
	if err != nil {
		return err
	}
	if !isT0 {
		return fmt.Errorf("VPN is supported only on Tier-0 with ACTIVE-STANDBY HA mode")
	}

	connector := getPolicyConnector(m)

	var err1 error
	client := ipsec_vpn_services.NewSessionsClient(connector)
	err1 = client.Delete(gwID, localeServiceID, serviceID, id)

	if err1 != nil {
		return handleDeleteError("IPSecVpnSession", id, err1)
	}

	return nil
}
