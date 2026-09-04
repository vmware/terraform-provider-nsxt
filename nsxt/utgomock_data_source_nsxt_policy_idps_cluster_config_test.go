//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// dataSourceNsxtPolicyIdpsClusterConfigRead uses
// intrusion_services.NewClusterConfigsClient(connector) directly rather than
// an injectable wrapper variable, so full mock-based Read tests are not
// possible. These tests cover the guard conditions that run before any
// client call.

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockDataSourceNsxtPolicyIdpsClusterConfigRead(t *testing.T) {
	ds := dataSourceNsxtPolicyIdpsClusterConfig()

	t.Run("Read fails on global manager", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": "cfg-1",
		})
		c := newGoMockProviderClient()
		c.PolicyGlobalManager = true

		err := dataSourceNsxtPolicyIdpsClusterConfigRead(d, c)
		require.Error(t, err)
	})

	t.Run("Read fails when id and display_name are empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsClusterConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "obtaining IdsClusterConfig")
	})
}
