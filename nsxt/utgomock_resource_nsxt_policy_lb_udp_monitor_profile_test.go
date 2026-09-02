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
	lbUdpMonitorID          = "lb-udp-monitor-1"
	lbUdpMonitorDisplayName = "Test LB UDP Monitor Profile"
)

func minimalLBUdpMonitorData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbUdpMonitorDisplayName,
		"nsx_id":       lbUdpMonitorID,
	}
}

func lbUdpMonitorStructValue(t *testing.T) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	displayName := lbUdpMonitorDisplayName
	description := "Test LB UDP Monitor Profile description"
	path := "/infra/lb-monitor-profiles/" + lbUdpMonitorID
	revision := int64(1)
	receive := "OK"
	send := "PING"
	fallCount := int64(3)
	interval := int64(5)
	monitorPort := int64(80)
	riseCount := int64(3)
	timeout := int64(15)

	obj := model.LBUdpMonitorProfile{
		DisplayName:  &displayName,
		Description:  &description,
		Path:         &path,
		Revision:     &revision,
		Receive:      &receive,
		Send:         &send,
		FallCount:    &fallCount,
		Interval:     &interval,
		MonitorPort:  &monitorPort,
		RiseCount:    &riseCount,
		Timeout:      &timeout,
		ResourceType: model.LBMonitorProfile_RESOURCE_TYPE_LBUDPMONITORPROFILE,
	}
	val, errs := converter.ConvertToVapi(obj, model.LBUdpMonitorProfileBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func assertLBUdpMonitorFields(t *testing.T, d *schema.ResourceData) {
	t.Helper()
	assert.Equal(t, lbUdpMonitorDisplayName, d.Get("display_name"))
	assert.Equal(t, "Test LB UDP Monitor Profile description", d.Get("description"))
	assert.Equal(t, "OK", d.Get("receive"))
	assert.Equal(t, "PING", d.Get("send"))
	assert.Equal(t, 3, d.Get("fall_count"))
	assert.Equal(t, 5, d.Get("interval"))
	assert.Equal(t, 80, d.Get("monitor_port"))
	assert.Equal(t, 3, d.Get("rise_count"))
	assert.Equal(t, 15, d.Get("timeout"))
	assert.Equal(t, "/infra/lb-monitor-profiles/"+lbUdpMonitorID, d.Get("path"))
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

	t.Run("Create success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		sv := lbUdpMonitorStructValue(t)
		mockSDK.EXPECT().Get(gomock.Any()).Return(sv, nil)

		res := resourceNsxtPolicyLBUdpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": lbUdpMonitorDisplayName,
		})

		err := resourceNsxtPolicyLBUdpMonitorProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assertLBUdpMonitorFields(t, d)
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

	t.Run("Read success", func(t *testing.T) {
		sv := lbUdpMonitorStructValue(t)
		mockSDK.EXPECT().Get(lbUdpMonitorID).Return(sv, nil)

		res := resourceNsxtPolicyLBUdpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBUdpMonitorData())
		d.SetId(lbUdpMonitorID)

		err := resourceNsxtPolicyLBUdpMonitorProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBUdpMonitorFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBUdpMonitorProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBUdpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBUdpMonitorData())

		err := resourceNsxtPolicyLBUdpMonitorProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(lbUdpMonitorID, gomock.Any()).Return(nil)
		sv := lbUdpMonitorStructValue(t)
		mockSDK.EXPECT().Get(lbUdpMonitorID).Return(sv, nil)

		res := resourceNsxtPolicyLBUdpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBUdpMonitorData())
		d.SetId(lbUdpMonitorID)

		err := resourceNsxtPolicyLBUdpMonitorProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBUdpMonitorFields(t, d)
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
