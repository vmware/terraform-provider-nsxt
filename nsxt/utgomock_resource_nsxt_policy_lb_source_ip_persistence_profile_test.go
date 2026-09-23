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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"
)

var (
	lbSourceIpID          = "lb-source-ip-1"
	lbSourceIpDisplayName = "Test LB Source IP Persistence Profile"
)

func minimalLBSourceIpData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbSourceIpDisplayName,
		"nsx_id":       lbSourceIpID,
	}
}

func TestMockResourceNsxtPolicyLBSourceIpPersistenceProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbSourceIpID).Return(nil, nil)

		res := resourceNsxtPolicyLBSourceIpPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBSourceIpData())

		err := resourceNsxtPolicyLBSourceIpPersistenceProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create succeeds and re-reads the profile", func(t *testing.T) {
		displayName := lbSourceIpDisplayName
		sv := lbPersistenceProfileStructValue(t, model.LBSourceIpPersistenceProfile{
			DisplayName:  &displayName,
			ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBSOURCEIPPERSISTENCEPROFILE,
		}, model.LBSourceIpPersistenceProfileBindingType())

		mockSDK.EXPECT().Get(lbSourceIpID).Return(nil, vapiErrors.NotFound{})
		mockSDK.EXPECT().Patch(lbSourceIpID, gomock.Any()).Return(nil)
		mockSDK.EXPECT().Get(lbSourceIpID).Return(sv, nil)

		res := resourceNsxtPolicyLBSourceIpPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBSourceIpData())

		err := resourceNsxtPolicyLBSourceIpPersistenceProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbSourceIpID, d.Id())
		assert.Equal(t, lbSourceIpDisplayName, d.Get("display_name"))
	})
}

func TestMockResourceNsxtPolicyLBSourceIpPersistenceProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbSourceIpID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBSourceIpPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBSourceIpData())
		d.SetId(lbSourceIpID)

		err := resourceNsxtPolicyLBSourceIpPersistenceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbSourceIpID).Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyLBSourceIpPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBSourceIpData())
		d.SetId(lbSourceIpID)

		err := resourceNsxtPolicyLBSourceIpPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBSourceIpPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBSourceIpData())

		err := resourceNsxtPolicyLBSourceIpPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read success populates schema fields", func(t *testing.T) {
		displayName := lbSourceIpDisplayName
		description := "a source-ip profile"
		sv := lbPersistenceProfileStructValue(t, model.LBSourceIpPersistenceProfile{
			DisplayName:  &displayName,
			Description:  &description,
			ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBSOURCEIPPERSISTENCEPROFILE,
			Purge:        strPtr(model.LBSourceIpPersistenceProfile_PURGE_FULL),
			Timeout:      int64Ptr(300),
		}, model.LBSourceIpPersistenceProfileBindingType())
		mockSDK.EXPECT().Get(lbSourceIpID).Return(sv, nil)

		res := resourceNsxtPolicyLBSourceIpPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBSourceIpData())
		d.SetId(lbSourceIpID)

		err := resourceNsxtPolicyLBSourceIpPersistenceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbSourceIpDisplayName, d.Get("display_name"))
		assert.Equal(t, description, d.Get("description"))
	})
}

func TestMockResourceNsxtPolicyLBSourceIpPersistenceProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBSourceIpPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBSourceIpData())

		err := resourceNsxtPolicyLBSourceIpPersistenceProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update succeeds and re-reads the profile", func(t *testing.T) {
		displayName := lbSourceIpDisplayName
		sv := lbPersistenceProfileStructValue(t, model.LBSourceIpPersistenceProfile{
			DisplayName:  &displayName,
			ResourceType: model.LBPersistenceProfile_RESOURCE_TYPE_LBSOURCEIPPERSISTENCEPROFILE,
		}, model.LBSourceIpPersistenceProfileBindingType())

		mockSDK.EXPECT().Update(lbSourceIpID, gomock.Any()).Return(nil, nil)
		mockSDK.EXPECT().Get(lbSourceIpID).Return(sv, nil)

		res := resourceNsxtPolicyLBSourceIpPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBSourceIpData())
		d.SetId(lbSourceIpID)

		err := resourceNsxtPolicyLBSourceIpPersistenceProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyLBSourceIpPersistenceProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbSourceIpID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBSourceIpPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBSourceIpData())
		d.SetId(lbSourceIpID)

		err := resourceNsxtPolicyLBSourceIpPersistenceProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBSourceIpPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBSourceIpData())

		err := resourceNsxtPolicyLBSourceIpPersistenceProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
