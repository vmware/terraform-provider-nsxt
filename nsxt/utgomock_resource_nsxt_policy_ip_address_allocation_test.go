//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	ippoolsapi "github.com/vmware/terraform-provider-nsxt/api/infra/ip_pools"
	realizedstateapi "github.com/vmware/terraform-provider-nsxt/api/infra/realized_state"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	ippoolmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/ip_pools"
	realizedstatemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/realized_state"
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

func minimalNsxtIpAllocDataNoIP() map[string]interface{} {
	data := minimalNsxtIpAllocData()
	delete(data, "allocation_ip")
	data["timeout"] = 5
	return data
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

func setupNsxtIpAllocRealizedMock(t *testing.T, ctrl *gomock.Controller) (*realizedstatemocks.MockRealizedEntitiesClient, func()) {
	mockSDK := realizedstatemocks.NewMockRealizedEntitiesClient(ctrl)
	mockWrapper := &realizedstateapi.RealizedEntityClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliRealizedEntitiesClient
	cliRealizedEntitiesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *realizedstateapi.RealizedEntityClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliRealizedEntitiesClient = original }
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

func TestMockResourceNsxtPolicyIPAddressAllocationCreateRealizationCleanup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restoreAlloc := setupNsxtIpAllocMock(t, ctrl)
	defer restoreAlloc()
	mockRealizedSDK, restoreRealized := setupNsxtIpAllocRealizedMock(t, ctrl)
	defer restoreRealized()

	t.Run("Create cleans up and clears ID when realization fails", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(nsxtIpAllocPoolID, nsxtIpAllocID).Return(nsxModel.IpAddressAllocation{}, notFoundErr),
			mockSDK.EXPECT().Patch(nsxtIpAllocPoolID, nsxtIpAllocID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(nsxtIpAllocPoolID, nsxtIpAllocID).Return(nsxtIpAllocAPIResponse(), nil),
		)
		mockRealizedSDK.EXPECT().List(nsxtIpAllocPath, nil).Return(nsxModel.GenericPolicyRealizedResourceListResult{}, errors.New("realization query failed"))
		// Cleanup must delete the already-created NSX object, keyed by the
		// resource ID directly (not read back from d.Id(), which isn't set yet).
		mockSDK.EXPECT().Delete(nsxtIpAllocPoolID, nsxtIpAllocID).Return(nil)

		res := resourceNsxtPolicyIPAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNsxtIpAllocDataNoIP())

		err := resourceNsxtPolicyIPAddressAllocationCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Create cleans up and clears ID when realized IP attribute is missing", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		realizedState := "REALIZED"
		otherAttrKey := "other_attr"
		gomock.InOrder(
			mockSDK.EXPECT().Get(nsxtIpAllocPoolID, nsxtIpAllocID).Return(nsxModel.IpAddressAllocation{}, notFoundErr),
			mockSDK.EXPECT().Patch(nsxtIpAllocPoolID, nsxtIpAllocID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(nsxtIpAllocPoolID, nsxtIpAllocID).Return(nsxtIpAllocAPIResponse(), nil),
		)
		mockRealizedSDK.EXPECT().List(nsxtIpAllocPath, nil).Return(nsxModel.GenericPolicyRealizedResourceListResult{
			Results: []nsxModel.GenericPolicyRealizedResource{
				{
					State: &realizedState,
					ExtendedAttributes: []nsxModel.AttributeVal{
						{Key: &otherAttrKey, Values: []string{"x"}},
					},
				},
			},
		}, nil)
		mockSDK.EXPECT().Delete(nsxtIpAllocPoolID, nsxtIpAllocID).Return(nil)

		res := resourceNsxtPolicyIPAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalNsxtIpAllocDataNoIP())

		err := resourceNsxtPolicyIPAddressAllocationCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Equal(t, "", d.Id())
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
