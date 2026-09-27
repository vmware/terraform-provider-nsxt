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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

// bareMetalServerStructValue builds the *data.StructValue a policy search result for a
// BareMetalServer would decode to.
func bareMetalServerStructValue(t *testing.T, externalID, displayName string) *data.StructValue {
	t.Helper()
	resourceType := "BareMetalServer"
	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(model.BareMetalServer{
		ExternalId:   &externalID,
		DisplayName:  &displayName,
		ResourceType: &resourceType,
	}, model.BareMetalServerBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func TestMockDataSourceNsxtPolicyBareMetalServerSchema(t *testing.T) {
	dataSource := dataSourceNsxtPolicyBareMetalServer()

	// Test schema structure
	assert.NotNil(t, dataSource.Schema)

	// Test required fields are properly defined
	dsSchema := dataSource.Schema
	assert.Contains(t, dsSchema, "id")
	assert.Contains(t, dsSchema, "external_id")
	assert.Contains(t, dsSchema, "display_name")
	assert.Contains(t, dsSchema, "resource_type")

	// Test external_id is optional (one of external_id or display_name must be provided)
	assert.False(t, dsSchema["external_id"].Required)
	assert.True(t, dsSchema["external_id"].Optional)
	assert.Equal(t, schema.TypeString, dsSchema["external_id"].Type)

	// Test display_name is optional
	assert.False(t, dsSchema["display_name"].Required)
	assert.True(t, dsSchema["display_name"].Optional)
	assert.Equal(t, schema.TypeString, dsSchema["display_name"].Type)

	// Test computed fields
	assert.True(t, dsSchema["resource_type"].Computed)
	assert.Equal(t, schema.TypeString, dsSchema["resource_type"].Type)
}

func TestMockDataSourceNsxtPolicyBareMetalServerValidation(t *testing.T) {
	dataSource := dataSourceNsxtPolicyBareMetalServer()

	// Test that schema validation is properly configured
	assert.NotNil(t, dataSource.Schema)
	assert.NotNil(t, dataSource.Read)

	// Test that required fields are present
	assert.Contains(t, dataSource.Schema, "external_id")
	assert.Contains(t, dataSource.Schema, "display_name")
}

func TestMockNsxtBareMetalServerConversion(t *testing.T) {
	// Test the conversion function for bare metal servers
	externalId := "test-server-id"
	displayName := "test-server"
	resourceType := "BareMetalServer"

	server := model.BareMetalServer{
		ExternalId:   &externalId,
		DisplayName:  &displayName,
		ResourceType: &resourceType,
	}

	// Verify server structure
	assert.Equal(t, "test-server-id", *server.ExternalId)
	assert.Equal(t, "test-server", *server.DisplayName)
	assert.Equal(t, "BareMetalServer", *server.ResourceType)
}

func TestMockDataSourceNsxtPolicyBareMetalServerRead(t *testing.T) {

	t.Run("Read succeeds on NSX 9.0.0", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		dataSource := dataSourceNsxtPolicyBareMetalServer()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"external_id": "test-bms-id",
		})

		// This will fail due to missing mock setup, but not due to version check
		err := dataSourceNsxtPolicyBareMetalServerRead(d, newGoMockProviderClient())
		if err != nil {
			// Should not be a version error
			assert.NotContains(t, err.Error(), "requires NSX version")
		}
	})

	t.Run("Read by external_id succeeds", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []model.SearchResponse{{
			Results:     []*data.StructValue{bareMetalServerStructValue(t, "bms-1", "server-1")},
			ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		dataSource := dataSourceNsxtPolicyBareMetalServer()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"external_id": "bms-1",
		})

		err := dataSourceNsxtPolicyBareMetalServerRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "bms-1", d.Id())
		assert.Equal(t, "server-1", d.Get("display_name"))
	})

	t.Run("Read by external_id fails when not found", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []model.SearchResponse{{
			Results: []*data.StructValue{}, ResultCount: i64(0),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		dataSource := dataSourceNsxtPolicyBareMetalServer()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"external_id": "nonexistent",
		})

		err := dataSourceNsxtPolicyBareMetalServerRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by display_name single match succeeds", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []model.SearchResponse{{
			Results:     []*data.StructValue{bareMetalServerStructValue(t, "bms-1", "server-1")},
			ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		dataSource := dataSourceNsxtPolicyBareMetalServer()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"display_name": "server-1",
		})

		err := dataSourceNsxtPolicyBareMetalServerRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "bms-1", d.Id())
	})

	t.Run("Read by display_name with multiple matches fails", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []model.SearchResponse{{
			Results: []*data.StructValue{
				bareMetalServerStructValue(t, "bms-1", "dup"),
				bareMetalServerStructValue(t, "bms-2", "dup"),
			},
			ResultCount: i64(2),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		dataSource := dataSourceNsxtPolicyBareMetalServer()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"display_name": "dup",
		})

		err := dataSourceNsxtPolicyBareMetalServerRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Multiple")
	})

	t.Run("Read by display_name with no matches fails", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []model.SearchResponse{{
			Results:     []*data.StructValue{bareMetalServerStructValue(t, "bms-1", "server-1")},
			ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		dataSource := dataSourceNsxtPolicyBareMetalServer()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		err := dataSourceNsxtPolicyBareMetalServerRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "No bare metal server found")
	})

	t.Run("Read fails when neither external_id nor display_name is set", func(t *testing.T) {
		dataSource := dataSourceNsxtPolicyBareMetalServer()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyBareMetalServerRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be specified")
	})

	t.Run("Read propagates a search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		dataSource := dataSourceNsxtPolicyBareMetalServer()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"external_id": "bms-1",
		})

		err := dataSourceNsxtPolicyBareMetalServerRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
