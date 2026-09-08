//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// dataSourceNsxtPolicyIdpsClusterConfigRead reads through the cliIdsClusterConfigsClient
// package-level var (declared in resource_nsxt_policy_idps_cluster_config.go, see
// utgomock_resource_nsxt_policy_idps_cluster_config_test.go for the mock setup helper), so
// both the guard conditions and the mocked lookups are covered here.

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
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

func TestMockDataSourceNsxtPolicyIdpsClusterConfigReadMocked(t *testing.T) {
	ds := dataSourceNsxtPolicyIdpsClusterConfig()

	t.Run("Read by id succeeds", func(t *testing.T) {
		mockClient := setupIdpsClusterConfigMock(t)
		id := "cfg-1"
		targetID, targetType := "domain-c1", "VC_Cluster"
		mockClient.EXPECT().Get("cfg-1", nil).Return(model.IdsClusterConfig{
			Id:      &id,
			Cluster: &model.PolicyResourceReference{TargetId: &targetID, TargetType: &targetType},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "cfg-1"})

		err := dataSourceNsxtPolicyIdpsClusterConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "cfg-1", d.Id())
	})

	t.Run("Read by id propagates API error", func(t *testing.T) {
		mockClient := setupIdpsClusterConfigMock(t)
		mockClient.EXPECT().Get("cfg-1", nil).Return(model.IdsClusterConfig{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "cfg-1"})

		err := dataSourceNsxtPolicyIdpsClusterConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by display_name matches a perfect match over a prefix match", func(t *testing.T) {
		mockClient := setupIdpsClusterConfigMock(t)
		perfectID, prefixID := "cfg-1", "cfg-2"
		perfectName, prefixName := "myconfig", "myconfig-extra"
		mockClient.EXPECT().List(nil, nil, nil, nil, nil, nil, nil, nil).Return(
			model.IdsClusterConfigListResult{Results: []model.IdsClusterConfig{
				{Id: &prefixID, DisplayName: &prefixName},
				{Id: &perfectID, DisplayName: &perfectName},
			}}, nil,
		)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "myconfig"})

		err := dataSourceNsxtPolicyIdpsClusterConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "cfg-1", d.Id())
	})

	t.Run("Read by display_name with only a prefix match succeeds", func(t *testing.T) {
		mockClient := setupIdpsClusterConfigMock(t)
		prefixID := "cfg-2"
		prefixName := "myconfig-extra"
		mockClient.EXPECT().List(nil, nil, nil, nil, nil, nil, nil, nil).Return(
			model.IdsClusterConfigListResult{Results: []model.IdsClusterConfig{{Id: &prefixID, DisplayName: &prefixName}}}, nil,
		)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "myconfig"})

		err := dataSourceNsxtPolicyIdpsClusterConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "cfg-2", d.Id())
	})

	t.Run("Read by display_name with no match fails", func(t *testing.T) {
		mockClient := setupIdpsClusterConfigMock(t)
		mockClient.EXPECT().List(nil, nil, nil, nil, nil, nil, nil, nil).Return(model.IdsClusterConfigListResult{}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "myconfig"})

		err := dataSourceNsxtPolicyIdpsClusterConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("List error propagates", func(t *testing.T) {
		mockClient := setupIdpsClusterConfigMock(t)
		mockClient.EXPECT().List(nil, nil, nil, nil, nil, nil, nil, nil).Return(
			model.IdsClusterConfigListResult{}, vapiErrors.InternalServerError{},
		)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "myconfig"})

		err := dataSourceNsxtPolicyIdpsClusterConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
