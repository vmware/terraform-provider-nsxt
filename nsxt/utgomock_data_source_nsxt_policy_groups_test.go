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
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestMockDataSourceNsxtPolicyGroupsRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGroupMock(t, ctrl)
	defer restore()

	boolFalse := false

	t.Run("single page aggregates items", func(t *testing.T) {
		g2Name := "group-two"
		g2Path := "/infra/domains/default/groups/group-two"
		g2ID := "group-2"
		g2 := nsxModel.Group{
			Id:          &g2ID,
			DisplayName: &g2Name,
			Path:        &g2Path,
		}

		mockSDK.EXPECT().List(groupDomain, nil, nil, nil, nil, nil, &boolFalse, nil).Return(nsxModel.GroupListResult{
			Results: []nsxModel.Group{groupAPIResponse(), g2},
			Cursor:  nil,
		}, nil)

		ds := dataSourceNsxtPolicyGroups()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain": groupDomain,
		})

		err := dataSourceNsxtPolicyGroupsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		assert.Equal(t, groupPath, items[groupDisplayName])
		assert.Equal(t, g2Path, items[g2Name])
		assert.NotEmpty(t, d.Id())
	})

	t.Run("two pages merge results", func(t *testing.T) {
		cursor1 := "next"
		g2Name := "page2-group"
		g2Path := "/infra/domains/default/groups/page2"
		g2ID := "g-page-2"
		g2 := nsxModel.Group{
			Id:          &g2ID,
			DisplayName: &g2Name,
			Path:        &g2Path,
		}

		gomock.InOrder(
			mockSDK.EXPECT().List(groupDomain, nil, nil, nil, nil, nil, &boolFalse, nil).Return(nsxModel.GroupListResult{
				Results: []nsxModel.Group{groupAPIResponse()},
				Cursor:  &cursor1,
			}, nil),
			mockSDK.EXPECT().List(groupDomain, &cursor1, nil, nil, nil, nil, &boolFalse, nil).Return(nsxModel.GroupListResult{
				Results: []nsxModel.Group{g2},
				Cursor:  nil,
			}, nil),
		)

		ds := dataSourceNsxtPolicyGroups()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain": groupDomain,
		})

		err := dataSourceNsxtPolicyGroupsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		assert.Equal(t, groupPath, items[groupDisplayName])
		assert.Equal(t, g2Path, items[g2Name])
	})

	t.Run("List API error", func(t *testing.T) {
		mockSDK.EXPECT().List(groupDomain, nil, nil, nil, nil, nil, &boolFalse, nil).Return(nsxModel.GroupListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyGroups()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain": groupDomain,
		})

		err := dataSourceNsxtPolicyGroupsRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
