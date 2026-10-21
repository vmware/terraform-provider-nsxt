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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
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
	util.NsxVersion = "9.0.0"
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

	t.Run("Create fails on NSX version lower than 9.0.0", func(t *testing.T) {
		util.NsxVersion = "4.2.0"
		defer func() { util.NsxVersion = "9.0.0" }()

		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbPoolData())
		err := res.Create(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "requires NSX version 9.0.0")
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
	util.NsxVersion = "9.0.0"
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
		assert.Contains(t, err.Error(), "/vpc-lb-pools/[pool]")
	})

	t.Run("project path without VPC is rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/orgs/default/projects/proj1/vpc-lb-pools/pool1")
		_, err := res.Importer.State(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("import fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.2.0"
		defer func() { util.NsxVersion = "9.0.0" }()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vpcLbPoolPath)
		_, err := res.Importer.State(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "requires NSX version 9.0.0")
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

func TestMockResourceNsxtVpcLBPoolEdgeCases(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()
	m := newGoMockProviderClient()

	t.Run("Exists propagates API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(nsxModel.LBPool{}, vapiErrors.InternalServerError{})
		sc := utl.SessionContext{ClientType: utl.VPC, ProjectID: "proj1", VPCID: "vpc1"}
		exists, err := resourceNsxtVpcLBPoolExists(sc, vpcLbPoolID, nil)
		require.Error(t, err)
		assert.False(t, exists)
	})

	t.Run("Create fails when exists check fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(nsxModel.LBPool{}, vapiErrors.InternalServerError{})
		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbPoolData())
		require.Error(t, res.Create(d, m))
	})

	t.Run("Create fails on invalid snat", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(nsxModel.LBPool{}, vapiErrors.NotFound{})
		res := resourceNsxtVpcLBPool()
		data := minimalVpcLbPoolData()
		data["snat"] = []interface{}{map[string]interface{}{
			"type":              "AUTOMAP",
			"ip_pool_addresses": []interface{}{"10.0.0.0/24"},
		}}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		err := res.Create(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ip_pool_addresses")
	})

	t.Run("Update fails on invalid snat", func(t *testing.T) {
		res := resourceNsxtVpcLBPool()
		data := minimalVpcLbPoolData()
		data["snat"] = []interface{}{map[string]interface{}{
			"type":              "AUTOMAP",
			"ip_pool_addresses": []interface{}{"10.0.0.0/24"},
		}}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(vpcLbPoolID)
		require.Error(t, res.Update(d, m))
	})

	t.Run("Empty ID is rejected by Read, Update and Delete", func(t *testing.T) {
		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbPoolData())
		assert.Contains(t, res.Read(d, m).Error(), "VPC LBPool ID")
		assert.Contains(t, res.Update(d, m).Error(), "VPC LBPool ID")
		assert.Contains(t, res.Delete(d, m).Error(), "VPC LBPool ID")
	})

	t.Run("Update sends revision and deprecated monitor path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update("default", "proj1", "vpc1", vpcLbPoolID, gomock.Any()).
				DoAndReturn(func(_, _, _, _ string, obj nsxModel.LBPool) (nsxModel.LBPool, error) {
					assert.Equal(t, vpcLbPoolRevision, *obj.Revision)
					assert.Equal(t, []string{vpcLbPoolMonitorPath}, obj.ActiveMonitorPaths)
					require.NotNil(t, obj.MemberGroup)
					return vpcLbPoolAPIResponse(), nil
				}),
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(vpcLbPoolAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBPool()
		data := minimalVpcLbPoolData()
		data["revision"] = int(vpcLbPoolRevision)
		data["active_monitor_path"] = vpcLbPoolMonitorPath
		data["member_group"] = []interface{}{map[string]interface{}{
			"group_path": "/orgs/default/projects/proj1/infra/domains/default/groups/g1",
			"allow_ipv4": true,
			"allow_ipv6": true,
		}}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(vpcLbPoolID)
		require.NoError(t, res.Update(d, m))
	})

	t.Run("Read populates members, member group, snat and deprecated monitor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		ip := "10.0.0.5"
		port := "80"
		resp := vpcLbPoolAPIResponse()
		resp.MinActiveMembers = nil
		resp.ActiveMonitorPaths = []string{vpcLbPoolMonitorPath}
		resp.Members = []nsxModel.LBPoolMember{{IpAddress: &ip, Port: &port}}
		groupPath := "/orgs/default/projects/proj1/infra/domains/default/groups/g1"
		filter := nsxModel.LBPoolMemberGroup_IP_REVISION_FILTER_IPV4
		resp.MemberGroup = &nsxModel.LBPoolMemberGroup{GroupPath: &groupPath, IpRevisionFilter: &filter}
		snat, err := getPolicyPoolSnatFromSchema(schema.TestResourceDataRaw(t, resourceNsxtVpcLBPool().Schema, map[string]interface{}{
			"snat": []interface{}{map[string]interface{}{"type": "DISABLED"}},
		}))
		require.NoError(t, err)
		resp.SnatTranslation = snat
		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(resp, nil)

		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"context":             vpcLbPoolContext(),
			"active_monitor_path": vpcLbPoolMonitorPath,
		})
		d.SetId(vpcLbPoolID)
		require.NoError(t, res.Read(d, m))
		assert.Equal(t, vpcLbPoolMonitorPath, d.Get("active_monitor_path"))
		assert.Equal(t, 1, d.Get("min_active_members"))
		assert.Len(t, d.Get("member"), 1)
		assert.Len(t, d.Get("member_group"), 1)
		assert.Equal(t, "DISABLED", d.Get("snat.0.type"))
	})

	t.Run("Read fails on unrecognized snat type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbPoolMock(ctrl)
		defer restore()

		resp := vpcLbPoolAPIResponse()
		resp.SnatTranslation = data.NewStructValue("", map[string]data.DataValue{"type": data.NewStringValue("BOGUS")})
		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbPoolID).Return(resp, nil)

		res := resourceNsxtVpcLBPool()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"context": vpcLbPoolContext()})
		d.SetId(vpcLbPoolID)
		require.Error(t, res.Read(d, m))
	})
}
