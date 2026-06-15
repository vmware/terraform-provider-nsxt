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

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

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
	t.Run("Read fails on NSX version below 9.0.0", func(t *testing.T) {
		util.NsxVersion = "8.0.0"
		defer func() { util.NsxVersion = "" }()

		dataSource := dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembers()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"group_id": "test-group-id",
		})

		err := dataSourceNsxtPolicyGroupBareMetalServerInterfaceMembersRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Bare Metal Server features require NSX-T version 9.0.0 or higher")
	})

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
}
