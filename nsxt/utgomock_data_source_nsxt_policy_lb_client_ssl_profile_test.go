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

// setupLBClientSslMock is defined in utgomock_resource_nsxt_policy_lb_client_ssl_profile_test.go

func TestMockDataSourceNsxtPolicyLBClientSslProfileGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBClientSslMock(t, ctrl)
	defer restore()

	t.Run("Read by ID success", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbClientSslID).Return(lbClientSslAPIResponse(), nil)

		ds := dataSourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbClientSslID,
		})

		err := dataSourceNsxtPolicyLBClientSslProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbClientSslID, d.Id())
		assert.Equal(t, lbClientSslDisplayName, d.Get("display_name"))
	})

	t.Run("Read by ID not found returns error", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbClientSslID).Return(nsxModel.LBClientSslProfile{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbClientSslID,
		})

		err := dataSourceNsxtPolicyLBClientSslProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by ID API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbClientSslID).Return(nsxModel.LBClientSslProfile{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbClientSslID,
		})

		err := dataSourceNsxtPolicyLBClientSslProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockDataSourceNsxtPolicyLBClientSslProfileListByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBClientSslMock(t, ctrl)
	defer restore()

	t.Run("Read by name perfect match", func(t *testing.T) {
		listResult := nsxModel.LBClientSslProfileListResult{
			Results: []nsxModel.LBClientSslProfile{lbClientSslAPIResponse()},
		}
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(listResult, nil)

		ds := dataSourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": lbClientSslDisplayName,
		})

		err := dataSourceNsxtPolicyLBClientSslProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbClientSslID, d.Id())
		assert.Equal(t, lbClientSslDisplayName, d.Get("display_name"))
	})

	t.Run("Read by name prefix match", func(t *testing.T) {
		listResult := nsxModel.LBClientSslProfileListResult{
			Results: []nsxModel.LBClientSslProfile{lbClientSslAPIResponse()},
		}
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(listResult, nil)

		ds := dataSourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "Test LB",
		})

		err := dataSourceNsxtPolicyLBClientSslProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbClientSslID, d.Id())
	})

	t.Run("Read by name not found returns error", func(t *testing.T) {
		listResult := nsxModel.LBClientSslProfileListResult{
			Results: []nsxModel.LBClientSslProfile{},
		}
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(listResult, nil)

		ds := dataSourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		err := dataSourceNsxtPolicyLBClientSslProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Read by name multiple perfect matches returns error", func(t *testing.T) {
		dup := lbClientSslAPIResponse()
		listResult := nsxModel.LBClientSslProfileListResult{
			Results: []nsxModel.LBClientSslProfile{lbClientSslAPIResponse(), dup},
		}
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(listResult, nil)

		ds := dataSourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": lbClientSslDisplayName,
		})

		err := dataSourceNsxtPolicyLBClientSslProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("Read by name List API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.LBClientSslProfileListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": lbClientSslDisplayName,
		})

		err := dataSourceNsxtPolicyLBClientSslProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockDataSourceNsxtPolicyLBClientSslProfileNoIDOrName(t *testing.T) {
	t.Run("No ID or display_name returns error", func(t *testing.T) {
		ds := dataSourceNsxtPolicyLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyLBClientSslProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining")
	})
}
