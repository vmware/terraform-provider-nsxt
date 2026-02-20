// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/SpoofguardProfilesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/SpoofguardProfilesClient.go SpoofguardProfilesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	spoofmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	spoofDisplayName             = "spoof-fooname"
	spoofDescription             = "spoof mock description"
	spoofPath                    = "/infra/spoofguard-profiles/my-spoof-profile"
	spoofRevision                = int64(1)
	spoofID                      = "my-spoof-profile"
	spoofAddressBindingAllowlist = true
)

func TestMockResourceNsxtPolicySpoofGuardProfileCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpoofSDK := spoofmocks.NewMockSpoofguardProfilesClient(ctrl)
	mockWrapper := &cliinfra.SpoofGuardProfileClientContext{
		Client:     mockSpoofSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSpoofguardProfilesClient
	defer func() { cliSpoofguardProfilesClient = originalCli }()
	cliSpoofguardProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SpoofGuardProfileClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockSpoofSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockSpoofSDK.EXPECT().Get(gomock.Any()).Return(model.SpoofGuardProfile{
			DisplayName:             &spoofDisplayName,
			Description:             &spoofDescription,
			Path:                    &spoofPath,
			Revision:                &spoofRevision,
			AddressBindingAllowlist: &spoofAddressBindingAllowlist,
		}, nil)

		res := resourceNsxtPolicySpoofGuardProfile("")
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":              spoofDisplayName,
			"description":               spoofDescription,
			"address_binding_allowlist": spoofAddressBindingAllowlist,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySpoofGuardProfileCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockSpoofSDK.EXPECT().Get("existing-id").Return(model.SpoofGuardProfile{Id: &spoofID}, nil)

		res := resourceNsxtPolicySpoofGuardProfile("")
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"nsx_id":                    "existing-id",
			"display_name":              spoofDisplayName,
			"description":               spoofDescription,
			"address_binding_allowlist": spoofAddressBindingAllowlist,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySpoofGuardProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicySpoofGuardProfileRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpoofSDK := spoofmocks.NewMockSpoofguardProfilesClient(ctrl)
	mockWrapper := &cliinfra.SpoofGuardProfileClientContext{
		Client:     mockSpoofSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSpoofguardProfilesClient
	defer func() { cliSpoofguardProfilesClient = originalCli }()
	cliSpoofguardProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SpoofGuardProfileClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockSpoofSDK.EXPECT().Get(spoofID).Return(model.SpoofGuardProfile{
			DisplayName:             &spoofDisplayName,
			Description:             &spoofDescription,
			Path:                    &spoofPath,
			Revision:                &spoofRevision,
			AddressBindingAllowlist: &spoofAddressBindingAllowlist,
		}, nil)

		res := resourceNsxtPolicySpoofGuardProfile("")
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(spoofID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySpoofGuardProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, spoofDisplayName, d.Get("display_name"))
		assert.Equal(t, spoofDescription, d.Get("description"))
		assert.Equal(t, spoofPath, d.Get("path"))
		assert.Equal(t, int(spoofRevision), d.Get("revision"))
		assert.Equal(t, spoofAddressBindingAllowlist, d.Get("address_binding_allowlist"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySpoofGuardProfile("")
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySpoofGuardProfileRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining SpoofGuardProfile ID")
	})
}

func TestMockResourceNsxtPolicySpoofGuardProfileUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpoofSDK := spoofmocks.NewMockSpoofguardProfilesClient(ctrl)
	mockWrapper := &cliinfra.SpoofGuardProfileClientContext{
		Client:     mockSpoofSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSpoofguardProfilesClient
	defer func() { cliSpoofguardProfilesClient = originalCli }()
	cliSpoofguardProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SpoofGuardProfileClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockSpoofSDK.EXPECT().Patch(spoofID, gomock.Any(), gomock.Any()).Return(nil)
		mockSpoofSDK.EXPECT().Get(spoofID).Return(model.SpoofGuardProfile{
			DisplayName:             &spoofDisplayName,
			Description:             &spoofDescription,
			Path:                    &spoofPath,
			Revision:                &spoofRevision,
			AddressBindingAllowlist: &spoofAddressBindingAllowlist,
		}, nil)

		res := resourceNsxtPolicySpoofGuardProfile("")
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":              spoofDisplayName,
			"description":               spoofDescription,
			"address_binding_allowlist": spoofAddressBindingAllowlist,
		})
		d.SetId(spoofID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySpoofGuardProfileUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySpoofGuardProfile("")
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":              spoofDisplayName,
			"description":               spoofDescription,
			"address_binding_allowlist": spoofAddressBindingAllowlist,
		})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySpoofGuardProfileUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining SpoofGuardProfile ID")
	})
}

func TestMockResourceNsxtPolicySpoofGuardProfileDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSpoofSDK := spoofmocks.NewMockSpoofguardProfilesClient(ctrl)
	mockWrapper := &cliinfra.SpoofGuardProfileClientContext{
		Client:     mockSpoofSDK,
		ClientType: utl.Local,
	}

	originalCli := cliSpoofguardProfilesClient
	defer func() { cliSpoofguardProfilesClient = originalCli }()
	cliSpoofguardProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.SpoofGuardProfileClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockSpoofSDK.EXPECT().Delete(spoofID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicySpoofGuardProfile("")
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(spoofID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySpoofGuardProfileDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySpoofGuardProfile("")
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySpoofGuardProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining SpoofGuardProfile ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockSpoofSDK.EXPECT().Delete(spoofID, gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtPolicySpoofGuardProfile("")
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(spoofID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicySpoofGuardProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}
