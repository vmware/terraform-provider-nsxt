//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

func TestMockDataSourceNsxtComputeManagerRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupComputeManagerMock(t, ctrl)
	defer restore()

	t.Run("by id success", func(t *testing.T) {
		mockSDK.EXPECT().Get(cmID).Return(computeManagerAPIResponse(), nil)

		ds := dataSourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": cmID,
		})

		err := dataSourceNsxtComputeManagerRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, cmID, d.Id())
		assert.Equal(t, cmDisplayName, d.Get("display_name"))
		assert.Equal(t, cmServer, d.Get("server"))
	})

	t.Run("by id API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(cmID).Return(nsxModel.ComputeManager{}, errors.New("get failed"))

		ds := dataSourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": cmID,
		})

		err := dataSourceNsxtComputeManagerRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read ComputeManager")
	})

	t.Run("by display_name single exact match", func(t *testing.T) {
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ComputeManagerListResult{
			Results: []nsxModel.ComputeManager{computeManagerAPIResponse()},
		}, nil)

		ds := dataSourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": cmDisplayName,
		})

		err := dataSourceNsxtComputeManagerRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, cmID, d.Id())
	})

	t.Run("by display_name prefix single match", func(t *testing.T) {
		otherName := "vcenter-other"
		otherID := "cm-2"
		other := computeManagerAPIResponse()
		other.DisplayName = &otherName
		other.Id = &otherID

		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ComputeManagerListResult{
			Results: []nsxModel.ComputeManager{other},
		}, nil)

		ds := dataSourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "vcenter-oth",
		})

		err := dataSourceNsxtComputeManagerRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, otherID, d.Id())
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ComputeManagerListResult{}, errors.New("list failed"))

		ds := dataSourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": cmDisplayName,
		})

		err := dataSourceNsxtComputeManagerRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read Compute Managers")
	})

	t.Run("no match from list", func(t *testing.T) {
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ComputeManagerListResult{
			Results: []nsxModel.ComputeManager{},
		}, nil)

		ds := dataSourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		err := dataSourceNsxtComputeManagerRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "No Compute Manager matches")
	})
}
