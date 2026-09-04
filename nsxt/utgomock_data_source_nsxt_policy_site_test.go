//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// dataSourceNsxtPolicySiteRead delegates to policyDataSourceResourceRead, which is backed by
// the search (cliQueryClient) machinery already covered by policy_search_unit_test.go. This
// file reuses seqQueryListClient / setupCliQueryClientStub / policyResourceToStructValue from
// that file instead of introducing new mocks.

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

func TestMockDataSourceNsxtPolicySiteReadGuard(t *testing.T) {
	t.Run("Read fails when not global manager", func(t *testing.T) {
		ds := dataSourceNsxtPolicySite()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": "site-1",
		})

		err := dataSourceNsxtPolicySiteRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Global Manager")
	})
}

func TestMockDataSourceNsxtPolicySiteRead(t *testing.T) {
	t.Run("success by id", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("site-1"), DisplayName: str("Default"), Path: str("/infra/sites/site-1"),
			ResourceType: str("Site"),
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicySite()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": "site-1",
		})

		err := dataSourceNsxtPolicySiteRead(d, newGoMockGlobalProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "site-1", d.Id())
		assert.Equal(t, "Default", d.Get("display_name"))
	})

	t.Run("success by display_name", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("site-2"), DisplayName: str("Second Site"), Path: str("/infra/sites/site-2"),
			ResourceType: str("Site"),
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicySite()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "Second Site",
		})

		err := dataSourceNsxtPolicySiteRead(d, newGoMockGlobalProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "site-2", d.Id())
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicySite()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": "site-3",
		})

		err := dataSourceNsxtPolicySiteRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "search boom")
	})
}
