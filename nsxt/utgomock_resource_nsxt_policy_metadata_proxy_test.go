//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/MetadataProxiesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/MetadataProxiesClient.go MetadataProxiesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	metaproxy "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	metadataProxyID          = "metadata-proxy-1"
	metadataProxyDisplayName = "metadata-proxy-fooname"
	metadataProxyDescription = "metadata proxy mock"
	metadataProxyPath        = "/infra/metadata-proxies/metadata-proxy-1"
	metadataProxyRevision    = int64(1)
	metadataProxyEdgePath    = "/infra/sites/default/enforcement-points/default/edge-clusters/ec1"
	metadataProxySecret      = "secret"
	metadataProxyServerAddr  = "192.168.1.1"
)

func TestMockResourceNsxtPolicyMetadataProxyCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetaSDK := metaproxy.NewMockMetadataProxiesClient(ctrl)
	mockWrapper := &cliinfra.MetadataProxyConfigClientContext{
		Client:     mockMetaSDK,
		ClientType: utl.Local,
	}

	originalCli := cliMetadataProxiesClient
	defer func() { cliMetadataProxiesClient = originalCli }()
	cliMetadataProxiesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.MetadataProxyConfigClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockMetaSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockMetaSDK.EXPECT().Get(gomock.Any()).Return(model.MetadataProxyConfig{
			Id:              &metadataProxyID,
			DisplayName:     &metadataProxyDisplayName,
			Description:     &metadataProxyDescription,
			Path:            &metadataProxyPath,
			Revision:        &metadataProxyRevision,
			EdgeClusterPath: &metadataProxyEdgePath,
			ServerAddress:   &metadataProxyServerAddr,
		}, nil)

		res := resourceNsxtPolicyMetadataProxy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalMetadataProxyData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMetadataProxyCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockMetaSDK.EXPECT().Get("existing-id").Return(model.MetadataProxyConfig{Id: &metadataProxyID}, nil)

		res := resourceNsxtPolicyMetadataProxy()
		data := minimalMetadataProxyData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMetadataProxyCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyMetadataProxyRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetaSDK := metaproxy.NewMockMetadataProxiesClient(ctrl)
	mockWrapper := &cliinfra.MetadataProxyConfigClientContext{
		Client:     mockMetaSDK,
		ClientType: utl.Local,
	}

	originalCli := cliMetadataProxiesClient
	defer func() { cliMetadataProxiesClient = originalCli }()
	cliMetadataProxiesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.MetadataProxyConfigClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockMetaSDK.EXPECT().Get(metadataProxyID).Return(model.MetadataProxyConfig{
			Id:              &metadataProxyID,
			DisplayName:     &metadataProxyDisplayName,
			Description:     &metadataProxyDescription,
			Path:            &metadataProxyPath,
			Revision:        &metadataProxyRevision,
			EdgeClusterPath: &metadataProxyEdgePath,
			ServerAddress:   &metadataProxyServerAddr,
		}, nil)

		res := resourceNsxtPolicyMetadataProxy()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(metadataProxyID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMetadataProxyRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, metadataProxyDisplayName, d.Get("display_name"))
		assert.Equal(t, metadataProxyDescription, d.Get("description"))
		assert.Equal(t, metadataProxyPath, d.Get("path"))
		assert.Equal(t, int(metadataProxyRevision), d.Get("revision"))
		assert.Equal(t, metadataProxyID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyMetadataProxy()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMetadataProxyRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining PolicyMetadataProxy ID")
	})
}

func TestMockResourceNsxtPolicyMetadataProxyUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetaSDK := metaproxy.NewMockMetadataProxiesClient(ctrl)
	mockWrapper := &cliinfra.MetadataProxyConfigClientContext{
		Client:     mockMetaSDK,
		ClientType: utl.Local,
	}

	originalCli := cliMetadataProxiesClient
	defer func() { cliMetadataProxiesClient = originalCli }()
	cliMetadataProxiesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.MetadataProxyConfigClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockMetaSDK.EXPECT().Update(metadataProxyID, gomock.Any()).Return(model.MetadataProxyConfig{
			Id:              &metadataProxyID,
			DisplayName:     &metadataProxyDisplayName,
			Description:     &metadataProxyDescription,
			Path:            &metadataProxyPath,
			Revision:        &metadataProxyRevision,
			EdgeClusterPath: &metadataProxyEdgePath,
			ServerAddress:   &metadataProxyServerAddr,
		}, nil)
		mockMetaSDK.EXPECT().Get(metadataProxyID).Return(model.MetadataProxyConfig{
			Id:              &metadataProxyID,
			DisplayName:     &metadataProxyDisplayName,
			Description:     &metadataProxyDescription,
			Path:            &metadataProxyPath,
			Revision:        &metadataProxyRevision,
			EdgeClusterPath: &metadataProxyEdgePath,
			ServerAddress:   &metadataProxyServerAddr,
		}, nil)

		res := resourceNsxtPolicyMetadataProxy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalMetadataProxyData())
		d.Set("revision", 1)
		d.SetId(metadataProxyID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMetadataProxyUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyMetadataProxy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalMetadataProxyData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMetadataProxyUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining PolicyMetadataProxy ID")
	})
}

func TestMockResourceNsxtPolicyMetadataProxyDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetaSDK := metaproxy.NewMockMetadataProxiesClient(ctrl)
	mockWrapper := &cliinfra.MetadataProxyConfigClientContext{
		Client:     mockMetaSDK,
		ClientType: utl.Local,
	}

	originalCli := cliMetadataProxiesClient
	defer func() { cliMetadataProxiesClient = originalCli }()
	cliMetadataProxiesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.MetadataProxyConfigClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockMetaSDK.EXPECT().Delete(metadataProxyID).Return(nil)

		res := resourceNsxtPolicyMetadataProxy()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(metadataProxyID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMetadataProxyDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyMetadataProxy()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMetadataProxyDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining PolicyMetadataProxy ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockMetaSDK.EXPECT().Delete(metadataProxyID).Return(errors.New("API error"))

		res := resourceNsxtPolicyMetadataProxy()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(metadataProxyID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMetadataProxyDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalMetadataProxyData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":      metadataProxyDisplayName,
		"description":       metadataProxyDescription,
		"edge_cluster_path": metadataProxyEdgePath,
		"secret":            metadataProxySecret,
		"server_address":    metadataProxyServerAddr,
	}
}
