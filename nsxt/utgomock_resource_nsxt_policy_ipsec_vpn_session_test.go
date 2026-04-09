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
	ipsecSessionID          = "ipsec-session-1"
	ipsecSessionServicePath = "/infra/tier-1s/t1-gw-1/ipsec-vpn-services/svc-1"
	ipsecSessionGwID        = "t1-gw-1"
	ipsecSessionSvcID       = "svc-1"
)

func minimalIPSecSessionData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":               "Test IPSec Session",
		"description":                "Test ipsec vpn session",
		"nsx_id":                     ipsecSessionID,
		"service_path":               ipsecSessionServicePath,
		"vpn_type":                   routeBasedIPSecVpnSession,
		"peer_id":                    "10.20.30.40",
		"peer_address":               "10.20.30.40",
		"ip_addresses":               []interface{}{"192.168.10.1"},
		"prefix_length":              24,
		"enabled":                    true,
		"compliance_suite":           nsxModel.IPSecVpnSession_COMPLIANCE_SUITE_NONE,
		"authentication_mode":        nsxModel.IPSecVpnSession_AUTHENTICATION_MODE_PSK,
		"connection_initiation_mode": nsxModel.IPSecVpnSession_CONNECTION_INITIATION_MODE_INITIATOR,
	}
}

func setupIPSecSessionMock(t *testing.T, ctrl *gomock.Controller) (*t1IpsecMocks.MockSessionsClient, func()) {
	mockSDK := t1IpsecMocks.NewMockSessionsClient(ctrl)
	mockWrapper := &ipsecvpnapi.IpsecVpnSessionClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliTier1IpsecVpnSessionsClient
	cliTier1IpsecVpnSessionsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *ipsecvpnapi.IpsecVpnSessionClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier1IpsecVpnSessionsClient = original }
}

func TestMockResourceNsxtPolicyIPSecVpnSessionRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSessionMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID).Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())

		err := resourceNsxtPolicyIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPSecVpnSessionDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPSecSessionMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(ipsecSessionGwID, ipsecSessionSvcID, ipsecSessionID).Return(nil)

		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())
		d.SetId(ipsecSessionID)

		err := resourceNsxtPolicyIPSecVpnSessionDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecSessionData())

		err := resourceNsxtPolicyIPSecVpnSessionDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
