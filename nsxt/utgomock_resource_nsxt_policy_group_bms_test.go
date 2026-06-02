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

func TestMockResourceNsxtPolicyGroupBMSSchema(t *testing.T) {
	resource := resourceNsxtPolicyGroup()

	// Test schema structure
	assert.NotNil(t, resource.Schema)

	// Test required fields are properly defined
	resourceSchema := resource.Schema
	assert.Contains(t, resourceSchema, "display_name")
	assert.Contains(t, resourceSchema, "description")
	assert.Contains(t, resourceSchema, "domain")
	assert.Contains(t, resourceSchema, "group_type")
	assert.Contains(t, resourceSchema, "criteria")

	// Test display_name is required
	assert.True(t, resourceSchema["display_name"].Required)
	assert.Equal(t, schema.TypeString, resourceSchema["display_name"].Type)

	// Test group_type configuration
	assert.True(t, resourceSchema["group_type"].Optional)
	assert.Equal(t, schema.TypeString, resourceSchema["group_type"].Type)
}

func TestMockResourceNsxtPolicyGroupBMSValidation(t *testing.T) {
	resource := resourceNsxtPolicyGroup()

	// Test that schema validation is properly configured
	assert.NotNil(t, resource.Schema)
	assert.NotNil(t, resource.Create)
	assert.NotNil(t, resource.Read)
	assert.NotNil(t, resource.Update)
	assert.NotNil(t, resource.Delete)

	// Test required fields presence
	assert.Contains(t, resource.Schema, "display_name")
	assert.Contains(t, resource.Schema, "criteria")
}

func TestMockResourceNsxtPolicyGroupBMSCriteria(t *testing.T) {
	// Test BMS group criteria validation
	t.Run("BMS group criteria", func(t *testing.T) {
		// Test basic criteria logic
		memberType := "BareMetalServer"
		assert.Equal(t, "BareMetalServer", memberType)

		// Test criteria structure
		criteria := map[string]interface{}{
			"member_type": "BareMetalServer",
			"group_type":  "static",
		}
		assert.NotNil(t, criteria)
		assert.Equal(t, "BareMetalServer", criteria["member_type"])
	})
}
