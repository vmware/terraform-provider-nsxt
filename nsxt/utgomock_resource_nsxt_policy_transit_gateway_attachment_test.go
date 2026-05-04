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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	apitgw "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	tgwattmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways"
)

var (
	tgwAttachmentID          = "tgw-att-id"
	tgwAttachmentDisplayName = "test-tgw-attachment"
	tgwAttachmentDescription = "Test TGW Attachment"
	tgwAttachmentRevision    = int64(1)
	tgwAttachmentParentPath  = "/orgs/default/projects/project1/transit-gateways/tgw1"
	tgwAttachmentOrgID       = "default"
	tgwAttachmentProjectID   = "project1"
	tgwAttachmentTGWID       = "tgw1"
	tgwAttachmentConnPath    = "/orgs/default/projects/project1/vpcs/vpc1"
)

func tgwAttachmentAPIResponse() nsxModel.TransitGatewayAttachment {
	return nsxModel.TransitGatewayAttachment{
		Id:             &tgwAttachmentID,
		DisplayName:    &tgwAttachmentDisplayName,
		Description:    &tgwAttachmentDescription,
		Revision:       &tgwAttachmentRevision,
		ConnectionPath: &tgwAttachmentConnPath,
	}
}

func minimalTGWAttachmentData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    tgwAttachmentDisplayName,
		"description":     tgwAttachmentDescription,
		"nsx_id":          tgwAttachmentID,
		"parent_path":     tgwAttachmentParentPath,
		"connection_path": tgwAttachmentConnPath,
	}
}

func setupTGWAttachmentMock(t *testing.T, ctrl *gomock.Controller) (*tgwattmocks.MockAttachmentsClient, func()) {
	mockSDK := tgwattmocks.NewMockAttachmentsClient(ctrl)
	mockWrapper := &apitgw.TransitGatewayAttachmentClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwAttachmentProjectID,
	}

	original := cliTransitGatewayAttachmentsClient
	cliTransitGatewayAttachmentsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apitgw.TransitGatewayAttachmentClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTransitGatewayAttachmentsClient = original }
}

func TestMockResourceNsxtPolicyTransitGatewayAttachmentCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWAttachmentMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID).Return(nsxModel.TransitGatewayAttachment{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID).Return(tgwAttachmentAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWAttachmentData())

		err := resourceNsxtPolicyTransitGatewayAttachmentCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwAttachmentID, d.Id())
	})

	t.Run("Create fails with invalid parent_path", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayAttachment()
		data := minimalTGWAttachmentData()
		data["parent_path"] = "/orgs/default"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayAttachmentCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Create fails when Patch returns error", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID).Return(nsxModel.TransitGatewayAttachment{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID, gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtPolicyTransitGatewayAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWAttachmentData())

		err := resourceNsxtPolicyTransitGatewayAttachmentCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayAttachmentRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWAttachmentMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID).Return(tgwAttachmentAPIResponse(), nil)

		res := resourceNsxtPolicyTransitGatewayAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWAttachmentData())
		d.SetId(tgwAttachmentID)

		err := resourceNsxtPolicyTransitGatewayAttachmentRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwAttachmentDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID).Return(nsxModel.TransitGatewayAttachment{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTransitGatewayAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWAttachmentData())
		d.SetId(tgwAttachmentID)

		err := resourceNsxtPolicyTransitGatewayAttachmentRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWAttachmentData())

		err := resourceNsxtPolicyTransitGatewayAttachmentRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayAttachmentUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWAttachmentMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Update(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID, gomock.Any()).Return(tgwAttachmentAPIResponse(), nil),
			mockSDK.EXPECT().Get(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID).Return(tgwAttachmentAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWAttachmentData())
		d.SetId(tgwAttachmentID)

		err := resourceNsxtPolicyTransitGatewayAttachmentUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWAttachmentData())

		err := resourceNsxtPolicyTransitGatewayAttachmentUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayAttachmentNewAttributes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWAttachmentMock(t, ctrl)
	defer restore()

	t.Run("Create with cna_path instead of connection_path", func(t *testing.T) {
		cnaPath := "/orgs/default/projects/project1/vpcs/vpc1/subnets/sub1"
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID).Return(nsxModel.TransitGatewayAttachment{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID).Return(nsxModel.TransitGatewayAttachment{
				Id:          &tgwAttachmentID,
				DisplayName: &tgwAttachmentDisplayName,
				Description: &tgwAttachmentDescription,
				Revision:    &tgwAttachmentRevision,
				CnaPath:     &cnaPath,
			}, nil),
		)

		data := map[string]interface{}{
			"display_name": tgwAttachmentDisplayName,
			"description":  tgwAttachmentDescription,
			"nsx_id":       tgwAttachmentID,
			"parent_path":  tgwAttachmentParentPath,
			"cna_path":     cnaPath,
		}

		res := resourceNsxtPolicyTransitGatewayAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayAttachmentCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwAttachmentID, d.Id())
		assert.Equal(t, cnaPath, d.Get("cna_path"))
	})

	t.Run("Read sets admin_state, urpf_mode and route_advertisement_types from API", func(t *testing.T) {
		adminState := nsxModel.TransitGatewayAttachment_ADMIN_STATE_UP
		urpfMode := nsxModel.TransitGatewayAttachment_URPF_MODE_STRICT
		advTypePublic := nsxModel.TransitGatewayRouteAdvertisementRule_ROUTE_ADVERTISEMENT_TYPE_PUBLIC
		advTypePrivate := nsxModel.TransitGatewayRouteAdvertisementRule_ROUTE_ADVERTISEMENT_TYPE_TGW_PRIVATE
		mockSDK.EXPECT().Get(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID).Return(nsxModel.TransitGatewayAttachment{
			Id:             &tgwAttachmentID,
			DisplayName:    &tgwAttachmentDisplayName,
			Description:    &tgwAttachmentDescription,
			Revision:       &tgwAttachmentRevision,
			ConnectionPath: &tgwAttachmentConnPath,
			AdminState:     &adminState,
			UrpfMode:       &urpfMode,
			RouteAdvertisementRules: []nsxModel.TransitGatewayRouteAdvertisementRule{
				{RouteAdvertisementType: &advTypePublic},
				{RouteAdvertisementType: &advTypePrivate},
			},
		}, nil)

		res := resourceNsxtPolicyTransitGatewayAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWAttachmentData())
		d.SetId(tgwAttachmentID)

		err := resourceNsxtPolicyTransitGatewayAttachmentRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, adminState, d.Get("admin_state"))
		assert.Equal(t, urpfMode, d.Get("urpf_mode"))
		types := d.Get("route_advertisement_types").([]interface{})
		require.Len(t, types, 2)
		assert.Equal(t, advTypePublic, types[0])
		assert.Equal(t, advTypePrivate, types[1])
	})

	t.Run("Create with route_advertisement_types sends rules to API", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		var capturedObj nsxModel.TransitGatewayAttachment
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID).Return(nsxModel.TransitGatewayAttachment{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID, gomock.Any()).DoAndReturn(
				func(_, _, _, _ string, obj nsxModel.TransitGatewayAttachment) error {
					capturedObj = obj
					return nil
				}),
			mockSDK.EXPECT().Get(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID).Return(tgwAttachmentAPIResponse(), nil),
		)

		data := minimalTGWAttachmentData()
		data["route_advertisement_types"] = []interface{}{"PUBLIC", "TGW_PRIVATE"}

		res := resourceNsxtPolicyTransitGatewayAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayAttachmentCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, capturedObj.RouteAdvertisementRules, 2)
		assert.Equal(t, "PUBLIC", *capturedObj.RouteAdvertisementRules[0].RouteAdvertisementType)
		assert.Equal(t, "TGW_PRIVATE", *capturedObj.RouteAdvertisementRules[1].RouteAdvertisementType)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayAttachmentDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWAttachmentMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(tgwAttachmentOrgID, tgwAttachmentProjectID, tgwAttachmentTGWID, tgwAttachmentID).Return(nil)

		res := resourceNsxtPolicyTransitGatewayAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWAttachmentData())
		d.SetId(tgwAttachmentID)

		err := resourceNsxtPolicyTransitGatewayAttachmentDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWAttachmentData())

		err := resourceNsxtPolicyTransitGatewayAttachmentDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
