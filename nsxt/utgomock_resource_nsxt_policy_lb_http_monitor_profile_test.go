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
	lbHttpMonitorID          = "lb-http-monitor-1"
	lbHttpMonitorDisplayName = "Test LB HTTP Monitor Profile"
)

func minimalLBHttpMonitorData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbHttpMonitorDisplayName,
		"nsx_id":       lbHttpMonitorID,
	}
}

func setupLBMonitorProfileMock(t *testing.T, ctrl *gomock.Controller) (*infraMocks.MockLbMonitorProfilesClient, func()) {
	mockSDK := infraMocks.NewMockLbMonitorProfilesClient(ctrl)
	mockWrapper := &infraapi.LBMonitorProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliLbMonitorProfilesClient
	cliLbMonitorProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.LBMonitorProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliLbMonitorProfilesClient = original }
}

func lbHttpMonitorStructValue(t *testing.T) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	displayName := lbHttpMonitorDisplayName
	description := "Test LB HTTP Monitor Profile description"
	path := "/infra/lb-monitor-profiles/" + lbHttpMonitorID
	revision := int64(1)
	requestBody := "req-body"
	headerName := "X-Custom"
	headerValue := "abc"
	requestMethod := "GET"
	requestURL := "/healthz"
	requestVersion := "HTTP_VERSION_1_1"
	responseBody := "OK"
	fallCount := int64(3)
	interval := int64(5)
	monitorPort := int64(80)
	riseCount := int64(3)
	timeout := int64(15)

	obj := model.LBHttpMonitorProfile{
		DisplayName:         &displayName,
		Description:         &description,
		Path:                &path,
		Revision:            &revision,
		RequestBody:         &requestBody,
		RequestHeaders:      []model.LbHttpRequestHeader{{HeaderName: &headerName, HeaderValue: &headerValue}},
		RequestMethod:       &requestMethod,
		RequestUrl:          &requestURL,
		RequestVersion:      &requestVersion,
		ResponseBody:        &responseBody,
		ResponseStatusCodes: []int64{200, 201},
		FallCount:           &fallCount,
		Interval:            &interval,
		MonitorPort:         &monitorPort,
		RiseCount:           &riseCount,
		Timeout:             &timeout,
		ResourceType:        model.LBMonitorProfile_RESOURCE_TYPE_LBHTTPMONITORPROFILE,
	}
	val, errs := converter.ConvertToVapi(obj, model.LBHttpMonitorProfileBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func assertLBHttpMonitorFields(t *testing.T, d *schema.ResourceData) {
	t.Helper()
	assert.Equal(t, lbHttpMonitorDisplayName, d.Get("display_name"))
	assert.Equal(t, "Test LB HTTP Monitor Profile description", d.Get("description"))
	assert.Equal(t, "req-body", d.Get("request_body"))
	assert.Equal(t, "GET", d.Get("request_method"))
	assert.Equal(t, "/healthz", d.Get("request_url"))
	assert.Equal(t, "HTTP_VERSION_1_1", d.Get("request_version"))
	assert.Equal(t, "OK", d.Get("response_body"))
	assert.Equal(t, 3, d.Get("fall_count"))
	assert.Equal(t, 5, d.Get("interval"))
	assert.Equal(t, 80, d.Get("monitor_port"))
	assert.Equal(t, 3, d.Get("rise_count"))
	assert.Equal(t, 15, d.Get("timeout"))
	assert.Equal(t, "/infra/lb-monitor-profiles/"+lbHttpMonitorID, d.Get("path"))

	headers := d.Get("request_header").(*schema.Set).List()
	require.Len(t, headers, 1)
	h := headers[0].(map[string]interface{})
	assert.Equal(t, "X-Custom", h["name"])
	assert.Equal(t, "abc", h["value"])

	codes := d.Get("response_status_codes").([]interface{})
	assert.ElementsMatch(t, []interface{}{200, 201}, codes)
}

func TestMockResourceNsxtPolicyLBHttpMonitorProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpMonitorID).Return(nil, nil)

		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())

		err := resourceNsxtPolicyLBHttpMonitorProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		sv := lbHttpMonitorStructValue(t)
		mockSDK.EXPECT().Get(gomock.Any()).Return(sv, nil)

		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": lbHttpMonitorDisplayName,
		})

		err := resourceNsxtPolicyLBHttpMonitorProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assertLBHttpMonitorFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBHttpMonitorProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpMonitorID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())
		d.SetId(lbHttpMonitorID)

		err := resourceNsxtPolicyLBHttpMonitorProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpMonitorID).Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())
		d.SetId(lbHttpMonitorID)

		err := resourceNsxtPolicyLBHttpMonitorProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())

		err := resourceNsxtPolicyLBHttpMonitorProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read success", func(t *testing.T) {
		sv := lbHttpMonitorStructValue(t)
		mockSDK.EXPECT().Get(lbHttpMonitorID).Return(sv, nil)

		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())
		d.SetId(lbHttpMonitorID)

		err := resourceNsxtPolicyLBHttpMonitorProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBHttpMonitorFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBHttpMonitorProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())

		err := resourceNsxtPolicyLBHttpMonitorProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(lbHttpMonitorID, gomock.Any()).Return(nil)
		sv := lbHttpMonitorStructValue(t)
		mockSDK.EXPECT().Get(lbHttpMonitorID).Return(sv, nil)

		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())
		d.SetId(lbHttpMonitorID)

		err := resourceNsxtPolicyLBHttpMonitorProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBHttpMonitorFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBHttpMonitorProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbHttpMonitorID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())
		d.SetId(lbHttpMonitorID)

		err := resourceNsxtPolicyLBHttpMonitorProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())

		err := resourceNsxtPolicyLBHttpMonitorProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
