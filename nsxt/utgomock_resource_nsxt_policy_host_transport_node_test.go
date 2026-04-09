//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	enforcement_points "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points"
)

var (
	htnID          = "htn-1"
	htnDisplayName = "Test HTN"
	htnDescription = "test host transport node"
	htnSitePath    = "/infra/sites/default"
	htnSiteID      = "default"
	htnEPID        = "default"
	htnDiscNodeID  = "host-1"
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
