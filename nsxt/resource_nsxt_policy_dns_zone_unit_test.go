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

func TestUnitNsxt_policyDnsZoneFromSchema(t *testing.T) {
	res := resourceNsxtPolicyDnsZone()

	t.Run("basic fields and empty soa are converted", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":     "zone1",
			"description":      "a zone",
			"dns_domain_name":  "example.com",
			"ttl":              600,
			"resolution_scope": []interface{}{"/orgs/default/projects/proj1/vpcs/vpc1"},
		})
		obj := policyDnsZoneFromSchema(d)
		assert.Equal(t, "zone1", *obj.DisplayName)
		assert.Equal(t, "a zone", *obj.Description)
		assert.Equal(t, "example.com", *obj.DnsDomainName)
		assert.EqualValues(t, 600, *obj.Ttl)
		assert.ElementsMatch(t, []string{"/orgs/default/projects/proj1/vpcs/vpc1"}, obj.ResolutionScope)
		assert.Nil(t, obj.Soa, "all-zero-value soa block should not be sent")
	})

	t.Run("soa fields are converted when set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":    "zone1",
			"dns_domain_name": "example.com",
			"soa": []interface{}{
				map[string]interface{}{
					"primary_nameserver": "ns1.example.com.",
					"responsible_party":  "admin.example.com.",
					"serial_number":      5,
					"refresh_interval":   3600,
					"retry_interval":     600,
					"expire_time":        86400,
					"negative_cache_ttl": 60,
				},
			},
		})
		obj := policyDnsZoneFromSchema(d)
		require.NotNil(t, obj.Soa)
		assert.Equal(t, "ns1.example.com.", *obj.Soa.PrimaryNameserver)
		assert.Equal(t, "admin.example.com.", *obj.Soa.ResponsibleParty)
		assert.EqualValues(t, 5, *obj.Soa.SerialNumber)
		assert.EqualValues(t, 3600, *obj.Soa.RefreshInterval)
		assert.EqualValues(t, 600, *obj.Soa.RetryInterval)
		assert.EqualValues(t, 86400, *obj.Soa.ExpireTime)
		assert.EqualValues(t, 60, *obj.Soa.NegativeCacheTtl)
	})

	t.Run("soa with only negative_cache_ttl zero is still sent", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":    "zone1",
			"dns_domain_name": "example.com",
			"soa": []interface{}{
				map[string]interface{}{
					"primary_nameserver": "ns1.example.com.",
				},
			},
		})
		obj := policyDnsZoneFromSchema(d)
		require.NotNil(t, obj.Soa)
		assert.Equal(t, "ns1.example.com.", *obj.Soa.PrimaryNameserver)
	})
}

func TestUnitNsxt_policyDnsZoneToSchema(t *testing.T) {
	res := resourceNsxtPolicyDnsZone()

	t.Run("populates schema fields including soa", func(t *testing.T) {
		displayName := "zone1"
		description := "a zone"
		dnsDomain := "example.com"
		ttl := int64(600)
		primaryNS := "ns1.example.com."
		serial := int64(5)

		obj := model.DnsZone{
			DisplayName:     &displayName,
			Description:     &description,
			DnsDomainName:   &dnsDomain,
			Ttl:             &ttl,
			ResolutionScope: []string{"/orgs/default/projects/proj1/vpcs/vpc1"},
			Soa: &model.DnsZoneSoa{
				PrimaryNameserver: &primaryNS,
				SerialNumber:      &serial,
			},
		}

		d := res.TestResourceData()
		policyDnsZoneToSchema(d, obj)

		assert.Equal(t, displayName, d.Get("display_name"))
		assert.Equal(t, description, d.Get("description"))
		assert.Equal(t, dnsDomain, d.Get("dns_domain_name"))
		assert.Equal(t, 600, d.Get("ttl"))
		soa := d.Get("soa").([]interface{})
		require.Len(t, soa, 1)
		soaMap := soa[0].(map[string]interface{})
		assert.Equal(t, primaryNS, soaMap["primary_nameserver"])
		assert.Equal(t, 5, soaMap["serial_number"])
	})

	t.Run("nil soa leaves schema soa unset", func(t *testing.T) {
		displayName := "zone1"
		obj := model.DnsZone{DisplayName: &displayName}

		d := res.TestResourceData()
		policyDnsZoneToSchema(d, obj)

		assert.Empty(t, d.Get("soa").([]interface{}))
	})
}
