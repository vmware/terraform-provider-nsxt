//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func testDHCPBindingSegmentVariants(t *testing.T, v6 bool) {
	t.Helper()

	res := resourceNsxtPolicyDhcpV4StaticBinding()
	kind := "DhcpV4StaticBindingConfig"
	if v6 {
		res = resourceNsxtPolicyDhcpV6StaticBinding()
		kind = "DhcpV6StaticBindingConfig"
	}
	for _, fixed := range []bool{false, true} {
		t.Run(fmt.Sprintf("ipv6_%v_fixed_%v", v6, fixed), func(t *testing.T) {
			segmentPath := "/infra/segments/segment"
			if fixed {
				segmentPath = "/infra/tier-1s/gateway/segments/segment"
			}
			raw := map[string]interface{}{"segment_path": segmentPath, "mac_address": "00:11:22:33:44:55", "lease_time": 3600}
			if v6 {
				raw["ip_addresses"] = []interface{}{"2001:db8::10"}
				raw["domain_names"] = []interface{}{"example.com"}
				raw["dns_nameservers"] = []interface{}{"2001:db8::53"}
				raw["sntp_servers"] = []interface{}{"2001:db8::123"}
				raw["preferred_time"] = 1800
			} else {
				raw["ip_address"] = "192.0.2.10"
				raw["gateway_address"] = "192.0.2.1"
			}
			d := schema.TestResourceDataRaw(t, res.Schema, raw)
			d.SetId("binding")
			m := newIsolatedMockProviderClientCacheEnabled(t)
			m.Host = "https://nsx.invalid"
			methods := []string{}
			m.PolicyHTTPClient = &http.Client{Transport: unitPolicyTransport(func(r *http.Request) (*http.Response, error) {
				methods = append(methods, r.Method)
				assert.Contains(t, r.URL.Path, segmentPath+"/dhcp-static-binding-configs/binding")
				if r.Method == http.MethodPatch {
					var payload map[string]interface{}
					require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
					assert.Equal(t, kind, payload["resource_type"])
					assert.Equal(t, "00:11:22:33:44:55", payload["mac_address"])
				}
				body := `{"resource_type":"` + kind + `","id":"binding","mac_address":"00:11:22:33:44:55","ip_address":"192.0.2.10","ip_addresses":["2001:db8::10"],"lease_time":3600,"preferred_time":1800,"domain_names":["example.com"],"dns_nameservers":["2001:db8::53"],"sntp_servers":["2001:db8::123"]}`
				return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			})}
			var err error
			if v6 {
				err = policyDhcpV6StaticBindingConvertAndPatch(d, segmentPath, "binding", m)
			} else {
				err = policyDhcpV4StaticBindingConvertAndPatch(d, segmentPath, "binding", m)
			}
			require.NoError(t, err)
			typ := resourceTypeDhcpV4StaticBindingConfig
			key := CacheKeyForResourceID(typ, d)
			if v6 {
				typ = resourceTypeDhcpV6StaticBindingConfig
				key = d.Id()
			}
			MarkPostWriteForResourceTypeKey(typ, key, m)
			t.Cleanup(func() { postWriteByKey.Delete(postWriteKey(typ, key)) })
			require.NoError(t, res.Read(d, m))
			assert.Equal(t, []string{"PATCH", "GET", "PATCH"}, methods)
			assert.Equal(t, "00:11:22:33:44:55", d.Get("mac_address"))
			if v6 {
				assert.Equal(t, []interface{}{"example.com"}, d.Get("domain_names"))
				assert.Equal(t, 1800, d.Get("preferred_time"))
			} else {
				assert.Equal(t, "192.0.2.10", d.Get("ip_address"))
			}
		})
	}
	for _, global := range []bool{false, true} {
		t.Run(fmt.Sprintf("ipv6_%v_rejected_global_%v", v6, global), func(t *testing.T) {
			path := "/infra/tier-0s/gateway/segments/segment"
			if global {
				path = "/infra/tier-1s/gateway/segments/segment"
			}
			d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"segment_path": path})
			m := newIsolatedMockProviderClient(t)
			m.PolicyGlobalManager = global
			var err error
			if v6 {
				err = policyDhcpV6StaticBindingConvertAndPatch(d, path, "binding", m)
			} else {
				err = policyDhcpV4StaticBindingConvertAndPatch(d, path, "binding", m)
			}
			require.ErrorContains(t, err, "not applicable")
		})
	}
}
