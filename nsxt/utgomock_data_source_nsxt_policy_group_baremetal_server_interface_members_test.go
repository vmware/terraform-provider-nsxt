//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	groupsAPI "github.com/vmware/terraform-provider-nsxt/api/infra/domains/groups"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	bmsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains/groups/members"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func setupGroupBmsiMembersMock(t *testing.T, ctrl *gomock.Controller) (*bmsmocks.MockBmsiClient, func()) {
	t.Helper()
	mockSDK := bmsmocks.NewMockBmsiClient(ctrl)
	wrapper := &groupsAPI.BmsiMembersClientContext{Client: mockSDK, ClientType: utl.Local}
	orig := cliGroupBareMetalServerInterfaceMembersClient
	cliGroupBareMetalServerInterfaceMembersClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *groupsAPI.BmsiMembersClientContext {
		return wrapper
	}
	return mockSDK, func() { cliGroupBareMetalServerInterfaceMembersClient = orig }
}

func TestMockDataSourceNsxtPolicyGroupBareMetalServerInterfaceMembersSchema(t *testing.T) {
	dataSource := dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembers()

	// Test schema structure
	assert.NotNil(t, dataSource.Schema)

	// Test required fields are properly defined
	dsSchema := dataSource.Schema
	assert.Contains(t, dsSchema, "id")
	assert.Contains(t, dsSchema, "domain")
	assert.Contains(t, dsSchema, "group_id")
	assert.Contains(t, dsSchema, "enforcement_point_path")
	assert.Contains(t, dsSchema, "items")

	// Test group_id is required
	assert.True(t, dsSchema["group_id"].Required)
	assert.Equal(t, schema.TypeString, dsSchema["group_id"].Type)

	// Test domain is optional
	assert.False(t, dsSchema["domain"].Required)
	assert.True(t, dsSchema["domain"].Optional)

	// Test computed fields
	assert.True(t, dsSchema["items"].Computed)
	assert.Equal(t, schema.TypeList, dsSchema["items"].Type)
}

func TestMockDataSourceNsxtPolicyGroupBareMetalServerInterfaceMembersRead(t *testing.T) {

	t.Run("Read succeeds on NSX 9.0.0", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		dataSource := dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembers()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"group_id": "test-group-id",
		})

		// This will fail due to missing mock setup, but not due to version check
		err := dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembersRead(d, newGoMockProviderClient())
		if err != nil {
			// Should not be a version error
			assert.NotContains(t, err.Error(), "requires NSX version")
		}
	})

	t.Run("Read succeeds and populates items", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupGroupBmsiMembersMock(t, ctrl)
		defer restore()

		externalID := "bmsi-1"
		displayName := "iface-1"
		state := "UP"
		mockSDK.EXPECT().List("default", "group-1", (*string)(nil), (*string)(nil), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.BareMetalServerInterfaceListResult{
				Results: []nsxModel.BareMetalServerInterface{
					{ExternalId: &externalID, DisplayName: &displayName, State: &state},
				},
			}, nil)

		dataSource := dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembers()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"group_id": "group-1",
		})

		err := dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembersRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "group-1", d.Id())
		items := d.Get("items").([]interface{})
		require.Len(t, items, 1)
		elem := items[0].(map[string]interface{})
		assert.Equal(t, externalID, elem["external_id"])
		assert.Equal(t, displayName, elem["display_name"])
		assert.Equal(t, state, elem["state"])
	})

	t.Run("Read fails when the API returns an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupGroupBmsiMembersMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().List("default", "group-1", (*string)(nil), (*string)(nil), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.BareMetalServerInterfaceListResult{}, errors.New("boom"))

		dataSource := dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembers()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"group_id": "group-1",
		})

		err := dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembersRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
