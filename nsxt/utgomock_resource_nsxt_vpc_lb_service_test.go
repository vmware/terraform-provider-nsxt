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
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	apivpcs "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	vpcsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs"
)

// vpcLbPrfContext, vpcLbPrfSetNsxVersion, vpcLbPrf* constants and runVpcLbImporterSuite are defined in
// utgomock_vpc_lb_common_test.go

var (
	vpcLbServiceID          = "vpc-lb-service-1"
	vpcLbServiceDisplayName = "Test VPC LB Service"
	vpcLbServiceDescription = "Test vpc lb service"
	vpcLbServicePath        = "/orgs/default/projects/proj1/vpcs/vpc1/vpc-lbs/vpc-lb-service-1"
	vpcLbServiceRevision    = int64(1)
	vpcLbServiceEnabled     = true
	vpcLbServiceLogLevel    = nsxModel.LBService_ERROR_LOG_LEVEL_ERROR
	vpcLbServiceSize        = nsxModel.LBService_SIZE_SMALL
)

func vpcLbServiceAPIResponse() nsxModel.LBService {
	return nsxModel.LBService{
		Id:            &vpcLbServiceID,
		DisplayName:   &vpcLbServiceDisplayName,
		Description:   &vpcLbServiceDescription,
		Path:          &vpcLbServicePath,
		Revision:      &vpcLbServiceRevision,
		Enabled:       &vpcLbServiceEnabled,
		ErrorLogLevel: &vpcLbServiceLogLevel,
		Size:          &vpcLbServiceSize,
	}
}

func vpcLbServiceData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    vpcLbServiceDisplayName,
		"description":     vpcLbServiceDescription,
		"nsx_id":          vpcLbServiceID,
		"context":         vpcLbPrfContext(),
		"enabled":         true,
		"error_log_level": "ERROR",
		"size":            "SMALL",
	}
}

func setupVpcLbServiceMock(t *testing.T, ctrl *gomock.Controller) (*vpcsmocks.MockVpcLbsClient, func()) {
	mockSDK := vpcsmocks.NewMockVpcLbsClient(ctrl)
	mockWrapper := &apivpcs.VpcLbClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  vpcLbPrfProject,
		VPCID:      vpcLbPrfVpc,
	}
	original := cliVpcLbsClient
	cliVpcLbsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcLbClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcLbsClient = original }
}

func TestMockResourceNsxtVpcLBServiceCreate(t *testing.T) {
	newMock := func(t *testing.T) (*vpcsmocks.MockVpcLbsClient, func()) {
		vpcLbPrfSetNsxVersion(t, "9.0.0")
		ctrl := gomock.NewController(t)
		mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
		return mockSDK, func() { restore(); ctrl.Finish() }
	}

	t.Run("success", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		gomock.InOrder(
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(nsxModel.LBService{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID, gomock.Any(), nil).DoAndReturn(
				func(_, _, _, _ string, obj nsxModel.LBService, action *string) error {
					assert.Equal(t, vpcLbServiceDisplayName, *obj.DisplayName)
					assert.True(t, *obj.Enabled)
					assert.Equal(t, "ERROR", *obj.ErrorLogLevel)
					assert.Equal(t, "SMALL", *obj.Size)
					return nil
				}),
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(vpcLbServiceAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		require.NoError(t, res.Create(d, newGoMockProviderClient()))
		assert.Equal(t, vpcLbServiceID, d.Id())
		assert.Equal(t, vpcLbServiceDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcLbServicePath, d.Get("path"))
		assert.Equal(t, "ERROR", d.Get("error_log_level"))
	})

	t.Run("already exists", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(vpcLbServiceAPIResponse(), nil)

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		err := res.Create(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("NSX version lower than 9.0.0", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "4.2.0")
		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		err := res.Create(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.0.0")
	})

	t.Run("patch error", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		gomock.InOrder(
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(nsxModel.LBService{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID, gomock.Any(), nil).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		require.Error(t, res.Create(d, newGoMockProviderClient()))
		assert.Equal(t, "", d.Id())
	})
}

func TestMockResourceNsxtVpcLBServiceRead(t *testing.T) {
	newMock := func(t *testing.T) (*vpcsmocks.MockVpcLbsClient, func()) {
		ctrl := gomock.NewController(t)
		mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
		return mockSDK, func() { restore(); ctrl.Finish() }
	}

	t.Run("success", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(vpcLbServiceAPIResponse(), nil)

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		d.SetId(vpcLbServiceID)
		require.NoError(t, res.Read(d, newGoMockProviderClient()))
		assert.Equal(t, vpcLbServiceDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcLbServiceDescription, d.Get("description"))
		assert.Equal(t, vpcLbServicePath, d.Get("path"))
		assert.Equal(t, 1, d.Get("revision"))
		assert.Equal(t, true, d.Get("enabled"))
		assert.Equal(t, "ERROR", d.Get("error_log_level"))
		assert.Equal(t, "SMALL", d.Get("size"))
	})

	t.Run("not found clears ID", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(nsxModel.LBService{}, vapiErrors.NotFound{})

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		d.SetId(vpcLbServiceID)
		require.NoError(t, res.Read(d, newGoMockProviderClient()))
		assert.Equal(t, "", d.Id())
	})

	t.Run("API error is propagated", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(nsxModel.LBService{}, vapiErrors.InternalServerError{})

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		d.SetId(vpcLbServiceID)
		require.Error(t, res.Read(d, newGoMockProviderClient()))
	})

	t.Run("empty ID", func(t *testing.T) {
		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		require.Error(t, res.Read(d, newGoMockProviderClient()))
	})

}

func TestMockResourceNsxtVpcLBServiceUpdate(t *testing.T) {
	newMock := func(t *testing.T) (*vpcsmocks.MockVpcLbsClient, func()) {
		ctrl := gomock.NewController(t)
		mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
		return mockSDK, func() { restore(); ctrl.Finish() }
	}

	t.Run("success", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		gomock.InOrder(
			mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID, gomock.Any(), nil).Return(nil),
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(vpcLbServiceAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		d.SetId(vpcLbServiceID)
		require.NoError(t, res.Update(d, newGoMockProviderClient()))
		assert.Equal(t, vpcLbServiceDisplayName, d.Get("display_name"))
	})

	t.Run("API error", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID, gomock.Any(), nil).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		d.SetId(vpcLbServiceID)
		require.Error(t, res.Update(d, newGoMockProviderClient()))
	})

	t.Run("empty ID", func(t *testing.T) {
		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		require.Error(t, res.Update(d, newGoMockProviderClient()))
	})
}

func TestMockResourceNsxtVpcLBServiceDelete(t *testing.T) {
	newMock := func(t *testing.T) (*vpcsmocks.MockVpcLbsClient, func()) {
		ctrl := gomock.NewController(t)
		mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
		return mockSDK, func() { restore(); ctrl.Finish() }
	}

	t.Run("success", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Delete(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID, gomock.Any()).DoAndReturn(
			func(_, _, _, _ string, force *bool) error {
				require.NotNil(t, force)
				assert.True(t, *force)
				return nil
			})

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		d.SetId(vpcLbServiceID)
		require.NoError(t, res.Delete(d, newGoMockProviderClient()))
	})

	t.Run("API error", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Delete(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		d.SetId(vpcLbServiceID)
		require.Error(t, res.Delete(d, newGoMockProviderClient()))
	})

	t.Run("empty ID", func(t *testing.T) {
		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		require.Error(t, res.Delete(d, newGoMockProviderClient()))
	})
}

func TestMockResourceNsxtVpcLBServiceSchemaAndImporter(t *testing.T) {
	res := resourceNsxtVpcLBService()
	ctxSchema, ok := res.Schema["context"]
	require.True(t, ok)
	assert.True(t, ctxSchema.Required)
	require.NotNil(t, res.Importer)

	runVpcLbImporterSuite(t, res, "vpc-lbs")
}
