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

func vpcDsContext() map[string]interface{} {
	return map[string]interface{}{
		"context": []interface{}{
			map[string]interface{}{
				"project_id":  "proj-1",
				"vpc_id":      "",
				"from_global": false,
			},
		},
	}
}

func TestUnitNsxt_dataSourceNsxtVPCRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	rt := "Vpc"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("vpc-ds-1"), DisplayName: str("vpc-ds-name"), Path: str("/orgs/default/projects/proj-1/vpcs/vpc-ds-1"), ResourceType: &rt,
	})

	t.Run("by id", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		raw := vpcDsContext()
		raw["id"] = "vpc-ds-1"

		ds := dataSourceNsxtVPC()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtVPCRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "vpc-ds-1", d.Id())
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		raw := vpcDsContext()
		raw["id"] = "vpc-ds-1"

		ds := dataSourceNsxtVPC()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtVPCRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
