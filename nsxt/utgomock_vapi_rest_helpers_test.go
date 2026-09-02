//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// Shared helpers for data sources that build their vAPI client directly
// (e.g. via a package-level constructor call rather than a swappable
// package-level client var), so there is no seam to inject a gomock mock.
// For those, we stand up a real httptest.Server and point the provider's
// Host/PolicyHTTPClient at it: the generated vAPI REST client will make a
// real (but fully local, offline) HTTP call, and we control the JSON it
// gets back by encoding the expected Go model through the SDK's own
// bindings type converter + cleanjson encoder, so we don't have to guess
// the wire field names by hand.

package nsxt

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data/serializers/cleanjson"
)

// vapiWireJSON encodes a Go vAPI model struct into the exact JSON wire
// format the generated REST clients expect on the response body.
func vapiWireJSON(t *testing.T, goObj interface{}, bindingType bindings.BindingType) string {
	t.Helper()
	converter := bindings.NewTypeConverter()
	dataVal, errs := converter.ConvertToVapi(goObj, bindingType)
	require.Empty(t, errs)
	encoder := cleanjson.NewDataValueToJsonEncoder()
	jsonStr, err := encoder.Encode(dataVal)
	require.NoError(t, err)
	return jsonStr
}

// newVapiRestTestServer starts an httptest server and returns a
// nsxtClients-ready base client whose Host/PolicyHTTPClient point at it.
// The caller supplies a handler that inspects method+path and writes the
// desired JSON response (see vapiWireJSON) and status code.
func newVapiRestTestServer(t *testing.T, handler http.HandlerFunc) (nsxtClients, func()) {
	t.Helper()
	server := httptest.NewServer(handler)
	m := newGoMockProviderClient()
	m.Host = server.URL
	m.PolicyHTTPClient = server.Client()
	return m, server.Close
}
