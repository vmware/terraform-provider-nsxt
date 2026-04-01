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

	ipsecvpnapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/ipsec_vpn_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t1IpsecMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s/ipsec_vpn_services"
)

var (
	ipsecEPID          = "ep-1"
	ipsecEPDisplayName = "Test IPSec Local Endpoint"
	ipsecEPDescription = "Test ipsec local endpoint"
	ipsecEPRevision    = int64(1)
	ipsecEPServicePath = "/infra/tier-1s/t1-gw-1/ipsec-vpn-services/svc-1"
	ipsecEPGwID        = "t1-gw-1"
	ipsecEPSvcID       = "svc-1"
	ipsecEPLocalAddr   = "192.168.1.1"
)

func ipsecEPAPIResponse() nsxModel.IPSecVpnLocalEndpoint {
	return nsxModel.IPSecVpnLocalEndpoint{
		Id:           &ipsecEPID,
		DisplayName:  &ipsecEPDisplayName,
		Description:  &ipsecEPDescription,
		Revision:     &ipsecEPRevision,
		LocalAddress: &ipsecEPLocalAddr,
	}
}

func minimalIPSecEPData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":  ipsecEPDisplayName,
		"description":   ipsecEPDescription,
		"nsx_id":        ipsecEPID,
		"service_path":  ipsecEPServicePath,
		"local_address": ipsecEPLocalAddr,
	}
}

func setupIPSecEPMock(t *testing.T, ctrl *gomock.Controller) (*t1IpsecMocks.MockLocalEndpointsClient, func()) {
	mockSDK := t1IpsecMocks.NewMockLocalEndpointsClient(ctrl)
	mockWrapper := &ipsecvpnapi.IPSecVpnLocalEndpointClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliTier1IpsecVpnLocalEndpointsClient
	cliTier1IpsecVpnLocalEndpointsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *ipsecvpnapi.IPSecVpnLocalEndpointClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier1IpsecVpnLocalEndpointsClient = original }
}

func TestMockResourceNsxtPolicyIPSecVpnLocalEndpointCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecEPMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(ipsecEPGwID, ipsecEPSvcID, ipsecEPID).Return(nsxModel.IPSecVpnLocalEndpoint{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(ipsecEPGwID, ipsecEPSvcID, ipsecEPID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(ipsecEPGwID, ipsecEPSvcID, ipsecEPID).Return(ipsecEPAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())

		err := resourceNsxtPolicyIPSecVpnLocalEndpointCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipsecEPID, d.Id())
	})
}

func TestMockResourceNsxtPolicyIPSecVpnLocalEndpointRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecEPMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipsecEPGwID, ipsecEPSvcID, ipsecEPID).Return(ipsecEPAPIResponse(), nil)

		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())
		d.SetId(ipsecEPID)

		err := resourceNsxtPolicyIPSecVpnLocalEndpointRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipsecEPDisplayName, d.Get("display_name"))
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipsecEPGwID, ipsecEPSvcID, ipsecEPID).Return(nsxModel.IPSecVpnLocalEndpoint{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())
		d.SetId(ipsecEPID)

		err := resourceNsxtPolicyIPSecVpnLocalEndpointRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())

		err := resourceNsxtPolicyIPSecVpnLocalEndpointRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnLocalEndpointUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecEPMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(ipsecEPGwID, ipsecEPSvcID, ipsecEPID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(ipsecEPGwID, ipsecEPSvcID, ipsecEPID).Return(ipsecEPAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())
		d.SetId(ipsecEPID)

		err := resourceNsxtPolicyIPSecVpnLocalEndpointUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())

		err := resourceNsxtPolicyIPSecVpnLocalEndpointUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnLocalEndpointDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecEPMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(ipsecEPGwID, ipsecEPSvcID, ipsecEPID).Return(nil)

		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())
		d.SetId(ipsecEPID)

		err := resourceNsxtPolicyIPSecVpnLocalEndpointDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecEPData())

		err := resourceNsxtPolicyIPSecVpnLocalEndpointDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
