//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/settings/firewall/security/intrusion_services/ClusterConfigsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services/ClusterConfigsClient.go ClusterConfigsClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/settings/firewall/security/intrusion_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	idpsclusterconfigmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/settings/firewall/security/intrusion_services"
)

func setupIdpsClusterConfigMock(t *testing.T) *idpsclusterconfigmocks.MockClusterConfigsClient {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockClient := idpsclusterconfigmocks.NewMockClusterConfigsClient(ctrl)

	originalCli := cliIdsClusterConfigsClient
	t.Cleanup(func() { cliIdsClusterConfigsClient = originalCli })
	cliIdsClusterConfigsClient = func(connector client.Connector) intrusion_services.ClusterConfigsClient {
		return mockClient
	}

	return mockClient
}

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

func TestMockResourceNsxtPolicyIdpsClusterConfigCreate(t *testing.T) {
	mockClient := setupIdpsClusterConfigMock(t)

	t.Run("Create success", func(t *testing.T) {
		displayName := "idps-cfg"
		description := "desc"
		targetID := "domain-c1"
		targetType := "VC_Cluster"
		idsEnabled := true

		mockClient.EXPECT().
			Patch("domain-c1", gomock.Any()).
			Return(nil)
		mockClient.EXPECT().
			Get("domain-c1", gomock.Any()).
			Return(model.IdsClusterConfig{
				DisplayName: &displayName,
				Description: &description,
				IdsEnabled:  &idsEnabled,
				Cluster: &model.PolicyResourceReference{
					TargetId:   &targetID,
					TargetType: &targetType,
				},
			}, nil)

		res := resourceNsxtPolicyIdpsClusterConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsClusterConfigData())

		err := resourceNsxtPolicyIdpsClusterConfigCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "domain-c1", d.Id())
		assert.Equal(t, "domain-c1", d.Get("nsx_id"))
	})

	t.Run("Create fails when Patch API errors", func(t *testing.T) {
		mockClient.EXPECT().
			Patch("domain-c1", gomock.Any()).
			Return(errors.New("API error"))

		res := resourceNsxtPolicyIdpsClusterConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsClusterConfigData())

		err := resourceNsxtPolicyIdpsClusterConfigCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIdpsClusterConfigRead(t *testing.T) {
	mockClient := setupIdpsClusterConfigMock(t)

	t.Run("Read success sets all fields", func(t *testing.T) {
		displayName := "idps-cfg"
		description := "desc"
		path := "/infra/settings/firewall/security/intrusion-services/cluster-configs/domain-c1"
		revision := int64(2)
		targetID := "domain-c1"
		targetType := "VC_Cluster"
		idsEnabled := true

		mockClient.EXPECT().
			Get("domain-c1", gomock.Any()).
			Return(model.IdsClusterConfig{
				DisplayName: &displayName,
				Description: &description,
				Path:        &path,
				Revision:    &revision,
				IdsEnabled:  &idsEnabled,
				Cluster: &model.PolicyResourceReference{
					TargetId:   &targetID,
					TargetType: &targetType,
				},
			}, nil)

		res := resourceNsxtPolicyIdpsClusterConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("domain-c1")

		err := resourceNsxtPolicyIdpsClusterConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, displayName, d.Get("display_name"))
		assert.Equal(t, description, d.Get("description"))
		assert.Equal(t, path, d.Get("path"))
		assert.Equal(t, int(revision), d.Get("revision"))
		assert.Equal(t, "domain-c1", d.Get("nsx_id"))
		assert.True(t, d.Get("ids_enabled").(bool))

		clusterList := d.Get("cluster").([]interface{})
		require.Len(t, clusterList, 1)
		clusterMap := clusterList[0].(map[string]interface{})
		assert.Equal(t, targetID, clusterMap["target_id"])
		assert.Equal(t, targetType, clusterMap["target_type"])
	})

	t.Run("Read API error", func(t *testing.T) {
		mockClient.EXPECT().
			Get("domain-c1", gomock.Any()).
			Return(model.IdsClusterConfig{}, errors.New("API error"))

		res := resourceNsxtPolicyIdpsClusterConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("domain-c1")

		err := resourceNsxtPolicyIdpsClusterConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIdpsClusterConfigUpdate(t *testing.T) {
	mockClient := setupIdpsClusterConfigMock(t)

	t.Run("Update success", func(t *testing.T) {
		displayName := "idps-cfg-updated"
		targetID := "domain-c1"
		targetType := "VC_Cluster"
		idsEnabled := false

		mockClient.EXPECT().
			Patch("domain-c1", gomock.Any()).
			Return(nil)
		mockClient.EXPECT().
			Get("domain-c1", gomock.Any()).
			Return(model.IdsClusterConfig{
				DisplayName: &displayName,
				IdsEnabled:  &idsEnabled,
				Cluster: &model.PolicyResourceReference{
					TargetId:   &targetID,
					TargetType: &targetType,
				},
			}, nil)

		res := resourceNsxtPolicyIdpsClusterConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"ids_enabled": false,
			"cluster": []interface{}{
				map[string]interface{}{
					"target_id":   "domain-c1",
					"target_type": "VC_Cluster",
				},
			},
		})
		d.SetId("domain-c1")

		err := resourceNsxtPolicyIdpsClusterConfigUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, displayName, d.Get("display_name"))
		assert.False(t, d.Get("ids_enabled").(bool))
	})

	t.Run("Update fails when Patch API errors", func(t *testing.T) {
		mockClient.EXPECT().
			Patch("domain-c1", gomock.Any()).
			Return(errors.New("API error"))

		res := resourceNsxtPolicyIdpsClusterConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsClusterConfigData())
		d.SetId("domain-c1")

		err := resourceNsxtPolicyIdpsClusterConfigUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIdpsClusterConfigDelete(t *testing.T) {
	mockClient := setupIdpsClusterConfigMock(t)

	t.Run("Delete disables IDPS via Patch", func(t *testing.T) {
		mockClient.EXPECT().
			Patch("domain-c1", gomock.Any()).
			DoAndReturn(func(_ string, config model.IdsClusterConfig) error {
				require.NotNil(t, config.IdsEnabled)
				assert.False(t, *config.IdsEnabled)
				return nil
			})

		res := resourceNsxtPolicyIdpsClusterConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsClusterConfigData())
		d.SetId("domain-c1")

		err := resourceNsxtPolicyIdpsClusterConfigDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when Patch API errors", func(t *testing.T) {
		mockClient.EXPECT().
			Patch("domain-c1", gomock.Any()).
			Return(errors.New("API error"))

		res := resourceNsxtPolicyIdpsClusterConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdpsClusterConfigData())
		d.SetId("domain-c1")

		err := resourceNsxtPolicyIdpsClusterConfigDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
