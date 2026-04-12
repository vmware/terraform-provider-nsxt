//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/IpBlocksClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/IpBlocksClient.go IpBlocksClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	ipblockmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	ipBlockID          = "ip-block-1"
	ipBlockDisplayName = "ip-block-fooname"
	ipBlockDescription = "ip block mock"
	ipBlockPath        = "/infra/ip-blocks/ip-block-1"
	ipBlockRevision    = int64(1)
	ipBlockCidr        = "192.168.1.0/24"
)

func TestMockResourceNsxtPolicyIPBlockCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlockSDK := ipblockmocks.NewMockIpBlocksClient(ctrl)
	mockWrapper := &cliinfra.IpAddressBlockClientContext{
		Client:     mockBlockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpBlocksClient
	defer func() { cliIpBlocksClient = originalCli }()
	cliIpBlocksClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IpAddressBlockClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		// When nsx_id is empty, getOrGenerateID2 returns newUUID() without calling exists
		mockBlockSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockBlockSDK.EXPECT().Get(gomock.Any(), nil).Return(model.IpAddressBlock{
			Id:          &ipBlockID,
			DisplayName: &ipBlockDisplayName,
			Description: &ipBlockDescription,
			Path:        &ipBlockPath,
			Revision:    &ipBlockRevision,
			Cidr:        &ipBlockCidr,
		}, nil)

		res := resourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPBlockData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPBlockCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		// Read sets nsx_id from API response (block.Id), which may differ from d.Id() after Create
		assert.NotEmpty(t, d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockBlockSDK.EXPECT().Get("existing-id", nil).Return(model.IpAddressBlock{Id: &ipBlockID}, nil)

		res := resourceNsxtPolicyIPBlock()
		data := minimalIPBlockData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPBlockCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyIPBlockRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlockSDK := ipblockmocks.NewMockIpBlocksClient(ctrl)
	mockWrapper := &cliinfra.IpAddressBlockClientContext{
		Client:     mockBlockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpBlocksClient
	defer func() { cliIpBlocksClient = originalCli }()
	cliIpBlocksClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IpAddressBlockClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockBlockSDK.EXPECT().Get(ipBlockID, nil).Return(model.IpAddressBlock{
			Id:          &ipBlockID,
			DisplayName: &ipBlockDisplayName,
			Description: &ipBlockDescription,
			Path:        &ipBlockPath,
			Revision:    &ipBlockRevision,
			Cidr:        &ipBlockCidr,
		}, nil)

		res := resourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ipBlockID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPBlockRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, ipBlockDisplayName, d.Get("display_name"))
		assert.Equal(t, ipBlockDescription, d.Get("description"))
		assert.Equal(t, ipBlockPath, d.Get("path"))
		assert.Equal(t, int(ipBlockRevision), d.Get("revision"))
		assert.Equal(t, ipBlockID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPBlockRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IP Block ID")
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockBlockSDK.EXPECT().Get(ipBlockID, nil).Return(model.IpAddressBlock{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ipBlockID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPBlockRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})
}

func TestMockResourceNsxtPolicyIPBlockUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlockSDK := ipblockmocks.NewMockIpBlocksClient(ctrl)
	mockWrapper := &cliinfra.IpAddressBlockClientContext{
		Client:     mockBlockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpBlocksClient
	defer func() { cliIpBlocksClient = originalCli }()
	cliIpBlocksClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IpAddressBlockClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockBlockSDK.EXPECT().Update(ipBlockID, gomock.Any()).Return(model.IpAddressBlock{
			Id:          &ipBlockID,
			DisplayName: &ipBlockDisplayName,
			Description: &ipBlockDescription,
			Path:        &ipBlockPath,
			Revision:    &ipBlockRevision,
			Cidr:        &ipBlockCidr,
		}, nil)
		mockBlockSDK.EXPECT().Get(ipBlockID, nil).Return(model.IpAddressBlock{
			Id:          &ipBlockID,
			DisplayName: &ipBlockDisplayName,
			Description: &ipBlockDescription,
			Path:        &ipBlockPath,
			Revision:    &ipBlockRevision,
			Cidr:        &ipBlockCidr,
		}, nil)

		res := resourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPBlockData())
		d.Set("revision", 1)
		d.SetId(ipBlockID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPBlockUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPBlockData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPBlockUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IP Block ID")
	})
}

func TestMockResourceNsxtPolicyIPBlockDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBlockSDK := ipblockmocks.NewMockIpBlocksClient(ctrl)
	mockWrapper := &cliinfra.IpAddressBlockClientContext{
		Client:     mockBlockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpBlocksClient
	defer func() { cliIpBlocksClient = originalCli }()
	cliIpBlocksClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IpAddressBlockClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockBlockSDK.EXPECT().Delete(ipBlockID).Return(nil)

		res := resourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ipBlockID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPBlockDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPBlockDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IP Block ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockBlockSDK.EXPECT().Delete(ipBlockID).Return(errors.New("API error"))

		res := resourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ipBlockID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPBlockDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalIPBlockData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": ipBlockDisplayName,
		"description":  ipBlockDescription,
		"cidr":         ipBlockCidr,
	}
}
