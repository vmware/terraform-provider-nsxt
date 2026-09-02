//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// Reuses the VersionClient mock and setupVersionMock helper already defined
// for getNSXVersion in utgomock_utils_version_test.go:
// mockgen -destination=mocks/nsx/node/VersionClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/node/VersionClient.go VersionClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"go.uber.org/mock/gomock"
)

func TestMockDataSourceNsxtManagerInfoRead(t *testing.T) {
	t.Run("success sets version and id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVersionMock(ctrl)
		defer restore()

		nodeVersion := "4.1.0.0.0"
		productVersion := "4.1.0"
		mockSDK.EXPECT().Get().Return(nsxModel.NodeVersion{NodeVersion: &nodeVersion, ProductVersion: &productVersion}, nil)

		ds := dataSourceNsxtManagerInfo()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtManagerInfoRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, nsxVersionID, d.Id())
		assert.Equal(t, productVersion, d.Get("version"))
	})

	t.Run("API error is wrapped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVersionMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(nsxModel.NodeVersion{}, errors.New("connection error"))

		ds := dataSourceNsxtManagerInfo()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtManagerInfoRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to retrieve NSX version")
	})
}
