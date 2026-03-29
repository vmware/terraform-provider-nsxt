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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	tier0localeservices "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	localeservices "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/locale_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t0mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	localeServicesMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
	localeSvcBgpMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/locale_services"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	bgpCfgGwPath    = "/infra/tier-0s/t0-bgp"
	bgpCfgGwID      = "t0-bgp"
	bgpCfgServiceID = "default"
	bgpCfgRevision  = int64(1)
	bgpCfgPath      = "/infra/tier-0s/t0-bgp/locale-services/default/bgp"
	bgpCfgEnabled   = true
	bgpCfgEcmp      = true
)

func TestMockResourceNsxtPolicyBgpConfigRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBgpSDK := localeSvcBgpMocks.NewMockBgpClient(ctrl)
	bgpWrapper := &localeservices.BgpRoutingConfigClientContext{
		Client:     mockBgpSDK,
		ClientType: utl.Local,
	}

	originalBgp := cliBgpClient
	defer func() { cliBgpClient = originalBgp }()
	cliBgpClient = func(sessionContext utl.SessionContext, connector client.Connector) *localeservices.BgpRoutingConfigClientContext {
		return bgpWrapper
	}

	res := resourceNsxtPolicyBgpConfig()

	t.Run("Read_success", func(t *testing.T) {
		restartMode := model.BgpGracefulRestartConfig_MODE_HELPER_ONLY
		restartTimer := int64(180)
		staleTimer := int64(600)
		restartTimerStruct := model.BgpGracefulRestartTimer{
			RestartTimer:    &restartTimer,
			StaleRouteTimer: &staleTimer,
		}
		restartCfg := model.BgpGracefulRestartConfig{
			Mode:  &restartMode,
			Timer: &restartTimerStruct,
		}
		mockBgpSDK.EXPECT().Get(bgpCfgGwID, bgpCfgServiceID).Return(model.BgpRoutingConfig{
			Enabled:               &bgpCfgEnabled,
			Ecmp:                  &bgpCfgEcmp,
			Path:                  &bgpCfgPath,
			Revision:              &bgpCfgRevision,
			GracefulRestartConfig: &restartCfg,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path":      bgpCfgGwPath,
			"locale_service_id": bgpCfgServiceID,
		})
		d.SetId("bgp-uuid-1")
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpConfigRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, bgpCfgEnabled, d.Get("enabled"))
		assert.Equal(t, bgpCfgEcmp, d.Get("ecmp"))
	})

	t.Run("Read_fails_when_gateway_path_not_tier0", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path":      "/infra/tier-1s/t1-id",
			"locale_service_id": bgpCfgServiceID,
		})
		d.SetId("bgp-uuid-2")
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpConfigRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Tier0 Gateway path expected")
	})

	t.Run("Read_clears_id_when_not_found", func(t *testing.T) {
		mockBgpSDK.EXPECT().Get(bgpCfgGwID, bgpCfgServiceID).Return(model.BgpRoutingConfig{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path":      bgpCfgGwPath,
			"locale_service_id": bgpCfgServiceID,
		})
		d.SetId("bgp-uuid-3")
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpConfigRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyBgpConfigCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTier0sSDK := t0mocks.NewMockTier0sClient(ctrl)
	mockLocaleServicesSDK := localeServicesMocks.NewMockLocaleServicesClient(ctrl)
	mockBgpSDK := localeSvcBgpMocks.NewMockBgpClient(ctrl)

	tier0Wrapper := &cliinfra.Tier0ClientContext{
		Client:     mockTier0sSDK,
		ClientType: utl.Local,
	}
	localeServicesWrapper := &tier0localeservices.LocaleServicesClientContext{
		Client:     mockLocaleServicesSDK,
		ClientType: utl.Local,
	}
	bgpWrapper := &localeservices.BgpRoutingConfigClientContext{
		Client:     mockBgpSDK,
		ClientType: utl.Local,
	}

	originalTier0s := cliTier0sClient
	originalLocSvc := cliTier0LocaleServicesClient
	originalBgp := cliBgpClient
	defer func() {
		cliTier0sClient = originalTier0s
		cliTier0LocaleServicesClient = originalLocSvc
		cliBgpClient = originalBgp
	}()

	cliTier0sClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.Tier0ClientContext {
		return tier0Wrapper
	}
	cliTier0LocaleServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *tier0localeservices.LocaleServicesClientContext {
		return localeServicesWrapper
	}
	cliBgpClient = func(sessionContext utl.SessionContext, connector client.Connector) *localeservices.BgpRoutingConfigClientContext {
		return bgpWrapper
	}

	res := resourceNsxtPolicyBgpConfig()

	t.Run("Create_success", func(t *testing.T) {
		mockTier0sSDK.EXPECT().Get(bgpCfgGwID).Return(model.Tier0{}, nil)
		mockLocaleServicesSDK.EXPECT().Get(bgpCfgGwID, defaultPolicyLocaleServiceID).Return(model.LocaleServices{
			Id: &bgpCfgServiceID,
		}, nil)
		mockBgpSDK.EXPECT().Patch(bgpCfgGwID, bgpCfgServiceID, gomock.Any(), gomock.Any()).Return(nil)
		restartMode := model.BgpGracefulRestartConfig_MODE_HELPER_ONLY
		restartTimer := int64(180)
		staleTimer := int64(600)
		restartTimerStruct := model.BgpGracefulRestartTimer{
			RestartTimer:    &restartTimer,
			StaleRouteTimer: &staleTimer,
		}
		restartCfg := model.BgpGracefulRestartConfig{
			Mode:  &restartMode,
			Timer: &restartTimerStruct,
		}
		mockBgpSDK.EXPECT().Get(bgpCfgGwID, bgpCfgServiceID).Return(model.BgpRoutingConfig{
			Enabled:               &bgpCfgEnabled,
			Ecmp:                  &bgpCfgEcmp,
			Path:                  &bgpCfgPath,
			Revision:              &bgpCfgRevision,
			GracefulRestartConfig: &restartCfg,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path":                       bgpCfgGwPath,
			"enabled":                            true,
			"ecmp":                               true,
			"graceful_restart_mode":              model.BgpGracefulRestartConfig_MODE_HELPER_ONLY,
			"graceful_restart_timer":             180,
			"graceful_restart_stale_route_timer": 600,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpConfigCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, bgpCfgServiceID, d.Get("locale_service_id"))
	})

	t.Run("Create_fails_when_gateway_path_not_tier0", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path": "/infra/tier-1s/t1-id",
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpConfigCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Tier0 Gateway path expected")
	})

	t.Run("Create_fails_when_no_locale_service", func(t *testing.T) {
		mockTier0sSDK.EXPECT().Get(bgpCfgGwID).Return(model.Tier0{}, nil)
		mockLocaleServicesSDK.EXPECT().Get(bgpCfgGwID, defaultPolicyLocaleServiceID).Return(model.LocaleServices{},
			vapiErrors.NotFound{})
		resultCount := int64(0)
		mockLocaleServicesSDK.EXPECT().List(bgpCfgGwID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(model.LocaleServicesListResult{Results: []model.LocaleServices{}, ResultCount: &resultCount}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path":                       bgpCfgGwPath,
			"graceful_restart_mode":              model.BgpGracefulRestartConfig_MODE_HELPER_ONLY,
			"graceful_restart_timer":             180,
			"graceful_restart_stale_route_timer": 600,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpConfigCreate(d, m)
		require.Error(t, err)
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockTier0sSDK.EXPECT().Get(bgpCfgGwID).Return(model.Tier0{}, nil)
		mockLocaleServicesSDK.EXPECT().Get(bgpCfgGwID, defaultPolicyLocaleServiceID).Return(model.LocaleServices{
			Id: &bgpCfgServiceID,
		}, nil)
		mockBgpSDK.EXPECT().Patch(bgpCfgGwID, bgpCfgServiceID, gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path":                       bgpCfgGwPath,
			"graceful_restart_mode":              model.BgpGracefulRestartConfig_MODE_HELPER_ONLY,
			"graceful_restart_timer":             180,
			"graceful_restart_stale_route_timer": 600,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpConfigCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyBgpConfigUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTier0sSDK := t0mocks.NewMockTier0sClient(ctrl)
	mockBgpSDK := localeSvcBgpMocks.NewMockBgpClient(ctrl)

	tier0Wrapper := &cliinfra.Tier0ClientContext{
		Client:     mockTier0sSDK,
		ClientType: utl.Local,
	}
	bgpWrapper := &localeservices.BgpRoutingConfigClientContext{
		Client:     mockBgpSDK,
		ClientType: utl.Local,
	}

	originalTier0s := cliTier0sClient
	originalBgp := cliBgpClient
	defer func() {
		cliTier0sClient = originalTier0s
		cliBgpClient = originalBgp
	}()
	cliTier0sClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.Tier0ClientContext {
		return tier0Wrapper
	}
	cliBgpClient = func(sessionContext utl.SessionContext, connector client.Connector) *localeservices.BgpRoutingConfigClientContext {
		return bgpWrapper
	}

	res := resourceNsxtPolicyBgpConfig()

	t.Run("Update_success", func(t *testing.T) {
		mockTier0sSDK.EXPECT().Get(bgpCfgGwID).Return(model.Tier0{}, nil)
		mockBgpSDK.EXPECT().Update(bgpCfgGwID, bgpCfgServiceID, gomock.Any(), gomock.Any()).Return(model.BgpRoutingConfig{}, nil)
		restartMode := model.BgpGracefulRestartConfig_MODE_HELPER_ONLY
		restartTimer := int64(180)
		staleTimer := int64(600)
		restartTimerStruct := model.BgpGracefulRestartTimer{
			RestartTimer:    &restartTimer,
			StaleRouteTimer: &staleTimer,
		}
		restartCfg := model.BgpGracefulRestartConfig{
			Mode:  &restartMode,
			Timer: &restartTimerStruct,
		}
		mockBgpSDK.EXPECT().Get(bgpCfgGwID, bgpCfgServiceID).Return(model.BgpRoutingConfig{
			Enabled:               &bgpCfgEnabled,
			Ecmp:                  &bgpCfgEcmp,
			Path:                  &bgpCfgPath,
			Revision:              &bgpCfgRevision,
			GracefulRestartConfig: &restartCfg,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path":                       bgpCfgGwPath,
			"locale_service_id":                  bgpCfgServiceID,
			"enabled":                            true,
			"ecmp":                               true,
			"graceful_restart_mode":              model.BgpGracefulRestartConfig_MODE_HELPER_ONLY,
			"graceful_restart_timer":             180,
			"graceful_restart_stale_route_timer": 600,
		})
		d.SetId("bgp-uuid-update")
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpConfigUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_Update_returns_error", func(t *testing.T) {
		mockTier0sSDK.EXPECT().Get(bgpCfgGwID).Return(model.Tier0{}, nil)
		mockBgpSDK.EXPECT().Update(bgpCfgGwID, bgpCfgServiceID, gomock.Any(), gomock.Any()).Return(model.BgpRoutingConfig{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path":                       bgpCfgGwPath,
			"locale_service_id":                  bgpCfgServiceID,
			"graceful_restart_mode":              model.BgpGracefulRestartConfig_MODE_HELPER_ONLY,
			"graceful_restart_timer":             180,
			"graceful_restart_stale_route_timer": 600,
		})
		d.SetId("bgp-uuid-update-fail")
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpConfigUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyBgpConfigDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	res := resourceNsxtPolicyBgpConfig()

	t.Run("Delete_is_noop", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"gateway_path":      bgpCfgGwPath,
			"locale_service_id": bgpCfgServiceID,
		})
		d.SetId("bgp-uuid-delete")
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyBgpConfigDelete(d, m)
		require.NoError(t, err)
	})
}
