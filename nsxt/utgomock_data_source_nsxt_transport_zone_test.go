//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/go-vmware-nsxt/manager"
)

// dataSourceNsxtTransportZoneRead uses the legacy go-vmware-nsxt swagger
// client (nsxtClients.NsxtClient), which issues plain HTTP+JSON requests via
// a concrete (non-interface) *api.APIClient. That type can't be swapped for
// a gomock double the way the newer vAPI SDK clients are, so this test
// instead spins up a real httptest.Server and points a real NsxtClient at
// it (mirroring how TestUnitNsxt_configureNsxtClient in
// utgomock_provider_test.go builds a client without touching the network).
func setupTransportZoneHTTPTestClient(t *testing.T, handler http.HandlerFunc) (nsxtClients, func()) {
	t.Helper()
	server := httptest.NewServer(handler)

	res := Provider()
	pd := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"session_auth": false}))

	clients := &nsxtClients{}
	err := configureNsxtClient(pd, clients)
	require.NoError(t, err)
	require.NotNil(t, clients.NsxtClient)

	// Redirect the already-constructed client at our test server.
	clients.NsxtClientConfig.Host = server.Listener.Addr().String()
	clients.NsxtClientConfig.Scheme = "http"

	return *clients, server.Close
}

func transportZoneJSONHandler(t *testing.T, byID map[string]manager.TransportZone, list manager.TransportZoneListResult, listStatus int) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v1/transport-zones" {
			if listStatus != 0 && listStatus != http.StatusOK {
				w.WriteHeader(listStatus)
				return
			}
			require.NoError(t, json.NewEncoder(w).Encode(list))
			return
		}
		// /api/v1/transport-zones/{id}
		id := r.URL.Path[len("/api/v1/transport-zones/"):]
		tz, ok := byID[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		require.NoError(t, json.NewEncoder(w).Encode(tz))
	}
}

func TestUnitNsxt_DataSourceNsxtTransportZoneRead(t *testing.T) {
	t.Run("missing id and display_name errors without a network call", func(t *testing.T) {
		ds := dataSourceNsxtTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		m, closeServer := setupTransportZoneHTTPTestClient(t, transportZoneJSONHandler(t, nil, manager.TransportZoneListResult{}, http.StatusOK))
		defer closeServer()

		err := dataSourceNsxtTransportZoneRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining transport zone ID or name")
	})

	t.Run("by id success", func(t *testing.T) {
		tz := manager.TransportZone{Id: "tz-1", DisplayName: "tz-one", TransportType: "OVERLAY", HostSwitchName: "hs-1"}
		m, closeServer := setupTransportZoneHTTPTestClient(t, transportZoneJSONHandler(t, map[string]manager.TransportZone{"tz-1": tz}, manager.TransportZoneListResult{}, http.StatusOK))
		defer closeServer()

		ds := dataSourceNsxtTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "tz-1"})

		err := dataSourceNsxtTransportZoneRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, "tz-1", d.Id())
		assert.Equal(t, "tz-one", d.Get("display_name"))
		assert.Equal(t, "OVERLAY", d.Get("transport_type"))
	})

	t.Run("by id not found", func(t *testing.T) {
		m, closeServer := setupTransportZoneHTTPTestClient(t, transportZoneJSONHandler(t, map[string]manager.TransportZone{}, manager.TransportZoneListResult{}, http.StatusOK))
		defer closeServer()

		ds := dataSourceNsxtTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "no-such-tz"})

		err := dataSourceNsxtTransportZoneRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by display_name perfect match", func(t *testing.T) {
		list := manager.TransportZoneListResult{Results: []manager.TransportZone{
			{Id: "tz-1", DisplayName: "tz-one", TransportType: "OVERLAY"},
		}}
		m, closeServer := setupTransportZoneHTTPTestClient(t, transportZoneJSONHandler(t, nil, list, http.StatusOK))
		defer closeServer()

		ds := dataSourceNsxtTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "tz-one"})

		err := dataSourceNsxtTransportZoneRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, "tz-1", d.Id())
	})

	t.Run("by display_name prefix match", func(t *testing.T) {
		list := manager.TransportZoneListResult{Results: []manager.TransportZone{
			{Id: "tz-1", DisplayName: "tz-one-prod", TransportType: "OVERLAY"},
		}}
		m, closeServer := setupTransportZoneHTTPTestClient(t, transportZoneJSONHandler(t, nil, list, http.StatusOK))
		defer closeServer()

		ds := dataSourceNsxtTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "tz-one"})

		err := dataSourceNsxtTransportZoneRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, "tz-1", d.Id())
	})

	t.Run("by display_name multiple perfect matches errors", func(t *testing.T) {
		list := manager.TransportZoneListResult{Results: []manager.TransportZone{
			{Id: "tz-1", DisplayName: "dup", TransportType: "OVERLAY"},
			{Id: "tz-2", DisplayName: "dup", TransportType: "VLAN"},
		}}
		m, closeServer := setupTransportZoneHTTPTestClient(t, transportZoneJSONHandler(t, nil, list, http.StatusOK))
		defer closeServer()

		ds := dataSourceNsxtTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "dup"})

		err := dataSourceNsxtTransportZoneRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Found multiple transport zones")
	})

	t.Run("by display_name not found", func(t *testing.T) {
		m, closeServer := setupTransportZoneHTTPTestClient(t, transportZoneJSONHandler(t, nil, manager.TransportZoneListResult{}, http.StatusOK))
		defer closeServer()

		ds := dataSourceNsxtTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "nonexistent"})

		err := dataSourceNsxtTransportZoneRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("List API error is wrapped", func(t *testing.T) {
		m, closeServer := setupTransportZoneHTTPTestClient(t, transportZoneJSONHandler(t, nil, manager.TransportZoneListResult{}, http.StatusInternalServerError))
		defer closeServer()

		ds := dataSourceNsxtTransportZone()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"display_name": "tz-one"})

		err := dataSourceNsxtTransportZoneRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading transport zones")
	})
}
