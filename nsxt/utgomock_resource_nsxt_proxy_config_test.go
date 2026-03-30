// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/nsx/proxy/ConfigClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/proxy/ConfigClient.go ConfigClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	mpModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/proxy"

	proxymocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/proxy"
)

var (
	proxyID          = "TelemetryConfigIdentifier"
	proxyDisplayName = "NSX Proxy Config"
	proxyRevision    = int64(1)
	proxyHost        = "proxy.example.com"
	proxyPort        = int64(3128)
	proxyScheme      = "HTTP"
	proxyEnabled     = false
)

func proxyConfigAPIResponse() mpModel.Proxy {
	testURL := "https://www.vmware.com"
	return mpModel.Proxy{
		Id:                &proxyID,
		DisplayName:       &proxyDisplayName,
		Revision:          &proxyRevision,
		Host:              &proxyHost,
		Port:              &proxyPort,
		Scheme:            &proxyScheme,
		Enabled:           &proxyEnabled,
		TestConnectionUrl: &testURL,
	}
}

func minimalProxyConfigData() map[string]interface{} {
	return map[string]interface{}{
		"enabled":             proxyEnabled,
		"scheme":              proxyScheme,
		"host":                proxyHost,
		"port":                int(proxyPort),
		"test_connection_url": "https://www.vmware.com",
	}
}

func setupProxyConfigMock(ctrl *gomock.Controller) (*proxymocks.MockConfigClient, func()) {
	mockSDK := proxymocks.NewMockConfigClient(ctrl)

	original := cliProxyConfigClient
	cliProxyConfigClient = func(_ vapiProtocolClient.Connector) proxy.ConfigClient {
		return mockSDK
	}
	return mockSDK, func() { cliProxyConfigClient = original }
}

func TestMockResourceNsxtProxyConfigCreate(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Create success (calls update)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupProxyConfigMock(ctrl)
		defer restore()

		// Create calls Update which calls Get then Update, then Read which calls Get
		gomock.InOrder(
			mockSDK.EXPECT().Get().Return(proxyConfigAPIResponse(), nil),
			mockSDK.EXPECT().Update(gomock.Any()).Return(proxyConfigAPIResponse(), nil),
			mockSDK.EXPECT().Get().Return(proxyConfigAPIResponse(), nil),
		)

		res := resourceNsxtProxyConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProxyConfigData())

		err := resourceNsxtProxyConfigCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, proxyID, d.Id())
		assert.Equal(t, proxyHost, d.Get("host"))
	})

	t.Run("Create fails on API error during update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupProxyConfigMock(ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Get().Return(proxyConfigAPIResponse(), nil),
			mockSDK.EXPECT().Update(gomock.Any()).Return(mpModel.Proxy{}, errors.New("update failed")),
		)

		res := resourceNsxtProxyConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProxyConfigData())

		err := resourceNsxtProxyConfigCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtProxyConfigRead(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupProxyConfigMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(proxyConfigAPIResponse(), nil)

		res := resourceNsxtProxyConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProxyConfigData())
		d.SetId(proxyID)

		err := resourceNsxtProxyConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, proxyHost, d.Get("host"))
		assert.Equal(t, proxyScheme, d.Get("scheme"))
		assert.Equal(t, int(proxyPort), d.Get("port"))
	})

	t.Run("Read API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupProxyConfigMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(mpModel.Proxy{}, errors.New("API error"))

		res := resourceNsxtProxyConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProxyConfigData())
		d.SetId(proxyID)

		err := resourceNsxtProxyConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtProxyConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProxyConfigData())

		err := resourceNsxtProxyConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Proxy Config ID")
	})
}

func TestMockResourceNsxtProxyConfigUpdate(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupProxyConfigMock(ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Get().Return(proxyConfigAPIResponse(), nil),
			mockSDK.EXPECT().Update(gomock.Any()).Return(proxyConfigAPIResponse(), nil),
			mockSDK.EXPECT().Get().Return(proxyConfigAPIResponse(), nil),
		)

		res := resourceNsxtProxyConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProxyConfigData())
		d.SetId(proxyID)

		err := resourceNsxtProxyConfigUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when host required but missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupProxyConfigMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(proxyConfigAPIResponse(), nil)

		enabled := true
		data := minimalProxyConfigData()
		data["enabled"] = enabled
		data["host"] = ""

		res := resourceNsxtProxyConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(proxyID)

		err := resourceNsxtProxyConfigUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "host is required")
	})
}

func TestMockResourceNsxtProxyConfigDelete(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Delete success (disables proxy)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupProxyConfigMock(ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Get().Return(proxyConfigAPIResponse(), nil),
			mockSDK.EXPECT().Update(gomock.Any()).Return(proxyConfigAPIResponse(), nil),
		)

		res := resourceNsxtProxyConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProxyConfigData())
		d.SetId(proxyID)

		err := resourceNsxtProxyConfigDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Delete succeeds even when Get fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupProxyConfigMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(mpModel.Proxy{}, errors.New("API error"))

		res := resourceNsxtProxyConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProxyConfigData())
		d.SetId(proxyID)

		err := resourceNsxtProxyConfigDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}
