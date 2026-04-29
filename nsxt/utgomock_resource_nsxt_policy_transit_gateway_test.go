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

	apiRoot "github.com/vmware/terraform-provider-nsxt/api"
	apiprojects "github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	transitgateways "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	nsxtmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsxt"
	tgwmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects"
	ccmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

const (
	tgwBgpLocalAsNum         = "65001"
	tgwBgpGracefulMode       = "HELPER_ONLY"
	tgwBgpGracefulTimer      = 180
	tgwBgpGracefulStaleTimer = 600
)

var (
	tgwID          = "tgw-test-id"
	tgwDisplayName = "test-transit-gateway"
	tgwDescription = "Test Transit Gateway"
	tgwRevision    = int64(1)
	tgwPath        = "/orgs/default/projects/project1/transit-gateways/tgw-test-id"
	tgwOrgID       = "default"
	tgwProjectID   = "project1"
)

func tgwAPIResponse() nsxModel.TransitGateway {
	return nsxModel.TransitGateway{
		Id:          &tgwID,
		DisplayName: &tgwDisplayName,
		Description: &tgwDescription,
		Revision:    &tgwRevision,
		Path:        &tgwPath,
	}
}

func minimalTGWData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": tgwDisplayName,
		"description":  tgwDescription,
		"nsx_id":       tgwID,
		"context": []interface{}{
			map[string]interface{}{
				"project_id":  tgwProjectID,
				"vpc_id":      "",
				"from_global": false,
			},
		},
	}
}

type tgwMocks struct {
	tgw     *tgwmocks.MockTransitGatewaysClient
	orgRoot *nsxtmocks.MockOrgRootClient
	cc      *ccmocks.MockCentralizedConfigsClient
	rc      *ccmocks.MockRoutingClient
	bgp     *ccmocks.MockBgpClient
}

func setupTransitGatewayMock(t *testing.T, ctrl *gomock.Controller) (*tgwmocks.MockTransitGatewaysClient, *nsxtmocks.MockOrgRootClient, *ccmocks.MockCentralizedConfigsClient, *ccmocks.MockRoutingClient, func()) {
	m := setupTransitGatewayMockFull(t, ctrl)
	return m.tgw, m.orgRoot, m.cc, m.rc, func() {}
}

func setupTransitGatewayMockFull(t *testing.T, ctrl *gomock.Controller) *tgwMocks {
	mockTGWSDK := tgwmocks.NewMockTransitGatewaysClient(ctrl)
	tgwWrapper := &apiprojects.TransitGatewayClientContext{
		Client:     mockTGWSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwProjectID,
	}

	mockOrgRoot := nsxtmocks.NewMockOrgRootClient(ctrl)
	orgRootWrapper := &apiRoot.OrgRootClientContext{
		Client:     mockOrgRoot,
		ClientType: utl.Local,
	}

	mockCCSDK := ccmocks.NewMockCentralizedConfigsClient(ctrl)
	ccWrapper := &transitgateways.CentralizedConfigsClientContext{
		Client:     mockCCSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwProjectID,
	}

	mockRCSDK := ccmocks.NewMockRoutingClient(ctrl)
	rcWrapper := &transitgateways.RoutingConfigClientContext{
		Client:     mockRCSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwProjectID,
	}

	mockBGPSDK := ccmocks.NewMockBgpClient(ctrl)
	bgpWrapper := &transitgateways.BgpClientContext{
		Client:     mockBGPSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwProjectID,
	}

	originalTGW := cliTransitGatewaysClient
	originalOrgRoot := cliOrgRootClient
	originalCC := cliTransitGatewayCentralizedConfigsClient
	originalRC := cliTransitGatewayRoutingConfigsClient
	originalBGP := cliTransitGatewayBgpClient

	cliTransitGatewaysClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiprojects.TransitGatewayClientContext {
		return tgwWrapper
	}
	cliOrgRootClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiRoot.OrgRootClientContext {
		return orgRootWrapper
	}
	cliTransitGatewayCentralizedConfigsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *transitgateways.CentralizedConfigsClientContext {
		return ccWrapper
	}
	cliTransitGatewayRoutingConfigsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *transitgateways.RoutingConfigClientContext {
		return rcWrapper
	}
	cliTransitGatewayBgpClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *transitgateways.BgpClientContext {
		return bgpWrapper
	}

	t.Cleanup(func() {
		cliTransitGatewaysClient = originalTGW
		cliOrgRootClient = originalOrgRoot
		cliTransitGatewayCentralizedConfigsClient = originalCC
		cliTransitGatewayRoutingConfigsClient = originalRC
		cliTransitGatewayBgpClient = originalBGP
	})

	return &tgwMocks{
		tgw:     mockTGWSDK,
		orgRoot: mockOrgRoot,
		cc:      mockCCSDK,
		rc:      mockRCSDK,
		bgp:     mockBGPSDK,
	}
}

func tgwBgpAPIResponse() nsxModel.BgpRoutingConfig {
	enabled := true
	ecmp := true
	localAsNum := tgwBgpLocalAsNum
	restartMode := tgwBgpGracefulMode
	restartTimer := int64(tgwBgpGracefulTimer)
	staleTimer := int64(tgwBgpGracefulStaleTimer)
	revision := int64(0)
	id := "bgp"
	return nsxModel.BgpRoutingConfig{
		Id:         &id,
		Revision:   &revision,
		Enabled:    &enabled,
		Ecmp:       &ecmp,
		LocalAsNum: &localAsNum,
		GracefulRestartConfig: &nsxModel.BgpGracefulRestartConfig{
			Mode: &restartMode,
			Timer: &nsxModel.BgpGracefulRestartTimer{
				RestartTimer:    &restartTimer,
				StaleRouteTimer: &staleTimer,
			},
		},
	}
}

func minimalBgpConfigData() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"revision":                           0,
			"path":                               "",
			"ecmp":                               true,
			"enabled":                            true,
			"inter_sr_ibgp":                      false,
			"local_as_num":                       tgwBgpLocalAsNum,
			"multipath_relax":                    false,
			"graceful_restart_mode":              tgwBgpGracefulMode,
			"graceful_restart_timer":             tgwBgpGracefulTimer,
			"graceful_restart_stale_route_timer": tgwBgpGracefulStaleTimer,
			"route_aggregation":                  []interface{}{},
			"tag":                                []interface{}{},
		},
	}
}

func TestMockResourceNsxtPolicyTransitGatewayCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := setupTransitGatewayMockFull(t, ctrl)

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGateway{}, notFoundErr),
			m.orgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			m.cc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
			m.rc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGatewayRoutingConfig{}, vapiErrors.NotFound{}),
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.BgpRoutingConfig{}, vapiErrors.NotFound{}),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())

		err := resourceNsxtPolicyTransitGatewayCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwID, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())

		err := resourceNsxtPolicyTransitGatewayCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when OrgRoot Patch returns error", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGateway{}, notFoundErr),
			m.orgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())

		err := resourceNsxtPolicyTransitGatewayCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Create fails when failover_mode set below NSX 9.2", func(t *testing.T) {
		util.NsxVersion = "9.1.0"
		defer func() { util.NsxVersion = "" }()

		notFoundErr := vapiErrors.NotFound{}
		m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGateway{}, notFoundErr)

		data := minimalTGWData()
		data["centralized_config"] = []interface{}{
			map[string]interface{}{
				"ha_mode":            "ACTIVE_STANDBY",
				"failover_mode":      "PREEMPTIVE",
				"edge_cluster_paths": []interface{}{"/infra/edge-clusters/ec"},
			},
		}

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})
}

func TestMockResourceNsxtPolicyTransitGatewayRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := setupTransitGatewayMockFull(t, ctrl)

	t.Run("Read success", func(t *testing.T) {
		gomock.InOrder(
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			m.cc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
			m.rc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGatewayRoutingConfig{}, vapiErrors.NotFound{}),
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.BgpRoutingConfig{}, vapiErrors.NotFound{}),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwDisplayName, d.Get("display_name"))
		assert.Equal(t, tgwDescription, d.Get("description"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGateway{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())

		err := resourceNsxtPolicyTransitGatewayRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := setupTransitGatewayMockFull(t, ctrl)

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			m.orgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			m.cc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
			m.rc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGatewayRoutingConfig{}, vapiErrors.NotFound{}),
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.BgpRoutingConfig{}, vapiErrors.NotFound{}),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())

		err := resourceNsxtPolicyTransitGatewayUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayRedistributionConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := setupTransitGatewayMockFull(t, ctrl)

	redistributionConfigData := []interface{}{
		map[string]interface{}{
			"rule": []interface{}{
				map[string]interface{}{
					"types":          []interface{}{"PUBLIC", "TGW_STATIC_ROUTE"},
					"route_map_path": "/orgs/default/projects/project1/transit-gateways/tgw-test-id/route-maps/rm1",
				},
			},
		},
	}

	t.Run("Create with redistribution_config patches routing config via H-API", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGateway{}, notFoundErr),
			m.orgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			m.cc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
			m.rc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGatewayRoutingConfig{}, vapiErrors.NotFound{}),
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.BgpRoutingConfig{}, vapiErrors.NotFound{}),
		)

		data := minimalTGWData()
		data["redistribution_config"] = redistributionConfigData

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwID, d.Id())
	})

	t.Run("Read sets redistribution_config from API response", func(t *testing.T) {
		routeMapPath := "/orgs/default/projects/project1/transit-gateways/tgw-test-id/route-maps/rm1"
		rcFromAPI := nsxModel.TransitGatewayRoutingConfig{
			RedistributionConfig: &nsxModel.TransitGatewayRedistributionConfig{
				Rules: []nsxModel.TransitGatewayRedistributionRule{
					{
						RouteRedistributionTypes: []string{"PUBLIC", "TGW_PRIVATE"},
						RouteMapPath:             &routeMapPath,
					},
				},
			},
		}
		gomock.InOrder(
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			m.cc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
			m.rc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(rcFromAPI, nil),
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.BgpRoutingConfig{}, vapiErrors.NotFound{}),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)

		rc := d.Get("redistribution_config").([]interface{})
		require.Len(t, rc, 1)
		rules := rc[0].(map[string]interface{})["rule"].([]interface{})
		require.Len(t, rules, 1)
		rule := rules[0].(map[string]interface{})
		assert.Equal(t, []interface{}{"PUBLIC", "TGW_PRIVATE"}, rule["types"])
		assert.Equal(t, routeMapPath, rule["route_map_path"])
	})
}

func TestMockResourceNsxtPolicyTransitGatewayBgpConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := setupTransitGatewayMockFull(t, ctrl)

	t.Run("Create with bgp_config includes ChildBgpRoutingConfig in H-API call", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGateway{}, notFoundErr),
			m.orgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			m.cc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
			m.rc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGatewayRoutingConfig{}, vapiErrors.NotFound{}),
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwBgpAPIResponse(), nil),
		)

		data := minimalTGWData()
		data["bgp_config"] = minimalBgpConfigData()

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwID, d.Id())

		// Verify bgp_config was set from API response
		bgpList := d.Get("bgp_config").([]interface{})
		require.Len(t, bgpList, 1)
		bgpMap := bgpList[0].(map[string]interface{})
		assert.Equal(t, tgwBgpLocalAsNum, bgpMap["local_as_num"])
	})

	t.Run("Read sets bgp_config from API response", func(t *testing.T) {
		gomock.InOrder(
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			m.cc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
			m.rc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGatewayRoutingConfig{}, vapiErrors.NotFound{}),
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwBgpAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)

		bgpList := d.Get("bgp_config").([]interface{})
		require.Len(t, bgpList, 1)
		bgpMap := bgpList[0].(map[string]interface{})
		assert.Equal(t, true, bgpMap["enabled"])
		assert.Equal(t, tgwBgpLocalAsNum, bgpMap["local_as_num"])
		assert.Equal(t, tgwBgpGracefulMode, bgpMap["graceful_restart_mode"])
		assert.Equal(t, tgwBgpGracefulTimer, bgpMap["graceful_restart_timer"])
		assert.Equal(t, tgwBgpGracefulStaleTimer, bgpMap["graceful_restart_stale_route_timer"])
	})

	t.Run("Read with no bgp_config sets empty list", func(t *testing.T) {
		gomock.InOrder(
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			m.cc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
			m.rc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGatewayRoutingConfig{}, vapiErrors.NotFound{}),
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.BgpRoutingConfig{}, vapiErrors.NotFound{}),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		bgpList := d.Get("bgp_config").([]interface{})
		assert.Empty(t, bgpList)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := setupTransitGatewayMockFull(t, ctrl)

	t.Run("Delete success", func(t *testing.T) {
		m.orgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())

		err := resourceNsxtPolicyTransitGatewayDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
