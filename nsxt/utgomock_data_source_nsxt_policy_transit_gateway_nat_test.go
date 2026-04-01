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

func TestMockDataSourceNsxtPolicyTransitGatewayNatInvalidPath(t *testing.T) {
	t.Run("Invalid transit_gateway_path returns error immediately", func(t *testing.T) {
		ds := dataSourceNsxtPolicyTransitGatewayNat()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"transit_gateway_path": "/invalid/path",
		})

		err := dataSourceNsxtPolicyTransitGatewayNatRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid transit_gateway_path")
	})

	t.Run("Too-short transit_gateway_path returns error immediately", func(t *testing.T) {
		ds := dataSourceNsxtPolicyTransitGatewayNat()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"transit_gateway_path": "/orgs/default",
		})

		err := dataSourceNsxtPolicyTransitGatewayNatRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid transit_gateway_path")
	})
}
