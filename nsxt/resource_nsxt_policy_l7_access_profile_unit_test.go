//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func l7AccessEntryMap(nsxID, displayName string) map[string]interface{} {
	return map[string]interface{}{
		"nsx_id":          nsxID,
		"display_name":    displayName,
		"description":     "an entry",
		"path":            "",
		"revision":        0,
		"action":          model.L7AccessEntry_ACTION_ALLOW,
		"disabled":        false,
		"logged":          false,
		"sequence_number": 1,
		"attribute": []interface{}{
			map[string]interface{}{
				"attribute_source":         "SYSTEM",
				"custom_url_partial_match": false,
				"description":              "attr desc",
				"is_alg_type":              false,
				"key":                      model.L7AccessAttributes_KEY_APP_ID,
				"metadata": []interface{}{
					map[string]interface{}{"key": "meta-key", "value": "meta-value"},
				},
				"sub_attribute": []interface{}{
					map[string]interface{}{"key": "sub-key", "values": []interface{}{"v1", "v2"}},
				},
				"values": []interface{}{"HTTP"},
			},
		},
	}
}

func TestUnitNsxt_childL7AccessEntryFromSchema(t *testing.T) {
	t.Run("converts a fully populated entry", func(t *testing.T) {
		sv, err := childL7AccessEntryFromSchema("entry-1", l7AccessEntryMap("entry-1", "my-entry"), false)
		require.NoError(t, err)
		require.NotNil(t, sv)

		converter := bindings.NewTypeConverter()
		golang, errs := converter.ConvertToGolang(sv, model.ChildL7AccessEntryBindingType())
		require.Empty(t, errs)
		child := golang.(model.ChildL7AccessEntry)
		assert.Equal(t, "entry-1", *child.Id)
		assert.False(t, *child.MarkedForDelete)
		require.NotNil(t, child.L7AccessEntry)
		assert.Equal(t, "my-entry", *child.L7AccessEntry.DisplayName)
		require.Len(t, child.L7AccessEntry.Attributes, 1)
		assert.Equal(t, model.L7AccessAttributes_KEY_APP_ID, *child.L7AccessEntry.Attributes[0].Key)
		require.Len(t, child.L7AccessEntry.Attributes[0].SubAttributes, 1)
		assert.Equal(t, "sub-key", *child.L7AccessEntry.Attributes[0].SubAttributes[0].Key)
	})

	t.Run("marks the entry for deletion", func(t *testing.T) {
		sv, err := childL7AccessEntryFromSchema("entry-1", l7AccessEntryMap("entry-1", "my-entry"), true)
		require.NoError(t, err)

		converter := bindings.NewTypeConverter()
		golang, errs := converter.ConvertToGolang(sv, model.ChildL7AccessEntryBindingType())
		require.Empty(t, errs)
		child := golang.(model.ChildL7AccessEntry)
		assert.True(t, *child.MarkedForDelete)
	})
}

func TestUnitNsxt_populateL7AccessProfileStructChildren(t *testing.T) {
	t.Run("new entries are added and stale entries are marked for deletion", func(t *testing.T) {
		oldEntries := []interface{}{
			l7AccessEntryMap("entry-old", "old-entry"),
		}
		newEntries := []interface{}{
			l7AccessEntryMap("entry-new", "new-entry"),
		}

		children, err := populateL7AccessProfileStructChildren(oldEntries, newEntries)
		require.NoError(t, err)
		require.Len(t, children, 2, "one new entry plus one stale entry marked for deletion")

		converter := bindings.NewTypeConverter()

		golangNew, errs := converter.ConvertToGolang(children[0], model.ChildL7AccessEntryBindingType())
		require.Empty(t, errs)
		newChild := golangNew.(model.ChildL7AccessEntry)
		assert.Equal(t, "entry-new", *newChild.Id)
		assert.False(t, *newChild.MarkedForDelete)

		golangOld, errs := converter.ConvertToGolang(children[1], model.ChildL7AccessEntryBindingType())
		require.Empty(t, errs)
		oldChild := golangOld.(model.ChildL7AccessEntry)
		assert.Equal(t, "entry-old", *oldChild.Id)
		assert.True(t, *oldChild.MarkedForDelete)
	})

	t.Run("entries that persist are not marked for deletion", func(t *testing.T) {
		entries := []interface{}{l7AccessEntryMap("entry-1", "entry-1")}
		children, err := populateL7AccessProfileStructChildren(entries, entries)
		require.NoError(t, err)
		require.Len(t, children, 1)
	})

	t.Run("empty old and new returns empty list", func(t *testing.T) {
		children, err := populateL7AccessProfileStructChildren([]interface{}{}, []interface{}{})
		require.NoError(t, err)
		assert.Empty(t, children)
	})
}
