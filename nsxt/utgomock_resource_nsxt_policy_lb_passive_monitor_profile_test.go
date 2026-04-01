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
	lbPassiveMonitorID          = "lb-passive-monitor-1"            //nolint:gosec
	lbPassiveMonitorDisplayName = "Test LB Passive Monitor Profile" //nolint:gosec
)

func minimalLBPassiveMonitorData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbPassiveMonitorDisplayName,
		"nsx_id":       lbPassiveMonitorID,
	}
}

func TestMockResourceNsxtPolicyLBPassiveMonitorProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbPassiveMonitorID).Return(nil, nil)

		res := resourceNsxtPolicyLBPassiveMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPassiveMonitorData())

		err := resourceNsxtPolicyLBPassiveMonitorProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyLBPassiveMonitorProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbPassiveMonitorID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBPassiveMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPassiveMonitorData())
		d.SetId(lbPassiveMonitorID)

		err := resourceNsxtPolicyLBPassiveMonitorProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBPassiveMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPassiveMonitorData())

		err := resourceNsxtPolicyLBPassiveMonitorProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBPassiveMonitorProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBPassiveMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPassiveMonitorData())

		err := resourceNsxtPolicyLBPassiveMonitorProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBPassiveMonitorProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbPassiveMonitorID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBPassiveMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPassiveMonitorData())
		d.SetId(lbPassiveMonitorID)

		err := resourceNsxtPolicyLBPassiveMonitorProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBPassiveMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPassiveMonitorData())

		err := resourceNsxtPolicyLBPassiveMonitorProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
