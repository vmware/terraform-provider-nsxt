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
	lbHttpAppID          = "lb-http-app-1"
	lbHttpAppDisplayName = "Test LB HTTP App Profile"
)

func minimalLBHttpAppData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbHttpAppDisplayName,
		"nsx_id":       lbHttpAppID,
	}
}

func lbHttpAppStructValue(t *testing.T) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	displayName := lbHttpAppDisplayName
	description := "Test LB HTTP App Profile description"
	path := "/infra/lb-app-profiles/" + lbHttpAppID
	revision := int64(1)
	httpRedirectTo := "https://example.com"
	httpRedirectToHTTPS := true
	idleTimeout := int64(30)
	requestBodySize := int64(1024)
	requestHeaderSize := int64(2048)
	responseBuffering := true
	responseHeaderSize := int64(4096)
	responseTimeout := int64(60)
	serverKeepAlive := true
	xForwardedFor := model.LBHttpProfile_X_FORWARDED_FOR_INSERT

	obj := model.LBHttpProfile{
		DisplayName:         &displayName,
		Description:         &description,
		Path:                &path,
		Revision:            &revision,
		HttpRedirectTo:      &httpRedirectTo,
		HttpRedirectToHttps: &httpRedirectToHTTPS,
		IdleTimeout:         &idleTimeout,
		RequestBodySize:     &requestBodySize,
		RequestHeaderSize:   &requestHeaderSize,
		ResponseBuffering:   &responseBuffering,
		ResponseHeaderSize:  &responseHeaderSize,
		ResponseTimeout:     &responseTimeout,
		ServerKeepAlive:     &serverKeepAlive,
		XForwardedFor:       &xForwardedFor,
		ResourceType:        model.LBAppProfile_RESOURCE_TYPE_LBHTTPPROFILE,
	}
	val, errs := converter.ConvertToVapi(obj, model.LBHttpProfileBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func assertLBHttpAppFields(t *testing.T, d *schema.ResourceData) {
	t.Helper()
	assert.Equal(t, lbHttpAppDisplayName, d.Get("display_name"))
	assert.Equal(t, "Test LB HTTP App Profile description", d.Get("description"))
	assert.Equal(t, "https://example.com", d.Get("http_redirect_to"))
	assert.Equal(t, true, d.Get("http_redirect_to_https"))
	assert.Equal(t, 30, d.Get("idle_timeout"))
	assert.Equal(t, 1024, d.Get("request_body_size"))
	assert.Equal(t, 2048, d.Get("request_header_size"))
	assert.Equal(t, true, d.Get("response_buffering"))
	assert.Equal(t, 4096, d.Get("response_header_size"))
	assert.Equal(t, 60, d.Get("response_timeout"))
	assert.Equal(t, true, d.Get("server_keep_alive"))
	assert.Equal(t, model.LBHttpProfile_X_FORWARDED_FOR_INSERT, d.Get("x_forwarded_for"))
	assert.Equal(t, "/infra/lb-app-profiles/"+lbHttpAppID, d.Get("path"))
}

func TestMockResourceNsxtPolicyLBHttpApplicationProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpAppID).Return(nil, nil)

		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())

		err := resourceNsxtPolicyLBHttpApplicationProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		sv := lbHttpAppStructValue(t)
		mockSDK.EXPECT().Get(gomock.Any()).Return(sv, nil)

		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": lbHttpAppDisplayName,
		})

		err := resourceNsxtPolicyLBHttpApplicationProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assertLBHttpAppFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBHttpApplicationProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpAppID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())
		d.SetId(lbHttpAppID)

		err := resourceNsxtPolicyLBHttpApplicationProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpAppID).Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())
		d.SetId(lbHttpAppID)

		err := resourceNsxtPolicyLBHttpApplicationProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())

		err := resourceNsxtPolicyLBHttpApplicationProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read success", func(t *testing.T) {
		sv := lbHttpAppStructValue(t)
		mockSDK.EXPECT().Get(lbHttpAppID).Return(sv, nil)

		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())
		d.SetId(lbHttpAppID)

		err := resourceNsxtPolicyLBHttpApplicationProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBHttpAppFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBHttpApplicationProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())

		err := resourceNsxtPolicyLBHttpApplicationProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(lbHttpAppID, gomock.Any()).Return(nil)
		sv := lbHttpAppStructValue(t)
		mockSDK.EXPECT().Get(lbHttpAppID).Return(sv, nil)

		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())
		d.SetId(lbHttpAppID)

		err := resourceNsxtPolicyLBHttpApplicationProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBHttpAppFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBHttpApplicationProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbHttpAppID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())
		d.SetId(lbHttpAppID)

		err := resourceNsxtPolicyLBHttpApplicationProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpAppData())

		err := resourceNsxtPolicyLBHttpApplicationProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
