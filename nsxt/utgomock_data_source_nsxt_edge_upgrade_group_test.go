//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// Reuses the existing UpgradeUnitGroupsClient mock generated for
// resource_nsxt_upgrade_run.go (see utgomock_resource_nsxt_upgrade_run_test.go):
// mockgen -destination=mocks/nsx/upgrade/UpgradeUnitGroupsClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/upgrade/UpgradeUnitGroupsClient.go UpgradeUnitGroupsClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"go.uber.org/mock/gomock"

	upgrademocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/upgrade"
)

func setupEdgeUpgradeGroupMock(ctrl *gomock.Controller) (*upgrademocks.MockUpgradeUnitGroupsClient, func()) {
	mockSDK := upgrademocks.NewMockUpgradeUnitGroupsClient(ctrl)
	original := cliUpgradeUnitGroupsClient
	cliUpgradeUnitGroupsClient = func(_ vapiProtocolClient.Connector) upgradeGroupOps {
		return mockSDK
	}
	return mockSDK, func() { cliUpgradeUnitGroupsClient = original }
}

func edgeUpgradeGroupAPIResponse(id, name, groupType string) nsxModel.UpgradeUnitGroup {
	return nsxModel.UpgradeUnitGroup{
		Id:          &id,
		DisplayName: &name,
		Description: &name,
		Type_:       &groupType,
	}
}

func TestMockDataSourceNsxtEdgeUpgradeGroupRead(t *testing.T) {
	groupID := "edge-group-1"
	groupName := "edge-group-name"

	t.Run("by id success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeUpgradeGroupMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get(groupID, nil).Return(edgeUpgradeGroupAPIResponse(groupID, groupName, edgeUpgradeGroup), nil)

		ds := dataSourceNsxtEdgeUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"id":                 groupID,
		})

		err := dataSourceNsxtEdgeUpgradeGroupRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, groupID, d.Id())
		assert.Equal(t, groupName, d.Get("display_name"))
	})

	t.Run("by id wrong group type is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeUpgradeGroupMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get(groupID, nil).Return(edgeUpgradeGroupAPIResponse(groupID, groupName, hostUpgradeGroup), nil)

		ds := dataSourceNsxtEdgeUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"id":                 groupID,
		})

		err := dataSourceNsxtEdgeUpgradeGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by id API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeUpgradeGroupMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get(groupID, nil).Return(nsxModel.UpgradeUnitGroup{}, errors.New("get failed"))

		ds := dataSourceNsxtEdgeUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"id":                 groupID,
		})

		err := dataSourceNsxtEdgeUpgradeGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error while reading")
	})

	t.Run("missing id and display_name errors", func(t *testing.T) {
		ds := dataSourceNsxtEdgeUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
		})

		err := dataSourceNsxtEdgeUpgradeGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining")
	})

	t.Run("by display_name single exact match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeUpgradeGroupMock(ctrl)
		defer restore()

		groupType := edgeUpgradeGroup
		mockSDK.EXPECT().List(&groupType, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.UpgradeUnitGroupListResult{
			Results: []nsxModel.UpgradeUnitGroup{edgeUpgradeGroupAPIResponse(groupID, groupName, edgeUpgradeGroup)},
		}, nil)

		ds := dataSourceNsxtEdgeUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"display_name":       groupName,
		})

		err := dataSourceNsxtEdgeUpgradeGroupRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, groupID, d.Id())
	})

	t.Run("by display_name prefix single match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeUpgradeGroupMock(ctrl)
		defer restore()

		groupType := edgeUpgradeGroup
		otherName := "edge-group-other"
		mockSDK.EXPECT().List(&groupType, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.UpgradeUnitGroupListResult{
			Results: []nsxModel.UpgradeUnitGroup{edgeUpgradeGroupAPIResponse(groupID, otherName, edgeUpgradeGroup)},
		}, nil)

		ds := dataSourceNsxtEdgeUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"display_name":       "edge-group-oth",
		})

		err := dataSourceNsxtEdgeUpgradeGroupRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, groupID, d.Id())
	})

	t.Run("by display_name multiple matches errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeUpgradeGroupMock(ctrl)
		defer restore()

		groupType := edgeUpgradeGroup
		mockSDK.EXPECT().List(&groupType, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.UpgradeUnitGroupListResult{
			Results: []nsxModel.UpgradeUnitGroup{
				edgeUpgradeGroupAPIResponse("id-1", groupName, edgeUpgradeGroup),
				edgeUpgradeGroupAPIResponse("id-2", groupName, edgeUpgradeGroup),
			},
		}, nil)

		ds := dataSourceNsxtEdgeUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"display_name":       groupName,
		})

		err := dataSourceNsxtEdgeUpgradeGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "found multiple")
	})

	t.Run("list error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeUpgradeGroupMock(ctrl)
		defer restore()

		groupType := edgeUpgradeGroup
		mockSDK.EXPECT().List(&groupType, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.UpgradeUnitGroupListResult{}, errors.New("list failed"))

		ds := dataSourceNsxtEdgeUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"display_name":       groupName,
		})

		err := dataSourceNsxtEdgeUpgradeGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error while reading")
	})

	t.Run("no match from list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeUpgradeGroupMock(ctrl)
		defer restore()

		groupType := edgeUpgradeGroup
		mockSDK.EXPECT().List(&groupType, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.UpgradeUnitGroupListResult{
			Results: []nsxModel.UpgradeUnitGroup{},
		}, nil)

		ds := dataSourceNsxtEdgeUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"display_name":       "nonexistent",
		})

		err := dataSourceNsxtEdgeUpgradeGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})
}
