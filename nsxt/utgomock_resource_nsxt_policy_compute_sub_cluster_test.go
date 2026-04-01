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
	subClusterID       = "sub-cluster-1"
	subClusterSitePath = "/infra/sites/default"
	subClusterSiteID   = "default"
	subClusterEPID     = "default"
	subClusterCollID   = "collection-1"
	subClusterName     = "sub-cluster-fooname"
	subClusterRevision = int64(1)
	subClusterPath     = "/infra/sites/default/enforcement-points/default/sub-clusters/sub-cluster-1"
)

func TestMockResourceNsxtPolicyComputeSubClusterCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubClustersSDK := epmocks.NewMockSubClustersClient(ctrl)
	mockSitesSDK := inframocks.NewMockSitesClient(ctrl)

	subClusterWrapper := &enforcementpoints.SubClusterClientContext{
		Client:     mockSubClustersSDK,
		ClientType: utl.Local,
	}
	siteWrapper := &cliinfra.SiteClientContext{
		Client:     mockSitesSDK,
		ClientType: utl.Local,
	}

	originalSubClusters := cliSubClustersClient
	originalSites := cliSitesClient
	defer func() {
		cliSubClustersClient = originalSubClusters
		cliSitesClient = originalSites
	}()
	cliSubClustersClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.SubClusterClientContext {
		return subClusterWrapper
	}
	cliSitesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SiteClientContext {
		return siteWrapper
	}

	res := resourceNsxtPolicyComputeSubCluster()

	t.Run("Create_success", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(subClusterSiteID).Return(model.Site{}, nil)
		mockSubClustersSDK.EXPECT().Get(subClusterSiteID, subClusterEPID, gomock.Any()).Return(model.SubCluster{}, vapiErrors.NotFound{})
		mockSubClustersSDK.EXPECT().Patch(subClusterSiteID, subClusterEPID, gomock.Any(), gomock.Any()).Return(model.SubCluster{}, nil)
		mockSubClustersSDK.EXPECT().Get(subClusterSiteID, subClusterEPID, gomock.Any()).Return(model.SubCluster{
			DisplayName:         &subClusterName,
			ComputeCollectionId: &subClusterCollID,
			Path:                &subClusterPath,
			Revision:            &subClusterRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":          subClusterName,
			"site_path":             subClusterSitePath,
			"enforcement_point":     subClusterEPID,
			"compute_collection_id": subClusterCollID,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyComputeSubClusterCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, subClusterName, d.Get("display_name"))
	})

	t.Run("Create_fails_when_already_exists", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(subClusterSiteID).Return(model.Site{}, nil)
		mockSubClustersSDK.EXPECT().Get(subClusterSiteID, subClusterEPID, subClusterID).Return(model.SubCluster{}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"nsx_id":                subClusterID,
			"site_path":             subClusterSitePath,
			"enforcement_point":     subClusterEPID,
			"compute_collection_id": subClusterCollID,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyComputeSubClusterCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create_fails_when_site_path_invalid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":             "/infra/no-sites/default",
			"enforcement_point":     subClusterEPID,
			"compute_collection_id": subClusterCollID,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyComputeSubClusterCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Site ID")
	})

	t.Run("Create_fails_when_Patch_returns_error", func(t *testing.T) {
		mockSitesSDK.EXPECT().Get(subClusterSiteID).Return(model.Site{}, nil)
		mockSubClustersSDK.EXPECT().Get(subClusterSiteID, subClusterEPID, gomock.Any()).Return(model.SubCluster{}, vapiErrors.NotFound{})
		mockSubClustersSDK.EXPECT().Patch(subClusterSiteID, subClusterEPID, gomock.Any(), gomock.Any()).Return(model.SubCluster{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":             subClusterSitePath,
			"enforcement_point":     subClusterEPID,
			"compute_collection_id": subClusterCollID,
		})
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyComputeSubClusterCreate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyComputeSubClusterRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubClustersSDK := epmocks.NewMockSubClustersClient(ctrl)
	subClusterWrapper := &enforcementpoints.SubClusterClientContext{
		Client:     mockSubClustersSDK,
		ClientType: utl.Local,
	}

	originalSubClusters := cliSubClustersClient
	defer func() { cliSubClustersClient = originalSubClusters }()
	cliSubClustersClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.SubClusterClientContext {
		return subClusterWrapper
	}

	res := resourceNsxtPolicyComputeSubCluster()

	t.Run("Read_success", func(t *testing.T) {
		mockSubClustersSDK.EXPECT().Get(subClusterSiteID, subClusterEPID, subClusterID).Return(model.SubCluster{
			DisplayName:         &subClusterName,
			ComputeCollectionId: &subClusterCollID,
			Path:                &subClusterPath,
			Revision:            &subClusterRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         subClusterSitePath,
			"enforcement_point": subClusterEPID,
		})
		d.SetId(subClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyComputeSubClusterRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, subClusterName, d.Get("display_name"))
		assert.Equal(t, subClusterCollID, d.Get("compute_collection_id"))
	})

	t.Run("Read_fails_when_Get_returns_error", func(t *testing.T) {
		mockSubClustersSDK.EXPECT().Get(subClusterSiteID, subClusterEPID, subClusterID).Return(model.SubCluster{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         subClusterSitePath,
			"enforcement_point": subClusterEPID,
		})
		d.SetId(subClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyComputeSubClusterRead(d, m)
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyComputeSubClusterUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubClustersSDK := epmocks.NewMockSubClustersClient(ctrl)
	subClusterWrapper := &enforcementpoints.SubClusterClientContext{
		Client:     mockSubClustersSDK,
		ClientType: utl.Local,
	}

	originalSubClusters := cliSubClustersClient
	defer func() { cliSubClustersClient = originalSubClusters }()
	cliSubClustersClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.SubClusterClientContext {
		return subClusterWrapper
	}

	res := resourceNsxtPolicyComputeSubCluster()

	t.Run("Update_success", func(t *testing.T) {
		mockSubClustersSDK.EXPECT().Update(subClusterSiteID, subClusterEPID, subClusterID, gomock.Any()).Return(model.SubCluster{}, nil)
		mockSubClustersSDK.EXPECT().Get(subClusterSiteID, subClusterEPID, subClusterID).Return(model.SubCluster{
			DisplayName:         &subClusterName,
			ComputeCollectionId: &subClusterCollID,
			Path:                &subClusterPath,
			Revision:            &subClusterRevision,
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":          subClusterName,
			"site_path":             subClusterSitePath,
			"enforcement_point":     subClusterEPID,
			"compute_collection_id": subClusterCollID,
		})
		d.SetId(subClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyComputeSubClusterUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update_fails_when_Update_returns_error", func(t *testing.T) {
		mockSubClustersSDK.EXPECT().Update(subClusterSiteID, subClusterEPID, subClusterID, gomock.Any()).Return(model.SubCluster{}, vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":             subClusterSitePath,
			"enforcement_point":     subClusterEPID,
			"compute_collection_id": subClusterCollID,
		})
		d.SetId(subClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyComputeSubClusterUpdate(d, m)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyComputeSubClusterDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubClustersSDK := epmocks.NewMockSubClustersClient(ctrl)
	subClusterWrapper := &enforcementpoints.SubClusterClientContext{
		Client:     mockSubClustersSDK,
		ClientType: utl.Local,
	}

	originalSubClusters := cliSubClustersClient
	defer func() { cliSubClustersClient = originalSubClusters }()
	cliSubClustersClient = func(sessionContext utl.SessionContext, connector client.Connector) *enforcementpoints.SubClusterClientContext {
		return subClusterWrapper
	}

	res := resourceNsxtPolicyComputeSubCluster()

	t.Run("Delete_success", func(t *testing.T) {
		mockSubClustersSDK.EXPECT().Delete(subClusterSiteID, subClusterEPID, subClusterID).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         subClusterSitePath,
			"enforcement_point": subClusterEPID,
		})
		d.SetId(subClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyComputeSubClusterDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete_fails_when_Delete_returns_error", func(t *testing.T) {
		mockSubClustersSDK.EXPECT().Delete(subClusterSiteID, subClusterEPID, subClusterID).Return(vapiErrors.InternalServerError{})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"site_path":         subClusterSitePath,
			"enforcement_point": subClusterEPID,
		})
		d.SetId(subClusterID)
		m := newGoMockProviderClient()
		err := resourceNsxtPolicyComputeSubClusterDelete(d, m)
		require.Error(t, err)
	})
}
