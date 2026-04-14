//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	ipblockmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

var (
	dsIPBlockID          = "ip-block-ds-1"
	dsIPBlockName        = "test-ip-block-ds"
	dsIPBlockPath        = "/infra/ip-blocks/ip-block-ds-1"
	dsIPBlockDescription = "ds test"
)

func dsIPBlockModel() nsxModel.IpAddressBlock {
	return nsxModel.IpAddressBlock{
		Id:          &dsIPBlockID,
		DisplayName: &dsIPBlockName,
		Description: &dsIPBlockDescription,
		Path:        &dsIPBlockPath,
	}
}

func setupIPBlockDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*ipblockmocks.MockIpBlocksClient, func()) {
	t.Helper()
	mockSDK := ipblockmocks.NewMockIpBlocksClient(ctrl)
	wrapper := &cliinfra.IpAddressBlockClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliIpBlocksClient
	cliIpBlocksClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *cliinfra.IpAddressBlockClientContext {
		return wrapper
	}
	return mockSDK, func() { cliIpBlocksClient = orig }
}

func TestMockDataSourceNsxtPolicyIPBlockRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIPBlockDataSourceMock(t, ctrl)
	defer restore()

	t.Run("by ID success", func(t *testing.T) {
		mockSDK.EXPECT().Get(dsIPBlockID, nil).Return(dsIPBlockModel(), nil)

		ds := dataSourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": dsIPBlockID,
		})

		err := dataSourceNsxtPolicyIPBlockRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dsIPBlockID, d.Id())
		assert.Equal(t, dsIPBlockName, d.Get("display_name"))
	})

	t.Run("by ID not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(dsIPBlockID, nil).Return(nsxModel.IpAddressBlock{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": dsIPBlockID,
		})

		err := dataSourceNsxtPolicyIPBlockRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("by ID API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(dsIPBlockID, nil).Return(nsxModel.IpAddressBlock{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": dsIPBlockID,
		})

		err := dataSourceNsxtPolicyIPBlockRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("by display_name single match", func(t *testing.T) {
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil).Return(nsxModel.IpAddressBlockListResult{
			Results: []nsxModel.IpAddressBlock{dsIPBlockModel()},
		}, nil)

		ds := dataSourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": dsIPBlockName,
		})

		err := dataSourceNsxtPolicyIPBlockRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dsIPBlockID, d.Id())
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil).Return(nsxModel.IpAddressBlockListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": dsIPBlockName,
		})

		err := dataSourceNsxtPolicyIPBlockRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to read IpAddressBlocks")
	})

	t.Run("missing id and name", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyIPBlockRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ID or name")
	})

	t.Run("by name not found", func(t *testing.T) {
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil).Return(nsxModel.IpAddressBlockListResult{
			Results: []nsxModel.IpAddressBlock{},
		}, nil)

		ds := dataSourceNsxtPolicyIPBlock()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "missing",
		})

		err := dataSourceNsxtPolicyIPBlockRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})
}
