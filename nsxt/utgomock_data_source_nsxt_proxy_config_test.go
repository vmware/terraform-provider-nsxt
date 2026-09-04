//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// Reuses the existing ConfigClient mock generated for resource_nsxt_proxy_config.go
// (see utgomock_resource_nsxt_proxy_config_test.go):
// mockgen -destination=mocks/nsx/proxy/ConfigClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/proxy/ConfigClient.go ConfigClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mpModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"go.uber.org/mock/gomock"
)

func TestMockDataSourceNsxtProxyConfigRead(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupProxyConfigMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(proxyConfigAPIResponse(), nil)

		ds := dataSourceNsxtProxyConfig()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtProxyConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, proxyID, d.Id())
		assert.Equal(t, proxyDisplayName, d.Get("display_name"))
		assert.Equal(t, proxyHost, d.Get("host"))
		assert.Equal(t, proxyScheme, d.Get("scheme"))
		assert.Equal(t, int(proxyPort), d.Get("port"))
		assert.Equal(t, proxyEnabled, d.Get("enabled"))
		assert.Equal(t, "/api/v1/proxy/config", d.Get("path"))
	})

	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupProxyConfigMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(mpModel.Proxy{}, errors.New("API error"))

		ds := dataSourceNsxtProxyConfig()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtProxyConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
