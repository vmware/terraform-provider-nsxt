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

func setupHostUpgradeGroupMock(ctrl *gomock.Controller) (*upgrademocks.MockUpgradeUnitGroupsClient, func()) {
	mockSDK := upgrademocks.NewMockUpgradeUnitGroupsClient(ctrl)
	original := cliUpgradeUnitGroupsClient
	cliUpgradeUnitGroupsClient = func(_ vapiProtocolClient.Connector) upgradeGroupOps {
		return mockSDK
	}
	return mockSDK, func() { cliUpgradeUnitGroupsClient = original }
}

func hostUpgradeGroupAPIResponse(id, name, groupType string) nsxModel.UpgradeUnitGroup {
	return nsxModel.UpgradeUnitGroup{
		Id:          &id,
		DisplayName: &name,
		Description: &name,
		Type_:       &groupType,
	}
}

func TestMockDataSourceNsxtHostUpgradeGroupRead(t *testing.T) {
	groupID := "host-group-1"
	groupName := "host-group-name"

	t.Run("by id success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHostUpgradeGroupMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get(groupID, nil).Return(hostUpgradeGroupAPIResponse(groupID, groupName, hostUpgradeGroup), nil)

		ds := dataSourceNsxtHostUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"id":                 groupID,
		})

		err := dataSourceNsxtHostUpgradeGroupRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, groupID, d.Id())
		assert.Equal(t, groupName, d.Get("display_name"))
	})

	t.Run("by id wrong group type is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHostUpgradeGroupMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get(groupID, nil).Return(hostUpgradeGroupAPIResponse(groupID, groupName, edgeUpgradeGroup), nil)

		ds := dataSourceNsxtHostUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"id":                 groupID,
		})

		err := dataSourceNsxtHostUpgradeGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by id API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHostUpgradeGroupMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get(groupID, nil).Return(nsxModel.UpgradeUnitGroup{}, errors.New("get failed"))

		ds := dataSourceNsxtHostUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"id":                 groupID,
		})

		err := dataSourceNsxtHostUpgradeGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error while reading")
	})

	t.Run("missing id and display_name errors", func(t *testing.T) {
		ds := dataSourceNsxtHostUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
		})

		err := dataSourceNsxtHostUpgradeGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining")
	})

	t.Run("by display_name single exact match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHostUpgradeGroupMock(ctrl)
		defer restore()

		groupType := hostUpgradeGroup
		mockSDK.EXPECT().List(&groupType, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.UpgradeUnitGroupListResult{
			Results: []nsxModel.UpgradeUnitGroup{hostUpgradeGroupAPIResponse(groupID, groupName, hostUpgradeGroup)},
		}, nil)

		ds := dataSourceNsxtHostUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"display_name":       groupName,
		})

		err := dataSourceNsxtHostUpgradeGroupRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, groupID, d.Id())
	})

	t.Run("by display_name prefix single match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHostUpgradeGroupMock(ctrl)
		defer restore()

		groupType := hostUpgradeGroup
		otherName := "host-group-other"
		mockSDK.EXPECT().List(&groupType, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.UpgradeUnitGroupListResult{
			Results: []nsxModel.UpgradeUnitGroup{hostUpgradeGroupAPIResponse(groupID, otherName, hostUpgradeGroup)},
		}, nil)

		ds := dataSourceNsxtHostUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"display_name":       "host-group-oth",
		})

		err := dataSourceNsxtHostUpgradeGroupRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, groupID, d.Id())
	})

	t.Run("by display_name multiple matches errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHostUpgradeGroupMock(ctrl)
		defer restore()

		groupType := hostUpgradeGroup
		mockSDK.EXPECT().List(&groupType, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.UpgradeUnitGroupListResult{
			Results: []nsxModel.UpgradeUnitGroup{
				hostUpgradeGroupAPIResponse("id-1", groupName, hostUpgradeGroup),
				hostUpgradeGroupAPIResponse("id-2", groupName, hostUpgradeGroup),
			},
		}, nil)

		ds := dataSourceNsxtHostUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"display_name":       groupName,
		})

		err := dataSourceNsxtHostUpgradeGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "found multiple")
	})

	t.Run("list error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHostUpgradeGroupMock(ctrl)
		defer restore()

		groupType := hostUpgradeGroup
		mockSDK.EXPECT().List(&groupType, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.UpgradeUnitGroupListResult{}, errors.New("list failed"))

		ds := dataSourceNsxtHostUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"display_name":       groupName,
		})

		err := dataSourceNsxtHostUpgradeGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error while reading")
	})

	t.Run("no match from list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupHostUpgradeGroupMock(ctrl)
		defer restore()

		groupType := hostUpgradeGroup
		mockSDK.EXPECT().List(&groupType, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.UpgradeUnitGroupListResult{
			Results: []nsxModel.UpgradeUnitGroup{},
		}, nil)

		ds := dataSourceNsxtHostUpgradeGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"upgrade_prepare_id": "prep-1",
			"display_name":       "nonexistent",
		})

		err := dataSourceNsxtHostUpgradeGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})
}
