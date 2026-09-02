//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// dataSourceNsxtPolicyEdgeTransportNodeRead is a "thin wrapper" data source
// that resolves via the search API (policyDataSourceResourceRead), not via a
// directly-mockable SDK client. It is tested here using the seqQueryListClient
// stub and setupCliQueryClientStub helper defined in policy_search_unit_test.go.

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func policyEdgeTransportNodeToStructValue(t *testing.T, etn nsxModel.PolicyEdgeTransportNode) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(etn, nsxModel.PolicyEdgeTransportNodeBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func TestUnitNsxt_dataSourceNsxtPolicyEdgeTransportNodeRead(t *testing.T) {
	rt := "PolicyEdgeTransportNode"
	ds := dataSourceNsxtPolicyEdgeTransportNode()

	t.Run("success by id", func(t *testing.T) {
		uniqueID := "unique-1"
		sv := policyEdgeTransportNodeToStructValue(t, nsxModel.PolicyEdgeTransportNode{
			Id:           str("etn-1"),
			DisplayName:  str("etn-name"),
			Path:         str("/infra/sites/default/enforcement-points/default/edge-transport-nodes/etn-1"),
			ResourceType: &rt,
			UniqueId:     &uniqueID,
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": "etn-1",
		})
		err := dataSourceNsxtPolicyEdgeTransportNodeRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "etn-1", d.Id())
		assert.Equal(t, "etn-name", d.Get("display_name"))
		assert.Equal(t, uniqueID, d.Get("unique_id"))
	})

	t.Run("not found", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{}, ResultCount: i64(0),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": "nonexistent",
		})
		err := dataSourceNsxtPolicyEdgeTransportNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search-fail")}}
		defer setupCliQueryClientStub(t, stub)()

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": "etn-1",
		})
		err := dataSourceNsxtPolicyEdgeTransportNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("no id or display_name", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
		err := dataSourceNsxtPolicyEdgeTransportNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "No 'id' or 'display_name'")
	})
}
