//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
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
	vnsClusterID       = "vns-cluster-1"
	vnsClusterSitePath = "/infra/sites/default"
	vnsClusterSiteID   = "default"
	vnsClusterEPID     = "default"
	vnsClusterName     = "vns-cluster-fooname"
	vnsClusterRevision = int64(1)
	vnsClusterPath     = "/infra/sites/default/enforcement-points/default/virtual-network-service-clusters/vns-cluster-1"
	vnsFormFactor      = model.VirtualNetworkServiceCluster_NODE_FORM_FACTOR_MEDIUM
	vnsNodes           = int64(3)
	vnsBfdPath         = "/infra/bfd-profiles/bfd1"
	vnsComputeManager  = "cm-1"
)

func vnsClusterFullModel() model.VirtualNetworkServiceCluster {
	alloc := true
	return model.VirtualNetworkServiceCluster{
		Id:             &vnsClusterID,
		DisplayName:    &vnsClusterName,
		Path:           &vnsClusterPath,
		Revision:       &vnsClusterRevision,
		NodeFormFactor: &vnsFormFactor,
		NumberOfNodes:  &vnsNodes,
		BfdProfilePath: &vnsBfdPath,
		AllocationConfig: &model.VirtualNetworkServiceAllocationConfig{
			DefaultProjectServiceAllocation: &alloc,
		},
		VsphereDeploymentConfig: []model.VirtualNetworkServiceDeploymentConfig{
			{
				ComputeManagerId: &vnsComputeManager,
				ClusterIds:       []string{"domain-c1", "domain-c2"},
				FaultDomainIds:   []string{"fd-1"},
			},
		},
	}
}

func setupVNSClientOverride(
	mockVNS *epmocks.MockVirtualNetworkServiceClustersClient,
	mockSites *inframocks.MockSitesClient,
) func() {
	vnsWrapper := &enforcementpoints.VirtualNetworkServiceClusterClientContext{
		Client:     mockVNS,
		ClientType: utl.Local,
	}

	origVNS := cliVNSClustersClient
	origSites := cliSitesClient

	cliVNSClustersClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.VirtualNetworkServiceClusterClientContext {
		return vnsWrapper
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
		cliVNSClustersClient = origVNS
		cliSitesClient = origSites
	}
}

func setupVNSTest(t *testing.T) (*epmocks.MockVirtualNetworkServiceClustersClient, *inframocks.MockSitesClient, func()) {
	t.Helper()
	util.NsxVersion = "9.2.0"
	ctrl := gomock.NewController(t)
	mockVNS := epmocks.NewMockVirtualNetworkServiceClustersClient(ctrl)
	mockSites := inframocks.NewMockSitesClient(ctrl)
	restore := setupVNSClientOverride(mockVNS, mockSites)
	return mockVNS, mockSites, func() {
		restore()
		ctrl.Finish()
		util.NsxVersion = ""
	}
}

func TestMockResourceNsxtPolicyVirtualNetworkServiceClusterCreate(t *testing.T) {
	mockVNS, mockSites, cleanup := setupVNSTest(t)
	defer cleanup()

	res := resourceNsxtPolicyVirtualNetworkServiceCluster()

	t.Run("Create_success", func(t *testing.T) {
		mockSites.EXPECT().Get(vnsClusterSiteID).Return(model.Site{}, nil)
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, gomock.Any()).Return(model.VirtualNetworkServiceCluster{}, vapiErrors.NotFound{})
		mockVNS.EXPECT().Patch(vnsClusterSiteID, vnsClusterEPID, gomock.Any(), gomock.Any()).DoAndReturn(
			func(_, _, _ string, obj model.VirtualNetworkServiceCluster) error {
				require.NotNil(t, obj.NodeFormFactor)
				assert.Equal(t, vnsFormFactor, *obj.NodeFormFactor)
				require.NotNil(t, obj.NumberOfNodes)
				assert.Equal(t, vnsNodes, *obj.NumberOfNodes)
				require.NotNil(t, obj.BfdProfilePath)
				require.NotNil(t, obj.AllocationConfig)
				assert.True(t, *obj.AllocationConfig.DefaultProjectServiceAllocation)
				require.Len(t, obj.VsphereDeploymentConfig, 1)
				assert.Equal(t, vnsComputeManager, *obj.VsphereDeploymentConfig[0].ComputeManagerId)
				assert.ElementsMatch(t, []string{"domain-c1", "domain-c2"}, obj.VsphereDeploymentConfig[0].ClusterIds)
				assert.Equal(t, []string{"fd-1"}, obj.VsphereDeploymentConfig[0].FaultDomainIds)
				return nil
			})
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, gomock.Any()).Return(vnsClusterFullModel(), nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":                       vnsClusterName,
			"site_path":                          vnsClusterSitePath,
			"enforcement_point":                  vnsClusterEPID,
			"node_form_factor":                   vnsFormFactor,
			"number_of_nodes":                    int(vnsNodes),
			"bfd_profile_path":                   vnsBfdPath,
			"default_project_service_allocation": true,
			"vsphere_deployment_config": []interface{}{
				map[string]interface{}{
					"compute_manager_id": vnsComputeManager,
					"cluster_ids":        []interface{}{"domain-c1", "domain-c2"},
					"fault_domain_ids":   []interface{}{"fd-1"},
				},
			},
		})
		err := resourceNsxtPolicyVirtualNetworkServiceClusterCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, vnsClusterName, d.Get("display_name"))
		assert.Equal(t, vnsFormFactor, d.Get("node_form_factor"))
		assert.Equal(t, int(vnsNodes), d.Get("number_of_nodes"))
		assert.Equal(t, true, d.Get("default_project_service_allocation"))
		assert.Len(t, d.Get("vsphere_deployment_config").([]interface{}), 1)
	})

	t.Run("Create_minimal_omits_optional_fields", func(t *testing.T) {
		mockSites.EXPECT().Get(vnsClusterSiteID).Return(model.Site{}, nil)
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, gomock.Any()).Return(model.VirtualNetworkServiceCluster{}, vapiErrors.NotFound{})
		mockVNS.EXPECT().Patch(vnsClusterSiteID, vnsClusterEPID, gomock.Any(), gomock.Any()).DoAndReturn(
			func(_, _, _ string, obj model.VirtualNetworkServiceCluster) error {
				assert.Nil(t, obj.NodeFormFactor)
				assert.Nil(t, obj.NumberOfNodes)
				assert.Nil(t, obj.BfdProfilePath)
				assert.Nil(t, obj.AllocationConfig)
				assert.Empty(t, obj.VsphereDeploymentConfig)
				return nil
			})
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, gomock.Any()).Return(model.VirtualNetworkServiceCluster{
			DisplayName: &vnsClusterName,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":      vnsClusterName,
			"site_path":         vnsClusterSitePath,
			"enforcement_point": vnsClusterEPID,
		})
		require.NoError(t, resourceNsxtPolicyVirtualNetworkServiceClusterCreate(d, newGoMockProviderClient()))
		assert.Empty(t, d.Get("vsphere_deployment_config"))
	})

	t.Run("Create_falls_back_to_provider_enforcement_point_when_empty", func(t *testing.T) {
		mockSites.EXPECT().Get(vnsClusterSiteID).Return(model.Site{}, nil)
		mockVNS.EXPECT().Get(vnsClusterSiteID, gomock.Any(), gomock.Any()).Return(model.VirtualNetworkServiceCluster{}, vapiErrors.NotFound{})
		mockVNS.EXPECT().Patch(vnsClusterSiteID, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockVNS.EXPECT().Get(vnsClusterSiteID, gomock.Any(), gomock.Any()).Return(model.VirtualNetworkServiceCluster{}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path": vnsClusterSitePath,
		})
		d.Set("enforcement_point", "")
		require.NoError(t, resourceNsxtPolicyVirtualNetworkServiceClusterCreate(d, newGoMockProviderClient()))
	})

	t.Run("Create_fails_when_already_exists", func(t *testing.T) {
		mockSites.EXPECT().Get(vnsClusterSiteID).Return(model.Site{}, nil)
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, vnsClusterID).Return(model.VirtualNetworkServiceCluster{}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"nsx_id":            vnsClusterID,
			"site_path":         vnsClusterSitePath,
			"enforcement_point": vnsClusterEPID,
		})
		err := resourceNsxtPolicyVirtualNetworkServiceClusterCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create_fails_when_site_path_invalid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path": "/infra/no-sites/default",
		})
		err := resourceNsxtPolicyVirtualNetworkServiceClusterCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Site ID")
	})

	t.Run("Create_fails_when_site_lookup_fails", func(t *testing.T) {
		mockSites.EXPECT().Get(vnsClusterSiteID).Return(model.Site{}, vapiErrors.ServiceUnavailable{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path": vnsClusterSitePath,
		})
		err := resourceNsxtPolicyVirtualNetworkServiceClusterCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Create_fails_when_exists_check_fails", func(t *testing.T) {
		mockSites.EXPECT().Get(vnsClusterSiteID).Return(model.Site{}, nil)
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, gomock.Any()).Return(model.VirtualNetworkServiceCluster{}, vapiErrors.ServiceUnavailable{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path": vnsClusterSitePath,
		})
		err := resourceNsxtPolicyVirtualNetworkServiceClusterCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockSites.EXPECT().Get(vnsClusterSiteID).Return(model.Site{}, nil)
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, gomock.Any()).Return(model.VirtualNetworkServiceCluster{}, vapiErrors.NotFound{})
		mockVNS.EXPECT().Patch(vnsClusterSiteID, vnsClusterEPID, gomock.Any(), gomock.Any()).Return(errors.New("patch failed"))

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path": vnsClusterSitePath,
		})
		err := resourceNsxtPolicyVirtualNetworkServiceClusterCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Create_fails_when_version_too_low", func(t *testing.T) {
		util.NsxVersion = "9.1.1"
		defer func() { util.NsxVersion = "9.2.0" }()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path": vnsClusterSitePath,
		})
		err := resourceNsxtPolicyVirtualNetworkServiceClusterCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})
}

func TestMockResourceNsxtPolicyVirtualNetworkServiceClusterRead(t *testing.T) {
	mockVNS, _, cleanup := setupVNSTest(t)
	defer cleanup()

	res := resourceNsxtPolicyVirtualNetworkServiceCluster()
	newData := func(t *testing.T) *schema.ResourceData {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         vnsClusterSitePath,
			"enforcement_point": vnsClusterEPID,
		})
		d.SetId(vnsClusterID)
		return d
	}

	t.Run("Read_success_all_fields", func(t *testing.T) {
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, vnsClusterID).Return(vnsClusterFullModel(), nil)

		d := newData(t)
		require.NoError(t, resourceNsxtPolicyVirtualNetworkServiceClusterRead(d, newGoMockProviderClient()))
		assert.Equal(t, vnsClusterName, d.Get("display_name"))
		assert.Equal(t, vnsClusterPath, d.Get("path"))
		assert.Equal(t, vnsFormFactor, d.Get("node_form_factor"))
		assert.Equal(t, int(vnsNodes), d.Get("number_of_nodes"))
		assert.Equal(t, vnsBfdPath, d.Get("bfd_profile_path"))
		assert.Equal(t, true, d.Get("default_project_service_allocation"))
		cfgs := d.Get("vsphere_deployment_config").([]interface{})
		require.Len(t, cfgs, 1)
		cfg := cfgs[0].(map[string]interface{})
		assert.Equal(t, vnsComputeManager, cfg["compute_manager_id"])
		assert.Equal(t, 2, cfg["cluster_ids"].(*schema.Set).Len())
		assert.Equal(t, 1, cfg["fault_domain_ids"].(*schema.Set).Len())
	})

	t.Run("Read_success_nil_optional_fields", func(t *testing.T) {
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, vnsClusterID).Return(model.VirtualNetworkServiceCluster{
			DisplayName: &vnsClusterName,
			Path:        &vnsClusterPath,
			VsphereDeploymentConfig: []model.VirtualNetworkServiceDeploymentConfig{
				{ClusterIds: []string{"domain-c1"}},
			},
		}, nil)

		d := newData(t)
		require.NoError(t, resourceNsxtPolicyVirtualNetworkServiceClusterRead(d, newGoMockProviderClient()))
		assert.Equal(t, vnsClusterName, d.Get("display_name"))
		assert.Equal(t, 0, d.Get("number_of_nodes"))
		assert.Equal(t, false, d.Get("default_project_service_allocation"))
	})

	t.Run("Read_fails_without_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		err := resourceNsxtPolicyVirtualNetworkServiceClusterRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Resource ID")
	})

	t.Run("Read_NotFound_clears_id", func(t *testing.T) {
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, vnsClusterID).Return(model.VirtualNetworkServiceCluster{}, vapiErrors.NotFound{})

		d := newData(t)
		require.NoError(t, resourceNsxtPolicyVirtualNetworkServiceClusterRead(d, newGoMockProviderClient()))
		assert.Empty(t, d.Id())
	})

	t.Run("Read_fails_on_other_error", func(t *testing.T) {
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, vnsClusterID).Return(model.VirtualNetworkServiceCluster{}, vapiErrors.ServiceUnavailable{})

		d := newData(t)
		require.Error(t, resourceNsxtPolicyVirtualNetworkServiceClusterRead(d, newGoMockProviderClient()))
		assert.Equal(t, vnsClusterID, d.Id())
	})
}

func TestMockResourceNsxtPolicyVirtualNetworkServiceClusterUpdate(t *testing.T) {
	mockVNS, _, cleanup := setupVNSTest(t)
	defer cleanup()

	res := resourceNsxtPolicyVirtualNetworkServiceCluster()
	newData := func(t *testing.T) *schema.ResourceData {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":      "vns-cluster-updated",
			"site_path":         vnsClusterSitePath,
			"enforcement_point": vnsClusterEPID,
			"node_form_factor":  model.VirtualNetworkServiceCluster_NODE_FORM_FACTOR_LARGE,
		})
		d.SetId(vnsClusterID)
		d.Set("revision", 7)
		return d
	}

	t.Run("Update_success", func(t *testing.T) {
		updatedName := "vns-cluster-updated"
		large := model.VirtualNetworkServiceCluster_NODE_FORM_FACTOR_LARGE
		mockVNS.EXPECT().Update(vnsClusterSiteID, vnsClusterEPID, vnsClusterID, gomock.Any()).DoAndReturn(
			func(_, _, _ string, obj model.VirtualNetworkServiceCluster) (model.VirtualNetworkServiceCluster, error) {
				require.NotNil(t, obj.Revision)
				assert.Equal(t, int64(7), *obj.Revision)
				assert.Equal(t, large, *obj.NodeFormFactor)
				return model.VirtualNetworkServiceCluster{}, nil
			})
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, vnsClusterID).Return(model.VirtualNetworkServiceCluster{
			DisplayName:    &updatedName,
			NodeFormFactor: &large,
		}, nil)

		d := newData(t)
		require.NoError(t, resourceNsxtPolicyVirtualNetworkServiceClusterUpdate(d, newGoMockProviderClient()))
		assert.Equal(t, updatedName, d.Get("display_name"))
		assert.Equal(t, large, d.Get("node_form_factor"))
	})

	t.Run("Update_fails_when_Update_returns_error", func(t *testing.T) {
		mockVNS.EXPECT().Update(vnsClusterSiteID, vnsClusterEPID, vnsClusterID, gomock.Any()).Return(model.VirtualNetworkServiceCluster{}, errors.New("update failed"))

		require.Error(t, resourceNsxtPolicyVirtualNetworkServiceClusterUpdate(newData(t), newGoMockProviderClient()))
	})

	t.Run("Update_fails_without_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		require.Error(t, resourceNsxtPolicyVirtualNetworkServiceClusterUpdate(d, newGoMockProviderClient()))
	})

}

func TestMockResourceNsxtPolicyVirtualNetworkServiceClusterDelete(t *testing.T) {
	mockVNS, _, cleanup := setupVNSTest(t)
	defer cleanup()

	res := resourceNsxtPolicyVirtualNetworkServiceCluster()
	newData := func(t *testing.T) *schema.ResourceData {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         vnsClusterSitePath,
			"enforcement_point": vnsClusterEPID,
		})
		d.SetId(vnsClusterID)
		return d
	}

	t.Run("Delete_success", func(t *testing.T) {
		mockVNS.EXPECT().Delete(vnsClusterSiteID, vnsClusterEPID, vnsClusterID).Return(nil)
		require.NoError(t, resourceNsxtPolicyVirtualNetworkServiceClusterDelete(newData(t), newGoMockProviderClient()))
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockVNS.EXPECT().Delete(vnsClusterSiteID, vnsClusterEPID, vnsClusterID).Return(errors.New("delete failed"))
		require.Error(t, resourceNsxtPolicyVirtualNetworkServiceClusterDelete(newData(t), newGoMockProviderClient()))
	})

	t.Run("Delete_fails_without_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		require.Error(t, resourceNsxtPolicyVirtualNetworkServiceClusterDelete(d, newGoMockProviderClient()))
	})
}

func TestUnitNsxt_resourceNsxtPolicyVirtualNetworkServiceClusterSchema(t *testing.T) {
	res := resourceNsxtPolicyVirtualNetworkServiceCluster()
	require.NoError(t, res.InternalValidate(nil, true))

	validate := res.Schema["node_form_factor"].ValidateFunc
	_, errs := validate("SMALL", "node_form_factor")
	assert.Empty(t, errs)
	_, errs = validate("HUGE", "node_form_factor")
	assert.NotEmpty(t, errs)

	validate = res.Schema["number_of_nodes"].ValidateFunc
	_, errs = validate(0, "number_of_nodes")
	assert.NotEmpty(t, errs)
	_, errs = validate(2, "number_of_nodes")
	assert.Empty(t, errs)
}

func TestUnitNsxt_resourceNsxtPolicyVirtualNetworkServiceClusterImporter(t *testing.T) {
	util.NsxVersion = "9.2.0"
	defer func() { util.NsxVersion = "" }()

	res := resourceNsxtPolicyVirtualNetworkServiceCluster()

	t.Run("version too low fails", func(t *testing.T) {
		util.NsxVersion = "9.1.1"
		defer func() { util.NsxVersion = "9.2.0" }()

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vnsClusterPath)

		_, err := res.Importer.State(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})

	t.Run("version supported imports via wrapped importer", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vnsClusterPath)

		out, err := res.Importer.State(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, vnsClusterID, d.Id())
	})

	t.Run("valid path sets enforcement_point and site_path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(vnsClusterPath)

		out, err := resourceNsxtPolicyVirtualNetworkServiceClusterImporter(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, vnsClusterID, d.Id())
		assert.Equal(t, vnsClusterEPID, d.Get("enforcement_point"))
		assert.Equal(t, vnsClusterSitePath, d.Get("site_path"))
	})

	t.Run("non-policy-path id fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("plain-id")

		_, err := resourceNsxtPolicyVirtualNetworkServiceClusterImporter(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("path of other resource type fails", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("/infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/c1")

		_, err := resourceNsxtPolicyVirtualNetworkServiceClusterImporter(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
