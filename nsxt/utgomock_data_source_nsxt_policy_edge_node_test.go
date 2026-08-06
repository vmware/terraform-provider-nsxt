//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// TestMockDataSourceNsxtPolicyEdgeNodeReadSearchPath exercises the
// search-based code path. It also documents the known "id" misbehavior: the
// value supplied as input (an nsx_id) is not what ends up in the id
// attribute - NSX reports back its own id for the resource, which in
// practice is the node's member index.
func TestMockDataSourceNsxtPolicyEdgeNodeReadSearchPath(t *testing.T) {
	resourceType := "PolicyEdgeNode"
	inputNsxID := "edge-transport-node-uuid"
	returnedID := "0" // NSX reports the member index as this resource's id
	displayName := "edge-node-2"
	memberIndex := int64(0)
	uniqueID := "unique-node-2"
	realizationID := "realization-node-2"

	converter := bindings.NewTypeConverter()
	val, errs := converter.ConvertToVapi(model.PolicyEdgeNode{
		Id:            &returnedID,
		DisplayName:   &displayName,
		ResourceType:  &resourceType,
		MemberIndex:   &memberIndex,
		UniqueId:      &uniqueID,
		RealizationId: &realizationID,
	}, model.PolicyEdgeNodeBindingType())
	require.Empty(t, errs)
	sv := val.(*data.StructValue)

	stub := &seqQueryListClient{responses: []model.SearchResponse{{
		Results:     []*data.StructValue{sv},
		ResultCount: i64(1),
	}}}
	defer setupCliQueryClientStub(t, stub)()

	ds := dataSourceNsxtPolicyEdgeNode()
	d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
		"edge_cluster_path": "/infra/sites/default/enforcement-points/default/edge-clusters/cluster1",
		"id":                inputNsxID,
	})

	m := newGoMockProviderClient()
	err := dataSourceNsxtPolicyEdgeNodeRead(d, m)
	require.NoError(t, err)

	assert.Equal(t, returnedID, d.Id())
	assert.NotEqual(t, inputNsxID, d.Id())
	assert.Equal(t, uniqueID, d.Get("unique_id"))
	assert.Equal(t, realizationID, d.Get("realization_id"))
}
