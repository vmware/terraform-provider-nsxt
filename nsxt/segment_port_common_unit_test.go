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

	portsapi "github.com/vmware/terraform-provider-nsxt/api/infra/segments"
	t1portsapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/segments"
	t1profilesapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/segments/ports"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	portmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/segments"
	t1portmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s/segments"
	t1profilemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s/segments/ports"
)

func setupSegmentPortExistsMocks(t *testing.T) (*portmocks.MockPortsClient, *t1portmocks.MockPortsClient) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockPortSDK := portmocks.NewMockPortsClient(ctrl)
	origPorts := cliPortsClient
	t.Cleanup(func() { cliPortsClient = origPorts })
	cliPortsClient = func(_ utl.SessionContext, _ client.Connector) *portsapi.SegmentPortClientContext {
		return &portsapi.SegmentPortClientContext{Client: mockPortSDK, ClientType: utl.Local}
	}

	mockT1PortSDK := t1portmocks.NewMockPortsClient(ctrl)
	origT1Ports := cliT1SegmentPortsClient
	t.Cleanup(func() { cliT1SegmentPortsClient = origT1Ports })
	cliT1SegmentPortsClient = func(_ utl.SessionContext, _ client.Connector) *t1portsapi.SegmentPortClientContext {
		return &t1portsapi.SegmentPortClientContext{Client: mockT1PortSDK, ClientType: utl.Local}
	}

	return mockPortSDK, mockT1PortSDK
}

func setupT1PortProfileMocks(t *testing.T) (*t1profilemocks.MockPortDiscoveryProfileBindingMapsClient, *t1profilemocks.MockPortQosProfileBindingMapsClient, *t1profilemocks.MockPortSecurityProfileBindingMapsClient) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockDiscovery := t1profilemocks.NewMockPortDiscoveryProfileBindingMapsClient(ctrl)
	origDiscovery := cliT1PortDiscoveryProfileBindingMapsClient
	t.Cleanup(func() { cliT1PortDiscoveryProfileBindingMapsClient = origDiscovery })
	cliT1PortDiscoveryProfileBindingMapsClient = func(_ utl.SessionContext, _ client.Connector) *t1profilesapi.PortDiscoveryProfileBindingMapClientContext {
		return &t1profilesapi.PortDiscoveryProfileBindingMapClientContext{Client: mockDiscovery, ClientType: utl.Local}
	}

	mockQos := t1profilemocks.NewMockPortQosProfileBindingMapsClient(ctrl)
	origQos := cliT1PortQosProfileBindingMapsClient
	t.Cleanup(func() { cliT1PortQosProfileBindingMapsClient = origQos })
	cliT1PortQosProfileBindingMapsClient = func(_ utl.SessionContext, _ client.Connector) *t1profilesapi.PortQosProfileBindingMapClientContext {
		return &t1profilesapi.PortQosProfileBindingMapClientContext{Client: mockQos, ClientType: utl.Local}
	}

	mockSecurity := t1profilemocks.NewMockPortSecurityProfileBindingMapsClient(ctrl)
	origSecurity := cliT1PortSecurityProfileBindingMapsClient
	t.Cleanup(func() { cliT1PortSecurityProfileBindingMapsClient = origSecurity })
	cliT1PortSecurityProfileBindingMapsClient = func(_ utl.SessionContext, _ client.Connector) *t1profilesapi.PortSecurityProfileBindingMapClientContext {
		return &t1profilesapi.PortSecurityProfileBindingMapClientContext{Client: mockSecurity, ClientType: utl.Local}
	}

	return mockDiscovery, mockQos, mockSecurity
}

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

func TestMockNsxt_resourceNsxtPolicySegmentPortExists(t *testing.T) {
	portSchema := segmentPortTestSchema()
	ctx := utl.SessionContext{ClientType: utl.Local}

	t.Run("non-tier1 segment: Get succeeds means it exists", func(t *testing.T) {
		mockPortSDK, _ := setupSegmentPortExistsMocks(t)
		mockPortSDK.EXPECT().Get("seg-1", "port-1").Return(model.SegmentPort{}, nil)

		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{"segment_path": "/infra/segments/seg-1"})
		exists, err := resourceNsxtPolicySegmentPortExists(d)(ctx, "port-1", nil)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("non-tier1 segment: NotFound means it does not exist", func(t *testing.T) {
		mockPortSDK, _ := setupSegmentPortExistsMocks(t)
		mockPortSDK.EXPECT().Get("seg-1", "port-1").Return(model.SegmentPort{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{"segment_path": "/infra/segments/seg-1"})
		exists, err := resourceNsxtPolicySegmentPortExists(d)(ctx, "port-1", nil)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("non-tier1 segment: other errors propagate", func(t *testing.T) {
		mockPortSDK, _ := setupSegmentPortExistsMocks(t)
		mockPortSDK.EXPECT().Get("seg-1", "port-1").Return(model.SegmentPort{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{"segment_path": "/infra/segments/seg-1"})
		_, err := resourceNsxtPolicySegmentPortExists(d)(ctx, "port-1", nil)
		require.Error(t, err)
	})

	t.Run("tier-1 segment: Get succeeds means it exists", func(t *testing.T) {
		_, mockT1PortSDK := setupSegmentPortExistsMocks(t)
		mockT1PortSDK.EXPECT().Get("gw-1", "seg-1", "port-1").Return(model.SegmentPort{}, nil)

		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{"segment_path": "/infra/tier-1s/gw-1/segments/seg-1"})
		exists, err := resourceNsxtPolicySegmentPortExists(d)(ctx, "port-1", nil)
		require.NoError(t, err)
		assert.True(t, exists)
	})
}

func TestMockNsxt_tier1SegmentPortProfileReads(t *testing.T) {
	c := tier1SegmentPort{
		tier1GatewayId: "gw-1",
		ids:            &segmentPort{segmentId: "seg-1", portId: "port-1"},
	}
	portSchema := segmentPortTestSchema()

	t.Run("discovery profile read populates schema", func(t *testing.T) {
		mockDiscovery, _, _ := setupT1PortProfileMocks(t)
		ipPath := "/infra/ip-discovery-profiles/p1"
		mockDiscovery.EXPECT().List("gw-1", "seg-1", "port-1", nil, nil, nil, nil, nil, nil, nil).Return(
			model.PortDiscoveryProfileBindingMapListResult{Results: []model.PortDiscoveryProfileBindingMap{{IpDiscoveryProfilePath: &ipPath}}}, nil,
		)

		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{})
		err := c.nsxtPolicySegmentPortDiscoveryProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		got := d.Get("discovery_profile").([]interface{})
		require.Len(t, got, 1)
		assert.Equal(t, ipPath, got[0].(map[string]interface{})["ip_discovery_profile_path"])
	})

	t.Run("discovery profile read propagates API error", func(t *testing.T) {
		mockDiscovery, _, _ := setupT1PortProfileMocks(t)
		mockDiscovery.EXPECT().List("gw-1", "seg-1", "port-1", nil, nil, nil, nil, nil, nil, nil).Return(
			model.PortDiscoveryProfileBindingMapListResult{}, vapiErrors.InternalServerError{},
		)

		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{})
		err := c.nsxtPolicySegmentPortDiscoveryProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("qos profile read populates schema", func(t *testing.T) {
		_, mockQos, _ := setupT1PortProfileMocks(t)
		qosPath := "/infra/qos-profiles/p1"
		mockQos.EXPECT().List("gw-1", "seg-1", "port-1", nil, nil, nil, nil, nil).Return(
			model.PortQosProfileBindingMapListResult{Results: []model.PortQosProfileBindingMap{{QosProfilePath: &qosPath}}}, nil,
		)

		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{})
		err := c.nsxtPolicySegmentPortQosProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		got := d.Get("qos_profile").([]interface{})
		require.Len(t, got, 1)
		assert.Equal(t, qosPath, got[0].(map[string]interface{})["qos_profile_path"])
	})

	t.Run("qos profile read skips entries with no qos_profile_path", func(t *testing.T) {
		_, mockQos, _ := setupT1PortProfileMocks(t)
		empty := ""
		mockQos.EXPECT().List("gw-1", "seg-1", "port-1", nil, nil, nil, nil, nil).Return(
			model.PortQosProfileBindingMapListResult{Results: []model.PortQosProfileBindingMap{{QosProfilePath: &empty}}}, nil,
		)

		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{})
		err := c.nsxtPolicySegmentPortQosProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Get("qos_profile").([]interface{}))
	})

	t.Run("security profile read populates schema", func(t *testing.T) {
		_, _, mockSecurity := setupT1PortProfileMocks(t)
		secPath := "/infra/segment-security-profiles/p1"
		mockSecurity.EXPECT().List("gw-1", "seg-1", "port-1", nil, nil, nil, nil, nil).Return(
			model.PortSecurityProfileBindingMapListResult{Results: []model.PortSecurityProfileBindingMap{{SegmentSecurityProfilePath: &secPath}}}, nil,
		)

		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{})
		err := c.nsxtPolicyPortSegmentSecurityProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		got := d.Get("security_profile").([]interface{})
		require.Len(t, got, 1)
		assert.Equal(t, secPath, got[0].(map[string]interface{})["security_profile_path"])
	})

	t.Run("security profile read propagates API error", func(t *testing.T) {
		_, _, mockSecurity := setupT1PortProfileMocks(t)
		mockSecurity.EXPECT().List("gw-1", "seg-1", "port-1", nil, nil, nil, nil, nil).Return(
			model.PortSecurityProfileBindingMapListResult{}, vapiErrors.InternalServerError{},
		)

		d := schema.TestResourceDataRaw(t, portSchema, map[string]interface{}{})
		err := c.nsxtPolicyPortSegmentSecurityProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
