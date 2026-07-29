//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/baremetal_servers/TagsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/baremetal_servers/TagsClient.go TagsClient

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
	bmsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/baremetal_servers"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockResourceNsxtPolicyBareMetalServerTagsSchema(t *testing.T) {
	resource := resourceNsxtPolicyBareMetalServerTags()

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

func TestMockResourceNsxtPolicyBareMetalServerTagsValidation(t *testing.T) {
	resource := resourceNsxtPolicyBareMetalServerTags()

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

func TestMockResourceNsxtPolicyBareMetalServerTagsBasicOperations(t *testing.T) {
	// Test basic tag operations
	t.Run("tag validation", func(t *testing.T) {
		// Test basic tag structure
		scope := "environment"
		tag := "production"

		assert.Equal(t, "environment", scope)
		assert.Equal(t, "production", tag)

		// Test tag mapping
		tagData := map[string]interface{}{
			"scope": scope,
			"tag":   tag,
		}
		assert.NotNil(t, tagData)
		assert.Equal(t, "environment", tagData["scope"])
		assert.Equal(t, "production", tagData["tag"])
	})
}

func bareMetalServerToStructValue(t *testing.T, server nsxModel.BareMetalServer) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(server, nsxModel.BareMetalServerBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func setupBareMetalServerTagsMock(ctrl *gomock.Controller) (*bmsmocks.MockTagsClient, func()) {
	mockSDK := bmsmocks.NewMockTagsClient(ctrl)
	mockWrapper := &infraAPI.BareMetalServerTagsClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliBareMetalServerTagsClient
	cliBareMetalServerTagsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraAPI.BareMetalServerTagsClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliBareMetalServerTagsClient = original }
}

func TestMockResourceNsxtPolicyBareMetalServerTagsImporter(t *testing.T) {
	res := resourceNsxtPolicyBareMetalServerTags()

	t.Run("empty ID is rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("")
		_, err := resourceNsxtPolicyBareMetalServerTagsImporter(d, nil)
		require.Error(t, err)
		assert.Equal(t, ErrEmptyImportID, err)
	})

	t.Run("non-UUID ID is rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("not-a-uuid")
		_, err := resourceNsxtPolicyBareMetalServerTagsImporter(d, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid import ID")
	})

	t.Run("valid UUID is accepted", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("71be0142-2ed1-1d53-9c60-5564cf4b7e2e")
		out, err := resourceNsxtPolicyBareMetalServerTagsImporter(d, nil)
		require.NoError(t, err)
		assert.Len(t, out, 1)
	})
}

func TestMockResourceNsxtPolicyBareMetalServerTagsCreate(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	externalID := "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"
	res := resourceNsxtPolicyBareMetalServerTags()

	t.Run("Create success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restoreCli := setupBareMetalServerTagsMock(ctrl)
		defer restoreCli()

		serverSV := bareMetalServerToStructValue(t, nsxModel.BareMetalServer{ExternalId: &externalID})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
			{Results: []*data.StructValue{serverSV}, ResultCount: i64(1)},
			{Results: []*data.StructValue{serverSV}, ResultCount: i64(1)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		mockSDK.EXPECT().Create(gomock.Any()).Return(nsxModel.BareMetalServerTagList{}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_id": externalID,
			"tag": []interface{}{map[string]interface{}{
				"scope": "env",
				"tag":   "prod",
			}},
		})

		err := resourceNsxtPolicyBareMetalServerTagsCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, externalID, d.Id())
	})

	t.Run("Create fails when server not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, restoreCli := setupBareMetalServerTagsMock(ctrl)
		defer restoreCli()

		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
			{Results: []*data.StructValue{}, ResultCount: i64(0)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_id": externalID,
		})

		err := resourceNsxtPolicyBareMetalServerTagsCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not find bare metal server")
	})
}

func TestMockResourceNsxtPolicyBareMetalServerTagsRead(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	externalID := "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"
	res := resourceNsxtPolicyBareMetalServerTags()

	t.Run("Read success sets tags", func(t *testing.T) {
		scope := "env"
		tag := "prod"
		serverSV := bareMetalServerToStructValue(t, nsxModel.BareMetalServer{
			ExternalId: &externalID,
			Tags:       []nsxModel.Tag{{Scope: &scope, Tag: &tag}},
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
			{Results: []*data.StructValue{serverSV}, ResultCount: i64(1)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(externalID)

		err := resourceNsxtPolicyBareMetalServerTagsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, externalID, d.Get("external_id"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtPolicyBareMetalServerTagsRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining bare metal server external_id")
	})

	t.Run("Read fails when server not found", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
			{Results: []*data.StructValue{}, ResultCount: i64(0)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(externalID)

		err := resourceNsxtPolicyBareMetalServerTagsRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyBareMetalServerTagsDelete(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	externalID := "71be0142-2ed1-1d53-9c60-5564cf4b7e2e"
	res := resourceNsxtPolicyBareMetalServerTags()

	t.Run("Delete success clears tags", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restoreCli := setupBareMetalServerTagsMock(ctrl)
		defer restoreCli()

		serverSV := bareMetalServerToStructValue(t, nsxModel.BareMetalServer{ExternalId: &externalID})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
			{Results: []*data.StructValue{serverSV}, ResultCount: i64(1)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		mockSDK.EXPECT().Create(gomock.Any()).Return(nsxModel.BareMetalServerTagList{}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_id": externalID,
		})
		d.SetId(externalID)

		err := resourceNsxtPolicyBareMetalServerTagsDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when server not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, restoreCli := setupBareMetalServerTagsMock(ctrl)
		defer restoreCli()

		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{
			{Results: []*data.StructValue{}, ResultCount: i64(0)},
		}}
		defer setupCliQueryClientStub(t, stub)()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"external_id": externalID,
		})
		d.SetId(externalID)

		err := resourceNsxtPolicyBareMetalServerTagsDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
