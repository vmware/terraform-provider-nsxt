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

// dataSourceNsxtPolicyIdpsSignatureVersionRead calls
// intrusion_services.NewSignatureVersionsClient(connector) directly rather than through an
// injectable wrapper variable, so only guard/validation branches that don't require an API
// call can be exercised here (same limitation as the sibling idps data sources).

func TestMockDataSourceNsxtPolicyIdpsSignatureVersionReadGuard(t *testing.T) {
	t.Run("Read fails for global manager", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Global Manager")
	})
}

func TestMockDataSourceNsxtPolicyIdpsSignatureVersionReadMissingSelector(t *testing.T) {
	t.Run("Read fails when neither id nor display_name is set", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIdpsSignatureVersion()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsSignatureVersionRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "'id' or 'display_name'")
	})
}
