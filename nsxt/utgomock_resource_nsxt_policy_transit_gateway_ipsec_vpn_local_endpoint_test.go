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

	apiipsec "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways/ipsec_vpn_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	lepmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways/ipsec_vpn_services"
)

var (
	tgwLEID          = "tgw-le-id"
	tgwLEDisplayName = "test-tgw-local-endpoint"
	tgwLEDescription = "Test TGW Local Endpoint"
	tgwLERevision    = int64(1)
	tgwLEParentPath  = "/orgs/default/projects/project1/transit-gateways/tgw1/ipsec-vpn-services/svc1"
	tgwLEOrgID       = "default"
	tgwLEProjectID   = "project1"
	tgwLETGWID       = "tgw1"
	tgwLESvcID       = "svc1"
	tgwLEAddress     = "192.168.1.1"
)

func tgwLEAPIResponse() nsxModel.IPSecVpnLocalEndpoint {
	return nsxModel.IPSecVpnLocalEndpoint{
		Id:           &tgwLEID,
		DisplayName:  &tgwLEDisplayName,
		Description:  &tgwLEDescription,
		Revision:     &tgwLERevision,
		LocalAddress: &tgwLEAddress,
	}
}

func minimalTGWLEData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":  tgwLEDisplayName,
		"description":   tgwLEDescription,
		"nsx_id":        tgwLEID,
		"parent_path":   tgwLEParentPath,
		"local_address": tgwLEAddress,
	}
}

func setupTGWLocalEndpointMock(t *testing.T, ctrl *gomock.Controller) (*lepmocks.MockLocalEndpointsClient, func()) {
	mockSDK := lepmocks.NewMockLocalEndpointsClient(ctrl)
	mockWrapper := &apiipsec.TransitGatewayIpsecVpnLocalEndpointClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwLEProjectID,
	}

	original := cliTransitGatewayIpsecVpnLocalEndpointsClient
	cliTransitGatewayIpsecVpnLocalEndpointsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiipsec.TransitGatewayIpsecVpnLocalEndpointClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTransitGatewayIpsecVpnLocalEndpointsClient = original }
}

func TestMockResourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpointCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWLocalEndpointMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwLEOrgID, tgwLEProjectID, tgwLETGWID, tgwLESvcID, tgwLEID).Return(nsxModel.IPSecVpnLocalEndpoint{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwLEOrgID, tgwLEProjectID, tgwLETGWID, tgwLESvcID, tgwLEID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwLEOrgID, tgwLEProjectID, tgwLETGWID, tgwLESvcID, tgwLEID).Return(tgwLEAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWLEData())

		err := resourceNsxtPolicyTGWIPSecVpnLocalEndpointCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwLEID, d.Id())
	})

	t.Run("Create fails with invalid parent_path", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint()
		data := minimalTGWLEData()
		data["parent_path"] = "/orgs/default"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTGWIPSecVpnLocalEndpointCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Create fails when Patch returns error", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwLEOrgID, tgwLEProjectID, tgwLETGWID, tgwLESvcID, tgwLEID).Return(nsxModel.IPSecVpnLocalEndpoint{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwLEOrgID, tgwLEProjectID, tgwLETGWID, tgwLESvcID, tgwLEID, gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWLEData())

		err := resourceNsxtPolicyTGWIPSecVpnLocalEndpointCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpointRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWLocalEndpointMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwLEOrgID, tgwLEProjectID, tgwLETGWID, tgwLESvcID, tgwLEID).Return(tgwLEAPIResponse(), nil)

		res := resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWLEData())
		d.SetId(tgwLEID)

		err := resourceNsxtPolicyTGWIPSecVpnLocalEndpointRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwLEDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwLEOrgID, tgwLEProjectID, tgwLETGWID, tgwLESvcID, tgwLEID).Return(nsxModel.IPSecVpnLocalEndpoint{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWLEData())
		d.SetId(tgwLEID)

		err := resourceNsxtPolicyTGWIPSecVpnLocalEndpointRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWLEData())

		err := resourceNsxtPolicyTGWIPSecVpnLocalEndpointRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpointUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWLocalEndpointMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(tgwLEOrgID, tgwLEProjectID, tgwLETGWID, tgwLESvcID, tgwLEID, gomock.Any()).Return(tgwLEAPIResponse(), nil),
			mockSDK.EXPECT().Get(tgwLEOrgID, tgwLEProjectID, tgwLETGWID, tgwLESvcID, tgwLEID).Return(tgwLEAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWLEData())
		d.SetId(tgwLEID)

		err := resourceNsxtPolicyTGWIPSecVpnLocalEndpointUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWLEData())

		err := resourceNsxtPolicyTGWIPSecVpnLocalEndpointUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpointDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWLocalEndpointMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(tgwLEOrgID, tgwLEProjectID, tgwLETGWID, tgwLESvcID, tgwLEID).Return(nil)

		res := resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWLEData())
		d.SetId(tgwLEID)

		err := resourceNsxtPolicyTGWIPSecVpnLocalEndpointDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayIPSecVpnLocalEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWLEData())

		err := resourceNsxtPolicyTGWIPSecVpnLocalEndpointDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
