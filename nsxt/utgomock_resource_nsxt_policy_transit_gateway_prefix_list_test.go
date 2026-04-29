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

	transitgateways "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	tgwmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways"
)

var (
	tgwPLID          = "prefix-list-1"
	tgwPLDisplayName = "Test TGW Prefix List"
	tgwPLDescription = "Test transit gateway prefix list"
	tgwPLRevision    = int64(1)
	tgwPLParentPath  = "/orgs/default/projects/project1/transit-gateways/tgw1"
	tgwPLOrgID       = "default"
	tgwPLProjectID   = "project1"
	tgwPLTGWID       = "tgw1"
	tgwPLPath        = "/orgs/default/projects/project1/transit-gateways/tgw1/prefix-lists/prefix-list-1"
	tgwPLNetwork     = "10.0.0.0/8"
	tgwPLAction      = nsxModel.PrefixEntry_ACTION_PERMIT
)

func tgwPrefixListAPIResponse() nsxModel.PrefixList {
	return nsxModel.PrefixList{
		Id:          &tgwPLID,
		DisplayName: &tgwPLDisplayName,
		Description: &tgwPLDescription,
		Revision:    &tgwPLRevision,
		Path:        &tgwPLPath,
		Prefixes: []nsxModel.PrefixEntry{
			{
				Action:  &tgwPLAction,
				Network: &tgwPLNetwork,
			},
		},
	}
}

func minimalTGWPrefixListData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": tgwPLDisplayName,
		"description":  tgwPLDescription,
		"nsx_id":       tgwPLID,
		"parent_path":  tgwPLParentPath,
		"prefix": []interface{}{
			map[string]interface{}{
				"action":  tgwPLAction,
				"network": tgwPLNetwork,
				"ge":      0,
				"le":      0,
			},
		},
	}
}

func setupTGWPrefixListMock(t *testing.T, ctrl *gomock.Controller) (*tgwmocks.MockPrefixListsClient, func()) {
	mockSDK := tgwmocks.NewMockPrefixListsClient(ctrl)
	mockWrapper := &transitgateways.TGWPrefixListClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwPLProjectID,
	}

	original := cliTGWPrefixListsClient
	cliTGWPrefixListsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *transitgateways.TGWPrefixListClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTGWPrefixListsClient = original }
}

func TestMockResourceNsxtPolicyTransitGatewayPrefixListCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWPrefixListMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwPLOrgID, tgwPLProjectID, tgwPLTGWID, tgwPLID).Return(nsxModel.PrefixList{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwPLOrgID, tgwPLProjectID, tgwPLTGWID, tgwPLID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwPLOrgID, tgwPLProjectID, tgwPLTGWID, tgwPLID).Return(tgwPrefixListAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWPrefixListData())

		err := resourceNsxtPolicyTransitGatewayPrefixListCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwPLID, d.Id())
		assert.Equal(t, tgwPLDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when parent_path is missing", func(t *testing.T) {
		data := minimalTGWPrefixListData()
		data["parent_path"] = ""

		res := resourceNsxtPolicyTransitGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayPrefixListCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayPrefixListRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWPrefixListMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwPLOrgID, tgwPLProjectID, tgwPLTGWID, tgwPLID).Return(tgwPrefixListAPIResponse(), nil)

		res := resourceNsxtPolicyTransitGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWPrefixListData())
		d.SetId(tgwPLID)

		err := resourceNsxtPolicyTransitGatewayPrefixListRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwPLDisplayName, d.Get("display_name"))

		prefixes := d.Get("prefix").([]interface{})
		require.Len(t, prefixes, 1)
		p := prefixes[0].(map[string]interface{})
		assert.Equal(t, tgwPLNetwork, p["network"])
		assert.Equal(t, tgwPLAction, p["action"])
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwPLOrgID, tgwPLProjectID, tgwPLTGWID, tgwPLID).Return(nsxModel.PrefixList{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTransitGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWPrefixListData())
		d.SetId(tgwPLID)

		err := resourceNsxtPolicyTransitGatewayPrefixListRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWPrefixListData())

		err := resourceNsxtPolicyTransitGatewayPrefixListRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayPrefixListUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWPrefixListMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(tgwPLOrgID, tgwPLProjectID, tgwPLTGWID, tgwPLID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwPLOrgID, tgwPLProjectID, tgwPLTGWID, tgwPLID).Return(tgwPrefixListAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWPrefixListData())
		d.SetId(tgwPLID)

		err := resourceNsxtPolicyTransitGatewayPrefixListUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwPLID, d.Id())
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWPrefixListData())

		err := resourceNsxtPolicyTransitGatewayPrefixListUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayPrefixListDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWPrefixListMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(tgwPLOrgID, tgwPLProjectID, tgwPLTGWID, tgwPLID).Return(nil)

		res := resourceNsxtPolicyTransitGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWPrefixListData())
		d.SetId(tgwPLID)

		err := resourceNsxtPolicyTransitGatewayPrefixListDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWPrefixListData())

		err := resourceNsxtPolicyTransitGatewayPrefixListDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayPrefixListMultiplePrefixes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWPrefixListMock(t, ctrl)
	defer restore()

	t.Run("Create with multiple prefixes including ge/le", func(t *testing.T) {
		denyAction := nsxModel.PrefixEntry_ACTION_DENY
		network2 := "172.16.0.0/12"
		ge := int64(16)
		le := int64(24)

		response := nsxModel.PrefixList{
			Id:          &tgwPLID,
			DisplayName: &tgwPLDisplayName,
			Description: &tgwPLDescription,
			Revision:    &tgwPLRevision,
			Path:        &tgwPLPath,
			Prefixes: []nsxModel.PrefixEntry{
				{Action: &tgwPLAction, Network: &tgwPLNetwork},
				{Action: &denyAction, Network: &network2, Ge: &ge, Le: &le},
			},
		}

		data := minimalTGWPrefixListData()
		data["prefix"] = []interface{}{
			map[string]interface{}{
				"action":  tgwPLAction,
				"network": tgwPLNetwork,
				"ge":      0,
				"le":      0,
			},
			map[string]interface{}{
				"action":  nsxModel.PrefixEntry_ACTION_DENY,
				"network": "172.16.0.0/12",
				"ge":      16,
				"le":      24,
			},
		}

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwPLOrgID, tgwPLProjectID, tgwPLTGWID, tgwPLID).Return(nsxModel.PrefixList{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwPLOrgID, tgwPLProjectID, tgwPLTGWID, tgwPLID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwPLOrgID, tgwPLProjectID, tgwPLTGWID, tgwPLID).Return(response, nil),
		)

		res := resourceNsxtPolicyTransitGatewayPrefixList()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayPrefixListCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwPLID, d.Id())

		prefixes := d.Get("prefix").([]interface{})
		require.Len(t, prefixes, 2)
	})
}
