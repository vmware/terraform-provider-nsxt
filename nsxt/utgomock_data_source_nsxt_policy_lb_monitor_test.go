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

// setupLBMonitorProfileMock is defined in utgomock_resource_nsxt_policy_lb_http_monitor_profile_test.go

// lbMonitorProfileStructValue builds a generic *data.StructValue for a base LBMonitorProfile.
func lbMonitorProfileStructValue(t *testing.T, id, displayName, resourceType string) *data.StructValue {
	t.Helper()
	converter := bindings.NewTypeConverter()
	path := "/infra/lb-monitor-profiles/" + id
	profile := nsxModel.LBMonitorProfile{
		Id:           &id,
		DisplayName:  &displayName,
		Path:         &path,
		ResourceType: resourceType,
	}
	val, errs := converter.ConvertToVapi(profile, nsxModel.LBMonitorProfileBindingType())
	require.Empty(t, errs)
	return val.(*data.StructValue)
}

func TestMockDataSourceNsxtPolicyLBMonitorErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Read by ID - API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpMonitorID).Return(nil, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyLBMonitor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbHttpMonitorID,
		})

		err := dataSourceNsxtPolicyLBMonitorRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by ID - not found returns error", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpMonitorID).Return(nil, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyLBMonitor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbHttpMonitorID,
		})

		err := dataSourceNsxtPolicyLBMonitorRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by name - List API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.LBMonitorProfileListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyLBMonitor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "test-monitor",
		})

		err := dataSourceNsxtPolicyLBMonitorRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by name - not found returns error", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.LBMonitorProfileListResult{}, nil)

		ds := dataSourceNsxtPolicyLBMonitor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		err := dataSourceNsxtPolicyLBMonitorRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestMockDataSourceNsxtPolicyLBMonitorSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Read by ID with type ANY succeeds", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpMonitorID).Return(lbMonitorProfileStructValue(t, lbHttpMonitorID, lbHttpMonitorDisplayName, nsxModel.LBMonitorProfile_RESOURCE_TYPE_LBHTTPMONITORPROFILE), nil)

		ds := dataSourceNsxtPolicyLBMonitor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":   lbHttpMonitorID,
			"type": "ANY",
		})

		err := dataSourceNsxtPolicyLBMonitorRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbHttpMonitorDisplayName, d.Get("display_name"))
		assert.Equal(t, "HTTP", d.Get("type"))
		assert.Equal(t, lbHttpMonitorID, d.Id())
	})

	t.Run("Read by ID with a matching explicit type succeeds", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpMonitorID).Return(lbMonitorProfileStructValue(t, lbHttpMonitorID, lbHttpMonitorDisplayName, nsxModel.LBMonitorProfile_RESOURCE_TYPE_LBHTTPMONITORPROFILE), nil)

		ds := dataSourceNsxtPolicyLBMonitor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":   lbHttpMonitorID,
			"type": "HTTP",
		})

		err := dataSourceNsxtPolicyLBMonitorRead(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Read by ID with a mismatched type is treated as not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpMonitorID).Return(lbMonitorProfileStructValue(t, lbHttpMonitorID, lbHttpMonitorDisplayName, nsxModel.LBMonitorProfile_RESOURCE_TYPE_LBHTTPMONITORPROFILE), nil)

		ds := dataSourceNsxtPolicyLBMonitor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":   lbHttpMonitorID,
			"type": "TCP",
		})

		err := dataSourceNsxtPolicyLBMonitorRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("Read by ID with an unrecognized resource type fails conversion", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpMonitorID).Return(lbMonitorProfileStructValue(t, lbHttpMonitorID, lbHttpMonitorDisplayName, "LBUnknownMonitorProfile"), nil)

		ds := dataSourceNsxtPolicyLBMonitor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbHttpMonitorID,
		})

		err := dataSourceNsxtPolicyLBMonitorRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while converting")
	})

	t.Run("Read by name - single perfect match succeeds", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBMonitorProfileListResult{Results: []*data.StructValue{lbMonitorProfileStructValue(t, lbHttpMonitorID, lbHttpMonitorDisplayName, nsxModel.LBMonitorProfile_RESOURCE_TYPE_LBHTTPMONITORPROFILE)}}, nil)

		ds := dataSourceNsxtPolicyLBMonitor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": lbHttpMonitorDisplayName,
		})

		err := dataSourceNsxtPolicyLBMonitorRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbHttpMonitorID, d.Id())
	})

	t.Run("Read by name - multiple perfect matches is an error", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBMonitorProfileListResult{Results: []*data.StructValue{
				lbMonitorProfileStructValue(t, "dup-1", "dup", nsxModel.LBMonitorProfile_RESOURCE_TYPE_LBHTTPMONITORPROFILE),
				lbMonitorProfileStructValue(t, "dup-2", "dup", nsxModel.LBMonitorProfile_RESOURCE_TYPE_LBHTTPMONITORPROFILE),
			}}, nil)

		ds := dataSourceNsxtPolicyLBMonitor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "dup",
		})

		err := dataSourceNsxtPolicyLBMonitorRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("Read by name - single prefix match succeeds", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBMonitorProfileListResult{Results: []*data.StructValue{lbMonitorProfileStructValue(t, lbHttpMonitorID, lbHttpMonitorDisplayName, nsxModel.LBMonitorProfile_RESOURCE_TYPE_LBHTTPMONITORPROFILE)}}, nil)

		ds := dataSourceNsxtPolicyLBMonitor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "Test LB HTTP",
		})

		err := dataSourceNsxtPolicyLBMonitorRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbHttpMonitorID, d.Id())
	})

	t.Run("Read by name - multiple prefix matches is an error", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBMonitorProfileListResult{Results: []*data.StructValue{
				lbMonitorProfileStructValue(t, "pfx-1", "prefix-a", nsxModel.LBMonitorProfile_RESOURCE_TYPE_LBHTTPMONITORPROFILE),
				lbMonitorProfileStructValue(t, "pfx-2", "prefix-b", nsxModel.LBMonitorProfile_RESOURCE_TYPE_LBHTTPMONITORPROFILE),
			}}, nil)

		ds := dataSourceNsxtPolicyLBMonitor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "prefix",
		})

		err := dataSourceNsxtPolicyLBMonitorRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("Read by type only (no display_name) matches by type", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), nil, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nsxModel.LBMonitorProfileListResult{Results: []*data.StructValue{lbMonitorProfileStructValue(t, lbHttpMonitorID, lbHttpMonitorDisplayName, nsxModel.LBMonitorProfile_RESOURCE_TYPE_LBHTTPMONITORPROFILE)}}, nil)

		ds := dataSourceNsxtPolicyLBMonitor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"type": "HTTP",
		})

		err := dataSourceNsxtPolicyLBMonitorRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbHttpMonitorID, d.Id())
	})

	t.Run("Read with none of id, name or type set is an error", func(t *testing.T) {
		ds := dataSourceNsxtPolicyLBMonitor()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})
		d.Set("type", "")

		err := dataSourceNsxtPolicyLBMonitorRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining")
	})
}
