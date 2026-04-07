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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	apitgw "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	sroutmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	tgwStaticRouteID          = "tgw-static-route-id"
	tgwStaticRouteDisplayName = "test-tgw-static-route"
	tgwStaticRouteDescription = "Test TGW Static Route"
	tgwStaticRouteRevision    = int64(1)
	tgwStaticRouteParentPath  = "/orgs/default/projects/project1/transit-gateways/tgw1"
	tgwStaticRouteOrgID       = "default"
	tgwStaticRouteProjectID   = "project1"
	tgwStaticRouteTGWID       = "tgw1"
	tgwStaticRouteNetwork     = "10.0.0.0/24"
)

func tgwStaticRouteAPIResponse() nsxModel.StaticRoutes {
	return nsxModel.StaticRoutes{
		Id:          &tgwStaticRouteID,
		DisplayName: &tgwStaticRouteDisplayName,
		Description: &tgwStaticRouteDescription,
		Revision:    &tgwStaticRouteRevision,
		Network:     &tgwStaticRouteNetwork,
	}
}

func minimalTGWStaticRouteData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": tgwStaticRouteDisplayName,
		"description":  tgwStaticRouteDescription,
		"nsx_id":       tgwStaticRouteID,
		"parent_path":  tgwStaticRouteParentPath,
		"network":      tgwStaticRouteNetwork,
	}
}

func setupTGWStaticRouteMock(t *testing.T, ctrl *gomock.Controller) (*sroutmocks.MockStaticRoutesClient, func()) {
	mockSDK := sroutmocks.NewMockStaticRoutesClient(ctrl)
	mockWrapper := &apitgw.TransitGatewayStaticRoutesClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwStaticRouteProjectID,
	}

	original := cliTransitGatewayStaticRoutesClient
	cliTransitGatewayStaticRoutesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apitgw.TransitGatewayStaticRoutesClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTransitGatewayStaticRoutesClient = original }
}

func TestMockResourceNsxtPolicyTransitGatewayStaticRouteCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWStaticRouteMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.1.0"
		defer func() { util.NsxVersion = "" }()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwStaticRouteOrgID, tgwStaticRouteProjectID, tgwStaticRouteTGWID, tgwStaticRouteID).Return(nsxModel.StaticRoutes{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwStaticRouteOrgID, tgwStaticRouteProjectID, tgwStaticRouteTGWID, tgwStaticRouteID, gomock.Any()).Return(tgwStaticRouteAPIResponse(), nil),
			mockSDK.EXPECT().Get(tgwStaticRouteOrgID, tgwStaticRouteProjectID, tgwStaticRouteTGWID, tgwStaticRouteID).Return(tgwStaticRouteAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWStaticRouteData())

		err := resourceNsxtPolicyTransitGatewayStaticRouteCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwStaticRouteID, d.Id())
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtPolicyTransitGatewayStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWStaticRouteData())

		err := resourceNsxtPolicyTransitGatewayStaticRouteCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.1.0")
	})

	t.Run("Create fails with invalid parent_path", func(t *testing.T) {
		util.NsxVersion = "9.1.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtPolicyTransitGatewayStaticRoute()
		data := minimalTGWStaticRouteData()
		data["parent_path"] = "/orgs/default"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayStaticRouteCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayStaticRouteRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWStaticRouteMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwStaticRouteOrgID, tgwStaticRouteProjectID, tgwStaticRouteTGWID, tgwStaticRouteID).Return(tgwStaticRouteAPIResponse(), nil)

		res := resourceNsxtPolicyTransitGatewayStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWStaticRouteData())
		d.SetId(tgwStaticRouteID)

		err := resourceNsxtPolicyTransitGatewayStaticRouteRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwStaticRouteDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwStaticRouteOrgID, tgwStaticRouteProjectID, tgwStaticRouteTGWID, tgwStaticRouteID).Return(nsxModel.StaticRoutes{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTransitGatewayStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWStaticRouteData())
		d.SetId(tgwStaticRouteID)

		err := resourceNsxtPolicyTransitGatewayStaticRouteRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWStaticRouteData())

		err := resourceNsxtPolicyTransitGatewayStaticRouteRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayStaticRouteUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWStaticRouteMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(tgwStaticRouteOrgID, tgwStaticRouteProjectID, tgwStaticRouteTGWID, tgwStaticRouteID, gomock.Any()).Return(tgwStaticRouteAPIResponse(), nil),
			mockSDK.EXPECT().Get(tgwStaticRouteOrgID, tgwStaticRouteProjectID, tgwStaticRouteTGWID, tgwStaticRouteID).Return(tgwStaticRouteAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWStaticRouteData())
		d.SetId(tgwStaticRouteID)

		err := resourceNsxtPolicyTransitGatewayStaticRouteUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWStaticRouteData())

		err := resourceNsxtPolicyTransitGatewayStaticRouteUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayStaticRouteDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWStaticRouteMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(tgwStaticRouteOrgID, tgwStaticRouteProjectID, tgwStaticRouteTGWID, tgwStaticRouteID).Return(nil)

		res := resourceNsxtPolicyTransitGatewayStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWStaticRouteData())
		d.SetId(tgwStaticRouteID)

		err := resourceNsxtPolicyTransitGatewayStaticRouteDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayStaticRoute()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWStaticRouteData())

		err := resourceNsxtPolicyTransitGatewayStaticRouteDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
