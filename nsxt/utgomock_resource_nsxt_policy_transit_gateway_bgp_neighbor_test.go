//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	tgwbgp "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways/bgp"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	tgwbgpmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways/bgp"
)

var (
	tgwBgpNbrID          = "bgp-neighbor-1"
	tgwBgpNbrDisplayName = "Test TGW BGP Neighbor"
	tgwBgpNbrDescription = "Test transit gateway BGP neighbor"
	tgwBgpNbrRevision    = int64(1)
	tgwBgpNbrParentPath  = "/orgs/default/projects/project1/transit-gateways/tgw1"
	tgwBgpNbrOrgID       = "default"
	tgwBgpNbrProjectID   = "project1"
	tgwBgpNbrTGWID       = "tgw1"
	tgwBgpNbrAddress     = "192.168.20.1"
	tgwBgpNbrRemoteAsNum = "65000"
	tgwBgpNbrPath        = "/orgs/default/projects/project1/transit-gateways/tgw1/bgp/neighbors/bgp-neighbor-1"
	tgwBgpNbrAttachment  = "/orgs/default/projects/project1/transit-gateways/tgw1/attachments/att1"
)

func tgwBgpNeighborAPIResponse() nsxModel.BgpNeighborConfig {
	holdDown := int64(180)
	keepAlive := int64(60)
	maxHop := int64(1)
	allowAsIn := false
	grMode := nsxModel.BgpNeighborConfig_GRACEFUL_RESTART_MODE_HELPER_ONLY
	return nsxModel.BgpNeighborConfig{
		Id:                  &tgwBgpNbrID,
		DisplayName:         &tgwBgpNbrDisplayName,
		Description:         &tgwBgpNbrDescription,
		Revision:            &tgwBgpNbrRevision,
		Path:                &tgwBgpNbrPath,
		NeighborAddress:     &tgwBgpNbrAddress,
		RemoteAsNum:         &tgwBgpNbrRemoteAsNum,
		HoldDownTime:        &holdDown,
		KeepAliveTime:       &keepAlive,
		MaximumHopLimit:     &maxHop,
		AllowAsIn:           &allowAsIn,
		GracefulRestartMode: &grMode,
		SourceAttachment:    []string{tgwBgpNbrAttachment},
	}
}

func minimalTGWBgpNeighborData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":          tgwBgpNbrDisplayName,
		"description":           tgwBgpNbrDescription,
		"nsx_id":                tgwBgpNbrID,
		"parent_path":           tgwBgpNbrParentPath,
		"neighbor_address":      tgwBgpNbrAddress,
		"remote_as_num":         tgwBgpNbrRemoteAsNum,
		"source_attachment":     []interface{}{tgwBgpNbrAttachment},
		"allow_as_in":           false,
		"graceful_restart_mode": nsxModel.BgpNeighborConfig_GRACEFUL_RESTART_MODE_HELPER_ONLY,
		"hold_down_time":        180,
		"keep_alive_time":       60,
		"maximum_hop_limit":     1,
	}
}

func setupTGWBgpNeighborMock(t *testing.T, ctrl *gomock.Controller) (*tgwbgpmocks.MockNeighborsClient, func()) {
	mockSDK := tgwbgpmocks.NewMockNeighborsClient(ctrl)
	mockWrapper := &tgwbgp.TGWBgpNeighborClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwBgpNbrProjectID,
	}

	original := cliTGWBgpNeighborsClient
	cliTGWBgpNeighborsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tgwbgp.TGWBgpNeighborClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTGWBgpNeighborsClient = original }
}

func TestMockResourceNsxtPolicyTransitGatewayBgpNeighborCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWBgpNeighborMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwBgpNbrOrgID, tgwBgpNbrProjectID, tgwBgpNbrTGWID, tgwBgpNbrID).Return(nsxModel.BgpNeighborConfig{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwBgpNbrOrgID, tgwBgpNbrProjectID, tgwBgpNbrTGWID, tgwBgpNbrID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwBgpNbrOrgID, tgwBgpNbrProjectID, tgwBgpNbrTGWID, tgwBgpNbrID).Return(tgwBgpNeighborAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBgpNeighborData())

		err := resourceNsxtPolicyTransitGatewayBgpNeighborCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwBgpNbrID, d.Id())
		assert.Equal(t, tgwBgpNbrDisplayName, d.Get("display_name"))
		assert.Equal(t, tgwBgpNbrAddress, d.Get("neighbor_address"))
		assert.Equal(t, tgwBgpNbrRemoteAsNum, d.Get("remote_as_num"))
	})

	t.Run("Create fails when parent_path is missing", func(t *testing.T) {
		data := minimalTGWBgpNeighborData()
		data["parent_path"] = ""

		res := resourceNsxtPolicyTransitGatewayBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayBgpNeighborCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayBgpNeighborRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWBgpNeighborMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwBgpNbrOrgID, tgwBgpNbrProjectID, tgwBgpNbrTGWID, tgwBgpNbrID).Return(tgwBgpNeighborAPIResponse(), nil)

		res := resourceNsxtPolicyTransitGatewayBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBgpNeighborData())
		d.SetId(tgwBgpNbrID)

		err := resourceNsxtPolicyTransitGatewayBgpNeighborRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwBgpNbrDisplayName, d.Get("display_name"))
		assert.Equal(t, tgwBgpNbrAddress, d.Get("neighbor_address"))
		assert.Equal(t, tgwBgpNbrRemoteAsNum, d.Get("remote_as_num"))
		assert.Equal(t, int(180), d.Get("hold_down_time"))
		assert.Equal(t, int(60), d.Get("keep_alive_time"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwBgpNbrOrgID, tgwBgpNbrProjectID, tgwBgpNbrTGWID, tgwBgpNbrID).Return(nsxModel.BgpNeighborConfig{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTransitGatewayBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBgpNeighborData())
		d.SetId(tgwBgpNbrID)

		err := resourceNsxtPolicyTransitGatewayBgpNeighborRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBgpNeighborData())

		err := resourceNsxtPolicyTransitGatewayBgpNeighborRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayBgpNeighborUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWBgpNeighborMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(tgwBgpNbrOrgID, tgwBgpNbrProjectID, tgwBgpNbrTGWID, tgwBgpNbrID, gomock.Any()).Return(tgwBgpNeighborAPIResponse(), nil),
			mockSDK.EXPECT().Get(tgwBgpNbrOrgID, tgwBgpNbrProjectID, tgwBgpNbrTGWID, tgwBgpNbrID).Return(tgwBgpNeighborAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBgpNeighborData())
		d.SetId(tgwBgpNbrID)

		err := resourceNsxtPolicyTransitGatewayBgpNeighborUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwBgpNbrID, d.Id())
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBgpNeighborData())

		err := resourceNsxtPolicyTransitGatewayBgpNeighborUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayBgpNeighborDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWBgpNeighborMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(tgwBgpNbrOrgID, tgwBgpNbrProjectID, tgwBgpNbrTGWID, tgwBgpNbrID).Return(nil)

		res := resourceNsxtPolicyTransitGatewayBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBgpNeighborData())
		d.SetId(tgwBgpNbrID)

		err := resourceNsxtPolicyTransitGatewayBgpNeighborDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWBgpNeighborData())

		err := resourceNsxtPolicyTransitGatewayBgpNeighborDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayBgpNeighborWithRouteFiltering(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWBgpNeighborMock(t, ctrl)
	defer restore()

	t.Run("Create with route_filtering succeeds", func(t *testing.T) {
		addrFamilyIPv4 := nsxModel.BgpRouteFiltering_ADDRESS_FAMILY_IPV4
		filterEnabled := true
		holdDown := int64(180)
		keepAlive := int64(60)
		maxHop := int64(1)
		allowAsIn := false
		grMode := nsxModel.BgpNeighborConfig_GRACEFUL_RESTART_MODE_HELPER_ONLY

		response := nsxModel.BgpNeighborConfig{
			Id:                  &tgwBgpNbrID,
			DisplayName:         &tgwBgpNbrDisplayName,
			Description:         &tgwBgpNbrDescription,
			Revision:            &tgwBgpNbrRevision,
			Path:                &tgwBgpNbrPath,
			NeighborAddress:     &tgwBgpNbrAddress,
			RemoteAsNum:         &tgwBgpNbrRemoteAsNum,
			HoldDownTime:        &holdDown,
			KeepAliveTime:       &keepAlive,
			MaximumHopLimit:     &maxHop,
			AllowAsIn:           &allowAsIn,
			GracefulRestartMode: &grMode,
			SourceAttachment:    []string{tgwBgpNbrAttachment},
			RouteFiltering: []nsxModel.BgpRouteFiltering{
				{AddressFamily: &addrFamilyIPv4, Enabled: &filterEnabled},
			},
		}

		data := minimalTGWBgpNeighborData()
		data["route_filtering"] = []interface{}{
			map[string]interface{}{
				"address_family":   nsxModel.BgpRouteFiltering_ADDRESS_FAMILY_IPV4,
				"enabled":          true,
				"in_route_filter":  "",
				"out_route_filter": "",
				"maximum_routes":   0,
			},
		}

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwBgpNbrOrgID, tgwBgpNbrProjectID, tgwBgpNbrTGWID, tgwBgpNbrID).Return(nsxModel.BgpNeighborConfig{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwBgpNbrOrgID, tgwBgpNbrProjectID, tgwBgpNbrTGWID, tgwBgpNbrID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwBgpNbrOrgID, tgwBgpNbrProjectID, tgwBgpNbrTGWID, tgwBgpNbrID).Return(response, nil),
		)

		res := resourceNsxtPolicyTransitGatewayBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayBgpNeighborCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwBgpNbrID, d.Id())
	})
}
