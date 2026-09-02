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

func TestUnitNsxt_dataSourceNsxtPolicyRouteControllerInterfaceRead(t *testing.T) {
	rt := "RouteControllerInterface"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("rc-if-1"), DisplayName: str("rc-if-name"), Path: str("/infra/route-controllers/rc-1/interfaces/rc-if-1"), ResourceType: &rt,
	})

	t.Run("by id, no parent_path", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "rc-if-1"})

		err := dataSourceNsxtPolicyRouteControllerInterfaceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "rc-if-1", d.Id())
	})

	t.Run("by id with valid parent_path", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":          "rc-if-1",
			"parent_path": "/infra/route-controllers/rc-1",
		})

		err := dataSourceNsxtPolicyRouteControllerInterfaceRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "rc-if-1", d.Id())
	})

	t.Run("invalid parent_path", func(t *testing.T) {
		ds := dataSourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":          "rc-if-1",
			"parent_path": "/x",
		})

		err := dataSourceNsxtPolicyRouteControllerInterfaceRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid parent_path")
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyRouteControllerInterface()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "rc-if-1"})

		err := dataSourceNsxtPolicyRouteControllerInterfaceRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
