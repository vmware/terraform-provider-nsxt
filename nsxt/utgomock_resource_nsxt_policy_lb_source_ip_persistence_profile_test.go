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
}

func TestMockResourceNsxtPolicyLBSourceIpPersistenceProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBSourceIpPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBSourceIpData())

		err := resourceNsxtPolicyLBSourceIpPersistenceProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
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
