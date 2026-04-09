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
	predefinedSecPolicyID          = "default"
	predefinedSecPolicyDomain      = "default"
	predefinedSecPolicyPath        = "/infra/domains/default/security-policies/default"
	predefinedSecPolicyDescription = "Predefined Security Policy"
	predefinedSecPolicyRevision    = int64(1)
	predefinedSecPolicyCategory    = "Default"
)

func predefinedSecPolicyAPIResponse() nsxModel.SecurityPolicy {
	return nsxModel.SecurityPolicy{
		Id:          &predefinedSecPolicyID,
		Description: &predefinedSecPolicyDescription,
		Revision:    &predefinedSecPolicyRevision,
		Path:        &predefinedSecPolicyPath,
		Category:    &predefinedSecPolicyCategory,
	}
}

func minimalPredefinedSecPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"path": predefinedSecPolicyPath,
	}
}

func setupPredefinedSecPolicyMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockSecurityPoliciesClient, *inframocks.MockInfraClient, func()) {
	mockSDK := domainmocks.NewMockSecurityPoliciesClient(ctrl)
	mockWrapper := &apidomains.SecurityPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalSecPolicy := cliSecurityPoliciesClient
	cliSecurityPoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.SecurityPolicyClientContext {
		return mockWrapper
	}

	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)
	originalInfra := cliInfraClient
	cliInfraClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	return mockSDK, mockInfraSDK, func() {
		cliSecurityPoliciesClient = originalSecPolicy
		cliInfraClient = originalInfra
	}
}

func TestMockResourceNsxtPolicyPredefinedSecurityPolicyCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupPredefinedSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())

		err := resourceNsxtPolicyPredefinedSecurityPolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, predefinedSecPolicyID, d.Id())
	})

	t.Run("Create fails when path is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"path": "",
		})

		err := resourceNsxtPolicyPredefinedSecurityPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyPredefinedSecurityPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupPredefinedSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())
		d.SetId(predefinedSecPolicyID)

		err := resourceNsxtPolicyPredefinedSecurityPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, predefinedSecPolicyDescription, d.Get("description"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())

		err := resourceNsxtPolicyPredefinedSecurityPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails on API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(nsxModel.SecurityPolicy{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())
		d.SetId(predefinedSecPolicyID)

		err := resourceNsxtPolicyPredefinedSecurityPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyPredefinedSecurityPolicyUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupPredefinedSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())
		d.SetId(predefinedSecPolicyID)

		err := resourceNsxtPolicyPredefinedSecurityPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())

		err := resourceNsxtPolicyPredefinedSecurityPolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyPredefinedSecurityPolicyDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupPredefinedSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Delete (revert) success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
		)

		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())
		d.SetId(predefinedSecPolicyID)

		err := resourceNsxtPolicyPredefinedSecurityPolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())

		err := resourceNsxtPolicyPredefinedSecurityPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
