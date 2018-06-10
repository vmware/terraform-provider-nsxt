/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt"
)

var defaultRetryOnStatusCodes = []int{429, 503}

// Provider for VMWare NSX-T. Returns terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"allow_unverified_ssl": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_ALLOW_UNVERIFIED_SSL", false),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_PASSWORD", nil),
				Sensitive:   true,
			},
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_MANAGER_HOST", nil),
			},
			"client_auth_cert_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_CLIENT_AUTH_CERT_FILE", nil),
			},
			"client_auth_key_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_CLIENT_AUTH_KEY_FILE", nil),
			},
			"ca_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSXT_CA_FILE", nil),
			},
			"max_retries": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum number of HTTP client retries",
				DefaultFunc: schema.EnvDefaultFunc("NSXT_MAX_RETRIES", 50),
			},
			"retry_min_delay": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum delay in milliseconds between retries of a request",
				DefaultFunc: schema.EnvDefaultFunc("NSXT_RETRY_MIN_DELAY", 500),
			},
			"retry_max_delay": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum delay in milliseconds between retries of a request",
				DefaultFunc: schema.EnvDefaultFunc("NSXT_RETRY_MAX_DELAY", 5000),
			},
			"retry_on_status_codes": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "HTTP replies status codes to retry on",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				// There is no support for default values/func for list, so it will be handled later
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"nsxt_transport_zone":       dataSourceNsxtTransportZone(),
			"nsxt_switching_profile":    dataSourceNsxtSwitchingProfile(),
			"nsxt_logical_tier0_router": dataSourceNsxtLogicalTier0Router(),
			"nsxt_ns_service":           dataSourceNsxtNsService(),
			"nsxt_edge_cluster":         dataSourceNsxtEdgeCluster(),
			"nsxt_certificate":          dataSourceNsxtCertificate(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"nsxt_dhcp_relay_profile":                resourceNsxtDhcpRelayProfile(),
			"nsxt_dhcp_relay_service":                resourceNsxtDhcpRelayService(),
			"nsxt_dhcp_server_profile":               resourceNsxtDhcpServerProfile(),
			"nsxt_logical_dhcp_server":               resourceNsxtLogicalDhcpServer(),
			"nsxt_logical_switch":                    resourceNsxtLogicalSwitch(),
			"nsxt_logical_dhcp_port":                 resourceNsxtLogicalDhcpPort(),
			"nsxt_logical_port":                      resourceNsxtLogicalPort(),
			"nsxt_logical_tier0_router":              resourceNsxtLogicalTier0Router(),
			"nsxt_logical_tier1_router":              resourceNsxtLogicalTier1Router(),
			"nsxt_logical_router_downlink_port":      resourceNsxtLogicalRouterDownLinkPort(),
			"nsxt_logical_router_link_port_on_tier0": resourceNsxtLogicalRouterLinkPortOnTier0(),
			"nsxt_logical_router_link_port_on_tier1": resourceNsxtLogicalRouterLinkPortOnTier1(),
			"nsxt_l4_port_set_ns_service":            resourceNsxtL4PortSetNsService(),
			"nsxt_algorithm_type_ns_service":         resourceNsxtAlgorithmTypeNsService(),
			"nsxt_icmp_type_ns_service":              resourceNsxtIcmpTypeNsService(),
			"nsxt_igmp_type_ns_service":              resourceNsxtIgmpTypeNsService(),
			"nsxt_ether_type_ns_service":             resourceNsxtEtherTypeNsService(),
			"nsxt_ip_protocol_ns_service":            resourceNsxtIPProtocolNsService(),
			"nsxt_ns_group":                          resourceNsxtNsGroup(),
			"nsxt_firewall_section":                  resourceNsxtFirewallSection(),
			"nsxt_nat_rule":                          resourceNsxtNatRule(),
			"nsxt_ip_block":                          resourceNsxtIPBlock(),
			"nsxt_ip_block_subnet":                   resourceNsxtIPBlockSubnet(),
			"nsxt_ip_pool":                           resourceNsxtIPPool(),
			"nsxt_ip_set":                            resourceNsxtIPSet(),
			"nsxt_static_route":                      resourceNsxtStaticRoute(),
			"nsxt_vm_tags":                           resourceNsxtVMTags(),
			"nsxt_lb_icmp_monitor":                   resourceNsxtLbIcmpMonitor(),
			"nsxt_lb_tcp_monitor":                    resourceNsxtLbTCPMonitor(),
			"nsxt_lb_udp_monitor":                    resourceNsxtLbUDPMonitor(),
			"nsxt_lb_http_monitor":                   resourceNsxtLbHTTPMonitor(),
			"nsxt_lb_https_monitor":                  resourceNsxtLbHTTPSMonitor(),
			"nsxt_lb_passive_monitor":                resourceNsxtLbPassiveMonitor(),
			"nsxt_lb_pool":                           resourceNsxtLbPool(),
			"nsxt_lb_cookie_persistence_profile":     resourceNsxtLbCookiePersistenceProfile(),
			"nsxt_lb_client_ssl_profile":             resourceNsxtLbClientSslProfile(),
			"nsxt_lb_server_ssl_profile":             resourceNsxtLbServerSslProfile(),
			"nsxt_lb_http_application_profile":       resourceNsxtLbHttpApplicationProfile(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConnectivityCheck(nsxClient *nsxt.APIClient) error {
	// Connectivity check - get a random object to check connectivity and credentials
	// TODO(asarfaty): Use a list command which returns the full body, when the go vendor has one.
	_, httpResponse, err := nsxClient.ServicesApi.ReadLoadBalancerPool(nsxClient.Context, "Dummy")
	if err != nil {
		if httpResponse == nil || (httpResponse.StatusCode == 401 || httpResponse.StatusCode == 403) {
			return fmt.Errorf("NSXT provider connectivity check failed: %s", err)
		}
	}
	return nil
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	clientAuthCertFile := d.Get("client_auth_cert_file").(string)
	clientAuthKeyFile := d.Get("client_auth_key_file").(string)

	needCreds := true
	if len(clientAuthCertFile) > 0 && len(clientAuthKeyFile) > 0 {
		needCreds = false
	}

	insecure := d.Get("allow_unverified_ssl").(bool)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	if needCreds {
		if username == "" {
			return nil, fmt.Errorf("username must be provided")
		}

		if password == "" {
			return nil, fmt.Errorf("password must be provided")
		}
	}

	host := d.Get("host").(string)

	if host == "" {
		return nil, fmt.Errorf("host must be provided")
	}

	caFile := d.Get("ca_file").(string)

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

	retriesConfig := nsxt.ClientRetriesConfiguration{
		MaxRetries:      maxRetries,
		RetryMinDelay:   retryMinDelay,
		RetryMaxDelay:   retryMaxDelay,
		RetryOnStatuses: retryStatuses,
	}

	cfg := nsxt.Configuration{
		BasePath:             "/api/v1",
		Host:                 host,
		Scheme:               "https",
		UserAgent:            "terraform-provider-nsxt/1.0",
		UserName:             username,
		Password:             password,
		ClientAuthCertFile:   clientAuthCertFile,
		ClientAuthKeyFile:    clientAuthKeyFile,
		CAFile:               caFile,
		Insecure:             insecure,
		RetriesConfiguration: retriesConfig,
	}

	nsxClient, err := nsxt.NewAPIClient(&cfg)
	if err != nil {
		return nil, err
	}
	// Check provider connectivity
	err = providerConnectivityCheck(nsxClient)
	if err != nil {
		return nil, err
	}

	return nsxClient, nil
}
