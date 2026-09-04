//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// dataSourceNsxtPolicyIdpsCustomSignatureRead uses
// custom_signature_versions.NewCustomSignaturesClient(connector) directly
// rather than an injectable wrapper variable, so full mock-based Read tests
// are not possible. These tests cover the guard/parsing conditions that run
// before any client call.

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
