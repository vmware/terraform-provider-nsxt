//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// data_source_nsxt_edge_cluster.go talks to the legacy go-vmware-nsxt
// NetworkTransportApi, a concrete *api.APIClient with no interface seam for
// gomock. See utgomock_mp_rest_helpers_test.go for the httptest based
// approach used instead.

package nsxt

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/go-vmware-nsxt/manager"
)

func TestMockDataSourceNsxtEdgeClusterRead(t *testing.T) {
	edgeClusterID := "ec-1"
	edgeClusterName := "my-edge-cluster"

	t.Run("by id success", func(t *testing.T) {
		nsxClient, closeServer := newMPRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/edge-clusters/"+edgeClusterID, r.URL.Path)
			writeJSON(t, w, manager.EdgeCluster{
				Id:             edgeClusterID,
				DisplayName:    edgeClusterName,
				Description:    "desc",
				DeploymentType: "VIRTUAL_MACHINE",
				MemberNodeType: "EDGE_NODE",
			})
		})
		defer closeServer()

		ds := dataSourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": edgeClusterID,
		})

		m := newGoMockProviderClient()
		m.NsxtClient = nsxClient

		err := dataSourceNsxtEdgeClusterRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, edgeClusterID, d.Id())
		assert.Equal(t, edgeClusterName, d.Get("display_name"))
		assert.Equal(t, "VIRTUAL_MACHINE", d.Get("deployment_type"))
	})

	t.Run("by id not found", func(t *testing.T) {
		nsxClient, closeServer := newMPRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		defer closeServer()

		ds := dataSourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": edgeClusterID,
		})

		m := newGoMockProviderClient()
		m.NsxtClient = nsxClient

		err := dataSourceNsxtEdgeClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("missing id and display_name errors before any API call", func(t *testing.T) {
		nsxClient, closeServer := newMPRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			t.Fatalf("unexpected API call to %s", r.URL.Path)
		})
		defer closeServer()

		ds := dataSourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		m.NsxtClient = nsxClient

		err := dataSourceNsxtEdgeClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining edge cluster ID or name")
	})

	t.Run("by display_name perfect match preferred over prefix", func(t *testing.T) {
		nsxClient, closeServer := newMPRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/edge-clusters", r.URL.Path)
			writeJSON(t, w, manager.EdgeClusterListResult{
				Results: []manager.EdgeCluster{
					{Id: "prefix-id", DisplayName: edgeClusterName + "-other"},
					{Id: edgeClusterID, DisplayName: edgeClusterName},
				},
			})
		})
		defer closeServer()

		ds := dataSourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": edgeClusterName,
		})

		m := newGoMockProviderClient()
		m.NsxtClient = nsxClient

		err := dataSourceNsxtEdgeClusterRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, edgeClusterID, d.Id())
	})

	t.Run("by display_name prefix single match", func(t *testing.T) {
		nsxClient, closeServer := newMPRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, manager.EdgeClusterListResult{
				Results: []manager.EdgeCluster{
					{Id: edgeClusterID, DisplayName: "edge-cluster-other"},
				},
			})
		})
		defer closeServer()

		ds := dataSourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "edge-cluster-oth",
		})

		m := newGoMockProviderClient()
		m.NsxtClient = nsxClient

		err := dataSourceNsxtEdgeClusterRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, edgeClusterID, d.Id())
	})

	t.Run("by display_name multiple perfect matches errors", func(t *testing.T) {
		nsxClient, closeServer := newMPRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, manager.EdgeClusterListResult{
				Results: []manager.EdgeCluster{
					{Id: "id-1", DisplayName: edgeClusterName},
					{Id: "id-2", DisplayName: edgeClusterName},
				},
			})
		})
		defer closeServer()

		ds := dataSourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": edgeClusterName,
		})

		m := newGoMockProviderClient()
		m.NsxtClient = nsxClient

		err := dataSourceNsxtEdgeClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple edge clusters")
	})

	t.Run("by display_name no match", func(t *testing.T) {
		nsxClient, closeServer := newMPRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, manager.EdgeClusterListResult{Results: []manager.EdgeCluster{}})
		})
		defer closeServer()

		ds := dataSourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		m := newGoMockProviderClient()
		m.NsxtClient = nsxClient

		err := dataSourceNsxtEdgeClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("list API error", func(t *testing.T) {
		nsxClient, closeServer := newMPRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		defer closeServer()

		ds := dataSourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": edgeClusterName,
		})

		m := newGoMockProviderClient()
		m.NsxtClient = nsxClient

		err := dataSourceNsxtEdgeClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading edge clusters")
	})
}
