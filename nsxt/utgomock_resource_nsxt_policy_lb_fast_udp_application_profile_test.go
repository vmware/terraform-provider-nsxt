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
	lbFastUdpID          = "lb-fast-udp-1"
	lbFastUdpDisplayName = "Test LB Fast UDP Profile"
)

func minimalLBFastUdpData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbFastUdpDisplayName,
		"nsx_id":       lbFastUdpID,
	}
}

func TestMockResourceNsxtPolicyLBFastUdpApplicationProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbFastUdpID).Return(nil, nil)

		res := resourceNsxtPolicyLBFastUdpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastUdpData())

		err := resourceNsxtPolicyLBUdpApplicationProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyLBFastUdpApplicationProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbFastUdpID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBFastUdpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastUdpData())
		d.SetId(lbFastUdpID)

		err := resourceNsxtPolicyLBUdpApplicationProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbFastUdpID).Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyLBFastUdpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastUdpData())
		d.SetId(lbFastUdpID)

		err := resourceNsxtPolicyLBUdpApplicationProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBFastUdpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastUdpData())

		err := resourceNsxtPolicyLBUdpApplicationProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBFastUdpApplicationProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBFastUdpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastUdpData())

		err := resourceNsxtPolicyLBUdpApplicationProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBFastUdpApplicationProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbFastUdpID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBFastUdpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastUdpData())
		d.SetId(lbFastUdpID)

		err := resourceNsxtPolicyLBUdpApplicationProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBFastUdpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastUdpData())

		err := resourceNsxtPolicyLBUdpApplicationProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
