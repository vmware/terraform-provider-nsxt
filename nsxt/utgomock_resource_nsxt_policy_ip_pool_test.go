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

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	infraMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

var (
	ipPoolID          = "ip-pool-1"
	ipPoolDisplayName = "Test IP Pool"
	ipPoolDescription = "Test ip pool"
	ipPoolRevision    = int64(1)
	ipPoolPath        = "/infra/ip-pools/ip-pool-1"
)

func ipPoolAPIResponse() nsxModel.IpAddressPool {
	return nsxModel.IpAddressPool{
		Id:          &ipPoolID,
		DisplayName: &ipPoolDisplayName,
		Description: &ipPoolDescription,
		Revision:    &ipPoolRevision,
		Path:        &ipPoolPath,
	}
}

func minimalIPPoolData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": ipPoolDisplayName,
		"description":  ipPoolDescription,
		"nsx_id":       ipPoolID,
	}
}

func setupIPPoolMock(t *testing.T, ctrl *gomock.Controller) (*infraMocks.MockIpPoolsClient, func()) {
	mockSDK := infraMocks.NewMockIpPoolsClient(ctrl)
	mockWrapper := &infraapi.IpAddressPoolClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliIpPoolsClient
	cliIpPoolsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.IpAddressPoolClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliIpPoolsClient = original }
}

func TestMockResourceNsxtPolicyIPPoolCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPPoolMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(ipPoolID).Return(nsxModel.IpAddressPool{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(ipPoolID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(ipPoolID).Return(ipPoolAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIPPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPPoolData())

		err := resourceNsxtPolicyIPPoolCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipPoolID, d.Id())
		assert.Equal(t, ipPoolDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipPoolID).Return(ipPoolAPIResponse(), nil)

		res := resourceNsxtPolicyIPPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPPoolData())

		err := resourceNsxtPolicyIPPoolCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPPoolRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPPoolMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipPoolID).Return(ipPoolAPIResponse(), nil)

		res := resourceNsxtPolicyIPPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPPoolData())
		d.SetId(ipPoolID)

		err := resourceNsxtPolicyIPPoolRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipPoolDisplayName, d.Get("display_name"))
		assert.Equal(t, ipPoolID, d.Id())
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipPoolID).Return(nsxModel.IpAddressPool{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIPPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPPoolData())
		d.SetId(ipPoolID)

		err := resourceNsxtPolicyIPPoolRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPPoolData())

		err := resourceNsxtPolicyIPPoolRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPPoolUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPPoolMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(ipPoolID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(ipPoolID).Return(ipPoolAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIPPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPPoolData())
		d.SetId(ipPoolID)

		err := resourceNsxtPolicyIPPoolUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPPoolData())

		err := resourceNsxtPolicyIPPoolUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPPoolDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPPoolMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(ipPoolID).Return(nil)

		res := resourceNsxtPolicyIPPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPPoolData())
		d.SetId(ipPoolID)

		err := resourceNsxtPolicyIPPoolDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPPoolData())

		err := resourceNsxtPolicyIPPoolDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
