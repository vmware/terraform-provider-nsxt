//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// dataSourceNsxtPolicyHostTransportNodeRead is a "thin wrapper" data source
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

func policyHostTransportNodeToStructValue(t *testing.T, htn nsxModel.HostTransportNode) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(htn, nsxModel.HostTransportNodeBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func TestUnitNsxt_dataSourceNsxtPolicyHostTransportNodeRead(t *testing.T) {
	rt := "HostTransportNode"
	ds := dataSourceNsxtPolicyHostTransportNode()

	t.Run("success by id", func(t *testing.T) {
		uniqueID := "htn-unique-1"
		sv := policyHostTransportNodeToStructValue(t, nsxModel.HostTransportNode{
			Id:           str("htn-1"),
			DisplayName:  str("htn-name"),
			Path:         str("/infra/sites/default/enforcement-points/default/host-transport-nodes/htn-1"),
			ResourceType: &rt,
			UniqueId:     &uniqueID,
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": "htn-1",
		})
		err := dataSourceNsxtPolicyHostTransportNodeRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "htn-1", d.Id())
		assert.Equal(t, "htn-name", d.Get("display_name"))
		assert.Equal(t, uniqueID, d.Get("unique_id"))
	})

	t.Run("success by display_name", func(t *testing.T) {
		uniqueID := "htn-unique-2"
		sv := policyHostTransportNodeToStructValue(t, nsxModel.HostTransportNode{
			Id:           str("htn-2"),
			DisplayName:  str("htn-name-2"),
			Path:         str("/infra/sites/default/enforcement-points/default/host-transport-nodes/htn-2"),
			ResourceType: &rt,
			UniqueId:     &uniqueID,
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "htn-name-2",
		})
		err := dataSourceNsxtPolicyHostTransportNodeRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "htn-2", d.Id())
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
		err := dataSourceNsxtPolicyHostTransportNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search-fail")}}
		defer setupCliQueryClientStub(t, stub)()

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": "htn-1",
		})
		err := dataSourceNsxtPolicyHostTransportNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("no id or display_name", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
		err := dataSourceNsxtPolicyHostTransportNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "No 'id' or 'display_name'")
	})
}
