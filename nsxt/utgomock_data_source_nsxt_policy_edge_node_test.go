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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	edge_clusters "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points/edge_clusters"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	edgenodemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points/edge_clusters"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func setupEdgeNodeDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*edgenodemocks.MockEdgeNodesClient, func()) {
	t.Helper()
	mockSDK := edgenodemocks.NewMockEdgeNodesClient(ctrl)
	wrapper := &edge_clusters.PolicyEdgeNodeClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliEdgeNodesClient
	cliEdgeNodesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *edge_clusters.PolicyEdgeNodeClientContext {
		return wrapper
	}
	return mockSDK, func() { cliEdgeNodesClient = orig }
}

// TestMockDataSourceNsxtPolicyEdgeNodeReadLocalManager exercises the legacy
// (pre-3.2.0, Local Manager) code path and verifies unique_id/realization_id
// are surfaced from the SDK model.
func TestMockDataSourceNsxtPolicyEdgeNodeReadLocalManager(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupEdgeNodeDataSourceMock(t, ctrl)
	defer restore()
	m := newGoMockProviderClient()
	ep := getPolicyEnforcementPoint(m)

	edgeClusterPath := "/infra/sites/default/enforcement-points/default/edge-clusters/cluster1"
	nodeID := "node-1"
	memberIndex := int64(0)
	uniqueID := "unique-node-1"
	realizationID := "realization-node-1"

	t.Run("by id exposes unique_id and realization_id", func(t *testing.T) {
		mockSDK.EXPECT().Get(defaultSite, ep, "cluster1", nodeID).Return(model.PolicyEdgeNode{
			Id:            &nodeID,
			DisplayName:   &nodeID,
			MemberIndex:   &memberIndex,
			UniqueId:      &uniqueID,
			RealizationId: &realizationID,
		}, nil)

		ds := dataSourceNsxtPolicyEdgeNode()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"edge_cluster_path": edgeClusterPath,
			"id":                nodeID,
		})

		err := dataSourceNsxtPolicyEdgeNodeRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, nodeID, d.Id())
		assert.Equal(t, uniqueID, d.Get("unique_id"))
		assert.Equal(t, realizationID, d.Get("realization_id"))
	})
}

// TestMockDataSourceNsxtPolicyEdgeNodeReadSearchPath exercises the modern
// (GM / NSX >= 3.2.0) search-based code path. It also documents the known
// "id" misbehavior: the value supplied as input (an nsx_id) is not what ends
// up in the id attribute - NSX reports back its own id for the resource,
// which in practice is the node's member index.
func TestMockDataSourceNsxtPolicyEdgeNodeReadSearchPath(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

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
