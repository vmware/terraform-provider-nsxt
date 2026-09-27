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
)

func TestUnitNsxt_etNodeShallowCopyMap(t *testing.T) {
	orig := map[string]interface{}{"a": 1, "b": "x"}
	cp := etNodeShallowCopyMap(orig)
	assert.Equal(t, orig, cp)

	cp["a"] = 2
	assert.Equal(t, 1, orig["a"], "mutating the copy must not affect the original")
}

func TestUnitNsxt_etNodeStringVal(t *testing.T) {
	m := map[string]interface{}{"key": "value", "num": 5}
	assert.Equal(t, "value", etNodeStringVal(m, "key"))
	assert.Equal(t, "", etNodeStringVal(m, "missing"))
	assert.Equal(t, "", etNodeStringVal(m, "num"))
}

func TestUnitNsxt_etNodeNestedList(t *testing.T) {
	list := []interface{}{"a", "b"}
	m := map[string]interface{}{"list": list, "notlist": "x"}
	assert.Equal(t, list, etNodeNestedList(m, "list"))
	assert.Nil(t, etNodeNestedList(m, "missing"))
	assert.Nil(t, etNodeNestedList(m, "notlist"))
}

func TestUnitNsxt_etNodeNormalizeStringList(t *testing.T) {
	t.Run("prior policy path replaces bare API value", func(t *testing.T) {
		prior := map[string]interface{}{"host_switch_profile": []interface{}{"/infra/host-switch-profiles/hsp1"}}
		newSw := map[string]interface{}{"host_switch_profile": []interface{}{"hsp-uuid-1"}}

		modified := etNodeNormalizeStringList(prior, newSw, "host_switch_profile")
		require.True(t, modified)
		assert.Equal(t, []interface{}{"/infra/host-switch-profiles/hsp1"}, newSw["host_switch_profile"])
	})

	t.Run("prior bare value is not treated as a policy path", func(t *testing.T) {
		prior := map[string]interface{}{"host_switch_profile": []interface{}{"hsp-uuid-old"}}
		newSw := map[string]interface{}{"host_switch_profile": []interface{}{"hsp-uuid-new"}}

		modified := etNodeNormalizeStringList(prior, newSw, "host_switch_profile")
		require.False(t, modified)
		assert.Equal(t, []interface{}{"hsp-uuid-new"}, newSw["host_switch_profile"])
	})

	t.Run("length mismatch is a no-op", func(t *testing.T) {
		prior := map[string]interface{}{"host_switch_profile": []interface{}{"/infra/host-switch-profiles/hsp1"}}
		newSw := map[string]interface{}{"host_switch_profile": []interface{}{}}

		modified := etNodeNormalizeStringList(prior, newSw, "host_switch_profile")
		require.False(t, modified)
	})
}

func TestUnitNsxt_etNodeNormalizeTZEndpoints(t *testing.T) {
	t.Run("prior policy path replaces bare API transport_zone", func(t *testing.T) {
		prior := map[string]interface{}{
			"transport_zone_endpoint": []interface{}{
				map[string]interface{}{"transport_zone": "/infra/sites/default/enforcement-points/default/transport-zones/tz1"},
			},
		}
		newSw := map[string]interface{}{
			"transport_zone_endpoint": []interface{}{
				map[string]interface{}{"transport_zone": "tz-uuid-1"},
			},
		}

		modified := etNodeNormalizeTZEndpoints(prior, newSw)
		require.True(t, modified)
		endpoints := newSw["transport_zone_endpoint"].([]interface{})
		ep := endpoints[0].(map[string]interface{})
		assert.Equal(t, "/infra/sites/default/enforcement-points/default/transport-zones/tz1", ep["transport_zone"])
	})

	t.Run("length mismatch is a no-op", func(t *testing.T) {
		prior := map[string]interface{}{
			"transport_zone_endpoint": []interface{}{
				map[string]interface{}{"transport_zone": "/infra/x"},
			},
		}
		newSw := map[string]interface{}{"transport_zone_endpoint": []interface{}{}}

		modified := etNodeNormalizeTZEndpoints(prior, newSw)
		require.False(t, modified)
	})
}

func TestUnitNsxt_etNodeNormalizeHostSwitchesInState(t *testing.T) {
	res := resourceNsxtEdgeTransportNode()

	t.Run("empty prior switches is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		etNodeNormalizeHostSwitchesInState(d, nil)
		assert.Empty(t, d.Get("standard_host_switch").([]interface{}))
	})

	t.Run("restores policy paths for uplink_profile, vtep_ha_profile, host_switch_profile and transport_zone", func(t *testing.T) {
		uplinkPolicyPath := "/infra/host-switch-profiles/uplink1"
		haPolicyPath := "/infra/host-switch-profiles/ha1"
		hspPolicyPath := "/infra/host-switch-profiles/hsp1"
		tzPolicyPath := "/infra/sites/default/enforcement-points/default/transport-zones/tz1"

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"standard_host_switch": []interface{}{
				map[string]interface{}{
					"uplink_profile":      "uplink-uuid-1",
					"vtep_ha_profile":     "ha-uuid-1",
					"host_switch_profile": []interface{}{"hsp-uuid-1"},
					"transport_zone_endpoint": []interface{}{
						map[string]interface{}{"transport_zone": "tz-uuid-1"},
					},
				},
			},
		})

		priorSwitches := []interface{}{
			map[string]interface{}{
				"uplink_profile":      uplinkPolicyPath,
				"vtep_ha_profile":     haPolicyPath,
				"host_switch_profile": []interface{}{hspPolicyPath},
				"transport_zone_endpoint": []interface{}{
					map[string]interface{}{"transport_zone": tzPolicyPath},
				},
			},
		}

		etNodeNormalizeHostSwitchesInState(d, priorSwitches)

		switches := d.Get("standard_host_switch").([]interface{})
		require.Len(t, switches, 1)
		sw := switches[0].(map[string]interface{})
		assert.Equal(t, uplinkPolicyPath, sw["uplink_profile"])
		assert.Equal(t, haPolicyPath, sw["vtep_ha_profile"])
		assert.Equal(t, []interface{}{hspPolicyPath}, sw["host_switch_profile"])
		endpoints := sw["transport_zone_endpoint"].([]interface{})
		ep := endpoints[0].(map[string]interface{})
		assert.Equal(t, tzPolicyPath, ep["transport_zone"])
	})

	t.Run("length mismatch between prior and API state is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"standard_host_switch": []interface{}{
				map[string]interface{}{"uplink_profile": "uplink-uuid-1"},
			},
		})
		priorSwitches := []interface{}{
			map[string]interface{}{"uplink_profile": "/infra/host-switch-profiles/uplink1"},
			map[string]interface{}{"uplink_profile": "/infra/host-switch-profiles/uplink2"},
		}

		etNodeNormalizeHostSwitchesInState(d, priorSwitches)

		switches := d.Get("standard_host_switch").([]interface{})
		require.Len(t, switches, 1)
		sw := switches[0].(map[string]interface{})
		assert.Equal(t, "uplink-uuid-1", sw["uplink_profile"])
	})
}
