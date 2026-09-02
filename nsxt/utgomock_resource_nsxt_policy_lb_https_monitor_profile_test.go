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
	lbHttpsMonitorID          = "lb-https-monitor-1"
	lbHttpsMonitorDisplayName = "Test LB HTTPS Monitor Profile"
)

func minimalLBHttpsMonitorData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbHttpsMonitorDisplayName,
		"nsx_id":       lbHttpsMonitorID,
	}
}

func lbHttpsMonitorStructValue(t *testing.T) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	displayName := lbHttpsMonitorDisplayName
	description := "Test LB HTTPS Monitor Profile description"
	path := "/infra/lb-monitor-profiles/" + lbHttpsMonitorID
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
	certChainDepth := int64(3)
	clientCertPath := "/infra/certificates/client-cert"
	serverAuth := model.LBServerSslProfileBinding_SERVER_AUTH_REQUIRED
	sslProfilePath := "/infra/ls-ssl-profiles/default"

	obj := model.LBHttpsMonitorProfile{
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
		ServerSslProfileBinding: &model.LBServerSslProfileBinding{
			CertificateChainDepth: &certChainDepth,
			ClientCertificatePath: &clientCertPath,
			ServerAuth:            &serverAuth,
			SslProfilePath:        &sslProfilePath,
		},
		FallCount:    &fallCount,
		Interval:     &interval,
		MonitorPort:  &monitorPort,
		RiseCount:    &riseCount,
		Timeout:      &timeout,
		ResourceType: model.LBMonitorProfile_RESOURCE_TYPE_LBHTTPSMONITORPROFILE,
	}
	val, errs := converter.ConvertToVapi(obj, model.LBHttpsMonitorProfileBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func assertLBHttpsMonitorFields(t *testing.T, d *schema.ResourceData) {
	t.Helper()
	assert.Equal(t, lbHttpsMonitorDisplayName, d.Get("display_name"))
	assert.Equal(t, "Test LB HTTPS Monitor Profile description", d.Get("description"))
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
	assert.Equal(t, "/infra/lb-monitor-profiles/"+lbHttpsMonitorID, d.Get("path"))

	headers := d.Get("request_header").(*schema.Set).List()
	require.Len(t, headers, 1)
	h := headers[0].(map[string]interface{})
	assert.Equal(t, "X-Custom", h["name"])
	assert.Equal(t, "abc", h["value"])

	codes := d.Get("response_status_codes").([]interface{})
	assert.ElementsMatch(t, []interface{}{200, 201}, codes)

	sslList := d.Get("server_ssl").([]interface{})
	require.Len(t, sslList, 1)
	sslElem := sslList[0].(map[string]interface{})
	assert.Equal(t, 3, sslElem["certificate_chain_depth"])
	assert.Equal(t, "/infra/certificates/client-cert", sslElem["client_certificate_path"])
	assert.Equal(t, model.LBServerSslProfileBinding_SERVER_AUTH_REQUIRED, sslElem["server_auth"])
	assert.Equal(t, "/infra/ls-ssl-profiles/default", sslElem["ssl_profile_path"])
}

func TestMockResourceNsxtPolicyLBHttpsMonitorProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpsMonitorID).Return(nil, nil)

		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpsMonitorData())

		err := resourceNsxtPolicyLBHttpsMonitorProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		sv := lbHttpsMonitorStructValue(t)
		mockSDK.EXPECT().Get(gomock.Any()).Return(sv, nil)

		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name": lbHttpsMonitorDisplayName,
		})

		err := resourceNsxtPolicyLBHttpsMonitorProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assertLBHttpsMonitorFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBHttpsMonitorProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpsMonitorID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpsMonitorData())
		d.SetId(lbHttpsMonitorID)

		err := resourceNsxtPolicyLBHttpsMonitorProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpsMonitorData())

		err := resourceNsxtPolicyLBHttpsMonitorProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read success", func(t *testing.T) {
		sv := lbHttpsMonitorStructValue(t)
		mockSDK.EXPECT().Get(lbHttpsMonitorID).Return(sv, nil)

		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpsMonitorData())
		d.SetId(lbHttpsMonitorID)

		err := resourceNsxtPolicyLBHttpsMonitorProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBHttpsMonitorFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBHttpsMonitorProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpsMonitorData())

		err := resourceNsxtPolicyLBHttpsMonitorProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update success", func(t *testing.T) {
		mockSDK.EXPECT().Patch(lbHttpsMonitorID, gomock.Any()).Return(nil)
		sv := lbHttpsMonitorStructValue(t)
		mockSDK.EXPECT().Get(lbHttpsMonitorID).Return(sv, nil)

		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpsMonitorData())
		d.SetId(lbHttpsMonitorID)

		err := resourceNsxtPolicyLBHttpsMonitorProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assertLBHttpsMonitorFields(t, d)
	})
}

func TestMockResourceNsxtPolicyLBHttpsMonitorProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbHttpsMonitorID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpsMonitorData())
		d.SetId(lbHttpsMonitorID)

		err := resourceNsxtPolicyLBHttpsMonitorProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpsMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpsMonitorData())

		err := resourceNsxtPolicyLBHttpsMonitorProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
