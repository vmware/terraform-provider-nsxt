//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/tier_0s/LocaleServicesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s/LocaleServicesClient.go LocaleServicesClient

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	tier0sAPI "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	tier0smocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
)

func gatewayLocaleServiceTestResource() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"locale_service":     getPolicyLocaleServiceSchema(false),
		"redistribution_set": {Type: schema.TypeBool, Optional: true, Computed: true},
	}}
}

func TestUnitNsxt_setLocaleServiceRedistributionConfig(t *testing.T) {
	t.Run("empty config is a no-op", func(t *testing.T) {
		serviceStruct := model.LocaleServices{}
		setLocaleServiceRedistributionConfig(nil, &serviceStruct)
		assert.Nil(t, serviceStruct.RouteRedistributionConfig)
	})

	t.Run("sets bgp and ospf flags", func(t *testing.T) {
		serviceStruct := model.LocaleServices{}
		setLocaleServiceRedistributionConfig([]interface{}{
			map[string]interface{}{"enabled": true, "ospf_enabled": true, "rule": []interface{}{}},
		}, &serviceStruct)

		require.NotNil(t, serviceStruct.RouteRedistributionConfig)
		assert.True(t, *serviceStruct.RouteRedistributionConfig.BgpEnabled)
		require.NotNil(t, serviceStruct.RouteRedistributionConfig.OspfEnabled)
		assert.True(t, *serviceStruct.RouteRedistributionConfig.OspfEnabled)
	})
}

func TestUnitNsxt_initGatewayLocaleServices(t *testing.T) {
	res := gatewayLocaleServiceTestResource()
	noopLister := func(_ utl.SessionContext, _ vapiProtocolClient.Connector, _ string) ([]model.LocaleServices, error) {
		return nil, nil
	}

	t.Run("create generates a new id when nsx_id is unset", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"locale_service": []interface{}{
				map[string]interface{}{
					"edge_cluster_path":     "/infra/sites/default/enforcement-points/default/edge-clusters/ec1",
					"preferred_edge_paths":  []interface{}{},
					"redistribution_config": []interface{}{},
				},
			},
		})

		out, err := initGatewayLocaleServices(utl.SessionContext{ClientType: utl.Local}, d, nil, noopLister)
		require.NoError(t, err)
		require.Len(t, out, 1)
	})

	t.Run("duplicate nsx_id across locale services errors", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"locale_service": []interface{}{
				map[string]interface{}{
					"nsx_id":                "ls-1",
					"edge_cluster_path":     "/infra/sites/default/enforcement-points/default/edge-clusters/ec1",
					"preferred_edge_paths":  []interface{}{},
					"redistribution_config": []interface{}{},
				},
				map[string]interface{}{
					"nsx_id":                "ls-1",
					"edge_cluster_path":     "/infra/sites/default/enforcement-points/default/edge-clusters/ec2",
					"preferred_edge_paths":  []interface{}{},
					"redistribution_config": []interface{}{},
				},
			},
		})

		_, err := initGatewayLocaleServices(utl.SessionContext{ClientType: utl.Local}, d, nil, noopLister)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Duplicate nsx_id")
	})

	t.Run("update deletes locale services no longer present in intent", func(t *testing.T) {
		existingID := "existing-ls"
		lister := func(_ utl.SessionContext, _ vapiProtocolClient.Connector, _ string) ([]model.LocaleServices, error) {
			return []model.LocaleServices{{Id: &existingID}}, nil
		}

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"locale_service": []interface{}{
				map[string]interface{}{
					"nsx_id":                "new-ls",
					"edge_cluster_path":     "/infra/sites/default/enforcement-points/default/edge-clusters/ec1",
					"preferred_edge_paths":  []interface{}{},
					"redistribution_config": []interface{}{},
				},
			},
		})
		d.SetId("gw-1")

		out, err := initGatewayLocaleServices(utl.SessionContext{ClientType: utl.Local}, d, nil, lister)
		require.NoError(t, err)
		// one create instruction for new-ls, one delete instruction for existing-ls
		assert.Len(t, out, 2)
	})

	t.Run("listLocaleServicesFunc error is propagated", func(t *testing.T) {
		failingLister := func(_ utl.SessionContext, _ vapiProtocolClient.Connector, _ string) ([]model.LocaleServices, error) {
			return nil, assertErr("list failed")
		}

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"locale_service": []interface{}{},
		})
		d.SetId("gw-1")

		_, err := initGatewayLocaleServices(utl.SessionContext{ClientType: utl.Local}, d, nil, failingLister)
		require.Error(t, err)
	})
}

type assertErr string

func (e assertErr) Error() string { return string(e) }

func setupTier0LocaleServicesMock(ctrl *gomock.Controller) (*tier0smocks.MockLocaleServicesClient, func()) {
	mockSDK := tier0smocks.NewMockLocaleServicesClient(ctrl)
	mockWrapper := &tier0sAPI.LocaleServicesClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliTier0LocaleServicesClient
	cliTier0LocaleServicesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier0sAPI.LocaleServicesClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier0LocaleServicesClient = original }
}

func TestMockNsxtFindTier0LocaleServiceForSite(t *testing.T) {
	ctx := utl.SessionContext{ClientType: utl.Local}

	t.Run("finds locale service matching site prefix", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTier0LocaleServicesMock(ctrl)
		defer restore()

		lsID := "ls-1"
		edgeClusterPath := "/infra/sites/default/enforcement-points/default/edge-clusters/ec1"
		count := int64(1)
		mockSDK.EXPECT().List("gw1", (*string)(nil), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.LocaleServicesListResult{
				Results:     []model.LocaleServices{{Id: &lsID, EdgeClusterPath: &edgeClusterPath}},
				ResultCount: &count,
			}, nil)

		id, err := findTier0LocaleServiceForSite(ctx, nil, "gw1", "/infra/sites/default")
		require.NoError(t, err)
		assert.Equal(t, "ls-1", id)
	})

	t.Run("no matching site errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTier0LocaleServicesMock(ctrl)
		defer restore()

		count := int64(0)
		mockSDK.EXPECT().List("gw1", (*string)(nil), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.LocaleServicesListResult{Results: []model.LocaleServices{}, ResultCount: &count}, nil)

		_, err := findTier0LocaleServiceForSite(ctx, nil, "gw1", "/infra/sites/other")
		require.Error(t, err)
	})

	t.Run("List error is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTier0LocaleServicesMock(ctrl)
		defer restore()

		mockSDK.EXPECT().List("gw1", (*string)(nil), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			model.LocaleServicesListResult{}, assertErr("list error"))

		_, err := findTier0LocaleServiceForSite(ctx, nil, "gw1", "/infra/sites/default")
		require.Error(t, err)
	})
}

func TestMockNsxtPolicyTier0GetLocaleService(t *testing.T) {
	ctx := utl.SessionContext{ClientType: utl.Local}

	t.Run("success returns locale service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTier0LocaleServicesMock(ctrl)
		defer restore()

		lsID := "ls-1"
		mockSDK.EXPECT().Get("gw1", "ls-1").Return(model.LocaleServices{Id: &lsID}, nil)

		out := testAccPolicyTier0GetLocaleService(ctx, "gw1", "ls-1", nil)
		require.NotNil(t, out)
		assert.Equal(t, "ls-1", *out.Id)
	})

	t.Run("error returns nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupTier0LocaleServicesMock(ctrl)
		defer restore()

		mockSDK.EXPECT().Get("gw1", "ls-1").Return(model.LocaleServices{}, assertErr("not found"))

		out := testAccPolicyTier0GetLocaleService(ctx, "gw1", "ls-1", nil)
		assert.Nil(t, out)
	})
}
