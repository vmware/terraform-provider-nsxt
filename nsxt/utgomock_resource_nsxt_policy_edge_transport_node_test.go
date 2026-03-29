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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	enforcementpoints "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	policyETNID         = "etn-1"
	policyETNSitePath   = "/infra/sites/default"
	policyETNSiteID     = "default"
	policyETNEPID       = "default"
	policyETNName       = "etn-fooname"
	policyETNRevision   = int64(1)
	policyETNPath       = "/infra/sites/default/enforcement-points/default/edge-transport-nodes/etn-1"
	policyETNParentPath = "/infra/sites/default/enforcement-points/default"
	policyETNHostname   = "edge-host-1.example.com"
)

// minimalPolicyETNSwitchData returns a minimal valid switch schema map for
// schema.TestResourceDataRaw. The pnic list is required with at least one entry.
func minimalPolicyETNSwitchData() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"overlay_transport_zone_path": "",
			"pnic": []interface{}{
				map[string]interface{}{
					"datapath_network_id": "",
					"device_name":         "fp-eth0",
					"uplink_name":         "uplink1",
				},
			},
			"uplink_host_switch_profile_path": "",
			"lldp_host_switch_profile_path":   "",
			"switch_name":                     "nsxDefaultHostSwitch",
			"tunnel_endpoint":                 []interface{}{},
			"vlan_transport_zone_paths":       []interface{}{},
		},
	}
}

// policyETNAPIResponse builds a PolicyEdgeTransportNode object suitable for
// mock Get responses used in Read tests (ParentPath must be set).
func policyETNAPIResponse() model.PolicyEdgeTransportNode {
	hostname := policyETNHostname
	parentPath := policyETNParentPath
	switchName := "nsxDefaultHostSwitch"
	return model.PolicyEdgeTransportNode{
		DisplayName: &policyETNName,
		Path:        &policyETNPath,
		Revision:    &policyETNRevision,
		Hostname:    &hostname,
		ParentPath:  &parentPath,
		SwitchSpec: &model.PolicyEdgeTransportNodeSwitchSpec{
			Switches: []model.PolicyEdgeTransportNodeSwitch{
				{SwitchName: &switchName},
			},
		},
	}
}

func TestMockResourceNsxtPolicyEdgeTransportNodeCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockETNSDK := epmocks.NewMockEdgeTransportNodesClient(ctrl)
	mockSitesSDK := inframocks.NewMockSitesClient(ctrl)

	etnWrapper := &enforcementpoints.PolicyEdgeTransportNodeClientContext{
		Client:     mockETNSDK,
		ClientType: utl.Local,
	}
	siteWrapper := &cliinfra.SiteClientContext{
		Client:     mockSitesSDK,
		ClientType: utl.Local,
	}

	originalETN := cliEdgeTransportNodesClient
	originalSites := cliSitesClient
	defer func() {
		cliEdgeTransportNodesClient = originalETN
		cliSitesClient = originalSites
	}()
	cliEdgeTransportNodesClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyEdgeTransportNodeClientContext {
		return etnWrapper
	}
	cliSitesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SiteClientContext {
		return siteWrapper
	}

	res := resourceNsxtPolicyEdgeTransportNode()

	t.Run("Create_success", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(policyETNSiteID).Return(model.Site{}, nil)
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, gomock.Any()).Return(model.PolicyEdgeTransportNode{}, vapiErrors.NotFound{})
		mockETNSDK.EXPECT().Patch(policyETNSiteID, policyETNEPID, gomock.Any(), gomock.Any()).Return(nil)
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, gomock.Any()).Return(policyETNAPIResponse(), nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":      policyETNName,
			"hostname":          policyETNHostname,
			"site_path":         policyETNSitePath,
			"enforcement_point": policyETNEPID,
			"switch":            minimalPolicyETNSwitchData(),
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, policyETNName, d.Get("display_name"))
	})

	t.Run("Create_fails_when_already_exists", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(policyETNSiteID).Return(model.Site{}, nil)
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(model.PolicyEdgeTransportNode{}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"nsx_id":            policyETNID,
			"hostname":          policyETNHostname,
			"site_path":         policyETNSitePath,
			"enforcement_point": policyETNEPID,
			"switch":            minimalPolicyETNSwitchData(),
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create_fails_when_site_path_invalid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"hostname":          policyETNHostname,
			"site_path":         "/infra/no-sites/default",
			"enforcement_point": policyETNEPID,
			"switch":            minimalPolicyETNSwitchData(),
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Site ID")
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(policyETNSiteID).Return(model.Site{}, nil)
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, gomock.Any()).Return(model.PolicyEdgeTransportNode{}, vapiErrors.NotFound{})
		mockETNSDK.EXPECT().Patch(policyETNSiteID, policyETNEPID, gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"hostname":          policyETNHostname,
			"site_path":         policyETNSitePath,
			"enforcement_point": policyETNEPID,
			"switch":            minimalPolicyETNSwitchData(),
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyEdgeTransportNodeRead(t *testing.T) {
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

	res := resourceNsxtPolicyEdgeTransportNode()

	t.Run("Read_success", func(t *testing.T) {
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(policyETNAPIResponse(), nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         policyETNSitePath,
			"enforcement_point": policyETNEPID,
			"hostname":          policyETNHostname,
			"switch":            minimalPolicyETNSwitchData(),
		})
		d.SetId(policyETNID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, policyETNName, d.Get("display_name"))
		assert.Equal(t, policyETNHostname, d.Get("hostname"))
	})

	t.Run("Read_clears_id_when_not_found", func(t *testing.T) {
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(model.PolicyEdgeTransportNode{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         policyETNSitePath,
			"enforcement_point": policyETNEPID,
			"hostname":          policyETNHostname,
			"switch":            minimalPolicyETNSwitchData(),
		})
		d.SetId(policyETNID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyEdgeTransportNodeUpdate(t *testing.T) {
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

	res := resourceNsxtPolicyEdgeTransportNode()

	t.Run("Update_success", func(t *testing.T) {
		mockETNSDK.EXPECT().Patch(policyETNSiteID, policyETNEPID, policyETNID, gomock.Any()).Return(nil)
		mockETNSDK.EXPECT().Get(policyETNSiteID, policyETNEPID, policyETNID).Return(policyETNAPIResponse(), nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":      policyETNName,
			"hostname":          policyETNHostname,
			"site_path":         policyETNSitePath,
			"enforcement_point": policyETNEPID,
			"switch":            minimalPolicyETNSwitchData(),
		})
		d.SetId(policyETNID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_Patch_returns_error", func(t *testing.T) {
		mockETNSDK.EXPECT().Patch(policyETNSiteID, policyETNEPID, policyETNID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"hostname":          policyETNHostname,
			"site_path":         policyETNSitePath,
			"enforcement_point": policyETNEPID,
			"switch":            minimalPolicyETNSwitchData(),
		})
		d.SetId(policyETNID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyEdgeTransportNodeDelete(t *testing.T) {
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

	res := resourceNsxtPolicyEdgeTransportNode()

	t.Run("Delete_success", func(t *testing.T) {
		mockETNSDK.EXPECT().Delete(policyETNSiteID, policyETNEPID, policyETNID, gomock.Any()).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         policyETNSitePath,
			"enforcement_point": policyETNEPID,
			"hostname":          policyETNHostname,
			"switch":            minimalPolicyETNSwitchData(),
		})
		d.SetId(policyETNID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockETNSDK.EXPECT().Delete(policyETNSiteID, policyETNEPID, policyETNID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         policyETNSitePath,
			"enforcement_point": policyETNEPID,
			"hostname":          policyETNHostname,
			"switch":            minimalPolicyETNSwitchData(),
		})
		d.SetId(policyETNID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeTransportNodeDelete(d, m)
		require.Error(t, err)
	})
}
