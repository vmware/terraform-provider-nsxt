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
	t0mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
)

var (
	haVipTier0ID         = "t0-gw-1"
	haVipLocaleServiceID = "ls-1"
	haVipIntfPath1       = "/infra/tier-0s/t0-gw-1/locale-services/ls-1/interfaces/intf-1"
	haVipIntfPath2       = "/infra/tier-0s/t0-gw-1/locale-services/ls-1/interfaces/intf-2"
	haVipSubnet          = "192.168.100.10/24"
)

func haVipLocaleServicesResponse() nsxModel.LocaleServices {
	lsID := haVipLocaleServiceID
	lsType := "LocaleServices"
	vipEnabled := true
	prefix := int64(24)
	return nsxModel.LocaleServices{
		Id:           &lsID,
		ResourceType: &lsType,
		HaVipConfigs: []nsxModel.Tier0HaVipConfig{
			{
				Enabled:                &vipEnabled,
				ExternalInterfacePaths: []string{haVipIntfPath1, haVipIntfPath2},
				VipSubnets: []nsxModel.InterfaceSubnet{
					{
						IpAddresses: []string{"192.168.100.10"},
						PrefixLen:   &prefix,
					},
				},
			},
		},
	}
}

func minimalHAVipData() map[string]interface{} {
	return map[string]interface{}{
		"config": []interface{}{
			map[string]interface{}{
				"enabled": true,
				"external_interface_paths": []interface{}{
					haVipIntfPath1,
					haVipIntfPath2,
				},
				"vip_subnets": []interface{}{haVipSubnet},
			},
		},
	}
}

func haVipDataWithComputed() map[string]interface{} {
	d := minimalHAVipData()
	d["tier0_id"] = haVipTier0ID
	d["locale_service_id"] = haVipLocaleServiceID
	return d
}

func setupHAVipMock(t *testing.T, ctrl *gomock.Controller) (*t0mocks.MockLocaleServicesClient, func()) {
	mockSDK := t0mocks.NewMockLocaleServicesClient(ctrl)
	mockWrapper := &tier0sapi.LocaleServicesClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliTier0LocaleServicesClient
	cliTier0LocaleServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier0sapi.LocaleServicesClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliTier0LocaleServicesClient = original }
}

func TestMockResourceNsxtPolicyTier0GatewayHAVipConfigCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHAVipMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(haVipTier0ID, haVipLocaleServiceID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(haVipTier0ID, haVipLocaleServiceID).Return(haVipLocaleServicesResponse(), nil),
		)

		res := resourceNsxtPolicyTier0GatewayHAVipConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHAVipData())

		err := resourceNsxtPolicyTier0GatewayHAVipConfigCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, haVipTier0ID, d.Get("tier0_id"))
		assert.Equal(t, haVipLocaleServiceID, d.Get("locale_service_id"))
	})
}

func TestMockResourceNsxtPolicyTier0GatewayHAVipConfigRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHAVipMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(haVipTier0ID, haVipLocaleServiceID).Return(haVipLocaleServicesResponse(), nil)

		res := resourceNsxtPolicyTier0GatewayHAVipConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, haVipDataWithComputed())
		d.SetId("some-ha-vip-uuid")

		err := resourceNsxtPolicyTier0GatewayHAVipConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(haVipTier0ID, haVipLocaleServiceID).Return(nsxModel.LocaleServices{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTier0GatewayHAVipConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, haVipDataWithComputed())
		d.SetId("some-ha-vip-uuid")

		err := resourceNsxtPolicyTier0GatewayHAVipConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when IDs are empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier0GatewayHAVipConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHAVipData())

		err := resourceNsxtPolicyTier0GatewayHAVipConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier0GatewayHAVipConfigUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHAVipMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(haVipTier0ID, haVipLocaleServiceID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(haVipTier0ID, haVipLocaleServiceID).Return(haVipLocaleServicesResponse(), nil),
		)

		res := resourceNsxtPolicyTier0GatewayHAVipConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, haVipDataWithComputed())
		d.SetId("some-ha-vip-uuid")

		err := resourceNsxtPolicyTier0GatewayHAVipConfigUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when IDs are empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier0GatewayHAVipConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHAVipData())

		err := resourceNsxtPolicyTier0GatewayHAVipConfigUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier0GatewayHAVipConfigDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHAVipMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		lsResponse := haVipLocaleServicesResponse()
		gomock.InOrder(
			mockSDK.EXPECT().Get(haVipTier0ID, haVipLocaleServiceID).Return(lsResponse, nil),
			mockSDK.EXPECT().Update(haVipTier0ID, haVipLocaleServiceID, gomock.Any()).Return(nsxModel.LocaleServices{}, nil),
		)

		res := resourceNsxtPolicyTier0GatewayHAVipConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, haVipDataWithComputed())
		d.SetId("some-ha-vip-uuid")

		err := resourceNsxtPolicyTier0GatewayHAVipConfigDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when IDs are empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier0GatewayHAVipConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHAVipData())

		err := resourceNsxtPolicyTier0GatewayHAVipConfigDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
