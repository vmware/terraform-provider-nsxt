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
	vapiData "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	enforcement_points "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	host_transport_nodes "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points/host_transport_nodes"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points"
	htnstatemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points/host_transport_nodes"
)

var (
	htnID          = "htn-1"
	htnDisplayName = "Test HTN"
	htnDescription = "test host transport node"
	htnSitePath    = "/infra/sites/default"
	htnSiteID      = "default"
	htnEPID        = "default"
	htnDiscNodeID  = "host-1"
	htnRevision    = int64(1)
)

func minimalHtnData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":          htnDisplayName,
		"description":           htnDescription,
		"nsx_id":                htnID,
		"site_path":             htnSitePath,
		"enforcement_point":     htnEPID,
		"discovered_node_id":    htnDiscNodeID,
		"remove_nsx_on_destroy": false,
	}
}

func setupHtnMock(t *testing.T, ctrl *gomock.Controller) (*epmocks.MockHostTransportNodesClient, func()) {
	mockSDK := epmocks.NewMockHostTransportNodesClient(ctrl)
	mockWrapper := &enforcement_points.HostTransportNodeClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliHostTransportNodesClient
	cliHostTransportNodesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *enforcement_points.HostTransportNodeClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliHostTransportNodesClient = original }
}

func setupHtnSiteMock(t *testing.T, ctrl *gomock.Controller) (*inframocks.MockSitesClient, func()) {
	mockSDK := inframocks.NewMockSitesClient(ctrl)
	mockWrapper := &infraapi.SiteClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliSitesClient
	cliSitesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.SiteClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliSitesClient = original }
}

// TestMockResourceNsxtPolicyHostTransportNodeRead tests the Read function.
// The success path is not fully testable here because setHostSwitchSpecInSchema
// requires a non-nil, properly serialized HostSwitchSpec StructValue that is
// complex to construct in unit tests. We cover the not-found path instead.
func TestMockResourceNsxtPolicyHostTransportNodeRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHtnMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(htnSiteID, htnEPID, htnID).Return(nsxModel.HostTransportNode{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyHostTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtnData())
		d.SetId(htnID)

		err := resourceNsxtPolicyHostTransportNodeRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(htnSiteID, htnEPID, htnID).Return(nsxModel.HostTransportNode{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyHostTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtnData())
		d.SetId(htnID)

		err := resourceNsxtPolicyHostTransportNodeRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

// TestMockResourceNsxtPolicyHostTransportNodeCreate tests the Create function.
// The full success path (including Read) is not tested here because it requires
// a serialized HostSwitchSpec StructValue that is complex to construct in unit tests.
func TestMockResourceNsxtPolicyHostTransportNodeCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHTNSDK, restoreHTN := setupHtnMock(t, ctrl)
	defer restoreHTN()
	mockSiteSDK, restoreSite := setupHtnSiteMock(t, ctrl)
	defer restoreSite()

	t.Run("Create fails when already exists", func(t *testing.T) {
		gomock.InOrder(
			mockSiteSDK.EXPECT().Get(htnSiteID).Return(nsxModel.Site{}, nil),
			mockHTNSDK.EXPECT().Get(htnSiteID, htnEPID, htnID).Return(nsxModel.HostTransportNode{
				Id: &htnID,
			}, nil),
		)

		res := resourceNsxtPolicyHostTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalHtnData())

		err := resourceNsxtPolicyHostTransportNodeCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyHostTransportNodeDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHtnMock(t, ctrl)
	defer restore()

	t.Run("Delete success without NSX removal", func(t *testing.T) {
		mockSDK.EXPECT().Delete(htnSiteID, htnEPID, htnID, gomock.Any(), gomock.Any()).Return(nil)

		res := resourceNsxtPolicyHostTransportNode()
		data := minimalHtnData()
		data["remove_nsx_on_destroy"] = false
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(htnID)

		err := resourceNsxtPolicyHostTransportNodeDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Delete(htnSiteID, htnEPID, htnID, gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyHostTransportNode()
		data := minimalHtnData()
		data["remove_nsx_on_destroy"] = false
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(htnID)

		err := resourceNsxtPolicyHostTransportNodeDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

// htnDataWithHostSwitch returns HTN test data with a minimal (but schema-valid)
// standard_host_switch entry, needed to exercise setHostSwitchSpecInSchema successfully.
func htnDataWithHostSwitch() map[string]interface{} {
	data := minimalHtnData()
	data["standard_host_switch"] = []interface{}{
		map[string]interface{}{
			"host_switch_id": "hostswitch-1",
		},
	}
	return data
}

// buildHtnHostSwitchSpec builds a real *data.StructValue for the standard_host_switch
// configured on d, using the same production helper the resource itself uses when
// building the Patch request body. This lets Read/Create/Update success tests round-trip
// a realistic HostSwitchSpec without hand-authoring the low-level StructValue.
func buildHtnHostSwitchSpec(t *testing.T, d *schema.ResourceData) *vapiData.StructValue {
	spec, err := getHostSwitchSpecFromSchema(d, newGoMockProviderClient(), nodeTypeHost)
	require.NoError(t, err)
	require.NotNil(t, spec)
	return spec
}

func htnAPIResponse(t *testing.T, d *schema.ResourceData) nsxModel.HostTransportNode {
	parentPath := htnSitePath + "/enforcement-points/" + htnEPID
	return nsxModel.HostTransportNode{
		DisplayName:    &htnDisplayName,
		Description:    &htnDescription,
		Path:           &htnSitePath,
		Revision:       &htnRevision,
		ParentPath:     &parentPath,
		HostSwitchSpec: buildHtnHostSwitchSpec(t, d),
		NodeDeploymentInfo: &nsxModel.FabricHostNode{
			DiscoveredNodeId: &htnDiscNodeID,
		},
	}
}

func TestMockResourceNsxtPolicyHostTransportNodeReadSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupHtnMock(t, ctrl)
	defer restore()

	t.Run("Read success sets fields", func(t *testing.T) {
		res := resourceNsxtPolicyHostTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, htnDataWithHostSwitch())
		d.SetId(htnID)

		mockSDK.EXPECT().Get(htnSiteID, htnEPID, htnID).Return(htnAPIResponse(t, d), nil)

		err := resourceNsxtPolicyHostTransportNodeRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, htnDisplayName, d.Get("display_name"))
		assert.Equal(t, htnDescription, d.Get("description"))
		assert.Equal(t, htnDiscNodeID, d.Get("discovered_node_id"))
		assert.Equal(t, htnSitePath, d.Get("site_path"))
		assert.Equal(t, htnEPID, d.Get("enforcement_point"))
		assert.Equal(t, htnID, d.Get("nsx_id"))

		switches := d.Get("standard_host_switch").([]interface{})
		require.Len(t, switches, 1)
		sw := switches[0].(map[string]interface{})
		assert.Equal(t, "hostswitch-1", sw["host_switch_id"])
	})
}

func TestMockResourceNsxtPolicyHostTransportNodeCreateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHTNSDK, restoreHTN := setupHtnMock(t, ctrl)
	defer restoreHTN()
	mockSiteSDK, restoreSite := setupHtnSiteMock(t, ctrl)
	defer restoreSite()

	t.Run("Create success", func(t *testing.T) {
		res := resourceNsxtPolicyHostTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, htnDataWithHostSwitch())

		resp := htnAPIResponse(t, d)

		gomock.InOrder(
			mockSiteSDK.EXPECT().Get(htnSiteID).Return(nsxModel.Site{}, nil),
			mockHTNSDK.EXPECT().Get(htnSiteID, htnEPID, htnID).Return(nsxModel.HostTransportNode{}, vapiErrors.NotFound{}),
			mockHTNSDK.EXPECT().Patch(htnSiteID, htnEPID, htnID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
			mockHTNSDK.EXPECT().Get(htnSiteID, htnEPID, htnID).Return(resp, nil),
		)

		err := resourceNsxtPolicyHostTransportNodeCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, htnID, d.Id())
		assert.Equal(t, htnID, d.Get("nsx_id"))
	})

	t.Run("Create fails when Patch API errors", func(t *testing.T) {
		res := resourceNsxtPolicyHostTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, htnDataWithHostSwitch())

		gomock.InOrder(
			mockSiteSDK.EXPECT().Get(htnSiteID).Return(nsxModel.Site{}, nil),
			mockHTNSDK.EXPECT().Get(htnSiteID, htnEPID, htnID).Return(nsxModel.HostTransportNode{}, vapiErrors.NotFound{}),
			mockHTNSDK.EXPECT().Patch(htnSiteID, htnEPID, htnID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		err := resourceNsxtPolicyHostTransportNodeCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyHostTransportNodeUpdateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHTNSDK, restore := setupHtnMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		res := resourceNsxtPolicyHostTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, htnDataWithHostSwitch())
		d.SetId(htnID)

		resp := htnAPIResponse(t, d)

		gomock.InOrder(
			mockHTNSDK.EXPECT().Patch(htnSiteID, htnEPID, htnID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
			mockHTNSDK.EXPECT().Get(htnSiteID, htnEPID, htnID).Return(resp, nil),
		)

		err := resourceNsxtPolicyHostTransportNodeUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, htnDisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when Patch API errors", func(t *testing.T) {
		res := resourceNsxtPolicyHostTransportNode()
		d := schema.TestResourceDataRaw(t, res.Schema, htnDataWithHostSwitch())
		d.SetId(htnID)

		mockHTNSDK.EXPECT().Patch(htnSiteID, htnEPID, htnID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		err := resourceNsxtPolicyHostTransportNodeUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyHostTransportNodeDeleteWithNsxRemoval(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHTNSDK, restore := setupHtnMock(t, ctrl)
	defer restore()

	mockStateSDK := htnstatemocks.NewMockStateClient(ctrl)
	stateWrapper := &host_transport_nodes.StateClientContext{Client: mockStateSDK, ClientType: utl.Local}
	originalState := cliHostTransportNodeStateClient
	cliHostTransportNodeStateClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *host_transport_nodes.StateClientContext {
		return stateWrapper
	}
	defer func() { cliHostTransportNodeStateClient = originalState }()

	origDelay := hostTransportNodeStatePollDelay
	origInterval := hostTransportNodeStatePollInterval
	origTimeout := hostTransportNodeStatePollTimeout
	hostTransportNodeStatePollDelay = 0
	hostTransportNodeStatePollInterval = 1
	hostTransportNodeStatePollTimeout = 5
	defer func() {
		hostTransportNodeStatePollDelay = origDelay
		hostTransportNodeStatePollInterval = origInterval
		hostTransportNodeStatePollTimeout = origTimeout
	}()

	t.Run("Delete waits for NSX removal to complete", func(t *testing.T) {
		mockHTNSDK.EXPECT().Delete(htnSiteID, htnEPID, htnID, gomock.Any(), gomock.Any()).Return(nil)
		mockStateSDK.EXPECT().Get(htnSiteID, htnEPID, htnID).Return(nsxModel.TransportNodeState{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyHostTransportNode()
		data := minimalHtnData()
		data["remove_nsx_on_destroy"] = true
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(htnID)

		err := resourceNsxtPolicyHostTransportNodeDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
