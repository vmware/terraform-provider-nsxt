//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	tier0sapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	tier0Mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
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

func TestMockResourceNsxtPolicyTier0GatewayImporter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTier0sSDK := tier0Mocks.NewMockTier0sClient(ctrl)
	tier0Wrapper := &cliinfra.Tier0ClientContext{
		Client:     mockTier0sSDK,
		ClientType: utl.Local,
	}
	original := cliTier0sClient
	defer func() { cliTier0sClient = original }()
	cliTier0sClient = func(sessionContext utl.SessionContext, connector vapiProtocolClient.Connector) *cliinfra.Tier0ClientContext {
		return tier0Wrapper
	}

	res := resourceNsxtPolicyGatewayPrefixList()

	t.Run("valid <gateway-id>/<id> sets gateway_path and id", func(t *testing.T) {
		gw := nsxModel.Tier0{Path: &prefixListGwPath}
		mockTier0sSDK.EXPECT().Get(prefixListGwID).Return(gw, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(prefixListGwID + "/" + prefixListID)

		out, err := resourceNsxtPolicyTier0GatewayImporter(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, prefixListID, d.Id())
		assert.Equal(t, prefixListGwPath, d.Get("gateway_path"))
	})

	t.Run("malformed id fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("no-slash-here")

		_, err := resourceNsxtPolicyTier0GatewayImporter(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("gateway lookup failure is propagated", func(t *testing.T) {
		mockTier0sSDK.EXPECT().Get(prefixListGwID).Return(nsxModel.Tier0{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(prefixListGwID + "/" + prefixListID)

		_, err := resourceNsxtPolicyTier0GatewayImporter(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
