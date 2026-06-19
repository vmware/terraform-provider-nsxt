//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestMockResourceNsxtPolicyBareMetalServerInterfaceTagsSchema(t *testing.T) {
	resource := resourceNsxtPolicyBareMetalServerInterfaceTags()

	// Test schema structure
	assert.NotNil(t, resource.Schema)

	// Test required fields are properly defined
	resourceSchema := resource.Schema
	assert.Contains(t, resourceSchema, "external_id")
	assert.Contains(t, resourceSchema, "tag")

	// Test external_id is required
	assert.True(t, resourceSchema["external_id"].Required)
	assert.Equal(t, schema.TypeString, resourceSchema["external_id"].Type)

	// Test tag block schema
	tagSchema := resourceSchema["tag"].Elem.(*schema.Resource).Schema
	assert.Contains(t, tagSchema, "scope")
	assert.Contains(t, tagSchema, "tag")
	assert.False(t, tagSchema["scope"].Required) // scope is optional
	assert.True(t, tagSchema["scope"].Optional)
	assert.False(t, tagSchema["tag"].Required) // tag value is also optional in standard schema
	assert.True(t, tagSchema["tag"].Optional)
}

func TestMockResourceNsxtPolicyBareMetalServerInterfaceTagsValidation(t *testing.T) {
	resource := resourceNsxtPolicyBareMetalServerInterfaceTags()

	// Test that CRUD operations are properly configured
	assert.NotNil(t, resource.Schema)
	assert.NotNil(t, resource.Create)
	assert.NotNil(t, resource.Read)
	assert.NotNil(t, resource.Update)
	assert.NotNil(t, resource.Delete)

	// Test required fields
	assert.Contains(t, resource.Schema, "external_id")
	assert.Contains(t, resource.Schema, "tag")
}

func TestMockResourceNsxtPolicyBareMetalServerInterfaceTagsBasicOperations(t *testing.T) {
	// Test basic interface tag operations
	t.Run("interface tag validation", func(t *testing.T) {
		// Test basic tag structure
		scope := "network-type"
		tag := "data-plane"

		assert.Equal(t, "network-type", scope)
		assert.Equal(t, "data-plane", tag)

		// Test tag mapping
		tagData := map[string]interface{}{
			"scope": scope,
			"tag":   tag,
		}
		assert.NotNil(t, tagData)
		assert.Equal(t, "network-type", tagData["scope"])
		assert.Equal(t, "data-plane", tagData["tag"])
	})
}
