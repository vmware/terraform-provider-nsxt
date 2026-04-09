//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"errors"
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
	vpcAttachmentID          = "attachment-test-id"
	vpcAttachmentDisplayName = "test-attachment"
	vpcAttachmentDescription = "Test VPC Attachment"
	vpcAttachmentRevision    = int64(1)
	vpcAttachmentParentPath  = "/orgs/default/projects/project1/vpcs/vpc1"
	vpcConnProfilePath       = "/orgs/default/projects/project1/vpc-connectivity-profiles/profile1"
)

func vpcAttachmentAPIResponse() nsxModel.VpcAttachment {
	return nsxModel.VpcAttachment{
		Id:                     &vpcAttachmentID,
		DisplayName:            &vpcAttachmentDisplayName,
		Description:            &vpcAttachmentDescription,
		Revision:               &vpcAttachmentRevision,
		VpcConnectivityProfile: &vpcConnProfilePath,
	}
}

func minimalVpcAttachmentData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":             vpcAttachmentDisplayName,
		"description":              vpcAttachmentDescription,
		"nsx_id":                   vpcAttachmentID,
		"parent_path":              vpcAttachmentParentPath,
		"vpc_connectivity_profile": vpcConnProfilePath,
	}
}

func setupVpcAttachmentMock(t *testing.T, ctrl *gomock.Controller) (*vpcsmocks.MockAttachmentsClient, func()) {
	mockSDK := vpcsmocks.NewMockAttachmentsClient(ctrl)
	mockWrapper := &apivpcs.VpcAttachmentClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  "project1",
		VPCID:      "vpc1",
	}

	original := cliVpcAttachmentsClient
	cliVpcAttachmentsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcAttachmentClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcAttachmentsClient = original }
}

func TestMockResourceNsxtVpcAttachmentCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcAttachmentMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcAttachmentID).Return(nsxModel.VpcAttachment{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), vpcAttachmentID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcAttachmentID).Return(vpcAttachmentAPIResponse(), nil),
		)

		res := resourceNsxtVpcAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcAttachmentData())

		err := resourceNsxtVpcAttachmentCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcAttachmentID, d.Id())
	})

	t.Run("Create fails on NSX version below 9.0.0", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVpcAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcAttachmentData())

		err := resourceNsxtVpcAttachmentCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.0.0")
	})
}

func TestMockResourceNsxtVpcAttachmentRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcAttachmentMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcAttachmentID).Return(vpcAttachmentAPIResponse(), nil)

		res := resourceNsxtVpcAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcAttachmentData())
		d.SetId(vpcAttachmentID)

		err := resourceNsxtVpcAttachmentRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vpcAttachmentDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcAttachmentMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcAttachmentID).Return(nsxModel.VpcAttachment{}, notFoundErr)

		res := resourceNsxtVpcAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcAttachmentData())
		d.SetId(vpcAttachmentID)

		err := resourceNsxtVpcAttachmentRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcAttachmentUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcAttachmentMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), vpcAttachmentID, gomock.Any()).Return(vpcAttachmentAPIResponse(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), vpcAttachmentID).Return(vpcAttachmentAPIResponse(), nil),
		)

		res := resourceNsxtVpcAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcAttachmentData())
		d.SetId(vpcAttachmentID)

		err := resourceNsxtVpcAttachmentUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtVpcAttachmentDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcAttachmentMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), vpcAttachmentID).Return(nil)

		res := resourceNsxtVpcAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcAttachmentData())
		d.SetId(vpcAttachmentID)

		err := resourceNsxtVpcAttachmentDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete API error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcAttachmentMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), vpcAttachmentID).Return(errors.New("delete failed"))

		res := resourceNsxtVpcAttachment()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVpcAttachmentData())
		d.SetId(vpcAttachmentID)

		err := resourceNsxtVpcAttachmentDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
