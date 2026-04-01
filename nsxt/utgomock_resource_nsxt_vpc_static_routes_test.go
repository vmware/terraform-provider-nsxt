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

	apivpcs "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	vpcsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	srID          = "static-route-id"
	srDisplayName = "test-static-route"
	srDescription = "Test VPC Static Route"
	srRevision    = int64(1)
	srNetwork     = "10.0.0.0/24"
)

func srAPIResponse() nsxModel.StaticRoutes {
	return nsxModel.StaticRoutes{
		Id:          &srID,
		DisplayName: &srDisplayName,
		Description: &srDescription,
		Revision:    &srRevision,
		Network:     &srNetwork,
	}
}

func minimalSrData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": srDisplayName,
		"description":  srDescription,
		"nsx_id":       srID,
		"network":      srNetwork,
		"next_hop":     []interface{}{},
	}
}

func setupSrMock(t *testing.T, ctrl *gomock.Controller) (*vpcsmocks.MockStaticRoutesClient, func()) {
	mockSDK := vpcsmocks.NewMockStaticRoutesClient(ctrl)
	mockWrapper := &apivpcs.VpcStaticRoutesClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  "project1",
		VPCID:      "vpc1",
	}

	original := cliVpcStaticRoutesClient
	cliVpcStaticRoutesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcStaticRoutesClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcStaticRoutesClient = original }
}

func TestMockResourceNsxtVpcStaticRoutesCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupSrMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), srID).Return(nsxModel.StaticRoutes{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), srID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), srID).Return(srAPIResponse(), nil),
		)

		res := resourceNsxtVpcStaticRoutes()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSrData())

		err := resourceNsxtVpcStaticRoutesCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, srID, d.Id())
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVpcStaticRoutes()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSrData())

		err := resourceNsxtVpcStaticRoutesCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcStaticRoutesRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupSrMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), srID).Return(srAPIResponse(), nil)

		res := resourceNsxtVpcStaticRoutes()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSrData())
		d.SetId(srID)

		err := resourceNsxtVpcStaticRoutesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, srDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupSrMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), srID).Return(nsxModel.StaticRoutes{}, vapiErrors.NotFound{})

		res := resourceNsxtVpcStaticRoutes()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSrData())
		d.SetId(srID)

		err := resourceNsxtVpcStaticRoutesRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcStaticRoutesUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupSrMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), srID, gomock.Any()).Return(srAPIResponse(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), srID).Return(srAPIResponse(), nil),
		)

		res := resourceNsxtVpcStaticRoutes()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSrData())
		d.SetId(srID)

		err := resourceNsxtVpcStaticRoutesUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtVpcStaticRoutesDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupSrMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), srID).Return(nil)

		res := resourceNsxtVpcStaticRoutes()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSrData())
		d.SetId(srID)

		err := resourceNsxtVpcStaticRoutesDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
