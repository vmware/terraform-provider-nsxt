// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/IpsecVpnDpdProfilesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/IpsecVpnDpdProfilesClient.go IpsecVpnDpdProfilesClient

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
	dpdmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	dpdProfileID            = "dpd-profile-1"
	dpdProfileDisplayName   = "dpd-profile-fooname"
	dpdProfileDescription   = "dpd profile mock"
	dpdProfilePath          = "/infra/ipsec-vpn-dpd-profiles/dpd-profile-1"
	dpdProfileRevision      = int64(1)
	dpdProfileProbeInterval = int64(60)
	dpdProfileProbeMode     = model.IPSecVpnDpdProfile_DPD_PROBE_MODE_PERIODIC
	dpdProfileEnabled       = true
	dpdProfileRetryCount    = int64(10)
)

func TestMockResourceNsxtPolicyIPSecVpnDpdProfileCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDpdSDK := dpdmocks.NewMockIpsecVpnDpdProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPSecVpnDpdProfileClientContext{
		Client:     mockDpdSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpsecVpnDpdProfilesClient
	defer func() { cliIpsecVpnDpdProfilesClient = originalCli }()
	cliIpsecVpnDpdProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPSecVpnDpdProfileClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockDpdSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockDpdSDK.EXPECT().Get(gomock.Any()).Return(model.IPSecVpnDpdProfile{
			Id:               &dpdProfileID,
			DisplayName:      &dpdProfileDisplayName,
			Description:      &dpdProfileDescription,
			Path:             &dpdProfilePath,
			Revision:         &dpdProfileRevision,
			DpdProbeInterval: &dpdProfileProbeInterval,
			DpdProbeMode:     &dpdProfileProbeMode,
			Enabled:          &dpdProfileEnabled,
			RetryCount:       &dpdProfileRetryCount,
		}, nil)

		res := resourceNsxtPolicyIPSecVpnDpdProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecVpnDpdProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnDpdProfileCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockDpdSDK.EXPECT().Get("existing-id").Return(model.IPSecVpnDpdProfile{Id: &dpdProfileID}, nil)

		res := resourceNsxtPolicyIPSecVpnDpdProfile()
		data := minimalIPSecVpnDpdProfileData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnDpdProfileCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyIPSecVpnDpdProfileRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDpdSDK := dpdmocks.NewMockIpsecVpnDpdProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPSecVpnDpdProfileClientContext{
		Client:     mockDpdSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpsecVpnDpdProfilesClient
	defer func() { cliIpsecVpnDpdProfilesClient = originalCli }()
	cliIpsecVpnDpdProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPSecVpnDpdProfileClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockDpdSDK.EXPECT().Get(dpdProfileID).Return(model.IPSecVpnDpdProfile{
			Id:               &dpdProfileID,
			DisplayName:      &dpdProfileDisplayName,
			Description:      &dpdProfileDescription,
			Path:             &dpdProfilePath,
			Revision:         &dpdProfileRevision,
			DpdProbeInterval: &dpdProfileProbeInterval,
			DpdProbeMode:     &dpdProfileProbeMode,
			Enabled:          &dpdProfileEnabled,
			RetryCount:       &dpdProfileRetryCount,
		}, nil)

		res := resourceNsxtPolicyIPSecVpnDpdProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dpdProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnDpdProfileRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, dpdProfileDisplayName, d.Get("display_name"))
		assert.Equal(t, dpdProfileDescription, d.Get("description"))
		assert.Equal(t, dpdProfilePath, d.Get("path"))
		assert.Equal(t, int(dpdProfileRevision), d.Get("revision"))
		assert.Equal(t, dpdProfileID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnDpdProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnDpdProfileRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IPSecVpnDpdProfile ID")
	})
}

func TestMockResourceNsxtPolicyIPSecVpnDpdProfileUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDpdSDK := dpdmocks.NewMockIpsecVpnDpdProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPSecVpnDpdProfileClientContext{
		Client:     mockDpdSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpsecVpnDpdProfilesClient
	defer func() { cliIpsecVpnDpdProfilesClient = originalCli }()
	cliIpsecVpnDpdProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPSecVpnDpdProfileClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockDpdSDK.EXPECT().Patch(dpdProfileID, gomock.Any()).Return(nil)
		mockDpdSDK.EXPECT().Get(dpdProfileID).Return(model.IPSecVpnDpdProfile{
			Id:               &dpdProfileID,
			DisplayName:      &dpdProfileDisplayName,
			Description:      &dpdProfileDescription,
			Path:             &dpdProfilePath,
			Revision:         &dpdProfileRevision,
			DpdProbeInterval: &dpdProfileProbeInterval,
			DpdProbeMode:     &dpdProfileProbeMode,
			Enabled:          &dpdProfileEnabled,
			RetryCount:       &dpdProfileRetryCount,
		}, nil)

		res := resourceNsxtPolicyIPSecVpnDpdProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecVpnDpdProfileData())
		d.SetId(dpdProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnDpdProfileUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnDpdProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIPSecVpnDpdProfileData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnDpdProfileUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IPSecVpnDpdProfile ID")
	})
}

func TestMockResourceNsxtPolicyIPSecVpnDpdProfileDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDpdSDK := dpdmocks.NewMockIpsecVpnDpdProfilesClient(ctrl)
	mockWrapper := &cliinfra.IPSecVpnDpdProfileClientContext{
		Client:     mockDpdSDK,
		ClientType: utl.Local,
	}

	originalCli := cliIpsecVpnDpdProfilesClient
	defer func() { cliIpsecVpnDpdProfilesClient = originalCli }()
	cliIpsecVpnDpdProfilesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.IPSecVpnDpdProfileClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockDpdSDK.EXPECT().Delete(dpdProfileID).Return(nil)

		res := resourceNsxtPolicyIPSecVpnDpdProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dpdProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnDpdProfileDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIPSecVpnDpdProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnDpdProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining IPSecVpnDpdProfile ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockDpdSDK.EXPECT().Delete(dpdProfileID).Return(errors.New("API error"))

		res := resourceNsxtPolicyIPSecVpnDpdProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dpdProfileID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyIPSecVpnDpdProfileDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalIPSecVpnDpdProfileData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":       dpdProfileDisplayName,
		"description":        dpdProfileDescription,
		"dpd_probe_interval": int(dpdProfileProbeInterval),
		"dpd_probe_mode":     dpdProfileProbeMode,
		"enabled":            dpdProfileEnabled,
		"retry_count":        int(dpdProfileRetryCount),
	}
}
