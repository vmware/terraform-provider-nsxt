// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mocks for this test, run:
// mockgen -destination=mocks/nsx/upgrade/UpgradeChecksInfoClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/upgrade/UpgradeChecksInfoClient.go UpgradeChecksInfoClient
// mockgen -destination=mocks/nsx/upgrade/PreUpgradeChecksClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/upgrade/PreUpgradeChecksClient.go PreUpgradeChecksClient
// mockgen -destination=mocks/nsx/upgrade/pre_upgrade_checks/FailuresClient.go -package=mocks -source=<sdk>/services/nsxt-mp/nsx/upgrade/pre_upgrade_checks/FailuresClient.go FailuresClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/upgrade"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/upgrade/pre_upgrade_checks"

	upgrademocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/upgrade"
	preupgcheckmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx/upgrade/pre_upgrade_checks"
)

var precheckID = "check-id-1"

func precheckWarningItem(id string, acked bool) nsxModel.UpgradeCheckFailure {
	warningType := nsxModel.UpgradeCheckFailure_TYPE_WARNING
	msg := "test warning"
	return nsxModel.UpgradeCheckFailure{
		Id:    &id,
		Type_: &warningType,
		Acked: &acked,
		Message: &nsxModel.UpgradeCheckFailureMessage{
			Message: &msg,
		},
	}
}

func checksInfoResult(id string) nsxModel.ComponentUpgradeChecksInfoListResult {
	ciID := id
	ci := nsxModel.UpgradeCheckInfo{Id: &ciID}
	comp := nsxModel.ComponentUpgradeChecksInfo{
		PreUpgradeChecksInfo: []nsxModel.UpgradeCheckInfo{ci},
	}
	return nsxModel.ComponentUpgradeChecksInfoListResult{
		Results: []nsxModel.ComponentUpgradeChecksInfo{comp},
	}
}

func minimalPrecheckAcknowledgeData() map[string]interface{} {
	return map[string]interface{}{
		"precheck_ids":   []interface{}{precheckID},
		"target_version": "4.1.0",
	}
}

func setupPrecheckAcknowledgeMocks(ctrl *gomock.Controller) (
	*upgrademocks.MockUpgradeChecksInfoClient,
	*upgrademocks.MockPreUpgradeChecksClient,
	*preupgcheckmocks.MockFailuresClient,
	func(),
) {
	mockChecksInfo := upgrademocks.NewMockUpgradeChecksInfoClient(ctrl)
	mockPreChecks := upgrademocks.NewMockPreUpgradeChecksClient(ctrl)
	mockFailures := preupgcheckmocks.NewMockFailuresClient(ctrl)

	origChecksInfo := cliUpgradeChecksInfoClient
	cliUpgradeChecksInfoClient = func(_ vapiProtocolClient.Connector) upgrade.UpgradeChecksInfoClient {
		return mockChecksInfo
	}

	origPreChecks := cliPreUpgradeChecksClient
	cliPreUpgradeChecksClient = func(_ vapiProtocolClient.Connector) upgrade.PreUpgradeChecksClient {
		return mockPreChecks
	}

	origFailures := cliPreUpgradeChecksFailuresClient
	cliPreUpgradeChecksFailuresClient = func(_ vapiProtocolClient.Connector) pre_upgrade_checks.FailuresClient {
		return mockFailures
	}

	restore := func() {
		cliUpgradeChecksInfoClient = origChecksInfo
		cliPreUpgradeChecksClient = origPreChecks
		cliPreUpgradeChecksFailuresClient = origFailures
	}
	return mockChecksInfo, mockPreChecks, mockFailures, restore
}

func TestMockResourceNsxtUpgradePrecheckAcknowledgeCreate(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Create success with acknowledgment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockChecksInfo, mockPreChecks, mockFailures, restore := setupPrecheckAcknowledgeMocks(ctrl)
		defer restore()

		warningType := nsxModel.UpgradeCheckFailure_TYPE_WARNING
		acked := false
		gomock.InOrder(
			// validatePrecheckIDs
			mockChecksInfo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(checksInfoResult(precheckID), nil),
			// acknowledgePrecheckWarnings -> getPrecheckErrors (WARNING type)
			mockFailures.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), &warningType, gomock.Any(), gomock.Any()).Return(
				nsxModel.UpgradeCheckFailureListResult{
					Results: []nsxModel.UpgradeCheckFailure{precheckWarningItem(precheckID, acked)},
				}, nil,
			),
			// acknowledge the warning
			mockPreChecks.EXPECT().Acknowledge(precheckID).Return(nil),
			// Read -> getPrecheckErrors (WARNING type, same as acknowledgePrecheckWarnings)
			mockFailures.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
				nsxModel.UpgradeCheckFailureListResult{
					Results: []nsxModel.UpgradeCheckFailure{},
				}, nil,
			),
		)

		res := resourceNsxtUpgradePrecheckAcknowledge()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrecheckAcknowledgeData())

		err := resourceNsxtUpgradePrecheckAcknowledgeCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})

	t.Run("Create fails when precheck ID invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, _, restore := setupPrecheckAcknowledgeMocks(ctrl)
		defer restore()

		mockChecksInfo2 := upgrademocks.NewMockUpgradeChecksInfoClient(ctrl)
		origChecksInfo := cliUpgradeChecksInfoClient
		cliUpgradeChecksInfoClient = func(_ vapiProtocolClient.Connector) upgrade.UpgradeChecksInfoClient {
			return mockChecksInfo2
		}
		defer func() { cliUpgradeChecksInfoClient = origChecksInfo }()

		// Return empty results - the precheckID won't be in valid checks
		mockChecksInfo2.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nsxModel.ComponentUpgradeChecksInfoListResult{Results: []nsxModel.ComponentUpgradeChecksInfo{}}, nil,
		)

		data := map[string]interface{}{
			"precheck_ids":   []interface{}{"invalid-check-id"},
			"target_version": "4.1.0",
		}
		res := resourceNsxtUpgradePrecheckAcknowledge()
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtUpgradePrecheckAcknowledgeCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid-check-id")
	})
}

func TestMockResourceNsxtUpgradePrecheckAcknowledgeRead(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, mockFailures, restore := setupPrecheckAcknowledgeMocks(ctrl)
		defer restore()

		acked := true
		mockFailures.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeCheckFailureListResult{
				Results: []nsxModel.UpgradeCheckFailure{precheckWarningItem(precheckID, acked)},
			}, nil,
		)

		res := resourceNsxtUpgradePrecheckAcknowledge()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrecheckAcknowledgeData())
		d.SetId("some-id")

		err := resourceNsxtUpgradePrecheckAcknowledgeRead(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Read fails on API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, mockFailures, restore := setupPrecheckAcknowledgeMocks(ctrl)
		defer restore()

		mockFailures.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
			nsxModel.UpgradeCheckFailureListResult{}, errors.New("API error"),
		)

		res := resourceNsxtUpgradePrecheckAcknowledge()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrecheckAcknowledgeData())
		d.SetId("some-id")

		err := resourceNsxtUpgradePrecheckAcknowledgeRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtUpgradePrecheckAcknowledgeDelete(t *testing.T) {
	util.NsxVersion = "3.0.0"
	defer func() { util.NsxVersion = "" }()
	t.Run("Delete is a no-op", func(t *testing.T) {
		res := resourceNsxtUpgradePrecheckAcknowledge()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPrecheckAcknowledgeData())
		d.SetId("some-id")

		err := resourceNsxtUpgradePrecheckAcknowledgeDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
