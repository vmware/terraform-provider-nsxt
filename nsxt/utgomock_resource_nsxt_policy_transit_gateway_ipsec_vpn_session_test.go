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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	apiipsecvpn "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways/ipsec_vpn_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	sessionmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways/ipsec_vpn_services"
)

var (
	tgwVpnSessionID          = "tgw-vpn-session-id"
	tgwVpnSessionDisplayName = "test-tgw-vpn-session"
	tgwVpnSessionDescription = "Test TGW VPN Session"
	tgwVpnSessionRevision    = int64(1)
	tgwVpnSessionParentPath  = "/orgs/default/projects/project1/transit-gateways/tgw1/ipsec-vpn-services/svc1"
	tgwVpnSessionOrgID       = "default"
	tgwVpnSessionProjectID   = "project1"
	tgwVpnSessionTGWID       = "tgw1"
	tgwVpnSessionSvcID       = "svc1"
	tgwVpnSessionPeerID      = "10.0.0.2"
	tgwVpnSessionPeerAddress = "10.0.0.2"
	tgwVpnSessionLEPath      = "/orgs/default/projects/project1/transit-gateways/tgw1/ipsec-vpn-local-endpoints/le1"
)

func minimalTGWVpnSessionData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":               tgwVpnSessionDisplayName,
		"description":                tgwVpnSessionDescription,
		"nsx_id":                     tgwVpnSessionID,
		"parent_path":                tgwVpnSessionParentPath,
		"vpn_type":                   policyBasedIPSecVpnSession,
		"peer_id":                    tgwVpnSessionPeerID,
		"peer_address":               tgwVpnSessionPeerAddress,
		"local_endpoint_path":        tgwVpnSessionLEPath,
		"enabled":                    true,
		"authentication_mode":        nsxModel.IPSecVpnSession_AUTHENTICATION_MODE_PSK,
		"compliance_suite":           nsxModel.IPSecVpnSession_COMPLIANCE_SUITE_NONE,
		"connection_initiation_mode": nsxModel.IPSecVpnSession_CONNECTION_INITIATION_MODE_INITIATOR,
	}
}

func setupTGWVpnSessionMock(t *testing.T, ctrl *gomock.Controller) (*sessionmocks.MockSessionsClient, func()) {
	mockSDK := sessionmocks.NewMockSessionsClient(ctrl)
	mockWrapper := &apiipsecvpn.TransitGatewayIpsecVpnSessionClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwVpnSessionProjectID,
	}

	original := cliTransitGatewayIpsecVpnSessionsClient
	cliTransitGatewayIpsecVpnSessionsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiipsecvpn.TransitGatewayIpsecVpnSessionClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTransitGatewayIpsecVpnSessionsClient = original }
}

func TestMockResourceNsxtPolicyTransitGatewayIPSecVpnSessionCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWVpnSessionMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		converter := bindings.NewTypeConverter()
		rsType := nsxModel.IPSecVpnSession_RESOURCE_TYPE_POLICYBASEDIPSECVPNSESSION
		session := nsxModel.PolicyBasedIPSecVpnSession{
			Id:           &tgwVpnSessionID,
			DisplayName:  &tgwVpnSessionDisplayName,
			ResourceType: rsType,
			PeerAddress:  &tgwVpnSessionPeerAddress,
			PeerId:       &tgwVpnSessionPeerID,
		}
		sv, errs := converter.ConvertToVapi(session, nsxModel.PolicyBasedIPSecVpnSessionBindingType())
		require.Nil(t, errs)
		structVal := sv.(*data.StructValue)

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwVpnSessionOrgID, tgwVpnSessionProjectID, tgwVpnSessionTGWID, tgwVpnSessionSvcID, tgwVpnSessionID).Return(nil, notFoundErr),
			mockSDK.EXPECT().Patch(tgwVpnSessionOrgID, tgwVpnSessionProjectID, tgwVpnSessionTGWID, tgwVpnSessionSvcID, tgwVpnSessionID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwVpnSessionOrgID, tgwVpnSessionProjectID, tgwVpnSessionTGWID, tgwVpnSessionSvcID, tgwVpnSessionID).Return(structVal, nil),
		)

		res := resourceNsxtPolicyTransitGatewayIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWVpnSessionData())

		err := resourceNsxtPolicyTGWIPSecVpnSessionCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwVpnSessionID, d.Id())
	})

	t.Run("Create fails with invalid parent_path", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayIPSecVpnSession()
		data := minimalTGWVpnSessionData()
		data["parent_path"] = "/orgs/default"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTGWIPSecVpnSessionCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Create fails when Patch returns error", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwVpnSessionOrgID, tgwVpnSessionProjectID, tgwVpnSessionTGWID, tgwVpnSessionSvcID, tgwVpnSessionID).Return(nil, notFoundErr),
			mockSDK.EXPECT().Patch(tgwVpnSessionOrgID, tgwVpnSessionProjectID, tgwVpnSessionTGWID, tgwVpnSessionSvcID, tgwVpnSessionID, gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtPolicyTransitGatewayIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWVpnSessionData())

		err := resourceNsxtPolicyTGWIPSecVpnSessionCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayIPSecVpnSessionRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWVpnSessionMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		converter := bindings.NewTypeConverter()
		rsType := nsxModel.IPSecVpnSession_RESOURCE_TYPE_POLICYBASEDIPSECVPNSESSION
		session := nsxModel.PolicyBasedIPSecVpnSession{
			Id:           &tgwVpnSessionID,
			DisplayName:  &tgwVpnSessionDisplayName,
			Revision:     &tgwVpnSessionRevision,
			ResourceType: rsType,
			PeerAddress:  &tgwVpnSessionPeerAddress,
			PeerId:       &tgwVpnSessionPeerID,
		}
		sv, errs := converter.ConvertToVapi(session, nsxModel.PolicyBasedIPSecVpnSessionBindingType())
		require.Nil(t, errs)
		structVal := sv.(*data.StructValue)

		mockSDK.EXPECT().Get(tgwVpnSessionOrgID, tgwVpnSessionProjectID, tgwVpnSessionTGWID, tgwVpnSessionSvcID, tgwVpnSessionID).Return(structVal, nil)

		res := resourceNsxtPolicyTransitGatewayIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWVpnSessionData())
		d.SetId(tgwVpnSessionID)

		err := resourceNsxtPolicyTGWIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwVpnSessionDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwVpnSessionOrgID, tgwVpnSessionProjectID, tgwVpnSessionTGWID, tgwVpnSessionSvcID, tgwVpnSessionID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTransitGatewayIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWVpnSessionData())
		d.SetId(tgwVpnSessionID)

		err := resourceNsxtPolicyTGWIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWVpnSessionData())

		err := resourceNsxtPolicyTGWIPSecVpnSessionRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayIPSecVpnSessionUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWVpnSessionMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		converter := bindings.NewTypeConverter()
		rsType := nsxModel.IPSecVpnSession_RESOURCE_TYPE_POLICYBASEDIPSECVPNSESSION
		session := nsxModel.PolicyBasedIPSecVpnSession{
			Id:           &tgwVpnSessionID,
			DisplayName:  &tgwVpnSessionDisplayName,
			Revision:     &tgwVpnSessionRevision,
			ResourceType: rsType,
			PeerAddress:  &tgwVpnSessionPeerAddress,
			PeerId:       &tgwVpnSessionPeerID,
		}
		sv, errs := converter.ConvertToVapi(session, nsxModel.PolicyBasedIPSecVpnSessionBindingType())
		require.Nil(t, errs)
		structVal := sv.(*data.StructValue)

		gomock.InOrder(
			mockSDK.EXPECT().Patch(tgwVpnSessionOrgID, tgwVpnSessionProjectID, tgwVpnSessionTGWID, tgwVpnSessionSvcID, tgwVpnSessionID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwVpnSessionOrgID, tgwVpnSessionProjectID, tgwVpnSessionTGWID, tgwVpnSessionSvcID, tgwVpnSessionID).Return(structVal, nil),
		)

		res := resourceNsxtPolicyTransitGatewayIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWVpnSessionData())
		d.SetId(tgwVpnSessionID)

		err := resourceNsxtPolicyTGWIPSecVpnSessionUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWVpnSessionData())

		err := resourceNsxtPolicyTGWIPSecVpnSessionUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayIPSecVpnSessionDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWVpnSessionMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(tgwVpnSessionOrgID, tgwVpnSessionProjectID, tgwVpnSessionTGWID, tgwVpnSessionSvcID, tgwVpnSessionID).Return(nil)

		res := resourceNsxtPolicyTransitGatewayIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWVpnSessionData())
		d.SetId(tgwVpnSessionID)

		err := resourceNsxtPolicyTGWIPSecVpnSessionDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayIPSecVpnSession()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWVpnSessionData())

		err := resourceNsxtPolicyTGWIPSecVpnSessionDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
