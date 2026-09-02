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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	vapiData "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	ippoolsapi "github.com/vmware/terraform-provider-nsxt/api/infra/ip_pools"
	realizedstate "github.com/vmware/terraform-provider-nsxt/api/infra/realized_state"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	ipSubnetMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/ip_pools"
	realizedmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/realized_state"
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

func setupBlockSubnetRealizedEntitiesMock(t *testing.T, ctrl *gomock.Controller) (*realizedmocks.MockRealizedEntitiesClient, func()) {
	mockSDK := realizedmocks.NewMockRealizedEntitiesClient(ctrl)
	wrapper := &realizedstate.RealizedEntityClientContext{Client: mockSDK, ClientType: utl.Local}
	original := cliRealizedEntitiesClient
	cliRealizedEntitiesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *realizedstate.RealizedEntityClientContext {
		return wrapper
	}
	return mockSDK, func() { cliRealizedEntitiesClient = original }
}

// blockSubnetStructValue builds a real *data.StructValue for an IpAddressPoolBlockSubnet,
// mimicking what the SDK Get() call returns on the wire.
func blockSubnetStructValue(t *testing.T, displayName, description string, size int64, autoAssignGateway bool) *vapiData.StructValue {
	converter := bindings.NewTypeConverter()
	obj := nsxModel.IpAddressPoolBlockSubnet{
		Id:                &blockSubnetID,
		DisplayName:       &displayName,
		Description:       &description,
		Size:              &size,
		AutoAssignGateway: &autoAssignGateway,
		ResourceType:      "IpAddressPoolBlockSubnet",
		IpBlockPath:       &blockSubnetBlock,
	}
	dataValue, errs := converter.ConvertToVapi(obj, nsxModel.IpAddressPoolBlockSubnetBindingType())
	require.Nil(t, errs)
	return dataValue.(*vapiData.StructValue)
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

func TestMockResourceNsxtPolicyIPPoolBlockSubnetReadSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupBlockSubnetMock(t, ctrl)
	defer restore()

	t.Run("Read success sets all fields", func(t *testing.T) {
		sv := blockSubnetStructValue(t, "Test Block Subnet", "Test block subnet", 16, true)
		mockSDK.EXPECT().Get(blockSubnetPoolID, blockSubnetID).Return(sv, nil)

		res := resourceNsxtPolicyIPPoolBlockSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBlockSubnetData())
		d.SetId(blockSubnetID)

		err := resourceNsxtPolicyIPPoolBlockSubnetRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "Test Block Subnet", d.Get("display_name"))
		assert.Equal(t, "Test block subnet", d.Get("description"))
		assert.Equal(t, blockSubnetID, d.Get("nsx_id"))
		assert.Equal(t, blockSubnetPoolPath, d.Get("pool_path"))
		assert.Equal(t, blockSubnetBlock, d.Get("block_path"))
		assert.Equal(t, 16, d.Get("size"))
		assert.True(t, d.Get("auto_assign_gateway").(bool))
	})
}

func TestMockResourceNsxtPolicyIPPoolBlockSubnetCreateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupBlockSubnetMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		sv := blockSubnetStructValue(t, "Test Block Subnet", "Test block subnet", 16, true)

		gomock.InOrder(
			mockSDK.EXPECT().Get(blockSubnetPoolID, blockSubnetID).Return(nil, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(blockSubnetPoolID, blockSubnetID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(blockSubnetPoolID, blockSubnetID).Return(sv, nil),
		)

		res := resourceNsxtPolicyIPPoolBlockSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBlockSubnetData())

		err := resourceNsxtPolicyIPPoolBlockSubnetCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, blockSubnetID, d.Id())
		assert.Equal(t, blockSubnetID, d.Get("nsx_id"))
	})

	t.Run("Create fails when Patch API errors", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(blockSubnetPoolID, blockSubnetID).Return(nil, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(blockSubnetPoolID, blockSubnetID, gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtPolicyIPPoolBlockSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBlockSubnetData())

		err := resourceNsxtPolicyIPPoolBlockSubnetCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPPoolBlockSubnetUpdateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupBlockSubnetMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		sv := blockSubnetStructValue(t, "Updated name", "Test block subnet", 16, true)

		gomock.InOrder(
			mockSDK.EXPECT().Patch(blockSubnetPoolID, blockSubnetID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(blockSubnetPoolID, blockSubnetID).Return(sv, nil),
		)

		res := resourceNsxtPolicyIPPoolBlockSubnet()
		data := minimalBlockSubnetData()
		data["display_name"] = "Updated name"
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(blockSubnetID)

		err := resourceNsxtPolicyIPPoolBlockSubnetUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "Updated name", d.Get("display_name"))
	})

	t.Run("Update fails when Patch API errors", func(t *testing.T) {
		mockSDK.EXPECT().Patch(blockSubnetPoolID, blockSubnetID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIPPoolBlockSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBlockSubnetData())
		d.SetId(blockSubnetID)

		err := resourceNsxtPolicyIPPoolBlockSubnetUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPPoolBlockSubnetDeleteSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupBlockSubnetMock(t, ctrl)
	defer restore()
	mockRealizedSDK, restoreRealized := setupBlockSubnetRealizedEntitiesMock(t, ctrl)
	defer restoreRealized()

	t.Run("Delete success waits for realization to clear", func(t *testing.T) {
		mockSDK.EXPECT().Delete(blockSubnetPoolID, blockSubnetID, gomock.Any()).Return(nil)
		mockRealizedSDK.EXPECT().
			List(gomock.Any(), gomock.Any()).
			Return(nsxModel.GenericPolicyRealizedResourceListResult{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIPPoolBlockSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBlockSubnetData())
		d.SetId(blockSubnetID)

		err := resourceNsxtPolicyIPPoolBlockSubnetDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when Delete API errors", func(t *testing.T) {
		mockSDK.EXPECT().Delete(blockSubnetPoolID, blockSubnetID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIPPoolBlockSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalBlockSubnetData())
		d.SetId(blockSubnetID)

		err := resourceNsxtPolicyIPPoolBlockSubnetDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
