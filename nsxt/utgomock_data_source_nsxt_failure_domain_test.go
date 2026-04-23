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
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
)

func TestMockDataSourceNsxtFailureDomainRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupFailureDomainMock(t, ctrl)
	defer restore()

	t.Run("by id success", func(t *testing.T) {
		mockSDK.EXPECT().Get(fdID).Return(failureDomainAPIResponse(), nil)

		ds := dataSourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": fdID,
		})

		err := dataSourceNsxtFailureDomainRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, fdID, d.Id())
		assert.Equal(t, fdDisplayName, d.Get("display_name"))
	})

	t.Run("by id not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(fdID).Return(nsxModel.FailureDomain{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": fdID,
		})

		err := dataSourceNsxtFailureDomainRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by display_name via list", func(t *testing.T) {
		mockSDK.EXPECT().List().Return(nsxModel.FailureDomainListResult{
			Results: []nsxModel.FailureDomain{failureDomainAPIResponse()},
		}, nil)

		ds := dataSourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": fdDisplayName,
		})

		err := dataSourceNsxtFailureDomainRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, fdID, d.Id())
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List().Return(nsxModel.FailureDomainListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": fdDisplayName,
		})

		err := dataSourceNsxtFailureDomainRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
