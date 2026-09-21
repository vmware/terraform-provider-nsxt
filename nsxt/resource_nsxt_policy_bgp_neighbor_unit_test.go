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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func bgpNeighborMinimalData() map[string]interface{} {
	return map[string]interface{}{
		"neighbor_address": "10.0.0.1",
		"remote_as_num":    "65000",
	}
}

func TestUnitNsxt_resourceNsxtPolicyBgpNeighborResourceDataToStruct(t *testing.T) {
	res := resourceNsxtPolicyBgpNeighbor()

	t.Run("basic fields are converted", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, bgpNeighborMinimalData())
		obj, err := resourceNsxtPolicyBgpNeighborResourceDataToStruct(d, "neighbor-1")
		require.NoError(t, err)
		assert.Equal(t, "neighbor-1", *obj.Id)
		assert.Equal(t, "10.0.0.1", *obj.NeighborAddress)
		assert.Equal(t, "65000", *obj.RemoteAsNum)
	})

	t.Run("bfd_config, route_filtering, and local_as_config are converted", func(t *testing.T) {
		data := bgpNeighborMinimalData()
		data["bfd_config"] = []interface{}{
			map[string]interface{}{"enabled": true, "interval": 1000, "multiple": 5},
		}
		data["route_filtering"] = []interface{}{
			map[string]interface{}{
				"address_family":   "IPV4",
				"enabled":          true,
				"in_route_filter":  "/infra/tier-0s/t0/prefix-lists/pl1",
				"out_route_filter": "/infra/tier-0s/t0/prefix-lists/pl2",
				"maximum_routes":   100,
			},
		}
		data["neighbor_local_as_config"] = []interface{}{
			map[string]interface{}{"local_as_num": "65001", "as_path_modifier_type": "NO_PREPEND_REPLACE_AS"},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		obj, err := resourceNsxtPolicyBgpNeighborResourceDataToStruct(d, "neighbor-1")
		require.NoError(t, err)

		require.NotNil(t, obj.Bfd)
		assert.True(t, *obj.Bfd.Enabled)
		assert.EqualValues(t, 1000, *obj.Bfd.Interval)
		assert.EqualValues(t, 5, *obj.Bfd.Multiple)

		require.Len(t, obj.RouteFiltering, 1)
		assert.Equal(t, "IPV4", *obj.RouteFiltering[0].AddressFamily)
		assert.Equal(t, []string{"/infra/tier-0s/t0/prefix-lists/pl1"}, obj.RouteFiltering[0].InRouteFilters)
		assert.Equal(t, []string{"/infra/tier-0s/t0/prefix-lists/pl2"}, obj.RouteFiltering[0].OutRouteFilters)
		require.NotNil(t, obj.RouteFiltering[0].MaximumRoutes)
		assert.EqualValues(t, 100, *obj.RouteFiltering[0].MaximumRoutes)

		require.NotNil(t, obj.NeighborLocalAsConfig)
		assert.Equal(t, "65001", *obj.NeighborLocalAsConfig.LocalAsNum)
		assert.Equal(t, "NO_PREPEND_REPLACE_AS", *obj.NeighborLocalAsConfig.AsPathModifierType)
	})

	t.Run("source_attachment requires NSX 9.2.0 or higher", func(t *testing.T) {
		data := bgpNeighborMinimalData()
		data["source_attachment"] = []interface{}{"/infra/tier-0s/t0/tg-attachments/a1"}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		util.NsxVersion = "9.1.0"
		defer func() { util.NsxVersion = "" }()

		_, err := resourceNsxtPolicyBgpNeighborResourceDataToStruct(d, "neighbor-1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})

	t.Run("source_attachment is included when NSX version supports it", func(t *testing.T) {
		data := bgpNeighborMinimalData()
		data["source_attachment"] = []interface{}{"/infra/tier-0s/t0/tg-attachments/a1"}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		obj, err := resourceNsxtPolicyBgpNeighborResourceDataToStruct(d, "neighbor-1")
		require.NoError(t, err)
		assert.Equal(t, []string{"/infra/tier-0s/t0/tg-attachments/a1"}, obj.SourceAttachment)
	})
}

func TestUnitNsxt_setBgpNeighborConfigInSchema(t *testing.T) {
	res := resourceNsxtPolicyBgpNeighbor()

	t.Run("populates schema including nested bfd/route_filtering/local_as_config", func(t *testing.T) {
		d := res.TestResourceData()
		holdDown := int64(180)
		keepAlive := int64(60)
		maxHop := int64(1)
		neighborAddr := "10.0.0.1"
		remoteAs := "65000"
		gracefulMode := model.BgpNeighborConfig_GRACEFUL_RESTART_MODE_HELPER_ONLY
		bfdEnabled := true
		bfdInterval := int64(1000)
		bfdMultiple := int64(5)
		addrFamily := "IPV4"
		filterEnabled := true
		maxRoutes := int64(100)
		localAsNum := "65001"

		obj := model.BgpNeighborConfig{
			HoldDownTime:        &holdDown,
			KeepAliveTime:       &keepAlive,
			MaximumHopLimit:     &maxHop,
			NeighborAddress:     &neighborAddr,
			RemoteAsNum:         &remoteAs,
			GracefulRestartMode: &gracefulMode,
			Bfd: &model.BgpBfdConfig{
				Enabled:  &bfdEnabled,
				Interval: &bfdInterval,
				Multiple: &bfdMultiple,
			},
			RouteFiltering: []model.BgpRouteFiltering{
				{
					AddressFamily:   &addrFamily,
					Enabled:         &filterEnabled,
					InRouteFilters:  []string{"/infra/tier-0s/t0/prefix-lists/pl1"},
					OutRouteFilters: []string{"/infra/tier-0s/t0/prefix-lists/pl2"},
					MaximumRoutes:   &maxRoutes,
				},
			},
			NeighborLocalAsConfig: &model.BgpNeighborLocalAsConfig{
				LocalAsNum: &localAsNum,
			},
		}

		setBgpNeighborConfigInSchema(d, obj, false)

		assert.Equal(t, neighborAddr, d.Get("neighbor_address"))
		assert.Equal(t, 180, d.Get("hold_down_time"))

		bfd := d.Get("bfd_config").([]interface{})
		require.Len(t, bfd, 1)
		bfdMap := bfd[0].(map[string]interface{})
		assert.Equal(t, true, bfdMap["enabled"])
		assert.Equal(t, 1000, bfdMap["interval"])

		filters := d.Get("route_filtering").([]interface{})
		require.Len(t, filters, 1)
		filterMap := filters[0].(map[string]interface{})
		assert.Equal(t, "/infra/tier-0s/t0/prefix-lists/pl1", filterMap["in_route_filter"])
		assert.Equal(t, "/infra/tier-0s/t0/prefix-lists/pl2", filterMap["out_route_filter"])
		assert.Equal(t, 100, filterMap["maximum_routes"])

		localAsConfigs := d.Get("neighbor_local_as_config").([]interface{})
		require.Len(t, localAsConfigs, 1)
		assert.Equal(t, localAsNum, localAsConfigs[0].(map[string]interface{})["local_as_num"])
	})

	t.Run("nil bfd and local_as_config leave those schema fields empty", func(t *testing.T) {
		d := res.TestResourceData()
		holdDown := int64(180)
		keepAlive := int64(60)
		maxHop := int64(1)
		neighborAddr := "10.0.0.1"
		remoteAs := "65000"

		obj := model.BgpNeighborConfig{
			HoldDownTime:    &holdDown,
			KeepAliveTime:   &keepAlive,
			MaximumHopLimit: &maxHop,
			NeighborAddress: &neighborAddr,
			RemoteAsNum:     &remoteAs,
		}

		setBgpNeighborConfigInSchema(d, obj, false)

		assert.Empty(t, d.Get("bfd_config").([]interface{}))
		assert.Empty(t, d.Get("neighbor_local_as_config").([]interface{}))
	})
}
