//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func upgradePostcheckData(overrides map[string]interface{}) map[string]interface{} {
	data := map[string]interface{}{
		"upgrade_run_id": "run-1",
		"type":           edgeUpgradeGroup,
		"timeout":        5,
		"interval":       1,
		"delay":          0,
	}
	for k, v := range overrides {
		data[k] = v
	}
	return data
}

// NOTE: dataSourceNsxtUpgradePostCheckRead's client accessor,
// cliUpgradeUnitGroupsAggregateInfoClient, is declared as
//
//	var cliUpgradeUnitGroupsAggregateInfoClient = upgrade_unit_groups.NewAggregateInfoClient
//
// Unlike every sibling client var in this file/package (e.g.
// cliUpgradeChecksInfoClient, cliPreUpgradeChecksFailuresClient,
// cliUpgradeBundlesClient, ...), which are all wrapped in a func literal typed
// to return the *exported* SDK interface (upgrade.SummaryClient,
// pre_upgrade_checks.FailuresClient, etc.), this one is a bare assignment of
// the raw constructor. Its inferred static type is
// func(client.Connector) *upgrade_unit_groups.aggregateInfoClient, where
// aggregateInfoClient is an *unexported* concrete struct in the SDK's
// upgrade_unit_groups package. Go's type identity rules make it impossible to
// name or construct a value of that type from this (nsxt) package, so the
// var cannot be swapped for a gomock double without editing the non-test
// source file to wrap it in an interface-typed literal like its siblings
// (which is out of scope here). As a result only the parameter-validation
// branch, which returns before touching the client, is unit-testable today.
func TestUnitNsxt_DataSourceNsxtUpgradePostCheckRead(t *testing.T) {
	t.Run("invalid upgrade_run_id is rejected", func(t *testing.T) {
		ds := dataSourceNsxtUpgradePostCheck()
		d := schema.TestResourceDataRaw(t, ds.Schema, upgradePostcheckData(map[string]interface{}{
			"upgrade_run_id": "not-a-valid-id",
		}))

		err := dataSourceNsxtUpgradePostCheckRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "value for upgrade_run_id is invalid")
	})
}
