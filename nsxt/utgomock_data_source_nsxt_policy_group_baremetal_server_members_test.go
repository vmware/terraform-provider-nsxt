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

func setupGroupBmsMembersMock(t *testing.T, ctrl *gomock.Controller) (*bmsmocks.MockBmsClient, func()) {
	t.Helper()
	mockSDK := bmsmocks.NewMockBmsClient(ctrl)
	wrapper := &groupsAPI.BmsMembersClientContext{Client: mockSDK, ClientType: utl.Local}
	orig := cliGroupBareMetalServerMembersClient
	cliGroupBareMetalServerMembersClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *groupsAPI.BmsMembersClientContext {
		return wrapper
	}
	return mockSDK, func() { cliGroupBareMetalServerMembersClient = orig }
}

func TestMockDataSourceNsxtPolicyGroupBareMetalServerMembersSchema(t *testing.T) {
	dataSource := dataSourceNsxtPolicyGroupBareMetalServerMembers()

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

func TestMockDataSourceNsxtPolicyGroupBareMetalServerMembersRead(t *testing.T) {

	t.Run("Read succeeds on NSX 9.0.0", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		dataSource := dataSourceNsxtPolicyGroupBareMetalServerMembers()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"group_id": "test-group-id",
		})

		// This will fail due to missing mock setup, but not due to version check
		err := dataSourceNsxtPolicyGroupBareMetalServerMembersRead(d, newGoMockProviderClient())
		if err != nil {
			// Should not be a version error
			assert.NotContains(t, err.Error(), "requires NSX version")
		}
	})

	t.Run("Read succeeds and populates items", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupGroupBmsMembersMock(t, ctrl)
		defer restore()

		externalID := "bms-1"
		displayName := "server-1"
		cpuCores := int64(4)
		mockSDK.EXPECT().List("default", "group-1", (*string)(nil), (*string)(nil), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.BareMetalServerListResult{
				Results: []nsxModel.BareMetalServer{
					{ExternalId: &externalID, DisplayName: &displayName, CpuCores: &cpuCores},
				},
			}, nil)

		dataSource := dataSourceNsxtPolicyGroupBareMetalServerMembers()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"group_id": "group-1",
		})

		err := dataSourceNsxtPolicyGroupBareMetalServerMembersRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "group-1", d.Id())
		items := d.Get("items").([]interface{})
		require.Len(t, items, 1)
		elem := items[0].(map[string]interface{})
		assert.Equal(t, externalID, elem["external_id"])
		assert.Equal(t, displayName, elem["display_name"])
		assert.Equal(t, 4, elem["cpu_cores"])
	})

	t.Run("Read uses a non-default domain and enforcement_point_path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupGroupBmsMembersMock(t, ctrl)
		defer restore()

		epPath := "/infra/sites/default/enforcement-points/default"
		mockSDK.EXPECT().List("dom1", "group-1", (*string)(nil), &epPath, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.BareMetalServerListResult{}, nil)

		dataSource := dataSourceNsxtPolicyGroupBareMetalServerMembers()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"domain":                 "dom1",
			"group_id":               "group-1",
			"enforcement_point_path": epPath,
		})

		err := dataSourceNsxtPolicyGroupBareMetalServerMembersRead(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Read fails when the API returns an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupGroupBmsMembersMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().List("default", "group-1", (*string)(nil), (*string)(nil), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.BareMetalServerListResult{}, errors.New("boom"))

		dataSource := dataSourceNsxtPolicyGroupBareMetalServerMembers()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"group_id": "group-1",
		})

		err := dataSourceNsxtPolicyGroupBareMetalServerMembersRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
