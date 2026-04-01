//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/IpsecVpnIkeProfilesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/IpsecVpnIkeProfilesClient.go IpsecVpnIkeProfilesClient

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
	ikemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	ikeProfileID          = "ike-profile-1"
	ikeProfileDisplayName = "ike-profile-fooname"
	ikeProfileDescription = "ike profile mock"
	ikeProfilePath        = "/infra/ipsec-vpn-ike-profiles/ike-profile-1"
	ikeProfileRevision    = int64(1)
	ikeProfileDhGroups    = []string{model.IPSecVpnIkeProfile_DH_GROUPS_GROUP2}
	ikeProfileEncAlgos    = []string{model.IPSecVpnIkeProfile_ENCRYPTION_ALGORITHMS_256}
	ikeProfileIkeVersion  = model.IPSecVpnIkeProfile_IKE_VERSION_V2
	ikeProfileSaLifeTime  = int64(86400)
)

func TestMockResourceNsxtPolicyIPSecVpnIkeProfileCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIkeSDK := ikemocks.NewMockIpsecVpnIkeProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPSecVpnIkeProfileClientContext{
		Client:     mockIkeSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpsecVpnIkeProfilesClient
	defer func() { cliIpsecVpnIkeProfilesClient = originalCli }()
	cliIpsecVpnIkeProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPSecVpnIkeProfileClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockIkeSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockIkeSDK.EXPECT().Get(gomock.Any()).Return(model.IPSecVpnIkeProfile{
			Id:                   &ikeProfileID,
			DisplayName:          &ikeProfileDisplayName,
			Description:          &ikeProfileDescription,
			Path:                 &ikeProfilePath,
			Revision:             &ikeProfileRevision,
			DhGroups:             ikeProfileDhGroups,
			EncryptionAlgorithms: ikeProfileEncAlgos,
			IkeVersion:           &ikeProfileIkeVersion,
			SaLifeTime:           &ikeProfileSaLifeTime,
		}, nil)

		res := resourceNsxtPolicyIPSecVpnIkeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecVpnIkeProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnIkeProfileCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockIkeSDK.EXPECT().Get("existing-id").Return(model.IPSecVpnIkeProfile{Id: &ikeProfileID}, nil)

		res := resourceNsxtPolicyIPSecVpnIkeProfile()
		data := minimalIPSecVpnIkeProfileData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnIkeProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyIPSecVpnIkeProfileRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIkeSDK := ikemocks.NewMockIpsecVpnIkeProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPSecVpnIkeProfileClientContext{
		Client:     mockIkeSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpsecVpnIkeProfilesClient
	defer func() { cliIpsecVpnIkeProfilesClient = originalCli }()
	cliIpsecVpnIkeProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPSecVpnIkeProfileClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockIkeSDK.EXPECT().Get(ikeProfileID).Return(model.IPSecVpnIkeProfile{
			Id:                   &ikeProfileID,
			DisplayName:          &ikeProfileDisplayName,
			Description:          &ikeProfileDescription,
			Path:                 &ikeProfilePath,
			Revision:             &ikeProfileRevision,
			DhGroups:             ikeProfileDhGroups,
			EncryptionAlgorithms: ikeProfileEncAlgos,
			IkeVersion:           &ikeProfileIkeVersion,
			SaLifeTime:           &ikeProfileSaLifeTime,
		}, nil)

		res := resourceNsxtPolicyIPSecVpnIkeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ikeProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnIkeProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, ikeProfileDisplayName, d.Get("display_name"))
		assert.Equal(t, ikeProfileDescription, d.Get("description"))
		assert.Equal(t, ikeProfilePath, d.Get("path"))
		assert.Equal(t, int(ikeProfileRevision), d.Get("revision"))
		assert.Equal(t, ikeProfileID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnIkeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnIkeProfileRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IPSecVpnIkeProfile ID")
	})
}

func TestMockResourceNsxtPolicyIPSecVpnIkeProfileUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIkeSDK := ikemocks.NewMockIpsecVpnIkeProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPSecVpnIkeProfileClientContext{
		Client:     mockIkeSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpsecVpnIkeProfilesClient
	defer func() { cliIpsecVpnIkeProfilesClient = originalCli }()
	cliIpsecVpnIkeProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPSecVpnIkeProfileClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockIkeSDK.EXPECT().Patch(ikeProfileID, gomock.Any()).Return(nil)
		mockIkeSDK.EXPECT().Get(ikeProfileID).Return(model.IPSecVpnIkeProfile{
			Id:                   &ikeProfileID,
			DisplayName:          &ikeProfileDisplayName,
			Description:          &ikeProfileDescription,
			Path:                 &ikeProfilePath,
			Revision:             &ikeProfileRevision,
			DhGroups:             ikeProfileDhGroups,
			EncryptionAlgorithms: ikeProfileEncAlgos,
			IkeVersion:           &ikeProfileIkeVersion,
			SaLifeTime:           &ikeProfileSaLifeTime,
		}, nil)

		res := resourceNsxtPolicyIPSecVpnIkeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecVpnIkeProfileData())
		d.SetId(ikeProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnIkeProfileUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnIkeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecVpnIkeProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnIkeProfileUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IPSecVpnIkeProfile ID")
	})
}

func TestMockResourceNsxtPolicyIPSecVpnIkeProfileDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIkeSDK := ikemocks.NewMockIpsecVpnIkeProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPSecVpnIkeProfileClientContext{
		Client:     mockIkeSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpsecVpnIkeProfilesClient
	defer func() { cliIpsecVpnIkeProfilesClient = originalCli }()
	cliIpsecVpnIkeProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPSecVpnIkeProfileClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockIkeSDK.EXPECT().Delete(ikeProfileID).Return(nil)

		res := resourceNsxtPolicyIPSecVpnIkeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ikeProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnIkeProfileDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnIkeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnIkeProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IPSecVpnIkeProfile ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockIkeSDK.EXPECT().Delete(ikeProfileID).Return(errors.New("API error"))

		res := resourceNsxtPolicyIPSecVpnIkeProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ikeProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnIkeProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalIPSecVpnIkeProfileData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":          ikeProfileDisplayName,
		"description":           ikeProfileDescription,
		"dh_groups":             []interface{}{model.IPSecVpnIkeProfile_DH_GROUPS_GROUP2},
		"encryption_algorithms": []interface{}{model.IPSecVpnIkeProfile_ENCRYPTION_ALGORITHMS_256},
		"ike_version":           ikeProfileIkeVersion,
		"sa_life_time":          int(ikeProfileSaLifeTime),
	}
}
