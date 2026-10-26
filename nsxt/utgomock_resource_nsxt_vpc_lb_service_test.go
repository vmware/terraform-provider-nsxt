//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
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

func TestMockResourceNsxtVpcLBServiceRelaxScaleValidation(t *testing.T) {
	boolPtr := func(b bool) *bool { return &b }
	apiResponseWithRelax := func(relax *bool) nsxModel.LBService {
		obj := vpcLbServiceAPIResponse()
		obj.RelaxScaleValidation = relax
		return obj
	}

	t.Run("create without value does not send it", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.0")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
		defer restore()
		gomock.InOrder(
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(nsxModel.LBService{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID, gomock.Any(), nil).DoAndReturn(
				func(_, _, _, _ string, obj nsxModel.LBService, _ *string) error {
					assert.Nil(t, obj.RelaxScaleValidation)
					return nil
				}),
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(apiResponseWithRelax(boolPtr(true)), nil),
		)

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		require.NoError(t, res.Create(d, newGoMockProviderClient()))
		assert.Equal(t, true, d.Get("relax_scale_validation"))
	})

	t.Run("create with explicit false sends false", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.0")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
		defer restore()
		gomock.InOrder(
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(nsxModel.LBService{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID, gomock.Any(), nil).DoAndReturn(
				func(_, _, _, _ string, obj nsxModel.LBService, _ *string) error {
					require.NotNil(t, obj.RelaxScaleValidation)
					assert.False(t, *obj.RelaxScaleValidation)
					return nil
				}),
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(apiResponseWithRelax(boolPtr(false)), nil),
		)

		data := vpcLbServiceData()
		data["relax_scale_validation"] = false
		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		require.NoError(t, res.Create(d, newGoMockProviderClient()))
		assert.Equal(t, false, d.Get("relax_scale_validation"))
	})

	for _, relax := range []bool{false, true} {
		t.Run(fmt.Sprintf("update preserves value read from NSX (%t)", relax), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
			defer restore()
			gomock.InOrder(
				mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID, gomock.Any(), nil).DoAndReturn(
					func(_, _, _, _ string, obj nsxModel.LBService, _ *string) error {
						require.NotNil(t, obj.RelaxScaleValidation)
						assert.Equal(t, relax, *obj.RelaxScaleValidation)
						return nil
					}),
				mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(apiResponseWithRelax(boolPtr(relax)), nil),
			)

			res := resourceNsxtVpcLBService()
			// The value is not configured, it is known from the refresh preceding the update
			d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
			d.SetId(vpcLbServiceID)
			require.NoError(t, d.Set("relax_scale_validation", relax))
			require.NoError(t, res.Update(d, newGoMockProviderClient()))
			assert.Equal(t, relax, d.Get("relax_scale_validation"))
		})
	}

	t.Run("update without a known value does not send it", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
		defer restore()
		gomock.InOrder(
			mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID, gomock.Any(), nil).DoAndReturn(
				func(_, _, _, _ string, obj nsxModel.LBService, _ *string) error {
					assert.Nil(t, obj.RelaxScaleValidation)
					return nil
				}),
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(vpcLbServiceAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		d.SetId(vpcLbServiceID)
		require.NoError(t, res.Update(d, newGoMockProviderClient()))
	})

	t.Run("read keeps previous value if NSX does not return it", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
		defer restore()
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(vpcLbServiceAPIResponse(), nil)

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		d.SetId(vpcLbServiceID)
		// Seed a previous value, as if it was set by an earlier read
		require.NoError(t, d.Set("relax_scale_validation", true))
		require.NoError(t, res.Read(d, newGoMockProviderClient()))
		v, exists := d.GetOkExists("relax_scale_validation")
		assert.True(t, exists)
		assert.Equal(t, true, v)
	})
}

func TestMockResourceNsxtVpcLBServiceVersionedFields(t *testing.T) {
	snatPaths := []string{
		"/orgs/default/projects/proj1/vpcs/vpc1/ip-address-allocations/ipv4",
		"/orgs/default/projects/proj1/vpcs/vpc1/ip-address-allocations/ipv6",
	}
	boolPtr := func(b bool) *bool { return &b }

	t.Run("create sends snat paths on 9.2.0", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.2.0")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
		defer restore()
		response := vpcLbServiceAPIResponse()
		response.SnatIpAllocationPaths = snatPaths
		response.NetworkApplianceLoadBalancer = boolPtr(true)
		gomock.InOrder(
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(nsxModel.LBService{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID, gomock.Any(), nil).DoAndReturn(
				func(_, _, _, _ string, obj nsxModel.LBService, _ *string) error {
					assert.ElementsMatch(t, snatPaths, obj.SnatIpAllocationPaths)
					assert.Nil(t, obj.NetworkApplianceLoadBalancer)
					return nil
				}),
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(response, nil),
		)

		data := vpcLbServiceData()
		data["snat_ip_allocation_paths"] = []interface{}{snatPaths[0], snatPaths[1]}
		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		require.NoError(t, res.Create(d, newGoMockProviderClient()))
		assert.Equal(t, true, d.Get("network_appliance_load_balancer"))
		assert.Equal(t, 2, d.Get("snat_ip_allocation_paths").(*schema.Set).Len())
	})

	t.Run("read ignores the order of snat paths", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.2.0")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
		defer restore()
		response := vpcLbServiceAPIResponse()
		// NSX returns the paths in a different order than they are configured
		response.SnatIpAllocationPaths = []string{snatPaths[1], snatPaths[0]}
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(response, nil)

		data := vpcLbServiceData()
		data["snat_ip_allocation_paths"] = []interface{}{snatPaths[0], snatPaths[1]}
		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(vpcLbServiceID)
		configured := d.Get("snat_ip_allocation_paths").(*schema.Set)
		require.NoError(t, res.Read(d, newGoMockProviderClient()))
		assert.True(t, configured.Equal(d.Get("snat_ip_allocation_paths")))
	})

	t.Run("read leaves snat paths unset if NSX returns none", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.2.0")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
		defer restore()
		response := vpcLbServiceAPIResponse()
		response.SnatIpAllocationPaths = []string{}
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(response, nil)

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		d.SetId(vpcLbServiceID)
		require.NoError(t, res.Read(d, newGoMockProviderClient()))
		assert.Equal(t, 0, d.Get("snat_ip_allocation_paths").(*schema.Set).Len())
	})

	t.Run("create with snat paths fails before 9.2.0", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.0")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
		defer restore()
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(nsxModel.LBService{}, vapiErrors.NotFound{})

		data := vpcLbServiceData()
		data["snat_ip_allocation_paths"] = []interface{}{snatPaths[0]}
		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		err := res.Create(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})

	t.Run("update on 9.2.0 sends network appliance value and clears snat paths", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.2.0")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
		defer restore()
		gomock.InOrder(
			mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID, gomock.Any(), nil).DoAndReturn(
				func(_, _, _, _ string, obj nsxModel.LBService, _ *string) error {
					require.NotNil(t, obj.NetworkApplianceLoadBalancer)
					assert.False(t, *obj.NetworkApplianceLoadBalancer)
					assert.NotNil(t, obj.SnatIpAllocationPaths)
					assert.Empty(t, obj.SnatIpAllocationPaths)
					return nil
				}),
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(vpcLbServiceAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		d.SetId(vpcLbServiceID)
		require.NoError(t, d.Set("network_appliance_load_balancer", false))
		require.NoError(t, res.Update(d, newGoMockProviderClient()))
	})

	t.Run("update before 9.2.0 does not send the fields", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.0")
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVpcLbServiceMock(t, ctrl)
		defer restore()
		gomock.InOrder(
			mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID, gomock.Any(), nil).DoAndReturn(
				func(_, _, _, _ string, obj nsxModel.LBService, _ *string) error {
					assert.Nil(t, obj.NetworkApplianceLoadBalancer)
					assert.Nil(t, obj.SnatIpAllocationPaths)
					return nil
				}),
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServiceID).Return(vpcLbServiceAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServiceData())
		d.SetId(vpcLbServiceID)
		require.NoError(t, res.Update(d, newGoMockProviderClient()))
	})

	t.Run("schema defines snat paths as a set limited to two", func(t *testing.T) {
		res := resourceNsxtVpcLBService()
		assert.Equal(t, schema.TypeSet, res.Schema["snat_ip_allocation_paths"].Type)
		assert.Equal(t, 2, res.Schema["snat_ip_allocation_paths"].MaxItems)
		assert.True(t, res.Schema["network_appliance_load_balancer"].Computed)
		assert.False(t, res.Schema["network_appliance_load_balancer"].Optional)
	})
}
