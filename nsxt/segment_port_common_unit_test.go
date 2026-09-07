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

func TestUnitNsxt_isT1Segment(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"tier-1 segment port path", "/infra/tier-1s/gw-1/segments/seg-1", true},
		{"infra segment path", "/infra/segments/seg-1", false},
		{"too short", "seg-1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isT1Segment(tt.path))
		})
	}
}

func TestUnitNsxt_getSegmentIdFromSegPath(t *testing.T) {
	assert.Equal(t, "seg-1", getSegmentIdFromSegPath("/infra/segments/seg-1"))
	assert.Equal(t, "seg-1", getSegmentIdFromSegPath("/infra/tier-1s/gw-1/segments/seg-1"))
}

func TestUnitNsxt_getT1IdFromSegPath(t *testing.T) {
	assert.Equal(t, "gw-1", getT1IdFromSegPath("/infra/tier-1s/gw-1/segments/seg-1"))
	assert.Equal(t, "", getT1IdFromSegPath("/infra/segments/seg-1"))
}

func TestUnitNsxt_getPolicySegmentPathFromPortPath(t *testing.T) {
	t.Run("valid port path", func(t *testing.T) {
		got, err := getPolicySegmentPathFromPortPath("/infra/segments/seg-1/ports/port-1")
		require.NoError(t, err)
		assert.Equal(t, "/infra/segments/seg-1", got)
	})

	t.Run("path without /ports/ errors", func(t *testing.T) {
		_, err := getPolicySegmentPathFromPortPath("/infra/segments/seg-1")
		require.Error(t, err)
	})
}

func segmentPortTestSchema() map[string]*schema.Schema {
	return resourceNsxtPolicySegmentPort().Schema
}

func TestUnitNsxt_nsxtPolicySegmentPortAttachmentConfigSetInStruct(t *testing.T) {
	portSchema := segmentPortTestSchema()

	t.Run("no attachment is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{})
		obj := model.SegmentPort{}
		nsxtPolicySegmentPortAttachmentConfigSetInStruct(d, &obj)
		assert.Nil(t, obj.Attachment)
	})

	t.Run("attachment fields are set from schema", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{
			"attachment": []interface{}{
				map[string]interface{}{
					"id":                 "att-1",
					"allocate_addresses": "BOTH",
					"app_id":             "app-1",
					"context_id":         "ctx-1",
					"context_type":       "VIF",
					"evpn_vlans":         []interface{}{"100", "200"},
					"hyperbus_mode":      "NONE",
					"type":               "STATIC",
					"traffic_tag":        5,
				},
			},
		})
		obj := model.SegmentPort{}
		nsxtPolicySegmentPortAttachmentConfigSetInStruct(d, &obj)
		require.NotNil(t, obj.Attachment)
		assert.Equal(t, "att-1", *obj.Attachment.Id)
		assert.Equal(t, "app-1", *obj.Attachment.AppId)
		assert.Equal(t, []string{"100", "200"}, obj.Attachment.EvpnVlans)
		assert.EqualValues(t, 5, *obj.Attachment.TrafficTag)
	})
}

func TestUnitNsxt_nsxtPolicyPortProfileSetInStruct(t *testing.T) {
	portSchema := segmentPortTestSchema()

	t.Run("discovery profile: no old, no new is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{})
		val, err := nsxtPolicyPortDiscoveryProfileSetInStruct(d)
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("discovery profile: new profile configured", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{
			"discovery_profile": []interface{}{
				map[string]interface{}{
					"ip_discovery_profile_path":  "/infra/ip-discovery-profiles/p1",
					"mac_discovery_profile_path": "",
					"binding_map_path":           "",
					"revision":                   0,
				},
			},
		})
		val, err := nsxtPolicyPortDiscoveryProfileSetInStruct(d)
		require.NoError(t, err)
		assert.NotNil(t, val)
	})

	t.Run("qos profile: new profile configured", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{
			"qos_profile": []interface{}{
				map[string]interface{}{
					"qos_profile_path": "/infra/qos-profiles/p1",
					"binding_map_path": "",
					"revision":         0,
				},
			},
		})
		val, err := nsxtPolicyPortQosProfileSetInStruct(d)
		require.NoError(t, err)
		assert.NotNil(t, val)
	})

	t.Run("qos profile: no old, no new is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{})
		val, err := nsxtPolicyPortQosProfileSetInStruct(d)
		require.NoError(t, err)
		assert.Nil(t, val)
	})

	t.Run("security profile: new profile configured", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{
			"security_profile": []interface{}{
				map[string]interface{}{
					"spoofguard_profile_path": "/infra/spoofguard-profiles/p1",
					"security_profile_path":   "/infra/segment-security-profiles/p1",
					"binding_map_path":        "",
					"revision":                0,
				},
			},
		})
		val, err := nsxtPolicyPortSecurityProfileSetInStruct(d)
		require.NoError(t, err)
		assert.NotNil(t, val)
	})

	t.Run("security profile: no old, no new is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{})
		val, err := nsxtPolicyPortSecurityProfileSetInStruct(d)
		require.NoError(t, err)
		assert.Nil(t, val)
	})
}

func TestUnitNsxt_policySegmentPortResourceToInfraStruct(t *testing.T) {
	portSchema := segmentPortTestSchema()

	t.Run("infra segment port builds an Infra struct", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{
			"display_name": "port-1",
			"segment_path": "/infra/segments/seg-1",
		})
		obj, err := policySegmentPortResourceToInfraStruct("port-1", d, false)
		require.NoError(t, err)
		assert.NotNil(t, obj.Children)
	})

	t.Run("tier-1 segment port builds a nested Infra struct", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{
			"display_name": "port-1",
			"segment_path": "/infra/tier-1s/gw-1/segments/seg-1",
		})
		obj, err := policySegmentPortResourceToInfraStruct("port-1", d, false)
		require.NoError(t, err)
		assert.NotNil(t, obj.Children)
	})

	t.Run("isDestroy marks the child for delete", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{
			"display_name": "port-1",
			"segment_path": "/infra/segments/seg-1",
		})
		obj, err := policySegmentPortResourceToInfraStructWithTags("port-1", d, true, nil)
		require.NoError(t, err)
		assert.NotNil(t, obj.Children)
	})
}
