//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// data_source_nsxt_policy_baremetal_server_interface_group_associations.go
// builds its BareMetalServerInterfaceGroupAssociationsClient via a direct
// constructor call (api/infra.NewBareMetalServerInterfaceGroupAssociationsClient),
// not a swappable package-level client var, so there is no seam to inject a
// gomock mock. This test instead exercises the client at the HTTP layer via
// newVapiRestTestServer (see utgomock_vapi_rest_helpers_test.go).

package nsxt

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyBareMetalServerInterfaceGroupAssociationsRead(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	externalID := "bmsi-ext-1"

	t.Run("success with results", func(t *testing.T) {
		m, closeServer := newVapiRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/policy/api/v1/infra/bmsi-group-associations", r.URL.Path)
			list := nsxModel.PolicyResourceReferenceForEPListResult{
				Results: []nsxModel.PolicyResourceReferenceForEP{
					bmsGroupAssociationAPIResponse("group-1", "group-one", "Group", true),
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(vapiWireJSON(t, list, nsxModel.PolicyResourceReferenceForEPListResultBindingType())))
		})
		defer closeServer()

		ds := dataSourceNsxtPolicyBareMetalServerInterfaceGroupAssociations()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"external_id": externalID,
		})

		err := dataSourceNsxtPolicyBareMetalServerInterfaceGroupAssociationsRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, externalID, d.Id())
		groups := d.Get("groups").([]interface{})
		require.Len(t, groups, 1)
		group := groups[0].(map[string]interface{})
		assert.Equal(t, "group-one", group["display_name"])
		assert.Equal(t, true, group["is_valid"])
	})

	t.Run("empty results", func(t *testing.T) {
		m, closeServer := newVapiRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(vapiWireJSON(t, nsxModel.PolicyResourceReferenceForEPListResult{}, nsxModel.PolicyResourceReferenceForEPListResultBindingType())))
		})
		defer closeServer()

		ds := dataSourceNsxtPolicyBareMetalServerInterfaceGroupAssociations()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"external_id": externalID,
		})

		err := dataSourceNsxtPolicyBareMetalServerInterfaceGroupAssociationsRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, externalID, d.Id())
	})

	t.Run("API error", func(t *testing.T) {
		m, closeServer := newVapiRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		defer closeServer()

		ds := dataSourceNsxtPolicyBareMetalServerInterfaceGroupAssociations()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"external_id": externalID,
		})

		err := dataSourceNsxtPolicyBareMetalServerInterfaceGroupAssociationsRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to read group associations")
	})

	t.Run("missing external_id errors before any API call", func(t *testing.T) {
		ds := dataSourceNsxtPolicyBareMetalServerInterfaceGroupAssociations()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyBareMetalServerInterfaceGroupAssociationsRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "external_id is required")
	})

	t.Run("with enforcement_point_path", func(t *testing.T) {
		m, closeServer := newVapiRestTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.URL.RawQuery, "enforcement_point_path")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(vapiWireJSON(t, nsxModel.PolicyResourceReferenceForEPListResult{}, nsxModel.PolicyResourceReferenceForEPListResultBindingType())))
		})
		defer closeServer()

		ds := dataSourceNsxtPolicyBareMetalServerInterfaceGroupAssociations()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"external_id":            externalID,
			"enforcement_point_path": "/infra/sites/default/enforcement-points/default",
		})

		err := dataSourceNsxtPolicyBareMetalServerInterfaceGroupAssociationsRead(d, m)
		require.NoError(t, err)
	})
}
