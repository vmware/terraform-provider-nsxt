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

// vnaPath is defined in utgomock_resource_nsxt_policy_virtual_network_appliance_test.go:
// /infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/vna-cluster-1/virtual-network-appliances/vna-1

func TestMockDataSourceNsxtPolicyVirtualNetworkApplianceRealization(t *testing.T) {
	util.NsxVersion = "9.1.1"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVNAState := epmocks.NewMockVirtualNetworkApplianceStateClient(ctrl)

	origState := cliVNAStateClient
	wrapper := &enforcementpoints.VirtualNetworkApplianceStateClientContext{
		Client:     mockVNAState,
		ClientType: utl.Local,
	}
	cliVNAStateClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *enforcementpoints.VirtualNetworkApplianceStateClientContext {
		return wrapper
	}
	defer func() { cliVNAStateClient = origState }()

	ds := dataSourceNsxtPolicyVirtualNetworkApplianceRealization()

	statusSuccess := model.VirtualNetworkApplianceConfigurationState_CONSOLIDATED_STATUS_SUCCESS
	statusError := model.VirtualNetworkApplianceConfigurationState_CONSOLIDATED_STATUS_ERROR

	t.Run("Read_success", func(t *testing.T) {
		mockVNAState.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID).Return(
			model.NetworkApplianceState{
				VirtualNetworkApplianceState: &model.VirtualNetworkApplianceState{
					ConfigurationState: &model.VirtualNetworkApplianceConfigurationState{
						ConsolidatedStatus: &statusSuccess,
					},
				},
			}, nil,
		)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path": vnaPath,
		})
		err := dataSourceNsxtPolicyVirtualNetworkApplianceRealizationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, statusSuccess, d.Get("state"))
	})

	t.Run("Read_error_state", func(t *testing.T) {
		mockVNAState.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID).Return(
			model.NetworkApplianceState{
				VirtualNetworkApplianceState: &model.VirtualNetworkApplianceState{
					ConfigurationState: &model.VirtualNetworkApplianceConfigurationState{
						ConsolidatedStatus: &statusError,
					},
				},
			}, nil,
		)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path": vnaPath,
		})
		err := dataSourceNsxtPolicyVirtualNetworkApplianceRealizationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, statusError, d.Get("state"))
	})

	t.Run("Read_api_error", func(t *testing.T) {
		mockVNAState.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID).Return(
			model.NetworkApplianceState{}, vapiErrors.InternalServerError{},
		)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path":    vnaPath,
			"timeout": 1,
		})
		err := dataSourceNsxtPolicyVirtualNetworkApplianceRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read_invalid_path_no_vna_segment", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path": "/infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/vna-cluster-1",
		})
		err := dataSourceNsxtPolicyVirtualNetworkApplianceRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "VNA ID")
	})

	t.Run("Read_invalid_path_no_site_segment", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"path": "/infra/enforcement-points/default/virtual-network-appliance-clusters/vna-cluster-1/virtual-network-appliances/vna-1",
		})
		err := dataSourceNsxtPolicyVirtualNetworkApplianceRealizationRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "site ID")
	})
}
