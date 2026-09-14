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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	gmModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func dnsServiceDsContext() map[string]interface{} {
	return map[string]interface{}{
		"context": []interface{}{
			map[string]interface{}{
				"project_id":  "proj-1",
				"vpc_id":      "",
				"from_global": false,
			},
		},
	}
}

func TestUnitNsxt_dataSourceNsxtPolicyDnsServiceRead(t *testing.T) {
	rt := "DnsService"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("dns-service-1"), DisplayName: str("service-name"), Path: str("/orgs/default/projects/proj-1/dns-services/dns-service-1"), ResourceType: &rt,
	})

	t.Run("by id", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		raw := dnsServiceDsContext()
		raw["id"] = "dns-service-1"

		ds := dataSourceNsxtPolicyDnsService()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtPolicyDnsServiceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "dns-service-1", d.Id())
	})

	t.Run("by display name", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		raw := dnsServiceDsContext()
		raw["display_name"] = "service-name"

		ds := dataSourceNsxtPolicyDnsService()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtPolicyDnsServiceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "dns-service-1", d.Id())
	})

	t.Run("no id or display_name", func(t *testing.T) {
		raw := dnsServiceDsContext()

		ds := dataSourceNsxtPolicyDnsService()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtPolicyDnsServiceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		raw := dnsServiceDsContext()
		raw["id"] = "dns-service-1"

		ds := dataSourceNsxtPolicyDnsService()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtPolicyDnsServiceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
