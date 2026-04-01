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
	vapiData_ "github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

var (
	ipbqID          = "ipbq-1"
	ipbqDisplayName = "Test IP Block Quota"
	ipbqDescription = "test ip block quota"
	ipbqRevision    = int64(1)
	ipbqPath        = "/infra/limits/ipbq-1"
)

func ipbqAPIResponse() nsxModel.Limit {
	visibility := nsxModel.IpBlockQuota_IP_BLOCK_VISIBILITY_PRIVATE
	addrType := nsxModel.IpBlockQuota_IP_BLOCK_ADDRESS_TYPE_IPV4
	quota := nsxModel.IpBlockQuota{
		ResourceType:       nsxModel.Quota_RESOURCE_TYPE_IPBLOCKQUOTA,
		IpBlockVisibility:  &visibility,
		IpBlockAddressType: &addrType,
	}
	dataVal, _ := quota.GetDataValue__()
	var quotaStructVal *vapiData_.StructValue
	if sv, ok := dataVal.(*vapiData_.StructValue); ok {
		quotaStructVal = sv
	}
	return nsxModel.Limit{
		Id:          &ipbqID,
		DisplayName: &ipbqDisplayName,
		Description: &ipbqDescription,
		Revision:    &ipbqRevision,
		Path:        &ipbqPath,
		Quota:       quotaStructVal,
	}
}

func minimalIpbqData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": ipbqDisplayName,
		"description":  ipbqDescription,
		"nsx_id":       ipbqID,
		"quota": []interface{}{
			map[string]interface{}{
				"ip_block_visibility":   "PRIVATE",
				"ip_block_address_type": "IPV4",
				"other_cidrs": []interface{}{
					map[string]interface{}{
						"mask":        "",
						"total_count": -1,
					},
				},
				"single_ip_cidrs": -1,
			},
		},
	}
}

func setupIpbqMock(t *testing.T, ctrl *gomock.Controller) (*inframocks.MockLimitsClient, func()) {
	mockSDK := inframocks.NewMockLimitsClient(ctrl)
	mockWrapper := &infraapi.LimitClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliLimitsClient
	cliLimitsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.LimitClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliLimitsClient = original }
}

func TestMockResourceNsxtPolicyIpBlockQuotaCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIpbqMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(ipbqID).Return(nsxModel.Limit{}, notFoundErr),
			mockSDK.EXPECT().Patch(ipbqID, gomock.Any()).Return(nsxModel.Limit{}, nil),
			mockSDK.EXPECT().Get(ipbqID).Return(ipbqAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIpBlockQuota()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpbqData())

		err := resourceNsxtPolicyIpBlockQuotaCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipbqID, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipbqID).Return(ipbqAPIResponse(), nil)

		res := resourceNsxtPolicyIpBlockQuota()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpbqData())

		err := resourceNsxtPolicyIpBlockQuotaCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIpBlockQuotaRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIpbqMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipbqID).Return(ipbqAPIResponse(), nil)

		res := resourceNsxtPolicyIpBlockQuota()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpbqData())
		d.SetId(ipbqID)

		err := resourceNsxtPolicyIpBlockQuotaRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, ipbqDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(ipbqID).Return(nsxModel.Limit{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIpBlockQuota()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpbqData())
		d.SetId(ipbqID)

		err := resourceNsxtPolicyIpBlockQuotaRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIpBlockQuota()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpbqData())

		err := resourceNsxtPolicyIpBlockQuotaRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIpBlockQuotaUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIpbqMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(ipbqID, gomock.Any()).Return(nsxModel.Limit{}, nil),
			mockSDK.EXPECT().Get(ipbqID).Return(ipbqAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIpBlockQuota()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpbqData())
		d.SetId(ipbqID)

		err := resourceNsxtPolicyIpBlockQuotaUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIpBlockQuota()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpbqData())

		err := resourceNsxtPolicyIpBlockQuotaUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIpBlockQuotaDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIpbqMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(ipbqID).Return(nil)

		res := resourceNsxtPolicyIpBlockQuota()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpbqData())
		d.SetId(ipbqID)

		err := resourceNsxtPolicyIpBlockQuotaDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIpBlockQuota()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIpbqData())

		err := resourceNsxtPolicyIpBlockQuotaDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
