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
)

func rcBgpNeighborMinimalData() map[string]interface{} {
	return map[string]interface{}{
		"neighbor_address": "10.0.0.1",
		"remote_as_num":    "65000",
	}
}

func TestUnitNsxt_rcBgpNeighborToStruct(t *testing.T) {
	res := resourceNsxtPolicyRouteControllerBgpNeighbor()

	t.Run("basic fields are converted", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, rcBgpNeighborMinimalData())
		obj, err := rcBgpNeighborToStruct(d, "neighbor-1")
		require.NoError(t, err)
		assert.Equal(t, "neighbor-1", *obj.Id)
		assert.Equal(t, "10.0.0.1", *obj.NeighborAddress)
		assert.Equal(t, "65000", *obj.RemoteAsNum)
		assert.Empty(t, obj.GatewayIps)
	})

	t.Run("bfd_config, route_filtering, and gateway_ips are converted", func(t *testing.T) {
		data := rcBgpNeighborMinimalData()
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
		data["gateway_ips"] = []interface{}{"10.0.0.2"}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		obj, err := rcBgpNeighborToStruct(d, "neighbor-1")
		require.NoError(t, err)

		require.NotNil(t, obj.Bfd)
		assert.True(t, *obj.Bfd.Enabled)
		assert.EqualValues(t, 1000, *obj.Bfd.Interval)

		require.Len(t, obj.RouteFiltering, 1)
		assert.Equal(t, "IPV4", *obj.RouteFiltering[0].AddressFamily)
		assert.Equal(t, []string{"/infra/tier-0s/t0/prefix-lists/pl1"}, obj.RouteFiltering[0].InRouteFilters)
		require.NotNil(t, obj.RouteFiltering[0].MaximumRoutes)
		assert.EqualValues(t, 100, *obj.RouteFiltering[0].MaximumRoutes)

		assert.Equal(t, []string{"10.0.0.2"}, obj.GatewayIps)
	})
}

func TestUnitNsxt_setRCBgpNeighborConfigInSchema(t *testing.T) {
	res := resourceNsxtPolicyRouteControllerBgpNeighbor()

	t.Run("populates schema including nested bfd/route_filtering", func(t *testing.T) {
		d := res.TestResourceData()
		holdDown := int64(180)
		keepAlive := int64(60)
		maxHop := int64(1)
		neighborAddr := "10.0.0.1"
		remoteAs := "65000"
		bfdEnabled := true
		bfdInterval := int64(1000)
		bfdMultiple := int64(5)
		addrFamily := "IPV4"
		filterEnabled := true
		maxRoutes := int64(100)

		obj := model.RouteControllerBgpNeighborConfig{
			HoldDownTime:    &holdDown,
			KeepAliveTime:   &keepAlive,
			MaximumHopLimit: &maxHop,
			NeighborAddress: &neighborAddr,
			RemoteAsNum:     &remoteAs,
			GatewayIps:      []string{"10.0.0.2"},
			Bfd: &model.BgpBfdConfig{
				Enabled:  &bfdEnabled,
				Interval: &bfdInterval,
				Multiple: &bfdMultiple,
			},
			RouteFiltering: []model.RouteControllerBgpRouteFiltering{
				{
					AddressFamily:   &addrFamily,
					Enabled:         &filterEnabled,
					InRouteFilters:  []string{"/infra/tier-0s/t0/prefix-lists/pl1"},
					OutRouteFilters: []string{"/infra/tier-0s/t0/prefix-lists/pl2"},
					MaximumRoutes:   &maxRoutes,
				},
			},
		}

		setRCBgpNeighborConfigInSchema(d, obj)

		assert.Equal(t, neighborAddr, d.Get("neighbor_address"))
		assert.Equal(t, 180, d.Get("hold_down_time"))
		assert.Equal(t, []interface{}{"10.0.0.2"}, d.Get("gateway_ips"))

		bfd := d.Get("bfd_config").([]interface{})
		require.Len(t, bfd, 1)
		bfdMap := bfd[0].(map[string]interface{})
		assert.Equal(t, true, bfdMap["enabled"])
		assert.Equal(t, 1000, bfdMap["interval"])

		filters := d.Get("route_filtering").([]interface{})
		require.Len(t, filters, 1)
		filterMap := filters[0].(map[string]interface{})
		assert.Equal(t, "/infra/tier-0s/t0/prefix-lists/pl1", filterMap["in_route_filter"])
		assert.Equal(t, 100, filterMap["maximum_routes"])
	})

	t.Run("nil bfd leaves bfd_config empty", func(t *testing.T) {
		d := res.TestResourceData()
		holdDown := int64(180)
		keepAlive := int64(60)
		maxHop := int64(1)
		neighborAddr := "10.0.0.1"
		remoteAs := "65000"

		obj := model.RouteControllerBgpNeighborConfig{
			HoldDownTime:    &holdDown,
			KeepAliveTime:   &keepAlive,
			MaximumHopLimit: &maxHop,
			NeighborAddress: &neighborAddr,
			RemoteAsNum:     &remoteAs,
		}

		setRCBgpNeighborConfigInSchema(d, obj)

		assert.Empty(t, d.Get("bfd_config").([]interface{}))
	})
}
