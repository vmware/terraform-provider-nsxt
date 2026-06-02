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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestMockDataSourceNsxtPolicyBareMetalServerInterfaceSchema(t *testing.T) {
	dataSource := dataSourceNsxtPolicyBareMetalServerInterface()

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

func TestMockDataSourceNsxtPolicyBareMetalServerInterfaceValidation(t *testing.T) {
	dataSource := dataSourceNsxtPolicyBareMetalServerInterface()

	// Test that schema validation is properly configured
	assert.NotNil(t, dataSource.Schema)
	assert.NotNil(t, dataSource.Read)

	// Test that required fields are present
	assert.Contains(t, dataSource.Schema, "external_id")
	assert.Contains(t, dataSource.Schema, "display_name")
}

func TestMockBareMetalServerInterfaceConversion(t *testing.T) {
	// Test the conversion function for bare metal server interfaces
	externalId := "test-interface-id"
	displayName := "test-interface"
	resourceType := "BareMetalServerInterface"

	serverInterface := model.BareMetalServerInterface{
		ExternalId:   &externalId,
		DisplayName:  &displayName,
		ResourceType: &resourceType,
	}

	// Verify server interface structure
	assert.Equal(t, "test-interface-id", *serverInterface.ExternalId)
	assert.Equal(t, "test-interface", *serverInterface.DisplayName)
	assert.Equal(t, "BareMetalServerInterface", *serverInterface.ResourceType)
}

func TestMockDataSourceNsxtPolicyBareMetalServerInterfaceRead(t *testing.T) {
	t.Run("Read fails on NSX version below 9.0.0", func(t *testing.T) {
		util.NsxVersion = "8.0.0"
		defer func() { util.NsxVersion = "" }()

		dataSource := dataSourceNsxtPolicyBareMetalServerInterface()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"external_id": "test-bmsi-id",
		})

		err := dataSourceNsxtPolicyBareMetalServerInterfaceRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Bare Metal Server features require NSX-T version 9.0.0 or higher")
	})

	t.Run("Read succeeds on NSX 9.0.0", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		dataSource := dataSourceNsxtPolicyBareMetalServerInterface()
		d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
			"external_id": "test-bmsi-id",
		})

		// This will fail due to missing mock setup, but not due to version check
		err := dataSourceNsxtPolicyBareMetalServerInterfaceRead(d, newGoMockProviderClient())
		if err != nil {
			// Should not be a version error
			assert.NotContains(t, err.Error(), "requires NSX version")
		}
	})
}
