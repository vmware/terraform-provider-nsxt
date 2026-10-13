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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyVirtualNetworkServiceClusterRead(t *testing.T) {
	mockVNS, _, cleanup := setupVNSTest(t)
	defer cleanup()

	m := newGoMockProviderClient()
	ds := dataSourceNsxtPolicyVirtualNetworkServiceCluster()
	newData := func(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
		raw["site_path"] = vnsClusterSitePath
		raw["enforcement_point"] = vnsClusterEPID
		return schema.TestResourceDataRaw(t, ds.Schema, raw)
	}
	named := func(id, name string) model.VirtualNetworkServiceCluster {
		return model.VirtualNetworkServiceCluster{Id: &id, DisplayName: &name}
	}
	resultCount := func(n int) *int64 { c := int64(n); return &c }

	t.Run("by_id_success", func(t *testing.T) {
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, vnsClusterID).Return(vnsClusterFullModel(), nil)

		d := newData(t, map[string]interface{}{"id": vnsClusterID})
		require.NoError(t, dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m))
		assert.Equal(t, vnsClusterID, d.Id())
		assert.Equal(t, vnsClusterName, d.Get("display_name"))
		assert.Equal(t, vnsClusterPath, d.Get("path"))
		assert.Equal(t, vnsFormFactor, d.Get("node_form_factor"))
		assert.Equal(t, int(vnsNodes), d.Get("number_of_nodes"))
		assert.Equal(t, vnsBfdPath, d.Get("bfd_profile_path"))
		assert.Equal(t, true, d.Get("default_project_service_allocation"))
		cfgs := d.Get("vsphere_deployment_config").([]interface{})
		require.Len(t, cfgs, 1)
		assert.Equal(t, vnsComputeManager, cfgs[0].(map[string]interface{})["compute_manager_id"])
	})

	t.Run("by_id_without_allocation_config", func(t *testing.T) {
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, vnsClusterID).Return(named(vnsClusterID, vnsClusterName), nil)

		d := newData(t, map[string]interface{}{"id": vnsClusterID})
		require.NoError(t, dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m))
		assert.Equal(t, false, d.Get("default_project_service_allocation"))
		assert.Empty(t, d.Get("vsphere_deployment_config"))
	})

	t.Run("by_id_not_found", func(t *testing.T) {
		mockVNS.EXPECT().Get(vnsClusterSiteID, vnsClusterEPID, vnsClusterID).Return(model.VirtualNetworkServiceCluster{}, vapiErrors.NotFound{})

		d := newData(t, map[string]interface{}{"id": vnsClusterID})
		require.Error(t, dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m))
	})

	t.Run("by_display_name_exact_match", func(t *testing.T) {
		mockVNS.EXPECT().List(vnsClusterSiteID, vnsClusterEPID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.VirtualNetworkServiceClusterListResult{
				Results:     []model.VirtualNetworkServiceCluster{named("other", vnsClusterName+"-2"), named(vnsClusterID, vnsClusterName)},
				ResultCount: resultCount(2),
			}, nil)

		d := newData(t, map[string]interface{}{"display_name": vnsClusterName})
		require.NoError(t, dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m))
		assert.Equal(t, vnsClusterID, d.Id())
	})

	t.Run("by_display_name_prefix_match", func(t *testing.T) {
		mockVNS.EXPECT().List(vnsClusterSiteID, vnsClusterEPID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.VirtualNetworkServiceClusterListResult{
				Results:     []model.VirtualNetworkServiceCluster{named(vnsClusterID, vnsClusterName), {}},
				ResultCount: resultCount(2),
			}, nil)

		d := newData(t, map[string]interface{}{"display_name": "vns-cluster"})
		require.NoError(t, dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m))
		assert.Equal(t, vnsClusterID, d.Id())
	})

	t.Run("by_display_name_multiple_exact", func(t *testing.T) {
		mockVNS.EXPECT().List(vnsClusterSiteID, vnsClusterEPID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.VirtualNetworkServiceClusterListResult{
				Results:     []model.VirtualNetworkServiceCluster{named("a", vnsClusterName), named("b", vnsClusterName)},
				ResultCount: resultCount(2),
			}, nil)

		d := newData(t, map[string]interface{}{"display_name": vnsClusterName})
		err := dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("by_display_name_multiple_prefix", func(t *testing.T) {
		mockVNS.EXPECT().List(vnsClusterSiteID, vnsClusterEPID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.VirtualNetworkServiceClusterListResult{
				Results:     []model.VirtualNetworkServiceCluster{named("a", "vns-a"), named("b", "vns-b")},
				ResultCount: resultCount(2),
			}, nil)

		d := newData(t, map[string]interface{}{"display_name": "vns"})
		err := dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "starting with")
	})

	t.Run("by_display_name_not_found", func(t *testing.T) {
		mockVNS.EXPECT().List(vnsClusterSiteID, vnsClusterEPID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.VirtualNetworkServiceClusterListResult{
				Results:     []model.VirtualNetworkServiceCluster{named("a", "different")},
				ResultCount: resultCount(1),
			}, nil)

		d := newData(t, map[string]interface{}{"display_name": vnsClusterName})
		err := dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by_display_name_paginated", func(t *testing.T) {
		cursor := "next"
		gomock.InOrder(
			mockVNS.EXPECT().List(vnsClusterSiteID, vnsClusterEPID, nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
				model.VirtualNetworkServiceClusterListResult{
					Results:     []model.VirtualNetworkServiceCluster{named("a", "different")},
					ResultCount: resultCount(2),
					Cursor:      &cursor,
				}, nil),
			mockVNS.EXPECT().List(vnsClusterSiteID, vnsClusterEPID, &cursor, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
				model.VirtualNetworkServiceClusterListResult{
					Results:     []model.VirtualNetworkServiceCluster{named(vnsClusterID, vnsClusterName)},
					ResultCount: resultCount(2),
				}, nil),
		)

		d := newData(t, map[string]interface{}{"display_name": vnsClusterName})
		require.NoError(t, dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m))
		assert.Equal(t, vnsClusterID, d.Id())
	})

	t.Run("list_error", func(t *testing.T) {
		mockVNS.EXPECT().List(vnsClusterSiteID, vnsClusterEPID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.VirtualNetworkServiceClusterListResult{}, errors.New("list failed"))

		d := newData(t, map[string]interface{}{"display_name": vnsClusterName})
		require.Error(t, dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m))
	})

	t.Run("version_too_low", func(t *testing.T) {
		util.NsxVersion = "9.1.1"
		defer func() { util.NsxVersion = "9.2.0" }()

		d := newData(t, map[string]interface{}{"id": vnsClusterID})
		err := dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})

	t.Run("missing_id_and_name", func(t *testing.T) {
		d := newData(t, map[string]interface{}{})
		err := dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ID or name")
	})

	t.Run("invalid_site_path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":        vnsClusterID,
			"site_path": "/infra/no-sites/default",
		})
		err := dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Site ID")
	})

	t.Run("empty_enforcement_point_falls_back_to_provider", func(t *testing.T) {
		mockVNS.EXPECT().Get(vnsClusterSiteID, gomock.Any(), vnsClusterID).Return(named(vnsClusterID, vnsClusterName), nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":        vnsClusterID,
			"site_path": vnsClusterSitePath,
		})
		d.Set("enforcement_point", "")
		require.NoError(t, dataSourceNsxtPolicyVirtualNetworkServiceClusterRead(d, m))
	})
}
