// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/DnsForwarderZonesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/DnsForwarderZonesClient.go DnsForwarderZonesClient

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
	dnsmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	dnsZoneID          = "dns-zone-1"
	dnsZoneDisplayName = "dns-zone-fooname"
	dnsZoneDescription = "dns zone mock"
	dnsZonePath        = "/infra/dns-forwarder-zones/dns-zone-1"
	dnsZoneRevision    = int64(1)
	dnsZoneUpstream    = []string{"192.168.1.1"}
)

func TestMockResourceNsxtPolicyDNSForwarderZoneCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDnsSDK := dnsmocks.NewMockDnsForwarderZonesClient(ctrl)
	mockWrapper := &cliinfra.PolicyDnsForwarderZoneClientContext{
		Client:     mockDnsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDnsForwarderZonesClient
	defer func() { cliDnsForwarderZonesClient = originalCli }()
	cliDnsForwarderZonesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.PolicyDnsForwarderZoneClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockDnsSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockDnsSDK.EXPECT().Get(gomock.Any()).Return(model.PolicyDnsForwarderZone{
			Id:              &dnsZoneID,
			DisplayName:     &dnsZoneDisplayName,
			Description:     &dnsZoneDescription,
			Path:            &dnsZonePath,
			Revision:        &dnsZoneRevision,
			UpstreamServers: dnsZoneUpstream,
		}, nil)

		res := resourceNsxtPolicyDNSForwarderZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDNSForwarderZoneData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDNSForwarderZoneCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockDnsSDK.EXPECT().Get("existing-id").Return(model.PolicyDnsForwarderZone{Id: &dnsZoneID}, nil)

		res := resourceNsxtPolicyDNSForwarderZone()
		data := minimalDNSForwarderZoneData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDNSForwarderZoneCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyDNSForwarderZoneRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDnsSDK := dnsmocks.NewMockDnsForwarderZonesClient(ctrl)
	mockWrapper := &cliinfra.PolicyDnsForwarderZoneClientContext{
		Client:     mockDnsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDnsForwarderZonesClient
	defer func() { cliDnsForwarderZonesClient = originalCli }()
	cliDnsForwarderZonesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.PolicyDnsForwarderZoneClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockDnsSDK.EXPECT().Get(dnsZoneID).Return(model.PolicyDnsForwarderZone{
			Id:              &dnsZoneID,
			DisplayName:     &dnsZoneDisplayName,
			Description:     &dnsZoneDescription,
			Path:            &dnsZonePath,
			Revision:        &dnsZoneRevision,
			UpstreamServers: dnsZoneUpstream,
		}, nil)

		res := resourceNsxtPolicyDNSForwarderZone()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dnsZoneID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDNSForwarderZoneRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, dnsZoneDisplayName, d.Get("display_name"))
		assert.Equal(t, dnsZoneDescription, d.Get("description"))
		assert.Equal(t, dnsZonePath, d.Get("path"))
		assert.Equal(t, int(dnsZoneRevision), d.Get("revision"))
		assert.Equal(t, dnsZoneID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDNSForwarderZone()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDNSForwarderZoneRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Dns Forwarder Zone ID")
	})
}

func TestMockResourceNsxtPolicyDNSForwarderZoneUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDnsSDK := dnsmocks.NewMockDnsForwarderZonesClient(ctrl)
	mockWrapper := &cliinfra.PolicyDnsForwarderZoneClientContext{
		Client:     mockDnsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDnsForwarderZonesClient
	defer func() { cliDnsForwarderZonesClient = originalCli }()
	cliDnsForwarderZonesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.PolicyDnsForwarderZoneClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockDnsSDK.EXPECT().Patch(dnsZoneID, gomock.Any()).Return(nil)
		mockDnsSDK.EXPECT().Get(dnsZoneID).Return(model.PolicyDnsForwarderZone{
			Id:              &dnsZoneID,
			DisplayName:     &dnsZoneDisplayName,
			Description:     &dnsZoneDescription,
			Path:            &dnsZonePath,
			Revision:        &dnsZoneRevision,
			UpstreamServers: dnsZoneUpstream,
		}, nil)

		res := resourceNsxtPolicyDNSForwarderZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDNSForwarderZoneData())
		d.SetId(dnsZoneID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDNSForwarderZoneUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDNSForwarderZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDNSForwarderZoneData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDNSForwarderZoneUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Dns Forwarder Zone ID")
	})
}

func TestMockResourceNsxtPolicyDNSForwarderZoneDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDnsSDK := dnsmocks.NewMockDnsForwarderZonesClient(ctrl)
	mockWrapper := &cliinfra.PolicyDnsForwarderZoneClientContext{
		Client:     mockDnsSDK,
		ClientType: utl.Local,
	}

	originalCli := cliDnsForwarderZonesClient
	defer func() { cliDnsForwarderZonesClient = originalCli }()
	cliDnsForwarderZonesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.PolicyDnsForwarderZoneClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockDnsSDK.EXPECT().Delete(dnsZoneID).Return(nil)

		res := resourceNsxtPolicyDNSForwarderZone()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dnsZoneID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDNSForwarderZoneDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyDNSForwarderZone()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDNSForwarderZoneDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining Dns Forwarder Zone ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockDnsSDK.EXPECT().Delete(dnsZoneID).Return(errors.New("API error"))

		res := resourceNsxtPolicyDNSForwarderZone()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(dnsZoneID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyDNSForwarderZoneDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalDNSForwarderZoneData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":     dnsZoneDisplayName,
		"description":      dnsZoneDescription,
		"upstream_servers": []interface{}{"192.168.1.1"},
	}
}
