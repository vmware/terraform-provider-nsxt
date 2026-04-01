//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/fabric/ComputeManagersClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/fabric/ComputeManagersClient.go ComputeManagersClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"

	apifabric "github.com/vmware/terraform-provider-nsxt/api/nsx/fabric"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	fabricmocks "github.com/vmware/terraform-provider-nsxt/mocks/fabric"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	cmID          = "compute-manager-1"
	cmDisplayName = "vcenter-1"
	cmDescription = "Test compute manager"
	cmServer      = "192.168.1.10"
	cmOriginType  = "vCenter"
	cmRevision    = int64(1)
	cmPort        = int64(443)

	cmPassword   = "admin-password" //nolint:gosec
	cmThumbprint = "AA:BB:CC:DD"    //nolint:gosec
	cmUsername   = "administrator@vsphere.local"
)

// usernamePasswordCredentialStructValue builds the *data.StructValue that the
// resource's Read path expects to receive from the API for a
// UsernamePasswordLoginCredential.
func usernamePasswordCredentialStructValue() *data.StructValue {
	converter := bindings.NewTypeConverter()
	cred := nsxModel.UsernamePasswordLoginCredential{
		Password:       &cmPassword,
		Thumbprint:     &cmThumbprint,
		Username:       &cmUsername,
		CredentialType: nsxModel.UsernamePasswordLoginCredential__TYPE_IDENTIFIER,
	}
	dataValue, errs := converter.ConvertToVapi(cred, nsxModel.UsernamePasswordLoginCredentialBindingType())
	if errs != nil {
		panic(errs[0])
	}
	return dataValue.(*data.StructValue)
}

func minimalComputeManagerData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":             cmDisplayName,
		"description":              cmDescription,
		"server":                   cmServer,
		"origin_type":              cmOriginType,
		"access_level_for_oidc":    nsxModel.ComputeManager_ACCESS_LEVEL_FOR_OIDC_FULL,
		"multi_nsx":                false,
		"reverse_proxy_https_port": 443,
		"credential": []interface{}{
			map[string]interface{}{
				"username_password_login": []interface{}{
					map[string]interface{}{
						"password":   cmPassword,
						"thumbprint": cmThumbprint,
						"username":   cmUsername,
					},
				},
				"saml_login":                  []interface{}{},
				"session_login":               []interface{}{},
				"verifiable_asymmetric_login": []interface{}{},
			},
		},
	}
}

func computeManagerAPIResponse() nsxModel.ComputeManager {
	accessLevel := nsxModel.ComputeManager_ACCESS_LEVEL_FOR_OIDC_FULL
	multiNsx := false
	setOidc := false
	port := cmPort
	return nsxModel.ComputeManager{
		Id:                    &cmID,
		DisplayName:           &cmDisplayName,
		Description:           &cmDescription,
		Server:                &cmServer,
		OriginType:            &cmOriginType,
		Revision:              &cmRevision,
		AccessLevelForOidc:    &accessLevel,
		MultiNsx:              &multiNsx,
		SetAsOidcProvider:     &setOidc,
		ReverseProxyHttpsPort: &port,
		Credential:            usernamePasswordCredentialStructValue(),
	}
}

func setupComputeManagerMock(t *testing.T, ctrl *gomock.Controller) (*fabricmocks.MockComputeManagersClient, func()) {
	mockSDK := fabricmocks.NewMockComputeManagersClient(ctrl)
	mockWrapper := &apifabric.ComputeManagerClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliComputeManagersClient
	cliComputeManagersClient = func(_ utl.SessionContext, _ client.Connector) *apifabric.ComputeManagerClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliComputeManagersClient = originalCli }
}

func TestMockResourceNsxtComputeManagerCreate(t *testing.T) {
	t.Run("Create success with NSX < 9.0.0 (sets CreateServiceAccount)", func(t *testing.T) {
		util.NsxVersion = "8.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupComputeManagerMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Create(gomock.Any()).Return(computeManagerAPIResponse(), nil)
		mockSDK.EXPECT().Get(cmID).Return(computeManagerAPIResponse(), nil)

		res := resourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalComputeManagerData())

		err := resourceNsxtComputeManagerCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, cmID, d.Id())
		assert.Equal(t, cmDisplayName, d.Get("display_name"))
		assert.Equal(t, cmServer, d.Get("server"))
	})

	t.Run("Create success with NSX >= 9.0.0 (ignores CreateServiceAccount)", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupComputeManagerMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Create(gomock.Any()).Return(computeManagerAPIResponse(), nil)
		mockSDK.EXPECT().Get(cmID).Return(computeManagerAPIResponse(), nil)

		res := resourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalComputeManagerData())

		err := resourceNsxtComputeManagerCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, cmID, d.Id())
	})

	t.Run("Create fails when API returns error", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupComputeManagerMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Create(gomock.Any()).Return(nsxModel.ComputeManager{}, errors.New("create API error"))

		res := resourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalComputeManagerData())

		err := resourceNsxtComputeManagerCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "create API error")
	})
}

func TestMockResourceNsxtComputeManagerRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupComputeManagerMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(cmID).Return(computeManagerAPIResponse(), nil)

		res := resourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(cmID)

		err := resourceNsxtComputeManagerRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, cmDisplayName, d.Get("display_name"))
		assert.Equal(t, cmDescription, d.Get("description"))
		assert.Equal(t, cmServer, d.Get("server"))
		assert.Equal(t, int(cmRevision), d.Get("revision"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtComputeManagerRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Read fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupComputeManagerMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(cmID).Return(nsxModel.ComputeManager{}, errors.New("read API error"))

		res := resourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(cmID)

		err := resourceNsxtComputeManagerRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read API error")
	})
}

func TestMockResourceNsxtComputeManagerUpdate(t *testing.T) {
	t.Run("Update success with NSX < 9.0.0", func(t *testing.T) {
		util.NsxVersion = "8.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupComputeManagerMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Update(cmID, gomock.Any()).Return(computeManagerAPIResponse(), nil)
		mockSDK.EXPECT().Get(cmID).Return(computeManagerAPIResponse(), nil)

		res := resourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalComputeManagerData())
		d.SetId(cmID)

		err := resourceNsxtComputeManagerUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, cmDisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalComputeManagerData())

		err := resourceNsxtComputeManagerUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Update fails when API returns error", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupComputeManagerMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Update(cmID, gomock.Any()).Return(nsxModel.ComputeManager{}, errors.New("update API error"))

		res := resourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalComputeManagerData())
		d.SetId(cmID)

		err := resourceNsxtComputeManagerUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "update API error")
	})
}

func TestMockResourceNsxtComputeManagerDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupComputeManagerMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(cmID).Return(nil)

		res := resourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(cmID)

		err := resourceNsxtComputeManagerDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtComputeManagerDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error obtaining logical object id")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupComputeManagerMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(cmID).Return(errors.New("delete API error"))

		res := resourceNsxtComputeManager()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(cmID)

		err := resourceNsxtComputeManagerDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete API error")
	})
}
