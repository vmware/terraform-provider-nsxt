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
	vpcSvcEpID          = "svc-ep-id"
	vpcSvcEpDisplayName = "test-svc-ep"
	vpcSvcEpDescription = "Test VPC Service Endpoint"
	vpcSvcEpRevision    = int64(1)
	vpcSvcEpIP          = "192.168.1.100"
	vpcSvcEpIPType      = nsxModel.VpcServiceEndpoint_SERVICE_ENDPOINT_IP_TYPE_WORKLOAD
)

func vpcSvcEpAPIResponse() nsxModel.VpcServiceEndpoint {
	return nsxModel.VpcServiceEndpoint{
		Id:                    &vpcSvcEpID,
		DisplayName:           &vpcSvcEpDisplayName,
		Description:           &vpcSvcEpDescription,
		Revision:              &vpcSvcEpRevision,
		ServiceEndpointIp:     &vpcSvcEpIP,
		ServiceEndpointIpType: &vpcSvcEpIPType,
	}
}

func minimalVpcSvcEpData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":             vpcSvcEpDisplayName,
		"description":              vpcSvcEpDescription,
		"nsx_id":                   vpcSvcEpID,
		"service_endpoint_ip":      vpcSvcEpIP,
		"service_endpoint_ip_type": vpcSvcEpIPType,
	}
}

func setupVpcSvcEpMock(t *testing.T, ctrl *gomock.Controller) (*vpcsmocks.MockVpcServiceEndpointsClient, func()) {
	mockSDK := vpcsmocks.NewMockVpcServiceEndpointsClient(ctrl)
	mockWrapper := &apivpcs.VpcServiceEndpointClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  "project1",
		VPCID:      "vpc1",
	}

	original := cliVpcServiceEndpointsClient
	cliVpcServiceEndpointsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcServiceEndpointClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcServiceEndpointsClient = original }
}

func TestMockResourceNsxtVpcServiceEndpointCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcSvcEpMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcSvcEpID).Return(nsxModel.VpcServiceEndpoint{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), vpcSvcEpID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcSvcEpID).Return(vpcSvcEpAPIResponse(), nil),
		)

		res := resourceNsxtVpcServiceEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSvcEpData())

		err := resourceNsxtVpcServiceEndpointCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcSvcEpID, d.Id())
		assert.Equal(t, vpcSvcEpDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVpcServiceEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSvcEpData())

		err := resourceNsxtVpcServiceEndpointCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcServiceEndpointRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcSvcEpMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcSvcEpID).Return(vpcSvcEpAPIResponse(), nil)

		res := resourceNsxtVpcServiceEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSvcEpData())
		d.SetId(vpcSvcEpID)

		err := resourceNsxtVpcServiceEndpointRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcSvcEpDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcSvcEpIP, d.Get("service_endpoint_ip"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcSvcEpMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcSvcEpID).Return(nsxModel.VpcServiceEndpoint{}, vapiErrors.NotFound{})

		res := resourceNsxtVpcServiceEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSvcEpData())
		d.SetId(vpcSvcEpID)

		err := resourceNsxtVpcServiceEndpointRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcServiceEndpointUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcSvcEpMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), vpcSvcEpID, gomock.Any()).Return(vpcSvcEpAPIResponse(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcSvcEpID).Return(vpcSvcEpAPIResponse(), nil),
		)

		res := resourceNsxtVpcServiceEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSvcEpData())
		d.SetId(vpcSvcEpID)

		err := resourceNsxtVpcServiceEndpointUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtVpcServiceEndpointDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcSvcEpMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), vpcSvcEpID).Return(nil)

		res := resourceNsxtVpcServiceEndpoint()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcSvcEpData())
		d.SetId(vpcSvcEpID)

		err := resourceNsxtVpcServiceEndpointDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
