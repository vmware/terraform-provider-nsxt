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

// dataSourceNsxtPolicyIdpsSignatureDiffRead calls
// custom_signature_versions.NewCustomSignaturesDiffClient(connector) directly rather than
// through an injectable wrapper variable, so only the guard condition that doesn't require
// an API call can be exercised here (same limitation as the sibling idps data sources).

func TestMockDataSourceNsxtPolicyIdpsSignatureDiffReadGuard(t *testing.T) {
	t.Run("Read fails for global manager", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIdpsSignatureDiff()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"signature_version_id": "default",
		})

		err := dataSourceNsxtPolicyIdpsSignatureDiffRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Global Manager")
	})
}
