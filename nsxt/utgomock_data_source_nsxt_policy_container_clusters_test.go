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
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"
)

func TestMockDataSourceNsxtPolicyContainerClustersRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupClusterControlPlanesMock(t, ctrl)
	defer restore()

	t.Run("success", func(t *testing.T) {
		mockSDK.EXPECT().List(defaultSite, defaultEnforcementPoint, nil, nil, nil, nil, nil, gomock.Not(gomock.Nil()), nil).Return(nsxModel.ClusterControlPlaneListResult{
			Results:     []nsxModel.ClusterControlPlane{containerClusterControlPlaneAPIResponse()},
			ResultCount: int64Ptr(1),
		}, nil)

		ds := dataSourceNsxtPolicyContainerClusters()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyContainerClustersRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		items := d.Get("items").(map[string]interface{})
		assert.Equal(t, ccpPath, items[ccpDisplayName])
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List(defaultSite, defaultEnforcementPoint, nil, nil, nil, nil, nil, gomock.Not(gomock.Nil()), nil).Return(nsxModel.ClusterControlPlaneListResult{}, errors.New("list failed"))

		ds := dataSourceNsxtPolicyContainerClusters()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyContainerClustersRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
