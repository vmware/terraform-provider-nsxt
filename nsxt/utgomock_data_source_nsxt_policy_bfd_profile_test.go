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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	apinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

var (
	bfdProfileTestID   = "bfd-profile-id-1"
	bfdProfileTestName = "test-bfd-profile"
	bfdProfileTestPath = "/infra/bfd-profiles/bfd-profile-id-1"
	bfdProfileTestDesc = "test description"
)

func bfdProfileModel() nsxModel.BfdProfile {
	return nsxModel.BfdProfile{
		Id:          &bfdProfileTestID,
		DisplayName: &bfdProfileTestName,
		Description: &bfdProfileTestDesc,
		Path:        &bfdProfileTestPath,
	}
}

func setupBfdProfileDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*inframocks.MockBfdProfilesClient, func()) {
	t.Helper()
	mockSDK := inframocks.NewMockBfdProfilesClient(ctrl)
	wrapper := &apinfra.BfdProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliBfdProfilesClient
	cliBfdProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apinfra.BfdProfileClientContext {
		return wrapper
	}
	return mockSDK, func() { cliBfdProfilesClient = orig }
}

func TestMockDataSourceNsxtPolicyBfdProfileRead(t *testing.T) {
	t.Run("by ID success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupBfdProfileDataSourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(bfdProfileTestID).Return(bfdProfileModel(), nil)

		ds := dataSourceNsxtPolicyBfdProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": bfdProfileTestID,
		})

		err := dataSourceNsxtPolicyBfdProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, bfdProfileTestID, d.Id())
		assert.Equal(t, bfdProfileTestName, d.Get("display_name"))
		assert.Equal(t, bfdProfileTestPath, d.Get("path"))
	})

	t.Run("by ID not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupBfdProfileDataSourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(bfdProfileTestID).Return(nsxModel.BfdProfile{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyBfdProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": bfdProfileTestID,
		})

		err := dataSourceNsxtPolicyBfdProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by ID API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupBfdProfileDataSourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(bfdProfileTestID).Return(nsxModel.BfdProfile{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyBfdProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": bfdProfileTestID,
		})

		err := dataSourceNsxtPolicyBfdProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error reading BfdProfile")
	})

	t.Run("by display_name single match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupBfdProfileDataSourceMock(t, ctrl)
		defer restore()

		inc := false
		mockSDK.EXPECT().List(nil, &inc, nil, nil, nil, nil).Return(nsxModel.BfdProfileListResult{
			Results: []nsxModel.BfdProfile{bfdProfileModel()},
		}, nil)

		ds := dataSourceNsxtPolicyBfdProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": bfdProfileTestName,
		})

		err := dataSourceNsxtPolicyBfdProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, bfdProfileTestID, d.Id())
	})

	t.Run("by display_name list error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupBfdProfileDataSourceMock(t, ctrl)
		defer restore()

		inc := false
		mockSDK.EXPECT().List(nil, &inc, nil, nil, nil, nil).Return(nsxModel.BfdProfileListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyBfdProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": bfdProfileTestName,
		})

		err := dataSourceNsxtPolicyBfdProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error listing BfdProfiles")
	})

	t.Run("by name not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupBfdProfileDataSourceMock(t, ctrl)
		defer restore()

		inc := false
		mockSDK.EXPECT().List(nil, &inc, nil, nil, nil, nil).Return(nsxModel.BfdProfileListResult{
			Results: []nsxModel.BfdProfile{},
		}, nil)

		ds := dataSourceNsxtPolicyBfdProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "does-not-exist",
		})

		err := dataSourceNsxtPolicyBfdProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("missing id and name", func(t *testing.T) {
		ds := dataSourceNsxtPolicyBfdProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyBfdProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "id or display_name must be specified")
	})
}
