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
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	infraapi "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	infraMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
)

var (
	lbVsID          = "lb-vs-1"
	lbVsDisplayName = "Test LB Virtual Server"
	lbVsDescription = "Test lb virtual server"
	lbVsRevision    = int64(1)
	lbVsIPAddr      = "10.0.0.1"
)

func lbVsAPIResponse() nsxModel.LBVirtualServer {
	return nsxModel.LBVirtualServer{
		Id:          &lbVsID,
		DisplayName: &lbVsDisplayName,
		Description: &lbVsDescription,
		Revision:    &lbVsRevision,
		IpAddress:   &lbVsIPAddr,
	}
}

func minimalLBVsData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": lbVsDisplayName,
		"description":  lbVsDescription,
		"nsx_id":       lbVsID,
		"ip_address":   lbVsIPAddr,
		"ports":        []interface{}{"80"},
	}
}

func setupLBVsMock(t *testing.T, ctrl *gomock.Controller) (*infraMocks.MockLbVirtualServersClient, func()) {
	mockSDK := infraMocks.NewMockLbVirtualServersClient(ctrl)
	mockWrapper := &infraapi.LBVirtualServerClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	original := cliLbVirtualServersClient
	cliLbVirtualServersClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *infraapi.LBVirtualServerClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliLbVirtualServersClient = original }
}

func TestMockResourceNsxtPolicyLBVirtualServerCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBVsMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(lbVsID).Return(nsxModel.LBVirtualServer{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(lbVsID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(lbVsID).Return(lbVsAPIResponse(), nil),
		)

		res := resourceNsxtPolicyLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBVsData())

		err := resourceNsxtPolicyLBVirtualServerCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbVsID, d.Id())
	})

	t.Run("Create fails when already exists", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbVsID).Return(lbVsAPIResponse(), nil)

		res := resourceNsxtPolicyLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBVsData())

		err := resourceNsxtPolicyLBVirtualServerCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBVirtualServerRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBVsMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbVsID).Return(lbVsAPIResponse(), nil)

		res := resourceNsxtPolicyLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBVsData())
		d.SetId(lbVsID)

		err := resourceNsxtPolicyLBVirtualServerRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, lbVsDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(lbVsID).Return(nsxModel.LBVirtualServer{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBVsData())
		d.SetId(lbVsID)

		err := resourceNsxtPolicyLBVirtualServerRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBVsData())

		err := resourceNsxtPolicyLBVirtualServerRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBVirtualServerUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBVsMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		// Update: Get (for rule check when no rule change), Patch, then Get (read after update)
		gomock.InOrder(
			mockSDK.EXPECT().Get(lbVsID).Return(lbVsAPIResponse(), nil),
			mockSDK.EXPECT().Patch(lbVsID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(lbVsID).Return(lbVsAPIResponse(), nil),
		)

		res := resourceNsxtPolicyLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBVsData())
		d.SetId(lbVsID)

		err := resourceNsxtPolicyLBVirtualServerUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBVsData())

		err := resourceNsxtPolicyLBVirtualServerUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyLBVirtualServerDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupLBVsMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		mockSDK.EXPECT().Delete(lbVsID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBVsData())
		d.SetId(lbVsID)

		err := resourceNsxtPolicyLBVirtualServerDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBVirtualServer()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBVsData())

		err := resourceNsxtPolicyLBVirtualServerDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
