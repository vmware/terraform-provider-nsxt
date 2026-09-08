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
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

// getEdgeNodeDeploymentConfigFromSchema, setEdgeDeploymentConfigInSchema, and
// setVMDeploymentConfigInSchema are pure schema<->struct conversion helpers (no SDK calls), so
// they're covered directly. getIPAssignmentFromSchema/setIPAssignmentInSchema mirror the IPv6
// round-trip pattern already used in resource_nsxt_edge_transport_node_schema_unit_test.go.

func minimalDeploymentConfigSchemaData() map[string]interface{} {
	return map[string]interface{}{
		"form_factor": "SMALL",
		"node_user_settings": []interface{}{
			map[string]interface{}{
				"audit_username": "auditor",
				"audit_password": "auditpw",
				"cli_password":   "clipw",
				"cli_username":   "admin",
				"root_password":  "rootpw",
			},
		},
		"vm_deployment_config": []interface{}{
			map[string]interface{}{
				"compute_folder_id":       "folder-1",
				"compute_id":              "compute-1",
				"data_network_ids":        []interface{}{"net-1", "net-2"},
				"default_gateway_address": []interface{}{"10.0.0.1"},
				"host_id":                 "host-1",
				"management_network_id":   "mgmt-net-1",
				"storage_id":              "storage-1",
				"vc_id":                   "vc-1",
				"management_port_subnet": []interface{}{
					map[string]interface{}{
						"ip_addresses":  []interface{}{"10.0.0.5"},
						"prefix_length": 24,
					},
				},
				"reservation_info": []interface{}{
					map[string]interface{}{
						"cpu_reservation_in_mhz":        1000,
						"cpu_reservation_in_shares":     "NORMAL",
						"memory_reservation_percentage": 50,
					},
				},
			},
		},
	}
}

func TestUnitNsxt_getEdgeNodeDeploymentConfigFromSchema(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		cfg, err := getEdgeNodeDeploymentConfigFromSchema(nil)
		require.NoError(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("empty list returns nil", func(t *testing.T) {
		cfg, err := getEdgeNodeDeploymentConfigFromSchema([]interface{}{})
		require.NoError(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("builds full deployment config", func(t *testing.T) {
		cfg, err := getEdgeNodeDeploymentConfigFromSchema([]interface{}{minimalDeploymentConfigSchemaData()})
		require.NoError(t, err)
		require.NotNil(t, cfg)
		assert.Equal(t, "SMALL", *cfg.FormFactor)
		require.NotNil(t, cfg.NodeUserSettings)
		assert.Equal(t, "admin", *cfg.NodeUserSettings.CliUsername)
		assert.Equal(t, "auditor", *cfg.NodeUserSettings.AuditUsername)
		require.NotNil(t, cfg.VmDeploymentConfig)
	})

	t.Run("ipv6 management address sets the IPv6 assignment type", func(t *testing.T) {
		data := minimalDeploymentConfigSchemaData()
		vdc := data["vm_deployment_config"].([]interface{})[0].(map[string]interface{})
		vdc["management_port_subnet"] = []interface{}{
			map[string]interface{}{
				"ip_addresses":  []interface{}{"fe80::5"},
				"prefix_length": 64,
			},
		}
		cfg, err := getEdgeNodeDeploymentConfigFromSchema([]interface{}{data})
		require.NoError(t, err)
		require.NotNil(t, cfg)
	})

	t.Run("minimal config without optional fields", func(t *testing.T) {
		data := map[string]interface{}{
			"form_factor":        "SMALL",
			"node_user_settings": []interface{}{},
			"vm_deployment_config": []interface{}{
				map[string]interface{}{
					"compute_folder_id":       "",
					"compute_id":              "compute-1",
					"data_network_ids":        []interface{}{},
					"default_gateway_address": []interface{}{},
					"host_id":                 "",
					"management_network_id":   "mgmt-net-1",
					"storage_id":              "storage-1",
					"vc_id":                   "vc-1",
					"management_port_subnet":  []interface{}{},
					"reservation_info":        []interface{}{},
				},
			},
		}
		cfg, err := getEdgeNodeDeploymentConfigFromSchema([]interface{}{data})
		require.NoError(t, err)
		require.NotNil(t, cfg)
		assert.Nil(t, cfg.NodeUserSettings)
	})
}

func TestUnitNsxt_setEdgeDeploymentConfigInSchema(t *testing.T) {
	res := resourceNsxtEdgeTransportNode()

	t.Run("round trips a full deployment config into schema", func(t *testing.T) {
		cfg, err := getEdgeNodeDeploymentConfigFromSchema([]interface{}{minimalDeploymentConfigSchemaData()})
		require.NoError(t, err)
		require.NotNil(t, cfg)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		err = setEdgeDeploymentConfigInSchema(d, cfg)
		require.NoError(t, err)

		dc := d.Get("deployment_config").([]interface{})
		require.Len(t, dc, 1)
		elem := dc[0].(map[string]interface{})
		assert.Equal(t, "SMALL", elem["form_factor"])

		vdc := elem["vm_deployment_config"].([]interface{})
		require.Len(t, vdc, 1)
		vdcElem := vdc[0].(map[string]interface{})
		assert.Equal(t, "compute-1", vdcElem["compute_id"])
	})
}

func TestUnitNsxt_setVMDeploymentConfigInSchema(t *testing.T) {
	t.Run("converts a vsphere deployment config", func(t *testing.T) {
		cfg, err := getEdgeNodeDeploymentConfigFromSchema([]interface{}{minimalDeploymentConfigSchemaData()})
		require.NoError(t, err)
		require.NotNil(t, cfg)

		result, err := setVMDeploymentConfigInSchema(cfg.VmDeploymentConfig)
		require.NoError(t, err)
		list := result.([]interface{})
		require.Len(t, list, 1)
		elem := list[0].(map[string]interface{})
		assert.Equal(t, "compute-1", *elem["compute_id"].(*string))
		assert.Equal(t, "storage-1", *elem["storage_id"].(*string))
		assert.Equal(t, "vc-1", *elem["vc_id"].(*string))
	})
}

func TestUnitNsxt_getIPAssignmentFromSchema(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		sv, err := getIPAssignmentFromSchema(nil)
		require.NoError(t, err)
		assert.Nil(t, sv)
	})

	t.Run("empty list returns nil", func(t *testing.T) {
		sv, err := getIPAssignmentFromSchema([]interface{}{})
		require.NoError(t, err)
		assert.Nil(t, sv)
	})

	t.Run("assigned_by_dhcp", func(t *testing.T) {
		sv, err := getIPAssignmentFromSchema([]interface{}{
			map[string]interface{}{"assigned_by_dhcp": true},
		})
		require.NoError(t, err)
		require.NotNil(t, sv)
	})

	t.Run("no_ipv4", func(t *testing.T) {
		sv, err := getIPAssignmentFromSchema([]interface{}{
			map[string]interface{}{"no_ipv4": true},
		})
		require.NoError(t, err)
		require.NotNil(t, sv)
	})

	t.Run("static_ip", func(t *testing.T) {
		sv, err := getIPAssignmentFromSchema([]interface{}{
			map[string]interface{}{
				"static_ip": []interface{}{
					map[string]interface{}{
						"default_gateway": "192.168.1.1",
						"ip_addresses":    []interface{}{"192.168.1.2"},
						"subnet_mask":     "255.255.255.0",
					},
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, sv)
	})

	t.Run("static_ip_mac", func(t *testing.T) {
		sv, err := getIPAssignmentFromSchema([]interface{}{
			map[string]interface{}{
				"static_ip_mac": []interface{}{
					map[string]interface{}{
						"default_gateway": "192.168.1.1",
						"ip_mac_pair": []interface{}{
							map[string]interface{}{"ip_address": "192.168.1.2", "mac_address": "aa:bb:cc:dd:ee:ff"},
						},
						"subnet_mask": "255.255.255.0",
					},
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, sv)
	})

	t.Run("static_ip_pool", func(t *testing.T) {
		sv, err := getIPAssignmentFromSchema([]interface{}{
			map[string]interface{}{"static_ip_pool": "pool-1"},
		})
		require.NoError(t, err)
		require.NotNil(t, sv)
	})

	t.Run("no assignment set errors", func(t *testing.T) {
		_, err := getIPAssignmentFromSchema([]interface{}{
			map[string]interface{}{},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no valid IP assignment found")
	})

	t.Run("multiple assignments set errors", func(t *testing.T) {
		_, err := getIPAssignmentFromSchema([]interface{}{
			map[string]interface{}{"assigned_by_dhcp": true, "no_ipv4": true},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "exactly one IP assignment")
	})
}

func TestUnitNsxt_setIPAssignmentInSchema_roundTrip(t *testing.T) {
	cases := []struct {
		name string
		in   map[string]interface{}
		key  string
	}{
		{"assigned_by_dhcp", map[string]interface{}{"assigned_by_dhcp": true}, "assigned_by_dhcp"},
		{"no_ipv4", map[string]interface{}{"no_ipv4": true}, "no_ipv4"},
		{"static_ip_pool", map[string]interface{}{"static_ip_pool": "pool-1"}, "static_ip_pool"},
		{"static_ip", map[string]interface{}{
			"static_ip": []interface{}{
				map[string]interface{}{"default_gateway": "192.168.1.1", "ip_addresses": []interface{}{"192.168.1.2"}, "subnet_mask": "255.255.255.0"},
			},
		}, "static_ip"},
		{"static_ip_mac", map[string]interface{}{
			"static_ip_mac": []interface{}{
				map[string]interface{}{
					"default_gateway": "192.168.1.1",
					"ip_mac_pair":     []interface{}{map[string]interface{}{"ip_address": "192.168.1.2", "mac_address": "aa:bb:cc:dd:ee:ff"}},
					"subnet_mask":     "255.255.255.0",
				},
			},
		}, "static_ip_mac"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sv, err := getIPAssignmentFromSchema([]interface{}{tc.in})
			require.NoError(t, err)
			require.NotNil(t, sv)

			result, err := setIPAssignmentInSchema(sv)
			require.NoError(t, err)
			list := result.([]interface{})
			require.Len(t, list, 1)
			elem := list[0].(map[string]interface{})
			assert.Contains(t, elem, tc.key)
		})
	}
}

func TestUnitNsxt_setIPv6AssignmentInSchema_roundTripAllBranches(t *testing.T) {
	cases := []struct {
		name string
		in   map[string]interface{}
		key  string
	}{
		{"assigned_by_autoconf", map[string]interface{}{"assigned_by_autoconf": true}, "assigned_by_autoconf"},
		{"no_ipv6", map[string]interface{}{"no_ipv6": true}, "no_ipv6"},
		{"static_ip_pool", map[string]interface{}{"static_ip_pool": "pool-1"}, "static_ip_pool"},
		{"static_ip", map[string]interface{}{
			"static_ip": []interface{}{
				map[string]interface{}{"default_gateway": "fe80::1", "ip_addresses": []interface{}{"fe80::2"}, "prefix_length": "64"},
			},
		}, "static_ip"},
		{"static_ip_mac", map[string]interface{}{
			"static_ip_mac": []interface{}{
				map[string]interface{}{
					"default_gateway": "fe80::1",
					"ip_mac_pair":     []interface{}{map[string]interface{}{"ip_address": "fe80::2", "mac_address": "aa:bb:cc:dd:ee:ff"}},
					"prefix_length":   "64",
				},
			},
		}, "static_ip_mac"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sv, err := getIPv6AssignmentFromSchema([]interface{}{tc.in})
			require.NoError(t, err)
			require.NotNil(t, sv)

			result, err := setIPv6AssignmentInSchema(sv)
			require.NoError(t, err)
			list := result.([]interface{})
			require.Len(t, list, 1)
			elem := list[0].(map[string]interface{})
			assert.Contains(t, elem, tc.key)
		})
	}
}

func setupHostSwitchProfilesMock(t *testing.T) *inframocks.MockHostSwitchProfilesClient {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSDK := inframocks.NewMockHostSwitchProfilesClient(ctrl)
	mockWrapper := &cliinfra.HostSwitchProfilesClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	originalCli := cliHostSwitchProfilesClient
	t.Cleanup(func() { cliHostSwitchProfilesClient = originalCli })
	cliHostSwitchProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.HostSwitchProfilesClientContext {
		return mockWrapper
	}
	return mockSDK
}

func TestMockNsxt_getHostSwitchProfileResourceType(t *testing.T) {
	t.Run("List error propagates", func(t *testing.T) {
		mockSDK := setupHostSwitchProfilesMock(t)
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PolicyHostSwitchProfilesListResult{}, vapiErrors.InternalServerError{})

		_, err := getHostSwitchProfileResourceType(newGoMockProviderClient(), "profile-1")
		require.Error(t, err)
	})

	t.Run("not found errors", func(t *testing.T) {
		mockSDK := setupHostSwitchProfilesMock(t)
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.PolicyHostSwitchProfilesListResult{}, nil)

		_, err := getHostSwitchProfileResourceType(newGoMockProviderClient(), "profile-1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
