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
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyClusterSecurityConfigRead(t *testing.T) {
	ds := dataSourceNsxtPolicyClusterSecurityConfig()

	t.Run("Read_fails_below_version_9_1_0", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"cluster_id": "cluster-1",
		})
		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyClusterSecurityConfigRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.1.0")
	})

	t.Run("Read_fails_when_cluster_id_empty", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"cluster_id": "",
		})
		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyClusterSecurityConfigRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cluster_id is required")
	})

	t.Run("Read_succeeds_and_sets_dfw_enabled", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()
		mockClient := setupClusterSecurityConfigMock(t)

		feature := "DFW"
		enabled := true
		displayName := "cluster-1-security-config"
		description := "cluster security config"
		path := "/infra/sites/default/enforcement-points/default/cluster-configs/cluster-1"
		mockClient.EXPECT().Get("cluster-1", gomock.Any()).Return(nsxModel.ClusterSecurityConfiguration{
			DisplayName: &displayName,
			Description: &description,
			Path:        &path,
			Features: []nsxModel.ClusterSecurityFeature{
				{Feature: &feature, Enabled: &enabled},
			},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"cluster_id": "cluster-1",
		})
		err := dataSourceNsxtPolicyClusterSecurityConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "cluster-1", d.Id())
		assert.Equal(t, true, d.Get("dfw_enabled"))
		assert.Equal(t, displayName, d.Get("display_name"))
		assert.Equal(t, description, d.Get("description"))
		assert.Equal(t, path, d.Get("path"))
	})

	t.Run("Read_fails_when_API_returns_an_error", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()
		mockClient := setupClusterSecurityConfigMock(t)

		mockClient.EXPECT().Get("cluster-1", gomock.Any()).Return(nsxModel.ClusterSecurityConfiguration{}, assert.AnError)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"cluster_id": "cluster-1",
		})
		err := dataSourceNsxtPolicyClusterSecurityConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
