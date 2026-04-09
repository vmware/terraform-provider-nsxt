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

	domainsapi "github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	domainmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains"
)

var (
	groupID          = "group-1"
	groupDisplayName = "Test Group"
	groupDescription = "test group"
	groupRevision    = int64(1)
	groupDomain      = "default"
	groupPath        = "/infra/domains/default/groups/group-1"
)

func groupAPIResponse() nsxModel.Group {
	return nsxModel.Group{
		Id:          &groupID,
		DisplayName: &groupDisplayName,
		Description: &groupDescription,
		Revision:    &groupRevision,
		Path:        &groupPath,
	}
}

func minimalGroupData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": groupDisplayName,
		"description":  groupDescription,
		"nsx_id":       groupID,
		"domain":       groupDomain,
	}
}

func setupGroupMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockGroupsClient, func()) {
	mockSDK := domainmocks.NewMockGroupsClient(ctrl)
	mockWrapper := &domainsapi.GroupClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliGroupsClient
	cliGroupsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *domainsapi.GroupClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliGroupsClient = original }
}

func TestMockResourceNsxtPolicyGroupCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGroupMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(groupDomain, groupID).Return(nsxModel.Group{}, notFoundErr),
			mockSDK.EXPECT().Patch(groupDomain, groupID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(groupDomain, groupID).Return(groupAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGroupData())

		err := resourceNsxtPolicyGroupCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, groupID, d.Id())
		assert.Equal(t, groupDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(groupDomain, groupID).Return(groupAPIResponse(), nil)

		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGroupData())

		err := resourceNsxtPolicyGroupCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGroupRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGroupMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(groupDomain, groupID).Return(groupAPIResponse(), nil)

		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGroupData())
		d.SetId(groupID)

		err := resourceNsxtPolicyGroupRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, groupDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(groupDomain, groupID).Return(nsxModel.Group{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGroupData())
		d.SetId(groupID)

		err := resourceNsxtPolicyGroupRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGroupData())

		err := resourceNsxtPolicyGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGroupUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGroupMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(groupDomain, groupID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(groupDomain, groupID).Return(groupAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGroupData())
		d.SetId(groupID)

		err := resourceNsxtPolicyGroupUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGroupData())

		err := resourceNsxtPolicyGroupUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGroupDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGroupMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(groupDomain, groupID, gomock.Any(), gomock.Any()).Return(nil)

		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGroupData())
		d.SetId(groupID)

		err := resourceNsxtPolicyGroupDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGroupData())

		err := resourceNsxtPolicyGroupDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
