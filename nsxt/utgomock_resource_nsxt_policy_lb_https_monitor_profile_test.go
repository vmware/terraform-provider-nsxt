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
	lbHttpsMonitorID          = "lb-https-monitor-1"
	lbHttpsMonitorDisplayName = "Test LB HTTPS Monitor Profile"
)

func minimalLBHttpsMonitorData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbHttpsMonitorDisplayName,
		"nsx_id":       lbHttpsMonitorID,
	}
}

func TestMockResourceNsxtPolicyLBHttpsMonitorProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpsMonitorID).Return(nil, nil)

		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpsMonitorData())

		err := resourceNsxtPolicyLBHttpsMonitorProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyLBHttpsMonitorProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpsMonitorID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpsMonitorData())
		d.SetId(lbHttpsMonitorID)

		err := resourceNsxtPolicyLBHttpsMonitorProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpsMonitorData())

		err := resourceNsxtPolicyLBHttpsMonitorProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBHttpsMonitorProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpsMonitorData())

		err := resourceNsxtPolicyLBHttpsMonitorProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBHttpsMonitorProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbHttpsMonitorID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpsMonitorData())
		d.SetId(lbHttpsMonitorID)

		err := resourceNsxtPolicyLBHttpsMonitorProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpsMonitorData())

		err := resourceNsxtPolicyLBHttpsMonitorProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
