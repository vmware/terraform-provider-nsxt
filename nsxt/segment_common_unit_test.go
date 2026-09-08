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

	segmentsapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	segmentsprofilesapi "github.com/vmware/terraform-provider-nsxt/api/infra/segments"
	tier1sapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	segmentmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	segmentprofilemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/segments"
	tier1segmentmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s"
)

func setupSegmentExistsMocks(t *testing.T) (*segmentmocks.MockSegmentsClient, *tier1segmentmocks.MockSegmentsClient) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSegSDK := segmentmocks.NewMockSegmentsClient(ctrl)
	origSeg := cliSegmentsClient
	t.Cleanup(func() { cliSegmentsClient = origSeg })
	cliSegmentsClient = func(_ utl.SessionContext, _ client.Connector) *segmentsapi.SegmentClientContext {
		return &segmentsapi.SegmentClientContext{Client: mockSegSDK, ClientType: utl.Local}
	}

	mockT1SegSDK := tier1segmentmocks.NewMockSegmentsClient(ctrl)
	origT1Seg := cliTier1SegmentsClient
	t.Cleanup(func() { cliTier1SegmentsClient = origT1Seg })
	cliTier1SegmentsClient = func(_ utl.SessionContext, _ client.Connector) *tier1sapi.SegmentClientContext {
		return &tier1sapi.SegmentClientContext{Client: mockT1SegSDK, ClientType: utl.Local}
	}

	return mockSegSDK, mockT1SegSDK
}

func setupSegmentProfileReadMocks(t *testing.T) (*segmentprofilemocks.MockSegmentDiscoveryProfileBindingMapsClient, *segmentprofilemocks.MockSegmentQosProfileBindingMapsClient, *segmentprofilemocks.MockSegmentSecurityProfileBindingMapsClient) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockDiscoverySDK := segmentprofilemocks.NewMockSegmentDiscoveryProfileBindingMapsClient(ctrl)
	origDiscovery := cliSegmentDiscoveryProfileBindingMapsClient
	t.Cleanup(func() { cliSegmentDiscoveryProfileBindingMapsClient = origDiscovery })
	cliSegmentDiscoveryProfileBindingMapsClient = func(_ utl.SessionContext, _ client.Connector) *segmentsprofilesapi.SegmentDiscoveryProfileBindingMapClientContext {
		return &segmentsprofilesapi.SegmentDiscoveryProfileBindingMapClientContext{Client: mockDiscoverySDK, ClientType: utl.Local}
	}

	mockQosSDK := segmentprofilemocks.NewMockSegmentQosProfileBindingMapsClient(ctrl)
	origQos := cliSegmentQosProfileBindingMapsClient
	t.Cleanup(func() { cliSegmentQosProfileBindingMapsClient = origQos })
	cliSegmentQosProfileBindingMapsClient = func(_ utl.SessionContext, _ client.Connector) *segmentsprofilesapi.SegmentQosProfileBindingMapClientContext {
		return &segmentsprofilesapi.SegmentQosProfileBindingMapClientContext{Client: mockQosSDK, ClientType: utl.Local}
	}

	mockSecuritySDK := segmentprofilemocks.NewMockSegmentSecurityProfileBindingMapsClient(ctrl)
	origSecurity := cliSegmentSecurityProfileBindingMapsClient
	t.Cleanup(func() { cliSegmentSecurityProfileBindingMapsClient = origSecurity })
	cliSegmentSecurityProfileBindingMapsClient = func(_ utl.SessionContext, _ client.Connector) *segmentsprofilesapi.SegmentSecurityProfileBindingMapClientContext {
		return &segmentsprofilesapi.SegmentSecurityProfileBindingMapClientContext{Client: mockSecuritySDK, ClientType: utl.Local}
	}

	return mockDiscoverySDK, mockQosSDK, mockSecuritySDK
}

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

	t.Run("overlay segment with subnets, advanced_config, l2_extension and profiles builds an Infra struct", func(t *testing.T) {
		segSchema := getPolicyCommonSegmentSchema(false, false)
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{
			"display_name":        "seg-1",
			"description":         "a segment",
			"domain_name":         "example.com",
			"transport_zone_path": "/infra/transport-zones/tz-1",
			"connectivity_path":   "/infra/tier-1s/gw-1",
			"overlay_id":          100,
			"replication_mode":    model.Segment_REPLICATION_MODE_MTEP,
			"subnet": []interface{}{
				map[string]interface{}{
					"cidr":        "10.0.0.1/24",
					"network":     "10.0.0.0/24",
					"dhcp_ranges": []interface{}{"10.0.0.100-10.0.0.200"},
					"dhcp_v4_config": []interface{}{
						map[string]interface{}{
							"server_address": "10.0.0.2/24",
							"lease_time":     3600,
						},
					},
				},
			},
			"advanced_config": []interface{}{
				map[string]interface{}{
					"connectivity":          "ON",
					"hybrid":                true,
					"local_egress":          false,
					"multicast":             true,
					"address_pool_path":     "/infra/ip-pools/p1",
					"uplink_teaming_policy": "policy-1",
					"urpf_mode":             "STRICT",
				},
			},
			"l2_extension": []interface{}{
				map[string]interface{}{
					"l2vpn_paths": []interface{}{"/infra/tier-0s/gw-1/l2vpn-services/svc-1"},
					"tunnel_id":   5,
				},
			},
			"discovery_profile": []interface{}{
				map[string]interface{}{
					"ip_discovery_profile_path": "/infra/ip-discovery-profiles/p1",
				},
			},
			"bridge_config": []interface{}{
				map[string]interface{}{
					"profile_path": "/infra/bridge-profiles/bp-1",
				},
			},
		})
		obj, err := policySegmentResourceToInfraStruct(getSessionContext(d, m), "seg-1", d, m, false, false)
		require.NoError(t, err)
		require.NotNil(t, obj.Children)
		require.Len(t, obj.Children, 1)
	})

	t.Run("fixed segment builds an Infra struct via the gateway child path", func(t *testing.T) {
		segSchema := getPolicyCommonSegmentSchema(false, true)
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{
			"display_name":      "seg-1",
			"connectivity_path": "/infra/tier-1s/gw-1",
			"replication_mode":  model.Segment_REPLICATION_MODE_MTEP,
		})
		obj, err := policySegmentResourceToInfraStruct(getSessionContext(d, m), "seg-1", d, m, false, true)
		require.NoError(t, err)
		require.Len(t, obj.Children, 1)
	})

	t.Run("multitenancy segment cannot change transport_zone_path", func(t *testing.T) {
		segSchema := getPolicyCommonSegmentSchema(false, false)
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{
			"display_name":        "seg-1",
			"connectivity_path":   "/infra/tier-1s/gw-1",
			"transport_zone_path": "/infra/transport-zones/tz-1",
			"replication_mode":    model.Segment_REPLICATION_MODE_MTEP,
		})
		d.Set("transport_zone_path", "/infra/transport-zones/tz-2")
		mtM := newGoMockProviderClient()
		ctx := utl.SessionContext{ClientType: utl.Multitenancy, ProjectID: "proj-1"}
		_, err := policySegmentResourceToInfraStruct(ctx, "seg-1", d, mtM, false, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be specified for project based segments")
	})
}

func TestMockNsxt_nsxtPolicySegmentProfileReads(t *testing.T) {
	segSchema := getPolicyCommonSegmentSchema(false, false)

	t.Run("discovery profile read populates schema when a result is found", func(t *testing.T) {
		mockDiscoverySDK, _, _ := setupSegmentProfileReadMocks(t)
		path := "/infra/ip-discovery-profiles/p1"
		mockDiscoverySDK.EXPECT().List("seg-1", nil, nil, nil, nil, nil, nil, nil).Return(
			model.SegmentDiscoveryProfileBindingMapListResult{Results: []model.SegmentDiscoveryProfileBindingMap{{IpDiscoveryProfilePath: &path}}}, nil,
		)

		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{})
		d.SetId("seg-1")
		err := nsxtPolicySegmentDiscoveryProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		list := d.Get("discovery_profile").([]interface{})
		require.Len(t, list, 1)
	})

	t.Run("discovery profile read propagates List error", func(t *testing.T) {
		mockDiscoverySDK, _, _ := setupSegmentProfileReadMocks(t)
		mockDiscoverySDK.EXPECT().List("seg-1", nil, nil, nil, nil, nil, nil, nil).Return(
			model.SegmentDiscoveryProfileBindingMapListResult{}, vapiErrors.InternalServerError{},
		)

		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{})
		d.SetId("seg-1")
		err := nsxtPolicySegmentDiscoveryProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("qos profile read populates schema only when path is non-empty", func(t *testing.T) {
		mockDiscoverySDK, mockQosSDK, mockSecuritySDK := setupSegmentProfileReadMocks(t)
		mockDiscoverySDK.EXPECT().List("seg-1", nil, nil, nil, nil, nil, nil, nil).Return(
			model.SegmentDiscoveryProfileBindingMapListResult{}, nil,
		)
		empty := ""
		path := "/infra/qos-profiles/p1"
		mockQosSDK.EXPECT().List("seg-1", nil, nil, nil, nil, nil).Return(
			model.SegmentQosProfileBindingMapListResult{Results: []model.SegmentQosProfileBindingMap{{QosProfilePath: &empty}, {QosProfilePath: &path}}}, nil,
		)
		mockSecuritySDK.EXPECT().List("seg-1", nil, nil, nil, nil, nil).Return(
			model.SegmentSecurityProfileBindingMapListResult{}, nil,
		)

		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{})
		d.SetId("seg-1")
		err := nsxtPolicySegmentProfilesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		list := d.Get("qos_profile").([]interface{})
		require.Len(t, list, 1)
	})

	t.Run("security profile read populates schema when a result is found", func(t *testing.T) {
		mockDiscoverySDK, mockQosSDK, mockSecuritySDK := setupSegmentProfileReadMocks(t)
		mockDiscoverySDK.EXPECT().List("seg-1", nil, nil, nil, nil, nil, nil, nil).Return(
			model.SegmentDiscoveryProfileBindingMapListResult{}, nil,
		)
		mockQosSDK.EXPECT().List("seg-1", nil, nil, nil, nil, nil).Return(
			model.SegmentQosProfileBindingMapListResult{}, nil,
		)
		secPath := "/infra/segment-security-profiles/p1"
		mockSecuritySDK.EXPECT().List("seg-1", nil, nil, nil, nil, nil).Return(
			model.SegmentSecurityProfileBindingMapListResult{Results: []model.SegmentSecurityProfileBindingMap{{SegmentSecurityProfilePath: &secPath}}}, nil,
		)

		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{})
		d.SetId("seg-1")
		err := nsxtPolicySegmentProfilesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		list := d.Get("security_profile").([]interface{})
		require.Len(t, list, 1)
	})

	t.Run("security profile read propagates List error", func(t *testing.T) {
		mockDiscoverySDK, mockQosSDK, mockSecuritySDK := setupSegmentProfileReadMocks(t)
		mockDiscoverySDK.EXPECT().List("seg-1", nil, nil, nil, nil, nil, nil, nil).Return(
			model.SegmentDiscoveryProfileBindingMapListResult{}, nil,
		)
		mockQosSDK.EXPECT().List("seg-1", nil, nil, nil, nil, nil).Return(
			model.SegmentQosProfileBindingMapListResult{}, nil,
		)
		mockSecuritySDK.EXPECT().List("seg-1", nil, nil, nil, nil, nil).Return(
			model.SegmentSecurityProfileBindingMapListResult{}, vapiErrors.InternalServerError{},
		)

		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{})
		d.SetId("seg-1")
		err := nsxtPolicySegmentProfilesRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockNsxt_nsxtPolicySegmentRead(t *testing.T) {
	segSchema := getPolicyCommonSegmentSchema(false, false)

	t.Run("populates advanced_config, l2_extension, subnet and profile blocks", func(t *testing.T) {
		mockSegSDK, _ := setupSegmentExistsMocks(t)
		mockDiscoverySDK, mockQosSDK, mockSecuritySDK := setupSegmentProfileReadMocks(t)

		connPath := "/infra/tier-1s/gw-1"
		revision := int64(1)
		urpf := "STRICT"
		egress := false
		hybrid := true
		multicast := true
		gwAddr := "10.0.0.1/24"
		network := "10.0.0.0/24"
		mockSegSDK.EXPECT().Get("seg-1").Return(model.Segment{
			ConnectivityPath: &connPath,
			Revision:         &revision,
			AdvancedConfig: &model.SegmentAdvancedConfig{
				UrpfMode:    &urpf,
				LocalEgress: &egress,
				Hybrid:      &hybrid,
				Multicast:   &multicast,
			},
			L2Extension: &model.L2Extension{TunnelId: &revision},
			Subnets: []model.SegmentSubnet{
				{GatewayAddress: &gwAddr, Network: &network},
			},
		}, nil)
		mockDiscoverySDK.EXPECT().List("seg-1", nil, nil, nil, nil, nil, nil, nil).Return(model.SegmentDiscoveryProfileBindingMapListResult{}, nil)
		mockQosSDK.EXPECT().List("seg-1", nil, nil, nil, nil, nil).Return(model.SegmentQosProfileBindingMapListResult{}, nil)
		mockSecuritySDK.EXPECT().List("seg-1", nil, nil, nil, nil, nil).Return(model.SegmentSecurityProfileBindingMapListResult{}, nil)

		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{})
		d.SetId("seg-1")
		err := nsxtPolicySegmentRead(d, newGoMockProviderClient(), false, false)
		require.NoError(t, err)
		assert.Equal(t, connPath, d.Get("connectivity_path"))
		advList := d.Get("advanced_config").([]interface{})
		require.Len(t, advList, 1)
		l2List := d.Get("l2_extension").([]interface{})
		require.Len(t, l2List, 1)
		subnetList := d.Get("subnet").([]interface{})
		require.Len(t, subnetList, 1)
	})

	t.Run("empty ID is an error", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{})
		err := nsxtPolicySegmentRead(d, newGoMockProviderClient(), false, false)
		require.Error(t, err)
	})

	t.Run("Get error propagates", func(t *testing.T) {
		mockSegSDK, _ := setupSegmentExistsMocks(t)
		mockSegSDK.EXPECT().Get("seg-1").Return(model.Segment{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, segSchema, map[string]interface{}{})
		d.SetId("seg-1")
		err := nsxtPolicySegmentRead(d, newGoMockProviderClient(), false, false)
		require.Error(t, err)
	})
}

func TestMockNsxt_resourceNsxtPolicySegmentExists(t *testing.T) {
	ctx := utl.SessionContext{ClientType: utl.Local}

	t.Run("non-fixed segment: Get succeeds means it exists", func(t *testing.T) {
		mockSegSDK, _ := setupSegmentExistsMocks(t)
		mockSegSDK.EXPECT().Get("seg-1").Return(model.Segment{}, nil)

		exists, err := resourceNsxtPolicySegmentExists(ctx, "", false)(ctx, "seg-1", nil)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("non-fixed segment: NotFound means it does not exist", func(t *testing.T) {
		mockSegSDK, _ := setupSegmentExistsMocks(t)
		mockSegSDK.EXPECT().Get("seg-1").Return(model.Segment{}, vapiErrors.NotFound{})

		exists, err := resourceNsxtPolicySegmentExists(ctx, "", false)(ctx, "seg-1", nil)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("non-fixed segment: other errors propagate", func(t *testing.T) {
		mockSegSDK, _ := setupSegmentExistsMocks(t)
		mockSegSDK.EXPECT().Get("seg-1").Return(model.Segment{}, vapiErrors.InternalServerError{})

		_, err := resourceNsxtPolicySegmentExists(ctx, "", false)(ctx, "seg-1", nil)
		require.Error(t, err)
	})

	t.Run("fixed tier-1 segment: Get succeeds means it exists", func(t *testing.T) {
		_, mockT1SegSDK := setupSegmentExistsMocks(t)
		mockT1SegSDK.EXPECT().Get("gw-1", "seg-1").Return(model.Segment{}, nil)

		exists, err := resourceNsxtPolicySegmentExists(ctx, "/infra/tier-1s/gw-1", true)(ctx, "seg-1", nil)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("fixed tier-0 segment path is rejected", func(t *testing.T) {
		setupSegmentExistsMocks(t)
		_, err := resourceNsxtPolicySegmentExists(ctx, "/infra/tier-0s/gw-1", true)(ctx, "seg-1", nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Tier-0")
	})

	t.Run("fixed segment with invalid gwPath is rejected", func(t *testing.T) {
		setupSegmentExistsMocks(t)
		_, err := resourceNsxtPolicySegmentExists(ctx, "", true)(ctx, "seg-1", nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not a valid gateway path")
	})
}
