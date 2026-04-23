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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	gmModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func i64(v int64) *int64 { return &v }

func str(s string) *string { return &s }

func policyResourceToStructValue(t *testing.T, pr gmModel.PolicyResource) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(pr, gmModel.PolicyResourceBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

// seqQueryListClient returns configured SearchResponse values in order for each List call.
type seqQueryListClient struct {
	responses []nsxModel.SearchResponse
	errs      []error
	call      int
}

func (s *seqQueryListClient) List(queryParam string, cursorParam *string, includedFieldsParam *string, pageSizeParam *int64, sortAscendingParam *bool, sortByParam *string) (nsxModel.SearchResponse, error) {
	i := s.call
	s.call++
	if s.errs != nil && i < len(s.errs) && s.errs[i] != nil {
		return nsxModel.SearchResponse{}, s.errs[i]
	}
	if i < len(s.responses) {
		return s.responses[i], nil
	}
	return nsxModel.SearchResponse{}, nil
}

func setupCliQueryClientStub(t *testing.T, stub queryListClient) func() {
	t.Helper()
	orig := cliQueryClient
	cliQueryClient = func(_ utl.SessionContext, _ client.Connector) queryListClient {
		return stub
	}
	return func() { cliQueryClient = orig }
}

func TestUnitNsxt_escapeSpecialCharacters(t *testing.T) {
	assert.Equal(t, "plain", escapeSpecialCharacters("plain"))
	assert.Equal(t, `foo\(bar\)`, escapeSpecialCharacters("foo(bar)"))
	assert.Equal(t, `a\+b`, escapeSpecialCharacters("a+b"))
}

func TestUnitNsxt_buildPolicyResourcesQuery(t *testing.T) {
	q := "base"
	add := "extra"
	buildPolicyResourcesQuery(&q, &add)
	assert.Equal(t, "base AND extra", q)

	q2 := "only"
	empty := ""
	buildPolicyResourcesQuery(&q2, &empty)
	assert.Equal(t, "only", q2)

	q3 := "x"
	buildPolicyResourcesQuery(&q3, nil)
	assert.Equal(t, "x", q3)
}

func TestUnitNsxt_searchByContextInvalidClientType(t *testing.T) {
	ctx := utl.SessionContext{ClientType: utl.ClientType(99)}
	defer setupCliQueryClientStub(t, &seqQueryListClient{})()
	_, err := searchByContext(nil, ctx, "q", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid ClientType")
}

func TestUnitNsxt_searchLMListError(t *testing.T) {
	stub := &seqQueryListClient{
		errs: []error{errors.New("boom")},
	}
	defer setupCliQueryClientStub(t, stub)()
	ctx := utl.SessionContext{ClientType: utl.Local}
	_, err := searchLM(ctx, nil, "resource_type:Foo AND marked_for_delete:false AND path:\\/infra*")
	require.Error(t, err)
	assert.Equal(t, "boom", err.Error())
}

func TestUnitNsxt_searchLMPolicyResourcesGlobalInfraPathWhenIsGlobal(t *testing.T) {
	res := nsxModel.SearchResponse{Results: []*data.StructValue{}, ResultCount: i64(0)}
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{res}}
	defer setupCliQueryClientStub(t, stub)()
	ctx := utl.SessionContext{ClientType: utl.Local, FromGlobal: true}
	_, err := searchLMPolicyResources(ctx, nil, "id:x AND marked_for_delete:false", true)
	require.NoError(t, err)
}

func TestUnitNsxt_searchMultitenancyVPCPathSuffix(t *testing.T) {
	res := nsxModel.SearchResponse{Results: []*data.StructValue{}, ResultCount: i64(0)}
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{res}}
	defer setupCliQueryClientStub(t, stub)()
	ctx := utl.SessionContext{ClientType: utl.VPC, ProjectID: "p1", VPCID: "vpc1"}
	_, err := searchMultitenancyResources(nil, ctx, "resource_type:Z AND marked_for_delete:false")
	require.NoError(t, err)
}

func TestUnitNsxt_listPolicyResourcesByNsxID(t *testing.T) {
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("nid"), DisplayName: str("n"), Path: str("/p"), ResourceType: str("PolicyEdgeNode"),
	})
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
		Results: []*data.StructValue{sv}, ResultCount: i64(1),
	}}}
	defer setupCliQueryClientStub(t, stub)()
	nsxID := "nsx-uuid-1"
	add := ""
	ctx := utl.SessionContext{ClientType: utl.Local}
	out, err := listPolicyResourcesByNsxID(nil, ctx, &nsxID, &add)
	require.NoError(t, err)
	assert.Len(t, out, 1)
}

func TestUnitNsxt_listInventoryResourcesByAnyFieldAndType(t *testing.T) {
	rt := "FooType"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("i1"), DisplayName: str("n"), Path: str("/p"), ResourceType: &rt,
	})
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
		Results: []*data.StructValue{sv}, ResultCount: i64(1),
	}}}
	defer setupCliQueryClientStub(t, stub)()
	ctx := utl.SessionContext{ClientType: utl.Local}
	_, err := listInventoryResourcesByAnyFieldAndType(nil, ctx, "host_id:abc", rt, nil)
	require.NoError(t, err)
}

func TestUnitNsxt_searchLMPolicyResourcesAppendsInfraPathWhenNotGlobal(t *testing.T) {
	res := nsxModel.SearchResponse{
		Results:     []*data.StructValue{},
		ResultCount: i64(0),
		Cursor:      nil,
	}
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{res}}
	defer setupCliQueryClientStub(t, stub)()
	ctx := utl.SessionContext{ClientType: utl.Local, FromGlobal: false}
	_, err := searchLMPolicyResources(ctx, nil, "id:x AND marked_for_delete:false", false)
	require.NoError(t, err)
}

func TestUnitNsxt_searchGMAppendsGlobalInfraPath(t *testing.T) {
	res := nsxModel.SearchResponse{
		Results:     []*data.StructValue{},
		ResultCount: i64(0),
		Cursor:      nil,
	}
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{res}}
	defer setupCliQueryClientStub(t, stub)()
	ctx := utl.SessionContext{ClientType: utl.Global}
	_, err := searchGMPolicyResources(ctx, nil, "id:y AND marked_for_delete:false")
	require.NoError(t, err)
}

func TestUnitNsxt_searchMultitenancyResourcesAppendsProjectPath(t *testing.T) {
	res := nsxModel.SearchResponse{Results: []*data.StructValue{}, ResultCount: i64(0)}
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{res}}
	defer setupCliQueryClientStub(t, stub)()
	ctx := utl.SessionContext{ClientType: utl.Multitenancy, ProjectID: "proj1"}
	_, err := searchMultitenancyResources(nil, ctx, "resource_type:Z AND marked_for_delete:false")
	require.NoError(t, err)
}

func TestUnitNsxt_searchLMPagination(t *testing.T) {
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("id1"), DisplayName: str("n"), Path: str("/infra/x"), ResourceType: str("Tier1"),
	})
	page1 := nsxModel.SearchResponse{
		Results:     []*data.StructValue{sv},
		ResultCount: i64(2),
		Cursor:      str("c1"),
	}
	page2 := nsxModel.SearchResponse{
		Results:     []*data.StructValue{sv},
		ResultCount: i64(2),
		Cursor:      nil,
	}
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{page1, page2}}
	defer setupCliQueryClientStub(t, stub)()
	ctx := utl.SessionContext{ClientType: utl.Local}
	out, err := searchLM(ctx, nil, "q")
	require.NoError(t, err)
	assert.Len(t, out, 2)
	assert.Equal(t, 2, stub.call)
}

func TestUnitNsxt_policyDataSourceResourceFilterAndSet(t *testing.T) {
	rt := "Tier1"
	ds := schema.Resource{Schema: map[string]*schema.Schema{
		"id":           {Type: schema.TypeString, Optional: true},
		"display_name": {Type: schema.TypeString, Optional: true},
		"description":  {Type: schema.TypeString, Optional: true, Computed: true},
		"path":         {Type: schema.TypeString, Optional: true, Computed: true},
	}}

	t.Run("skips wrong resource type", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("a"), DisplayName: str("x"), Path: str("/p"), ResourceType: str("Other"),
		})
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "x"})
		_, err := policyDataSourceResourceFilterAndSet(d, []*data.StructValue{sv}, rt)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by display_name perfect match", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("tid"), DisplayName: str("my-tier1"), Path: str("/infra/t0"), Description: str("d"),
			ResourceType: &rt,
		})
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "my-tier1"})
		_, err := policyDataSourceResourceFilterAndSet(d, []*data.StructValue{sv}, rt)
		require.NoError(t, err)
		assert.Equal(t, "tid", d.Id())
	})

	t.Run("multiple perfect matches by name", func(t *testing.T) {
		sv1 := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("1"), DisplayName: str("dup"), Path: str("/p1"), ResourceType: &rt,
		})
		sv2 := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("2"), DisplayName: str("dup"), Path: str("/p2"), ResourceType: &rt,
		})
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "dup"})
		_, err := policyDataSourceResourceFilterAndSet(d, []*data.StructValue{sv1, sv2}, rt)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("prefix match single", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("pfx"), DisplayName: str("prefix-name"), Path: str("/p"), ResourceType: &rt,
		})
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "prefix"})
		_, err := policyDataSourceResourceFilterAndSet(d, []*data.StructValue{sv}, rt)
		require.NoError(t, err)
		assert.Equal(t, "pfx", d.Id())
	})

	t.Run("multiple prefix matches", func(t *testing.T) {
		sv1 := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("a"), DisplayName: str("pre-one"), Path: str("/p1"), ResourceType: &rt,
		})
		sv2 := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("b"), DisplayName: str("pre-two"), Path: str("/p2"), ResourceType: &rt,
		})
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "pre"})
		_, err := policyDataSourceResourceFilterAndSet(d, []*data.StructValue{sv1, sv2}, rt)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("not found by id when no matching type", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("want"), DisplayName: str("n"), Path: str("/p"), ResourceType: str("Other"),
		})
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "want"})
		_, err := policyDataSourceResourceFilterAndSet(d, []*data.StructValue{sv}, rt)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})
}

func TestUnitNsxt_policyDataSourceResourceReadWithValidation(t *testing.T) {
	rt := "Tier1"
	ds := schema.Resource{Schema: map[string]*schema.Schema{
		"id": {Type: schema.TypeString, Optional: true}, "display_name": {Type: schema.TypeString, Optional: true},
		"description": {Type: schema.TypeString, Optional: true, Computed: true},
		"path":        {Type: schema.TypeString, Optional: true, Computed: true},
	}}

	t.Run("policyDataSourceResourceRead delegates to WithValidation", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("rid2"), DisplayName: str("nm2"), Path: str("/infra/r2"), ResourceType: &rt,
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "rid2"})
		_, err := policyDataSourceResourceRead(d, nil, utl.SessionContext{ClientType: utl.Local, FromGlobal: false}, rt, nil)
		require.NoError(t, err)
		assert.Equal(t, "rid2", d.Id())
	})

	t.Run("missing id and display_name when validation on", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
		_, err := policyDataSourceResourceReadWithValidation(d, nil, utl.SessionContext{ClientType: utl.Local}, rt, nil, true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "No 'id' or 'display_name'")
	})

	t.Run("success by id via search", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("rid"), DisplayName: str("nm"), Path: str("/infra/r"), ResourceType: &rt,
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "rid"})
		_, err := policyDataSourceResourceReadWithValidation(d, nil, utl.SessionContext{ClientType: utl.Local, FromGlobal: false}, rt, nil, true)
		require.NoError(t, err)
		assert.Equal(t, "rid", d.Id())
	})

	t.Run("success by display_name via listPolicyResourcesByNameAndType", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("byname"), DisplayName: str("exact-name"), Path: str("/infra/e"), ResourceType: &rt,
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "exact-name"})
		_, err := policyDataSourceResourceReadWithValidation(d, nil, utl.SessionContext{ClientType: utl.Local, FromGlobal: false}, rt, nil, true)
		require.NoError(t, err)
		assert.Equal(t, "byname", d.Id())
	})

	t.Run("search error wraps message for display_name", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search-fail")}}
		defer setupCliQueryClientStub(t, stub)()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "x"})
		_, err := policyDataSourceResourceReadWithValidation(d, nil, utl.SessionContext{ClientType: utl.Local, FromGlobal: false}, rt, nil, true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error getting resource")
		assert.Contains(t, err.Error(), "name")
	})
}

func TestUnitNsxt_policyDataSourceReadWithCustomField(t *testing.T) {
	rt := "SomeType"
	ds := schema.Resource{Schema: map[string]*schema.Schema{
		"id":           {Type: schema.TypeString, Optional: true, Computed: true},
		"display_name": {Type: schema.TypeString, Optional: true, Computed: true},
		"description":  {Type: schema.TypeString, Optional: true, Computed: true},
		"path":         {Type: schema.TypeString, Optional: true, Computed: true},
	}}
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("c1"), DisplayName: str("cn"), Path: str("/p"), ResourceType: &rt,
	})
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
		Results: []*data.StructValue{sv}, ResultCount: i64(1),
	}}}
	defer setupCliQueryClientStub(t, stub)()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	_, err := policyDataSourceReadWithCustomField(d, nil, utl.SessionContext{ClientType: utl.Local}, rt, map[string]string{"foo": "bar"})
	require.NoError(t, err)
	assert.Equal(t, "c1", d.Id())
}

func TestUnitNsxt_policyDataSourceReadWithCustomFieldWrongType(t *testing.T) {
	rt := "WantedType"
	wrongRT := "OtherType"
	ds := schema.Resource{Schema: map[string]*schema.Schema{
		"id":           {Type: schema.TypeString, Optional: true, Computed: true},
		"display_name": {Type: schema.TypeString, Optional: true, Computed: true},
		"description":  {Type: schema.TypeString, Optional: true, Computed: true},
		"path":         {Type: schema.TypeString, Optional: true, Computed: true},
	}}
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("c1"), DisplayName: str("cn"), Path: str("/p"), ResourceType: &wrongRT,
	})
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
		Results: []*data.StructValue{sv}, ResultCount: i64(1),
	}}}
	defer setupCliQueryClientStub(t, stub)()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	_, err := policyDataSourceReadWithCustomField(d, nil, utl.SessionContext{ClientType: utl.Local}, rt, map[string]string{"k": "v"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "found no resources")
}

func TestUnitNsxt_policyDataSourceReadWithCustomFieldNotExactlyOne(t *testing.T) {
	rt := "SomeType"
	ds := schema.Resource{Schema: map[string]*schema.Schema{
		"id":           {Type: schema.TypeString, Optional: true, Computed: true},
		"display_name": {Type: schema.TypeString, Optional: true, Computed: true},
		"description":  {Type: schema.TypeString, Optional: true, Computed: true},
		"path":         {Type: schema.TypeString, Optional: true, Computed: true},
	}}
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("a"), DisplayName: str("n"), Path: str("/p"), ResourceType: &rt,
	})
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
		Results: []*data.StructValue{sv, sv}, ResultCount: i64(2),
	}}}
	defer setupCliQueryClientStub(t, stub)()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
	_, err := policyDataSourceReadWithCustomField(d, nil, utl.SessionContext{ClientType: utl.Local}, rt, map[string]string{"k": "v"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "none or multiple")
}

func TestUnitNsxt_policyDataSourceCreateMap(t *testing.T) {
	rt := "Tier1"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("m1"), DisplayName: str("MapName"), Path: str("/p"), ResourceType: &rt,
	})
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
		Results: []*data.StructValue{sv}, ResultCount: i64(1),
	}}}
	defer setupCliQueryClientStub(t, stub)()
	out := map[string]string{}
	err := policyDataSourceCreateMap(nil, utl.SessionContext{ClientType: utl.Local}, rt, out, nil)
	require.NoError(t, err)
	assert.Equal(t, "MapName", out["m1"])

	t.Run("list error", func(t *testing.T) {
		stub2 := &seqQueryListClient{errs: []error{errors.New("map-fail")}}
		defer setupCliQueryClientStub(t, stub2)()
		out2 := map[string]string{}
		err := policyDataSourceCreateMap(nil, utl.SessionContext{ClientType: utl.Local}, rt, out2, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch")
	})
}

func TestUnitNsxt_listInventoryResourcesByNameAndType(t *testing.T) {
	rt := "HostTransportNode"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("h1"), DisplayName: str("host-a"), Path: str("/p"), ResourceType: &rt,
	})
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
		Results: []*data.StructValue{sv}, ResultCount: i64(1),
	}}}
	defer setupCliQueryClientStub(t, stub)()
	ctx := utl.SessionContext{ClientType: utl.Local}
	_, err := listInventoryResourcesByNameAndType(nil, ctx, "host", rt, nil)
	require.NoError(t, err)
}

func TestUnitNsxt_initNSXVersionVMCUsesCliQueryClient(t *testing.T) {
	stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{Results: nil, ResultCount: i64(0)}}}
	defer setupCliQueryClientStub(t, stub)()
	clients := newGoMockProviderClient()
	initNSXVersionVMC(clients)
	assert.Equal(t, 1, stub.call)
}
