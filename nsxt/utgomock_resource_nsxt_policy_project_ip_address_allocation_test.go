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

	projectsapi "github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	projectsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
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

	t.Run("Create success (defaults omitted on NSX 9.1)", func(t *testing.T) {
		util.NsxVersion = "9.1.0"
		defer func() { util.NsxVersion = "" }()
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(ipAllocOrgID, ipAllocProjectID, ipAllocID).Return(nsxModel.ProjectIpAddressAllocation{}, notFoundErr),
			mockSDK.EXPECT().Patch(ipAllocOrgID, ipAllocProjectID, ipAllocID, gomock.Any()).DoAndReturn(
				func(_ string, _ string, _ string, req nsxModel.ProjectIpAddressAllocation) error {
					assert.Nil(t, req.IpAddressType)
					assert.Nil(t, req.Ipv6AllocationPrefixLength)
					return nil
				}),
			mockSDK.EXPECT().Get(ipAllocOrgID, ipAllocProjectID, ipAllocID).Return(ipAllocAPIResponse(), nil),
		)

		res := resourceNsxtPolicyProjectIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpAllocData())

		err := resourceNsxtPolicyProjectIpAddressAllocationCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipAllocID, d.Id())
		assert.Equal(t, ipAllocDisplayName, d.Get("display_name"))
	})

	t.Run("Create success IPv6 on NSX 9.2", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()
		notFoundErr := vapiErrors.NotFound{}
		ipType := nsxModel.ProjectIpAddressAllocation_IP_ADDRESS_TYPE_IPV6
		prefixLen := int64(64)
		ipBlock := "/infra/ip-blocks/block-v6"
		v6Response := nsxModel.ProjectIpAddressAllocation{
			Id:                         &ipAllocID,
			DisplayName:                &ipAllocDisplayName,
			Description:                &ipAllocDescription,
			Revision:                   &ipAllocRevision,
			IpAddressType:              &ipType,
			Ipv6AllocationPrefixLength: &prefixLen,
			IpBlock:                    &ipBlock,
		}
		gomock.InOrder(
			mockSDK.EXPECT().Get(ipAllocOrgID, ipAllocProjectID, ipAllocID).Return(nsxModel.ProjectIpAddressAllocation{}, notFoundErr),
			mockSDK.EXPECT().Patch(ipAllocOrgID, ipAllocProjectID, ipAllocID, gomock.Any()).DoAndReturn(
				func(_ string, _ string, _ string, req nsxModel.ProjectIpAddressAllocation) error {
					assert.Equal(t, &ipType, req.IpAddressType)
					assert.Equal(t, &prefixLen, req.Ipv6AllocationPrefixLength)
					assert.Nil(t, req.AllocationSize)
					assert.Equal(t, &ipBlock, req.IpBlock)
					return nil
				}),
			mockSDK.EXPECT().Get(ipAllocOrgID, ipAllocProjectID, ipAllocID).Return(v6Response, nil),
		)

		res := resourceNsxtPolicyProjectIpAddressAllocation()
		data := minimalIpAllocData()
		data["ip_address_type"] = "IPV6"
		data["ipv6_allocation_prefix_length"] = 64
		data["ip_block"] = "/infra/ip-blocks/block-v6"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyProjectIpAddressAllocationCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipAllocID, d.Id())
		assert.Equal(t, "IPV6", d.Get("ip_address_type"))
		assert.Equal(t, 64, d.Get("ipv6_allocation_prefix_length"))
		assert.Equal(t, 0, d.Get("allocation_size"))
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

	t.Run("Update success IPv6 on NSX 9.2", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()
		ipType := nsxModel.ProjectIpAddressAllocation_IP_ADDRESS_TYPE_IPV6
		prefixLen := int64(64)
		v6Response := nsxModel.ProjectIpAddressAllocation{
			Id:                         &ipAllocID,
			DisplayName:                &ipAllocDisplayName,
			Description:                &ipAllocDescription,
			Revision:                   &ipAllocRevision,
			IpAddressType:              &ipType,
			Ipv6AllocationPrefixLength: &prefixLen,
		}
		gomock.InOrder(
			mockSDK.EXPECT().Update(ipAllocOrgID, ipAllocProjectID, ipAllocID, gomock.Any()).DoAndReturn(
				func(_ string, _ string, _ string, req nsxModel.ProjectIpAddressAllocation) (nsxModel.ProjectIpAddressAllocation, error) {
					assert.Equal(t, &ipType, req.IpAddressType)
					assert.Equal(t, &prefixLen, req.Ipv6AllocationPrefixLength)
					return v6Response, nil
				}),
			mockSDK.EXPECT().Get(ipAllocOrgID, ipAllocProjectID, ipAllocID).Return(v6Response, nil),
		)

		res := resourceNsxtPolicyProjectIpAddressAllocation()
		data := minimalIpAllocData()
		data["ip_address_type"] = "IPV6"
		data["ipv6_allocation_prefix_length"] = 64
		d := schema.TestResourceDataRaw(t, res.Schema, data)
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
