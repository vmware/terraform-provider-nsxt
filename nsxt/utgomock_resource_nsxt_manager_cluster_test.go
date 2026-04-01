//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/nsx/ClusterClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/ClusterClient.go ClusterClient
//
// Note: resourceNsxtManagerClusterCreate and resourceNsxtManagerClusterUpdate involve
// network connectivity probing and joining cluster nodes, which requires live NSX
// infrastructure. Only Read and Delete are unit-testable via mock injection.

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"

	nsxmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx"
)

var (
	mcClusterID  = "cluster-uuid-1"
	mcNodeUUID1  = "node-uuid-1"
	mcNodeIP1    = "192.168.1.10"
	mcNodeFQDN1  = "nsxt-node1.example.com"
	mcNodeStatus = "active"
	mcRevision   = int64(3)
)

func clusterConfigAPIResponse() nsxModel.ClusterConfig {
	entityIP := mcNodeIP1
	entityPort := int64(443)
	nodeUUID := mcNodeUUID1
	nodeFQDN := mcNodeFQDN1
	nodeStatus := mcNodeStatus
	clusterID := mcClusterID

	return nsxModel.ClusterConfig{
		ClusterId: &clusterID,
		Revision:  &mcRevision,
		Nodes: []nsxModel.ClusterNodeInfo{
			{
				NodeUuid: &nodeUUID,
				Fqdn:     &nodeFQDN,
				Status:   &nodeStatus,
				Entities: []nsxModel.NodeEntityInfo{
					{
						IpAddress: &entityIP,
						Port:      &entityPort,
					},
				},
			},
		},
	}
}

func setupClusterMock(t *testing.T, ctrl *gomock.Controller) (*nsxmocks.MockClusterClient, func()) {
	mockSDK := nsxmocks.NewMockClusterClient(ctrl)

	originalCli := cliClusterClient
	cliClusterClient = func(_ client.Connector) nsx.ClusterClient {
		return mockSDK
	}
	return mockSDK, func() { cliClusterClient = originalCli }
}

func clusterNodeSchemaData() map[string]interface{} {
	return map[string]interface{}{
		"node": []interface{}{
			map[string]interface{}{
				"id":         mcNodeUUID1,
				"ip_address": mcNodeIP1,
				"username":   "admin",
				"password":   "admin-pass", //nolint:gosec
				"fqdn":       mcNodeFQDN1,
				"status":     mcNodeStatus,
			},
		},
	}
}

func TestMockResourceNsxtManagerClusterRead(t *testing.T) {
	t.Run("Read success - populates node computed fields", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(clusterConfigAPIResponse(), nil)

		res := resourceNsxtManagerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, clusterNodeSchemaData())
		d.SetId(mcClusterID)

		err := resourceNsxtManagerClusterRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, int(mcRevision), d.Get("revision"))
	})

	t.Run("Read fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(nsxModel.ClusterConfig{}, errors.New("read API error"))

		res := resourceNsxtManagerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(mcClusterID)

		err := resourceNsxtManagerClusterRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read API error")
	})
}

func TestMockResourceNsxtManagerClusterDelete(t *testing.T) {
	t.Run("Delete success - removes all cluster nodes", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupClusterMock(t, ctrl)
		defer restore()

		force := "true"
		graceful := "true"
		ignoreRepo := "false"
		mockSDK.EXPECT().Removenode(mcNodeUUID1, &force, &graceful, &ignoreRepo).Return(nsxModel.ClusterConfiguration{}, nil)

		res := resourceNsxtManagerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, clusterNodeSchemaData())
		d.SetId(mcClusterID)

		err := resourceNsxtManagerClusterDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when Removenode returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Removenode(mcNodeUUID1, gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.ClusterConfiguration{}, errors.New("remove API error"))

		res := resourceNsxtManagerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, clusterNodeSchemaData())
		d.SetId(mcClusterID)

		err := resourceNsxtManagerClusterDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "remove API error")
	})
}
