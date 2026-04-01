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
	parentGwPolicyID          = "parent-gw-policy-id"
	parentGwPolicyDisplayName = "test-parent-gw-policy"
	parentGwPolicyDescription = "Test Parent Gateway Policy"
	parentGwPolicyRevision    = int64(1)
	parentGwPolicyDomain      = "default"
	parentGwPolicyPath        = "/infra/domains/default/gateway-policies/parent-gw-policy-id"
	parentGwPolicyCategory    = "LocalGatewayRules"
)

func parentGwPolicyAPIResponse() nsxModel.GatewayPolicy {
	return nsxModel.GatewayPolicy{
		Id:          &parentGwPolicyID,
		DisplayName: &parentGwPolicyDisplayName,
		Description: &parentGwPolicyDescription,
		Revision:    &parentGwPolicyRevision,
		Path:        &parentGwPolicyPath,
		Category:    &parentGwPolicyCategory,
	}
}

func minimalParentGwPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    parentGwPolicyDisplayName,
		"description":     parentGwPolicyDescription,
		"nsx_id":          parentGwPolicyID,
		"domain":          parentGwPolicyDomain,
		"locked":          false,
		"stateful":        true,
		"sequence_number": 0,
		"category":        parentGwPolicyCategory,
	}
}

func setupParentGwPolicyMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockGatewayPoliciesClient, *inframocks.MockInfraClient, func()) {
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

func TestMockResourceNsxtPolicyParentGatewayPolicyCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupParentGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(parentGwPolicyDomain, parentGwPolicyID).Return(nsxModel.GatewayPolicy{}, notFoundErr),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(parentGwPolicyDomain, parentGwPolicyID).Return(parentGwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyParentGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentGwPolicyData())

		err := resourceNsxtPolicyParentGatewayPolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, parentGwPolicyID, d.Id())
		assert.Equal(t, parentGwPolicyDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(parentGwPolicyDomain, parentGwPolicyID).Return(parentGwPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyParentGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentGwPolicyData())

		err := resourceNsxtPolicyParentGatewayPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyParentGatewayPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupParentGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(parentGwPolicyDomain, parentGwPolicyID).Return(parentGwPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyParentGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentGwPolicyData())
		d.SetId(parentGwPolicyID)

		err := resourceNsxtPolicyParentGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, parentGwPolicyDisplayName, d.Get("display_name"))
		assert.Equal(t, parentGwPolicyDescription, d.Get("description"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(parentGwPolicyDomain, parentGwPolicyID).Return(nsxModel.GatewayPolicy{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyParentGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentGwPolicyData())
		d.SetId(parentGwPolicyID)

		err := resourceNsxtPolicyParentGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyParentGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentGwPolicyData())

		err := resourceNsxtPolicyParentGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyParentGatewayPolicyUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupParentGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(parentGwPolicyDomain, parentGwPolicyID).Return(parentGwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyParentGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentGwPolicyData())
		d.SetId(parentGwPolicyID)

		err := resourceNsxtPolicyParentGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyParentGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentGwPolicyData())

		err := resourceNsxtPolicyParentGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyParentGatewayPolicyDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupParentGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(parentGwPolicyDomain, parentGwPolicyID).Return(nil)

		res := resourceNsxtPolicyParentGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentGwPolicyData())
		d.SetId(parentGwPolicyID)

		err := resourceNsxtPolicyGatewayPolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyParentGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentGwPolicyData())

		err := resourceNsxtPolicyGatewayPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
