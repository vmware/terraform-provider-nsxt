//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/DistributedVlanConnectionsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/DistributedVlanConnectionsClient.go DistributedVlanConnectionsClient

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
	dvcmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	distributedVlanConnID          = "dvc-1"
	distributedVlanConnDisplayName = "dvc-fooname"
	distributedVlanConnDescription = "distributed vlan connection mock"
	distributedVlanConnPath        = "/infra/distributed-vlan-connections/dvc-1"
	distributedVlanConnRevision    = int64(1)
	distributedVlanConnVlanID      = int64(100)
)

func TestMockResourceNsxtPolicyDistributedVlanConnectionCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDvcSDK := dvcmocks.NewMockDistributedVlanConnectionsClient(ctrl)
	mockWrapper := &cliinfra.DistributedVlanConnectionClientContext{
		Client:     mockDvcSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDistributedVlanConnectionsClient
	defer func() { cliDistributedVlanConnectionsClient = originalCli }()
	cliDistributedVlanConnectionsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DistributedVlanConnectionClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockDvcSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockDvcSDK.EXPECT().Get(gomock.Any()).Return(model.DistributedVlanConnection{
			Id:          &distributedVlanConnID,
			DisplayName: &distributedVlanConnDisplayName,
			Description: &distributedVlanConnDescription,
			Path:        &distributedVlanConnPath,
			Revision:    &distributedVlanConnRevision,
			VlanId:      &distributedVlanConnVlanID,
		}, nil)

		res := resourceNsxtPolicyDistributedVlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDistributedVlanConnectionData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVlanConnectionCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockDvcSDK.EXPECT().Get("existing-id").Return(model.DistributedVlanConnection{Id: &distributedVlanConnID}, nil)

		res := resourceNsxtPolicyDistributedVlanConnection()
		data := minimalDistributedVlanConnectionData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVlanConnectionCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyDistributedVlanConnectionRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDvcSDK := dvcmocks.NewMockDistributedVlanConnectionsClient(ctrl)
	mockWrapper := &cliinfra.DistributedVlanConnectionClientContext{
		Client:     mockDvcSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDistributedVlanConnectionsClient
	defer func() { cliDistributedVlanConnectionsClient = originalCli }()
	cliDistributedVlanConnectionsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DistributedVlanConnectionClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockDvcSDK.EXPECT().Get(distributedVlanConnID).Return(model.DistributedVlanConnection{
			Id:          &distributedVlanConnID,
			DisplayName: &distributedVlanConnDisplayName,
			Description: &distributedVlanConnDescription,
			Path:        &distributedVlanConnPath,
			Revision:    &distributedVlanConnRevision,
			VlanId:      &distributedVlanConnVlanID,
		}, nil)

		res := resourceNsxtPolicyDistributedVlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(distributedVlanConnID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVlanConnectionRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, distributedVlanConnDisplayName, d.Get("display_name"))
		assert.Equal(t, distributedVlanConnDescription, d.Get("description"))
		assert.Equal(t, distributedVlanConnPath, d.Get("path"))
		assert.Equal(t, int(distributedVlanConnRevision), d.Get("revision"))
		assert.Equal(t, distributedVlanConnID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDistributedVlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVlanConnectionRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining DistributedVlanConnection ID")
	})
}

func TestMockResourceNsxtPolicyDistributedVlanConnectionUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDvcSDK := dvcmocks.NewMockDistributedVlanConnectionsClient(ctrl)
	mockWrapper := &cliinfra.DistributedVlanConnectionClientContext{
		Client:     mockDvcSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDistributedVlanConnectionsClient
	defer func() { cliDistributedVlanConnectionsClient = originalCli }()
	cliDistributedVlanConnectionsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DistributedVlanConnectionClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockDvcSDK.EXPECT().Update(distributedVlanConnID, gomock.Any()).Return(model.DistributedVlanConnection{
			Id:          &distributedVlanConnID,
			DisplayName: &distributedVlanConnDisplayName,
			Description: &distributedVlanConnDescription,
			Path:        &distributedVlanConnPath,
			Revision:    &distributedVlanConnRevision,
			VlanId:      &distributedVlanConnVlanID,
		}, nil)
		mockDvcSDK.EXPECT().Get(distributedVlanConnID).Return(model.DistributedVlanConnection{
			Id:          &distributedVlanConnID,
			DisplayName: &distributedVlanConnDisplayName,
			Description: &distributedVlanConnDescription,
			Path:        &distributedVlanConnPath,
			Revision:    &distributedVlanConnRevision,
			VlanId:      &distributedVlanConnVlanID,
		}, nil)

		res := resourceNsxtPolicyDistributedVlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDistributedVlanConnectionData())
		d.Set("revision", 1)
		d.SetId(distributedVlanConnID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVlanConnectionUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDistributedVlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDistributedVlanConnectionData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVlanConnectionUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining DistributedVlanConnection ID")
	})
}

func TestMockResourceNsxtPolicyDistributedVlanConnectionDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDvcSDK := dvcmocks.NewMockDistributedVlanConnectionsClient(ctrl)
	mockWrapper := &cliinfra.DistributedVlanConnectionClientContext{
		Client:     mockDvcSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDistributedVlanConnectionsClient
	defer func() { cliDistributedVlanConnectionsClient = originalCli }()
	cliDistributedVlanConnectionsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DistributedVlanConnectionClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockDvcSDK.EXPECT().Delete(distributedVlanConnID).Return(nil)

		res := resourceNsxtPolicyDistributedVlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(distributedVlanConnID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVlanConnectionDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDistributedVlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVlanConnectionDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining DistributedVlanConnection ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockDvcSDK.EXPECT().Delete(distributedVlanConnID).Return(errors.New("API error"))

		res := resourceNsxtPolicyDistributedVlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(distributedVlanConnID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVlanConnectionDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalDistributedVlanConnectionData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": distributedVlanConnDisplayName,
		"description":  distributedVlanConnDescription,
		"vlan_id":      int(distributedVlanConnVlanID),
	}
}
