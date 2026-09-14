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

func TestUnitNsxt_dataSourceNsxtPolicyTransitGatewayCommunityListRead(t *testing.T) {
	parentPath := "/orgs/default/projects/proj-1/transit-gateways/tgw-1"
	rt := "CommunityList"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("community-list-1"), DisplayName: str("community-list-name"), Path: str(parentPath + "/community-lists/community-list-1"), ResourceType: &rt,
	})

	t.Run("by id", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyTransitGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":          "community-list-1",
			"parent_path": parentPath,
		})

		err := dataSourceNsxtPolicyTransitGatewayCommunityListRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "community-list-1", d.Id())
	})

	t.Run("by display name", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyTransitGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "community-list-name",
			"parent_path":  parentPath,
		})

		err := dataSourceNsxtPolicyTransitGatewayCommunityListRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "community-list-1", d.Id())
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyTransitGatewayCommunityList()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":          "community-list-1",
			"parent_path": parentPath,
		})

		err := dataSourceNsxtPolicyTransitGatewayCommunityListRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
