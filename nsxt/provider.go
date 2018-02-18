/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_INSECURE", false),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_PASSWORD", nil),
				Sensitive:   true,
			},
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_MANAGER_HOST", nil),
			},
			"client_auth_cert_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_CLIENT_AUTH_CERT_FILE", nil),
			},
			"client_auth_key_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_CLIENT_AUTH_KEY_FILE", nil),
			},
			"ca_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NSX_CA_FILE", nil),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"nsxt_transport_zone":       dataSourceTransportZone(),
			"nsxt_switching_profile":    dataSourceSwitchingProfile(),
			"nsxt_logical_tier0_router": dataSourceLogicalTier0Router(),
			"nsxt_edge_cluster":         dataSourceEdgeCluster(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"nsxt_dhcp_relay_profile":                resourceDhcpRelayProfile(),
			"nsxt_dhcp_relay_service":                resourceDhcpRelayService(),
			"nsxt_logical_switch":                    resourceLogicalSwitch(),
			"nsxt_logical_port":                      resourceLogicalPort(),
			"nsxt_logical_tier1_router":              resourceLogicalTier1Router(),
			"nsxt_logical_router_downlink_port":      resourceLogicalRouterDownLinkPort(),
			"nsxt_logical_router_link_port_on_tier0": resourceLogicalRouterLinkPortOnTier0(),
			"nsxt_logical_router_link_port_on_tier1": resourceLogicalRouterLinkPortOnTier1(),
			"nsxt_l4_port_set_ns_service":            resourceL4PortSetNsService(),
			"nsxt_algorithm_type_ns_service":         resourceAlgorithmTypeNsService(),
			"nsxt_icmp_type_ns_service":              resourceIcmpTypeNsService(),
			"nsxt_igmp_type_ns_service":              resourceIgmpTypeNsService(),
			"nsxt_ether_type_ns_service":             resourceEtherTypeNsService(),
			"nsxt_ip_protocol_ns_service":            resourceIpProtocolNsService(),
			"nsxt_ns_group":                          resourceNsGroup(),
			"nsxt_firewall_section":                  resourceFirewallSection(),
			"nsxt_nat_rule":                          resourceNatRule(),
			"nsxt_ip_set":                            resourceIpSet(),
			"nsxt_static_route":                      resourceStaticRoute(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	insecure := d.Get("insecure").(bool)
	username := d.Get("username").(string)

	if username == "" {
		return nil, fmt.Errorf("username must be provided")
	}

	password := d.Get("password").(string)

	if password == "" {
		return nil, fmt.Errorf("password must be provided")
	}

	host := d.Get("host").(string)

	if host == "" {
		return nil, fmt.Errorf("host must be provided")
	}

	client_auth_cert_file := d.Get("client_auth_cert_file").(string)
	client_auth_key_file := d.Get("client_auth_key_file").(string)
	ca_file := d.Get("ca_file").(string)

	cfg := nsxt.Configuration{
		BasePath:           "/api/v1",
		Host:               host,
		Scheme:             "https",
		UserAgent:          "terraform-provider-nsxt/1.0",
		UserName:           username,
		Password:           password,
		ClientAuthCertFile: client_auth_cert_file,
		ClientAuthKeyFile:  client_auth_key_file,
		CAFile:             ca_file,
		Insecure:           insecure,
	}

	nsxClient, err := nsxt.NewAPIClient(&cfg)
	return nsxClient, err
}
