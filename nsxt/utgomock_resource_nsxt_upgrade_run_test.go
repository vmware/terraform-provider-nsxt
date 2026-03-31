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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/upgrade"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/upgrade/plan"

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
	cliUpgradeUnitGroupsClient = func(_ vapiProtocolClient.Connector) upgrade.UpgradeUnitGroupsClient {
		return mockGroups
	}

	origSettings := cliUpgradeSettingsClient
	cliUpgradeSettingsClient = func(_ vapiProtocolClient.Connector) plan.SettingsClient {
		return mockSettings
	}

	origPlan := cliUpgradePlanClient
	cliUpgradePlanClient = func(_ vapiProtocolClient.Connector) upgrade.PlanClient {
		return mockPlan
	}

	origStatus := cliUpgradeStatusSummaryClient
	cliUpgradeStatusSummaryClient = func(_ vapiProtocolClient.Connector) upgrade.StatusSummaryClient {
		return mockStatus
	}

	origUpgrade := cliUpgradeClient
	cliUpgradeClient = func(_ vapiProtocolClient.Connector) nsx.UpgradeClient {
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
