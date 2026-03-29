// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// Note: ClusterSecurityConfig uses security.NewClusterConfigsClient(connector) directly
// (not via a package-level var), so the SDK client cannot be injected for mocking.
// These tests exercise the version-guard logic that runs before any client calls.

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockResourceNsxtPolicyClusterSecurityConfigVersionGuard(t *testing.T) {
	res := resourceNsxtPolicyClusterSecurityConfig()

	t.Run("Create_fails_below_version_9_1_0", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_id":  "cluster-1",
			"dfw_enabled": true,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyClusterSecurityConfigCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.1.0")
	})

	t.Run("Read_fails_below_version_9_1_0", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_id": "cluster-1",
		})
		d.SetId("cluster-1")
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyClusterSecurityConfigRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.1.0")
	})

	t.Run("Update_fails_below_version_9_1_0", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_id":  "cluster-1",
			"dfw_enabled": true,
		})
		d.SetId("cluster-1")
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyClusterSecurityConfigUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.1.0")
	})

	t.Run("Delete_fails_below_version_9_1_0", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_id": "cluster-1",
		})
		d.SetId("cluster-1")
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyClusterSecurityConfigDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.1.0")
	})
}
