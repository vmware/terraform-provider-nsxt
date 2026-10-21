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
// utgomock_resource_nsxt_vpc_lb_cookie_persistence_profile_test.go

var (
	vpcLbClientSslID          = "vpc-lb-clientssl-1"
	vpcLbClientSslDisplayName = "Test VPC LB ClientSsl Profile"
	vpcLbClientSslDescription = "Test vpc lb clientssl profile"
	vpcLbClientSslRevision    = int64(1)
)

func vpcLbClientSslAPIResponse() nsxModel.LBClientSslProfile {
	return nsxModel.LBClientSslProfile{
		Id:          &vpcLbClientSslID,
		DisplayName: &vpcLbClientSslDisplayName,
		Description: &vpcLbClientSslDescription,
		Revision:    &vpcLbClientSslRevision,
		Ciphers:     []string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"},
		Protocols:   []string{"TLS_V1_2"},
	}
}

func vpcLbClientSslData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": vpcLbClientSslDisplayName,
		"description":  vpcLbClientSslDescription,
		"nsx_id":       vpcLbClientSslID,
		"context":      vpcLbPrfContext(),
		"ciphers":      []interface{}{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"},
		"protocols":    []interface{}{"TLS_V1_2"},
	}
}

func setupVpcLbClientSslMock(t *testing.T, ctrl *gomock.Controller) (*vpcsmocks.MockVpcLbClientSslProfilesClient, func()) {
	mockSDK := vpcsmocks.NewMockVpcLbClientSslProfilesClient(ctrl)
	mockWrapper := &apivpcs.VpcLbClientSslProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.VPC,
		ProjectID:  vpcLbPrfProject,
		VPCID:      vpcLbPrfVpc,
	}
	original := cliVpcLbClientSslProfilesClient
	cliVpcLbClientSslProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apivpcs.VpcLbClientSslProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcLbClientSslProfilesClient = original }
}

func TestMockResourceNsxtVpcLBClientSslProfileCreate(t *testing.T) {
	newMock := func(t *testing.T) (*vpcsmocks.MockVpcLbClientSslProfilesClient, func()) {
		vpcLbPrfSetNsxVersion(t, "9.0.0")
		ctrl := gomock.NewController(t)
		mockSDK, restore := setupVpcLbClientSslMock(t, ctrl)
		return mockSDK, func() { restore(); ctrl.Finish() }
	}

	t.Run("success", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		gomock.InOrder(
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID).Return(nsxModel.LBClientSslProfile{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID, gomock.Any()).DoAndReturn(
				func(_, _, _, _ string, obj nsxModel.LBClientSslProfile) error {
					assert.Equal(t, vpcLbClientSslDisplayName, *obj.DisplayName)
					assert.Equal(t, []string{"TLS_V1_2"}, obj.Protocols)
					return nil
				}),
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID).Return(vpcLbClientSslAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
		require.NoError(t, res.Create(d, newGoMockProviderClient()))
		assert.Equal(t, vpcLbClientSslID, d.Id())
		assert.Equal(t, vpcLbClientSslDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcLbClientSslDescription, d.Get("description"))
	})

	t.Run("already exists", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID).Return(vpcLbClientSslAPIResponse(), nil)

		res := resourceNsxtVpcLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
		err := res.Create(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("NSX version lower than 9.0.0", func(t *testing.T) {
		vpcLbPrfSetNsxVersion(t, "4.2.0")
		res := resourceNsxtVpcLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
		err := res.Create(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.0.0")
	})

	t.Run("patch error", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		gomock.InOrder(
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID).Return(nsxModel.LBClientSslProfile{}, vapiErrors.NotFound{}),
			mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID, gomock.Any()).Return(vapiErrors.InternalServerError{}),
		)

		res := resourceNsxtVpcLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
		require.Error(t, res.Create(d, newGoMockProviderClient()))
		assert.Equal(t, "", d.Id())
	})
}

func TestMockResourceNsxtVpcLBClientSslProfileRead(t *testing.T) {
	newMock := func(t *testing.T) (*vpcsmocks.MockVpcLbClientSslProfilesClient, func()) {
		ctrl := gomock.NewController(t)
		mockSDK, restore := setupVpcLbClientSslMock(t, ctrl)
		return mockSDK, func() { restore(); ctrl.Finish() }
	}

	t.Run("success", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID).Return(vpcLbClientSslAPIResponse(), nil)

		res := resourceNsxtVpcLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
		d.SetId(vpcLbClientSslID)
		require.NoError(t, res.Read(d, newGoMockProviderClient()))
		assert.Equal(t, vpcLbClientSslDisplayName, d.Get("display_name"))
		assert.Equal(t, vpcLbClientSslDescription, d.Get("description"))
		assert.Equal(t, 1, d.Get("revision"))
		assert.Equal(t, 1, d.Get("protocols.#"))
	})

	t.Run("not found clears ID", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID).Return(nsxModel.LBClientSslProfile{}, vapiErrors.NotFound{})

		res := resourceNsxtVpcLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
		d.SetId(vpcLbClientSslID)
		require.NoError(t, res.Read(d, newGoMockProviderClient()))
		assert.Equal(t, "", d.Id())
	})

	t.Run("API error is propagated", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID).Return(nsxModel.LBClientSslProfile{}, vapiErrors.InternalServerError{})

		res := resourceNsxtVpcLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
		d.SetId(vpcLbClientSslID)
		require.Error(t, res.Read(d, newGoMockProviderClient()))
	})

	t.Run("empty ID", func(t *testing.T) {
		res := resourceNsxtVpcLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
		require.Error(t, res.Read(d, newGoMockProviderClient()))
	})
}

func TestMockResourceNsxtVpcLBClientSslProfileUpdate(t *testing.T) {
	newMock := func(t *testing.T) (*vpcsmocks.MockVpcLbClientSslProfilesClient, func()) {
		ctrl := gomock.NewController(t)
		mockSDK, restore := setupVpcLbClientSslMock(t, ctrl)
		return mockSDK, func() { restore(); ctrl.Finish() }
	}

	t.Run("success", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		gomock.InOrder(
			// the client SSL profile is updated with PATCH
			mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID).Return(vpcLbClientSslAPIResponse(), nil),
		)

		res := resourceNsxtVpcLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
		d.SetId(vpcLbClientSslID)
		require.NoError(t, res.Update(d, newGoMockProviderClient()))
		assert.Equal(t, vpcLbClientSslDisplayName, d.Get("display_name"))
	})

	t.Run("API error", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtVpcLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
		d.SetId(vpcLbClientSslID)
		require.Error(t, res.Update(d, newGoMockProviderClient()))
	})

	t.Run("empty ID", func(t *testing.T) {
		res := resourceNsxtVpcLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
		require.Error(t, res.Update(d, newGoMockProviderClient()))
	})
}

func TestMockResourceNsxtVpcLBClientSslProfileDelete(t *testing.T) {
	newMock := func(t *testing.T) (*vpcsmocks.MockVpcLbClientSslProfilesClient, func()) {
		ctrl := gomock.NewController(t)
		mockSDK, restore := setupVpcLbClientSslMock(t, ctrl)
		return mockSDK, func() { restore(); ctrl.Finish() }
	}

	t.Run("success", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Delete(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID, gomock.Any()).Return(nil)

		res := resourceNsxtVpcLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
		d.SetId(vpcLbClientSslID)
		require.NoError(t, res.Delete(d, newGoMockProviderClient()))
	})

	t.Run("API error", func(t *testing.T) {
		mockSDK, done := newMock(t)
		defer done()
		mockSDK.EXPECT().Delete(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, vpcLbClientSslID, gomock.Any()).Return(vapiErrors.InternalServerError{})

		res := resourceNsxtVpcLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
		d.SetId(vpcLbClientSslID)
		require.Error(t, res.Delete(d, newGoMockProviderClient()))
	})

	t.Run("empty ID", func(t *testing.T) {
		res := resourceNsxtVpcLBClientSslProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
		require.Error(t, res.Delete(d, newGoMockProviderClient()))
	})
}

func TestMockResourceNsxtVpcLBClientSslProfileSchemaAndImporter(t *testing.T) {
	res := resourceNsxtVpcLBClientSslProfile()
	ctxSchema, ok := res.Schema["context"]
	require.True(t, ok)
	assert.True(t, ctxSchema.Required)
	require.NotNil(t, res.Importer)

	runVpcLbImporterSuite(t, res, "vpc-lb-client-ssl-profiles")
}

func TestMockResourceNsxtVpcLBClientSslProfileExists(t *testing.T) {
	sessionContext := utl.SessionContext{ClientType: utl.VPC, ProjectID: vpcLbPrfProject, VPCID: vpcLbPrfVpc}
	cases := []struct {
		name       string
		err        error
		wantExists bool
		wantErr    bool
	}{
		{name: "exists", wantExists: true},
		{name: "not found", err: vapiErrors.NotFound{}},
		{name: "error", err: vapiErrors.InternalServerError{}, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSDK, restore := setupVpcLbClientSslMock(t, ctrl)
			defer restore()
			mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, "some-id").Return(nsxModel.LBClientSslProfile{}, tc.err)

			exists, err := resourceNsxtVpcLBClientSslProfileExists(sessionContext, "some-id", nil)
			assert.Equal(t, tc.wantExists, exists)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestMockResourceNsxtVpcLBClientSslProfileCreateReadFails(t *testing.T) {
	vpcLbPrfSetNsxVersion(t, "9.0.0")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupVpcLbClientSslMock(t, ctrl)
	defer restore()

	id := vpcLbClientSslID
	gomock.InOrder(
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, id).Return(nsxModel.LBClientSslProfile{}, vapiErrors.NotFound{}),
		mockSDK.EXPECT().Patch(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, id, gomock.Any()).Return(nil),
		mockSDK.EXPECT().Get(vpcLbPrfOrg, vpcLbPrfProject, vpcLbPrfVpc, id).Return(nsxModel.LBClientSslProfile{}, vapiErrors.InternalServerError{}),
	)

	res := resourceNsxtVpcLBClientSslProfile()
	d := schema.TestResourceDataRaw(t, res.Schema, vpcLbClientSslData())
	require.Error(t, res.Create(d, newGoMockProviderClient()))
}
