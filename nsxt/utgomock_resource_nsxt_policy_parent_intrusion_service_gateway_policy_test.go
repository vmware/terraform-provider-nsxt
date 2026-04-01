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
	parentIdsGwPolicyID          = "parent-ids-gw-policy-id"
	parentIdsGwPolicyDisplayName = "test-parent-ids-gw-policy"
	parentIdsGwPolicyDescription = "Test Parent IDS Gateway Policy"
	parentIdsGwPolicyRevision    = int64(1)
	parentIdsGwPolicyDomain      = "default"
	parentIdsGwPolicyPath        = "/infra/domains/default/ids-gateway-policies/parent-ids-gw-policy-id"
	parentIdsGwPolicyCategory    = "LocalGatewayRules"
)

func parentIdsGwPolicyAPIResponse() nsxModel.IdsGatewayPolicy {
	resourceType := "IdsGatewayPolicy"
	seq := int64(0)
	stateful := true
	locked := false
	comments := ""
	return nsxModel.IdsGatewayPolicy{
		Id:             &parentIdsGwPolicyID,
		DisplayName:    &parentIdsGwPolicyDisplayName,
		Description:    &parentIdsGwPolicyDescription,
		Revision:       &parentIdsGwPolicyRevision,
		Path:           &parentIdsGwPolicyPath,
		Category:       &parentIdsGwPolicyCategory,
		ResourceType:   &resourceType,
		SequenceNumber: &seq,
		Stateful:       &stateful,
		Locked:         &locked,
		Comments:       &comments,
	}
}

func minimalParentIdsGwPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    parentIdsGwPolicyDisplayName,
		"description":     parentIdsGwPolicyDescription,
		"nsx_id":          parentIdsGwPolicyID,
		"domain":          parentIdsGwPolicyDomain,
		"locked":          false,
		"stateful":        true,
		"sequence_number": 0,
		"category":        parentIdsGwPolicyCategory,
	}
}

func setupParentIdsGwPolicyMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockIntrusionServiceGatewayPoliciesClient, *inframocks.MockInfraClient, func()) {
	t.Helper()
	mockSDK := domainmocks.NewMockIntrusionServiceGatewayPoliciesClient(ctrl)
	mockWrapper := &apidomains.IntrusionServiceGatewayPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliIntrusionServiceGatewayPoliciesClient
	cliIntrusionServiceGatewayPoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.IntrusionServiceGatewayPolicyClientContext {
		return mockWrapper
	}

	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)
	originalInfra := cliInfraClient
	cliInfraClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	return mockSDK, mockInfraSDK, func() {
		cliIntrusionServiceGatewayPoliciesClient = original
		cliInfraClient = originalInfra
	}
}

func TestMockResourceNsxtPolicyParentIntrusionServiceGatewayPolicyCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupParentIdsGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(parentIdsGwPolicyDomain, parentIdsGwPolicyID).Return(nsxModel.IdsGatewayPolicy{}, notFoundErr),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(parentIdsGwPolicyDomain, parentIdsGwPolicyID).Return(parentIdsGwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsGwPolicyData())

		err := resourceNsxtPolicyParentIntrusionServiceGatewayPolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, parentIdsGwPolicyID, d.Id())
		assert.Equal(t, parentIdsGwPolicyDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(parentIdsGwPolicyDomain, parentIdsGwPolicyID).Return(parentIdsGwPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsGwPolicyData())

		err := resourceNsxtPolicyParentIntrusionServiceGatewayPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupParentIdsGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(parentIdsGwPolicyDomain, parentIdsGwPolicyID).Return(parentIdsGwPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsGwPolicyData())
		d.SetId(parentIdsGwPolicyID)

		err := resourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, parentIdsGwPolicyDisplayName, d.Get("display_name"))
		assert.Equal(t, parentIdsGwPolicyDescription, d.Get("description"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(parentIdsGwPolicyDomain, parentIdsGwPolicyID).Return(nsxModel.IdsGatewayPolicy{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsGwPolicyData())
		d.SetId(parentIdsGwPolicyID)

		err := resourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsGwPolicyData())

		err := resourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyParentIntrusionServiceGatewayPolicyUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupParentIdsGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(parentIdsGwPolicyDomain, parentIdsGwPolicyID).Return(parentIdsGwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsGwPolicyData())
		d.SetId(parentIdsGwPolicyID)

		err := resourceNsxtPolicyParentIntrusionServiceGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsGwPolicyData())

		err := resourceNsxtPolicyParentIntrusionServiceGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyParentIntrusionServiceGatewayPolicyDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupParentIdsGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(parentIdsGwPolicyDomain, parentIdsGwPolicyID).Return(nil)

		res := resourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsGwPolicyData())
		d.SetId(parentIdsGwPolicyID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsGwPolicyData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
