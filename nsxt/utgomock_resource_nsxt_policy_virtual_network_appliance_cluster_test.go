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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	enforcementpoints "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vnaClusterID       = "vna-cluster-1"
	vnaClusterSitePath = "/infra/sites/default"
	vnaClusterSiteID   = "default"
	vnaClusterEPID     = "default"
	vnaClusterName     = "vna-cluster-fooname"
	vnaClusterRevision = int64(1)
	vnaClusterPath     = "/infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/vna-cluster-1"
	vnaFormFactor      = model.VirtualNetworkApplianceCluster_APPLIANCE_FORM_FACTOR_MEDIUM
	vnaServiceType     = model.VirtualNetworkApplianceCluster_SERVICE_TYPE_VPC_SERVICES
)

func setupVNAClusterMocks(ctrl *gomock.Controller) (
	*epmocks.MockVirtualNetworkApplianceClustersClient,
	*inframocks.MockSitesClient,
) {
	mockVNAClusters := epmocks.NewMockVirtualNetworkApplianceClustersClient(ctrl)
	mockSites := inframocks.NewMockSitesClient(ctrl)
	return mockVNAClusters, mockSites
}

func setupVNAClientOverride(
	mockVNAClusters *epmocks.MockVirtualNetworkApplianceClustersClient,
	mockSites *inframocks.MockSitesClient,
) func() {
	vnaWrapper := &enforcementpoints.VirtualNetworkApplianceClusterClientContext{
		Client:     mockVNAClusters,
		ClientType: utl.Local,
	}

	origVNA := cliVNAClustersClient
	origSites := cliSitesClient

	cliVNAClustersClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.VirtualNetworkApplianceClusterClientContext {
		return vnaWrapper
	}
	if mockSites != nil {
		siteWrapper := &cliinfra.SiteClientContext{
			Client:     mockSites,
			ClientType: utl.Local,
		}
		cliSitesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SiteClientContext {
			return siteWrapper
		}
	}

	return func() {
		cliVNAClustersClient = origVNA
		cliSitesClient = origSites
	}
}

func TestMockResourceNsxtPolicyVirtualNetworkApplianceClusterCreate(t *testing.T) {
	util.NsxVersion = "9.1.1"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVNAClusters, mockSites := setupVNAClusterMocks(ctrl)
	restore := setupVNAClientOverride(mockVNAClusters, mockSites)
	defer restore()

	res := resourceNsxtPolicyVirtualNetworkApplianceCluster()

	t.Run("Create_success", func(t *testing.T) {
		mockSites.EXPECT().Get(vnaClusterSiteID).Return(model.Site{}, nil)
		mockVNAClusters.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, gomock.Any()).Return(model.VirtualNetworkApplianceCluster{}, vapiErrors.NotFound{})
		mockVNAClusters.EXPECT().Patch(vnaClusterSiteID, vnaClusterEPID, gomock.Any(), gomock.Any()).Return(nil)
		mockVNAClusters.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, gomock.Any()).Return(model.VirtualNetworkApplianceCluster{
			DisplayName:         &vnaClusterName,
			ApplianceFormFactor: &vnaFormFactor,
			ServiceType:         &vnaServiceType,
			Path:                &vnaClusterPath,
			Revision:            &vnaClusterRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":          vnaClusterName,
			"site_path":             vnaClusterSitePath,
			"enforcement_point":     vnaClusterEPID,
			"appliance_form_factor": vnaFormFactor,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceClusterCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, vnaClusterName, d.Get("display_name"))
		assert.Equal(t, vnaFormFactor, d.Get("appliance_form_factor"))
	})

	t.Run("Create_fails_when_already_exists", func(t *testing.T) {
		mockSites.EXPECT().Get(vnaClusterSiteID).Return(model.Site{}, nil)
		mockVNAClusters.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID).Return(model.VirtualNetworkApplianceCluster{}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"nsx_id":            vnaClusterID,
			"site_path":         vnaClusterSitePath,
			"enforcement_point": vnaClusterEPID,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceClusterCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create_fails_when_site_path_invalid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         "/infra/no-sites/default",
			"enforcement_point": vnaClusterEPID,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceClusterCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Site ID")
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockSites.EXPECT().Get(vnaClusterSiteID).Return(model.Site{}, nil)
		mockVNAClusters.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, gomock.Any()).Return(model.VirtualNetworkApplianceCluster{}, vapiErrors.NotFound{})
		mockVNAClusters.EXPECT().Patch(vnaClusterSiteID, vnaClusterEPID, gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         vnaClusterSitePath,
			"enforcement_point": vnaClusterEPID,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceClusterCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyVirtualNetworkApplianceClusterRead(t *testing.T) {
	util.NsxVersion = "9.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVNAClusters, _ := setupVNAClusterMocks(ctrl)

	vnaWrapper := &enforcementpoints.VirtualNetworkApplianceClusterClientContext{
		Client:     mockVNAClusters,
		ClientType: utl.Local,
	}
	origVNA := cliVNAClustersClient
	defer func() { cliVNAClustersClient = origVNA }()
	cliVNAClustersClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.VirtualNetworkApplianceClusterClientContext {
		return vnaWrapper
	}

	res := resourceNsxtPolicyVirtualNetworkApplianceCluster()

	t.Run("Read_success", func(t *testing.T) {
		mockVNAClusters.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID).Return(model.VirtualNetworkApplianceCluster{
			DisplayName:         &vnaClusterName,
			ApplianceFormFactor: &vnaFormFactor,
			ServiceType:         &vnaServiceType,
			Path:                &vnaClusterPath,
			Revision:            &vnaClusterRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         vnaClusterSitePath,
			"enforcement_point": vnaClusterEPID,
		})
		d.SetId(vnaClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceClusterRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, vnaClusterName, d.Get("display_name"))
		assert.Equal(t, vnaFormFactor, d.Get("appliance_form_factor"))
		assert.Equal(t, vnaServiceType, d.Get("service_type"))
	})

	t.Run("Read_success_with_members_and_advanced_config", func(t *testing.T) {
		edgePath := "/infra/sites/default/enforcement-points/default/edge-transport-nodes/edge-1"
		appPath := "/infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/vna-cluster-1/virtual-network-appliances/vna-1"
		uniqueID := "unique-id-123"
		coreProfile := model.VirtualNetworkApplianceClusterAdvancedConfiguration_CORE_ALLOCATION_PROFILE_L4LBSERVICE
		mockVNAClusters.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID).Return(model.VirtualNetworkApplianceCluster{
			DisplayName: &vnaClusterName,
			Path:        &vnaClusterPath,
			Revision:    &vnaClusterRevision,
			Members: []model.VirtualNetworkApplianceClusterMember{
				{
					AppliancePath:         &appPath,
					ApplianceUniqueId:     &uniqueID,
					EdgeTransportNodePath: &edgePath,
				},
			},
			AdvancedConfiguration: &model.VirtualNetworkApplianceClusterAdvancedConfiguration{
				CoreAllocationProfile: &coreProfile,
			},
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         vnaClusterSitePath,
			"enforcement_point": vnaClusterEPID,
		})
		d.SetId(vnaClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceClusterRead(d, m)
		require.NoError(t, err)
		members := d.Get("member").([]interface{})
		require.Len(t, members, 1)
		assert.Equal(t, edgePath, members[0].(map[string]interface{})["edge_transport_node_path"])
		advCfg := d.Get("advanced_configuration").([]interface{})
		require.Len(t, advCfg, 1)
		assert.Equal(t, coreProfile, advCfg[0].(map[string]interface{})["core_allocation_profile"])
	})

	t.Run("Read_fails_when_Get_returns_NotFound", func(t *testing.T) {
		mockVNAClusters.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID).Return(model.VirtualNetworkApplianceCluster{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         vnaClusterSitePath,
			"enforcement_point": vnaClusterEPID,
		})
		d.SetId(vnaClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceClusterRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyVirtualNetworkApplianceClusterUpdate(t *testing.T) {
	util.NsxVersion = "9.1.1"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVNAClusters, _ := setupVNAClusterMocks(ctrl)
	restore := setupVNAClientOverride(mockVNAClusters, nil)
	defer restore()

	res := resourceNsxtPolicyVirtualNetworkApplianceCluster()

	t.Run("Update_success", func(t *testing.T) {
		updatedName := "vna-cluster-updated"
		mockVNAClusters.EXPECT().Update(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, gomock.Any()).Return(model.VirtualNetworkApplianceCluster{}, nil)
		mockVNAClusters.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID).Return(model.VirtualNetworkApplianceCluster{
			DisplayName: &updatedName,
			Path:        &vnaClusterPath,
			Revision:    &vnaClusterRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":      updatedName,
			"site_path":         vnaClusterSitePath,
			"enforcement_point": vnaClusterEPID,
		})
		d.SetId(vnaClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceClusterUpdate(d, m)
		require.NoError(t, err)
		assert.Equal(t, updatedName, d.Get("display_name"))
	})

	t.Run("Update_fails_when_Update_returns_error", func(t *testing.T) {
		mockVNAClusters.EXPECT().Update(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, gomock.Any()).Return(model.VirtualNetworkApplianceCluster{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         vnaClusterSitePath,
			"enforcement_point": vnaClusterEPID,
		})
		d.SetId(vnaClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceClusterUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyVirtualNetworkApplianceClusterDelete(t *testing.T) {
	util.NsxVersion = "9.1.1"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVNAClusters, _ := setupVNAClusterMocks(ctrl)

	vnaWrapper := &enforcementpoints.VirtualNetworkApplianceClusterClientContext{
		Client:     mockVNAClusters,
		ClientType: utl.Local,
	}
	origVNA := cliVNAClustersClient
	defer func() { cliVNAClustersClient = origVNA }()
	cliVNAClustersClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.VirtualNetworkApplianceClusterClientContext {
		return vnaWrapper
	}

	res := resourceNsxtPolicyVirtualNetworkApplianceCluster()

	t.Run("Delete_success", func(t *testing.T) {
		mockVNAClusters.EXPECT().Delete(vnaClusterSiteID, vnaClusterEPID, vnaClusterID).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         vnaClusterSitePath,
			"enforcement_point": vnaClusterEPID,
		})
		d.SetId(vnaClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceClusterDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockVNAClusters.EXPECT().Delete(vnaClusterSiteID, vnaClusterEPID, vnaClusterID).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         vnaClusterSitePath,
			"enforcement_point": vnaClusterEPID,
		})
		d.SetId(vnaClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyVirtualNetworkApplianceClusterDelete(d, m)
		require.Error(t, err)
	})
}

func TestUnitNsxtPolicyVirtualNetworkApplianceCluster_getVNAClusterMembersFromSchema(t *testing.T) {
	t.Run("empty_members", func(t *testing.T) {
		res := resourceNsxtPolicyVirtualNetworkApplianceCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		members := getVNAClusterMembersFromSchema(d)
		assert.Empty(t, members)
	})

	t.Run("with_members", func(t *testing.T) {
		edgePath := "/infra/sites/default/enforcement-points/default/edge-transport-nodes/edge-1"
		res := resourceNsxtPolicyVirtualNetworkApplianceCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"member": []interface{}{
				map[string]interface{}{
					"edge_transport_node_path": edgePath,
					"appliance_path":           "",
					"appliance_unique_id":      "",
				},
			},
		})
		members := getVNAClusterMembersFromSchema(d)
		require.Len(t, members, 1)
		require.NotNil(t, members[0].EdgeTransportNodePath)
		assert.Equal(t, edgePath, *members[0].EdgeTransportNodePath)
	})
}

func TestUnitNsxtPolicyVirtualNetworkApplianceCluster_getVNAClusterAdvancedConfigFromSchema(t *testing.T) {
	t.Run("nil_returns_nil", func(t *testing.T) {
		res := resourceNsxtPolicyVirtualNetworkApplianceCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		cfg, err := getVNAClusterAdvancedConfigFromSchema(d)
		require.NoError(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("with_core_allocation_profile", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtPolicyVirtualNetworkApplianceCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"advanced_configuration": []interface{}{
				map[string]interface{}{
					"core_allocation_profile":     model.VirtualNetworkApplianceClusterAdvancedConfiguration_CORE_ALLOCATION_PROFILE_L4LBSERVICE,
					"high_availability_profile":   "",
					"overlay_transport_zone_path": "",
				},
			},
		})
		cfg, err := getVNAClusterAdvancedConfigFromSchema(d)
		require.NoError(t, err)
		require.NotNil(t, cfg)
		require.NotNil(t, cfg.CoreAllocationProfile)
		assert.Equal(t, model.VirtualNetworkApplianceClusterAdvancedConfiguration_CORE_ALLOCATION_PROFILE_L4LBSERVICE, *cfg.CoreAllocationProfile)
	})

	t.Run("core_allocation_profile_version_too_low", func(t *testing.T) {
		util.NsxVersion = "9.1.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtPolicyVirtualNetworkApplianceCluster()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"advanced_configuration": []interface{}{
				map[string]interface{}{
					"core_allocation_profile":     model.VirtualNetworkApplianceClusterAdvancedConfiguration_CORE_ALLOCATION_PROFILE_L4LBSERVICE,
					"high_availability_profile":   "",
					"overlay_transport_zone_path": "",
				},
			},
		})
		_, err := getVNAClusterAdvancedConfigFromSchema(d)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})
}
