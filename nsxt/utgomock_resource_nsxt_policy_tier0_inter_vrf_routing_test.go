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

	tier0sapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	intervrfmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
)

var (
	interVrfID          = "inter-vrf-1"
	interVrfDisplayName = "Test InterVRF Route"
	interVrfDescription = "Test inter VRF routing"
	interVrfRevision    = int64(1)
	interVrfGwPath      = "/infra/tier-0s/t0-gw-1"
	interVrfGwID        = "t0-gw-1"
	interVrfTargetPath  = "/infra/tier-0s/t0-gw-2"
)

func interVrfAPIResponse() nsxModel.PolicyInterVrfRoutingConfig {
	return nsxModel.PolicyInterVrfRoutingConfig{
		Id:          &interVrfID,
		DisplayName: &interVrfDisplayName,
		Description: &interVrfDescription,
		Revision:    &interVrfRevision,
		TargetPath:  &interVrfTargetPath,
	}
}

func minimalInterVrfData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": interVrfDisplayName,
		"description":  interVrfDescription,
		"nsx_id":       interVrfID,
		"gateway_path": interVrfGwPath,
		"target_path":  interVrfTargetPath,
	}
}

func setupInterVrfMock(t *testing.T, ctrl *gomock.Controller) (*intervrfmocks.MockInterVrfRoutingClient, func()) {
	mockSDK := intervrfmocks.NewMockInterVrfRoutingClient(ctrl)
	mockWrapper := &tier0sapi.PolicyInterVrfRoutingConfigClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliInterVrfRoutingClient
	cliInterVrfRoutingClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *tier0sapi.PolicyInterVrfRoutingConfigClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliInterVrfRoutingClient = original }
}

func TestMockResourceNsxtPolicyTier0InterVRFRoutingCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupInterVrfMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(interVrfGwID, interVrfID).Return(nsxModel.PolicyInterVrfRoutingConfig{}, notFoundErr),
			mockSDK.EXPECT().Patch(interVrfGwID, interVrfID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(interVrfGwID, interVrfID).Return(interVrfAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())

		err := resourceNsxtPolicyTier0InterVRFRoutingCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, interVrfID, d.Id())
		assert.Equal(t, interVrfDisplayName, d.Get("display_name"))
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(interVrfGwID, interVrfID).Return(interVrfAPIResponse(), nil)

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())

		err := resourceNsxtPolicyTier0InterVRFRoutingCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Create fails for non-Tier0 gateway path", func(t *testing.T) {
		t1Data := minimalInterVrfData()
		t1Data["gateway_path"] = "/infra/tier-1s/t1-gw-1"

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, t1Data)

		err := resourceNsxtPolicyTier0InterVRFRoutingCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "tier0")
	})
}

func TestMockResourceNsxtPolicyTier0InterVRFRoutingRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupInterVrfMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(interVrfGwID, interVrfID).Return(interVrfAPIResponse(), nil)

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())
		d.SetId(interVrfID)

		err := resourceNsxtPolicyTier0InterVRFRoutingRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, interVrfDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(interVrfGwID, interVrfID).Return(nsxModel.PolicyInterVrfRoutingConfig{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())
		d.SetId(interVrfID)

		err := resourceNsxtPolicyTier0InterVRFRoutingRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails for non-Tier0 gateway path", func(t *testing.T) {
		t1Data := minimalInterVrfData()
		t1Data["gateway_path"] = "/infra/tier-1s/t1-gw-1"

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, t1Data)
		d.SetId(interVrfID)

		err := resourceNsxtPolicyTier0InterVRFRoutingRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Tier0")
	})
}

func TestMockResourceNsxtPolicyTier0InterVRFRoutingUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupInterVrfMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(interVrfGwID, interVrfID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(interVrfGwID, interVrfID).Return(interVrfAPIResponse(), nil),
		)

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())
		d.SetId(interVrfID)

		err := resourceNsxtPolicyTier0InterVRFRoutingUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails for non-Tier0 gateway path", func(t *testing.T) {
		t1Data := minimalInterVrfData()
		t1Data["gateway_path"] = "/infra/tier-1s/t1-gw-1"

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, t1Data)
		d.SetId(interVrfID)

		err := resourceNsxtPolicyTier0InterVRFRoutingUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Tier0")
	})

	t.Run("Update fails with API error", func(t *testing.T) {
		mockSDK.EXPECT().Patch(interVrfGwID, interVrfID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())
		d.SetId(interVrfID)

		err := resourceNsxtPolicyTier0InterVRFRoutingUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier0InterVRFRoutingDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupInterVrfMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(interVrfGwID, interVrfID).Return(nil)

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())
		d.SetId(interVrfID)

		err := resourceNsxtPolicyTier0InterVRFRoutingDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails with API error", func(t *testing.T) {
		mockSDK.EXPECT().Delete(interVrfGwID, interVrfID).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyTier0InterVRFRouting()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())
		d.SetId(interVrfID)

		err := resourceNsxtPolicyTier0InterVRFRoutingDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyTier0InterVRFRoutingGetFromSchema(t *testing.T) {
	res := resourceNsxtPolicyTier0InterVRFRouting()

	t.Run("minimal fields produce empty optional collections", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, minimalInterVrfData())
		obj := getPolicyInterVRFRoutingFromSchema(d)
		assert.Equal(t, interVrfDisplayName, *obj.DisplayName)
		assert.Equal(t, interVrfDescription, *obj.Description)
		assert.Equal(t, interVrfTargetPath, *obj.TargetPath)
		assert.Empty(t, obj.BgpRouteLeaking)
		assert.Nil(t, obj.StaticRouteAdvertisement)
	})

	t.Run("bgp_route_leaking IPV4 with in and out filters is parsed", func(t *testing.T) {
		data := minimalInterVrfData()
		inPath := "/infra/tier-0s/t0-gw-1/route-maps/rm-in"
		outPath := "/infra/tier-0s/t0-gw-1/route-maps/rm-out"
		data["bgp_route_leaking"] = []interface{}{
			map[string]interface{}{
				"address_family": "IPV4",
				"in_filter":      []interface{}{inPath},
				"out_filter":     []interface{}{outPath},
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		obj := getPolicyInterVRFRoutingFromSchema(d)
		require.Len(t, obj.BgpRouteLeaking, 1)
		assert.Equal(t, "IPV4", *obj.BgpRouteLeaking[0].AddressFamily)
		assert.Equal(t, []string{inPath}, obj.BgpRouteLeaking[0].InFilter)
		assert.Equal(t, []string{outPath}, obj.BgpRouteLeaking[0].OutFilter)
	})

	t.Run("two bgp_route_leaking entries for IPV4 and IPV6 are both parsed", func(t *testing.T) {
		data := minimalInterVrfData()
		data["bgp_route_leaking"] = []interface{}{
			map[string]interface{}{
				"address_family": "IPV4",
				"in_filter":      []interface{}{},
				"out_filter":     []interface{}{},
			},
			map[string]interface{}{
				"address_family": "IPV6",
				"in_filter":      []interface{}{},
				"out_filter":     []interface{}{},
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		obj := getPolicyInterVRFRoutingFromSchema(d)
		require.Len(t, obj.BgpRouteLeaking, 2)
		assert.Equal(t, "IPV4", *obj.BgpRouteLeaking[0].AddressFamily)
		assert.Equal(t, "IPV6", *obj.BgpRouteLeaking[1].AddressFamily)
	})

	t.Run("static_route_advertisement with advertisement_rule and in_filter_prefix_list is parsed", func(t *testing.T) {
		data := minimalInterVrfData()
		plPath := "/infra/tier-0s/t0-gw-1/prefix-lists/pl-1"
		data["static_route_advertisement"] = []interface{}{
			map[string]interface{}{
				"advertisement_rule": []interface{}{
					map[string]interface{}{
						"action":                    "PERMIT",
						"name":                      "rule-1",
						"prefix_operator":           "GE",
						"route_advertisement_types": []interface{}{"TIER0_CONNECTED", "TIER0_NAT"},
						"subnets":                   []interface{}{"192.168.1.0/24", "192.168.2.0/24"},
					},
				},
				"in_filter_prefix_list": []interface{}{plPath},
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		obj := getPolicyInterVRFRoutingFromSchema(d)
		require.NotNil(t, obj.StaticRouteAdvertisement)
		require.Len(t, obj.StaticRouteAdvertisement.AdvertisementRules, 1)
		rule := obj.StaticRouteAdvertisement.AdvertisementRules[0]
		assert.Equal(t, "PERMIT", *rule.Action)
		assert.Equal(t, "rule-1", *rule.Name)
		assert.Equal(t, "GE", *rule.PrefixOperator)
		assert.Equal(t, []string{"TIER0_CONNECTED", "TIER0_NAT"}, rule.RouteAdvertisementTypes)
		assert.Equal(t, []string{"192.168.1.0/24", "192.168.2.0/24"}, rule.Subnets)
		assert.Equal(t, []string{plPath}, obj.StaticRouteAdvertisement.InFilterPrefixList)
	})

	t.Run("static_route_advertisement without advertisement_rule is parsed", func(t *testing.T) {
		data := minimalInterVrfData()
		plPath := "/infra/tier-0s/t0-gw-1/prefix-lists/pl-1"
		data["static_route_advertisement"] = []interface{}{
			map[string]interface{}{
				"advertisement_rule":    []interface{}{},
				"in_filter_prefix_list": []interface{}{plPath},
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		obj := getPolicyInterVRFRoutingFromSchema(d)
		require.NotNil(t, obj.StaticRouteAdvertisement)
		assert.Empty(t, obj.StaticRouteAdvertisement.AdvertisementRules)
		assert.Equal(t, []string{plPath}, obj.StaticRouteAdvertisement.InFilterPrefixList)
	})
}
