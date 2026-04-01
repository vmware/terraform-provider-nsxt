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

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	infraMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

var (
	lbClientSslID          = "lb-client-ssl-1"
	lbClientSslDisplayName = "Test LB Client SSL Profile"
	lbClientSslDescription = "Test lb client ssl profile"
	lbClientSslRevision    = int64(1)
)

func lbClientSslAPIResponse() nsxModel.LBClientSslProfile {
	return nsxModel.LBClientSslProfile{
		Id:          &lbClientSslID,
		DisplayName: &lbClientSslDisplayName,
		Description: &lbClientSslDescription,
		Revision:    &lbClientSslRevision,
	}
}

func minimalLBClientSslData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbClientSslDisplayName,
		"description":  lbClientSslDescription,
		"nsx_id":       lbClientSslID,
	}
}

func setupLBClientSslMock(t *testing.T, ctrl *gomock.Controller) (*infraMocks.MockLbClientSslProfilesClient, func()) {
	mockSDK := infraMocks.NewMockLbClientSslProfilesClient(ctrl)
	mockWrapper := &infraapi.LBClientSslProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliLbClientSslProfilesClient
	cliLbClientSslProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.LBClientSslProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliLbClientSslProfilesClient = original }
}

func TestMockResourceNsxtPolicyLBClientSslProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBClientSslMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(lbClientSslID).Return(nsxModel.LBClientSslProfile{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(lbClientSslID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(lbClientSslID).Return(lbClientSslAPIResponse(), nil),
		)

		res := resourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBClientSslData())

		err := resourceNsxtPolicyLBClientSslProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbClientSslID, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbClientSslID).Return(lbClientSslAPIResponse(), nil)

		res := resourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBClientSslData())

		err := resourceNsxtPolicyLBClientSslProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBClientSslProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBClientSslMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbClientSslID).Return(lbClientSslAPIResponse(), nil)

		res := resourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBClientSslData())
		d.SetId(lbClientSslID)

		err := resourceNsxtPolicyLBClientSslProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbClientSslDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbClientSslID).Return(nsxModel.LBClientSslProfile{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBClientSslData())
		d.SetId(lbClientSslID)

		err := resourceNsxtPolicyLBClientSslProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBClientSslData())

		err := resourceNsxtPolicyLBClientSslProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBClientSslProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBClientSslMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(lbClientSslID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(lbClientSslID).Return(lbClientSslAPIResponse(), nil),
		)

		res := resourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBClientSslData())
		d.SetId(lbClientSslID)

		err := resourceNsxtPolicyLBClientSslProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBClientSslData())

		err := resourceNsxtPolicyLBClientSslProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBClientSslProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBClientSslMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbClientSslID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBClientSslData())
		d.SetId(lbClientSslID)

		err := resourceNsxtPolicyLBClientSslProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBClientSslData())

		err := resourceNsxtPolicyLBClientSslProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
