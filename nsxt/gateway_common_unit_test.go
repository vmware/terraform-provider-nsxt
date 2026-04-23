//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	defer func() { util.NsxVersion = "" }()
	util.NsxVersion = "3.2.0"
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
