//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/orgs/projects/VpcsClient.go -package=mocks -source=<sdk>/services/nsxt/orgs/projects/VpcsClient.go VpcsClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	apiprojects "github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	projectmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vpcID          = "vpc-test-id"
	vpcDisplayName = "test-vpc"
	vpcDescription = "Test VPC"
	vpcRevision    = int64(1)
	vpcProjectID   = "project1"
)

func vpcAPIResponse() nsxModel.Vpc {
	ipType := nsxModel.Vpc_IP_ADDRESS_TYPE_IPV4
	return nsxModel.Vpc{
		Id:            &vpcID,
		DisplayName:   &vpcDisplayName,
		Description:   &vpcDescription,
		Revision:      &vpcRevision,
		IpAddressType: &ipType,
	}
}

func minimalVpcData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": vpcDisplayName,
		"description":  vpcDescription,
		"nsx_id":       vpcID,
	}
}

func setupVpcMock(t *testing.T, ctrl *gomock.Controller) (*projectmocks.MockVpcsClient, func()) {
	mockSDK := projectmocks.NewMockVpcsClient(ctrl)
	mockWrapper := &apiprojects.VpcClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  vpcProjectID,
		VPCID:      "parent-vpc",
	}

	original := cliVpcsClient
	cliVpcsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiprojects.VpcClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcsClient = original }
}

func TestMockResourceNsxtVpcCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vpcID).Return(nsxModel.Vpc{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), vpcID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vpcID).Return(vpcAPIResponse(), nil),
		)

		res := resourceNsxtVpc()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcData())

		err := resourceNsxtVpcCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcID, d.Id())
		assert.Equal(t, vpcDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vpcID).Return(vpcAPIResponse(), nil)

		res := resourceNsxtVpc()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcData())

		err := resourceNsxtVpcCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create API error", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		apiErr := errors.New("internal server error")
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vpcID).Return(nsxModel.Vpc{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), vpcID, gomock.Any()).Return(apiErr),
		)

		res := resourceNsxtVpc()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcData())

		err := resourceNsxtVpcCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vpcID).Return(vpcAPIResponse(), nil)

		res := resourceNsxtVpc()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcData())
		d.SetId(vpcID)

		err := resourceNsxtVpcRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcDescription, d.Get("description"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vpcID).Return(nsxModel.Vpc{}, notFoundErr)

		res := resourceNsxtVpc()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcData())
		d.SetId(vpcID)

		err := resourceNsxtVpcRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcMock(t, ctrl)
		defer restore()

		updatedDisplayName := "updated-vpc"
		updatedObj := vpcAPIResponse()
		updatedObj.DisplayName = &updatedDisplayName
		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), vpcID, gomock.Any()).Return(updatedObj, nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vpcID).Return(updatedObj, nil),
		)

		res := resourceNsxtVpc()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcData())
		d.SetId(vpcID)

		err := resourceNsxtVpcUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), vpcID, gomock.Any()).Return(nsxModel.Vpc{}, errors.New("update failed"))

		res := resourceNsxtVpc()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcData())
		d.SetId(vpcID)

		err := resourceNsxtVpcUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), vpcID, gomock.Any()).Return(nil)

		res := resourceNsxtVpc()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcData())
		d.SetId(vpcID)

		err := resourceNsxtVpcDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), vpcID, gomock.Any()).Return(errors.New("delete failed"))

		res := resourceNsxtVpc()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcData())
		d.SetId(vpcID)

		err := resourceNsxtVpcDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
