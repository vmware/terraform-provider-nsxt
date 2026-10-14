//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run:
// mockgen -destination=mocks/nsx/upgrade/upgrade_unit_groups/AggregateInfoClient.go -package=mocks github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/upgrade/upgrade_unit_groups AggregateInfoClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/upgrade/upgrade_unit_groups"
	"go.uber.org/mock/gomock"

	upgradeunitgroupsmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/upgrade/upgrade_unit_groups"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func upgradePostcheckData(overrides map[string]interface{}) map[string]interface{} {
	data := map[string]interface{}{
		"upgrade_run_id": util.GetVerifiableID("run-1", "nsxt_upgrade_run"),
		"type":           edgeUpgradeGroup,
		"timeout":        5,
		// Zero (below the schema minimum, which TestResourceDataRaw doesn't
		// enforce) makes WaitForState use its short default backoff.
		"interval": 0,
		"delay":    0,
	}
	for k, v := range overrides {
		data[k] = v
	}
	return data
}

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

func setupUpgradePostcheckMock(t *testing.T) *upgradeunitgroupsmocks.MockAggregateInfoClient {
	ctrl := gomock.NewController(t)
	mockClient := upgradeunitgroupsmocks.NewMockAggregateInfoClient(ctrl)
	original := cliUpgradeUnitGroupsAggregateInfoClient
	t.Cleanup(func() { cliUpgradeUnitGroupsAggregateInfoClient = original })
	cliUpgradeUnitGroupsAggregateInfoClient = func(_ vapiProtocolClient.Connector) upgrade_unit_groups.AggregateInfoClient {
		return mockClient
	}
	return mockClient
}

func postcheckGroup(id *string, status string, details *string, failureCount, nodeWithIssues *int64) nsxModel.UpgradeUnitGroupAggregateInfo {
	return nsxModel.UpgradeUnitGroupAggregateInfo{
		Id: id,
		PostUpgradeStatus: &nsxModel.UpgradeChecksExecutionStatus{
			Status:              &status,
			Details:             details,
			FailureCount:        failureCount,
			NodeWithIssuesCount: nodeWithIssues,
		},
	}
}

func postcheckList(groups ...nsxModel.UpgradeUnitGroupAggregateInfo) nsxModel.UpgradeUnitGroupAggregateInfoListResult {
	return nsxModel.UpgradeUnitGroupAggregateInfoListResult{Results: groups}
}

func TestMockDataSourceNsxtUpgradePostCheckRead(t *testing.T) {
	completed := nsxModel.UpgradeChecksExecutionStatus_STATUS_COMPLETED
	inProgress := nsxModel.UpgradeChecksExecutionStatus_STATUS_IN_PROGRESS
	runID := util.GetVerifiableID("run-1", "nsxt_upgrade_run")

	t.Run("all groups pass", func(t *testing.T) {
		mockClient := setupUpgradePostcheckMock(t)
		list := postcheckList(postcheckGroup(strPtr("g1"), completed, nil, int64Ptr(0), nil))
		mockClient.EXPECT().List(gomock.Any(), nil, nil, nil, nil, nil, nil, nil).Return(list, nil).Times(2)

		ds := dataSourceNsxtUpgradePostCheck()
		d := schema.TestResourceDataRaw(t, ds.Schema, upgradePostcheckData(nil))
		require.NoError(t, dataSourceNsxtUpgradePostCheckRead(d, newGoMockProviderClient()))
		assert.Equal(t, runID+"/"+edgeUpgradeGroup, d.Id())
		assert.Empty(t, d.Get("failed_group").([]interface{}))
	})

	t.Run("waits for in-progress groups and records failures", func(t *testing.T) {
		mockClient := setupUpgradePostcheckMock(t)
		final := postcheckList(
			postcheckGroup(strPtr("clean"), completed, nil, nil, nil),
			postcheckGroup(strPtr("failures"), completed, strPtr("2 checks failed"), int64Ptr(2), int64Ptr(9)),
			postcheckGroup(strPtr("nodes"), completed, nil, nil, int64Ptr(3)),
			postcheckGroup(nil, nsxModel.UpgradeChecksExecutionStatus_STATUS_ABORTED, strPtr("aborted"), nil, nil),
		)
		gomock.InOrder(
			mockClient.EXPECT().List(gomock.Any(), nil, nil, nil, nil, nil, nil, nil).
				Return(postcheckList(postcheckGroup(strPtr("clean"), inProgress, nil, nil, nil)), nil),
			mockClient.EXPECT().List(gomock.Any(), nil, nil, nil, nil, nil, nil, nil).Return(final, nil).Times(2),
		)

		ds := dataSourceNsxtUpgradePostCheck()
		d := schema.TestResourceDataRaw(t, ds.Schema, upgradePostcheckData(map[string]interface{}{
			"type": hostUpgradeGroup,
		}))
		require.NoError(t, dataSourceNsxtUpgradePostCheckRead(d, newGoMockProviderClient()))
		assert.Equal(t, runID+"/"+hostUpgradeGroup, d.Id())

		failed := d.Get("failed_group").([]interface{})
		require.Len(t, failed, 3)
		first := failed[0].(map[string]interface{})
		assert.Equal(t, "failures", first["id"])
		assert.Equal(t, completed, first["status"])
		assert.Equal(t, "2 checks failed", first["details"])
		// FailureCount takes precedence over NodeWithIssuesCount.
		assert.Equal(t, 2, first["failure_count"])
		second := failed[1].(map[string]interface{})
		assert.Equal(t, "nodes", second["id"])
		assert.Equal(t, 3, second["failure_count"])
		third := failed[2].(map[string]interface{})
		assert.Equal(t, "", third["id"])
		assert.Equal(t, nsxModel.UpgradeChecksExecutionStatus_STATUS_ABORTED, third["status"])
		assert.Equal(t, "aborted", third["details"])
		assert.Equal(t, 0, third["failure_count"])
	})

	t.Run("list error while polling", func(t *testing.T) {
		mockClient := setupUpgradePostcheckMock(t)
		mockClient.EXPECT().List(gomock.Any(), nil, nil, nil, nil, nil, nil, nil).
			Return(nsxModel.UpgradeUnitGroupAggregateInfoListResult{}, errors.New("boom"))

		ds := dataSourceNsxtUpgradePostCheck()
		d := schema.TestResourceDataRaw(t, ds.Schema, upgradePostcheckData(nil))
		err := dataSourceNsxtUpgradePostCheckRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to retrieve EDGE upgrade post-check results")
		assert.Empty(t, d.Id())
	})

	t.Run("list error after polling", func(t *testing.T) {
		mockClient := setupUpgradePostcheckMock(t)
		gomock.InOrder(
			mockClient.EXPECT().List(gomock.Any(), nil, nil, nil, nil, nil, nil, nil).Return(postcheckList(), nil),
			mockClient.EXPECT().List(gomock.Any(), nil, nil, nil, nil, nil, nil, nil).
				Return(nsxModel.UpgradeUnitGroupAggregateInfoListResult{}, errors.New("boom")),
		)

		ds := dataSourceNsxtUpgradePostCheck()
		d := schema.TestResourceDataRaw(t, ds.Schema, upgradePostcheckData(nil))
		err := dataSourceNsxtUpgradePostCheckRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to retrieve EDGE upgrade post-check result:")
		assert.Empty(t, d.Id())
	})
}
