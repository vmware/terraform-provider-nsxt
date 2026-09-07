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
	lbFastUdpID          = "lb-fast-udp-1"
	lbFastUdpDisplayName = "Test LB Fast UDP Profile"
)

func minimalLBFastUdpData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbFastUdpDisplayName,
		"nsx_id":       lbFastUdpID,
	}
}

func lbFastUdpStructValue(t *testing.T) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	displayName := lbFastUdpDisplayName
	description := "Test LB Fast UDP Profile description"
	path := "/infra/lb-app-profiles/" + lbFastUdpID
	revision := int64(1)
	flowMirroringEnabled := true
	idleTimeout := int64(150)

	obj := model.LBFastUdpProfile{
		DisplayName:          &displayName,
		Description:          &description,
		Path:                 &path,
		Revision:             &revision,
		FlowMirroringEnabled: &flowMirroringEnabled,
		IdleTimeout:          &idleTimeout,
		ResourceType:         model.LBAppProfile_RESOURCE_TYPE_LBFASTUDPPROFILE,
	}
	val, errs := converter.ConvertToVapi(obj, model.LBFastUdpProfileBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func assertLBFastUdpFields(t *testing.T, d *schema.ResourceData) {
	t.Helper()
	assert.Equal(t, lbFastUdpDisplayName, d.Get("display_name"))
	assert.Equal(t, "Test LB Fast UDP Profile description", d.Get("description"))
	assert.Equal(t, true, d.Get("flow_mirroring_enabled"))
	assert.Equal(t, 150, d.Get("idle_timeout"))
	assert.Equal(t, "/infra/lb-app-profiles/"+lbFastUdpID, d.Get("path"))
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

	t.Run("Create success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		sv := lbFastUdpStructValue(t)
		mockSDK.EXPECT().Get(gomock.Any()).Return(sv, nil)

		res := resourceNsxtPolicyLBFastUdpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": lbFastUdpDisplayName,
		})

		err := resourceNsxtPolicyLBUdpApplicationProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assertLBFastUdpFields(t, d)
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

	t.Run("Read success", func(t *testing.T) {
		sv := lbFastUdpStructValue(t)
		mockSDK.EXPECT().Get(lbFastUdpID).Return(sv, nil)

		res := resourceNsxtPolicyLBFastUdpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastUdpData())
		d.SetId(lbFastUdpID)

		err := resourceNsxtPolicyLBUdpApplicationProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBFastUdpFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBFastUdpApplicationProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBFastUdpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastUdpData())

		err := resourceNsxtPolicyLBUdpApplicationProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(lbFastUdpID, gomock.Any()).Return(nil)
		sv := lbFastUdpStructValue(t)
		mockSDK.EXPECT().Get(lbFastUdpID).Return(sv, nil)

		res := resourceNsxtPolicyLBFastUdpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastUdpData())
		d.SetId(lbFastUdpID)

		err := resourceNsxtPolicyLBUdpApplicationProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBFastUdpFields(t, d)
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
