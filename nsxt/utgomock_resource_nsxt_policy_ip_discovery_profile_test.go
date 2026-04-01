//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/IpDiscoveryProfilesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/IpDiscoveryProfilesClient.go IpDiscoveryProfilesClient

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
	ipdiscMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	ipDiscoveryProfileID          = "ip-discovery-1"
	ipDiscoveryProfileDisplayName = "ip-discovery-fooname"
	ipDiscoveryProfileDescription = "ip discovery profile mock"
	ipDiscoveryProfilePath        = "/infra/ip-discovery-profiles/ip-discovery-1"
	ipDiscoveryProfileRevision    = int64(1)
	ipDiscoveryProfileArpTimeout  = int64(10)
	ipDiscoveryProfileArpLimit    = int64(1)
	ipDiscoveryProfileNdLimit     = int64(3)
)

func minimalIPDiscoveryProfileModel() model.IPDiscoveryProfile {
	dupFalse := false
	arpLimit := ipDiscoveryProfileArpLimit
	ndLimit := ipDiscoveryProfileNdLimit
	return model.IPDiscoveryProfile{
		Id:                   &ipDiscoveryProfileID,
		DisplayName:          &ipDiscoveryProfileDisplayName,
		Description:          &ipDiscoveryProfileDescription,
		Path:                 &ipDiscoveryProfilePath,
		Revision:             &ipDiscoveryProfileRevision,
		ArpNdBindingTimeout:  &ipDiscoveryProfileArpTimeout,
		DuplicateIpDetection: &model.DuplicateIPDetectionOptions{DuplicateIpDetectionEnabled: &dupFalse},
		IpV4DiscoveryOptions: &model.IPv4DiscoveryOptions{
			ArpSnoopingConfig: &model.ArpSnoopingConfig{
				ArpBindingLimit:    &arpLimit,
				ArpSnoopingEnabled: &dupFalse,
			},
			DhcpSnoopingEnabled: &dupFalse,
			VmtoolsEnabled:      &dupFalse,
		},
		IpV6DiscoveryOptions: &model.IPv6DiscoveryOptions{
			DhcpSnoopingV6Enabled: &dupFalse,
			NdSnoopingConfig: &model.NdSnoopingConfig{
				NdSnoopingEnabled: &dupFalse,
				NdSnoopingLimit:   &ndLimit,
			},
			VmtoolsV6Enabled: &dupFalse,
		},
		TofuEnabled: &dupFalse,
	}
}

func TestMockResourceNsxtPolicyIPDiscoveryProfileCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIPDiscSDK := ipdiscMocks.NewMockIpDiscoveryProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPDiscoveryProfileClientContext{
		Client:     mockIPDiscSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpDiscoveryProfilesClient
	defer func() { cliIpDiscoveryProfilesClient = originalCli }()
	cliIpDiscoveryProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPDiscoveryProfileClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockIPDiscSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockIPDiscSDK.EXPECT().Get(gomock.Any()).Return(minimalIPDiscoveryProfileModel(), nil)

		res := resourceNsxtPolicyIPDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPDiscoveryProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPDiscoveryProfileCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockIPDiscSDK.EXPECT().Get("existing-id").Return(model.IPDiscoveryProfile{Id: &ipDiscoveryProfileID}, nil)

		res := resourceNsxtPolicyIPDiscoveryProfile()
		data := minimalIPDiscoveryProfileData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPDiscoveryProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyIPDiscoveryProfileRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIPDiscSDK := ipdiscMocks.NewMockIpDiscoveryProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPDiscoveryProfileClientContext{
		Client:     mockIPDiscSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpDiscoveryProfilesClient
	defer func() { cliIpDiscoveryProfilesClient = originalCli }()
	cliIpDiscoveryProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPDiscoveryProfileClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockIPDiscSDK.EXPECT().Get(ipDiscoveryProfileID).Return(minimalIPDiscoveryProfileModel(), nil)

		res := resourceNsxtPolicyIPDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ipDiscoveryProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPDiscoveryProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, ipDiscoveryProfileDisplayName, d.Get("display_name"))
		assert.Equal(t, ipDiscoveryProfileDescription, d.Get("description"))
		assert.Equal(t, ipDiscoveryProfilePath, d.Get("path"))
		assert.Equal(t, int(ipDiscoveryProfileRevision), d.Get("revision"))
		assert.Equal(t, ipDiscoveryProfileID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPDiscoveryProfileRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IPDiscoveryProfile ID")
	})
}

func TestMockResourceNsxtPolicyIPDiscoveryProfileUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIPDiscSDK := ipdiscMocks.NewMockIpDiscoveryProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPDiscoveryProfileClientContext{
		Client:     mockIPDiscSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpDiscoveryProfilesClient
	defer func() { cliIpDiscoveryProfilesClient = originalCli }()
	cliIpDiscoveryProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPDiscoveryProfileClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockIPDiscSDK.EXPECT().Patch(ipDiscoveryProfileID, gomock.Any(), gomock.Any()).Return(nil)
		mockIPDiscSDK.EXPECT().Get(ipDiscoveryProfileID).Return(minimalIPDiscoveryProfileModel(), nil)

		res := resourceNsxtPolicyIPDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPDiscoveryProfileData())
		d.SetId(ipDiscoveryProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPDiscoveryProfileUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPDiscoveryProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPDiscoveryProfileUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IPDiscoveryProfile ID")
	})
}

func TestMockResourceNsxtPolicyIPDiscoveryProfileDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIPDiscSDK := ipdiscMocks.NewMockIpDiscoveryProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPDiscoveryProfileClientContext{
		Client:     mockIPDiscSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpDiscoveryProfilesClient
	defer func() { cliIpDiscoveryProfilesClient = originalCli }()
	cliIpDiscoveryProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPDiscoveryProfileClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockIPDiscSDK.EXPECT().Delete(ipDiscoveryProfileID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyIPDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ipDiscoveryProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPDiscoveryProfileDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPDiscoveryProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IPDiscoveryProfile ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockIPDiscSDK.EXPECT().Delete(ipDiscoveryProfileID, gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtPolicyIPDiscoveryProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(ipDiscoveryProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPDiscoveryProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalIPDiscoveryProfileData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": ipDiscoveryProfileDisplayName,
		"description":  ipDiscoveryProfileDescription,
	}
}
