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

	transitgateways "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	tgwmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways"
)

var (
	tgwCLID          = "community-list-1"
	tgwCLDisplayName = "Test TGW Community List"
	tgwCLDescription = "Test transit gateway community list"
	tgwCLRevision    = int64(1)
	tgwCLParentPath  = "/orgs/default/projects/project1/transit-gateways/tgw1"
	tgwCLOrgID       = "default"
	tgwCLProjectID   = "project1"
	tgwCLTGWID       = "tgw1"
	tgwCLPath        = "/orgs/default/projects/project1/transit-gateways/tgw1/community-lists/community-list-1"
	tgwCLCommunity1  = "65000:100"
	tgwCLCommunity2  = "65000:200"
)

func tgwCommunityListAPIResponse() nsxModel.CommunityList {
	return nsxModel.CommunityList{
		Id:          &tgwCLID,
		DisplayName: &tgwCLDisplayName,
		Description: &tgwCLDescription,
		Revision:    &tgwCLRevision,
		Path:        &tgwCLPath,
		Communities: []string{tgwCLCommunity1, tgwCLCommunity2},
	}
}

func minimalTGWCommunityListData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": tgwCLDisplayName,
		"description":  tgwCLDescription,
		"nsx_id":       tgwCLID,
		"parent_path":  tgwCLParentPath,
		"communities":  []interface{}{tgwCLCommunity1, tgwCLCommunity2},
	}
}

func setupTGWCommunityListMock(t *testing.T, ctrl *gomock.Controller) (*tgwmocks.MockCommunityListsClient, func()) {
	mockSDK := tgwmocks.NewMockCommunityListsClient(ctrl)
	mockWrapper := &transitgateways.TGWCommunityListClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwCLProjectID,
	}

	original := cliTGWCommunityListsClient
	cliTGWCommunityListsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *transitgateways.TGWCommunityListClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTGWCommunityListsClient = original }
}

func TestMockResourceNsxtPolicyTransitGatewayCommunityListCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWCommunityListMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwCLOrgID, tgwCLProjectID, tgwCLTGWID, tgwCLID).Return(nsxModel.CommunityList{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwCLOrgID, tgwCLProjectID, tgwCLTGWID, tgwCLID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwCLOrgID, tgwCLProjectID, tgwCLTGWID, tgwCLID).Return(tgwCommunityListAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWCommunityListData())

		err := resourceNsxtPolicyTransitGatewayCommunityListCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwCLID, d.Id())
		assert.Equal(t, tgwCLDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when parent_path is missing", func(t *testing.T) {
		data := minimalTGWCommunityListData()
		data["parent_path"] = ""

		res := resourceNsxtPolicyTransitGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayCommunityListCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayCommunityListRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWCommunityListMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwCLOrgID, tgwCLProjectID, tgwCLTGWID, tgwCLID).Return(tgwCommunityListAPIResponse(), nil)

		res := resourceNsxtPolicyTransitGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWCommunityListData())
		d.SetId(tgwCLID)

		err := resourceNsxtPolicyTransitGatewayCommunityListRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwCLDisplayName, d.Get("display_name"))
		assert.Equal(t, tgwCLPath, d.Get("path"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwCLOrgID, tgwCLProjectID, tgwCLTGWID, tgwCLID).Return(nsxModel.CommunityList{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTransitGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWCommunityListData())
		d.SetId(tgwCLID)

		err := resourceNsxtPolicyTransitGatewayCommunityListRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWCommunityListData())

		err := resourceNsxtPolicyTransitGatewayCommunityListRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayCommunityListUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWCommunityListMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(tgwCLOrgID, tgwCLProjectID, tgwCLTGWID, tgwCLID, gomock.Any()).Return(tgwCommunityListAPIResponse(), nil),
			mockSDK.EXPECT().Get(tgwCLOrgID, tgwCLProjectID, tgwCLTGWID, tgwCLID).Return(tgwCommunityListAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWCommunityListData())
		d.SetId(tgwCLID)

		err := resourceNsxtPolicyTransitGatewayCommunityListUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwCLID, d.Id())
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWCommunityListData())

		err := resourceNsxtPolicyTransitGatewayCommunityListUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayCommunityListDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWCommunityListMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(tgwCLOrgID, tgwCLProjectID, tgwCLTGWID, tgwCLID).Return(nil)

		res := resourceNsxtPolicyTransitGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWCommunityListData())
		d.SetId(tgwCLID)

		err := resourceNsxtPolicyTransitGatewayCommunityListDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWCommunityListData())

		err := resourceNsxtPolicyTransitGatewayCommunityListDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
