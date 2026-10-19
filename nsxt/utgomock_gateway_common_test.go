//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestUnitNsxt_parseGatewayInterfacePolicyPath(t *testing.T) {
	// First return value is isTier0, not a generic "ok" flag.
	isT0, gw, ls, ifID := parseGatewayInterfacePolicyPath("/infra/tier-0s/gw1/locale-services/ls1/interfaces/if1")
	assert.True(t, isT0)
	assert.Equal(t, "gw1", gw)
	assert.Equal(t, "ls1", ls)
	assert.Equal(t, "if1", ifID)

	isT0, gw, ls, ifID = parseGatewayInterfacePolicyPath("/infra/tier-1s/t1/locale-services/ls2/interfaces/if2")
	assert.False(t, isT0)
	assert.Equal(t, "t1", gw)
	assert.Equal(t, "ls2", ls)
	assert.Equal(t, "if2", ifID)

	isT0, gw, ls, ifID = parseGatewayInterfacePolicyPath("/infra/tier-0s/gw")
	assert.False(t, isT0)
	assert.Empty(t, gw)

	isT0, _, _, _ = parseGatewayInterfacePolicyPath("/infra/tier-0s/gw/locale-services/ls/interfaces/extra/seg")
	assert.False(t, isT0)
}

func TestUnitNsxt_getGlobalPolicyGatewayLocaleServiceIDWithSite(t *testing.T) {
	_, err := getGlobalPolicyGatewayLocaleServiceIDWithSite(nil, "/site/x", "gw1")
	require.Error(t, err)

	ec := "/global-infra/sites/loc1/edge-clusters/ec1"
	id1 := "ls-1"
	lsID, err := getGlobalPolicyGatewayLocaleServiceIDWithSite([]model.LocaleServices{
		{Id: &id1, EdgeClusterPath: &ec},
	}, "/global-infra/sites/loc1", "gw1")
	require.NoError(t, err)
	assert.Equal(t, "ls-1", lsID)

	ecOther := "/global-infra/sites/other/edge-clusters/ec1"
	_, err = getGlobalPolicyGatewayLocaleServiceIDWithSite([]model.LocaleServices{
		{Id: &id1, EdgeClusterPath: &ecOther},
	}, "/global-infra/sites/loc1", "gw1")
	require.Error(t, err)
}

func TestUnitNsxt_listPolicyGatewayLocaleServices_pagination(t *testing.T) {
	rc1 := int64(2)
	rc2 := int64(2)
	ls1, ls2 := "a", "b"
	calls := 0
	listFn := func(_ utl.SessionContext, _ client.Connector, _ string, cursor *string) (model.LocaleServicesListResult, error) {
		calls++
		if calls == 1 {
			cur := "next"
			return model.LocaleServicesListResult{
				Results:     []model.LocaleServices{{Id: &ls1}},
				ResultCount: &rc1,
				Cursor:      &cur,
			}, nil
		}
		return model.LocaleServicesListResult{
			Results:     []model.LocaleServices{{Id: &ls2}},
			ResultCount: &rc2,
			Cursor:      nil,
		}, nil
	}
	out, err := listPolicyGatewayLocaleServices(utl.SessionContext{}, nil, "gw", listFn)
	require.NoError(t, err)
	require.Len(t, out, 2)
	assert.Equal(t, "a", *out[0].Id)
	assert.Equal(t, "b", *out[1].Id)
	assert.Equal(t, 2, calls)
}

func TestUnitNsxt_listPolicyGatewayLocaleServices_error(t *testing.T) {
	listFn := func(_ utl.SessionContext, _ client.Connector, _ string, _ *string) (model.LocaleServicesListResult, error) {
		return model.LocaleServicesListResult{}, errors.New("list boom")
	}
	_, err := listPolicyGatewayLocaleServices(utl.SessionContext{}, nil, "gw", listFn)
	require.Error(t, err)
}

func TestUnitNsxt_initChildLocaleService(t *testing.T) {
	id := "ls1"
	sv, err := initChildLocaleService(&model.LocaleServices{Id: &id}, true)
	require.NoError(t, err)
	require.NotNil(t, sv)
}

func TestUnitNsxt_setIpv6ProfilePathsInSchema(t *testing.T) {
	sch := map[string]*schema.Schema{
		"ipv6_ndra_profile_path": {Type: schema.TypeString, Optional: true},
		"ipv6_dad_profile_path":  {Type: schema.TypeString, Optional: true},
	}
	d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{})
	require.NoError(t, setIpv6ProfilePathsInSchema(d, []string{
		"/infra/ipv6-ndra-profiles/p1",
		"/infra/ipv6-dad-profiles/p2",
	}))
	assert.Equal(t, "/infra/ipv6-ndra-profiles/p1", d.Get("ipv6_ndra_profile_path"))
	assert.Equal(t, "/infra/ipv6-dad-profiles/p2", d.Get("ipv6_dad_profile_path"))

	d2 := schema.TestResourceDataRaw(t, sch, map[string]interface{}{})
	mtPath := "/orgs/o1/projects/p1/infra/ipv6-ndra-profiles/p9"
	require.NoError(t, setIpv6ProfilePathsInSchema(d2, []string{mtPath}))
	assert.Equal(t, mtPath, d2.Get("ipv6_ndra_profile_path"))
}

func TestUnitNsxt_getIpv6ProfilePathsFromSchema(t *testing.T) {
	sch := map[string]*schema.Schema{
		"ipv6_ndra_profile_path": {Type: schema.TypeString, Optional: true},
		"ipv6_dad_profile_path":  {Type: schema.TypeString, Optional: true},
	}
	d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{
		"ipv6_ndra_profile_path": "/ndra",
		"ipv6_dad_profile_path":  "/dad",
	})
	out := getIpv6ProfilePathsFromSchema(d)
	assert.Equal(t, []string{"/ndra", "/dad"}, out)
}

func TestUnitNsxt_gatewayIntersiteConfigSchemaRoundTrip(t *testing.T) {
	sch := map[string]*schema.Schema{
		"intersite_config": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"transit_subnet": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"primary_site_path": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"fallback_site_paths": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
	}
	d := schema.TestResourceDataRaw(t, sch, map[string]interface{}{})
	require.NoError(t, d.Set("intersite_config", []interface{}{
		map[string]interface{}{
			"transit_subnet":      "10.0.0.0/24",
			"primary_site_path":   "/global-infra/sites/s1",
			"fallback_site_paths": []interface{}{"/global-infra/sites/s2"},
		},
	}))
	cfg := getPolicyGatewayIntersiteConfigFromSchema(d)
	require.NotNil(t, cfg)
	assert.Equal(t, "10.0.0.0/24", *cfg.IntersiteTransitSubnet)
	assert.Equal(t, "/global-infra/sites/s1", *cfg.PrimarySitePath)
	require.Len(t, cfg.FallbackSites, 1)

	d2 := schema.TestResourceDataRaw(t, sch, map[string]interface{}{})
	require.NoError(t, setPolicyGatewayIntersiteConfigInSchema(d2, cfg))
	lst := d2.Get("intersite_config").([]interface{})
	require.Len(t, lst, 1)
}

func TestUnitNsxt_setLocaleServiceRedistributionRulesConfig(t *testing.T) {
	rule := map[string]interface{}{
		"name":           "r1",
		"route_map_path": "/infra/rm1",
		"types":          schema.NewSet(schema.HashString, []interface{}{"TIER0_STATIC"}),
		"bgp":            true,
		"ospf":           false,
	}
	var cfg model.Tier0RouteRedistributionConfig
	setLocaleServiceRedistributionRulesConfig([]interface{}{rule}, &cfg)
	require.NotNil(t, cfg.RedistributionRules)
	assert.GreaterOrEqual(t, len(cfg.RedistributionRules), 1)
}

func TestUnitNsxt_getLocaleServiceRedistributionConfig_empty(t *testing.T) {
	out := getLocaleServiceRedistributionConfig(&model.LocaleServices{})
	assert.Empty(t, out)
}

func TestUnitNsxt_LocaleServiceRedistributionSchema(t *testing.T) {
	assert.Empty(t, getLocaleServiceRedistributionConfig(&model.LocaleServices{}))
	config := &model.Tier0RouteRedistributionConfig{BgpEnabled: boolPtr(true), OspfEnabled: boolPtr(false), RedistributionRules: []model.Tier0RouteRedistributionRule{
		{Name: strPtr("both"), RouteMapPath: strPtr("/infra/route-maps/map"), RouteRedistributionTypes: []string{"TIER0_CONNECTED"}, Destinations: []string{"BGP", "OSPF"}},
		{Name: strPtr("neither")},
	}}
	result := getLocaleServiceRedistributionConfig(&model.LocaleServices{RouteRedistributionConfig: config})
	require.Len(t, result, 1)
	assert.Equal(t, config.BgpEnabled, result[0]["enabled"])
	assert.Equal(t, config.OspfEnabled, result[0]["ospf_enabled"])
	rules := result[0]["rule"].([]map[string]interface{})
	require.Len(t, rules, 2)
	assert.Equal(t, true, rules[0]["bgp"])
	assert.Equal(t, true, rules[0]["ospf"])
	assert.Equal(t, false, rules[1]["bgp"])
	assert.Equal(t, false, rules[1]["ospf"])
	assert.Equal(t, config.RedistributionRules[0].RouteMapPath, rules[0]["route_map_path"])
	assert.Equal(t, []string{"TIER0_CONNECTED"}, rules[0]["types"])
}

func TestUnitNsxt_GatewayLocaleServiceRemovalAndIdentity(t *testing.T) {
	converter := bindings.NewTypeConverter()
	for _, tc := range []struct {
		name        string
		services    []interface{}
		existing    []model.LocaleServices
		wantError   string
		wantDeleted int
	}{
		{"remove_obsolete", []interface{}{map[string]interface{}{"nsx_id": "keep", "edge_cluster_path": "/infra/edge-clusters/keep"}}, []model.LocaleServices{{Id: strPtr("keep")}, {Id: strPtr("remove")}}, "", 1},
		{"reuse_single_service", []interface{}{map[string]interface{}{"edge_cluster_path": "/infra/edge-clusters/keep"}}, []model.LocaleServices{{Id: strPtr("keep")}}, "", 0},
		{"generate_service_id", []interface{}{map[string]interface{}{"edge_cluster_path": "/infra/edge-clusters/new"}}, nil, "", 0},
		{"duplicate_ids", []interface{}{map[string]interface{}{"nsx_id": "duplicate", "edge_cluster_path": "/infra/edge-clusters/one"}, map[string]interface{}{"nsx_id": "duplicate", "edge_cluster_path": "/infra/edge-clusters/two"}}, nil, "Duplicate nsx_id", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceNsxtPolicyTier0Gateway().Schema, map[string]interface{}{"locale_service": tc.services})
			d.SetId("gateway")
			children, err := initGatewayLocaleServices(utl.SessionContext{ClientType: utl.Local}, d, nil, func(_ utl.SessionContext, _ client.Connector, id string) ([]model.LocaleServices, error) {
				assert.Equal(t, "gateway", id)
				return tc.existing, nil
			})
			if tc.wantError != "" {
				require.ErrorContains(t, err, tc.wantError)
				return
			}
			require.NoError(t, err)
			deleted := 0
			active := 0
			for _, child := range children {
				value, errs := converter.ConvertToGolang(child, model.ChildLocaleServicesBindingType())
				require.Empty(t, errs)
				ls := value.(model.ChildLocaleServices)
				require.NotNil(t, ls.LocaleServices)
				require.NotNil(t, ls.LocaleServices.Id)
				if ls.MarkedForDelete != nil && *ls.MarkedForDelete {
					deleted++
					assert.Equal(t, "remove", *ls.LocaleServices.Id)
				} else {
					active++
					assert.NotEmpty(t, *ls.LocaleServices.Id)
					if len(tc.existing) > 0 {
						assert.Equal(t, "keep", *ls.LocaleServices.Id)
					}
				}
			}
			assert.Equal(t, tc.wantDeleted, deleted)
			assert.Equal(t, 1, active)
		})
	}
	t.Run("list_error", func(t *testing.T) {
		d := resourceNsxtPolicyTier0Gateway().TestResourceData()
		d.SetId("gateway")
		expected := errors.New("list unavailable")
		_, err := initGatewayLocaleServices(utl.SessionContext{}, d, nil, func(utl.SessionContext, client.Connector, string) ([]model.LocaleServices, error) {
			return nil, expected
		})
		require.ErrorIs(t, err, expected)
	})
}

func testGatewayOptionalReadFields(t *testing.T, resource *schema.Resource, tier string, iface bool) {
	oldVersion := util.NsxVersion
	t.Cleanup(func() { util.NsxVersion = oldVersion })
	util.NsxVersion = "9.2.0"

	for _, global := range []bool{false, true} {
		t.Run(fmt.Sprintf("global_%v", global), func(t *testing.T) {
			res := resource
			raw := map[string]interface{}{}
			if iface {
				raw["gateway_path"] = "/infra/" + tier + "/gateway"
				raw["locale_service_id"] = "locale"
			} else {
				raw["locale_service"] = []interface{}{map[string]interface{}{"nsx_id": "locale", "edge_cluster_path": "/infra/sites/site/enforcement-points/default/edge-clusters/cluster"}}
			}
			d := schema.TestResourceDataRaw(t, res.Schema, raw)
			d.SetId("gateway")
			m := newIsolatedMockProviderClient(t)
			m.PolicyGlobalManager = global
			m.Host = "https://nsx.invalid"
			m.CommonConfig.MaxRetries = 0
			requests := 0
			m.PolicyHTTPClient = &http.Client{Transport: unitPolicyTransport(func(r *http.Request) (*http.Response, error) {
				requests++
				body := `{"id":"gateway","display_name":"Gateway","path":"/infra/` + tier + `/gateway","disable_firewall":false,"type":"EXTERNAL","ha_mode":"ACTIVE_STANDBY","pool_allocation":"ROUTING","dhcp_config_paths":["/infra/dhcp-server-configs/dhcp"],"qos_profile":{"ingress_qos_profile_path":"/infra/qos-profiles/in","egress_qos_profile_path":"/infra/qos-profiles/out"},"intersite_config":{"primary_site_path":"/infra/sites/site","intersite_transit_subnet":"100.64.0.0/16"},"ipv6_profile_paths":["/infra/ipv6-ndra-profiles/profile"],"mtu":1600,"multicast":{"enabled":true},"access_vlan_id":123,"edge_cluster_member_index":1,"urpf_mode":"LOOSE","subnets":[{"ip_addresses":["192.0.2.1"],"prefix_len":24}]}`
				service := `{"id":"locale","path":"/infra/` + tier + `/gateway/locale-services/locale","edge_cluster_path":"/infra/sites/site/enforcement-points/default/edge-clusters/cluster","preferred_edge_paths":["/infra/sites/site/enforcement-points/default/edge-clusters/cluster/edge-nodes/1"],"revision":2}`
				if strings.HasSuffix(r.URL.Path, "/locale-services") {
					body = `{"results":[` + service + `],"result_count":1}`
				} else if strings.HasSuffix(r.URL.Path, "/locale-services/locale") {
					body = service
				} else if strings.HasSuffix(r.URL.Path, "/bgp") {
					body = `{}`
				}
				return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			})}
			require.NoError(t, res.Read(d, m))
			require.Positive(t, requests)
			assert.Equal(t, "Gateway", d.Get("display_name"))
			assert.Equal(t, "gateway", d.Id())
			if iface {
				assert.Equal(t, 1600, d.Get("mtu"))
				assert.Equal(t, "LOOSE", d.Get("urpf_mode"))
				assert.Contains(t, d.Get("subnets"), "192.0.2.1/24")
				if global {
					assert.Equal(t, "/infra/sites/site", d.Get("site_path"))
				}
			} else {
				assert.Equal(t, "/infra/dhcp-server-configs/dhcp", d.Get("dhcp_config_path"))
				services := d.Get("locale_service").(*schema.Set).List()
				require.Len(t, services, 1)
				assert.Equal(t, "locale", services[0].(map[string]interface{})["nsx_id"])
			}
		})
	}
}

func testGatewayBuilderChangedConfiguration(t *testing.T, tier0 bool) {
	oldVersion := util.NsxVersion
	t.Cleanup(func() { util.NsxVersion = oldVersion })
	util.NsxVersion = "9.2.0"

	name := "tier1"
	res := resourceNsxtPolicyTier1Gateway()
	if tier0 {
		name = "tier0"
		res = resourceNsxtPolicyTier0Gateway()
	}
	t.Run(name, func(t *testing.T) {
		raw := map[string]interface{}{"display_name": "updated", "dhcp_config_path": "/infra/dhcp-server-configs/dhcp", "edge_cluster_path": "/infra/sites/default/enforcement-points/default/edge-clusters/cluster"}
		if tier0 {
			raw["bgp_config"] = []interface{}{map[string]interface{}{"enabled": false, "local_as_num": "65000"}}
		} else {
			raw["tier0_path"] = "/infra/tier-0s/parent"
			raw["type"] = "ISOLATED"
		}
		old := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"display_name": "old"})
		old.SetId("gateway")
		diff, diffErr := res.Diff(context.Background(), old.State(), terraform.NewResourceConfigRaw(raw), nil)
		require.NoError(t, diffErr)
		d, dataErr := schema.InternalMap(res.Schema).Data(old.State(), diff)
		require.NoError(t, dataErr)
		m := newIsolatedMockProviderClient(t)
		m.Host = "https://nsx.invalid"
		m.PolicyHTTPClient = &http.Client{Transport: unitPolicyTransport(func(r *http.Request) (*http.Response, error) {
			body := `{"id":"default","resource_type":"LocaleServices","edge_cluster_path":"/infra/sites/default/enforcement-points/default/edge-clusters/cluster","revision":3}`
			if strings.HasSuffix(r.URL.Path, "/locale-services") {
				body = `{"results":[` + body + `],"result_count":1}`
			}
			return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
		})}
		ctx := utl.SessionContext{ClientType: utl.Local}
		connector := getPolicyConnector(m)
		var infra model.Infra
		var err error
		converter := bindings.NewTypeConverter()
		if tier0 {
			infra, err = policyTier0GatewayResourceToInfraStruct(ctx, d, connector, "gateway")
		} else {
			infra, err = policyTier1GatewayResourceToInfraStruct(ctx, d, m, connector, "gateway")
		}
		require.NoError(t, err)
		require.Len(t, infra.Children, 1)
		if tier0 {
			v, errs := converter.ConvertToGolang(infra.Children[0], model.ChildTier0BindingType())
			require.Empty(t, errs)
			gateway := v.(model.ChildTier0).Tier0
			assert.Equal(t, []string{"/infra/dhcp-server-configs/dhcp"}, gateway.DhcpConfigPaths)
			require.Len(t, gateway.Children, 1)
			v, errs = converter.ConvertToGolang(gateway.Children[0], model.ChildLocaleServicesBindingType())
			require.Empty(t, errs)
			require.Len(t, v.(model.ChildLocaleServices).LocaleServices.Children, 1)
		} else {
			v, errs := converter.ConvertToGolang(infra.Children[0], model.ChildTier1BindingType())
			require.Empty(t, errs)
			gateway := v.(model.ChildTier1).Tier1
			assert.Equal(t, "/infra/tier-0s/parent", *gateway.Tier0Path)
			assert.Equal(t, "ISOLATED", *gateway.Type_)
			assert.Equal(t, []string{"/infra/dhcp-server-configs/dhcp"}, gateway.DhcpConfigPaths)
			require.Len(t, gateway.Children, 1)
		}
	})
}

func testGlobalGatewayInterfaceCreation(t *testing.T, tier0 bool) {
	t.Helper()

	for _, scenario := range []string{"success", "site_missing", "site_not_matched", "list_error", "patch_error"} {
		t.Run(fmt.Sprintf("tier0_%v_%s", tier0, scenario), func(t *testing.T) {
			res := resourceNsxtPolicyTier1GatewayInterface()
			tier := "tier-1s"
			if tier0 {
				res = resourceNsxtPolicyTier0GatewayInterface()
				tier = "tier-0s"
			}
			raw := map[string]interface{}{"nsx_id": "interface", "gateway_path": "/infra/" + tier + "/gateway", "site_path": "/infra/sites/site", "segment_path": "/infra/segments/segment", "subnets": []interface{}{"192.0.2.1/24"}, "ipv6_ndra_profile_path": "/infra/ipv6-ndra-profiles/profile", "dhcp_relay_path": "/infra/dhcp-relay-configs/relay", "mtu": 1600}
			if scenario == "site_missing" {
				delete(raw, "site_path")
			}
			if scenario == "site_not_matched" {
				raw["site_path"] = "/infra/sites/other"
			}
			d := schema.TestResourceDataRaw(t, res.Schema, raw)
			m := newIsolatedMockProviderClient(t)
			m.Host = "https://nsx.invalid"
			m.PolicyGlobalManager = true
			m.CommonConfig.MaxRetries = 0
			patched := false
			m.PolicyHTTPClient = &http.Client{Transport: unitPolicyTransport(func(r *http.Request) (*http.Response, error) {
				if scenario == "site_missing" {
					t.Fatal("missing site must fail before API access")
				}
				service := `{"id":"locale","path":"/infra/tier-0s/gateway/locale-services/locale","edge_cluster_path":"/infra/sites/site/enforcement-points/default/edge-clusters/cluster"}`
				status := 200
				body := `{"id":"interface","display_name":"Interface","subnets":[{"ip_addresses":["192.0.2.1"],"prefix_len":24}],"mtu":1600}`
				if strings.HasSuffix(r.URL.Path, "/locale-services") {
					body = `{"results":[` + service + `],"result_count":1}`
					if scenario == "list_error" {
						status = 503
					}
				} else if strings.HasSuffix(r.URL.Path, "/locale-services/locale") {
					body = service
				} else if r.Method == http.MethodPatch {
					patched = true
					if scenario == "patch_error" {
						status = 503
					}
				} else if !patched {
					status = 404
				}
				if status != 200 {
					body = `{"error_code":` + fmt.Sprint(status) + `,"error_message":"unavailable","module_name":"test"}`
				}
				return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
			})}
			err := res.Create(d, m)
			if scenario != "success" {
				require.Error(t, err)
				assert.Empty(t, d.Id())
				return
			}
			require.NoError(t, err)
			assert.True(t, patched)
			assert.Equal(t, "interface", d.Id())
			assert.Equal(t, "locale", d.Get("locale_service_id"))
			assert.Equal(t, "/infra/sites/site", d.Get("site_path"))
		})
	}
}
