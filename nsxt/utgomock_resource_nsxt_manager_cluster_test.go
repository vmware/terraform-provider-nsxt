//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/nsx/ClusterClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/ClusterClient.go ClusterClient
//
// Note: Create/Update are unit-testable as long as node API probing is disabled
// (api_probing.enabled = false), since otherwise they wait on live connectivity
// polling. With probing disabled, and the client's Host set to a literal IP so
// resolveHostIPs takes its non-DNS branch, the rest of the join flow (security
// context construction, cluster client calls) is pure local logic.

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"go.uber.org/mock/gomock"

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
	cliClusterClient = func(_ client.Connector) clusterOps {
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

func TestUnitNsxt_getClusterNodesFromSchema(t *testing.T) {
	res := resourceNsxtManagerCluster()
	d := schema.TestResourceDataRaw(t, res.Schema, clusterNodeSchemaData())

	nodes := getClusterNodesFromSchema(d)
	require.Len(t, nodes, 1)
	assert.Equal(t, mcNodeIP1, nodes[0].IPAddress)
	assert.Equal(t, "admin", nodes[0].UserName)
	assert.Equal(t, "admin-pass", nodes[0].Password)
}

func TestUnitNsxt_getClusterNodesIPs(t *testing.T) {
	nodes := []interface{}{
		map[string]interface{}{"ip_address": "10.0.0.1"},
		map[string]interface{}{"ip_address": "10.0.0.2"},
	}
	assert.Equal(t, []string{"10.0.0.1", "10.0.0.2"}, getClusterNodesIPs(nodes))
}

func TestUnitNsxt_getMatchingIPVersion(t *testing.T) {
	assert.Equal(t, "10.0.0.5", getMatchingIPVersion("10.0.0.1", []string{"10.0.0.5", "::1"}))
	assert.Equal(t, "::1", getMatchingIPVersion("::2", []string{"10.0.0.5", "::1"}))
	assert.Equal(t, "", getMatchingIPVersion("10.0.0.1", []string{"::1"}))
}

func TestUnitNsxt_isMatchingNode(t *testing.T) {
	ip := mcNodeIP1
	port := int64(443)
	node := nsxModel.ClusterNodeInfo{Entities: []nsxModel.NodeEntityInfo{{IpAddress: &ip, Port: &port}}}

	assert.True(t, isMatchingNode(node, mcNodeIP1))
	assert.False(t, isMatchingNode(node, "10.10.10.10"))
	assert.False(t, isMatchingNode(nsxModel.ClusterNodeInfo{}, mcNodeIP1))
}

func TestUnitNsxt_resolveHostIPs(t *testing.T) {
	c := newGoMockProviderClient()
	c.Host = "https://203.0.113.9"

	ips, err := resolveHostIPs(c)
	require.NoError(t, err)
	assert.Equal(t, []string{"203.0.113.9"}, ips)
}

func TestUnitNsxt_getHostCredential(t *testing.T) {
	username, password := getHostCredential(newGoMockProviderClient())
	assert.Equal(t, "username", username)
	assert.Equal(t, "password", password)
}

// clusterConfigWithListenAddr builds a ClusterConfig whose single node has both
// an ApiListenAddr (needed by getClusterInfoFromHostNode) and a matching
// Entities[0].IpAddress (needed by isMatchingNode during Read).
func clusterConfigWithListenAddr(nodeIP, thumbprint string) nsxModel.ClusterConfig {
	clusterID := mcClusterID
	nodeUUID := mcNodeUUID1
	entityIP := nodeIP
	entityPort := int64(443)
	thumb := thumbprint
	return nsxModel.ClusterConfig{
		ClusterId: &clusterID,
		Revision:  &mcRevision,
		Nodes: []nsxModel.ClusterNodeInfo{
			{
				NodeUuid:      &nodeUUID,
				ApiListenAddr: &nsxModel.ServiceEndpoint{CertificateSha256Thumbprint: &thumb},
				Entities:      []nsxModel.NodeEntityInfo{{IpAddress: &entityIP, Port: &entityPort}},
			},
		},
	}
}

func TestMockResourceNsxtManagerClusterCreate(t *testing.T) {
	t.Run("Create fails when no nodes provided", func(t *testing.T) {
		res := resourceNsxtManagerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtManagerClusterCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be provided to form a cluster")
	})

	t.Run("Create success with probing disabled joins node and reads back state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupClusterMock(t, ctrl)
		defer restore()

		guestIP := "192.168.50.20"
		hostIP := "203.0.113.9"
		thumbprint := "AA:BB:CC:DD"

		mockSDK.EXPECT().Get().Return(clusterConfigWithListenAddr(guestIP, thumbprint), nil).Times(2)
		mockSDK.EXPECT().Joincluster(gomock.Any()).Return(nsxModel.ClusterConfiguration{}, nil)

		res := resourceNsxtManagerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"api_probing": []interface{}{map[string]interface{}{"enabled": false}},
			"node": []interface{}{
				map[string]interface{}{
					"ip_address": guestIP,
					"username":   "admin",
					"password":   "admin-pass",
				},
			},
		})

		m := newGoMockProviderClient()
		m.Host = "https://" + hostIP

		err := resourceNsxtManagerClusterCreate(d, m)
		require.NoError(t, err)
		assert.Equal(t, mcClusterID, d.Id())
	})

	t.Run("Create fails when join fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupClusterMock(t, ctrl)
		defer restore()

		guestIP := "192.168.50.21"
		hostIP := "203.0.113.10"
		thumbprint := "EE:FF:00:11"

		mockSDK.EXPECT().Get().Return(clusterConfigWithListenAddr(guestIP, thumbprint), nil)
		mockSDK.EXPECT().Joincluster(gomock.Any()).Return(nsxModel.ClusterConfiguration{}, errors.New("join API error"))

		res := resourceNsxtManagerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"api_probing": []interface{}{map[string]interface{}{"enabled": false}},
			"node": []interface{}{
				map[string]interface{}{
					"ip_address": guestIP,
					"username":   "admin",
					"password":   "admin-pass",
				},
			},
		})

		m := newGoMockProviderClient()
		m.Host = "https://" + hostIP

		err := resourceNsxtManagerClusterCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "join API error")
	})
}

func TestMockResourceNsxtManagerClusterUpdate(t *testing.T) {
	// Note: schema.TestResourceDataRaw has no prior state to diff against, so
	// d.HasChange("node") is always true here; resourceNsxtManagerClusterUpdate
	// therefore always proceeds to re-fetch cluster info via cliClusterClient.
	t.Run("Update fails when cluster info cannot be read", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get().Return(nsxModel.ClusterConfig{}, errors.New("get API error"))

		res := resourceNsxtManagerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, clusterNodeSchemaData())
		d.SetId(mcClusterID)

		err := resourceNsxtManagerClusterUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "get API error")
	})

	t.Run("Update joins newly added node", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupClusterMock(t, ctrl)
		defer restore()

		guestIP := "192.168.60.30"
		hostIP := "203.0.113.11"
		thumbprint := "12:34:56:78"

		mockSDK.EXPECT().Get().Return(clusterConfigWithListenAddr(guestIP, thumbprint), nil).Times(2)
		mockSDK.EXPECT().Joincluster(gomock.Any()).Return(nsxModel.ClusterConfiguration{}, nil)

		res := resourceNsxtManagerCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"node": []interface{}{
				map[string]interface{}{
					"ip_address": guestIP,
					"username":   "admin",
					"password":   "admin-pass",
				},
			},
		})
		d.SetId(mcClusterID)

		m := newGoMockProviderClient()
		m.Host = "https://" + hostIP

		err := resourceNsxtManagerClusterUpdate(d, m)
		require.NoError(t, err)
	})
}
