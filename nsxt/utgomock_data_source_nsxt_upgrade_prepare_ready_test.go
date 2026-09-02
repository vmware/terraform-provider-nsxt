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
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"go.uber.org/mock/gomock"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

// setupPrecheckAcknowledgeMocks, precheckWarningItem and checksInfoResult are
// defined in utgomock_resource_nsxt_upgrade_precheck_acknowledge_test.go and
// reused here since dataSourceNsxtUpgradePrepareReadyRead calls the same
// getPrecheckErrors/getPrechecksText helpers via the same client vars
// (cliPreUpgradeChecksFailuresClient, cliUpgradeChecksInfoClient).

func precheckFailureItem(id string) nsxModel.UpgradeCheckFailure {
	failureType := nsxModel.UpgradeCheckFailure_TYPE_FAILURE
	msg := "test failure"
	return nsxModel.UpgradeCheckFailure{
		Id:    &id,
		Type_: &failureType,
		Message: &nsxModel.UpgradeCheckFailureMessage{
			Message: &msg,
		},
	}
}

// checksInfoResultWithDetails is like checksInfoResult but also populates
// ComponentType and Description, which getPrechecksText dereferences
// unconditionally for any matching precheck ID.
func checksInfoResultWithDetails(id string) nsxModel.ComponentUpgradeChecksInfoListResult {
	ciID := id
	desc := "test description for " + id
	componentType := edgeUpgradeGroup
	ci := nsxModel.UpgradeCheckInfo{Id: &ciID, Description: &desc}
	comp := nsxModel.ComponentUpgradeChecksInfo{
		ComponentType:        &componentType,
		PreUpgradeChecksInfo: []nsxModel.UpgradeCheckInfo{ci},
	}
	return nsxModel.ComponentUpgradeChecksInfoListResult{
		Results: []nsxModel.ComponentUpgradeChecksInfo{comp},
	}
}

func upgradePrepareReadyData(prepareID string) map[string]interface{} {
	return map[string]interface{}{
		"upgrade_prepare_id": prepareID,
	}
}

func TestUnitNsxt_DataSourceNsxtUpgradePrepareReadyRead(t *testing.T) {
	validPrepareID := util.GetVerifiableID("prepare-1", "nsxt_upgrade_prepare")

	t.Run("invalid upgrade_prepare_id is rejected", func(t *testing.T) {
		ds := dataSourceNsxtUpgradePrepareReady()
		d := schema.TestResourceDataRaw(t, ds.Schema, upgradePrepareReadyData("not-a-valid-id"))

		err := dataSourceNsxtUpgradePrepareReadyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "value for upgrade_prepare_id is invalid")
	})

	t.Run("no precheck errors succeeds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, mockFailures, restore := setupPrecheckAcknowledgeMocks(ctrl)
		defer restore()

		mockFailures.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any()).
			Return(nsxModel.UpgradeCheckFailureListResult{Results: []nsxModel.UpgradeCheckFailure{}}, nil)

		ds := dataSourceNsxtUpgradePrepareReady()
		d := schema.TestResourceDataRaw(t, ds.Schema, upgradePrepareReadyData(validPrepareID))

		err := dataSourceNsxtUpgradePrepareReadyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})

	t.Run("acknowledged warning does not block", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, mockFailures, restore := setupPrecheckAcknowledgeMocks(ctrl)
		defer restore()

		mockFailures.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any()).
			Return(nsxModel.UpgradeCheckFailureListResult{Results: []nsxModel.UpgradeCheckFailure{
				precheckWarningItem(precheckID, true),
			}}, nil)

		ds := dataSourceNsxtUpgradePrepareReady()
		d := schema.TestResourceDataRaw(t, ds.Schema, upgradePrepareReadyData(validPrepareID))

		err := dataSourceNsxtUpgradePrepareReadyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})

	t.Run("unacknowledged warning blocks", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockChecksInfo, _, mockFailures, restore := setupPrecheckAcknowledgeMocks(ctrl)
		defer restore()

		mockFailures.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any()).
			Return(nsxModel.UpgradeCheckFailureListResult{Results: []nsxModel.UpgradeCheckFailure{
				precheckWarningItem(precheckID, false),
			}}, nil)
		mockChecksInfo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(checksInfoResultWithDetails(precheckID), nil)

		ds := dataSourceNsxtUpgradePrepareReady()
		d := schema.TestResourceDataRaw(t, ds.Schema, upgradePrepareReadyData(validPrepareID))

		err := dataSourceNsxtUpgradePrepareReadyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unacknowledged warnings")
	})

	t.Run("failure blocks", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockChecksInfo, _, mockFailures, restore := setupPrecheckAcknowledgeMocks(ctrl)
		defer restore()

		failID := "failure-check-1"
		mockFailures.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any()).
			Return(nsxModel.UpgradeCheckFailureListResult{Results: []nsxModel.UpgradeCheckFailure{
				precheckFailureItem(failID),
			}}, nil)
		mockChecksInfo.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(checksInfoResultWithDetails(failID), nil)

		ds := dataSourceNsxtUpgradePrepareReady()
		d := schema.TestResourceDataRaw(t, ds.Schema, upgradePrepareReadyData(validPrepareID))

		err := dataSourceNsxtUpgradePrepareReadyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failures in prechecks")
	})

	t.Run("getPrecheckErrors API error is wrapped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		_, _, mockFailures, restore := setupPrecheckAcknowledgeMocks(ctrl)
		defer restore()

		mockFailures.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), nil, gomock.Any(), gomock.Any()).
			Return(nsxModel.UpgradeCheckFailureListResult{}, errors.New("list failed"))

		ds := dataSourceNsxtUpgradePrepareReady()
		d := schema.TestResourceDataRaw(t, ds.Schema, upgradePrepareReadyData(validPrepareID))

		err := dataSourceNsxtUpgradePrepareReadyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading precheck failures")
	})
}
