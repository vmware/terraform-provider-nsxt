// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	tf_api "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	api "github.com/vmware/go-vmware-nsxt"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client/middleware/retry"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"golang.org/x/exp/slices"
)

var defaultRetryOnStatusCodes = []int{400, 409, 429, 500, 503, 504}

// Provider configuration that is shared for policy and MP
type commonProviderConfig struct {
	RemoteAuth             bool
	BearerToken            string
	ToleratePartialSuccess bool
	MaxRetries             int
	MinRetryInterval       int
	MaxRetryInterval       int
	RetryStatusCodes       []int
	Username               string
	Password               string
	LicenseKeys            []string
}

type nsxtClients struct {
	CommonConfig commonProviderConfig
	// NSX Manager client - based on go-vmware-nsxt SDK
	NsxtClient *api.APIClient
	// Config for the above client
	NsxtClientConfig *api.Configuration
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
			"session_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_SESSION_AUTH", true),
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
				DefaultFunc: schema.EnvDefaultFunc("NSXT_MAX_RETRIES", 4),
			},
			"retry_min_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum delay in milliseconds between retries of a request",
				DefaultFunc: schema.EnvDefaultFunc("NSXT_RETRY_MIN_DELAY", 0),
			},
			"retry_max_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum delay in milliseconds between retries of a request",
				DefaultFunc: schema.EnvDefaultFunc("NSXT_RETRY_MAX_DELAY", 500),
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
				DefaultFunc: schema.EnvDefaultFunc("NSXT_VMC_AUTH_HOST", nil),
			},
			"vmc_token": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Long-living API token for VMC authorization",
				DefaultFunc:   schema.EnvDefaultFunc("NSXT_VMC_TOKEN", nil),
				ConflictsWith: []string{"vmc_client_id", "vmc_client_secret"},
			},
			"vmc_client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "ID of OAuth App associated with the VMC organization",
				DefaultFunc:   schema.EnvDefaultFunc("NSXT_VMC_CLIENT_ID", nil),
				ConflictsWith: []string{"vmc_token"},
				RequiredWith:  []string{"vmc_client_secret"},
			},
			"vmc_client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Secret of OAuth App associated with the VMC organization",
				DefaultFunc:   schema.EnvDefaultFunc("NSXT_VMC_CLIENT_SECRET", nil),
				ConflictsWith: []string{"vmc_token"},
				RequiredWith:  []string{"vmc_client_id"},
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
				ConflictsWith: []string{"vmc_token", "vmc_client_id", "vmc_client_secret"},
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
			"on_demand_connection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Avoid initializing NSX connection on startup",
				DefaultFunc: schema.EnvDefaultFunc("NSXT_ON_DEMAND_CONNECTION", false),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			// Ensure alphabetical order for organization and readability.
			"nsxt_certificate":                                       dataSourceNsxtCertificate(),
			"nsxt_compute_collection":                                dataSourceNsxtComputeCollection(),
			"nsxt_compute_manager":                                   dataSourceNsxtComputeManager(),
			"nsxt_compute_manager_realization":                       dataSourceNsxtComputeManagerRealization(),
			"nsxt_discovered_node":                                   dataSourceNsxtDiscoveredNode(),
			"nsxt_edge_cluster":                                      dataSourceNsxtEdgeCluster(),
			"nsxt_edge_upgrade_group":                                dataSourceNsxtEdgeUpgradeGroup(),
			"nsxt_failure_domain":                                    dataSourceNsxtFailureDomain(),
			"nsxt_firewall_section":                                  removedDataSourceWrapper(dataSourceNsxtFirewallSection, "nsxt_firewall_section"),
			"nsxt_host_upgrade_group":                                dataSourceNsxtHostUpgradeGroup(),
			"nsxt_ip_pool":                                           removedDataSourceWrapper(dataSourceNsxtIPPool, "nsxt_ip_pool"),
			"nsxt_logical_tier0_router":                              removedDataSourceWrapper(dataSourceNsxtLogicalTier0Router, "nsxt_logical_tier0_router"),
			"nsxt_logical_tier1_router":                              removedDataSourceWrapper(dataSourceNsxtLogicalTier1Router, "nsxt_logical_tier1_router"),
			"nsxt_mac_pool":                                          removedDataSourceWrapper(dataSourceNsxtMacPool, "nsxt_mac_pool"),
			"nsxt_management_cluster":                                dataSourceNsxtManagementCluster(),
			"nsxt_manager_cluster_node":                              dataSourceNsxtManagerClusterNode(),
			"nsxt_manager_info":                                      dataSourceNsxtManagerInfo(),
			"nsxt_ns_group":                                          removedDataSourceWrapper(dataSourceNsxtNsGroup, "nsxt_ns_group"),
			"nsxt_ns_groups":                                         removedDataSourceWrapper(dataSourceNsxtNsGroups, "nsxt_ns_groups"),
			"nsxt_ns_service":                                        removedDataSourceWrapper(dataSourceNsxtNsService, "nsxt_ns_service"),
			"nsxt_ns_services":                                       removedDataSourceWrapper(dataSourceNsxtNsServices, "nsxt_ns_services"),
			"nsxt_policy_bfd_profile":                                dataSourceNsxtPolicyBfdProfile(),
			"nsxt_policy_bridge_profile":                             dataSourceNsxtPolicyBridgeProfile(),
			"nsxt_policy_certificate":                                dataSourceNsxtPolicyCertificate(),
			"nsxt_policy_context_profile":                            dataSourceNsxtPolicyContextProfile(),
			"nsxt_policy_dhcp_server":                                dataSourceNsxtPolicyDhcpServer(),
			"nsxt_policy_distributed_flood_protection_profile":       dataSourceNsxtPolicyDistributedFloodProtectionProfile(),
			"nsxt_policy_distributed_vlan_connection":                dataSourceNsxtPolicyDistributedVlanConnection(),
			"nsxt_policy_edge_cluster":                               dataSourceNsxtPolicyEdgeCluster(),
			"nsxt_policy_edge_high_availability_profile":             dataSourceNsxtPolicyEdgeHighAvailabilityProfile(),
			"nsxt_policy_edge_node":                                  dataSourceNsxtPolicyEdgeNode(),
			"nsxt_policy_edge_transport_node":                        dataSourceNsxtPolicyEdgeTransportNode(),
			"nsxt_policy_edge_transport_node_realization":            dataSourceNsxtPolicyEdgeTransportNodeRealization(),
			"nsxt_policy_gateway_connection":                         dataSourceNsxtPolicyGatewayConnection(),
			"nsxt_policy_gateway_dns_forwarder":                      dataSourceNsxtPolicyGatewayDNSForwarder(),
			"nsxt_policy_gateway_flood_protection_profile":           dataSourceNsxtPolicyGatewayFloodProtectionProfile(),
			"nsxt_policy_gateway_interface":                          dataSourceNsxtPolicyGatewayInterface(),
			"nsxt_policy_gateway_interface_realization":              dataSourceNsxtPolicyGatewayInterfaceRealization(),
			"nsxt_policy_gateway_locale_service":                     dataSourceNsxtPolicyGatewayLocaleService(),
			"nsxt_policy_gateway_policy":                             dataSourceNsxtPolicyGatewayPolicy(),
			"nsxt_policy_gateway_prefix_list":                        dataSourceNsxtPolicyGatewayPrefixList(),
			"nsxt_policy_gateway_qos_profile":                        dataSourceNsxtPolicyGatewayQosProfile(),
			"nsxt_policy_gateway_route_map":                          dataSourceNsxtPolicyGatewayRouteMap(),
			"nsxt_policy_group":                                      dataSourceNsxtPolicyGroup(),
			"nsxt_policy_groups":                                     dataSourceNsxtPolicyGroups(),
			"nsxt_policy_host_transport_node":                        dataSourceNsxtPolicyHostTransportNode(),
			"nsxt_policy_host_transport_node_collection":             dataSourceNsxtPolicyHostTransportNodeCollection(),
			"nsxt_policy_host_transport_node_collection_realization": dataSourceNsxtPolicyHostTransportNodeCollectionRealization(),
			"nsxt_policy_host_transport_node_profile":                dataSourceNsxtPolicyHostTransportNodeProfile(),
			"nsxt_policy_intrusion_service_profile":                  dataSourceNsxtPolicyIntrusionServiceProfile(),
			"nsxt_policy_ip_block":                                   dataSourceNsxtPolicyIPBlock(),
			"nsxt_policy_ip_discovery_profile":                       dataSourceNsxtPolicyIPDiscoveryProfile(),
			"nsxt_policy_ip_pool":                                    dataSourceNsxtPolicyIPPool(),
			"nsxt_policy_ipsec_vpn_dpd_profile":                      dataSourceNsxtPolicyIPSecVpnDpdProfile(),
			"nsxt_policy_ipsec_vpn_ike_profile":                      dataSourceNsxtPolicyIPSecVpnIkeProfile(),
			"nsxt_policy_ipsec_vpn_local_endpoint":                   dataSourceNsxtPolicyIPSecVpnLocalEndpoint(),
			"nsxt_policy_ipsec_vpn_service":                          dataSourceNsxtPolicyIPSecVpnService(),
			"nsxt_policy_ipsec_vpn_tunnel_profile":                   dataSourceNsxtPolicyIPSecVpnTunnelProfile(),
			"nsxt_policy_ipv6_dad_profile":                           dataSourceNsxtPolicyIpv6DadProfile(),
			"nsxt_policy_ipv6_ndra_profile":                          dataSourceNsxtPolicyIpv6NdraProfile(),
			"nsxt_policy_l2_vpn_service":                             dataSourceNsxtPolicyL2VpnService(),
			"nsxt_policy_lb_app_profile":                             dataSourceNsxtPolicyLBAppProfile(),
			"nsxt_policy_lb_client_ssl_profile":                      dataSourceNsxtPolicyLBClientSslProfile(),
			"nsxt_policy_lb_monitor":                                 dataSourceNsxtPolicyLBMonitor(),
			"nsxt_policy_lb_persistence_profile":                     dataSourceNsxtPolicyLbPersistenceProfile(),
			"nsxt_policy_lb_server_ssl_profile":                      dataSourceNsxtPolicyLBServerSslProfile(),
			"nsxt_policy_lb_service":                                 dataSourceNsxtPolicyLbService(),
			"nsxt_policy_mac_discovery_profile":                      dataSourceNsxtPolicyMacDiscoveryProfile(),
			"nsxt_policy_network_span":                               dataSourceNsxtPolicyNetworkSpan(),
			"nsxt_policy_project":                                    dataSourceNsxtPolicyProject(),
			"nsxt_policy_project_ip_address_allocation":              dataSourceNsxtProjectIpAddressAllocation(),
			"nsxt_policy_qos_profile":                                dataSourceNsxtPolicyQosProfile(),
			"nsxt_policy_realization_info":                           dataSourceNsxtPolicyRealizationInfo(),
			"nsxt_policy_security_policy":                            dataSourceNsxtPolicySecurityPolicy(),
			"nsxt_policy_segment":                                    dataSourceNsxtPolicySegment(),
			"nsxt_policy_segment_realization":                        dataSourceNsxtPolicySegmentRealization(),
			"nsxt_policy_segment_security_profile":                   dataSourceNsxtPolicySegmentSecurityProfile(),
			"nsxt_policy_service":                                    dataSourceNsxtPolicyService(),
			"nsxt_policy_services":                                   dataSourceNsxtPolicyServices(),
			"nsxt_policy_share":                                      dataSourceNsxtPolicyShare(),
			"nsxt_policy_site":                                       dataSourceNsxtPolicySite(),
			"nsxt_policy_spoofguard_profile":                         dataSourceNsxtPolicySpoofGuardProfile(),
			"nsxt_policy_tags":                                       dataSourceNsxtTags(),
			"nsxt_policy_tier0_gateway":                              dataSourceNsxtPolicyTier0Gateway(),
			"nsxt_policy_tier0_gateways":                             dataSourceNsxtPolicyTier0Gateways(),
			"nsxt_policy_tier1_gateway":                              dataSourceNsxtPolicyTier1Gateway(),
			"nsxt_policy_tier1_gateways":                             dataSourceNsxtPolicyTier1Gateways(),
			"nsxt_policy_transit_gateway":                            dataSourceNsxtPolicyTransitGateway(),
			"nsxt_policy_transit_gateway_nat":                        dataSourceNsxtPolicyTransitGatewayNat(),
			"nsxt_policy_transport_zone":                             dataSourceNsxtPolicyTransportZone(),
			"nsxt_policy_uplink_host_switch_profile":                 dataSourceNsxtUplinkHostSwitchProfile(),
			"nsxt_policy_vm":                                         dataSourceNsxtPolicyVM(),
			"nsxt_policy_vms":                                        dataSourceNsxtPolicyVMs(),
			"nsxt_policy_vni_pool":                                   dataSourceNsxtPolicyVniPool(),
			"nsxt_policy_vtep_ha_host_switch_profile":                dataSourceNsxtVtepHAHostSwitchProfile(),
			"nsxt_provider_info":                                     dataSourceNsxtProviderInfo(),
			"nsxt_switching_profile":                                 removedDataSourceWrapper(dataSourceNsxtSwitchingProfile, "nsxt_switching_profile"),
			"nsxt_transport_node":                                    dataSourceNsxtTransportNode(),
			"nsxt_transport_node_realization":                        dataSourceNsxtTransportNodeRealization(),
			"nsxt_transport_zone":                                    dataSourceNsxtTransportZone(),
			"nsxt_upgrade_postcheck":                                 dataSourceNsxtUpgradePostCheck(),
			"nsxt_upgrade_prepare_ready":                             dataSourceNsxtUpgradePrepareReady(),
			"nsxt_vpc":                                               dataSourceNsxtVPC(),
			"nsxt_vpc_connectivity_profile":                          dataSourceNsxtVpcConnectivityProfile(),
			"nsxt_vpc_group":                                         dataSourceNsxtVpcGroup(),
			"nsxt_vpc_ip_address_allocation":                         dataSourceNsxtVpcIpAddressAllocation(),
			"nsxt_vpc_nat":                                           dataSourceNsxtVpcNat(),
			"nsxt_vpc_service_profile":                               dataSourceNsxtVpcServiceProfile(),
			"nsxt_vpc_subnet":                                        dataSourceNsxtVpcSubnet(),
			"nsxt_vpc_subnet_port":                                   dataSourceNsxtVpcSubnetPort(),
		},

		ResourcesMap: map[string]*schema.Resource{
			// Ensure alphabetical order for organization and readability.
			"nsxt_algorithm_type_ns_service":                           removedResourceWrapper(resourceNsxtAlgorithmTypeNsService, "nsxt_algorithm_type_ns_service"),
			"nsxt_cluster_virtual_ip":                                  resourceNsxtClusterVirtualIP(),
			"nsxt_compute_manager":                                     resourceNsxtComputeManager(),
			"nsxt_dhcp_relay_profile":                                  removedResourceWrapper(resourceNsxtDhcpRelayProfile, "nsxt_dhcp_relay_profile"),
			"nsxt_dhcp_relay_service":                                  removedResourceWrapper(resourceNsxtDhcpRelayService, "nsxt_dhcp_relay_service"),
			"nsxt_dhcp_server_ip_pool":                                 removedResourceWrapper(resourceNsxtDhcpServerIPPool, "nsxt_dhcp_server_ip_pool"),
			"nsxt_dhcp_server_profile":                                 removedResourceWrapper(resourceNsxtDhcpServerProfile, "nsxt_dhcp_server_profile"),
			"nsxt_edge_cluster":                                        resourceNsxtEdgeCluster(),
			"nsxt_edge_high_availability_profile":                      resourceNsxtEdgeHighAvailabilityProfile(),
			"nsxt_edge_transport_node":                                 resourceNsxtEdgeTransportNode(),
			"nsxt_edge_transport_node_rtep":                            resourceNsxtEdgeTransportNodeRTEP(),
			"nsxt_ether_type_ns_service":                               removedResourceWrapper(resourceNsxtEtherTypeNsService, "nsxt_ether_type_ns_service"),
			"nsxt_failure_domain":                                      resourceNsxtFailureDomain(),
			"nsxt_firewall_section":                                    removedResourceWrapper(resourceNsxtFirewallSection, "nsxt_firewall_section"),
			"nsxt_icmp_type_ns_service":                                removedResourceWrapper(resourceNsxtIcmpTypeNsService, "nsxt_icmp_type_ns_service"),
			"nsxt_igmp_type_ns_service":                                removedResourceWrapper(resourceNsxtIgmpTypeNsService, "nsxt_igmp_type_ns_service"),
			"nsxt_ip_block":                                            removedResourceWrapper(resourceNsxtIPBlock, "nsxt_ip_block"),
			"nsxt_ip_block_subnet":                                     removedResourceWrapper(resourceNsxtIPBlockSubnet, "nsxt_ip_block_subnet"),
			"nsxt_ip_discovery_switching_profile":                      removedResourceWrapper(resourceNsxtIPDiscoverySwitchingProfile, "nsxt_ip_discovery_switching_profile"),
			"nsxt_ip_pool":                                             removedResourceWrapper(resourceNsxtIPPool, "nsxt_ip_pool"),
			"nsxt_ip_pool_allocation_ip_address":                       removedResourceWrapper(resourceNsxtIPPoolAllocationIPAddress, "nsxt_ip_pool_allocation_ip_address"),
			"nsxt_ip_protocol_ns_service":                              removedResourceWrapper(resourceNsxtIPProtocolNsService, "nsxt_ip_protocol_ns_service"),
			"nsxt_ip_set":                                              removedResourceWrapper(resourceNsxtIPSet, "nsxt_ip_set"),
			"nsxt_l4_port_set_ns_service":                              removedResourceWrapper(resourceNsxtL4PortSetNsService, "nsxt_l4_port_set_ns_service"),
			"nsxt_lb_client_ssl_profile":                               removedResourceWrapper(resourceNsxtLbClientSslProfile, "nsxt_lb_client_ssl_profile"),
			"nsxt_lb_cookie_persistence_profile":                       removedResourceWrapper(resourceNsxtLbCookiePersistenceProfile, "nsxt_lb_cookie_persistence_profile"),
			"nsxt_lb_fast_tcp_application_profile":                     removedResourceWrapper(resourceNsxtLbFastTCPApplicationProfile, "nsxt_lb_fast_tcp_application_profile"),
			"nsxt_lb_fast_udp_application_profile":                     removedResourceWrapper(resourceNsxtLbFastUDPApplicationProfile, "nsxt_lb_fast_udp_application_profile"),
			"nsxt_lb_http_application_profile":                         removedResourceWrapper(resourceNsxtLbHTTPApplicationProfile, "nsxt_lb_http_application_profile"),
			"nsxt_lb_http_forwarding_rule":                             removedResourceWrapper(resourceNsxtLbHTTPForwardingRule, "nsxt_lb_http_forwarding_rule"),
			"nsxt_lb_http_monitor":                                     removedResourceWrapper(resourceNsxtLbHTTPMonitor, "nsxt_lb_http_monitor"),
			"nsxt_lb_http_request_rewrite_rule":                        removedResourceWrapper(resourceNsxtLbHTTPRequestRewriteRule, "nsxt_lb_http_request_rewrite_rule"),
			"nsxt_lb_http_response_rewrite_rule":                       removedResourceWrapper(resourceNsxtLbHTTPResponseRewriteRule, "nsxt_lb_http_response_rewrite_rule"),
			"nsxt_lb_http_virtual_server":                              removedResourceWrapper(resourceNsxtLbHTTPVirtualServer, "nsxt_lb_http_virtual_server"),
			"nsxt_lb_https_monitor":                                    removedResourceWrapper(resourceNsxtLbHTTPSMonitor, "nsxt_lb_https_monitor"),
			"nsxt_lb_icmp_monitor":                                     removedResourceWrapper(resourceNsxtLbIcmpMonitor, "nsxt_lb_icmp_monitor"),
			"nsxt_lb_passive_monitor":                                  removedResourceWrapper(resourceNsxtLbPassiveMonitor, "nsxt_lb_passive_monitor"),
			"nsxt_lb_pool":                                             removedResourceWrapper(resourceNsxtLbPool, "nsxt_lb_pool"),
			"nsxt_lb_server_ssl_profile":                               removedResourceWrapper(resourceNsxtLbServerSslProfile, "nsxt_lb_server_ssl_profile"),
			"nsxt_lb_service":                                          removedResourceWrapper(resourceNsxtLbService, "nsxt_lb_service"),
			"nsxt_lb_source_ip_persistence_profile":                    removedResourceWrapper(resourceNsxtLbSourceIPPersistenceProfile, "nsxt_lb_source_ip_persistence_profile"),
			"nsxt_lb_tcp_monitor":                                      removedResourceWrapper(resourceNsxtLbTCPMonitor, "nsxt_lb_tcp_monitor"),
			"nsxt_lb_tcp_virtual_server":                               removedResourceWrapper(resourceNsxtLbTCPVirtualServer, "nsxt_lb_tcp_virtual_server"),
			"nsxt_lb_udp_monitor":                                      removedResourceWrapper(resourceNsxtLbUDPMonitor, "nsxt_lb_udp_monitor"),
			"nsxt_lb_udp_virtual_server":                               removedResourceWrapper(resourceNsxtLbUDPVirtualServer, "nsxt_lb_udp_virtual_server"),
			"nsxt_logical_dhcp_port":                                   removedResourceWrapper(resourceNsxtLogicalDhcpPort, "nsxt_logical_dhcp_port"),
			"nsxt_logical_dhcp_server":                                 removedResourceWrapper(resourceNsxtLogicalDhcpServer, "nsxt_logical_dhcp_server"),
			"nsxt_logical_port":                                        removedResourceWrapper(resourceNsxtLogicalPort, "nsxt_logical_port"),
			"nsxt_logical_router_centralized_service_port":             removedResourceWrapper(resourceNsxtLogicalRouterCentralizedServicePort, "nsxt_logical_router_centralized_service_port"),
			"nsxt_logical_router_downlink_port":                        removedResourceWrapper(resourceNsxtLogicalRouterDownLinkPort, "nsxt_logical_router_downlink_port"),
			"nsxt_logical_router_link_port_on_tier0":                   removedResourceWrapper(resourceNsxtLogicalRouterLinkPortOnTier0, "nsxt_logical_router_link_port_on_tier0"),
			"nsxt_logical_router_link_port_on_tier1":                   removedResourceWrapper(resourceNsxtLogicalRouterLinkPortOnTier1, "nsxt_logical_router_link_port_on_tier1"),
			"nsxt_logical_switch":                                      removedResourceWrapper(resourceNsxtLogicalSwitch, "nsxt_logical_switch"),
			"nsxt_logical_tier0_router":                                removedResourceWrapper(resourceNsxtLogicalTier0Router, "nsxt_logical_tier0_router"),
			"nsxt_logical_tier1_router":                                removedResourceWrapper(resourceNsxtLogicalTier1Router, "nsxt_logical_tier1_router"),
			"nsxt_mac_management_switching_profile":                    removedResourceWrapper(resourceNsxtMacManagementSwitchingProfile, "nsxt_mac_management_switching_profile"),
			"nsxt_manager_cluster":                                     resourceNsxtManagerCluster(),
			"nsxt_nat_rule":                                            removedResourceWrapper(resourceNsxtNatRule, "nsxt_nat_rule"),
			"nsxt_node_user":                                           resourceNsxtUsers(),
			"nsxt_ns_group":                                            removedResourceWrapper(resourceNsxtNsGroup, "nsxt_ns_group"),
			"nsxt_ns_service_group":                                    removedResourceWrapper(resourceNsxtNsServiceGroup, "nsxt_ns_service_group"),
			"nsxt_policy_bgp_config":                                   resourceNsxtPolicyBgpConfig(),
			"nsxt_policy_bgp_neighbor":                                 resourceNsxtPolicyBgpNeighbor(),
			"nsxt_policy_compute_sub_cluster":                          resourceNsxtPolicyComputeSubCluster(),
			"nsxt_policy_connectivity_policy":                          resourceNsxtPolicyConnectivityPolicy(),
			"nsxt_policy_constraint":                                   resourceNsxtPolicyConstraint(),
			"nsxt_policy_context_profile":                              resourceNsxtPolicyContextProfile(),
			"nsxt_policy_context_profile_custom_attribute":             resourceNsxtPolicyContextProfileCustomAttribute(),
			"nsxt_policy_dhcp_relay":                                   resourceNsxtPolicyDhcpRelayConfig(),
			"nsxt_policy_dhcp_server":                                  resourceNsxtPolicyDhcpServer(),
			"nsxt_policy_dhcp_v4_static_binding":                       resourceNsxtPolicyDhcpV4StaticBinding(),
			"nsxt_policy_dhcp_v6_static_binding":                       resourceNsxtPolicyDhcpV6StaticBinding(),
			"nsxt_policy_distributed_flood_protection_profile":         resourceNsxtPolicyDistributedFloodProtectionProfile(),
			"nsxt_policy_distributed_flood_protection_profile_binding": resourceNsxtPolicyDistributedFloodProtectionProfileBinding(),
			"nsxt_policy_distributed_vlan_connection":                  resourceNsxtPolicyDistributedVlanConnection(),
			"nsxt_policy_dns_forwarder_zone":                           resourceNsxtPolicyDNSForwarderZone(),
			"nsxt_policy_domain":                                       resourceNsxtPolicyDomain(),
			"nsxt_policy_edge_cluster":                                 resourceNsxtPolicyEdgeCluster(),
			"nsxt_policy_edge_high_availability_profile":               resourceNsxtPolicyEdgeHighAvailabilityProfile(),
			"nsxt_policy_edge_transport_node":                          resourceNsxtPolicyEdgeTransportNode(),
			"nsxt_policy_edge_transport_node_rtep":                     resourceNsxtPolicyEdgeTransportNodeRTEP(),
			"nsxt_policy_evpn_config":                                  resourceNsxtPolicyEvpnConfig(),
			"nsxt_policy_evpn_tenant":                                  resourceNsxtPolicyEvpnTenant(),
			"nsxt_policy_evpn_tunnel_endpoint":                         resourceNsxtPolicyEvpnTunnelEndpoint(),
			"nsxt_policy_firewall_exclude_list_member":                 resourceNsxtPolicyFirewallExcludeListMember(),
			"nsxt_policy_fixed_segment":                                resourceNsxtPolicyFixedSegment(),
			"nsxt_policy_gateway_community_list":                       resourceNsxtPolicyGatewayCommunityList(),
			"nsxt_policy_gateway_connection":                           resourceNsxtPolicyGatewayConnection(),
			"nsxt_policy_gateway_dns_forwarder":                        resourceNsxtPolicyGatewayDNSForwarder(),
			"nsxt_policy_gateway_flood_protection_profile":             resourceNsxtPolicyGatewayFloodProtectionProfile(),
			"nsxt_policy_gateway_flood_protection_profile_binding":     resourceNsxtPolicyGatewayFloodProtectionProfileBinding(),
			"nsxt_policy_gateway_policy":                               resourceNsxtPolicyGatewayPolicy(),
			"nsxt_policy_gateway_prefix_list":                          resourceNsxtPolicyGatewayPrefixList(),
			"nsxt_policy_gateway_qos_profile":                          resourceNsxtPolicyGatewayQosProfile(),
			"nsxt_policy_gateway_redistribution_config":                resourceNsxtPolicyGatewayRedistributionConfig(),
			"nsxt_policy_gateway_route_map":                            resourceNsxtPolicyGatewayRouteMap(),
			"nsxt_policy_global_manager":                               resourceNsxtPolicyGlobalManager(),
			"nsxt_policy_group":                                        resourceNsxtPolicyGroup(),
			"nsxt_policy_host_transport_node":                          resourceNsxtPolicyHostTransportNode(),
			"nsxt_policy_host_transport_node_collection":               resourceNsxtPolicyHostTransportNodeCollection(),
			"nsxt_policy_host_transport_node_profile":                  resourceNsxtPolicyHostTransportNodeProfile(),
			"nsxt_policy_intrusion_service_policy":                     resourceNsxtPolicyIntrusionServicePolicy(),
			"nsxt_policy_intrusion_service_profile":                    resourceNsxtPolicyIntrusionServiceProfile(),
			"nsxt_policy_ip_address_allocation":                        resourceNsxtPolicyIPAddressAllocation(),
			"nsxt_policy_ip_block":                                     resourceNsxtPolicyIPBlock(),
			"nsxt_policy_ip_block_quota":                               resourceNsxtPolicyIpBlockQuota(),
			"nsxt_policy_ip_discovery_profile":                         resourceNsxtPolicyIPDiscoveryProfile(),
			"nsxt_policy_ip_pool":                                      resourceNsxtPolicyIPPool(),
			"nsxt_policy_ip_pool_block_subnet":                         resourceNsxtPolicyIPPoolBlockSubnet(),
			"nsxt_policy_ip_pool_static_subnet":                        resourceNsxtPolicyIPPoolStaticSubnet(),
			"nsxt_policy_ipsec_vpn_dpd_profile":                        resourceNsxtPolicyIPSecVpnDpdProfile(),
			"nsxt_policy_ipsec_vpn_ike_profile":                        resourceNsxtPolicyIPSecVpnIkeProfile(),
			"nsxt_policy_ipsec_vpn_local_endpoint":                     resourceNsxtPolicyIPSecVpnLocalEndpoint(),
			"nsxt_policy_ipsec_vpn_service":                            resourceNsxtPolicyIPSecVpnService(),
			"nsxt_policy_ipsec_vpn_session":                            resourceNsxtPolicyIPSecVpnSession(),
			"nsxt_policy_ipsec_vpn_tunnel_profile":                     resourceNsxtPolicyIPSecVpnTunnelProfile(),
			"nsxt_policy_l2_vpn_service":                               resourceNsxtPolicyL2VpnService(),
			"nsxt_policy_l2_vpn_session":                               resourceNsxtPolicyL2VPNSession(),
			"nsxt_policy_l7_access_profile":                            resourceNsxtPolicyL7AccessProfile(),
			"nsxt_policy_lb_client_ssl_profile":                        resourceNsxtPolicyLBClientSslProfile(),
			"nsxt_policy_lb_cookie_persistence_profile":                resourceNsxtPolicyLBCookiePersistenceProfile(),
			"nsxt_policy_lb_fast_tcp_application_profile":              resourceNsxtPolicyLBFastTcpApplicationProfile(),
			"nsxt_policy_lb_fast_udp_application_profile":              resourceNsxtPolicyLBFastUdpApplicationProfile(),
			"nsxt_policy_lb_generic_persistence_profile":               resourceNsxtPolicyLBGenericPersistenceProfile(),
			"nsxt_policy_lb_http_application_profile":                  resourceNsxtPolicyLBHttpApplicationProfile(),
			"nsxt_policy_lb_http_monitor_profile":                      resourceNsxtPolicyLBHttpMonitorProfile(),
			"nsxt_policy_lb_https_monitor_profile":                     resourceNsxtPolicyLBHttpsMonitorProfile(),
			"nsxt_policy_lb_icmp_monitor_profile":                      resourceNsxtPolicyLBIcmpMonitorProfile(),
			"nsxt_policy_lb_passive_monitor_profile":                   resourceNsxtPolicyLBPassiveMonitorProfile(),
			"nsxt_policy_lb_pool":                                      resourceNsxtPolicyLBPool(),
			"nsxt_policy_lb_server_ssl_profile":                        resourceNsxtPolicyLBServerSslProfile(),
			"nsxt_policy_lb_service":                                   resourceNsxtPolicyLBService(),
			"nsxt_policy_lb_source_ip_persistence_profile":             resourceNsxtPolicyLBSourceIpPersistenceProfile(),
			"nsxt_policy_lb_tcp_monitor_profile":                       resourceNsxtPolicyLBTcpMonitorProfile(),
			"nsxt_policy_lb_udp_monitor_profile":                       resourceNsxtPolicyLBUdpMonitorProfile(),
			"nsxt_policy_lb_virtual_server":                            resourceNsxtPolicyLBVirtualServer(),
			"nsxt_policy_ldap_identity_source":                         resourceNsxtPolicyLdapIdentitySource(),
			"nsxt_policy_mac_discovery_profile":                        resourceNsxtPolicyMacDiscoveryProfile(),
			"nsxt_policy_metadata_proxy":                               resourceNsxtPolicyMetadataProxy(),
			"nsxt_policy_nat_rule":                                     resourceNsxtPolicyNATRule(),
			"nsxt_policy_network_span":                                 resourceNsxtPolicyNetworkSpan(),
			"nsxt_policy_ospf_area":                                    resourceNsxtPolicyOspfArea(),
			"nsxt_policy_ospf_config":                                  resourceNsxtPolicyOspfConfig(),
			"nsxt_policy_parent_security_policy":                       resourceNsxtPolicyParentSecurityPolicy(),
			"nsxt_policy_predefined_gateway_policy":                    resourceNsxtPolicyPredefinedGatewayPolicy(),
			"nsxt_policy_predefined_security_policy":                   resourceNsxtPolicyPredefinedSecurityPolicy(),
			"nsxt_policy_project":                                      resourceNsxtPolicyProject(),
			"nsxt_policy_project_ip_address_allocation":                resourceNsxtPolicyProjectIpAddressAllocation(),
			"nsxt_policy_qos_profile":                                  resourceNsxtPolicyQosProfile(),
			"nsxt_policy_security_policy":                              resourceNsxtPolicySecurityPolicy(),
			"nsxt_policy_security_policy_rule":                         resourceNsxtPolicySecurityPolicyRule(),
			"nsxt_policy_segment":                                      resourceNsxtPolicySegment(),
			"nsxt_policy_segment_security_profile":                     resourceNsxtPolicySegmentSecurityProfile(),
			"nsxt_policy_service":                                      resourceNsxtPolicyService(),
			"nsxt_policy_share":                                        resourceNsxtPolicyShare(),
			"nsxt_policy_shared_resource":                              resourceNsxtPolicySharedResource(),
			"nsxt_policy_site":                                         resourceNsxtPolicySite(),
			"nsxt_policy_spoof_guard_profile":                          resourceNsxtPolicySpoofGuardProfile("Please use nsxt_policy_spoofguard_profile resource instead"),
			"nsxt_policy_spoofguard_profile":                           resourceNsxtPolicySpoofGuardProfile(""),
			"nsxt_policy_static_route":                                 resourceNsxtPolicyStaticRoute(),
			"nsxt_policy_static_route_bfd_peer":                        resourceNsxtPolicyStaticRouteBfdPeer(),
			"nsxt_policy_tier0_gateway":                                resourceNsxtPolicyTier0Gateway(),
			"nsxt_policy_tier0_gateway_gre_tunnel":                     resourceNsxtPolicyTier0GatewayGRETunnel(),
			"nsxt_policy_tier0_gateway_ha_vip_config":                  resourceNsxtPolicyTier0GatewayHAVipConfig(),
			"nsxt_policy_tier0_gateway_interface":                      resourceNsxtPolicyTier0GatewayInterface(),
			"nsxt_policy_tier0_inter_vrf_routing":                      resourceNsxtPolicyTier0InterVRFRouting(),
			"nsxt_policy_tier1_gateway":                                resourceNsxtPolicyTier1Gateway(),
			"nsxt_policy_tier1_gateway_interface":                      resourceNsxtPolicyTier1GatewayInterface(),
			"nsxt_policy_transit_gateway":                              resourceNsxtPolicyTransitGateway(),
			"nsxt_policy_transit_gateway_attachment":                   resourceNsxtPolicyTransitGatewayAttachment(),
			"nsxt_policy_transit_gateway_nat_rule":                     resourceNsxtPolicyTransitGatewayNatRule(),
			"nsxt_policy_transit_gateway_static_route":                 resourceNsxtPolicyTransitGatewayStaticRoute(),
			"nsxt_policy_transport_zone":                               resourceNsxtPolicyTransportZone(),
			"nsxt_policy_uplink_host_switch_profile":                   resourceNsxtUplinkHostSwitchProfile(),
			"nsxt_policy_user_management_role":                         resourceNsxtPolicyUserManagementRole(),
			"nsxt_policy_user_management_role_binding":                 resourceNsxtPolicyUserManagementRoleBinding(),
			"nsxt_policy_vlan_segment":                                 resourceNsxtPolicyVlanSegment(),
			"nsxt_policy_vm_tags":                                      resourceNsxtPolicyVMTags(),
			"nsxt_policy_vni_pool":                                     resourceNsxtPolicyVniPool(),
			"nsxt_policy_vtep_ha_host_switch_profile":                  resourceNsxtVtepHAHostSwitchProfile(),
			"nsxt_principal_identity":                                  resourceNsxtPrincipalIdentity(),
			"nsxt_qos_switching_profile":                               removedResourceWrapper(resourceNsxtQosSwitchingProfile, "nsxt_qos_switching_profile"),
			"nsxt_spoofguard_switching_profile":                        removedResourceWrapper(resourceNsxtSpoofGuardSwitchingProfile, "nsxt_spoofguard_switching_profile"),
			"nsxt_static_route":                                        removedResourceWrapper(resourceNsxtStaticRoute, "nsxt_static_route"),
			"nsxt_switch_security_switching_profile":                   removedResourceWrapper(resourceNsxtSwitchSecuritySwitchingProfile, "nsxt_switch_security_switching_profile"),
			"nsxt_upgrade_precheck_acknowledge":                        resourceNsxtUpgradePrecheckAcknowledge(),
			"nsxt_upgrade_prepare":                                     resourceNsxtUpgradePrepare(),
			"nsxt_upgrade_run":                                         resourceNsxtUpgradeRun(),
			"nsxt_vlan_logical_switch":                                 removedResourceWrapper(resourceNsxtVlanLogicalSwitch, "nsxt_vlan_logical_switch"),
			"nsxt_vm_tags":                                             removedResourceWrapper(resourceNsxtVMTags, "nsxt_vm_tags"),
			"nsxt_vpc":                                                 resourceNsxtVpc(),
			"nsxt_vpc_attachment":                                      resourceNsxtVpcAttachment(),
			"nsxt_vpc_connectivity_profile":                            resourceNsxtVpcConnectivityProfile(),
			"nsxt_vpc_dhcp_v4_static_binding":                          resourceNsxtVpcSubnetDhcpV4StaticBindingConfig(),
			"nsxt_vpc_external_address":                                resourceNsxtVpcExternalAddress(),
			"nsxt_vpc_gateway_policy":                                  resourceNsxtVPCGatewayPolicy(),
			"nsxt_vpc_group":                                           resourceNsxtVPCGroup(),
			"nsxt_vpc_ip_address_allocation":                           resourceNsxtVpcIpAddressAllocation(),
			"nsxt_vpc_nat_rule":                                        resourceNsxtPolicyVpcNatRule(),
			"nsxt_vpc_security_policy":                                 resourceNsxtVPCSecurityPolicy(),
			"nsxt_vpc_service_profile":                                 resourceNsxtVpcServiceProfile(),
			"nsxt_vpc_static_route":                                    resourceNsxtVpcStaticRoutes(),
			"nsxt_vpc_subnet":                                          resourceNsxtVpcSubnet(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func isVMCCredentialSet(d *schema.ResourceData) bool {
	// Refresh token
	vmcToken := d.Get("vmc_token").(string)
	if len(vmcToken) > 0 {
		return true
	}

	// Oauth app
	vmcClientID := d.Get("vmc_client_id").(string)
	vmcClientSecret := d.Get("vmc_client_secret").(string)
	if len(vmcClientSecret) > 0 && len(vmcClientID) > 0 {
		return true
	}

	return false
}

func configureNsxtClient(d *schema.ResourceData, clients *nsxtClients) error {
	onDemandConn := d.Get("on_demand_connection").(bool)
	clientAuthCertFile := d.Get("client_auth_cert_file").(string)
	clientAuthKeyFile := d.Get("client_auth_key_file").(string)
	clientAuthCert := d.Get("client_auth_cert").(string)
	clientAuthKey := d.Get("client_auth_key").(string)
	vmcAuthMode := d.Get("vmc_auth_mode").(string)

	if onDemandConn {
		// On demand connection option is not supported with old SDK
		return nil
	}

	if (vmcAuthMode == "Basic") || isVMCCredentialSet(d) {
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

	// The correct place to escape special chars would be inside the SDK
	// However since the SDK is deprecated, we implement escaping here
	// TODO implement this functionality with new mp-sdk
	username = url.QueryEscape(username)
	password = url.QueryEscape(password)

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
	sessionAuth := d.Get("session_auth").(bool)
	skipSessionAuth := !sessionAuth

	retriesConfig := api.ClientRetriesConfiguration{
		MaxRetries:      clients.CommonConfig.MaxRetries,
		RetryMinDelay:   clients.CommonConfig.MinRetryInterval,
		RetryMaxDelay:   clients.CommonConfig.MaxRetryInterval,
		RetryOnStatuses: clients.CommonConfig.RetryStatusCodes,
	}

	clients.NsxtClientConfig = &api.Configuration{
		BasePath:             "/api/v1",
		Host:                 host,
		Scheme:               "https",
		UserAgent:            "terraform-provider-nsxt",
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
		SkipSessionAuth:      skipSessionAuth,
	}

	nsxClient, err := api.NewAPIClient(clients.NsxtClientConfig)
	if err != nil {
		return err
	}

	clients.NsxtClient = nsxClient

	return nil
}

type jwtToken struct {
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type vmcAuthInfo struct {
	authHost     string
	authMode     string
	accessToken  string
	clientID     string
	clientSecret string
}

func getVmcAuthInfo(d *schema.ResourceData) *vmcAuthInfo {
	vmcInfo := vmcAuthInfo{
		authHost:     d.Get("vmc_auth_host").(string),
		authMode:     d.Get("vmc_auth_mode").(string),
		accessToken:  d.Get("vmc_token").(string),
		clientID:     d.Get("vmc_client_id").(string),
		clientSecret: d.Get("vmc_client_secret").(string),
	}
	if len(vmcInfo.authHost) > 0 {
		return &vmcInfo
	}

	// Fill in default auth host + url based on auth method
	if len(vmcInfo.accessToken) > 0 {
		vmcInfo.authHost = "console.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize"
	} else if len(vmcInfo.clientSecret) > 0 && len(vmcInfo.clientID) > 0 {
		vmcInfo.authHost = "console.cloud.vmware.com/csp/gateway/am/api/auth/authorize"
	}
	return &vmcInfo
}

func (v *vmcAuthInfo) IsZero() bool {
	return len(v.accessToken) == 0 && len(v.clientID) == 0 && len(v.clientSecret) == 0
}

func (v *vmcAuthInfo) getAPIToken() (string, error) {
	var req *http.Request

	// Access token
	if len(v.accessToken) > 0 {
		payload := strings.NewReader("refresh_token=" + v.accessToken)
		req, _ = http.NewRequest("POST", "https://"+v.authHost, payload)
	}
	// Oauth app
	if len(v.clientSecret) > 0 && len(v.clientID) > 0 {
		payload := strings.NewReader("grant_type=client_credentials")
		req, _ = http.NewRequest("POST", "https://"+v.authHost, payload)
		req.SetBasicAuth(v.clientID, v.clientSecret)
	}
	if req == nil {
		return "", fmt.Errorf("invalid VMC auth input")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("unexpected status code %d trying to get auth token. %s", res.StatusCode, string(b))
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
		caCert, err := os.ReadFile(caFile)
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
	onDemandConn := d.Get("on_demand_connection").(bool)
	host := d.Get("host").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	clientAuthCertFile := d.Get("client_auth_cert_file").(string)
	clientAuthCert := d.Get("client_auth_cert").(string)
	clientAuthDefined := (len(clientAuthCertFile) > 0) || (len(clientAuthCert) > 0)
	policyEnforcementPoint := d.Get("enforcement_point").(string)
	policyGlobalManager := d.Get("global_manager").(bool)
	vmcInfo := getVmcAuthInfo(d)

	isVMC := false
	if (vmcInfo.authMode == "Basic") || isVMCCredentialSet(d) {
		isVMC = true
		if onDemandConn {
			return fmt.Errorf("on demand connection option is not supported with VMC")
		}
	}

	if host == "" {
		return fmt.Errorf("host must be provided")
	}

	if !strings.HasPrefix(host, "https://") {
		host = fmt.Sprintf("https://%s", host)
	}

	securityContextNeeded := !(clientAuthDefined && !clients.CommonConfig.RemoteAuth)

	if securityContextNeeded {
		securityCtx, err := getConfiguredSecurityContext(clients, vmcInfo, username, password)
		if err != nil {
			return err
		}
		clients.PolicySecurityContext = securityCtx
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
	clients.Host = host
	clients.PolicyEnforcementPoint = policyEnforcementPoint
	clients.PolicyGlobalManager = policyGlobalManager

	if onDemandConn {
		// version init will happen on demand
		return nil
	}

	if !isVMC {
		err = configureLicenses(getStandalonePolicyConnector(*clients, true), clients.CommonConfig.LicenseKeys)
		if err != nil {
			return err
		}
	}

	err = initNSXVersion(getStandalonePolicyConnector(*clients, true))
	if err != nil && isVMC {
		// In case version API does not work for VMC, we workaround by testing version-specific APIs
		// TODO - remove this when /node/version API works for all auth methods on VMC
		initNSXVersionVMC(*clients)
		return nil
	}
	return err
}

func getConfiguredSecurityContext(clients *nsxtClients, vmcInfo *vmcAuthInfo, username string, password string) (*core.SecurityContextImpl, error) {
	securityCtx := core.NewSecurityContextImpl()
	if vmcInfo == nil || vmcInfo.IsZero() {
		if username == "" {
			return nil, fmt.Errorf("username must be provided")
		}

		if password == "" {
			return nil, fmt.Errorf("password must be provided")
		}

		securityCtx.SetProperty(security.AUTHENTICATION_SCHEME_ID, security.USER_PASSWORD_SCHEME_ID)
		securityCtx.SetProperty(security.USER_KEY, username)
		securityCtx.SetProperty(security.PASSWORD_KEY, password)
	} else {
		apiToken, err := vmcInfo.getAPIToken()
		if err != nil {
			return nil, err
		}

		// We'll be sending Bearer token anyway even with scp-auth-token auth
		// For now, node API is not working on VMC without Bearer token present
		clients.CommonConfig.BearerToken = apiToken
		if vmcInfo.authMode != "Bearer" {
			securityCtx.SetProperty(security.AUTHENTICATION_SCHEME_ID, security.OAUTH_SCHEME_ID)
			securityCtx.SetProperty(security.ACCESS_TOKEN, apiToken)
		}
	}
	return securityCtx, nil
}

type customHeaderProcessor struct {
	customHeaders *map[string]string
}

func (processor customHeaderProcessor) Process(req *http.Request) error {
	for header, value := range *processor.customHeaders {
		req.Header.Set(header, value)
	}
	return nil
}

func newCustomHeaderProcessor(customHeaders *map[string]string) *customHeaderProcessor {
	return &customHeaderProcessor{customHeaders: customHeaders}
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

type logRequestProcessor struct {
}

func newLogRequestProcessor() *logRequestProcessor {
	return &logRequestProcessor{}
}

func (processor logRequestProcessor) Process(req *http.Request) error {

	filter := os.Getenv("TF_LOG_PROVIDER_NSX_HTTP_URI_FILTER")
	if filter != "" {
		// Only log for requests that match URL filter
		if !strings.Contains(req.URL.Path, filter) {
			return nil
		}
	}
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatal(err)
	}

	// Replace sensitive information in HTTP headers
	authHeaderRegexp := regexp.MustCompile(`(?i)Authorization:.*`)
	cspHeaderRegexp := regexp.MustCompile(`(?i)Csp-Auth-Token:.*`)
	replaced := authHeaderRegexp.ReplaceAllString(string(reqDump), "<Omitted Authorization header>")
	replaced = cspHeaderRegexp.ReplaceAllString(replaced, "<Omitted Csp-Auth-Token header>")

	log.Printf("Issuing request towards NSX:\n%s", replaced)
	return nil
}

type logResponseAcceptor struct {
}

func newLogResponseAcceptor() *logResponseAcceptor {
	return &logResponseAcceptor{}
}

func (processor logResponseAcceptor) Accept(resp *http.Response) {
	filter := os.Getenv("TF_LOG_PROVIDER_NSX_HTTP_URI_FILTER")
	if filter != "" && resp.Request != nil {
		// Only log for requests that match URI filter
		if !strings.Contains(resp.Request.URL.Path, filter) {
			return
		}
	}

	dumpResponse, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Received NSX response:\n%s", dumpResponse)
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

type sessionHeaderProcessor struct {
	cookie string
	xsrf   string
}

func newSessionHeaderProcessor(cookie string, xsrf string) *sessionHeaderProcessor {
	return &sessionHeaderProcessor{
		cookie: cookie,
		xsrf:   xsrf,
	}
}

func (processor sessionHeaderProcessor) Process(req *http.Request) error {
	req.Header.Set("Cookie", processor.cookie)
	req.Header.Set("X-XSRF-TOKEN", processor.xsrf)
	return nil
}

func getLicenses(connector client.Connector) ([]string, error) {
	var licenseList []string
	client := nsx.NewLicensesClient(connector)
	list, err := client.List()
	if err != nil {
		return licenseList, fmt.Errorf("Error during license create: %v", err)
	}

	defaultLicenseMarkers := []string{"NSX for vShield Endpoint"}
	for _, item := range list.Results {
		// Ignore default licenses
		if item.Description != nil && slices.Contains(defaultLicenseMarkers, *item.Description) {
			continue
		}
		licenseList = append(licenseList, *item.LicenseKey)
	}

	return licenseList, nil
}

func applyLicense(connector client.Connector, licenseKey string) error {
	client := nsx.NewLicensesClient(connector)
	license := model.License{LicenseKey: &licenseKey}
	_, err := client.Create(license)
	if err != nil {
		return fmt.Errorf("Error during license create: %v", err)
	}

	return nil
}

// license keys are applied on terraform plan and are not removed
func configureLicenses(connector client.Connector, intentLicenses []string) error {
	if len(intentLicenses) == 0 {
		// Since we never remove licenses, nothing to do here
		return nil
	}
	existingLicenses, err := getLicenses(connector)
	if err != nil {
		return err
	}
	// Apply new licenses
	for _, license := range intentLicenses {
		if slices.Contains(existingLicenses, license) {
			continue
		}
		err := applyLicense(connector, license)
		if err != nil {
			return fmt.Errorf("error applying license key: %s. %s", license, err.Error())
		}
	}

	return nil
}

func initCommonConfig(d *schema.ResourceData) commonProviderConfig {
	remoteAuth := d.Get("remote_auth").(bool)
	toleratePartialSuccess := d.Get("tolerate_partial_success").(bool)
	maxRetries := d.Get("max_retries").(int)
	retryMinDelay := d.Get("retry_min_delay").(int)
	retryMaxDelay := d.Get("retry_max_delay").(int)
	if retryMaxDelay < retryMinDelay {
		retryMaxDelay = retryMinDelay + 1
	}
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	statuses := d.Get("retry_on_status_codes").([]interface{})
	retryStatuses := make([]int, 0, len(statuses))
	for _, s := range statuses {
		retryStatuses = append(retryStatuses, s.(int))
	}

	if len(retryStatuses) == 0 {
		// Set to the defaults if empty
		retryStatuses = append(retryStatuses, defaultRetryOnStatusCodes...)
	}

	licenses := interfaceListToStringList(d.Get("license_keys").([]interface{}))
	return commonProviderConfig{
		RemoteAuth:             remoteAuth,
		ToleratePartialSuccess: toleratePartialSuccess,
		MaxRetries:             maxRetries,
		MinRetryInterval:       retryMinDelay,
		MaxRetryInterval:       retryMaxDelay,
		RetryStatusCodes:       retryStatuses,
		Username:               username,
		Password:               password,
		LicenseKeys:            licenses,
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

	return clients, nil
}

// Standard policy connection that initializes global connection settings on demand
func getPolicyConnector(clients interface{}) client.Connector {
	return getPolicyConnectorWithHeaders(clients, nil, false, true)
}

// Standalone policy connector, possibly for different endpoint,
// for the purpose of special tasks (such as joining manager cluster node)
// Does not initialize global connection settings
func getStandalonePolicyConnector(clients interface{}, withRetry bool) client.Connector {
	return getPolicyConnectorWithHeaders(clients, nil, true, withRetry)
}

func getPolicyConnectorWithHeaders(clients interface{}, customHeaders *map[string]string, standaloneFlow bool, withRetry bool) client.Connector {
	c := clients.(nsxtClients)

	retryFunc := func(retryContext retry.RetryContext) bool {
		shouldRetry := false
		if retryContext.Response != nil {
			for _, code := range c.CommonConfig.RetryStatusCodes {
				if retryContext.Response.StatusCode == code {
					log.Printf("[DEBUG]: Retrying request due to error code %d", code)
					shouldRetry = true
					break
				}
			}
		} else {
			shouldRetry = true
			log.Printf("[DEBUG]: Retrying request due to error")
		}

		if !shouldRetry {
			return false
		}

		minRetryInterval := c.CommonConfig.MinRetryInterval
		maxRetryInterval := c.CommonConfig.MaxRetryInterval
		if maxRetryInterval > 0 {
			interval := (rand.Intn(maxRetryInterval-minRetryInterval) + minRetryInterval)
			time.Sleep(time.Duration(interval) * time.Millisecond)
			log.Printf("[DEBUG]: Waited %d ms before retrying", interval)
		}

		return true
	}

	connectorOptions := []client.ConnectorOption{client.UsingRest(nil), client.WithHttpClient(c.PolicyHTTPClient)}
	var requestProcessors []core.RequestProcessor
	var responseAcceptors []core.ResponseAcceptor

	if withRetry {
		connectorOptions = append(connectorOptions, client.WithDecorators(retry.NewRetryDecorator(uint(c.CommonConfig.MaxRetries), retryFunc)))
	}

	if c.PolicySecurityContext != nil {
		connectorOptions = append(connectorOptions, client.WithSecurityContext(c.PolicySecurityContext))
	}
	if c.CommonConfig.RemoteAuth {
		requestProcessors = append(requestProcessors, newRemoteAuthHeaderProcessor().Process)
	}
	if len(c.CommonConfig.BearerToken) > 0 {
		requestProcessors = append(requestProcessors, newBearerAuthHeaderProcessor(c.CommonConfig.BearerToken).Process)
	}
	if customHeaders != nil {
		requestProcessors = append(requestProcessors, newCustomHeaderProcessor(customHeaders).Process)
	}

	// Session support for policy resources (main rationale - vIDM environment where auth is slow)
	// Currently session creation is done via old MP sdk.
	// TODO - when MP resources are removed, switch to official SDK to initiate session/create API
	// TODO - re-trigger session/create when token is expired
	if c.NsxtClientConfig != nil && len(c.NsxtClientConfig.DefaultHeader["Cookie"]) > 0 {
		cookie := c.NsxtClientConfig.DefaultHeader["Cookie"]
		xsrf := ""
		if len(c.NsxtClientConfig.DefaultHeader["X-XSRF-TOKEN"]) > 0 {
			xsrf = c.NsxtClientConfig.DefaultHeader["X-XSRF-TOKEN"]
		}
		requestProcessors = append(requestProcessors, newSessionHeaderProcessor(cookie, xsrf).Process)
		log.Printf("[INFO]: Session headers configured for policy objects")
	}

	if os.Getenv("TF_LOG_PROVIDER_NSX_HTTP") != "" {
		requestProcessors = append(requestProcessors, newLogRequestProcessor().Process)
		responseAcceptors = append(responseAcceptors, newLogResponseAcceptor().Accept)
	}

	if len(requestProcessors) > 0 {
		connectorOptions = append(connectorOptions, client.WithRequestProcessors(requestProcessors...))
	}
	if len(responseAcceptors) > 0 {
		connectorOptions = append(connectorOptions, client.WithResponseAcceptors(responseAcceptors...))
	}
	connector := client.NewConnector(c.Host, connectorOptions...)
	// Init NSX version on demand if not done yet
	// This is also our indication to apply licenses, in case of delayed connection
	// This step is skipped if the connector is for special purpose, or for different endpoint
	if util.NsxVersion == "" && !standaloneFlow {
		initNSXVersion(connector)
		err := configureLicenses(connector, c.CommonConfig.LicenseKeys)
		if err != nil {
			log.Printf("[ERROR]: Failed to apply NSX licenses")
		}
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

func getContextDataFromSchema(d *schema.ResourceData) (string, string, bool) {
	ctxPtr := d.Get("context")
	if ctxPtr != nil {
		contexts := ctxPtr.([]interface{})
		for _, context := range contexts {
			data := context.(map[string]interface{})
			vpcID := ""
			fromGlobal := false
			if data["vpc_id"] != nil {
				vpcID = data["vpc_id"].(string)
			}
			if data["from_global"] != nil {
				fromGlobal = data["from_global"].(bool)
			}
			return data["project_id"].(string), vpcID, fromGlobal
		}
	}
	return "", "", false
}

func getContextDataFromParentPath(parentPath string) (string, string) {
	segments := strings.Split(parentPath, "/")
	var projectID, vpcID string
	if len(segments) > 4 && segments[1] == "orgs" && segments[3] == "projects" {
		projectID = segments[4]

		if len(segments) > 6 && segments[5] == "vpcs" {
			vpcID = segments[6]
		}
	}
	return projectID, vpcID
}

func getSessionContext(d *schema.ResourceData, m interface{}) tf_api.SessionContext {
	return getSessionContextHelper(d, m, "")
}

func getParentContext(d *schema.ResourceData, m interface{}, parentPath string) tf_api.SessionContext {
	return getSessionContextHelper(d, m, parentPath)
}

func getSessionContextHelper(d *schema.ResourceData, m interface{}, parentPath string) tf_api.SessionContext {
	var clientType tf_api.ClientType
	var projectID, vpcID string
	fromGlobal := false
	if parentPath == "" {
		projectID, vpcID, fromGlobal = getContextDataFromSchema(d)
	} else {
		projectID, vpcID = getContextDataFromParentPath(parentPath)
	}
	// if fromGlobal = True values pertaining multitenancy should be ignored
	// and clientType must be tf_api.Local
	if fromGlobal {
		projectID = ""
		vpcID = ""
	}
	if projectID != "" {
		clientType = tf_api.Multitenancy
		if vpcID != "" {
			clientType = tf_api.VPC
		}
	} else if isPolicyGlobalManager(m) {
		clientType = tf_api.Global
	} else {
		clientType = tf_api.Local
	}
	return tf_api.SessionContext{ProjectID: projectID, VPCID: vpcID, ClientType: clientType, FromGlobal: fromGlobal}
}
