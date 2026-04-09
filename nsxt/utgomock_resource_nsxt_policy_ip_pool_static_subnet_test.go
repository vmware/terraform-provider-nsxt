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
}

func TestMockResourceNsxtPolicyIPPoolStaticSubnetUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupStaticSubnetMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())

		err := resourceNsxtPolicyIPPoolStaticSubnetUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPPoolStaticSubnetDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupStaticSubnetMock(t, ctrl)
	defer restore()

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPPoolStaticSubnet()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalStaticSubnetData())

		err := resourceNsxtPolicyIPPoolStaticSubnetDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
