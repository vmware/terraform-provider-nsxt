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

func dnsRecordDsContext() map[string]interface{} {
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

func TestUnitNsxt_dataSourceNsxtPolicyDnsRecordRead(t *testing.T) {
	rt := "DnsRecord"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("dns-record-1"), DisplayName: str("record-name"), Path: str("/orgs/default/projects/proj-1/dns-records/dns-record-1"), ResourceType: &rt,
	})

	t.Run("by id", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		raw := dnsRecordDsContext()
		raw["id"] = "dns-record-1"

		ds := dataSourceNsxtPolicyDnsRecord()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtPolicyDnsRecordRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "dns-record-1", d.Id())
	})

	t.Run("by display name with zone filter", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		raw := dnsRecordDsContext()
		raw["display_name"] = "record-name"
		raw["zone_path"] = "/orgs/default/projects/proj-1/dns-services/svc-1"

		ds := dataSourceNsxtPolicyDnsRecord()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtPolicyDnsRecordRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "dns-record-1", d.Id())
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		raw := dnsRecordDsContext()
		raw["id"] = "dns-record-1"

		ds := dataSourceNsxtPolicyDnsRecord()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtPolicyDnsRecordRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
