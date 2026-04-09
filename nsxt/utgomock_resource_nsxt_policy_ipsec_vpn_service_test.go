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
	t1Mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s"
)

var (
	ipsecSvcID          = "ipsec-svc-1"
	ipsecSvcDisplayName = "Test IPSec VPN Service"
	ipsecSvcDescription = "Test ipsec vpn service"
	ipsecSvcRevision    = int64(1)
	ipsecSvcGwPath      = "/infra/tier-1s/t1-gw-1"
	ipsecSvcGwID        = "t1-gw-1"
	ipsecSvcPath        = "/infra/tier-1s/t1-gw-1/ipsec-vpn-services/ipsec-svc-1"
)

func ipsecSvcAPIResponse() nsxModel.IPSecVpnService {
	enabled := true
	haSync := true
	logLevel := nsxModel.IPSecVpnService_IKE_LOG_LEVEL_INFO
	return nsxModel.IPSecVpnService{
		Id:          &ipsecSvcID,
		DisplayName: &ipsecSvcDisplayName,
		Description: &ipsecSvcDescription,
		Revision:    &ipsecSvcRevision,
		Path:        &ipsecSvcPath,
		Enabled:     &enabled,
		HaSync:      &haSync,
		IkeLogLevel: &logLevel,
	}
}

func minimalIPSecSvcData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":  ipsecSvcDisplayName,
		"description":   ipsecSvcDescription,
		"nsx_id":        ipsecSvcID,
		"gateway_path":  ipsecSvcGwPath,
		"enabled":       true,
		"ha_sync":       true,
		"ike_log_level": nsxModel.IPSecVpnService_IKE_LOG_LEVEL_INFO,
	}
}

func setupIPSecSvcMock(t *testing.T, ctrl *gomock.Controller) (*t1Mocks.MockIpsecVpnServicesClient, func()) {
	mockSDK := t1Mocks.NewMockIpsecVpnServicesClient(ctrl)
	mockWrapper := &ipsecvpnapi.IPSecVpnServiceClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliTier1IpsecVpnServicesClient
	cliTier1IpsecVpnServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *ipsecvpnapi.IPSecVpnServiceClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier1IpsecVpnServicesClient = original }
}

func TestMockResourceNsxtPolicyIPSecVpnServiceCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSvcMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(ipsecSvcGwID, ipsecSvcID).Return(nsxModel.IPSecVpnService{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(ipsecSvcGwID, ipsecSvcID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(ipsecSvcGwID, ipsecSvcID).Return(ipsecSvcAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())

		err := resourceNsxtPolicyIPSecVpnServiceCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipsecSvcID, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipsecSvcGwID, ipsecSvcID).Return(ipsecSvcAPIResponse(), nil)

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())

		err := resourceNsxtPolicyIPSecVpnServiceCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyIPSecVpnServiceRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSvcMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipsecSvcGwID, ipsecSvcID).Return(ipsecSvcAPIResponse(), nil)

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())
		d.SetId(ipsecSvcID)

		err := resourceNsxtPolicyIPSecVpnServiceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipsecSvcDisplayName, d.Get("display_name"))
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipsecSvcGwID, ipsecSvcID).Return(nsxModel.IPSecVpnService{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())
		d.SetId(ipsecSvcID)

		err := resourceNsxtPolicyIPSecVpnServiceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())

		err := resourceNsxtPolicyIPSecVpnServiceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnServiceUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSvcMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(ipsecSvcGwID, ipsecSvcID, gomock.Any()).Return(ipsecSvcAPIResponse(), nil),
			mockSDK.EXPECT().Get(ipsecSvcGwID, ipsecSvcID).Return(ipsecSvcAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())
		d.SetId(ipsecSvcID)

		err := resourceNsxtPolicyIPSecVpnServiceUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())

		err := resourceNsxtPolicyIPSecVpnServiceUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnServiceDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSvcMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(ipsecSvcGwID, ipsecSvcID).Return(nil)

		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())
		d.SetId(ipsecSvcID)

		err := resourceNsxtPolicyIPSecVpnServiceDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSvcData())

		err := resourceNsxtPolicyIPSecVpnServiceDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
