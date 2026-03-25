// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/orgs/projects/vpcs/GroupsClient.go -package=mocks -source=<sdk>/services/nsxt/orgs/projects/vpcs/GroupsClient.go GroupsClient

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

	apidomains "github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	vpcmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vpcGroupID          = "group-id"
	vpcGroupDisplayName = "test-group"
	vpcGroupDescription = "Test VPC Group"
	vpcGroupRevision    = int64(1)
)

func vpcGroupAPIResponse() nsxModel.Group {
	path := "/orgs/default/projects/project1/vpcs/vpc1/groups/group-id"
	return nsxModel.Group{
		Id:          &vpcGroupID,
		DisplayName: &vpcGroupDisplayName,
		Description: &vpcGroupDescription,
		Revision:    &vpcGroupRevision,
		Path:        &path,
	}
}

func minimalVpcGroupData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": vpcGroupDisplayName,
		"description":  vpcGroupDescription,
		"nsx_id":       vpcGroupID,
		"group_type":   "",
	}
}

func setupVpcGroupMock(t *testing.T, ctrl *gomock.Controller) (*vpcmocks.MockGroupsClient, func()) {
	mockSDK := vpcmocks.NewMockGroupsClient(ctrl)
	mockWrapper := &apidomains.GroupClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  "project1",
		VPCID:      "vpc1",
	}

	original := cliGroupsClient
	cliGroupsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.GroupClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliGroupsClient = original }
}

func TestMockResourceNsxtVpcGroupCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcGroupMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcGroupID).Return(nsxModel.Group{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), vpcGroupID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcGroupID).Return(vpcGroupAPIResponse(), nil),
		)

		res := resourceNsxtVPCGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGroupData())

		err := resourceNsxtVPCGroupCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcGroupID, d.Id())
		assert.Equal(t, vpcGroupDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails on NSX version check", func(t *testing.T) {
		util.NsxVersion = "8.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVPCGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGroupData())

		err := resourceNsxtVPCGroupCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.0.0")
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcGroupMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcGroupID).Return(vpcGroupAPIResponse(), nil)

		res := resourceNsxtVPCGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGroupData())

		err := resourceNsxtVPCGroupCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtVpcGroupRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcGroupMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcGroupID).Return(vpcGroupAPIResponse(), nil)

		res := resourceNsxtVPCGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGroupData())
		d.SetId(vpcGroupID)

		err := resourceNsxtVPCGroupRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcGroupDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcGroupDescription, d.Get("description"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcGroupMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcGroupID).Return(nsxModel.Group{}, notFoundErr)

		res := resourceNsxtVPCGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGroupData())
		d.SetId(vpcGroupID)

		err := resourceNsxtVPCGroupRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtVPCGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGroupData())

		err := resourceNsxtVPCGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcGroupUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcGroupMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), vpcGroupID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcGroupID).Return(vpcGroupAPIResponse(), nil),
		)

		res := resourceNsxtVPCGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGroupData())
		d.SetId(vpcGroupID)

		err := resourceNsxtVPCGroupUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtVpcGroupDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcGroupMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), vpcGroupID).Return(nil)

		res := resourceNsxtVPCGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGroupData())
		d.SetId(vpcGroupID)

		err := resourceNsxtVPCGroupDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcGroupMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), vpcGroupID).Return(errors.New("delete failed"))

		res := resourceNsxtVPCGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcGroupData())
		d.SetId(vpcGroupID)

		err := resourceNsxtVPCGroupDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
	})
}
