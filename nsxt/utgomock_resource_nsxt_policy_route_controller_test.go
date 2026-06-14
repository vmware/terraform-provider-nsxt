//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run:
// mockgen -destination=mocks/infra/RouteControllersClient.go -package=mocks -source=<sdk>/services/nsxt/infra/RouteControllersClient.go RouteControllersClient
// mockgen -destination=mocks/infra/route_controllers/BgpClient.go -package=mocks -source=<sdk>/services/nsxt/infra/route_controllers/BgpClient.go BgpClient

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

	apipkg "github.com/vmware/terraform-provider-nsxt/api"
	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	rcbgpmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/route_controllers"
)

var (
	routeControllerID          = "rc-1"
	routeControllerDisplayName = "rc-fooname"
	routeControllerDescription = "route controller mock"
	routeControllerPath        = "/infra/route-controllers/rc-1"
	routeControllerRevision    = int64(1)
	routeControllerHaMode      = model.RouteController_HA_MODE_STANDBY
)

func rcAPIResponse() model.RouteController {
	return model.RouteController{
		Id:          &routeControllerID,
		DisplayName: &routeControllerDisplayName,
		Description: &routeControllerDescription,
		Path:        &routeControllerPath,
		Revision:    &routeControllerRevision,
		HaMode:      &routeControllerHaMode,
	}
}

func rcBgpAPIResponse() model.RouteControllerBgpRoutingConfig {
	bgpPath := routeControllerPath + "/bgp"
	bgpRevision := int64(0)
	ecmp := true
	mode := model.BgpGracefulRestartConfig_MODE_HELPER_ONLY
	timer := int64(policyBGPGracefulRestartTimerDefault)
	restartTimer := model.BgpGracefulRestartTimer{RestartTimer: &timer}
	restartConfig := model.BgpGracefulRestartConfig{Mode: &mode, Timer: &restartTimer}
	return model.RouteControllerBgpRoutingConfig{
		Path:                  &bgpPath,
		Revision:              &bgpRevision,
		Ecmp:                  &ecmp,
		GracefulRestartConfig: &restartConfig,
	}
}

func minimalRouteControllerData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": routeControllerDisplayName,
		"description":  routeControllerDescription,
	}
}

func setupRCMocks(t *testing.T, ctrl *gomock.Controller) (
	*inframocks.MockRouteControllersClient,
	*inframocks.MockInfraClient,
	*rcbgpmocks.MockBgpClient,
	func(),
) {
	mockRcSDK := inframocks.NewMockRouteControllersClient(ctrl)
	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)
	mockBgpSDK := rcbgpmocks.NewMockBgpClient(ctrl)

	rcWrapper := &cliinfra.RouteControllerClientContext{
		Client:     mockRcSDK,
		ClientType: utl.Local,
	}
	bgpWrapper := &cliinfra.RouteControllerBgpClientContext{
		Client:     mockBgpSDK,
		ClientType: utl.Local,
	}

	origRC := cliRouteControllersClient
	origInfra := cliInfraClient
	origBgp := cliRouteControllerBgpClient

	cliRouteControllersClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.RouteControllerClientContext {
		return rcWrapper
	}
	cliInfraClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}
	cliRouteControllerBgpClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.RouteControllerBgpClientContext {
		return bgpWrapper
	}

	return mockRcSDK, mockInfraSDK, mockBgpSDK, func() {
		cliRouteControllersClient = origRC
		cliInfraClient = origInfra
		cliRouteControllerBgpClient = origBgp
	}
}

func TestMockResourceNsxtPolicyRouteControllerCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRcSDK, mockInfraSDK, mockBgpSDK, restore := setupRCMocks(t, ctrl)
	defer restore()

	t.Run("Create success without bgp_config", func(t *testing.T) {
		gomock.InOrder(
			mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockRcSDK.EXPECT().Get(gomock.Any()).Return(rcAPIResponse(), nil),
			mockBgpSDK.EXPECT().Get(gomock.Any()).Return(model.RouteControllerBgpRoutingConfig{}, vapiErrors.NotFound{}),
		)

		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteControllerData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockRcSDK.EXPECT().Get("existing-id").Return(model.RouteController{Id: &routeControllerID}, nil)

		res := resourceNsxtPolicyRouteController()
		data := minimalRouteControllerData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create success with bgp_config", func(t *testing.T) {
		bgpResp := rcBgpAPIResponse()
		gomock.InOrder(
			mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockRcSDK.EXPECT().Get(gomock.Any()).Return(rcAPIResponse(), nil),
			mockBgpSDK.EXPECT().Get(gomock.Any()).Return(bgpResp, nil),
		)

		res := resourceNsxtPolicyRouteController()
		data := minimalRouteControllerData()
		data["bgp_config"] = []interface{}{map[string]interface{}{
			"local_as_num":                 "65001",
			"ecmp":                         true,
			"multipath_relax":              false,
			"graceful_restart_mode":        model.BgpGracefulRestartConfig_MODE_HELPER_ONLY,
			"graceful_restart_timer":       policyBGPGracefulRestartTimerDefault,
			"peer_route_convergence_timer": 0,
		}}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})

	t.Run("Create fails when H-API returns error", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteControllerData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestMockResourceNsxtPolicyRouteControllerRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRcSDK, _, mockBgpSDK, restore := setupRCMocks(t, ctrl)
	defer restore()

	t.Run("Read success without bgp_config", func(t *testing.T) {
		mockRcSDK.EXPECT().Get(routeControllerID).Return(rcAPIResponse(), nil)
		mockBgpSDK.EXPECT().Get(routeControllerID).Return(model.RouteControllerBgpRoutingConfig{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(routeControllerID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, routeControllerDisplayName, d.Get("display_name"))
		assert.Equal(t, routeControllerDescription, d.Get("description"))
		assert.Equal(t, routeControllerPath, d.Get("path"))
		assert.Equal(t, int(routeControllerRevision), d.Get("revision"))
		assert.Equal(t, routeControllerHaMode, d.Get("ha_mode"))
		assert.Equal(t, routeControllerID, d.Id())
	})

	t.Run("Read success with bgp_config", func(t *testing.T) {
		mockRcSDK.EXPECT().Get(routeControllerID).Return(rcAPIResponse(), nil)
		mockBgpSDK.EXPECT().Get(routeControllerID).Return(rcBgpAPIResponse(), nil)

		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(routeControllerID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, routeControllerDisplayName, d.Get("display_name"))
		bgpList := d.Get("bgp_config").([]interface{})
		require.Len(t, bgpList, 1)
		bgpMap := bgpList[0].(map[string]interface{})
		assert.Equal(t, true, bgpMap["ecmp"])
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining RouteController ID")
	})
}

func TestMockResourceNsxtPolicyRouteControllerUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRcSDK, mockInfraSDK, mockBgpSDK, restore := setupRCMocks(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockRcSDK.EXPECT().Get(routeControllerID).Return(rcAPIResponse(), nil),
			mockBgpSDK.EXPECT().Get(routeControllerID).Return(model.RouteControllerBgpRoutingConfig{}, vapiErrors.NotFound{}),
		)

		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteControllerData())
		d.SetId(routeControllerID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update success with bgp_config", func(t *testing.T) {
		bgpResp := rcBgpAPIResponse()
		gomock.InOrder(
			mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockRcSDK.EXPECT().Get(routeControllerID).Return(rcAPIResponse(), nil),
			mockBgpSDK.EXPECT().Get(routeControllerID).Return(bgpResp, nil),
		)

		res := resourceNsxtPolicyRouteController()
		data := minimalRouteControllerData()
		data["bgp_config"] = []interface{}{map[string]interface{}{
			"local_as_num":                 "65001",
			"ecmp":                         true,
			"multipath_relax":              false,
			"graceful_restart_mode":        model.BgpGracefulRestartConfig_MODE_HELPER_ONLY,
			"graceful_restart_timer":       policyBGPGracefulRestartTimerDefault,
			"peer_route_convergence_timer": 0,
		}}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(routeControllerID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteControllerData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining RouteController ID")
	})

	t.Run("Update fails when H-API returns error", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRouteControllerData())
		d.SetId(routeControllerID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestMockResourceNsxtPolicyRouteControllerDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_, mockInfraSDK, _, restore := setupRCMocks(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)

		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(routeControllerID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining RouteController ID")
	})

	t.Run("Delete fails when H-API returns error", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtPolicyRouteController()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(routeControllerID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyRouteControllerDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}
