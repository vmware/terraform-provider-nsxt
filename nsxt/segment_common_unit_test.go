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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestUnitNsxt_parseSegmentPolicyPath(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantIsT0  bool
		wantGwID  string
		wantSegID string
	}{
		{"infra segment", "/infra/segments/seg-1", false, "", "seg-1"},
		{"project segment", "/orgs/o1/projects/p1/infra/segments/seg-1", false, "", "seg-1"},
		{"fixed tier-1 segment", "/infra/tier-1s/gw-1/segments/seg-1", false, "gw-1", "seg-1"},
		{"fixed tier-0 segment", "/infra/tier-0s/gw-1/segments/seg-1", true, "gw-1", "seg-1"},
		{"not a segment path", "/infra/tier-1s/gw-1", false, "", ""},
		{"too short", "seg-1", false, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isT0, gwID, segID := parseSegmentPolicyPath(tt.path)
			assert.Equal(t, tt.wantIsT0, isT0)
			assert.Equal(t, tt.wantGwID, gwID)
			assert.Equal(t, tt.wantSegID, segID)
		})
	}
}

func TestUnitNsxt_setBridgeConfigInStructAndSchema(t *testing.T) {
	segSchema := getPolicyCommonSegmentSchema(false, false)

	t.Run("no bridge config is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{})
		obj := model.Segment{}
		setBridgeConfigInStruct(d, &obj)
		assert.Empty(t, obj.BridgeProfiles)
	})

	t.Run("bridge config round-trips through struct and back to schema", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{
			"bridge_config": []interface{}{
				map[string]interface{}{
					"profile_path":          "/infra/bridge-profiles/bp-1",
					"uplink_teaming_policy": "policy-1",
					"vlan_ids":              []interface{}{"100"},
					"transport_zone_path":   "/infra/transport-zones/tz-1",
				},
			},
		})
		obj := model.Segment{}
		setBridgeConfigInStruct(d, &obj)
		require.Len(t, obj.BridgeProfiles, 1)
		assert.Equal(t, "/infra/bridge-profiles/bp-1", *obj.BridgeProfiles[0].BridgeProfilePath)

		d2 := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{})
		setSegmentBridgeConfigInSchema(d2, &obj)
		got := d2.Get("bridge_config").([]interface{})
		require.Len(t, got, 1)
		assert.Equal(t, "/infra/bridge-profiles/bp-1", got[0].(map[string]interface{})["profile_path"])
	})
}

func TestUnitNsxt_nsxtPolicySegmentAddGatewayToInfraStruct(t *testing.T) {
	segSchema := getPolicyCommonSegmentSchema(true, true)
	dummySegID := "dummy-seg"
	dummyTargetType := "Segment"
	dummyChild, err := vAPIConversion(model.ChildResourceReference{
		Id:           &dummySegID,
		ResourceType: "ChildResourceReference",
		TargetType:   &dummyTargetType,
	}, model.ChildResourceReferenceBindingType())
	require.NoError(t, err)

	t.Run("valid tier-1 connectivity_path succeeds", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{
			"connectivity_path": "/infra/tier-1s/gw-1",
		})
		val, err := nsxtPolicySegmentAddGatewayToInfraStruct(d, dummyChild)
		require.NoError(t, err)
		assert.NotNil(t, val)
	})

	t.Run("tier-0 connectivity_path is rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{
			"connectivity_path": "/infra/tier-0s/gw-1",
		})
		_, err := nsxtPolicySegmentAddGatewayToInfraStruct(d, dummyChild)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Tier0")
	})

	t.Run("empty connectivity_path is rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{})
		_, err := nsxtPolicySegmentAddGatewayToInfraStruct(d, dummyChild)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not a valid gateway path")
	})
}

func TestUnitNsxt_nsxtPolicySegmentProfileSetInStruct(t *testing.T) {
	segSchema := getPolicyCommonSegmentSchema(false, false)

	t.Run("discovery profile: no old, no new is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{})
		val, err := nsxtPolicySegmentDiscoveryProfileSetInStruct(d)
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("discovery profile: new profile configured", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{
			"discovery_profile": []interface{}{
				map[string]interface{}{
					"ip_discovery_profile_path":  "/infra/ip-discovery-profiles/p1",
					"mac_discovery_profile_path": "",
					"binding_map_path":           "",
					"revision":                   0,
				},
			},
		})
		val, err := nsxtPolicySegmentDiscoveryProfileSetInStruct(d)
		require.NoError(t, err)
		assert.NotNil(t, val)
	})

	t.Run("qos profile: new profile configured", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{
			"qos_profile": []interface{}{
				map[string]interface{}{
					"qos_profile_path": "/infra/qos-profiles/p1",
					"binding_map_path": "",
					"revision":         0,
				},
			},
		})
		val, err := nsxtPolicySegmentQosProfileSetInStruct(d)
		require.NoError(t, err)
		assert.NotNil(t, val)
	})

	t.Run("qos profile: no old, no new is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{})
		val, err := nsxtPolicySegmentQosProfileSetInStruct(d)
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("security profile: new profile configured", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{
			"security_profile": []interface{}{
				map[string]interface{}{
					"spoofguard_profile_path": "/infra/spoofguard-profiles/p1",
					"security_profile_path":   "/infra/segment-security-profiles/p1",
					"binding_map_path":        "",
					"revision":                0,
				},
			},
		})
		val, err := nsxtPolicySegmentSecurityProfileSetInStruct(d)
		require.NoError(t, err)
		assert.NotNil(t, val)
	})

	t.Run("security profile: no old, no new is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{})
		val, err := nsxtPolicySegmentSecurityProfileSetInStruct(d)
		require.NoError(t, err)
		assert.Nil(t, val)
	})
}

func TestUnitNsxt_policySegmentResourceToInfraStruct(t *testing.T) {
	m := newGoMockProviderClient()

	t.Run("vlan segment builds an Infra struct", func(t *testing.T) {
		segSchema := getPolicyCommonSegmentSchema(true, false)
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{
			"display_name":        "seg-1",
			"transport_zone_path": "/infra/transport-zones/tz-1",
			"vlan_ids":            []interface{}{"100"},
			"replication_mode":    model.Segment_REPLICATION_MODE_MTEP,
		})
		obj, err := policySegmentResourceToInfraStruct(getSessionContext(d, m), "seg-1", d, m, true, false)
		require.NoError(t, err)
		assert.NotNil(t, obj.Children)
	})

	t.Run("missing transport_zone_path on local manager overlay segment is an error", func(t *testing.T) {
		segSchema := getPolicyCommonSegmentSchema(false, false)
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{
			"display_name":      "seg-1",
			"connectivity_path": "/infra/tier-1s/gw-1",
			"replication_mode":  model.Segment_REPLICATION_MODE_MTEP,
		})
		_, err := policySegmentResourceToInfraStruct(getSessionContext(d, m), "seg-1", d, m, false, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "transport_zone_path")
	})
}
