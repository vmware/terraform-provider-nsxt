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

// setupLBServerSslMock is defined in utgomock_resource_nsxt_policy_lb_server_ssl_profile_test.go

func TestMockDataSourceNsxtPolicyLBServerSslProfileGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBServerSslMock(t, ctrl)
	defer restore()

	t.Run("Read by ID success", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbServerSslID).Return(lbServerSslAPIResponse(), nil)

		ds := dataSourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbServerSslID,
		})

		err := dataSourceNsxtPolicyLBServerSslProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbServerSslID, d.Id())
		assert.Equal(t, lbServerSslDisplayName, d.Get("display_name"))
	})

	t.Run("Read by ID not found returns error", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbServerSslID).Return(nsxModel.LBServerSslProfile{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbServerSslID,
		})

		err := dataSourceNsxtPolicyLBServerSslProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read by ID API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbServerSslID).Return(nsxModel.LBServerSslProfile{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": lbServerSslID,
		})

		err := dataSourceNsxtPolicyLBServerSslProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockDataSourceNsxtPolicyLBServerSslProfileListByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBServerSslMock(t, ctrl)
	defer restore()

	t.Run("Read by name perfect match", func(t *testing.T) {
		listResult := nsxModel.LBServerSslProfileListResult{
			Results: []nsxModel.LBServerSslProfile{lbServerSslAPIResponse()},
		}
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(listResult, nil)

		ds := dataSourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": lbServerSslDisplayName,
		})

		err := dataSourceNsxtPolicyLBServerSslProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbServerSslID, d.Id())
		assert.Equal(t, lbServerSslDisplayName, d.Get("display_name"))
	})

	t.Run("Read by name prefix match", func(t *testing.T) {
		listResult := nsxModel.LBServerSslProfileListResult{
			Results: []nsxModel.LBServerSslProfile{lbServerSslAPIResponse()},
		}
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(listResult, nil)

		ds := dataSourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "Test LB",
		})

		err := dataSourceNsxtPolicyLBServerSslProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbServerSslID, d.Id())
	})

	t.Run("Read by name not found returns error", func(t *testing.T) {
		listResult := nsxModel.LBServerSslProfileListResult{
			Results: []nsxModel.LBServerSslProfile{},
		}
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(listResult, nil)

		ds := dataSourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": "nonexistent",
		})

		err := dataSourceNsxtPolicyLBServerSslProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Read by name multiple perfect matches returns error", func(t *testing.T) {
		dup := lbServerSslAPIResponse()
		listResult := nsxModel.LBServerSslProfileListResult{
			Results: []nsxModel.LBServerSslProfile{lbServerSslAPIResponse(), dup},
		}
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(listResult, nil)

		ds := dataSourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": lbServerSslDisplayName,
		})

		err := dataSourceNsxtPolicyLBServerSslProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("Read by name List API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nsxModel.LBServerSslProfileListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": lbServerSslDisplayName,
		})

		err := dataSourceNsxtPolicyLBServerSslProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockDataSourceNsxtPolicyLBServerSslProfileNoIDOrName(t *testing.T) {
	t.Run("No ID or display_name returns error", func(t *testing.T) {
		ds := dataSourceNsxtPolicyLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{})

		err := dataSourceNsxtPolicyLBServerSslProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining")
	})
}
