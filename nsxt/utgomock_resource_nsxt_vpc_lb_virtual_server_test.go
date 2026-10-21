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

const (
	vpcLbVsProfilePath     = "/orgs/default/projects/proj1/vpcs/vpc1/vpc-lb-app-profiles/app1"
	vpcLbVsPersistencePath = "/orgs/default/projects/proj1/vpcs/vpc1/vpc-lb-persistence-profiles/pers1"
	vpcLbVsPoolPath        = "/orgs/default/projects/proj1/vpcs/vpc1/vpc-lb-pools/pool1"
	vpcLbVsSorryPoolPath   = "/orgs/default/projects/proj1/vpcs/vpc1/vpc-lb-pools/sorry1"
)

var (
	vpcLbVsID          = "vpc-lb-vs-1"
	vpcLbVsDisplayName = "Test VPC LB Virtual Server"
	vpcLbVsDescription = "Test vpc lb virtual server"
	vpcLbVsRevision    = int64(1)
	vpcLbVsIPAddr      = "10.0.0.1"
	vpcLbVsPath        = "/orgs/default/projects/proj1/vpcs/vpc1/vpc-lb-virtual-servers/vpc-lb-vs-1"
)

func vpcLbVsAPIResponse() nsxModel.LBVirtualServer {
	return nsxModel.LBVirtualServer{
		Id:          &vpcLbVsID,
		DisplayName: &vpcLbVsDisplayName,
		Description: &vpcLbVsDescription,
		Revision:    &vpcLbVsRevision,
		Path:        &vpcLbVsPath,
		IpAddress:   &vpcLbVsIPAddr,
		Ports:       []string{"80"},
	}
}

func minimalVpcLbVsData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": vpcLbVsDisplayName,
		"description":  vpcLbVsDescription,
		"nsx_id":       vpcLbVsID,
		"ip_address":   vpcLbVsIPAddr,
		"ports":        []interface{}{"80"},
		"context":      []interface{}{map[string]interface{}{"project_id": "proj1", "vpc_id": "vpc1"}},
	}
}

func vpcLbVsRules() []interface{} {
	return []interface{}{map[string]interface{}{
		"display_name": "rule-1",
		"phase":        "HTTP_FORWARDING",
		"action": []interface{}{map[string]interface{}{
			"select_pool": []interface{}{map[string]interface{}{"pool_id": vpcLbVsPoolPath}},
		}},
	}}
}

func setupVpcLbVsMock(ctrl *gomock.Controller) (*vpcsmocks.MockVpcLbVirtualServersClient, func()) {
	mockSDK := vpcsmocks.NewMockVpcLbVirtualServersClient(ctrl)
	mockWrapper := &apivpcs.VpcLbVirtualServerClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  "proj1",
		VPCID:      "vpc1",
	}
	original := cliVpcLbVirtualServersClient
	cliVpcLbVirtualServersClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcLbVirtualServerClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcLbVirtualServersClient = original }
}

func TestMockResourceNsxtVpcLBVirtualServerCreate(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	t.Run("Create success with paths and rules", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(nsxModel.LBVirtualServer{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch("default", "proj1", "vpc1", vpcLbVsID, gomock.Any()).
				DoAndReturn(func(_, _, _, _ string, obj nsxModel.LBVirtualServer) error {
					assert.Equal(t, vpcLbVsProfilePath, *obj.ApplicationProfilePath)
					assert.Equal(t, vpcLbVsPersistencePath, *obj.LbPersistenceProfilePath)
					assert.Equal(t, vpcLbVsPoolPath, *obj.PoolPath)
					assert.Equal(t, vpcLbVsSorryPoolPath, *obj.SorryPoolPath)
					require.Len(t, obj.Rules, 1)
					assert.Equal(t, "rule-1", *obj.Rules[0].DisplayName)
					require.Len(t, obj.Rules[0].Actions, 1)
					return nil
				}),
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(vpcLbVsAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBVirtualServer()
		data := minimalVpcLbVsData()
		data["application_profile_path"] = vpcLbVsProfilePath
		data["persistence_profile_path"] = vpcLbVsPersistencePath
		data["pool_path"] = vpcLbVsPoolPath
		data["sorry_pool_path"] = vpcLbVsSorryPoolPath
		data["rule"] = vpcLbVsRules()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := res.Create(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcLbVsID, d.Id())
		assert.Equal(t, vpcLbVsDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcLbVsPath, d.Get("path"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(vpcLbVsAPIResponse(), nil)

		res := resourceNsxtVpcLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbVsData())
		err := res.Create(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails on NSX version lower than 9.0.0", func(t *testing.T) {
		util.NsxVersion = "4.2.0"
		defer func() { util.NsxVersion = "9.0.0" }()

		res := resourceNsxtVpcLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbVsData())
		err := res.Create(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "requires NSX version 9.0.0")
	})

	t.Run("Create fails when patch fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(nsxModel.LBVirtualServer{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch("default", "proj1", "vpc1", vpcLbVsID, gomock.Any()).Return(errors.New("patch failed")),
		)

		res := resourceNsxtVpcLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbVsData())
		err := res.Create(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "patch failed")
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcLBVirtualServerRead(t *testing.T) {
	ctxData := map[string]interface{}{"context": minimalVpcLbVsData()["context"]}

	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		resp := vpcLbVsAPIResponse()
		resp.ApplicationProfilePath = &[]string{vpcLbVsProfilePath}[0]
		resp.PoolPath = &[]string{vpcLbVsPoolPath}[0]
		resp.LbPersistenceProfilePath = &[]string{vpcLbVsPersistencePath}[0]
		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(resp, nil)

		res := resourceNsxtVpcLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, ctxData)
		d.SetId(vpcLbVsID)

		require.NoError(t, res.Read(d, newGoMockProviderClient()))
		assert.Equal(t, vpcLbVsDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcLbVsIPAddr, d.Get("ip_address"))
		assert.Equal(t, vpcLbVsProfilePath, d.Get("application_profile_path"))
		assert.Equal(t, vpcLbVsPoolPath, d.Get("pool_path"))
		assert.Equal(t, vpcLbVsPersistencePath, d.Get("persistence_profile_path"))
	})

	t.Run("Read imports existing rules into empty state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		// Build the rule model through the create path, then feed it back via Read
		var rules []nsxModel.LBRule
		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(nsxModel.LBVirtualServer{}, vapiErrors.NotFound{})
		mockSDK.EXPECT().Patch("default", "proj1", "vpc1", vpcLbVsID, gomock.Any()).
			DoAndReturn(func(_, _, _, _ string, obj nsxModel.LBVirtualServer) error {
				rules = obj.Rules
				return nil
			})
		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).DoAndReturn(
			func(_, _, _, _ string) (nsxModel.LBVirtualServer, error) {
				resp := vpcLbVsAPIResponse()
				resp.Rules = rules
				return resp, nil
			}).Times(2)

		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()
		res := resourceNsxtVpcLBVirtualServer()
		data := minimalVpcLbVsData()
		data["rule"] = vpcLbVsRules()
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		require.NoError(t, res.Create(d, newGoMockProviderClient()))
		require.Len(t, d.Get("rule").([]interface{}), 1)

		d2 := schema.TestResourceDataRaw(t, res.Schema, ctxData)
		d2.SetId(vpcLbVsID)
		require.NoError(t, res.Read(d2, newGoMockProviderClient()))
		assert.Len(t, d2.Get("rule").([]interface{}), 1)
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(nsxModel.LBVirtualServer{}, vapiErrors.NotFound{})

		res := resourceNsxtVpcLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, ctxData)
		d.SetId(vpcLbVsID)

		require.NoError(t, res.Read(d, newGoMockProviderClient()))
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails on API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(nsxModel.LBVirtualServer{}, vapiErrors.InternalServerError{})

		res := resourceNsxtVpcLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, ctxData)
		d.SetId(vpcLbVsID)

		require.Error(t, res.Read(d, newGoMockProviderClient()))
	})
}

func TestMockResourceNsxtVpcLBVirtualServerUpdate(t *testing.T) {
	t.Run("Update success without rule change keeps existing rules", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		existingRule := nsxModel.LBRule{DisplayName: &[]string{"live-rule"}[0]}
		existing := vpcLbVsAPIResponse()
		existing.Rules = []nsxModel.LBRule{existingRule}
		gomock.InOrder(
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(existing, nil),
			mockSDK.EXPECT().Patch("default", "proj1", "vpc1", vpcLbVsID, gomock.Any()).
				DoAndReturn(func(_, _, _, _ string, obj nsxModel.LBVirtualServer) error {
					require.Len(t, obj.Rules, 1)
					assert.Equal(t, "live-rule", *obj.Rules[0].DisplayName)
					return nil
				}),
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(vpcLbVsAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbVsData())
		d.SetId(vpcLbVsID)

		require.NoError(t, res.Update(d, newGoMockProviderClient()))
	})

	t.Run("Update success with rule change", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Patch("default", "proj1", "vpc1", vpcLbVsID, gomock.Any()).
				DoAndReturn(func(_, _, _, _ string, obj nsxModel.LBVirtualServer) error {
					require.Len(t, obj.Rules, 1)
					return nil
				}),
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(vpcLbVsAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBVirtualServer()
		data := minimalVpcLbVsData()
		data["rule"] = vpcLbVsRules()
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(vpcLbVsID)

		require.NoError(t, res.Update(d, newGoMockProviderClient()))
	})

	t.Run("Update fails when reading existing rules fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(nsxModel.LBVirtualServer{}, errors.New("get failed"))

		res := resourceNsxtVpcLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbVsData())
		d.SetId(vpcLbVsID)

		require.Error(t, res.Update(d, newGoMockProviderClient()))
	})

	t.Run("Update fails when patch fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(vpcLbVsAPIResponse(), nil),
			mockSDK.EXPECT().Patch("default", "proj1", "vpc1", vpcLbVsID, gomock.Any()).Return(errors.New("patch failed")),
		)

		res := resourceNsxtVpcLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbVsData())
		d.SetId(vpcLbVsID)

		err := res.Update(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "patch failed")
	})
}

func TestMockResourceNsxtVpcLBVirtualServerDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Delete("default", "proj1", "vpc1", vpcLbVsID, gomock.Any()).Return(nil)

		res := resourceNsxtVpcLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbVsData())
		d.SetId(vpcLbVsID)

		require.NoError(t, res.Delete(d, newGoMockProviderClient()))
	})

	t.Run("Delete fails on API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Delete("default", "proj1", "vpc1", vpcLbVsID, gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtVpcLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbVsData())
		d.SetId(vpcLbVsID)

		err := res.Delete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func TestUnitNsxt_vpcLBVirtualServerImporter(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()
	res := resourceNsxtVpcLBVirtualServer()

	t.Run("valid VPC path sets context and ID", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vpcLbVsPath)
		out, err := res.Importer.State(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, vpcLbVsID, out[0].Id())
		projectID, vpcID, _ := getContextDataFromSchema(out[0])
		assert.Equal(t, "proj1", projectID)
		assert.Equal(t, "vpc1", vpcID)
	})

	t.Run("non-path ID is rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("some-id")
		_, err := res.Importer.State(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "/vpc-lb-virtual-servers/[virtual server]")
	})

	t.Run("project path without VPC is rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/orgs/default/projects/proj1/vpc-lb-virtual-servers/vs1")
		_, err := res.Importer.State(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("import fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.2.0"
		defer func() { util.NsxVersion = "9.0.0" }()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vpcLbVsPath)
		_, err := res.Importer.State(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "requires NSX version 9.0.0")
	})
}

func TestUnitNsxt_vpcLBVirtualServerSchemaHasRequiredContext(t *testing.T) {
	res := resourceNsxtVpcLBVirtualServer()
	ctx, ok := res.Schema["context"]
	require.True(t, ok)
	assert.True(t, ctx.Required)
	assert.Contains(t, res.Schema, "rule")
	assert.Contains(t, res.Schema, "pool_path")
}

func TestMockResourceNsxtVpcLBVirtualServerEdgeCases(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()
	m := newGoMockProviderClient()

	t.Run("Exists propagates API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(nsxModel.LBVirtualServer{}, vapiErrors.InternalServerError{})
		sc := utl.SessionContext{ClientType: utl.VPC, ProjectID: "proj1", VPCID: "vpc1"}
		exists, err := resourceNsxtVpcLBVirtualServerExists(sc, vpcLbVsID, nil)
		require.Error(t, err)
		assert.False(t, exists)
	})

	t.Run("Create fails when exists check fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(nsxModel.LBVirtualServer{}, vapiErrors.InternalServerError{})
		res := resourceNsxtVpcLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbVsData())
		require.Error(t, res.Create(d, m))
	})

	t.Run("Empty ID is rejected by Read, Update and Delete", func(t *testing.T) {
		res := resourceNsxtVpcLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcLbVsData())
		assert.Contains(t, res.Read(d, m).Error(), "VPC LBVirtualServer ID")
		assert.Contains(t, res.Update(d, m).Error(), "VPC LBVirtualServer ID")
		assert.Contains(t, res.Delete(d, m).Error(), "VPC LBVirtualServer ID")
	})

	t.Run("Create sends connection limits, ssl bindings and access list control", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		certPath := "/orgs/default/projects/proj1/infra/certificates/cert1"
		gomock.InOrder(
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(nsxModel.LBVirtualServer{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch("default", "proj1", "vpc1", vpcLbVsID, gomock.Any()).
				DoAndReturn(func(_, _, _, _ string, obj nsxModel.LBVirtualServer) error {
					assert.Equal(t, int64(100), *obj.MaxConcurrentConnections)
					assert.Equal(t, int64(10), *obj.MaxNewConnectionRate)
					require.NotNil(t, obj.ClientSslProfileBinding)
					assert.Equal(t, certPath, *obj.ClientSslProfileBinding.DefaultCertificatePath)
					require.NotNil(t, obj.ServerSslProfileBinding)
					require.NotNil(t, obj.AccessListControl)
					assert.Equal(t, "DROP", *obj.AccessListControl.Action)
					return nil
				}),
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).DoAndReturn(
				func(_, _, _, _ string) (nsxModel.LBVirtualServer, error) {
					resp := vpcLbVsAPIResponse()
					depth := int64(3)
					auth := "IGNORE"
					resp.ClientSslProfileBinding = &nsxModel.LBClientSslProfileBinding{
						DefaultCertificatePath: &certPath, CertificateChainDepth: &depth, ClientAuth: &auth,
					}
					sAuth := "AUTO_APPLY"
					resp.ServerSslProfileBinding = &nsxModel.LBServerSslProfileBinding{
						CertificateChainDepth: &depth, ServerAuth: &sAuth,
					}
					action, grp, en := "DROP", "/orgs/default/projects/proj1/infra/domains/default/groups/g1", true
					resp.AccessListControl = &nsxModel.LBAccessListControl{Action: &action, GroupPath: &grp, Enabled: &en}
					return resp, nil
				}),
		)

		res := resourceNsxtVpcLBVirtualServer()
		data := minimalVpcLbVsData()
		data["max_concurrent_connections"] = 100
		data["max_new_connection_rate"] = 10
		data["client_ssl"] = []interface{}{map[string]interface{}{"default_certificate_path": certPath}}
		data["server_ssl"] = []interface{}{map[string]interface{}{"server_auth": "AUTO_APPLY"}}
		data["access_list_control"] = []interface{}{map[string]interface{}{
			"action":     "DROP",
			"group_path": "/orgs/default/projects/proj1/infra/domains/default/groups/g1",
		}}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		require.NoError(t, res.Create(d, m))
		assert.Equal(t, certPath, d.Get("client_ssl.0.default_certificate_path"))
		assert.Equal(t, "DROP", d.Get("access_list_control.0.action"))
	})

	t.Run("Update sends connection limits", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbVsMock(ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(vpcLbVsAPIResponse(), nil),
			mockSDK.EXPECT().Patch("default", "proj1", "vpc1", vpcLbVsID, gomock.Any()).
				DoAndReturn(func(_, _, _, _ string, obj nsxModel.LBVirtualServer) error {
					assert.Equal(t, int64(100), *obj.MaxConcurrentConnections)
					assert.Equal(t, int64(10), *obj.MaxNewConnectionRate)
					return nil
				}),
			mockSDK.EXPECT().Get("default", "proj1", "vpc1", vpcLbVsID).Return(vpcLbVsAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBVirtualServer()
		data := minimalVpcLbVsData()
		data["max_concurrent_connections"] = 100
		data["max_new_connection_rate"] = 10
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(vpcLbVsID)
		require.NoError(t, res.Update(d, m))
	})
}
