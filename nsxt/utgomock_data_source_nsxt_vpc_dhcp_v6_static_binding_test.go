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

func TestUnitNsxt_dataSourceNsxtVpcDhcpV6StaticBindingRead(t *testing.T) {
	parentPath := "/orgs/default/projects/proj-1/vpcs/vpc-1/subnets/subnet-1"
	rt := "DhcpV6StaticBindingConfig"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("binding-1"), DisplayName: str("binding-name"), Path: str(parentPath + "/dhcp-static-binding-configs/binding-1"), ResourceType: &rt,
	})

	t.Run("by id", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtVpcDhcpV6StaticBinding()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":          "binding-1",
			"parent_path": parentPath,
		})

		err := dataSourceNsxtVpcDhcpV6StaticBindingRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "binding-1", d.Id())
	})

	t.Run("by display name", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtVpcDhcpV6StaticBinding()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "binding-name",
			"parent_path":  parentPath,
		})

		err := dataSourceNsxtVpcDhcpV6StaticBindingRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "binding-1", d.Id())
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtVpcDhcpV6StaticBinding()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":          "binding-1",
			"parent_path": parentPath,
		})

		err := dataSourceNsxtVpcDhcpV6StaticBindingRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
