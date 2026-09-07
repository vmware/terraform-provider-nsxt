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

func TestUnitNsxt_dataSourceNsxtPolicyEdgeClusterRead(t *testing.T) {
	edgeClusterDsID := "edge-cluster-1"

	t.Run("local manager success", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str(edgeClusterDsID), DisplayName: str("ec-name"), Path: str("/infra/ec"), ResourceType: str("PolicyEdgeCluster"),
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyEdgeCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": edgeClusterDsID})

		err := dataSourceNsxtPolicyEdgeClusterRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, edgeClusterDsID, d.Id())
	})

	t.Run("local manager error", func(t *testing.T) {
		stub := &seqQueryListClient{errs: []error{errors.New("search failed")}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyEdgeCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": edgeClusterDsID})

		err := dataSourceNsxtPolicyEdgeClusterRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("site_path on local manager is rejected", func(t *testing.T) {
		ds := dataSourceNsxtPolicyEdgeCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":        edgeClusterDsID,
			"site_path": "/infra/sites/default",
		})

		err := dataSourceNsxtPolicyEdgeClusterRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("global manager without site_path is rejected", func(t *testing.T) {
		ds := dataSourceNsxtPolicyEdgeCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{"id": edgeClusterDsID})

		m := newGoMockProviderClient()
		m.PolicyGlobalManager = true
		err := dataSourceNsxtPolicyEdgeClusterRead(d, m)
		require.Error(t, err)
	})

	t.Run("global manager with site_path success", func(t *testing.T) {
		sv := policyResourceToStructValue(t, gmModel.PolicyResource{
			Id: str(edgeClusterDsID), DisplayName: str("ec-name"), Path: str("/infra/ec"), ResourceType: str("PolicyEdgeCluster"),
		})
		stub := &seqQueryListClient{responses: []nsxModel.SearchResponse{{
			Results: []*data.StructValue{sv}, ResultCount: i64(1),
		}}}
		defer setupCliQueryClientStub(t, stub)()

		ds := dataSourceNsxtPolicyEdgeCluster()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":        edgeClusterDsID,
			"site_path": "/infra/sites/default",
		})

		m := newGoMockProviderClient()
		m.PolicyGlobalManager = true
		err := dataSourceNsxtPolicyEdgeClusterRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, edgeClusterDsID, d.Id())
	})
}
