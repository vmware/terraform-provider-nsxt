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
	"testing"

	"github.com/golang/mock/gomock"
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
