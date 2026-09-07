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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"
)

// dataSourceNsxtPolicyIdpsSettingsRead reads through the cliIdsSettingsClient and
// cliIdsCustomSigSettingsClient package-level vars (see
// utgomock_resource_nsxt_policy_idps_settings_test.go for the mock setup helper), so both the
// guard condition and the mocked lookups are covered here.

func TestMockDataSourceNsxtPolicyIdpsSettingsReadGuard(t *testing.T) {
	t.Run("Read fails for global manager", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIdpsSettings()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsSettingsRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Global Manager")
	})
}

func TestMockDataSourceNsxtPolicyIdpsSettingsReadMocked(t *testing.T) {
	t.Run("Read succeeds without a custom signature version set", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSettings, _, restore := setupIdpsSettingsMocks(ctrl)
		defer restore()
		oversub := "DROPPED"
		mockSettings.EXPECT().Get().Return(model.IdsSettings{Oversubscription: &oversub}, nil)

		ds := dataSourceNsxtPolicyIdpsSettings()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsSettingsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "DROPPED", d.Get("oversubscription"))
		assert.Equal(t, false, d.Get("enable_custom_signatures"))
	})

	t.Run("Read also fetches custom signature settings when a version id is set", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSettings, mockCustomSig, restore := setupIdpsSettingsMocks(ctrl)
		defer restore()
		mockSettings.EXPECT().Get().Return(model.IdsSettings{}, nil)
		enabled := true
		mockCustomSig.EXPECT().Get("default").Return(model.IdsCustomSignatureSettings{EnableCustomSignatures: &enabled}, nil)

		ds := dataSourceNsxtPolicyIdpsSettings()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"custom_signature_version_id": "default"})

		err := dataSourceNsxtPolicyIdpsSettingsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, true, d.Get("enable_custom_signatures"))
	})

	t.Run("Read falls back to false when custom signature settings lookup errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSettings, mockCustomSig, restore := setupIdpsSettingsMocks(ctrl)
		defer restore()
		mockSettings.EXPECT().Get().Return(model.IdsSettings{}, nil)
		mockCustomSig.EXPECT().Get("default").Return(model.IdsCustomSignatureSettings{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIdpsSettings()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"custom_signature_version_id": "default"})

		err := dataSourceNsxtPolicyIdpsSettingsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, false, d.Get("enable_custom_signatures"))
	})

	t.Run("Read propagates the main settings API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSettings, _, restore := setupIdpsSettingsMocks(ctrl)
		defer restore()
		mockSettings.EXPECT().Get().Return(model.IdsSettings{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIdpsSettings()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIdpsSettingsRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
