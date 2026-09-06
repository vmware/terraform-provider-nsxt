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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

// cacheTestSchema is a resource schema covering every attribute the cache helpers under
// test read from *schema.ResourceData: tags (buildTagQuery), parent_path/context
// (getEffectiveCacheContext), and path (CacheKeyForResourceID).
func cacheTestSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id":           {Type: schema.TypeString, Optional: true, Computed: true},
		"display_name": {Type: schema.TypeString, Optional: true, Computed: true},
		"description":  {Type: schema.TypeString, Optional: true, Computed: true},
		"path":         {Type: schema.TypeString, Optional: true, Computed: true},
		"tag": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"scope": {Type: schema.TypeString, Optional: true},
				"tag":   {Type: schema.TypeString, Optional: true},
			}},
		},
		"parent_path": {Type: schema.TypeString, Optional: true},
		"context":     getContextSchema(false, false, true),
	}
}

func groupStructValue(t *testing.T, id, displayName, path string) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(nsxModel.Group{
		Id: strPtr(id), DisplayName: strPtr(displayName), Path: strPtr(path),
	}, nsxModel.GroupBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func TestUnitNsxt_postWriteKey(t *testing.T) {
	assert.Equal(t, "Group::g1", postWriteKey("Group", "g1"))
	assert.Equal(t, "::", postWriteKey("", ""))
}

func TestUnitNsxt_indexCacheMapKey(t *testing.T) {
	ret := map[string]*data.StructValue{}
	first := groupStructValue(t, "g1", "first", "/p1")
	second := groupStructValue(t, "g1", "second", "/p2")

	indexCacheMapKey(ret, "k", first)
	indexCacheMapKey(ret, "k", second)

	assert.Same(t, first, ret["k"], "first-seen value must win")
}

func TestUnitNsxt_getStructValueStringField(t *testing.T) {
	t.Run("nil object", func(t *testing.T) {
		_, ok := getStructValueStringField(nil, "id")
		assert.False(t, ok)
	})

	t.Run("missing field", func(t *testing.T) {
		sv := data.NewStructValue("X", map[string]data.DataValue{})
		_, ok := getStructValueStringField(sv, "id")
		assert.False(t, ok)
	})

	t.Run("plain string field", func(t *testing.T) {
		sv := data.NewStructValue("X", map[string]data.DataValue{"id": data.NewStringValue("abc")})
		v, ok := getStructValueStringField(sv, "id")
		require.True(t, ok)
		assert.Equal(t, "abc", v)
	})

	t.Run("set optional string field", func(t *testing.T) {
		sv := data.NewStructValue("X", map[string]data.DataValue{
			"id": data.NewOptionalValue(data.NewStringValue("opt")),
		})
		v, ok := getStructValueStringField(sv, "id")
		require.True(t, ok)
		assert.Equal(t, "opt", v)
	})

	t.Run("unset optional field", func(t *testing.T) {
		sv := data.NewStructValue("X", map[string]data.DataValue{
			"id": data.NewOptionalValue(nil),
		})
		_, ok := getStructValueStringField(sv, "id")
		assert.False(t, ok)
	})

	t.Run("non-string field type", func(t *testing.T) {
		sv := data.NewStructValue("X", map[string]data.DataValue{"count": data.NewIntegerValue(3)})
		_, ok := getStructValueStringField(sv, "count")
		assert.False(t, ok)
	})
}

func TestUnitNsxt_getQueryString(t *testing.T) {
	cases := []struct {
		name string
		ctx  utl.SessionContext
		want string
	}{
		{"global", utl.SessionContext{ClientType: utl.Global}, "resource_type:Group AND marked_for_delete:false AND context:Global"},
		{"local", utl.SessionContext{ClientType: utl.Local}, "resource_type:Group AND marked_for_delete:false AND context:Local"},
		{"vpc", utl.SessionContext{ClientType: utl.VPC, ProjectID: "p1", VPCID: "v1"}, "resource_type:Group AND marked_for_delete:false AND context:p1"},
		{"multitenancy", utl.SessionContext{ClientType: utl.Multitenancy, ProjectID: "p1", VPCID: ""}, "resource_type:Group AND marked_for_delete:false AND context:p1"},
		{"default", utl.SessionContext{ClientType: utl.ClientType(99)}, "resource_type:Group AND marked_for_delete:false"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, getQueryString("Group", c.ctx))
		})
	}
}

func TestUnitNsxt_typeScopedCacheGetTypeCache(t *testing.T) {
	c := &typeScopedCache{byTyp: make(map[string]*resourceTypeCache)}
	tc1 := c.getTypeCache("Group")
	tc2 := c.getTypeCache("Group")
	assert.Same(t, tc1, tc2, "repeated calls for the same type must return the same bucket")

	tc3 := c.getTypeCache("Service")
	assert.NotSame(t, tc1, tc3)
}

func TestUnitNsxt_getEffectiveCacheContext(t *testing.T) {
	m := nsxtClients{}

	t.Run("parent_path drives context when set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{
			"parent_path": "/orgs/o1/projects/p1/vpcs/v1/gateway-policies/x",
		})
		ctx := getEffectiveCacheContext(d, m)
		assert.Equal(t, "p1", ctx.ProjectID)
		assert.Equal(t, "v1", ctx.VPCID)
		assert.EqualValues(t, utl.VPC, ctx.ClientType)
	})

	t.Run("falls back to session context when parent_path is empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
		ctx := getEffectiveCacheContext(d, m)
		assert.EqualValues(t, utl.Local, ctx.ClientType)
	})
}

func TestUnitNsxt_getCacheQueryKey(t *testing.T) {
	m := nsxtClients{Host: "10.0.0.1", CommonConfig: commonProviderConfig{CacheMode: "config_scope"}}

	t.Run("no additional query", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
		key := getCacheQueryKey(resourceTypeGroup, d, m)
		assert.Equal(t, "[10.0.0.1]resource_type:Group AND marked_for_delete:false AND context:Local", key)
	})

	t.Run("appends tag query", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{
			"tag": []interface{}{map[string]interface{}{"scope": "env", "tag": "dev"}},
		})
		key := getCacheQueryKey(resourceTypeGroup, d, m)
		assert.Contains(t, key, "[10.0.0.1]resource_type:Group AND marked_for_delete:false AND context:Local AND ")
		assert.Contains(t, key, "tags.scope:env")
		assert.Contains(t, key, "tags.tag:dev")
	})
}

func TestUnitNsxt_resourceTypeCacheGetQueryResult(t *testing.T) {
	sv := groupStructValue(t, "g1", "grp", "/infra/domains/default/groups/g1")

	t.Run("bucket absent", func(t *testing.T) {
		c := &resourceTypeCache{data: map[string]map[string]*data.StructValue{}}
		_, err := c.getQueryResult("q", "g1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "is not found")
	})

	t.Run("bucket present, key hit", func(t *testing.T) {
		c := &resourceTypeCache{data: map[string]map[string]*data.StructValue{"q": {"g1": sv}}}
		got, err := c.getQueryResult("q", "g1")
		require.NoError(t, err)
		assert.Same(t, sv, got)
	})

	t.Run("bucket present, key not present is not found via any fallback scan", func(t *testing.T) {
		// converListToMapByType indexes every object under its raw "id" field at populate
		// time, so getQueryResult itself no longer scans bucket values to find a match by
		// a field other than the lookup key; a miss on the key falls straight to the caller.
		raw := data.NewStructValue("X", map[string]data.DataValue{"id": data.NewStringValue("scan-id")})
		c := &resourceTypeCache{data: map[string]map[string]*data.StructValue{"q": {"other-key": raw}}}
		_, err := c.getQueryResult("q", "scan-id")
		assert.True(t, errors.Is(err, errCacheUseBackendDirect))
	})

	t.Run("bucket present, id missing everywhere", func(t *testing.T) {
		c := &resourceTypeCache{data: map[string]map[string]*data.StructValue{"q": {"g1": sv}}}
		_, err := c.getQueryResult("q", "nope")
		assert.True(t, errors.Is(err, errCacheUseBackendDirect))
	})
}

func TestUnitNsxt_resourceTypeCacheWriteCacheSkipsExistingBucket(t *testing.T) {
	d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
	m := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "config_scope"}}
	c := &resourceTypeCache{data: map[string]map[string]*data.StructValue{"q1": {}}}

	_, err := c.writeCache("q1", resourceTypeGroup, d, m, nil)
	require.NoError(t, err)
	assert.Len(t, c.data, 1, "existing bucket must not be touched or replaced")
}

func TestUnitNsxt_readCachePopulatesOnMissAndHitsThereafter(t *testing.T) {
	sv := groupStructValue(t, "g1", "my-group", "/infra/domains/default/groups/g1")
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
		{Results: []*data.StructValue{sv}, ResultCount: i64(1)},
	}}
	defer setupCliQueryClientStub(t, stub)()

	d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
	m := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "config_scope"}}
	c := &typeScopedCache{byTyp: make(map[string]*resourceTypeCache)}

	val, err := c.readCache("g1", resourceTypeGroup, d, m, nil)
	require.NoError(t, err)
	require.NotNil(t, val)
	sv2, ok := val.(*data.StructValue)
	require.True(t, ok)
	gotID, ok := getStructValueStringField(sv2, "id")
	require.True(t, ok)
	assert.Equal(t, "g1", gotID)

	// Second read must hit the now-populated bucket without another search call.
	_, err = c.readCache("g1", resourceTypeGroup, d, m, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, stub.call, "second read should be served from cache")
}

func TestUnitNsxt_readCacheMissingIDBypassesBackend(t *testing.T) {
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
		{Results: []*data.StructValue{}, ResultCount: i64(0)},
	}}
	defer setupCliQueryClientStub(t, stub)()

	d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
	m := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "config_scope"}}
	c := &typeScopedCache{byTyp: make(map[string]*resourceTypeCache)}

	_, err := c.readCache("missing", resourceTypeGroup, d, m, nil)
	assert.True(t, errors.Is(err, errCacheUseBackendDirect))
}

func TestUnitNsxt_getListOfPolicyResourcesCompositeMergesRules(t *testing.T) {
	parentPath := "/infra/domains/default/security-policies/sp1"
	converter := bindings.NewTypeConverter()

	parentVal, errs := converter.ConvertToVapi(nsxModel.SecurityPolicy{
		Id: strPtr("sp1"), DisplayName: strPtr("sp-name"), Path: strPtr(parentPath),
	}, nsxModel.SecurityPolicyBindingType())
	require.Empty(t, errs)

	childVal, errs := converter.ConvertToVapi(nsxModel.Rule{
		Id: strPtr("r1"), DisplayName: strPtr("rule1"), ParentPath: strPtr(parentPath), SequenceNumber: int64Ptr(1),
	}, nsxModel.RuleBindingType())
	require.Empty(t, errs)

	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
		{Results: []*data.StructValue{parentVal.(*data.StructValue)}, ResultCount: i64(1)},
		{Results: []*data.StructValue{childVal.(*data.StructValue)}, ResultCount: i64(1)},
	}}
	defer setupCliQueryClientStub(t, stub)()

	d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
	m := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "config_scope"}}
	c := &typeScopedCache{byTyp: make(map[string]*resourceTypeCache)}

	val, err := c.readCache("sp1", resourceTypeSecurityPolicy, d, m, nil)
	require.NoError(t, err)
	require.NotNil(t, val)

	sv := val.(*data.StructValue)
	golangVal, errs := converter.ConvertToGolang(sv, nsxModel.SecurityPolicyBindingType())
	require.Empty(t, errs)
	sp := golangVal.(nsxModel.SecurityPolicy)
	require.Len(t, sp.Rules, 1)
	assert.Equal(t, "r1", *sp.Rules[0].Id)
}

func TestUnitNsxt_getListOfPolicyResourcesErrorWithNoResults(t *testing.T) {
	stub := &seqQueryListClient{errs: []error{errors.New("search failed")}}
	defer setupCliQueryClientStub(t, stub)()

	d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
	m := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "config_scope"}}
	c := &resourceTypeCache{data: map[string]map[string]*data.StructValue{}}

	_, err := c.getListOfPolicyResources("q", d, m, nil, utl.SessionContext{ClientType: utl.Local}, resourceTypeGroup, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error listing resource")
}

func TestUnitNsxt_convertCachedValue(t *testing.T) {
	m := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "config_scope"}}

	t.Run("non struct-value input fails", func(t *testing.T) {
		_, ok := convertCachedValue[nsxModel.Group]("not-a-struct-value", "Group", "g1", nsxModel.GroupBindingType(), m)
		assert.False(t, ok)
	})

	t.Run("successful conversion strips provider-managed tags outside global mode", func(t *testing.T) {
		converter := bindings.NewTypeConverter()
		val, errs := converter.ConvertToVapi(nsxModel.Group{
			Id: strPtr("g1"), DisplayName: strPtr("grp"),
			Tags: []nsxModel.Tag{
				{Scope: strPtr(managedDefaultTagScope), Tag: strPtr("run-1")},
				{Scope: strPtr("env"), Tag: strPtr("dev")},
			},
		}, nsxModel.GroupBindingType())
		require.Empty(t, errs)

		typed, ok := convertCachedValue[nsxModel.Group](val, "Group", "g1", nsxModel.GroupBindingType(), m)
		require.True(t, ok)
		require.Len(t, typed.Tags, 1)
		assert.Equal(t, "env", *typed.Tags[0].Scope)
	})

	t.Run("global mode keeps provider-managed tags", func(t *testing.T) {
		globalM := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "global"}}
		converter := bindings.NewTypeConverter()
		val, errs := converter.ConvertToVapi(nsxModel.Group{
			Id: strPtr("g1"),
			Tags: []nsxModel.Tag{
				{Scope: strPtr(managedDefaultTagScope), Tag: strPtr("run-1")},
			},
		}, nsxModel.GroupBindingType())
		require.Empty(t, errs)

		typed, ok := convertCachedValue[nsxModel.Group](val, "Group", "g1", nsxModel.GroupBindingType(), globalM)
		require.True(t, ok)
		require.Len(t, typed.Tags, 1)
	})
}

func TestUnitNsxt_structValuesAndModelsRoundTrip(t *testing.T) {
	group := nsxModel.Group{Id: strPtr("g1"), DisplayName: strPtr("grp"), Path: strPtr("/infra/domains/default/groups/g1")}

	svs, err := modelsToStructValues([]nsxModel.Group{group}, nsxModel.GroupBindingType())
	require.NoError(t, err)
	require.Len(t, svs, 1)

	models, err := structValuesToModels[nsxModel.Group](svs, nsxModel.GroupBindingType())
	require.NoError(t, err)
	require.Len(t, models, 1)
	assert.Equal(t, "g1", *models[0].Id)
	assert.Equal(t, "grp", *models[0].DisplayName)
}

func TestUnitNsxt_structValuesToRules(t *testing.T) {
	converter := bindings.NewTypeConverter()
	sv, errs := converter.ConvertToVapi(nsxModel.Rule{
		Id: strPtr("r1"), ParentPath: strPtr("/infra/domains/default/security-policies/sp1"), SequenceNumber: int64Ptr(2),
	}, nsxModel.RuleBindingType())
	require.Empty(t, errs)

	rules := structValuesToRules([]*data.StructValue{sv.(*data.StructValue)})
	require.Len(t, rules, 1)
	assert.Equal(t, "r1", *rules[0].Id)
	assert.EqualValues(t, 2, *rules[0].SequenceNumber)
}

func TestUnitNsxt_mergeGatewayPolicyCacheSearchResults(t *testing.T) {
	parentPath := "/infra/domains/default/gateway-policies/gp1"
	converter := bindings.NewTypeConverter()

	parentVal, errs := converter.ConvertToVapi(nsxModel.GatewayPolicy{
		Id: strPtr("gp1"), Path: strPtr(parentPath),
	}, nsxModel.GatewayPolicyBindingType())
	require.Empty(t, errs)

	childVal, errs := converter.ConvertToVapi(nsxModel.Rule{
		Id: strPtr("r1"), ParentPath: strPtr(parentPath), SequenceNumber: int64Ptr(1),
	}, nsxModel.RuleBindingType())
	require.Empty(t, errs)

	merged, err := mergeGatewayPolicyCacheSearchResults(
		[]*data.StructValue{parentVal.(*data.StructValue)},
		[]*data.StructValue{childVal.(*data.StructValue)},
	)
	require.NoError(t, err)
	require.Len(t, merged, 1)

	golangVal, errs := converter.ConvertToGolang(merged[0], nsxModel.GatewayPolicyBindingType())
	require.Empty(t, errs)
	gp := golangVal.(nsxModel.GatewayPolicy)
	require.Len(t, gp.Rules, 1)
	assert.Equal(t, "r1", *gp.Rules[0].Id)
}

func TestUnitNsxt_converListToMapByType(t *testing.T) {
	sv := groupStructValue(t, "g1", "my-group", "/infra/domains/default/groups/g1")
	ret := converListToMapByType([]*data.StructValue{sv}, resourceTypeGroup)
	require.NotNil(t, ret)
	assert.Same(t, sv, ret["g1"])
	assert.Same(t, sv, ret["my-group"])
}

func TestUnitNsxt_tryCacheRead(t *testing.T) {
	d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
	enabled := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "config_scope"}}
	disabled := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "disabled"}}

	t.Run("not refresh phase: no attempt", func(t *testing.T) {
		_, used, attempted, err := TryCacheRead[nsxModel.Group](d, enabled, nil, "g1", resourceTypeGroup, nsxModel.GroupBindingType())
		require.NoError(t, err)
		assert.False(t, used)
		assert.False(t, attempted)
	})

	t.Run("cache disabled: no attempt", func(t *testing.T) {
		d.SetId("g1")
		_, used, attempted, err := TryCacheRead[nsxModel.Group](d, disabled, nil, "g1", resourceTypeGroup, nsxModel.GroupBindingType())
		require.NoError(t, err)
		assert.False(t, used)
		assert.False(t, attempted)
	})

	t.Run("post-write bypass consumes the marker once", func(t *testing.T) {
		d.SetId("g1")
		postWriteByKey.Store(postWriteKey(resourceTypeGroup, "g1"), struct{}{})

		_, used, attempted, err := TryCacheRead[nsxModel.Group](d, enabled, nil, "g1", resourceTypeGroup, nsxModel.GroupBindingType())
		require.NoError(t, err)
		assert.False(t, used)
		assert.True(t, attempted)

		_, stillMarked := postWriteByKey.Load(postWriteKey(resourceTypeGroup, "g1"))
		assert.False(t, stillMarked, "post-write marker must be consumed on first read")
	})

	t.Run("cache hit returns the typed value", func(t *testing.T) {
		d2 := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
		d2.SetId("g1")
		resourceType := "TestOnlyGroupForTryCacheRead"
		sv := groupStructValue(t, "g1", "my-group", "/infra/domains/default/groups/g1")

		query := getCacheQueryKey(resourceType, d2, enabled)
		tc := gcache.getTypeCache(resourceType)
		tc.data[query] = map[string]*data.StructValue{"g1": sv}
		defer func() { delete(gcache.byTyp, resourceType) }()

		typed, used, attempted, err := TryCacheRead[nsxModel.Group](d2, enabled, nil, "g1", resourceType, nsxModel.GroupBindingType())
		require.NoError(t, err)
		assert.True(t, used)
		assert.True(t, attempted)
		require.NotNil(t, typed)
		assert.Equal(t, "g1", *typed.Id)
	})
}

func TestUnitNsxt_markPostWriteForResourceTypeKey(t *testing.T) {
	enabled := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "config_scope"}}
	disabled := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "disabled"}}

	t.Run("cache disabled: nothing stored", func(t *testing.T) {
		key := postWriteKey("MPWTest", "id-disabled")
		MarkPostWriteForResourceTypeKey("MPWTest", "id-disabled", disabled)
		_, ok := postWriteByKey.Load(key)
		assert.False(t, ok)
	})

	t.Run("empty resourceType or resourceID: nothing stored", func(t *testing.T) {
		MarkPostWriteForResourceTypeKey("", "id", enabled)
		MarkPostWriteForResourceTypeKey("MPWTest", "", enabled)
		_, ok := postWriteByKey.Load(postWriteKey("", "id"))
		assert.False(t, ok)
		_, ok = postWriteByKey.Load(postWriteKey("MPWTest", ""))
		assert.False(t, ok)
	})

	t.Run("enabled with both keys: marker stored", func(t *testing.T) {
		key := postWriteKey("MPWTest", "id-ok")
		MarkPostWriteForResourceTypeKey("MPWTest", "id-ok", enabled)
		_, ok := postWriteByKey.LoadAndDelete(key)
		assert.True(t, ok)
	})
}

func TestUnitNsxt_parseCacheModeString(t *testing.T) {
	cases := []struct {
		raw  string
		want cacheMode
	}{
		{"", cacheDisabled},
		{"disabled", cacheDisabled},
		{"off", cacheDisabled},
		{"  config_scope  ", cacheConfigScoped},
		{"GLOBAL", cacheGlobal},
		{"garbage", cacheDisabled},
	}
	for _, c := range cases {
		t.Run(c.raw, func(t *testing.T) {
			assert.Equal(t, c.want, parseCacheModeString(c.raw))
		})
	}
	// Calling again with the same invalid value exercises the "already warned" branch.
	assert.Equal(t, cacheDisabled, parseCacheModeString("garbage"))
}

func TestUnitNsxt_cacheKeyForResourceID(t *testing.T) {
	t.Run("path-indexed type uses path when set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{"path": "/infra/x"})
		d.SetId("fallback-id")
		assert.Equal(t, "/infra/x", CacheKeyForResourceID(resourceTypeSegmentPort, d))
	})

	t.Run("path-indexed type falls back to id when path is empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
		d.SetId("fallback-id")
		assert.Equal(t, "fallback-id", CacheKeyForResourceID(resourceTypeSegmentPort, d))
	})

	t.Run("non path-indexed type always uses id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{"path": "/infra/x"})
		d.SetId("g1")
		assert.Equal(t, "g1", CacheKeyForResourceID(resourceTypeGroup, d))
	})
}

func TestUnitNsxt_invalidateCacheForResourceType(t *testing.T) {
	const resourceType = "TestOnlyInvalidateType"
	defer delete(gcache.byTyp, resourceType)

	enabled := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "config_scope"}}
	disabled := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "disabled"}}

	tc := gcache.getTypeCache(resourceType)
	sv := groupStructValue(t, "g1", "grp", "/p")
	tc.data["q"] = map[string]*data.StructValue{"g1": sv}

	InvalidateCacheForResourceType(resourceType, disabled)
	assert.Len(t, tc.data, 1, "disabled cache must not be touched")

	InvalidateCacheForResourceType(resourceType, enabled)
	assert.Empty(t, tc.data, "enabled cache must be cleared")
}

func TestUnitNsxt_cacheAwareDataSourceReadByID(t *testing.T) {
	enabled := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "config_scope"}}
	disabled := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "disabled"}}

	t.Run("empty objID short-circuits", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
		_, ok := cacheAwareDataSourceReadByID[nsxModel.Group](d, enabled, nil, "", resourceTypeGroup, nsxModel.GroupBindingType())
		assert.False(t, ok)
	})

	t.Run("cache disabled short-circuits", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
		_, ok := cacheAwareDataSourceReadByID[nsxModel.Group](d, disabled, nil, "g1", resourceTypeGroup, nsxModel.GroupBindingType())
		assert.False(t, ok)
	})

	t.Run("post-write bypass short-circuits", func(t *testing.T) {
		const resourceType = "TestOnlyDSReadByIDBypass"
		postWriteByKey.Store(postWriteKey(resourceType, "g1"), struct{}{})
		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
		_, ok := cacheAwareDataSourceReadByID[nsxModel.Group](d, enabled, nil, "g1", resourceType, nsxModel.GroupBindingType())
		assert.False(t, ok)
		_, stillMarked := postWriteByKey.Load(postWriteKey(resourceType, "g1"))
		assert.False(t, stillMarked)
	})

	t.Run("cache hit populates schema attributes", func(t *testing.T) {
		const resourceType = "TestOnlyDSReadByIDHit"
		defer delete(gcache.byTyp, resourceType)

		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
		sv := groupStructValue(t, "g1", "my-group", "/infra/domains/default/groups/g1")
		query := getCacheQueryKey(resourceType, d, enabled)
		tc := gcache.getTypeCache(resourceType)
		tc.data[query] = map[string]*data.StructValue{"g1": sv}

		typed, ok := cacheAwareDataSourceReadByID[nsxModel.Group](d, enabled, nil, "g1", resourceType, nsxModel.GroupBindingType())
		require.True(t, ok)
		require.NotNil(t, typed)
		assert.Equal(t, "g1", d.Id())
		assert.Equal(t, "my-group", d.Get("display_name"))
		assert.Equal(t, "/infra/domains/default/groups/g1", d.Get("path"))
	})
}

func TestUnitNsxt_cacheAwareResourceRead(t *testing.T) {
	enabled := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "config_scope", contextID: "run-1"}}
	disabled := nsxtClients{CommonConfig: commonProviderConfig{CacheMode: "disabled"}}
	managedTag := nsxModel.Tag{Scope: strPtr(managedDefaultTagScope), Tag: strPtr("run-1")}
	userTag := nsxModel.Tag{Scope: strPtr("env"), Tag: strPtr("dev")}

	t.Run("cache disabled calls backend directly and strips managed tags", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
		d.SetId("g1")
		backendCalls := 0
		backendRead := func() (*nsxModel.Group, error) {
			backendCalls++
			return &nsxModel.Group{Id: strPtr("g1"), Tags: []nsxModel.Tag{managedTag, userTag}}, nil
		}
		patchCalls := 0
		patchFunc := func(obj *nsxModel.Group) error { patchCalls++; return nil }

		obj, used, attempted, err := CacheAwareResourceRead[nsxModel.Group](d, disabled, nil, "g1", resourceTypeGroup, nsxModel.GroupBindingType(), backendRead, patchFunc)
		require.NoError(t, err)
		assert.False(t, used)
		assert.False(t, attempted)
		assert.Equal(t, 1, backendCalls)
		assert.Equal(t, 0, patchCalls)
		require.Len(t, obj.Tags, 1)
		assert.Equal(t, "env", *obj.Tags[0].Scope)
	})

	t.Run("post-write bypass calls backend and patches missing managed tags", func(t *testing.T) {
		const resourceType = "TestOnlyCARRBypass"
		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
		d.SetId("g1")
		postWriteByKey.Store(postWriteKey(resourceType, "g1"), struct{}{})

		backendCalls := 0
		backendRead := func() (*nsxModel.Group, error) {
			backendCalls++
			return &nsxModel.Group{Id: strPtr("g1"), Tags: []nsxModel.Tag{userTag}}, nil
		}
		patchCalls := 0
		patchFunc := func(obj *nsxModel.Group) error { patchCalls++; return nil }

		obj, used, attempted, err := CacheAwareResourceRead[nsxModel.Group](d, enabled, nil, "g1", resourceType, nsxModel.GroupBindingType(), backendRead, patchFunc)
		require.NoError(t, err)
		assert.False(t, used)
		assert.True(t, attempted)
		assert.Equal(t, 1, backendCalls)
		assert.Equal(t, 1, patchCalls, "missing provider-managed tag must trigger a patch")
		// Patched tag is stamped on NSX then immediately stripped from the returned/state object.
		require.Len(t, obj.Tags, 1)
		assert.Equal(t, "env", *obj.Tags[0].Scope)
	})

	t.Run("cache hit returns typed value without calling backend", func(t *testing.T) {
		const resourceType = "TestOnlyCARRHit"
		defer delete(gcache.byTyp, resourceType)

		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
		d.SetId("g1")
		sv := groupStructValue(t, "g1", "my-group", "/p")
		query := getCacheQueryKey(resourceType, d, enabled)
		tc := gcache.getTypeCache(resourceType)
		tc.data[query] = map[string]*data.StructValue{"g1": sv}

		backendCalls := 0
		backendRead := func() (*nsxModel.Group, error) {
			backendCalls++
			return nil, errors.New("must not be called")
		}

		obj, used, attempted, err := CacheAwareResourceRead[nsxModel.Group](d, enabled, nil, "g1", resourceType, nsxModel.GroupBindingType(), backendRead, nil)
		require.NoError(t, err)
		assert.True(t, used)
		assert.True(t, attempted)
		assert.Equal(t, 0, backendCalls)
		require.NotNil(t, obj)
		assert.Equal(t, "g1", *obj.Id)
	})

	t.Run("cache miss falls through and surfaces backend errors", func(t *testing.T) {
		const resourceType = "TestOnlyCARRMissError"
		defer delete(gcache.byTyp, resourceType)

		d := schema.TestResourceDataRaw(t, cacheTestSchema(), map[string]interface{}{})
		d.SetId("g1")
		// Empty bucket: getQueryResult returns errCacheUseBackendDirect without a backend call.
		tc := gcache.getTypeCache(resourceType)
		query := getCacheQueryKey(resourceType, d, enabled)
		tc.data[query] = map[string]*data.StructValue{}

		backendRead := func() (*nsxModel.Group, error) { return nil, errors.New("backend down") }

		obj, used, attempted, err := CacheAwareResourceRead[nsxModel.Group](d, enabled, nil, "g1", resourceType, nsxModel.GroupBindingType(), backendRead, nil)
		require.Error(t, err)
		assert.Equal(t, "backend down", err.Error())
		assert.False(t, used)
		assert.True(t, attempted)
		assert.Nil(t, obj)
	})
}
