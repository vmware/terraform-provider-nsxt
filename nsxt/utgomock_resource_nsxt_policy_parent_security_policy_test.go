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
	parentSecPolicyID          = "parent-sec-policy-id"
	parentSecPolicyDisplayName = "test-parent-sec-policy"
	parentSecPolicyDescription = "Test Parent Security Policy"
	parentSecPolicyRevision    = int64(1)
	parentSecPolicyDomain      = "default"
	parentSecPolicyPath        = "/infra/domains/default/security-policies/parent-sec-policy-id"
	parentSecPolicyCategory    = "Application"
)

func parentSecPolicyAPIResponse() nsxModel.SecurityPolicy {
	return nsxModel.SecurityPolicy{
		Id:          &parentSecPolicyID,
		DisplayName: &parentSecPolicyDisplayName,
		Description: &parentSecPolicyDescription,
		Revision:    &parentSecPolicyRevision,
		Path:        &parentSecPolicyPath,
		Category:    &parentSecPolicyCategory,
	}
}

func minimalParentSecPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    parentSecPolicyDisplayName,
		"description":     parentSecPolicyDescription,
		"nsx_id":          parentSecPolicyID,
		"domain":          parentSecPolicyDomain,
		"locked":          false,
		"stateful":        true,
		"sequence_number": 0,
		"category":        parentSecPolicyCategory,
	}
}

func setupParentSecPolicyMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockSecurityPoliciesClient, *inframocks.MockInfraClient, func()) {
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

func TestMockResourceNsxtPolicyParentSecurityPolicyCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupParentSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(parentSecPolicyDomain, parentSecPolicyID).Return(nsxModel.SecurityPolicy{}, notFoundErr),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(parentSecPolicyDomain, parentSecPolicyID).Return(parentSecPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyParentSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentSecPolicyData())

		err := resourceNsxtPolicyParentSecurityPolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, parentSecPolicyID, d.Id())
		assert.Equal(t, parentSecPolicyDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(parentSecPolicyDomain, parentSecPolicyID).Return(parentSecPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyParentSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentSecPolicyData())

		err := resourceNsxtPolicyParentSecurityPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyParentSecurityPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupParentSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(parentSecPolicyDomain, parentSecPolicyID).Return(parentSecPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyParentSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentSecPolicyData())
		d.SetId(parentSecPolicyID)

		err := resourceNsxtPolicyParentSecurityPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, parentSecPolicyDisplayName, d.Get("display_name"))
		assert.Equal(t, parentSecPolicyDescription, d.Get("description"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(parentSecPolicyDomain, parentSecPolicyID).Return(nsxModel.SecurityPolicy{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyParentSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentSecPolicyData())
		d.SetId(parentSecPolicyID)

		err := resourceNsxtPolicyParentSecurityPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyParentSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentSecPolicyData())

		err := resourceNsxtPolicyParentSecurityPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyParentSecurityPolicyUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupParentSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(parentSecPolicyDomain, parentSecPolicyID).Return(parentSecPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyParentSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentSecPolicyData())
		d.SetId(parentSecPolicyID)

		err := resourceNsxtPolicyParentSecurityPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyParentSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentSecPolicyData())

		err := resourceNsxtPolicyParentSecurityPolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyParentSecurityPolicyDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupParentSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(parentSecPolicyDomain, parentSecPolicyID).Return(nil)

		res := resourceNsxtPolicyParentSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentSecPolicyData())
		d.SetId(parentSecPolicyID)

		err := resourceNsxtPolicyParentSecurityPolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyParentSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentSecPolicyData())

		err := resourceNsxtPolicyParentSecurityPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
