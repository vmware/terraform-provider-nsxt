//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitNsxt_addTrailingSlash(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", "/"},
		{"already slash", "/foo/bar/", "/foo/bar/"},
		{"no slash", "/foo/bar", "/foo/bar/"},
		{"single segment", "path", "path/"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, addTrailingSlash(tc.in))
		})
	}
}

func TestUnitNsxt_pathPrefixDiffSuppress(t *testing.T) {
	cases := []struct {
		name   string
		oldVal string
		newVal string
		want   bool
	}{
		{"equal with slash mismatch on old", "/a/b/", "/a/b", true},
		{"equal with slash mismatch on new", "/a/b", "/a/b/", true},
		{"both normalized equal", "/a/b/", "/a/b/", true},
		{"different paths", "/a/b", "/a/c", false},
		{"empty vs slash", "", "/", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := pathPrefixDiffSupress("path_prefix", tc.oldVal, tc.newVal, &schema.ResourceData{})
			assert.Equal(t, tc.want, got)
		})
	}
}
