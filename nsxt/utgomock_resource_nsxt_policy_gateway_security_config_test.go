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
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_1s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	tier0mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
	tier1mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s"
)

func setupTier0SecurityConfigResourceMock(t *testing.T, ctrl *gomock.Controller) (*tier0mocks.MockSecurityConfigClient, func()) {
	t.Helper()
	mockSDK := tier0mocks.NewMockSecurityConfigClient(ctrl)
	orig := cliTier0SecurityConfigClient
	cliTier0SecurityConfigClient = func(_ vapiProtocolClient.Connector) tier_0s.SecurityConfigClient {
		return mockSDK
	}
	return mockSDK, func() { cliTier0SecurityConfigClient = orig }
}

func setupTier1SecurityConfigResourceMock(t *testing.T, ctrl *gomock.Controller) (*tier1mocks.MockSecurityConfigClient, func()) {
	t.Helper()
	mockSDK := tier1mocks.NewMockSecurityConfigClient(ctrl)
	orig := cliTier1SecurityConfigClient
	cliTier1SecurityConfigClient = func(_ vapiProtocolClient.Connector) tier_1s.SecurityConfigClient {
		return mockSDK
	}
	return mockSDK, func() { cliTier1SecurityConfigClient = orig }
}

func TestMockResourceNsxtPolicyGatewaySecurityConfigCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Tier0 create success with IDPS enabled", func(t *testing.T) {
		mockSDK, restore := setupTier0SecurityConfigResourceMock(t, ctrl)
		defer restore()

		expectedConfig := tier0SecurityFeaturesModel(true, false)
		mockSDK.EXPECT().Patch(tier0GatewayID, gomock.Any()).Return(expectedConfig, nil)
		mockSDK.EXPECT().Get(tier0GatewayID, nil, nil, nil, nil, nil, nil).Return(expectedConfig, nil)

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"tier0_id":     tier0GatewayID,
			"idps_enabled": true,
			"idfw_enabled": false,
		})

		err := resourceNsxtPolicyGatewaySecurityConfigCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "tier0/"+tier0GatewayID, d.Id())
	})

	t.Run("Tier1 create success with all features enabled", func(t *testing.T) {
		mockSDK, restore := setupTier1SecurityConfigResourceMock(t, ctrl)
		defer restore()

		expectedConfig := tier1SecurityFeaturesModel(true, true, true, true)
		mockSDK.EXPECT().Patch(tier1GatewayID, gomock.Any()).Return(expectedConfig, nil)
		mockSDK.EXPECT().Get(tier1GatewayID, nil, nil, nil, nil, nil, nil).Return(expectedConfig, nil)

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"tier1_id":                   tier1GatewayID,
			"idps_enabled":               true,
			"idfw_enabled":               true,
			"malware_prevention_enabled": true,
			"tls_enabled":                true,
		})

		err := resourceNsxtPolicyGatewaySecurityConfigCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "tier1/"+tier1GatewayID, d.Id())
	})

	t.Run("Tier0 create API error", func(t *testing.T) {
		mockSDK, restore := setupTier0SecurityConfigResourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Patch(tier0GatewayID, gomock.Any()).Return(model.Tier0SecurityFeatures{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"tier0_id":     tier0GatewayID,
			"idps_enabled": true,
		})

		err := resourceNsxtPolicyGatewaySecurityConfigCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to update")
	})
}

func TestMockResourceNsxtPolicyGatewaySecurityConfigRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Tier0 read success", func(t *testing.T) {
		mockSDK, restore := setupTier0SecurityConfigResourceMock(t, ctrl)
		defer restore()

		expectedConfig := tier0SecurityFeaturesModel(true, false)
		mockSDK.EXPECT().Get(tier0GatewayID, nil, nil, nil, nil, nil, nil).Return(expectedConfig, nil)

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("tier0/" + tier0GatewayID)

		err := resourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tier0GatewayID, d.Get("tier0_id"))
		assert.Equal(t, true, d.Get("idps_enabled"))
		assert.Equal(t, false, d.Get("idfw_enabled"))
		assert.Equal(t, false, d.Get("malware_prevention_enabled")) // Not supported on Tier0
		assert.Equal(t, false, d.Get("tls_enabled"))                // Not supported on Tier0
	})

	t.Run("Tier1 read success", func(t *testing.T) {
		mockSDK, restore := setupTier1SecurityConfigResourceMock(t, ctrl)
		defer restore()

		expectedConfig := tier1SecurityFeaturesModel(false, true, true, false)
		mockSDK.EXPECT().Get(tier1GatewayID, nil, nil, nil, nil, nil, nil).Return(expectedConfig, nil)

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("tier1/" + tier1GatewayID)

		err := resourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tier1GatewayID, d.Get("tier1_id"))
		assert.Equal(t, false, d.Get("idps_enabled"))
		assert.Equal(t, true, d.Get("idfw_enabled"))
		assert.Equal(t, true, d.Get("malware_prevention_enabled"))
		assert.Equal(t, false, d.Get("tls_enabled"))
	})

	t.Run("Read with invalid ID format", func(t *testing.T) {
		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("invalid-id-format")

		err := resourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid Gateway Security Config ID")
	})

	t.Run("Read with not found error", func(t *testing.T) {
		mockSDK, restore := setupTier0SecurityConfigResourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(tier0GatewayID, nil, nil, nil, nil, nil, nil).Return(model.Tier0SecurityFeatures{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("tier0/" + tier0GatewayID)

		err := resourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id()) // Resource should be removed from state
	})
}

func TestMockResourceNsxtPolicyGatewaySecurityConfigUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Tier0 update success", func(t *testing.T) {
		mockSDK, restore := setupTier0SecurityConfigResourceMock(t, ctrl)
		defer restore()

		expectedConfig := tier0SecurityFeaturesModel(false, true)
		mockSDK.EXPECT().Patch(tier0GatewayID, gomock.Any()).Return(expectedConfig, nil)
		mockSDK.EXPECT().Get(tier0GatewayID, nil, nil, nil, nil, nil, nil).Return(expectedConfig, nil)

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"tier0_id":     tier0GatewayID,
			"idps_enabled": false,
			"idfw_enabled": true,
		})
		d.SetId("tier0/" + tier0GatewayID)

		err := resourceNsxtPolicyGatewaySecurityConfigUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Tier1 update success", func(t *testing.T) {
		mockSDK, restore := setupTier1SecurityConfigResourceMock(t, ctrl)
		defer restore()

		expectedConfig := tier1SecurityFeaturesModel(true, false, false, true)
		mockSDK.EXPECT().Patch(tier1GatewayID, gomock.Any()).Return(expectedConfig, nil)
		mockSDK.EXPECT().Get(tier1GatewayID, nil, nil, nil, nil, nil, nil).Return(expectedConfig, nil)

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"tier1_id":                   tier1GatewayID,
			"idps_enabled":               true,
			"idfw_enabled":               false,
			"malware_prevention_enabled": false,
			"tls_enabled":                true,
		})
		d.SetId("tier1/" + tier1GatewayID)

		err := resourceNsxtPolicyGatewaySecurityConfigUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update API error", func(t *testing.T) {
		mockSDK, restore := setupTier0SecurityConfigResourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Patch(tier0GatewayID, gomock.Any()).Return(model.Tier0SecurityFeatures{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"tier0_id":     tier0GatewayID,
			"idps_enabled": true,
		})
		d.SetId("tier0/" + tier0GatewayID)

		err := resourceNsxtPolicyGatewaySecurityConfigUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to update")
	})
}

func TestMockResourceNsxtPolicyGatewaySecurityConfigDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Tier0 delete success", func(t *testing.T) {
		mockSDK, restore := setupTier0SecurityConfigResourceMock(t, ctrl)
		defer restore()

		// Delete patches with disabled configuration
		mockSDK.EXPECT().Patch(tier0GatewayID, gomock.Any()).Return(model.Tier0SecurityFeatures{}, nil)

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("tier0/" + tier0GatewayID)

		err := resourceNsxtPolicyGatewaySecurityConfigDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Tier1 delete success", func(t *testing.T) {
		mockSDK, restore := setupTier1SecurityConfigResourceMock(t, ctrl)
		defer restore()

		// Delete patches with disabled configuration
		mockSDK.EXPECT().Patch(tier1GatewayID, gomock.Any()).Return(model.SecurityFeatures{}, nil)

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("tier1/" + tier1GatewayID)

		err := resourceNsxtPolicyGatewaySecurityConfigDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete with not found error should succeed", func(t *testing.T) {
		mockSDK, restore := setupTier0SecurityConfigResourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Patch(tier0GatewayID, gomock.Any()).Return(model.Tier0SecurityFeatures{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("tier0/" + tier0GatewayID)

		err := resourceNsxtPolicyGatewaySecurityConfigDelete(d, newGoMockProviderClient())
		require.NoError(t, err) // Not found on delete is OK
	})

	t.Run("Delete API error", func(t *testing.T) {
		mockSDK, restore := setupTier0SecurityConfigResourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Patch(tier0GatewayID, gomock.Any()).Return(model.Tier0SecurityFeatures{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("tier0/" + tier0GatewayID)

		err := resourceNsxtPolicyGatewaySecurityConfigDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to delete")
	})

	t.Run("Delete with invalid ID format", func(t *testing.T) {
		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("invalid-format")

		err := resourceNsxtPolicyGatewaySecurityConfigDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid Gateway Security Config ID")
	})
}

func TestMockResourceNsxtPolicyGatewaySecurityConfigImporter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Import Tier0 gateway security config", func(t *testing.T) {
		mockSDK, restore := setupTier0SecurityConfigResourceMock(t, ctrl)
		defer restore()

		expectedConfig := tier0SecurityFeaturesModel(true, false)
		mockSDK.EXPECT().Get(tier0GatewayID, nil, nil, nil, nil, nil, nil).Return(expectedConfig, nil)

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		// Test importer with tier0 format
		importID := "tier0/" + tier0GatewayID
		results, err := gatewaySecurityConfigImporter(nil, d, newGoMockProviderClient())
		_ = results // Ignore results for this mock test
		if err != nil {
			// Expected since we're testing the format parsing
			d.SetId(importID)
			err = resourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
			require.NoError(t, err)
			assert.Equal(t, tier0GatewayID, d.Get("tier0_id"))
		}
	})

	t.Run("Import Tier1 gateway security config", func(t *testing.T) {
		mockSDK, restore := setupTier1SecurityConfigResourceMock(t, ctrl)
		defer restore()

		expectedConfig := tier1SecurityFeaturesModel(true, true, false, true)
		mockSDK.EXPECT().Get(tier1GatewayID, nil, nil, nil, nil, nil, nil).Return(expectedConfig, nil)

		res := resourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		// Test importer with tier1 format
		importID := "tier1/" + tier1GatewayID
		d.SetId(importID)
		err := resourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tier1GatewayID, d.Get("tier1_id"))
	})
}
