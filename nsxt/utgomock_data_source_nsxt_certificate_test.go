//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// data_source_nsxt_certificate.go talks to the legacy go-vmware-nsxt
// NsxComponentAdministrationApi, a concrete *api.APIClient with no interface
// seam for gomock. See utgomock_mp_rest_helpers_test.go for the httptest
// based approach used instead.

package nsxt

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/go-vmware-nsxt/trust"
)

func writeJSON(t *testing.T, w http.ResponseWriter, v interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	require.NoError(t, json.NewEncoder(w).Encode(v))
}

func TestMockDataSourceNsxtCertificateRead(t *testing.T) {
	certID := "cert-1"
	certName := "my-cert"

	t.Run("by id success", func(t *testing.T) {
		nsxClient, closeServer := newMPRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/trust-management/certificates/"+certID, r.URL.Path)
			writeJSON(t, w, trust.Certificate{Id: certID, DisplayName: certName, Description: "desc"})
		})
		defer closeServer()

		ds := dataSourceNsxtCertificate()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": certID,
		})

		m := newGoMockProviderClient()
		m.NsxtClient = nsxClient

		err := dataSourceNsxtCertificateRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, certID, d.Id())
		assert.Equal(t, certName, d.Get("display_name"))
	})

	t.Run("by id not found", func(t *testing.T) {
		nsxClient, closeServer := newMPRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		defer closeServer()

		ds := dataSourceNsxtCertificate()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": certID,
		})

		m := newGoMockProviderClient()
		m.NsxtClient = nsxClient

		err := dataSourceNsxtCertificateRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("missing id and display_name errors before any API call", func(t *testing.T) {
		nsxClient, closeServer := newMPRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			t.Fatalf("unexpected API call to %s", r.URL.Path)
		})
		defer closeServer()

		ds := dataSourceNsxtCertificate()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		m.NsxtClient = nsxClient

		err := dataSourceNsxtCertificateRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining certificate ID or name")
	})

	t.Run("by display_name single match", func(t *testing.T) {
		nsxClient, closeServer := newMPRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/trust-management/certificates", r.URL.Path)
			writeJSON(t, w, trust.CertificateList{
				Results:     []trust.Certificate{{Id: certID, DisplayName: certName}},
				ResultCount: 1,
			})
		})
		defer closeServer()

		ds := dataSourceNsxtCertificate()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": certName,
		})

		m := newGoMockProviderClient()
		m.NsxtClient = nsxClient

		err := dataSourceNsxtCertificateRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, certID, d.Id())
	})

	t.Run("by display_name multiple matches errors", func(t *testing.T) {
		nsxClient, closeServer := newMPRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, trust.CertificateList{
				Results: []trust.Certificate{
					{Id: "id-1", DisplayName: certName},
					{Id: "id-2", DisplayName: certName},
				},
				ResultCount: 2,
			})
		})
		defer closeServer()

		ds := dataSourceNsxtCertificate()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": certName,
		})

		m := newGoMockProviderClient()
		m.NsxtClient = nsxClient

		err := dataSourceNsxtCertificateRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Found multiple certificates")
	})

	t.Run("by display_name no match", func(t *testing.T) {
		nsxClient, closeServer := newMPRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			writeJSON(t, w, trust.CertificateList{Results: []trust.Certificate{}, ResultCount: 0})
		})
		defer closeServer()

		ds := dataSourceNsxtCertificate()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		m := newGoMockProviderClient()
		m.NsxtClient = nsxClient

		err := dataSourceNsxtCertificateRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})
}
