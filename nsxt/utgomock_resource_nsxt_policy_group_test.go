//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	domainsapi "github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	domainmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
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

// ─── BMS helpers ─────────────────────────────────────────────────────────────

// groupTypeMatcher is a gomock.Matcher that verifies the GroupType field of a
// nsxModel.Group passed to the mock Patch call.
type groupTypeMatcher struct {
	expected string
}

func (m groupTypeMatcher) Matches(x interface{}) bool {
	obj, ok := x.(nsxModel.Group)
	return ok && len(obj.GroupType) == 1 && obj.GroupType[0] == m.expected
}

func (m groupTypeMatcher) String() string {
	return fmt.Sprintf("Group.GroupType == [%s]", m.expected)
}

// groupAPIResponseWithType returns a Group model that has GroupType set.
func groupAPIResponseWithType(gt string) nsxModel.Group {
	resp := groupAPIResponse()
	resp.GroupType = []string{gt}
	return resp
}

// bmsGroupData returns schema input data for a BareMetalServer group.
func bmsGroupData() map[string]interface{} {
	d := minimalGroupData()
	d["group_type"] = nsxModel.Group_GROUP_TYPE_BAREMETALSERVER
	return d
}

// ─── BMS Create tests ─────────────────────────────────────────────────────────

func TestMockResourceNsxtPolicyGroupBMSCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGroupMock(t, ctrl)
	defer restore()

	t.Run("BMS group_type rejected on NSX 4.2", func(t *testing.T) {
		util.NsxVersion = "4.2.0"
		defer func() { util.NsxVersion = "" }()

		// Existence check still happens before the version guard
		mockSDK.EXPECT().Get(groupDomain, groupID).Return(nsxModel.Group{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, bmsGroupData())

		err := resourceNsxtPolicyGroupCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Bare Metal Server features require NSX-T version 9.0.0 or higher")
	})

	t.Run("BMS group_type create succeeds on NSX 9.0", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		gomock.InOrder(
			mockSDK.EXPECT().Get(groupDomain, groupID).Return(nsxModel.Group{}, vapiErrors.NotFound{}),
			// Verify group_type is passed to the API
			mockSDK.EXPECT().Patch(groupDomain, groupID, groupTypeMatcher{nsxModel.Group_GROUP_TYPE_BAREMETALSERVER}).Return(nil),
			mockSDK.EXPECT().Get(groupDomain, groupID).Return(groupAPIResponseWithType(nsxModel.Group_GROUP_TYPE_BAREMETALSERVER), nil),
		)

		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, bmsGroupData())

		err := resourceNsxtPolicyGroupCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, nsxModel.Group_GROUP_TYPE_BAREMETALSERVER, d.Get("group_type"))
	})

	t.Run("non-BMS group_type create still works on NSX 3.2", func(t *testing.T) {
		util.NsxVersion = "3.2.0"
		defer func() { util.NsxVersion = "" }()

		ipGroupData := minimalGroupData()
		ipGroupData["group_type"] = nsxModel.Group_GROUP_TYPE_IPADDRESS

		gomock.InOrder(
			mockSDK.EXPECT().Get(groupDomain, groupID).Return(nsxModel.Group{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(groupDomain, groupID, groupTypeMatcher{nsxModel.Group_GROUP_TYPE_IPADDRESS}).Return(nil),
			mockSDK.EXPECT().Get(groupDomain, groupID).Return(groupAPIResponseWithType(nsxModel.Group_GROUP_TYPE_IPADDRESS), nil),
		)

		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, ipGroupData)

		err := resourceNsxtPolicyGroupCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, nsxModel.Group_GROUP_TYPE_IPADDRESS, d.Get("group_type"))
	})
}

// ─── BMS Read tests ───────────────────────────────────────────────────────────

func TestMockResourceNsxtPolicyGroupBMSRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGroupMock(t, ctrl)
	defer restore()

	t.Run("BMS group_type is populated in state on read", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		mockSDK.EXPECT().Get(groupDomain, groupID).Return(
			groupAPIResponseWithType(nsxModel.Group_GROUP_TYPE_BAREMETALSERVER), nil,
		)

		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGroupData())
		d.SetId(groupID)

		err := resourceNsxtPolicyGroupRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, nsxModel.Group_GROUP_TYPE_BAREMETALSERVER, d.Get("group_type"))
	})

	t.Run("group without group_type has empty group_type in state", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		mockSDK.EXPECT().Get(groupDomain, groupID).Return(groupAPIResponse(), nil)

		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGroupData())
		d.SetId(groupID)

		err := resourceNsxtPolicyGroupRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Get("group_type"))
	})

	t.Run("group_type not set when NSX version below 3.2", func(t *testing.T) {
		util.NsxVersion = "3.1.0"
		defer func() { util.NsxVersion = "" }()

		// API returns a group with group_type but provider must ignore it below 3.2
		mockSDK.EXPECT().Get(groupDomain, groupID).Return(
			groupAPIResponseWithType(nsxModel.Group_GROUP_TYPE_BAREMETALSERVER), nil,
		)

		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGroupData())
		d.SetId(groupID)

		err := resourceNsxtPolicyGroupRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		// group_type is guarded by NsxVersionHigherOrEqual("3.2.0") in the read path
		assert.Equal(t, "", d.Get("group_type"))
	})
}

// ─── BMS Update tests ─────────────────────────────────────────────────────────

func TestMockResourceNsxtPolicyGroupBMSUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGroupMock(t, ctrl)
	defer restore()

	t.Run("BMS group_type update rejected on NSX 4.2", func(t *testing.T) {
		util.NsxVersion = "4.2.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, bmsGroupData())
		d.SetId(groupID)

		err := resourceNsxtPolicyGroupUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Bare Metal Server features require NSX-T version 9.0.0 or higher")
	})

	t.Run("BMS group_type update succeeds on NSX 9.0", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		gomock.InOrder(
			mockSDK.EXPECT().Patch(groupDomain, groupID, groupTypeMatcher{nsxModel.Group_GROUP_TYPE_BAREMETALSERVER}).Return(nil),
			mockSDK.EXPECT().Get(groupDomain, groupID).Return(groupAPIResponseWithType(nsxModel.Group_GROUP_TYPE_BAREMETALSERVER), nil),
		)

		res := resourceNsxtPolicyGroup()
		d := schema.TestResourceDataRaw(t, res.Schema, bmsGroupData())
		d.SetId(groupID)

		err := resourceNsxtPolicyGroupUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, nsxModel.Group_GROUP_TYPE_BAREMETALSERVER, d.Get("group_type"))
	})
}
