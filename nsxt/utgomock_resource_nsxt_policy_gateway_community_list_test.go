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

	tier0sapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t0mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
)

var (
	commListID          = "comm-list-1"
	commListDisplayName = "Test Community List"
	commListDescription = "Test community list"
	commListRevision    = int64(1)
	commListGwPath      = "/infra/tier-0s/t0-gw-1"
	commListGwID        = "t0-gw-1"
	commListPath        = "/infra/tier-0s/t0-gw-1/community-lists/comm-list-1"
	commListCommunity   = "65001:100"
)

func commListAPIResponse() nsxModel.CommunityList {
	return nsxModel.CommunityList{
		Id:          &commListID,
		DisplayName: &commListDisplayName,
		Description: &commListDescription,
		Revision:    &commListRevision,
		Path:        &commListPath,
		Communities: []string{commListCommunity},
	}
}

func minimalCommListData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": commListDisplayName,
		"description":  commListDescription,
		"nsx_id":       commListID,
		"gateway_path": commListGwPath,
		"communities":  []interface{}{commListCommunity},
	}
}

func setupCommListMock(t *testing.T, ctrl *gomock.Controller) (*t0mocks.MockCommunityListsClient, func()) {
	mockSDK := t0mocks.NewMockCommunityListsClient(ctrl)
	mockWrapper := &tier0sapi.CommunityListClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliCommunityListsClient
	cliCommunityListsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier0sapi.CommunityListClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliCommunityListsClient = original }
}

func TestMockResourceNsxtPolicyGatewayCommunityListCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCommListMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(commListGwID, commListID).Return(nsxModel.CommunityList{}, notFoundErr),
			mockSDK.EXPECT().Patch(commListGwID, commListID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(commListGwID, commListID).Return(commListAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCommListData())

		err := resourceNsxtPolicyGatewayCommunityListCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, commListID, d.Id())
		assert.Equal(t, commListDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(commListGwID, commListID).Return(commListAPIResponse(), nil)

		res := resourceNsxtPolicyGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCommListData())

		err := resourceNsxtPolicyGatewayCommunityListCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Create fails for non-Tier0 gateway path", func(t *testing.T) {
		data := minimalCommListData()
		data["gateway_path"] = "/infra/tier-1s/t1-gw-1"

		res := resourceNsxtPolicyGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyGatewayCommunityListCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Tier0")
	})
}

func TestMockResourceNsxtPolicyGatewayCommunityListRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCommListMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(commListGwID, commListID).Return(commListAPIResponse(), nil)

		res := resourceNsxtPolicyGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCommListData())
		d.SetId(commListID)

		err := resourceNsxtPolicyGatewayCommunityListRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, commListDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(commListGwID, commListID).Return(nsxModel.CommunityList{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCommListData())
		d.SetId(commListID)

		err := resourceNsxtPolicyGatewayCommunityListRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCommListData())

		err := resourceNsxtPolicyGatewayCommunityListRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayCommunityListUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCommListMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(commListGwID, commListID, gomock.Any()).Return(commListAPIResponse(), nil),
			mockSDK.EXPECT().Get(commListGwID, commListID).Return(commListAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCommListData())
		d.SetId(commListID)

		err := resourceNsxtPolicyGatewayCommunityListUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCommListData())

		err := resourceNsxtPolicyGatewayCommunityListUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayCommunityListDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupCommListMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(commListGwID, commListID).Return(nil)

		res := resourceNsxtPolicyGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCommListData())
		d.SetId(commListID)

		err := resourceNsxtPolicyGatewayCommunityListDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCommListData())

		err := resourceNsxtPolicyGatewayCommunityListDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
