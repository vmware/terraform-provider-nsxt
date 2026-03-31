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

	apipkg "github.com/vmware/terraform-provider-nsxt/api"
	apidomains "github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	domainmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains"
)

var (
	gwPolicyID          = "gw-policy-1"
	gwPolicyDisplayName = "Test Gateway Policy"
	gwPolicyDescription = "Test gateway policy"
	gwPolicyRevision    = int64(1)
	gwPolicyDomain      = "default"
	gwPolicyPath        = "/infra/domains/default/gateway-policies/gw-policy-1"
	gwPolicyCategory    = "LocalGatewayRules"
)

func gwPolicyAPIResponse() nsxModel.GatewayPolicy {
	return nsxModel.GatewayPolicy{
		Id:          &gwPolicyID,
		DisplayName: &gwPolicyDisplayName,
		Description: &gwPolicyDescription,
		Revision:    &gwPolicyRevision,
		Path:        &gwPolicyPath,
		Category:    &gwPolicyCategory,
	}
}

func minimalGwPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    gwPolicyDisplayName,
		"description":     gwPolicyDescription,
		"nsx_id":          gwPolicyID,
		"domain":          gwPolicyDomain,
		"locked":          false,
		"stateful":        true,
		"sequence_number": 0,
		"category":        gwPolicyCategory,
	}
}

func setupGwPolicyMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockGatewayPoliciesClient, *inframocks.MockInfraClient, func()) {
	mockSDK := domainmocks.NewMockGatewayPoliciesClient(ctrl)
	mockWrapper := &apidomains.GatewayPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalGwPolicy := cliGatewayPoliciesClient
	cliGatewayPoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.GatewayPolicyClientContext {
		return mockWrapper
	}

	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)
	originalInfra := cliInfraClient
	cliInfraClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	return mockSDK, mockInfraSDK, func() {
		cliGatewayPoliciesClient = originalGwPolicy
		cliInfraClient = originalInfra
	}
}

func TestMockResourceNsxtPolicyGatewayPolicyCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gwPolicyDomain, gwPolicyID).Return(nsxModel.GatewayPolicy{}, notFoundErr),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gwPolicyDomain, gwPolicyID).Return(gwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwPolicyData())

		err := resourceNsxtPolicyGatewayPolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, gwPolicyID, d.Id())
		assert.Equal(t, gwPolicyDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(gwPolicyDomain, gwPolicyID).Return(gwPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwPolicyData())

		err := resourceNsxtPolicyGatewayPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(gwPolicyDomain, gwPolicyID).Return(gwPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwPolicyData())
		d.SetId(gwPolicyID)

		err := resourceNsxtPolicyGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, gwPolicyDisplayName, d.Get("display_name"))
		assert.Equal(t, gwPolicyCategory, d.Get("category"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(gwPolicyDomain, gwPolicyID).Return(nsxModel.GatewayPolicy{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwPolicyData())
		d.SetId(gwPolicyID)

		err := resourceNsxtPolicyGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwPolicyData())

		err := resourceNsxtPolicyGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayPolicyUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gwPolicyDomain, gwPolicyID).Return(gwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwPolicyData())
		d.SetId(gwPolicyID)

		err := resourceNsxtPolicyGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwPolicyData())

		err := resourceNsxtPolicyGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayPolicyDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(gwPolicyDomain, gwPolicyID).Return(nil)

		res := resourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwPolicyData())
		d.SetId(gwPolicyID)

		err := resourceNsxtPolicyGatewayPolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwPolicyData())

		err := resourceNsxtPolicyGatewayPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
