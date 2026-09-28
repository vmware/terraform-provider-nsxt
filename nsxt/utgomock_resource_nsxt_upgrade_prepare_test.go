//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run:
// mockgen -destination=mocks/nsx/upgrade/BundlesClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/upgrade/BundlesClient.go BundlesClient
// mockgen -destination=mocks/nsx/upgrade/eula/AcceptClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/upgrade/eula/AcceptClient.go AcceptClient
// mockgen -destination=mocks/nsx/upgrade/bundles/UploadStatusClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/upgrade/bundles/UploadStatusClient.go UploadStatusClient
// mockgen -destination=mocks/nsx/upgrade/UcUpgradeStatusClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/upgrade/UcUpgradeStatusClient.go UcUpgradeStatusClient
// mockgen -destination=mocks/nsx/upgrade/SummaryClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/upgrade/SummaryClient.go SummaryClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/upgrade"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/upgrade/bundles"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/upgrade/eula"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/upgrade/pre_upgrade_checks"
	"go.uber.org/mock/gomock"

	upgrademocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/upgrade"
	bundlesmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/upgrade/bundles"
	eulamocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/upgrade/eula"
	preupgchecks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/upgrade/pre_upgrade_checks"
)

func minimalUpgradePrepareData() map[string]interface{} {
	return map[string]interface{}{
		"accept_user_agreement": true,
	}
}

func upgradeSummaryNotStarted() nsxModel.UpgradeSummary {
	notStarted := nsxModel.UpgradeSummary_UPGRADE_STATUS_NOT_STARTED
	ucUpdated := false
	targetVersion := "4.1.0"
	return nsxModel.UpgradeSummary{
		UpgradeStatus:             &notStarted,
		UpgradeCoordinatorUpdated: &ucUpdated,
		TargetVersion:             &targetVersion,
	}
}

func upgradeSummaryUCUpdated() nsxModel.UpgradeSummary {
	notStarted := nsxModel.UpgradeSummary_UPGRADE_STATUS_NOT_STARTED
	ucUpdated := true
	targetVersion := "4.1.0"
	return nsxModel.UpgradeSummary{
		UpgradeStatus:             &notStarted,
		UpgradeCoordinatorUpdated: &ucUpdated,
		TargetVersion:             &targetVersion,
	}
}

func setupUpgradePrepareMocks(ctrl *gomock.Controller) (
	*upgrademocks.MockSummaryClient,
	*upgrademocks.MockBundlesClient,
	*eulamocks.MockAcceptClient,
	*bundlesmocks.MockUploadStatusClient,
	*upgrademocks.MockUcUpgradeStatusClient,
	*preupgchecks.MockFailuresClient,
	func(),
) {
	mockSummary := upgrademocks.NewMockSummaryClient(ctrl)
	mockBundles := upgrademocks.NewMockBundlesClient(ctrl)
	mockEula := eulamocks.NewMockAcceptClient(ctrl)
	mockUploadStatus := bundlesmocks.NewMockUploadStatusClient(ctrl)
	mockUcStatus := upgrademocks.NewMockUcUpgradeStatusClient(ctrl)
	mockFailures := preupgchecks.NewMockFailuresClient(ctrl)

	origSummary := cliUpgradeSummaryClient
	cliUpgradeSummaryClient = func(_ vapiProtocolClient.Connector) upgrade.SummaryClient {
		return mockSummary
	}

	origBundles := cliUpgradeBundlesClient
	cliUpgradeBundlesClient = func(_ vapiProtocolClient.Connector) upgrade.BundlesClient {
		return mockBundles
	}

	origEula := cliEulaAcceptClient
	cliEulaAcceptClient = func(_ vapiProtocolClient.Connector) eula.AcceptClient {
		return mockEula
	}

	origUploadStatus := cliBundlesUploadStatusClient
	cliBundlesUploadStatusClient = func(_ vapiProtocolClient.Connector) bundles.UploadStatusClient {
		return mockUploadStatus
	}

	origUcStatus := cliUcUpgradeStatusClient
	cliUcUpgradeStatusClient = func(_ vapiProtocolClient.Connector) upgrade.UcUpgradeStatusClient {
		return mockUcStatus
	}

	origFailures := cliPreUpgradeChecksFailuresClient
	cliPreUpgradeChecksFailuresClient = func(_ vapiProtocolClient.Connector) pre_upgrade_checks.FailuresClient {
		return mockFailures
	}

	restore := func() {
		cliUpgradeSummaryClient = origSummary
		cliUpgradeBundlesClient = origBundles
		cliEulaAcceptClient = origEula
		cliBundlesUploadStatusClient = origUploadStatus
		cliUcUpgradeStatusClient = origUcStatus
		cliPreUpgradeChecksFailuresClient = origFailures
	}
	return mockSummary, mockBundles, mockEula, mockUploadStatus, mockUcStatus, mockFailures, restore
}

func TestMockResourceNsxtUpgradePrepareDelete(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Delete is a no-op", func(t *testing.T) {
		res := resourceNsxtUpgradePrepare()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradePrepareData())
		d.SetId("some-id")

		err := resourceNsxtUpgradePrepareDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtUpgradePrepareRead(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Read when UC not upgraded (skip precheck)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, mockFailures, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		// getSummaryInfo: UpgradeCoordinatorUpdated=false -> precheckNeeded=false
		mockSummary.EXPECT().Get().Return(upgradeSummaryNotStarted(), nil)
		// getPrecheckErrors(nil)
		mockFailures.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeCheckFailureListResult{Results: []nsxModel.UpgradeCheckFailure{}}, nil,
		)

		res := resourceNsxtUpgradePrepare()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradePrepareData())
		d.SetId("some-id")

		err := resourceNsxtUpgradePrepareRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "4.1.0", d.Get("target_version"))
	})
}

func TestMockResourceNsxtUpgradePrepareCreate(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	// Create invokes prepareForUpgrade which uploads bundle, accepts EULA, upgrades UC,
	// then calls Read. When HasChange is false (TestResourceDataRaw), bundle upload returns early
	// after getting summary. UC is simulated as already upgraded.
	t.Run("Create with UC already upgraded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, mockEula, _, _, mockFailures, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		// uploadUpgradeBundle(UPGRADE): always calls summaryClient.Get() first
		// HasChange("upgrade_bundle_url") = false for TestResourceDataRaw -> returns early
		mockSummary.EXPECT().Get().Return(upgradeSummaryNotStarted(), nil)

		// acceptUserAgreement
		mockEula.EXPECT().Create().Return(nil)

		// upgradeUc: Get() -> UpgradeCoordinatorUpdated=true -> return early
		mockSummary.EXPECT().Get().Return(upgradeSummaryUCUpdated(), nil)

		// Read -> getSummaryInfo: UC not updated -> precheckNeeded=false
		mockSummary.EXPECT().Get().Return(upgradeSummaryNotStarted(), nil)

		// getPrecheckErrors(nil)
		mockFailures.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeCheckFailureListResult{Results: []nsxModel.UpgradeCheckFailure{}}, nil,
		)

		res := resourceNsxtUpgradePrepare()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradePrepareData())

		err := resourceNsxtUpgradePrepareCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})

	t.Run("Create fails when accept_user_agreement is false", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		// uploadUpgradeBundle(UPGRADE): summaryClient.Get() first, then HasChange=false -> return early
		mockSummary.EXPECT().Get().Return(upgradeSummaryNotStarted(), nil)

		data := minimalUpgradePrepareData()
		data["accept_user_agreement"] = false

		res := resourceNsxtUpgradePrepare()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtUpgradePrepareCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user agreement")
	})
}

func TestMockResourceNsxtUpgradePrepareUpdate(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()

	t.Run("Update with UC already upgraded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, mockEula, _, _, mockFailures, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		mockSummary.EXPECT().Get().Return(upgradeSummaryNotStarted(), nil)
		mockEula.EXPECT().Create().Return(nil)
		mockSummary.EXPECT().Get().Return(upgradeSummaryUCUpdated(), nil)
		mockSummary.EXPECT().Get().Return(upgradeSummaryNotStarted(), nil)
		mockFailures.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeCheckFailureListResult{Results: []nsxModel.UpgradeCheckFailure{}}, nil,
		)

		res := resourceNsxtUpgradePrepare()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradePrepareData())
		d.SetId("existing-id")

		err := resourceNsxtUpgradePrepareUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "4.1.0", d.Get("target_version"))
	})

	t.Run("Update fails when accept_user_agreement is false", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		mockSummary.EXPECT().Get().Return(upgradeSummaryNotStarted(), nil)

		data := minimalUpgradePrepareData()
		data["accept_user_agreement"] = false
		res := resourceNsxtUpgradePrepare()
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId("existing-id")

		err := resourceNsxtUpgradePrepareUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user agreement")
	})
}

func TestUnitNsxt_isVCF9HostUpgrade_belowVCF9(t *testing.T) {
	// No API call should happen when target version is below 9.0.0.
	ok, err := isVCF9HostUpgrade(newGoMockProviderClient(), "8.9.9")
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestMockNsxtIsVCF9HostUpgrade(t *testing.T) {
	t.Run("Get error is propagated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		mockStatus.EXPECT().Get(nil, nil, nil).Return(nsxModel.UpgradeStatus{}, errors.New("mock API error"))

		_, err := isVCF9HostUpgrade(newGoMockProviderClient(), "9.0.0")
		require.Error(t, err)
	})

	t.Run("not paused returns false", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		inProgress := nsxModel.UpgradeStatus_OVERALL_UPGRADE_STATUS_IN_PROGRESS
		mockStatus.EXPECT().Get(nil, nil, nil).Return(nsxModel.UpgradeStatus{OverallUpgradeStatus: &inProgress}, nil)

		ok, err := isVCF9HostUpgrade(newGoMockProviderClient(), "9.0.0")
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("paused with host pre-upgraded returns true", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		paused := nsxModel.UpgradeStatus_OVERALL_UPGRADE_STATUS_PAUSED
		success := nsxModel.ComponentUpgradeStatus_STATUS_SUCCESS
		notStarted := nsxModel.ComponentUpgradeStatus_STATUS_NOT_STARTED
		hostType := hostUpgradeGroup
		finalizeType := finalizeUpgradeGroup
		edgeType := "EDGE"
		mockStatus.EXPECT().Get(nil, nil, nil).Return(nsxModel.UpgradeStatus{
			OverallUpgradeStatus: &paused,
			ComponentStatus: []nsxModel.ComponentUpgradeStatus{
				{ComponentType: &hostType, Status: &success},
				{ComponentType: &finalizeType, Status: &notStarted},
				{ComponentType: &edgeType, Status: &notStarted},
			},
		}, nil)

		ok, err := isVCF9HostUpgrade(newGoMockProviderClient(), "9.0.0")
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("paused but host not yet upgraded returns false", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockStatus, _, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		paused := nsxModel.UpgradeStatus_OVERALL_UPGRADE_STATUS_PAUSED
		notStarted := nsxModel.ComponentUpgradeStatus_STATUS_NOT_STARTED
		hostType := hostUpgradeGroup
		mockStatus.EXPECT().Get(nil, nil, nil).Return(nsxModel.UpgradeStatus{
			OverallUpgradeStatus: &paused,
			ComponentStatus: []nsxModel.ComponentUpgradeStatus{
				{ComponentType: &hostType, Status: &notStarted},
			},
		}, nil)

		ok, err := isVCF9HostUpgrade(newGoMockProviderClient(), "9.0.0")
		require.NoError(t, err)
		assert.False(t, ok)
	})
}

func TestMockNsxtWaitForBundleUpload(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockUploadStatus, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		status := nsxModel.UpgradeBundleUploadStatus_STATUS_SUCCESS
		mockUploadStatus.EXPECT().Get("bundle-1").Return(nsxModel.UpgradeBundleUploadStatus{Status: &status}, nil)

		err := waitForBundleUpload(newGoMockProviderClient(), "bundle-1", 5)
		require.NoError(t, err)
	})

	t.Run("failure surfaces detailed status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockUploadStatus, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		status := nsxModel.UpgradeBundleUploadStatus_STATUS_FAILED
		detail := "disk full"
		mockUploadStatus.EXPECT().Get("bundle-1").Return(nsxModel.UpgradeBundleUploadStatus{Status: &status, DetailedStatus: &detail}, nil)

		err := waitForBundleUpload(newGoMockProviderClient(), "bundle-1", 5)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "disk full")
	})
}

func TestMockNsxtWaitForUcUpgrade(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, _, mockUcStatus, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		state := nsxModel.UcUpgradeStatus_STATE_SUCCESS
		mockUcStatus.EXPECT().Get().Return(nsxModel.UcUpgradeStatus{State: &state}, nil)

		err := waitForUcUpgrade(newGoMockProviderClient(), 5)
		require.NoError(t, err)
	})

	t.Run("FAILED state is a reached target state, not an error", func(t *testing.T) {
		// Unlike waitForBundleUpload, waitForUcUpgrade's Refresh func does not
		// turn a FAILED state into an error when Get() itself succeeds -- FAILED
		// is simply one of the configured target states.
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, _, mockUcStatus, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		state := nsxModel.UcUpgradeStatus_STATE_FAILED
		mockUcStatus.EXPECT().Get().Return(nsxModel.UcUpgradeStatus{State: &state}, nil)

		err := waitForUcUpgrade(newGoMockProviderClient(), 5)
		require.NoError(t, err)
	})

	t.Run("Get error propagates once retries are exhausted", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, _, mockUcStatus, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		mockUcStatus.EXPECT().Get().Return(nsxModel.UcUpgradeStatus{}, errors.New("connection reset")).AnyTimes()

		err := waitForUcUpgrade(newGoMockProviderClient(), 2)
		require.Error(t, err)
	})
}

func TestMockNsxtExecutePreupgradeChecks(t *testing.T) {
	t.Run("Executepreupgradechecks failure short-circuits", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, _, mockUpgrade, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		mockUpgrade.EXPECT().Executepreupgradechecks(nil, nil, nil, nil, nil, nil).Return(errors.New("mock API error"))

		res := resourceNsxtUpgradePrepare()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradePrepareData())

		err := executePreupgradeChecks(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("waits for each component type to complete", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, mockStatus, mockUpgrade, _, restore := setupUpgradeRunMocks(ctrl)
		defer restore()

		mockUpgrade.EXPECT().Executepreupgradechecks(nil, nil, nil, nil, nil, nil).Return(nil)

		completed := nsxModel.UpgradeChecksExecutionStatus_STATUS_COMPLETED
		for _, ct := range precheckComponentTypes {
			componentType := ct
			mockStatus.EXPECT().Get(&componentType, nil, nil).Return(nsxModel.UpgradeStatus{
				ComponentStatus: []nsxModel.ComponentUpgradeStatus{
					{PreUpgradeStatus: &nsxModel.UpgradeChecksExecutionStatus{Status: &completed}},
				},
			}, nil)
		}

		res := resourceNsxtUpgradePrepare()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradePrepareData())

		err := executePreupgradeChecks(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockNsxtUploadUpgradeBundle(t *testing.T) {
	res := resourceNsxtUpgradePrepare()
	upgradeType := nsxModel.UpgradeBundleFetchRequest_BUNDLE_TYPE_UPGRADE
	precheckType := nsxModel.UpgradeBundleFetchRequest_BUNDLE_TYPE_PRE_UPGRADE

	newData := func(t *testing.T, extra map[string]interface{}) *schema.ResourceData {
		data := minimalUpgradePrepareData()
		data["bundle_upload_timeout"] = 5
		for k, v := range extra {
			data[k] = v
		}
		return schema.TestResourceDataRaw(t, res.Schema, data)
	}

	t.Run("upgrade bundle is fetched and waited on", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, mockBundles, _, mockUploadStatus, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		bundleID := "bundle-1"
		success := nsxModel.UpgradeBundleUploadStatus_STATUS_SUCCESS
		mockSummary.EXPECT().Get().Return(upgradeSummaryNotStarted(), nil)
		mockBundles.EXPECT().Create(gomock.Any(), nil).
			DoAndReturn(func(req nsxModel.UpgradeBundleFetchRequest, _ *bool) (nsxModel.UpgradeBundleId, error) {
				assert.Equal(t, "https://bundles.example.com/upgrade.mub", *req.Url)
				assert.Equal(t, upgradeType, *req.BundleType)
				return nsxModel.UpgradeBundleId{BundleId: &bundleID}, nil
			})
		mockUploadStatus.EXPECT().Get(bundleID).Return(nsxModel.UpgradeBundleUploadStatus{Status: &success}, nil)

		d := newData(t, map[string]interface{}{"upgrade_bundle_url": "https://bundles.example.com/upgrade.mub"})
		require.NoError(t, uploadUpgradeBundle(d, newGoMockProviderClient(), upgradeType))
	})

	t.Run("bundle Create error is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, mockBundles, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		mockSummary.EXPECT().Get().Return(upgradeSummaryNotStarted(), nil)
		mockBundles.EXPECT().Create(gomock.Any(), nil).Return(nsxModel.UpgradeBundleId{}, errors.New("fetch refused"))

		d := newData(t, map[string]interface{}{"upgrade_bundle_url": "https://bundles.example.com/upgrade.mub"})
		err := uploadUpgradeBundle(d, newGoMockProviderClient(), upgradeType)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "fetch refused")
	})

	t.Run("missing bundle ID is reported as invalid bundle", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, mockBundles, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		mockSummary.EXPECT().Get().Return(upgradeSummaryNotStarted(), nil)
		mockBundles.EXPECT().Create(gomock.Any(), nil).Return(nsxModel.UpgradeBundleId{}, nil)

		d := newData(t, map[string]interface{}{"upgrade_bundle_url": "https://bundles.example.com/upgrade.mub"})
		err := uploadUpgradeBundle(d, newGoMockProviderClient(), upgradeType)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "apparently invalid")
	})

	t.Run("summary Get error is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		mockSummary.EXPECT().Get().Return(nsxModel.UpgradeSummary{}, errors.New("summary unavailable"))

		d := newData(t, nil)
		require.Error(t, uploadUpgradeBundle(d, newGoMockProviderClient(), upgradeType))
	})

	t.Run("precheck bundle already uploaded for the same version is skipped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		summary := upgradeSummaryNotStarted()
		existing := "4.2.0.0.0.12345"
		summary.PreUpgradeBundleVersion = &existing
		mockSummary.EXPECT().Get().Return(summary, nil)

		d := newData(t, map[string]interface{}{"version": "4.2.0", "precheck_bundle_url": "https://bundles.example.com/pre.pub"})
		require.NoError(t, uploadUpgradeBundle(d, newGoMockProviderClient(), precheckType))
	})

	t.Run("precheck bundle with unparsable version is still uploaded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, mockBundles, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		summary := upgradeSummaryNotStarted()
		existing := "4.2.0.0.0.12345"
		summary.PreUpgradeBundleVersion = &existing
		mockSummary.EXPECT().Get().Return(summary, nil)
		mockBundles.EXPECT().Create(gomock.Any(), nil).
			DoAndReturn(func(req nsxModel.UpgradeBundleFetchRequest, _ *bool) (nsxModel.UpgradeBundleId, error) {
				assert.Equal(t, precheckType, *req.BundleType)
				assert.Equal(t, "https://bundles.example.com/pre.pub", *req.Url)
				return nsxModel.UpgradeBundleId{}, errors.New("stop here")
			})

		d := newData(t, map[string]interface{}{"version": "4.2", "precheck_bundle_url": "https://bundles.example.com/pre.pub"})
		require.Error(t, uploadUpgradeBundle(d, newGoMockProviderClient(), precheckType))
	})
}

func TestMockNsxtUploadPrecheckAndUpgradeBundle(t *testing.T) {
	res := resourceNsxtUpgradePrepare()

	t.Run("precheck bundle upload failure is wrapped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, mockBundles, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		mockSummary.EXPECT().Get().Return(upgradeSummaryNotStarted(), nil)
		mockBundles.EXPECT().Create(gomock.Any(), nil).Return(nsxModel.UpgradeBundleId{}, errors.New("unreachable"))

		data := minimalUpgradePrepareData()
		data["precheck_bundle_url"] = "https://bundles.example.com/pre.pub"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := uploadPrecheckAndUpgradeBundle(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "precheck bundle")
	})

	t.Run("upgrade bundle upload failure is wrapped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		mockSummary.EXPECT().Get().Return(nsxModel.UpgradeSummary{}, errors.New("summary unavailable"))

		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradePrepareData())

		err := uploadPrecheckAndUpgradeBundle(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "upgrade bundle")
	})
}

func TestMockNsxtUpgradeUc(t *testing.T) {
	res := resourceNsxtUpgradePrepare()
	newData := func(t *testing.T) *schema.ResourceData {
		data := minimalUpgradePrepareData()
		data["uc_upgrade_timeout"] = 5
		return schema.TestResourceDataRaw(t, res.Schema, data)
	}

	t.Run("triggers UC upgrade and waits for it", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, mockUcStatus, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()
		_, _, _, _, mockUpgrade, _, restoreRun := setupUpgradeRunMocks(ctrl)
		defer restoreRun()

		success := nsxModel.UcUpgradeStatus_STATE_SUCCESS
		gomock.InOrder(
			mockSummary.EXPECT().Get().Return(upgradeSummaryNotStarted(), nil),
			mockUpgrade.EXPECT().Upgradeuc().Return(nil),
			mockUcStatus.EXPECT().Get().Return(nsxModel.UcUpgradeStatus{State: &success}, nil),
		)

		require.NoError(t, upgradeUc(newData(t), newGoMockProviderClient()))
	})

	t.Run("Upgradeuc error is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()
		_, _, _, _, mockUpgrade, _, restoreRun := setupUpgradeRunMocks(ctrl)
		defer restoreRun()

		mockSummary.EXPECT().Get().Return(upgradeSummaryNotStarted(), nil)
		mockUpgrade.EXPECT().Upgradeuc().Return(errors.New("uc upgrade rejected"))

		require.Error(t, upgradeUc(newData(t), newGoMockProviderClient()))
	})

	t.Run("summary Get error is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		mockSummary.EXPECT().Get().Return(nsxModel.UpgradeSummary{}, errors.New("summary unavailable"))

		require.Error(t, upgradeUc(newData(t), newGoMockProviderClient()))
	})
}

func TestMockNsxtGetSummaryInfo(t *testing.T) {
	ucUpdated := true

	t.Run("UC upgraded and upgrade not started needs precheck", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		mockSummary.EXPECT().Get().Return(upgradeSummaryUCUpdated(), nil)

		version, needed, err := getSummaryInfo(newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "4.1.0", version)
		assert.True(t, needed)
	})

	t.Run("upgrade in progress below 9.0 skips precheck", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		inProgress := nsxModel.UpgradeSummary_UPGRADE_STATUS_IN_PROGRESS
		target := "4.1.0"
		mockSummary.EXPECT().Get().Return(nsxModel.UpgradeSummary{
			UpgradeCoordinatorUpdated: &ucUpdated,
			UpgradeStatus:             &inProgress,
			TargetVersion:             &target,
		}, nil)

		_, needed, err := getSummaryInfo(newGoMockProviderClient())
		require.NoError(t, err)
		assert.False(t, needed)
	})

	t.Run("VCF9 host pre-upgrade still needs precheck", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()
		_, _, _, mockStatus, _, _, restoreRun := setupUpgradeRunMocks(ctrl)
		defer restoreRun()

		paused := nsxModel.UpgradeSummary_UPGRADE_STATUS_PAUSED
		target := "9.0.0"
		mockSummary.EXPECT().Get().Return(nsxModel.UpgradeSummary{
			UpgradeCoordinatorUpdated: &ucUpdated,
			UpgradeStatus:             &paused,
			TargetVersion:             &target,
		}, nil)
		overallPaused := nsxModel.UpgradeStatus_OVERALL_UPGRADE_STATUS_PAUSED
		mockStatus.EXPECT().Get(nil, nil, nil).Return(nsxModel.UpgradeStatus{OverallUpgradeStatus: &overallPaused}, nil)

		_, needed, err := getSummaryInfo(newGoMockProviderClient())
		require.NoError(t, err)
		assert.True(t, needed)
	})

	t.Run("VCF9 status lookup error is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()
		_, _, _, mockStatus, _, _, restoreRun := setupUpgradeRunMocks(ctrl)
		defer restoreRun()

		target := "9.0.0"
		mockSummary.EXPECT().Get().Return(nsxModel.UpgradeSummary{UpgradeCoordinatorUpdated: &ucUpdated, TargetVersion: &target}, nil)
		mockStatus.EXPECT().Get(nil, nil, nil).Return(nsxModel.UpgradeStatus{}, errors.New("status unavailable"))

		_, _, err := getSummaryInfo(newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("summary Get error is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		mockSummary.EXPECT().Get().Return(nsxModel.UpgradeSummary{}, errors.New("summary unavailable"))

		_, _, err := getSummaryInfo(newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtUpgradePrepareReadWithPrecheck(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	res := resourceNsxtUpgradePrepare()
	anyFailuresList := func(m *preupgchecks.MockFailuresClient) *gomock.Call {
		return m.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	}

	t.Run("runs prechecks and records the failures", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()
		_, _, _, mockStatus, mockUpgrade, _, restoreRun := setupUpgradeRunMocks(ctrl)
		defer restoreRun()
		// Registered last so its FailuresClient mock is the one wired in.
		_, _, mockFailures, restoreAck := setupPrecheckAcknowledgeMocks(ctrl)
		defer restoreAck()

		mockSummary.EXPECT().Get().Return(upgradeSummaryUCUpdated(), nil)
		// An already-acknowledged warning: remembered before the precheck run, and needs no re-ack.
		anyFailuresList(mockFailures).Return(nsxModel.UpgradeCheckFailureListResult{
			Results: []nsxModel.UpgradeCheckFailure{precheckWarningItem(precheckID, true)},
		}, nil).AnyTimes()
		mockUpgrade.EXPECT().Executepreupgradechecks(nil, nil, nil, nil, nil, nil).Return(nil)
		completed := nsxModel.UpgradeChecksExecutionStatus_STATUS_COMPLETED
		mockStatus.EXPECT().Get(gomock.Any(), nil, nil).Return(nsxModel.UpgradeStatus{
			ComponentStatus: []nsxModel.ComponentUpgradeStatus{
				{PreUpgradeStatus: &nsxModel.UpgradeChecksExecutionStatus{Status: &completed}},
			},
		}, nil).Times(len(precheckComponentTypes))

		data := minimalUpgradePrepareData()
		data["precheck_timeout"] = 5
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId("some-id")

		require.NoError(t, resourceNsxtUpgradePrepareRead(d, newGoMockProviderClient()))
		failed := d.Get("failed_prechecks").([]interface{})
		require.Len(t, failed, 1)
		assert.Equal(t, precheckID, failed[0].(map[string]interface{})["id"])
	})

	t.Run("previous-ack lookup error is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, mockFailures, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		mockSummary.EXPECT().Get().Return(upgradeSummaryUCUpdated(), nil)
		anyFailuresList(mockFailures).Return(nsxModel.UpgradeCheckFailureListResult{}, errors.New("list failed"))

		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradePrepareData())
		d.SetId("some-id")
		require.Error(t, resourceNsxtUpgradePrepareRead(d, newGoMockProviderClient()))
	})

	t.Run("precheck execution error is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, mockFailures, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()
		_, _, _, _, mockUpgrade, _, restoreRun := setupUpgradeRunMocks(ctrl)
		defer restoreRun()

		mockSummary.EXPECT().Get().Return(upgradeSummaryUCUpdated(), nil)
		anyFailuresList(mockFailures).Return(nsxModel.UpgradeCheckFailureListResult{}, nil)
		mockUpgrade.EXPECT().Executepreupgradechecks(nil, nil, nil, nil, nil, nil).Return(errors.New("precheck rejected"))

		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradePrepareData())
		d.SetId("some-id")
		require.Error(t, resourceNsxtUpgradePrepareRead(d, newGoMockProviderClient()))
	})

	t.Run("summary error is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, _, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		mockSummary.EXPECT().Get().Return(nsxModel.UpgradeSummary{}, errors.New("summary unavailable"))

		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradePrepareData())
		d.SetId("some-id")
		require.Error(t, resourceNsxtUpgradePrepareRead(d, newGoMockProviderClient()))
	})

	t.Run("failure listing error is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSummary, _, _, _, _, mockFailures, restore := setupUpgradePrepareMocks(ctrl)
		defer restore()

		mockSummary.EXPECT().Get().Return(upgradeSummaryNotStarted(), nil)
		anyFailuresList(mockFailures).Return(nsxModel.UpgradeCheckFailureListResult{}, errors.New("list failed"))

		d := schema.TestResourceDataRaw(t, res.Schema, minimalUpgradePrepareData())
		d.SetId("some-id")
		require.Error(t, resourceNsxtUpgradePrepareRead(d, newGoMockProviderClient()))
	})
}
