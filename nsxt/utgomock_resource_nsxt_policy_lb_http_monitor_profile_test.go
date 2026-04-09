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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	infraMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

var (
	lbHttpMonitorID          = "lb-http-monitor-1"
	lbHttpMonitorDisplayName = "Test LB HTTP Monitor Profile"
)

func minimalLBHttpMonitorData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbHttpMonitorDisplayName,
		"nsx_id":       lbHttpMonitorID,
	}
}

func setupLBMonitorProfileMock(t *testing.T, ctrl *gomock.Controller) (*infraMocks.MockLbMonitorProfilesClient, func()) {
	mockSDK := infraMocks.NewMockLbMonitorProfilesClient(ctrl)
	mockWrapper := &infraapi.LBMonitorProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliLbMonitorProfilesClient
	cliLbMonitorProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.LBMonitorProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliLbMonitorProfilesClient = original }
}

func TestMockResourceNsxtPolicyLBHttpMonitorProfileCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpMonitorID).Return(nil, nil)

		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())

		err := resourceNsxtPolicyLBHttpMonitorProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyLBHttpMonitorProfileRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpMonitorID).Return(nil, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())
		d.SetId(lbHttpMonitorID)

		err := resourceNsxtPolicyLBHttpMonitorProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read API error is propagated", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbHttpMonitorID).Return(nil, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())
		d.SetId(lbHttpMonitorID)

		err := resourceNsxtPolicyLBHttpMonitorProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())

		err := resourceNsxtPolicyLBHttpMonitorProfileRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBHttpMonitorProfileUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())

		err := resourceNsxtPolicyLBHttpMonitorProfileUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBHttpMonitorProfileDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBMonitorProfileMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbHttpMonitorID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())
		d.SetId(lbHttpMonitorID)

		err := resourceNsxtPolicyLBHttpMonitorProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBHttpMonitorProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBHttpMonitorData())

		err := resourceNsxtPolicyLBHttpMonitorProfileDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
