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

// dataSourceNsxtPolicyIdpsSystemSignaturesRead calls
// intrusion_services.NewSignatureVersionsClient(connector) and
// signature_versions.NewSignaturesClient(connector) directly rather than through injectable
// wrapper variables, so only the guard condition that doesn't require an API call can be
// exercised here (same limitation as the sibling idps data sources).

func TestMockDataSourceNsxtPolicyIdpsSystemSignaturesReadGuard(t *testing.T) {
	t.Run("Read fails for global manager", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIdpsSystemSignatures()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsSystemSignaturesRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Global Manager")
	})
}
