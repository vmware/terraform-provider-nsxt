//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// Reuses the ClusterClient mock and setupClusterMock helper already defined
// for resource_nsxt_manager_cluster.go in utgomock_resource_nsxt_manager_cluster_test.go:
// mockgen -destination=mocks/nsx/ClusterClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/ClusterClient.go ClusterClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"go.uber.org/mock/gomock"
)

func managementClusterAPIResponse(clusterID, ip, thumbprint string) nsxModel.ClusterConfig {
	return nsxModel.ClusterConfig{
		ClusterId: &clusterID,
		Nodes: []nsxModel.ClusterNodeInfo{
			{
				ApiListenAddr: &nsxModel.ServiceEndpoint{
					IpAddress:                   &ip,
					CertificateSha256Thumbprint: &thumbprint,
				},
			},
		},
	}
}

func TestMockDataSourceNsxtManagementClusterRead(t *testing.T) {
	clusterID := "mgmt-cluster-1"
	nodeIP := "192.0.2.1"
	thumbprint := "AA:BB:CC"

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(managementClusterAPIResponse(clusterID, nodeIP, thumbprint), nil)

		ds := dataSourceNsxtManagementCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		m.Host = "https://" + nodeIP

		err := dataSourceNsxtManagementClusterRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, clusterID, d.Id())
		assert.Equal(t, thumbprint, d.Get("node_sha256_thumbprint"))
	})

	t.Run("Get API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(nsxModel.ClusterConfig{}, errors.New("get failed"))

		ds := dataSourceNsxtManagementCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		m.Host = "https://" + nodeIP

		err := dataSourceNsxtManagementClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error while reading cluster configuration")
	})

	t.Run("no matching node thumbprint errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupClusterMock(t, ctrl)
		defer restore()

		otherIP := "192.0.2.99"
		mockSDK.EXPECT().Get().Return(managementClusterAPIResponse(clusterID, otherIP, thumbprint), nil)

		ds := dataSourceNsxtManagementCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		m.Host = "https://" + nodeIP

		err := dataSourceNsxtManagementClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "thumbprint not found")
	})
}
