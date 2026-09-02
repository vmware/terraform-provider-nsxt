//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// Reuses the cliQueryClient seam and seqQueryListClient stub already defined
// for the search-backed data sources in policy_search_unit_test.go.

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

// bareMetalServerInterfaceToStructValue is reused from
// utgomock_resource_nsxt_policy_baremetal_server_interface_tags_test.go.

func TestMockDataSourceNsxtPolicyBareMetalServerInterfacesRead(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	t.Run("success with filters", func(t *testing.T) {
		extID := "iface-1"
		name := "mgmt-nic"
		bmsID := "bms-1"
		sourceID := "src-1"
		state := "UP"
		resourceType := "BareMetalServerInterface"
		isMgmt := true
		lastSync := int64(100)
		sv := bareMetalServerInterfaceToStructValue(t, nsxModel.BareMetalServerInterface{
			ExternalId:      &extID,
			DisplayName:     &name,
			BmsExternalId:   &bmsID,
			SourceId:        &sourceID,
			State:           &state,
			ResourceType:    &resourceType,
			IsMgmtInterface: &isMgmt,
			LastSyncTime:    &lastSync,
			IpAddresses:     []string{"10.0.0.5"},
			Tags:            []nsxModel.Tag{{Scope: strPtr("env"), Tag: strPtr("prod")}},
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyBareMetalServerInterfaces()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name":    "mgmt",
			"bms_external_id": bmsID,
			"source_id":       sourceID,
		})

		err := dataSourceNsxtPolicyBareMetalServerInterfacesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		results := d.Get("results").([]interface{})
		require.Len(t, results, 1)
		item := results[0].(map[string]interface{})
		assert.Equal(t, extID, item["external_id"])
		assert.Equal(t, name, item["display_name"])
		assert.Equal(t, true, item["is_mgmt_interface"])
	})

	t.Run("display_name regex filters out non-matching", func(t *testing.T) {
		extID := "iface-2"
		name := "data-nic"
		sv := bareMetalServerInterfaceToStructValue(t, nsxModel.BareMetalServerInterface{
			ExternalId:  &extID,
			DisplayName: &name,
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyBareMetalServerInterfaces()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "^mgmt",
		})

		err := dataSourceNsxtPolicyBareMetalServerInterfacesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		results := d.Get("results").([]interface{})
		assert.Len(t, results, 0)
	})

	t.Run("invalid display_name regex errors", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{Results: []*data.StructValue{}, ResultCount: i64(0)}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyBareMetalServerInterfaces()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "(unterminated",
		})

		err := dataSourceNsxtPolicyBareMetalServerInterfacesRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid regex")
	})

	t.Run("search error is wrapped", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyBareMetalServerInterfaces()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyBareMetalServerInterfacesRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error listing bare metal server interfaces")
	})
}
