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
	l7ProfileID          = "l7-profile-1"
	l7ProfileDisplayName = "Test L7 Access Profile"
	l7ProfileDescription = "Test l7 access profile"
	l7ProfileRevision    = int64(1)
)

func l7ProfileAPIResponse() nsxModel.L7AccessProfile {
	defaultAction := nsxModel.L7AccessProfile_DEFAULT_ACTION_ALLOW
	return nsxModel.L7AccessProfile{
		Id:            &l7ProfileID,
		DisplayName:   &l7ProfileDisplayName,
		Description:   &l7ProfileDescription,
		Revision:      &l7ProfileRevision,
		DefaultAction: &defaultAction,
	}
}

func minimalL7ProfileData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":   l7ProfileDisplayName,
		"description":    l7ProfileDescription,
		"nsx_id":         l7ProfileID,
		"default_action": nsxModel.L7AccessProfile_DEFAULT_ACTION_ALLOW,
	}
}

func setupL7ProfileMock(t *testing.T, ctrl *gomock.Controller) (*infraMocks.MockL7AccessProfilesClient, func()) {
	mockSDK := infraMocks.NewMockL7AccessProfilesClient(ctrl)
	mockWrapper := &infraapi.L7AccessProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliL7AccessProfilesClient
	cliL7AccessProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.L7AccessProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliL7AccessProfilesClient = original }
}

func TestMockResourceNsxtPolicyL7AccessProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL7ProfileMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(l7ProfileID).Return(l7ProfileAPIResponse(), nil)

		res := resourceNsxtPolicyL7AccessProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL7ProfileData())
		d.SetId(l7ProfileID)

		err := resourceNsxtPolicyL7AccessProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, l7ProfileDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(l7ProfileID).Return(nsxModel.L7AccessProfile{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyL7AccessProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL7ProfileData())
		d.SetId(l7ProfileID)

		err := resourceNsxtPolicyL7AccessProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyL7AccessProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL7ProfileData())

		err := resourceNsxtPolicyL7AccessProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyL7AccessProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL7ProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(l7ProfileID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyL7AccessProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL7ProfileData())
		d.SetId(l7ProfileID)

		err := resourceNsxtPolicyL7AccessProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyL7AccessProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL7ProfileData())

		err := resourceNsxtPolicyL7AccessProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
