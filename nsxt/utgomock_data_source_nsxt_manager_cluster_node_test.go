//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// data_source_nsxt_manager_cluster_node.go builds its cluster.NodesClient via
// a direct `var cliClusterNodesClient = cluster.NewNodesClient` alias, whose
// inferred type pins the return value to the SDK's unexported concrete
// *nodesClient type. That leaves no swappable seam for a gomock mock, so
// this test instead exercises the client at the HTTP layer via
// newVapiRestTestServer (see utgomock_vapi_rest_helpers_test.go).

package nsxt

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

func clusterNodeConfigAPIResponse(id, name, addr string) nsxModel.ClusterNodeConfig {
	return nsxModel.ClusterNodeConfig{
		Id:                      &id,
		DisplayName:             &name,
		Description:             &name,
		ApplianceMgmtListenAddr: &addr,
	}
}

func TestMockDataSourceNsxtManagerClusterNodeRead(t *testing.T) {
	nodeID := "node-1"
	nodeName := "node-name"
	addr := "10.0.0.1:443"

	t.Run("by id success", func(t *testing.T) {
		m, closeServer := newVapiRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/cluster/nodes/"+nodeID, r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(vapiWireJSON(t, clusterNodeConfigAPIResponse(nodeID, nodeName, addr), nsxModel.ClusterNodeConfigBindingType())))
		})
		defer closeServer()

		ds := dataSourceNsxtManagerClusterNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": nodeID,
		})

		err := dataSourceNsxtManagerClusterNodeRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, nodeID, d.Id())
		assert.Equal(t, nodeName, d.Get("display_name"))
		assert.Equal(t, addr, d.Get("appliance_mgmt_listen_address"))
	})

	t.Run("by id API error", func(t *testing.T) {
		m, closeServer := newVapiRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		defer closeServer()

		ds := dataSourceNsxtManagerClusterNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": nodeID,
		})

		err := dataSourceNsxtManagerClusterNodeRead(d, m)
		require.Error(t, err)
	})

	t.Run("missing id and display_name errors", func(t *testing.T) {
		ds := dataSourceNsxtManagerClusterNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtManagerClusterNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining ClusterNode")
	})

	t.Run("by display_name single exact match", func(t *testing.T) {
		m, closeServer := newVapiRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/cluster/nodes", r.URL.Path)
			list := nsxModel.ClusterNodeConfigListResult{
				Results: []nsxModel.ClusterNodeConfig{clusterNodeConfigAPIResponse(nodeID, nodeName, addr)},
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(vapiWireJSON(t, list, nsxModel.ClusterNodeConfigListResultBindingType())))
		})
		defer closeServer()

		ds := dataSourceNsxtManagerClusterNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": nodeName,
		})

		err := dataSourceNsxtManagerClusterNodeRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, nodeID, d.Id())
	})

	t.Run("by display_name prefix single match", func(t *testing.T) {
		otherName := "node-name-other"
		m, closeServer := newVapiRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			list := nsxModel.ClusterNodeConfigListResult{
				Results: []nsxModel.ClusterNodeConfig{clusterNodeConfigAPIResponse(nodeID, otherName, addr)},
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(vapiWireJSON(t, list, nsxModel.ClusterNodeConfigListResultBindingType())))
		})
		defer closeServer()

		ds := dataSourceNsxtManagerClusterNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "node-name-oth",
		})

		err := dataSourceNsxtManagerClusterNodeRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, nodeID, d.Id())
	})

	t.Run("by display_name multiple matches errors", func(t *testing.T) {
		m, closeServer := newVapiRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			list := nsxModel.ClusterNodeConfigListResult{
				Results: []nsxModel.ClusterNodeConfig{
					clusterNodeConfigAPIResponse("id-1", nodeName, addr),
					clusterNodeConfigAPIResponse("id-2", nodeName, addr),
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(vapiWireJSON(t, list, nsxModel.ClusterNodeConfigListResultBindingType())))
		})
		defer closeServer()

		ds := dataSourceNsxtManagerClusterNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": nodeName,
		})

		err := dataSourceNsxtManagerClusterNodeRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "found multiple ClusterNode")
	})

	t.Run("no match from list", func(t *testing.T) {
		m, closeServer := newVapiRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			list := nsxModel.ClusterNodeConfigListResult{Results: []nsxModel.ClusterNodeConfig{}}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(vapiWireJSON(t, list, nsxModel.ClusterNodeConfigListResultBindingType())))
		})
		defer closeServer()

		ds := dataSourceNsxtManagerClusterNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		err := dataSourceNsxtManagerClusterNodeRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})
}
