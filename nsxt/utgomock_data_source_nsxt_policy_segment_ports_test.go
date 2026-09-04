//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// dataSourceNsxtPolicySegmentPortsRead delegates to listPolicyResources, which is backed by
// the search (cliQueryClient) machinery already covered by policy_search_unit_test.go. This
// file reuses seqQueryListClient / setupCliQueryClientStub from that file, plus the
// segmentPortToStructValue helper defined in utgomock_data_source_nsxt_policy_segment_port_test.go.

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestMockDataSourceNsxtPolicySegmentPortsRead(t *testing.T) {
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

		ds := dataSourceNsxtPolicySegmentPorts()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"vif_id": vifID,
		})

		err := dataSourceNsxtPolicySegmentPortsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").([]interface{})
		require.Len(t, items, 1)
		item := items[0].(map[string]interface{})
		assert.Equal(t, "port-1", item["id"])
		assert.Equal(t, "/infra/segments/seg-1", item["segment_path"])
	})

	t.Run("success by segment_path and display_name filters items id", func(t *testing.T) {
		match := segmentPortToStructValue(t, model.SegmentPort{
			Id: str("port-a"), DisplayName: str("keep"), Path: str("/infra/segments/seg-1/ports/port-a"),
			ResourceType: &resourceType,
		})
		stub := &seqQueryListClient{responses: []model.SearchResponse{{
			Results: []*data.StructValue{match}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicySegmentPorts()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"segment_path": "/infra/segments/seg-1",
			"display_name": "keep",
		})

		err := dataSourceNsxtPolicySegmentPortsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").([]interface{})
		require.Len(t, items, 1)
		assert.Equal(t, "port-a", items[0].(map[string]interface{})["id"])
	})

	t.Run("success by display_name only, filters client-side for exact match", func(t *testing.T) {
		exact := segmentPortToStructValue(t, model.SegmentPort{
			Id: str("port-exact"), DisplayName: str("web"), Path: str("/infra/segments/seg-1/ports/port-exact"),
			ResourceType: &resourceType,
		})
		stub := &seqQueryListClient{responses: []model.SearchResponse{{
			Results: []*data.StructValue{exact}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicySegmentPorts()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "web",
		})

		err := dataSourceNsxtPolicySegmentPortsRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").([]interface{})
		require.Len(t, items, 1)
	})

	t.Run("none of the selectors set", func(t *testing.T) {
		ds := dataSourceNsxtPolicySegmentPorts()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicySegmentPortsRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "At least one of vif_id, segment_path, or display_name")
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicySegmentPorts()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "err",
		})

		err := dataSourceNsxtPolicySegmentPortsRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "search boom")
	})
}
