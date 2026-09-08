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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"
)

// dataSourceNsxtPolicyIdpsSignatureVersionRead reads through the cliIdsSignatureVersionsClient
// package-level var (see utgomock_resource_nsxt_policy_idps_signature_version_test.go for the
// mock setup helper), so both the guard/validation branches and the mocked lookups are covered.

func TestMockDataSourceNsxtPolicyIdpsSignatureVersionReadGuard(t *testing.T) {
	t.Run("Read fails for global manager", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Global Manager")
	})
}

func TestMockDataSourceNsxtPolicyIdpsSignatureVersionReadMissingSelector(t *testing.T) {
	t.Run("Read fails when neither id nor display_name is set", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "'id' or 'display_name'")
	})
}

func TestMockDataSourceNsxtPolicyIdpsSignatureVersionReadByID(t *testing.T) {
	t.Run("Read by id succeeds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient, restore := setupIdpsSignatureVersionMock(ctrl)
		defer restore()
		id := "sig-ver-1"
		mockClient.EXPECT().Get("sig-ver-1").Return(model.IdsSignatureVersion{Id: &id}, nil)

		ds := dataSourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "sig-ver-1"})

		err := dataSourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "sig-ver-1", d.Id())
	})

	t.Run("Read by id propagates API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient, restore := setupIdpsSignatureVersionMock(ctrl)
		defer restore()
		mockClient.EXPECT().Get("sig-ver-1").Return(model.IdsSignatureVersion{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "sig-ver-1"})

		err := dataSourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockDataSourceNsxtPolicyIdpsSignatureVersionReadByName(t *testing.T) {
	t.Run("Read by display_name matches a single result", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient, restore := setupIdpsSignatureVersionMock(ctrl)
		defer restore()
		id, name := "sig-ver-1", "myversion"
		mockClient.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(
			model.IdsSignatureVersionListResult{Results: []model.IdsSignatureVersion{{Id: &id, DisplayName: &name}}}, nil,
		)

		ds := dataSourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "myversion"})

		err := dataSourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "sig-ver-1", d.Id())
	})

	t.Run("Read by display_name with no match fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient, restore := setupIdpsSignatureVersionMock(ctrl)
		defer restore()
		mockClient.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(
			model.IdsSignatureVersionListResult{}, nil,
		)

		ds := dataSourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "myversion"})

		err := dataSourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Read by display_name with multiple matches fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient, restore := setupIdpsSignatureVersionMock(ctrl)
		defer restore()
		id1, id2, name := "sig-ver-1", "sig-ver-2", "myversion"
		mockClient.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(
			model.IdsSignatureVersionListResult{Results: []model.IdsSignatureVersion{
				{Id: &id1, DisplayName: &name},
				{Id: &id2, DisplayName: &name},
			}}, nil,
		)

		ds := dataSourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "myversion"})

		err := dataSourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("List error propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient, restore := setupIdpsSignatureVersionMock(ctrl)
		defer restore()
		mockClient.EXPECT().List(nil, nil, nil, nil, nil, nil, nil).Return(
			model.IdsSignatureVersionListResult{}, vapiErrors.InternalServerError{},
		)

		ds := dataSourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "myversion"})

		err := dataSourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
