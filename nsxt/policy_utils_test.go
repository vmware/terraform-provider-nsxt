//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sdkerrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

type policyPathTest struct {
	path    string
	parents []string
}

func TestUnitNsxt_ParseStandardPolicyPath(t *testing.T) {

	testData := []policyPathTest{
		{
			path:    "/infra/tier-1s/mygw1",
			parents: []string{"mygw1"},
		},
		{
			path:    "/global-infra/tier-0s/mygw1",
			parents: []string{"mygw1"},
		},
		{
			path:    "/orgs/myorg/projects/myproj/infra/tier-1s/mygw1",
			parents: []string{"myorg", "myproj", "mygw1"},
		},
		{
			path:    "/orgs/myorg/projects/myproj/vpcs/myvpc/tier-1s/mygw1",
			parents: []string{"myorg", "myproj", "myvpc", "mygw1"},
		},
		{
			path:    "/orgs/myorg/projects/myproj/infra/domains/d1/groups/g1",
			parents: []string{"myorg", "myproj", "d1", "g1"},
		},
		{
			path:    "/global-infra/tier-1s/t15/ipsec-vpn-services/default/sessions/xxx-yyy-xxx",
			parents: []string{"t15", "default", "xxx-yyy-xxx"},
		},
		{
			path:    "/infra/evpn-tenant-configs/{config-id}",
			parents: []string{"{config-id}"},
		},
		{
			path:    "/infra/tier-0s/{tier-0-id}/locale-services/{locale-service-id}/service-interfaces/{interface-id}",
			parents: []string{"{tier-0-id}", "{locale-service-id}", "{interface-id}"},
		},
	}

	for _, test := range testData {
		parents, err := parseStandardPolicyPath(test.path)
		assert.Nil(t, err)
		assert.Equal(t, test.parents, parents)
	}
}

func TestUnitNsxt_IsPolicyPath(t *testing.T) {

	testData := []string{
		"/infra/tier-1s/mygw1",
		"/global-infra/tier-1s/mygw1",
		"/orgs/infra/tier-1s/mygw1",
		"/orgs/myorg/projects/myproj/domains/d",
		"/orgs/myorg/projects/myproj/vpcs/nicevpc",
	}

	for _, test := range testData {
		pp := isPolicyPath(test)
		assert.True(t, pp)
	}
}

func TestUnitNsxt_NegativeParseStandardPolicyPath(t *testing.T) {

	testData := []string{
		"/a",
		"orgs/infra/tier-1s/mygw1-1",
		"/orgs/myorg",
	}

	for _, test := range testData {
		_, err := parseStandardPolicyPath(test)
		assert.NotNil(t, err)
	}
}

func TestUnitNsxt_ParseStandardPolicyPathVerifySize(t *testing.T) {

	_, err := parseStandardPolicyPathVerifySize("/infra/things/thing1/sub-things/sub-thing1", 3, "sample")
	assert.NotNil(t, err)

	parents, err := parseStandardPolicyPathVerifySize("/infra/things/thing1/sub-things/sub-thing1", 2, "sample")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(parents))

	_, err = parseStandardPolicyPathVerifySize("/global-infra/things/1/sub-things/2/fine-tuned-thing/3", 1, "sample")
	assert.NotNil(t, err)
}

func TestUnitNsxt_GetResourceIDFromResourcePath(t *testing.T) {
	assert.Equal(t, "default", getResourceIDFromResourcePath("/infra/domains/default/groups/g1", "domains"))
	assert.Equal(t, "myproj", getResourceIDFromResourcePath("/orgs/acme/projects/myproj/infra/tier-1s/t1", "projects"))
	assert.Equal(t, "", getResourceIDFromResourcePath("/infra/tier-1s/t1", "domains"))
}

func TestUnitNsxt_GetDomainFromResourcePath(t *testing.T) {
	assert.Equal(t, "default", getDomainFromResourcePath("/infra/domains/default/gateway-policies/gw1"))
}

func TestUnitNsxt_GetProjectIDFromResourcePath(t *testing.T) {
	assert.Equal(t, "p1", getProjectIDFromResourcePath("/orgs/o1/projects/p1/infra/domains/d1"))
}

func TestUnitNsxt_GetPolicyIDFromPath(t *testing.T) {
	assert.Equal(t, "rule-a", getPolicyIDFromPath("/infra/domains/default/gateway-policies/pol/rules/rule-a"))
}

func TestUnitNsxt_GetParameterFromPolicyPath(t *testing.T) {
	v, err := getParameterFromPolicyPath("pre-", "-post", "pre-value-post")
	assert.NoError(t, err)
	assert.Equal(t, "value", v)

	_, err = getParameterFromPolicyPath("x", "y", "nope")
	assert.Error(t, err)
}

func TestUnitNsxt_ShouldIgnoreScope(t *testing.T) {
	assert.True(t, shouldIgnoreScope("s1", []string{"s0", "s1", "s2"}))
	assert.False(t, shouldIgnoreScope("sx", []string{"s0", "s1"}))
	assert.False(t, shouldIgnoreScope("s1", nil))
}

func TestUnitNsxt_CollectSeparatedStringListToMap(t *testing.T) {
	m := collectSeparatedStringListToMap([]string{"a:1", "b:2", "no-sep", "c:"}, ":")
	assert.Equal(t, "1", m["a"])
	assert.Equal(t, "2", m["b"])
	assert.NotContains(t, m, "no-sep")
}

func TestUnitNsxt_StringListCommaSeparatedRoundTrip(t *testing.T) {
	s := stringListToCommaSeparatedString([]string{"a", "b", "c"})
	assert.NotNil(t, s)
	assert.Equal(t, "a,b,c", *s)
	assert.Nil(t, stringListToCommaSeparatedString(nil))

	assert.Empty(t, commaSeparatedStringToStringList(""))
	assert.Equal(t, []string{"x", "y"}, commaSeparatedStringToStringList("x,y"))
	assert.Equal(t, []string{"x", "y"}, commaSeparatedStringToStringList("x,,y"))
}

func TestUnitNsxt_InterfaceListToStringList(t *testing.T) {
	assert.Equal(t, []string{"a", "b"}, interfaceListToStringList([]interface{}{"a", "b"}))
}

func TestUnitNsxt_IsValidResourceID(t *testing.T) {
	assert.True(t, isValidResourceID("abc-123"))
	assert.False(t, isValidResourceID(""))
	assert.False(t, isValidResourceID("   "))
	assert.False(t, isValidResourceID("a/b"))
	assert.False(t, isValidResourceID(string(make([]byte, 1025))))
}

func TestUnitNsxt_IsValidID(t *testing.T) {
	assert.True(t, isValidID("id-1"))
	assert.False(t, isValidID("a/b"))
	assert.False(t, isValidID("a&b"))
}

func TestUnitNsxt_IsSpaceString(t *testing.T) {
	assert.True(t, isSpaceString(""))
	assert.True(t, isSpaceString("  \t "))
	assert.False(t, isSpaceString("x"))
}

func TestUnitNsxt_RetryUponPreconditionFailed(t *testing.T) {
	t.Run("stops on success", func(t *testing.T) {
		calls := 0
		err := retryUponPreconditionFailed(func() error {
			calls++
			return nil
		}, 3)
		assert.NoError(t, err)
		assert.Equal(t, 1, calls)
	})

	t.Run("retries InvalidRequest then succeeds", func(t *testing.T) {
		calls := 0
		err := retryUponPreconditionFailed(func() error {
			calls++
			if calls < 2 {
				return sdkerrors.InvalidRequest{}
			}
			return nil
		}, 3)
		assert.NoError(t, err)
		assert.Equal(t, 2, calls)
	})

	t.Run("non-retryable error returns immediately", func(t *testing.T) {
		boom := errors.New("boom")
		calls := 0
		err := retryUponPreconditionFailed(func() error {
			calls++
			return boom
		}, 3)
		assert.Equal(t, boom, err)
		assert.Equal(t, 1, calls)
	})

	t.Run("returns last error after max attempts", func(t *testing.T) {
		calls := 0
		err := retryUponPreconditionFailed(func() error {
			calls++
			return sdkerrors.InvalidRequest{}
		}, 1)
		assert.Error(t, err)
		assert.Equal(t, 2, calls)
	})
}

func TestUnitNsxt_initPolicyTagsSet(t *testing.T) {
	s1, t1 := "sc", "tg"
	out := initPolicyTagsSet([]model.Tag{{Scope: &s1, Tag: &t1}})
	require.Len(t, out, 1)
	assert.Equal(t, &s1, out[0]["scope"])
	assert.Equal(t, &t1, out[0]["tag"])
}

func TestUnitNsxt_getIgnoredTagsFromSchema_emptyList(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"ignore_tags": getIgnoreTagsSchema(),
	}, map[string]interface{}{
		"ignore_tags": []interface{}{},
	})
	assert.Nil(t, getIgnoredTagsFromSchema(d))
}

func TestUnitNsxt_setIgnoredTagsInSchema(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"ignore_tags": getIgnoreTagsSchema(),
	}, map[string]interface{}{})
	setIgnoredTagsInSchema(d, []string{"s1"}, []map[string]interface{}{
		{"scope": "s1", "tag": "t1"},
	})
	v := d.Get("ignore_tags").([]interface{})
	require.Len(t, v, 1)
	m := v[0].(map[string]interface{})
	assert.Equal(t, []interface{}{"s1"}, m["scopes"])
}

func TestUnitNsxt_getPathListFromMap_and_setPathListInMap(t *testing.T) {
	t.Run("non-empty set", func(t *testing.T) {
		st := schema.NewSet(schema.HashString, []interface{}{"/p1", "/p2"})
		data := map[string]interface{}{"paths": st}
		assert.Equal(t, []string{"/p1", "/p2"}, getPathListFromMap(data, "paths"))
	})
	t.Run("empty set becomes ANY", func(t *testing.T) {
		st := schema.NewSet(schema.HashString, []interface{}{})
		data := map[string]interface{}{"paths": st}
		assert.Equal(t, []string{"ANY"}, getPathListFromMap(data, "paths"))
	})
	t.Run("setPathListInMap ANY clears", func(t *testing.T) {
		data := map[string]interface{}{"paths": "placeholder"}
		setPathListInMap(data, "paths", []string{"ANY"})
		assert.Nil(t, data["paths"])
	})
	t.Run("setPathListInMap paths", func(t *testing.T) {
		data := map[string]interface{}{}
		setPathListInMap(data, "paths", []string{"/a"})
		assert.Equal(t, []string{"/a"}, data["paths"])
	})
}

func TestUnitNsxt_validateImportPolicyPath(t *testing.T) {
	assert.ErrorIs(t, validateImportPolicyPath("   "), ErrEmptyImportID)
	assert.ErrorIs(t, validateImportPolicyPath("not-a-valid-uri-for-request-parser"), ErrNotAPolicyPath)
	assert.ErrorIs(t, validateImportPolicyPath("/infra//tier-0s/x"), ErrNotAPolicyPath)
	assert.ErrorIs(t, validateImportPolicyPath("/wrong/tier-0s/x"), ErrNotAPolicyPath)
	require.NoError(t, validateImportPolicyPath("/infra/tier-0s/mygw"))
}

func TestUnitNsxt_getMultitenancyPathExample(t *testing.T) {
	out := getMultitenancyPathExample("/infra/foo")
	assert.Contains(t, out, "/infra/foo")
	assert.Contains(t, out, "/orgs/[org]/projects/[project]")
}

func TestUnitNsxt_getGlobalPolicyEnforcementPointPathWithLocation(t *testing.T) {
	c := newGoMockProviderClient()
	c.PolicyEnforcementPoint = "ep99"
	assert.Equal(t, "/global-infra/sites/site-a/enforcement-points/ep99", getGlobalPolicyEnforcementPointPathWithLocation(c, "site-a"))
}

func TestUnitNsxt_lbHTTPHeaderSchemaRoundTrip(t *testing.T) {
	elem := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":  {Type: schema.TypeString, Required: true},
			"value": {Type: schema.TypeString, Required: true},
		},
	}
	sch := map[string]*schema.Schema{
		"headers": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     elem,
		},
	}
	hn, hv := "X-Custom", "abc"
	// Avoid embedding *schema.Set in TestResourceDataRaw (SDK/HCL2 panic); populate via Set.
	d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{})
	require.NoError(t, d.Set("headers", []interface{}{
		map[string]interface{}{"name": hn, "value": hv},
	}))
	out := getPolicyLbHTTPHeaderFromSchema(d, "headers")
	require.Len(t, out, 1)
	assert.Equal(t, hn, *out[0].HeaderName)
	assert.Equal(t, hv, *out[0].HeaderValue)

	d2 := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"headers": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name":  {Type: schema.TypeString, Required: true},
					"value": {Type: schema.TypeString, Required: true},
				},
			},
		},
	}, map[string]interface{}{})
	setPolicyLbHTTPHeaderInSchema(d2, "headers", out)
	lst := d2.Get("headers").([]interface{})
	require.Len(t, lst, 1)
	assert.Equal(t, hn, lst[0].(map[string]interface{})["name"])
}

func TestUnitNsxt_nsxtParentPathResourceImporter(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"parent_path": {Type: schema.TypeString, Optional: true, Computed: true},
	}, map[string]interface{}{})
	d.SetId("/infra/domains/default/gateway-policies/pol-1")
	rd, err := nsxtParentPathResourceImporter(d, nil)
	require.NoError(t, err)
	require.Len(t, rd, 1)
	assert.Equal(t, "pol-1", rd[0].Id())
	// Importer strips the last two path segments (resource id + type) from the import path.
	assert.Equal(t, "/infra/domains/default", rd[0].Get("parent_path"))
}

func TestUnitNsxt_getVpcParentsFromContext(t *testing.T) {
	out := getVpcParentsFromContext(utl.SessionContext{ProjectID: "proj1", VPCID: "vpc1"})
	assert.Equal(t, []string{utl.DefaultOrgID, "proj1", "vpc1"}, out)
}
