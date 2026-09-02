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

func projectIPAllocDsContext() map[string]interface{} {
	return map[string]interface{}{
		"allocation_ips": "10.0.0.1",
		"context": []interface{}{
			map[string]interface{}{
				"project_id":  "proj-1",
				"vpc_id":      "",
				"from_global": false,
			},
		},
	}
}

func TestUnitNsxt_dataSourceNsxtProjectIpAddressAllocationRead(t *testing.T) {
	rt := "ProjectIpAddressAllocation"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("palloc-1"), DisplayName: str("palloc-name"), Path: str("/orgs/default/projects/proj-1/ip-address-allocations/palloc-1"), ResourceType: &rt,
	})

	t.Run("by id", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		raw := projectIPAllocDsContext()
		raw["id"] = "palloc-1"

		ds := dataSourceNsxtProjectIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtProjectIpAddressAllocationRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "palloc-1", d.Id())
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		raw := projectIPAllocDsContext()
		raw["id"] = "palloc-1"

		ds := dataSourceNsxtProjectIpAddressAllocation()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtProjectIpAddressAllocationRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
