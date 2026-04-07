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

	tier0sapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t0mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
)

var (
	routeMapID          = "route-map-1"
	routeMapDisplayName = "Test Route Map"
	routeMapDescription = "Test route map"
	routeMapRevision    = int64(1)
	routeMapGwPath      = "/infra/tier-0s/t0-rm-gw-1"
	routeMapGwID        = "t0-rm-gw-1"
	routeMapPath        = "/infra/tier-0s/t0-rm-gw-1/route-maps/route-map-1"
	routeMapEntryAction = nsxModel.RouteMapEntry_ACTION_PERMIT
)

func routeMapAPIResponse() nsxModel.Tier0RouteMap {
	return nsxModel.Tier0RouteMap{
		Id:          &routeMapID,
		DisplayName: &routeMapDisplayName,
		Description: &routeMapDescription,
		Revision:    &routeMapRevision,
		Path:        &routeMapPath,
		Entries: []nsxModel.RouteMapEntry{
			{Action: &routeMapEntryAction},
		},
	}
}

func minimalRouteMapData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": routeMapDisplayName,
		"description":  routeMapDescription,
		"nsx_id":       routeMapID,
		"gateway_path": routeMapGwPath,
		"entry": []interface{}{
			map[string]interface{}{
				"action":               routeMapEntryAction,
				"community_list_match": []interface{}{},
				"prefix_list_matches":  []interface{}{},
				"set":                  []interface{}{},
			},
		},
	}
}

func setupRouteMapMock(t *testing.T, ctrl *gomock.Controller) (*t0mocks.MockRouteMapsClient, func()) {
	mockSDK := t0mocks.NewMockRouteMapsClient(ctrl)
	mockWrapper := &tier0sapi.Tier0RouteMapClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliRouteMapsClient
	cliRouteMapsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier0sapi.Tier0RouteMapClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliRouteMapsClient = original }
}

func TestMockResourceNsxtPolicyGatewayRouteMapCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRouteMapMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(routeMapGwID, routeMapID).Return(nsxModel.Tier0RouteMap{}, notFoundErr),
			mockSDK.EXPECT().Patch(routeMapGwID, routeMapID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(routeMapGwID, routeMapID).Return(routeMapAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteMapData())

		err := resourceNsxtPolicyGatewayRouteMapCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, routeMapID, d.Id())
		assert.Equal(t, routeMapDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(routeMapGwID, routeMapID).Return(routeMapAPIResponse(), nil)

		res := resourceNsxtPolicyGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteMapData())

		err := resourceNsxtPolicyGatewayRouteMapCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Create fails for non-Tier0 gateway path", func(t *testing.T) {
		data := minimalRouteMapData()
		data["gateway_path"] = "/infra/tier-1s/t1-gw-1"

		res := resourceNsxtPolicyGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyGatewayRouteMapCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Tier0")
	})
}

func TestMockResourceNsxtPolicyGatewayRouteMapRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRouteMapMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(routeMapGwID, routeMapID).Return(routeMapAPIResponse(), nil)

		res := resourceNsxtPolicyGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteMapData())
		d.SetId(routeMapID)

		err := resourceNsxtPolicyGatewayRouteMapRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, routeMapDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(routeMapGwID, routeMapID).Return(nsxModel.Tier0RouteMap{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteMapData())
		d.SetId(routeMapID)

		err := resourceNsxtPolicyGatewayRouteMapRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteMapData())

		err := resourceNsxtPolicyGatewayRouteMapRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayRouteMapUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRouteMapMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(routeMapGwID, routeMapID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(routeMapGwID, routeMapID).Return(routeMapAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteMapData())
		d.SetId(routeMapID)

		err := resourceNsxtPolicyGatewayRouteMapUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteMapData())

		err := resourceNsxtPolicyGatewayRouteMapUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayRouteMapDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRouteMapMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(routeMapGwID, routeMapID).Return(nil)

		res := resourceNsxtPolicyGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteMapData())
		d.SetId(routeMapID)

		err := resourceNsxtPolicyGatewayRouteMapDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteMapData())

		err := resourceNsxtPolicyGatewayRouteMapDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
