//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// Reuses the cliQueryClient seam and seqQueryListClient stub already defined
// for the search-backed data sources in policy_search_unit_test.go, and the
// bareMetalServerToStructValue helper already defined in
// utgomock_resource_nsxt_policy_baremetal_server_tags_test.go.

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

func TestMockDataSourceNsxtPolicyBareMetalServersRead(t *testing.T) {
	util.NsxVersion = "9.0.0"
	defer func() { util.NsxVersion = "" }()

	t.Run("success with filters", func(t *testing.T) {
		extID := "server-1"
		name := "esx-01"
		sourceID := "src-1"
		resourceType := "BareMetalServer"
		cpuCores := int64(32)
		osName := "Ubuntu"
		osVersion := "22.04"
		lastSync := int64(200)
		sv := bareMetalServerToStructValue(t, nsxModel.BareMetalServer{
			ExternalId:   &extID,
			DisplayName:  &name,
			SourceId:     &sourceID,
			ResourceType: &resourceType,
			CpuCores:     &cpuCores,
			LastSyncTime: &lastSync,
			OsInfo: &nsxModel.OsInfo{
				OsName:    &osName,
				OsVersion: &osVersion,
			},
			Tags: []nsxModel.Tag{{Scope: strPtr("env"), Tag: strPtr("prod")}},
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyBareMetalServers()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "esx",
			"source_id":    sourceID,
			"os_name":      "ubuntu",
			"os_version":   "22.04",
		})

		err := dataSourceNsxtPolicyBareMetalServersRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		results := d.Get("results").([]interface{})
		require.Len(t, results, 1)
		item := results[0].(map[string]interface{})
		assert.Equal(t, extID, item["external_id"])
		assert.Equal(t, name, item["display_name"])
		assert.Equal(t, 32, item["cpu_cores"])
		assert.Equal(t, osName, item["os_name"])
	})

	t.Run("os_name filter excludes non-matching", func(t *testing.T) {
		extID := "server-2"
		name := "esx-02"
		osName := "Windows"
		sv := bareMetalServerToStructValue(t, nsxModel.BareMetalServer{
			ExternalId:  &extID,
			DisplayName: &name,
			OsInfo:      &nsxModel.OsInfo{OsName: &osName},
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyBareMetalServers()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"os_name": "ubuntu",
		})

		err := dataSourceNsxtPolicyBareMetalServersRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		results := d.Get("results").([]interface{})
		assert.Len(t, results, 0)
	})

	t.Run("invalid display_name regex errors", func(t *testing.T) {
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{Results: []*data.StructValue{}, ResultCount: i64(0)}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyBareMetalServers()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "(unterminated",
		})

		err := dataSourceNsxtPolicyBareMetalServersRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid regex")
	})

	t.Run("search error is wrapped", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search boom")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyBareMetalServers()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyBareMetalServersRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error listing bare metal servers")
	})
}
