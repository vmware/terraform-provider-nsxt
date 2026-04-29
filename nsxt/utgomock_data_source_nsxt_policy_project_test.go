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
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

func TestMockDataSourceNsxtPolicyProjectRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupProjectMock(t, ctrl)
	defer restore()

	t.Run("by id success", func(t *testing.T) {
		p := projectAPIResponse()
		sitePath := "/infra/sites/default"
		p.SiteInfos = []nsxModel.SiteInfo{
			{EdgeClusterPaths: []string{"/infra/edge-clusters/ec1"}, SitePath: &sitePath},
		}
		p.Tier0s = []string{"/infra/tier-0s/t0"}
		p.ExternalIpv4Blocks = []string{"/infra/ip-blocks/block1"}

		mockSDK.EXPECT().Get(utl.DefaultOrgID, projectID, gomock.Any()).Return(p, nil)

		ds := dataSourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": projectID,
		})

		err := dataSourceNsxtPolicyProjectRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, projectID, d.Id())
		assert.Equal(t, projectDisplayName, d.Get("display_name"))
		si := d.Get("site_info").([]interface{})
		require.Len(t, si, 1)
	})

	t.Run("by id API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(utl.DefaultOrgID, projectID, gomock.Any()).Return(nsxModel.Project{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": projectID,
		})

		err := dataSourceNsxtPolicyProjectRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("missing id and display_name", func(t *testing.T) {
		ds := dataSourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyProjectRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ID or name")
	})

	t.Run("by display_name single exact match", func(t *testing.T) {
		mockSDK.EXPECT().List(utl.DefaultOrgID, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ProjectListResult{
			Results: []nsxModel.Project{projectAPIResponse()},
		}, nil)

		ds := dataSourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": projectDisplayName,
		})

		err := dataSourceNsxtPolicyProjectRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, projectID, d.Id())
	})

	t.Run("by display_name not found", func(t *testing.T) {
		mockSDK.EXPECT().List(utl.DefaultOrgID, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ProjectListResult{
			Results: []nsxModel.Project{},
		}, nil)

		ds := dataSourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "missing-name",
		})

		err := dataSourceNsxtPolicyProjectRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by display_name multiple exact matches", func(t *testing.T) {
		mockSDK.EXPECT().List(utl.DefaultOrgID, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ProjectListResult{
			Results: []nsxModel.Project{projectAPIResponse(), projectAPIResponse()},
		}, nil)

		ds := dataSourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": projectDisplayName,
		})

		err := dataSourceNsxtPolicyProjectRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("list API error", func(t *testing.T) {
		mockSDK.EXPECT().List(utl.DefaultOrgID, nil, nil, nil, nil, nil, nil, nil).Return(nsxModel.ProjectListResult{}, vapiErrors.ServiceUnavailable{})

		ds := dataSourceNsxtPolicyProject()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": projectDisplayName,
		})

		err := dataSourceNsxtPolicyProjectRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
