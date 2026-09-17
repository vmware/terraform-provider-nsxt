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
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	edgetransportnodes "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points/edge_transport_nodes"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	etnmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points/edge_transport_nodes"
)

var etnRealizationPath = "/infra/sites/default/enforcement-points/default/edge-transport-nodes/etn-1"

func setupEdgeTransportNodeStateMock(t *testing.T, ctrl *gomock.Controller) (*etnmocks.MockStateClient, func()) {
	mockSDK := etnmocks.NewMockStateClient(ctrl)
	mockWrapper := &edgetransportnodes.StateClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliEdgeTransportNodeStateClient
	cliEdgeTransportNodeStateClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *edgetransportnodes.StateClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliEdgeTransportNodeStateClient = original }
}

func TestMockDataSourceNsxtPolicyEdgeTransportNodeRealizationRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupEdgeTransportNodeStateMock(t, ctrl)
	defer restore()

	ds := dataSourceNsxtPolicyEdgeTransportNodeRealization()

	t.Run("Read success", func(t *testing.T) {
		status := nsxModel.EdgeTnState_CONSOLIDATED_STATUS_SUCCESS
		mockSDK.EXPECT().Get("default", "default", "etn-1").Return(nsxModel.PolicyEdgeTransportNodeState{
			EdgeTnState: &nsxModel.EdgeTnState{
				ConsolidatedStatus: &status,
			},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path":  etnRealizationPath,
			"delay": 0,
		})

		err := dataSourceNsxtPolicyEdgeTransportNodeRealizationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Read failure state", func(t *testing.T) {
		status := nsxModel.EdgeTnState_CONSOLIDATED_STATUS_ERROR
		failureMsg := "boom"
		mockSDK.EXPECT().Get("default", "default", "etn-1").Return(nsxModel.PolicyEdgeTransportNodeState{
			EdgeTnState: &nsxModel.EdgeTnState{
				ConsolidatedStatus: &status,
				FailureMessage:     &failureMsg,
				DeploymentState:    &nsxModel.ConfigurationState{},
			},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path":    etnRealizationPath,
			"timeout": 5,
			"delay":   0,
		})

		err := dataSourceNsxtPolicyEdgeTransportNodeRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to realize. boom")
	})

	t.Run("Read failure state with deployment details failure message", func(t *testing.T) {
		status := nsxModel.EdgeTnState_CONSOLIDATED_STATUS_ERROR
		detailMsg := "[Fabric] Send Edge path=[/infra/sites/default/enforcement-points/default/edge-transport-nodes/etn-1] configuration failed."
		emptyMsg := ""
		mockSDK.EXPECT().Get("default", "default", "etn-1").Return(nsxModel.PolicyEdgeTransportNodeState{
			EdgeTnState: &nsxModel.EdgeTnState{
				ConsolidatedStatus: &status,
				DeploymentState: &nsxModel.ConfigurationState{
					FailureMessage: &emptyMsg,
					Details: []nsxModel.ConfigurationStateElement{
						{
							FailureMessage: &detailMsg,
						},
					},
				},
			},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path":    etnRealizationPath,
			"timeout": 5,
			"delay":   0,
		})

		err := dataSourceNsxtPolicyEdgeTransportNodeRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), detailMsg)
	})

	t.Run("Read failure state with nil deployment state", func(t *testing.T) {
		status := nsxModel.EdgeTnState_CONSOLIDATED_STATUS_ERROR
		failureMsg := "generic failure"
		mockSDK.EXPECT().Get("default", "default", "etn-1").Return(nsxModel.PolicyEdgeTransportNodeState{
			EdgeTnState: &nsxModel.EdgeTnState{
				ConsolidatedStatus: &status,
				FailureMessage:     &failureMsg,
				DeploymentState:    nil,
			},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path":    etnRealizationPath,
			"timeout": 5,
			"delay":   0,
		})

		err := dataSourceNsxtPolicyEdgeTransportNodeRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to realize. generic failure")
	})

	t.Run("Read failure state without messages", func(t *testing.T) {
		status := nsxModel.EdgeTnState_CONSOLIDATED_STATUS_ERROR
		mockSDK.EXPECT().Get("default", "default", "etn-1").Return(nsxModel.PolicyEdgeTransportNodeState{
			EdgeTnState: &nsxModel.EdgeTnState{
				ConsolidatedStatus: &status,
			},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path":    etnRealizationPath,
			"timeout": 5,
			"delay":   0,
		})

		err := dataSourceNsxtPolicyEdgeTransportNodeRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Equal(t, "transport node etn-1 failed to realize", err.Error())
	})

	t.Run("Read API error", func(t *testing.T) {
		mockSDK.EXPECT().Get("default", "default", "etn-1").Return(nsxModel.PolicyEdgeTransportNodeState{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path":    etnRealizationPath,
			"timeout": 5,
			"delay":   0,
		})

		err := dataSourceNsxtPolicyEdgeTransportNodeRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestUnitGetEdgeTnRealizationFailureMessage(t *testing.T) {
	t.Run("nil state returns empty", func(t *testing.T) {
		assert.Equal(t, "", getEdgeTnRealizationFailureMessage(nil))
	})

	t.Run("empty state returns empty", func(t *testing.T) {
		assert.Equal(t, "", getEdgeTnRealizationFailureMessage(&nsxModel.EdgeTnState{}))
	})

	t.Run("deduplicate identical messages across levels", func(t *testing.T) {
		msg := "duplicate error"
		diffMsg := "other error"
		edgeState := &nsxModel.EdgeTnState{
			FailureMessage: &msg,
			DeploymentState: &nsxModel.ConfigurationState{
				FailureMessage: &msg,
				Details: []nsxModel.ConfigurationStateElement{
					{FailureMessage: &msg},
					{FailureMessage: &diffMsg},
				},
			},
			TransportNodeState: &nsxModel.ConfigurationState{
				FailureMessage: &diffMsg,
			},
		}
		result := getEdgeTnRealizationFailureMessage(edgeState)
		assert.Equal(t, "duplicate error, other error", result)
	})

	t.Run("whitespace-only messages are ignored", func(t *testing.T) {
		spaces := "   "
		valid := "real error"
		edgeState := &nsxModel.EdgeTnState{
			FailureMessage: &spaces,
			Details: []nsxModel.ConfigurationStateElement{
				{FailureMessage: &valid},
				{FailureMessage: &spaces},
			},
		}
		result := getEdgeTnRealizationFailureMessage(edgeState)
		assert.Equal(t, "real error", result)
	})
}
