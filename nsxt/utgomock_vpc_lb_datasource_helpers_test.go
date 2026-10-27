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
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"go.uber.org/mock/gomock"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

// vpcLbProfileDataSourceSpec describes a VPC LB profile data source (app profile or monitor) under test.
type vpcLbProfileDataSourceSpec struct {
	newDataSource func() *schema.Resource
	setupMock     func(ctrl *gomock.Controller) *vpcLbProfileMock
	// structValue builds a profile of the given resource type
	structValue func(t *testing.T, id, displayName, resourceType string) *data.StructValue
	// resourceType / typeName: a profile type and its data source "type" value
	resourceType string
	typeName     string
	// otherTypeName is a different valid data source "type" value
	otherTypeName string
	apiSegment    string
}

func runVpcLbProfileDataSourceTests(t *testing.T, spec vpcLbProfileDataSourceSpec) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	const id = "ds-profile-1"
	const name = "ds-profile-name"
	sessionClient := newGoMockProviderClient()

	setup := func(t *testing.T) (*vpcLbProfileMock, func()) {
		ctrl := gomock.NewController(t)
		m := spec.setupMock(ctrl)
		return m, func() { m.restore(); ctrl.Finish() }
	}
	read := func(t *testing.T, raw map[string]interface{}) (*schema.ResourceData, error) {
		ds := spec.newDataSource()
		d := schema.TestResourceDataRaw(t, ds.Schema, vpcLbWithContext(raw))
		return d, ds.Read(d, sessionClient)
	}
	sv := func(t *testing.T, id, name string) *data.StructValue {
		return spec.structValue(t, id, name, spec.resourceType)
	}

	t.Run("Schema", func(t *testing.T) {
		ds := spec.newDataSource()
		require.NotNil(t, ds.Schema["context"])
		assert.True(t, ds.Schema["context"].Required)
	})

	t.Run("Fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.2.0"
		defer func() { util.NsxVersion = "9.0.0" }()
		_, cleanup := setup(t)
		defer cleanup()
		_, err := read(t, map[string]interface{}{"id": id})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "requires NSX version 9.0.0 or higher")
	})

	t.Run("Read by ID", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.get(id, sv(t, id, name), nil)
		d, err := read(t, map[string]interface{}{"id": id})
		require.NoError(t, err)
		assert.Equal(t, id, d.Id())
		assert.Equal(t, name, d.Get("display_name"))
		assert.Equal(t, spec.typeName, d.Get("type"))
		assert.Contains(t, d.Get("path"), "/"+id)
	})

	t.Run("Read by ID with matching type", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.get(id, sv(t, id, name), nil)
		d, err := read(t, map[string]interface{}{"id": id, "type": spec.typeName})
		require.NoError(t, err)
		assert.Equal(t, id, d.Id())
	})

	t.Run("Read by ID with mismatched type is not found", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.get(id, sv(t, id, name), nil)
		_, err := read(t, map[string]interface{}{"id": id, "type": spec.otherTypeName})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("Read by ID with unknown resource type", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.get(id, spec.structValue(t, id, name, "LBUnknownProfile"), nil)
		_, err := read(t, map[string]interface{}{"id": id})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while converting")
	})

	t.Run("Read by ID API error", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.get(id, nil, vapiErrors.InternalServerError{})
		_, err := read(t, map[string]interface{}{"id": id})
		require.Error(t, err)
	})

	t.Run("Read by ID not found", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.get(id, nil, vapiErrors.NotFound{})
		_, err := read(t, map[string]interface{}{"id": id})
		require.Error(t, err)
	})

	t.Run("Read by perfect name match", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.list([]*data.StructValue{sv(t, "other", name+"-longer"), sv(t, id, name)}, nil)
		d, err := read(t, map[string]interface{}{"display_name": name})
		require.NoError(t, err)
		assert.Equal(t, id, d.Id())
		assert.Equal(t, name, d.Get("display_name"))
	})

	t.Run("Read by prefix name match", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.list([]*data.StructValue{sv(t, id, name)}, nil)
		d, err := read(t, map[string]interface{}{"display_name": "ds-profile"})
		require.NoError(t, err)
		assert.Equal(t, id, d.Id())
	})

	t.Run("Read by name multiple perfect matches", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.list([]*data.StructValue{sv(t, "dup-1", "dup"), sv(t, "dup-2", "dup")}, nil)
		_, err := read(t, map[string]interface{}{"display_name": "dup"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("Read by name multiple prefix matches", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.list([]*data.StructValue{sv(t, "p-1", "prefix-a"), sv(t, "p-2", "prefix-b")}, nil)
		_, err := read(t, map[string]interface{}{"display_name": "prefix"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("Read by name not found", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.list(nil, nil)
		_, err := read(t, map[string]interface{}{"display_name": "nonexistent"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Read by name filters by type", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.list([]*data.StructValue{sv(t, id, name)}, nil)
		_, err := read(t, map[string]interface{}{"display_name": name, "type": spec.otherTypeName})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Read by name and matching type", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.list([]*data.StructValue{sv(t, id, name)}, nil)
		d, err := read(t, map[string]interface{}{"display_name": name, "type": spec.typeName})
		require.NoError(t, err)
		assert.Equal(t, id, d.Id())
	})

	t.Run("Read by type only", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.list([]*data.StructValue{sv(t, id, name)}, nil)
		d, err := read(t, map[string]interface{}{"type": spec.typeName})
		require.NoError(t, err)
		assert.Equal(t, id, d.Id())
	})

	t.Run("Read by name with unknown resource type in list", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.list([]*data.StructValue{spec.structValue(t, id, name, "LBUnknownProfile")}, nil)
		_, err := read(t, map[string]interface{}{"display_name": name})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while converting")
		assert.Contains(t, err.Error(), name)
	})

	t.Run("Read by name list error", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.list(nil, vapiErrors.InternalServerError{})
		_, err := read(t, map[string]interface{}{"display_name": name})
		require.Error(t, err)
	})

	t.Run("Read without id, name or type", func(t *testing.T) {
		_, cleanup := setup(t)
		defer cleanup()
		ds := spec.newDataSource()
		d := schema.TestResourceDataRaw(t, ds.Schema, vpcLbWithContext(map[string]interface{}{}))
		require.NoError(t, d.Set("type", ""))
		err := ds.Read(d, sessionClient)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining")
	})
}
