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
	t0staticmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
)

var (
	staticRouteID          = "sr-001"
	staticRouteDisplayName = "Test Static Route"
	staticRouteDescription = "Test static route"
	staticRouteRevision    = int64(1)
	staticRouteGwPath      = "/infra/tier-0s/t0-gw-1"
	staticRouteGwID        = "t0-gw-1"
	staticRouteNetwork     = "10.0.0.0/24"
	staticRoutePath        = "/infra/tier-0s/t0-gw-1/static-routes/sr-001"
)

func staticRouteAPIResponse() nsxModel.StaticRoutes {
	return nsxModel.StaticRoutes{
		Id:          &staticRouteID,
		DisplayName: &staticRouteDisplayName,
		Description: &staticRouteDescription,
		Revision:    &staticRouteRevision,
		Network:     &staticRouteNetwork,
		Path:        &staticRoutePath,
	}
}

func minimalStaticRouteData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": staticRouteDisplayName,
		"description":  staticRouteDescription,
		"nsx_id":       staticRouteID,
		"gateway_path": staticRouteGwPath,
		"network":      staticRouteNetwork,
		"next_hop": []interface{}{
			map[string]interface{}{
				"admin_distance": 1,
				"ip_address":     "192.168.1.1",
				"interface":      "",
			},
		},
	}
}

func setupStaticRouteMock(t *testing.T, ctrl *gomock.Controller) (*t0staticmocks.MockStaticRoutesClient, func()) {
	mockSDK := t0staticmocks.NewMockStaticRoutesClient(ctrl)
	mockWrapper := &tier0sapi.StaticRoutesClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalT0 := cliTier0StaticRoutesClient
	cliTier0StaticRoutesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier0sapi.StaticRoutesClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliTier0StaticRoutesClient = originalT0 }
}

func TestMockResourceNsxtPolicyStaticRouteCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupStaticRouteMock(t, ctrl)
	defer restore()

	t.Run("Create success for Tier0 gateway", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(staticRouteGwID, staticRouteID).Return(nsxModel.StaticRoutes{}, notFoundErr),
			mockSDK.EXPECT().Patch(staticRouteGwID, staticRouteID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(staticRouteGwID, staticRouteID).Return(staticRouteAPIResponse(), nil),
		)

		res := resourceNsxtPolicyStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticRouteData())

		err := resourceNsxtPolicyStaticRouteCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, staticRouteID, d.Id())
		assert.Equal(t, staticRouteDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(staticRouteGwID, staticRouteID).Return(staticRouteAPIResponse(), nil)

		res := resourceNsxtPolicyStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticRouteData())

		err := resourceNsxtPolicyStaticRouteCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyStaticRouteRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupStaticRouteMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(staticRouteGwID, staticRouteID).Return(staticRouteAPIResponse(), nil)

		res := resourceNsxtPolicyStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticRouteData())
		d.SetId(staticRouteID)

		err := resourceNsxtPolicyStaticRouteRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, staticRouteDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(staticRouteGwID, staticRouteID).Return(nsxModel.StaticRoutes{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticRouteData())
		d.SetId(staticRouteID)

		err := resourceNsxtPolicyStaticRouteRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticRouteData())

		err := resourceNsxtPolicyStaticRouteRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyStaticRouteUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupStaticRouteMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(staticRouteGwID, staticRouteID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(staticRouteGwID, staticRouteID).Return(staticRouteAPIResponse(), nil),
		)

		res := resourceNsxtPolicyStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticRouteData())
		d.SetId(staticRouteID)

		err := resourceNsxtPolicyStaticRouteUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticRouteData())

		err := resourceNsxtPolicyStaticRouteUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyStaticRouteDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupStaticRouteMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(staticRouteGwID, staticRouteID).Return(nil)

		res := resourceNsxtPolicyStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticRouteData())
		d.SetId(staticRouteID)

		err := resourceNsxtPolicyStaticRouteDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticRouteData())

		err := resourceNsxtPolicyStaticRouteDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
