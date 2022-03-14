/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"log"

	//"strconv"
	//"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"

	//"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	//"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/locale_services/bgp"
	ipsec_vpn_services "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/locale_services/ipsec_vpn_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var IPSecVpnSessionResourceType = []string{
	model.IPSecVpnSession_RESOURCE_TYPE_POLICYBASEDIPSECVPNSESSION,
	model.IPSecVpnSession_RESOURCE_TYPE_ROUTEBASEDIPSECVPNSESSION,
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
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"vpn_type": {
				Type:         schema.TypeString,
				Description:  "A Policy Based VPN requires to define protect rules that match local and peer subnets. IPSec security associations is negotiated for each pair of local and peer subnet. A Route Based VPN is more flexible, more powerful and recommended over policy based VPN. IP Tunnel port is created and all traffic routed via tunnel port is protected. Routes can be configured statically or can be learned through BGP. A route based VPN is must for establishing redundant VPN session to remote site.",
				ValidateFunc: validation.StringInSlice(IPSecVpnSessionResourceType, false),
				Optional:     true,
			},
			"compliance_suite": {
				Type:         schema.TypeString,
				Description:  "Connection initiation mode used by local endpoint to establish ike connection with peer site. INITIATOR - In this mode local endpoint initiates tunnel setup and will also respond to incoming tunnel setup requests from peer gateway. RESPOND_ONLY - In this mode, local endpoint shall only respond to incoming tunnel setup requests. It shall not initiate the tunnel setup. ON_DEMAND - In this mode local endpoint will initiate tunnel creation once first packet matching the policy rule is received and will also respond to incoming initiation request.",
				ValidateFunc: validation.StringInSlice(IPSecVpnSessionComplianceSuite, false),
				Optional:     true,
			},
			"connection_initiation_mode": {
				Type:         schema.TypeString,
				Description:  "Compliance suite.",
				ValidateFunc: validation.StringInSlice(IPSecVpnSessionConnectionInitiationMode, false),
				Optional:     true,
				Default:      "INITIATOR",
			},
			"authentication_mode": {
				Type:         schema.TypeString,
				Description:  "Peer authentication mode. PSK - In this mode a secret key shared between local and peer sites is to be used for authentication. The secret key can be a string with a maximum length of 128 characters. CERTIFICATE - In this mode a certificate defined at the global level is to be used for authentication.",
				ValidateFunc: validation.StringInSlice(IPSecVpnSessionAuthenticationMode, false),
				Optional:     true,
				Default:      "PSK",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Enable/Disable IPSec VPN session.",
				Optional:    true,
				Default:     true,
			},
			"dpd_profile_path": {
				Type:        schema.TypeString,
				Description: "Policy path referencing Dead Peer Detection (DPD) profile. Default is set to system default profile.",
				Optional:    true,
				Default:     "/infra/ipsec-vpn-dpd-profiles/nsx-default-l3vpn-dpd-profile",
			},
			"psk": {
				Type:        schema.TypeString,
				Description: "IPSec Pre-shared key. Maximum length of this field is 128 characters.",
				Optional:    true,
			},
			"peer_id": {
				Type:        schema.TypeString,
				Description: "Peer ID to uniquely identify the peer site. The peer ID is the public IP address of the remote device terminating the VPN tunnel. When NAT is configured for the peer, enter the private IP address of the peer.",
				Optional:    true,
			},
			"peer_address": {
				Type:        schema.TypeString,
				Description: "Public IPV4 address of the remote device terminating the VPN connection.",
				Optional:    true,
			},
			"tunnel_profile_path": {
				Type:        schema.TypeString,
				Description: "Policy path referencing Tunnel profile to be used. Default is set to system default profile.",
				Required:    true,
			},
			"local_endpoint_path": {
				Type:        schema.TypeString,
				Description: "Policy path referencing Local endpoint.",
				Optional:    true,
			},
			"ike_profile_path": {
				Type:        schema.TypeString,
				Description: "Policy path referencing IKE profile.",
				Optional:    true,
			},
			"tier0_id": {
				Type:        schema.TypeString,
				Description: "Unique identifier of the T0 resource.",
				Optional:    true,
				Default:     "vmc",
			},
			"locale_service": {
				Type:        schema.TypeString,
				Description: "Unique identifier of the Locale Service resource.",
				Optional:    true,
				Default:     "default",
			},
			"service_id": {
				Type:        schema.TypeString,
				Description: "Unique identifier of the Service.",
				Optional:    true,
				Default:     "default",
			},
			"subnets": {
				Type:        schema.TypeList,
				Description: "IP Tunnel interface (commonly referred as VTI) subnet.",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateSingleIP(),
				},
			},
			/*
				"sources": {
					Type:        schema.TypeList,
					Description: "List of local subnets. Specifying no value is interpreted as 0.0.0.0/0.",
					Optional:    true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validateCidrOrIPOrRange(),
					},
				},
				"destinations": {
					Type:        schema.TypeList,
					Description: "List of remote subnets used in policy-based L3Vpn.",
					Optional:    true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validateCidrOrIPOrRange(),
					},
				},
			*/
			"rule": getIPSecVPNRulesSchema(),
			"prefix_length": {
				Type:         schema.TypeInt,
				Description:  "Subnet Prefix Length.",
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 255),
			},
		},
	}
}

func getIPSecVPNSessionFromSchema(d *schema.ResourceData) (*data.StructValue, error) {
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

	Psk := d.Get("psk").(string)
	PeerID := d.Get("peer_id").(string)
	PeerAddress := d.Get("peer_address").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	IkeProfilePath := d.Get("ike_profile_path").(string)
	ResourceType := d.Get("vpn_type").(string)
	LocalEndpointPath := d.Get("local_endpoint_path").(string)
	//var LocalEndpointPath string
	//if LocalEndpointVar == "Private" {
	//	LocalEndpointPath = "/infra/tier-0s/vmc/locale-services/default/ipsec-vpn-services/default/local-endpoints/Private-IP1"
	// else {
	//	LocalEndpointPath = "/infra/tier-0s/vmc/locale-services/default/ipsec-vpn-services/default/local-endpoints/Public-IP1"
	//}
	DpdProfilePath := d.Get("dpd_profile_path").(string)
	TunnelProfilePath := d.Get("tunnel_profile_path").(string)
	ConnectionInitiationMode := d.Get("connection_initiation_mode").(string)
	AuthenticationMode := d.Get("authentication_mode").(string)
	ComplianceSuite := d.Get("compliance_suite").(string)
	PrefixLengh := int64(d.Get("prefix_length").(int))
	Enabled := d.Get("enabled").(bool)

	if ResourceType == "RouteBasedIPSecVpnSession" {
		TunnelInterface := interfaceListToStringList(d.Get("subnets").([]interface{}))
		var IPSubnets []model.TunnelInterfaceIPSubnet
		IPSubnet := model.TunnelInterfaceIPSubnet{
			IpAddresses:  TunnelInterface,
			PrefixLength: &PrefixLengh,
		}
		IPSubnets = append(IPSubnets, IPSubnet)
		var VTIlist []model.IPSecVpnTunnelInterface

		vti := model.IPSecVpnTunnelInterface{
			IpSubnets:   IPSubnets,
			DisplayName: &displayName,
		}

		VTIlist = append(VTIlist, vti)

		routeObj := model.RouteBasedIPSecVpnSession{
			DisplayName:              &displayName,
			Description:              &description,
			IkeProfilePath:           &IkeProfilePath,
			LocalEndpointPath:        &LocalEndpointPath,
			TunnelProfilePath:        &TunnelProfilePath,
			DpdProfilePath:           &DpdProfilePath,
			ConnectionInitiationMode: &ConnectionInitiationMode,
			ComplianceSuite:          &ComplianceSuite,
			AuthenticationMode:       &AuthenticationMode,
			ResourceType:             ResourceType,
			Enabled:                  &Enabled,
			TunnelInterfaces:         VTIlist,
			PeerAddress:              &PeerAddress,
			PeerId:                   &PeerID,
			Psk:                      &Psk,
		}

		dataValue, err := converter.ConvertToVapi(routeObj, model.RouteBasedIPSecVpnSessionBindingType())
		if err != nil {
			return nil, err[0]
		}
		return dataValue.(*data.StructValue), nil
	}

	if ResourceType == "PolicyBasedIPSecVpnSession" {
		log.Println("#################################################1")
		IPSecVpnRules := getIPSecVPNRulesFromSchema(d)
		log.Println("#################################################10")
		/*

			var exampleSubnet string
			exampleSubnet = "10.10.10.0/24"
			log.Println("#################################################")
			log.Println(exampleSubnet)
			DestinationsObj := model.IPSecVpnSubnet{
				Subnet: &exampleSubnet,
			}
			log.Println(DestinationsObj)
			var DestinationsObjects []model.IPSecVpnSubnet
			DestinationsObjects = append(DestinationsObjects, DestinationsObj)
			var sourceSubnet string
			sourceSubnet = "10.10.11.0/24"
			SourcesObj := model.IPSecVpnSubnet{
				Subnet: &sourceSubnet,
			}
			var SourcesObjects []model.IPSecVpnSubnet
			SourcesObjects = append(SourcesObjects, SourcesObj)

			rule_id := newUUID()
			var resourceType string
			resourceType = "IPSecVpnRule"
			//var action string
			//action = "PROTECT"
			IPSecVpnRule := model.IPSecVpnRule{
				Sources:      SourcesObjects,
				Destinations: DestinationsObjects,
				UniqueId:     &rule_id,
				Id:           &rule_id,
				ResourceType: &resourceType,
				//Action:       &action,
			}
			log.Println(IPSecVpnRule)
			var IPSecVpnRules []model.IPSecVpnRule
			IPSecVpnRules = append(IPSecVpnRules, IPSecVpnRule)
		*/
		policyObj := model.PolicyBasedIPSecVpnSession{
			DisplayName:              &displayName,
			Description:              &description,
			IkeProfilePath:           &IkeProfilePath,
			LocalEndpointPath:        &LocalEndpointPath,
			TunnelProfilePath:        &TunnelProfilePath,
			DpdProfilePath:           &DpdProfilePath,
			ConnectionInitiationMode: &ConnectionInitiationMode,
			ComplianceSuite:          &ComplianceSuite,
			AuthenticationMode:       &AuthenticationMode,
			ResourceType:             ResourceType,
			Enabled:                  &Enabled,
			Rules:                    IPSecVpnRules,
			PeerAddress:              &PeerAddress,
			PeerId:                   &PeerID,
			Psk:                      &Psk,
		}
		dataValue, err := converter.ConvertToVapi(policyObj, model.PolicyBasedIPSecVpnSessionBindingType())
		if err != nil {
			return nil, err[0]
		}
		return dataValue.(*data.StructValue), nil

	}

	var err error
	return nil, err
}

func getIPSecVPNRulesSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "For policy-based IPsec VPNs, a security policy specifies as its action the VPN tunnel to be used for transit traffic that meets the policy’s match criteria.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
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

//func getSubnetsFromSchema(d *schema.ResourceData) []manager.IpPoolSubnet {
//subnets := d.Get("subnet").([]interface{})
//var subnetsList []manager.IpPoolSubnet
//for _, subnet := range subnets {
//	data := subnet.(map[string]interface{})
//	elem := manager.IpPoolSubnet{
//		Cidr:             data["cidr"].(string),
//		DnsSuffix:        data["dns_suffix"].(string),
//		GatewayIp:        data["gateway_ip"].(string),
//		DnsNameservers:   interface2StringList(data["dns_nameservers"].([]interface{})),
//		AllocationRanges: getAllocationRangesFromRanges(data["allocation_ranges"].([]interface{})),
//	}
//
//	subnetsList = append(subnetsList, elem)
//}
//	return subnetsList
//}

func getSourcesSubnetsFromSchema(d *schema.ResourceData) []model.IPSecVpnSubnet {
	IPSecVpnSubnets := d.Get("sources").(*schema.Set).List()
	IPSecVpnSubnetList := make([]model.IPSecVpnSubnet, 0)
	for _, subnet := range IPSecVpnSubnets {
		data := subnet.(map[string]interface{})
		subnet := data["sources"].(string)
		elem := model.IPSecVpnSubnet{
			Subnet: &subnet,
		}
		IPSecVpnSubnetList = append(IPSecVpnSubnetList, elem)
	}
	return IPSecVpnSubnetList
}

func getIPSecVPNRulesFromSchema(d *schema.ResourceData) []model.IPSecVpnRule {
	rules := d.Get("rule").([]interface{})
	log.Println("#################################################2")
	var ruleList []model.IPSecVpnRule
	for _, rule := range rules {
		data := rule.(map[string]interface{})
		log.Println("#################################################3")
		action := data["action"].(string)
		log.Println(action)
		log.Println("#################################################4")
		sourceRanges := interface2StringList(data["sources"].(*schema.Set).List())
		destinationRanges := interface2StringList(data["destinations"].(*schema.Set).List())

		/// Source Subnets

		log.Println("#################################################5")
		SourceIPSecVpnSubnetList := make([]model.IPSecVpnSubnet, 0)
		if len(sourceRanges) > 0 {
			for _, element := range sourceRanges {
				subnet := element
				log.Println("#################################################6")
				log.Println(subnet)
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
				log.Println("#################################################6")
				log.Println(subnet)
				IPSecVpnSubnet := model.IPSecVpnSubnet{
					Subnet: &subnet,
				}
				DestinationIPSecVpnSubnetList = append(DestinationIPSecVpnSubnetList, IPSecVpnSubnet)
			}
		}

		rule_id := newUUID()
		log.Println("#################################################")
		elem := model.IPSecVpnRule{
			Action:       &action,
			Sources:      SourceIPSecVpnSubnetList,
			Destinations: DestinationIPSecVpnSubnetList,
			UniqueId:     &rule_id,
			Id:           &rule_id,
		}
		ruleList = append(ruleList, elem)
		log.Println("#################################################")
	}
	return ruleList
}

func resourceNsxtPolicyIPSecVpnSessionCreate(d *schema.ResourceData, m interface{}) error {

	Tier0ID := d.Get("tier0_id").(string)
	LocaleService := d.Get("locale_service").(string)
	ServiceID := d.Get("service_id").(string)

	// Initialize resource Id and verify this ID is not yet used
	id := newUUID()

	connector := getPolicyConnector(m)

	obj, err := getIPSecVPNSessionFromSchema(d)
	if err != nil {
		return err
	}

	client := ipsec_vpn_services.NewDefaultSessionsClient(connector)

	// Create the resource using PATCH
	log.Printf("[INFO] Creating IPSecVpnSession with ID %s", id)

	err2 := client.Patch(Tier0ID, LocaleService, ServiceID, id, obj)

	if err2 != nil {
		return handleCreateError("IPSecVpnSession", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyIPSecVpnSessionRead(d, m)
}

func resourceNsxtPolicyIPSecVpnSessionRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPSecVpnSession ID")
	}

	Tier0ID := d.Get("tier0_id").(string)
	LocaleService := d.Get("locale_service").(string)
	ServiceID := d.Get("service_id").(string)

	client := ipsec_vpn_services.NewDefaultSessionsClient(connector)

	obj, err := client.Get(Tier0ID, LocaleService, ServiceID, id)

	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			log.Printf("[DEBUG] VPNSession %s not found", id)
			return nil
		}
		return handleReadError(d, "VPN Session", id, err)
	}

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

	return nil
}

func resourceNsxtPolicyIPSecVpnSessionUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPSecVpnSession ID")
	}

	Tier0ID := d.Get("tier0_id").(string)
	LocaleService := d.Get("locale_service").(string)
	ServiceID := d.Get("service_id").(string)

	obj, err := getIPSecVPNSessionFromSchema(d)
	if err != nil {
		return err
	}

	client := ipsec_vpn_services.NewDefaultSessionsClient(connector)

	// Create the resource using PATCH
	log.Printf("[INFO] Creating IPSecVpnSession with ID %s", id)

	err2 := client.Patch(Tier0ID, LocaleService, ServiceID, id, obj)

	if err2 != nil {
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
	Tier0ID := d.Get("tier0_id").(string)
	LocaleService := d.Get("locale_service").(string)
	ServiceID := d.Get("service_id").(string)

	connector := getPolicyConnector(m)

	var err error
	client := ipsec_vpn_services.NewDefaultSessionsClient(connector)
	err = client.Delete(Tier0ID, LocaleService, ServiceID, id)

	if err != nil {
		return handleDeleteError("IPSecVpnSession", id, err)
	}

	return nil
}
