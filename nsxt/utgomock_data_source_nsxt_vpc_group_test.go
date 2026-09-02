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

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func TestUnitNsxt_DataSourceNsxtVpcGroupRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	rt := "Group"

	t.Run("by display_name success", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str("grp-id-1"), DisplayName: str("grp1"), Path: str("/orgs/default/projects/p1/vpcs/vpc1/groups/grp-id-1"),
			ResourceType: &rt,
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtVpcGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "grp1",
			"context": []interface{}{
				map[string]interface{}{"project_id": "p1", "vpc_id": "vpc1"},
			},
		})

		err := dataSourceNsxtVpcGroupRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "grp-id-1", d.Id())
	})

	t.Run("search error is wrapped", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search-fail")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtVpcGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "grp1",
			"context": []interface{}{
				map[string]interface{}{"project_id": "p1", "vpc_id": "vpc1"},
			},
		})

		err := dataSourceNsxtVpcGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "search-fail")
	})

	t.Run("not found", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{}, ResultCount: i64(0),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtVpcGroup()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
			"context": []interface{}{
				map[string]interface{}{"project_id": "p1", "vpc_id": "vpc1"},
			},
		})

		err := dataSourceNsxtVpcGroupRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})
}
