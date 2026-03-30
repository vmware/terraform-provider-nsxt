// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/nsx/EdgeClustersClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/EdgeClustersClient.go EdgeClustersClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"

	apinsx "github.com/vmware/terraform-provider-nsxt/api/nsx"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	nsxmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx"
)

var (
	ecID              = "edge-cluster-1"
	ecDisplayName     = "test-edge-cluster"
	ecDescription     = "Test edge cluster"
	ecRevision        = int64(1)
	ecHaProfileID     = "ha-profile-uuid"
	ecMemberNodeType  = "EdgeNode"
	ecTransportNodeID = "transport-node-uuid-1"
)

// enabledFailureDomainAllocationRules builds the AllocationRules slice that
// the resource's Read path expects when failure_domain_allocation = "enable".
func enabledFailureDomainAllocationRules() []nsxModel.AllocationRule {
	enabled := true
	action := nsxModel.AllocationBasedOnFailureDomain{
		ActionType: nsxModel.AllocationRuleAction_ACTION_TYPE_ALLOCATIONBASEDONFAILUREDOMAIN,
		Enabled:    &enabled,
	}
	converter := bindings.NewTypeConverter()
	actionValue, errs := converter.ConvertToVapi(action, nsxModel.AllocationBasedOnFailureDomainBindingType())
	if errs != nil {
		panic(errs[0])
	}
	return []nsxModel.AllocationRule{{Action: actionValue.(*data.StructValue)}}
}

func edgeClusterAPIResponse() nsxModel.EdgeCluster {
	haProfileResourceType := nsxModel.ClusterProfileTypeIdEntry_RESOURCE_TYPE_EDGEHIGHAVAILABILITYPROFILE
	memberIndex := int64(0)
	memberDisplayName := "edge-node-1"
	memberDescription := "first edge node"
	return nsxModel.EdgeCluster{
		Id:             &ecID,
		DisplayName:    &ecDisplayName,
		Description:    &ecDescription,
		Revision:       &ecRevision,
		MemberNodeType: &ecMemberNodeType,
		ClusterProfileBindings: []nsxModel.ClusterProfileTypeIdEntry{
			{
				ProfileId:    &ecHaProfileID,
				ResourceType: &haProfileResourceType,
			},
		},
		Members: []nsxModel.EdgeClusterMember{
			{
				TransportNodeId: &ecTransportNodeID,
				DisplayName:     &memberDisplayName,
				Description:     &memberDescription,
				MemberIndex:     &memberIndex,
			},
		},
	}
}

func edgeClusterAPIResponseWithFailureDomain() nsxModel.EdgeCluster {
	resp := edgeClusterAPIResponse()
	resp.AllocationRules = enabledFailureDomainAllocationRules()
	return resp
}

func minimalEdgeClusterData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":              ecDisplayName,
		"description":               ecDescription,
		"edge_ha_profile_id":        ecHaProfileID,
		"failure_domain_allocation": "",
		"member": []interface{}{
			map[string]interface{}{
				"transport_node_id": ecTransportNodeID,
				"display_name":      "edge-node-1",
				"description":       "first edge node",
			},
		},
	}
}

func setupEdgeClusterMock(t *testing.T, ctrl *gomock.Controller) (*nsxmocks.MockEdgeClustersClient, func()) {
	mockSDK := nsxmocks.NewMockEdgeClustersClient(ctrl)
	mockWrapper := &apinsx.EdgeClusterClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliEdgeClustersClient
	cliEdgeClustersClient = func(_ utl.SessionContext, _ client.Connector) *apinsx.EdgeClusterClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliEdgeClustersClient = originalCli }
}

func TestMockResourceNsxtEdgeClusterCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Create(gomock.Any()).Return(edgeClusterAPIResponse(), nil)
		mockSDK.EXPECT().Get(ecID).Return(edgeClusterAPIResponse(), nil)

		res := resourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalEdgeClusterData())

		err := resourceNsxtEdgeClusterCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ecID, d.Id())
		assert.Equal(t, ecDisplayName, d.Get("display_name"))
		assert.Equal(t, ecHaProfileID, d.Get("edge_ha_profile_id"))
	})

	t.Run("Create success with failure_domain_allocation enabled", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Create(gomock.Any()).Return(edgeClusterAPIResponseWithFailureDomain(), nil)
		mockSDK.EXPECT().Get(ecID).Return(edgeClusterAPIResponseWithFailureDomain(), nil)

		res := resourceNsxtEdgeCluster()
		data := minimalEdgeClusterData()
		data["failure_domain_allocation"] = "enable"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtEdgeClusterCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ecID, d.Id())
		assert.Equal(t, "enable", d.Get("failure_domain_allocation"))
	})

	t.Run("Create fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Create(gomock.Any()).Return(nsxModel.EdgeCluster{}, errors.New("create API error"))

		res := resourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalEdgeClusterData())

		err := resourceNsxtEdgeClusterCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "create API error")
	})
}

func TestMockResourceNsxtEdgeClusterRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(ecID).Return(edgeClusterAPIResponse(), nil)

		res := resourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ecID)

		err := resourceNsxtEdgeClusterRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ecDisplayName, d.Get("display_name"))
		assert.Equal(t, ecDescription, d.Get("description"))
		assert.Equal(t, int(ecRevision), d.Get("revision"))
		assert.Equal(t, ecHaProfileID, d.Get("edge_ha_profile_id"))
		assert.Equal(t, ecMemberNodeType, d.Get("member_node_type"))
	})

	t.Run("Read populates failure_domain_allocation from AllocationRules", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(ecID).Return(edgeClusterAPIResponseWithFailureDomain(), nil)

		res := resourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ecID)

		err := resourceNsxtEdgeClusterRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "enable", d.Get("failure_domain_allocation"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtEdgeClusterRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Read fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(ecID).Return(nsxModel.EdgeCluster{}, errors.New("read API error"))

		res := resourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ecID)

		err := resourceNsxtEdgeClusterRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read API error")
	})
}

func TestMockResourceNsxtEdgeClusterUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Update(ecID, gomock.Any()).Return(edgeClusterAPIResponse(), nil)
		mockSDK.EXPECT().Get(ecID).Return(edgeClusterAPIResponse(), nil)

		res := resourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalEdgeClusterData())
		d.SetId(ecID)

		err := resourceNsxtEdgeClusterUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ecDisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalEdgeClusterData())

		err := resourceNsxtEdgeClusterUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Update fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Update(ecID, gomock.Any()).Return(nsxModel.EdgeCluster{}, errors.New("update API error"))

		res := resourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalEdgeClusterData())
		d.SetId(ecID)

		err := resourceNsxtEdgeClusterUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "update API error")
	})
}

func TestMockResourceNsxtEdgeClusterDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(ecID).Return(nil)

		res := resourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ecID)

		err := resourceNsxtEdgeClusterDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtEdgeClusterDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupEdgeClusterMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(ecID).Return(errors.New("delete API error"))

		res := resourceNsxtEdgeCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ecID)

		err := resourceNsxtEdgeClusterDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete API error")
	})
}
