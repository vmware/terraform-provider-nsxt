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
	prefixListID          = "prefix-list-1"
	prefixListDisplayName = "Test Prefix List"
	prefixListDescription = "Test prefix list"
	prefixListRevision    = int64(1)
	prefixListGwPath      = "/infra/tier-0s/t0-pl-gw-1"
	prefixListGwID        = "t0-pl-gw-1"
	prefixListPath        = "/infra/tier-0s/t0-pl-gw-1/prefix-lists/prefix-list-1"
	prefixListNetwork     = "192.168.0.0/24"
	prefixListAction      = nsxModel.PrefixEntry_ACTION_PERMIT
)

func prefixListAPIResponse() nsxModel.PrefixList {
	return nsxModel.PrefixList{
		Id:          &prefixListID,
		DisplayName: &prefixListDisplayName,
		Description: &prefixListDescription,
		Revision:    &prefixListRevision,
		Path:        &prefixListPath,
		Prefixes: []nsxModel.PrefixEntry{
			{
				Action:  &prefixListAction,
				Network: &prefixListNetwork,
			},
		},
	}
}

func minimalPrefixListData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": prefixListDisplayName,
		"description":  prefixListDescription,
		"nsx_id":       prefixListID,
		"gateway_path": prefixListGwPath,
		"prefix": []interface{}{
			map[string]interface{}{
				"action":  prefixListAction,
				"network": prefixListNetwork,
				"ge":      0,
				"le":      0,
			},
		},
	}
}

func setupPrefixListMock(t *testing.T, ctrl *gomock.Controller) (*t0mocks.MockPrefixListsClient, func()) {
	mockSDK := t0mocks.NewMockPrefixListsClient(ctrl)
	mockWrapper := &tier0sapi.PrefixListClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliPrefixListsClient
	cliPrefixListsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier0sapi.PrefixListClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliPrefixListsClient = original }
}

func TestMockResourceNsxtPolicyGatewayPrefixListCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupPrefixListMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(prefixListGwID, prefixListID).Return(nsxModel.PrefixList{}, notFoundErr),
			mockSDK.EXPECT().Patch(prefixListGwID, prefixListID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(prefixListGwID, prefixListID).Return(prefixListAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrefixListData())

		err := resourceNsxtPolicyGatewayPrefixListCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, prefixListID, d.Id())
		assert.Equal(t, prefixListDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(prefixListGwID, prefixListID).Return(prefixListAPIResponse(), nil)

		res := resourceNsxtPolicyGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrefixListData())

		err := resourceNsxtPolicyGatewayPrefixListCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayPrefixListRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupPrefixListMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(prefixListGwID, prefixListID).Return(prefixListAPIResponse(), nil)

		res := resourceNsxtPolicyGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrefixListData())
		d.SetId(prefixListID)

		err := resourceNsxtPolicyGatewayPrefixListRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, prefixListDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(prefixListGwID, prefixListID).Return(nsxModel.PrefixList{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrefixListData())
		d.SetId(prefixListID)

		err := resourceNsxtPolicyGatewayPrefixListRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrefixListData())

		err := resourceNsxtPolicyGatewayPrefixListRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayPrefixListUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupPrefixListMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(prefixListGwID, prefixListID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(prefixListGwID, prefixListID).Return(prefixListAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrefixListData())
		d.SetId(prefixListID)

		err := resourceNsxtPolicyGatewayPrefixListUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrefixListData())

		err := resourceNsxtPolicyGatewayPrefixListUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayPrefixListDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupPrefixListMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(prefixListGwID, prefixListID).Return(nil)

		res := resourceNsxtPolicyGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrefixListData())
		d.SetId(prefixListID)

		err := resourceNsxtPolicyGatewayPrefixListDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrefixListData())

		err := resourceNsxtPolicyGatewayPrefixListDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
