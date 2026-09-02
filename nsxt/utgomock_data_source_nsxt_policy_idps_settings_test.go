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

// dataSourceNsxtPolicyIdpsSettingsRead calls security.NewIntrusionServicesClient(connector)
// directly rather than through an injectable wrapper variable (same as
// resourceNsxtPolicyIdpsSettings, see utgomock_resource_nsxt_policy_idps_settings_test.go),
// so only the guard condition that doesn't require an API call can be exercised here.

func TestMockDataSourceNsxtPolicyIdpsSettingsReadGuard(t *testing.T) {
	t.Run("Read fails for global manager", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIdpsSettings()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsSettingsRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Global Manager")
	})
}
