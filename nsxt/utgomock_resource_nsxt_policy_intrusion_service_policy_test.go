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
	ispID          = "isp-1"
	ispDisplayName = "Test IDS Policy"
	ispDescription = "test intrusion service policy"
	ispRevision    = int64(1)
	ispDomain      = "default"
	ispPath        = "/infra/domains/default/intrusion-service-policies/isp-1"
)

func ispAPIResponse() nsxModel.IdsSecurityPolicy {
	stateful := true
	seqNum := int64(1)
	return nsxModel.IdsSecurityPolicy{
		Id:             &ispID,
		DisplayName:    &ispDisplayName,
		Description:    &ispDescription,
		Path:           &ispPath,
		Stateful:       &stateful,
		SequenceNumber: &seqNum,
	}
}

func minimalIspData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": ispDisplayName,
		"description":  ispDescription,
		"nsx_id":       ispID,
		"domain":       ispDomain,
		"stateful":     true,
	}
}

func minimalIspDataUpdate() map[string]interface{} {
	return map[string]interface{}{
		"display_name": ispDisplayName,
		"description":  ispDescription,
		"nsx_id":       ispID,
		"domain":       ispDomain,
	}
}

func setupIspMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockIntrusionServicePoliciesClient, *inframocks.MockInfraClient, func()) {
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

func TestMockResourceNsxtPolicyIntrusionServicePolicyCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupIspMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(ispDomain, ispID).Return(nsxModel.IdsSecurityPolicy{}, notFoundErr),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(ispDomain, ispID).Return(ispAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIspData())

		err := resourceNsxtPolicyIntrusionServicePolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ispID, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(ispDomain, ispID).Return(ispAPIResponse(), nil)

		res := resourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIspData())

		err := resourceNsxtPolicyIntrusionServicePolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServicePolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupIspMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(ispDomain, ispID).Return(ispAPIResponse(), nil)

		res := resourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIspData())
		d.SetId(ispID)

		err := resourceNsxtPolicyIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ispDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(ispDomain, ispID).Return(nsxModel.IdsSecurityPolicy{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIspData())
		d.SetId(ispID)

		err := resourceNsxtPolicyIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIspData())

		err := resourceNsxtPolicyIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServicePolicyUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupIspMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(ispDomain, ispID).Return(ispAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIspDataUpdate())
		d.SetId(ispID)

		err := resourceNsxtPolicyIntrusionServicePolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIspDataUpdate())

		err := resourceNsxtPolicyIntrusionServicePolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServicePolicyDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupIspMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(ispDomain, ispID).Return(nil)

		res := resourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIspData())
		d.SetId(ispID)

		err := resourceNsxtPolicyIntrusionServicePolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIspData())

		err := resourceNsxtPolicyIntrusionServicePolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
