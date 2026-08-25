package nsxt

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
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
		m := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "config_scope"}}

		d := schema.TestResourceDataRaw(t, tagSchema, map[string]interface{}{})
		q := buildTagQuery(d, runID, m)

		assertContainsAll(t, q,
			expectScope("nsx-tf/tf-run-id"),
			expectTag(runID),
		)
	})

	// 2) Tag-mode: if user tags are present and managed tags are not, provider-managed tags are appended.
	t.Run("tag-mode/user-tags-appends-provider-tags", func(t *testing.T) {
		m := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "config_scope"}}

		d := schema.TestResourceDataRaw(t, tagSchema, map[string]interface{}{
			"tag": []interface{}{
				map[string]interface{}{"scope": "env", "tag": "dev"},
			},
		})
		q := buildTagQuery(d, runID, m)

		assertContainsAll(t, q,
			expectScope("env"),
			expectTag("dev"),
			expectScope("nsx-tf/tf-run-id"),
			expectTag(runID),
		)
	})

	// 3) Global-search mode: no tag attribute means no additional query.
	t.Run("global-search/no-user-tags", func(t *testing.T) {
		m := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "global"}}

		d := schema.TestResourceDataRaw(t, tagSchema, map[string]interface{}{})
		q := buildTagQuery(d, runID, m)
		if q != "" {
			t.Fatalf("expected empty query in global-search mode when no tags exist, got %q", q)
		}
	})
}

func TestAttachRulesByParentPathSecurityPolicy(t *testing.T) {
	policyPathA := "/infra/domains/default/security-policies/pol-a"
	policyPathB := "/infra/domains/default/security-policies/pol-b"

	t.Run("happy-path-partitions-by-parent-path", func(t *testing.T) {
		parents := []model.SecurityPolicy{
			{Path: strPtr(policyPathA), Id: strPtr("pol-a")},
			{Path: strPtr(policyPathB), Id: strPtr("pol-b")},
		}
		rules := []model.Rule{
			{Id: strPtr("r1"), ParentPath: strPtr(policyPathA), SequenceNumber: int64Ptr(1)},
			{Id: strPtr("r2"), ParentPath: strPtr(policyPathB), SequenceNumber: int64Ptr(1)},
		}
		got := attachRulesToSecurityPoliciesForTest(parents, rules)
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

	t.Run("orphan-rule-discarded", func(t *testing.T) {
		parents := []model.SecurityPolicy{{Path: strPtr(policyPathA), Id: strPtr("pol-a")}}
		rules := []model.Rule{
			{Id: strPtr("ok"), ParentPath: strPtr(policyPathA), SequenceNumber: int64Ptr(1)},
			{Id: strPtr("orphan"), ParentPath: strPtr("/infra/domains/default/security-policies/other"), SequenceNumber: int64Ptr(1)},
		}
		got := attachRulesToSecurityPoliciesForTest(parents, rules)
		if len(got[0].Rules) != 1 || got[0].Rules[0].Id == nil || *got[0].Rules[0].Id != "ok" {
			t.Fatalf("rules: %+v", got[0].Rules)
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

	t.Run("partitions-sorted-by-sequence-number", func(t *testing.T) {
		// NSX evaluates rules in ascending sequence_number order, which may not match
		// the order the Search API returned them in, so buckets must be sorted explicitly.
		rules := []model.Rule{
			{Id: strPtr("b"), ParentPath: strPtr(pathA), SequenceNumber: int64Ptr(2)},
			{Id: strPtr("a"), ParentPath: strPtr(pathA), SequenceNumber: int64Ptr(1)},
			{Id: strPtr("x"), ParentPath: strPtr(pathB), SequenceNumber: int64Ptr(1)},
		}
		got := groupRulesByValidParentPath(valid, rules)
		if len(got[pathA]) != 2 || *got[pathA][0].Id != "a" || *got[pathA][1].Id != "b" {
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

func attachRulesToSecurityPoliciesForTest(parents []model.SecurityPolicy, rules []model.Rule) []model.SecurityPolicy {
	return attachRulesByParentPath(parents, rules,
		func(p model.SecurityPolicy) *string { return p.Path },
		func(p *model.SecurityPolicy, r []model.Rule) { p.Rules = r },
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

	t.Run("sorted-by-sequence-number-ties-preserve-input-order", func(t *testing.T) {
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
		// a (seq 1) sorts first; b and c tie at seq 2 so their relative input order (b before c) is preserved.
		if ids[0] != "a" || ids[1] != "b" || ids[2] != "c" {
			t.Fatalf("expected sorted order a,b,c got %v", ids)
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

func TestGetQueryStringVPCScopedToProjectNotVPC(t *testing.T) {
	// VPCID must be omitted from the cache bucket key: NSX policy paths/IDs are unique
	// within a project across all VPCs, and narrowing the key (and the underlying search)
	// to a single VPC caused a fresh cache bucket per VPC, regressing cache mode below
	// no-cache performance for VPC-scoped resource types touching many VPCs.
	for _, clientType := range []utl.ClientType{utl.VPC, utl.Multitenancy} {
		context := utl.SessionContext{ClientType: clientType, ProjectID: "proj-1", VPCID: "vpc-1"}
		got := getQueryString(resourceTypeVpcAttachment, context)
		if strings.Contains(got, "vpc-1") {
			t.Fatalf("clientType=%v: query %q must not be scoped to a specific VPCID", clientType, got)
		}
		if !strings.Contains(got, "proj-1") {
			t.Fatalf("clientType=%v: query %q must still be scoped to the project", got, got)
		}

		otherVPC := context
		otherVPC.VPCID = "vpc-2"
		if got2 := getQueryString(resourceTypeVpcAttachment, otherVPC); got2 != got {
			t.Fatalf("clientType=%v: query must be identical across VPCs in the same project so the cache bucket is shared; got %q vs %q", clientType, got, got2)
		}
	}
}

func TestProjectScopedSearchContextStripsVPCID(t *testing.T) {
	for _, clientType := range []utl.ClientType{utl.VPC, utl.Multitenancy} {
		in := utl.SessionContext{ClientType: clientType, ProjectID: "proj-1", VPCID: "vpc-1"}
		out := projectScopedSearchContext(in)
		if out.VPCID != "" {
			t.Fatalf("clientType=%v: expected VPCID stripped, got %q", clientType, out.VPCID)
		}
		if out.ProjectID != "proj-1" {
			t.Fatalf("clientType=%v: ProjectID must be preserved, got %q", clientType, out.ProjectID)
		}
	}

	// Non-VPC-scoped contexts must be returned unchanged.
	local := utl.SessionContext{ClientType: utl.Local, ProjectID: "", VPCID: ""}
	if got := projectScopedSearchContext(local); got != local {
		t.Fatalf("Local context should be unchanged, got %+v", got)
	}
}

func TestShouldIndexByPathForVPCScopedTypes(t *testing.T) {
	// VPC-scoped types must key by path, not short id: NSX ids for these types (often
	// user-chosen via nsx_id) are only guaranteed unique within their own VPC, but the cache
	// populate search/bucket for these types is now shared across all VPCs in a project.
	vpcScopedTypes := []string{
		resourceTypeVpc, resourceTypeVpcAttachment, resourceTypeVpcConnectivityProfile,
		resourceTypeVpcIpAddressAllocation, resourceTypeVpcServiceProfile, resourceTypeVpcSubnet,
		resourceTypeTransitGateway, resourceTypeTransitGatewayAttachment,
		resourceTypeProjectIpAddressAllocation, resourceTypePolicyVpcNatRule,
	}
	for _, rt := range vpcScopedTypes {
		if !shouldIndexByPath(rt) {
			t.Errorf("shouldIndexByPath(%q) = false, want true", rt)
		}
	}
}

func TestConverListToMapByTypeVpcScopedResourcesIndexedByPath(t *testing.T) {
	// Two different VPCs in the same project can legitimately have a VpcSubnet with the same
	// user-chosen short id (getOrGenerateID2 only checks uniqueness within the current VPC).
	// Since the cache bucket for VPC-scoped types is now shared project-wide, both objects
	// land in the same map; this test confirms each remains independently retrievable via its
	// distinct (project/system-unique) path, even though they share a colliding short id.
	pathA := "/orgs/o/projects/p/vpcs/vpcA/subnets/subnet1"
	pathB := "/orgs/o/projects/p/vpcs/vpcB/subnets/subnet1"
	subnetA := model.VpcSubnet{Id: strPtr("subnet1"), DisplayName: strPtr("subnet1-a"), Path: strPtr(pathA)}
	subnetB := model.VpcSubnet{Id: strPtr("subnet1"), DisplayName: strPtr("subnet1-b"), Path: strPtr(pathB)}

	svs, err := modelsToStructValues([]model.VpcSubnet{subnetA, subnetB}, model.VpcSubnetBindingType())
	if err != nil {
		t.Fatalf("modelsToStructValues: %v", err)
	}

	got := converListToMapByType(svs, resourceTypeVpcSubnet)
	if got == nil {
		t.Fatal("converListToMapByType returned nil")
	}
	if got[pathA] == nil {
		t.Errorf("VPC A's subnet not retrievable by its path %q", pathA)
	}
	if got[pathB] == nil {
		t.Errorf("VPC B's subnet not retrievable by its path %q", pathB)
	}
}

func TestErrCacheUseBackendDirect(t *testing.T) {
	if !errors.Is(errCacheUseBackendDirect, errCacheUseBackendDirect) {
		t.Fatal("errors.Is should match sentinel to itself")
	}
	wrapped := fmt.Errorf("wrap: %w", errCacheUseBackendDirect)
	if !errors.Is(wrapped, errCacheUseBackendDirect) {
		t.Fatal("errors.Is should unwrap to sentinel")
	}
}

func TestReflectStringField(t *testing.T) {
	t.Run("returns-pointer-value", func(t *testing.T) {
		obj := &model.Group{DisplayName: strPtr("g1")}
		got := reflectStringField(obj, "DisplayName")
		if got == nil || *got != "g1" {
			t.Fatalf("expected \"g1\", got %v", got)
		}
	})

	t.Run("nil-field-returns-nil", func(t *testing.T) {
		obj := &model.Group{}
		if got := reflectStringField(obj, "DisplayName"); got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
	})

	t.Run("missing-field-returns-nil", func(t *testing.T) {
		obj := &model.Group{DisplayName: strPtr("g1")}
		if got := reflectStringField(obj, "NoSuchField"); got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
	})

	t.Run("non-pointer-string-field-returns-nil", func(t *testing.T) {
		// SequenceNumber is *int64, not *string.
		obj := &model.Rule{SequenceNumber: int64Ptr(5)}
		if got := reflectStringField(obj, "SequenceNumber"); got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
	})

	t.Run("nil-obj-returns-nil", func(t *testing.T) {
		if got := reflectStringField(nil, "DisplayName"); got != nil {
			t.Fatalf("expected nil, got %v", got)
		}
	})
}
