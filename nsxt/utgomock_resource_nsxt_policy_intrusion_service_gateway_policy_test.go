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
	idsGwPolicyID          = "ids-gw-policy-1"
	idsGwPolicyDisplayName = "Test IDS Gateway Policy"
	idsGwPolicyDescription = "Test IDS gateway policy"
	idsGwPolicyRevision    = int64(1)
	idsGwPolicyDomain      = "default"
	idsGwPolicyPath        = "/infra/domains/default/ids-gateway-policies/ids-gw-policy-1"
	idsGwPolicyCategory    = "LocalGatewayRules"
)

func idsGwPolicyAPIResponse() nsxModel.IdsGatewayPolicy {
	resourceType := "IdsGatewayPolicy"
	seq := int64(0)
	stateful := true
	locked := false
	comments := ""
	return nsxModel.IdsGatewayPolicy{
		Id:             &idsGwPolicyID,
		DisplayName:    &idsGwPolicyDisplayName,
		Description:    &idsGwPolicyDescription,
		Revision:       &idsGwPolicyRevision,
		Path:           &idsGwPolicyPath,
		Category:       &idsGwPolicyCategory,
		ResourceType:   &resourceType,
		SequenceNumber: &seq,
		Stateful:       &stateful,
		Locked:         &locked,
		Comments:       &comments,
		Rules:          nil,
	}
}

func minimalIdsGwPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    idsGwPolicyDisplayName,
		"description":     idsGwPolicyDescription,
		"nsx_id":          idsGwPolicyID,
		"domain":          idsGwPolicyDomain,
		"locked":          false,
		"stateful":        true,
		"sequence_number": 0,
		"category":        idsGwPolicyCategory,
	}
}

func setupIdsGwPolicyMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockIntrusionServiceGatewayPoliciesClient, *inframocks.MockInfraClient, func()) {
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

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupIdsGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(idsGwPolicyDomain, idsGwPolicyID).Return(nsxModel.IdsGatewayPolicy{}, notFoundErr),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(idsGwPolicyDomain, idsGwPolicyID).Return(idsGwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsGwPolicyID, d.Id())
		assert.Equal(t, idsGwPolicyDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwPolicyDomain, idsGwPolicyID).Return(idsGwPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupIdsGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwPolicyDomain, idsGwPolicyID).Return(idsGwPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())
		d.SetId(idsGwPolicyID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsGwPolicyDisplayName, d.Get("display_name"))
		assert.Equal(t, idsGwPolicyCategory, d.Get("category"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwPolicyDomain, idsGwPolicyID).Return(nsxModel.IdsGatewayPolicy{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())
		d.SetId(idsGwPolicyID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupIdsGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(idsGwPolicyDomain, idsGwPolicyID).Return(idsGwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())
		d.SetId(idsGwPolicyID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupIdsGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(idsGwPolicyDomain, idsGwPolicyID).Return(nil)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())
		d.SetId(idsGwPolicyID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwPolicyData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
