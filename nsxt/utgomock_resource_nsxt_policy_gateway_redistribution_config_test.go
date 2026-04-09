//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	t0lsapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t0lsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
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
	util.NsxVersion = "3.1.0"
	defer func() { util.NsxVersion = "" }()

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
	util.NsxVersion = "3.1.0"
	defer func() { util.NsxVersion = "" }()

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
	util.NsxVersion = "3.1.0"
	defer func() { util.NsxVersion = "" }()

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
	util.NsxVersion = "3.1.0"
	defer func() { util.NsxVersion = "" }()

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
