// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

// resourceNsxtPolicyIdpsClusterConfig uses intrusion_services.NewClusterConfigsClient(connector)
// directly rather than an injectable wrapper variable, so full mock-based CRUD tests are not
// possible. These tests cover guard conditions and schema validation.

func minimalIdpsClusterConfigData() map[string]interface{} {
	return map[string]interface{}{
		"ids_enabled": true,
		"cluster": []interface{}{
			map[string]interface{}{
				"target_id":   "domain-c1",
				"target_type": "VC_Cluster",
			},
		},
	}
}

func TestMockResourceNsxtPolicyIdpsClusterConfigReadEmptyID(t *testing.T) {
	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsClusterConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsClusterConfigData())

		err := resourceNsxtPolicyIdpsClusterConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIdpsClusterConfigUpdateEmptyID(t *testing.T) {
	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsClusterConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsClusterConfigData())

		err := resourceNsxtPolicyIdpsClusterConfigUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIdpsClusterConfigDeleteEmptyID(t *testing.T) {
	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIdpsClusterConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsClusterConfigData())

		err := resourceNsxtPolicyIdpsClusterConfigDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
