//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/DistributedVxlanConnectionsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/DistributedVxlanConnectionsClient.go DistributedVxlanConnectionsClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	dvxcmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	distributedVxlanConnID          = "dvxc-1"
	distributedVxlanConnDisplayName = "dvxc-fooname"
	distributedVxlanConnDescription = "distributed vxlan connection mock"
	distributedVxlanConnPath        = "/infra/distributed-vxlan-connections/dvxc-1"
	distributedVxlanConnRevision    = int64(1)
	distributedVxlanConnL3Vni       = int64(5000)
)

func TestMockResourceNsxtPolicyDistributedVxlanConnectionCreate(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDvxcSDK := dvxcmocks.NewMockDistributedVxlanConnectionsClient(ctrl)
	mockWrapper := &cliinfra.DistributedVxlanConnectionClientContext{
		Client:     mockDvxcSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDistributedVxlanConnectionsClient
	defer func() { cliDistributedVxlanConnectionsClient = originalCli }()
	cliDistributedVxlanConnectionsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DistributedVxlanConnectionClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockDvxcSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockDvxcSDK.EXPECT().Get(gomock.Any()).Return(model.DistributedVxlanConnection{
			Id:          &distributedVxlanConnID,
			DisplayName: &distributedVxlanConnDisplayName,
			Description: &distributedVxlanConnDescription,
			Path:        &distributedVxlanConnPath,
			Revision:    &distributedVxlanConnRevision,
			L3Vni:       &distributedVxlanConnL3Vni,
		}, nil)

		res := resourceNsxtPolicyDistributedVxlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDistributedVxlanConnectionData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVxlanConnectionCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockDvxcSDK.EXPECT().Get("existing-id").Return(model.DistributedVxlanConnection{Id: &distributedVxlanConnID}, nil)

		res := resourceNsxtPolicyDistributedVxlanConnection()
		data := minimalDistributedVxlanConnectionData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVxlanConnectionCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyDistributedVxlanConnectionRead(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDvxcSDK := dvxcmocks.NewMockDistributedVxlanConnectionsClient(ctrl)
	mockWrapper := &cliinfra.DistributedVxlanConnectionClientContext{
		Client:     mockDvxcSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDistributedVxlanConnectionsClient
	defer func() { cliDistributedVxlanConnectionsClient = originalCli }()
	cliDistributedVxlanConnectionsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DistributedVxlanConnectionClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockDvxcSDK.EXPECT().Get(distributedVxlanConnID).Return(model.DistributedVxlanConnection{
			Id:          &distributedVxlanConnID,
			DisplayName: &distributedVxlanConnDisplayName,
			Description: &distributedVxlanConnDescription,
			Path:        &distributedVxlanConnPath,
			Revision:    &distributedVxlanConnRevision,
			L3Vni:       &distributedVxlanConnL3Vni,
		}, nil)

		res := resourceNsxtPolicyDistributedVxlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(distributedVxlanConnID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVxlanConnectionRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, distributedVxlanConnDisplayName, d.Get("display_name"))
		assert.Equal(t, distributedVxlanConnDescription, d.Get("description"))
		assert.Equal(t, distributedVxlanConnPath, d.Get("path"))
		assert.Equal(t, int(distributedVxlanConnRevision), d.Get("revision"))
		assert.Equal(t, distributedVxlanConnID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDistributedVxlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVxlanConnectionRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining DistributedVxlanConnection ID")
	})
}

func TestMockResourceNsxtPolicyDistributedVxlanConnectionUpdate(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDvxcSDK := dvxcmocks.NewMockDistributedVxlanConnectionsClient(ctrl)
	mockWrapper := &cliinfra.DistributedVxlanConnectionClientContext{
		Client:     mockDvxcSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDistributedVxlanConnectionsClient
	defer func() { cliDistributedVxlanConnectionsClient = originalCli }()
	cliDistributedVxlanConnectionsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DistributedVxlanConnectionClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockDvxcSDK.EXPECT().Update(distributedVxlanConnID, gomock.Any()).Return(model.DistributedVxlanConnection{
			Id:          &distributedVxlanConnID,
			DisplayName: &distributedVxlanConnDisplayName,
			Description: &distributedVxlanConnDescription,
			Path:        &distributedVxlanConnPath,
			Revision:    &distributedVxlanConnRevision,
			L3Vni:       &distributedVxlanConnL3Vni,
		}, nil)
		mockDvxcSDK.EXPECT().Get(distributedVxlanConnID).Return(model.DistributedVxlanConnection{
			Id:          &distributedVxlanConnID,
			DisplayName: &distributedVxlanConnDisplayName,
			Description: &distributedVxlanConnDescription,
			Path:        &distributedVxlanConnPath,
			Revision:    &distributedVxlanConnRevision,
			L3Vni:       &distributedVxlanConnL3Vni,
		}, nil)

		res := resourceNsxtPolicyDistributedVxlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDistributedVxlanConnectionData())
		d.Set("revision", 1)
		d.SetId(distributedVxlanConnID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVxlanConnectionUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDistributedVxlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDistributedVxlanConnectionData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVxlanConnectionUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining DistributedVxlanConnection ID")
	})
}

func TestMockResourceNsxtPolicyDistributedVxlanConnectionDelete(t *testing.T) {
	util.NsxVersion = "4.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDvxcSDK := dvxcmocks.NewMockDistributedVxlanConnectionsClient(ctrl)
	mockWrapper := &cliinfra.DistributedVxlanConnectionClientContext{
		Client:     mockDvxcSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDistributedVxlanConnectionsClient
	defer func() { cliDistributedVxlanConnectionsClient = originalCli }()
	cliDistributedVxlanConnectionsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DistributedVxlanConnectionClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockDvxcSDK.EXPECT().Delete(distributedVxlanConnID).Return(nil)

		res := resourceNsxtPolicyDistributedVxlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(distributedVxlanConnID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVxlanConnectionDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDistributedVxlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVxlanConnectionDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining DistributedVxlanConnection ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockDvxcSDK.EXPECT().Delete(distributedVxlanConnID).Return(errors.New("API error"))

		res := resourceNsxtPolicyDistributedVxlanConnection()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(distributedVxlanConnID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDistributedVxlanConnectionDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalDistributedVxlanConnectionData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": distributedVxlanConnDisplayName,
		"description":  distributedVxlanConnDescription,
		"l3_vni":       int(distributedVxlanConnL3Vni),
	}
}
