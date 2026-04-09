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
	lbUdpMonitorID          = "lb-udp-monitor-1"
	lbUdpMonitorDisplayName = "Test LB UDP Monitor Profile"
)

func minimalLBUdpMonitorData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbUdpMonitorDisplayName,
		"nsx_id":       lbUdpMonitorID,
	}
}

func TestMockResourceNsxtPolicyLBUdpMonitorProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbUdpMonitorID).Return(nil, nil)

		res := resourceNsxtPolicyLBUdpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBUdpMonitorData())

		err := resourceNsxtPolicyLBUdpMonitorProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyLBUdpMonitorProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbUdpMonitorID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBUdpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBUdpMonitorData())
		d.SetId(lbUdpMonitorID)

		err := resourceNsxtPolicyLBUdpMonitorProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBUdpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBUdpMonitorData())

		err := resourceNsxtPolicyLBUdpMonitorProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBUdpMonitorProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBUdpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBUdpMonitorData())

		err := resourceNsxtPolicyLBUdpMonitorProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBUdpMonitorProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbUdpMonitorID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBUdpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBUdpMonitorData())
		d.SetId(lbUdpMonitorID)

		err := resourceNsxtPolicyLBUdpMonitorProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBUdpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBUdpMonitorData())

		err := resourceNsxtPolicyLBUdpMonitorProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
