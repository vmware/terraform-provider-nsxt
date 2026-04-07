// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	t0dnsapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s"
	t1dnsapi "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	t0dnsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_0s"
	t1dnsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/tier_1s"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	dnsT0GwID          = "t0-dns-gw-1"
	dnsT0GwPath        = "/infra/tier-0s/t0-dns-gw-1"
	dnsT1GwID          = "t1-dns-gw-1"
	dnsT1GwPath        = "/infra/tier-1s/t1-dns-gw-1"
	dnsDisplayName     = "Test DNS Forwarder"
	dnsDescription     = "Test DNS forwarder"
	dnsRevision        = int64(1)
	dnsListenerIP      = "10.10.10.1"
	dnsDefaultZonePath = "/infra/dns-forwarder-zones/default-zone"
	dnsPath            = "/infra/tier-0s/t0-dns-gw-1/dns-forwarder"
)

func dnsForwarderAPIResponse() nsxModel.PolicyDnsForwarder {
	enabled := true
	logLevel := nsxModel.PolicyDnsForwarder_LOG_LEVEL_INFO
	cacheSize := int64(1024)
	return nsxModel.PolicyDnsForwarder{
		DisplayName:              &dnsDisplayName,
		Description:              &dnsDescription,
		Revision:                 &dnsRevision,
		Path:                     &dnsPath,
		ListenerIp:               &dnsListenerIP,
		DefaultForwarderZonePath: &dnsDefaultZonePath,
		Enabled:                  &enabled,
		LogLevel:                 &logLevel,
		CacheSize:                &cacheSize,
	}
}

func minimalDnsForwarderData(gwPath string) map[string]interface{} {
	return map[string]interface{}{
		"display_name":                dnsDisplayName,
		"description":                 dnsDescription,
		"gateway_path":                gwPath,
		"listener_ip":                 dnsListenerIP,
		"default_forwarder_zone_path": dnsDefaultZonePath,
		"enabled":                     true,
		"log_level":                   nsxModel.PolicyDnsForwarder_LOG_LEVEL_INFO,
		"cache_size":                  1024,
	}
}

func setupT0DnsForwarderMock(t *testing.T, ctrl *gomock.Controller) (*t0dnsmocks.MockDnsForwarderClient, func()) {
	mockSDK := t0dnsmocks.NewMockDnsForwarderClient(ctrl)
	mockWrapper := &t0dnsapi.PolicyDnsForwarderClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliTier0DnsForwarderClient
	cliTier0DnsForwarderClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *t0dnsapi.PolicyDnsForwarderClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier0DnsForwarderClient = original }
}

func setupT1DnsForwarderMock(t *testing.T, ctrl *gomock.Controller) (*t1dnsmocks.MockDnsForwarderClient, func()) {
	mockSDK := t1dnsmocks.NewMockDnsForwarderClient(ctrl)
	mockWrapper := &t1dnsapi.PolicyDnsForwarderClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliTier1DnsForwarderClient
	cliTier1DnsForwarderClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *t1dnsapi.PolicyDnsForwarderClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliTier1DnsForwarderClient = original }
}

func TestMockResourceNsxtPolicyGatewayDNSForwarderCreate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockT0SDK, restoreT0 := setupT0DnsForwarderMock(t, ctrl)
	defer restoreT0()

	t.Run("Create success for Tier0 gateway", func(t *testing.T) {
		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockT0SDK.EXPECT().Get(dnsT0GwID).Return(nsxModel.PolicyDnsForwarder{}, notFoundErr),
			mockT0SDK.EXPECT().Patch(dnsT0GwID, gomock.Any()).Return(nil),
			mockT0SDK.EXPECT().Get(dnsT0GwID).Return(dnsForwarderAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGatewayDNSForwarder()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsForwarderData(dnsT0GwPath))

		err := resourceNsxtPolicyGatewayDNSForwarderCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dnsT0GwID, d.Id())
	})

	t.Run("Create fails when forwarder already exists", func(t *testing.T) {
		mockT0SDK.EXPECT().Get(dnsT0GwID).Return(dnsForwarderAPIResponse(), nil)

		res := resourceNsxtPolicyGatewayDNSForwarder()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsForwarderData(dnsT0GwPath))

		err := resourceNsxtPolicyGatewayDNSForwarderCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayDNSForwarderRead(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Read success for Tier0 gateway", func(t *testing.T) {
		mockT0SDK, restoreT0 := setupT0DnsForwarderMock(t, ctrl)
		defer restoreT0()

		mockT0SDK.EXPECT().Get(dnsT0GwID).Return(dnsForwarderAPIResponse(), nil)

		res := resourceNsxtPolicyGatewayDNSForwarder()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsForwarderData(dnsT0GwPath))
		d.SetId(dnsT0GwID)

		err := resourceNsxtPolicyGatewayDNSForwarderRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dnsListenerIP, d.Get("listener_ip"))
	})

	t.Run("Read success for Tier1 gateway", func(t *testing.T) {
		mockT1SDK, restoreT1 := setupT1DnsForwarderMock(t, ctrl)
		defer restoreT1()

		mockT1SDK.EXPECT().Get(dnsT1GwID).Return(dnsForwarderAPIResponse(), nil)

		res := resourceNsxtPolicyGatewayDNSForwarder()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsForwarderData(dnsT1GwPath))
		d.SetId(dnsT1GwID)

		err := resourceNsxtPolicyGatewayDNSForwarderRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dnsListenerIP, d.Get("listener_ip"))
	})
}

func TestMockResourceNsxtPolicyGatewayDNSForwarderUpdate(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockT0SDK, restoreT0 := setupT0DnsForwarderMock(t, ctrl)
	defer restoreT0()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockT0SDK.EXPECT().Patch(dnsT0GwID, gomock.Any()).Return(nil),
			mockT0SDK.EXPECT().Get(dnsT0GwID).Return(dnsForwarderAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGatewayDNSForwarder()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsForwarderData(dnsT0GwPath))
		d.SetId(dnsT0GwID)

		err := resourceNsxtPolicyGatewayDNSForwarderUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayDNSForwarderDelete(t *testing.T) {
	util.NsxVersion = "3.2.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockT0SDK, restoreT0 := setupT0DnsForwarderMock(t, ctrl)
	defer restoreT0()

	t.Run("Delete success for Tier0", func(t *testing.T) {
		mockT0SDK.EXPECT().Delete(dnsT0GwID).Return(nil)

		res := resourceNsxtPolicyGatewayDNSForwarder()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsForwarderData(dnsT0GwPath))
		d.SetId(dnsT0GwID)

		err := resourceNsxtPolicyGatewayDNSForwarderDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
