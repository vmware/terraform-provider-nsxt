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
		assert.Contains(t, err.Error(), "failed to realize")
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
