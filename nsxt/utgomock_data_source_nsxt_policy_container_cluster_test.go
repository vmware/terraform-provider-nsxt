//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/sites/enforcement_points/cluster_control_planes/ClusterControlPlanesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points/ClusterControlPlanesClient.go ClusterControlPlanesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	enforcementpoints "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	ccpmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points/cluster_control_planes"
)

var (
	ccpID          = "ccp-1"
	ccpDisplayName = "Test Container Cluster"
	ccpDescription = "Test container cluster control plane"
	ccpPath        = "/infra/sites/default/enforcement-points/default/cluster-control-planes/ccp-1"
)

func containerClusterControlPlaneAPIResponse() nsxModel.ClusterControlPlane {
	return nsxModel.ClusterControlPlane{
		Id:          &ccpID,
		DisplayName: &ccpDisplayName,
		Description: &ccpDescription,
		Path:        &ccpPath,
	}
}

func setupClusterControlPlanesMock(t *testing.T, ctrl *gomock.Controller) (*ccpmocks.MockClusterControlPlanesClient, func()) {
	mockSDK := ccpmocks.NewMockClusterControlPlanesClient(ctrl)
	mockWrapper := &enforcementpoints.ClusterControlPlanesClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliClusterControlPlanesClient
	cliClusterControlPlanesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *enforcementpoints.ClusterControlPlanesClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliClusterControlPlanesClient = original }
}

func TestMockDataSourceNsxtPolicyContainerClusterRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupClusterControlPlanesMock(t, ctrl)
	defer restore()

	t.Run("by id success", func(t *testing.T) {
		mockSDK.EXPECT().Get(defaultSite, defaultEnforcementPoint, ccpID).Return(containerClusterControlPlaneAPIResponse(), nil)

		ds := dataSourceNsxtPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": ccpID,
		})

		err := dataSourceNsxtPolicyContainerClusterRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ccpID, d.Id())
		assert.Equal(t, ccpDisplayName, d.Get("display_name"))
		assert.Equal(t, ccpPath, d.Get("path"))
	})

	t.Run("by id API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(defaultSite, defaultEnforcementPoint, ccpID).Return(nsxModel.ClusterControlPlane{}, errors.New("get failed"))

		ds := dataSourceNsxtPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": ccpID,
		})

		err := dataSourceNsxtPolicyContainerClusterRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("by display_name single exact match", func(t *testing.T) {
		mockSDK.EXPECT().List(defaultSite, defaultEnforcementPoint, nil, nil, nil, nil, nil, gomock.Not(gomock.Nil()), nil).Return(nsxModel.ClusterControlPlaneListResult{
			Results:     []nsxModel.ClusterControlPlane{containerClusterControlPlaneAPIResponse()},
			ResultCount: int64Ptr(1),
		}, nil)

		ds := dataSourceNsxtPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": ccpDisplayName,
		})

		err := dataSourceNsxtPolicyContainerClusterRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ccpID, d.Id())
	})

	t.Run("by display_name multiple exact matches", func(t *testing.T) {
		mockSDK.EXPECT().List(defaultSite, defaultEnforcementPoint, nil, nil, nil, nil, nil, gomock.Not(gomock.Nil()), nil).Return(nsxModel.ClusterControlPlaneListResult{
			Results:     []nsxModel.ClusterControlPlane{containerClusterControlPlaneAPIResponse(), containerClusterControlPlaneAPIResponse()},
			ResultCount: int64Ptr(2),
		}, nil)

		ds := dataSourceNsxtPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": ccpDisplayName,
		})

		err := dataSourceNsxtPolicyContainerClusterRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("no id or name provided", func(t *testing.T) {
		ds := dataSourceNsxtPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyContainerClusterRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "obtaining ClusterControlPlane")
	})

	t.Run("no match from list", func(t *testing.T) {
		mockSDK.EXPECT().List(defaultSite, defaultEnforcementPoint, nil, nil, nil, nil, nil, gomock.Not(gomock.Nil()), nil).Return(nsxModel.ClusterControlPlaneListResult{
			Results:     []nsxModel.ClusterControlPlane{},
			ResultCount: int64Ptr(0),
		}, nil)

		ds := dataSourceNsxtPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		err := dataSourceNsxtPolicyContainerClusterRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List(defaultSite, defaultEnforcementPoint, nil, nil, nil, nil, nil, gomock.Not(gomock.Nil()), nil).Return(nsxModel.ClusterControlPlaneListResult{}, errors.New("list failed"))

		ds := dataSourceNsxtPolicyContainerCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": ccpDisplayName,
		})

		err := dataSourceNsxtPolicyContainerClusterRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
