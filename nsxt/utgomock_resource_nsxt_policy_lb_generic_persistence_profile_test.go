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
	lbGenericPersistID          = "lb-generic-persist-1"
	lbGenericPersistDisplayName = "Test LB Generic Persistence Profile"
)

func minimalLBGenericPersistData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbGenericPersistDisplayName,
		"nsx_id":       lbGenericPersistID,
	}
}

func TestMockResourceNsxtPolicyLBGenericPersistenceProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbGenericPersistID).Return(nil, nil)

		res := resourceNsxtPolicyLBGenericPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBGenericPersistData())

		err := resourceNsxtPolicyLBGenericPersistenceProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyLBGenericPersistenceProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbGenericPersistID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBGenericPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBGenericPersistData())
		d.SetId(lbGenericPersistID)

		err := resourceNsxtPolicyLBGenericPersistenceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbGenericPersistID).Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyLBGenericPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBGenericPersistData())
		d.SetId(lbGenericPersistID)

		err := resourceNsxtPolicyLBGenericPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBGenericPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBGenericPersistData())

		err := resourceNsxtPolicyLBGenericPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBGenericPersistenceProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBGenericPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBGenericPersistData())

		err := resourceNsxtPolicyLBGenericPersistenceProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBGenericPersistenceProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbGenericPersistID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBGenericPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBGenericPersistData())
		d.SetId(lbGenericPersistID)

		err := resourceNsxtPolicyLBGenericPersistenceProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBGenericPersistenceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBGenericPersistData())

		err := resourceNsxtPolicyLBGenericPersistenceProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
