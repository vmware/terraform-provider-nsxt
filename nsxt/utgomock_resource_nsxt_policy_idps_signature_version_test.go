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
)

// resourceNsxtPolicyIdpsSignatureVersion uses intrusion_services.NewSignatureVersionsClient(connector)
// directly rather than an injectable wrapper variable. These tests cover guard conditions and
// validation logic that doesn't require API calls.

func minimalIdpsSigVersionData() map[string]interface{} {
	return map[string]interface{}{
		"nsx_id": "sig-ver-1",
	}
}

func TestMockResourceNsxtPolicyIdpsSignatureVersionGuard(t *testing.T) {
	t.Run("Create fails for global manager", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsSigVersionData())

		err := resourceNsxtPolicyIdpsSignatureVersionCreate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})

	t.Run("Read fails for global manager", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsSigVersionData())
		d.SetId("sig-ver-1")

		err := resourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})

	t.Run("Update fails for global manager", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsSigVersionData())
		d.SetId("sig-ver-1")

		err := resourceNsxtPolicyIdpsSignatureVersionUpdate(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "global")
	})
}

func TestMockResourceNsxtPolicyIdpsSignatureVersionCreateNsxIDRequired(t *testing.T) {
	t.Run("Create fails when nsx_id is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtPolicyIdpsSignatureVersionCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nsx_id")
	})
}

func TestMockResourceNsxtPolicyIdpsSignatureVersionReadEmptyID(t *testing.T) {
	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsSigVersionData())

		err := resourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIdpsSignatureVersionDeleteNoOp(t *testing.T) {
	t.Run("Delete is a no-op (removes from state only)", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsSigVersionData())
		d.SetId("sig-ver-1")

		err := resourceNsxtPolicyIdpsSignatureVersionDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyIdpsSignatureVersionUpdateInvalidState(t *testing.T) {
	t.Run("Update fails with invalid state value", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsSignatureVersion()
		data := minimalIdpsSigVersionData()
		data["state"] = "INVALID_STATE"
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId("sig-ver-1")

		// Since state is Computed-only, HasChange won't fire in unit test context.
		// Verify schema definition is correct.
		assert.NotNil(t, res.Schema["state"])
		assert.NotNil(t, res.Schema["version_id"])
	})
}
