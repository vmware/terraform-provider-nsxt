//go:build unittest

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

	apiprojects "github.com/vmware/terraform-provider-nsxt/api/orgs/projects"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	projectmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	vspID          = "service-profile-id"
	vspDisplayName = "test-service-profile"
	vspDescription = "Test VPC Service Profile"
	vspRevision    = int64(1)
)

func vspAPIResponse() nsxModel.VpcServiceProfile {
	return nsxModel.VpcServiceProfile{
		Id:          &vspID,
		DisplayName: &vspDisplayName,
		Description: &vspDescription,
		Revision:    &vspRevision,
	}
}

func vspAPIResponse920() nsxModel.VpcServiceProfile {
	leaseV6 := int64(86400)
	prefV6 := int64(48000)
	cache := int64(1024)
	logLevel := nsxModel.PolicyVpcDnsForwarder_LOG_LEVEL_INFO
	dist := false
	return nsxModel.VpcServiceProfile{
		Id:          &vspID,
		DisplayName: &vspDisplayName,
		Description: &vspDescription,
		Revision:    &vspRevision,
		Dhcpv6Config: &nsxModel.VpcProfileDhcpV6Config{
			Dhcpv6ServerConfig: &nsxModel.VpcDhcpv6ServerConfig{
				LeaseTime:     &leaseV6,
				PreferredTime: &prefV6,
				NtpServers:    []string{"2001:db8:1::1"},
				SntpServers:   []string{"2001:db8:1::2"},
				DnsClientConfig: &nsxModel.DnsClientConfig{
					DnsServerIps: []string{"2001:db8:1::53"},
				},
				AdvancedConfig: &nsxModel.VpcDhcpAdvancedConfig{
					IsDistributedDhcp: &dist,
				},
			},
		},
		DnsForwarderConfig: &nsxModel.PolicyVpcDnsForwarder{
			CacheSize: &cache,
			LogLevel:  &logLevel,
		},
		Ipv6ProfilePaths:   []string{"/orgs/default/projects/project1/infra/ipv6-dad-profiles/dad1"},
		ServiceSubnetCidrs: []string{"fd12:3456:789a::/64"},
	}
}

func minimalVspData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": vspDisplayName,
		"description":  vspDescription,
		"nsx_id":       vspID,
		"dhcp_config": []interface{}{
			map[string]interface{}{},
		},
	}
}

func minimalVspData920Ipv6() map[string]interface{} {
	data := minimalVspData()
	data["dhcpv6_config"] = []interface{}{map[string]interface{}{
		"dhcpv6_server_config": []interface{}{map[string]interface{}{
			"lease_time":     86400,
			"preferred_time": 48000,
			"ntp_servers":    []interface{}{"2001:db8:1::1"},
			"sntp_servers":   []interface{}{"2001:db8:1::2"},
			"dns_client_config": []interface{}{map[string]interface{}{
				"dns_server_ips": []interface{}{"2001:db8:1::53"},
			}},
			"advanced_config": []interface{}{map[string]interface{}{
				"is_distributed_dhcp": false,
			}},
		}},
	}}
	data["dns_forwarder_config"] = []interface{}{map[string]interface{}{
		"log_level":  "INFO",
		"cache_size": 1024,
	}}
	data["ipv6_profile_paths"] = []interface{}{"/orgs/default/projects/project1/infra/ipv6-dad-profiles/dad1"}
	data["service_subnet_cidrs"] = []interface{}{"fd12:3456:789a::/64"}
	return data
}

func setupVspMock(t *testing.T, ctrl *gomock.Controller) (*projectmocks.MockVpcServiceProfilesClient, func()) {
	mockSDK := projectmocks.NewMockVpcServiceProfilesClient(ctrl)
	mockWrapper := &apiprojects.VpcServiceProfileClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  "project1",
	}

	original := cliVpcServiceProfilesClient
	cliVpcServiceProfilesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiprojects.VpcServiceProfileClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliVpcServiceProfilesClient = original }
}

func TestMockResourceNsxtVpcServiceProfileCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.0.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVspMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vspID).Return(nsxModel.VpcServiceProfile{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), vspID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vspID).Return(vspAPIResponse(), nil),
		)

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData())

		err := resourceNsxtVpcServiceProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vspID, d.Id())
	})

	t.Run("Create success with IPv6 fields on NSX 9.2", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVspMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vspID).Return(nsxModel.VpcServiceProfile{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), vspID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vspID).Return(vspAPIResponse920(), nil),
		)

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData920Ipv6())

		err := resourceNsxtVpcServiceProfileCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vspID, d.Id())
		assert.Equal(t, 1, d.Get("dhcpv6_config.#"))
		assert.Equal(t, 1, d.Get("service_subnet_cidrs.#"))
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "4.0.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData())

		err := resourceNsxtVpcServiceProfileCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtVpcServiceProfileRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVspMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vspID).Return(vspAPIResponse(), nil)

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData())
		d.SetId(vspID)

		err := resourceNsxtVpcServiceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vspDisplayName, d.Get("display_name"))
	})

	t.Run("Read maps IPv6 and DNS forwarder fields from API", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVspMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vspID).Return(vspAPIResponse920(), nil)

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData920Ipv6())
		d.SetId(vspID)

		err := resourceNsxtVpcServiceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, vspDisplayName, d.Get("display_name"))
		assert.Equal(t, 1, d.Get("dhcpv6_config.#"))
		assert.Equal(t, 86400, d.Get("dhcpv6_config.0.dhcpv6_server_config.0.lease_time"))
		assert.Equal(t, 48000, d.Get("dhcpv6_config.0.dhcpv6_server_config.0.preferred_time"))
		assert.Equal(t, "2001:db8:1::1", d.Get("dhcpv6_config.0.dhcpv6_server_config.0.ntp_servers.0"))
		assert.Equal(t, "2001:db8:1::2", d.Get("dhcpv6_config.0.dhcpv6_server_config.0.sntp_servers.0"))
		assert.Equal(t, "2001:db8:1::53", d.Get("dhcpv6_config.0.dhcpv6_server_config.0.dns_client_config.0.dns_server_ips.0"))
		assert.Equal(t, false, d.Get("dhcpv6_config.0.dhcpv6_server_config.0.advanced_config.0.is_distributed_dhcp"))
		assert.Equal(t, 1, d.Get("dns_forwarder_config.#"))
		assert.Equal(t, "INFO", d.Get("dns_forwarder_config.0.log_level"))
		assert.Equal(t, 1024, d.Get("dns_forwarder_config.0.cache_size"))
		assert.Equal(t, "/orgs/default/projects/project1/infra/ipv6-dad-profiles/dad1", d.Get("ipv6_profile_paths.0"))
		assert.Equal(t, "fd12:3456:789a::/64", d.Get("service_subnet_cidrs.0"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVspMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vspID).Return(nsxModel.VpcServiceProfile{}, vapiErrors.NotFound{})

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData())
		d.SetId(vspID)

		err := resourceNsxtVpcServiceProfileRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtVpcServiceProfileUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVspMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), vspID, gomock.Any()).Return(vspAPIResponse(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vspID).Return(vspAPIResponse(), nil),
		)

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData())
		d.SetId(vspID)

		err := resourceNsxtVpcServiceProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update success with IPv6 fields on NSX 9.2", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVspMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), vspID, gomock.Any()).Return(vspAPIResponse920(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), vspID).Return(vspAPIResponse920(), nil),
		)

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData920Ipv6())
		d.SetId(vspID)
		d.Set("revision", 1)

		err := resourceNsxtVpcServiceProfileUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, 1, d.Get("dhcpv6_config.#"))
	})
}

func TestMockResourceNsxtVpcServiceProfileDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupVspMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), vspID).Return(nil)

		res := resourceNsxtVpcServiceProfile()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalVspData())
		d.SetId(vspID)

		err := resourceNsxtVpcServiceProfileDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
