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
	cpID          = "conn-profile-id"
	cpDisplayName = "test-connectivity-profile"
	cpDescription = "Test Connectivity Profile"
	cpRevision    = int64(1)
)

func cpAPIResponse() nsxModel.VpcConnectivityProfile {
	return nsxModel.VpcConnectivityProfile{
		Id:          &cpID,
		DisplayName: &cpDisplayName,
		Description: &cpDescription,
		Revision:    &cpRevision,
	}
}

func minimalCpData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": cpDisplayName,
		"description":  cpDescription,
		"nsx_id":       cpID,
	}
}

func setupCpMock(t *testing.T, ctrl *gomock.Controller) (*projectmocks.MockVpcConnectivityProfilesClient, func()) {
	mockSDK := projectmocks.NewMockVpcConnectivityProfilesClient(ctrl)
	mockWrapper := &apiprojects.VpcConnectivityProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  "project1",
	}

	original := cliVpcConnectivityProfilesClient
	cliVpcConnectivityProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiprojects.VpcConnectivityProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcConnectivityProfilesClient = original }
}

func TestMockResourceNsxtVpcConnectivityProfileCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupCpMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), cpID).Return(nsxModel.VpcConnectivityProfile{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), cpID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), cpID).Return(cpAPIResponse(), nil),
		)

		res := resourceNsxtVpcConnectivityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCpData())

		err := resourceNsxtVpcConnectivityProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, cpID, d.Id())
		assert.Equal(t, cpDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVpcConnectivityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCpData())

		err := resourceNsxtVpcConnectivityProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcConnectivityProfileRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupCpMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), cpID).Return(cpAPIResponse(), nil)

		res := resourceNsxtVpcConnectivityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCpData())
		d.SetId(cpID)

		err := resourceNsxtVpcConnectivityProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, cpDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupCpMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), cpID).Return(nsxModel.VpcConnectivityProfile{}, vapiErrors.NotFound{})

		res := resourceNsxtVpcConnectivityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCpData())
		d.SetId(cpID)

		err := resourceNsxtVpcConnectivityProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcConnectivityProfileUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupCpMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), cpID, gomock.Any()).Return(cpAPIResponse(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), cpID).Return(cpAPIResponse(), nil),
		)

		res := resourceNsxtVpcConnectivityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCpData())
		d.SetId(cpID)

		err := resourceNsxtVpcConnectivityProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtVpcConnectivityProfileDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupCpMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), cpID).Return(nil)

		res := resourceNsxtVpcConnectivityProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalCpData())
		d.SetId(cpID)

		err := resourceNsxtVpcConnectivityProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
