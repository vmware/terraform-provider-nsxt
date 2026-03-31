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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	localesvcapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/locale_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	tunnelsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s/locale_services"
)

var (
	greTunnelID          = "gre-tunnel-1"
	greTunnelDisplayName = "Test GRE Tunnel"
	greTunnelDescription = "Test GRE tunnel"
	greTunnelRevision    = int64(1)
	greTunnelTier0ID     = "t0-gw-1"
	greTunnelLocSvcID    = "ls-1"
	greTunnelLocSvcPath  = "/infra/tier-0s/t0-gw-1/locale-services/ls-1"
	greTunnelDstAddr     = "192.0.2.100"
)

func greTunnelStructValue(t *testing.T) *data.StructValue {
	t.Helper()
	deadTime := int64(3)
	enableAck := true
	kaEnabled := false
	interval := int64(10)
	mtu := int64(1476)
	enabled := true

	greTunnel := nsxModel.GreTunnel{
		Id:                 &greTunnelID,
		DisplayName:        &greTunnelDisplayName,
		Description:        &greTunnelDescription,
		Revision:           &greTunnelRevision,
		DestinationAddress: &greTunnelDstAddr,
		Enabled:            &enabled,
		Mtu:                &mtu,
		ResourceType:       nsxModel.GreTunnel__TYPE_IDENTIFIER,
		TunnelKeepalive: &nsxModel.TunnelKeepAlive{
			DeadTimeMultiplier: &deadTime,
			EnableKeepaliveAck: &enableAck,
			Enabled:            &kaEnabled,
			KeepaliveInterval:  &interval,
		},
	}
	converter := bindings.NewTypeConverter()
	dataValue, errs := converter.ConvertToVapi(greTunnel, nsxModel.GreTunnelBindingType())
	if errs != nil {
		t.Fatalf("failed to convert GreTunnel to StructValue: %v", errs[0])
	}
	return dataValue.(*data.StructValue)
}

func minimalGreTunnelData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":        greTunnelDisplayName,
		"description":         greTunnelDescription,
		"nsx_id":              greTunnelID,
		"locale_service_path": greTunnelLocSvcPath,
		"destination_address": greTunnelDstAddr,
		"enabled":             true,
		"mtu":                 1476,
		"tunnel_address": []interface{}{
			map[string]interface{}{
				"edge_path":      "/infra/sites/default/enforcement-points/default/edge-nodes/edge-1",
				"source_address": "10.0.0.1",
				"tunnel_interface_subnet": []interface{}{
					map[string]interface{}{
						"ip_addresses": []interface{}{"10.1.1.1"},
						"prefix_len":   24,
					},
				},
			},
		},
	}
}

func setupGreTunnelMock(t *testing.T, ctrl *gomock.Controller) (*tunnelsmocks.MockTunnelsClient, func()) {
	mockSDK := tunnelsmocks.NewMockTunnelsClient(ctrl)
	mockWrapper := &localesvcapi.TunnelsClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliTunnelsClient
	cliTunnelsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *localesvcapi.TunnelsClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliTunnelsClient = original }
}

func TestMockResourceNsxtPolicyTier0GatewayGRETunnelCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGreTunnelMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(greTunnelTier0ID, greTunnelLocSvcID, greTunnelID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(greTunnelTier0ID, greTunnelLocSvcID, greTunnelID).Return(greTunnelStructValue(t), nil),
		)

		res := resourceNsxtPolicyTier0GatewayGRETunnel()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGreTunnelData())

		err := resourceNsxtPolicyTier0GatewayGRETunnelCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, greTunnelID, d.Id())
		assert.Equal(t, greTunnelDisplayName, d.Get("display_name"))
	})
}

func TestMockResourceNsxtPolicyTier0GatewayGRETunnelRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGreTunnelMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(greTunnelTier0ID, greTunnelLocSvcID, greTunnelID).Return(greTunnelStructValue(t), nil)

		res := resourceNsxtPolicyTier0GatewayGRETunnel()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGreTunnelData())
		d.SetId(greTunnelID)

		err := resourceNsxtPolicyTier0GatewayGRETunnelRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, greTunnelDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(greTunnelTier0ID, greTunnelLocSvcID, greTunnelID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTier0GatewayGRETunnel()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGreTunnelData())
		d.SetId(greTunnelID)

		err := resourceNsxtPolicyTier0GatewayGRETunnelRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTier0GatewayGRETunnel()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGreTunnelData())

		err := resourceNsxtPolicyTier0GatewayGRETunnelRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier0GatewayGRETunnelUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGreTunnelMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(greTunnelTier0ID, greTunnelLocSvcID, greTunnelID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(greTunnelTier0ID, greTunnelLocSvcID, greTunnelID).Return(greTunnelStructValue(t), nil),
		)

		res := resourceNsxtPolicyTier0GatewayGRETunnel()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGreTunnelData())
		d.SetId(greTunnelID)

		err := resourceNsxtPolicyTier0GatewayGRETunnelUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyTier0GatewayGRETunnelDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGreTunnelMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(greTunnelTier0ID, greTunnelLocSvcID, greTunnelID).Return(nil)

		res := resourceNsxtPolicyTier0GatewayGRETunnel()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGreTunnelData())
		d.SetId(greTunnelID)

		err := resourceNsxtPolicyTier0GatewayGRETunnelDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
