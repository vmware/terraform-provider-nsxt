//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestUnitNsxt_getPolicyDhcpOptions121(t *testing.T) {
	t.Run("no routes returns zero value", func(t *testing.T) {
		opt := getPolicyDhcpOptions121(nil)
		assert.Nil(t, opt.StaticRoutes)
	})

	t.Run("builds static routes", func(t *testing.T) {
		opt := getPolicyDhcpOptions121([]interface{}{
			map[string]interface{}{"network": "10.0.0.0/24", "next_hop": "10.0.0.1"},
		})
		require.Len(t, opt.StaticRoutes, 1)
		assert.Equal(t, "10.0.0.0/24", *opt.StaticRoutes[0].Network)
		assert.Equal(t, "10.0.0.1", *opt.StaticRoutes[0].NextHop)
	})
}

func TestUnitNsxt_getPolicyDhcpOptions121FromStruct(t *testing.T) {
	network := "10.0.0.0/24"
	nextHop := "10.0.0.1"
	opt := &model.DhcpOption121{StaticRoutes: []model.ClasslessStaticRoute{{Network: &network, NextHop: &nextHop}}}

	out := getPolicyDhcpOptions121FromStruct(opt)
	require.Len(t, out, 1)
	assert.Equal(t, &network, out[0]["network"])
	assert.Equal(t, &nextHop, out[0]["next_hop"])
}

func TestUnitNsxt_getPolicyDhcpGenericOptions(t *testing.T) {
	opts := getPolicyDhcpGenericOptions([]interface{}{
		map[string]interface{}{"code": 66, "values": []interface{}{"v1", "v2"}},
	})
	require.Len(t, opts, 1)
	assert.EqualValues(t, 66, *opts[0].Code)
	assert.Equal(t, []string{"v1", "v2"}, opts[0].Values)
}

func TestUnitNsxt_getPolicyDhcpGenericOptionsFromStruct(t *testing.T) {
	code := int64(66)
	out := getPolicyDhcpGenericOptionsFromStruct([]model.GenericDhcpOption{
		{Code: &code, Values: []string{"v1"}},
	})
	require.Len(t, out, 1)
	assert.Equal(t, &code, out[0]["code"])
	assert.Equal(t, []string{"v1"}, out[0]["values"])
}

func TestUnitNsxt_getDhcpOptsFromMap(t *testing.T) {
	t.Run("no options returns nil", func(t *testing.T) {
		out := getDhcpOptsFromMap(map[string]interface{}{
			"dhcp_option_121":     []interface{}{},
			"dhcp_generic_option": []interface{}{},
		})
		assert.Nil(t, out)
	})

	t.Run("combines option_121 and generic options", func(t *testing.T) {
		out := getDhcpOptsFromMap(map[string]interface{}{
			"dhcp_option_121": []interface{}{
				map[string]interface{}{"network": "10.0.0.0/24", "next_hop": "10.0.0.1"},
			},
			"dhcp_generic_option": []interface{}{
				map[string]interface{}{"code": 66, "values": []interface{}{"v1"}},
			},
		})
		require.NotNil(t, out)
		require.NotNil(t, out.Option121)
		require.Len(t, out.Option121.StaticRoutes, 1)
		require.Len(t, out.Others, 1)
	})
}

func TestUnitNsxt_getOldProfileDataForRemoval(t *testing.T) {
	t.Run("nil input returns zero values", func(t *testing.T) {
		id, rev := getOldProfileDataForRemoval(nil)
		assert.Equal(t, "", id)
		assert.EqualValues(t, 0, rev)
	})

	t.Run("empty list returns zero values", func(t *testing.T) {
		id, rev := getOldProfileDataForRemoval([]interface{}{})
		assert.Equal(t, "", id)
		assert.EqualValues(t, 0, rev)
	})

	t.Run("extracts id and revision", func(t *testing.T) {
		id, rev := getOldProfileDataForRemoval([]interface{}{
			map[string]interface{}{
				"binding_map_path": "/infra/segments/seg1/segment-discovery-profile-binding-maps/map1",
				"revision":         3,
			},
		})
		assert.Equal(t, "map1", id)
		assert.EqualValues(t, 3, rev)
	})
}

func TestUnitNsxt_getSegmentSubnetDhcpConfigFromSchema(t *testing.T) {
	util.NsxVersion = "3.1.0"
	defer func() { util.NsxVersion = "" }()

	t.Run("version too low returns nil", func(t *testing.T) {
		util.NsxVersion = "2.5.0"
		defer func() { util.NsxVersion = "3.1.0" }()

		out, err := getSegmentSubnetDhcpConfigFromSchema(map[string]interface{}{
			"dhcp_v4_config": []interface{}{},
			"dhcp_v6_config": []interface{}{},
		})
		require.NoError(t, err)
		assert.Nil(t, out)
	})

	t.Run("both v4 and v6 configured errors", func(t *testing.T) {
		_, err := getSegmentSubnetDhcpConfigFromSchema(map[string]interface{}{
			"dhcp_v4_config": []interface{}{map[string]interface{}{}},
			"dhcp_v6_config": []interface{}{map[string]interface{}{}},
		})
		require.Error(t, err)
	})

	t.Run("neither configured returns nil", func(t *testing.T) {
		out, err := getSegmentSubnetDhcpConfigFromSchema(map[string]interface{}{
			"dhcp_v4_config": []interface{}{},
			"dhcp_v6_config": []interface{}{},
		})
		require.NoError(t, err)
		assert.Nil(t, out)
	})

	t.Run("v4 config builds a StructValue", func(t *testing.T) {
		out, err := getSegmentSubnetDhcpConfigFromSchema(map[string]interface{}{
			"dhcp_v4_config": []interface{}{
				map[string]interface{}{
					"server_address":      "10.0.0.1/24",
					"dns_servers":         []interface{}{"8.8.8.8"},
					"lease_time":          3600,
					"dhcp_option_121":     []interface{}{},
					"dhcp_generic_option": []interface{}{},
				},
			},
			"dhcp_v6_config": []interface{}{},
		})
		require.NoError(t, err)
		require.NotNil(t, out)

		converter := bindings.NewTypeConverter()
		golang, errs := converter.ConvertToGolang(out, model.SegmentDhcpV4ConfigBindingType())
		require.Empty(t, errs)
		v4 := golang.(model.SegmentDhcpV4Config)
		assert.Equal(t, "10.0.0.1/24", *v4.ServerAddress)
		assert.EqualValues(t, 3600, *v4.LeaseTime)
	})

	t.Run("v6 config builds a StructValue", func(t *testing.T) {
		out, err := getSegmentSubnetDhcpConfigFromSchema(map[string]interface{}{
			"dhcp_v4_config": []interface{}{},
			"dhcp_v6_config": []interface{}{
				map[string]interface{}{
					"server_address": "fe80::1/64",
					"dns_servers":    []interface{}{"2001:4860:4860::8888"},
					"sntp_servers":   []interface{}{},
					"domain_names":   []interface{}{},
					"excluded_range": []interface{}{
						map[string]interface{}{"start": "fe80::10", "end": "fe80::20"},
					},
					"lease_time":     3600,
					"preferred_time": 1800,
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, out)

		converter := bindings.NewTypeConverter()
		golang, errs := converter.ConvertToGolang(out, model.SegmentDhcpV6ConfigBindingType())
		require.Empty(t, errs)
		v6 := golang.(model.SegmentDhcpV6Config)
		assert.Equal(t, "fe80::1/64", *v6.ServerAddress)
		require.Len(t, v6.ExcludedRanges, 1)
		assert.Equal(t, "fe80::10-fe80::20", v6.ExcludedRanges[0])
	})
}

func TestUnitNsxt_setSegmentSubnetDhcpConfigInSchema(t *testing.T) {
	t.Run("nil DhcpConfig is a no-op", func(t *testing.T) {
		schemaConfig := map[string]interface{}{}
		err := setSegmentSubnetDhcpConfigInSchema(schemaConfig, model.SegmentSubnet{})
		require.NoError(t, err)
		assert.Empty(t, schemaConfig)
	})

	t.Run("round trips a v4 config", func(t *testing.T) {
		serverAddress := "10.0.0.1/24"
		leaseTime := int64(3600)
		converter := bindings.NewTypeConverter()
		v4Config := model.SegmentDhcpV4Config{
			ResourceType:  model.SegmentDhcpConfig_RESOURCE_TYPE_SEGMENTDHCPV4CONFIG,
			ServerAddress: &serverAddress,
			LeaseTime:     &leaseTime,
			DnsServers:    []string{"8.8.8.8"},
		}
		sv, errs := converter.ConvertToVapi(v4Config, model.SegmentDhcpV4ConfigBindingType())
		require.Empty(t, errs)

		schemaConfig := map[string]interface{}{}
		err := setSegmentSubnetDhcpConfigInSchema(schemaConfig, model.SegmentSubnet{DhcpConfig: sv.(*data.StructValue)})
		require.NoError(t, err)

		v4Result := schemaConfig["dhcp_v4_config"].([]map[string]interface{})
		require.Len(t, v4Result, 1)
		assert.Equal(t, &serverAddress, v4Result[0]["server_address"])
		assert.Equal(t, &leaseTime, v4Result[0]["lease_time"])
	})

	t.Run("round trips a v6 config with excluded ranges", func(t *testing.T) {
		serverAddress := "fe80::1/64"
		leaseTime := int64(3600)
		preferredTime := int64(1800)
		converter := bindings.NewTypeConverter()
		v6Config := model.SegmentDhcpV6Config{
			ResourceType:   model.SegmentDhcpConfig_RESOURCE_TYPE_SEGMENTDHCPV6CONFIG,
			ServerAddress:  &serverAddress,
			LeaseTime:      &leaseTime,
			PreferredTime:  &preferredTime,
			DnsServers:     []string{"2001:4860:4860::8888"},
			ExcludedRanges: []string{"fe80::10-fe80::20"},
		}
		sv, errs := converter.ConvertToVapi(v6Config, model.SegmentDhcpV6ConfigBindingType())
		require.Empty(t, errs)

		schemaConfig := map[string]interface{}{}
		err := setSegmentSubnetDhcpConfigInSchema(schemaConfig, model.SegmentSubnet{DhcpConfig: sv.(*data.StructValue)})
		require.NoError(t, err)

		v6Result := schemaConfig["dhcp_v6_config"].([]map[string]interface{})
		require.Len(t, v6Result, 1)
		assert.Equal(t, &serverAddress, v6Result[0]["server_address"])
		assert.Equal(t, &preferredTime, v6Result[0]["preferred_time"])
		excludedRanges := v6Result[0]["excluded_range"].([]map[string]interface{})
		require.Len(t, excludedRanges, 1)
		assert.Equal(t, "fe80::10", excludedRanges[0]["start"])
		assert.Equal(t, "fe80::20", excludedRanges[0]["end"])
	})

	t.Run("unrecognized resource type errors", func(t *testing.T) {
		converter := bindings.NewTypeConverter()
		badConfig := model.SegmentDhcpConfig{ResourceType: "SomeUnknownType"}
		sv, errs := converter.ConvertToVapi(badConfig, model.SegmentDhcpConfigBindingType())
		require.Empty(t, errs)

		schemaConfig := map[string]interface{}{}
		err := setSegmentSubnetDhcpConfigInSchema(schemaConfig, model.SegmentSubnet{DhcpConfig: sv.(*data.StructValue)})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Unrecognized DHCP Config Resource Type")
	})
}
