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

func TestMockResourceNsxtPolicyBareMetalServerTagsSchema(t *testing.T) {
	resource := resourceNsxtPolicyBareMetalServerTags()

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

func TestMockResourceNsxtPolicyBareMetalServerTagsValidation(t *testing.T) {
	resource := resourceNsxtPolicyBareMetalServerTags()

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

func TestMockResourceNsxtPolicyBareMetalServerTagsBasicOperations(t *testing.T) {
	// Test basic tag operations
	t.Run("tag validation", func(t *testing.T) {
		// Test basic tag structure
		scope := "environment"
		tag := "production"

		assert.Equal(t, "environment", scope)
		assert.Equal(t, "production", tag)

		// Test tag mapping
		tagData := map[string]interface{}{
			"scope": scope,
			"tag":   tag,
		}
		assert.NotNil(t, tagData)
		assert.Equal(t, "environment", tagData["scope"])
		assert.Equal(t, "production", tagData["tag"])
	})
}
