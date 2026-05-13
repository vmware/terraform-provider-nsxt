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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	enforcementpoints "github.com/vmware/terraform-provider-nsxt/api/infra/sites/enforcement_points"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	epmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/sites/enforcement_points"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func setupVNAClusterDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*epmocks.MockVirtualNetworkApplianceClustersClient, func()) {
	t.Helper()
	mockSDK := epmocks.NewMockVirtualNetworkApplianceClustersClient(ctrl)
	wrapper := &enforcementpoints.VirtualNetworkApplianceClusterClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliVNAClustersClient
	cliVNAClustersClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *enforcementpoints.VirtualNetworkApplianceClusterClientContext {
		return wrapper
	}
	return mockSDK, func() { cliVNAClustersClient = orig }
}

func TestMockDataSourceNsxtPolicyVirtualNetworkApplianceClusterRead(t *testing.T) {
	util.NsxVersion = "9.1.1"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSDK, restore := setupVNAClusterDataSourceMock(t, ctrl)
	defer restore()

	m := newGoMockProviderClient()
	formFactor := model.VirtualNetworkApplianceCluster_APPLIANCE_FORM_FACTOR_MEDIUM
	svcType := model.VirtualNetworkApplianceCluster_SERVICE_TYPE_VPC_SERVICES
	clusterModel := model.VirtualNetworkApplianceCluster{
		Id:                  &vnaClusterID,
		DisplayName:         &vnaClusterName,
		Path:                &vnaClusterPath,
		Revision:            &vnaClusterRevision,
		ApplianceFormFactor: &formFactor,
		ServiceType:         &svcType,
	}
	resultCount := int64(1)

	t.Run("by_id_success", func(t *testing.T) {
		mockSDK.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID).Return(clusterModel, nil)

		ds := dataSourceNsxtPolicyVirtualNetworkApplianceCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":                vnaClusterID,
			"site_path":         vnaClusterSitePath,
			"enforcement_point": vnaClusterEPID,
		})

		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, vnaClusterID, d.Id())
		assert.Equal(t, vnaClusterName, d.Get("display_name"))
		assert.Equal(t, formFactor, d.Get("appliance_form_factor"))
		assert.Equal(t, svcType, d.Get("service_type"))
	})

	t.Run("by_id_not_found", func(t *testing.T) {
		mockSDK.EXPECT().Get(vnaClusterSiteID, vnaClusterEPID, vnaClusterID).Return(model.VirtualNetworkApplianceCluster{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyVirtualNetworkApplianceCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":                vnaClusterID,
			"site_path":         vnaClusterSitePath,
			"enforcement_point": vnaClusterEPID,
		})

		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRead(d, m)
		require.Error(t, err)
	})

	t.Run("by_display_name_success", func(t *testing.T) {
		mockSDK.EXPECT().List(vnaClusterSiteID, vnaClusterEPID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.VirtualNetworkApplianceClusterListResult{
				Results:     []model.VirtualNetworkApplianceCluster{clusterModel},
				ResultCount: &resultCount,
			}, nil)

		ds := dataSourceNsxtPolicyVirtualNetworkApplianceCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name":      vnaClusterName,
			"site_path":         vnaClusterSitePath,
			"enforcement_point": vnaClusterEPID,
		})

		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, vnaClusterID, d.Id())
		assert.Equal(t, vnaClusterName, d.Get("display_name"))
	})

	t.Run("by_display_name_not_found", func(t *testing.T) {
		zeroCount := int64(0)
		mockSDK.EXPECT().List(vnaClusterSiteID, vnaClusterEPID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.VirtualNetworkApplianceClusterListResult{
				Results:     []model.VirtualNetworkApplianceCluster{},
				ResultCount: &zeroCount,
			}, nil)

		ds := dataSourceNsxtPolicyVirtualNetworkApplianceCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name":      "nonexistent",
			"site_path":         vnaClusterSitePath,
			"enforcement_point": vnaClusterEPID,
		})

		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("missing_id_and_name", func(t *testing.T) {
		ds := dataSourceNsxtPolicyVirtualNetworkApplianceCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"site_path":         vnaClusterSitePath,
			"enforcement_point": vnaClusterEPID,
		})

		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ID or name")
	})

	t.Run("invalid_site_path", func(t *testing.T) {
		ds := dataSourceNsxtPolicyVirtualNetworkApplianceCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":                vnaClusterID,
			"site_path":         "/infra/no-sites/default",
			"enforcement_point": vnaClusterEPID,
		})

		err := dataSourceNsxtPolicyVirtualNetworkApplianceClusterRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Site ID")
	})
}
