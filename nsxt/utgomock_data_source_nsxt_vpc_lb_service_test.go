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

// vpcLbService* fixtures are defined in utgomock_resource_nsxt_vpc_lb_service_test.go

func vpcLbDsServiceItem(t *testing.T, id, name string) *data.StructValue {
	return policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str(id), DisplayName: str(name), ResourceType: str("LBService"),
		Path: str("/orgs/default/projects/" + vpcLbPrfProject + "/vpcs/" + vpcLbPrfVpc + "/vpc-lbs/" + id),
	})
}

func vpcLbDsServiceRead(t *testing.T, mc nsxtClients, raw map[string]interface{}) (*schema.ResourceData, error) {
	ds := dataSourceNsxtVpcLBService()
	in := map[string]interface{}{"context": vpcLbPrfContext()}
	for k, v := range raw {
		in[k] = v
	}
	d := schema.TestResourceDataRaw(t, ds.Schema, in)
	return d, ds.Read(d, mc)
}

func TestMockDataSourceNsxtVpcLBService(t *testing.T) {
	searchOne := func(t *testing.T, items ...*data.StructValue) *seqQueryListClient {
		return &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: items, ResultCount: i64(int64(len(items))),
		}}}
	}

	t.Run("by ID", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.0")
		stub := searchOne(t, vpcLbDsServiceItem(t, vpcLbServiceID, vpcLbServiceDisplayName))
		defer setupCliQueryClientStub(t, stub)()

		d, err := vpcLbDsServiceRead(t, newGoMockProviderClient(), map[string]interface{}{"id": vpcLbServiceID})
		require.NoError(t, err)
		assert.Equal(t, vpcLbServiceID, d.Id())
		assert.Equal(t, vpcLbServiceDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcLbServicePath, d.Get("path"))
	})

	t.Run("by name perfect match", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.0")
		stub := searchOne(t,
			vpcLbDsServiceItem(t, "other", vpcLbServiceDisplayName+" extra"),
			vpcLbDsServiceItem(t, vpcLbServiceID, vpcLbServiceDisplayName))
		defer setupCliQueryClientStub(t, stub)()

		d, err := vpcLbDsServiceRead(t, newGoMockProviderClient(), map[string]interface{}{"display_name": vpcLbServiceDisplayName})
		require.NoError(t, err)
		assert.Equal(t, vpcLbServiceID, d.Id())
	})

	t.Run("by name prefix match", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.0")
		stub := searchOne(t, vpcLbDsServiceItem(t, vpcLbServiceID, vpcLbServiceDisplayName))
		defer setupCliQueryClientStub(t, stub)()

		d, err := vpcLbDsServiceRead(t, newGoMockProviderClient(), map[string]interface{}{"display_name": "Test VPC"})
		require.NoError(t, err)
		assert.Equal(t, vpcLbServiceID, d.Id())
	})

	t.Run("multiple perfect matches", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.0")
		stub := searchOne(t, vpcLbDsServiceItem(t, "d1", "dup"), vpcLbDsServiceItem(t, "d2", "dup"))
		defer setupCliQueryClientStub(t, stub)()

		_, err := vpcLbDsServiceRead(t, newGoMockProviderClient(), map[string]interface{}{"display_name": "dup"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("not found", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.0")
		stub := searchOne(t)
		defer setupCliQueryClientStub(t, stub)()

		_, err := vpcLbDsServiceRead(t, newGoMockProviderClient(), map[string]interface{}{"display_name": "nonexistent"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("search error", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.0")
		stub := &seqQueryListClient{errs: []error{errors.New("search failed")}}
		defer setupCliQueryClientStub(t, stub)()

		_, err := vpcLbDsServiceRead(t, newGoMockProviderClient(), map[string]interface{}{"id": vpcLbServiceID})
		require.Error(t, err)
	})

	t.Run("NSX version lower than 9.0.0", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "4.2.0")
		stub := searchOne(t)
		defer setupCliQueryClientStub(t, stub)()

		_, err := vpcLbDsServiceRead(t, newGoMockProviderClient(), map[string]interface{}{"id": vpcLbServiceID})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.0.0")
		assert.Equal(t, 0, stub.call)
	})

}
