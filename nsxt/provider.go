/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/go-vmware-nsxt/licensing"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
)

var defaultRetryOnStatusCodes = []int{429, 503}

// Provider configuration that is shared for policy and MP
type commonProviderConfig struct {
	RemoteAuth             bool
	BearerToken            string
	ToleratePartialSuccess bool
}

type nsxtClients struct {
	CommonConfig commonProviderConfig
	// NSX Manager client - based on go-vmware-nsxt SDK
	NsxtClient *api.APIClient
	// Data for NSX Policy client - based on vsphere-automation-sdk-go SDK
	// First offering of Policy SDK does not support concurrent
	// operations in single connector. In order to avoid heavy locks,
	// we are allocating connector per provider operation.
	// TODO: when concurrency support is introduced policy client,
	// change this code to allocate single connector for all provider
	// operations.
	PolicySecurityContext  *core.SecurityContextImpl
	PolicyHTTPClient       *http.Client
	Host                   string
	PolicyEnforcementPoint string
	PolicyGlobalManager    bool
}

// Provider for VMWare NSX-T
func Provider() *schema.Provider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"allow_unverified_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_ALLOW_UNVERIFIED_SSL", false),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_PASSWORD", nil),
				Sensitive:   true,
			},
			"remote_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_REMOTE_AUTH", false),
			},
			"host": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("NSXT_MANAGER_HOST", nil),
				ValidateFunc: validateNsxtProviderHostFormat(),
				Description:  "The hostname or IP address of the NSX manager.",
			},
			"client_auth_cert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_CLIENT_AUTH_CERT_FILE", nil),
			},
			"client_auth_key_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_CLIENT_AUTH_KEY_FILE", nil),
			},
			"ca_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_CA_FILE", nil),
			},
			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum number of HTTP client retries",
				DefaultFunc: schema.EnvDefaultFunc("NSXT_MAX_RETRIES", 8),
			},
			"retry_min_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum delay in milliseconds between retries of a request",
				DefaultFunc: schema.EnvDefaultFunc("NSXT_RETRY_MIN_DELAY", 500),
			},
			"retry_max_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum delay in milliseconds between retries of a request",
				DefaultFunc: schema.EnvDefaultFunc("NSXT_RETRY_MAX_DELAY", 5000),
			},
			"retry_on_status_codes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "HTTP replies status codes to retry on",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				// There is no support for default values/func for list, so it will be handled later
			},
			"tolerate_partial_success": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Treat partial success status as success",
				DefaultFunc: schema.EnvDefaultFunc("NSXT_TOLERATE_PARTIAL_SUCCESS", false),
			},
			"vmc_auth_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL for VMC authorization service (CSP)",
				DefaultFunc: schema.EnvDefaultFunc("NSXT_VMC_AUTH_HOST", "console.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize"),
			},
			"vmc_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Long-living API token for VMC authorization",
				DefaultFunc: schema.EnvDefaultFunc("NSXT_VMC_TOKEN", nil),
			},
			"vmc_auth_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("NSXT_VMC_AUTH_MODE", "Default"),
				ValidateFunc: validation.StringInSlice([]string{"Default", "Bearer", "Basic"}, false),
				Description:  "Mode for VMC authorization",
			},
			"enforcement_point": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enforcement Point for NSXT Policy",
				DefaultFunc: schema.EnvDefaultFunc("NSXT_POLICY_ENFORCEMENT_POINT", "default"),
			},
			"global_manager": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is this a policy global manager endpoint",
				DefaultFunc: schema.EnvDefaultFunc("NSXT_GLOBAL_MANAGER", false),
			},
			"license_keys": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "license keys",
				ConflictsWith: []string{"vmc_token"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(
							"^[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}$"),
						"Must be a valid nsx license key matching: ^[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}-[A-Z0-9]{5}$"),
				},
			},
			"client_auth_cert": {
				Type:        schema.TypeString,
				Description: "Client certificate passed as string",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_CLIENT_AUTH_CERT", nil),
			},
			"client_auth_key": {
				Type:        schema.TypeString,
				Description: "Client certificate key passed as string",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_CLIENT_AUTH_KEY", nil),
			},
			"ca": {
				Type:        schema.TypeString,
				Description: "CA certificate passed as string",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_CA", nil),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"nsxt_provider_info":                    dataSourceNsxtProviderInfo(),
			"nsxt_transport_zone":                   dataSourceNsxtTransportZone(),
			"nsxt_switching_profile":                dataSourceNsxtSwitchingProfile(),
			"nsxt_logical_tier0_router":             dataSourceNsxtLogicalTier0Router(),
			"nsxt_logical_tier1_router":             dataSourceNsxtLogicalTier1Router(),
			"nsxt_mac_pool":                         dataSourceNsxtMacPool(),
			"nsxt_ns_group":                         dataSourceNsxtNsGroup(),
			"nsxt_ns_service":                       dataSourceNsxtNsService(),
			"nsxt_edge_cluster":                     dataSourceNsxtEdgeCluster(),
			"nsxt_certificate":                      dataSourceNsxtCertificate(),
			"nsxt_ip_pool":                          dataSourceNsxtIPPool(),
			"nsxt_firewall_section":                 dataSourceNsxtFirewallSection(),
			"nsxt_management_cluster":               dataSourceNsxtManagementCluster(),
			"nsxt_policy_edge_cluster":              dataSourceNsxtPolicyEdgeCluster(),
			"nsxt_policy_edge_node":                 dataSourceNsxtPolicyEdgeNode(),
			"nsxt_policy_tier0_gateway":             dataSourceNsxtPolicyTier0Gateway(),
			"nsxt_policy_tier1_gateway":             dataSourceNsxtPolicyTier1Gateway(),
			"nsxt_policy_service":                   dataSourceNsxtPolicyService(),
			"nsxt_policy_realization_info":          dataSourceNsxtPolicyRealizationInfo(),
			"nsxt_policy_segment_realization":       dataSourceNsxtPolicySegmentRealization(),
			"nsxt_policy_transport_zone":            dataSourceNsxtPolicyTransportZone(),
			"nsxt_policy_ip_discovery_profile":      dataSourceNsxtPolicyIPDiscoveryProfile(),
			"nsxt_policy_spoofguard_profile":        dataSourceNsxtPolicySpoofGuardProfile(),
			"nsxt_policy_qos_profile":               dataSourceNsxtPolicyQosProfile(),
			"nsxt_policy_ipv6_ndra_profile":         dataSourceNsxtPolicyIpv6NdraProfile(),
			"nsxt_policy_ipv6_dad_profile":          dataSourceNsxtPolicyIpv6DadProfile(),
			"nsxt_policy_gateway_qos_profile":       dataSourceNsxtPolicyGatewayQosProfile(),
			"nsxt_policy_segment_security_profile":  dataSourceNsxtPolicySegmentSecurityProfile(),
			"nsxt_policy_mac_discovery_profile":     dataSourceNsxtPolicyMacDiscoveryProfile(),
			"nsxt_policy_vm":                        dataSourceNsxtPolicyVM(),
			"nsxt_policy_lb_app_profile":            dataSourceNsxtPolicyLBAppProfile(),
			"nsxt_policy_lb_client_ssl_profile":     dataSourceNsxtPolicyLBClientSslProfile(),
			"nsxt_policy_lb_server_ssl_profile":     dataSourceNsxtPolicyLBServerSslProfile(),
			"nsxt_policy_lb_monitor":                dataSourceNsxtPolicyLBMonitor(),
			"nsxt_policy_certificate":               dataSourceNsxtPolicyCertificate(),
			"nsxt_policy_lb_persistence_profile":    dataSourceNsxtPolicyLbPersistenceProfile(),
			"nsxt_policy_vni_pool":                  dataSourceNsxtPolicyVniPool(),
			"nsxt_policy_ip_block":                  dataSourceNsxtPolicyIPBlock(),
			"nsxt_policy_ip_pool":                   dataSourceNsxtPolicyIPPool(),
			"nsxt_policy_site":                      dataSourceNsxtPolicySite(),
			"nsxt_policy_gateway_policy":            dataSourceNsxtPolicyGatewayPolicy(),
			"nsxt_policy_security_policy":           dataSourceNsxtPolicySecurityPolicy(),
			"nsxt_policy_group":                     dataSourceNsxtPolicyGroup(),
			"nsxt_policy_context_profile":           dataSourceNsxtPolicyContextProfile(),
			"nsxt_policy_dhcp_server":               dataSourceNsxtPolicyDhcpServer(),
			"nsxt_policy_bfd_profile":               dataSourceNsxtPolicyBfdProfile(),
			"nsxt_policy_intrusion_service_profile": dataSourceNsxtPolicyIntrusionServiceProfile(),
			"nsxt_policy_lb_service":                dataSourceNsxtPolicyLbService(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"nsxt_dhcp_relay_profile":                      resourceNsxtDhcpRelayProfile(),
			"nsxt_dhcp_relay_service":                      resourceNsxtDhcpRelayService(),
			"nsxt_dhcp_server_profile":                     resourceNsxtDhcpServerProfile(),
			"nsxt_logical_dhcp_server":                     resourceNsxtLogicalDhcpServer(),
			"nsxt_dhcp_server_ip_pool":                     resourceNsxtDhcpServerIPPool(),
			"nsxt_logical_switch":                          resourceNsxtLogicalSwitch(),
			"nsxt_vlan_logical_switch":                     resourceNsxtVlanLogicalSwitch(),
			"nsxt_logical_dhcp_port":                       resourceNsxtLogicalDhcpPort(),
			"nsxt_logical_port":                            resourceNsxtLogicalPort(),
			"nsxt_logical_tier0_router":                    resourceNsxtLogicalTier0Router(),
			"nsxt_logical_tier1_router":                    resourceNsxtLogicalTier1Router(),
			"nsxt_logical_router_centralized_service_port": resourceNsxtLogicalRouterCentralizedServicePort(),
			"nsxt_logical_router_downlink_port":            resourceNsxtLogicalRouterDownLinkPort(),
			"nsxt_logical_router_link_port_on_tier0":       resourceNsxtLogicalRouterLinkPortOnTier0(),
			"nsxt_logical_router_link_port_on_tier1":       resourceNsxtLogicalRouterLinkPortOnTier1(),
			"nsxt_ip_discovery_switching_profile":          resourceNsxtIPDiscoverySwitchingProfile(),
			"nsxt_mac_management_switching_profile":        resourceNsxtMacManagementSwitchingProfile(),
			"nsxt_qos_switching_profile":                   resourceNsxtQosSwitchingProfile(),
			"nsxt_spoofguard_switching_profile":            resourceNsxtSpoofGuardSwitchingProfile(),
			"nsxt_switch_security_switching_profile":       resourceNsxtSwitchSecuritySwitchingProfile(),
			"nsxt_l4_port_set_ns_service":                  resourceNsxtL4PortSetNsService(),
			"nsxt_algorithm_type_ns_service":               resourceNsxtAlgorithmTypeNsService(),
			"nsxt_icmp_type_ns_service":                    resourceNsxtIcmpTypeNsService(),
			"nsxt_igmp_type_ns_service":                    resourceNsxtIgmpTypeNsService(),
			"nsxt_ether_type_ns_service":                   resourceNsxtEtherTypeNsService(),
			"nsxt_ip_protocol_ns_service":                  resourceNsxtIPProtocolNsService(),
			"nsxt_ns_service_group":                        resourceNsxtNsServiceGroup(),
			"nsxt_ns_group":                                resourceNsxtNsGroup(),
			"nsxt_firewall_section":                        resourceNsxtFirewallSection(),
			"nsxt_nat_rule":                                resourceNsxtNatRule(),
			"nsxt_ip_block":                                resourceNsxtIPBlock(),
			"nsxt_ip_block_subnet":                         resourceNsxtIPBlockSubnet(),
			"nsxt_ip_pool":                                 resourceNsxtIPPool(),
			"nsxt_ip_pool_allocation_ip_address":           resourceNsxtIPPoolAllocationIPAddress(),
			"nsxt_ip_set":                                  resourceNsxtIPSet(),
			"nsxt_static_route":                            resourceNsxtStaticRoute(),
			"nsxt_vm_tags":                                 resourceNsxtVMTags(),
			"nsxt_lb_icmp_monitor":                         resourceNsxtLbIcmpMonitor(),
			"nsxt_lb_tcp_monitor":                          resourceNsxtLbTCPMonitor(),
			"nsxt_lb_udp_monitor":                          resourceNsxtLbUDPMonitor(),
			"nsxt_lb_http_monitor":                         resourceNsxtLbHTTPMonitor(),
			"nsxt_lb_https_monitor":                        resourceNsxtLbHTTPSMonitor(),
			"nsxt_lb_passive_monitor":                      resourceNsxtLbPassiveMonitor(),
			"nsxt_lb_pool":                                 resourceNsxtLbPool(),
			"nsxt_lb_tcp_virtual_server":                   resourceNsxtLbTCPVirtualServer(),
			"nsxt_lb_udp_virtual_server":                   resourceNsxtLbUDPVirtualServer(),
			"nsxt_lb_http_virtual_server":                  resourceNsxtLbHTTPVirtualServer(),
			"nsxt_lb_http_forwarding_rule":                 resourceNsxtLbHTTPForwardingRule(),
			"nsxt_lb_http_request_rewrite_rule":            resourceNsxtLbHTTPRequestRewriteRule(),
			"nsxt_lb_http_response_rewrite_rule":           resourceNsxtLbHTTPResponseRewriteRule(),
			"nsxt_lb_cookie_persistence_profile":           resourceNsxtLbCookiePersistenceProfile(),
			"nsxt_lb_source_ip_persistence_profile":        resourceNsxtLbSourceIPPersistenceProfile(),
			"nsxt_lb_client_ssl_profile":                   resourceNsxtLbClientSslProfile(),
			"nsxt_lb_server_ssl_profile":                   resourceNsxtLbServerSslProfile(),
			"nsxt_lb_service":                              resourceNsxtLbService(),
			"nsxt_lb_fast_tcp_application_profile":         resourceNsxtLbFastTCPApplicationProfile(),
			"nsxt_lb_fast_udp_application_profile":         resourceNsxtLbFastUDPApplicationProfile(),
			"nsxt_lb_http_application_profile":             resourceNsxtLbHTTPApplicationProfile(),
			"nsxt_policy_tier1_gateway":                    resourceNsxtPolicyTier1Gateway(),
			"nsxt_policy_tier1_gateway_interface":          resourceNsxtPolicyTier1GatewayInterface(),
			"nsxt_policy_tier0_gateway":                    resourceNsxtPolicyTier0Gateway(),
			"nsxt_policy_tier0_gateway_interface":          resourceNsxtPolicyTier0GatewayInterface(),
			"nsxt_policy_tier0_gateway_ha_vip_config":      resourceNsxtPolicyTier0GatewayHAVipConfig(),
			"nsxt_policy_group":                            resourceNsxtPolicyGroup(),
			"nsxt_policy_domain":                           resourceNsxtPolicyDomain(),
			"nsxt_policy_security_policy":                  resourceNsxtPolicySecurityPolicy(),
			"nsxt_policy_service":                          resourceNsxtPolicyService(),
			"nsxt_policy_gateway_policy":                   resourceNsxtPolicyGatewayPolicy(),
			"nsxt_policy_predefined_gateway_policy":        resourceNsxtPolicyPredefinedGatewayPolicy(),
			"nsxt_policy_predefined_security_policy":       resourceNsxtPolicyPredefinedSecurityPolicy(),
			"nsxt_policy_segment":                          resourceNsxtPolicySegment(),
			"nsxt_policy_vlan_segment":                     resourceNsxtPolicyVlanSegment(),
			"nsxt_policy_fixed_segment":                    resourceNsxtPolicyFixedSegment(),
			"nsxt_policy_static_route":                     resourceNsxtPolicyStaticRoute(),
			"nsxt_policy_gateway_prefix_list":              resourceNsxtPolicyGatewayPrefixList(),
			"nsxt_policy_vm_tags":                          resourceNsxtPolicyVMTags(),
			"nsxt_policy_nat_rule":                         resourceNsxtPolicyNATRule(),
			"nsxt_policy_ip_block":                         resourceNsxtPolicyIPBlock(),
			"nsxt_policy_lb_pool":                          resourceNsxtPolicyLBPool(),
			"nsxt_policy_ip_pool":                          resourceNsxtPolicyIPPool(),
			"nsxt_policy_ip_pool_block_subnet":             resourceNsxtPolicyIPPoolBlockSubnet(),
			"nsxt_policy_ip_pool_static_subnet":            resourceNsxtPolicyIPPoolStaticSubnet(),
			"nsxt_policy_lb_service":                       resourceNsxtPolicyLBService(),
			"nsxt_policy_lb_virtual_server":                resourceNsxtPolicyLBVirtualServer(),
			"nsxt_policy_ip_address_allocation":            resourceNsxtPolicyIPAddressAllocation(),
			"nsxt_policy_bgp_neighbor":                     resourceNsxtPolicyBgpNeighbor(),
			"nsxt_policy_bgp_config":                       resourceNsxtPolicyBgpConfig(),
			"nsxt_policy_dhcp_relay":                       resourceNsxtPolicyDhcpRelayConfig(),
			"nsxt_policy_dhcp_server":                      resourceNsxtPolicyDhcpServer(),
			"nsxt_policy_context_profile":                  resourceNsxtPolicyContextProfile(),
			"nsxt_policy_dhcp_v4_static_binding":           resourceNsxtPolicyDhcpV4StaticBinding(),
			"nsxt_policy_dhcp_v6_static_binding":           resourceNsxtPolicyDhcpV6StaticBinding(),
			"nsxt_policy_dns_forwarder_zone":               resourceNsxtPolicyDNSForwarderZone(),
			"nsxt_policy_gateway_dns_forwarder":            resourceNsxtPolicyGatewayDNSForwarder(),
			"nsxt_policy_gateway_community_list":           resourceNsxtPolicyGatewayCommunityList(),
			"nsxt_policy_gateway_route_map":                resourceNsxtPolicyGatewayRouteMap(),
			"nsxt_policy_intrusion_service_policy":         resourceNsxtPolicyIntrusionServicePolicy(),
			"nsxt_policy_static_route_bfd_peer":            resourceNsxtPolicyStaticRouteBfdPeer(),
			"nsxt_policy_intrusion_service_profile":        resourceNsxtPolicyIntrusionServiceProfile(),
			"nsxt_policy_evpn_tenant":                      resourceNsxtPolicyEvpnTenant(),
			"nsxt_policy_evpn_config":                      resourceNsxtPolicyEvpnConfig(),
			"nsxt_policy_evpn_tunnel_endpoint":             resourceNsxtPolicyEvpnTunnelEndpoint(),
			"nsxt_policy_qos_profile":                      resourceNsxtPolicyQosProfile(),
			"nsxt_policy_ospf_config":                      resourceNsxtPolicyOspfConfig(),
			"nsxt_policy_ospf_area":                        resourceNsxtPolicyOspfArea(),
			"nsxt_policy_gateway_redistribution_config":    resourceNsxtPolicyGatewayRedistributionConfig(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func configureNsxtClient(d *schema.ResourceData, clients *nsxtClients) error {
	clientAuthCertFile := d.Get("client_auth_cert_file").(string)
	clientAuthKeyFile := d.Get("client_auth_key_file").(string)
	clientAuthCert := d.Get("client_auth_cert").(string)
	clientAuthKey := d.Get("client_auth_key").(string)
	vmcToken := d.Get("vmc_token").(string)
	vmcAuthMode := d.Get("vmc_auth_mode").(string)

	if (len(vmcToken) > 0) || (vmcAuthMode == "Basic") {
		// VMC can operate without token with basic auth, however MP API is not
		// available for cloud admin user
		return nil
	}

	needCreds := true
	if len(clientAuthCertFile) > 0 {
		if len(clientAuthKeyFile) == 0 {
			return fmt.Errorf("Please provide key file for client certificate")
		}
		needCreds = false
	}

	if len(clientAuthCert) > 0 {
		if len(clientAuthKey) == 0 {
			return fmt.Errorf("Please provide key for client certificate")
		}
		// only supported for policy resources
		needCreds = false
	}

	insecure := d.Get("allow_unverified_ssl").(bool)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	if needCreds {
		if username == "" {
			return fmt.Errorf("username must be provided")
		}

		if password == "" {
			return fmt.Errorf("password must be provided")
		}
	}

	host := d.Get("host").(string)
	// Remove schema
	host = strings.TrimPrefix(host, "https://")

	if host == "" {
		return fmt.Errorf("host must be provided")
	}

	caFile := d.Get("ca_file").(string)
	caString := d.Get("ca").(string)

	maxRetries := d.Get("max_retries").(int)
	retryMinDelay := d.Get("retry_min_delay").(int)
	retryMaxDelay := d.Get("retry_max_delay").(int)

	statuses := d.Get("retry_on_status_codes").([]interface{})
	if len(statuses) == 0 {
		// Set to the defaults if empty
		for _, val := range defaultRetryOnStatusCodes {
			statuses = append(statuses, val)
		}
	}
	retryStatuses := make([]int, 0, len(statuses))
	for _, s := range statuses {
		retryStatuses = append(retryStatuses, s.(int))
	}

	retriesConfig := api.ClientRetriesConfiguration{
		MaxRetries:      maxRetries,
		RetryMinDelay:   retryMinDelay,
		RetryMaxDelay:   retryMaxDelay,
		RetryOnStatuses: retryStatuses,
	}

	cfg := api.Configuration{
		BasePath:             "/api/v1",
		Host:                 host,
		Scheme:               "https",
		UserAgent:            "terraform-provider-nsxt/1.0",
		UserName:             username,
		Password:             password,
		RemoteAuth:           clients.CommonConfig.RemoteAuth,
		ClientAuthCertFile:   clientAuthCertFile,
		ClientAuthKeyFile:    clientAuthKeyFile,
		CAFile:               caFile,
		ClientAuthCertString: clientAuthCert,
		ClientAuthKeyString:  clientAuthKey,
		CAString:             caString,
		Insecure:             insecure,
		RetriesConfiguration: retriesConfig,
	}

	nsxClient, err := api.NewAPIClient(&cfg)
	if err != nil {
		return err
	}

	clients.NsxtClient = nsxClient

	return initNSXVersion(nsxClient)
}

type jwtToken struct {
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func getAPIToken(vmcAuthHost string, vmcAccessToken string) (string, error) {

	payload := strings.NewReader("refresh_token=" + vmcAccessToken)
	req, _ := http.NewRequest("POST", "https://"+vmcAuthHost, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		b, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("Unexpected status code %d trying to get auth token. %s", res.StatusCode, string(b))
	}

	defer res.Body.Close()
	token := jwtToken{}
	err = json.NewDecoder(res.Body).Decode(&token)
	if err != nil {
		// Not fatal
		log.Printf("[WARNING]: Failed to decode access token from response: %v", err)
	}

	return token.AccessToken, nil
}

func getConnectorTLSConfig(d *schema.ResourceData) (*tls.Config, error) {

	insecure := d.Get("allow_unverified_ssl").(bool)
	clientAuthCertFile := d.Get("client_auth_cert_file").(string)
	clientAuthKeyFile := d.Get("client_auth_key_file").(string)
	caFile := d.Get("ca_file").(string)
	clientAuthCert := d.Get("client_auth_cert").(string)
	clientAuthKey := d.Get("client_auth_key").(string)
	caCert := d.Get("ca").(string)
	tlsConfig := tls.Config{InsecureSkipVerify: insecure}

	if len(clientAuthCertFile) > 0 {

		// cert and key are passed via filesystem
		if len(clientAuthKeyFile) == 0 {
			return nil, fmt.Errorf("Please provide key file for client certificate")
		}

		cert, err := tls.LoadX509KeyPair(clientAuthCertFile, clientAuthKeyFile)

		if err != nil {
			return nil, fmt.Errorf("Failed to load client cert/key pair: %v", err)
		}

		tlsConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return &cert, nil
		}
	}

	if len(clientAuthCert) > 0 {
		// cert and key are passed as strings
		if len(clientAuthKey) == 0 {
			return nil, fmt.Errorf("Please provide key for client certificate")
		}

		cert, err := tls.X509KeyPair([]byte(clientAuthCert), []byte(clientAuthKey))

		if err != nil {
			return nil, fmt.Errorf("Failed to load client cert/key pair: %v", err)
		}

		tlsConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return &cert, nil
		}
	}

	if len(caFile) > 0 {
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig.RootCAs = caCertPool
	}

	if len(caCert) > 0 {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))

		tlsConfig.RootCAs = caCertPool
	}

	return &tlsConfig, nil
}

func configurePolicyConnectorData(d *schema.ResourceData, clients *nsxtClients) error {
	host := d.Get("host").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	vmcAccessToken := d.Get("vmc_token").(string)
	vmcAuthHost := d.Get("vmc_auth_host").(string)
	clientAuthCertFile := d.Get("client_auth_cert_file").(string)
	clientAuthCert := d.Get("client_auth_cert").(string)
	clientAuthDefined := (len(clientAuthCertFile) > 0) || (len(clientAuthCert) > 0)
	policyEnforcementPoint := d.Get("enforcement_point").(string)
	policyGlobalManager := d.Get("global_manager").(bool)
	vmcAuthMode := d.Get("vmc_auth_mode").(string)

	if host == "" {
		return fmt.Errorf("host must be provided")
	}

	if !strings.HasPrefix(host, "https://") {
		host = fmt.Sprintf("https://%s", host)
	}

	securityCtx := core.NewSecurityContextImpl()
	securityContextNeeded := true
	if clientAuthDefined && !clients.CommonConfig.RemoteAuth {
		securityContextNeeded = false
	}

	if securityContextNeeded {
		if len(vmcAccessToken) > 0 {
			if vmcAuthHost == "" {
				return fmt.Errorf("vmc auth host must be provided if auth token is provided")
			}

			apiToken, err := getAPIToken(vmcAuthHost, vmcAccessToken)
			if err != nil {
				return err
			}

			if vmcAuthMode == "Bearer" {
				clients.CommonConfig.BearerToken = apiToken
			} else {

				securityCtx.SetProperty(security.AUTHENTICATION_SCHEME_ID, security.OAUTH_SCHEME_ID)
				securityCtx.SetProperty(security.ACCESS_TOKEN, apiToken)
			}
		} else {
			if username == "" {
				return fmt.Errorf("username must be provided")
			}

			if password == "" {
				return fmt.Errorf("password must be provided")
			}

			securityCtx.SetProperty(security.AUTHENTICATION_SCHEME_ID, security.USER_PASSWORD_SCHEME_ID)
			securityCtx.SetProperty(security.USER_KEY, username)
			securityCtx.SetProperty(security.PASSWORD_KEY, password)
		}
	}

	tlsConfig, err := getConnectorTLSConfig(d)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: tlsConfig,
	}

	httpClient := http.Client{Transport: tr}
	clients.PolicyHTTPClient = &httpClient
	if securityContextNeeded {
		clients.PolicySecurityContext = securityCtx
	}
	clients.Host = host
	clients.PolicyEnforcementPoint = policyEnforcementPoint
	clients.PolicyGlobalManager = policyGlobalManager

	if (len(vmcAccessToken) > 0) || (vmcAuthMode == "Basic") {
		// Special treatment for VMC since MP API is not available there
		initNSXVersionVMC(*clients)
	}
	return nil
}

type remoteAuthHeaderProcessor struct {
}

func newRemoteAuthHeaderProcessor() *remoteAuthHeaderProcessor {
	return &remoteAuthHeaderProcessor{}
}

func (processor remoteAuthHeaderProcessor) Process(req *http.Request) error {
	oldAuthHeader := req.Header.Get("Authorization")
	newAuthHeader := strings.Replace(oldAuthHeader, "Basic", "Remote", 1)
	req.Header.Set("Authorization", newAuthHeader)
	return nil
}

type bearerAuthHeaderProcessor struct {
	Token string
}

func newBearerAuthHeaderProcessor(token string) *bearerAuthHeaderProcessor {
	return &bearerAuthHeaderProcessor{Token: token}
}

func (processor bearerAuthHeaderProcessor) Process(req *http.Request) error {
	newAuthHeader := fmt.Sprintf("Bearer %s", processor.Token)
	req.Header.Set("Authorization", newAuthHeader)
	return nil
}

func applyLicense(c *api.APIClient, licenseKey string) error {
	if c == nil {
		return fmt.Errorf("API client not configured")
	}

	license := licensing.License{LicenseKey: licenseKey}
	_, resp, err := c.LicensingApi.CreateLicense(c.Context, license)
	if err != nil {
		return fmt.Errorf("Error during license create: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status returned during license create: %v", resp.StatusCode)
	}

	return nil
}

// license keys are applied on terraform plan and are not removed
func configureLicenses(d *schema.ResourceData, clients *nsxtClients) error {
	for _, licKey := range d.Get("license_keys").([]interface{}) {
		err := applyLicense(clients.NsxtClient, licKey.(string))
		if err != nil {
			return fmt.Errorf("Error applying license key: %s. %s", licKey, err.Error())
		}
	}
	return nil
}

func initCommonConfig(d *schema.ResourceData) commonProviderConfig {
	remoteAuth := d.Get("remote_auth").(bool)
	toleratePartialSuccess := d.Get("tolerate_partial_success").(bool)

	return commonProviderConfig{
		RemoteAuth:             remoteAuth,
		ToleratePartialSuccess: toleratePartialSuccess,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	commonConfig := initCommonConfig(d)
	clients := nsxtClients{
		CommonConfig: commonConfig,
	}

	err := configureNsxtClient(d, &clients)
	if err != nil {
		return nil, err
	}

	err = configurePolicyConnectorData(d, &clients)
	if err != nil {
		return nil, err
	}

	err = configureLicenses(d, &clients)
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func getPolicyConnector(clients interface{}) *client.RestConnector {
	c := clients.(nsxtClients)
	connector := client.NewRestConnector(c.Host, *c.PolicyHTTPClient)
	if c.PolicySecurityContext != nil {
		connector.SetSecurityContext(c.PolicySecurityContext)
	}
	if c.CommonConfig.RemoteAuth {
		connector.AddRequestProcessor(newRemoteAuthHeaderProcessor())
	}
	if len(c.CommonConfig.BearerToken) > 0 {
		connector.AddRequestProcessor(newBearerAuthHeaderProcessor(c.CommonConfig.BearerToken))
	}

	return connector
}

func getPolicyEnforcementPoint(clients interface{}) string {
	return clients.(nsxtClients).PolicyEnforcementPoint
}

func isPolicyGlobalManager(clients interface{}) bool {
	return clients.(nsxtClients).PolicyGlobalManager
}

func getCommonProviderConfig(clients interface{}) commonProviderConfig {
	return clients.(nsxtClients).CommonConfig
}

func getGlobalPolicyEnforcementPointPath(m interface{}, sitePath *string) string {
	return fmt.Sprintf("%s/enforcement-points/%s", *sitePath, getPolicyEnforcementPoint(m))
}
