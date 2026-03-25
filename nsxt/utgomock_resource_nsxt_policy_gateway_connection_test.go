// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/GatewayConnectionsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/GatewayConnectionsClient.go GatewayConnectionsClient

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
	gwconnMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	gatewayConnID          = "gateway-conn-1"
	gatewayConnDisplayName = "gateway-conn-fooname"
	gatewayConnDescription = "gateway connection mock"
	gatewayConnPath        = "/infra/gateway-connections/gateway-conn-1"
	gatewayConnRevision    = int64(1)
	gatewayConnTier0Path   = "/infra/tier-0s/t0-1"
)

func TestMockResourceNsxtPolicyGatewayConnectionCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGwSDK := gwconnMocks.NewMockGatewayConnectionsClient(ctrl)
	mockWrapper := &cliinfra.GatewayConnectionClientContext{
		Client:     mockGwSDK,
		ClientType: utl.Local,
	}

	originalCli := cliGatewayConnectionsClient
	defer func() { cliGatewayConnectionsClient = originalCli }()
	cliGatewayConnectionsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.GatewayConnectionClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockGwSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockGwSDK.EXPECT().Get(gomock.Any()).Return(model.GatewayConnection{
			Id:          &gatewayConnID,
			DisplayName: &gatewayConnDisplayName,
			Description: &gatewayConnDescription,
			Path:        &gatewayConnPath,
			Revision:    &gatewayConnRevision,
			Tier0Path:   &gatewayConnTier0Path,
		}, nil)

		res := resourceNsxtPolicyGatewayConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGatewayConnectionData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayConnectionCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockGwSDK.EXPECT().Get("existing-id").Return(model.GatewayConnection{Id: &gatewayConnID}, nil)

		res := resourceNsxtPolicyGatewayConnection()
		data := minimalGatewayConnectionData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayConnectionCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails when NSX version below 9.0.0", func(t *testing.T) {
		util.NsxVersion = "8.0.0"
		defer func() { util.NsxVersion = "9.1.0" }()

		res := resourceNsxtPolicyGatewayConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGatewayConnectionData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayConnectionCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.0.0")
	})
}

func TestMockResourceNsxtPolicyGatewayConnectionRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGwSDK := gwconnMocks.NewMockGatewayConnectionsClient(ctrl)
	mockWrapper := &cliinfra.GatewayConnectionClientContext{
		Client:     mockGwSDK,
		ClientType: utl.Local,
	}

	originalCli := cliGatewayConnectionsClient
	defer func() { cliGatewayConnectionsClient = originalCli }()
	cliGatewayConnectionsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.GatewayConnectionClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockGwSDK.EXPECT().Get(gatewayConnID).Return(model.GatewayConnection{
			Id:          &gatewayConnID,
			DisplayName: &gatewayConnDisplayName,
			Description: &gatewayConnDescription,
			Path:        &gatewayConnPath,
			Revision:    &gatewayConnRevision,
			Tier0Path:   &gatewayConnTier0Path,
		}, nil)

		res := resourceNsxtPolicyGatewayConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(gatewayConnID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayConnectionRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, gatewayConnDisplayName, d.Get("display_name"))
		assert.Equal(t, gatewayConnDescription, d.Get("description"))
		assert.Equal(t, gatewayConnPath, d.Get("path"))
		assert.Equal(t, int(gatewayConnRevision), d.Get("revision"))
		assert.Equal(t, gatewayConnID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayConnectionRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining GatewayConnection ID")
	})
}

func TestMockResourceNsxtPolicyGatewayConnectionUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGwSDK := gwconnMocks.NewMockGatewayConnectionsClient(ctrl)
	mockWrapper := &cliinfra.GatewayConnectionClientContext{
		Client:     mockGwSDK,
		ClientType: utl.Local,
	}

	originalCli := cliGatewayConnectionsClient
	defer func() { cliGatewayConnectionsClient = originalCli }()
	cliGatewayConnectionsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.GatewayConnectionClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockGwSDK.EXPECT().Update(gatewayConnID, gomock.Any()).Return(model.GatewayConnection{
			Id:          &gatewayConnID,
			DisplayName: &gatewayConnDisplayName,
			Description: &gatewayConnDescription,
			Path:        &gatewayConnPath,
			Revision:    &gatewayConnRevision,
			Tier0Path:   &gatewayConnTier0Path,
		}, nil)
		mockGwSDK.EXPECT().Get(gatewayConnID).Return(model.GatewayConnection{
			Id:          &gatewayConnID,
			DisplayName: &gatewayConnDisplayName,
			Description: &gatewayConnDescription,
			Path:        &gatewayConnPath,
			Revision:    &gatewayConnRevision,
			Tier0Path:   &gatewayConnTier0Path,
		}, nil)

		res := resourceNsxtPolicyGatewayConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGatewayConnectionData())
		d.Set("revision", 1)
		d.SetId(gatewayConnID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayConnectionUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGatewayConnectionData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayConnectionUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining GatewayConnection ID")
	})
}

func TestMockResourceNsxtPolicyGatewayConnectionDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGwSDK := gwconnMocks.NewMockGatewayConnectionsClient(ctrl)
	mockWrapper := &cliinfra.GatewayConnectionClientContext{
		Client:     mockGwSDK,
		ClientType: utl.Local,
	}

	originalCli := cliGatewayConnectionsClient
	defer func() { cliGatewayConnectionsClient = originalCli }()
	cliGatewayConnectionsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.GatewayConnectionClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockGwSDK.EXPECT().Delete(gatewayConnID).Return(nil)

		res := resourceNsxtPolicyGatewayConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(gatewayConnID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayConnectionDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayConnectionDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining GatewayConnection ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockGwSDK.EXPECT().Delete(gatewayConnID).Return(errors.New("API error"))

		res := resourceNsxtPolicyGatewayConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(gatewayConnID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyGatewayConnectionDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalGatewayConnectionData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": gatewayConnDisplayName,
		"description":  gatewayConnDescription,
		"tier0_path":   gatewayConnTier0Path,
	}
}
