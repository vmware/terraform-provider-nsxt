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
	lbHttpAppID          = "lb-http-app-1"
	lbHttpAppDisplayName = "Test LB HTTP App Profile"
)

func minimalLBHttpAppData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbHttpAppDisplayName,
		"nsx_id":       lbHttpAppID,
	}
}

func TestMockResourceNsxtPolicyLBHttpApplicationProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpAppID).Return(nil, nil)

		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())

		err := resourceNsxtPolicyLBHttpApplicationProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyLBHttpApplicationProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpAppID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())
		d.SetId(lbHttpAppID)

		err := resourceNsxtPolicyLBHttpApplicationProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpAppID).Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())
		d.SetId(lbHttpAppID)

		err := resourceNsxtPolicyLBHttpApplicationProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())

		err := resourceNsxtPolicyLBHttpApplicationProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBHttpApplicationProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())

		err := resourceNsxtPolicyLBHttpApplicationProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBHttpApplicationProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbHttpAppID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())
		d.SetId(lbHttpAppID)

		err := resourceNsxtPolicyLBHttpApplicationProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())

		err := resourceNsxtPolicyLBHttpApplicationProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
