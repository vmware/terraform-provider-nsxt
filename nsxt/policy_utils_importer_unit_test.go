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

func domainImporterTestResource() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"domain": {Type: schema.TypeString, Optional: true},
	}}
}

func TestUnitNsxt_nsxtDomainResourceImporter(t *testing.T) {
	res := domainImporterTestResource()

	t.Run("empty ID errors", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("")
		_, err := nsxtDomainResourceImporter(d, nil)
		require.Error(t, err)
		assert.Equal(t, ErrEmptyImportID, err)
	})

	t.Run("plain ID falls back to default domain", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("mygroup")
		out, err := nsxtDomainResourceImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, "mygroup", d.Id())
		assert.Equal(t, "default", d.Get("domain"))
	})

	t.Run("domain/id shorthand sets both domain and id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("mydomain/mygroup")
		out, err := nsxtDomainResourceImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, "mygroup", d.Id())
		assert.Equal(t, "mydomain", d.Get("domain"))
	})

	t.Run("full policy path extracts domain from path segments", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/infra/domains/mydomain/groups/mygroup")
		out, err := nsxtDomainResourceImporter(d, nil)
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, "mygroup", d.Id())
		assert.Equal(t, "mydomain", d.Get("domain"))
	})
}
