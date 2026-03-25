// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/DhcpServerConfigsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/DhcpServerConfigsClient.go DhcpServerConfigsClient

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
	dhcpsvcmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	dhcpServerID          = "dhcp-server-1"
	dhcpServerDisplayName = "dhcp-server-fooname"
	dhcpServerDescription = "dhcp server mock"
	dhcpServerPath        = "/infra/dhcp-server-configs/dhcp-server-1"
	dhcpServerRevision    = int64(1)
	dhcpServerLeaseTime   = int64(86400)
	dhcpServerEdgePath    = "/infra/sites/default/enforcement-points/default/edge-clusters/ec1"
)

func TestMockResourceNsxtPolicyDhcpServerCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDhcpSDK := dhcpsvcmocks.NewMockDhcpServerConfigsClient(ctrl)
	mockWrapper := &cliinfra.DhcpServerConfigClientContext{
		Client:     mockDhcpSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDhcpServerConfigsClient
	defer func() { cliDhcpServerConfigsClient = originalCli }()
	cliDhcpServerConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DhcpServerConfigClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockDhcpSDK.EXPECT().Get(gomock.Any()).Return(model.DhcpServerConfig{
			Id:              &dhcpServerID,
			DisplayName:     &dhcpServerDisplayName,
			Description:     &dhcpServerDescription,
			Path:            &dhcpServerPath,
			Revision:        &dhcpServerRevision,
			LeaseTime:       &dhcpServerLeaseTime,
			EdgeClusterPath: &dhcpServerEdgePath,
		}, nil)

		res := resourceNsxtPolicyDhcpServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpServerData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpServerCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Get("existing-id").Return(model.DhcpServerConfig{Id: &dhcpServerID}, nil)

		res := resourceNsxtPolicyDhcpServer()
		data := minimalDhcpServerData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpServerCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyDhcpServerRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDhcpSDK := dhcpsvcmocks.NewMockDhcpServerConfigsClient(ctrl)
	mockWrapper := &cliinfra.DhcpServerConfigClientContext{
		Client:     mockDhcpSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDhcpServerConfigsClient
	defer func() { cliDhcpServerConfigsClient = originalCli }()
	cliDhcpServerConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DhcpServerConfigClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Get(dhcpServerID).Return(model.DhcpServerConfig{
			Id:              &dhcpServerID,
			DisplayName:     &dhcpServerDisplayName,
			Description:     &dhcpServerDescription,
			Path:            &dhcpServerPath,
			Revision:        &dhcpServerRevision,
			LeaseTime:       &dhcpServerLeaseTime,
			EdgeClusterPath: &dhcpServerEdgePath,
		}, nil)

		res := resourceNsxtPolicyDhcpServer()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dhcpServerID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpServerRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, dhcpServerDisplayName, d.Get("display_name"))
		assert.Equal(t, dhcpServerDescription, d.Get("description"))
		assert.Equal(t, dhcpServerPath, d.Get("path"))
		assert.Equal(t, int(dhcpServerRevision), d.Get("revision"))
		assert.Equal(t, dhcpServerID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDhcpServer()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpServerRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining DhcpServer ID")
	})
}

func TestMockResourceNsxtPolicyDhcpServerUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDhcpSDK := dhcpsvcmocks.NewMockDhcpServerConfigsClient(ctrl)
	mockWrapper := &cliinfra.DhcpServerConfigClientContext{
		Client:     mockDhcpSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDhcpServerConfigsClient
	defer func() { cliDhcpServerConfigsClient = originalCli }()
	cliDhcpServerConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DhcpServerConfigClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Patch(dhcpServerID, gomock.Any()).Return(nil)
		mockDhcpSDK.EXPECT().Get(dhcpServerID).Return(model.DhcpServerConfig{
			Id:              &dhcpServerID,
			DisplayName:     &dhcpServerDisplayName,
			Description:     &dhcpServerDescription,
			Path:            &dhcpServerPath,
			Revision:        &dhcpServerRevision,
			LeaseTime:       &dhcpServerLeaseTime,
			EdgeClusterPath: &dhcpServerEdgePath,
		}, nil)

		res := resourceNsxtPolicyDhcpServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpServerData())
		d.SetId(dhcpServerID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpServerUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDhcpServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpServerData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpServerUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining DhcpServer ID")
	})
}

func TestMockResourceNsxtPolicyDhcpServerDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDhcpSDK := dhcpsvcmocks.NewMockDhcpServerConfigsClient(ctrl)
	mockWrapper := &cliinfra.DhcpServerConfigClientContext{
		Client:     mockDhcpSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDhcpServerConfigsClient
	defer func() { cliDhcpServerConfigsClient = originalCli }()
	cliDhcpServerConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DhcpServerConfigClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Delete(dhcpServerID).Return(nil)

		res := resourceNsxtPolicyDhcpServer()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dhcpServerID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpServerDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDhcpServer()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpServerDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining DhcpServer ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockDhcpSDK.EXPECT().Delete(dhcpServerID).Return(errors.New("API error"))

		res := resourceNsxtPolicyDhcpServer()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dhcpServerID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpServerDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalDhcpServerData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":      dhcpServerDisplayName,
		"description":       dhcpServerDescription,
		"edge_cluster_path": dhcpServerEdgePath,
		"lease_time":        86400,
	}
}
