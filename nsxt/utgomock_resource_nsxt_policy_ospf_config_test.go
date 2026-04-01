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

	ospfConfigapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/locale_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	lsMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/locale_services"
)

var (
	ospfCfgGwPath   = "/infra/tier-0s/t0-gw-1"
	ospfCfgGwID     = "t0-gw-1"
	ospfCfgLsID     = "ls-1"
	ospfCfgRevision = int64(1)
)

func ospfConfigAPIResponse() nsxModel.OspfRoutingConfig {
	enabled := true
	return nsxModel.OspfRoutingConfig{
		DisplayName: nil,
		Revision:    &ospfCfgRevision,
		Enabled:     &enabled,
	}
}

func minimalOSPFConfigData() map[string]interface{} {
	return map[string]interface{}{
		"gateway_path":          ospfCfgGwPath,
		"enabled":               true,
		"ecmp":                  true,
		"default_originate":     false,
		"graceful_restart_mode": nsxModel.OspfRoutingConfig_GRACEFUL_RESTART_MODE_HELPER_ONLY,
	}
}

func setupOSPFConfigMock(t *testing.T, ctrl *gomock.Controller) (*lsMocks.MockOspfClient, func()) {
	mockSDK := lsMocks.NewMockOspfClient(ctrl)
	mockWrapper := &ospfConfigapi.OspfRoutingConfigClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliOspfClient
	cliOspfClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *ospfConfigapi.OspfRoutingConfigClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliOspfClient = original }
}

func TestMockResourceNsxtPolicyOspfConfigRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupOSPFConfigMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(ospfCfgGwID, ospfCfgLsID).Return(ospfConfigAPIResponse(), nil)

		res := resourceNsxtPolicyOspfConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalOSPFConfigData())
		d.SetId("ospf-cfg-1")
		d.Set("locale_service_id", ospfCfgLsID)

		err := resourceNsxtPolicyOspfConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, true, d.Get("enabled"))
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(ospfCfgGwID, ospfCfgLsID).Return(nsxModel.OspfRoutingConfig{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyOspfConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalOSPFConfigData())
		d.SetId("ospf-cfg-1")
		d.Set("locale_service_id", ospfCfgLsID)

		err := resourceNsxtPolicyOspfConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyOspfConfigUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupOSPFConfigMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(ospfCfgGwID, ospfCfgLsID, gomock.Any()).Return(ospfConfigAPIResponse(), nil),
			mockSDK.EXPECT().Get(ospfCfgGwID, ospfCfgLsID).Return(ospfConfigAPIResponse(), nil),
		)

		res := resourceNsxtPolicyOspfConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalOSPFConfigData())
		d.SetId("ospf-cfg-1")
		d.Set("gateway_id", ospfCfgGwID)
		d.Set("locale_service_id", ospfCfgLsID)

		err := resourceNsxtPolicyOspfConfigUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyOspfConfigDelete(t *testing.T) {
	t.Run("Delete is a no-op", func(t *testing.T) {
		res := resourceNsxtPolicyOspfConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalOSPFConfigData())
		d.SetId("ospf-cfg-1")

		err := resourceNsxtPolicyOspfConfigDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
