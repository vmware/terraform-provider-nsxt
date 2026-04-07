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
	secPolicyID          = "sec-policy-id"
	secPolicyDisplayName = "test-sec-policy"
	secPolicyDescription = "Test Security Policy"
	secPolicyRevision    = int64(1)
	secPolicyDomain      = "default"
	secPolicyPath        = "/infra/domains/default/security-policies/sec-policy-id"
	secPolicyCategory    = "Application"
)

func secPolicyAPIResponse() nsxModel.SecurityPolicy {
	return nsxModel.SecurityPolicy{
		Id:          &secPolicyID,
		DisplayName: &secPolicyDisplayName,
		Description: &secPolicyDescription,
		Revision:    &secPolicyRevision,
		Path:        &secPolicyPath,
		Category:    &secPolicyCategory,
	}
}

func minimalSecPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    secPolicyDisplayName,
		"description":     secPolicyDescription,
		"nsx_id":          secPolicyID,
		"domain":          secPolicyDomain,
		"locked":          false,
		"stateful":        true,
		"sequence_number": 0,
		"category":        secPolicyCategory,
	}
}

func setupSecPolicyMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockSecurityPoliciesClient, *inframocks.MockInfraClient, func()) {
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

func TestMockResourceNsxtPolicySecurityPolicyCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(secPolicyDomain, secPolicyID).Return(nsxModel.SecurityPolicy{}, notFoundErr),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(secPolicyDomain, secPolicyID).Return(secPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSecPolicyData())

		err := resourceNsxtPolicySecurityPolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, secPolicyID, d.Id())
		assert.Equal(t, secPolicyDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(secPolicyDomain, secPolicyID).Return(secPolicyAPIResponse(), nil)

		res := resourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSecPolicyData())

		err := resourceNsxtPolicySecurityPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicySecurityPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(secPolicyDomain, secPolicyID).Return(secPolicyAPIResponse(), nil)

		res := resourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSecPolicyData())
		d.SetId(secPolicyID)

		err := resourceNsxtPolicySecurityPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, secPolicyDisplayName, d.Get("display_name"))
		assert.Equal(t, secPolicyDescription, d.Get("description"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(secPolicyDomain, secPolicyID).Return(nsxModel.SecurityPolicy{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSecPolicyData())
		d.SetId(secPolicyID)

		err := resourceNsxtPolicySecurityPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSecPolicyData())

		err := resourceNsxtPolicySecurityPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySecurityPolicyUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(secPolicyDomain, secPolicyID).Return(secPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSecPolicyData())
		d.SetId(secPolicyID)

		err := resourceNsxtPolicySecurityPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSecPolicyData())

		err := resourceNsxtPolicySecurityPolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySecurityPolicyDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(secPolicyDomain, secPolicyID).Return(nil)

		res := resourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSecPolicyData())
		d.SetId(secPolicyID)

		err := resourceNsxtPolicySecurityPolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSecPolicyData())

		err := resourceNsxtPolicySecurityPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
