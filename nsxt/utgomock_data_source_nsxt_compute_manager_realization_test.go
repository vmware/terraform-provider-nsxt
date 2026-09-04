//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// data_source_nsxt_compute_manager_realization.go builds its StateClient and
// StatusClient via direct `var cliComputeManagerStateClient = compute_managers.NewStateClient`
// / `NewStatusClient` aliases, whose inferred types pin the return value to
// the SDK's unexported concrete types. That leaves no swappable seam for a
// gomock mock, so this test instead exercises the clients at the HTTP layer
// via newVapiRestTestServer (see utgomock_vapi_rest_helpers_test.go).

package nsxt

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

func TestMockDataSourceNsxtComputeManagerRealizationRead(t *testing.T) {
	cmID := "cm-1"

	t.Run("realization and registration succeed", func(t *testing.T) {
		m, closeServer := newVapiRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Path {
			case "/api/v1/fabric/compute-managers/" + cmID + "/state":
				state := nsxModel.ConfigurationState_STATE_SUCCESS
				_, _ = w.Write([]byte(vapiWireJSON(t, nsxModel.ConfigurationState{State: &state}, nsxModel.ConfigurationStateBindingType())))
			case "/api/v1/fabric/compute-managers/" + cmID + "/status":
				status := nsxModel.ComputeManagerStatus_REGISTRATION_STATUS_REGISTERED
				_, _ = w.Write([]byte(vapiWireJSON(t, nsxModel.ComputeManagerStatus{RegistrationStatus: &status}, nsxModel.ComputeManagerStatusBindingType())))
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		})
		defer closeServer()

		ds := dataSourceNsxtComputeManagerRealization()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":      cmID,
			"timeout": 5,
			"delay":   0,
		})

		err := dataSourceNsxtComputeManagerRealizationRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, nsxModel.ConfigurationState_STATE_SUCCESS, d.Get("state"))
		assert.Equal(t, nsxModel.ComputeManagerStatus_REGISTRATION_STATUS_REGISTERED, d.Get("registration_status"))
	})

	t.Run("registration check skipped when check_registration is false", func(t *testing.T) {
		statusCalled := false
		m, closeServer := newVapiRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Path {
			case "/api/v1/fabric/compute-managers/" + cmID + "/state":
				state := nsxModel.ConfigurationState_STATE_SUCCESS
				_, _ = w.Write([]byte(vapiWireJSON(t, nsxModel.ConfigurationState{State: &state}, nsxModel.ConfigurationStateBindingType())))
			case "/api/v1/fabric/compute-managers/" + cmID + "/status":
				statusCalled = true
				w.WriteHeader(http.StatusInternalServerError)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		})
		defer closeServer()

		ds := dataSourceNsxtComputeManagerRealization()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":                 cmID,
			"timeout":            5,
			"delay":              0,
			"check_registration": false,
		})

		err := dataSourceNsxtComputeManagerRealizationRead(d, m)
		require.NoError(t, err)
		assert.False(t, statusCalled)
	})

	t.Run("realization failure surfaces error when check_registration is false", func(t *testing.T) {
		// The data source only returns the realization error directly when
		// check_registration is false; otherwise its return value is
		// entirely replaced by the (possibly successful) registration
		// check result.
		m, closeServer := newVapiRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			state := nsxModel.ConfigurationState_STATE_FAILED
			_, _ = w.Write([]byte(vapiWireJSON(t, nsxModel.ConfigurationState{State: &state}, nsxModel.ConfigurationStateBindingType())))
		})
		defer closeServer()

		ds := dataSourceNsxtComputeManagerRealization()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":                 cmID,
			"timeout":            5,
			"delay":              0,
			"check_registration": false,
		})

		err := dataSourceNsxtComputeManagerRealizationRead(d, m)
		require.Error(t, err)
	})

	t.Run("state API error is wrapped", func(t *testing.T) {
		m, closeServer := newVapiRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		defer closeServer()

		ds := dataSourceNsxtComputeManagerRealization()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":      cmID,
			"timeout": 5,
			"delay":   0,
		})

		err := dataSourceNsxtComputeManagerRealizationRead(d, m)
		require.Error(t, err)
	})
}
