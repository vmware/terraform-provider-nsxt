//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run:
// mockgen -destination=mocks/nsx/upgrade/UpgradeUnitGroupsClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/upgrade/UpgradeUnitGroupsClient.go UpgradeUnitGroupsClient
// mockgen -destination=mocks/nsx/upgrade/plan/SettingsClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/upgrade/plan/SettingsClient.go SettingsClient
// mockgen -destination=mocks/nsx/upgrade/PlanClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/upgrade/PlanClient.go PlanClient
// mockgen -destination=mocks/nsx/upgrade/StatusSummaryClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/upgrade/StatusSummaryClient.go StatusSummaryClient
// mockgen -destination=mocks/nsx/UpgradeClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/UpgradeClient.go UpgradeClient
// mockgen -destination=mocks/nsx/upgrade/UpgradeUnitGroupsStatusClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/upgrade/UpgradeUnitGroupsStatusClient.go UpgradeUnitGroupsStatusClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/upgrade"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/upgrade/plan"
	"go.uber.org/mock/gomock"

	nsxmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx"
	upgrademocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/upgrade"
	planmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/upgrade/plan"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

func validUpgradePrepareReadyID() string {
	return util.GetVerifiableID("prepare-id", "nsxt_upgrade_prepare_ready")
}

func minimalUpgradeRunData() map[string]interface{} {
	return map[string]interface{}{
		"upgrade_prepare_ready_id": validUpgradePrepareReadyID(),
		"timeout":                  1,
		"interval":                 1,
		"delay":                    0,
		"max_retries":              1,
	}
}

func setupUpgradeRunMocks(ctrl *gomock.Controller) (
	*upgrademocks.MockUpgradeUnitGroupsClient,
	*planmocks.MockSettingsClient,
	*upgrademocks.MockPlanClient,
	*upgrademocks.MockStatusSummaryClient,
	*nsxmocks.MockUpgradeClient,
	*upgrademocks.MockUpgradeUnitGroupsStatusClient,
	func(),
) {
	mockGroups := upgrademocks.NewMockUpgradeUnitGroupsClient(ctrl)
	mockSettings := planmocks.NewMockSettingsClient(ctrl)
	mockPlan := upgrademocks.NewMockPlanClient(ctrl)
	mockStatus := upgrademocks.NewMockStatusSummaryClient(ctrl)
	mockUpgrade := nsxmocks.NewMockUpgradeClient(ctrl)
	mockGroupStatus := upgrademocks.NewMockUpgradeUnitGroupsStatusClient(ctrl)

	origGroups := cliUpgradeUnitGroupsClient
	cliUpgradeUnitGroupsClient = func(_ vapiProtocolClient.Connector) upgradeGroupOps {
		return mockGroups
	}

	origSettings := cliUpgradeSettingsClient
	cliUpgradeSettingsClient = func(_ vapiProtocolClient.Connector) plan.SettingsClient {
		return mockSettings
	}

	origPlan := cliUpgradePlanClient
	cliUpgradePlanClient = func(_ vapiProtocolClient.Connector) upgradePlanOps {
		return mockPlan
	}

	origStatus := cliUpgradeStatusSummaryClient
	cliUpgradeStatusSummaryClient = func(_ vapiProtocolClient.Connector) upgrade.StatusSummaryClient {
		return mockStatus
	}

	origUpgrade := cliUpgradeClient
	cliUpgradeClient = func(_ vapiProtocolClient.Connector) upgradeOps {
		return mockUpgrade
	}

	origGroupStatus := cliUpgradeUnitGroupsStatusClient
	cliUpgradeUnitGroupsStatusClient = func(_ vapiProtocolClient.Connector) upgrade.UpgradeUnitGroupsStatusClient {
		return mockGroupStatus
	}

	restore := func() {
		cliUpgradeUnitGroupsClient = origGroups
		cliUpgradeSettingsClient = origSettings
		cliUpgradePlanClient = origPlan
		cliUpgradeStatusSummaryClient = origStatus
		cliUpgradeClient = origUpgrade
		cliUpgradeUnitGroupsStatusClient = origGroupStatus
	}
	return mockGroups, mockSettings, mockPlan, mockStatus, mockUpgrade, mockGroupStatus, restore
}

func TestMockResourceNsxtUpgradeRunDelete(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Delete is a no-op", func(t *testing.T) {
		res := resourceNsxtUpgradeRun()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradeRunData())
		d.SetId("some-id")

		err := resourceNsxtUpgradeRunDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtUpgradeRunRead(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Read with empty state", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGroups, _, _, mockStatus, _, mockGroupStatus, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		// setUpgradeRunOutput: GroupClient.List() -> empty
		mockGroups.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeUnitGroupListResult{Results: []nsxModel.UpgradeUnitGroup{}}, nil,
		)
		// StatusClient.Get() -> empty component status
		mockStatus.EXPECT().Get(nil, nil, nil).Return(
			nsxModel.UpgradeStatus{ComponentStatus: []nsxModel.ComponentUpgradeStatus{}}, nil,
		)
		// No GroupStatusClient.Getall() calls (no components)
		_ = mockGroupStatus

		res := resourceNsxtUpgradeRun()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradeRunData())
		d.SetId("some-id")

		err := resourceNsxtUpgradeRunRead(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Read with component status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGroups, _, _, mockStatus, _, mockGroupStatus, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		enabled := true
		parallel := false
		pauseAfter := false
		groupID := "edge-group-1"
		groupType := "EDGE"
		mockGroups.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeUnitGroupListResult{
				Results: []nsxModel.UpgradeUnitGroup{
					{
						Id:                        &groupID,
						Enabled:                   &enabled,
						Parallel:                  &parallel,
						PauseAfterEachUpgradeUnit: &pauseAfter,
						Type_:                     &groupType,
					},
				},
			}, nil,
		)

		edgeType := "EDGE"
		edgeStatus := nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS
		componentStatus := nsxModel.ComponentUpgradeStatus{
			ComponentType: &edgeType,
			Status:        &edgeStatus,
		}
		mockStatus.EXPECT().Get(nil, nil, nil).Return(
			nsxModel.UpgradeStatus{ComponentStatus: []nsxModel.ComponentUpgradeStatus{componentStatus}}, nil,
		)

		// GroupStatusClient.Getall() for EDGE component
		mockGroupStatus.EXPECT().Getall(&edgeType, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeUnitGroupStatusListResult{Results: []nsxModel.UpgradeUnitGroupStatus{}}, nil,
		)

		res := resourceNsxtUpgradeRun()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradeRunData())
		d.SetId("some-id")

		err := resourceNsxtUpgradeRunRead(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtUpgradeRunCreate(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Create fails with invalid upgrade_prepare_ready_id", func(t *testing.T) {
		res := resourceNsxtUpgradeRun()
		data := minimalUpgradeRunData()
		data["upgrade_prepare_ready_id"] = "invalid-id-without-hash"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtUpgradeRunCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("Create fails when getTargetVersion returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, _, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		// Also need to mock summary client for getTargetVersion
		origSummary := cliUpgradeSummaryClient
		mockSummary := upgrademocks.NewMockSummaryClient(ctrl)
		cliUpgradeSummaryClient = func(_ vapiProtocolClient.Connector) upgrade.SummaryClient {
			return mockSummary
		}
		defer func() { cliUpgradeSummaryClient = origSummary }()

		// getTargetVersion -> summaryClient.Get() returns summary without TargetVersion
		mockSummary.EXPECT().Get().Return(nsxModel.UpgradeSummary{}, nil)

		res := resourceNsxtUpgradeRun()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradeRunData())

		err := resourceNsxtUpgradeRunCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "target version")
	})
}

func TestMockResourceNsxtUpgradeRunUpdate(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	t.Run("Update fails with invalid upgrade_prepare_ready_id", func(t *testing.T) {
		res := resourceNsxtUpgradeRun()
		data := minimalUpgradeRunData()
		data["upgrade_prepare_ready_id"] = "invalid-id-without-hash"
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId("existing-id")

		err := resourceNsxtUpgradeRunUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid")
	})
}

func TestUnitNsxt_getUpgradeComponentList(t *testing.T) {
	assert.Equal(t, upgradeComponentList, getUpgradeComponentList("4.1.0"))
	assert.Equal(t, upgradeComponentListPost9, getUpgradeComponentList("9.0.0"))
}

func TestUnitNsxt_getPartialUpgradeMap(t *testing.T) {
	res := resourceNsxtUpgradeRun()

	t.Run("no groups defaults to full upgrade", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		m := getPartialUpgradeMap(d, "4.1.0")
		assert.False(t, m[edgeUpgradeGroup])
		assert.False(t, m[hostUpgradeGroup])
	})

	t.Run("disabled group marks partial upgrade", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_group": []interface{}{
				map[string]interface{}{"id": "g1", "enabled": false, "parallel": true, "pause_after_each_upgrade_unit": false},
			},
		})
		m := getPartialUpgradeMap(d, "4.1.0")
		assert.True(t, m[edgeUpgradeGroup])
	})

	t.Run("pause_after_each_upgrade_unit marks partial upgrade", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"host_group": []interface{}{
				map[string]interface{}{"enabled": true, "parallel": true, "pause_after_each_upgrade_unit": true},
			},
		})
		m := getPartialUpgradeMap(d, "4.1.0")
		assert.True(t, m[hostUpgradeGroup])
	})

	t.Run("stage_in_vlcm upgrade_mode marks partial upgrade", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"host_group": []interface{}{
				map[string]interface{}{"enabled": true, "parallel": true, "pause_after_each_upgrade_unit": false, "upgrade_mode": "stage_in_vlcm"},
			},
		})
		m := getPartialUpgradeMap(d, "4.1.0")
		assert.True(t, m[hostUpgradeGroup])
	})

	t.Run("fully enabled group is not partial", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"host_group": []interface{}{
				map[string]interface{}{"enabled": true, "parallel": true, "pause_after_each_upgrade_unit": false, "upgrade_mode": "in_place"},
			},
		})
		m := getPartialUpgradeMap(d, "4.1.0")
		assert.False(t, m[hostUpgradeGroup])
	})
}

func TestMockNsxtGetUpgradeStatus(t *testing.T) {
	t.Run("nil component returns overall status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		overall := nsxModel.ComponentUpgradeStatus_STATUS_IN_PROGRESS
		mockStatus.EXPECT().Get(nil, nil, nil).Return(nsxModel.UpgradeStatus{OverallUpgradeStatus: &overall}, nil)

		status, err := getUpgradeStatus(mockStatus, nil)
		require.NoError(t, err)
		assert.Equal(t, overall, status.Status)
		assert.Equal(t, overall, status.OverallStatus)
	})

	t.Run("component found returns its status and detail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		edge := "EDGE"
		overall := nsxModel.ComponentUpgradeStatus_STATUS_IN_PROGRESS
		compStatus := nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS
		detail := "done"
		mockStatus.EXPECT().Get(&edge, nil, nil).Return(nsxModel.UpgradeStatus{
			OverallUpgradeStatus: &overall,
			ComponentStatus: []nsxModel.ComponentUpgradeStatus{
				{ComponentType: &edge, Status: &compStatus, Details: &detail},
			},
		}, nil)

		status, err := getUpgradeStatus(mockStatus, &edge)
		require.NoError(t, err)
		assert.Equal(t, compStatus, status.Status)
		assert.Equal(t, detail, status.Detail)
	})

	t.Run("component not found errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		edge := "EDGE"
		mockStatus.EXPECT().Get(&edge, nil, nil).Return(nsxModel.UpgradeStatus{ComponentStatus: []nsxModel.ComponentUpgradeStatus{}}, nil)

		_, err := getUpgradeStatus(mockStatus, &edge)
		require.Error(t, err)
	})

	t.Run("Get error is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		edge := "EDGE"
		mockStatus.EXPECT().Get(&edge, nil, nil).Return(nsxModel.UpgradeStatus{}, errors.New("boom"))

		_, err := getUpgradeStatus(mockStatus, &edge)
		require.Error(t, err)
	})
}

func TestUnitNsxt_getPostResetGroupIDFromPreResetList(t *testing.T) {
	preList := nsxModel.UpgradeUnitGroupListResult{Results: []nsxModel.UpgradeUnitGroup{
		{Id: str("old-1"), DisplayName: str("group-a")},
	}}
	postList := nsxModel.UpgradeUnitGroupListResult{Results: []nsxModel.UpgradeUnitGroup{
		{Id: str("new-1"), DisplayName: str("group-a")},
	}}

	t.Run("matches by display name", func(t *testing.T) {
		id, err := getPostResetGroupIDFromPreResetList("old-1", preList, postList)
		require.NoError(t, err)
		assert.Equal(t, "new-1", id)
	})

	t.Run("group id not found in pre-reset list errors", func(t *testing.T) {
		_, err := getPostResetGroupIDFromPreResetList("missing", preList, postList)
		require.Error(t, err)
	})

	t.Run("multiple matches by display name errors", func(t *testing.T) {
		dupPost := nsxModel.UpgradeUnitGroupListResult{Results: []nsxModel.UpgradeUnitGroup{
			{Id: str("new-1"), DisplayName: str("group-a")},
			{Id: str("new-2"), DisplayName: str("group-a")},
		}}
		_, err := getPostResetGroupIDFromPreResetList("old-1", preList, dupPost)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple groups")
	})
}

func TestMockNsxtUpdateComponentUpgradePlanSetting(t *testing.T) {
	res := resourceNsxtUpgradeRun()

	t.Run("empty setting is a no-op", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, mockSettings, _, _, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()
		_ = mockSettings // no calls expected

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		err := updateComponentUpgradePlanSetting(mockSettings, d, edgeUpgradeGroup)
		require.NoError(t, err)
	})

	t.Run("edge component does not modify stop_on_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, mockSettings, _, _, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		mockSettings.EXPECT().Get(edgeUpgradeGroup).Return(nsxModel.UpgradePlanSettings{}, nil)
		mockSettings.EXPECT().Update(edgeUpgradeGroup, gomock.Any()).DoAndReturn(
			func(_ string, settings nsxModel.UpgradePlanSettings) (nsxModel.UpgradePlanSettings, error) {
				assert.True(t, *settings.Parallel)
				assert.Nil(t, settings.PauseOnError)
				return settings, nil
			})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_upgrade_setting": []interface{}{
				map[string]interface{}{"parallel": true, "post_upgrade_check": true},
			},
		})
		err := updateComponentUpgradePlanSetting(mockSettings, d, edgeUpgradeGroup)
		require.NoError(t, err)
	})

	t.Run("host component sets stop_on_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, mockSettings, _, _, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		mockSettings.EXPECT().Get(hostUpgradeGroup).Return(nsxModel.UpgradePlanSettings{}, nil)
		mockSettings.EXPECT().Update(hostUpgradeGroup, gomock.Any()).DoAndReturn(
			func(_ string, settings nsxModel.UpgradePlanSettings) (nsxModel.UpgradePlanSettings, error) {
				assert.False(t, *settings.Parallel)
				require.NotNil(t, settings.PauseOnError)
				assert.True(t, *settings.PauseOnError)
				return settings, nil
			})

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"host_upgrade_setting": []interface{}{
				map[string]interface{}{"parallel": false, "stop_on_error": true, "post_upgrade_check": true},
			},
		})
		err := updateComponentUpgradePlanSetting(mockSettings, d, hostUpgradeGroup)
		require.NoError(t, err)
	})

	t.Run("Get error is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, mockSettings, _, _, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		mockSettings.EXPECT().Get(hostUpgradeGroup).Return(nsxModel.UpgradePlanSettings{}, errors.New("get failed"))

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"host_upgrade_setting": []interface{}{
				map[string]interface{}{"parallel": false, "stop_on_error": true, "post_upgrade_check": true},
			},
		})
		err := updateComponentUpgradePlanSetting(mockSettings, d, hostUpgradeGroup)
		require.Error(t, err)
	})
}

func TestMockNsxtRunPostcheck(t *testing.T) {
	res := resourceNsxtUpgradeRun()

	t.Run("triggers postcheck for components with post_upgrade_check enabled", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUpgrade := nsxmocks.NewMockUpgradeClient(ctrl)

		mockUpgrade.EXPECT().Executepostupgradechecks(edgeUpgradeGroup).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_upgrade_setting": []interface{}{
				map[string]interface{}{"post_upgrade_check": true, "parallel": true},
			},
		})
		runPostcheck(mockUpgrade, d)
	})

	t.Run("skips components without post_upgrade_check", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUpgrade := nsxmocks.NewMockUpgradeClient(ctrl)
		// No EXPECT() calls set up: Executepostupgradechecks must not be called.

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_upgrade_setting": []interface{}{
				map[string]interface{}{"post_upgrade_check": false, "parallel": true},
			},
		})
		runPostcheck(mockUpgrade, d)
	})

	t.Run("skips components with no setting block", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUpgrade := nsxmocks.NewMockUpgradeClient(ctrl)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		runPostcheck(mockUpgrade, d)
	})
}

func TestMockNsxtWaitUpgradeForStatus(t *testing.T) {
	edgeType := edgeUpgradeGroup

	t.Run("reaches target status immediately", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockStatus := upgrademocks.NewMockStatusSummaryClient(ctrl)

		success := nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS
		overall := nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS
		mockStatus.EXPECT().Get(&edgeType, nil, nil).Return(nsxModel.UpgradeStatus{
			OverallUpgradeStatus: &overall,
			ComponentStatus: []nsxModel.ComponentUpgradeStatus{
				{ComponentType: &edgeType, Status: &success},
			},
		}, nil)

		ucs := &upgradeClientSet{StatusClient: mockStatus, Timeout: 2, Interval: 1, MaxRetries: 1}
		status, err := waitUpgradeForStatus(ucs, &edgeType,
			[]string{nsxModel.ComponentUpgradeStatus_STATUS_IN_PROGRESS},
			[]string{nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS}, false)
		require.NoError(t, err)
		assert.Equal(t, nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS, status)
	})

	t.Run("FAILED status is surfaced as an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockStatus := upgrademocks.NewMockStatusSummaryClient(ctrl)

		failed := nsxModel.ComponentUpgradeStatus_STATUS_FAILED
		overall := nsxModel.ComponentUpgradeStatus_STATUS_FAILED
		mockStatus.EXPECT().Get(&edgeType, nil, nil).Return(nsxModel.UpgradeStatus{
			OverallUpgradeStatus: &overall,
			ComponentStatus: []nsxModel.ComponentUpgradeStatus{
				{ComponentType: &edgeType, Status: &failed},
			},
		}, nil).AnyTimes()

		ucs := &upgradeClientSet{StatusClient: mockStatus, Timeout: 2, Interval: 1, MaxRetries: 1}
		_, err := waitUpgradeForStatus(ucs, &edgeType,
			[]string{nsxModel.ComponentUpgradeStatus_STATUS_IN_PROGRESS},
			[]string{nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS, nsxModel.ComponentUpgradeStatus_STATUS_FAILED}, false)
		require.Error(t, err)
	})

	t.Run("times out while stuck pending", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockStatus := upgrademocks.NewMockStatusSummaryClient(ctrl)

		inProgress := nsxModel.ComponentUpgradeStatus_STATUS_IN_PROGRESS
		mockStatus.EXPECT().Get(&edgeType, nil, nil).Return(nsxModel.UpgradeStatus{
			OverallUpgradeStatus: &inProgress,
			ComponentStatus: []nsxModel.ComponentUpgradeStatus{
				{ComponentType: &edgeType, Status: &inProgress},
			},
		}, nil).AnyTimes()

		ucs := &upgradeClientSet{StatusClient: mockStatus, Timeout: 1, Interval: 1, MaxRetries: 1}
		_, err := waitUpgradeForStatus(ucs, &edgeType,
			[]string{nsxModel.ComponentUpgradeStatus_STATUS_IN_PROGRESS},
			[]string{nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS}, false)
		require.Error(t, err)
	})
}

func TestMockNsxtUpdateUpgradeUnitGroups(t *testing.T) {
	res := resourceNsxtUpgradeRun()

	t.Run("predefined group is updated by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGroups, _, _, _, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		mockGroups.EXPECT().List(&edgeUpgradeGroup, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeUnitGroupListResult{}, nil,
		)
		mockGroups.EXPECT().Get("group-1", gomock.Any()).Return(nsxModel.UpgradeUnitGroup{Id: str("group-1")}, nil)
		mockGroups.EXPECT().Update("group-1", gomock.Any()).Return(nsxModel.UpgradeUnitGroup{Id: str("group-1")}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_group": []interface{}{
				map[string]interface{}{"id": "group-1", "enabled": true, "parallel": true, "pause_after_each_upgrade_unit": false},
			},
		})

		hasVLCM := false
		err := updateUpgradeUnitGroups(&upgradeClientSet{GroupClient: mockGroups}, d, edgeUpgradeGroup, nsxModel.UpgradeUnitGroupListResult{}, &hasVLCM)
		require.NoError(t, err)
	})

	t.Run("group not found by id falls back to post-reset lookup by name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGroups, _, _, _, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		preResetList := nsxModel.UpgradeUnitGroupListResult{Results: []nsxModel.UpgradeUnitGroup{
			{Id: str("old-id"), DisplayName: str("hosts-group")},
		}}
		postResetList := nsxModel.UpgradeUnitGroupListResult{Results: []nsxModel.UpgradeUnitGroup{
			{Id: str("new-id"), DisplayName: str("hosts-group")},
		}}

		mockGroups.EXPECT().List(&edgeUpgradeGroup, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(postResetList, nil)
		mockGroups.EXPECT().Get("old-id", gomock.Any()).Return(nsxModel.UpgradeUnitGroup{}, vapiErrors.NotFound{})
		mockGroups.EXPECT().Get("new-id", gomock.Any()).Return(nsxModel.UpgradeUnitGroup{Id: str("new-id")}, nil)
		mockGroups.EXPECT().Update("new-id", gomock.Any()).Return(nsxModel.UpgradeUnitGroup{Id: str("new-id")}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_group": []interface{}{
				map[string]interface{}{"id": "old-id", "enabled": true, "parallel": true, "pause_after_each_upgrade_unit": false},
			},
		})

		hasVLCM := false
		err := updateUpgradeUnitGroups(&upgradeClientSet{GroupClient: mockGroups}, d, edgeUpgradeGroup, preResetList, &hasVLCM)
		require.NoError(t, err)
	})

	t.Run("new custom group is created and reordered after the previous group", func(t *testing.T) {
		// display_name only exists in the host_group schema (getUpgradeGroupSchema(true)),
		// so custom (id-less) groups must be exercised via host_group, not edge_group.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGroups, _, _, _, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		mockGroups.EXPECT().List(&hostUpgradeGroup, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeUnitGroupListResult{}, nil,
		)
		mockGroups.EXPECT().Get("group-1", gomock.Any()).Return(nsxModel.UpgradeUnitGroup{Id: str("group-1")}, nil)
		mockGroups.EXPECT().Update("group-1", gomock.Any()).Return(nsxModel.UpgradeUnitGroup{Id: str("group-1")}, nil)
		mockGroups.EXPECT().Create(gomock.Any()).Return(nsxModel.UpgradeUnitGroup{Id: str("group-2")}, nil)
		mockGroups.EXPECT().Reorder("group-2", gomock.Any()).Return(nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"host_group": []interface{}{
				map[string]interface{}{"id": "group-1", "enabled": true, "parallel": true, "pause_after_each_upgrade_unit": false},
				map[string]interface{}{"id": "", "display_name": "new-custom-group", "enabled": true, "parallel": true, "pause_after_each_upgrade_unit": false},
			},
		})

		hasVLCM := false
		err := updateUpgradeUnitGroups(&upgradeClientSet{GroupClient: mockGroups}, d, hostUpgradeGroup, nsxModel.UpgradeUnitGroupListResult{}, &hasVLCM)
		require.NoError(t, err)
	})

	t.Run("custom group matching an existing display_name errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGroups, _, _, _, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		mockGroups.EXPECT().List(&hostUpgradeGroup, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeUnitGroupListResult{Results: []nsxModel.UpgradeUnitGroup{
				{Id: str("existing-id"), DisplayName: str("dup-name")},
			}}, nil,
		)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"host_group": []interface{}{
				map[string]interface{}{"id": "", "display_name": "dup-name", "enabled": true, "parallel": true, "pause_after_each_upgrade_unit": false},
			},
		})

		hasVLCM := false
		err := updateUpgradeUnitGroups(&upgradeClientSet{GroupClient: mockGroups}, d, hostUpgradeGroup, nsxModel.UpgradeUnitGroupListResult{}, &hasVLCM)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("host group with stage_in_vlcm upgrade_mode sets hasVLCM", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGroups, _, _, _, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		mockGroups.EXPECT().List(&hostUpgradeGroup, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeUnitGroupListResult{}, nil,
		)
		mockGroups.EXPECT().Get("host-group-1", gomock.Any()).Return(nsxModel.UpgradeUnitGroup{Id: str("host-group-1")}, nil)
		mockGroups.EXPECT().Update("host-group-1", gomock.Any()).Return(nsxModel.UpgradeUnitGroup{Id: str("host-group-1")}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"host_group": []interface{}{
				map[string]interface{}{
					"id": "host-group-1", "enabled": true, "parallel": true, "pause_after_each_upgrade_unit": false,
					"upgrade_mode": "stage_in_vlcm", "maintenance_mode_config_vsan_mode": "no_action",
					"maintenance_mode_config_evacuate_powered_off_vms": false, "rebootless_upgrade": true,
				},
			},
		})

		hasVLCM := false
		err := updateUpgradeUnitGroups(&upgradeClientSet{GroupClient: mockGroups}, d, hostUpgradeGroup, nsxModel.UpgradeUnitGroupListResult{}, &hasVLCM)
		require.NoError(t, err)
		assert.True(t, hasVLCM)
	})

	t.Run("List error is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGroups, _, _, _, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		mockGroups.EXPECT().List(&edgeUpgradeGroup, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeUnitGroupListResult{}, errors.New("list failed"),
		)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		hasVLCM := false
		err := updateUpgradeUnitGroups(&upgradeClientSet{GroupClient: mockGroups}, d, edgeUpgradeGroup, nsxModel.UpgradeUnitGroupListResult{}, &hasVLCM)
		require.Error(t, err)
	})
}

func TestMockNsxtPrepareUpgrade(t *testing.T) {
	res := resourceNsxtUpgradeRun()

	t.Run("no changes to any component is a no-op", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		hasVLCM := false
		err := prepareUpgrade(&upgradeClientSet{}, d, "4.1.0", &hasVLCM)
		require.NoError(t, err)
	})

	t.Run("component already succeeded is skipped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		success := nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS
		mockStatus.EXPECT().Get(&edgeUpgradeGroup, nil, nil).Return(nsxModel.UpgradeStatus{
			OverallUpgradeStatus: &success,
			ComponentStatus:      []nsxModel.ComponentUpgradeStatus{{ComponentType: &edgeUpgradeGroup, Status: &success}},
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_group": []interface{}{
				map[string]interface{}{"id": "group-1", "enabled": true, "parallel": true, "pause_after_each_upgrade_unit": false},
			},
		})

		ucs := &upgradeClientSet{StatusClient: mockStatus, Timeout: 2, Interval: 1, MaxRetries: 1}
		hasVLCM := false
		err := prepareUpgrade(ucs, d, "4.1.0", &hasVLCM)
		require.NoError(t, err)
	})

	t.Run("in-progress component is paused, reset, and its groups updated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGroups, mockSettings, mockPlan, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		inProgress := nsxModel.ComponentUpgradeStatus_STATUS_IN_PROGRESS
		notStarted := nsxModel.ComponentUpgradeStatus_STATUS_NOT_STARTED
		// First check: component is IN_PROGRESS
		mockStatus.EXPECT().Get(&edgeUpgradeGroup, nil, nil).Return(nsxModel.UpgradeStatus{
			OverallUpgradeStatus: &inProgress,
			ComponentStatus:      []nsxModel.ComponentUpgradeStatus{{ComponentType: &edgeUpgradeGroup, Status: &inProgress}},
		}, nil)
		mockPlan.EXPECT().Pause().Return(nil)
		// waitUpgradeForStatus polls until component reaches a static (non-pending) status.
		// Must resolve to something other than SUCCESS: prepareUpgrade treats a SUCCESS
		// result here as an error (a concurrent upgrade may have completed the component).
		mockStatus.EXPECT().Get(&edgeUpgradeGroup, nil, nil).Return(nsxModel.UpgradeStatus{
			OverallUpgradeStatus: &notStarted,
			ComponentStatus:      []nsxModel.ComponentUpgradeStatus{{ComponentType: &edgeUpgradeGroup, Status: &notStarted}},
		}, nil)
		// preResetGroupList
		mockGroups.EXPECT().List(&edgeUpgradeGroup, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeUnitGroupListResult{}, nil,
		)
		mockPlan.EXPECT().Reset(edgeUpgradeGroup).Return(nil)
		// updateUpgradeUnitGroups: List again (post reset)
		mockGroups.EXPECT().List(&edgeUpgradeGroup, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeUnitGroupListResult{}, nil,
		)
		mockGroups.EXPECT().Get("group-1", gomock.Any()).Return(nsxModel.UpgradeUnitGroup{Id: str("group-1")}, nil)
		mockGroups.EXPECT().Update("group-1", gomock.Any()).Return(nsxModel.UpgradeUnitGroup{Id: str("group-1")}, nil)
		// updateComponentUpgradePlanSetting: no setting block configured -> no-op
		_ = mockSettings

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_group": []interface{}{
				map[string]interface{}{"id": "group-1", "enabled": true, "parallel": true, "pause_after_each_upgrade_unit": false},
			},
		})

		ucs := &upgradeClientSet{StatusClient: mockStatus, GroupClient: mockGroups, PlanClient: mockPlan, SettingClient: mockSettings, Timeout: 2, Interval: 1, MaxRetries: 1}
		hasVLCM := false
		err := prepareUpgrade(ucs, d, "4.1.0", &hasVLCM)
		require.NoError(t, err)
	})

	t.Run("unexpected SUCCESS after wait errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		notStarted := nsxModel.ComponentUpgradeStatus_STATUS_NOT_STARTED
		success := nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS
		mockStatus.EXPECT().Get(&edgeUpgradeGroup, nil, nil).Return(nsxModel.UpgradeStatus{
			OverallUpgradeStatus: &notStarted,
			ComponentStatus:      []nsxModel.ComponentUpgradeStatus{{ComponentType: &edgeUpgradeGroup, Status: &notStarted}},
		}, nil)
		mockStatus.EXPECT().Get(&edgeUpgradeGroup, nil, nil).Return(nsxModel.UpgradeStatus{
			OverallUpgradeStatus: &success,
			ComponentStatus:      []nsxModel.ComponentUpgradeStatus{{ComponentType: &edgeUpgradeGroup, Status: &success}},
		}, nil)

		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"edge_group": []interface{}{
				map[string]interface{}{"id": "group-1", "enabled": true, "parallel": true, "pause_after_each_upgrade_unit": false},
			},
		})

		ucs := &upgradeClientSet{StatusClient: mockStatus, Timeout: 2, Interval: 1, MaxRetries: 1}
		hasVLCM := false
		err := prepareUpgrade(ucs, d, "4.1.0", &hasVLCM)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status")
	})
}

func TestMockNsxtRunUpgrade(t *testing.T) {
	partialMap := map[string]bool{edgeUpgradeGroup: false, hostUpgradeGroup: false, mpUpgradeGroup: false}

	t.Run("all components already succeeded is a no-op", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		success := nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS
		mockStatus.EXPECT().Get(gomock.Any(), nil, nil).Return(nsxModel.UpgradeStatus{
			OverallUpgradeStatus: &success,
			ComponentStatus: []nsxModel.ComponentUpgradeStatus{
				{ComponentType: &edgeUpgradeGroup, Status: &success},
				{ComponentType: &hostUpgradeGroup, Status: &success},
				{ComponentType: &mpUpgradeGroup, Status: &success},
			},
		}, nil).AnyTimes()

		ucs := &upgradeClientSet{StatusClient: mockStatus, Timeout: 2, Interval: 1, MaxRetries: 1}
		err := runUpgrade(ucs, partialMap, "4.1.0", false, true)
		require.NoError(t, err)
	})

	t.Run("upgrades a not-started component to success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, mockPlan, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		notStarted := nsxModel.ComponentUpgradeStatus_STATUS_NOT_STARTED
		success := nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS
		statusMap := map[string]*string{
			edgeUpgradeGroup: &notStarted,
			hostUpgradeGroup: &success,
			mpUpgradeGroup:   &success,
		}
		mockStatus.EXPECT().Get(gomock.Any(), nil, nil).DoAndReturn(
			func(component *string, _, _ interface{}) (nsxModel.UpgradeStatus, error) {
				if component == nil {
					// Overall-status check between components: report stable/SUCCESS so
					// runUpgrade's inter-component wait resolves immediately.
					return nsxModel.UpgradeStatus{OverallUpgradeStatus: &success}, nil
				}
				st := statusMap[*component]
				return nsxModel.UpgradeStatus{
					OverallUpgradeStatus: st,
					ComponentStatus:      []nsxModel.ComponentUpgradeStatus{{ComponentType: component, Status: st}},
				}, nil
			}).AnyTimes()
		mockPlan.EXPECT().Upgrade(&edgeUpgradeGroup).DoAndReturn(func(_ *string) error {
			notStarted = nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS
			return nil
		})

		ucs := &upgradeClientSet{StatusClient: mockStatus, PlanClient: mockPlan, Timeout: 2, Interval: 1, MaxRetries: 1}
		err := runUpgrade(ucs, partialMap, "4.1.0", false, true)
		require.NoError(t, err)
	})

	t.Run("PlanClient.Upgrade error is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, mockPlan, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		notStarted := nsxModel.ComponentUpgradeStatus_STATUS_NOT_STARTED
		mockStatus.EXPECT().Get(gomock.Any(), nil, nil).DoAndReturn(
			func(component *string, _, _ interface{}) (nsxModel.UpgradeStatus, error) {
				if component == nil {
					// Overall-status check between components resolves immediately.
					return nsxModel.UpgradeStatus{OverallUpgradeStatus: &notStarted}, nil
				}
				return nsxModel.UpgradeStatus{
					OverallUpgradeStatus: &notStarted,
					ComponentStatus:      []nsxModel.ComponentUpgradeStatus{{ComponentType: component, Status: &notStarted}},
				}, nil
			}).AnyTimes()
		mockPlan.EXPECT().Upgrade(&edgeUpgradeGroup).Return(errors.New("upgrade failed"))

		ucs := &upgradeClientSet{StatusClient: mockStatus, PlanClient: mockPlan, Timeout: 2, Interval: 1, MaxRetries: 1}
		err := runUpgrade(ucs, map[string]bool{edgeUpgradeGroup: false}, "8.9.9", false, true)
		require.Error(t, err)
	})

	t.Run("skips finalize component when finalizeUpgrade is false", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		success := nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS
		mockStatus.EXPECT().Get(gomock.Any(), nil, nil).Return(nsxModel.UpgradeStatus{
			OverallUpgradeStatus: &success,
			ComponentStatus: []nsxModel.ComponentUpgradeStatus{
				{ComponentType: &mpUpgradeGroup, Status: &success},
				{ComponentType: &edgeUpgradeGroup, Status: &success},
				{ComponentType: &hostUpgradeGroup, Status: &success},
			},
		}, nil).AnyTimes()

		ucs := &upgradeClientSet{StatusClient: mockStatus, Timeout: 2, Interval: 1, MaxRetries: 1}
		// getUpgradeComponentList("9.0.0") includes finalizeUpgradeGroup; with finalizeUpgrade=false it must never
		// be queried since the loop `continue`s before calling getUpgradeStatus for it.
		err := runUpgrade(ucs, map[string]bool{}, "9.0.0", false, false)
		require.NoError(t, err)
	})
}
