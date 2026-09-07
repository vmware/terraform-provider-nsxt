//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/settings/security/ClusterConfigsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/settings/security/ClusterConfigsClient.go ClusterConfigsClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/security"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	clustersecuritymocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/settings/security"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func setupClusterSecurityConfigMock(t *testing.T) *clustersecuritymocks.MockClusterConfigsClient {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockClient := clustersecuritymocks.NewMockClusterConfigsClient(ctrl)

	originalCli := cliClusterSecurityConfigsClient
	t.Cleanup(func() { cliClusterSecurityConfigsClient = originalCli })
	cliClusterSecurityConfigsClient = func(connector client.Connector) security.ClusterConfigsClient {
		return mockClient
	}

	return mockClient
}

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

func TestMockResourceNsxtPolicyClusterSecurityConfigCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	mockClient := setupClusterSecurityConfigMock(t)

	t.Run("Create success", func(t *testing.T) {
		displayName := "cluster-1"
		description := "auto-created"
		revision := int64(1)
		dfwFeature := "DFW"
		dfwEnabled := true

		mockClient.EXPECT().
			Patch("cluster-1", gomock.Any()).
			Return(nil)
		mockClient.EXPECT().
			Get("cluster-1", gomock.Any()).
			Return(model.ClusterSecurityConfiguration{
				DisplayName: &displayName,
				Description: &description,
				Revision:    &revision,
				Features: []model.ClusterSecurityFeature{
					{Feature: &dfwFeature, Enabled: &dfwEnabled},
				},
			}, nil)

		res := resourceNsxtPolicyClusterSecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_id":  "cluster-1",
			"dfw_enabled": true,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyClusterSecurityConfigCreate(d, m)
		require.NoError(t, err)
		assert.Equal(t, "cluster-1", d.Id())
		assert.True(t, d.Get("dfw_enabled").(bool))
	})
}

func TestMockResourceNsxtPolicyClusterSecurityConfigRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	mockClient := setupClusterSecurityConfigMock(t)

	t.Run("Read success sets all fields", func(t *testing.T) {
		displayName := "cluster-1"
		description := "auto-created"
		path := "/infra/settings/security/cluster-configs/cluster-1"
		revision := int64(3)
		dfwFeature := "DFW"
		dfwEnabled := true

		mockClient.EXPECT().
			Get("cluster-1", gomock.Any()).
			Return(model.ClusterSecurityConfiguration{
				DisplayName: &displayName,
				Description: &description,
				Path:        &path,
				Revision:    &revision,
				Features: []model.ClusterSecurityFeature{
					{Feature: &dfwFeature, Enabled: &dfwEnabled},
				},
			}, nil)

		res := resourceNsxtPolicyClusterSecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("cluster-1")

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyClusterSecurityConfigRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, "cluster-1", d.Get("cluster_id"))
		assert.True(t, d.Get("dfw_enabled").(bool))
		assert.Equal(t, displayName, d.Get("display_name"))
		assert.Equal(t, description, d.Get("description"))
		assert.Equal(t, path, d.Get("path"))
		assert.Equal(t, int(revision), d.Get("revision"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyClusterSecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		// d.Id() is empty

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyClusterSecurityConfigRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Cluster Security Config ID")
	})

	t.Run("Read defaults dfw_enabled false when DFW feature absent", func(t *testing.T) {
		displayName := "cluster-1"

		mockClient.EXPECT().
			Get("cluster-1", gomock.Any()).
			Return(model.ClusterSecurityConfiguration{
				DisplayName: &displayName,
				Features:    []model.ClusterSecurityFeature{},
			}, nil)

		res := resourceNsxtPolicyClusterSecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("cluster-1")

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyClusterSecurityConfigRead(d, m)
		require.NoError(t, err)
		assert.False(t, d.Get("dfw_enabled").(bool))
	})
}

func TestMockResourceNsxtPolicyClusterSecurityConfigUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	mockClient := setupClusterSecurityConfigMock(t)

	t.Run("Update success", func(t *testing.T) {
		displayName := "cluster-1"
		dfwFeature := "DFW"
		dfwEnabled := false

		mockClient.EXPECT().
			Patch("cluster-1", gomock.Any()).
			Return(nil)
		mockClient.EXPECT().
			Get("cluster-1", gomock.Any()).
			Return(model.ClusterSecurityConfiguration{
				DisplayName: &displayName,
				Features: []model.ClusterSecurityFeature{
					{Feature: &dfwFeature, Enabled: &dfwEnabled},
				},
			}, nil)

		res := resourceNsxtPolicyClusterSecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_id":  "cluster-1",
			"dfw_enabled": false,
		})
		d.SetId("cluster-1")

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyClusterSecurityConfigUpdate(d, m)
		require.NoError(t, err)
		assert.False(t, d.Get("dfw_enabled").(bool))
	})

	t.Run("Update fails when Patch API errors", func(t *testing.T) {
		mockClient.EXPECT().
			Patch("cluster-1", gomock.Any()).
			Return(errors.New("API error"))

		res := resourceNsxtPolicyClusterSecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"cluster_id":  "cluster-1",
			"dfw_enabled": true,
		})
		d.SetId("cluster-1")

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyClusterSecurityConfigUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyClusterSecurityConfigDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	mockClient := setupClusterSecurityConfigMock(t)

	t.Run("Delete disables all existing features", func(t *testing.T) {
		dfwFeature := "DFW"
		idsFeature := "IDS"
		enabled := true

		mockClient.EXPECT().
			Get("cluster-1", gomock.Any()).
			Return(model.ClusterSecurityConfiguration{
				Features: []model.ClusterSecurityFeature{
					{Feature: &dfwFeature, Enabled: &enabled},
					{Feature: &idsFeature, Enabled: &enabled},
				},
			}, nil)
		mockClient.EXPECT().
			Patch("cluster-1", gomock.Any()).
			DoAndReturn(func(_ string, config model.ClusterSecurityConfiguration) error {
				require.Len(t, config.Features, 2)
				for _, f := range config.Features {
					assert.False(t, *f.Enabled)
				}
				return nil
			})

		res := resourceNsxtPolicyClusterSecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("cluster-1")

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyClusterSecurityConfigDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete is no-op when Get returns NotFound", func(t *testing.T) {
		mockClient.EXPECT().
			Get("cluster-1", gomock.Any()).
			Return(model.ClusterSecurityConfiguration{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyClusterSecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("cluster-1")

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyClusterSecurityConfigDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete swallows other Get errors", func(t *testing.T) {
		mockClient.EXPECT().
			Get("cluster-1", gomock.Any()).
			Return(model.ClusterSecurityConfiguration{}, errors.New("transient error"))

		res := resourceNsxtPolicyClusterSecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("cluster-1")

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyClusterSecurityConfigDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyClusterSecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		// d.Id() is empty

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyClusterSecurityConfigDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Cluster Security Config ID")
	})
}
