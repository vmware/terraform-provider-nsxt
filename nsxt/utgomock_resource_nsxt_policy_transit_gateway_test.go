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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	apiRoot "github.com/vmware/terraform-provider-nsxt/api"
	apiprojects "github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	transitgateways "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	tgwrouting "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways/routing"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	nsxtmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsxt"
	tgwmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects"
	ccmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways"
	tgwroutingmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways/routing"
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
	bgp     *tgwroutingmocks.MockBgpClient
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

	mockBGPSDK := tgwroutingmocks.NewMockBgpClient(ctrl)
	bgpWrapper := &tgwrouting.BgpClientContext{
		Client:     mockBGPSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwProjectID,
	}

	originalTGW := cliTransitGatewaysClient
	originalOrgRoot := cliOrgRootClient
	originalCC := cliTransitGatewayCentralizedConfigsClient
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
	cliTransitGatewayBgpClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tgwrouting.BgpClientContext {
		return bgpWrapper
	}

	t.Cleanup(func() {
		cliTransitGatewaysClient = originalTGW
		cliOrgRootClient = originalOrgRoot
		cliTransitGatewayCentralizedConfigsClient = originalCC
		cliTransitGatewayBgpClient = originalBGP
	})

	return &tgwMocks{
		tgw:     mockTGWSDK,
		orgRoot: mockOrgRoot,
		cc:      mockCCSDK,
		bgp:     mockBGPSDK,
	}
}

func tgwBgpAPIResponse() nsxModel.TransitGatewayBgpRoutingConfig {
	enabled := true
	ecmp := true
	localAsNum := tgwBgpLocalAsNum
	restartMode := tgwBgpGracefulMode
	restartTimer := int64(tgwBgpGracefulTimer)
	staleTimer := int64(tgwBgpGracefulStaleTimer)
	revision := int64(0)
	id := "bgp"
	return nsxModel.TransitGatewayBgpRoutingConfig{
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
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGatewayBgpRoutingConfig{}, vapiErrors.NotFound{}),
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
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGatewayBgpRoutingConfig{}, vapiErrors.NotFound{}),
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
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGatewayBgpRoutingConfig{}, vapiErrors.NotFound{}),
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

func TestMockResourceNsxtPolicyTransitGatewayBgpConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := setupTransitGatewayMockFull(t, ctrl)

	t.Run("Create with bgp_config includes BGP in initial H-API OrgRoot call", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGateway{}, notFoundErr),
			m.orgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			m.cc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
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
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGatewayBgpRoutingConfig{}, vapiErrors.NotFound{}),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		bgpList := d.Get("bgp_config").([]interface{})
		assert.Empty(t, bgpList)
	})

	// Guards against BZ#3742235: NSX omits graceful_restart_config.timer entirely
	// until it's been explicitly configured (confirmed from a live NSX response -
	// {"mode": "HELPER_ONLY"}, no "timer" key), leaving GracefulRestartConfig
	// non-nil but its Timer nil. setTGWBgpConfigInSchema must still fall back to
	// the schema defaults for the two timer fields instead of zero-filling them,
	// since a cached 0 violates NSX's own minimum-value-1 validation if it's ever
	// resent.
	t.Run("Read defaults graceful restart timers when NSX omits Timer", func(t *testing.T) {
		enabled := true
		ecmp := true
		localAsNum := tgwBgpLocalAsNum
		restartMode := tgwBgpGracefulMode
		revision := int64(0)
		id := "bgp"
		bgpResponse := nsxModel.TransitGatewayBgpRoutingConfig{
			Id:         &id,
			Revision:   &revision,
			Enabled:    &enabled,
			Ecmp:       &ecmp,
			LocalAsNum: &localAsNum,
			GracefulRestartConfig: &nsxModel.BgpGracefulRestartConfig{
				Mode: &restartMode,
				// Timer intentionally nil - matches the real NSX GET response.
			},
		}

		gomock.InOrder(
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			m.cc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(bgpResponse, nil),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)

		bgpList := d.Get("bgp_config").([]interface{})
		require.Len(t, bgpList, 1)
		bgpMap := bgpList[0].(map[string]interface{})
		assert.Equal(t, tgwBgpGracefulTimer, bgpMap["graceful_restart_timer"])
		assert.Equal(t, tgwBgpGracefulStaleTimer, bgpMap["graceful_restart_stale_route_timer"])
	})
}

// TestGetSpanFromSchemaZoneBasedEmptyZones guards against a regression where
// span.zone_based_span with an explicitly empty zone_external_ids list was
// dropped from the create/update payload entirely. The SDKv2 diff engine
// emits no attribute entries for the nested block's contents when the list
// stays at zero elements, so the reconstructed ResourceData can hand
// getSpanFromSchema a zone_based_span element that is a bare nil instead of
// a map. getSpanFromSchema must still recognize the block as present (via
// the list length) and produce a ZoneBasedSpan rather than silently
// returning nil, which previously caused NSX to fall back to its own
// default ClusterBasedSpan and produced perpetual drift.
func TestGetSpanFromSchemaZoneBasedEmptyZones(t *testing.T) {
	span := []interface{}{
		map[string]interface{}{
			"cluster_based_span": []interface{}{},
			"zone_based_span":    []interface{}{nil},
		},
	}

	structVal, err := getSpanFromSchema(span)
	require.NoError(t, err)
	require.NotNil(t, structVal, "getSpanFromSchema must not drop a present-but-empty zone_based_span block")

	out, err := setSpanFromSchema(structVal)
	require.NoError(t, err)
	require.NotNil(t, out)

	spanList := out.([]interface{})
	require.Len(t, spanList, 1)
	spanMap := spanList[0].(map[string]interface{})
	zbs := spanMap["zone_based_span"].([]interface{})
	require.Len(t, zbs, 1)
}

// TestGetTGWBgpConfigFromSchemaAdvancedAndRedistribution guards against a
// regression where advanced_config and redistribution_config were never
// merged into the TransitGatewayBgpRoutingConfig sent to NSX. Both map onto
// fields of that same struct (ForwardingUpTimer, RouteRedistributionConfig)
// rather than separate API objects, so they must be read alongside
// bgp_config, not ignored.
func TestGetTGWBgpConfigFromSchemaAdvancedAndRedistribution(t *testing.T) {
	data := minimalTGWData()
	data["bgp_config"] = minimalBgpConfigData()
	data["advanced_config"] = []interface{}{
		map[string]interface{}{"forwarding_up_timer": 10},
	}
	data["redistribution_config"] = []interface{}{
		map[string]interface{}{
			"rule": []interface{}{
				map[string]interface{}{
					"types":          []interface{}{"PUBLIC", "TGW_STATIC_ROUTE"},
					"route_map_path": "",
				},
			},
		},
	}

	res := resourceNsxtPolicyTransitGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, data)

	cfg := getTGWBgpConfigFromSchema(d, true, true, true)
	require.NotNil(t, cfg)
	require.NotNil(t, cfg.ForwardingUpTimer)
	assert.Equal(t, int64(10), *cfg.ForwardingUpTimer)

	require.NotNil(t, cfg.RouteRedistributionConfig)
	require.Len(t, cfg.RouteRedistributionConfig.Rules, 1)
	assert.Equal(t, []string{"PUBLIC", "TGW_STATIC_ROUTE"}, cfg.RouteRedistributionConfig.Rules[0].RouteRedistributionTypes)
}

// TestGetTGWBgpConfigFromSchemaOnlyIncludesChangedBlocks guards against
// BZ#3742235's underlying cause: resourceNsxtPolicyTransitGatewayUpdate used to
// always resend all three of bgp_config/advanced_config/redistribution_config
// together whenever any single one of them changed, dragging along cached
// values (including a possibly-invalid graceful restart timer) for blocks the
// user never touched. NSX's PATCH .../routing/bgp is a documented partial
// update ("only the provided fields are updated; omitted fields retain their
// current values"), so Update now passes per-block include flags gated by
// d.HasChange, and getTGWBgpConfigFromSchema must leave the excluded blocks'
// fields nil (letting the VAPI binding omit them from the request) even when
// they're present in ResourceData.
func TestGetTGWBgpConfigFromSchemaOnlyIncludesChangedBlocks(t *testing.T) {
	data := minimalTGWData()
	data["bgp_config"] = minimalBgpConfigData()
	data["advanced_config"] = []interface{}{
		map[string]interface{}{"forwarding_up_timer": 10},
	}
	data["redistribution_config"] = []interface{}{
		map[string]interface{}{
			"rule": []interface{}{
				map[string]interface{}{
					"types":          []interface{}{"PUBLIC"},
					"route_map_path": "",
				},
			},
		},
	}

	res := resourceNsxtPolicyTransitGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, data)

	// Simulate only redistribution_config having changed: bgp_config and
	// advanced_config are present in ResourceData (as they would be for an
	// already-provisioned TGW) but excluded via their include flags.
	cfg := getTGWBgpConfigFromSchema(d, false, false, true)
	require.NotNil(t, cfg)

	assert.Nil(t, cfg.Ecmp, "bgp_config field must be omitted when bgp_config didn't change")
	assert.Nil(t, cfg.Enabled, "bgp_config field must be omitted when bgp_config didn't change")
	assert.Nil(t, cfg.LocalAsNum, "bgp_config field must be omitted when bgp_config didn't change")
	assert.Nil(t, cfg.GracefulRestartConfig, "bgp_config field must be omitted when bgp_config didn't change")
	assert.Nil(t, cfg.ForwardingUpTimer, "advanced_config field must be omitted when advanced_config didn't change")

	require.NotNil(t, cfg.RouteRedistributionConfig)
	require.Len(t, cfg.RouteRedistributionConfig.Rules, 1)
	assert.Equal(t, []string{"PUBLIC"}, cfg.RouteRedistributionConfig.Rules[0].RouteRedistributionTypes)
}

// TestMockResourceNsxtPolicyTransitGatewayAdvancedAndRedistributionConfig
// guards against a regression where Read unconditionally reset
// advanced_config and redistribution_config to nil regardless of what NSX
// actually returned, causing perpetual drift.
func TestMockResourceNsxtPolicyTransitGatewayAdvancedAndRedistributionConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := setupTransitGatewayMockFull(t, ctrl)

	t.Run("Read sets advanced_config and redistribution_config from bgp API response", func(t *testing.T) {
		bgpResp := tgwBgpAPIResponse()
		timer := int64(10)
		bgpResp.ForwardingUpTimer = &timer
		bgpResp.RouteRedistributionConfig = &nsxModel.TransitGatewayRouteRedistributionConfig{
			Rules: []nsxModel.TransitGatewayBgpRedistributionRule{
				{RouteRedistributionTypes: []string{"PUBLIC", "TGW_STATIC_ROUTE"}},
			},
		}

		gomock.InOrder(
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			m.cc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(bgpResp, nil),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)

		advList := d.Get("advanced_config").([]interface{})
		require.Len(t, advList, 1)
		assert.Equal(t, 10, advList[0].(map[string]interface{})["forwarding_up_timer"])

		redisList := d.Get("redistribution_config").([]interface{})
		require.Len(t, redisList, 1)
		rules := redisList[0].(map[string]interface{})["rule"].([]interface{})
		require.Len(t, rules, 1)
		assert.Equal(t, []interface{}{"PUBLIC", "TGW_STATIC_ROUTE"}, rules[0].(map[string]interface{})["types"])
	})

	t.Run("Read with no advanced or redistribution config sets empty lists", func(t *testing.T) {
		gomock.InOrder(
			m.tgw.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			m.cc.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
			m.bgp.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwBgpAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Get("advanced_config").([]interface{}))
		assert.Empty(t, d.Get("redistribution_config").([]interface{}))
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
