//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/realized_state/RealizedEntitiesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/realized_state/RealizedEntitiesClient.go RealizedEntitiesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	realizedstate "github.com/vmware/terraform-provider-nsxt/api/infra/realized_state"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	realizedmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/realized_state"
)

func setupRealizationInfoMock(t *testing.T, ctrl *gomock.Controller) (*realizedmocks.MockRealizedEntitiesClient, func()) {
	t.Helper()
	mockSDK := realizedmocks.NewMockRealizedEntitiesClient(ctrl)
	wrapper := &realizedstate.RealizedEntityClientContext{Client: mockSDK, ClientType: utl.Local}

	original := cliRealizedEntitiesClient
	cliRealizedEntitiesClient = func(sessionContext utl.SessionContext, connector client.Connector) *realizedstate.RealizedEntityClientContext {
		return wrapper
	}
	return mockSDK, func() { cliRealizedEntitiesClient = original }
}

func realizationInfoTestData(overrides map[string]interface{}) map[string]interface{} {
	data := map[string]interface{}{
		"path":    "/infra/tier-1s/tier1-1",
		"timeout": 5,
		"delay":   0,
	}
	for k, v := range overrides {
		data[k] = v
	}
	return data
}

func TestMockDataSourceNsxtPolicyRealizationInfoReadGuard(t *testing.T) {
	t.Run("site_path required on global manager", func(t *testing.T) {
		ds := dataSourceNsxtPolicyRealizationInfo()
		d := schema.TestResourceDataRaw(t, ds.Schema, realizationInfoTestData(nil))

		err := dataSourceNsxtPolicyRealizationInfoRead(d, newGoMockGlobalProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "site_path")
	})

	t.Run("site_path not allowed on local manager", func(t *testing.T) {
		ds := dataSourceNsxtPolicyRealizationInfo()
		d := schema.TestResourceDataRaw(t, ds.Schema, realizationInfoTestData(map[string]interface{}{
			"site_path": "/infra/sites/default",
		}))

		err := dataSourceNsxtPolicyRealizationInfoRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Global Manager")
	})
}

func TestMockDataSourceNsxtPolicyRealizationInfoRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupRealizationInfoMock(t, ctrl)
	defer restore()

	t.Run("success on first entry when entity_type not set", func(t *testing.T) {
		state := "REALIZED"
		entityType := "RealizedLogicalRouter"
		realizedID := "realized-1"
		mockSDK.EXPECT().List("/infra/tier-1s/tier1-1", (*string)(nil)).Return(model.GenericPolicyRealizedResourceListResult{
			Results: []model.GenericPolicyRealizedResource{
				{
					State:                         &state,
					EntityType:                    &entityType,
					RealizationSpecificIdentifier: &realizedID,
				},
			},
		}, nil)

		ds := dataSourceNsxtPolicyRealizationInfo()
		d := schema.TestResourceDataRaw(t, ds.Schema, realizationInfoTestData(nil))

		err := dataSourceNsxtPolicyRealizationInfoRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "REALIZED", d.Get("state"))
		assert.Equal(t, entityType, d.Get("entity_type"))
		assert.Equal(t, realizedID, d.Get("realized_id"))
	})

	t.Run("success matching specific entity_type", func(t *testing.T) {
		state := "REALIZED"
		otherType := "OtherType"
		wantedType := "WantedType"
		realizedID := "realized-2"
		mockSDK.EXPECT().List("/infra/tier-1s/tier1-1", (*string)(nil)).Return(model.GenericPolicyRealizedResourceListResult{
			Results: []model.GenericPolicyRealizedResource{
				{State: &state, EntityType: &otherType},
				{State: &state, EntityType: &wantedType, RealizationSpecificIdentifier: &realizedID},
			},
		}, nil)

		ds := dataSourceNsxtPolicyRealizationInfo()
		d := schema.TestResourceDataRaw(t, ds.Schema, realizationInfoTestData(map[string]interface{}{
			"entity_type": wantedType,
		}))

		err := dataSourceNsxtPolicyRealizationInfoRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, realizedID, d.Get("realized_id"))
	})

	t.Run("API error surfaces as failure", func(t *testing.T) {
		mockSDK.EXPECT().List("/infra/tier-1s/tier1-1", (*string)(nil)).Return(
			model.GenericPolicyRealizedResourceListResult{}, errors.New("list failed"),
		)

		ds := dataSourceNsxtPolicyRealizationInfo()
		d := schema.TestResourceDataRaw(t, ds.Schema, realizationInfoTestData(nil))

		err := dataSourceNsxtPolicyRealizationInfoRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to get realization information")
	})
}
