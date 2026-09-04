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

func vpcEndpointDsContext() map[string]interface{} {
	return map[string]interface{}{
		"context": []interface{}{
			map[string]interface{}{
				"project_id":  "proj-1",
				"vpc_id":      "vpc-1",
				"from_global": false,
			},
		},
	}
}

func TestUnitNsxt_dataSourceNsxtVpcEndpointRead(t *testing.T) {
	util.NsxVersion = "9.2.0"
	defer func() { util.NsxVersion = "" }()

	rt := "VpcEndpoint"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("vep-1"), DisplayName: str("vep-name"),
		Path: str("/orgs/default/projects/proj-1/vpcs/vpc-1/vpc-endpoints/vep-1"), ResourceType: &rt,
	})

	t.Run("by id", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		raw := vpcEndpointDsContext()
		raw["id"] = "vep-1"

		ds := dataSourceNsxtVpcEndpoint()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtVpcEndpointRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "vep-1", d.Id())
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		raw := vpcEndpointDsContext()
		raw["id"] = "vep-1"

		ds := dataSourceNsxtVpcEndpoint()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtVpcEndpointRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("version gate", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "9.2.0" }()

		raw := vpcEndpointDsContext()
		raw["id"] = "vep-1"

		ds := dataSourceNsxtVpcEndpoint()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtVpcEndpointRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "requires NSX version")
	})
}
