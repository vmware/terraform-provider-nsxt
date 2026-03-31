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

	ippoolsapi "github.com/vmware/terraform-provider-nsxt/api/infra/ip_pools"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	ippoolmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/ip_pools"
)

var (
	nsxtIpAllocID          = "alloc-1"
	nsxtIpAllocDisplayName = "Test Allocation"
	nsxtIpAllocDescription = "test ip allocation"
	nsxtIpAllocRevision    = int64(1)
	nsxtIpAllocPoolPath    = "/infra/ip-pools/pool-1"
	nsxtIpAllocPoolID      = "pool-1"
	nsxtIpAllocIP          = "192.168.1.10"
	nsxtIpAllocPath        = "/infra/ip-pools/pool-1/ip-allocations/alloc-1"
)

func nsxtIpAllocAPIResponse() nsxModel.IpAddressAllocation {
	return nsxModel.IpAddressAllocation{
		Id:           &nsxtIpAllocID,
		DisplayName:  &nsxtIpAllocDisplayName,
		Description:  &nsxtIpAllocDescription,
		Revision:     &nsxtIpAllocRevision,
		Path:         &nsxtIpAllocPath,
		ParentPath:   &nsxtIpAllocPoolPath,
		AllocationIp: &nsxtIpAllocIP,
	}
}

func minimalNsxtIpAllocData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":  nsxtIpAllocDisplayName,
		"description":   nsxtIpAllocDescription,
		"nsx_id":        nsxtIpAllocID,
		"pool_path":     nsxtIpAllocPoolPath,
		"allocation_ip": nsxtIpAllocIP,
	}
}

func setupNsxtIpAllocMock(t *testing.T, ctrl *gomock.Controller) (*ippoolmocks.MockIpAllocationsClient, func()) {
	mockSDK := ippoolmocks.NewMockIpAllocationsClient(ctrl)
	mockWrapper := &ippoolsapi.IpAddressAllocationClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliIpAllocationsClient
	cliIpAllocationsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *ippoolsapi.IpAddressAllocationClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliIpAllocationsClient = original }
}

func TestMockResourceNsxtPolicyIPAddressAllocationCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupNsxtIpAllocMock(t, ctrl)
	defer restore()

	t.Run("Create success with allocation_ip set", func(t *testing.T) {
		// When allocation_ip is set, create does not wait for realization
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(nsxtIpAllocPoolID, nsxtIpAllocID).Return(nsxModel.IpAddressAllocation{}, notFoundErr),
			mockSDK.EXPECT().Patch(nsxtIpAllocPoolID, nsxtIpAllocID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(nsxtIpAllocPoolID, nsxtIpAllocID).Return(nsxtIpAllocAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIPAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNsxtIpAllocData())

		err := resourceNsxtPolicyIPAddressAllocationCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, nsxtIpAllocID, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(nsxtIpAllocPoolID, nsxtIpAllocID).Return(nsxtIpAllocAPIResponse(), nil)

		res := resourceNsxtPolicyIPAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNsxtIpAllocData())

		err := resourceNsxtPolicyIPAddressAllocationCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPAddressAllocationRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupNsxtIpAllocMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(nsxtIpAllocPoolID, nsxtIpAllocID).Return(nsxtIpAllocAPIResponse(), nil)

		res := resourceNsxtPolicyIPAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNsxtIpAllocData())
		d.SetId(nsxtIpAllocID)

		err := resourceNsxtPolicyIPAddressAllocationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, nsxtIpAllocDisplayName, d.Get("display_name"))
		assert.Equal(t, nsxtIpAllocIP, d.Get("allocation_ip"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(nsxtIpAllocPoolID, nsxtIpAllocID).Return(nsxModel.IpAddressAllocation{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIPAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNsxtIpAllocData())
		d.SetId(nsxtIpAllocID)

		err := resourceNsxtPolicyIPAddressAllocationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNsxtIpAllocData())

		err := resourceNsxtPolicyIPAddressAllocationRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPAddressAllocationUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupNsxtIpAllocMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(nsxtIpAllocPoolID, nsxtIpAllocID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(nsxtIpAllocPoolID, nsxtIpAllocID).Return(nsxtIpAllocAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIPAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNsxtIpAllocData())
		d.SetId(nsxtIpAllocID)

		err := resourceNsxtPolicyIPAddressAllocationUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNsxtIpAllocData())

		err := resourceNsxtPolicyIPAddressAllocationUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIPAddressAllocationDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupNsxtIpAllocMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(nsxtIpAllocPoolID, nsxtIpAllocID).Return(nil)

		res := resourceNsxtPolicyIPAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNsxtIpAllocData())
		d.SetId(nsxtIpAllocID)

		err := resourceNsxtPolicyIPAddressAllocationDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNsxtIpAllocData())

		err := resourceNsxtPolicyIPAddressAllocationDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
