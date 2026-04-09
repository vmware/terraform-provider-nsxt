//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	apinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

var (
	vniPoolTestID          = "vni-pool-id-1"
	vniPoolTestName        = "test-vni-pool"
	vniPoolTestPath        = "/infra/vni-pool-configs/vni-pool-id-1"
	vniPoolTestDescription = "test description"
	vniPoolTestStart       = int64(70000)
	vniPoolTestEnd         = int64(71000)
)

func vniPoolModel() nsxModel.VniPoolConfig {
	return nsxModel.VniPoolConfig{
		Id:          &vniPoolTestID,
		DisplayName: &vniPoolTestName,
		Description: &vniPoolTestDescription,
		Path:        &vniPoolTestPath,
		Start:       &vniPoolTestStart,
		End:         &vniPoolTestEnd,
	}
}

func setupVniPoolDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*inframocks.MockVniPoolsClient, func()) {
	t.Helper()
	mockSDK := inframocks.NewMockVniPoolsClient(ctrl)
	wrapper := &apinfra.VniPoolConfigClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliVniPoolsClient
	cliVniPoolsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apinfra.VniPoolConfigClientContext {
		return wrapper
	}
	return mockSDK, func() { cliVniPoolsClient = orig }
}

func TestMockDataSourceNsxtPolicyVniPoolRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupVniPoolDataSourceMock(t, ctrl)
	defer restore()

	t.Run("by ID success", func(t *testing.T) {
		mockSDK.EXPECT().Get(vniPoolTestID).Return(vniPoolModel(), nil)

		ds := dataSourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": vniPoolTestID,
		})

		err := dataSourceNsxtPolicyVniPoolRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vniPoolTestID, d.Id())
		assert.Equal(t, vniPoolTestName, d.Get("display_name"))
		assert.Equal(t, int(vniPoolTestStart), d.Get("start"))
		assert.Equal(t, int(vniPoolTestEnd), d.Get("end"))
	})

	t.Run("by ID not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(vniPoolTestID).Return(nsxModel.VniPoolConfig{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": vniPoolTestID,
		})

		err := dataSourceNsxtPolicyVniPoolRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by ID API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(vniPoolTestID).Return(nsxModel.VniPoolConfig{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": vniPoolTestID,
		})

		err := dataSourceNsxtPolicyVniPoolRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading")
	})

	t.Run("by display_name single match", func(t *testing.T) {
		inc := false
		mockSDK.EXPECT().List(nil, &inc, nil, nil, nil, nil).Return(nsxModel.VniPoolConfigListResult{
			Results: []nsxModel.VniPoolConfig{vniPoolModel()},
		}, nil)

		ds := dataSourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": vniPoolTestName,
		})

		err := dataSourceNsxtPolicyVniPoolRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vniPoolTestID, d.Id())
	})

	t.Run("by display_name list error", func(t *testing.T) {
		inc := false
		mockSDK.EXPECT().List(nil, &inc, nil, nil, nil, nil).Return(nsxModel.VniPoolConfigListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": vniPoolTestName,
		})

		err := dataSourceNsxtPolicyVniPoolRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading VniPoolConfigs")
	})

	t.Run("missing id and name", func(t *testing.T) {
		ds := dataSourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyVniPoolRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ID or name")
	})

	t.Run("by name not found", func(t *testing.T) {
		inc := false
		mockSDK.EXPECT().List(nil, &inc, nil, nil, nil, nil).Return(nsxModel.VniPoolConfigListResult{
			Results: []nsxModel.VniPoolConfig{},
		}, nil)

		ds := dataSourceNsxtPolicyVniPool()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "does-not-exist",
		})

		err := dataSourceNsxtPolicyVniPoolRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})
}
