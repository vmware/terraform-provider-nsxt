//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run mockgen for Tier0sClient and LocaleServicesClient
// in api/infra and api/infra/tier_0s.

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

	apipkg "github.com/vmware/terraform-provider-nsxt/api"
	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	tier0localeservices "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	localeservices "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/locale_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t0mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	localeServicesMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
	bgpmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/locale_services"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	t0GatewayID       = "t0-gw-1"
	t0DisplayName     = "tier0-fooname"
	t0Description     = "tier0 mock description"
	t0Path            = "/infra/tier-0s/t0-gw-1"
	t0Revision        = int64(1)
	t0FailoverMode    = "PREEMPTIVE"
	t0HaMode          = "ACTIVE_STANDBY"
	t0DisableFirewall = false
)

func TestMockResourceNsxtPolicyTier0GatewayRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTier0sSDK := t0mocks.NewMockTier0sClient(ctrl)
	mockLocaleServicesSDK := localeServicesMocks.NewMockLocaleServicesClient(ctrl)

	tier0Wrapper := &cliinfra.Tier0ClientContext{
		Client:     mockTier0sSDK,
		ClientType: utl.Local,
	}
	localeServicesWrapper := &tier0localeservices.LocaleServicesClientContext{
		Client:     mockLocaleServicesSDK,
		ClientType: utl.Local,
	}

	originalTier0s := cliTier0sClient
	originalLocaleServices := cliTier0LocaleServicesClient
	defer func() {
		cliTier0sClient = originalTier0s
		cliTier0LocaleServicesClient = originalLocaleServices
	}()

	cliTier0sClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.Tier0ClientContext {
		return tier0Wrapper
	}
	cliTier0LocaleServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *tier0localeservices.LocaleServicesClientContext {
		return localeServicesWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockTier0sSDK.EXPECT().Get(t0GatewayID).Return(model.Tier0{
			DisplayName:     &t0DisplayName,
			Description:     &t0Description,
			Path:            &t0Path,
			Revision:        &t0Revision,
			FailoverMode:    &t0FailoverMode,
			HaMode:          &t0HaMode,
			DisableFirewall: &t0DisableFirewall,
		}, nil)
		resultCount := int64(0)
		mockLocaleServicesSDK.EXPECT().List(t0GatewayID, gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.LocaleServicesListResult{Results: []model.LocaleServices{}, ResultCount: &resultCount}, nil)

		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(t0GatewayID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTier0GatewayRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, t0DisplayName, d.Get("display_name"))
		assert.Equal(t, t0Description, d.Get("description"))
		assert.Equal(t, t0Path, d.Get("path"))
		assert.Equal(t, int(t0Revision), d.Get("revision"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTier0GatewayRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Tier0 ID")
	})
}

func TestMockResourceNsxtPolicyTier0GatewayCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTier0sSDK := t0mocks.NewMockTier0sClient(ctrl)
	mockLocaleServicesSDK := localeServicesMocks.NewMockLocaleServicesClient(ctrl)
	mockInfraSDK := t0mocks.NewMockInfraClient(ctrl)

	tier0Wrapper := &cliinfra.Tier0ClientContext{
		Client:     mockTier0sSDK,
		ClientType: utl.Local,
	}
	localeServicesWrapper := &tier0localeservices.LocaleServicesClientContext{
		Client:     mockLocaleServicesSDK,
		ClientType: utl.Local,
	}

	originalTier0s := cliTier0sClient
	originalLocaleServices := cliTier0LocaleServicesClient
	originalInfra := cliInfraClient
	defer func() {
		cliTier0sClient = originalTier0s
		cliTier0LocaleServicesClient = originalLocaleServices
		cliInfraClient = originalInfra
	}()

	cliTier0sClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.Tier0ClientContext {
		return tier0Wrapper
	}
	cliTier0LocaleServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *tier0localeservices.LocaleServicesClientContext {
		return localeServicesWrapper
	}
	cliInfraClient = func(sessionContext utl.SessionContext, connector client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	t.Run("Create success", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockTier0sSDK.EXPECT().Get(gomock.Any()).Return(model.Tier0{
			DisplayName:     &t0DisplayName,
			Description:     &t0Description,
			Path:            &t0Path,
			Revision:        &t0Revision,
			FailoverMode:    &t0FailoverMode,
			HaMode:          &t0HaMode,
			DisableFirewall: &t0DisableFirewall,
		}, nil)
		resultCount := int64(0)
		mockLocaleServicesSDK.EXPECT().List(gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.LocaleServicesListResult{Results: []model.LocaleServices{}, ResultCount: &resultCount}, nil)

		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": t0DisplayName,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTier0GatewayCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, t0DisplayName, d.Get("display_name"))
	})
}

func TestMockResourceNsxtPolicyTier0GatewayUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTier0sSDK := t0mocks.NewMockTier0sClient(ctrl)
	mockLocaleServicesSDK := localeServicesMocks.NewMockLocaleServicesClient(ctrl)
	mockInfraSDK := t0mocks.NewMockInfraClient(ctrl)

	tier0Wrapper := &cliinfra.Tier0ClientContext{
		Client:     mockTier0sSDK,
		ClientType: utl.Local,
	}
	localeServicesWrapper := &tier0localeservices.LocaleServicesClientContext{
		Client:     mockLocaleServicesSDK,
		ClientType: utl.Local,
	}

	originalTier0s := cliTier0sClient
	originalLocaleServices := cliTier0LocaleServicesClient
	originalInfra := cliInfraClient
	defer func() {
		cliTier0sClient = originalTier0s
		cliTier0LocaleServicesClient = originalLocaleServices
		cliInfraClient = originalInfra
	}()

	cliTier0sClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.Tier0ClientContext {
		return tier0Wrapper
	}
	cliTier0LocaleServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *tier0localeservices.LocaleServicesClientContext {
		return localeServicesWrapper
	}
	cliInfraClient = func(sessionContext utl.SessionContext, connector client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	t.Run("Update success", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockTier0sSDK.EXPECT().Get(t0GatewayID).Return(model.Tier0{
			DisplayName:     &t0DisplayName,
			Description:     &t0Description,
			Path:            &t0Path,
			Revision:        &t0Revision,
			FailoverMode:    &t0FailoverMode,
			HaMode:          &t0HaMode,
			DisableFirewall: &t0DisableFirewall,
		}, nil)
		resultCount := int64(0)
		mockLocaleServicesSDK.EXPECT().List(t0GatewayID, gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.LocaleServicesListResult{Results: []model.LocaleServices{}, ResultCount: &resultCount}, nil)

		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": t0DisplayName,
		})
		d.SetId(t0GatewayID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTier0GatewayUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, t0DisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": t0DisplayName,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTier0GatewayUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Tier0 ID")
	})
}

func TestMockResourceNsxtPolicyTier0GatewayDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInfraSDK := t0mocks.NewMockInfraClient(ctrl)

	originalInfra := cliInfraClient
	defer func() { cliInfraClient = originalInfra }()
	cliInfraClient = func(sessionContext utl.SessionContext, connector client.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	t.Run("Delete success", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)

		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": t0DisplayName,
		})
		d.SetId(t0GatewayID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTier0GatewayDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyTier0GatewayDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Tier0 ID")
	})
}

func TestMockResourceNsxtPolicyTier0GatewayExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTier0sSDK := t0mocks.NewMockTier0sClient(ctrl)
	tier0Wrapper := &cliinfra.Tier0ClientContext{Client: mockTier0sSDK, ClientType: utl.Local}
	original := cliTier0sClient
	defer func() { cliTier0sClient = original }()
	cliTier0sClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.Tier0ClientContext {
		return tier0Wrapper
	}
	ctx := utl.SessionContext{ClientType: utl.Local}

	t.Run("returns true when Get succeeds", func(t *testing.T) {
		mockTier0sSDK.EXPECT().Get(t0GatewayID).Return(model.Tier0{}, nil)
		exists, err := resourceNsxtPolicyTier0GatewayExists(ctx, t0GatewayID, nil)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("returns false on NotFound", func(t *testing.T) {
		mockTier0sSDK.EXPECT().Get(t0GatewayID).Return(model.Tier0{}, vapiErrors.NotFound{})
		exists, err := resourceNsxtPolicyTier0GatewayExists(ctx, t0GatewayID, nil)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("propagates other errors", func(t *testing.T) {
		mockTier0sSDK.EXPECT().Get(t0GatewayID).Return(model.Tier0{}, vapiErrors.InternalServerError{})
		exists, err := resourceNsxtPolicyTier0GatewayExists(ctx, t0GatewayID, nil)
		require.Error(t, err)
		assert.False(t, exists)
	})
}

func TestMockResourceNsxtPolicyTier0GatewayReadBGPConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBgpSDK := bgpmocks.NewMockBgpClient(ctrl)
	bgpWrapper := &localeservices.BgpRoutingConfigClientContext{Client: mockBgpSDK, ClientType: utl.Local}
	original := cliBgpClient
	defer func() { cliBgpClient = original }()
	cliBgpClient = func(sessionContext utl.SessionContext, connector client.Connector) *localeservices.BgpRoutingConfigClientContext {
		return bgpWrapper
	}

	localeServiceID := "default"
	localeService := model.LocaleServices{Id: &localeServiceID}

	t.Run("populates bgp_config on success", func(t *testing.T) {
		enabled := true
		mockBgpSDK.EXPECT().Get(t0GatewayID, localeServiceID).Return(model.BgpRoutingConfig{Enabled: &enabled}, nil)

		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(t0GatewayID)

		err := resourceNsxtPolicyTier0GatewayReadBGPConfig(d, newGoMockProviderClient(), nil, localeService)
		require.NoError(t, err)
		bgpConfigs := d.Get("bgp_config").([]interface{})
		require.Len(t, bgpConfigs, 1)
	})

	t.Run("clears bgp_config when not found", func(t *testing.T) {
		mockBgpSDK.EXPECT().Get(t0GatewayID, localeServiceID).Return(model.BgpRoutingConfig{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(t0GatewayID)

		err := resourceNsxtPolicyTier0GatewayReadBGPConfig(d, newGoMockProviderClient(), nil, localeService)
		require.NoError(t, err)
		assert.Empty(t, d.Get("bgp_config").([]interface{}))
	})

	t.Run("propagates other errors", func(t *testing.T) {
		mockBgpSDK.EXPECT().Get(t0GatewayID, localeServiceID).Return(model.BgpRoutingConfig{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(t0GatewayID)

		err := resourceNsxtPolicyTier0GatewayReadBGPConfig(d, newGoMockProviderClient(), nil, localeService)
		require.Error(t, err)
	})
}

func TestUnitNsxt_isInterSrIbgpSetInConfig(t *testing.T) {
	t.Run("returns false when raw config has no bgp_config block", func(t *testing.T) {
		res := resourceNsxtPolicyTier0Gateway()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		assert.False(t, isInterSrIbgpSetInConfig(d))
	})
}

func TestUnitNsxt_resourceNsxtPolicyTier0GatewayBGPConfigSchemaToStruct(t *testing.T) {
	tagElem := getTagsSchema().Elem.(*schema.Resource)
	tagSet := schema.NewSet(schema.HashResource(tagElem), []interface{}{
		map[string]interface{}{"scope": "s1", "tag": "t1"},
	})
	baseCfgMap := func() map[string]interface{} {
		return map[string]interface{}{
			"revision":                           1,
			"ecmp":                               true,
			"enabled":                            true,
			"local_as_num":                       "65000",
			"multipath_relax":                    true,
			"graceful_restart_mode":              "HELPER_ONLY",
			"graceful_restart_timer":             180,
			"graceful_restart_stale_route_timer": 600,
			"tag":                                tagSet,
			"inter_sr_ibgp":                      true,
			"route_aggregation": []interface{}{
				map[string]interface{}{"prefix": "10.0.0.0/24", "summary_only": true},
			},
		}
	}

	t.Run("non-VRF gateway sets multipath_relax and graceful restart config", func(t *testing.T) {
		result := resourceNsxtPolicyTier0GatewayBGPConfigSchemaToStruct(baseCfgMap(), false, t0GatewayID)
		require.NotNil(t, result.MultipathRelax)
		assert.True(t, *result.MultipathRelax)
		require.NotNil(t, result.GracefulRestartConfig)
		require.NotNil(t, result.InterSrIbgp)
		assert.True(t, *result.InterSrIbgp)
		require.Len(t, result.RouteAggregations, 1)
		assert.Equal(t, "10.0.0.0/24", *result.RouteAggregations[0].Prefix)
		require.NotNil(t, result.LocalAsNum)
		assert.Equal(t, "65000", *result.LocalAsNum)
	})

	t.Run("VRF gateway omits multipath_relax and graceful restart config", func(t *testing.T) {
		result := resourceNsxtPolicyTier0GatewayBGPConfigSchemaToStruct(baseCfgMap(), true, t0GatewayID)
		assert.Nil(t, result.MultipathRelax)
		assert.Nil(t, result.GracefulRestartConfig)
	})

	t.Run("empty local_as_num is omitted", func(t *testing.T) {
		cfgMap := baseCfgMap()
		cfgMap["local_as_num"] = ""
		result := resourceNsxtPolicyTier0GatewayBGPConfigSchemaToStruct(cfgMap, false, t0GatewayID)
		assert.Nil(t, result.LocalAsNum)
	})
}

func TestUnitNsxt_initPolicyTier0ChildBgpConfig(t *testing.T) {
	t.Run("converts a BgpRoutingConfig into a child struct value", func(t *testing.T) {
		enabled := true
		result, err := initPolicyTier0ChildBgpConfig(&model.BgpRoutingConfig{Enabled: &enabled})
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestUnitNsxt_getPolicyVRFConfigFromSchema(t *testing.T) {
	res := resourceNsxtPolicyTier0Gateway()

	t.Run("returns nil when vrf_config is unset", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		assert.Nil(t, getPolicyVRFConfigFromSchema(d))
	})

	t.Run("builds a Tier0VrfConfig from schema", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"vrf_config": []interface{}{
				map[string]interface{}{
					"gateway_path":        "/infra/tier-0s/t0-1",
					"route_distinguisher": "65000:1",
					"evpn_transit_vni":    100,
					"route_target": []interface{}{
						map[string]interface{}{
							"address_family": "IPV4",
							"import_targets": []interface{}{"target:1:1"},
							"export_targets": []interface{}{"target:1:1"},
						},
					},
				},
			},
		})

		config := getPolicyVRFConfigFromSchema(d)
		require.NotNil(t, config)
		assert.Equal(t, "/infra/tier-0s/t0-1", *config.Tier0Path)
		require.NotNil(t, config.RouteDistinguisher)
		assert.Equal(t, "65000:1", *config.RouteDistinguisher)
		require.NotNil(t, config.EvpnTransitVni)
		assert.Equal(t, int64(100), *config.EvpnTransitVni)
		require.Len(t, config.RouteTargets, 1)
		assert.Equal(t, "IPV4", *config.RouteTargets[0].AddressFamily)
	})
}

func TestUnitNsxt_setPolicyVRFConfigInSchema(t *testing.T) {
	res := resourceNsxtPolicyTier0Gateway()

	t.Run("no-op when config is nil", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		require.NoError(t, setPolicyVRFConfigInSchema(d, nil))
	})

	t.Run("sets vrf_config from a Tier0VrfConfig", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		gwPath := "/infra/tier-0s/t0-1"
		addressFamily := "IPV4"
		config := &model.Tier0VrfConfig{
			Tier0Path: &gwPath,
			RouteTargets: []model.VrfRouteTargets{
				{AddressFamily: &addressFamily, ImportRouteTargets: []string{"target:1:1"}},
			},
		}

		require.NoError(t, setPolicyVRFConfigInSchema(d, config))
		vrfConfigs := d.Get("vrf_config").([]interface{})
		require.Len(t, vrfConfigs, 1)
		elem := vrfConfigs[0].(map[string]interface{})
		assert.Equal(t, gwPath, elem["gateway_path"])
	})
}

func TestMockNsxt_initImplicitTier0GatewayLocaleService(t *testing.T) {
	res := resourceNsxtPolicyTier0Gateway()

	t.Run("create flow builds a locale service without fetching an existing one", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_cluster_path": "/infra/sites/default/enforcement-points/default/edge-clusters/ec-1",
		})

		result, err := initImplicitTier0GatewayLocaleService(utl.SessionContext{ClientType: utl.Local}, d, nil, nil)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}
