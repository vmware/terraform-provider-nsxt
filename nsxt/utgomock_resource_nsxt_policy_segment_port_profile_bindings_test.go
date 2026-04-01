//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func minimalSegmentPortProfileBindingsData() map[string]interface{} {
	return map[string]interface{}{
		"segment_port_path": "invalid-port-path",
	}
}

func TestMockResourceNsxtPolicySegmentPortProfileBindingsCreate(t *testing.T) {
	t.Run("Create fails with invalid segment_port_path", func(t *testing.T) {
		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSegmentPortProfileBindingsData())

		err := resourceNsxtPolicySegmentPortProfileBindingsCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySegmentPortProfileBindingsRead(t *testing.T) {
	t.Run("Read fails with invalid segment_port_path", func(t *testing.T) {
		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSegmentPortProfileBindingsData())
		d.SetId("port-1")

		err := resourceNsxtPolicySegmentPortProfileBindingsRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySegmentPortProfileBindingsUpdate(t *testing.T) {
	t.Run("Update fails with invalid segment_port_path", func(t *testing.T) {
		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSegmentPortProfileBindingsData())
		d.SetId("port-1")

		err := resourceNsxtPolicySegmentPortProfileBindingsUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySegmentPortProfileBindingsDelete(t *testing.T) {
	t.Run("Delete fails with invalid segment_port_path", func(t *testing.T) {
		res := resourceNsxtPolicySegmentPortProfileBindings()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSegmentPortProfileBindingsData())
		d.SetId("port-1")

		err := resourceNsxtPolicySegmentPortProfileBindingsDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
