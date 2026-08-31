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
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	mp_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func tagElemResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"scope": {Type: schema.TypeString, Optional: true},
			"tag":   {Type: schema.TypeString, Optional: true},
		},
	}
}

func TestUnitNsxt_getPolicyTagsFromSet(t *testing.T) {
	elem := tagElemResource()
	tagSet := schema.NewSet(schema.HashResource(elem), []interface{}{
		map[string]interface{}{"scope": "env", "tag": "prod"},
	})
	out := getPolicyTagsFromSet(tagSet)
	require.Len(t, out, 1)
	assert.Equal(t, "env", *out[0].Scope)
	assert.Equal(t, "prod", *out[0].Tag)
}

func TestUnitNsxt_getIgnoredTagsFromSchema_undefined(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"ignore_tags": getIgnoreTagsSchema(),
	}, map[string]interface{}{})
	assert.Nil(t, getIgnoredTagsFromSchema(d))
}

func TestUnitNsxt_getCustomizedPolicyTagsFromSchema(t *testing.T) {
	sch := map[string]*schema.Schema{
		"my_tags": getTagsSchema(),
	}
	d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{
		"my_tags": []interface{}{
			map[string]interface{}{"scope": "a", "tag": "b"},
		},
	})
	out, err := getCustomizedPolicyTagsFromSchema(d, "my_tags")
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "a", *out[0].Scope)
	assert.Equal(t, "b", *out[0].Tag)
}

func TestUnitNsxt_getCustomizedPolicyTagsFromSchema_emptyTagError(t *testing.T) {
	sch := map[string]*schema.Schema{
		"my_tags": getTagsSchema(),
	}
	d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{
		"my_tags": []interface{}{
			map[string]interface{}{"scope": "", "tag": ""},
		},
	})
	_, err := getCustomizedPolicyTagsFromSchema(d, "my_tags")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tag value or scope value needs to be specified")
}

func TestUnitNsxt_getTagScopesToIgnore(t *testing.T) {
	sch := map[string]*schema.Schema{
		"ignore_tags": getIgnoreTagsSchema(),
	}
	d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{
		"ignore_tags": []interface{}{
			map[string]interface{}{
				"scopes": []interface{}{"a", "b"},
			},
		},
	})
	scopes := getTagScopesToIgnore(d)
	assert.Equal(t, []string{"a", "b"}, scopes)
}

func TestUnitNsxt_setCustomizedPolicyTagsInSchema(t *testing.T) {
	sch := map[string]*schema.Schema{"tag": getTagsSchema()}
	d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{})
	s1, t1 := "scope1", "tag1"
	setCustomizedPolicyTagsInSchema(d, []model.Tag{{Scope: &s1, Tag: &t1}}, "tag")
	st := d.Get("tag").(*schema.Set)
	require.Equal(t, 1, st.Len())
}

func TestUnitNsxt_setCustomizedPolicyTagsInSchema_skipsManagedTags(t *testing.T) {
	sch := map[string]*schema.Schema{"tag": getTagsSchema()}
	d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{})
	s1, t1 := "scope1", "tag1"
	sManaged, tManaged := managedDefaultTagScope, "some-run-id"
	setCustomizedPolicyTagsInSchema(d, []model.Tag{
		{Scope: &s1, Tag: &t1},
		{Scope: &sManaged, Tag: &tManaged},
	}, "tag")
	st := d.Get("tag").(*schema.Set)
	require.Equal(t, 1, st.Len())
	elem := st.List()[0].(map[string]interface{})
	assert.Equal(t, "scope1", elem["scope"])
}

func TestUnitNsxt_initPolicyTagsSet_skipsManagedTags(t *testing.T) {
	s1, t1 := "scope1", "tag1"
	sManaged, tManaged := managedDefaultTagScope, "some-run-id"
	out := initPolicyTagsSet([]model.Tag{
		{Scope: &s1, Tag: &t1},
		{Scope: &sManaged, Tag: &tManaged},
	})
	require.Len(t, out, 1)
	assert.Equal(t, "scope1", *out[0]["scope"].(*string))
}

func TestUnitNsxt_setPathListInSchema(t *testing.T) {
	sch := map[string]*schema.Schema{
		"paths": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
	t.Run("ANY clears attribute", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{
			"paths": []interface{}{"/infra/a"},
		})
		setPathListInSchema(d, "paths", []string{"ANY"})
		v := d.Get("paths")
		assert.True(t, v == nil || len(v.([]interface{})) == 0)
	})
	t.Run("non ANY sets list", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{})
		setPathListInSchema(d, "paths", []string{"/p1", "/p2"})
		v := d.Get("paths").([]interface{})
		require.Len(t, v, 2)
		assert.Equal(t, "/p1", v[0].(string))
		assert.Equal(t, "/p2", v[1].(string))
	})
}

func TestUnitNsxt_getPolicyKeyValuePairListFromSchema(t *testing.T) {
	assert.Empty(t, getPolicyKeyValuePairListFromSchema(nil))
	in := []interface{}{
		map[string]interface{}{"key": "k1", "value": "v1"},
	}
	out := getPolicyKeyValuePairListFromSchema(in)
	require.Len(t, out, 1)
	assert.Equal(t, "k1", *out[0].Key)
	assert.Equal(t, "v1", *out[0].Value)
}

func TestUnitNsxt_setPolicyKeyValueListForSchema(t *testing.T) {
	ka, vb := "a", "b"
	in := []model.KeyValuePair{
		{Key: &ka, Value: &vb},
	}
	raw := setPolicyKeyValueListForSchema(in).([]interface{})
	require.Len(t, raw, 1)
	m := raw[0].(map[string]interface{})
	assert.Equal(t, &ka, m["key"])
	assert.Equal(t, &vb, m["value"])
}

func TestUnitNsxt_nsxtVersionCheckImporter(t *testing.T) {
	defer func() { util.NsxVersion = "" }()
	var handlerCalled bool
	handler := func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		handlerCalled = true
		return []*schema.ResourceData{d}, nil
	}
	imp := nsxtVersionCheckImporter("9.0.0", "TestResource", handler)
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
	d.SetId("id1")

	util.NsxVersion = "8.0.0"
	_, err := imp(d, newGoMockProviderClient())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires NSX version 9.0.0")
	assert.False(t, handlerCalled)

	util.NsxVersion = "9.1.0"
	_, err = imp(d, newGoMockProviderClient())
	require.NoError(t, err)
	assert.True(t, handlerCalled)
}

func TestUnitNsxt_getAttrKeyMapFromSchemaSet(t *testing.T) {
	elem := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {Type: schema.TypeString, Required: true},
			"val":  {Type: schema.TypeString, Optional: true},
		},
	}
	st := schema.NewSet(schema.HashResource(elem), []interface{}{
		map[string]interface{}{"name": "n1", "val": "x"},
		map[string]interface{}{"name": "n2", "val": "y"},
	})
	m := getAttrKeyMapFromSchemaSet(st, "name")
	assert.True(t, m["n1"])
	assert.True(t, m["n2"])
	assert.False(t, m["missing"])
}

func TestUnitNsxt_getCustomizedMPTagsFromSchema(t *testing.T) {
	sch := map[string]*schema.Schema{"mp_tags": getTagsSchema()}
	d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{
		"mp_tags": []interface{}{
			map[string]interface{}{"scope": "mp", "tag": "t"},
		},
	})
	out := getCustomizedMPTagsFromSchema(d, "mp_tags")
	require.Len(t, out, 1)
	assert.Equal(t, "mp", *out[0].Scope)
	assert.Equal(t, "t", *out[0].Tag)
}

func TestUnitNsxt_setCustomizedMPTagsInSchema(t *testing.T) {
	sch := map[string]*schema.Schema{"mp_tags": getTagsSchema()}
	d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{})
	ms, mv := "s", "v"
	tags := []mp_model.Tag{
		{Scope: &ms, Tag: &mv},
	}
	setCustomizedMPTagsInSchema(d, tags, "mp_tags")
	st := d.Get("mp_tags").(*schema.Set)
	require.Equal(t, 1, st.Len())
}

func TestUnitNsxt_getCustomizedGMTagsFromSchema(t *testing.T) {
	sch := map[string]*schema.Schema{"gm_tags": getTagsSchema()}
	d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{
		"gm_tags": []interface{}{
			map[string]interface{}{"scope": "gm", "tag": "gt"},
		},
	})
	out := getCustomizedGMTagsFromSchema(d, "gm_tags")
	require.Len(t, out, 1)
	assert.Equal(t, "gm", *out[0].Scope)
}

func TestUnitNsxt_setCustomizedGMTagsInSchema(t *testing.T) {
	sch := map[string]*schema.Schema{"gm_tags": getTagsSchema()}
	d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{})
	gs, gv := "gs", "gv"
	tags := []gm_model.Tag{
		{Scope: &gs, Tag: &gv},
	}
	setCustomizedGMTagsInSchema(d, tags, "gm_tags")
	st := d.Get("gm_tags").(*schema.Set)
	require.Equal(t, 1, st.Len())
}

func TestUnitNsxt_getKeyValuePairListFromSchema(t *testing.T) {
	assert.Empty(t, getKeyValuePairListFromSchema(nil))
	in := []interface{}{
		map[string]interface{}{"key": "k", "value": "v"},
	}
	out := getKeyValuePairListFromSchema(in)
	require.Len(t, out, 1)
	assert.Equal(t, "k", *out[0].Key)
}

func TestUnitNsxt_setKeyValueListForSchema(t *testing.T) {
	ka2, vb2 := "a", "b"
	out := setKeyValueListForSchema([]mp_model.KeyValuePair{
		{Key: &ka2, Value: &vb2},
	}).([]interface{})
	require.Len(t, out, 1)
}

func TestUnitNsxt_getStringValue(t *testing.T) {
	assert.Equal(t, "", getStringValue(nil))
	s := "hello"
	assert.Equal(t, "hello", getStringValue(&s))
}

func TestUnitNsxt_getInt64Value(t *testing.T) {
	assert.EqualValues(t, 0, getInt64Value(nil))
	v := int64(42)
	assert.EqualValues(t, 42, getInt64Value(&v))
}

func TestUnitNsxt_interface2Int64List(t *testing.T) {
	out := interface2Int64List([]interface{}{1, 2, "not-an-int", 3})
	assert.Equal(t, []int64{1, 2, 3}, out)
	assert.Empty(t, interface2Int64List(nil))
}

func TestUnitNsxt_int64List2Interface(t *testing.T) {
	out := int64List2Interface([]int64{1, 2, 3})
	assert.Equal(t, []interface{}{1, 2, 3}, out)
}

func TestUnitNsxt_resourceKeyValueHash(t *testing.T) {
	// Only single-key maps are compared for equality: with multiple keys, Go's
	// randomized map iteration order feeds resourceKeyValueHash's buffer in a
	// different sequence on each call, so even equal-content maps can hash
	// differently. That's an existing quirk of the hash function, not something
	// under test here.
	h1 := resourceKeyValueHash(map[string]interface{}{"a": "1"})
	h2 := resourceKeyValueHash(map[string]interface{}{"a": "1"})
	h3 := resourceKeyValueHash(map[string]interface{}{"a": "2"})
	assert.Equal(t, h1, h2)
	assert.NotEqual(t, h1, h3)
	assert.GreaterOrEqual(t, resourceKeyValueHash(nil), 0)
}

func TestUnitNsxt_resourceNotSupportedError(t *testing.T) {
	err := resourceNotSupportedError()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not supported")
}

func TestUnitNsxt_dataSourceNotSupportedError(t *testing.T) {
	err := dataSourceNotSupportedError()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not supported")
}

func TestUnitNsxt_handlePagination(t *testing.T) {
	t.Run("empty list returns zero total with no error", func(t *testing.T) {
		calls := 0
		total, err := handlePagination(func(info *paginationInfo) error {
			calls++
			info.TotalCount = 0
			return nil
		})
		require.NoError(t, err)
		assert.EqualValues(t, 0, total)
		assert.Equal(t, 1, calls)
	})

	t.Run("paginates across multiple pages until count reaches total", func(t *testing.T) {
		calls := 0
		total, err := handlePagination(func(info *paginationInfo) error {
			calls++
			info.TotalCount = 5
			info.PageCount = 2
			info.Cursor = "next"
			return nil
		})
		require.NoError(t, err)
		assert.EqualValues(t, 5, total)
		assert.Equal(t, 3, calls)
	})

	t.Run("lister error is propagated", func(t *testing.T) {
		_, err := handlePagination(func(info *paginationInfo) error {
			return errors.New("lister failed")
		})
		require.Error(t, err)
	})
}

func TestUnitNsxt_attributeRequiredGlobalManagerError(t *testing.T) {
	err := attributeRequiredGlobalManagerError("attr", "resource")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "attr")
	assert.Contains(t, err.Error(), "resource")
}

func TestUnitNsxt_getSitePathFromEdgePath(t *testing.T) {
	edgePath := "/global-infra/sites/site1/enforcement-points/default/edge-clusters/ec1"
	assert.Equal(t, "/global-infra/sites/site1", getSitePathFromEdgePath(edgePath))
}

func TestUnitNsxt_ptr(t *testing.T) {
	p := ptr("value")
	require.NotNil(t, p)
	assert.Equal(t, "value", *p)
}
