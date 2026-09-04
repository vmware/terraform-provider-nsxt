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

func TestUnitNsxt_dataSourceNsxtPolicyRouteControllerBgpNeighborRead(t *testing.T) {
	rt := "RouteControllerBgpNeighborConfig"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("rc-bgp-1"), DisplayName: str("rc-bgp-name"), Path: str("/infra/route-controllers/rc-1/bgp/neighbors/rc-bgp-1"), ResourceType: &rt,
	})

	t.Run("by id, no parent_path", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "rc-bgp-1"})

		err := dataSourceNsxtPolicyRouteControllerBgpNeighborRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "rc-bgp-1", d.Id())
	})

	t.Run("by id with valid parent_path", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":          "rc-bgp-1",
			"parent_path": "/infra/route-controllers/rc-1/bgp",
		})

		err := dataSourceNsxtPolicyRouteControllerBgpNeighborRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "rc-bgp-1", d.Id())
	})

	t.Run("invalid parent_path", func(t *testing.T) {
		ds := dataSourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":          "rc-bgp-1",
			"parent_path": "/x",
		})

		err := dataSourceNsxtPolicyRouteControllerBgpNeighborRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid parent_path")
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyRouteControllerBgpNeighbor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "rc-bgp-1"})

		err := dataSourceNsxtPolicyRouteControllerBgpNeighborRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
