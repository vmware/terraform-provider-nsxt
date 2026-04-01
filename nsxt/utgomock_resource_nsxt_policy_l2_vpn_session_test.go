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

	l2vpnapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/l2vpn_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	l2vpnMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s/l2vpn_services"
)

var (
	l2SessionID          = "l2-session-1"
	l2SessionDisplayName = "Test L2 VPN Session"
	l2SessionDescription = "Test l2 vpn session"
	l2SessionRevision    = int64(1)
	l2SessionServicePath = "/infra/tier-1s/t1-gw-1/l2vpn-services/l2-svc-1"
	l2SessionGwID        = "t1-gw-1"
	l2SessionSvcID       = "l2-svc-1"
)

func l2SessionAPIResponse() nsxModel.L2VPNSession {
	return nsxModel.L2VPNSession{
		Id:          &l2SessionID,
		DisplayName: &l2SessionDisplayName,
		Description: &l2SessionDescription,
		Revision:    &l2SessionRevision,
	}
}

func minimalL2SessionData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": l2SessionDisplayName,
		"description":  l2SessionDescription,
		"nsx_id":       l2SessionID,
		"service_path": l2SessionServicePath,
	}
}

func setupL2SessionMock(t *testing.T, ctrl *gomock.Controller) (*l2vpnMocks.MockSessionsClient, func()) {
	mockSDK := l2vpnMocks.NewMockSessionsClient(ctrl)
	mockWrapper := &l2vpnapi.L2VPNSessionClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliT1L2vpnSessionsClient
	cliT1L2vpnSessionsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *l2vpnapi.L2VPNSessionClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliT1L2vpnSessionsClient = original }
}

func TestMockResourceNsxtPolicyL2VPNSessionCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL2SessionMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(nsxModel.L2VPNSession{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(l2SessionGwID, l2SessionSvcID, l2SessionID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(l2SessionAPIResponse(), nil),
		)

		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())

		err := resourceNsxtPolicyL2VPNSessionCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, l2SessionID, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(l2SessionAPIResponse(), nil)

		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())

		err := resourceNsxtPolicyL2VPNSessionCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyL2VPNSessionRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL2SessionMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(l2SessionAPIResponse(), nil)

		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())
		d.SetId(l2SessionID)

		err := resourceNsxtPolicyL2VPNSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, l2SessionDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(nsxModel.L2VPNSession{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())
		d.SetId(l2SessionID)

		err := resourceNsxtPolicyL2VPNSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())

		err := resourceNsxtPolicyL2VPNSessionRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyL2VPNSessionUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL2SessionMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(l2SessionGwID, l2SessionSvcID, l2SessionID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(l2SessionAPIResponse(), nil),
		)

		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())
		d.SetId(l2SessionID)

		err := resourceNsxtPolicyL2VPNSessionUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())

		err := resourceNsxtPolicyL2VPNSessionUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyL2VPNSessionDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL2SessionMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(l2SessionGwID, l2SessionSvcID, l2SessionID).Return(nil)

		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())
		d.SetId(l2SessionID)

		err := resourceNsxtPolicyL2VPNSessionDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyL2VPNSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL2SessionData())

		err := resourceNsxtPolicyL2VPNSessionDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
