//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/MacDiscoveryProfilesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/MacDiscoveryProfilesClient.go MacDiscoveryProfilesClient

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
	macmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	macDiscoveryProfileID          = "mac-discovery-1"
	macDiscoveryProfileDisplayName = "mac-discovery-fooname"
	macDiscoveryProfileDescription = "mac discovery profile mock"
	macDiscoveryProfilePath        = "/infra/mac-discovery-profiles/mac-discovery-1"
	macDiscoveryProfileRevision    = int64(1)
)

func TestMockResourceNsxtPolicyMacDiscoveryProfileCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMacSDK := macmocks.NewMockMacDiscoveryProfilesClient(ctrl)
	mockWrapper := &cliinfra.MacDiscoveryProfileClientContext{
		Client:     mockMacSDK,
		ClientType: utl.Local,
	}

	originalCli := cliMacDiscoveryProfilesClient
	defer func() { cliMacDiscoveryProfilesClient = originalCli }()
	cliMacDiscoveryProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.MacDiscoveryProfileClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockMacSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockMacSDK.EXPECT().Get(gomock.Any()).Return(model.MacDiscoveryProfile{
			Id:          &macDiscoveryProfileID,
			DisplayName: &macDiscoveryProfileDisplayName,
			Description: &macDiscoveryProfileDescription,
			Path:        &macDiscoveryProfilePath,
			Revision:    &macDiscoveryProfileRevision,
		}, nil)

		res := resourceNsxtPolicyMacDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalMacDiscoveryProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMacDiscoveryProfileCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockMacSDK.EXPECT().Get("existing-id").Return(model.MacDiscoveryProfile{Id: &macDiscoveryProfileID}, nil)

		res := resourceNsxtPolicyMacDiscoveryProfile()
		data := minimalMacDiscoveryProfileData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMacDiscoveryProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyMacDiscoveryProfileRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMacSDK := macmocks.NewMockMacDiscoveryProfilesClient(ctrl)
	mockWrapper := &cliinfra.MacDiscoveryProfileClientContext{
		Client:     mockMacSDK,
		ClientType: utl.Local,
	}

	originalCli := cliMacDiscoveryProfilesClient
	defer func() { cliMacDiscoveryProfilesClient = originalCli }()
	cliMacDiscoveryProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.MacDiscoveryProfileClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockMacSDK.EXPECT().Get(macDiscoveryProfileID).Return(model.MacDiscoveryProfile{
			Id:          &macDiscoveryProfileID,
			DisplayName: &macDiscoveryProfileDisplayName,
			Description: &macDiscoveryProfileDescription,
			Path:        &macDiscoveryProfilePath,
			Revision:    &macDiscoveryProfileRevision,
		}, nil)

		res := resourceNsxtPolicyMacDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(macDiscoveryProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMacDiscoveryProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, macDiscoveryProfileDisplayName, d.Get("display_name"))
		assert.Equal(t, macDiscoveryProfileDescription, d.Get("description"))
		assert.Equal(t, macDiscoveryProfilePath, d.Get("path"))
		assert.Equal(t, int(macDiscoveryProfileRevision), d.Get("revision"))
		assert.Equal(t, macDiscoveryProfileID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyMacDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMacDiscoveryProfileRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining MacDiscoveryProfile ID")
	})
}

func TestMockResourceNsxtPolicyMacDiscoveryProfileUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMacSDK := macmocks.NewMockMacDiscoveryProfilesClient(ctrl)
	mockWrapper := &cliinfra.MacDiscoveryProfileClientContext{
		Client:     mockMacSDK,
		ClientType: utl.Local,
	}

	originalCli := cliMacDiscoveryProfilesClient
	defer func() { cliMacDiscoveryProfilesClient = originalCli }()
	cliMacDiscoveryProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.MacDiscoveryProfileClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockMacSDK.EXPECT().Update(macDiscoveryProfileID, gomock.Any(), gomock.Any()).Return(model.MacDiscoveryProfile{
			Id:          &macDiscoveryProfileID,
			DisplayName: &macDiscoveryProfileDisplayName,
			Description: &macDiscoveryProfileDescription,
			Path:        &macDiscoveryProfilePath,
			Revision:    &macDiscoveryProfileRevision,
		}, nil)
		mockMacSDK.EXPECT().Get(macDiscoveryProfileID).Return(model.MacDiscoveryProfile{
			Id:          &macDiscoveryProfileID,
			DisplayName: &macDiscoveryProfileDisplayName,
			Description: &macDiscoveryProfileDescription,
			Path:        &macDiscoveryProfilePath,
			Revision:    &macDiscoveryProfileRevision,
		}, nil)

		res := resourceNsxtPolicyMacDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalMacDiscoveryProfileData())
		d.Set("revision", 1)
		d.SetId(macDiscoveryProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMacDiscoveryProfileUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyMacDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalMacDiscoveryProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMacDiscoveryProfileUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining MacDiscoveryProfile ID")
	})
}

func TestMockResourceNsxtPolicyMacDiscoveryProfileDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMacSDK := macmocks.NewMockMacDiscoveryProfilesClient(ctrl)
	mockWrapper := &cliinfra.MacDiscoveryProfileClientContext{
		Client:     mockMacSDK,
		ClientType: utl.Local,
	}

	originalCli := cliMacDiscoveryProfilesClient
	defer func() { cliMacDiscoveryProfilesClient = originalCli }()
	cliMacDiscoveryProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.MacDiscoveryProfileClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockMacSDK.EXPECT().Delete(macDiscoveryProfileID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyMacDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(macDiscoveryProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMacDiscoveryProfileDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyMacDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMacDiscoveryProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining MacDiscoveryProfile ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockMacSDK.EXPECT().Delete(macDiscoveryProfileID, gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtPolicyMacDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(macDiscoveryProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyMacDiscoveryProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalMacDiscoveryProfileData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": macDiscoveryProfileDisplayName,
		"description":  macDiscoveryProfileDescription,
	}
}
