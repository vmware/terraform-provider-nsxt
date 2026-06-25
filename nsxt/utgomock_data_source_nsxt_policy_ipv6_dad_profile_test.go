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
	dadProfileTestID   = "dad-profile-id-1"
	dadProfileTestName = "test-dad-profile"
	dadProfileTestPath = "/infra/ipv6-dad-profiles/dad-profile-id-1"
	dadProfileTestDesc = "test description"
)

func ipv6DadProfileModel() nsxModel.Ipv6DadProfile {
	return nsxModel.Ipv6DadProfile{
		Id:          &dadProfileTestID,
		DisplayName: &dadProfileTestName,
		Description: &dadProfileTestDesc,
		Path:        &dadProfileTestPath,
	}
}

func setupIpv6DadProfileDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*inframocks.MockIpv6DadProfilesClient, func()) {
	t.Helper()
	mockSDK := inframocks.NewMockIpv6DadProfilesClient(ctrl)
	wrapper := &apinfra.Ipv6DadProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliIpv6DadProfilesClient
	cliIpv6DadProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apinfra.Ipv6DadProfileClientContext {
		return wrapper
	}
	return mockSDK, func() { cliIpv6DadProfilesClient = orig }
}

func TestMockDataSourceNsxtPolicyIpv6DadProfileRead(t *testing.T) {
	t.Run("by ID success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupIpv6DadProfileDataSourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(dadProfileTestID).Return(ipv6DadProfileModel(), nil)

		ds := dataSourceNsxtPolicyIpv6DadProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": dadProfileTestID,
		})

		err := dataSourceNsxtPolicyIpv6DadProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dadProfileTestID, d.Id())
		assert.Equal(t, dadProfileTestName, d.Get("display_name"))
		assert.Equal(t, dadProfileTestPath, d.Get("path"))
	})

	t.Run("by ID not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupIpv6DadProfileDataSourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(dadProfileTestID).Return(nsxModel.Ipv6DadProfile{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyIpv6DadProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": dadProfileTestID,
		})

		err := dataSourceNsxtPolicyIpv6DadProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by display_name single match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupIpv6DadProfileDataSourceMock(t, ctrl)
		defer restore()

		inc := false
		mockSDK.EXPECT().List(nil, &inc, nil, nil, nil, nil).Return(nsxModel.Ipv6DadProfileListResult{
			Results: []nsxModel.Ipv6DadProfile{ipv6DadProfileModel()},
		}, nil)

		ds := dataSourceNsxtPolicyIpv6DadProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": dadProfileTestName,
		})

		err := dataSourceNsxtPolicyIpv6DadProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dadProfileTestID, d.Id())
	})

	t.Run("by display_name list error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupIpv6DadProfileDataSourceMock(t, ctrl)
		defer restore()

		inc := false
		mockSDK.EXPECT().List(nil, &inc, nil, nil, nil, nil).Return(nsxModel.Ipv6DadProfileListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIpv6DadProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": dadProfileTestName,
		})

		err := dataSourceNsxtPolicyIpv6DadProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error listing Ipv6DadProfiles")
	})

	t.Run("by name not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupIpv6DadProfileDataSourceMock(t, ctrl)
		defer restore()

		inc := false
		mockSDK.EXPECT().List(nil, &inc, nil, nil, nil, nil).Return(nsxModel.Ipv6DadProfileListResult{
			Results: []nsxModel.Ipv6DadProfile{},
		}, nil)

		ds := dataSourceNsxtPolicyIpv6DadProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "does-not-exist",
		})

		err := dataSourceNsxtPolicyIpv6DadProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("missing id and name", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIpv6DadProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIpv6DadProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "id or display_name must be specified")
	})
}
