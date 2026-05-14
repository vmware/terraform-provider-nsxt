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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestMockDataSourceNsxtPolicyVirtualNetworkAppliance(t *testing.T) {
	util.NsxVersion = "9.1.1"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVNA := setupVNACRUDMock(ctrl)
	restore := setupVNACRUDClientOverride(mockVNA)
	defer restore()

	ds := dataSourceNsxtPolicyVirtualNetworkAppliance()
	oneCount := int64(1)

	t.Run("Read_by_id_success", func(t *testing.T) {
		returnSV := vnaStructValue(model.VirtualNetworkAppliance{
			Id:          &vnaID,
			DisplayName: &vnaName,
			Path:        &vnaPath,
			Hostname:    &vnaHostname,
		})
		mockVNA.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID).Return(returnSV, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":           vnaID,
			"cluster_path": vnaClusterPath,
		})
		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyVirtualNetworkApplianceRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, vnaID, d.Id())
		assert.Equal(t, vnaName, d.Get("display_name"))
		assert.Equal(t, vnaHostname, d.Get("hostname"))
	})

	t.Run("Read_by_id_not_found", func(t *testing.T) {
		mockVNA.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID, vnaID).Return(nil, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":           vnaID,
			"cluster_path": vnaClusterPath,
		})
		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyVirtualNetworkApplianceRead(d, m)
		require.Error(t, err)
	})

	t.Run("Read_by_display_name_success", func(t *testing.T) {
		returnSV := vnaStructValue(model.VirtualNetworkAppliance{
			Id:          &vnaID,
			DisplayName: &vnaName,
			Path:        &vnaPath,
		})
		mockVNA.EXPECT().List(vnaClusterSiteID, vnaClusterEPID, vnaClusterID,
			nil, nil, nil, nil, nil, nil, nil, nil, nil,
		).Return(model.VirtualNetworkApplianceListResult{
			ResultCount: &oneCount,
			Results:     []model.VirtualNetworkAppliance{{Id: &vnaID, DisplayName: &vnaName, Path: &vnaPath}},
		}, nil)
		_ = returnSV

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": vnaName,
			"cluster_path": vnaClusterPath,
		})
		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyVirtualNetworkApplianceRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, vnaID, d.Id())
		assert.Equal(t, vnaPath, d.Get("path"))
	})

	t.Run("Read_by_display_name_multiple_found", func(t *testing.T) {
		otherID := vnaID + "-2"
		mockVNA.EXPECT().List(vnaClusterSiteID, vnaClusterEPID, vnaClusterID,
			nil, nil, nil, nil, nil, nil, nil, nil, nil,
		).Return(model.VirtualNetworkApplianceListResult{
			ResultCount: func() *int64 { v := int64(2); return &v }(),
			Results: []model.VirtualNetworkAppliance{
				{Id: &vnaID, DisplayName: &vnaName},
				{Id: &otherID, DisplayName: &vnaName},
			},
		}, nil)

		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": vnaName,
			"cluster_path": vnaClusterPath,
		})
		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyVirtualNetworkApplianceRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("Read_fails_when_no_id_or_name", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"cluster_path": vnaClusterPath,
		})
		m := newGoMockProviderClient()
		err := dataSourceNsxtPolicyVirtualNetworkApplianceRead(d, m)
		require.Error(t, err)
	})
}

func TestMockDataSourceNsxtPolicyVirtualNetworkAppliance_epmocks(t *testing.T) {
	util.NsxVersion = "9.1.1"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	_ = epmocks.NewMockVirtualNetworkAppliancesInClusterClient(ctrl)
}
