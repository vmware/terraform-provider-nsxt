// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/nsx/TransportNodesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/TransportNodesClient.go TransportNodesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	mpmodel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

var (
	rtepEdgeID        = "edge-transport-node-uuid-1"
	rtepHostSwitch    = "nvds-1"
	rtepVlanID        = int64(100)
	rtepTeamingPolicy = "named-teaming-1"
)

// dhcpIpAssignmentStructValue builds the *data.StructValue for AssignedByDhcp
// so that RTEP Read tests can round-trip through setIPAssignmentInSchema.
func dhcpIpAssignmentStructValue() *data.StructValue {
	converter := bindings.NewTypeConverter()
	dhcp := mpmodel.AssignedByDhcp{
		ResourceType: mpmodel.IpAssignmentSpec_RESOURCE_TYPE_ASSIGNEDBYDHCP,
	}
	dataValue, errs := converter.ConvertToVapi(dhcp, mpmodel.AssignedByDhcpBindingType())
	if errs != nil {
		panic(errs[0])
	}
	return dataValue.(*data.StructValue)
}

func transportNodeWithRTEP() mpmodel.TransportNode {
	return mpmodel.TransportNode{
		Id:                 &rtepEdgeID,
		NodeDeploymentInfo: edgeNodeStructValue(),
		RemoteTunnelEndpoint: &mpmodel.TransportNodeRemoteTunnelEndpointConfig{
			HostSwitchName:     &rtepHostSwitch,
			RtepVlan:           &rtepVlanID,
			NamedTeamingPolicy: &rtepTeamingPolicy,
			IpAssignmentSpec:   dhcpIpAssignmentStructValue(),
		},
	}
}

func transportNodeWithoutRTEP() mpmodel.TransportNode {
	return mpmodel.TransportNode{
		Id:                   &rtepEdgeID,
		NodeDeploymentInfo:   edgeNodeStructValue(),
		RemoteTunnelEndpoint: nil,
	}
}

func minimalRTEPData() map[string]interface{} {
	return map[string]interface{}{
		"edge_id":          rtepEdgeID,
		"host_switch_name": rtepHostSwitch,
		"rtep_vlan":        int(rtepVlanID),
		"ip_assignment": []interface{}{
			map[string]interface{}{
				"assigned_by_dhcp": true,
				"no_ipv4":          false,
				"static_ip":        []interface{}{},
				"static_ip_pool":   "",
				"static_ip_mac":    []interface{}{},
			},
		},
	}
}

func TestMockResourceNsxtEdgeTransportNodeRTEPCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		// Create calls Get, then Update; then Read calls Get again
		mockSDK.EXPECT().Get(rtepEdgeID).Return(transportNodeWithoutRTEP(), nil)
		mockSDK.EXPECT().Update(rtepEdgeID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(transportNodeWithRTEP(), nil)
		mockSDK.EXPECT().Get(rtepEdgeID).Return(transportNodeWithRTEP(), nil)

		res := resourceNsxtEdgeTransportNodeRTEP()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRTEPData())

		err := resourceNsxtEdgeTransportNodeRTEPCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, rtepEdgeID, d.Id())
	})

	t.Run("Create fails when RTEP already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(rtepEdgeID).Return(transportNodeWithRTEP(), nil)

		res := resourceNsxtEdgeTransportNodeRTEP()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRTEPData())

		err := resourceNsxtEdgeTransportNodeRTEPCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when Get API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(rtepEdgeID).Return(mpmodel.TransportNode{}, errors.New("get API error"))

		res := resourceNsxtEdgeTransportNodeRTEP()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRTEPData())

		err := resourceNsxtEdgeTransportNodeRTEPCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "get API error")
	})
}

func TestMockResourceNsxtEdgeTransportNodeRTEPRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(rtepEdgeID).Return(transportNodeWithRTEP(), nil)

		res := resourceNsxtEdgeTransportNodeRTEP()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_id": rtepEdgeID,
		})

		err := resourceNsxtEdgeTransportNodeRTEPRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, rtepHostSwitch, d.Get("host_switch_name"))
		assert.Equal(t, int(rtepVlanID), d.Get("rtep_vlan"))
	})

	t.Run("Read fails when Get API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(rtepEdgeID).Return(mpmodel.TransportNode{}, errors.New("read API error"))

		res := resourceNsxtEdgeTransportNodeRTEP()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_id": rtepEdgeID,
		})

		err := resourceNsxtEdgeTransportNodeRTEPRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read API error")
	})
}

func TestMockResourceNsxtEdgeTransportNodeRTEPUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(rtepEdgeID).Return(transportNodeWithRTEP(), nil)
		mockSDK.EXPECT().Update(rtepEdgeID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(transportNodeWithRTEP(), nil)
		mockSDK.EXPECT().Get(rtepEdgeID).Return(transportNodeWithRTEP(), nil)

		res := resourceNsxtEdgeTransportNodeRTEP()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRTEPData())
		d.SetId(rtepEdgeID)

		err := resourceNsxtEdgeTransportNodeRTEPUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when RTEP not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(rtepEdgeID).Return(transportNodeWithoutRTEP(), nil)

		res := resourceNsxtEdgeTransportNodeRTEP()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRTEPData())

		err := resourceNsxtEdgeTransportNodeRTEPUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update fails when Update API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(rtepEdgeID).Return(transportNodeWithRTEP(), nil)
		mockSDK.EXPECT().Update(rtepEdgeID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mpmodel.TransportNode{}, errors.New("update API error"))

		res := resourceNsxtEdgeTransportNodeRTEP()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRTEPData())
		d.SetId(rtepEdgeID)

		err := resourceNsxtEdgeTransportNodeRTEPUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "update API error")
	})
}

func TestMockResourceNsxtEdgeTransportNodeRTEPDelete(t *testing.T) {
	t.Run("Delete success (clears RTEP via Update)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(rtepEdgeID).Return(transportNodeWithRTEP(), nil)
		mockSDK.EXPECT().Update(rtepEdgeID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(transportNodeWithoutRTEP(), nil)

		res := resourceNsxtEdgeTransportNodeRTEP()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_id": rtepEdgeID,
		})

		err := resourceNsxtEdgeTransportNodeRTEPDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when Get API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(rtepEdgeID).Return(mpmodel.TransportNode{}, errors.New("get API error"))

		res := resourceNsxtEdgeTransportNodeRTEP()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_id": rtepEdgeID,
		})

		err := resourceNsxtEdgeTransportNodeRTEPDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "get API error")
	})
}
