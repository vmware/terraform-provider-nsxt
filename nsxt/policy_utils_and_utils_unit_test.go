//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
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
