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
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	tgwrouting "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/transit_gateways/routing"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	tgwmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/transit_gateways/routing"
)

var (
	tgwRMID          = "route-map-1"
	tgwRMDisplayName = "Test TGW Route Map"
	tgwRMDescription = "Test transit gateway route map"
	tgwRMRevision    = int64(1)
	tgwRMParentPath  = "/orgs/default/projects/project1/transit-gateways/tgw1"
	tgwRMOrgID       = "default"
	tgwRMProjectID   = "project1"
	tgwRMTGWID       = "tgw1"
	tgwRMPath        = "/orgs/default/projects/project1/transit-gateways/tgw1/route-maps/route-map-1"
	tgwRMAction      = nsxModel.RouteMapEntry_ACTION_PERMIT
	tgwRMPrefixPath  = "/orgs/default/projects/project1/transit-gateways/tgw1/prefix-lists/pl-1"
)

func tgwRouteMapAPIResponse() nsxModel.TransitGatewayRouteMap {
	return nsxModel.TransitGatewayRouteMap{
		Id:          &tgwRMID,
		DisplayName: &tgwRMDisplayName,
		Description: &tgwRMDescription,
		Revision:    &tgwRMRevision,
		Path:        &tgwRMPath,
		Entries: []nsxModel.RouteMapEntry{
			{
				Action:            &tgwRMAction,
				PrefixListMatches: []string{tgwRMPrefixPath},
			},
		},
	}
}

func minimalTGWRouteMapData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": tgwRMDisplayName,
		"description":  tgwRMDescription,
		"nsx_id":       tgwRMID,
		"parent_path":  tgwRMParentPath,
		"entry": []interface{}{
			map[string]interface{}{
				"action":               tgwRMAction,
				"prefix_list_matches":  []interface{}{tgwRMPrefixPath},
				"community_list_match": []interface{}{},
				"set":                  []interface{}{},
			},
		},
	}
}

func setupTGWRouteMapMock(t *testing.T, ctrl *gomock.Controller) (*tgwmocks.MockRouteMapsClient, func()) {
	mockSDK := tgwmocks.NewMockRouteMapsClient(ctrl)
	mockWrapper := &tgwrouting.TGWRouteMapClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  tgwRMProjectID,
	}

	original := cliTGWRouteMapsClient
	cliTGWRouteMapsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tgwrouting.TGWRouteMapClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTGWRouteMapsClient = original }
}

func TestMockResourceNsxtPolicyTransitGatewayRouteMapCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWRouteMapMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwRMOrgID, tgwRMProjectID, tgwRMTGWID, tgwRMID).Return(nsxModel.TransitGatewayRouteMap{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwRMOrgID, tgwRMProjectID, tgwRMTGWID, tgwRMID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwRMOrgID, tgwRMProjectID, tgwRMTGWID, tgwRMID).Return(tgwRouteMapAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWRouteMapData())

		err := resourceNsxtPolicyTransitGatewayRouteMapCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwRMID, d.Id())
		assert.Equal(t, tgwRMDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when parent_path is missing", func(t *testing.T) {
		data := minimalTGWRouteMapData()
		data["parent_path"] = ""

		res := resourceNsxtPolicyTransitGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayRouteMapCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Create fails when prefix_list_matches and community_list_match are both set", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		mockSDK.EXPECT().Get(tgwRMOrgID, tgwRMProjectID, tgwRMTGWID, tgwRMID).Return(nsxModel.TransitGatewayRouteMap{}, notFoundErr)

		data := minimalTGWRouteMapData()
		data["entry"] = []interface{}{
			map[string]interface{}{
				"action":              tgwRMAction,
				"prefix_list_matches": []interface{}{tgwRMPrefixPath},
				"community_list_match": []interface{}{
					map[string]interface{}{
						"criteria":       "11:22",
						"match_operator": nsxModel.CommunityMatchCriteria_MATCH_OPERATOR_ANY,
					},
				},
				"set": []interface{}{},
			},
		}

		res := resourceNsxtPolicyTransitGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayRouteMapCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "mutually exclusive")
	})
}

func TestMockResourceNsxtPolicyTransitGatewayRouteMapRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWRouteMapMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwRMOrgID, tgwRMProjectID, tgwRMTGWID, tgwRMID).Return(tgwRouteMapAPIResponse(), nil)

		res := resourceNsxtPolicyTransitGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWRouteMapData())
		d.SetId(tgwRMID)

		err := resourceNsxtPolicyTransitGatewayRouteMapRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwRMDisplayName, d.Get("display_name"))

		entries := d.Get("entry").([]interface{})
		require.Len(t, entries, 1)
		e := entries[0].(map[string]interface{})
		assert.Equal(t, tgwRMAction, e["action"])
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(tgwRMOrgID, tgwRMProjectID, tgwRMTGWID, tgwRMID).Return(nsxModel.TransitGatewayRouteMap{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTransitGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWRouteMapData())
		d.SetId(tgwRMID)

		err := resourceNsxtPolicyTransitGatewayRouteMapRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWRouteMapData())

		err := resourceNsxtPolicyTransitGatewayRouteMapRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayRouteMapUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWRouteMapMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(tgwRMOrgID, tgwRMProjectID, tgwRMTGWID, tgwRMID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(tgwRMOrgID, tgwRMProjectID, tgwRMTGWID, tgwRMID).Return(tgwRouteMapAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTransitGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWRouteMapData())
		d.SetId(tgwRMID)

		err := resourceNsxtPolicyTransitGatewayRouteMapUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwRMID, d.Id())
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWRouteMapData())

		err := resourceNsxtPolicyTransitGatewayRouteMapUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayRouteMapDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWRouteMapMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(tgwRMOrgID, tgwRMProjectID, tgwRMTGWID, tgwRMID).Return(nil)

		res := resourceNsxtPolicyTransitGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWRouteMapData())
		d.SetId(tgwRMID)

		err := resourceNsxtPolicyTransitGatewayRouteMapDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyTransitGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalTGWRouteMapData())

		err := resourceNsxtPolicyTransitGatewayRouteMapDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTransitGatewayRouteMapWithSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupTGWRouteMapMock(t, ctrl)
	defer restore()

	t.Run("Create and update with set clause succeed", func(t *testing.T) {
		localPref := int64(200)
		med := int64(100)
		weight := int64(10)

		response := nsxModel.TransitGatewayRouteMap{
			Id:          &tgwRMID,
			DisplayName: &tgwRMDisplayName,
			Description: &tgwRMDescription,
			Revision:    &tgwRMRevision,
			Path:        &tgwRMPath,
			Entries: []nsxModel.RouteMapEntry{
				{
					Action:            &tgwRMAction,
					PrefixListMatches: []string{tgwRMPrefixPath},
					Set: &nsxModel.RouteMapEntrySet{
						LocalPreference: &localPref,
						Med:             &med,
						Weight:          &weight,
					},
				},
			},
		}

		data := map[string]interface{}{
			"display_name": tgwRMDisplayName,
			"description":  tgwRMDescription,
			"nsx_id":       tgwRMID,
			"parent_path":  tgwRMParentPath,
			"entry": []interface{}{
				map[string]interface{}{
					"action":               tgwRMAction,
					"prefix_list_matches":  []interface{}{tgwRMPrefixPath},
					"community_list_match": []interface{}{},
					"set": []interface{}{
						map[string]interface{}{
							"local_preference": 200,
							"med":              100,
							"weight":           10,
							"as_path_prepend":  "",
							"community":        "",
						},
					},
				},
			},
		}

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(tgwRMOrgID, tgwRMProjectID, tgwRMTGWID, tgwRMID).Return(nsxModel.TransitGatewayRouteMap{}, notFoundErr),
			mockSDK.EXPECT().Patch(tgwRMOrgID, tgwRMProjectID, tgwRMTGWID, tgwRMID, gomock.Any()).DoAndReturn(
				func(_, _, _, _ string, obj nsxModel.TransitGatewayRouteMap) error {
					require.Len(t, obj.Entries, 1)
					require.NotNil(t, obj.Entries[0].Set)
					assert.Nil(t, obj.Entries[0].Set.PreferGlobalV6NextHop)
					assert.Equal(t, &localPref, obj.Entries[0].Set.LocalPreference)
					assert.Equal(t, &med, obj.Entries[0].Set.Med)
					assert.Equal(t, &weight, obj.Entries[0].Set.Weight)
					return nil
				}),
			mockSDK.EXPECT().Get(tgwRMOrgID, tgwRMProjectID, tgwRMTGWID, tgwRMID).Return(response, nil),
		)

		res := resourceNsxtPolicyTransitGatewayRouteMap()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtPolicyTransitGatewayRouteMapCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, tgwRMID, d.Id())

		gomock.InOrder(
			mockSDK.EXPECT().Patch(tgwRMOrgID, tgwRMProjectID, tgwRMTGWID, tgwRMID, gomock.Any()).DoAndReturn(
				func(_, _, _, _ string, obj nsxModel.TransitGatewayRouteMap) error {
					require.Len(t, obj.Entries, 1)
					require.NotNil(t, obj.Entries[0].Set)
					assert.Nil(t, obj.Entries[0].Set.PreferGlobalV6NextHop)
					assert.Equal(t, &localPref, obj.Entries[0].Set.LocalPreference)
					assert.Equal(t, &med, obj.Entries[0].Set.Med)
					assert.Equal(t, &weight, obj.Entries[0].Set.Weight)
					return nil
				}),
			mockSDK.EXPECT().Get(tgwRMOrgID, tgwRMProjectID, tgwRMTGWID, tgwRMID).Return(response, nil),
		)
		require.NoError(t, resourceNsxtPolicyTransitGatewayRouteMapUpdate(d, newGoMockProviderClient()))
	})
}

func TestTGWRouteMapSchemaOmitsGlobalV6NextHop(t *testing.T) {
	tgw := resourceNsxtPolicyTransitGatewayRouteMap()
	tgwEntrySchema, ok := tgw.Schema["entry"].Elem.(*schema.Resource)
	require.True(t, ok)
	tgwSetSchema, ok := tgwEntrySchema.Schema["set"].Elem.(*schema.Resource)
	require.True(t, ok)
	assert.NotContains(t, tgwSetSchema.Schema, "prefer_global_v6_next_hop")

	gateway := resourceNsxtPolicyGatewayRouteMap()
	gatewayEntrySchema, ok := gateway.Schema["entry"].Elem.(*schema.Resource)
	require.True(t, ok)
	gatewaySetSchema, ok := gatewayEntrySchema.Schema["set"].Elem.(*schema.Resource)
	require.True(t, ok)
	assert.Contains(t, gatewaySetSchema.Schema, "prefer_global_v6_next_hop")

	data := minimalTGWRouteMapData()
	data["entry"] = []interface{}{map[string]interface{}{
		"action": "PERMIT",
		"set": []interface{}{map[string]interface{}{
			"local_preference": 200,
		}},
	}}
	d := schema.TestResourceDataRaw(t, tgw.Schema, data)
	obj, err := tgwRouteMapBuildObj(d, tgwRMID)
	require.NoError(t, err)
	require.Len(t, obj.Entries, 1)
	require.NotNil(t, obj.Entries[0].Set)
	assert.Nil(t, obj.Entries[0].Set.PreferGlobalV6NextHop)
}

func TestTGWRouteMapRefreshIgnoresGlobalV6NextHop(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value bool
	}{
		{name: "false"},
		{name: "true", value: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Start with empty state, as on import, before the API populates entries.
			res := resourceNsxtPolicyTransitGatewayRouteMap()
			d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
			localPreference := int64(200)
			entries := []nsxModel.RouteMapEntry{{
				Action: &tgwRMAction,
				Set: &nsxModel.RouteMapEntrySet{
					LocalPreference:       &localPreference,
					PreferGlobalV6NextHop: &tc.value,
				},
			}}
			require.NoError(t, setTGWRouteMapEntriesInSchema(d, entries))
			assert.Equal(t, 200, d.Get("entry.0.set.0.local_preference"))
			assert.Nil(t, d.Get("entry.0.set.0.prefer_global_v6_next_hop"))
			payload, err := tgwRouteMapBuildObj(d, tgwRMID)
			require.NoError(t, err)
			require.Len(t, payload.Entries, 1)
			require.NotNil(t, payload.Entries[0].Set)
			assert.Nil(t, payload.Entries[0].Set.PreferGlobalV6NextHop)
			assert.Equal(t, &localPreference, payload.Entries[0].Set.LocalPreference)

			entries[0].Set = nil
			require.NoError(t, setTGWRouteMapEntriesInSchema(d, entries))
			assert.Equal(t, 0, d.Get("entry.0.set.#"))
		})
	}
}
