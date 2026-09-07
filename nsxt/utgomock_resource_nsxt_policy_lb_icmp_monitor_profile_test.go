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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"
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

func lbIcmpMonitorStructValue(t *testing.T) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	displayName := lbIcmpMonitorDisplayName
	description := "Test LB ICMP Monitor Profile description"
	path := "/infra/lb-monitor-profiles/" + lbIcmpMonitorID
	revision := int64(1)
	dataLength := int64(56)
	fallCount := int64(3)
	interval := int64(5)
	riseCount := int64(3)
	timeout := int64(15)

	obj := model.LBIcmpMonitorProfile{
		DisplayName:  &displayName,
		Description:  &description,
		Path:         &path,
		Revision:     &revision,
		DataLength:   &dataLength,
		FallCount:    &fallCount,
		Interval:     &interval,
		RiseCount:    &riseCount,
		Timeout:      &timeout,
		ResourceType: model.LBMonitorProfile_RESOURCE_TYPE_LBICMPMONITORPROFILE,
	}
	val, errs := converter.ConvertToVapi(obj, model.LBIcmpMonitorProfileBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func assertLBIcmpMonitorFields(t *testing.T, d *schema.ResourceData) {
	t.Helper()
	assert.Equal(t, lbIcmpMonitorDisplayName, d.Get("display_name"))
	assert.Equal(t, "Test LB ICMP Monitor Profile description", d.Get("description"))
	assert.Equal(t, 56, d.Get("data_length"))
	assert.Equal(t, 3, d.Get("fall_count"))
	assert.Equal(t, 5, d.Get("interval"))
	assert.Equal(t, 3, d.Get("rise_count"))
	assert.Equal(t, 15, d.Get("timeout"))
	assert.Equal(t, "/infra/lb-monitor-profiles/"+lbIcmpMonitorID, d.Get("path"))
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

	t.Run("Create success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		sv := lbIcmpMonitorStructValue(t)
		mockSDK.EXPECT().Get(gomock.Any()).Return(sv, nil)

		res := resourceNsxtPolicyLBIcmpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": lbIcmpMonitorDisplayName,
		})

		err := resourceNsxtPolicyLBIcmpMonitorProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assertLBIcmpMonitorFields(t, d)
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

	t.Run("Read success", func(t *testing.T) {
		sv := lbIcmpMonitorStructValue(t)
		mockSDK.EXPECT().Get(lbIcmpMonitorID).Return(sv, nil)

		res := resourceNsxtPolicyLBIcmpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBIcmpMonitorData())
		d.SetId(lbIcmpMonitorID)

		err := resourceNsxtPolicyLBIcmpMonitorProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBIcmpMonitorFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBIcmpMonitorProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBIcmpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBIcmpMonitorData())

		err := resourceNsxtPolicyLBIcmpMonitorProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(lbIcmpMonitorID, gomock.Any()).Return(nil)
		sv := lbIcmpMonitorStructValue(t)
		mockSDK.EXPECT().Get(lbIcmpMonitorID).Return(sv, nil)

		res := resourceNsxtPolicyLBIcmpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBIcmpMonitorData())
		d.SetId(lbIcmpMonitorID)

		err := resourceNsxtPolicyLBIcmpMonitorProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBIcmpMonitorFields(t, d)
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
