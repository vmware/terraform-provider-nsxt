//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// setupLBPersistenceProfileMock is defined in utgomock_resource_nsxt_policy_lb_cookie_persistence_profile_test.go

func TestMockDataSourceNsxtPolicyLbPersistenceProfileErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Read by ID - API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbCookieID).Return(nil, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbCookieID,
		})

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by ID - not found returns error", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbCookieID).Return(nil, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbCookieID,
		})

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by name - List API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.LBPersistenceProfileListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "test-profile",
		})

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by name - not found returns error", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.LBPersistenceProfileListResult{}, nil)

		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
