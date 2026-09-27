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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"
)

// setupLBAppProfileMock is defined in utgomock_resource_nsxt_policy_lb_fast_tcp_application_profile_test.go

// lbAppProfileStructValue builds a generic *data.StructValue for a base LBAppProfile - lighter
// weight than lbFastTcpStructValue's full LBFastTcpProfile when the test only cares about id,
// display_name and resourceType (e.g. the search-by-name/type and unknown-type branches).
func lbAppProfileStructValue(t *testing.T, id, displayName, resourceType string) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	path := "/infra/lb-app-profiles/" + id
	profile := nsxModel.LBAppProfile{
		Id:           &id,
		DisplayName:  &displayName,
		Path:         &path,
		ResourceType: resourceType,
	}
	val, errs := converter.ConvertToVapi(profile, nsxModel.LBAppProfileBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func TestMockDataSourceNsxtPolicyLBAppProfileErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Read by ID - API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbFastTcpID).Return(nil, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyLBAppProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbFastTcpID,
		})

		err := dataSourceNsxtPolicyLBAppProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by ID - not found returns error", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbFastTcpID).Return(nil, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyLBAppProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbFastTcpID,
		})

		err := dataSourceNsxtPolicyLBAppProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by name - List API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.LBAppProfileListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyLBAppProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "test-profile",
		})

		err := dataSourceNsxtPolicyLBAppProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by name - not found returns error", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.LBAppProfileListResult{}, nil)

		ds := dataSourceNsxtPolicyLBAppProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		err := dataSourceNsxtPolicyLBAppProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestMockDataSourceNsxtPolicyLBAppProfileSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBAppProfileMock(t, ctrl)
	defer restore()

	t.Run("Read by ID with type ANY succeeds", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbFastTcpID).Return(lbAppProfileStructValue(t, lbFastTcpID, lbFastTcpDisplayName, nsxModel.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE), nil)

		ds := dataSourceNsxtPolicyLBAppProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":   lbFastTcpID,
			"type": "ANY",
		})

		err := dataSourceNsxtPolicyLBAppProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbFastTcpDisplayName, d.Get("display_name"))
		assert.Equal(t, "TCP", d.Get("type"))
		assert.Equal(t, lbFastTcpID, d.Id())
	})

	t.Run("Read by ID with a matching explicit type succeeds", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbFastTcpID).Return(lbAppProfileStructValue(t, lbFastTcpID, lbFastTcpDisplayName, nsxModel.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE), nil)

		ds := dataSourceNsxtPolicyLBAppProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":   lbFastTcpID,
			"type": "TCP",
		})

		err := dataSourceNsxtPolicyLBAppProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Read by ID with a mismatched type is treated as not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbFastTcpID).Return(lbAppProfileStructValue(t, lbFastTcpID, lbFastTcpDisplayName, nsxModel.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE), nil)

		ds := dataSourceNsxtPolicyLBAppProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":   lbFastTcpID,
			"type": "HTTP",
		})

		err := dataSourceNsxtPolicyLBAppProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("Read by ID with an unrecognized resource type fails conversion", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbFastTcpID).Return(lbAppProfileStructValue(t, lbFastTcpID, lbFastTcpDisplayName, "LBUnknownProfile"), nil)

		ds := dataSourceNsxtPolicyLBAppProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbFastTcpID,
		})

		err := dataSourceNsxtPolicyLBAppProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while converting")
	})

	t.Run("Read by name - single perfect match succeeds", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBAppProfileListResult{Results: []*data.StructValue{lbAppProfileStructValue(t, lbFastTcpID, lbFastTcpDisplayName, nsxModel.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE)}}, nil)

		ds := dataSourceNsxtPolicyLBAppProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": lbFastTcpDisplayName,
		})

		err := dataSourceNsxtPolicyLBAppProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbFastTcpID, d.Id())
	})

	t.Run("Read by name - multiple perfect matches is an error", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBAppProfileListResult{Results: []*data.StructValue{
				lbAppProfileStructValue(t, "dup-1", "dup", nsxModel.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE),
				lbAppProfileStructValue(t, "dup-2", "dup", nsxModel.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE),
			}}, nil)

		ds := dataSourceNsxtPolicyLBAppProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "dup",
		})

		err := dataSourceNsxtPolicyLBAppProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("Read by name - single prefix match succeeds", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBAppProfileListResult{Results: []*data.StructValue{lbAppProfileStructValue(t, lbFastTcpID, lbFastTcpDisplayName, nsxModel.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE)}}, nil)

		ds := dataSourceNsxtPolicyLBAppProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "Test LB Fast",
		})

		err := dataSourceNsxtPolicyLBAppProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbFastTcpID, d.Id())
	})

	t.Run("Read by name - multiple prefix matches is an error", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBAppProfileListResult{Results: []*data.StructValue{
				lbAppProfileStructValue(t, "pfx-1", "prefix-a", nsxModel.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE),
				lbAppProfileStructValue(t, "pfx-2", "prefix-b", nsxModel.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE),
			}}, nil)

		ds := dataSourceNsxtPolicyLBAppProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "prefix",
		})

		err := dataSourceNsxtPolicyLBAppProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("Read by type only (no display_name) matches by type", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBAppProfileListResult{Results: []*data.StructValue{lbAppProfileStructValue(t, lbFastTcpID, lbFastTcpDisplayName, nsxModel.LBAppProfile_RESOURCE_TYPE_LBFASTTCPPROFILE)}}, nil)

		ds := dataSourceNsxtPolicyLBAppProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"type": "TCP",
		})

		err := dataSourceNsxtPolicyLBAppProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbFastTcpID, d.Id())
	})

	t.Run("Read with none of id, name or type set is an error", func(t *testing.T) {
		ds := dataSourceNsxtPolicyLBAppProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
		d.Set("type", "")

		err := dataSourceNsxtPolicyLBAppProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining")
	})
}
