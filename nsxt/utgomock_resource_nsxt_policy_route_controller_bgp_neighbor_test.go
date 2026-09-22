//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run:
// mockgen -destination=mocks/infra/route_controllers/bgp/NeighborsClient.go -package=mocks -source=<sdk>/services/nsxt/infra/route_controllers/bgp/NeighborsClient.go NeighborsClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	rcbgpnbrmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers/bgp"
)

var (
	rcBgpNeighborID          = "nbr-1"
	rcBgpNeighborDisplayName = "nbr-fooname"
	rcBgpNeighborDescription = "bgp neighbor mock"
	rcBgpNeighborPath        = "/infra/route-controllers/rc-1/bgp/neighbors/nbr-1"
	rcBgpNeighborRevision    = int64(1)
	rcBgpNeighborAddress     = "192.168.1.1"
	rcBgpNeighborRemoteAsNum = "65001"
	rcBgpNeighborParentPath  = "/infra/route-controllers/rc-1/bgp"
)

func rcBgpNeighborAPIResponse() model.RouteControllerBgpNeighborConfig {
	enabled := true
	allowAsIn := false
	holdDownTime := int64(180)
	keepAliveTime := int64(60)
	maximumHopLimit := int64(1)
	gracefulRestartMode := model.BgpNeighborConfig_GRACEFUL_RESTART_MODE_HELPER_ONLY
	return model.RouteControllerBgpNeighborConfig{
		Id:                  &rcBgpNeighborID,
		DisplayName:         &rcBgpNeighborDisplayName,
		Description:         &rcBgpNeighborDescription,
		Path:                &rcBgpNeighborPath,
		Revision:            &rcBgpNeighborRevision,
		Enabled:             &enabled,
		AllowAsIn:           &allowAsIn,
		HoldDownTime:        &holdDownTime,
		KeepAliveTime:       &keepAliveTime,
		MaximumHopLimit:     &maximumHopLimit,
		GracefulRestartMode: &gracefulRestartMode,
		NeighborAddress:     &rcBgpNeighborAddress,
		RemoteAsNum:         &rcBgpNeighborRemoteAsNum,
	}
}

func minimalRCBgpNeighborData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":     rcBgpNeighborDisplayName,
		"description":      rcBgpNeighborDescription,
		"neighbor_address": rcBgpNeighborAddress,
		"remote_as_num":    rcBgpNeighborRemoteAsNum,
		"parent_path":      rcBgpNeighborParentPath,
	}
}

func setupRCBgpNeighborMocks(t *testing.T, ctrl *gomock.Controller) (
	*rcbgpnbrmocks.MockNeighborsClient,
	func(),
) {
	mockNbrSDK := rcbgpnbrmocks.NewMockNeighborsClient(ctrl)

	nbrWrapper := &cliinfra.RouteControllerBgpNeighborClientContext{
		Client:     mockNbrSDK,
		ClientType: utl.Local,
	}

	orig := cliRCBgpNeighborClient

	cliRCBgpNeighborClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.RouteControllerBgpNeighborClientContext {
		return nbrWrapper
	}

	return mockNbrSDK, func() {
		cliRCBgpNeighborClient = orig
	}
}

func TestMockResourceNsxtPolicyRouteControllerBgpNeighborCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNbrSDK, restore := setupRCBgpNeighborMocks(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		resp := rcBgpNeighborAPIResponse()
		gomock.InOrder(
			mockNbrSDK.EXPECT().Patch(routeControllerID, gomock.Any(), gomock.Any()).Return(nil),
			mockNbrSDK.EXPECT().Get(routeControllerID, gomock.Any()).Return(resp, nil),
		)

		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRCBgpNeighborData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockNbrSDK.EXPECT().Get(routeControllerID, rcBgpNeighborID).Return(rcBgpNeighborAPIResponse(), nil)

		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		data := minimalRCBgpNeighborData()
		data["nsx_id"] = rcBgpNeighborID
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when API returns error", func(t *testing.T) {
		mockNbrSDK.EXPECT().Patch(routeControllerID, gomock.Any(), gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRCBgpNeighborData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})

	t.Run("Create with valid route_filtering succeeds", func(t *testing.T) {
		resp := rcBgpNeighborAPIResponse()
		gomock.InOrder(
			mockNbrSDK.EXPECT().Patch(routeControllerID, gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ string, _ string, nbr model.RouteControllerBgpNeighborConfig) error {
					require.Len(t, nbr.RouteFiltering, 2)
					assert.Equal(t, "IPV4", *nbr.RouteFiltering[0].AddressFamily)
					assert.Nil(t, nbr.RouteFiltering[0].MaximumRoutes)
					assert.Equal(t, "L2VPN_EVPN", *nbr.RouteFiltering[1].AddressFamily)
					assert.NotNil(t, nbr.RouteFiltering[1].MaximumRoutes)
					assert.Equal(t, int64(500), *nbr.RouteFiltering[1].MaximumRoutes)
					return nil
				}),
			mockNbrSDK.EXPECT().Get(routeControllerID, gomock.Any()).Return(resp, nil),
		)

		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		data := minimalRCBgpNeighborData()
		data["route_filtering"] = []interface{}{
			map[string]interface{}{
				"address_family": "IPV4",
				"enabled":        true,
			},
			map[string]interface{}{
				"address_family": "L2VPN_EVPN",
				"enabled":        true,
				"maximum_routes": 500,
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})

	t.Run("Create fails when IPV4 route_filtering has maximum_routes", func(t *testing.T) {
		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		data := minimalRCBgpNeighborData()
		data["route_filtering"] = []interface{}{
			map[string]interface{}{
				"address_family": "IPV4",
				"enabled":        true,
				"maximum_routes": 500,
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "the property 'maximum_routes' is not supported for route_filtering configured with address family IPV4")
	})

	t.Run("Create fails when IPV4 route_filtering has in_route_filter", func(t *testing.T) {
		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		data := minimalRCBgpNeighborData()
		data["route_filtering"] = []interface{}{
			map[string]interface{}{
				"address_family":  "IPV4",
				"enabled":         true,
				"in_route_filter": "/infra/prefix-lists/test",
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "the property 'in_route_filter' is not supported for route_filtering configured with address family IPV4")
	})

	t.Run("Create fails when IPV4 route_filtering has out_route_filter", func(t *testing.T) {
		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		data := minimalRCBgpNeighborData()
		data["route_filtering"] = []interface{}{
			map[string]interface{}{
				"address_family":   "IPV4",
				"enabled":          true,
				"out_route_filter": "/infra/prefix-lists/test",
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "the property 'out_route_filter' is not supported for route_filtering configured with address family IPV4")
	})

	t.Run("Create fails when duplicate address_family configured", func(t *testing.T) {
		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		data := minimalRCBgpNeighborData()
		data["route_filtering"] = []interface{}{
			map[string]interface{}{
				"address_family": "IPV4",
				"enabled":        true,
			},
			map[string]interface{}{
				"address_family": "IPV4",
				"enabled":        true,
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate address_family 'IPV4' found in route_filtering")
	})
}

func TestMockResourceNsxtPolicyRouteControllerBgpNeighborRouteFilteringAddressFamilyValidation(t *testing.T) {
	res := resourceNsxtPolicyRouteControllerBgpNeighbor()
	elemSchema := res.Schema["route_filtering"].Elem.(*schema.Resource).Schema
	afSchema := elemSchema["address_family"]

	for _, validFamily := range []string{"IPV4", "L2VPN_EVPN"} {
		_, errs := afSchema.ValidateFunc(validFamily, "address_family")
		assert.Empty(t, errs, "expected %s to be valid", validFamily)
	}

	for _, invalidFamily := range []string{"IPV6", "INVALID_FAMILY"} {
		_, errs := afSchema.ValidateFunc(invalidFamily, "address_family")
		assert.NotEmpty(t, errs, "expected %s to be invalid", invalidFamily)
	}
}

func TestMockResourceNsxtPolicyRouteControllerBgpNeighborRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNbrSDK, restore := setupRCBgpNeighborMocks(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockNbrSDK.EXPECT().Get(routeControllerID, rcBgpNeighborID).Return(rcBgpNeighborAPIResponse(), nil)

		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": rcBgpNeighborParentPath,
		})
		d.SetId(rcBgpNeighborID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, rcBgpNeighborDisplayName, d.Get("display_name"))
		assert.Equal(t, rcBgpNeighborDescription, d.Get("description"))
		assert.Equal(t, rcBgpNeighborPath, d.Get("path"))
		assert.Equal(t, int(rcBgpNeighborRevision), d.Get("revision"))
		assert.Equal(t, rcBgpNeighborAddress, d.Get("neighbor_address"))
		assert.Equal(t, rcBgpNeighborRemoteAsNum, d.Get("remote_as_num"))
		assert.Equal(t, rcBgpNeighborID, d.Id())
	})

	t.Run("Read fails when not found", func(t *testing.T) {
		mockNbrSDK.EXPECT().Get(routeControllerID, rcBgpNeighborID).Return(model.RouteControllerBgpNeighborConfig{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": rcBgpNeighborParentPath,
		})
		d.SetId(rcBgpNeighborID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": rcBgpNeighborParentPath,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RouteControllerBgpNeighbor ID")
	})
}

func TestMockResourceNsxtPolicyRouteControllerBgpNeighborUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNbrSDK, restore := setupRCBgpNeighborMocks(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		resp := rcBgpNeighborAPIResponse()
		gomock.InOrder(
			mockNbrSDK.EXPECT().Update(routeControllerID, rcBgpNeighborID, gomock.Any()).Return(resp, nil),
			mockNbrSDK.EXPECT().Get(routeControllerID, rcBgpNeighborID).Return(resp, nil),
		)

		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRCBgpNeighborData())
		d.SetId(rcBgpNeighborID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRCBgpNeighborData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RouteControllerBgpNeighbor ID")
	})

	t.Run("Update fails when API returns error", func(t *testing.T) {
		mockNbrSDK.EXPECT().Update(routeControllerID, rcBgpNeighborID, gomock.Any()).Return(model.RouteControllerBgpNeighborConfig{}, errors.New("API error"))

		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRCBgpNeighborData())
		d.SetId(rcBgpNeighborID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestMockResourceNsxtPolicyRouteControllerBgpNeighborDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNbrSDK, restore := setupRCBgpNeighborMocks(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockNbrSDK.EXPECT().Delete(routeControllerID, rcBgpNeighborID).Return(nil)

		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": rcBgpNeighborParentPath,
		})
		d.SetId(rcBgpNeighborID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": rcBgpNeighborParentPath,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "RouteControllerBgpNeighbor ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockNbrSDK.EXPECT().Delete(routeControllerID, rcBgpNeighborID).Return(errors.New("API error"))

		res := resourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"parent_path": rcBgpNeighborParentPath,
		})
		d.SetId(rcBgpNeighborID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerBgpNeighborDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}
