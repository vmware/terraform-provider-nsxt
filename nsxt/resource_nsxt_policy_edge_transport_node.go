// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var policyEdgeNodeFormFactorValues = []string{
	model.PolicyEdgeTransportNode_FORM_FACTOR_SMALL,
	model.PolicyEdgeTransportNode_FORM_FACTOR_MEDIUM,
	model.PolicyEdgeTransportNode_FORM_FACTOR_LARGE,
	model.PolicyEdgeTransportNode_FORM_FACTOR_XLARGE,
}

var policySyslogLogLevelValues = []string{
	model.SyslogConfiguration_LOG_LEVEL_EMERGENCY,
	model.SyslogConfiguration_LOG_LEVEL_ALERT,
	model.SyslogConfiguration_LOG_LEVEL_CRITICAL,
	model.SyslogConfiguration_LOG_LEVEL_ERROR,
	model.SyslogConfiguration_LOG_LEVEL_WARNING,
	model.SyslogConfiguration_LOG_LEVEL_NOTICE,
	model.SyslogConfiguration_LOG_LEVEL_INFO,
	model.SyslogConfiguration_LOG_LEVEL_DEBUG,
}

var policySyslogProtocolValues = []string{
	model.SyslogConfiguration_PROTOCOL_TCP,
	model.SyslogConfiguration_PROTOCOL_UDP,
	model.SyslogConfiguration_PROTOCOL_TLS,
	model.SyslogConfiguration_PROTOCOL_LI,
	model.SyslogConfiguration_PROTOCOL_LI_TLS,
}

var policyCpuReservationValues = []string{
	model.CPUReservation_RESERVATION_IN_SHARES_EXTRA_HIGH_PRIORITY,
	model.CPUReservation_RESERVATION_IN_SHARES_HIGH_PRIORITY,
	model.CPUReservation_RESERVATION_IN_SHARES_NORMAL_PRIORITY,
	model.CPUReservation_RESERVATION_IN_SHARES_LOW_PRIORITY,
}

var managementAssignments = []string{
	"dhcp_v4",
	"static_ipv4",
	"static_ipv6",
}

var tepAssignments = []string{
	"dhcp_v4",
	"dhcp_v6",
	"static_ipv4_list",
	"static_ipv4_pool",
	"static_ipv6_list",
	"static_ipv6_pool",
}

func getPolicyIPAssignmentSchema(required bool, minItems, maxItems int, validAssignments []string) *schema.Schema {
	allAssignments := map[string]*schema.Schema{
		"auto_conf": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable Auto-conf based IPv6 assignment",
		},
		"dhcp_v4": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable DHCP based IPv4 assignment",
		},
		"dhcp_v6": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable DHCP based IPv6 assignment",
		},
		"no_assignment": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable no IP assignment",
		},
		"static_ipv4": {
			Type:        schema.TypeList,
			Description: "IP assignment specification for a Static IP",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"default_gateway": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Default IPv4 gateway for the node",
						ValidateFunc: validateSingleIP(),
					},
					"management_port_subnet": {
						Type:        schema.TypeList,
						Required:    true,
						MaxItems:    1,
						MinItems:    1,
						Description: "IPv4 Port subnet for management port",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ip_addresses": {
									Type:        schema.TypeList,
									Required:    true,
									Description: "IPv4 Addresses",
									MinItems:    1,
									MaxItems:    1,
									Elem: &schema.Schema{
										Type:         schema.TypeString,
										ValidateFunc: validation.IsIPv4Address,
									},
								},
								"prefix_length": {
									Type:         schema.TypeInt,
									Required:     true,
									Description:  "Subnet Prefix Length",
									ValidateFunc: validation.IntBetween(1, 32),
								},
							},
						},
					},
				},
			},
		},
		"static_ipv4_list": {
			Type:        schema.TypeList,
			Description: "IP assignment specification value for Static IPv4 List",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"default_gateway": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Gateway IP",
						ValidateFunc: validation.IsIPv4Address,
					},
					"ip_addresses": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "List of IPV4 addresses for edge transport node host switch virtual tunnel endpoints",
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.IsIPv4Address,
						},
					},
					"subnet_mask": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Subnet mask",
						ValidateFunc: validation.IsIPv4Address,
					},
				},
			},
		},
		"static_ipv4_mac_list": {
			Type:        schema.TypeList,
			Description: "IP and MAC assignment specification for Static IP List",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"default_gateway": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Gateway IP",
						ValidateFunc: validation.IsIPv4Address,
					},
					"ip_mac_pair": getIpMacPairSchema(),
					"subnet_mask": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Subnet mask",
						ValidateFunc: validation.IsIPv4Address,
					},
				},
			},
		},
		"static_ipv4_pool": getPolicyPathSchema(false, false, "IP assignment specification for Static IPv4 Pool. Input should be the policy path of IP pool"),
		"static_ipv6": {
			Type:        schema.TypeList,
			Description: "IP assignment specification for a Static IP",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"default_gateway": {
						Type:         schema.TypeString,
						Optional:     true,
						Description:  "Default IPv6 gateway for the node",
						ValidateFunc: validation.IsIPv6Address,
					},
					"management_port_subnet": {
						Type:        schema.TypeList,
						Required:    true,
						MaxItems:    1,
						MinItems:    1,
						Description: "Ipv6 Port subnet for management port",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ip_addresses": {
									Type:        schema.TypeList,
									Required:    true,
									Description: "IPv6 Addresses",
									MinItems:    1,
									MaxItems:    1,
									Elem: &schema.Schema{
										Type:         schema.TypeString,
										ValidateFunc: validation.IsIPv6Address,
									},
								},
								"prefix_length": {
									Type:         schema.TypeInt,
									Required:     true,
									Description:  "Subnet Prefix Length",
									ValidateFunc: validation.IntBetween(1, 128),
								},
							},
						},
					},
				},
			},
		},
		"static_ipv6_list": {
			Type:        schema.TypeList,
			Description: "IP assignment specification value for Static IPv4 List",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"default_gateway": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Gateway IP",
						ValidateFunc: validation.IsIPv6Address,
					},
					"ip_addresses": {
						Type:        schema.TypeList,
						Required:    true,
						Description: "List of IPv6 IPs for edge transport node host switch virtual tunnel endpoints",
						MaxItems:    32,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.IsIPv6Address,
						},
					},
					"prefix_length": {
						Type:         schema.TypeInt,
						Required:     true,
						Description:  "Prefix Length",
						ValidateFunc: validation.IntBetween(1, 128),
					},
				},
			},
		},
		"static_ipv6_mac_list": {
			Type:        schema.TypeList,
			Description: "IP and MAC assignment specification for Static IPv6 List",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"default_gateway": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "Gateway IP",
						ValidateFunc: validation.IsIPv6Address,
					},
					"ip_mac_pair": getIpMacPairSchema(),
					"prefix_length": {
						Type:        schema.TypeInt,
						Required:    true,
						Description: "Prefix Length",
					},
				},
			},
		},
		"static_ipv6_pool": getPolicyPathSchema(false, false, "IP assignment specification for Static IPv6 Pool. Input should be the policy path of IP pool"),
	}

	assignments := make(map[string]*schema.Schema)
	for _, validAssignment := range validAssignments {
		assignments[validAssignment] = allAssignments[validAssignment]
	}
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "IPv4 or Ipv6 Port subnets for management port",
		MinItems:    minItems,
		MaxItems:    maxItems,
		Required:    required,
		Optional:    !required,
		Elem: &schema.Resource{
			Schema: assignments,
		},
	}
}

func resourceNsxtPolicyEdgeTransportNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyEdgeTransportNodeCreate,
		Read:   resourceNsxtPolicyEdgeTransportNodeRead,
		Update: resourceNsxtPolicyEdgeTransportNodeUpdate,
		Delete: resourceNsxtPolicyEdgeTransportNodeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNsxtPolicyEdgeTransportNodeImporter,
		},
		Schema: map[string]*schema.Schema{
			"nsx_id":       getNsxIDSchema(),
			"path":         getPathSchema(),
			"display_name": getDisplayNameSchema(),
			"description":  getDescriptionSchema(),
			"revision":     getRevisionSchema(),
			"tag":          getTagsSchema(),
			"site_path": {
				Type:         schema.TypeString,
				Description:  "Path to the site this Edge Transport Node belongs to",
				Optional:     true,
				ForceNew:     true,
				Default:      defaultInfraSitePath,
				ValidateFunc: validatePolicyPath(),
			},
			"enforcement_point": {
				Type:        schema.TypeString,
				Description: "ID of the enforcement point this Edge Transport Node belongs to",
				Optional:    true,
				ForceNew:    true,
				Default:     "default",
			},
			"node_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "Unique Id of the fabric node",
				ConflictsWith: []string{"advanced_configuration", "appliance_config", "credentials", "form_factor", "management_interface", "vm_deployment_config"},
			},
			"advanced_configuration": getKeyValuePairListSchema(),
			"appliance_config": {
				Type:        schema.TypeList,
				Description: "Applicance configuration",
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_ssh_root_login": {
							Type:        schema.TypeBool,
							Default:     false,
							Description: "Allow root SSH logins",
							Optional:    true,
						},
						"dns_servers": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "DNS servers",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateSingleIP(),
							},
							Default: []string{},
						},
						"enable_ssh": {
							Type:        schema.TypeBool,
							Default:     false,
							Description: "Enable SSH",
							Optional:    true,
						},
						"enable_upt_mode": {
							Type:        schema.TypeBool,
							Default:     false,
							Description: "Enable Uniform Passthrough mode",
							Optional:    true,
						},
						"syslog_server": {
							Type:        schema.TypeList,
							MaxItems:    5,
							Optional:    true,
							Description: "Syslog servers",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"log_level": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Log level to be redirected",
										Default:      model.SyslogConfiguration_LOG_LEVEL_INFO,
										ValidateFunc: validation.StringInSlice(policySyslogLogLevelValues, false),
									},
									"port": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "514",
										ValidateFunc: validateSinglePort(),
										Description:  "Syslog server port",
									},
									"protocol": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Syslog protocol",
										Default:      model.SyslogConfiguration_PROTOCOL_UDP,
										ValidateFunc: validation.StringInSlice(policySyslogProtocolValues, false),
									},
									"server": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Server IP or fqdn",
									},
								},
							},
						},
					},
				},
			},
			"credentials": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Username and password settings for the node.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audit_password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Node audit user password",
						},
						"audit_username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CLI \"audit\" username",
						},
						"cli_password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Node cli password",
						},
						"cli_username": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "admin",
							Description: "CLI \"admin\" username",
						},
						"root_password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Node root user password",
						},
					},
				},
			},
			"failure_domain_path": {
				Type:        schema.TypeString,
				Description: "Path of the failure domain",
				Optional:    true,
				Computed:    true,
			},
			"form_factor": {
				Type:         schema.TypeString,
				Default:      model.PolicyEdgeTransportNode_FORM_FACTOR_MEDIUM,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(policyEdgeNodeFormFactorValues, false),
			},
			"hostname": {
				Type:        schema.TypeString,
				Description: "Host name or FQDN for edge node",
				Required:    true,
			},
			"management_interface": {
				Type:        schema.TypeList,
				Description: "Applicable For LCM managed Node and contains the management interface info",
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_assignment": getPolicyIPAssignmentSchema(true, 1, 2, managementAssignments),
						"network_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Portgroup, logical switch identifier or segment path for management network connectivity",
						},
					},
				},
			},
			"switch": {
				Type:        schema.TypeList,
				Description: "Edge Transport Node switches configuration",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"overlay_transport_zone_path": getPolicyPathSchema(false, false, "An overlay TransportZone path that is associated with the specified edge TN switch"),
						"pnic": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Physical NIC specification",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"datapath_network_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A portgroup, logical switch identifier or segment path for datapath connectivity",
									},
									"device_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Device name or key e.g. fp-eth0, fp-eth1 etc",
									},
									"uplink_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Uplink name for this Pnic. This name will be used to reference this Pnic in other configurations",
									},
								},
							},
						},
						"uplink_host_switch_profile_path": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "Uplink Host Switch Profile Path",
							ValidateFunc: validatePolicyPath(),
						},
						"lldp_host_switch_profile_path": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "LLDP Host Switch Profile Path",
							ValidateFunc: validatePolicyPath(),
						},
						"switch_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Edge Tn switch name. This name will be used to reference an edge TN switch",
							Default:     "nsxDefaultHostSwitch",
						},
						"tunnel_endpoint": {
							Type:        schema.TypeList,
							Description: "Tunnel Endpoint",
							MaxItems:    1,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_assignment": getPolicyIPAssignmentSchema(true, 1, 1, tepAssignments),
									"vlan": {
										Type:         schema.TypeInt,
										Optional:     true,
										Description:  "VLAN ID for tunnel endpoint",
										ValidateFunc: validation.IntBetween(0, 4094),
									},
								},
							},
						},
						"vlan_transport_zone_paths": {
							Type:        schema.TypeList,
							Description: "List of Vlan TransportZone paths that are to be associated with specified edge TN switch",
							Optional:    true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validatePolicyPath(),
							},
						},
					},
				},
			},
			"vm_deployment_config": {
				Type:        schema.TypeList,
				Description: "VM deployment configuration for LCM nodes",
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compute_folder_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Compute folder identifier in the specified vcenter server",
						},
						"compute_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cluster identifier for specified vcenter server",
						},
						"edge_host_affinity_config": {
							Type:        schema.TypeList,
							Description: "Edge VM to host affinity configuration",
							MaxItems:    1,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_group_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Host group name",
									},
								},
							},
						},
						"host_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Host identifier in the specified vcenter server",
						},
						"reservation_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Resource reservation settings",
							Optional:    true,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu_reservation_in_mhz": {
										Type:        schema.TypeInt,
										Description: "CPU reservation in MHz",
										Optional:    true,
									},
									"cpu_reservation_in_shares": {
										Type:         schema.TypeString,
										Description:  "CPU reservation in shares",
										Optional:     true,
										Default:      model.CPUReservation_RESERVATION_IN_SHARES_HIGH_PRIORITY,
										ValidateFunc: validation.StringInSlice(policyCpuReservationValues, false),
									},
									"memory_reservation_percentage": {
										Type:         schema.TypeInt,
										Optional:     true,
										Description:  "Memory reservation percentage",
										ValidateFunc: validation.IntBetween(0, 100),
										Default:      100,
									},
								},
							},
						},
						"storage_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Storage/datastore identifier in the specified vcenter server",
						},
						"compute_manager_id": {
							Type:        schema.TypeString,
							Description: "Vsphere compute identifier for identifying the vcenter server",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceNsxtPolicyEdgeTransportNodeExists(siteID, epID, tzID string, connector client.Connector) (bool, error) {
	var err error

	// Check site existence first
	siteClient := infra.NewSitesClient(connector)
	_, err = siteClient.Get(siteID)
	if err != nil {
		msg := fmt.Sprintf("failed to read site %s", siteID)
		return false, logAPIError(msg, err)
	}

	// Check (ep, htn) existence. In case of ep not found, NSX returns BAD_REQUEST
	htnClient := enforcement_points.NewEdgeTransportNodesClient(connector)
	_, err = htnClient.Get(siteID, epID, tzID)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func getApplianceConfigFromSchema(appCfg interface{}) *model.PolicyVmApplianceConfig {
	if appCfg == nil {
		return nil
	}
	var obj model.PolicyVmApplianceConfig

	for _, aCfg := range appCfg.([]interface{}) {
		appCfg := aCfg.(map[string]interface{})
		allowSshRootLogin := appCfg["allow_ssh_root_login"].(bool)
		dnsServers := interfaceListToStringList(appCfg["dns_servers"].([]interface{}))
		enableSSH := appCfg["enable_ssh"].(bool)
		enableUptMode := appCfg["enable_upt_mode"].(bool)
		var syslogServers []model.SyslogConfiguration
		if appCfg["syslog_servers"] != nil {
			for _, syslogSrvr := range appCfg["syslog_servers"].([]interface{}) {
				syslogServer := syslogSrvr.(map[string]interface{})
				logLevel := syslogServer["log_level"].(string)
				port := syslogServer["port"].(string)
				protocol := syslogServer["protocol"].(string)
				server := syslogServer["server"].(string)
				syslogServers = append(syslogServers, model.SyslogConfiguration{
					LogLevel: &logLevel,
					Port:     &port,
					Protocol: &protocol,
					Server:   &server,
				})
			}
		}

		obj = model.PolicyVmApplianceConfig{
			AllowSshRootLogin: &allowSshRootLogin,
			DnsServers:        dnsServers,
			EnableSsh:         &enableSSH,
			EnableUptMode:     &enableUptMode,
			SyslogServers:     syslogServers,
		}
	}
	return &obj
}

func getCredentialsFromSchema(creds interface{}) *model.PolicyEdgeTransportNodeCredential {
	if creds == nil {
		return nil
	}
	var obj model.PolicyEdgeTransportNodeCredential

	for _, iCred := range creds.([]interface{}) {
		cred := iCred.(map[string]interface{})
		auditPassword := cred["audit_password"].(string)

		auditUsername := cred["audit_username"].(string)
		cliPassword := cred["cli_password"].(string)
		cliUsername := cred["cli_username"].(string)
		rootPassword := cred["root_password"].(string)

		obj = model.PolicyEdgeTransportNodeCredential{
			CliPassword:  &cliPassword,
			CliUsername:  &cliUsername,
			RootPassword: &rootPassword,
		}
		if auditUsername != "" {
			obj.AuditUsername = &auditUsername
		}

		if auditPassword != "" {
			obj.AuditPassword = &auditPassword
		}
	}
	return &obj
}

func getPolicyIPAssignmentsFromSchema(iIPAssments interface{}) ([]*data.StructValue, error) {
	if iIPAssments == nil {
		return nil, nil
	}

	var iPAssignments []*data.StructValue
	converter := bindings.NewTypeConverter()

	for _, iIPAssment := range iIPAssments.([]interface{}) {
		iPAssignment := iIPAssment.(map[string]interface{})

		// Handle auto_conf assignments
		iAutoConf := iPAssignment["auto_conf"]
		if iAutoConf != nil && iAutoConf.(bool) {
			obj := model.AutoConf{
				IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_AUTOCONF,
			}
			dataValue, errs := converter.ConvertToVapi(obj, model.AutoConfBindingType())
			if errs != nil {
				return nil, errs[0]
			}
			iPAssignments = append(iPAssignments, dataValue.(*data.StructValue))
		}

		// Handle dhcp_v4 assignments
		iDhcpV4 := iPAssignment["dhcp_v4"]
		if iDhcpV4 != nil && iDhcpV4.(bool) {
			obj := model.Dhcpv4{
				IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_DHCPV4,
			}
			dataValue, errs := converter.ConvertToVapi(obj, model.Dhcpv4BindingType())
			if errs != nil {
				return nil, errs[0]
			}
			iPAssignments = append(iPAssignments, dataValue.(*data.StructValue))
		}

		// Handle dhcp_v6 assignments
		iDhcpV6 := iPAssignment["dhcp_v6"]
		if iDhcpV6 != nil && iDhcpV6.(bool) {
			obj := model.Dhcpv6{
				IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_DHCPV6,
			}
			dataValue, errs := converter.ConvertToVapi(obj, model.Dhcpv6BindingType())
			if errs != nil {
				return nil, errs[0]
			}
			iPAssignments = append(iPAssignments, dataValue.(*data.StructValue))
		}

		// Handle no_assignment assignments
		iNoAssignment := iPAssignment["no_assignment"]
		if iNoAssignment != nil && iNoAssignment.(bool) {
			obj := model.NoAssignment{
				IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_NOASSIGNMENT,
			}
			dataValue, errs := converter.ConvertToVapi(obj, model.NoAssignmentBindingType())
			if errs != nil {
				return nil, errs[0]
			}
			iPAssignments = append(iPAssignments, dataValue.(*data.StructValue))
		}

		// Handle static_ipv4 assignments
		iStaticIPv4 := iPAssignment["static_ipv4"]
		if iStaticIPv4 != nil {
			for _, e := range iStaticIPv4.([]interface{}) {
				elem := e.(map[string]interface{})
				defaultGateway := elem["default_gateway"].(string)
				var managementPortSubnet []model.IPv4Subnet
				impsList := elem["management_port_subnet"].([]interface{})
				for _, imps := range impsList {
					mps := imps.(map[string]interface{})
					ipAddresses := interfaceListToStringList(mps["ip_addresses"].([]interface{}))
					prefixLength := int64(mps["prefix_length"].(int))
					managementPortSubnet = append(managementPortSubnet, model.IPv4Subnet{
						IpAddresses:  ipAddresses,
						PrefixLength: &prefixLength,
					})
				}

				obj := model.StaticIpv4{
					DefaultGateway:        []string{defaultGateway},
					ManagementPortSubnets: managementPortSubnet,
					IpAssignmentType:      model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV4,
				}
				dataValue, errs := converter.ConvertToVapi(obj, model.StaticIpv4BindingType())
				if errs != nil {
					return nil, errs[0]
				}
				iPAssignments = append(iPAssignments, dataValue.(*data.StructValue))
			}
		}

		// Handle static_ipv4_list assignments
		staticIpv4List := iPAssignment["static_ipv4_list"]
		if staticIpv4List != nil {
			for _, e := range staticIpv4List.([]interface{}) {
				elem := e.(map[string]interface{})

				defaultGateway := elem["default_gateway"].(string)
				ipAddresses := interfaceListToStringList(elem["ip_addresses"].([]interface{}))
				subnetMask := elem["subnet_mask"].(string)

				obj := model.StaticIpv4List{
					DefaultGateway:   &defaultGateway,
					IpList:           ipAddresses,
					SubnetMask:       &subnetMask,
					IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV4LIST,
				}
				dataValue, errs := converter.ConvertToVapi(obj, model.StaticIpv4ListBindingType())
				if errs != nil {
					return nil, errs[0]
				}
				iPAssignments = append(iPAssignments, dataValue.(*data.StructValue))
			}
		}

		// Handle static_ipv4_mac_list assignments
		staticIpv4MacList := iPAssignment["static_ipv4_mac_list"]
		if staticIpv4MacList != nil {
			for _, e := range staticIpv4MacList.([]interface{}) {
				elem := e.(map[string]interface{})

				defaultGateway := elem["default_gateway"].(string)
				var ipMacList []model.IpMacPair
				iIpMacList := elem["ip_mac_pair"].([]interface{})
				for _, iIpMac := range iIpMacList {
					ipMac := iIpMac.(map[string]interface{})
					ipAddress := ipMac["ip_address"].(string)
					macAddress := ipMac["mac_address"].(string)
					ipMacList = append(ipMacList, model.IpMacPair{
						Ip:  &ipAddress,
						Mac: &macAddress,
					})
				}

				subnetMask := elem["subnet_mask"].(string)

				obj := model.StaticIpv4MacList{
					DefaultGateway:   &defaultGateway,
					IpMacList:        ipMacList,
					SubnetMask:       &subnetMask,
					IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV4MACLIST,
				}
				dataValue, errs := converter.ConvertToVapi(obj, model.StaticIpv4MacListBindingType())
				if errs != nil {
					return nil, errs[0]
				}
				iPAssignments = append(iPAssignments, dataValue.(*data.StructValue))
			}
		}

		// Handle static_ipv4_pool assignments
		staticIpv4Pool := iPAssignment["static_ipv4_pool"]
		if staticIpv4Pool != nil && staticIpv4Pool != "" {
			ipPool := staticIpv4Pool.(string)
			obj := model.StaticIpv4Pool{
				IpPool:           &ipPool,
				IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV4POOL,
			}
			dataValue, errs := converter.ConvertToVapi(obj, model.StaticIpv4PoolBindingType())
			if errs != nil {
				return nil, errs[0]
			}
			iPAssignments = append(iPAssignments, dataValue.(*data.StructValue))
		}

		// Handle static_ipv6 assignments
		iStaticIPv6 := iPAssignment["static_ipv6"]
		if iStaticIPv6 != nil {
			for _, e := range iStaticIPv6.([]interface{}) {
				elem := e.(map[string]interface{})
				defaultGateway := elem["default_gateway"].(string)
				var managementPortSubnet []model.IPv6Subnet
				impsList := elem["management_port_subnet"].([]interface{})
				for _, imps := range impsList {
					mps := imps.(map[string]interface{})
					ipAddresses := interfaceListToStringList(mps["ip_addresses"].([]interface{}))
					prefixLength := int64(mps["prefix_length"].(int))
					managementPortSubnet = append(managementPortSubnet, model.IPv6Subnet{
						IpAddresses:  ipAddresses,
						PrefixLength: &prefixLength,
					})
				}

				obj := model.StaticIpv6{
					DefaultGateway:        []string{defaultGateway},
					ManagementPortSubnets: managementPortSubnet,
					IpAssignmentType:      model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV6,
				}
				dataValue, errs := converter.ConvertToVapi(obj, model.StaticIpv6BindingType())
				if errs != nil {
					return nil, errs[0]
				}
				iPAssignments = append(iPAssignments, dataValue.(*data.StructValue))
			}
		}

		// Handle static_ipv6_list assignments
		staticIpv6List := iPAssignment["static_ipv6_list"]
		if staticIpv6List != nil {
			for _, e := range staticIpv6List.([]interface{}) {
				elem := e.(map[string]interface{})

				defaultGateway := elem["default_gateway"].(string)
				ipAddresses := interfaceListToStringList(elem["ip_addresses"].([]interface{}))
				prefixLength := strconv.Itoa(elem["prefix_length"].(int))

				obj := model.StaticIpv6List{
					DefaultGateway:   &defaultGateway,
					IpList:           ipAddresses,
					PrefixLength:     &prefixLength,
					IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV6LIST,
				}
				dataValue, errs := converter.ConvertToVapi(obj, model.StaticIpv6ListBindingType())
				if errs != nil {
					return nil, errs[0]
				}
				iPAssignments = append(iPAssignments, dataValue.(*data.StructValue))
			}
		}

		// Handle static_ipv6_mac_list assignments
		staticIpv6MacList := iPAssignment["static_ipv6_mac_list"]
		if staticIpv6MacList != nil {
			for _, e := range staticIpv6MacList.([]interface{}) {
				elem := e.(map[string]interface{})

				defaultGateway := elem["default_gateway"].(string)
				var ipMacList []model.Ipv6MacPair
				iIpMacList := elem["ip_mac_pair"].([]interface{})
				for _, iIpMac := range iIpMacList {
					ipMac := iIpMac.(map[string]interface{})
					ipAddress := ipMac["ip_address"].(string)
					macAddress := ipMac["mac_address"].(string)
					ipMacList = append(ipMacList, model.Ipv6MacPair{
						Ipv6: &ipAddress,
						Mac:  &macAddress,
					})
				}

				prefixLength := strconv.Itoa(elem["prefix_length"].(int))

				obj := model.StaticIpv6MacList{
					DefaultGateway:   &defaultGateway,
					IpMacList:        ipMacList,
					PrefixLength:     &prefixLength,
					IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV6MACLIST,
				}
				dataValue, errs := converter.ConvertToVapi(obj, model.StaticIpv6MacListBindingType())
				if errs != nil {
					return nil, errs[0]
				}
				iPAssignments = append(iPAssignments, dataValue.(*data.StructValue))
			}
		}

		// Handle static_ipv6_pool assignments
		staticIpv6Pool := iPAssignment["static_ipv6_pool"]
		if staticIpv6Pool != nil && staticIpv6Pool != "" {
			ipPool := staticIpv6Pool.(string)

			obj := model.StaticIpv6Pool{
				IpPool:           &ipPool,
				IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV6POOL,
			}
			dataValue, errs := converter.ConvertToVapi(obj, model.StaticIpv6PoolBindingType())
			if errs != nil {
				return nil, errs[0]
			}
			iPAssignments = append(iPAssignments, dataValue.(*data.StructValue))
		}
	}
	return iPAssignments, nil
}

func getManagementInterfaceFromSchema(mgtIntface interface{}) (*model.PolicyEdgeTransportManagementInterface, error) {
	if mgtIntface == nil {
		return nil, nil
	}
	var obj model.PolicyEdgeTransportManagementInterface
	for _, iMgtInt := range mgtIntface.([]interface{}) {
		mgtInt := iMgtInt.(map[string]interface{})
		iPAssignments, err := getPolicyIPAssignmentsFromSchema(mgtInt["ip_assignment"])
		if err != nil {
			return nil, err
		}
		networkID := mgtInt["network_id"].(string)

		obj = model.PolicyEdgeTransportManagementInterface{
			IpAssignmentSpecs: iPAssignments,
			NetworkId:         &networkID,
		}
	}
	return &obj, nil
}

func getSwitchFromSchema(iSwitchSpec interface{}) (*model.PolicyEdgeTransportNodeSwitchSpec, error) {
	if iSwitchSpec == nil {
		return nil, nil
	}

	var switchSpec model.PolicyEdgeTransportNodeSwitchSpec
	for _, iSwitch := range iSwitchSpec.([]interface{}) {
		sw := iSwitch.(map[string]interface{})
		overlayTransportZonePath := sw["overlay_transport_zone_path"].(string)

		var pnics []model.PolicyEdgeTransportNodePnic
		if sw["pnic"] != nil {
			for _, iPnic := range sw["pnic"].([]interface{}) {
				pnic := iPnic.(map[string]interface{})
				datapathNetworkId := pnic["datapath_network_id"].(string)
				deviceName := pnic["device_name"].(string)
				uplinkName := pnic["uplink_name"].(string)

				pnics = append(pnics, model.PolicyEdgeTransportNodePnic{
					DatapathNetworkId: &datapathNetworkId,
					DeviceName:        &deviceName,
					UplinkName:        &uplinkName,
				})
			}
		}
		var profilePaths []model.PolicyEdgeTransportNodeSwitchProfileTypePathEntry
		uplinkHostSwitchProfilePath := sw["uplink_host_switch_profile_path"].(string)
		if uplinkHostSwitchProfilePath != "" {
			profilePaths = append(profilePaths, model.PolicyEdgeTransportNodeSwitchProfileTypePathEntry{
				Key:   strPtr(model.PolicyEdgeTransportNodeSwitchProfileTypePathEntry_KEY_UPLINKHOSTSWITCHPROFILE),
				Value: &uplinkHostSwitchProfilePath,
			})
		}
		lldpHostSwitchProfilePath := sw["lldp_host_switch_profile_path"].(string)
		if lldpHostSwitchProfilePath != "" {
			profilePaths = append(profilePaths, model.PolicyEdgeTransportNodeSwitchProfileTypePathEntry{
				Key:   strPtr(model.PolicyEdgeTransportNodeSwitchProfileTypePathEntry_KEY_LLDPHOSTSWITCHPROFILE),
				Value: &lldpHostSwitchProfilePath,
			})
		}

		switchName := sw["switch_name"].(string)
		var tunnelEndpoints []model.PolicyEdgeTransportNodeSwitchTunnelEndPoint
		if sw["tunnel_endpoint"] != nil {
			for _, iTunnelEndpoint := range sw["tunnel_endpoint"].([]interface{}) {
				tunnelEndpoint := iTunnelEndpoint.(map[string]interface{})
				ipAssignment, err := getPolicyIPAssignmentsFromSchema(tunnelEndpoint["ip_assignment"])
				if err != nil {
					return nil, err
				}
				vlan := int64(tunnelEndpoint["vlan"].(int))
				tunnelEndpoints = append(tunnelEndpoints, model.PolicyEdgeTransportNodeSwitchTunnelEndPoint{
					IpAssignmentSpecs: ipAssignment,
					Vlan:              &vlan,
				})
			}
		}
		vlanTransportZonePaths := interfaceListToStringList(sw["vlan_transport_zone_paths"].([]interface{}))

		swit := model.PolicyEdgeTransportNodeSwitch{
			OverlayTransportZonePaths: []string{overlayTransportZonePath},
			Pnics:                     pnics,
			ProfilePaths:              profilePaths,
			SwitchName:                &switchName,
			TunnelEndpoints:           tunnelEndpoints,
			VlanTransportZonePaths:    vlanTransportZonePaths,
		}
		switchSpec.Switches = append(switchSpec.Switches, swit)
	}

	return &switchSpec, nil
}

func getVMDeploymentConfigFromSchema(iVmDeploymentCfg interface{}) (*data.StructValue, error) {
	if iVmDeploymentCfg == nil {
		return nil, nil
	}
	for _, iDeploymentCfg := range iVmDeploymentCfg.([]interface{}) {
		deploymentCfg := iDeploymentCfg.(map[string]interface{})
		computeFolderId := deploymentCfg["compute_folder_id"].(string)
		computeId := deploymentCfg["compute_id"].(string)
		var edgeHostAffinityConfig *model.EdgeHostAffinityConfig
		if deploymentCfg["edge_host_affinity_config"] != nil {
			eaCfg := deploymentCfg["edge_host_affinity_config"].([]interface{})
			for _, cfg := range eaCfg {
				c := cfg.(map[string]interface{})
				hostGroupName := c["host_group_name"].(string)
				edgeHostAffinityConfig = &model.EdgeHostAffinityConfig{
					HostGroupName: &hostGroupName,
				}
			}
		}
		hostID := deploymentCfg["host_id"].(string)

		var reservationInfo *model.ReservationInfo
		if deploymentCfg["reservation_info"] != nil {
			for _, ri := range deploymentCfg["reservation_info"].([]interface{}) {
				rInfo := ri.(map[string]interface{})
				cpuReservationInMhz := int64(rInfo["cpu_reservation_in_mhz"].(int))
				cpuReservationInShares := rInfo["cpu_reservation_in_shares"].(string)
				memoryReservationPercentage := int64(rInfo["memory_reservation_percentage"].(int))

				reservationInfo = &model.ReservationInfo{
					CpuReservation: &model.CPUReservation{
						ReservationInMhz:    &cpuReservationInMhz,
						ReservationInShares: &cpuReservationInShares,
					},
					MemoryReservation: &model.MemoryReservation{
						ReservationPercentage: &memoryReservationPercentage,
					},
				}
			}
		}
		storageID := deploymentCfg["storage_id"].(string)
		vcID := deploymentCfg["compute_manager_id"].(string)

		vmDeploymentCfg := model.PolicyVsphereDeploymentConfig{
			ComputeId:              &computeId,
			EdgeHostAffinityConfig: edgeHostAffinityConfig,
			HostId:                 &hostID,
			ReservationInfo:        reservationInfo,
			StorageId:              &storageID,
			VcId:                   &vcID,
			PlacementType:          model.PolicyVmDeploymentConfig_PLACEMENT_TYPE_POLICYVSPHEREDEPLOYMENTCONFIG,
		}
		if computeFolderId != "" {
			vmDeploymentCfg.ComputeFolderId = &computeFolderId
		}
		converter := bindings.NewTypeConverter()
		dataValue, errs := converter.ConvertToVapi(vmDeploymentCfg, model.PolicyVsphereDeploymentConfigBindingType())
		if errs != nil {
			return nil, errs[0]
		}
		return dataValue.(*data.StructValue), nil
	}
	return nil, nil
}

func policyEdgeTransportNodePredeployedPatch(siteID, epID, etnID string, d *schema.ResourceData, m interface{}) error {

	connector := getPolicyConnector(m)
	etnClient := enforcement_points.NewEdgeTransportNodesClient(connector)

	obj, err := etnClient.Get(siteID, epID, etnID)
	if err != nil {
		return err
	}

	description := d.Get("description").(string)
	obj.Description = &description

	displayName := d.Get("display_name").(string)
	obj.DisplayName = &displayName

	obj.Tags = getPolicyTagsFromSchema(d)

	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision

	failureDomainPath := d.Get("failure_domain_path").(string)
	obj.FailureDomainPath = &failureDomainPath

	hostname := d.Get("hostname").(string)
	obj.Hostname = &hostname

	switchSpec, err := getSwitchFromSchema(d.Get("switch"))
	if err != nil {
		return err
	}
	obj.SwitchSpec = switchSpec

	return etnClient.Patch(siteID, epID, etnID, obj)
}

func policyEdgeTransportNodePatch(siteID, epID, etnID string, d *schema.ResourceData, m interface{}) error {

	description := d.Get("description").(string)
	displayName := d.Get("display_name").(string)
	tags := getPolicyTagsFromSchema(d)
	revision := int64(d.Get("revision").(int))
	advancedConfiguration := getPolicyKeyValuePairListFromSchema(d.Get("advanced_configuration"))
	applianceConfig := getApplianceConfigFromSchema(d.Get("appliance_config"))
	credentials := getCredentialsFromSchema(d.Get("credentials"))
	failureDomainPath := d.Get("failure_domain_path").(string)
	formFactor := d.Get("form_factor").(string)
	hostname := d.Get("hostname").(string)
	managementInterface, err := getManagementInterfaceFromSchema(d.Get("management_interface"))
	if err != nil {
		return err
	}

	switchSpec, err := getSwitchFromSchema(d.Get("switch"))
	if err != nil {
		return err
	}

	vmDeploymentConfig, err := getVMDeploymentConfigFromSchema(d.Get("vm_deployment_config"))
	if err != nil {
		return err
	}

	obj := model.PolicyEdgeTransportNode{
		Description:           &description,
		DisplayName:           &displayName,
		Tags:                  tags,
		Revision:              &revision,
		AdvancedConfiguration: advancedConfiguration,
		ApplianceConfig:       applianceConfig,
		Credentials:           credentials,
		FailureDomainPath:     &failureDomainPath,
		FormFactor:            &formFactor,
		Hostname:              &hostname,
		ManagementInterface:   managementInterface,
		SwitchSpec:            switchSpec,
		VmDeploymentConfig:    vmDeploymentConfig,
	}

	connector := getPolicyConnector(m)
	etnClient := enforcement_points.NewEdgeTransportNodesClient(connector)
	return etnClient.Patch(siteID, epID, etnID, obj)
}

func resourceNsxtPolicyEdgeTransportNodeCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Get("nsx_id").(string)
	nodeID := d.Get("node_id").(string)
	if id == "" {
		if nodeID != "" {
			id = nodeID
		} else {
			id = newUUID()
		}
	}
	sitePath := d.Get("site_path").(string)
	siteID := getResourceIDFromResourcePath(sitePath, "sites")
	if siteID == "" {
		return fmt.Errorf("error obtaining Site ID from site path %s", sitePath)
	}
	epID := d.Get("enforcement_point").(string)
	if epID == "" {
		epID = getPolicyEnforcementPoint(m)
	}

	if nodeID == "" {
		exists, err := resourceNsxtPolicyEdgeTransportNodeExists(siteID, epID, id, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("resource with ID %s already exists", id)
		}

		// Create the resource using PATCH
		log.Printf("[INFO] Creating PolicyEdgeTransportNode with ID %s under site %s enforcement point %s", id, siteID, epID)
		err = policyEdgeTransportNodePatch(siteID, epID, id, d, m)
		if err != nil {
			return handleCreateError("EdgeTransportNode", id, err)
		}
	} else {
		log.Printf("Adding a pre-existing Edge appliance")
		if id != nodeID {
			return fmt.Errorf("cannot have both nsx_id and node_id attribute set, with different values")
		}
		err := policyEdgeTransportNodePredeployedPatch(siteID, epID, nodeID, d, m)
		if err != nil {
			return err
		}
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyEdgeTransportNodeRead(d, m)
}

func setApplianceConfigInSchema(d *schema.ResourceData, obj *model.PolicyVmApplianceConfig) error {
	if obj == nil {
		return nil
	}
	appCfg := make(map[string]interface{})
	appCfg["allow_ssh_root_login"] = obj.AllowSshRootLogin
	appCfg["dns_servers"] = obj.DnsServers
	appCfg["enable_ssh"] = obj.EnableSsh
	appCfg["enable_upt_mode"] = obj.EnableUptMode

	var syslogServers []interface{}
	if obj.SyslogServers != nil {
		for _, syslogServer := range obj.SyslogServers {
			server := make(map[string]interface{})
			server["log_level"] = syslogServer.LogLevel
			server["port"] = syslogServer.Port
			server["protocol"] = syslogServer.Protocol
			server["server"] = syslogServer.Server
			syslogServers = append(syslogServers, server)
		}
		appCfg["syslog_server"] = syslogServers
	}

	d.Set("appliance_config", []interface{}{appCfg})
	return nil
}

func setCredentialsInSchema(d *schema.ResourceData, obj *model.PolicyEdgeTransportNodeCredential) error {
	if obj == nil {
		return nil
	}

	var creds map[string]interface{}
	// Use credentials from intent, as passwords aren't returned from NSX
	c := d.Get("credentials").([]interface{})
	if len(c) > 0 {
		creds = c[0].(map[string]interface{})
	} else {
		creds = make(map[string]interface{})
	}
	if obj.AuditUsername != nil {
		creds["audit_username"] = obj.AuditUsername
	}
	if obj.CliUsername != nil {
		creds["cli_username"] = obj.CliUsername
	}

	d.Set("credentials", []interface{}{creds})
	return nil
}

func setPolicyIPAssignmentsInSchema(specs []*data.StructValue) (interface{}, error) {
	converter := bindings.NewTypeConverter()
	elem := make(map[string]interface{})

	for _, spec := range specs {
		base, errs := converter.ConvertToGolang(spec, model.PolicyIpAssignmentSpecBindingType())
		if errs != nil {
			return nil, errs[0]
		}
		assignmentType := base.(model.PolicyIpAssignmentSpec).IpAssignmentType
		switch assignmentType {
		case model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_AUTOCONF:
			elem["auto_conf"] = true

		case model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_DHCPV4:
			elem["dhcp_v4"] = true

		case model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_DHCPV6:
			elem["dhcp_v6"] = true

		case model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_NOASSIGNMENT:
			elem["no_assignment"] = true

		case model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV4:
			v, errs := converter.ConvertToGolang(spec, model.StaticIpv4BindingType())
			if errs != nil {
				return nil, errs[0]
			}
			assignment := v.(model.StaticIpv4)
			value := make(map[string]interface{})
			value["default_gateway"] = assignment.DefaultGateway

			var mpSubnets []interface{}
			for _, mSubnet := range assignment.ManagementPortSubnets {
				mpsubnet := make(map[string]interface{})
				mpsubnet["ip_addresses"] = mSubnet.IpAddresses
				mpsubnet["prefix_length"] = mSubnet.PrefixLength
				mpSubnets = append(mpSubnets, mpsubnet)
			}
			value["management_port_subnet"] = mpSubnets
			elem["static_ipv4"] = []interface{}{value}

		case model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV4LIST:
			v, errs := converter.ConvertToGolang(spec, model.StaticIpv4ListBindingType())
			if errs != nil {
				return nil, errs[0]
			}
			assignment := v.(model.StaticIpv4List)
			value := make(map[string]interface{})
			value["default_gateway"] = assignment.DefaultGateway
			value["ip_addresses"] = assignment.IpList
			value["subnet_mask"] = assignment.SubnetMask
			elem["static_ipv4_list"] = []interface{}{value}

		case model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV4MACLIST:
			v, errs := converter.ConvertToGolang(spec, model.StaticIpv4MacListBindingType())
			if errs != nil {
				return nil, errs[0]
			}
			assignment := v.(model.StaticIpv4MacList)
			value := make(map[string]interface{})
			value["default_gateway"] = assignment.DefaultGateway

			var ipMacList []interface{}
			for _, ipMac := range assignment.IpMacList {
				im := make(map[string]interface{})
				im["ip_address"] = ipMac.Ip
				im["mac_address"] = ipMac.Mac
				ipMacList = append(ipMacList, im)
			}
			value["ip_mac_pair"] = ipMacList

			value["subnet_mask"] = assignment.SubnetMask
			elem["static_ipv4_mac_list"] = []interface{}{value}

		case model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV4POOL:
			v, errs := converter.ConvertToGolang(spec, model.StaticIpv4PoolBindingType())
			if errs != nil {
				return nil, errs[0]
			}
			elem["static_ipv4_pool"] = v.(model.StaticIpv4Pool).IpPool

		case model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV6:
			v, errs := converter.ConvertToGolang(spec, model.StaticIpv6BindingType())
			if errs != nil {
				return nil, errs[0]
			}
			assignment := v.(model.StaticIpv6)
			value := make(map[string]interface{})
			value["default_gateway"] = assignment.DefaultGateway

			var mpSubnets []interface{}
			for _, mSubnet := range assignment.ManagementPortSubnets {
				mpsubnet := make(map[string]interface{})
				mpsubnet["ip_addresses"] = mSubnet.IpAddresses
				mpsubnet["prefix_length"] = mSubnet.PrefixLength
				mpSubnets = append(mpSubnets, mpsubnet)
			}
			value["management_port_subnet"] = mpSubnets
			elem["static_ipv6"] = []interface{}{value}

		case model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV6LIST:
			v, errs := converter.ConvertToGolang(spec, model.StaticIpv6ListBindingType())
			if errs != nil {
				return nil, errs[0]
			}
			assignment := v.(model.StaticIpv6List)
			value := make(map[string]interface{})
			value["default_gateway"] = assignment.DefaultGateway
			value["ip_addresses"] = assignment.IpList
			pfx, err := strconv.Atoi(*assignment.PrefixLength)
			if err != nil {
				return nil, err
			}
			value["prefix_length"] = pfx
			elem["static_ipv6_list"] = []interface{}{value}

		case model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV6MACLIST:
			v, errs := converter.ConvertToGolang(spec, model.StaticIpv6MacListBindingType())
			if errs != nil {
				return nil, errs[0]
			}
			assignment := v.(model.StaticIpv6MacList)
			value := make(map[string]interface{})
			value["default_gateway"] = assignment.DefaultGateway

			var ipMacList []interface{}
			for _, ipMac := range assignment.IpMacList {
				im := make(map[string]interface{})
				im["ip_address"] = ipMac.Ipv6
				im["mac_address"] = ipMac.Mac
				ipMacList = append(ipMacList, im)
			}
			value["ip_mac_pair"] = ipMacList

			pfx, err := strconv.Atoi(*assignment.PrefixLength)
			if err != nil {
				return nil, err
			}
			value["subnet_mask"] = pfx
			elem["static_ipv6_mac_list"] = []interface{}{value}

		case model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV6POOL:
			v, errs := converter.ConvertToGolang(spec, model.StaticIpv6PoolBindingType())
			if errs != nil {
				return nil, errs[0]
			}
			elem["static_ipv6_pool"] = v.(model.StaticIpv6Pool).IpPool
		}
	}

	return elem, nil
}

func resourceNsxtPolicyEdgeTransportNodeRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := enforcement_points.NewEdgeTransportNodesClient(connector)

	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	obj, err := client.Get(siteID, epID, id)
	if err != nil {
		return handleReadError(d, "EdgeTransportNode", id, err)
	}
	sitePath, err := getSitePathFromChildResourcePath(*obj.ParentPath)
	if err != nil {
		return handleReadError(d, "EdgeTransportNode", id, err)
	}

	d.Set("site_path", sitePath)
	d.Set("enforcement_point", epID)
	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("advanced_configuration", setPolicyKeyValueListForSchema(obj.AdvancedConfiguration))

	err = setApplianceConfigInSchema(d, obj.ApplianceConfig)
	if err != nil {
		return handleReadError(d, "EdgeTransportNode", id, err)
	}

	err = setCredentialsInSchema(d, obj.Credentials)
	if err != nil {
		return handleReadError(d, "EdgeTransportNode", id, err)
	}

	d.Set("failure_domain_path", obj.FailureDomainPath)
	d.Set("form_factor", obj.FormFactor)
	d.Set("hostname", obj.Hostname)

	if obj.ManagementInterface != nil {
		mgtInterface := make(map[string]interface{})
		mgtInterface["ip_assignment"], err = setPolicyIPAssignmentsInSchema(obj.ManagementInterface.IpAssignmentSpecs)
		if err != nil {
			return err
		}
		mgtInterface["network_id"] = obj.ManagementInterface.NetworkId
		d.Set("management_interface", []interface{}{mgtInterface})
	}

	var switches []interface{}
	if obj.SwitchSpec != nil && obj.SwitchSpec.Switches != nil {
		for _, switc := range obj.SwitchSpec.Switches {
			sw := make(map[string]interface{})
			if len(switc.OverlayTransportZonePaths) > 0 {
				sw["overlay_transport_zone_path"] = switc.OverlayTransportZonePaths[0]
			}

			var pnics []interface{}
			for _, pnic := range switc.Pnics {
				p := make(map[string]interface{})
				if pnic.DatapathNetworkId != nil {
					p["datapath_network_id"] = pnic.DatapathNetworkId
				}
				p["device_name"] = pnic.DeviceName
				p["uplink_name"] = pnic.UplinkName

				pnics = append(pnics, p)
			}
			sw["pnic"] = pnics

			for _, pp := range switc.ProfilePaths {
				if *pp.Key == model.PolicyEdgeTransportNodeSwitchProfileTypePathEntry_KEY_UPLINKHOSTSWITCHPROFILE {
					sw["uplink_host_switch_profile_path"] = pp.Value
				} else if *pp.Key == model.PolicyEdgeTransportNodeSwitchProfileTypePathEntry_KEY_LLDPHOSTSWITCHPROFILE {
					sw["lldp_host_switch_profile_path"] = pp.Value
				}
			}

			sw["switch_name"] = switc.SwitchName
			var tunnelEndpoints []interface{}
			for _, tep := range switc.TunnelEndpoints {
				t := make(map[string]interface{})
				t["ip_assignment"], err = setPolicyIPAssignmentsInSchema(tep.IpAssignmentSpecs)
				if err != nil {
					return err
				}
				t["vlan"] = tep.Vlan

				tunnelEndpoints = append(tunnelEndpoints, t)
			}
			sw["tunnel_endpoint"] = tunnelEndpoints
			d.Set("vlan_transport_zone_paths", switc.VlanTransportZonePaths)

			switches = append(switches, sw)
		}
	}
	d.Set("switch", switches)

	if obj.VmDeploymentConfig != nil {
		converter := bindings.NewTypeConverter()
		base, errs := converter.ConvertToGolang(obj.VmDeploymentConfig, model.PolicyVmDeploymentConfigBindingType())
		if errs != nil {
			return errs[0]
		}

		if base.(model.PolicyVmDeploymentConfig).PlacementType == model.PolicyVmDeploymentConfig_PLACEMENT_TYPE_POLICYVSPHEREDEPLOYMENTCONFIG {
			c, errs := converter.ConvertToGolang(obj.VmDeploymentConfig, model.PolicyVsphereDeploymentConfigBindingType())
			if errs != nil {
				return errs[0]
			}

			deploymentConfig := c.(model.PolicyVsphereDeploymentConfig)

			depCfg := make(map[string]interface{})
			depCfg["compute_folder_id"] = deploymentConfig.ComputeFolderId
			depCfg["compute_id"] = deploymentConfig.ComputeId

			if deploymentConfig.EdgeHostAffinityConfig != nil {
				ehAffCfg := make(map[string]interface{})
				ehAffCfg["host_group_name"] = deploymentConfig.EdgeHostAffinityConfig.HostGroupName
				depCfg["edge_host_affinity_config"] = []interface{}{ehAffCfg}
			}

			depCfg["host_id"] = deploymentConfig.HostId

			reservationInfo := make(map[string]interface{})
			if deploymentConfig.ReservationInfo.CpuReservation.ReservationInMhz != nil {
				reservationInfo["cpu_reservation_in_mhz"] = deploymentConfig.ReservationInfo.CpuReservation.ReservationInMhz
			}
			if deploymentConfig.ReservationInfo.CpuReservation.ReservationInShares != nil {
				reservationInfo["cpu_reservation_in_shares"] = deploymentConfig.ReservationInfo.CpuReservation.ReservationInShares
			}
			if deploymentConfig.ReservationInfo.MemoryReservation.ReservationPercentage != nil {
				reservationInfo["memory_reservation_percentage"] = deploymentConfig.ReservationInfo.MemoryReservation.ReservationPercentage
			}
			depCfg["reservation_info"] = []interface{}{reservationInfo}

			depCfg["storage_id"] = deploymentConfig.StorageId
			depCfg["compute_manager_id"] = deploymentConfig.VcId
		}
	}

	return nil
}

func resourceNsxtPolicyEdgeTransportNodeUpdate(d *schema.ResourceData, m interface{}) error {
	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating PolicyEdgeTransportNode with ID %s", id)
	err = policyEdgeTransportNodePatch(siteID, epID, id, d, m)
	if err != nil {
		return handleUpdateError("PolicyEdgeTransportNode", id, err)
	}

	return resourceNsxtPolicyEdgeTransportNodeRead(d, m)
}

func resourceNsxtPolicyEdgeTransportNodeDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	client := enforcement_points.NewEdgeTransportNodesClient(connector)

	id, siteID, epID, err := policyIDSiteEPTuple(d, m)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting PolicyEdgeTransportNode with ID %s", id)
	err = client.Delete(siteID, epID, id)
	if err != nil {
		return handleDeleteError("EdgeTransportNode", id, err)
	}

	return nil
}

func resourceNsxtPolicyEdgeTransportNodeImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	rd, err := nsxtPolicyPathResourceImporterHelper(d, m)
	if err != nil {
		return rd, err
	}

	epID, err := getParameterFromPolicyPath("/enforcement-points/", "/edge-transport-nodes/", importID)
	if err != nil {
		return nil, err
	}
	d.Set("enforcement_point", epID)
	sitePath, err := getSitePathFromChildResourcePath(importID)
	if err != nil {
		return rd, err
	}
	d.Set("site_path", sitePath)

	return rd, nil
}
