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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	tier0mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
	tier1mocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_1s"
)

var (
	tier0GatewayID  = "tier0-gw-1"
	tier1GatewayID  = "tier1-gw-1"
	tier0ConfigPath = "/infra/tier-0s/tier0-gw-1/security-config"
	tier1ConfigPath = "/infra/tier-1s/tier1-gw-1/security-config"
	configRevision  = int64(1)
)

func tier0SecurityFeaturesModel(idpsEnabled, idfwEnabled bool) model.Tier0SecurityFeatures {
	idpsFeature := model.Tier0SecurityFeature_FEATURE_IDPS
	idfwFeature := model.Tier0SecurityFeature_FEATURE_IDFW

	return model.Tier0SecurityFeatures{
		Path:     &tier0ConfigPath,
		Revision: &configRevision,
		Features: []model.Tier0SecurityFeature{
			{
				Feature: &idpsFeature,
				Enable:  &idpsEnabled,
			},
			{
				Feature: &idfwFeature,
				Enable:  &idfwEnabled,
			},
		},
	}
}

func tier1SecurityFeaturesModel(idpsEnabled, idfwEnabled, malwareEnabled, tlsEnabled bool) model.SecurityFeatures {
	idpsFeature := model.SecurityFeature_FEATURE_IDPS
	idfwFeature := model.SecurityFeature_FEATURE_IDFW
	malwareFeature := model.SecurityFeature_FEATURE_MALWAREPREVENTION
	tlsFeature := model.SecurityFeature_FEATURE_TLS

	return model.SecurityFeatures{
		Path:     &tier1ConfigPath,
		Revision: &configRevision,
		Features: []model.SecurityFeature{
			{
				Feature: &idpsFeature,
				Enable:  &idpsEnabled,
			},
			{
				Feature: &idfwFeature,
				Enable:  &idfwEnabled,
			},
			{
				Feature: &malwareFeature,
				Enable:  &malwareEnabled,
			},
			{
				Feature: &tlsFeature,
				Enable:  &tlsEnabled,
			},
		},
	}
}

func setupTier0SecurityConfigDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*tier0mocks.MockSecurityConfigClient, func()) {
	t.Helper()
	mockSDK := tier0mocks.NewMockSecurityConfigClient(ctrl)
	orig := cliTier0SecurityConfigClientDS
	cliTier0SecurityConfigClientDS = func(_ vapiProtocolClient.Connector) tier_0s.SecurityConfigClient {
		return mockSDK
	}
	return mockSDK, func() { cliTier0SecurityConfigClientDS = orig }
}

func setupTier1SecurityConfigDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*tier1mocks.MockSecurityConfigClient, func()) {
	t.Helper()
	mockSDK := tier1mocks.NewMockSecurityConfigClient(ctrl)
	orig := cliTier1SecurityConfigClientDS
	cliTier1SecurityConfigClientDS = func(_ vapiProtocolClient.Connector) tier_1s.SecurityConfigClient {
		return mockSDK
	}
	return mockSDK, func() { cliTier1SecurityConfigClientDS = orig }
}

func TestMockDataSourceNsxtPolicyGatewaySecurityConfigRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Tier0 gateway success with features enabled", func(t *testing.T) {
		mockSDK, restore := setupTier0SecurityConfigDataSourceMock(t, ctrl)
		defer restore()

		expectedConfig := tier0SecurityFeaturesModel(true, true)
		mockSDK.EXPECT().Get(tier0GatewayID, nil, nil, nil, nil, nil, nil).Return(expectedConfig, nil)

		ds := dataSourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"tier0_id": tier0GatewayID,
		})

		err := dataSourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "tier0/"+tier0GatewayID, d.Id())
		assert.Equal(t, true, d.Get("idps_enabled"))
		assert.Equal(t, true, d.Get("idfw_enabled"))
		assert.Equal(t, false, d.Get("malware_prevention_enabled")) // Not supported on Tier0
		assert.Equal(t, false, d.Get("tls_enabled"))                // Not supported on Tier0
		assert.Equal(t, tier0ConfigPath, d.Get("path"))
	})

	t.Run("Tier0 gateway success with features disabled", func(t *testing.T) {
		mockSDK, restore := setupTier0SecurityConfigDataSourceMock(t, ctrl)
		defer restore()

		expectedConfig := tier0SecurityFeaturesModel(false, false)
		mockSDK.EXPECT().Get(tier0GatewayID, nil, nil, nil, nil, nil, nil).Return(expectedConfig, nil)

		ds := dataSourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"tier0_id": tier0GatewayID,
		})

		err := dataSourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "tier0/"+tier0GatewayID, d.Id())
		assert.Equal(t, false, d.Get("idps_enabled"))
		assert.Equal(t, false, d.Get("idfw_enabled"))
		assert.Equal(t, false, d.Get("malware_prevention_enabled"))
		assert.Equal(t, false, d.Get("tls_enabled"))
	})

	t.Run("Tier0 gateway not found", func(t *testing.T) {
		mockSDK, restore := setupTier0SecurityConfigDataSourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(tier0GatewayID, nil, nil, nil, nil, nil, nil).Return(model.Tier0SecurityFeatures{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"tier0_id": tier0GatewayID,
		})

		err := dataSourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error reading Tier0 Gateway Security Config")
	})

	t.Run("Tier0 gateway API error", func(t *testing.T) {
		mockSDK, restore := setupTier0SecurityConfigDataSourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(tier0GatewayID, nil, nil, nil, nil, nil, nil).Return(model.Tier0SecurityFeatures{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"tier0_id": tier0GatewayID,
		})

		err := dataSourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error reading Tier0 Gateway Security Config")
	})

	t.Run("Tier1 gateway success with all features enabled", func(t *testing.T) {
		mockSDK, restore := setupTier1SecurityConfigDataSourceMock(t, ctrl)
		defer restore()

		expectedConfig := tier1SecurityFeaturesModel(true, true, true, true)
		mockSDK.EXPECT().Get(tier1GatewayID, nil, nil, nil, nil, nil, nil).Return(expectedConfig, nil)

		ds := dataSourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"tier1_id": tier1GatewayID,
		})

		err := dataSourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "tier1/"+tier1GatewayID, d.Id())
		assert.Equal(t, true, d.Get("idps_enabled"))
		assert.Equal(t, true, d.Get("idfw_enabled"))
		assert.Equal(t, true, d.Get("malware_prevention_enabled"))
		assert.Equal(t, true, d.Get("tls_enabled"))
		assert.Equal(t, tier1ConfigPath, d.Get("path"))
	})

	t.Run("Tier1 gateway success with mixed feature settings", func(t *testing.T) {
		mockSDK, restore := setupTier1SecurityConfigDataSourceMock(t, ctrl)
		defer restore()

		expectedConfig := tier1SecurityFeaturesModel(true, false, true, false)
		mockSDK.EXPECT().Get(tier1GatewayID, nil, nil, nil, nil, nil, nil).Return(expectedConfig, nil)

		ds := dataSourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"tier1_id": tier1GatewayID,
		})

		err := dataSourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "tier1/"+tier1GatewayID, d.Id())
		assert.Equal(t, true, d.Get("idps_enabled"))
		assert.Equal(t, false, d.Get("idfw_enabled"))
		assert.Equal(t, true, d.Get("malware_prevention_enabled"))
		assert.Equal(t, false, d.Get("tls_enabled"))
	})

	t.Run("Tier1 gateway success with all features disabled", func(t *testing.T) {
		mockSDK, restore := setupTier1SecurityConfigDataSourceMock(t, ctrl)
		defer restore()

		expectedConfig := tier1SecurityFeaturesModel(false, false, false, false)
		mockSDK.EXPECT().Get(tier1GatewayID, nil, nil, nil, nil, nil, nil).Return(expectedConfig, nil)

		ds := dataSourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"tier1_id": tier1GatewayID,
		})

		err := dataSourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "tier1/"+tier1GatewayID, d.Id())
		assert.Equal(t, false, d.Get("idps_enabled"))
		assert.Equal(t, false, d.Get("idfw_enabled"))
		assert.Equal(t, false, d.Get("malware_prevention_enabled"))
		assert.Equal(t, false, d.Get("tls_enabled"))
	})

	t.Run("Tier1 gateway not found", func(t *testing.T) {
		mockSDK, restore := setupTier1SecurityConfigDataSourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(tier1GatewayID, nil, nil, nil, nil, nil, nil).Return(model.SecurityFeatures{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"tier1_id": tier1GatewayID,
		})

		err := dataSourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error reading Tier1 Gateway Security Config")
	})

	t.Run("Tier1 gateway API error", func(t *testing.T) {
		mockSDK, restore := setupTier1SecurityConfigDataSourceMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(tier1GatewayID, nil, nil, nil, nil, nil, nil).Return(model.SecurityFeatures{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"tier1_id": tier1GatewayID,
		})

		err := dataSourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error reading Tier1 Gateway Security Config")
	})

	t.Run("missing gateway id parameters", func(t *testing.T) {
		ds := dataSourceNsxtPolicyGatewaySecurityConfig()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyGatewaySecurityConfigRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "one of tier0_id or tier1_id must be specified")
	})
}
