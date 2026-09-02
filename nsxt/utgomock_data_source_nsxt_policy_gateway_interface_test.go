//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

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

func gwInterfaceToStructValue(t *testing.T, golangValue interface{}, bindingType bindings.BindingType) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(golangValue, bindingType)
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func TestUnitNsxt_dataSourceNsxtPolicyGatewayInterfaceRead(t *testing.T) {
	t0GwPath := "/infra/tier-0s/t0-1"
	t1GwPath := "/infra/tier-1s/t1-1"

	t.Run("missing gateway_path and service_path errors", func(t *testing.T) {
		ds := dataSourceNsxtPolicyGatewayInterface()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyGatewayInterfaceRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "gateway_path or service_path")
	})

	t.Run("tier0 success via gateway_path", func(t *testing.T) {
		ifaceID := "iface-1"
		edgePath := "/infra/sites/default/enforcement-points/default/edge-clusters/ec1/edge-nodes/en1"
		segPath := "/infra/segments/seg1"
		desc := "t0 interface"
		sv := gwInterfaceToStructValue(t, nsxModel.Tier0Interface{
			Id: str(ifaceID), Path: str(t0GwPath + "/locale-services/default/interfaces/" + ifaceID),
			ResourceType: str("Tier0Interface"), EdgePath: str(edgePath), SegmentPath: str(segPath), Description: str(desc),
		}, nsxModel.Tier0InterfaceBindingType())
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyGatewayInterface()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":           ifaceID,
			"gateway_path": t0GwPath,
		})

		err := dataSourceNsxtPolicyGatewayInterfaceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ifaceID, d.Id())
		assert.Equal(t, edgePath, d.Get("edge_cluster_path"))
		assert.Equal(t, segPath, d.Get("segment_path"))
		assert.Equal(t, desc, d.Get("description"))
	})

	t.Run("tier1 success via service_path", func(t *testing.T) {
		ifaceID := "iface-2"
		segPath := "/infra/segments/seg2"
		desc := "t1 interface"
		servicePath := t1GwPath + "/locale-services/default"
		sv := gwInterfaceToStructValue(t, nsxModel.Tier1Interface{
			Id: str(ifaceID), Path: str(servicePath + "/interfaces/" + ifaceID),
			ResourceType: str("Tier1Interface"), SegmentPath: str(segPath), Description: str(desc),
		}, nsxModel.Tier1InterfaceBindingType())
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyGatewayInterface()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":           ifaceID,
			"service_path": servicePath,
		})

		err := dataSourceNsxtPolicyGatewayInterfaceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ifaceID, d.Id())
		assert.Equal(t, segPath, d.Get("segment_path"))
		assert.Equal(t, desc, d.Get("description"))
		assert.Equal(t, "", d.Get("edge_cluster_path"))
	})

	t.Run("invalid gateway_path errors", func(t *testing.T) {
		ds := dataSourceNsxtPolicyGatewayInterface()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":           "iface-x",
			"gateway_path": "bad",
		})

		err := dataSourceNsxtPolicyGatewayInterfaceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search failed")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyGatewayInterface()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":           "iface-1",
			"gateway_path": t0GwPath,
		})

		err := dataSourceNsxtPolicyGatewayInterfaceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
