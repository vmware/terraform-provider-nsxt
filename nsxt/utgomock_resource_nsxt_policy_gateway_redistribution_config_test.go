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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	t0lsapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t0lsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
)

var (
	redistGwPath        = "/infra/tier-0s/t0-redist-gw-1"
	redistGwID          = "t0-redist-gw-1"
	redistLocaleService = "default"
	redistEdgeCluster   = "/infra/sites/default/enforcement-points/default/edge-clusters/ec-1"
)

func redistLocaleServicesResponse() nsxModel.LocaleServices {
	lsType := "LocaleServices"
	return nsxModel.LocaleServices{
		Id:              &redistLocaleService,
		ResourceType:    &lsType,
		EdgeClusterPath: &redistEdgeCluster,
		Path:            &redistEdgeCluster,
	}
}

func redistLocaleServicesWithConfig() nsxModel.LocaleServices {
	lsType := "LocaleServices"
	bgpEnabled := true
	config := nsxModel.Tier0RouteRedistributionConfig{
		BgpEnabled: &bgpEnabled,
	}
	ls := nsxModel.LocaleServices{
		Id:                        &redistLocaleService,
		ResourceType:              &lsType,
		EdgeClusterPath:           &redistEdgeCluster,
		Path:                      &redistEdgeCluster,
		RouteRedistributionConfig: &config,
	}
	return ls
}

func minimalRedistConfigData() map[string]interface{} {
	return map[string]interface{}{
		"gateway_path": redistGwPath,
		"bgp_enabled":  true,
		"ospf_enabled": false,
		"rule":         []interface{}{},
	}
}

func setupRedistT0LSMock(t *testing.T, ctrl *gomock.Controller) (*t0lsmocks.MockLocaleServicesClient, func()) {
	mockSDK := t0lsmocks.NewMockLocaleServicesClient(ctrl)
	mockWrapper := &t0lsapi.LocaleServicesClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliTier0LocaleServicesClient
	cliTier0LocaleServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *t0lsapi.LocaleServicesClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier0LocaleServicesClient = original }
}

func TestMockResourceNsxtPolicyGatewayRedistributionConfigCreate(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRedistT0LSMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		lsWithConfig := redistLocaleServicesWithConfig()
		gomock.InOrder(
			// getPolicyTier0GatewayLocaleServiceWithEdgeCluster → Get(gwID, "default")
			mockSDK.EXPECT().Get(redistGwID, defaultPolicyLocaleServiceID).Return(redistLocaleServicesResponse(), nil),
			// policyGatewayRedistributionConfigPatch → Patch(gwID, localeServiceID, ...)
			mockSDK.EXPECT().Patch(redistGwID, redistLocaleService, gomock.Any()).Return(nil),
			// resourceNsxtPolicyGatewayRedistributionConfigRead → Get(gwID, localeServiceID)
			mockSDK.EXPECT().Get(redistGwID, redistLocaleService).Return(lsWithConfig, nil),
		)

		res := resourceNsxtPolicyGatewayRedistributionConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRedistConfigData())

		err := resourceNsxtPolicyGatewayRedistributionConfigCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Create fails for non-Tier0 gateway path", func(t *testing.T) {
		data := minimalRedistConfigData()
		data["gateway_path"] = "/infra/tier-1s/t1-gw-1"

		res := resourceNsxtPolicyGatewayRedistributionConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyGatewayRedistributionConfigCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		require.Contains(t, err.Error(), "Tier0")
	})
}

func TestMockResourceNsxtPolicyGatewayRedistributionConfigRead(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRedistT0LSMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		lsWithConfig := redistLocaleServicesWithConfig()
		mockSDK.EXPECT().Get(redistGwID, redistLocaleService).Return(lsWithConfig, nil)

		res := resourceNsxtPolicyGatewayRedistributionConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRedistConfigData())
		d.SetId("some-uuid")
		d.Set("gateway_id", redistGwID)
		d.Set("locale_service_id", redistLocaleService)

		err := resourceNsxtPolicyGatewayRedistributionConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Read fails when IDs are empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayRedistributionConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRedistConfigData())

		err := resourceNsxtPolicyGatewayRedistributionConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayRedistributionConfigUpdate(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRedistT0LSMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		lsWithConfig := redistLocaleServicesWithConfig()
		gomock.InOrder(
			mockSDK.EXPECT().Patch(redistGwID, redistLocaleService, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(redistGwID, redistLocaleService).Return(lsWithConfig, nil),
		)

		res := resourceNsxtPolicyGatewayRedistributionConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRedistConfigData())
		d.SetId("some-uuid")
		d.Set("gateway_id", redistGwID)
		d.Set("locale_service_id", redistLocaleService)

		err := resourceNsxtPolicyGatewayRedistributionConfigUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayRedistributionConfigDelete(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRedistT0LSMock(t, ctrl)
	defer restore()

	t.Run("Delete success (sets redistribution config to nil)", func(t *testing.T) {
		lsWithConfig := redistLocaleServicesWithConfig()
		gomock.InOrder(
			mockSDK.EXPECT().Get(redistGwID, redistLocaleService).Return(lsWithConfig, nil),
			mockSDK.EXPECT().Update(redistGwID, redistLocaleService, gomock.Any()).Return(nsxModel.LocaleServices{}, nil),
		)

		res := resourceNsxtPolicyGatewayRedistributionConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalRedistConfigData())
		d.SetId("some-uuid")
		d.Set("gateway_id", redistGwID)
		d.Set("locale_service_id", redistLocaleService)

		err := resourceNsxtPolicyGatewayRedistributionConfigDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestUnitNsxt_GlobalGatewayRedistributionCreation(t *testing.T) {
	for _, scenario := range []string{"success", "site_missing", "site_not_matched", "list_error"} {
		t.Run(scenario, func(t *testing.T) {
			raw := map[string]interface{}{"gateway_path": "/infra/tier-0s/gateway", "site_path": "/infra/sites/site"}
			if scenario == "site_missing" {
				delete(raw, "site_path")
			}
			if scenario == "site_not_matched" {
				raw["site_path"] = "/infra/sites/other"
			}
			d := schema.TestResourceDataRaw(t, resourceNsxtPolicyGatewayRedistributionConfig().Schema, raw)
			m := newIsolatedMockProviderClient(t)
			m.Host = "https://nsx.invalid"
			m.PolicyGlobalManager = true
			m.CommonConfig.MaxRetries = 0
			m.PolicyHTTPClient = &http.Client{Transport: unitPolicyTransport(func(r *http.Request) (*http.Response, error) {
				if scenario == "site_missing" {
					t.Fatal("missing site must fail before API access")
				}
				service := `{"id":"locale","path":"/infra/tier-0s/gateway/locale-services/locale","edge_cluster_path":"/infra/sites/site/enforcement-points/default/edge-clusters/cluster","route_redistribution_config":{"bgp_enabled":true}}`
				body := service
				status := 200
				if strings.HasSuffix(r.URL.Path, "/locale-services") {
					body = `{"results":[` + service + `],"result_count":1}`
					if scenario == "list_error" {
						status = 503
						body = `{"error_code":503,"error_message":"unavailable","module_name":"test"}`
					}
				}
				return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			})}
			err := resourceNsxtPolicyGatewayRedistributionConfigCreate(d, m)
			if scenario != "success" {
				require.Error(t, err)
				assert.Empty(t, d.Id())
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, d.Id())
			assert.Equal(t, "gateway", d.Get("gateway_id"))
			assert.Equal(t, "locale", d.Get("locale_service_id"))
		})
	}
}

func TestMockResourceNsxtPolicyRedistributionImport(t *testing.T) {
	for _, tc := range []struct {
		name      string
		config    *nsxModel.Tier0RouteRedistributionConfig
		apiErr    error
		wantError bool
	}{
		{"success", &nsxModel.Tier0RouteRedistributionConfig{}, nil, false},
		{"missing_config", nil, nil, true},
		{"API_error", nil, vapiErrors.ServiceUnavailable{}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			sdk, restore := setupRedistT0LSMock(t, ctrl)
			t.Cleanup(restore)
			sdk.EXPECT().Get("gateway", "locale").Return(nsxModel.LocaleServices{RouteRedistributionConfig: tc.config}, tc.apiErr)
			d := resourceNsxtPolicyGatewayRedistributionConfig().TestResourceData()
			d.SetId("gateway/locale")
			result, err := resourceNsxtPolicyGatewayRedistributionConfigImport(d, newIsolatedMockProviderClient(t))
			if tc.wantError {
				require.ErrorContains(t, err, "Failed to retrieve redistribution config")
				return
			}
			require.NoError(t, err)
			require.Len(t, result, 1)
			assert.NotEmpty(t, d.Id())
			assert.NotEqual(t, "gateway/locale", d.Id())
			assert.Equal(t, "gateway", d.Get("gateway_id"))
			assert.Equal(t, "locale", d.Get("locale_service_id"))
		})
	}
}
