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
	domainsapi "github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	domainmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains"
)

var (
	parentIdsPolicyID          = "parent-ids-policy-id"
	parentIdsPolicyDisplayName = "Test Parent IDS Policy"
	parentIdsPolicyDescription = "Test Parent Intrusion Service Policy"
	parentIdsPolicyDomain      = "default"
	parentIdsPolicyPath        = "/infra/domains/default/intrusion-service-policies/parent-ids-policy-id"
	parentIdsPolicySeqNum      = int64(0)
	parentIdsPolicyStateful    = true
	parentIdsPolicyLocked      = false
	parentIdsPolicyComments    = ""
	// Updated values
	parentIdsPolicyUpdatedSeqNum      = int64(1)
	parentIdsPolicyUpdatedLocked      = true
	parentIdsPolicyUpdatedComments    = "Update Policy"
	parentIdsPolicyUpdatedDisplayName = "Test Parent IDS Policy Update"
	parentIdsPolicyUpdatedDescription = "Test Parent Intrusion Service Policy Update"
)

func parentIdsPolicyAPIResponse() nsxModel.IdsSecurityPolicy {
	resourceType := "IdsSecurityPolicy"
	return nsxModel.IdsSecurityPolicy{
		Id:             &parentIdsPolicyID,
		DisplayName:    &parentIdsPolicyDisplayName,
		Description:    &parentIdsPolicyDescription,
		Path:           &parentIdsPolicyPath,
		ResourceType:   &resourceType,
		SequenceNumber: &parentIdsPolicySeqNum,
		Stateful:       &parentIdsPolicyStateful,
		Locked:         &parentIdsPolicyLocked,
		Comments:       &parentIdsPolicyComments,
		Rules:          nil, // Parent policy has no embedded rules
	}
}

func parentIdsPolicyAPIResponseUpdated() nsxModel.IdsSecurityPolicy {
	resourceType := "IdsSecurityPolicy"
	return nsxModel.IdsSecurityPolicy{
		Id:             &parentIdsPolicyID,
		DisplayName:    &parentIdsPolicyUpdatedDisplayName,
		Description:    &parentIdsPolicyUpdatedDescription,
		Path:           &parentIdsPolicyPath,
		ResourceType:   &resourceType,
		SequenceNumber: &parentIdsPolicyUpdatedSeqNum,
		Stateful:       &parentIdsPolicyStateful,
		Locked:         &parentIdsPolicyUpdatedLocked,
		Comments:       &parentIdsPolicyUpdatedComments,
		Rules:          nil, // Parent policy has no embedded rules
	}
}

func minimalParentIdsPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    parentIdsPolicyDisplayName,
		"description":     parentIdsPolicyDescription,
		"nsx_id":          parentIdsPolicyID,
		"domain":          parentIdsPolicyDomain,
		"sequence_number": int(parentIdsPolicySeqNum),
		"stateful":        parentIdsPolicyStateful,
		"locked":          parentIdsPolicyLocked,
		"comments":        parentIdsPolicyComments,
	}
}

func minimalParentIdsPolicyDataUpdate() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    parentIdsPolicyDisplayName,
		"description":     parentIdsPolicyDescription,
		"nsx_id":          parentIdsPolicyID,
		"domain":          parentIdsPolicyDomain,
		"sequence_number": int(parentIdsPolicySeqNum),
		"locked":          parentIdsPolicyLocked,
		"comments":        parentIdsPolicyComments,
	}
}

func setupParentIdsPolicyMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockIntrusionServicePoliciesClient, *inframocks.MockInfraClient, func()) {
	mockSDK := domainmocks.NewMockIntrusionServicePoliciesClient(ctrl)
	mockWrapper := &domainsapi.IdsSecurityPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)

	originalISP := cliIntrusionServicePoliciesClient
	originalInfra := cliInfraClient
	cliIntrusionServicePoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *domainsapi.IdsSecurityPolicyClientContext {
		return mockWrapper
	}
	cliInfraClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	return mockSDK, mockInfraSDK, func() {
		cliIntrusionServicePoliciesClient = originalISP
		cliInfraClient = originalInfra
	}
}

func TestMockResourceNsxtPolicyParentIntrusionServicePolicyCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupParentIdsPolicyMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(parentIdsPolicyDomain, parentIdsPolicyID).Return(nsxModel.IdsSecurityPolicy{}, notFoundErr),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(parentIdsPolicyDomain, parentIdsPolicyID).Return(parentIdsPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsPolicyData())

		err := resourceNsxtPolicyParentIntrusionServicePolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(parentIdsPolicyDomain, parentIdsPolicyID).Return(parentIdsPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsPolicyData())

		err := resourceNsxtPolicyParentIntrusionServicePolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyParentIntrusionServicePolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupParentIdsPolicyMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(parentIdsPolicyDomain, parentIdsPolicyID).Return(parentIdsPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsPolicyData())
		d.SetId(parentIdsPolicyID)

		err := resourceNsxtPolicyParentIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, parentIdsPolicyDomain, d.Get("domain"))
		assert.Equal(t, parentIdsPolicyDisplayName, d.Get("display_name"))
		assert.Equal(t, parentIdsPolicyDescription, d.Get("description"))
		assert.Equal(t, int(parentIdsPolicySeqNum), d.Get("sequence_number"))
		assert.Equal(t, parentIdsPolicyStateful, d.Get("stateful"))
		assert.Equal(t, parentIdsPolicyLocked, d.Get("locked"))
	})

	t.Run("Read not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(parentIdsPolicyDomain, parentIdsPolicyID).Return(nsxModel.IdsSecurityPolicy{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsPolicyData())
		d.SetId(parentIdsPolicyID)

		err := resourceNsxtPolicyParentIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsPolicyData())

		err := resourceNsxtPolicyParentIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyParentIntrusionServicePolicyUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupParentIdsPolicyMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(parentIdsPolicyDomain, parentIdsPolicyID).Return(parentIdsPolicyAPIResponseUpdated(), nil),
		)

		updatedData := minimalParentIdsPolicyDataUpdate()
		updatedData["sequence_number"] = int(parentIdsPolicyUpdatedSeqNum)
		updatedData["locked"] = parentIdsPolicyUpdatedLocked
		updatedData["comments"] = parentIdsPolicyUpdatedComments
		updatedData["display_name"] = parentIdsPolicyUpdatedDisplayName
		updatedData["description"] = parentIdsPolicyUpdatedDescription

		res := resourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, updatedData)
		d.SetId(parentIdsPolicyID)

		err := resourceNsxtPolicyParentIntrusionServicePolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, parentIdsPolicyDomain, d.Get("domain"))
		assert.Equal(t, parentIdsPolicyUpdatedDisplayName, d.Get("display_name"))
		assert.Equal(t, parentIdsPolicyUpdatedDescription, d.Get("description"))
		assert.Equal(t, int(parentIdsPolicyUpdatedSeqNum), d.Get("sequence_number"))
		assert.Equal(t, parentIdsPolicyStateful, d.Get("stateful"))
		assert.Equal(t, parentIdsPolicyUpdatedLocked, d.Get("locked"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsPolicyDataUpdate())
		// Don't set ID

		err := resourceNsxtPolicyParentIntrusionServicePolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyParentIntrusionServicePolicyDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupParentIdsPolicyMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(parentIdsPolicyDomain, parentIdsPolicyID).Return(nil)

		res := resourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsPolicyData())
		d.SetId(parentIdsPolicyID)

		err := resourceNsxtPolicyIntrusionServicePolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalParentIdsPolicyData())
		// Don't set ID

		err := resourceNsxtPolicyIntrusionServicePolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
