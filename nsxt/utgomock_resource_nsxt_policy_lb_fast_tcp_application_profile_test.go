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

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	infraMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

var (
	lbFastTcpID          = "lb-fast-tcp-1"
	lbFastTcpDisplayName = "Test LB Fast TCP Profile"
)

func minimalLBFastTcpData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbFastTcpDisplayName,
		"nsx_id":       lbFastTcpID,
	}
}

func setupLBAppProfileMock(t *testing.T, ctrl *gomock.Controller) (*infraMocks.MockLbAppProfilesClient, func()) {
	mockSDK := infraMocks.NewMockLbAppProfilesClient(ctrl)
	mockWrapper := &infraapi.LBAppProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliLbAppProfilesClient
	cliLbAppProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.LBAppProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliLbAppProfilesClient = original }
}

func TestMockResourceNsxtPolicyLBFastTcpApplicationProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		// Return nil error but non-nil struct to indicate object exists
		mockSDK.EXPECT().Get(lbFastTcpID).Return(nil, nil)

		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())

		err := resourceNsxtPolicyLBTcpApplicationProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyLBFastTcpApplicationProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbFastTcpID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())
		d.SetId(lbFastTcpID)

		err := resourceNsxtPolicyLBTcpApplicationProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbFastTcpID).Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())
		d.SetId(lbFastTcpID)

		err := resourceNsxtPolicyLBTcpApplicationProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())

		err := resourceNsxtPolicyLBTcpApplicationProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBFastTcpApplicationProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())

		err := resourceNsxtPolicyLBTcpApplicationProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBFastTcpApplicationProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbFastTcpID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())
		d.SetId(lbFastTcpID)

		err := resourceNsxtPolicyLBTcpApplicationProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBFastTcpApplicationProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBFastTcpData())

		err := resourceNsxtPolicyLBTcpApplicationProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
