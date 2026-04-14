package nsxt

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestProviderManagedTagsSearchQuery(t *testing.T) {
	if got := providerManagedTagsSearchQuery(""); got != "" {
		t.Fatalf("empty runID: want empty, got %q", got)
	}
	runID := "run-xyz"
	got := providerManagedTagsSearchQuery(runID)
	for _, part := range []string{
		"tags.scope:" + escapeSpecialCharacters("nsx-tf/tf-run-id"),
		"tags.tag:" + escapeSpecialCharacters(runID),
	} {
		if !strings.Contains(got, part) {
			t.Fatalf("expected query to contain %q, got %q", part, got)
		}
	}
}

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
		old := os.Getenv("NSXT_CACHE_MODE")
		t.Cleanup(func() { _ = os.Setenv("NSXT_CACHE_MODE", old) })
		_ = os.Setenv("NSXT_CACHE_MODE", "config-scope")

		d := schema.TestResourceDataRaw(t, tagSchema, map[string]interface{}{})
		q := buildTagQuery(d, runID)

		assertContainsAll(t, q,
			expectScope("nsx-tf/tf-run-id"),
			expectTag(runID),
		)
	})

	// 2) Tag-mode: if user tags are present and managed tags are not, provider-managed tags are appended.
	t.Run("tag-mode/user-tags-appends-provider-tags", func(t *testing.T) {
		old := os.Getenv("NSXT_CACHE_MODE")
		t.Cleanup(func() { _ = os.Setenv("NSXT_CACHE_MODE", old) })
		_ = os.Setenv("NSXT_CACHE_MODE", "config-scope")

		d := schema.TestResourceDataRaw(t, tagSchema, map[string]interface{}{
			"tag": []interface{}{
				map[string]interface{}{"scope": "env", "tag": "dev"},
			},
		})
		q := buildTagQuery(d, runID)

		assertContainsAll(t, q,
			expectScope("env"),
			expectTag("dev"),
			expectScope("nsx-tf/tf-run-id"),
			expectTag(runID),
		)
	})

	// 3) Global-search mode: no tag attribute means no additional query.
	t.Run("global-search/no-user-tags", func(t *testing.T) {
		old := os.Getenv("NSXT_CACHE_MODE")
		t.Cleanup(func() { _ = os.Setenv("NSXT_CACHE_MODE", old) })
		_ = os.Setenv("NSXT_CACHE_MODE", "global")

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
			{Scope: newString("nsx-tf/tf-run-id"), Tag: newString(runID)},
			{Scope: newString("env"), Tag: newString("dev")},
		}}
		called := false

		patched, err := ensureProviderManagedTagsWithPatchFunc(obj, m, func(o *testTagObj) error {
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
			{Scope: newString("nsx-tf/tf-run-id"), Tag: newString("old-run")},
			{Scope: newString("env"), Tag: newString("dev")},
		}}
		called := false

		patched, err := ensureProviderManagedTagsWithPatchFunc(obj, m, func(o *testTagObj) error {
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

		val, ok := findTag(obj.Tags, "nsx-tf/tf-run-id")
		if !ok {
			t.Fatalf("expected nsx-tf/tf-run-id tag to exist")
		}
		if val != runID {
			t.Fatalf("expected nsx-tf/tf-run-id tag to be %q, got %q", runID, val)
		}
	})
}

func int64Ptr(v int64) *int64 {
	return &v
}

func TestGroupRulesByValidParentPath(t *testing.T) {
	pathA := "/policies/a"
	pathB := "/policies/b"
	valid := map[string]struct{}{pathA: {}, pathB: {}}

	t.Run("partitions-preserves-input-order", func(t *testing.T) {
		rules := []model.Rule{
			{Id: strPtr("b"), ParentPath: strPtr(pathA), SequenceNumber: int64Ptr(2)},
			{Id: strPtr("a"), ParentPath: strPtr(pathA), SequenceNumber: int64Ptr(1)},
			{Id: strPtr("x"), ParentPath: strPtr(pathB), SequenceNumber: int64Ptr(1)},
		}
		got := groupRulesByValidParentPath(valid, rules)
		if len(got[pathA]) != 2 || *got[pathA][0].Id != "b" || *got[pathA][1].Id != "a" {
			t.Fatalf("pathA bucket: %+v", got[pathA])
		}
		if len(got[pathB]) != 1 || *got[pathB][0].Id != "x" {
			t.Fatalf("pathB bucket: %+v", got[pathB])
		}
	})

	t.Run("skips-nil-and-empty-parent-path", func(t *testing.T) {
		rules := []model.Rule{
			{Id: strPtr("no-parent")},
			{Id: strPtr("blank"), ParentPath: strPtr("  ")},
			{Id: strPtr("ok"), ParentPath: strPtr(pathA), SequenceNumber: int64Ptr(1)},
		}
		got := groupRulesByValidParentPath(valid, rules)
		if len(got[pathA]) != 1 || *got[pathA][0].Id != "ok" {
			t.Fatalf("got %v", got[pathA])
		}
	})

	t.Run("drops-orphan-parent-path", func(t *testing.T) {
		rules := []model.Rule{
			{Id: strPtr("orphan"), ParentPath: strPtr("/unknown/parent"), SequenceNumber: int64Ptr(1)},
		}
		got := groupRulesByValidParentPath(valid, rules)
		if len(got) != 0 {
			t.Fatalf("expected empty map, got %v", got)
		}
	})

	t.Run("empty-valid-set", func(t *testing.T) {
		got := groupRulesByValidParentPath(map[string]struct{}{}, []model.Rule{
			{Id: strPtr("r"), ParentPath: strPtr(pathA), SequenceNumber: int64Ptr(1)},
		})
		if len(got) != 0 {
			t.Fatalf("expected no buckets, got %v", got)
		}
	})
}

func attachRulesToGatewayPoliciesForTest(parents []model.GatewayPolicy, rules []model.Rule) []model.GatewayPolicy {
	return attachRulesByParentPath(parents, rules,
		func(p model.GatewayPolicy) *string { return p.Path },
		func(p *model.GatewayPolicy, r []model.Rule) { p.Rules = r },
	)
}

func TestAttachRulesByParentPathGatewayPolicy(t *testing.T) {
	policyPathA := "/orgs/p/proj/vpcs/vpc/gateway-policies/pol-a"
	policyPathB := "/orgs/p/proj/vpcs/vpc/gateway-policies/pol-b"

	t.Run("happy-path-partitions-by-parent-path", func(t *testing.T) {
		parents := []model.GatewayPolicy{
			{Path: strPtr(policyPathA), Id: strPtr("pol-a")},
			{Path: strPtr(policyPathB), Id: strPtr("pol-b")},
		}
		rules := []model.Rule{
			{Id: strPtr("r1"), ParentPath: strPtr(policyPathA), SequenceNumber: int64Ptr(1)},
			{Id: strPtr("r2"), ParentPath: strPtr(policyPathB), SequenceNumber: int64Ptr(1)},
		}
		got := attachRulesToGatewayPoliciesForTest(parents, rules)
		if len(got) != 2 {
			t.Fatalf("expected 2 policies, got %d", len(got))
		}
		if len(got[0].Rules) != 1 || got[0].Rules[0].Id == nil || *got[0].Rules[0].Id != "r1" {
			t.Fatalf("policy A rules: %+v", got[0].Rules)
		}
		if len(got[1].Rules) != 1 || got[1].Rules[0].Id == nil || *got[1].Rules[0].Id != "r2" {
			t.Fatalf("policy B rules: %+v", got[1].Rules)
		}
	})

	t.Run("preserves-input-order", func(t *testing.T) {
		parents := []model.GatewayPolicy{{Path: strPtr(policyPathA)}}
		rules := []model.Rule{
			{Id: strPtr("b"), ParentPath: strPtr(policyPathA), SequenceNumber: int64Ptr(2)},
			{Id: strPtr("a"), ParentPath: strPtr(policyPathA), SequenceNumber: int64Ptr(1)},
			{Id: strPtr("c"), ParentPath: strPtr(policyPathA), SequenceNumber: int64Ptr(2)},
		}
		got := attachRulesToGatewayPoliciesForTest(parents, rules)
		if len(got[0].Rules) != 3 {
			t.Fatalf("expected 3 rules, got %d", len(got[0].Rules))
		}
		ids := []string{*got[0].Rules[0].Id, *got[0].Rules[1].Id, *got[0].Rules[2].Id}
		if ids[0] != "b" || ids[1] != "a" || ids[2] != "c" {
			t.Fatalf("expected input order b,a,c got %v", ids)
		}
	})

	t.Run("orphan-rule-discarded", func(t *testing.T) {
		parents := []model.GatewayPolicy{{Path: strPtr(policyPathA)}}
		rules := []model.Rule{
			{Id: strPtr("ok"), ParentPath: strPtr(policyPathA), SequenceNumber: int64Ptr(1)},
			{Id: strPtr("orphan"), ParentPath: strPtr("/other/policy/path"), SequenceNumber: int64Ptr(1)},
		}
		got := attachRulesToGatewayPoliciesForTest(parents, rules)
		if len(got[0].Rules) != 1 || got[0].Rules[0].Id == nil || *got[0].Rules[0].Id != "ok" {
			t.Fatalf("rules: %+v", got[0].Rules)
		}
	})

	t.Run("rule-missing-parent-path-skipped", func(t *testing.T) {
		parents := []model.GatewayPolicy{{Path: strPtr(policyPathA), Id: strPtr("pol-a")}}
		rules := []model.Rule{
			{Id: strPtr("no-parent")},
			{Id: strPtr("ok"), ParentPath: strPtr(policyPathA), SequenceNumber: int64Ptr(1)},
		}
		got := attachRulesToGatewayPoliciesForTest(parents, rules)
		if len(got[0].Rules) != 1 || got[0].Rules[0].Id == nil || *got[0].Rules[0].Id != "ok" {
			t.Fatalf("rules: %+v", got[0].Rules)
		}
	})

	t.Run("rule-empty-parent-path-skipped", func(t *testing.T) {
		parents := []model.GatewayPolicy{{Path: strPtr(policyPathA)}}
		rules := []model.Rule{
			{Id: strPtr("blank-parent"), ParentPath: strPtr("   ")},
			{Id: strPtr("ok"), ParentPath: strPtr(policyPathA), SequenceNumber: int64Ptr(1)},
		}
		got := attachRulesToGatewayPoliciesForTest(parents, rules)
		if len(got[0].Rules) != 1 || *got[0].Rules[0].Id != "ok" {
			t.Fatalf("rules: %+v", got[0].Rules)
		}
	})

	t.Run("trims-path-and-parent-path", func(t *testing.T) {
		parents := []model.GatewayPolicy{{Path: strPtr("  " + policyPathA + "  ")}}
		rules := []model.Rule{{Id: strPtr("r1"), ParentPath: strPtr(policyPathA), SequenceNumber: int64Ptr(1)}}
		got := attachRulesToGatewayPoliciesForTest(parents, rules)
		if len(got[0].Rules) != 1 {
			t.Fatalf("expected 1 rule after trim, got %d", len(got[0].Rules))
		}
	})

	t.Run("policy-without-path-gets-empty-rules", func(t *testing.T) {
		parents := []model.GatewayPolicy{{Id: strPtr("no-path")}}
		rules := []model.Rule{{Id: strPtr("r1"), ParentPath: strPtr(policyPathA), SequenceNumber: int64Ptr(1)}}
		got := attachRulesToGatewayPoliciesForTest(parents, rules)
		if len(got[0].Rules) != 0 {
			t.Fatalf("expected empty rules, got %+v", got[0].Rules)
		}
	})

	t.Run("empty-parents", func(t *testing.T) {
		got := attachRulesToGatewayPoliciesForTest(nil, []model.Rule{{Id: strPtr("x"), ParentPath: strPtr(policyPathA)}})
		if len(got) != 0 {
			t.Fatalf("expected no policies, got %d", len(got))
		}
	})
}
