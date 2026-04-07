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

	sharesapi "github.com/vmware/terraform-provider-nsxt/api/infra/shares"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	sharesmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/shares"
)

var (
	sharedResID           = "sr-001"
	sharedResDisplayName  = "Test Shared Resource"
	sharedResDescription  = "Test shared resource description"
	sharedResRevision     = int64(1)
	sharedResShareID      = "share-1"
	sharedResSharePath    = "/infra/shares/share-1"
	sharedResPath         = "/infra/shares/share-1/resources/sr-001"
	sharedResResourcePath = "/infra/domains/default/groups/group-1"
)

func sharedResAPIResponse() nsxModel.SharedResource {
	return nsxModel.SharedResource{
		Id:              &sharedResID,
		DisplayName:     &sharedResDisplayName,
		Description:     &sharedResDescription,
		Revision:        &sharedResRevision,
		Path:            &sharedResPath,
		ResourceObjects: []nsxModel.ResourceObject{},
	}
}

func minimalSharedResData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": sharedResDisplayName,
		"description":  sharedResDescription,
		"share_path":   sharedResSharePath,
		"resource_object": []interface{}{
			map[string]interface{}{
				"include_children": false,
				"resource_path":    sharedResResourcePath,
			},
		},
	}
}

func setupSharedResMock(t *testing.T, ctrl *gomock.Controller) (*sharesmocks.MockResourcesClient, func()) {
	mockSDK := sharesmocks.NewMockResourcesClient(ctrl)
	mockWrapper := &sharesapi.SharedResourceClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliSharedResourcesClient
	cliSharedResourcesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *sharesapi.SharedResourceClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliSharedResourcesClient = original }
}

func TestMockResourceNsxtPolicySharedResourceCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupSharedResMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(sharedResShareID, gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(sharedResShareID, gomock.Any()).Return(sharedResAPIResponse(), nil),
		)

		res := resourceNsxtPolicySharedResource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSharedResData())

		err := resourceNsxtPolicySharedResourceCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicySharedResourceRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupSharedResMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(sharedResShareID, sharedResID).Return(sharedResAPIResponse(), nil)

		res := resourceNsxtPolicySharedResource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSharedResData())
		d.SetId(sharedResID)

		err := resourceNsxtPolicySharedResourceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, sharedResDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(sharedResShareID, sharedResID).Return(nsxModel.SharedResource{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicySharedResource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSharedResData())
		d.SetId(sharedResID)

		err := resourceNsxtPolicySharedResourceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySharedResource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSharedResData())

		err := resourceNsxtPolicySharedResourceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySharedResourceUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupSharedResMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(sharedResShareID, sharedResID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(sharedResShareID, sharedResID).Return(sharedResAPIResponse(), nil),
		)

		res := resourceNsxtPolicySharedResource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSharedResData())
		d.SetId(sharedResID)

		err := resourceNsxtPolicySharedResourceUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicySharedResourceDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupSharedResMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(sharedResShareID, sharedResID).Return(nil)

		res := resourceNsxtPolicySharedResource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSharedResData())
		d.SetId(sharedResID)

		err := resourceNsxtPolicySharedResourceDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySharedResource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSharedResData())

		err := resourceNsxtPolicySharedResourceDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
