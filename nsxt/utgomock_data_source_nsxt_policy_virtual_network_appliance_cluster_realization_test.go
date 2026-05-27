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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	enforcementpoints "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var vnaClusterRealizationPath = "/infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/vna-cluster-1"

// setupVNARealizationMocks sets up all three client mocks used by the
// realization data source and returns a restore function.
func setupVNARealizationMocks(
	ctrl *gomock.Controller,
	mockClusterState *epmocks.MockVirtualNetworkApplianceClusterStateClient,
	mockVNAs *epmocks.MockVirtualNetworkAppliancesInClusterClient,
	mockVNAState *epmocks.MockVirtualNetworkApplianceStateClient,
) func() {
	origClusterState := cliVNAClusterStateClient
	origVNAs := cliVNAsClient
	origVNAState := cliVNAStateClient

	clusterStateWrapper := &enforcementpoints.VirtualNetworkApplianceClusterStateClientContext{
		Client:     mockClusterState,
		ClientType: utl.Local,
	}
	vnasWrapper := &enforcementpoints.VirtualNetworkAppliancesClientContext{
		Client:     mockVNAs,
		ClientType: utl.Local,
	}
	vnaStateWrapper := &enforcementpoints.VirtualNetworkApplianceStateClientContext{
		Client:     mockVNAState,
		ClientType: utl.Local,
	}

	cliVNAClusterStateClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *enforcementpoints.VirtualNetworkApplianceClusterStateClientContext {
		return clusterStateWrapper
	}
	cliVNAsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *enforcementpoints.VirtualNetworkAppliancesClientContext {
		return vnasWrapper
	}
	cliVNAStateClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *enforcementpoints.VirtualNetworkApplianceStateClientContext {
		return vnaStateWrapper
	}

	return func() {
		cliVNAClusterStateClient = origClusterState
		cliVNAsClient = origVNAs
		cliVNAStateClient = origVNAState
	}
}

func TestMockDataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead(t *testing.T) {
	util.NsxVersion = "9.1.1"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClusterState := epmocks.NewMockVirtualNetworkApplianceClusterStateClient(ctrl)
	mockVNAs := epmocks.NewMockVirtualNetworkAppliancesInClusterClient(ctrl)
	mockVNAState := epmocks.NewMockVirtualNetworkApplianceStateClient(ctrl)

	restore := setupVNARealizationMocks(ctrl, mockClusterState, mockVNAs, mockVNAState)
	defer restore()

	ds := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealization()

	statusSuccess := model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_SUCCESS
	statusError := model.VirtualNetworkApplianceClusterState_CONSOLIDATED_STATUS_ERROR
	vnaStatusSuccess := model.VirtualNetworkApplianceConfigurationState_CONSOLIDATED_STATUS_SUCCESS
	zeroCount := int64(0)

	// Stage 1 success + Stage 2 no VNAs deployed (no-op)
	t.Run("Read_success_no_vnas", func(t *testing.T) {
		mockClusterState.EXPECT().Get(gomock.Any(), gomock.Any(), vnaClusterID).Return(
			model.VirtualNetworkApplianceClusterState{ConsolidatedStatus: &statusSuccess}, nil,
		)
		mockVNAs.EXPECT().List(gomock.Any(), gomock.Any(), vnaClusterID,
			nil, nil, nil, nil, nil, nil, nil, nil, nil).Return(
			model.VirtualNetworkApplianceListResult{ResultCount: &zeroCount, Results: []model.VirtualNetworkAppliance{}}, nil,
		)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path": vnaClusterRealizationPath,
		})
		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, statusSuccess, d.Get("state"))
		assert.Empty(t, d.Get("vna_paths"))
	})

	// Stage 1 success + Stage 2 with one VNA reaching SUCCESS
	t.Run("Read_success_with_vna", func(t *testing.T) {
		vnaID := "vna-appliance-1"
		vnaPath := vnaClusterRealizationPath + "/virtual-network-appliances/" + vnaID
		mockClusterState.EXPECT().Get(gomock.Any(), gomock.Any(), vnaClusterID).Return(
			model.VirtualNetworkApplianceClusterState{ConsolidatedStatus: &statusSuccess}, nil,
		)
		oneCount := int64(1)
		mockVNAs.EXPECT().List(gomock.Any(), gomock.Any(), vnaClusterID,
			nil, nil, nil, nil, nil, nil, nil, nil, nil).Return(
			model.VirtualNetworkApplianceListResult{
				ResultCount: &oneCount,
				Results:     []model.VirtualNetworkAppliance{{Id: &vnaID, Path: &vnaPath}},
			}, nil,
		)
		mockVNAState.EXPECT().Get(gomock.Any(), gomock.Any(), vnaClusterID, vnaID).Return(
			model.NetworkApplianceState{
				VirtualNetworkApplianceState: &model.VirtualNetworkApplianceState{
					ConfigurationState: &model.VirtualNetworkApplianceConfigurationState{
						ConsolidatedStatus: &vnaStatusSuccess,
					},
				},
			}, nil,
		)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path": vnaClusterRealizationPath,
		})
		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, statusSuccess, d.Get("state"))
		assert.Equal(t, []interface{}{vnaPath}, d.Get("vna_paths"))
	})

	// Stage 1 stops at ERROR (terminal), Stage 2 still lists VNAs (no-op here)
	t.Run("Read_stops_at_cluster_error_state", func(t *testing.T) {
		mockClusterState.EXPECT().Get(gomock.Any(), gomock.Any(), vnaClusterID).Return(
			model.VirtualNetworkApplianceClusterState{ConsolidatedStatus: &statusError}, nil,
		)
		mockVNAs.EXPECT().List(gomock.Any(), gomock.Any(), vnaClusterID,
			nil, nil, nil, nil, nil, nil, nil, nil, nil).Return(
			model.VirtualNetworkApplianceListResult{ResultCount: &zeroCount, Results: []model.VirtualNetworkAppliance{}}, nil,
		)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path": vnaClusterRealizationPath,
		})
		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, statusError, d.Get("state"))
	})

	// Stage 1 API error → WaitForState returns error
	t.Run("Read_fails_on_cluster_api_error", func(t *testing.T) {
		mockClusterState.EXPECT().Get(gomock.Any(), gomock.Any(), vnaClusterID).Return(
			model.VirtualNetworkApplianceClusterState{}, vapiErrors.InternalServerError{},
		)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path":    vnaClusterRealizationPath,
			"timeout": 1,
		})
		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	// Stage 2 VNA state API error
	t.Run("Read_fails_on_vna_state_api_error", func(t *testing.T) {
		vnaID := "vna-appliance-1"
		mockClusterState.EXPECT().Get(gomock.Any(), gomock.Any(), vnaClusterID).Return(
			model.VirtualNetworkApplianceClusterState{ConsolidatedStatus: &statusSuccess}, nil,
		)
		oneCount := int64(1)
		mockVNAs.EXPECT().List(gomock.Any(), gomock.Any(), vnaClusterID,
			nil, nil, nil, nil, nil, nil, nil, nil, nil).Return(
			model.VirtualNetworkApplianceListResult{
				ResultCount: &oneCount,
				Results:     []model.VirtualNetworkAppliance{{Id: &vnaID}},
			}, nil,
		)
		mockVNAState.EXPECT().Get(gomock.Any(), gomock.Any(), vnaClusterID, vnaID).Return(
			model.NetworkApplianceState{}, vapiErrors.InternalServerError{},
		)

		// timeout must be long enough for Stage 2's 2-second delay to fire.
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path":    vnaClusterRealizationPath,
			"timeout": 10,
		})
		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	// Invalid path (no sites segment)
	t.Run("Read_fails_when_path_has_no_site", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path": "/infra/no-sites/default/enforcement-points/default/virtual-network-appliance-clusters/id",
		})
		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "site ID")
	})
}
