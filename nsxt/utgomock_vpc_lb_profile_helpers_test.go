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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	apivpcs "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	vpcsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

// Shared helpers for the VPC LB application/monitor profile unit tests (resources and data sources).

const (
	vpcLbTestOrg     = "default"
	vpcLbTestProject = "proj1"
	vpcLbTestVPC     = "vpc1"
)

func vpcLbTestContext() []interface{} {
	return []interface{}{map[string]interface{}{
		"project_id": vpcLbTestProject,
		"vpc_id":     vpcLbTestVPC,
	}}
}

// vpcLbWithContext returns a copy of the given raw config extended with the VPC context.
func vpcLbWithContext(src map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{}, len(src)+1)
	for k, v := range src {
		res[k] = v
	}
	res["context"] = vpcLbTestContext()
	return res
}

// vpcLbProfileMock hides the difference between the app profile and monitor profile VPC SDK mocks.
type vpcLbProfileMock struct {
	get     func(id interface{}, sv *data.StructValue, err error) *gomock.Call
	patch   func(id interface{}, err error) *gomock.Call
	del     func(id interface{}, err error) *gomock.Call
	list    func(results []*data.StructValue, err error) *gomock.Call
	restore func()
}

func setupVpcLbAppProfileMock(ctrl *gomock.Controller) *vpcLbProfileMock {
	mockSDK := vpcsmocks.NewMockVpcLbAppProfilesClient(ctrl)
	wrapper := &apivpcs.VpcLbAppProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  vpcLbTestProject,
		VPCID:      vpcLbTestVPC,
	}
	original := cliVpcLbAppProfilesClient
	cliVpcLbAppProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcLbAppProfileClientContext {
		return wrapper
	}
	return &vpcLbProfileMock{
		get: func(id interface{}, sv *data.StructValue, err error) *gomock.Call {
			return mockSDK.EXPECT().Get(vpcLbTestOrg, vpcLbTestProject, vpcLbTestVPC, id).Return(sv, err)
		},
		patch: func(id interface{}, err error) *gomock.Call {
			return mockSDK.EXPECT().Patch(vpcLbTestOrg, vpcLbTestProject, vpcLbTestVPC, id, gomock.Any()).Return(err)
		},
		del: func(id interface{}, err error) *gomock.Call {
			return mockSDK.EXPECT().Delete(vpcLbTestOrg, vpcLbTestProject, vpcLbTestVPC, id, gomock.Any()).Return(err)
		},
		list: func(results []*data.StructValue, err error) *gomock.Call {
			return mockSDK.EXPECT().List(vpcLbTestOrg, vpcLbTestProject, vpcLbTestVPC, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nsxModel.LBAppProfileListResult{Results: results}, err)
		},
		restore: func() { cliVpcLbAppProfilesClient = original },
	}
}

func setupVpcLbMonitorProfileMock(ctrl *gomock.Controller) *vpcLbProfileMock {
	mockSDK := vpcsmocks.NewMockVpcLbMonitorProfilesClient(ctrl)
	wrapper := &apivpcs.VpcLbMonitorProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  vpcLbTestProject,
		VPCID:      vpcLbTestVPC,
	}
	original := cliVpcLbMonitorProfilesClient
	cliVpcLbMonitorProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcLbMonitorProfileClientContext {
		return wrapper
	}
	return &vpcLbProfileMock{
		get: func(id interface{}, sv *data.StructValue, err error) *gomock.Call {
			return mockSDK.EXPECT().Get(vpcLbTestOrg, vpcLbTestProject, vpcLbTestVPC, id).Return(sv, err)
		},
		patch: func(id interface{}, err error) *gomock.Call {
			return mockSDK.EXPECT().Patch(vpcLbTestOrg, vpcLbTestProject, vpcLbTestVPC, id, gomock.Any()).Return(err)
		},
		del: func(id interface{}, err error) *gomock.Call {
			return mockSDK.EXPECT().Delete(vpcLbTestOrg, vpcLbTestProject, vpcLbTestVPC, id, gomock.Any()).Return(err)
		},
		list: func(results []*data.StructValue, err error) *gomock.Call {
			return mockSDK.EXPECT().List(vpcLbTestOrg, vpcLbTestProject, vpcLbTestVPC, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nsxModel.LBMonitorProfileListResult{Results: results}, err)
		},
		restore: func() { cliVpcLbMonitorProfilesClient = original },
	}
}

// vpcLbProfileResourceSpec describes a VPC LB profile resource under test.
type vpcLbProfileResourceSpec struct {
	name         string
	apiSegment   string
	newResource  func() *schema.Resource
	minimalData  func() map[string]interface{}
	structValue  func(t *testing.T) *data.StructValue
	assertFields func(t *testing.T, d *schema.ResourceData)
	id           string
	setupMock    func(ctrl *gomock.Controller) *vpcLbProfileMock
}

func newVpcLbAppProfileResourceSpec(name string, newResource func() *schema.Resource, minimal func() map[string]interface{},
	sv func(t *testing.T) *data.StructValue, assertFields func(t *testing.T, d *schema.ResourceData), id string) vpcLbProfileResourceSpec {
	return vpcLbProfileResourceSpec{
		name: name, apiSegment: "vpc-lb-app-profiles", newResource: newResource, minimalData: minimal,
		structValue: sv, assertFields: assertFields, id: id, setupMock: setupVpcLbAppProfileMock,
	}
}

func newVpcLbMonitorProfileResourceSpec(name string, newResource func() *schema.Resource, minimal func() map[string]interface{},
	sv func(t *testing.T) *data.StructValue, assertFields func(t *testing.T, d *schema.ResourceData), id string) vpcLbProfileResourceSpec {
	return vpcLbProfileResourceSpec{
		name: name, apiSegment: "vpc-lb-monitor-profiles", newResource: newResource, minimalData: minimal,
		structValue: sv, assertFields: assertFields, id: id, setupMock: setupVpcLbMonitorProfileMock,
	}
}

func runVpcLbProfileResourceTests(t *testing.T, spec vpcLbProfileResourceSpec) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	newData := func(t *testing.T, res *schema.Resource) *schema.ResourceData {
		return schema.TestResourceDataRaw(t, res.Schema, vpcLbWithContext(spec.minimalData()))
	}
	// setup creates a controller and a mock, and returns the cleanup func
	setup := func(t *testing.T) (*vpcLbProfileMock, func()) {
		ctrl := gomock.NewController(t)
		m := spec.setupMock(ctrl)
		return m, func() { m.restore(); ctrl.Finish() }
	}
	sessionClient := newGoMockProviderClient()

	t.Run("Schema", func(t *testing.T) {
		res := spec.newResource()
		require.NotNil(t, res.Schema["context"])
		assert.True(t, res.Schema["context"].Required)
		require.NotNil(t, res.Importer)
		require.NotNil(t, res.Create)
		require.NotNil(t, res.Read)
		require.NotNil(t, res.Update)
		require.NotNil(t, res.Delete)
	})

	t.Run("Create success", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		gomock.InOrder(
			m.get(spec.id, nil, vapiErrors.NotFound{}),
			m.patch(spec.id, nil),
			m.get(spec.id, spec.structValue(t), nil),
		)
		res := spec.newResource()
		d := newData(t, res)
		require.NoError(t, res.Create(d, sessionClient))
		assert.Equal(t, spec.id, d.Id())
		spec.assertFields(t, d)
		assert.Equal(t, vpcLbTestProject, d.Get("context.0.project_id"))
		assert.Equal(t, vpcLbTestVPC, d.Get("context.0.vpc_id"))
	})

	t.Run("Create success with generated ID", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		gomock.InOrder(
			m.patch(gomock.Any(), nil),
			m.get(gomock.Any(), spec.structValue(t), nil),
		)
		res := spec.newResource()
		raw := spec.minimalData()
		delete(raw, "nsx_id")
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbWithContext(raw))
		require.NoError(t, res.Create(d, sessionClient))
		assert.NotEmpty(t, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.get(spec.id, spec.structValue(t), nil)
		res := spec.newResource()
		err := res.Create(newData(t, res), sessionClient)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails on existence check error", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.get(spec.id, nil, vapiErrors.InternalServerError{})
		res := spec.newResource()
		require.Error(t, res.Create(newData(t, res), sessionClient))
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.2.0"
		defer func() { util.NsxVersion = "9.0.0" }()
		// no API calls are expected
		_, cleanup := setup(t)
		defer cleanup()
		res := spec.newResource()
		err := res.Create(newData(t, res), sessionClient)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "requires NSX version 9.0.0 or higher")
	})

	t.Run("Create fails on patch error", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		gomock.InOrder(
			m.get(spec.id, nil, vapiErrors.NotFound{}),
			m.patch(spec.id, vapiErrors.InternalServerError{}),
		)
		res := spec.newResource()
		require.Error(t, res.Create(newData(t, res), sessionClient))
	})

	t.Run("Read success", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.get(spec.id, spec.structValue(t), nil)
		res := spec.newResource()
		d := newData(t, res)
		d.SetId(spec.id)
		require.NoError(t, res.Read(d, sessionClient))
		spec.assertFields(t, d)
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.get(spec.id, nil, vapiErrors.NotFound{})
		res := spec.newResource()
		d := newData(t, res)
		d.SetId(spec.id)
		require.NoError(t, res.Read(d, sessionClient))
		assert.Empty(t, d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.get(spec.id, nil, vapiErrors.InternalServerError{})
		res := spec.newResource()
		d := newData(t, res)
		d.SetId(spec.id)
		require.Error(t, res.Read(d, sessionClient))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		_, cleanup := setup(t)
		defer cleanup()
		res := spec.newResource()
		require.Error(t, res.Read(newData(t, res), sessionClient))
	})

	t.Run("Update success", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		gomock.InOrder(
			m.patch(spec.id, nil),
			m.get(spec.id, spec.structValue(t), nil),
		)
		res := spec.newResource()
		d := newData(t, res)
		d.SetId(spec.id)
		require.NoError(t, res.Update(d, sessionClient))
		spec.assertFields(t, d)
	})

	t.Run("Update fails on patch error", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.patch(spec.id, vapiErrors.InternalServerError{})
		res := spec.newResource()
		d := newData(t, res)
		d.SetId(spec.id)
		require.Error(t, res.Update(d, sessionClient))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		_, cleanup := setup(t)
		defer cleanup()
		res := spec.newResource()
		require.Error(t, res.Update(newData(t, res), sessionClient))
	})

	t.Run("Delete success", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.del(spec.id, nil)
		res := spec.newResource()
		d := newData(t, res)
		d.SetId(spec.id)
		require.NoError(t, res.Delete(d, sessionClient))
	})

	t.Run("Delete API error is propagated", func(t *testing.T) {
		m, cleanup := setup(t)
		defer cleanup()
		m.del(spec.id, vapiErrors.InternalServerError{})
		res := spec.newResource()
		d := newData(t, res)
		d.SetId(spec.id)
		require.Error(t, res.Delete(d, sessionClient))
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		_, cleanup := setup(t)
		defer cleanup()
		res := spec.newResource()
		require.Error(t, res.Delete(newData(t, res), sessionClient))
	})

	t.Run("Import by VPC path", func(t *testing.T) {
		res := spec.newResource()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/orgs/default/projects/" + vpcLbTestProject + "/vpcs/" + vpcLbTestVPC + "/" + spec.apiSegment + "/" + spec.id)
		rds, err := res.Importer.State(d, sessionClient)
		require.NoError(t, err)
		require.Len(t, rds, 1)
		assert.Equal(t, spec.id, rds[0].Id())
		assert.Equal(t, vpcLbTestProject, rds[0].Get("context.0.project_id"))
		assert.Equal(t, vpcLbTestVPC, rds[0].Get("context.0.vpc_id"))
	})

	t.Run("Import fails on non VPC path", func(t *testing.T) {
		res := spec.newResource()
		for _, id := range []string{
			"/infra/lb-profiles/" + spec.id,
			"/orgs/default/projects/" + vpcLbTestProject + "/" + spec.apiSegment + "/" + spec.id,
		} {
			d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
			d.SetId(id)
			_, err := res.Importer.State(d, sessionClient)
			require.Error(t, err, id)
		}
	})

	t.Run("Import fails on invalid ID", func(t *testing.T) {
		res := spec.newResource()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(spec.id)
		_, err := res.Importer.State(d, sessionClient)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected format")
		assert.Contains(t, err.Error(), spec.apiSegment)
	})

	t.Run("Import fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.2.0"
		defer func() { util.NsxVersion = "9.0.0" }()
		res := spec.newResource()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/orgs/default/projects/" + vpcLbTestProject + "/vpcs/" + vpcLbTestVPC + "/" + spec.apiSegment + "/" + spec.id)
		_, err := res.Importer.State(d, sessionClient)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "requires NSX version 9.0.0 or higher")
	})
}
