// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/DhcpRelayConfigsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/DhcpRelayConfigsClient.go DhcpRelayConfigsClient

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
	relaymocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	dhcpRelayID          = "dhcp-relay-1"
	dhcpRelayDisplayName = "dhcp-relay-fooname"
	dhcpRelayDescription = "dhcp relay mock"
	dhcpRelayPath        = "/infra/dhcp-relays/dhcp-relay-1"
	dhcpRelayRevision    = int64(1)
	dhcpRelayServers     = []string{"192.168.1.1"}
)

func TestMockResourceNsxtPolicyDhcpRelayConfigCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRelaySDK := relaymocks.NewMockDhcpRelayConfigsClient(ctrl)
	mockWrapper := &cliinfra.DhcpRelayConfigClientContext{
		Client:     mockRelaySDK,
		ClientType: utl.Local,
	}

	originalCli := cliDhcpRelayConfigsClient
	defer func() { cliDhcpRelayConfigsClient = originalCli }()
	cliDhcpRelayConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DhcpRelayConfigClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockRelaySDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockRelaySDK.EXPECT().Get(gomock.Any()).Return(model.DhcpRelayConfig{
			Id:              &dhcpRelayID,
			DisplayName:     &dhcpRelayDisplayName,
			Description:     &dhcpRelayDescription,
			Path:            &dhcpRelayPath,
			Revision:        &dhcpRelayRevision,
			ServerAddresses: dhcpRelayServers,
		}, nil)

		res := resourceNsxtPolicyDhcpRelayConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpRelayData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpRelayConfigCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockRelaySDK.EXPECT().Get("existing-id").Return(model.DhcpRelayConfig{Id: &dhcpRelayID}, nil)

		res := resourceNsxtPolicyDhcpRelayConfig()
		data := minimalDhcpRelayData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpRelayConfigCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyDhcpRelayConfigRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRelaySDK := relaymocks.NewMockDhcpRelayConfigsClient(ctrl)
	mockWrapper := &cliinfra.DhcpRelayConfigClientContext{
		Client:     mockRelaySDK,
		ClientType: utl.Local,
	}

	originalCli := cliDhcpRelayConfigsClient
	defer func() { cliDhcpRelayConfigsClient = originalCli }()
	cliDhcpRelayConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DhcpRelayConfigClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockRelaySDK.EXPECT().Get(dhcpRelayID).Return(model.DhcpRelayConfig{
			Id:              &dhcpRelayID,
			DisplayName:     &dhcpRelayDisplayName,
			Description:     &dhcpRelayDescription,
			Path:            &dhcpRelayPath,
			Revision:        &dhcpRelayRevision,
			ServerAddresses: dhcpRelayServers,
		}, nil)

		res := resourceNsxtPolicyDhcpRelayConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dhcpRelayID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpRelayConfigRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, dhcpRelayDisplayName, d.Get("display_name"))
		assert.Equal(t, dhcpRelayDescription, d.Get("description"))
		assert.Equal(t, dhcpRelayPath, d.Get("path"))
		assert.Equal(t, int(dhcpRelayRevision), d.Get("revision"))
		assert.Equal(t, dhcpRelayID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDhcpRelayConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpRelayConfigRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining DhcpRelayConfig ID")
	})
}

func TestMockResourceNsxtPolicyDhcpRelayConfigUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRelaySDK := relaymocks.NewMockDhcpRelayConfigsClient(ctrl)
	mockWrapper := &cliinfra.DhcpRelayConfigClientContext{
		Client:     mockRelaySDK,
		ClientType: utl.Local,
	}

	originalCli := cliDhcpRelayConfigsClient
	defer func() { cliDhcpRelayConfigsClient = originalCli }()
	cliDhcpRelayConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DhcpRelayConfigClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockRelaySDK.EXPECT().Update(dhcpRelayID, gomock.Any()).Return(model.DhcpRelayConfig{
			Id:              &dhcpRelayID,
			DisplayName:     &dhcpRelayDisplayName,
			Description:     &dhcpRelayDescription,
			Path:            &dhcpRelayPath,
			Revision:        &dhcpRelayRevision,
			ServerAddresses: dhcpRelayServers,
		}, nil)
		mockRelaySDK.EXPECT().Get(dhcpRelayID).Return(model.DhcpRelayConfig{
			Id:              &dhcpRelayID,
			DisplayName:     &dhcpRelayDisplayName,
			Description:     &dhcpRelayDescription,
			Path:            &dhcpRelayPath,
			Revision:        &dhcpRelayRevision,
			ServerAddresses: dhcpRelayServers,
		}, nil)

		res := resourceNsxtPolicyDhcpRelayConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpRelayData())
		d.Set("revision", 1)
		d.SetId(dhcpRelayID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpRelayConfigUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDhcpRelayConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDhcpRelayData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpRelayConfigUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining DhcpRelayConfig ID")
	})
}

func TestMockResourceNsxtPolicyDhcpRelayConfigDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRelaySDK := relaymocks.NewMockDhcpRelayConfigsClient(ctrl)
	mockWrapper := &cliinfra.DhcpRelayConfigClientContext{
		Client:     mockRelaySDK,
		ClientType: utl.Local,
	}

	originalCli := cliDhcpRelayConfigsClient
	defer func() { cliDhcpRelayConfigsClient = originalCli }()
	cliDhcpRelayConfigsClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.DhcpRelayConfigClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockRelaySDK.EXPECT().Delete(dhcpRelayID).Return(nil)

		res := resourceNsxtPolicyDhcpRelayConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dhcpRelayID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpRelayConfigDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDhcpRelayConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpRelayConfigDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining DhcpRelayConfig ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockRelaySDK.EXPECT().Delete(dhcpRelayID).Return(errors.New("API error"))

		res := resourceNsxtPolicyDhcpRelayConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dhcpRelayID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDhcpRelayConfigDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalDhcpRelayData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":     dhcpRelayDisplayName,
		"description":      dhcpRelayDescription,
		"server_addresses": []interface{}{"192.168.1.1"},
	}
}
