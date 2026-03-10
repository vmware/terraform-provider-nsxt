// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/IpsecVpnTunnelProfilesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/IpsecVpnTunnelProfilesClient.go IpsecVpnTunnelProfilesClient

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
	tunnelmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	tunnelProfileID   = "tunnel-profile-1"
	tunnelDisplayName = "tunnel-profile-fooname"
	tunnelDescription = "tunnel profile mock"
	tunnelPath        = "/infra/ipsec-vpn-tunnel-profiles/tunnel-profile-1"
	tunnelRevision    = int64(1)
	tunnelDfPolicy    = model.IPSecVpnTunnelProfile_DF_POLICY_COPY
	tunnelDhGroups    = []string{model.IPSecVpnTunnelProfile_DH_GROUPS_GROUP2}
	tunnelEncAlgos    = []string{model.IPSecVpnTunnelProfile_ENCRYPTION_ALGORITHMS_AES_128}
	tunnelSaLifeTime  = int64(3600)
	tunnelEnablePfs   = true
)

func TestMockResourceNsxtPolicyIPSecVpnTunnelProfileCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTunnelSDK := tunnelmocks.NewMockIpsecVpnTunnelProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPSecVpnTunnelProfileClientContext{
		Client:     mockTunnelSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpsecVpnTunnelProfilesClient
	defer func() { cliIpsecVpnTunnelProfilesClient = originalCli }()
	cliIpsecVpnTunnelProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPSecVpnTunnelProfileClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockTunnelSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockTunnelSDK.EXPECT().Get(gomock.Any()).Return(model.IPSecVpnTunnelProfile{
			Id:                          &tunnelProfileID,
			DisplayName:                 &tunnelDisplayName,
			Description:                 &tunnelDescription,
			Path:                        &tunnelPath,
			Revision:                    &tunnelRevision,
			DfPolicy:                    &tunnelDfPolicy,
			DhGroups:                    tunnelDhGroups,
			EncryptionAlgorithms:        tunnelEncAlgos,
			SaLifeTime:                  &tunnelSaLifeTime,
			EnablePerfectForwardSecrecy: &tunnelEnablePfs,
		}, nil)

		res := resourceNsxtPolicyIPSecVpnTunnelProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecVpnTunnelProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnTunnelProfileCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockTunnelSDK.EXPECT().Get("existing-id").Return(model.IPSecVpnTunnelProfile{Id: &tunnelProfileID}, nil)

		res := resourceNsxtPolicyIPSecVpnTunnelProfile()
		data := minimalIPSecVpnTunnelProfileData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnTunnelProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyIPSecVpnTunnelProfileRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTunnelSDK := tunnelmocks.NewMockIpsecVpnTunnelProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPSecVpnTunnelProfileClientContext{
		Client:     mockTunnelSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpsecVpnTunnelProfilesClient
	defer func() { cliIpsecVpnTunnelProfilesClient = originalCli }()
	cliIpsecVpnTunnelProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPSecVpnTunnelProfileClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockTunnelSDK.EXPECT().Get(tunnelProfileID).Return(model.IPSecVpnTunnelProfile{
			Id:                          &tunnelProfileID,
			DisplayName:                 &tunnelDisplayName,
			Description:                 &tunnelDescription,
			Path:                        &tunnelPath,
			Revision:                    &tunnelRevision,
			DfPolicy:                    &tunnelDfPolicy,
			DhGroups:                    tunnelDhGroups,
			EncryptionAlgorithms:        tunnelEncAlgos,
			SaLifeTime:                  &tunnelSaLifeTime,
			EnablePerfectForwardSecrecy: &tunnelEnablePfs,
		}, nil)

		res := resourceNsxtPolicyIPSecVpnTunnelProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(tunnelProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnTunnelProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, tunnelDisplayName, d.Get("display_name"))
		assert.Equal(t, tunnelDescription, d.Get("description"))
		assert.Equal(t, tunnelPath, d.Get("path"))
		assert.Equal(t, int(tunnelRevision), d.Get("revision"))
		assert.Equal(t, tunnelProfileID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnTunnelProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnTunnelProfileRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IPSecVpnTunnelProfile ID")
	})
}

func TestMockResourceNsxtPolicyIPSecVpnTunnelProfileUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTunnelSDK := tunnelmocks.NewMockIpsecVpnTunnelProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPSecVpnTunnelProfileClientContext{
		Client:     mockTunnelSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpsecVpnTunnelProfilesClient
	defer func() { cliIpsecVpnTunnelProfilesClient = originalCli }()
	cliIpsecVpnTunnelProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPSecVpnTunnelProfileClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockTunnelSDK.EXPECT().Patch(tunnelProfileID, gomock.Any()).Return(nil)
		mockTunnelSDK.EXPECT().Get(tunnelProfileID).Return(model.IPSecVpnTunnelProfile{
			Id:                          &tunnelProfileID,
			DisplayName:                 &tunnelDisplayName,
			Description:                 &tunnelDescription,
			Path:                        &tunnelPath,
			Revision:                    &tunnelRevision,
			DfPolicy:                    &tunnelDfPolicy,
			DhGroups:                    tunnelDhGroups,
			EncryptionAlgorithms:        tunnelEncAlgos,
			SaLifeTime:                  &tunnelSaLifeTime,
			EnablePerfectForwardSecrecy: &tunnelEnablePfs,
		}, nil)

		res := resourceNsxtPolicyIPSecVpnTunnelProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecVpnTunnelProfileData())
		d.SetId(tunnelProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnTunnelProfileUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnTunnelProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecVpnTunnelProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnTunnelProfileUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IPSecVpnTunnelProfile ID")
	})
}

func TestMockResourceNsxtPolicyIPSecVpnTunnelProfileDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTunnelSDK := tunnelmocks.NewMockIpsecVpnTunnelProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPSecVpnTunnelProfileClientContext{
		Client:     mockTunnelSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpsecVpnTunnelProfilesClient
	defer func() { cliIpsecVpnTunnelProfilesClient = originalCli }()
	cliIpsecVpnTunnelProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPSecVpnTunnelProfileClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockTunnelSDK.EXPECT().Delete(tunnelProfileID).Return(nil)

		res := resourceNsxtPolicyIPSecVpnTunnelProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(tunnelProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnTunnelProfileDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnTunnelProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnTunnelProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IPSecVpnTunnelProfile ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockTunnelSDK.EXPECT().Delete(tunnelProfileID).Return(errors.New("API error"))

		res := resourceNsxtPolicyIPSecVpnTunnelProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(tunnelProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnTunnelProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalIPSecVpnTunnelProfileData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":                   tunnelDisplayName,
		"description":                    tunnelDescription,
		"df_policy":                      tunnelDfPolicy,
		"dh_groups":                      []interface{}{model.IPSecVpnTunnelProfile_DH_GROUPS_GROUP2},
		"encryption_algorithms":          []interface{}{model.IPSecVpnTunnelProfile_ENCRYPTION_ALGORITHMS_AES_128},
		"enable_perfect_forward_secrecy": tunnelEnablePfs,
		"sa_life_time":                   int(tunnelSaLifeTime),
	}
}
