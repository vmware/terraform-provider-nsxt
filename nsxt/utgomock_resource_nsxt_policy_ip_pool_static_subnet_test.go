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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	ippoolsapi "github.com/vmware/terraform-provider-nsxt/api/infra/ip_pools"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	ipSubnetMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/ip_pools"
)

var (
	staticSubnetID       = "static-subnet-1"
	staticSubnetPoolPath = "/infra/ip-pools/pool-2"
	staticSubnetPoolID   = "pool-2"
	staticSubnetCIDR     = "10.0.0.0/24"
)

func minimalStaticSubnetData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": "Test Static Subnet",
		"description":  "Test static subnet",
		"nsx_id":       staticSubnetID,
		"pool_path":    staticSubnetPoolPath,
		"cidr":         staticSubnetCIDR,
		"allocation_range": []interface{}{
			map[string]interface{}{
				"start": "10.0.0.2",
				"end":   "10.0.0.254",
			},
		},
	}
}

func staticSubnetAPIStructValue(t *testing.T) *data.StructValue {
	displayName := "Test Static Subnet"
	description := "Test static subnet"
	gateway := "10.0.0.1"
	obj := model.IpAddressPoolStaticSubnet{
		Id:           &staticSubnetID,
		DisplayName:  &displayName,
		Description:  &description,
		ResourceType: "IpAddressPoolStaticSubnet",
		Cidr:         &staticSubnetCIDR,
		GatewayIp:    &gateway,
	}
	converter := bindings.NewTypeConverter()
	dataValue, errs := converter.ConvertToVapi(obj, model.IpAddressPoolStaticSubnetBindingType())
	require.Empty(t, errs)
	return dataValue.(*data.StructValue)
}

func setupStaticSubnetMock(t *testing.T, ctrl *gomock.Controller) (*ipSubnetMocks.MockIpSubnetsClient, func()) {
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

func TestMockResourceNsxtPolicyIPPoolStaticSubnetCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupStaticSubnetMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(staticSubnetPoolID, staticSubnetID).Return(nil, nil)

		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())

		err := resourceNsxtPolicyIPPoolStaticSubnetCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(staticSubnetPoolID, staticSubnetID).Return(nil, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(staticSubnetPoolID, staticSubnetID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(staticSubnetPoolID, staticSubnetID).Return(staticSubnetAPIStructValue(t), nil),
		)

		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())

		err := resourceNsxtPolicyIPPoolStaticSubnetCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, staticSubnetID, d.Id())
	})

	t.Run("Create propagates a Patch error", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(staticSubnetPoolID, staticSubnetID).Return(nil, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(staticSubnetPoolID, staticSubnetID, gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())

		err := resourceNsxtPolicyIPPoolStaticSubnetCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPPoolStaticSubnetRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupStaticSubnetMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(staticSubnetPoolID, staticSubnetID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())
		d.SetId(staticSubnetID)

		err := resourceNsxtPolicyIPPoolStaticSubnetRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(staticSubnetPoolID, staticSubnetID).Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())
		d.SetId(staticSubnetID)

		err := resourceNsxtPolicyIPPoolStaticSubnetRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())

		err := resourceNsxtPolicyIPPoolStaticSubnetRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read success populates fields from the API object", func(t *testing.T) {
		mockSDK.EXPECT().Get(staticSubnetPoolID, staticSubnetID).Return(staticSubnetAPIStructValue(t), nil)

		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())
		d.SetId(staticSubnetID)

		err := resourceNsxtPolicyIPPoolStaticSubnetRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, staticSubnetCIDR, d.Get("cidr"))
		assert.Equal(t, "10.0.0.1", d.Get("gateway"))
	})

	t.Run("Read succeeds when the ID is a path and gets normalized", func(t *testing.T) {
		mockSDK.EXPECT().Get(staticSubnetPoolID, staticSubnetID).Return(staticSubnetAPIStructValue(t), nil)

		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())
		d.SetId("/infra/ip-pools/pool-2/ip-subnets/" + staticSubnetID)

		err := resourceNsxtPolicyIPPoolStaticSubnetRead(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyIPPoolStaticSubnetUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupStaticSubnetMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())

		err := resourceNsxtPolicyIPPoolStaticSubnetUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(staticSubnetPoolID, staticSubnetID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(staticSubnetPoolID, staticSubnetID).Return(staticSubnetAPIStructValue(t), nil),
		)

		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())
		d.SetId(staticSubnetID)

		err := resourceNsxtPolicyIPPoolStaticSubnetUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update propagates a Patch error", func(t *testing.T) {
		mockSDK.EXPECT().Patch(staticSubnetPoolID, staticSubnetID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())
		d.SetId(staticSubnetID)

		err := resourceNsxtPolicyIPPoolStaticSubnetUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPPoolStaticSubnetDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupStaticSubnetMock(t, ctrl)
	defer restore()

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())

		err := resourceNsxtPolicyIPPoolStaticSubnetDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(staticSubnetPoolID, staticSubnetID, nil).Return(nil)

		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())
		d.SetId(staticSubnetID)

		err := resourceNsxtPolicyIPPoolStaticSubnetDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete propagates an API error", func(t *testing.T) {
		mockSDK.EXPECT().Delete(staticSubnetPoolID, staticSubnetID, nil).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())
		d.SetId(staticSubnetID)

		err := resourceNsxtPolicyIPPoolStaticSubnetDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestUnitNsxt_resourceNsxtPolicyIPPoolStaticSubnetSchemaToStructValue(t *testing.T) {
	t.Run("converts schema data to a StructValue round-trippable back to the model", func(t *testing.T) {
		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())

		sv, err := resourceNsxtPolicyIPPoolStaticSubnetSchemaToStructValue(d, staticSubnetID)
		require.NoError(t, err)
		require.NotNil(t, sv)

		converter := bindings.NewTypeConverter()
		golangValue, errs := converter.ConvertToGolang(sv, model.IpAddressPoolStaticSubnetBindingType())
		require.Empty(t, errs)
		obj := golangValue.(model.IpAddressPoolStaticSubnet)
		assert.Equal(t, staticSubnetCIDR, *obj.Cidr)
		assert.Equal(t, staticSubnetID, *obj.Id)
	})
}
