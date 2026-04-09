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

	apiRoot "github.com/vmware/terraform-provider-nsxt/api"
	apiprojects "github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	transitgateways "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	nsxtmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsxt"
	tgwmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects"
	ccmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways"
)

var (
	tgwID          = "tgw-test-id"
	tgwDisplayName = "test-transit-gateway"
	tgwDescription = "Test Transit Gateway"
	tgwRevision    = int64(1)
	tgwPath        = "/orgs/default/projects/project1/transit-gateways/tgw-test-id"
	tgwOrgID       = "default"
	tgwProjectID   = "project1"
)

func tgwAPIResponse() nsxModel.TransitGateway {
	return nsxModel.TransitGateway{
		Id:          &tgwID,
		DisplayName: &tgwDisplayName,
		Description: &tgwDescription,
		Revision:    &tgwRevision,
		Path:        &tgwPath,
	}
}

func minimalTGWData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": tgwDisplayName,
		"description":  tgwDescription,
		"nsx_id":       tgwID,
		"context": []interface{}{
			map[string]interface{}{
				"project_id":  tgwProjectID,
				"vpc_id":      "",
				"from_global": false,
			},
		},
	}
}

func setupTransitGatewayMock(t *testing.T, ctrl *gomock.Controller) (*tgwmocks.MockTransitGatewaysClient, *nsxtmocks.MockOrgRootClient, *ccmocks.MockCentralizedConfigsClient, func()) {
	mockTGWSDK := tgwmocks.NewMockTransitGatewaysClient(ctrl)
	tgwWrapper := &apiprojects.TransitGatewayClientContext{
		Client:     mockTGWSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwProjectID,
	}

	mockOrgRoot := nsxtmocks.NewMockOrgRootClient(ctrl)
	orgRootWrapper := &apiRoot.OrgRootClientContext{
		Client:     mockOrgRoot,
		ClientType: utl.Local,
	}

	mockCCSDK := ccmocks.NewMockCentralizedConfigsClient(ctrl)
	ccWrapper := &transitgateways.CentralizedConfigsClientContext{
		Client:     mockCCSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwProjectID,
	}

	originalTGW := cliTransitGatewaysClient
	originalOrgRoot := cliOrgRootClient
	originalCC := cliTransitGatewayCentralizedConfigsClient

	cliTransitGatewaysClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiprojects.TransitGatewayClientContext {
		return tgwWrapper
	}
	cliOrgRootClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiRoot.OrgRootClientContext {
		return orgRootWrapper
	}
	cliTransitGatewayCentralizedConfigsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *transitgateways.CentralizedConfigsClientContext {
		return ccWrapper
	}

	return mockTGWSDK, mockOrgRoot, mockCCSDK, func() {
		cliTransitGatewaysClient = originalTGW
		cliOrgRootClient = originalOrgRoot
		cliTransitGatewayCentralizedConfigsClient = originalCC
	}
}

func TestMockResourceNsxtPolicyTransitGatewayCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTGW, mockOrgRoot, mockCC, restore := setupTransitGatewayMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockTGW.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGateway{}, notFoundErr),
			mockOrgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockTGW.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			mockCC.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())

		err := resourceNsxtPolicyTransitGatewayCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwID, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockTGW.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())

		err := resourceNsxtPolicyTransitGatewayCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when OrgRoot Patch returns error", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockTGW.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGateway{}, notFoundErr),
			mockOrgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())

		err := resourceNsxtPolicyTransitGatewayCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTGW, _, mockCC, restore := setupTransitGatewayMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		gomock.InOrder(
			mockTGW.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			mockCC.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwDisplayName, d.Get("display_name"))
		assert.Equal(t, tgwDescription, d.Get("description"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockTGW.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(nsxModel.TransitGateway{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())

		err := resourceNsxtPolicyTransitGatewayRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTGW, mockOrgRoot, mockCC, restore := setupTransitGatewayMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockOrgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockTGW.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID).Return(tgwAPIResponse(), nil),
			mockCC.EXPECT().Get(tgwOrgID, tgwProjectID, tgwID, centralizedConfigID).Return(nsxModel.CentralizedConfig{}, vapiErrors.NotFound{}),
		)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())

		err := resourceNsxtPolicyTransitGatewayUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, mockOrgRoot, _, restore := setupTransitGatewayMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockOrgRoot.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)

		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())
		d.SetId(tgwID)

		err := resourceNsxtPolicyTransitGatewayDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGateway()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWData())

		err := resourceNsxtPolicyTransitGatewayDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
