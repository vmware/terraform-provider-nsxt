//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"go.uber.org/mock/gomock"

	fabricapi "github.com/vmware/terraform-provider-nsxt/api/nsx/fabric"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	fabricmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/fabric"
)

var (
	dnTestExternalID = "dn-ext-id-1"
	dnTestIPAddress  = "192.168.1.100"
)

func discoveredNodeModel() nsxModel.DiscoveredNode {
	return nsxModel.DiscoveredNode{
		ExternalId:  &dnTestExternalID,
		IpAddresses: []string{dnTestIPAddress},
	}
}

func setupDiscoveredNodeMock(t *testing.T, ctrl *gomock.Controller) (*fabricmocks.MockDiscoveredNodesClient, func()) {
	t.Helper()
	mockSDK := fabricmocks.NewMockDiscoveredNodesClient(ctrl)
	wrapper := &fabricapi.DiscoveredNodeClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliDiscoveredNodesClient
	cliDiscoveredNodesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *fabricapi.DiscoveredNodeClientContext {
		return wrapper
	}
	return mockSDK, func() { cliDiscoveredNodesClient = orig }
}

func TestMockDataSourceNsxtDiscoveredNodeRead(t *testing.T) {
	t.Run("by ID success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDiscoveredNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(dnTestExternalID).Return(discoveredNodeModel(), nil)

		ds := dataSourceNsxtDiscoveredNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": dnTestExternalID,
		})

		err := dataSourceNsxtDiscoveredNodeRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dnTestExternalID, d.Id())
		assert.Equal(t, dnTestIPAddress, d.Get("ip_address"))
	})

	t.Run("by ID not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDiscoveredNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(dnTestExternalID).Return(nsxModel.DiscoveredNode{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtDiscoveredNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": dnTestExternalID,
		})

		err := dataSourceNsxtDiscoveredNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by IP address success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDiscoveredNodeMock(t, ctrl)
		defer restore()

		ip := dnTestIPAddress
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, &ip, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.DiscoveredNodeListResult{
			Results: []nsxModel.DiscoveredNode{discoveredNodeModel()},
		}, nil)

		ds := dataSourceNsxtDiscoveredNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"ip_address": dnTestIPAddress,
		})

		err := dataSourceNsxtDiscoveredNodeRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dnTestExternalID, d.Id())
	})

	t.Run("by IP address list error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDiscoveredNodeMock(t, ctrl)
		defer restore()

		ip := dnTestIPAddress
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, &ip, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.DiscoveredNodeListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtDiscoveredNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"ip_address": dnTestIPAddress,
		})

		err := dataSourceNsxtDiscoveredNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading Discovered Node")
	})

	t.Run("by IP address not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDiscoveredNodeMock(t, ctrl)
		defer restore()

		ip := dnTestIPAddress
		mockSDK.EXPECT().List(nil, nil, nil, nil, nil, nil, &ip, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.DiscoveredNodeListResult{
			Results: []nsxModel.DiscoveredNode{},
		}, nil)

		ds := dataSourceNsxtDiscoveredNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"ip_address": dnTestIPAddress,
		})

		err := dataSourceNsxtDiscoveredNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("missing id and ip_address", func(t *testing.T) {
		ds := dataSourceNsxtDiscoveredNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtDiscoveredNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "external ID or IP address")
	})
}
