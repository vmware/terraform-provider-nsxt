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
	lbTcpMonitorID          = "lb-tcp-monitor-1"
	lbTcpMonitorDisplayName = "Test LB TCP Monitor Profile"
)

func minimalLBTcpMonitorData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbTcpMonitorDisplayName,
		"nsx_id":       lbTcpMonitorID,
	}
}

func TestMockResourceNsxtPolicyLBTcpMonitorProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbTcpMonitorID).Return(nil, nil)

		res := resourceNsxtPolicyLBTcpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBTcpMonitorData())

		err := resourceNsxtPolicyLBTcpMonitorProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyLBTcpMonitorProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbTcpMonitorID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBTcpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBTcpMonitorData())
		d.SetId(lbTcpMonitorID)

		err := resourceNsxtPolicyLBTcpMonitorProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBTcpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBTcpMonitorData())

		err := resourceNsxtPolicyLBTcpMonitorProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBTcpMonitorProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBTcpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBTcpMonitorData())

		err := resourceNsxtPolicyLBTcpMonitorProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBTcpMonitorProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbTcpMonitorID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBTcpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBTcpMonitorData())
		d.SetId(lbTcpMonitorID)

		err := resourceNsxtPolicyLBTcpMonitorProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBTcpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBTcpMonitorData())

		err := resourceNsxtPolicyLBTcpMonitorProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
