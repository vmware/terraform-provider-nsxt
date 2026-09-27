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
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"
)

// setupLBPersistenceProfileMock and lbPersistenceProfileStructValue(t, obj, bindingType) are
// defined in utgomock_resource_nsxt_policy_lb_cookie_persistence_profile_test.go

// genericLBPersistenceProfileValue builds a *data.StructValue for a base LBPersistenceProfile,
// using only id/display_name/resourceType - lighter weight than a full typed profile when the
// data source test only cares about those fields.
func genericLBPersistenceProfileValue(t *testing.T, id, displayName, resourceType string) *data.StructValue {
	t.Helper()
	path := "/infra/lb-persistence-profiles/" + id
	return lbPersistenceProfileStructValue(t, nsxModel.LBPersistenceProfile{
		Id:           &id,
		DisplayName:  &displayName,
		Path:         &path,
		ResourceType: resourceType,
	}, nsxModel.LBPersistenceProfileBindingType())
}

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
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.LBPersistenceProfileListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "test-profile",
		})

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by name - not found returns error", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.LBPersistenceProfileListResult{}, nil)

		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestMockDataSourceNsxtPolicyLbPersistenceProfileSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBPersistenceProfileMock(t, ctrl)
	defer restore()

	t.Run("Read by ID succeeds", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbCookieID).Return(genericLBPersistenceProfileValue(t, lbCookieID, lbCookieDisplayName, nsxModel.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE), nil)

		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbCookieID,
		})

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbCookieDisplayName, d.Get("display_name"))
		assert.Equal(t, "COOKIE", d.Get("type"))
		assert.Equal(t, lbCookieID, d.Id())
	})

	t.Run("Read by ID conversion error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbCookieID).Return(&data.StructValue{}, nil)

		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbCookieID,
		})

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by name - single perfect match succeeds", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBPersistenceProfileListResult{Results: []*data.StructValue{genericLBPersistenceProfileValue(t, lbCookieID, lbCookieDisplayName, nsxModel.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE)}}, nil)

		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": lbCookieDisplayName,
		})

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbCookieID, d.Id())
	})

	t.Run("Read by name - multiple perfect matches is an error", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBPersistenceProfileListResult{Results: []*data.StructValue{
				genericLBPersistenceProfileValue(t, "dup-1", "dup", nsxModel.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE),
				genericLBPersistenceProfileValue(t, "dup-2", "dup", nsxModel.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE),
			}}, nil)

		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "dup",
		})

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("Read by name - single prefix match succeeds when type also matches", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBPersistenceProfileListResult{Results: []*data.StructValue{genericLBPersistenceProfileValue(t, lbCookieID, lbCookieDisplayName, nsxModel.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE)}}, nil)

		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "Test LB Cookie",
			"type":         "COOKIE",
		})

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbCookieID, d.Id())
	})

	t.Run("Read by name - prefix match filtered out by type is not found", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBPersistenceProfileListResult{Results: []*data.StructValue{genericLBPersistenceProfileValue(t, lbCookieID, lbCookieDisplayName, nsxModel.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE)}}, nil)

		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "Test LB Cookie",
			"type":         "SOURCE_IP",
		})

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Read by name - multiple prefix matches is an error", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBPersistenceProfileListResult{Results: []*data.StructValue{
				genericLBPersistenceProfileValue(t, "pfx-1", "prefix-a", nsxModel.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE),
				genericLBPersistenceProfileValue(t, "pfx-2", "prefix-b", nsxModel.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE),
			}}, nil)

		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "prefix",
		})

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("Read by type only (no display_name) matches by type", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBPersistenceProfileListResult{Results: []*data.StructValue{genericLBPersistenceProfileValue(t, lbCookieID, lbCookieDisplayName, nsxModel.LBPersistenceProfile_RESOURCE_TYPE_LBCOOKIEPERSISTENCEPROFILE)}}, nil)

		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"type": "COOKIE",
		})

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbCookieID, d.Id())
	})

	t.Run("Read with none of id, name or type set is an error", func(t *testing.T) {
		ds := dataSourceNsxtPolicyLbPersistenceProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
		d.Set("type", "")

		err := dataSourceNsxtPolicyLbPersistenceProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining")
	})
}
