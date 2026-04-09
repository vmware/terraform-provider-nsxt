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
	lbServerSslID          = "lb-server-ssl-1"
	lbServerSslDisplayName = "Test LB Server SSL Profile"
	lbServerSslDescription = "Test lb server ssl profile"
	lbServerSslRevision    = int64(1)
)

func lbServerSslAPIResponse() nsxModel.LBServerSslProfile {
	return nsxModel.LBServerSslProfile{
		Id:          &lbServerSslID,
		DisplayName: &lbServerSslDisplayName,
		Description: &lbServerSslDescription,
		Revision:    &lbServerSslRevision,
	}
}

func minimalLBServerSslData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbServerSslDisplayName,
		"description":  lbServerSslDescription,
		"nsx_id":       lbServerSslID,
	}
}

func setupLBServerSslMock(t *testing.T, ctrl *gomock.Controller) (*infraMocks.MockLbServerSslProfilesClient, func()) {
	mockSDK := infraMocks.NewMockLbServerSslProfilesClient(ctrl)
	mockWrapper := &infraapi.LBServerSslProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliLbServerSslProfilesClient
	cliLbServerSslProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.LBServerSslProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliLbServerSslProfilesClient = original }
}

func TestMockResourceNsxtPolicyLBServerSslProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBServerSslMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(lbServerSslID).Return(nsxModel.LBServerSslProfile{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(lbServerSslID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(lbServerSslID).Return(lbServerSslAPIResponse(), nil),
		)

		res := resourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBServerSslData())

		err := resourceNsxtPolicyLBServerSslProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbServerSslID, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbServerSslID).Return(lbServerSslAPIResponse(), nil)

		res := resourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBServerSslData())

		err := resourceNsxtPolicyLBServerSslProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBServerSslProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBServerSslMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbServerSslID).Return(lbServerSslAPIResponse(), nil)

		res := resourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBServerSslData())
		d.SetId(lbServerSslID)

		err := resourceNsxtPolicyLBServerSslProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbServerSslDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbServerSslID).Return(nsxModel.LBServerSslProfile{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBServerSslData())
		d.SetId(lbServerSslID)

		err := resourceNsxtPolicyLBServerSslProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBServerSslData())

		err := resourceNsxtPolicyLBServerSslProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBServerSslProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBServerSslMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(lbServerSslID, gomock.Any()).Return(lbServerSslAPIResponse(), nil),
			mockSDK.EXPECT().Get(lbServerSslID).Return(lbServerSslAPIResponse(), nil),
		)

		res := resourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBServerSslData())
		d.SetId(lbServerSslID)

		err := resourceNsxtPolicyLBServerSslProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBServerSslData())

		err := resourceNsxtPolicyLBServerSslProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBServerSslProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBServerSslMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbServerSslID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBServerSslData())
		d.SetId(lbServerSslID)

		err := resourceNsxtPolicyLBServerSslProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBServerSslData())

		err := resourceNsxtPolicyLBServerSslProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
