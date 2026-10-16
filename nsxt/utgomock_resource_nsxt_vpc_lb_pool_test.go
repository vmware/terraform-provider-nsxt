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
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	apivpcs "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	vpcsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vpcLbPoolID          = "vpc-lb-pool-1"
	vpcLbPoolDisplayName = "vpc-lb-pool-fooname"
	vpcLbPoolDescription = "vpc lb pool mock"
	vpcLbPoolPath        = "/orgs/default/projects/proj1/vpcs/vpc1/vpc-lb-pools/vpc-lb-pool-1"
	vpcLbPoolRevision    = int64(1)
	vpcLbPoolAlgorithm   = nsxModel.LBPool_ALGORITHM_ROUND_ROBIN
	vpcLbPoolMinActive   = int64(1)
	vpcLbPoolMonitorPath = "/orgs/default/projects/proj1/vpcs/vpc1/vpc-lb-monitor-profiles/mon1"
)

func vpcLbPoolAPIResponse() nsxModel.LBPool {
	return nsxModel.LBPool{
		Id:               &vpcLbPoolID,
		DisplayName:      &vpcLbPoolDisplayName,
		Description:      &vpcLbPoolDescription,
		Path:             &vpcLbPoolPath,
		Revision:         &vpcLbPoolRevision,
		Algorithm:        &vpcLbPoolAlgorithm,
		MinActiveMembers: &vpcLbPoolMinActive,
	}
}

func vpcLbPoolContext() []interface{} {
	return []interface{}{map[string]interface{}{"project_id": "proj1", "vpc_id": "vpc1"}}
}

func minimalVpcLbPoolData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": vpcLbPoolDisplayName,
		"description":  vpcLbPoolDescription,
		"nsx_id":       vpcLbPoolID,
		"context":      vpcLbPoolContext(),
	}
}

func setupVpcLbPoolMock(ctrl *gomock.Controller) (*vpcsmocks.MockVpcLbPoolsClient, func()) {
	mockSDK := vpcsmocks.NewMockVpcLbPoolsClient(ctrl)
	mockWrapper := &apivpcs.VpcLbPoolClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  "proj1",
		VPCID:      "vpc1",
	}
	original := cliVpcLbPoolsClient
	cliVpcLbPoolsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcLbPoolClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcLbPoolsClient = original }
}

func TestMockResourceNsxtVpcLBPoolCreate(t *testing.T) {
	util.NsxVersion = "9.0.2"
	defer func() { util.NsxVersion = "" }()

	t.Run("Create success with nested blocks", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(nsxModel.LBPool{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch("default", "proj1", "vpc1", vpcLbPoolID, gomock.Any()).
				DoAndReturn(func(_, _, _, _ string, obj nsxModel.LBPool) error {
					assert.Equal(t, []string{vpcLbPoolMonitorPath}, obj.ActiveMonitorPaths)
					require.Len(t, obj.Members, 1)
					assert.Equal(t, "10.0.0.5", *obj.Members[0].IpAddress)
					assert.NotNil(t, obj.SnatTranslation)
					return nil
				}),
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(vpcLbPoolAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBPool()
		data := minimalVpcLbPoolData()
		data["active_monitor_paths"] = []interface{}{vpcLbPoolMonitorPath}
		data["member"] = []interface{}{map[string]interface{}{
			"display_name": "member-1",
			"admin_state":  "ENABLED",
			"ip_address":   "10.0.0.5",
			"port":         "80",
			"weight":       1,
		}}
		data["snat"] = []interface{}{map[string]interface{}{"type": "AUTOMAP"}}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := res.Create(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcLbPoolID, d.Id())
		assert.Equal(t, vpcLbPoolDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcLbPoolPath, d.Get("path"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(vpcLbPoolAPIResponse(), nil)

		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbPoolData())
		err := res.Create(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails on NSX version lower than 9.0.2", func(t *testing.T) {
		util.NsxVersion = "9.0.1"
		defer func() { util.NsxVersion = "9.0.2" }()

		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbPoolData())
		err := res.Create(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "requires NSX version 9.0.2")
	})

	t.Run("Create fails when patch fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(nsxModel.LBPool{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch("default", "proj1", "vpc1", vpcLbPoolID, gomock.Any()).Return(errors.New("patch failed")),
		)

		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbPoolData())
		err := res.Create(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "patch failed")
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcLBPoolRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		resp := vpcLbPoolAPIResponse()
		resp.ActiveMonitorPaths = []string{vpcLbPoolMonitorPath}
		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(resp, nil)

		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"context": vpcLbPoolContext()})
		d.SetId(vpcLbPoolID)

		require.NoError(t, res.Read(d, newGoMockProviderClient()))
		assert.Equal(t, vpcLbPoolDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcLbPoolPath, d.Get("path"))
		assert.Equal(t, int(vpcLbPoolRevision), d.Get("revision"))
		assert.Equal(t, []interface{}{vpcLbPoolMonitorPath}, d.Get("active_monitor_paths"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(nsxModel.LBPool{}, vapiErrors.NotFound{})

		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"context": vpcLbPoolContext()})
		d.SetId(vpcLbPoolID)

		require.NoError(t, res.Read(d, newGoMockProviderClient()))
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails on API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(nsxModel.LBPool{}, vapiErrors.InternalServerError{})

		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"context": vpcLbPoolContext()})
		d.SetId(vpcLbPoolID)

		require.Error(t, res.Read(d, newGoMockProviderClient()))
	})

	t.Run("Read bypasses the cache for VPC contexts when cache mode is on", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		mc := newGoMockProviderClientCacheEnabled()
		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"context": vpcLbPoolContext()})
		d.SetId(vpcLbPoolID)

		assert.True(t, isCacheEnabledForRead(d, mc))
		assert.False(t, isLbCacheEnabledForRead(d, mc))

		// A single Get and no tag Patch (which the cache-aware path would issue in config_scope mode)
		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(vpcLbPoolAPIResponse(), nil).Times(1)

		require.NoError(t, res.Read(d, mc))
		assert.Equal(t, vpcLbPoolDisplayName, d.Get("display_name"))
	})

	t.Run("Cache stays enabled for non-VPC contexts", func(t *testing.T) {
		mc := newGoMockProviderClientCacheEnabled()
		res := resourceNsxtPolicyLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vpcLbPoolID)
		assert.True(t, isLbCacheEnabledForRead(d, mc))
	})
}

func TestMockResourceNsxtVpcLBPoolUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update("default", "proj1", "vpc1", vpcLbPoolID, gomock.Any()).Return(vpcLbPoolAPIResponse(), nil),
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(vpcLbPoolAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbPoolData())
		d.SetId(vpcLbPoolID)

		require.NoError(t, res.Update(d, newGoMockProviderClient()))
	})

	t.Run("Update fails on API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Update("default", "proj1", "vpc1", vpcLbPoolID, gomock.Any()).Return(nsxModel.LBPool{}, errors.New("update failed"))

		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbPoolData())
		d.SetId(vpcLbPoolID)

		err := res.Update(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "update failed")
	})
}

func TestMockResourceNsxtVpcLBPoolDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Delete("default", "proj1", "vpc1", vpcLbPoolID, gomock.Any()).Return(nil)

		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbPoolData())
		d.SetId(vpcLbPoolID)

		require.NoError(t, res.Delete(d, newGoMockProviderClient()))
	})

	t.Run("Delete fails on API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Delete("default", "proj1", "vpc1", vpcLbPoolID, gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbPoolData())
		d.SetId(vpcLbPoolID)

		err := res.Delete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestUnitNsxt_vpcLBPoolImporter(t *testing.T) {
	util.NsxVersion = "9.0.2"
	defer func() { util.NsxVersion = "" }()
	res := resourceNsxtVpcLBPool()

	t.Run("valid VPC path sets context and ID", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vpcLbPoolPath)
		out, err := res.Importer.State(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, vpcLbPoolID, out[0].Id())
		projectID, vpcID, _ := getContextDataFromSchema(out[0])
		assert.Equal(t, "proj1", projectID)
		assert.Equal(t, "vpc1", vpcID)
	})

	t.Run("non-path ID is rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("some-id")
		_, err := res.Importer.State(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "/vpc-lb-pools/[id]")
	})

	t.Run("project path without VPC is rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/orgs/default/projects/proj1/vpc-lb-pools/pool1")
		_, err := res.Importer.State(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("import fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "9.0.1"
		defer func() { util.NsxVersion = "9.0.2" }()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vpcLbPoolPath)
		_, err := res.Importer.State(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "requires NSX version 9.0.2")
	})
}

func TestUnitNsxt_vpcLBPoolSchemaHasRequiredContext(t *testing.T) {
	res := resourceNsxtVpcLBPool()
	ctx, ok := res.Schema["context"]
	require.True(t, ok)
	assert.True(t, ctx.Required)
	assert.Contains(t, res.Schema, "member")
	assert.Contains(t, res.Schema, "snat")
}

func TestUnitNsxt_vpcLBPoolAdapter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupVpcLbPoolMock(ctrl)
	defer restore()

	sessionContext := utl.SessionContext{ClientType: utl.VPC, ProjectID: "proj1", VPCID: "vpc1"}
	api := getLbPoolAPI(sessionContext, nil)
	require.NotNil(t, api)

	t.Run("Get", func(t *testing.T) {
		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(vpcLbPoolAPIResponse(), nil)
		obj, err := api.Get(vpcLbPoolID)
		require.NoError(t, err)
		assert.Equal(t, vpcLbPoolID, *obj.Id)
	})

	t.Run("Patch", func(t *testing.T) {
		mockSDK.EXPECT().Patch("default", "proj1", "vpc1", vpcLbPoolID, gomock.Any()).Return(nil)
		require.NoError(t, api.Patch(vpcLbPoolID, vpcLbPoolAPIResponse()))
	})

	t.Run("Update", func(t *testing.T) {
		mockSDK.EXPECT().Update("default", "proj1", "vpc1", vpcLbPoolID, gomock.Any()).Return(vpcLbPoolAPIResponse(), nil)
		obj, err := api.Update(vpcLbPoolID, vpcLbPoolAPIResponse())
		require.NoError(t, err)
		assert.Equal(t, vpcLbPoolID, *obj.Id)
	})

	t.Run("Delete", func(t *testing.T) {
		force := true
		mockSDK.EXPECT().Delete("default", "proj1", "vpc1", vpcLbPoolID, &force).Return(nil)
		require.NoError(t, api.Delete(vpcLbPoolID, &force))
	})

	t.Run("List", func(t *testing.T) {
		results := []nsxModel.LBPool{vpcLbPoolAPIResponse()}
		mockSDK.EXPECT().List("default", "proj1", "vpc1", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBPoolListResult{Results: results}, nil)
		res, err := api.List(nil, nil, nil, nil, nil, nil)
		require.NoError(t, err)
		require.Len(t, res.Results, 1)
		assert.Equal(t, vpcLbPoolID, *res.Results[0].Id)
	})

	t.Run("nil VPC client yields nil API", func(t *testing.T) {
		original := cliVpcLbPoolsClient
		defer func() { cliVpcLbPoolsClient = original }()
		cliVpcLbPoolsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcLbPoolClientContext {
			return nil
		}
		assert.Nil(t, getLbPoolAPI(sessionContext, nil))

		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbPoolData())
		d.SetId(vpcLbPoolID)
		require.Error(t, res.Read(d, newGoMockProviderClient()))
	})
}
