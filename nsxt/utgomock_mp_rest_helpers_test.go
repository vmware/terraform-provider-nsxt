//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// Shared helper for data sources that talk to the legacy go-vmware-nsxt
// manager REST client (m.(nsxtClients).NsxtClient), which is a concrete
// *api.APIClient with no interface seam for gomock. Those clients just
// json.Decode the HTTP response body directly into the target struct, so a
// real (local, offline) httptest.Server serving plain encoding/json bodies
// is enough to exercise them.

package nsxt

import (
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/vmware/go-vmware-nsxt"
)

// newMPRestTestServer starts an httptest server and returns an
// api.APIClient wired to it (session auth disabled, so no login POST is
// made), ready to assign to nsxtClients.NsxtClient.
func newMPRestTestServer(t *testing.T, handler http.HandlerFunc) (*api.APIClient, func()) {
	t.Helper()
	server := httptest.NewServer(handler)
	cfg := &api.Configuration{
		BasePath:        server.URL + "/api/v1",
		HTTPClient:      server.Client(),
		SkipSessionAuth: true,
		RemoteAuth:      true,
	}
	client, err := api.NewAPIClient(cfg)
	if err != nil {
		t.Fatalf("failed to build MP API test client: %v", err)
	}
	return client, server.Close
}
