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

	apivpcs "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	vpcsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vpcEpID              = "ep-id"
	vpcEpDisplayName     = "test-ep"
	vpcEpDescription     = "Test VPC Endpoint"
	vpcEpRevision        = int64(1)
	vpcEpSvcEndpointPath = "/orgs/default/projects/proj-provider/vpcs/vpc-provider/vpc-service-endpoints/svc-ep-1"
	vpcEpIPAllocPath     = "/orgs/default/projects/project-1/vpcs/vpc-1/ip-address-allocations/alloc-1"
)

func vpcEpAPIResponse() nsxModel.VpcEndpoint {
	return nsxModel.VpcEndpoint{
		Id:                 &vpcEpID,
		DisplayName:        &vpcEpDisplayName,
		Description:        &vpcEpDescription,
		Revision:           &vpcEpRevision,
		VpcServiceEndpoint: &vpcEpSvcEndpointPath,
		IpAllocationPath:   &vpcEpIPAllocPath,
	}
}

func minimalVpcEpData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":         vpcEpDisplayName,
		"description":          vpcEpDescription,
		"nsx_id":               vpcEpID,
		"vpc_service_endpoint": vpcEpSvcEndpointPath,
		"ip_allocation_path":   vpcEpIPAllocPath,
	}
}

func setupVpcEpMock(t *testing.T, ctrl *gomock.Controller) (*vpcsmocks.MockVpcEndpointsClient, func()) {
	mockSDK := vpcsmocks.NewMockVpcEndpointsClient(ctrl)
	mockWrapper := &apivpcs.VpcEndpointClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  "project1",
		VPCID:      "vpc1",
	}

	original := cliVpcEndpointsClient
	cliVpcEndpointsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcEndpointClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcEndpointsClient = original }
}

func TestMockResourceNsxtVpcEndpointCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcEpMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcEpID).Return(nsxModel.VpcEndpoint{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), vpcEpID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcEpID).Return(vpcEpAPIResponse(), nil),
		)

		res := resourceNsxtVpcEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcEpData())

		err := resourceNsxtVpcEndpointCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcEpID, d.Id())
		assert.Equal(t, vpcEpDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVpcEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcEpData())

		err := resourceNsxtVpcEndpointCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcEndpointRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcEpMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcEpID).Return(vpcEpAPIResponse(), nil)

		res := resourceNsxtVpcEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcEpData())
		d.SetId(vpcEpID)

		err := resourceNsxtVpcEndpointRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcEpDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcEpSvcEndpointPath, d.Get("vpc_service_endpoint"))
		assert.Equal(t, vpcEpIPAllocPath, d.Get("ip_allocation_path"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcEpMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcEpID).Return(nsxModel.VpcEndpoint{}, vapiErrors.NotFound{})

		res := resourceNsxtVpcEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcEpData())
		d.SetId(vpcEpID)

		err := resourceNsxtVpcEndpointRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcEndpointUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcEpMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), vpcEpID, gomock.Any()).Return(vpcEpAPIResponse(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcEpID).Return(vpcEpAPIResponse(), nil),
		)

		res := resourceNsxtVpcEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcEpData())
		d.SetId(vpcEpID)

		err := resourceNsxtVpcEndpointUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtVpcEndpointDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcEpMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), vpcEpID).Return(nil)

		res := resourceNsxtVpcEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcEpData())
		d.SetId(vpcEpID)

		err := resourceNsxtVpcEndpointDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
