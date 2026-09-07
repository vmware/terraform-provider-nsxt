//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// dataSourceNsxtPolicyIdpsCustomSignatureRead reads through the cliIdsCustomSignaturesClient
// package-level var (see utgomock_resource_nsxt_policy_idps_custom_signature_test.go for the
// mock setup helper), so both the guard/parsing conditions and mocked lookups are covered here.

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

func TestMockDataSourceNsxtPolicyIdpsCustomSignatureRead(t *testing.T) {
	ds := dataSourceNsxtPolicyIdpsCustomSignature()

	t.Run("Read fails on global manager", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": "default/5000001",
		})
		c := newGoMockProviderClient()
		c.PolicyGlobalManager = true

		err := dataSourceNsxtPolicyIdpsCustomSignatureRead(d, c)
		require.Error(t, err)
	})

	t.Run("Read fails when bare signature_id has no signature_version_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": "5000001",
		})

		err := dataSourceNsxtPolicyIdpsCustomSignatureRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "signature_version_id must be set")
	})
}

func TestMockDataSourceNsxtPolicyIdpsCustomSignatureReadMocked(t *testing.T) {
	ds := dataSourceNsxtPolicyIdpsCustomSignature()

	t.Run("Read succeeds via composite ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		name := "sig-name"
		mockSigs.EXPECT().Get("default", "5000001").Return(model.IdsCustomSignature{DisplayName: &name}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "default/5000001"})

		err := dataSourceNsxtPolicyIdpsCustomSignatureRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "default/5000001", d.Id())
		assert.Equal(t, "sig-name", d.Get("display_name"))
	})

	t.Run("Read succeeds via signature_version_id plus bare signature_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		mockSigs.EXPECT().Get("default", "5000001").Return(model.IdsCustomSignature{}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":                   "5000001",
			"signature_version_id": "default",
		})

		err := dataSourceNsxtPolicyIdpsCustomSignatureRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "default/5000001", d.Id())
	})

	t.Run("Read falls back to list search when not found by Get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		sigID := "5000001"
		mockSigs.EXPECT().Get("default", "5000001").Return(model.IdsCustomSignature{}, vapiErrors.NotFound{})
		mockSigs.EXPECT().List("default", gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.IdsCustomSignatureListResult{Results: []model.IdsCustomSignature{{Id: &sigID}}}, nil,
		).AnyTimes()

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "default/5000001"})

		err := dataSourceNsxtPolicyIdpsCustomSignatureRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "default/5000001", d.Id())
	})

	t.Run("Read fails when neither Get nor list search find the signature", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		mockSigs.EXPECT().Get("default", "5000001").Return(model.IdsCustomSignature{}, vapiErrors.NotFound{})
		mockSigs.EXPECT().List("default", gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.IdsCustomSignatureListResult{}, nil,
		).AnyTimes()

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "default/5000001"})

		err := dataSourceNsxtPolicyIdpsCustomSignatureRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read propagates non-NotFound API errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSigs, _, restore := setupIdpsCustomSignatureMocks(ctrl)
		defer restore()
		mockSigs.EXPECT().Get("default", "5000001").Return(model.IdsCustomSignature{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "default/5000001"})

		err := dataSourceNsxtPolicyIdpsCustomSignatureRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
