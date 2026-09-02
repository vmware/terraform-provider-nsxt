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

func vpcServiceProfileDsContext() map[string]interface{} {
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

func TestUnitNsxt_dataSourceNsxtVpcServiceProfileRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	rt := "VpcServiceProfile"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("vsp-1"), DisplayName: str("vsp-name"),
		Path: str("/orgs/default/projects/proj-1/vpc-service-profiles/vsp-1"), ResourceType: &rt,
	})

	t.Run("by id", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		raw := vpcServiceProfileDsContext()
		raw["id"] = "vsp-1"

		ds := dataSourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtVpcServiceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "vsp-1", d.Id())
	})

	t.Run("is_default custom field", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		raw := vpcServiceProfileDsContext()
		raw["is_default"] = true

		ds := dataSourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtVpcServiceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "vsp-1", d.Id())
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		raw := vpcServiceProfileDsContext()
		raw["id"] = "vsp-1"

		ds := dataSourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtVpcServiceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
