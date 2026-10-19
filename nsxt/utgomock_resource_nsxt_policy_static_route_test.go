//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	tier0sapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
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

	t.Run("Read (cache enabled): cache miss falls through to backend read and config_scope mode patches the missing provider-managed tag", func(t *testing.T) {
		mc := newGoMockProviderClientCacheEnabled()
		res := resourceNsxtPolicyStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticRouteData())
		d.SetId(staticRouteID)

		query := getCacheQueryKey(resourceTypeStaticRoutes, d, mc)
		tc := gcache.getTypeCache(resourceTypeStaticRoutes)
		tc.data[query] = map[string]*data.StructValue{}
		defer delete(gcache.byTyp, resourceTypeStaticRoutes)

		gomock.InOrder(
			mockSDK.EXPECT().Get(staticRouteGwID, staticRouteID).Return(staticRouteAPIResponse(), nil),
			mockSDK.EXPECT().Patch(staticRouteGwID, staticRouteID, gomock.Any()).
				DoAndReturn(func(_, _ string, patched nsxModel.StaticRoutes) error {
					require.Len(t, patched.Tags, 1)
					assert.Equal(t, "ut-run-1", *patched.Tags[0].Tag)
					return nil
				}),
		)

		err := resourceNsxtPolicyStaticRouteRead(d, mc)
		require.NoError(t, err)
		assert.Equal(t, staticRouteDisplayName, d.Get("display_name"))
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

func TestUnitNsxt_resourceNsxtPolicyStaticRouteImport(t *testing.T) {
	res := resourceNsxtPolicyStaticRoute()

	t.Run("full policy path succeeds and sets gateway_path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/infra/tier-0s/gw-1/static-routes/route-1")

		out, err := resourceNsxtPolicyStaticRouteImport(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, "/infra/tier-0s/gw-1", d.Get("gateway_path"))
	})

	t.Run("legacy format missing a slash is rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("route-1")

		_, err := resourceNsxtPolicyStaticRouteImport(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "gateway-id")
	})

	t.Run("legacy gatewayID/routeID format resolves a tier0 gateway", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockT0SDK := inframocks.NewMockTier0sClient(ctrl)
		t0Wrapper := &cliinfra.Tier0ClientContext{Client: mockT0SDK, ClientType: utl.Local}
		originalT0 := cliTier0sClient
		defer func() { cliTier0sClient = originalT0 }()
		cliTier0sClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.Tier0ClientContext {
			return t0Wrapper
		}
		mockT0SDK.EXPECT().Get(staticRouteGwID).Return(nsxModel.Tier0{Path: &staticRouteGwPath}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(staticRouteGwID + "/" + staticRouteID)

		out, err := resourceNsxtPolicyStaticRouteImport(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, staticRouteGwPath, d.Get("gateway_path"))
		assert.Equal(t, staticRouteID, d.Id())
	})

	t.Run("legacy gatewayID/routeID format falls back to a tier1 gateway when tier0 is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockT0SDK := inframocks.NewMockTier0sClient(ctrl)
		mockT1SDK := inframocks.NewMockTier1sClient(ctrl)
		t0Wrapper := &cliinfra.Tier0ClientContext{Client: mockT0SDK, ClientType: utl.Local}
		t1Wrapper := &cliinfra.Tier1ClientContext{Client: mockT1SDK, ClientType: utl.Local}
		originalT0 := cliTier0sClient
		originalT1 := cliTier1sClient
		defer func() {
			cliTier0sClient = originalT0
			cliTier1sClient = originalT1
		}()
		cliTier0sClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.Tier0ClientContext {
			return t0Wrapper
		}
		cliTier1sClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.Tier1ClientContext {
			return t1Wrapper
		}
		t1Path := "/infra/tier-1s/t1-gw-1"
		mockT0SDK.EXPECT().Get(staticRouteGwID).Return(nsxModel.Tier0{}, vapiErrors.NotFound{})
		mockT1SDK.EXPECT().Get(staticRouteGwID).Return(nsxModel.Tier1{Path: &t1Path}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(staticRouteGwID + "/" + staticRouteID)

		out, err := resourceNsxtPolicyStaticRouteImport(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, t1Path, d.Get("gateway_path"))
	})

	t.Run("legacy gatewayID/routeID format fails when neither tier0 nor tier1 gateway is found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockT0SDK := inframocks.NewMockTier0sClient(ctrl)
		mockT1SDK := inframocks.NewMockTier1sClient(ctrl)
		t0Wrapper := &cliinfra.Tier0ClientContext{Client: mockT0SDK, ClientType: utl.Local}
		t1Wrapper := &cliinfra.Tier1ClientContext{Client: mockT1SDK, ClientType: utl.Local}
		originalT0 := cliTier0sClient
		originalT1 := cliTier1sClient
		defer func() {
			cliTier0sClient = originalT0
			cliTier1sClient = originalT1
		}()
		cliTier0sClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.Tier0ClientContext {
			return t0Wrapper
		}
		cliTier1sClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.Tier1ClientContext {
			return t1Wrapper
		}
		mockT0SDK.EXPECT().Get(staticRouteGwID).Return(nsxModel.Tier0{}, vapiErrors.NotFound{})
		mockT1SDK.EXPECT().Get(staticRouteGwID).Return(nsxModel.Tier1{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(staticRouteGwID + "/" + staticRouteID)

		_, err := resourceNsxtPolicyStaticRouteImport(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("legacy gatewayID/routeID format propagates a non-not-found tier0 error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockT0SDK := inframocks.NewMockTier0sClient(ctrl)
		t0Wrapper := &cliinfra.Tier0ClientContext{Client: mockT0SDK, ClientType: utl.Local}
		originalT0 := cliTier0sClient
		defer func() { cliTier0sClient = originalT0 }()
		cliTier0sClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.Tier0ClientContext {
			return t0Wrapper
		}
		mockT0SDK.EXPECT().Get(staticRouteGwID).Return(nsxModel.Tier0{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(staticRouteGwID + "/" + staticRouteID)

		_, err := resourceNsxtPolicyStaticRouteImport(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestUnitNsxt_StaticRouteTier1Operations(t *testing.T) {
	m := newIsolatedMockProviderClient(t)
	m.Host = "https://nsx.invalid"
	methods := []string{}
	m.PolicyHTTPClient = &http.Client{Transport: unitPolicyTransport(func(r *http.Request) (*http.Response, error) {
		methods = append(methods, r.Method)
		assert.Equal(t, "/policy/api/v1/infra/tier-1s/gateway/static-routes/route", r.URL.Path)
		body := `{"id":"route","display_name":"Route","network":"192.0.2.0/24","next_hops":[{"ip_address":"192.0.2.1","admin_distance":5,"scope":["/infra/tier-1s/gateway/locale-services/default/interfaces/interface"]}]}`
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
	})}
	res := resourceNsxtPolicyStaticRoute()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"gateway_path": "/infra/tier-1s/gateway"})
	d.SetId("route")
	require.NoError(t, res.Read(d, m))
	hops := d.Get("next_hop").([]interface{})
	require.Len(t, hops, 1)
	hop := hops[0].(map[string]interface{})
	assert.Equal(t, "192.0.2.1", hop["ip_address"])
	assert.Equal(t, 5, hop["admin_distance"])
	assert.Equal(t, "/infra/tier-1s/gateway/locale-services/default/interfaces/interface", hop["interface"])
	ctx := utl.SessionContext{ClientType: utl.Local}
	connector := getPolicyConnector(m)
	require.NoError(t, patchNsxtPolicyStaticRoute(ctx, connector, "gateway", nsxModel.StaticRoutes{Id: strPtr("route")}, false))
	require.NoError(t, deleteNsxtPolicyStaticRoute(ctx, connector, "gateway", false, "route"))
	assert.Equal(t, []string{"GET", "PATCH", "DELETE"}, methods)
	for _, tier0 := range []bool{false, true} {
		unsupported := utl.SessionContext{ClientType: utl.ClientType(99)}
		_, err := getNsxtPolicyStaticRouteByID(unsupported, connector, "gateway", tier0, "route")
		require.ErrorContains(t, err, "not supported")
		require.ErrorContains(t, patchNsxtPolicyStaticRoute(unsupported, connector, "gateway", nsxModel.StaticRoutes{Id: strPtr("route")}, tier0), "not supported")
		require.ErrorContains(t, deleteNsxtPolicyStaticRoute(unsupported, connector, "gateway", tier0, "route"), "not supported")
	}
}
