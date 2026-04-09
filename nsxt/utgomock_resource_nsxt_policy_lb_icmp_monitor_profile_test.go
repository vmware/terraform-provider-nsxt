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
	lbIcmpMonitorID          = "lb-icmp-monitor-1"
	lbIcmpMonitorDisplayName = "Test LB ICMP Monitor Profile"
)

func minimalLBIcmpMonitorData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbIcmpMonitorDisplayName,
		"nsx_id":       lbIcmpMonitorID,
	}
}

func TestMockResourceNsxtPolicyLBIcmpMonitorProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbIcmpMonitorID).Return(nil, nil)

		res := resourceNsxtPolicyLBIcmpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBIcmpMonitorData())

		err := resourceNsxtPolicyLBIcmpMonitorProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyLBIcmpMonitorProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbIcmpMonitorID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBIcmpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBIcmpMonitorData())
		d.SetId(lbIcmpMonitorID)

		err := resourceNsxtPolicyLBIcmpMonitorProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBIcmpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBIcmpMonitorData())

		err := resourceNsxtPolicyLBIcmpMonitorProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBIcmpMonitorProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBIcmpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBIcmpMonitorData())

		err := resourceNsxtPolicyLBIcmpMonitorProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBIcmpMonitorProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbIcmpMonitorID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBIcmpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBIcmpMonitorData())
		d.SetId(lbIcmpMonitorID)

		err := resourceNsxtPolicyLBIcmpMonitorProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBIcmpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBIcmpMonitorData())

		err := resourceNsxtPolicyLBIcmpMonitorProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
