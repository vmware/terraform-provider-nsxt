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

func TestUnitNsxt_getPolicyPathResourceImporter(t *testing.T) {
	m := newGoMockProviderClient()
	imp := getPolicyPathResourceImporter("/infra/domains/default/groups/example")

	t.Run("valid policy path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId("/infra/domains/default/groups/mygroup")

		rd, err := imp(d, m)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, "mygroup", rd[0].Id())
	})

	t.Run("invalid path format", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId("not-a-valid-uri")

		_, err := imp(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid policy path")
	})
}

func TestUnitNsxt_getPolicyPathOrIDResourceImporter(t *testing.T) {
	m := newGoMockProviderClient()
	imp := getPolicyPathOrIDResourceImporter("/infra/foo/bar")

	t.Run("valid id without path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId("plain-resource-id")

		rd, err := imp(d, m)
		require.NoError(t, err)
		require.Len(t, rd, 1)
	})

	t.Run("invalid id when not a path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		d.SetId("bad/id")

		_, err := imp(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid policy path or id")
	})
}

func TestUnitNsxt_getVpcPathResourceImporter(t *testing.T) {
	m := newGoMockProviderClient()
	imp := getVpcPathResourceImporter("/orgs/o/projects/p/vpcs/v/infra/example")

	t.Run("valid vpc path sets context", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"context": getContextSchema(false, false, true),
		}, map[string]interface{}{})
		d.SetId("/orgs/myorg/projects/myproj/vpcs/myvpc/infra/tier-1s/t1")

		rd, err := imp(d, m)
		require.NoError(t, err)
		require.Len(t, rd, 1)
	})

	t.Run("missing vpc in path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"context": getContextSchema(false, false, true),
		}, map[string]interface{}{})
		d.SetId("/infra/domains/default/groups/g1")

		_, err := imp(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "project_id and vpc_id")
	})
}
