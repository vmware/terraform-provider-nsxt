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

func cnaDsContext() map[string]interface{} {
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

func TestUnitNsxt_dataSourceNsxtPolicyProjectCentralizedNetworkAttachmentRead(t *testing.T) {
	rt := "CentralizedNetworkAttachment"
	sv := policyResourceToStructValue(t, gmModel.PolicyResource{
		Id: str("cna-1"), DisplayName: str("cna-name"), Path: str("/orgs/default/projects/proj-1/centralized-network-attachments/cna-1"), ResourceType: &rt,
	})

	t.Run("by id", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		raw := cnaDsContext()
		raw["id"] = "cna-1"

		ds := dataSourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtPolicyProjectCentralizedNetworkAttachmentRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "cna-1", d.Id())
	})

	t.Run("by display name", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		raw := cnaDsContext()
		raw["display_name"] = "cna-name"

		ds := dataSourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtPolicyProjectCentralizedNetworkAttachmentRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "cna-1", d.Id())
	})

	t.Run("no id or display_name", func(t *testing.T) {
		raw := cnaDsContext()

		ds := dataSourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtPolicyProjectCentralizedNetworkAttachmentRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("search error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("boom")}}
		defer setupCliQueryClientStub(t, stub)()

		raw := cnaDsContext()
		raw["id"] = "cna-1"

		ds := dataSourceNsxtPolicyProjectCentralizedNetworkAttachment()
		d := schema.TestResourceDataRaw(t, ds.Schema, raw)

		err := dataSourceNsxtPolicyProjectCentralizedNetworkAttachmentRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
