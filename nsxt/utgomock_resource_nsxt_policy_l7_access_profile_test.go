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

	apipkg "github.com/vmware/terraform-provider-nsxt/api"
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

func setupL7ProfileInfraMock(t *testing.T, ctrl *gomock.Controller) (*infraMocks.MockInfraClient, func()) {
	mockInfraSDK := infraMocks.NewMockInfraClient(ctrl)
	original := cliInfraClient
	cliInfraClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}
	return mockInfraSDK, func() { cliInfraClient = original }
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

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockSDK.EXPECT().Delete(l7ProfileID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyL7AccessProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL7ProfileData())
		d.SetId(l7ProfileID)

		err := resourceNsxtPolicyL7AccessProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func l7ProfileAPIResponseWithEntry() nsxModel.L7AccessProfile {
	obj := l7ProfileAPIResponse()

	entryID := "entry-1"
	entryDisplayName := "entry-1-display"
	entryAction := nsxModel.L7AccessEntry_ACTION_ALLOW
	disabled := false
	logged := true
	sequenceNumber := int64(0)

	attrSource := nsxModel.L7AccessAttributes_ATTRIBUTE_SOURCE_SYSTEM
	attrKey := nsxModel.L7AccessAttributes_KEY_APP_ID
	customUrlPartialMatch := false
	isALGType := false

	obj.L7AccessEntries = []nsxModel.L7AccessEntry{
		{
			Id:             &entryID,
			DisplayName:    &entryDisplayName,
			Action:         &entryAction,
			Disabled:       &disabled,
			Logged:         &logged,
			SequenceNumber: &sequenceNumber,
			Attributes: []nsxModel.L7AccessAttributes{
				{
					AttributeSource:       &attrSource,
					Key:                   &attrKey,
					Value:                 []string{"app1", "app2"},
					CustomUrlPartialMatch: &customUrlPartialMatch,
					IsALGType:             &isALGType,
				},
			},
		},
	}
	return obj
}

func TestMockResourceNsxtPolicyL7AccessProfileReadWithEntries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL7ProfileMock(t, ctrl)
	defer restore()

	t.Run("Read success sets l7_access_entry fields", func(t *testing.T) {
		mockSDK.EXPECT().Get(l7ProfileID).Return(l7ProfileAPIResponseWithEntry(), nil)

		res := resourceNsxtPolicyL7AccessProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL7ProfileData())
		d.SetId(l7ProfileID)

		err := resourceNsxtPolicyL7AccessProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, l7ProfileDisplayName, d.Get("display_name"))
		assert.Equal(t, l7ProfileDescription, d.Get("description"))
		assert.Equal(t, int(l7ProfileRevision), d.Get("revision"))

		entries := d.Get("l7_access_entry").([]interface{})
		require.Len(t, entries, 1)
		entry := entries[0].(map[string]interface{})
		assert.Equal(t, "entry-1", entry["nsx_id"])
		assert.Equal(t, "entry-1-display", entry["display_name"])
		assert.Equal(t, nsxModel.L7AccessEntry_ACTION_ALLOW, entry["action"])
		assert.False(t, entry["disabled"].(bool))
		assert.True(t, entry["logged"].(bool))

		attrs := entry["attribute"].([]interface{})
		require.Len(t, attrs, 1)
		attr := attrs[0].(map[string]interface{})
		assert.Equal(t, nsxModel.L7AccessAttributes_KEY_APP_ID, attr["key"])
		assert.Equal(t, []interface{}{"app1", "app2"}, attr["values"])
	})
}

func TestMockResourceNsxtPolicyL7AccessProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL7ProfileMock(t, ctrl)
	defer restore()
	mockInfraSDK, restoreInfra := setupL7ProfileInfraMock(t, ctrl)
	defer restoreInfra()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(l7ProfileID).Return(nsxModel.L7AccessProfile{}, vapiErrors.NotFound{}),
			mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(l7ProfileID).Return(l7ProfileAPIResponse(), nil),
		)

		res := resourceNsxtPolicyL7AccessProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL7ProfileData())

		err := resourceNsxtPolicyL7AccessProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, l7ProfileID, d.Id())
		assert.Equal(t, l7ProfileID, d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(l7ProfileID).Return(l7ProfileAPIResponse(), nil)

		res := resourceNsxtPolicyL7AccessProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL7ProfileData())

		err := resourceNsxtPolicyL7AccessProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when Patch API errors", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(l7ProfileID).Return(nsxModel.L7AccessProfile{}, vapiErrors.NotFound{}),
			mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtPolicyL7AccessProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL7ProfileData())

		err := resourceNsxtPolicyL7AccessProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyL7AccessProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupL7ProfileMock(t, ctrl)
	defer restore()
	mockInfraSDK, restoreInfra := setupL7ProfileInfraMock(t, ctrl)
	defer restoreInfra()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(l7ProfileID).Return(l7ProfileAPIResponse(), nil),
		)

		res := resourceNsxtPolicyL7AccessProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL7ProfileData())
		d.SetId(l7ProfileID)

		err := resourceNsxtPolicyL7AccessProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, l7ProfileDisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyL7AccessProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL7ProfileData())

		err := resourceNsxtPolicyL7AccessProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Update fails when Patch API errors", func(t *testing.T) {
		mockInfraSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyL7AccessProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalL7ProfileData())
		d.SetId(l7ProfileID)

		err := resourceNsxtPolicyL7AccessProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
