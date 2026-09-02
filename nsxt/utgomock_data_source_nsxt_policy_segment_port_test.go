//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// dataSourceNsxtPolicySegmentPortRead delegates to listPolicyResources, which is backed by
// the search (cliQueryClient) machinery already covered by policy_search_unit_test.go. This
// file reuses seqQueryListClient / setupCliQueryClientStub from that file instead of
// introducing new mocks, plus a local helper to build SegmentPort-shaped StructValues.

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func segmentPortToStructValue(t *testing.T, sp model.SegmentPort) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(sp, model.SegmentPortBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func TestMockDataSourceNsxtPolicySegmentPortRead(t *testing.T) {
	resourceType := "SegmentPort"

	t.Run("success by vif_id", func(t *testing.T) {
		vifID := "vif-1"
		sv := segmentPortToStructValue(t, model.SegmentPort{
			Id: str("port-1"), DisplayName: str("port-1-name"), Path: str("/infra/segments/seg-1/ports/port-1"),
			ResourceType: &resourceType, Attachment: &model.PortAttachment{Id: &vifID},
		})
		stub := &seqQueryListClient{responses: []model.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"vif_id": vifID,
		})

		err := dataSourceNsxtPolicySegmentPortRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "port-1", d.Id())
		assert.Equal(t, vifID, d.Get("vif_id"))
		assert.Equal(t, "/infra/segments/seg-1", d.Get("segment_path"))
	})

	t.Run("vif_id not found", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []model.SearchResponse{{
			Results: []*data.StructValue{}, ResultCount: i64(0),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"vif_id": "missing-vif",
		})

		err := dataSourceNsxtPolicySegmentPortRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("success by id", func(t *testing.T) {
		sv := segmentPortToStructValue(t, model.SegmentPort{
			Id: str("port-2"), DisplayName: str("port-2-name"), Path: str("/infra/segments/seg-1/ports/port-2"),
			ResourceType: &resourceType,
		})
		stub := &seqQueryListClient{responses: []model.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": "port-2",
		})

		err := dataSourceNsxtPolicySegmentPortRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "port-2", d.Id())
	})

	t.Run("success by display_name exact match preferred over prefix", func(t *testing.T) {
		exact := segmentPortToStructValue(t, model.SegmentPort{
			Id: str("port-exact"), DisplayName: str("web"), Path: str("/infra/segments/seg-1/ports/port-exact"),
			ResourceType: &resourceType,
		})
		prefix := segmentPortToStructValue(t, model.SegmentPort{
			Id: str("port-prefix"), DisplayName: str("web-2"), Path: str("/infra/segments/seg-1/ports/port-prefix"),
			ResourceType: &resourceType,
		})
		stub := &seqQueryListClient{responses: []model.SearchResponse{{
			Results: []*data.StructValue{prefix, exact}, ResultCount: i64(2),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "web",
		})

		err := dataSourceNsxtPolicySegmentPortRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "port-exact", d.Id())
	})

	t.Run("display_name multiple prefix matches error", func(t *testing.T) {
		p1 := segmentPortToStructValue(t, model.SegmentPort{
			Id: str("p1"), DisplayName: str("web-a"), Path: str("/infra/segments/seg-1/ports/p1"),
			ResourceType: &resourceType,
		})
		p2 := segmentPortToStructValue(t, model.SegmentPort{
			Id: str("p2"), DisplayName: str("web-b"), Path: str("/infra/segments/seg-1/ports/p2"),
			ResourceType: &resourceType,
		})
		stub := &seqQueryListClient{responses: []model.SearchResponse{{
			Results: []*data.StructValue{p1, p2}, ResultCount: i64(2),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "web",
		})

		err := dataSourceNsxtPolicySegmentPortRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("none of the selectors set", func(t *testing.T) {
		ds := dataSourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicySegmentPortRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Atleast one of vif_id, display_name or id")
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicySegmentPort()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": "port-err",
		})

		err := dataSourceNsxtPolicySegmentPortRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "search boom")
	})
}
