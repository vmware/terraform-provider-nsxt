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
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	apivpcs "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/vpcs"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	vpcsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/vpcs"
)

// vpcLbPrfContext, vpcLbPrfSetNsxVersion, vpcLbPrf* constants and runVpcLbImporterSuite are defined in
// utgomock_resource_nsxt_vpc_lb_cookie_persistence_profile_test.go

var (
	vpcLbServerSslID          = "vpc-lb-serverssl-1"
	vpcLbServerSslDisplayName = "Test VPC LB ServerSsl Profile"
	vpcLbServerSslDescription = "Test vpc lb serverssl profile"
	vpcLbServerSslRevision    = int64(1)
)

func vpcLbServerSslAPIResponse() nsxModel.LBServerSslProfile {
	return nsxModel.LBServerSslProfile{
		Id:          &vpcLbServerSslID,
		DisplayName: &vpcLbServerSslDisplayName,
		Description: &vpcLbServerSslDescription,
		Revision:    &vpcLbServerSslRevision,
		Ciphers:     []string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"},
		Protocols:   []string{"TLS_V1_2"},
	}
}

func vpcLbServerSslData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": vpcLbServerSslDisplayName,
		"description":  vpcLbServerSslDescription,
		"nsx_id":       vpcLbServerSslID,
		"context":      vpcLbPrfContext(),
		"ciphers":      []interface{}{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"},
		"protocols":    []interface{}{"TLS_V1_2"},
	}
}

func setupVpcLbServerSslMock(t *testing.T, ctrl *gomock.Controller) (*vpcsmocks.MockVpcLbServerSslProfilesClient, func()) {
	mockSDK := vpcsmocks.NewMockVpcLbServerSslProfilesClient(ctrl)
	mockWrapper := &apivpcs.VpcLbServerSslProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  vpcLbPrfProject,
		VPCID:      vpcLbPrfVpc,
	}
	original := cliVpcLbServerSslProfilesClient
	cliVpcLbServerSslProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcLbServerSslProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcLbServerSslProfilesClient = original }
}

func TestMockResourceNsxtVpcLBServerSslProfileCreate(t *testing.T) {
	newMock := func(t *testing.T) (*vpcsmocks.MockVpcLbServerSslProfilesClient, func()) {
		vpcLbPrfSetNsxVersion(t, "9.0.2")
		ctrl := gomock.NewController(t)
		mockSDK, restore := setupVpcLbServerSslMock(t, ctrl)
		return mockSDK, func() { restore(); ctrl.Finish() }
	}

	t.Run("success", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		gomock.InOrder(
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID).Return(nsxModel.LBServerSslProfile{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID, gomock.Any()).DoAndReturn(
				func(_, _, _, _ string, obj nsxModel.LBServerSslProfile) error {
					assert.Equal(t, vpcLbServerSslDisplayName, *obj.DisplayName)
					assert.Equal(t, []string{"TLS_V1_2"}, obj.Protocols)
					return nil
				}),
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID).Return(vpcLbServerSslAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		require.NoError(t, res.Create(d, newGoMockProviderClient()))
		assert.Equal(t, vpcLbServerSslID, d.Id())
		assert.Equal(t, vpcLbServerSslDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcLbServerSslDescription, d.Get("description"))
	})

	t.Run("already exists", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID).Return(vpcLbServerSslAPIResponse(), nil)

		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		err := res.Create(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("NSX version lower than 9.0.2", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "9.0.1")
		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		err := res.Create(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.0.2")
	})

	t.Run("patch error", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		gomock.InOrder(
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID).Return(nsxModel.LBServerSslProfile{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID, gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		require.Error(t, res.Create(d, newGoMockProviderClient()))
		assert.Equal(t, "", d.Id())
	})
}

func TestMockResourceNsxtVpcLBServerSslProfileRead(t *testing.T) {
	newMock := func(t *testing.T) (*vpcsmocks.MockVpcLbServerSslProfilesClient, func()) {
		ctrl := gomock.NewController(t)
		mockSDK, restore := setupVpcLbServerSslMock(t, ctrl)
		return mockSDK, func() { restore(); ctrl.Finish() }
	}

	t.Run("success", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID).Return(vpcLbServerSslAPIResponse(), nil)

		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		d.SetId(vpcLbServerSslID)
		require.NoError(t, res.Read(d, newGoMockProviderClient()))
		assert.Equal(t, vpcLbServerSslDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcLbServerSslDescription, d.Get("description"))
		assert.Equal(t, 1, d.Get("revision"))
		assert.Equal(t, 1, d.Get("protocols.#"))
	})

	t.Run("not found clears ID", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID).Return(nsxModel.LBServerSslProfile{}, vapiErrors.NotFound{})

		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		d.SetId(vpcLbServerSslID)
		require.NoError(t, res.Read(d, newGoMockProviderClient()))
		assert.Equal(t, "", d.Id())
	})

	t.Run("API error is propagated", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID).Return(nsxModel.LBServerSslProfile{}, vapiErrors.InternalServerError{})

		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		d.SetId(vpcLbServerSslID)
		require.Error(t, res.Read(d, newGoMockProviderClient()))
	})

	t.Run("empty ID", func(t *testing.T) {
		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		require.Error(t, res.Read(d, newGoMockProviderClient()))
	})

	t.Run("cache is bypassed for VPC context", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()

		mc := newGoMockProviderClientCacheEnabled()
		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		d.SetId(vpcLbServerSslID)
		require.True(t, isCacheEnabledForRead(d, mc))
		require.False(t, isLbCacheEnabledForRead(d, mc))

		// an empty cache entry would result in a tag patch if the cache path was used
		query := getCacheQueryKey(resourceTypeLBServerSslProfile, d, mc)
		tc := gcache.getTypeCache(resourceTypeLBServerSslProfile)
		tc.data[query] = map[string]*data.StructValue{}
		defer delete(gcache.byTyp, resourceTypeLBServerSslProfile)

		// exactly one Get, and no Patch
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID).Return(vpcLbServerSslAPIResponse(), nil).Times(1)
		require.NoError(t, res.Read(d, mc))
		assert.Equal(t, vpcLbServerSslDisplayName, d.Get("display_name"))
	})
}

func TestMockResourceNsxtVpcLBServerSslProfileUpdate(t *testing.T) {
	newMock := func(t *testing.T) (*vpcsmocks.MockVpcLbServerSslProfilesClient, func()) {
		ctrl := gomock.NewController(t)
		mockSDK, restore := setupVpcLbServerSslMock(t, ctrl)
		return mockSDK, func() { restore(); ctrl.Finish() }
	}

	t.Run("success", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		gomock.InOrder(
			mockSDK.EXPECT().Update(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID, gomock.Any()).Return(vpcLbServerSslAPIResponse(), nil),
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID).Return(vpcLbServerSslAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		d.SetId(vpcLbServerSslID)
		require.NoError(t, res.Update(d, newGoMockProviderClient()))
		assert.Equal(t, vpcLbServerSslDisplayName, d.Get("display_name"))
	})

	t.Run("API error", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Update(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID, gomock.Any()).Return(nsxModel.LBServerSslProfile{}, vapiErrors.InternalServerError{})

		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		d.SetId(vpcLbServerSslID)
		require.Error(t, res.Update(d, newGoMockProviderClient()))
	})

	t.Run("empty ID", func(t *testing.T) {
		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		require.Error(t, res.Update(d, newGoMockProviderClient()))
	})
}

func TestMockResourceNsxtVpcLBServerSslProfileDelete(t *testing.T) {
	newMock := func(t *testing.T) (*vpcsmocks.MockVpcLbServerSslProfilesClient, func()) {
		ctrl := gomock.NewController(t)
		mockSDK, restore := setupVpcLbServerSslMock(t, ctrl)
		return mockSDK, func() { restore(); ctrl.Finish() }
	}

	t.Run("success", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Delete(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID, gomock.Any()).Return(nil)

		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		d.SetId(vpcLbServerSslID)
		require.NoError(t, res.Delete(d, newGoMockProviderClient()))
	})

	t.Run("API error", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Delete(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbServerSslID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		d.SetId(vpcLbServerSslID)
		require.Error(t, res.Delete(d, newGoMockProviderClient()))
	})

	t.Run("empty ID", func(t *testing.T) {
		res := resourceNsxtVpcLBServerSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbServerSslData())
		require.Error(t, res.Delete(d, newGoMockProviderClient()))
	})
}

func TestMockResourceNsxtVpcLBServerSslProfileSchemaAndImporter(t *testing.T) {
	res := resourceNsxtVpcLBServerSslProfile()
	ctxSchema, ok := res.Schema["context"]
	require.True(t, ok)
	assert.True(t, ctxSchema.Required)
	require.NotNil(t, res.Importer)

	runVpcLbImporterSuite(t, res, "vpc-lb-server-ssl-profiles")
}
