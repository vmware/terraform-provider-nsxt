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
	edgeClusterID          = "edge-cluster-1"
	edgeClusterSitePath    = "/infra/sites/default"
	edgeClusterSiteID      = "default"
	edgeClusterEPID        = "default"
	edgeClusterName        = "edge-cluster-fooname"
	edgeClusterRevision    = int64(1)
	edgeClusterPath        = "/infra/sites/default/enforcement-points/default/edge-clusters/edge-cluster-1"
	edgeClusterProfilePath = "/infra/host-switch-profiles/edge-profile-1"
)

func TestMockResourceNsxtPolicyEdgeClusterCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEdgeClustersSDK := epmocks.NewMockEdgeClustersClient(ctrl)
	mockSitesSDK := inframocks.NewMockSitesClient(ctrl)

	edgeClusterWrapper := &enforcementpoints.PolicyEdgeClusterClientContext{
		Client:     mockEdgeClustersSDK,
		ClientType: utl.Local,
	}
	siteWrapper := &cliinfra.SiteClientContext{
		Client:     mockSitesSDK,
		ClientType: utl.Local,
	}

	originalEdgeClusters := cliPolicyEdgeClustersClient
	originalSites := cliSitesClient
	defer func() {
		cliPolicyEdgeClustersClient = originalEdgeClusters
		cliSitesClient = originalSites
	}()
	cliPolicyEdgeClustersClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyEdgeClusterClientContext {
		return edgeClusterWrapper
	}
	cliSitesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SiteClientContext {
		return siteWrapper
	}

	res := resourceNsxtPolicyEdgeCluster()

	t.Run("Create_success", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(edgeClusterSiteID).Return(model.Site{}, nil)
		mockEdgeClustersSDK.EXPECT().Get(edgeClusterSiteID, edgeClusterEPID, gomock.Any()).Return(model.PolicyEdgeCluster{}, vapiErrors.NotFound{})
		mockEdgeClustersSDK.EXPECT().Patch(edgeClusterSiteID, edgeClusterEPID, gomock.Any(), gomock.Any()).Return(nil)
		mockEdgeClustersSDK.EXPECT().Get(edgeClusterSiteID, edgeClusterEPID, gomock.Any()).Return(model.PolicyEdgeCluster{
			DisplayName:        &edgeClusterName,
			EdgeClusterProfile: &edgeClusterProfilePath,
			Path:               &edgeClusterPath,
			Revision:           &edgeClusterRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":              edgeClusterName,
			"site_path":                 edgeClusterSitePath,
			"enforcement_point":         edgeClusterEPID,
			"edge_cluster_profile_path": edgeClusterProfilePath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeClusterCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, edgeClusterName, d.Get("display_name"))
	})

	t.Run("Create_fails_when_already_exists", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(edgeClusterSiteID).Return(model.Site{}, nil)
		mockEdgeClustersSDK.EXPECT().Get(edgeClusterSiteID, edgeClusterEPID, edgeClusterID).Return(model.PolicyEdgeCluster{}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"nsx_id":                    edgeClusterID,
			"site_path":                 edgeClusterSitePath,
			"enforcement_point":         edgeClusterEPID,
			"edge_cluster_profile_path": edgeClusterProfilePath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeClusterCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create_fails_when_site_path_invalid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":                 "/infra/no-sites/default",
			"enforcement_point":         edgeClusterEPID,
			"edge_cluster_profile_path": edgeClusterProfilePath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeClusterCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Site ID")
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(edgeClusterSiteID).Return(model.Site{}, nil)
		mockEdgeClustersSDK.EXPECT().Get(edgeClusterSiteID, edgeClusterEPID, gomock.Any()).Return(model.PolicyEdgeCluster{}, vapiErrors.NotFound{})
		mockEdgeClustersSDK.EXPECT().Patch(edgeClusterSiteID, edgeClusterEPID, gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":                 edgeClusterSitePath,
			"enforcement_point":         edgeClusterEPID,
			"edge_cluster_profile_path": edgeClusterProfilePath,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeClusterCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyEdgeClusterRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEdgeClustersSDK := epmocks.NewMockEdgeClustersClient(ctrl)
	edgeClusterWrapper := &enforcementpoints.PolicyEdgeClusterClientContext{
		Client:     mockEdgeClustersSDK,
		ClientType: utl.Local,
	}

	originalEdgeClusters := cliPolicyEdgeClustersClient
	defer func() { cliPolicyEdgeClustersClient = originalEdgeClusters }()
	cliPolicyEdgeClustersClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyEdgeClusterClientContext {
		return edgeClusterWrapper
	}

	res := resourceNsxtPolicyEdgeCluster()

	t.Run("Read_success", func(t *testing.T) {
		mockEdgeClustersSDK.EXPECT().Get(edgeClusterSiteID, edgeClusterEPID, edgeClusterID).Return(model.PolicyEdgeCluster{
			DisplayName:        &edgeClusterName,
			EdgeClusterProfile: &edgeClusterProfilePath,
			Path:               &edgeClusterPath,
			Revision:           &edgeClusterRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         edgeClusterSitePath,
			"enforcement_point": edgeClusterEPID,
		})
		d.SetId(edgeClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeClusterRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, edgeClusterName, d.Get("display_name"))
		assert.Equal(t, edgeClusterProfilePath, d.Get("edge_cluster_profile_path"))
	})

	t.Run("Read_clears_id_when_not_found", func(t *testing.T) {
		mockEdgeClustersSDK.EXPECT().Get(edgeClusterSiteID, edgeClusterEPID, edgeClusterID).Return(model.PolicyEdgeCluster{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         edgeClusterSitePath,
			"enforcement_point": edgeClusterEPID,
		})
		d.SetId(edgeClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeClusterRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyEdgeClusterUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEdgeClustersSDK := epmocks.NewMockEdgeClustersClient(ctrl)
	edgeClusterWrapper := &enforcementpoints.PolicyEdgeClusterClientContext{
		Client:     mockEdgeClustersSDK,
		ClientType: utl.Local,
	}

	originalEdgeClusters := cliPolicyEdgeClustersClient
	defer func() { cliPolicyEdgeClustersClient = originalEdgeClusters }()
	cliPolicyEdgeClustersClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyEdgeClusterClientContext {
		return edgeClusterWrapper
	}

	res := resourceNsxtPolicyEdgeCluster()

	t.Run("Update_success", func(t *testing.T) {
		mockEdgeClustersSDK.EXPECT().Update(edgeClusterSiteID, edgeClusterEPID, edgeClusterID, gomock.Any()).Return(model.PolicyEdgeCluster{}, nil)
		mockEdgeClustersSDK.EXPECT().Get(edgeClusterSiteID, edgeClusterEPID, edgeClusterID).Return(model.PolicyEdgeCluster{
			DisplayName:        &edgeClusterName,
			EdgeClusterProfile: &edgeClusterProfilePath,
			Path:               &edgeClusterPath,
			Revision:           &edgeClusterRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":              edgeClusterName,
			"site_path":                 edgeClusterSitePath,
			"enforcement_point":         edgeClusterEPID,
			"edge_cluster_profile_path": edgeClusterProfilePath,
		})
		d.SetId(edgeClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeClusterUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_Update_returns_error", func(t *testing.T) {
		mockEdgeClustersSDK.EXPECT().Update(edgeClusterSiteID, edgeClusterEPID, edgeClusterID, gomock.Any()).Return(model.PolicyEdgeCluster{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":                 edgeClusterSitePath,
			"enforcement_point":         edgeClusterEPID,
			"edge_cluster_profile_path": edgeClusterProfilePath,
		})
		d.SetId(edgeClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeClusterUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyEdgeClusterDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEdgeClustersSDK := epmocks.NewMockEdgeClustersClient(ctrl)
	edgeClusterWrapper := &enforcementpoints.PolicyEdgeClusterClientContext{
		Client:     mockEdgeClustersSDK,
		ClientType: utl.Local,
	}

	originalEdgeClusters := cliPolicyEdgeClustersClient
	defer func() { cliPolicyEdgeClustersClient = originalEdgeClusters }()
	cliPolicyEdgeClustersClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.PolicyEdgeClusterClientContext {
		return edgeClusterWrapper
	}

	res := resourceNsxtPolicyEdgeCluster()

	t.Run("Delete_success", func(t *testing.T) {
		mockEdgeClustersSDK.EXPECT().Delete(edgeClusterSiteID, edgeClusterEPID, edgeClusterID, gomock.Any()).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         edgeClusterSitePath,
			"enforcement_point": edgeClusterEPID,
		})
		d.SetId(edgeClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeClusterDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockEdgeClustersSDK.EXPECT().Delete(edgeClusterSiteID, edgeClusterEPID, edgeClusterID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         edgeClusterSitePath,
			"enforcement_point": edgeClusterEPID,
		})
		d.SetId(edgeClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyEdgeClusterDelete(d, m)
		require.Error(t, err)
	})
}
