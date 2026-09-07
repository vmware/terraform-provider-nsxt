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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	infraMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

var (
	lbFastTcpID          = "lb-fast-tcp-1"
	lbFastTcpDisplayName = "Test LB Fast TCP Profile"
)

func minimalLBFastTcpData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbFastTcpDisplayName,
		"nsx_id":       lbFastTcpID,
	}
}

func setupLBAppProfileMock(t *testing.T, ctrl *gomock.Controller) (*infraMocks.MockLbAppProfilesClient, func()) {
	mockSDK := infraMocks.NewMockLbAppProfilesClient(ctrl)
	mockWrapper := &infraapi.LBAppProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliLbAppProfilesClient
	cliLbAppProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.LBAppProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliLbAppProfilesClient = original }
}

func lbFastTcpStructValue(t *testing.T) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	displayName := lbFastTcpDisplayName
	description := "Test LB Fast TCP Profile description"
	path := "/infra/lb-app-profiles/" + lbFastTcpID
	revision := int64(1)
	haFlowMirroringEnabled := true
	idleTimeout := int64(900)
	closeTimeout := int64(10)

	obj := model.LBFastTcpProfile{
		DisplayName:            &displayName,
		Description:            &description,
		Path:                   &path,
		Revision:               &revision,
		HaFlowMirroringEnabled: &haFlowMirroringEnabled,
		IdleTimeout:            &idleTimeout,
		CloseTimeout:           &closeTimeout,
		ResourceType:           model.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE,
	}
	val, errs := converter.ConvertToVapi(obj, model.LBFastTcpProfileBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func assertLBFastTcpFields(t *testing.T, d *schema.ResourceData) {
	t.Helper()
	assert.Equal(t, lbFastTcpDisplayName, d.Get("display_name"))
	assert.Equal(t, "Test LB Fast TCP Profile description", d.Get("description"))
	assert.Equal(t, true, d.Get("ha_flow_mirroring_enabled"))
	assert.Equal(t, 900, d.Get("idle_timeout"))
	assert.Equal(t, 10, d.Get("close_timeout"))
	assert.Equal(t, "/infra/lb-app-profiles/"+lbFastTcpID, d.Get("path"))
}

func TestMockResourceNsxtPolicyLBFastTcpApplicationProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		// Return nil error but non-nil struct to indicate object exists
		mockSDK.EXPECT().Get(lbFastTcpID).Return(nil, nil)

		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())

		err := resourceNsxtPolicyLBTcpApplicationProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		sv := lbFastTcpStructValue(t)
		mockSDK.EXPECT().Get(gomock.Any()).Return(sv, nil)

		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": lbFastTcpDisplayName,
		})

		err := resourceNsxtPolicyLBTcpApplicationProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assertLBFastTcpFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBFastTcpApplicationProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbFastTcpID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())
		d.SetId(lbFastTcpID)

		err := resourceNsxtPolicyLBTcpApplicationProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbFastTcpID).Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())
		d.SetId(lbFastTcpID)

		err := resourceNsxtPolicyLBTcpApplicationProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())

		err := resourceNsxtPolicyLBTcpApplicationProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read success", func(t *testing.T) {
		sv := lbFastTcpStructValue(t)
		mockSDK.EXPECT().Get(lbFastTcpID).Return(sv, nil)

		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())
		d.SetId(lbFastTcpID)

		err := resourceNsxtPolicyLBTcpApplicationProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBFastTcpFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBFastTcpApplicationProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())

		err := resourceNsxtPolicyLBTcpApplicationProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(lbFastTcpID, gomock.Any()).Return(nil)
		sv := lbFastTcpStructValue(t)
		mockSDK.EXPECT().Get(lbFastTcpID).Return(sv, nil)

		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())
		d.SetId(lbFastTcpID)

		err := resourceNsxtPolicyLBTcpApplicationProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBFastTcpFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBFastTcpApplicationProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbFastTcpID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())
		d.SetId(lbFastTcpID)

		err := resourceNsxtPolicyLBTcpApplicationProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())

		err := resourceNsxtPolicyLBTcpApplicationProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
