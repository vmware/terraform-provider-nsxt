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

	apivpcs "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	vpcsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vpcIpAllocID          = "ip-alloc-id"
	vpcIpAllocDisplayName = "test-ip-alloc"
	vpcIpAllocDescription = "Test VPC IP Address Allocation"
	vpcIpAllocRevision    = int64(1)
)

func vpcIpAllocAPIResponse() nsxModel.VpcIpAddressAllocation {
	ipType := nsxModel.VpcIpAddressAllocation_IP_ADDRESS_TYPE_IPV4
	return nsxModel.VpcIpAddressAllocation{
		Id:            &vpcIpAllocID,
		DisplayName:   &vpcIpAllocDisplayName,
		Description:   &vpcIpAllocDescription,
		Revision:      &vpcIpAllocRevision,
		IpAddressType: &ipType,
	}
}

func minimalVpcIpAllocData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": vpcIpAllocDisplayName,
		"description":  vpcIpAllocDescription,
		"nsx_id":       vpcIpAllocID,
	}
}

func setupVpcIpAllocMock(t *testing.T, ctrl *gomock.Controller) (*vpcsmocks.MockIpAddressAllocationsClient, func()) {
	mockSDK := vpcsmocks.NewMockIpAddressAllocationsClient(ctrl)
	mockWrapper := &apivpcs.VpcIpAddressAllocationClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  "project1",
		VPCID:      "vpc1",
	}

	original := cliVpcIpAddressAllocationsClient
	cliVpcIpAddressAllocationsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcIpAddressAllocationClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcIpAddressAllocationsClient = original }
}

func TestMockResourceNsxtVpcIpAddressAllocationCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcIpAllocMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcIpAllocID).Return(nsxModel.VpcIpAddressAllocation{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), vpcIpAllocID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcIpAllocID).Return(vpcIpAllocAPIResponse(), nil),
		)

		res := resourceNsxtVpcIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcIpAllocData())

		err := resourceNsxtVpcIpAddressAllocationCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcIpAllocID, d.Id())
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVpcIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcIpAllocData())

		err := resourceNsxtVpcIpAddressAllocationCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcIpAddressAllocationRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcIpAllocMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcIpAllocID).Return(vpcIpAllocAPIResponse(), nil)

		res := resourceNsxtVpcIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcIpAllocData())
		d.SetId(vpcIpAllocID)

		err := resourceNsxtVpcIpAddressAllocationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcIpAllocDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcIpAllocMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcIpAllocID).Return(nsxModel.VpcIpAddressAllocation{}, vapiErrors.NotFound{})

		res := resourceNsxtVpcIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcIpAllocData())
		d.SetId(vpcIpAllocID)

		err := resourceNsxtVpcIpAddressAllocationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcIpAddressAllocationUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcIpAllocMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), vpcIpAllocID, gomock.Any()).Return(vpcIpAllocAPIResponse(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcIpAllocID).Return(vpcIpAllocAPIResponse(), nil),
		)

		res := resourceNsxtVpcIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcIpAllocData())
		d.SetId(vpcIpAllocID)

		err := resourceNsxtVpcIpAddressAllocationUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtVpcIpAddressAllocationDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcIpAllocMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), vpcIpAllocID).Return(nil)

		res := resourceNsxtVpcIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcIpAllocData())
		d.SetId(vpcIpAllocID)

		err := resourceNsxtVpcIpAddressAllocationDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
