//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/nsx/node/UsersClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/node/UsersClient.go UsersClient

package nsxt

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/node"

	nodemocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/node"
)

var (
	nodeUserID            = int64(1001)
	nodeUserIDStr         = "1001"
	nodeUserFullName      = "Test User"
	nodeUserName          = "testuser"
	nodeUserStatus        = nsxModel.NodeUserProperties_STATUS_ACTIVE
	nodeUserLastPwdChange = int64(30)
	nodeUserPwdFreq       = int64(90)
	nodeUserPwdWarn       = int64(7)
	nodeUserPwdReset      = false
)

func nodeUserAPIResponse() nsxModel.NodeUserProperties {
	return nsxModel.NodeUserProperties{
		Userid:                  &nodeUserID,
		FullName:                &nodeUserFullName,
		Username:                &nodeUserName,
		Status:                  &nodeUserStatus,
		LastPasswordChange:      &nodeUserLastPwdChange,
		PasswordChangeFrequency: &nodeUserPwdFreq,
		PasswordChangeWarning:   &nodeUserPwdWarn,
		PasswordResetRequired:   &nodeUserPwdReset,
	}
}

func nodeUserNotActivatedAPIResponse() nsxModel.NodeUserProperties {
	status := nsxModel.NodeUserProperties_STATUS_NOT_ACTIVATED
	resp := nodeUserAPIResponse()
	resp.Status = &status
	return resp
}

func setupNodeUserMock(t *testing.T, ctrl *gomock.Controller) (*nodemocks.MockUsersClient, func()) {
	mockSDK := nodemocks.NewMockUsersClient(ctrl)

	originalCli := cliNodeUsersClient
	cliNodeUsersClient = func(_ vapiProtocolClient.Connector) node.UsersClient {
		return mockSDK
	}
	return mockSDK, func() { cliNodeUsersClient = originalCli }
}

func nodeUserSchemaData() map[string]interface{} {
	return map[string]interface{}{
		"full_name":                 nodeUserFullName,
		"username":                  nodeUserName,
		"password_change_frequency": int(nodeUserPwdFreq),
		"password_change_warning":   int(nodeUserPwdWarn),
		"active":                    true,
	}
}

func TestMockResourceNsxtNodeUserCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNodeUserMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Createuser(gomock.Any()).Return(nodeUserAPIResponse(), nil)
		mockSDK.EXPECT().Get(nodeUserIDStr).Return(nodeUserAPIResponse(), nil)

		res := resourceNsxtUsers()
		d := schema.TestResourceDataRaw(t, res.Schema, nodeUserSchemaData())

		err := resourceNsxtNodeUserCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, nodeUserIDStr, d.Id())
		assert.Equal(t, nodeUserFullName, d.Get("full_name"))
		assert.Equal(t, nodeUserName, d.Get("username"))
	})

	t.Run("Create with password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNodeUserMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Createuser(gomock.Any()).Return(nodeUserAPIResponse(), nil)
		mockSDK.EXPECT().Get(nodeUserIDStr).Return(nodeUserAPIResponse(), nil)

		res := resourceNsxtUsers()
		data := nodeUserSchemaData()
		data["password"] = "SecurePass1!" //nolint:gosec
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtNodeUserCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, nodeUserIDStr, d.Id())
	})

	t.Run("Create fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNodeUserMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Createuser(gomock.Any()).Return(nsxModel.NodeUserProperties{}, errors.New("create API error"))

		res := resourceNsxtUsers()
		d := schema.TestResourceDataRaw(t, res.Schema, nodeUserSchemaData())

		err := resourceNsxtNodeUserCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "create API error")
	})
}

func TestMockResourceNsxtNodeUserRead(t *testing.T) {
	t.Run("Read success for active user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNodeUserMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(nodeUserIDStr).Return(nodeUserAPIResponse(), nil)

		res := resourceNsxtUsers()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(nodeUserIDStr)

		err := resourceNsxtNodeUserRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, nodeUserFullName, d.Get("full_name"))
		assert.Equal(t, nodeUserName, d.Get("username"))
		assert.Equal(t, int(nodeUserLastPwdChange), d.Get("last_password_change"))
		assert.Equal(t, int(nodeUserPwdFreq), d.Get("password_change_frequency"))
		assert.Equal(t, int(nodeUserPwdWarn), d.Get("password_change_warning"))
		assert.Equal(t, nodeUserPwdReset, d.Get("password_reset_required"))
		assert.Equal(t, true, d.Get("active"))
	})

	t.Run("Read success for not-activated user sets active=false", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNodeUserMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(nodeUserIDStr).Return(nodeUserNotActivatedAPIResponse(), nil)

		res := resourceNsxtUsers()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(nodeUserIDStr)

		err := resourceNsxtNodeUserRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, false, d.Get("active"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtUsers()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtNodeUserRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Read fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNodeUserMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(nodeUserIDStr).Return(nsxModel.NodeUserProperties{}, errors.New("read API error"))

		res := resourceNsxtUsers()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(nodeUserIDStr)

		err := resourceNsxtNodeUserRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read API error")
	})
}

func TestMockResourceNsxtNodeUserUpdate(t *testing.T) {
	t.Run("Update success — no status or password change", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNodeUserMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Update(nodeUserIDStr, gomock.Any()).Return(nodeUserAPIResponse(), nil)
		mockSDK.EXPECT().Get(nodeUserIDStr).Return(nodeUserAPIResponse(), nil)

		res := resourceNsxtUsers()
		data := nodeUserSchemaData()
		data["status"] = nsxModel.NodeUserProperties_STATUS_ACTIVE
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(nodeUserIDStr)

		err := resourceNsxtNodeUserUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, nodeUserFullName, d.Get("full_name"))
	})

	t.Run("Update activates a NOT_ACTIVATED user with new password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNodeUserMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Activate(nodeUserIDStr, gomock.Any()).Return(nodeUserAPIResponse(), nil)
		mockSDK.EXPECT().Update(nodeUserIDStr, gomock.Any()).Return(nodeUserAPIResponse(), nil)
		mockSDK.EXPECT().Get(nodeUserIDStr).Return(nodeUserAPIResponse(), nil)

		res := resourceNsxtUsers()
		data := nodeUserSchemaData()
		data["status"] = nsxModel.NodeUserProperties_STATUS_NOT_ACTIVATED
		data["active"] = true
		data["password"] = "NewSecurePass1!" //nolint:gosec
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(nodeUserIDStr)

		err := resourceNsxtNodeUserUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update returns error when activating without password", func(t *testing.T) {
		res := resourceNsxtUsers()
		data := nodeUserSchemaData()
		data["status"] = nsxModel.NodeUserProperties_STATUS_NOT_ACTIVATED
		data["active"] = true
		// no password set
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(nodeUserIDStr)

		err := resourceNsxtNodeUserUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must specify password")
	})

	t.Run("Update deactivates an active user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNodeUserMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Deactivate(nodeUserIDStr).Return(nodeUserNotActivatedAPIResponse(), nil)
		mockSDK.EXPECT().Update(nodeUserIDStr, gomock.Any()).Return(nodeUserNotActivatedAPIResponse(), nil)
		mockSDK.EXPECT().Get(nodeUserIDStr).Return(nodeUserNotActivatedAPIResponse(), nil)

		res := resourceNsxtUsers()
		data := nodeUserSchemaData()
		data["status"] = nsxModel.NodeUserProperties_STATUS_ACTIVE
		data["active"] = false
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(nodeUserIDStr)

		err := resourceNsxtNodeUserUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, false, d.Get("active"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtUsers()
		d := schema.TestResourceDataRaw(t, res.Schema, nodeUserSchemaData())

		err := resourceNsxtNodeUserUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Update fails when Update API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNodeUserMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Update(nodeUserIDStr, gomock.Any()).Return(nsxModel.NodeUserProperties{}, errors.New("update API error"))

		res := resourceNsxtUsers()
		data := nodeUserSchemaData()
		data["status"] = nsxModel.NodeUserProperties_STATUS_ACTIVE
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(nodeUserIDStr)

		err := resourceNsxtNodeUserUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "update API error")
	})
}

func TestMockResourceNsxtNodeUserDelete(t *testing.T) {
	origSleep := nodeUserDeleteSleep
	nodeUserDeleteSleep = 0
	defer func() { nodeUserDeleteSleep = origSleep }()

	t.Run("Delete success on first attempt", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNodeUserMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(nodeUserIDStr).Return(nil)

		res := resourceNsxtUsers()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(nodeUserIDStr)

		err := resourceNsxtNodeUserDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete success via NotFound on second attempt (SDK content-type workaround)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNodeUserMock(t, ctrl)
		defer restore()

		// First attempt returns an internal server error (the SDK content-type bug).
		// Second attempt returns NotFound, which handleDeleteError converts to nil.
		mockSDK.EXPECT().Delete(nodeUserIDStr).Return(errors.New("internal server error")).Times(1)
		mockSDK.EXPECT().Delete(nodeUserIDStr).Return(vapiErrors.NotFound{}).Times(1)

		res := resourceNsxtUsers()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(nodeUserIDStr)

		err := resourceNsxtNodeUserDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtUsers()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtNodeUserDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Delete fails when both attempts return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNodeUserMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(nodeUserIDStr).Return(errors.New("persistent delete error")).Times(2)

		res := resourceNsxtUsers()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(nodeUserIDStr)

		err := resourceNsxtNodeUserDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "persistent delete error")
	})

	t.Run("Delete sleep duration is overrideable", func(t *testing.T) {
		// Verifies that the sleep override mechanism works as expected.
		start := time.Now()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupNodeUserMock(t, ctrl)
		defer restore()

		// Two failing attempts → two sleep calls (both 0 because of the override at test start).
		mockSDK.EXPECT().Delete(nodeUserIDStr).Return(errors.New("err")).Times(2)

		res := resourceNsxtUsers()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(nodeUserIDStr)

		resourceNsxtNodeUserDelete(d, newGoMockProviderClient()) //nolint:errcheck
		assert.Less(t, time.Since(start), 500*time.Millisecond)
	})
}
