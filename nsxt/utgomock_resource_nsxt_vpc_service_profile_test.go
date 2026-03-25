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

	apiprojects "github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	projectmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vspID          = "service-profile-id"
	vspDisplayName = "test-service-profile"
	vspDescription = "Test VPC Service Profile"
	vspRevision    = int64(1)
)

func vspAPIResponse() nsxModel.VpcServiceProfile {
	return nsxModel.VpcServiceProfile{
		Id:          &vspID,
		DisplayName: &vspDisplayName,
		Description: &vspDescription,
		Revision:    &vspRevision,
	}
}

func minimalVspData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": vspDisplayName,
		"description":  vspDescription,
		"nsx_id":       vspID,
	}
}

func setupVspMock(t *testing.T, ctrl *gomock.Controller) (*projectmocks.MockVpcServiceProfilesClient, func()) {
	mockSDK := projectmocks.NewMockVpcServiceProfilesClient(ctrl)
	mockWrapper := &apiprojects.VpcServiceProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  "project1",
	}

	original := cliVpcServiceProfilesClient
	cliVpcServiceProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiprojects.VpcServiceProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcServiceProfilesClient = original }
}

func TestMockResourceNsxtVpcServiceProfileCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVspMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vspID).Return(nsxModel.VpcServiceProfile{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), vspID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vspID).Return(vspAPIResponse(), nil),
		)

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData())

		err := resourceNsxtVpcServiceProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vspID, d.Id())
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData())

		err := resourceNsxtVpcServiceProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcServiceProfileRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVspMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vspID).Return(vspAPIResponse(), nil)

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData())
		d.SetId(vspID)

		err := resourceNsxtVpcServiceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vspDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVspMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vspID).Return(nsxModel.VpcServiceProfile{}, vapiErrors.NotFound{})

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData())
		d.SetId(vspID)

		err := resourceNsxtVpcServiceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcServiceProfileUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVspMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), vspID, gomock.Any()).Return(vspAPIResponse(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vspID).Return(vspAPIResponse(), nil),
		)

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData())
		d.SetId(vspID)

		err := resourceNsxtVpcServiceProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtVpcServiceProfileDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVspMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), vspID).Return(nil)

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData())
		d.SetId(vspID)

		err := resourceNsxtVpcServiceProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
