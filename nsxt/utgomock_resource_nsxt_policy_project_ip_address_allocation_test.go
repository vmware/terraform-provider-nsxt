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

	projectsapi "github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	projectsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects"
)

var (
	ipAllocID          = "ip-alloc-id"
	ipAllocDisplayName = "test-ip-alloc"
	ipAllocDescription = "Test IP Address Allocation"
	ipAllocRevision    = int64(1)
	ipAllocProjectID   = "project1"
	ipAllocOrgID       = "default"
)

func ipAllocAPIResponse() nsxModel.ProjectIpAddressAllocation {
	return nsxModel.ProjectIpAddressAllocation{
		Id:          &ipAllocID,
		DisplayName: &ipAllocDisplayName,
		Description: &ipAllocDescription,
		Revision:    &ipAllocRevision,
	}
}

func minimalIpAllocData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": ipAllocDisplayName,
		"description":  ipAllocDescription,
		"nsx_id":       ipAllocID,
		"context": []interface{}{
			map[string]interface{}{
				"project_id":  ipAllocProjectID,
				"vpc_id":      "",
				"from_global": false,
			},
		},
	}
}

func setupIpAllocMock(t *testing.T, ctrl *gomock.Controller) (*projectsmocks.MockIpAddressAllocationsClient, func()) {
	mockSDK := projectsmocks.NewMockIpAddressAllocationsClient(ctrl)
	mockWrapper := &projectsapi.ProjectIpAddressAllocationClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  ipAllocProjectID,
	}

	original := cliProjectIpAddressAllocationsClient
	cliProjectIpAddressAllocationsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *projectsapi.ProjectIpAddressAllocationClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliProjectIpAddressAllocationsClient = original }
}

func TestMockResourceNsxtPolicyProjectIpAddressAllocationCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIpAllocMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(ipAllocOrgID, ipAllocProjectID, ipAllocID).Return(nsxModel.ProjectIpAddressAllocation{}, notFoundErr),
			mockSDK.EXPECT().Patch(ipAllocOrgID, ipAllocProjectID, ipAllocID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(ipAllocOrgID, ipAllocProjectID, ipAllocID).Return(ipAllocAPIResponse(), nil),
		)

		res := resourceNsxtPolicyProjectIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpAllocData())

		err := resourceNsxtPolicyProjectIpAddressAllocationCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipAllocID, d.Id())
		assert.Equal(t, ipAllocDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipAllocOrgID, ipAllocProjectID, ipAllocID).Return(ipAllocAPIResponse(), nil)

		res := resourceNsxtPolicyProjectIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpAllocData())

		err := resourceNsxtPolicyProjectIpAddressAllocationCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyProjectIpAddressAllocationRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIpAllocMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipAllocOrgID, ipAllocProjectID, ipAllocID).Return(ipAllocAPIResponse(), nil)

		res := resourceNsxtPolicyProjectIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpAllocData())
		d.SetId(ipAllocID)

		err := resourceNsxtPolicyProjectIpAddressAllocationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipAllocDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipAllocOrgID, ipAllocProjectID, ipAllocID).Return(nsxModel.ProjectIpAddressAllocation{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyProjectIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpAllocData())
		d.SetId(ipAllocID)

		err := resourceNsxtPolicyProjectIpAddressAllocationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyProjectIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpAllocData())

		err := resourceNsxtPolicyProjectIpAddressAllocationRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyProjectIpAddressAllocationUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIpAllocMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(ipAllocOrgID, ipAllocProjectID, ipAllocID, gomock.Any()).Return(ipAllocAPIResponse(), nil),
			mockSDK.EXPECT().Get(ipAllocOrgID, ipAllocProjectID, ipAllocID).Return(ipAllocAPIResponse(), nil),
		)

		res := resourceNsxtPolicyProjectIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpAllocData())
		d.SetId(ipAllocID)

		err := resourceNsxtPolicyProjectIpAddressAllocationUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyProjectIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpAllocData())

		err := resourceNsxtPolicyProjectIpAddressAllocationUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyProjectIpAddressAllocationDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIpAllocMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(ipAllocOrgID, ipAllocProjectID, ipAllocID).Return(nil)

		res := resourceNsxtPolicyProjectIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpAllocData())
		d.SetId(ipAllocID)

		err := resourceNsxtPolicyProjectIpAddressAllocationDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyProjectIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpAllocData())

		err := resourceNsxtPolicyProjectIpAddressAllocationDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
