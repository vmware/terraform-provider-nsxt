package nsxt

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestBuildTagQuery(t *testing.T) {
	tagSchema := map[string]*schema.Schema{
		"tag": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"scope": {Type: schema.TypeString, Optional: true},
				"tag":   {Type: schema.TypeString, Optional: true},
			}},
		},
	}

	assertContainsAll := func(t *testing.T, s string, parts ...string) {
		t.Helper()
		for _, p := range parts {
			if !strings.Contains(s, p) {
				t.Fatalf("expected query to contain %q, got %q", p, s)
			}
		}
	}

	expectScope := func(scope string) string {
		return "tags.scope:" + escapeSpecialCharacters(scope)
	}
	expectTag := func(tag string) string {
		return "tags.tag:" + escapeSpecialCharacters(tag)
	}

	runID := "run-123"

	// 1) Tag-mode (default): if resource has no tag attribute, we still send provider-managed tags.
	t.Run("tag-mode/no-user-tags", func(t *testing.T) {
		old := os.Getenv("USE_GLOBAL_SEARCH_CACHE")
		t.Cleanup(func() { _ = os.Setenv("USE_GLOBAL_SEARCH_CACHE", old) })
		_ = os.Setenv("USE_GLOBAL_SEARCH_CACHE", "false")

		d := schema.TestResourceDataRaw(t, tagSchema, map[string]interface{}{})
		q := buildTagQuery(d, runID)

		assertContainsAll(t, q,
			expectScope("nsx-tf/managed-by"),
			expectTag("terraform"),
			expectScope("tf-run-id"),
			expectTag(runID),
		)
	})

	// 2) Tag-mode: if user tags are present and managed tags are not, provider-managed tags are appended.
	t.Run("tag-mode/user-tags-appends-provider-tags", func(t *testing.T) {
		old := os.Getenv("USE_GLOBAL_SEARCH_CACHE")
		t.Cleanup(func() { _ = os.Setenv("USE_GLOBAL_SEARCH_CACHE", old) })
		_ = os.Setenv("USE_GLOBAL_SEARCH_CACHE", "false")

		d := schema.TestResourceDataRaw(t, tagSchema, map[string]interface{}{
			"tag": []interface{}{
				map[string]interface{}{"scope": "env", "tag": "dev"},
			},
		})
		q := buildTagQuery(d, runID)

		assertContainsAll(t, q,
			expectScope("env"),
			expectTag("dev"),
			expectScope("nsx-tf/managed-by"),
			expectTag("terraform"),
			expectScope("tf-run-id"),
			expectTag(runID),
		)
	})

	// 3) Global-search mode: no tag attribute means no additional query.
	t.Run("global-search/no-user-tags", func(t *testing.T) {
		old := os.Getenv("USE_GLOBAL_SEARCH_CACHE")
		t.Cleanup(func() { _ = os.Setenv("USE_GLOBAL_SEARCH_CACHE", old) })
		_ = os.Setenv("USE_GLOBAL_SEARCH_CACHE", "true")

		d := schema.TestResourceDataRaw(t, tagSchema, map[string]interface{}{})
		q := buildTagQuery(d, runID)
		if q != "" {
			t.Fatalf("expected empty query in global-search mode when no tags exist, got %q", q)
		}
	})
}

func TestEnsureProviderManagedTagsWithPatchFunc(t *testing.T) {
	type testTagObj struct {
		Tags []model.Tag
	}

	findTag := func(tags []model.Tag, scope string) (string, bool) {
		for _, tg := range tags {
			if tg.Scope != nil && *tg.Scope == scope {
				if tg.Tag == nil {
					return "", true
				}
				return *tg.Tag, true
			}
		}
		return "", false
	}

	newString := func(s string) *string { return &s }

	runID := "run-abc"
	m := nsxtClients{CommonConfig: commonProviderConfig{contextID: runID}}

	t.Run("no-change-when-tags-match", func(t *testing.T) {
		obj := &testTagObj{Tags: []model.Tag{
			{Scope: newString("nsx-tf/managed-by"), Tag: newString("terraform")},
			{Scope: newString("tf-run-id"), Tag: newString(runID)},
			{Scope: newString("env"), Tag: newString("dev")},
		}}
		called := false

		patched, err := ensureProviderManagedTagsWithPatchFunc(obj, nil, m, func(o *testTagObj) error {
			called = true
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if patched != nil {
			t.Fatalf("expected no patch, got patched object")
		}
		if called {
			t.Fatalf("expected patchFunc not to be called")
		}
	})

	t.Run("patches-when-runid-mismatch", func(t *testing.T) {
		obj := &testTagObj{Tags: []model.Tag{
			{Scope: newString("nsx-tf/managed-by"), Tag: newString("terraform")},
			{Scope: newString("tf-run-id"), Tag: newString("old-run")},
			{Scope: newString("env"), Tag: newString("dev")},
		}}
		called := false

		patched, err := ensureProviderManagedTagsWithPatchFunc(obj, nil, m, func(o *testTagObj) error {
			called = true
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if patched == nil {
			t.Fatalf("expected patched object, got nil")
		}
		if !called {
			t.Fatalf("expected patchFunc to be called")
		}

		val, ok := findTag(obj.Tags, "tf-run-id")
		if !ok {
			t.Fatalf("expected tf-run-id tag to exist")
		}
		if val != runID {
			t.Fatalf("expected tf-run-id tag to be %q, got %q", runID, val)
		}
	})
}
