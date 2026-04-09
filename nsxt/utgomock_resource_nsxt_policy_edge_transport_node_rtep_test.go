//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	enforcementpoints "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	rtepPolicyETNPath    = "/infra/sites/default/enforcement-points/default/edge-transport-nodes/etn-1"
	rtepPolicySwitchName = "nsxDefaultHostSwitch"
	rtepPolicyVlan       = int64(100)
	rtepPolicyTeaming    = "teaming-policy-1"
	rtepPolicyPoolPath   = "/infra/ip-pools/rtep-pool"
)

// policyETNWithSwitch returns a PolicyEdgeTransportNode with the given switch name
// and optionally a single RTEP config on that switch.
func policyETNWithSwitch(switchName string, hasRTEP bool) model.PolicyEdgeTransportNode {
	hostname := policyETNHostname
	sw := model.PolicyEdgeTransportNodeSwitch{
		SwitchName: &switchName,
	}
	if hasRTEP {
		namedTeaming := rtepPolicyTeaming
		vlan := rtepPolicyVlan
		sw.RemoteTunnelEndpoint = []model.PolicyEdgeTransportNodeRtepConfig{
			{
				NamedTeamingPolicy: &namedTeaming,
				Vlan:               &vlan,
			},
		}
	}
	return model.PolicyEdgeTransportNode{
		Hostname: &hostname,
		SwitchSpec: &model.PolicyEdgeTransportNodeSwitchSpec{
			Switches: []model.PolicyEdgeTransportNodeSwitch{sw},
		},
	}
}

// minimalPolicyRTEPData returns schema raw data for the RTEP resource with a
// static_ipv4_pool ip_assignment.
func minimalPolicyRTEPData() map[string]interface{} {
	return map[string]interface{}{
		"edge_transport_node_path": rtepPolicyETNPath,
		"host_switch_name":         rtepPolicySwitchName,
		"named_teaming_policy":     rtepPolicyTeaming,
		"vlan":                     int(rtepPolicyVlan),
		"ip_assignment": []interface{}{
			map[string]interface{}{
				"static_ipv4_pool": rtepPolicyPoolPath,
				"static_ipv4_list": []interface{}{},
			},
		},
	}
}

func TestMockResourceNsxtPolicyEdgeTransportNodeRTEPCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockETNSDK := epmocks.NewMockEdgeTransportNodesClient(ctrl)
	etnWrapper := &enforcementpoints.PolicyEdgeTransportNodeClientContext{
		Client:     mockETNSDK,
		ClientType: utl.Local,
	}

	originalETN := cliEdgeTransportNodesClient
	defer func() { cliEdgeTransportNodesClient = originalETN }()
	cliEdgeTransportNodesClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyEdgeTransportNodeClientContext {
		return etnWrapper
	}

	res := resourceNsxtPolicyEdgeTransportNodeRTEP()

	t.Run("Create_success", func(t *testing.T) {
		// Get (exists check) -> no RTEP, Patch to add RTEP, Get (read) -> with RTEP
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(policyETNWithSwitch(rtepPolicySwitchName, false), nil)
		mockETNSDK.EXPECT().Patch(policyETNSiteID, policyETNEPID, policyETNID, gomock.Any()).Return(nil)
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(policyETNWithSwitch(rtepPolicySwitchName, true), nil)

		d := schema.TestResourceDataRaw(t, res.Schema, minimalPolicyRTEPData())
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRTEPCreate(d, m)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%s:%s", rtepPolicyETNPath, rtepPolicySwitchName), d.Id())
	})

	t.Run("Create_fails_when_RTEP_already_exists", func(t *testing.T) {
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(policyETNWithSwitch(rtepPolicySwitchName, true), nil)

		d := schema.TestResourceDataRaw(t, res.Schema, minimalPolicyRTEPData())
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRTEPCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create_fails_when_switch_not_found", func(t *testing.T) {
		// Return ETN with a different switch name
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(policyETNWithSwitch("other-switch", false), nil)

		d := schema.TestResourceDataRaw(t, res.Schema, minimalPolicyRTEPData())
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRTEPCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Create_fails_when_Get_returns_error", func(t *testing.T) {
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(model.PolicyEdgeTransportNode{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, minimalPolicyRTEPData())
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRTEPCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyEdgeTransportNodeRTEPRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockETNSDK := epmocks.NewMockEdgeTransportNodesClient(ctrl)
	etnWrapper := &enforcementpoints.PolicyEdgeTransportNodeClientContext{
		Client:     mockETNSDK,
		ClientType: utl.Local,
	}

	originalETN := cliEdgeTransportNodesClient
	defer func() { cliEdgeTransportNodesClient = originalETN }()
	cliEdgeTransportNodesClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyEdgeTransportNodeClientContext {
		return etnWrapper
	}

	res := resourceNsxtPolicyEdgeTransportNodeRTEP()

	t.Run("Read_success", func(t *testing.T) {
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(policyETNWithSwitch(rtepPolicySwitchName, true), nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_transport_node_path": rtepPolicyETNPath,
			"host_switch_name":         rtepPolicySwitchName,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRTEPRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, rtepPolicyTeaming, d.Get("named_teaming_policy"))
		assert.Equal(t, int(rtepPolicyVlan), d.Get("vlan"))
	})

	t.Run("Read_fails_when_no_RTEP_on_switch", func(t *testing.T) {
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(policyETNWithSwitch(rtepPolicySwitchName, false), nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_transport_node_path": rtepPolicyETNPath,
			"host_switch_name":         rtepPolicySwitchName,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRTEPRead(d, m)
		require.Error(t, err)
	})

	t.Run("Read_fails_when_Get_returns_error", func(t *testing.T) {
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(model.PolicyEdgeTransportNode{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_transport_node_path": rtepPolicyETNPath,
			"host_switch_name":         rtepPolicySwitchName,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRTEPRead(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyEdgeTransportNodeRTEPUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockETNSDK := epmocks.NewMockEdgeTransportNodesClient(ctrl)
	etnWrapper := &enforcementpoints.PolicyEdgeTransportNodeClientContext{
		Client:     mockETNSDK,
		ClientType: utl.Local,
	}

	originalETN := cliEdgeTransportNodesClient
	defer func() { cliEdgeTransportNodesClient = originalETN }()
	cliEdgeTransportNodesClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyEdgeTransportNodeClientContext {
		return etnWrapper
	}

	res := resourceNsxtPolicyEdgeTransportNodeRTEP()

	t.Run("Update_success", func(t *testing.T) {
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(policyETNWithSwitch(rtepPolicySwitchName, true), nil)
		mockETNSDK.EXPECT().Update(policyETNSiteID, policyETNEPID, policyETNID, gomock.Any()).Return(policyETNWithSwitch(rtepPolicySwitchName, true), nil)
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(policyETNWithSwitch(rtepPolicySwitchName, true), nil)

		d := schema.TestResourceDataRaw(t, res.Schema, minimalPolicyRTEPData())
		d.SetId(fmt.Sprintf("%s:%s", rtepPolicyETNPath, rtepPolicySwitchName))
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRTEPUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_switch_not_found", func(t *testing.T) {
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(policyETNWithSwitch("other-switch", true), nil)

		d := schema.TestResourceDataRaw(t, res.Schema, minimalPolicyRTEPData())
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRTEPUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Update_fails_when_Update_returns_error", func(t *testing.T) {
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(policyETNWithSwitch(rtepPolicySwitchName, true), nil)
		mockETNSDK.EXPECT().Update(policyETNSiteID, policyETNEPID, policyETNID, gomock.Any()).Return(model.PolicyEdgeTransportNode{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, minimalPolicyRTEPData())
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRTEPUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyEdgeTransportNodeRTEPDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockETNSDK := epmocks.NewMockEdgeTransportNodesClient(ctrl)
	etnWrapper := &enforcementpoints.PolicyEdgeTransportNodeClientContext{
		Client:     mockETNSDK,
		ClientType: utl.Local,
	}

	originalETN := cliEdgeTransportNodesClient
	defer func() { cliEdgeTransportNodesClient = originalETN }()
	cliEdgeTransportNodesClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyEdgeTransportNodeClientContext {
		return etnWrapper
	}

	res := resourceNsxtPolicyEdgeTransportNodeRTEP()

	t.Run("Delete_success_clears_RTEP_via_Update", func(t *testing.T) {
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(policyETNWithSwitch(rtepPolicySwitchName, true), nil)
		mockETNSDK.EXPECT().Update(policyETNSiteID, policyETNEPID, policyETNID, gomock.Any()).Return(policyETNWithSwitch(rtepPolicySwitchName, false), nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_transport_node_path": rtepPolicyETNPath,
			"host_switch_name":         rtepPolicySwitchName,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRTEPDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_Get_returns_error", func(t *testing.T) {
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(model.PolicyEdgeTransportNode{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_transport_node_path": rtepPolicyETNPath,
			"host_switch_name":         rtepPolicySwitchName,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRTEPDelete(d, m)
		require.Error(t, err)
	})

	t.Run("Delete_fails_when_Update_returns_error", func(t *testing.T) {
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(policyETNWithSwitch(rtepPolicySwitchName, true), nil)
		mockETNSDK.EXPECT().Update(policyETNSiteID, policyETNEPID, policyETNID, gomock.Any()).Return(model.PolicyEdgeTransportNode{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_transport_node_path": rtepPolicyETNPath,
			"host_switch_name":         rtepPolicySwitchName,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRTEPDelete(d, m)
		require.Error(t, err)
	})
}
