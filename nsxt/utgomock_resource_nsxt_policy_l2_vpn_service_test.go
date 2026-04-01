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

	t1sapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t1Mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s"
)

var (
	l2SvcID          = "l2-svc-1"
	l2SvcDisplayName = "Test L2 VPN Service"
	l2SvcDescription = "Test l2 vpn service"
	l2SvcRevision    = int64(1)
	l2SvcGwPath      = "/infra/tier-1s/t1-gw-1"
	l2SvcGwID        = "t1-gw-1"
	l2SvcPath        = "/infra/tier-1s/t1-gw-1/l2vpn-services/l2-svc-1"
)

func l2SvcAPIResponse() nsxModel.L2VPNService {
	mode := nsxModel.L2VPNService_MODE_SERVER
	return nsxModel.L2VPNService{
		Id:          &l2SvcID,
		DisplayName: &l2SvcDisplayName,
		Description: &l2SvcDescription,
		Revision:    &l2SvcRevision,
		Path:        &l2SvcPath,
		Mode:        &mode,
	}
}

func minimalL2SvcData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": l2SvcDisplayName,
		"description":  l2SvcDescription,
		"nsx_id":       l2SvcID,
		"gateway_path": l2SvcGwPath,
		"mode":         nsxModel.L2VPNService_MODE_SERVER,
	}
}

func setupL2SvcMock(t *testing.T, ctrl *gomock.Controller) (*t1Mocks.MockL2vpnServicesClient, func()) {
	mockSDK := t1Mocks.NewMockL2vpnServicesClient(ctrl)
	mockWrapper := &t1sapi.L2VPNServiceClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliT1L2vpnServicesClient
	cliT1L2vpnServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *t1sapi.L2VPNServiceClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliT1L2vpnServicesClient = original }
}

func TestMockResourceNsxtPolicyL2VpnServiceCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL2SvcMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(l2SvcGwID, l2SvcID).Return(nsxModel.L2VPNService{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(l2SvcGwID, l2SvcID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(l2SvcGwID, l2SvcID).Return(l2SvcAPIResponse(), nil),
		)

		res := resourceNsxtPolicyL2VpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SvcData())

		err := resourceNsxtPolicyL2VpnServiceCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, l2SvcID, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(l2SvcGwID, l2SvcID).Return(l2SvcAPIResponse(), nil)

		res := resourceNsxtPolicyL2VpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SvcData())

		err := resourceNsxtPolicyL2VpnServiceCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyL2VpnServiceRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL2SvcMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(l2SvcGwID, l2SvcID).Return(l2SvcAPIResponse(), nil)

		res := resourceNsxtPolicyL2VpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SvcData())
		d.SetId(l2SvcID)

		err := resourceNsxtPolicyL2VpnServiceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, l2SvcDisplayName, d.Get("display_name"))
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(l2SvcGwID, l2SvcID).Return(nsxModel.L2VPNService{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyL2VpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SvcData())
		d.SetId(l2SvcID)

		err := resourceNsxtPolicyL2VpnServiceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyL2VpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SvcData())

		err := resourceNsxtPolicyL2VpnServiceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyL2VpnServiceUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL2SvcMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(l2SvcGwID, l2SvcID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(l2SvcGwID, l2SvcID).Return(l2SvcAPIResponse(), nil),
		)

		res := resourceNsxtPolicyL2VpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SvcData())
		d.SetId(l2SvcID)

		err := resourceNsxtPolicyL2VpnServiceUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyL2VpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SvcData())

		err := resourceNsxtPolicyL2VpnServiceUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyL2VpnServiceDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL2SvcMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(l2SvcGwID, l2SvcID).Return(nil)

		res := resourceNsxtPolicyL2VpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SvcData())
		d.SetId(l2SvcID)

		err := resourceNsxtPolicyL2VpnServiceDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyL2VpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SvcData())

		err := resourceNsxtPolicyL2VpnServiceDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
