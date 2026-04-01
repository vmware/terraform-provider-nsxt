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

	ippoolsapi "github.com/vmware/terraform-provider-nsxt/api/infra/ip_pools"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	ipSubnetMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/ip_pools"
)

var (
	blockSubnetID       = "block-subnet-1"
	blockSubnetPoolPath = "/infra/ip-pools/pool-1"
	blockSubnetPoolID   = "pool-1"
	blockSubnetBlock    = "/infra/ip-blocks/block-1"
)

func minimalBlockSubnetData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": "Test Block Subnet",
		"description":  "Test block subnet",
		"nsx_id":       blockSubnetID,
		"pool_path":    blockSubnetPoolPath,
		"block_path":   blockSubnetBlock,
		"size":         16,
	}
}

func setupBlockSubnetMock(t *testing.T, ctrl *gomock.Controller) (*ipSubnetMocks.MockIpSubnetsClient, func()) {
	mockSDK := ipSubnetMocks.NewMockIpSubnetsClient(ctrl)
	mockWrapper := &ippoolsapi.StructValueClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliIpSubnetsClient
	cliIpSubnetsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *ippoolsapi.StructValueClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliIpSubnetsClient = original }
}

func TestMockResourceNsxtPolicyIPPoolBlockSubnetCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupBlockSubnetMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(blockSubnetPoolID, blockSubnetID).Return(nil, nil)

		res := resourceNsxtPolicyIPPoolBlockSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBlockSubnetData())

		err := resourceNsxtPolicyIPPoolBlockSubnetCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyIPPoolBlockSubnetRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupBlockSubnetMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(blockSubnetPoolID, blockSubnetID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIPPoolBlockSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBlockSubnetData())
		d.SetId(blockSubnetID)

		err := resourceNsxtPolicyIPPoolBlockSubnetRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(blockSubnetPoolID, blockSubnetID).Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIPPoolBlockSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBlockSubnetData())
		d.SetId(blockSubnetID)

		err := resourceNsxtPolicyIPPoolBlockSubnetRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPPoolBlockSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBlockSubnetData())

		err := resourceNsxtPolicyIPPoolBlockSubnetRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPPoolBlockSubnetUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupBlockSubnetMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPPoolBlockSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBlockSubnetData())

		err := resourceNsxtPolicyIPPoolBlockSubnetUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPPoolBlockSubnetDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupBlockSubnetMock(t, ctrl)
	defer restore()

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPPoolBlockSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBlockSubnetData())

		err := resourceNsxtPolicyIPPoolBlockSubnetDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
