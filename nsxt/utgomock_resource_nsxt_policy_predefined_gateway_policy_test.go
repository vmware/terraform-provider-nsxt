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

	apipkg "github.com/vmware/terraform-provider-nsxt/api"
	apidomains "github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	domainmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains"
)

var (
	predefinedGwPolicyID          = "default"
	predefinedGwPolicyDomain      = "default"
	predefinedGwPolicyPath        = "/infra/domains/default/gateway-policies/default"
	predefinedGwPolicyDescription = "Predefined Gateway Policy"
	predefinedGwPolicyRevision    = int64(1)
	predefinedGwPolicyCategory    = "Default"
)

func predefinedGwPolicyAPIResponse() nsxModel.GatewayPolicy {
	return nsxModel.GatewayPolicy{
		Id:          &predefinedGwPolicyID,
		Description: &predefinedGwPolicyDescription,
		Revision:    &predefinedGwPolicyRevision,
		Path:        &predefinedGwPolicyPath,
		Category:    &predefinedGwPolicyCategory,
	}
}

func minimalPredefinedGwPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"path": predefinedGwPolicyPath,
	}
}

func setupPredefinedGwPolicyMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockGatewayPoliciesClient, *inframocks.MockInfraClient, func()) {
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

func TestMockResourceNsxtPolicyPredefinedGatewayPolicyCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupPredefinedGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(predefinedGwPolicyAPIResponse(), nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(predefinedGwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())

		err := resourceNsxtPolicyPredefinedGatewayPolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, predefinedGwPolicyID, d.Id())
	})

	t.Run("Create fails when path is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"path": "",
		})

		err := resourceNsxtPolicyPredefinedGatewayPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyPredefinedGatewayPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupPredefinedGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(predefinedGwPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())
		d.SetId(predefinedGwPolicyID)

		err := resourceNsxtPolicyPredefinedGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, predefinedGwPolicyDescription, d.Get("description"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())

		err := resourceNsxtPolicyPredefinedGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails on API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(nsxModel.GatewayPolicy{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())
		d.SetId(predefinedGwPolicyID)

		err := resourceNsxtPolicyPredefinedGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyPredefinedGatewayPolicyUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupPredefinedGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(predefinedGwPolicyAPIResponse(), nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(predefinedGwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())
		d.SetId(predefinedGwPolicyID)

		err := resourceNsxtPolicyPredefinedGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())

		err := resourceNsxtPolicyPredefinedGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyPredefinedGatewayPolicyDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupPredefinedGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Delete (revert) success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(predefinedGwPolicyAPIResponse(), nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
		)

		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())
		d.SetId(predefinedGwPolicyID)

		err := resourceNsxtPolicyPredefinedGatewayPolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())

		err := resourceNsxtPolicyPredefinedGatewayPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
