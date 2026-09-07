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
	lbPassiveMonitorID          = "lb-passive-monitor-1"            //nolint:gosec
	lbPassiveMonitorDisplayName = "Test LB Passive Monitor Profile" //nolint:gosec
)

func minimalLBPassiveMonitorData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbPassiveMonitorDisplayName,
		"nsx_id":       lbPassiveMonitorID,
	}
}

func lbPassiveMonitorStructValue(t *testing.T) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	displayName := lbPassiveMonitorDisplayName
	description := "Test LB Passive Monitor Profile description"
	path := "/infra/lb-monitor-profiles/" + lbPassiveMonitorID
	revision := int64(1)
	maxFails := int64(7)
	timeout := int64(20)

	obj := model.LBPassiveMonitorProfile{
		DisplayName:  &displayName,
		Description:  &description,
		Path:         &path,
		Revision:     &revision,
		MaxFails:     &maxFails,
		Timeout:      &timeout,
		ResourceType: model.LBMonitorProfile_RESOURCE_TYPE_LBPASSIVEMONITORPROFILE,
	}
	val, errs := converter.ConvertToVapi(obj, model.LBPassiveMonitorProfileBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func assertLBPassiveMonitorFields(t *testing.T, d *schema.ResourceData) {
	t.Helper()
	assert.Equal(t, lbPassiveMonitorDisplayName, d.Get("display_name"))
	assert.Equal(t, "Test LB Passive Monitor Profile description", d.Get("description"))
	assert.Equal(t, 7, d.Get("max_fails"))
	assert.Equal(t, 20, d.Get("timeout"))
	assert.Equal(t, "/infra/lb-monitor-profiles/"+lbPassiveMonitorID, d.Get("path"))
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

	t.Run("Create success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		sv := lbPassiveMonitorStructValue(t)
		mockSDK.EXPECT().Get(gomock.Any()).Return(sv, nil)

		res := resourceNsxtPolicyLBPassiveMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": lbPassiveMonitorDisplayName,
		})

		err := resourceNsxtPolicyLBPassiveMonitorProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assertLBPassiveMonitorFields(t, d)
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

	t.Run("Read success", func(t *testing.T) {
		sv := lbPassiveMonitorStructValue(t)
		mockSDK.EXPECT().Get(lbPassiveMonitorID).Return(sv, nil)

		res := resourceNsxtPolicyLBPassiveMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPassiveMonitorData())
		d.SetId(lbPassiveMonitorID)

		err := resourceNsxtPolicyLBPassiveMonitorProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBPassiveMonitorFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBPassiveMonitorProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBPassiveMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPassiveMonitorData())

		err := resourceNsxtPolicyLBPassiveMonitorProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(lbPassiveMonitorID, gomock.Any()).Return(nil)
		sv := lbPassiveMonitorStructValue(t)
		mockSDK.EXPECT().Get(lbPassiveMonitorID).Return(sv, nil)

		res := resourceNsxtPolicyLBPassiveMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBPassiveMonitorData())
		d.SetId(lbPassiveMonitorID)

		err := resourceNsxtPolicyLBPassiveMonitorProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBPassiveMonitorFields(t, d)
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
