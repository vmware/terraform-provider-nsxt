// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// resourceNsxtPolicyIdpsCustomSignature uses SDK clients directly (not injected via package-level
// variables), so full mock-based CRUD tests are not possible. These tests cover guard conditions,
// composite ID parsing, and schema validation.

func minimalIdpsCustomSigData() map[string]interface{} {
	return map[string]interface{}{
		"signature_version_id": "default",
		"signature":            `alert tcp any any -> any 80 (msg:"Test"; sid:9000001; rev:1;)`,
	}
}

func TestMockResourceNsxtPolicyIdpsCustomSignatureGuard(t *testing.T) {
	t.Run("Create fails for global manager", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsCustomSignature()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsCustomSigData())

		err := resourceNsxtPolicyIdpsCustomSignatureCreate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})

	t.Run("Read fails for global manager", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsCustomSignature()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsCustomSigData())
		d.SetId("default/sig-1")

		err := resourceNsxtPolicyIdpsCustomSignatureRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})

	t.Run("Update fails for global manager", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsCustomSignature()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsCustomSigData())
		d.SetId("default/sig-1")

		err := resourceNsxtPolicyIdpsCustomSignatureUpdate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})

	t.Run("Delete fails for global manager", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsCustomSignature()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsCustomSigData())
		d.SetId("default/sig-1")

		err := resourceNsxtPolicyIdpsCustomSignatureDelete(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})
}

func TestMockResourceNsxtPolicyIdpsCustomSignatureIDParsing(t *testing.T) {
	t.Run("Composite ID parsed correctly", func(t *testing.T) {
		versionID, sigID, err := parseCustomSignatureCompositeID("myversion/sig-123")
		require.NoError(t, err)
		assert.Equal(t, "myversion", versionID)
		assert.Equal(t, "sig-123", sigID)
	})

	t.Run("Invalid composite ID returns error", func(t *testing.T) {
		_, _, err := parseCustomSignatureCompositeID("noseparator")
		require.Error(t, err)
	})

	t.Run("Legacy ID (no slash) treated as default version", func(t *testing.T) {
		versionID, sigID, legacy, err := parseCustomSignatureCompositeIDOrLegacy("legacy-sig-id")
		require.NoError(t, err)
		assert.True(t, legacy)
		assert.Equal(t, "default", versionID)
		assert.Equal(t, "legacy-sig-id", sigID)
	})

	t.Run("Empty ID returns error", func(t *testing.T) {
		_, _, _, err := parseCustomSignatureCompositeIDOrLegacy("")
		require.Error(t, err)
	})
}
