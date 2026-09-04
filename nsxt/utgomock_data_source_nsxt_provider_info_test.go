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
)

func TestUnitNsxt_dataSourceNsxtProviderInfoRead(t *testing.T) {
	origCommit := GitCommit
	GitCommit = "abc123"
	defer func() { GitCommit = origCommit }()

	ds := dataSourceNsxtProviderInfo()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

	err := dataSourceNsxtProviderInfoRead(d, nil)
	require.NoError(t, err)

	assert.Equal(t, "nsxt", d.Id())
	assert.Equal(t, "abc123", d.Get("commit"))
	assert.NotEmpty(t, d.Get("date"))
}
