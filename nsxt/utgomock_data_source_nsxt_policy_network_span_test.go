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

func TestUnitNsxt_dataSourceNsxtPolicyNetworkSpanRead(t *testing.T) {
	rt := "NetworkSpan"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("nspan-1"), DisplayName: str("nspan-name"), Path: str("/infra/network-spans/nspan-1"), ResourceType: &rt,
	})

	t.Run("by id", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyNetworkSpan()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "nspan-1"})

		err := dataSourceNsxtPolicyNetworkSpanRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "nspan-1", d.Id())
	})

	t.Run("is_default custom field", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyNetworkSpan()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"is_default": true})

		err := dataSourceNsxtPolicyNetworkSpanRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "nspan-1", d.Id())
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyNetworkSpan()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": "nspan-1"})

		err := dataSourceNsxtPolicyNetworkSpanRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
