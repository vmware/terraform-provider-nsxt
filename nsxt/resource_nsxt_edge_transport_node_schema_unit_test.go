//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mpmodel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

func TestUnitNsxt_getCPUConfigFromSchema(t *testing.T) {
	cfgs := getCPUConfigFromSchema([]interface{}{
		map[string]interface{}{"num_lcores": 4, "numa_node_index": 0},
		map[string]interface{}{"num_lcores": 2, "numa_node_index": 1},
	})
	require.Len(t, cfgs, 2)
	assert.EqualValues(t, 4, *cfgs[0].NumLcores)
	assert.EqualValues(t, 0, *cfgs[0].NumaNodeIndex)
	assert.EqualValues(t, 2, *cfgs[1].NumLcores)
	assert.EqualValues(t, 1, *cfgs[1].NumaNodeIndex)
}

func TestUnitNsxt_getHostSwitchProfileIDsFromSchema(t *testing.T) {
	t.Run("both uplink and ha profile set", func(t *testing.T) {
		parentMap := map[string]interface{}{
			"uplink_profile":      "/infra/host-switch-profiles/uplink1",
			"vtep_ha_profile":     "/infra/host-switch-profiles/ha1",
			"host_switch_profile": []interface{}{},
		}
		entries, err := getHostSwitchProfileIDsFromSchema(nil, parentMap)
		require.NoError(t, err)
		require.Len(t, entries, 2)
		assert.Equal(t, mpmodel.BaseHostSwitchProfile_RESOURCE_TYPE_UPLINKHOSTSWITCHPROFILE, *entries[0].Key)
		assert.Equal(t, "/infra/host-switch-profiles/uplink1", *entries[0].Value)
		assert.Equal(t, mpmodel.BaseHostSwitchProfile_RESOURCE_TYPE_VTEPHAHOSTSWITCHPROFILE, *entries[1].Key)
		assert.Equal(t, "/infra/host-switch-profiles/ha1", *entries[1].Value)
	})

	t.Run("only uplink profile set", func(t *testing.T) {
		parentMap := map[string]interface{}{
			"uplink_profile":      "/infra/host-switch-profiles/uplink1",
			"vtep_ha_profile":     "",
			"host_switch_profile": []interface{}{},
		}
		entries, err := getHostSwitchProfileIDsFromSchema(nil, parentMap)
		require.NoError(t, err)
		require.Len(t, entries, 1)
		assert.Equal(t, mpmodel.BaseHostSwitchProfile_RESOURCE_TYPE_UPLINKHOSTSWITCHPROFILE, *entries[0].Key)
	})

	t.Run("neither set and no deprecated list returns empty", func(t *testing.T) {
		parentMap := map[string]interface{}{
			"uplink_profile":      "",
			"vtep_ha_profile":     "",
			"host_switch_profile": []interface{}{},
		}
		entries, err := getHostSwitchProfileIDsFromSchema(nil, parentMap)
		require.NoError(t, err)
		assert.Empty(t, entries)
	})
}

func TestUnitNsxt_getTransportZoneEndpointsFromSchema(t *testing.T) {
	endpoints := getTransportZoneEndpointsFromSchema([]interface{}{
		map[string]interface{}{
			"transport_zone":          "/infra/sites/default/enforcement-points/default/transport-zones/tz1",
			"transport_zone_profiles": []interface{}{"profile-1"},
		},
	})
	require.Len(t, endpoints, 1)
	assert.Equal(t, "/infra/sites/default/enforcement-points/default/transport-zones/tz1", *endpoints[0].TransportZoneId)
	require.Len(t, endpoints[0].TransportZoneProfileIds, 1)
	assert.Equal(t, "profile-1", *endpoints[0].TransportZoneProfileIds[0].ProfileId)
}

func TestUnitNsxt_getUplinksFromSchema(t *testing.T) {
	uplinks := getUplinksFromSchema([]interface{}{
		map[string]interface{}{"uplink_name": "uplink-1", "vds_lag_name": "lag-1", "vds_uplink_name": "vds-uplink-1"},
	})
	require.Len(t, uplinks, 1)
	assert.Equal(t, "uplink-1", *uplinks[0].UplinkName)
	assert.Equal(t, "lag-1", *uplinks[0].VdsLagName)
	assert.Equal(t, "vds-uplink-1", *uplinks[0].VdsUplinkName)
}

func TestUnitNsxt_getTransportNodeSubProfileCfg(t *testing.T) {
	t.Run("nil iface returns empty list", func(t *testing.T) {
		cfgs, err := getTransportNodeSubProfileCfg(nil, nil)
		require.NoError(t, err)
		assert.Empty(t, cfgs)
	})

	t.Run("builds sub-profile config without needing the API client", func(t *testing.T) {
		iface := []interface{}{
			map[string]interface{}{
				"name": "sub-config-1",
				"host_switch_config_option": []interface{}{
					map[string]interface{}{
						"host_switch_id":      "hsw-1",
						"uplink_profile":      "/infra/host-switch-profiles/uplink1",
						"vtep_ha_profile":     "",
						"host_switch_profile": []interface{}{},
						"ip_assignment":       []interface{}{},
						"uplink":              []interface{}{},
					},
				},
			},
		}
		cfgs, err := getTransportNodeSubProfileCfg(nil, iface)
		require.NoError(t, err)
		require.Len(t, cfgs, 1)
		assert.Equal(t, "sub-config-1", *cfgs[0].Name)
		require.NotNil(t, cfgs[0].HostSwitchConfigOption)
		assert.Equal(t, "hsw-1", *cfgs[0].HostSwitchConfigOption.HostSwitchId)
		require.Len(t, cfgs[0].HostSwitchConfigOption.HostSwitchProfileIds, 1)
	})
}

func TestUnitNsxt_getIPv6AssignmentFromSchema(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		sv, err := getIPv6AssignmentFromSchema(nil)
		require.NoError(t, err)
		assert.Nil(t, sv)
	})

	t.Run("empty list returns nil", func(t *testing.T) {
		sv, err := getIPv6AssignmentFromSchema([]interface{}{})
		require.NoError(t, err)
		assert.Nil(t, sv)
	})

	t.Run("assigned_by_dhcpv6", func(t *testing.T) {
		sv, err := getIPv6AssignmentFromSchema([]interface{}{
			map[string]interface{}{"assigned_by_dhcpv6": true},
		})
		require.NoError(t, err)
		require.NotNil(t, sv)
	})

	t.Run("assigned_by_autoconf", func(t *testing.T) {
		sv, err := getIPv6AssignmentFromSchema([]interface{}{
			map[string]interface{}{"assigned_by_autoconf": true},
		})
		require.NoError(t, err)
		require.NotNil(t, sv)
	})

	t.Run("no_ipv6", func(t *testing.T) {
		sv, err := getIPv6AssignmentFromSchema([]interface{}{
			map[string]interface{}{"no_ipv6": true},
		})
		require.NoError(t, err)
		require.NotNil(t, sv)
	})

	t.Run("static_ip", func(t *testing.T) {
		sv, err := getIPv6AssignmentFromSchema([]interface{}{
			map[string]interface{}{
				"static_ip": []interface{}{
					map[string]interface{}{
						"default_gateway": "fe80::1",
						"ip_addresses":    []interface{}{"fe80::2"},
						"prefix_length":   "64",
					},
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, sv)
	})

	t.Run("static_ip_mac", func(t *testing.T) {
		sv, err := getIPv6AssignmentFromSchema([]interface{}{
			map[string]interface{}{
				"static_ip_mac": []interface{}{
					map[string]interface{}{
						"default_gateway": "fe80::1",
						"ip_mac_pair": []interface{}{
							map[string]interface{}{"ip_address": "fe80::2", "mac_address": "aa:bb:cc:dd:ee:ff"},
						},
						"prefix_length": "64",
					},
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, sv)
	})

	t.Run("static_ip_pool", func(t *testing.T) {
		sv, err := getIPv6AssignmentFromSchema([]interface{}{
			map[string]interface{}{"static_ip_pool": "pool-1"},
		})
		require.NoError(t, err)
		require.NotNil(t, sv)
	})

	t.Run("no assignment set errors", func(t *testing.T) {
		_, err := getIPv6AssignmentFromSchema([]interface{}{
			map[string]interface{}{},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no valid IP assignment found")
	})

	t.Run("multiple assignments set errors", func(t *testing.T) {
		_, err := getIPv6AssignmentFromSchema([]interface{}{
			map[string]interface{}{"assigned_by_dhcpv6": true, "no_ipv6": true},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "exactly one IP assignment")
	})
}

func TestUnitNsxt_setTransportZoneEndpointInSchema(t *testing.T) {
	tzID := "/infra/sites/default/enforcement-points/default/transport-zones/tz1"
	profileID := "profile-1"
	endpoints := []mpmodel.TransportZoneEndPoint{
		{
			TransportZoneId: &tzID,
			TransportZoneProfileIds: []mpmodel.TransportZoneProfileTypeIdEntry{
				{ProfileId: &profileID},
			},
		},
	}
	result := setTransportZoneEndpointInSchema(endpoints).([]map[string]interface{})
	require.Len(t, result, 1)
	assert.Equal(t, &tzID, result[0]["transport_zone"])
	assert.Equal(t, []string{profileID}, result[0]["transport_zone_profiles"])
}

func TestUnitNsxt_setUplinksFromSchema(t *testing.T) {
	name := "uplink-1"
	lag := "lag-1"
	vds := "vds-uplink-1"
	uplinks := []mpmodel.VdsUplink{{UplinkName: &name, VdsLagName: &lag, VdsUplinkName: &vds}}

	result := setUplinksFromSchema(uplinks).([]map[string]interface{})
	require.Len(t, result, 1)
	assert.Equal(t, &name, result[0]["uplink_name"])
	assert.Equal(t, &lag, result[0]["vds_lag_name"])
	assert.Equal(t, &vds, result[0]["vds_uplink_name"])
}

func TestUnitNsxt_setIPv6AssignmentInSchema_roundTrip(t *testing.T) {
	sv, err := getIPv6AssignmentFromSchema([]interface{}{
		map[string]interface{}{"assigned_by_dhcpv6": true},
	})
	require.NoError(t, err)

	result, err := setIPv6AssignmentInSchema(sv)
	require.NoError(t, err)
	list := result.([]interface{})
	require.Len(t, list, 1)
	elem := list[0].(map[string]interface{})
	assert.Equal(t, true, elem["assigned_by_dhcpv6"])
}

func TestUnitNsxt_setHostSwitchProfileIDsInSchema(t *testing.T) {
	uplinkType := mpmodel.BaseHostSwitchProfile_RESOURCE_TYPE_UPLINKHOSTSWITCHPROFILE
	haType := mpmodel.BaseHostSwitchProfile_RESOURCE_TYPE_VTEPHAHOSTSWITCHPROFILE
	uplinkVal := "/infra/host-switch-profiles/uplink1"
	haVal := "/infra/host-switch-profiles/ha1"

	entries := []mpmodel.HostSwitchProfileTypeIdEntry{
		{Key: &uplinkType, Value: &uplinkVal},
		{Key: &haType, Value: &haVal},
	}
	parentMap := map[string]interface{}{}
	setHostSwitchProfileIDsInSchema(entries, parentMap)

	assert.Equal(t, &uplinkVal, parentMap["uplink_profile"])
	assert.Equal(t, &haVal, parentMap["vtep_ha_profile"])
	assert.Equal(t, []interface{}{&uplinkVal, &haVal}, parentMap["host_switch_profile"])
}
