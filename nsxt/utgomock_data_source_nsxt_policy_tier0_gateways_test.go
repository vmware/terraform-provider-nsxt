//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// dataSourceNsxtPolicyTier0GatewaysRead delegates to policyDataSourceCreateMap, which is
// backed by the search (cliQueryClient) machinery already covered by
// policy_search_unit_test.go. This file reuses seqQueryListClient / setupCliQueryClientStub /
// policyResourceToStructValue from that file instead of introducing new mocks.

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

func TestMockDataSourceNsxtPolicyTier0GatewaysRead(t *testing.T) {
	t.Run("success without display_name filter", func(t *testing.T) {
		sv1 := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("t0-1"), DisplayName: str("gw-one"), Path: str("/infra/tier-0s/t0-1"), ResourceType: str("Tier0"),
		})
		sv2 := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("t0-2"), DisplayName: str("gw-two"), Path: str("/infra/tier-0s/t0-2"), ResourceType: str("Tier0"),
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv1, sv2}, ResultCount: i64(2),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyTier0Gateways()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyTier0GatewaysRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		assert.Equal(t, "gw-one", items["t0-1"])
		assert.Equal(t, "gw-two", items["t0-2"])
	})

	t.Run("success with display_name regex filter", func(t *testing.T) {
		sv1 := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("t0-1"), DisplayName: str("gw-one"), Path: str("/infra/tier-0s/t0-1"), ResourceType: str("Tier0"),
		})
		sv2 := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("t0-2"), DisplayName: str("other"), Path: str("/infra/tier-0s/t0-2"), ResourceType: str("Tier0"),
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv1, sv2}, ResultCount: i64(2),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyTier0Gateways()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "^gw-",
		})

		err := dataSourceNsxtPolicyTier0GatewaysRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		items := d.Get("items").(map[string]interface{})
		assert.Equal(t, "gw-one", items["t0-1"])
		_, hasOther := items["t0-2"]
		assert.False(t, hasOther)
	})

	t.Run("invalid regex", func(t *testing.T) {
		sv1 := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("t0-1"), DisplayName: str("gw-one"), Path: str("/infra/tier-0s/t0-1"), ResourceType: str("Tier0"),
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv1}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyTier0Gateways()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "[",
		})

		err := dataSourceNsxtPolicyTier0GatewaysRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("list error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("list boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyTier0Gateways()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyTier0GatewaysRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error in listing the Tier0 gateways items")
	})
}
