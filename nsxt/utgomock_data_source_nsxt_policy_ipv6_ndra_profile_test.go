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
	ndraProfileTestID   = "ndra-profile-id-1"
	ndraProfileTestName = "test-ndra-profile"
	ndraProfileTestPath = "/infra/ipv6-ndra-profiles/ndra-profile-id-1"
	ndraProfileTestDesc = "test description"
)

func ipv6NdraProfileModel() nsxModel.Ipv6NdraProfile {
	return nsxModel.Ipv6NdraProfile{
		Id:          &ndraProfileTestID,
		DisplayName: &ndraProfileTestName,
		Description: &ndraProfileTestDesc,
		Path:        &ndraProfileTestPath,
	}
}

func setupIpv6NdraProfileDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*inframocks.MockIpv6NdraProfilesClient, func()) {
	t.Helper()
	mockSDK := inframocks.NewMockIpv6NdraProfilesClient(ctrl)
	wrapper := &apinfra.Ipv6NdraProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliIpv6NdraProfilesClient
	cliIpv6NdraProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apinfra.Ipv6NdraProfileClientContext {
		return wrapper
	}
	return mockSDK, func() { cliIpv6NdraProfilesClient = orig }
}

func TestMockDataSourceNsxtPolicyIpv6NdraProfileRead(t *testing.T) {
	t.Run("by ID success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupIpv6NdraProfileDataSourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(ndraProfileTestID).Return(ipv6NdraProfileModel(), nil)

		ds := dataSourceNsxtPolicyIpv6NdraProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": ndraProfileTestID,
		})

		err := dataSourceNsxtPolicyIpv6NdraProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ndraProfileTestID, d.Id())
		assert.Equal(t, ndraProfileTestName, d.Get("display_name"))
		assert.Equal(t, ndraProfileTestPath, d.Get("path"))
	})

	t.Run("by ID not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupIpv6NdraProfileDataSourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(ndraProfileTestID).Return(nsxModel.Ipv6NdraProfile{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyIpv6NdraProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": ndraProfileTestID,
		})

		err := dataSourceNsxtPolicyIpv6NdraProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by display_name single match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupIpv6NdraProfileDataSourceMock(t, ctrl)
		defer restore()

		inc := false
		mockSDK.EXPECT().List(nil, &inc, nil, nil, nil, nil).Return(nsxModel.Ipv6NdraProfileListResult{
			Results: []nsxModel.Ipv6NdraProfile{ipv6NdraProfileModel()},
		}, nil)

		ds := dataSourceNsxtPolicyIpv6NdraProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": ndraProfileTestName,
		})

		err := dataSourceNsxtPolicyIpv6NdraProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ndraProfileTestID, d.Id())
	})

	t.Run("by display_name list error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupIpv6NdraProfileDataSourceMock(t, ctrl)
		defer restore()

		inc := false
		mockSDK.EXPECT().List(nil, &inc, nil, nil, nil, nil).Return(nsxModel.Ipv6NdraProfileListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIpv6NdraProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": ndraProfileTestName,
		})

		err := dataSourceNsxtPolicyIpv6NdraProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error listing Ipv6NdraProfiles")
	})

	t.Run("by name not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupIpv6NdraProfileDataSourceMock(t, ctrl)
		defer restore()

		inc := false
		mockSDK.EXPECT().List(nil, &inc, nil, nil, nil, nil).Return(nsxModel.Ipv6NdraProfileListResult{
			Results: []nsxModel.Ipv6NdraProfile{},
		}, nil)

		ds := dataSourceNsxtPolicyIpv6NdraProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "does-not-exist",
		})

		err := dataSourceNsxtPolicyIpv6NdraProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("missing id and name", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIpv6NdraProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIpv6NdraProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "id or display_name must be specified")
	})
}
