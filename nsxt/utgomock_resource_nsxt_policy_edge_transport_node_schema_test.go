//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// These cover the pure schema<->SDK-struct conversion helpers used by
// resourceNsxtPolicyEdgeTransportNode. They take plain Go values (not a live
// connector), so no mocking is required.

func TestUnitNsxt_getCredentialsFromSchema(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		assert.Nil(t, getCredentialsFromSchema(nil))
	})

	t.Run("required fields are always set, optional audit fields only when non-empty", func(t *testing.T) {
		creds := []interface{}{
			map[string]interface{}{
				"cli_username":   "admin",
				"cli_password":   "pw1",
				"root_password":  "pw2",
				"audit_username": "",
				"audit_password": "",
			},
		}
		obj := getCredentialsFromSchema(creds)
		require.NotNil(t, obj)
		assert.Equal(t, "admin", *obj.CliUsername)
		assert.Equal(t, "pw1", *obj.CliPassword)
		assert.Equal(t, "pw2", *obj.RootPassword)
		assert.Nil(t, obj.AuditUsername)
		assert.Nil(t, obj.AuditPassword)
	})

	t.Run("audit fields are set when non-empty", func(t *testing.T) {
		creds := []interface{}{
			map[string]interface{}{
				"cli_username":   "admin",
				"cli_password":   "pw1",
				"root_password":  "pw2",
				"audit_username": "auditor",
				"audit_password": "pw3",
			},
		}
		obj := getCredentialsFromSchema(creds)
		require.NotNil(t, obj)
		require.NotNil(t, obj.AuditUsername)
		require.NotNil(t, obj.AuditPassword)
		assert.Equal(t, "auditor", *obj.AuditUsername)
		assert.Equal(t, "pw3", *obj.AuditPassword)
	})
}

func TestUnitNsxt_getPolicyIPAssignmentsFromSchema(t *testing.T) {
	t.Run("nil input returns nil, nil", func(t *testing.T) {
		specs, err := getPolicyIPAssignmentsFromSchema(nil)
		require.NoError(t, err)
		assert.Nil(t, specs)
	})

	cases := []struct {
		name string
		elem map[string]interface{}
	}{
		{"auto_conf", map[string]interface{}{"auto_conf": true}},
		{"dhcp_v4", map[string]interface{}{"dhcp_v4": true}},
		{"dhcp_v6", map[string]interface{}{"dhcp_v6": true}},
		{"no_assignment", map[string]interface{}{"no_assignment": true}},
		{"static_ipv4", map[string]interface{}{
			"static_ipv4": []interface{}{
				map[string]interface{}{
					"default_gateway": "10.0.0.1",
					"management_port_subnet": []interface{}{
						map[string]interface{}{
							"ip_addresses":  []interface{}{"10.0.0.2"},
							"prefix_length": 24,
						},
					},
				},
			},
		}},
		{"static_ipv4_list", map[string]interface{}{
			"static_ipv4_list": []interface{}{
				map[string]interface{}{
					"default_gateway": "10.0.0.1",
					"ip_addresses":    []interface{}{"10.0.0.2"},
					"subnet_mask":     "255.255.255.0",
				},
			},
		}},
		{"static_ipv4_mac_list", map[string]interface{}{
			"static_ipv4_mac_list": []interface{}{
				map[string]interface{}{
					"default_gateway": "10.0.0.1",
					"ip_mac_pair": []interface{}{
						map[string]interface{}{"ip_address": "10.0.0.2", "mac_address": "00:11:22:33:44:55"},
					},
					"subnet_mask": "255.255.255.0",
				},
			},
		}},
		{"static_ipv4_pool", map[string]interface{}{"static_ipv4_pool": "/infra/ip-pools/pool-1"}},
		{"static_ipv6", map[string]interface{}{
			"static_ipv6": []interface{}{
				map[string]interface{}{
					"default_gateway": "fe80::1",
					"management_port_subnet": []interface{}{
						map[string]interface{}{
							"ip_addresses":  []interface{}{"fe80::2"},
							"prefix_length": 64,
						},
					},
				},
			},
		}},
		{"static_ipv6_list", map[string]interface{}{
			"static_ipv6_list": []interface{}{
				map[string]interface{}{
					"default_gateway": "fe80::1",
					"ip_addresses":    []interface{}{"fe80::2"},
					"prefix_length":   64,
				},
			},
		}},
		{"static_ipv6_mac_list", map[string]interface{}{
			"static_ipv6_mac_list": []interface{}{
				map[string]interface{}{
					"default_gateway": "fe80::1",
					"ip_mac_pair": []interface{}{
						map[string]interface{}{"ip_address": "fe80::2", "mac_address": "00:11:22:33:44:55"},
					},
					"prefix_length": 64,
				},
			},
		}},
		{"static_ipv6_pool", map[string]interface{}{"static_ipv6_pool": "/infra/ip-pools/pool-2"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			specs, err := getPolicyIPAssignmentsFromSchema([]interface{}{c.elem})
			require.NoError(t, err)
			require.Len(t, specs, 1)
		})
	}
}

func TestUnitNsxt_getManagementInterfaceFromSchema(t *testing.T) {
	t.Run("nil input returns nil, nil", func(t *testing.T) {
		obj, err := getManagementInterfaceFromSchema(nil)
		require.NoError(t, err)
		assert.Nil(t, obj)
	})

	t.Run("valid input sets network id and ip assignment", func(t *testing.T) {
		mgtIntf := []interface{}{
			map[string]interface{}{
				"network_id":    "network-1",
				"ip_assignment": []interface{}{map[string]interface{}{"dhcp_v4": true}},
			},
		}
		obj, err := getManagementInterfaceFromSchema(mgtIntf)
		require.NoError(t, err)
		require.NotNil(t, obj)
		assert.Equal(t, "network-1", *obj.NetworkId)
		assert.Len(t, obj.IpAssignmentSpecs, 1)
	})
}

func TestUnitNsxt_setApplianceConfigInSchema(t *testing.T) {
	res := resourceNsxtPolicyEdgeTransportNode()

	t.Run("nil object is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		err := setApplianceConfigInSchema(d, nil)
		require.NoError(t, err)
	})

	t.Run("populates appliance_config including syslog servers", func(t *testing.T) {
		allowSSHRoot, enableSSH, enableUpt := true, true, false
		logLevel, port, protocol, server := "INFO", "514", "UDP", "syslog.example.com"
		obj := &model.PolicyVmApplianceConfig{
			AllowSshRootLogin: &allowSSHRoot,
			DnsServers:        []string{"8.8.8.8"},
			EnableSsh:         &enableSSH,
			EnableUptMode:     &enableUpt,
			SyslogServers: []model.SyslogConfiguration{
				{LogLevel: &logLevel, Port: &port, Protocol: &protocol, Server: &server},
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		err := setApplianceConfigInSchema(d, obj)
		require.NoError(t, err)
		cfg := d.Get("appliance_config").([]interface{})
		require.Len(t, cfg, 1)
		cfgMap := cfg[0].(map[string]interface{})
		assert.Equal(t, true, cfgMap["allow_ssh_root_login"])
		servers := cfgMap["syslog_server"].([]interface{})
		require.Len(t, servers, 1)
		assert.Equal(t, "syslog.example.com", servers[0].(map[string]interface{})["server"])
	})
}

func TestUnitNsxt_setCredentialsInSchema(t *testing.T) {
	res := resourceNsxtPolicyEdgeTransportNode()

	t.Run("nil object is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		err := setCredentialsInSchema(d, nil)
		require.NoError(t, err)
	})

	t.Run("merges usernames into existing credentials block", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"credentials": []interface{}{
				map[string]interface{}{
					"cli_username":  "admin",
					"cli_password":  "pw1",
					"root_password": "pw2",
				},
			},
		})
		auditUsername, cliUsername := "auditor", "admin2"
		obj := &model.PolicyEdgeTransportNodeCredential{AuditUsername: &auditUsername, CliUsername: &cliUsername}
		err := setCredentialsInSchema(d, obj)
		require.NoError(t, err)
		creds := d.Get("credentials").([]interface{})
		require.Len(t, creds, 1)
		credMap := creds[0].(map[string]interface{})
		assert.Equal(t, "admin2", credMap["cli_username"])
		assert.Equal(t, "auditor", credMap["audit_username"])
	})

	t.Run("creates credentials block when none set yet", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		cliUsername := "admin"
		obj := &model.PolicyEdgeTransportNodeCredential{CliUsername: &cliUsername}
		err := setCredentialsInSchema(d, obj)
		require.NoError(t, err)
		creds := d.Get("credentials").([]interface{})
		require.Len(t, creds, 1)
	})
}

func TestUnitNsxt_setPolicyIPAssignmentsInSchema(t *testing.T) {
	converter := bindings.NewTypeConverter()

	toStructValue := func(t *testing.T, obj interface{}, bindingType bindings.BindingType) *data.StructValue {
		t.Helper()
		dv, errs := converter.ConvertToVapi(obj, bindingType)
		require.Empty(t, errs)
		return dv.(*data.StructValue)
	}

	t.Run("auto_conf", func(t *testing.T) {
		sv := toStructValue(t, model.AutoConf{IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_AUTOCONF}, model.AutoConfBindingType())
		out, err := setPolicyIPAssignmentsInSchema([]*data.StructValue{sv})
		require.NoError(t, err)
		elem := out.([]interface{})[0].(map[string]interface{})
		assert.Equal(t, true, elem["auto_conf"])
	})

	t.Run("static_ipv4", func(t *testing.T) {
		gw := "10.0.0.1"
		sv := toStructValue(t, model.StaticIpv4{
			DefaultGateway:   []string{gw},
			IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV4,
			ManagementPortSubnets: []model.IPv4Subnet{
				{IpAddresses: []string{"10.0.0.2"}, PrefixLength: int64Ptr(24)},
			},
		}, model.StaticIpv4BindingType())
		out, err := setPolicyIPAssignmentsInSchema([]*data.StructValue{sv})
		require.NoError(t, err)
		elem := out.([]interface{})[0].(map[string]interface{})
		require.Contains(t, elem, "static_ipv4")
	})

	t.Run("static_ipv4_pool", func(t *testing.T) {
		pool := "/infra/ip-pools/pool-1"
		sv := toStructValue(t, model.StaticIpv4Pool{IpPool: &pool, IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV4POOL}, model.StaticIpv4PoolBindingType())
		out, err := setPolicyIPAssignmentsInSchema([]*data.StructValue{sv})
		require.NoError(t, err)
		elem := out.([]interface{})[0].(map[string]interface{})
		assert.Equal(t, &pool, elem["static_ipv4_pool"])
	})

	t.Run("static_ipv6_list", func(t *testing.T) {
		gw := "fe80::1"
		prefixLength := "64"
		sv := toStructValue(t, model.StaticIpv6List{
			DefaultGateway:   &gw,
			IpList:           []string{"fe80::2"},
			PrefixLength:     &prefixLength,
			IpAssignmentType: model.PolicyIpAssignmentSpec_IP_ASSIGNMENT_TYPE_STATICIPV6LIST,
		}, model.StaticIpv6ListBindingType())
		out, err := setPolicyIPAssignmentsInSchema([]*data.StructValue{sv})
		require.NoError(t, err)
		elem := out.([]interface{})[0].(map[string]interface{})
		require.Contains(t, elem, "static_ipv6_list")
	})
}
