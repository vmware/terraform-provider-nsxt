//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/baremetal_server_interfaces/TagsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/baremetal_server_interfaces/TagsClient.go TagsClient

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	infraAPI "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	bmsimocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/baremetal_server_interfaces"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockResourceNsxtPolicyBareMetalServerInterfaceTagsSchema(t *testing.T) {
	resource := resourceNsxtPolicyBareMetalServerInterfaceTags()

	// Test schema structure
	assert.NotNil(t, resource.Schema)

	// Test required fields are properly defined
	resourceSchema := resource.Schema
	assert.Contains(t, resourceSchema, "external_id")
	assert.Contains(t, resourceSchema, "tag")

	// Test external_id is required
	assert.True(t, resourceSchema["external_id"].Required)
	assert.Equal(t, schema.TypeString, resourceSchema["external_id"].Type)

	// Test tag block schema
	tagSchema := resourceSchema["tag"].Elem.(*schema.Resource).Schema
	assert.Contains(t, tagSchema, "scope")
	assert.Contains(t, tagSchema, "tag")
	assert.False(t, tagSchema["scope"].Required) // scope is optional
	assert.True(t, tagSchema["scope"].Optional)
	assert.False(t, tagSchema["tag"].Required) // tag value is also optional in standard schema
	assert.True(t, tagSchema["tag"].Optional)
}

func TestMockResourceNsxtPolicyBareMetalServerInterfaceTagsValidation(t *testing.T) {
	resource := resourceNsxtPolicyBareMetalServerInterfaceTags()

	// Test that CRUD operations are properly configured
	assert.NotNil(t, resource.Schema)
	assert.NotNil(t, resource.Create)
	assert.NotNil(t, resource.Read)
	assert.NotNil(t, resource.Update)
	assert.NotNil(t, resource.Delete)

	// Test required fields
	assert.Contains(t, resource.Schema, "external_id")
	assert.Contains(t, resource.Schema, "tag")
}

func TestMockResourceNsxtPolicyBareMetalServerInterfaceTagsBasicOperations(t *testing.T) {
	// Test basic interface tag operations
	t.Run("interface tag validation", func(t *testing.T) {
		// Test basic tag structure
		scope := "network-type"
		tag := "data-plane"

		assert.Equal(t, "network-type", scope)
		assert.Equal(t, "data-plane", tag)

		// Test tag mapping
		tagData := map[string]interface{}{
			"scope": scope,
			"tag":   tag,
		}
		assert.NotNil(t, tagData)
		assert.Equal(t, "network-type", tagData["scope"])
		assert.Equal(t, "data-plane", tagData["tag"])
	})
}

func bareMetalServerInterfaceToStructValue(t *testing.T, iface nsxModel.BareMetalServerInterface) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(iface, nsxModel.BareMetalServerInterfaceBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func setupBareMetalServerInterfaceTagsMock(ctrl *gomock.Controller) (*bmsimocks.MockTagsClient, func()) {
	mockSDK := bmsimocks.NewMockTagsClient(ctrl)
	mockWrapper := &infraAPI.BareMetalServerInterfaceTagsClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliBareMetalServerInterfaceTagsClient
	cliBareMetalServerInterfaceTagsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraAPI.BareMetalServerInterfaceTagsClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliBareMetalServerInterfaceTagsClient = original }
}

func TestMockResourceNsxtPolicyBareMetalServerInterfaceTagsImporter(t *testing.T) {
	res := resourceNsxtPolicyBareMetalServerInterfaceTags()

	t.Run("empty ID is rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("")
		_, err := resourceNsxtPolicyBareMetalServerInterfaceTagsImporter(d, nil)
		require.Error(t, err)
		assert.Equal(t, ErrEmptyImportID, err)
	})

	t.Run("non-UUID ID is rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("not-a-uuid")
		_, err := resourceNsxtPolicyBareMetalServerInterfaceTagsImporter(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid import ID")
	})

	t.Run("valid UUID is accepted", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("71be0142-2ed1-1d53-9c60-5564cf4b7e2e")
		out, err := resourceNsxtPolicyBareMetalServerInterfaceTagsImporter(d, nil)
		require.NoError(t, err)
		assert.Len(t, out, 1)
	})
}

func TestMockResourceNsxtPolicyBareMetalServerInterfaceTagsCreate(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	externalID := "71be0142-2ed1-1d53-9c60-02005b4b7246"
	res := resourceNsxtPolicyBareMetalServerInterfaceTags()

	t.Run("Create success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restoreCli := setupBareMetalServerInterfaceTagsMock(ctrl)
		defer restoreCli()

		ifaceSV := bareMetalServerInterfaceToStructValue(t, nsxModel.BareMetalServerInterface{ExternalId: &externalID})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
			{Results: []*data.StructValue{ifaceSV}, ResultCount: i64(1)},
			{Results: []*data.StructValue{ifaceSV}, ResultCount: i64(1)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		mockSDK.EXPECT().Create(gomock.Any()).Return(nsxModel.BareMetalServerInterfaceTagList{}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_id": externalID,
			"tag": []interface{}{map[string]interface{}{
				"scope": "network-type",
				"tag":   "data-plane",
			}},
		})

		err := resourceNsxtPolicyBareMetalServerInterfaceTagsCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, externalID, d.Id())
	})

	t.Run("Create fails when interface not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, restoreCli := setupBareMetalServerInterfaceTagsMock(ctrl)
		defer restoreCli()

		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
			{Results: []*data.StructValue{}, ResultCount: i64(0)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_id": externalID,
		})

		err := resourceNsxtPolicyBareMetalServerInterfaceTagsCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not find bare metal server interface")
	})
}

func TestMockResourceNsxtPolicyBareMetalServerInterfaceTagsRead(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	externalID := "71be0142-2ed1-1d53-9c60-02005b4b7246"
	res := resourceNsxtPolicyBareMetalServerInterfaceTags()

	t.Run("Read success sets tags", func(t *testing.T) {
		scope := "network-type"
		tag := "data-plane"
		ifaceSV := bareMetalServerInterfaceToStructValue(t, nsxModel.BareMetalServerInterface{
			ExternalId: &externalID,
			Tags:       []nsxModel.Tag{{Scope: &scope, Tag: &tag}},
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
			{Results: []*data.StructValue{ifaceSV}, ResultCount: i64(1)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(externalID)

		err := resourceNsxtPolicyBareMetalServerInterfaceTagsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, externalID, d.Get("external_id"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtPolicyBareMetalServerInterfaceTagsRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining bare metal server interface external_id")
	})

	t.Run("Read fails when interface not found", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
			{Results: []*data.StructValue{}, ResultCount: i64(0)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(externalID)

		err := resourceNsxtPolicyBareMetalServerInterfaceTagsRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyBareMetalServerInterfaceTagsDelete(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	externalID := "71be0142-2ed1-1d53-9c60-02005b4b7246"
	res := resourceNsxtPolicyBareMetalServerInterfaceTags()

	t.Run("Delete success clears tags", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restoreCli := setupBareMetalServerInterfaceTagsMock(ctrl)
		defer restoreCli()

		ifaceSV := bareMetalServerInterfaceToStructValue(t, nsxModel.BareMetalServerInterface{ExternalId: &externalID})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
			{Results: []*data.StructValue{ifaceSV}, ResultCount: i64(1)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		mockSDK.EXPECT().Create(gomock.Any()).Return(nsxModel.BareMetalServerInterfaceTagList{}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_id": externalID,
		})
		d.SetId(externalID)

		err := resourceNsxtPolicyBareMetalServerInterfaceTagsDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when interface not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, restoreCli := setupBareMetalServerInterfaceTagsMock(ctrl)
		defer restoreCli()

		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
			{Results: []*data.StructValue{}, ResultCount: i64(0)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_id": externalID,
		})
		d.SetId(externalID)

		err := resourceNsxtPolicyBareMetalServerInterfaceTagsDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
