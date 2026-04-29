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
	parentIdsGwPolicyDisplayName = "Test Parent IDS Gateway Policy"
	parentIdsGwPolicyDescription = "Test Parent Intrusion Service Gateway Policy"
	parentIdsGwPolicyDomain      = "default"
	parentIdsGwPolicyPath        = "/infra/domains/default/intrusion-service-gateway-policies/parent-ids-gw-policy-id"
	parentIdsGwPolicyCategory    = "LocalGatewayRules"
	parentIdsGwPolicySeqNum      = int64(0)
	parentIdsGwPolicyStateful    = true
	parentIdsGwPolicyLocked      = false
	parentIdsGwPolicyComments    = ""

	// Updated values
	parentIdsGwPolicyUpdatedSeqNum = int64(1)
	parentIdsGwPolicyUpdatedLocked = true
	// parentIdsGwPolicyUpdatedTcpStrict removed - not used in unit tests per user request
	parentIdsGwPolicyUpdatedComments    = "Update Policy"
	parentIdsGwPolicyUpdatedDisplayName = "Test Parent IDS Gateway Policy Update"
	parentIdsGwPolicyUpdatedDescription = "Test Parent Intrusion Service Gateway Policy Update"
)

func parentIdsGwPolicyAPIResponse() nsxModel.IdsGatewayPolicy {
	resourceType := "IdsGatewayPolicy"
	return nsxModel.IdsGatewayPolicy{
		Id:             &parentIdsGwPolicyID,
		DisplayName:    &parentIdsGwPolicyDisplayName,
		Description:    &parentIdsGwPolicyDescription,
		Path:           &parentIdsGwPolicyPath,
		Category:       &parentIdsGwPolicyCategory,
		ResourceType:   &resourceType,
		SequenceNumber: &parentIdsGwPolicySeqNum,
		Stateful:       &parentIdsGwPolicyStateful,
		Locked:         &parentIdsGwPolicyLocked,
		Comments:       &parentIdsGwPolicyComments,
		Rules:          nil,
	}
}

func parentIdsGwPolicyAPIResponseUpdated() nsxModel.IdsGatewayPolicy {
	resourceType := "IdsGatewayPolicy"
	return nsxModel.IdsGatewayPolicy{
		Id:             &parentIdsGwPolicyID,
		DisplayName:    &parentIdsGwPolicyUpdatedDisplayName,
		Description:    &parentIdsGwPolicyUpdatedDescription,
		Path:           &parentIdsGwPolicyPath,
		Category:       &parentIdsGwPolicyCategory,
		ResourceType:   &resourceType,
		SequenceNumber: &parentIdsGwPolicyUpdatedSeqNum,
		Stateful:       &parentIdsGwPolicyStateful,
		Locked:         &parentIdsGwPolicyUpdatedLocked,
		Comments:       &parentIdsGwPolicyUpdatedComments,
		Rules:          nil,
	}
}

func minimalParentIdsGwPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"nsx_id":          parentIdsGwPolicyID,
		"display_name":    parentIdsGwPolicyDisplayName,
		"description":     parentIdsGwPolicyDescription,
		"domain":          parentIdsGwPolicyDomain,
		"category":        parentIdsGwPolicyCategory,
		"sequence_number": 0,
		"stateful":        true,
		"locked":          false,
	}
}

func minimalParentIdsGwPolicyDataUpdate() map[string]interface{} {
	return map[string]interface{}{
		"nsx_id":          parentIdsGwPolicyID,
		"display_name":    parentIdsGwPolicyDisplayName,
		"description":     parentIdsGwPolicyDescription,
		"domain":          parentIdsGwPolicyDomain,
		"sequence_number": 0,
		"locked":          false,
	}
}

func setupParentIdsGwPolicyMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockIntrusionServiceGatewayPoliciesClient, *inframocks.MockInfraClient, func()) {
	t.Helper()
	mockSDK := domainmocks.NewMockIntrusionServiceGatewayPoliciesClient(ctrl)
	mockWrapper := &apidomains.IntrusionServiceGatewayPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalGWP := cliIntrusionServiceGatewayPoliciesClient
	cliIntrusionServiceGatewayPoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.IntrusionServiceGatewayPolicyClientContext {
		return mockWrapper
	}

	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)
	originalInfra := cliInfraClient
	cliInfraClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	return mockSDK, mockInfraSDK, func() {
		cliIntrusionServiceGatewayPoliciesClient = originalGWP
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
		assert.NotEmpty(t, d.Id())
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
		assert.Equal(t, parentIdsGwPolicyStateful, d.Get("stateful"))
		assert.Equal(t, parentIdsGwPolicyLocked, d.Get("locked"))
		assert.Equal(t, int(parentIdsGwPolicySeqNum), d.Get("sequence_number"))
		assert.Equal(t, parentIdsGwPolicyCategory, d.Get("category"))
		assert.Equal(t, parentIdsGwPolicyDisplayName, d.Get("display_name"))
		assert.Equal(t, parentIdsGwPolicyDescription, d.Get("description"))
		assert.Equal(t, parentIdsGwPolicyDomain, d.Get("domain"))
	})

	t.Run("Read not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(parentIdsGwPolicyDomain, parentIdsGwPolicyID).Return(nsxModel.IdsGatewayPolicy{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsGwPolicyData())
		d.SetId(parentIdsGwPolicyID)

		err := resourceNsxtPolicyParentIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
		assert.Equal(t, "", d.Id())
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
			mockSDK.EXPECT().Get(parentIdsGwPolicyDomain, parentIdsGwPolicyID).Return(parentIdsGwPolicyAPIResponseUpdated(), nil),
		)

		updatedData := minimalParentIdsGwPolicyDataUpdate()
		updatedData["locked"] = parentIdsGwPolicyUpdatedLocked
		updatedData["sequence_number"] = int(parentIdsGwPolicyUpdatedSeqNum)
		updatedData["comments"] = parentIdsGwPolicyUpdatedComments
		updatedData["display_name"] = parentIdsGwPolicyUpdatedDisplayName
		updatedData["description"] = parentIdsGwPolicyUpdatedDescription

		res := resourceNsxtPolicyParentIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, updatedData)
		d.SetId(parentIdsGwPolicyID)

		err := resourceNsxtPolicyParentIntrusionServiceGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, parentIdsGwPolicyStateful, d.Get("stateful"))
		assert.Equal(t, parentIdsGwPolicyUpdatedLocked, d.Get("locked"))
		assert.Equal(t, int(parentIdsGwPolicyUpdatedSeqNum), d.Get("sequence_number"))
		assert.Equal(t, parentIdsGwPolicyCategory, d.Get("category"))
		assert.Equal(t, parentIdsGwPolicyUpdatedDisplayName, d.Get("display_name"))
		assert.Equal(t, parentIdsGwPolicyUpdatedDescription, d.Get("description"))
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
