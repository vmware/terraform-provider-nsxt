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

	apitgw "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	ipsecvcmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways"
)

var (
	tgwIpsecSvcID          = "tgw-vpnsvc-id"
	tgwIpsecSvcDisplayName = "test-tgw-ipsec-service"
	tgwIpsecSvcDescription = "Test TGW IPSec VPN Service"
	tgwIpsecSvcRevision    = int64(1)
	tgwIpsecSvcParentPath  = "/orgs/default/projects/project1/transit-gateways/tgw1"
	tgwIpsecSvcOrgID       = "default"
	tgwIpsecSvcProjectID   = "project1"
	tgwIpsecSvcTGWID       = "tgw1"
)

func tgwIpsecSvcAPIResponse() nsxModel.IPSecVpnService {
	return nsxModel.IPSecVpnService{
		Id:          &tgwIpsecSvcID,
		DisplayName: &tgwIpsecSvcDisplayName,
		Description: &tgwIpsecSvcDescription,
		Revision:    &tgwIpsecSvcRevision,
	}
}

func minimalTGWIpsecSvcData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":  tgwIpsecSvcDisplayName,
		"description":   tgwIpsecSvcDescription,
		"nsx_id":        tgwIpsecSvcID,
		"parent_path":   tgwIpsecSvcParentPath,
		"enabled":       true,
		"ha_sync":       true,
		"ike_log_level": "INFO",
	}
}

func setupTGWIpsecSvcMock(t *testing.T, ctrl *gomock.Controller) (*ipsecvcmocks.MockIpsecVpnServicesClient, func()) {
	mockSDK := ipsecvcmocks.NewMockIpsecVpnServicesClient(ctrl)
	mockWrapper := &apitgw.TransitGatewayIpsecVpnServiceClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwIpsecSvcProjectID,
	}

	original := cliTransitGatewayIpsecVpnServicesClient
	cliTransitGatewayIpsecVpnServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apitgw.TransitGatewayIpsecVpnServiceClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTransitGatewayIpsecVpnServicesClient = original }
}

func TestMockResourceNsxtPolicyTransitGatewayIpsecVpnServiceCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWIpsecSvcMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwIpsecSvcOrgID, tgwIpsecSvcProjectID, tgwIpsecSvcTGWID, tgwIpsecSvcID).Return(nsxModel.IPSecVpnService{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwIpsecSvcOrgID, tgwIpsecSvcProjectID, tgwIpsecSvcTGWID, tgwIpsecSvcID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwIpsecSvcOrgID, tgwIpsecSvcProjectID, tgwIpsecSvcTGWID, tgwIpsecSvcID).Return(tgwIpsecSvcAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayIpsecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWIpsecSvcData())

		err := resourceNsxtPolicyTGWIPSecVpnServicesCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwIpsecSvcID, d.Id())
	})

	t.Run("Create fails with invalid parent_path", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayIpsecVpnService()
		data := minimalTGWIpsecSvcData()
		data["parent_path"] = "/orgs/default"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTGWIPSecVpnServicesCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Create fails when Patch returns error", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwIpsecSvcOrgID, tgwIpsecSvcProjectID, tgwIpsecSvcTGWID, tgwIpsecSvcID).Return(nsxModel.IPSecVpnService{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwIpsecSvcOrgID, tgwIpsecSvcProjectID, tgwIpsecSvcTGWID, tgwIpsecSvcID, gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtPolicyTransitGatewayIpsecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWIpsecSvcData())

		err := resourceNsxtPolicyTGWIPSecVpnServicesCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayIpsecVpnServiceRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWIpsecSvcMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwIpsecSvcOrgID, tgwIpsecSvcProjectID, tgwIpsecSvcTGWID, tgwIpsecSvcID).Return(tgwIpsecSvcAPIResponse(), nil)

		res := resourceNsxtPolicyTransitGatewayIpsecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWIpsecSvcData())
		d.SetId(tgwIpsecSvcID)

		err := resourceNsxtPolicyTGWIPSecVpnServicesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwIpsecSvcDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwIpsecSvcOrgID, tgwIpsecSvcProjectID, tgwIpsecSvcTGWID, tgwIpsecSvcID).Return(nsxModel.IPSecVpnService{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTransitGatewayIpsecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWIpsecSvcData())
		d.SetId(tgwIpsecSvcID)

		err := resourceNsxtPolicyTGWIPSecVpnServicesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayIpsecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWIpsecSvcData())

		err := resourceNsxtPolicyTGWIPSecVpnServicesRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayIpsecVpnServiceUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWIpsecSvcMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(tgwIpsecSvcOrgID, tgwIpsecSvcProjectID, tgwIpsecSvcTGWID, tgwIpsecSvcID, gomock.Any()).Return(tgwIpsecSvcAPIResponse(), nil),
			mockSDK.EXPECT().Get(tgwIpsecSvcOrgID, tgwIpsecSvcProjectID, tgwIpsecSvcTGWID, tgwIpsecSvcID).Return(tgwIpsecSvcAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayIpsecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWIpsecSvcData())
		d.SetId(tgwIpsecSvcID)

		err := resourceNsxtPolicyTGWIPSecVpnServicesUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayIpsecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWIpsecSvcData())

		err := resourceNsxtPolicyTGWIPSecVpnServicesUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayIpsecVpnServiceDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWIpsecSvcMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(tgwIpsecSvcOrgID, tgwIpsecSvcProjectID, tgwIpsecSvcTGWID, tgwIpsecSvcID).Return(nil)

		res := resourceNsxtPolicyTransitGatewayIpsecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWIpsecSvcData())
		d.SetId(tgwIpsecSvcID)

		err := resourceNsxtPolicyTGWIPSecVpnServicesDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayIpsecVpnService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWIpsecSvcData())

		err := resourceNsxtPolicyTGWIPSecVpnServicesDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
