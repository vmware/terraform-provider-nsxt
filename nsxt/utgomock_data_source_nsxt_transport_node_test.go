//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mpmodel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"go.uber.org/mock/gomock"
)

// setupTransportNodeMock, transportNodeAPIResponse, tnID, tnDisplayName,
// tnDescription are all defined in utgomock_resource_nsxt_edge_transport_node_test.go
// and are reused here.

func TestUnitNsxt_DataSourceNsxtTransportNodeRead(t *testing.T) {
	t.Run("by id API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(tnID).Return(mpmodel.TransportNode{}, errors.New("get failed"))

		ds := dataSourceNsxtTransportNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": tnID,
		})

		err := dataSourceNsxtTransportNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read TransportNode")
	})

	t.Run("by id success (EdgeNode)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(tnID).Return(transportNodeAPIResponse(), nil)

		ds := dataSourceNsxtTransportNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": tnID,
		})

		err := dataSourceNsxtTransportNodeRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tnID, d.Id())
		assert.Equal(t, tnDisplayName, d.Get("display_name"))
	})

	t.Run("by display_name perfect match (filters to EdgeNodes)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mpmodel.TransportNodeListResult{Results: []mpmodel.TransportNode{transportNodeAPIResponse()}}, nil)

		ds := dataSourceNsxtTransportNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": tnDisplayName,
		})

		err := dataSourceNsxtTransportNodeRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tnID, d.Id())
		assert.Equal(t, tnDisplayName, d.Get("display_name"))
	})

	t.Run("by empty display_name grabs single result", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mpmodel.TransportNodeListResult{Results: []mpmodel.TransportNode{transportNodeAPIResponse()}}, nil)

		ds := dataSourceNsxtTransportNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtTransportNodeRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tnID, d.Id())
	})

	t.Run("List API error is wrapped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mpmodel.TransportNodeListResult{}, errors.New("list failed"))

		ds := dataSourceNsxtTransportNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": tnDisplayName,
		})

		err := dataSourceNsxtTransportNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read Transport Nodes")
	})

	t.Run("multiple perfect matches errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		second := transportNodeAPIResponse()
		secondID := "transport-node-uuid-2"
		second.Id = &secondID

		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mpmodel.TransportNodeListResult{Results: []mpmodel.TransportNode{transportNodeAPIResponse(), second}}, nil)

		ds := dataSourceNsxtTransportNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": tnDisplayName,
		})

		err := dataSourceNsxtTransportNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple Transport Nodes")
	})

	t.Run("no match errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTransportNodeMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mpmodel.TransportNodeListResult{Results: []mpmodel.TransportNode{transportNodeAPIResponse()}}, nil)

		ds := dataSourceNsxtTransportNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		err := dataSourceNsxtTransportNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no Transport Node matches")
	})
}
