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

func minimalLDAPData() map[string]interface{} {
	return map[string]interface{}{
		"nsx_id":      "ldap-src-1",
		"description": "Test LDAP Source",
		"type":        activeDirectoryType,
		"domain_name": "corp.example.com",
		"base_dn":     "dc=corp,dc=example,dc=com",
	}
}

func TestMockResourceNsxtPolicyLdapIdentitySourceRead(t *testing.T) {
	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLDAPData())

		err := resourceNsxtPolicyLdapIdentitySourceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLdapIdentitySourceUpdate(t *testing.T) {
	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLDAPData())

		err := resourceNsxtPolicyLdapIdentitySourceUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLdapIdentitySourceDelete(t *testing.T) {
	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLdapIdentitySource()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLDAPData())

		err := resourceNsxtPolicyLdapIdentitySourceDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
