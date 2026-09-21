//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	dnssvcsapi "github.com/vmware/terraform-provider-nsxt/api/orgs/projects/dns_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	dnssvcsmocks "github.com/vmware/terraform-provider-nsxt/mocks/orgs/projects/dns_services"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	projectDnsZoneID          = "dns-zone-id"
	projectDnsZoneDisplayName = "test-dns-zone"
	projectDnsZoneDescription = "Test Project DNS Zone"
	projectDnsZoneRevision    = int64(1)
	projectDnsZoneParentPath  = "/orgs/default/projects/project1/dns-services/dns-svc-1"
)

func projectDnsZoneAPIResponse() nsxModel.DnsZone {
	domain := "example.com"
	ttl := int64(300)
	return nsxModel.DnsZone{
		Id:            &projectDnsZoneID,
		DisplayName:   &projectDnsZoneDisplayName,
		Description:   &projectDnsZoneDescription,
		Revision:      &projectDnsZoneRevision,
		DnsDomainName: &domain,
		Ttl:           &ttl,
	}
}

func minimalProjectDnsZoneData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    projectDnsZoneDisplayName,
		"description":     projectDnsZoneDescription,
		"nsx_id":          projectDnsZoneID,
		"parent_path":     projectDnsZoneParentPath,
		"dns_domain_name": "example.com",
		"ttl":             300,
	}
}

func setupDnsZoneMock(t *testing.T, ctrl *gomock.Controller) (*dnssvcsmocks.MockZonesClient, func()) {
	mockSDK := dnssvcsmocks.NewMockZonesClient(ctrl)
	mockWrapper := &dnssvcsapi.DnsZoneClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  "project1",
	}

	original := cliDnsZonesClient
	cliDnsZonesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *dnssvcsapi.DnsZoneClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliDnsZonesClient = original }
}

func TestMockResourceNsxtPolicyDnsZoneCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsZoneMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID).Return(nsxModel.DnsZone{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID).Return(projectDnsZoneAPIResponse(), nil),
		)

		res := resourceNsxtPolicyDnsZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectDnsZoneData())

		err := resourceNsxtPolicyDnsZoneCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, projectDnsZoneID, d.Id())
	})

	t.Run("Create success with resolution_scope", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsZoneMock(t, ctrl)
		defer restore()

		scopeData := minimalProjectDnsZoneData()
		scopeData["resolution_scope"] = []interface{}{
			"/orgs/default/projects/project1/vpcs/vpc-1",
			"/orgs/default/projects/project1/vpcs/vpc-2",
		}

		notFoundErr := vapiErrors.NotFound{}
		var capturedObj nsxModel.DnsZone
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID).Return(nsxModel.DnsZone{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID, gomock.Any()).DoAndReturn(
				func(_ string, _ string, _ string, _ string, obj nsxModel.DnsZone) error {
					capturedObj = obj
					return nil
				}),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID).Return(projectDnsZoneAPIResponse(), nil),
		)

		res := resourceNsxtPolicyDnsZone()
		assert.Equal(t, schema.TypeSet, res.Schema["resolution_scope"].Type)

		d := schema.TestResourceDataRaw(t, res.Schema, scopeData)

		err := resourceNsxtPolicyDnsZoneCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, projectDnsZoneID, d.Id())
		assert.ElementsMatch(t, []string{
			"/orgs/default/projects/project1/vpcs/vpc-1",
			"/orgs/default/projects/project1/vpcs/vpc-2",
		}, capturedObj.ResolutionScope)
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "9.1.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtPolicyDnsZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectDnsZoneData())

		err := resourceNsxtPolicyDnsZoneCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})
}

func TestMockResourceNsxtPolicyDnsZoneRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsZoneMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID).Return(projectDnsZoneAPIResponse(), nil)

		res := resourceNsxtPolicyDnsZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectDnsZoneData())
		d.SetId(projectDnsZoneID)

		err := resourceNsxtPolicyDnsZoneRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, projectDnsZoneDisplayName, d.Get("display_name"))
	})

	t.Run("Read success with resolution_scope", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsZoneMock(t, ctrl)
		defer restore()

		apiResp := projectDnsZoneAPIResponse()
		apiResp.ResolutionScope = []string{
			"/orgs/default/projects/project1/vpcs/vpc-1",
			"/orgs/default/projects/project1/vpcs/vpc-2",
		}
		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID).Return(apiResp, nil)

		res := resourceNsxtPolicyDnsZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectDnsZoneData())
		d.SetId(projectDnsZoneID)

		err := resourceNsxtPolicyDnsZoneRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, projectDnsZoneDisplayName, d.Get("display_name"))
		resScopeSet := d.Get("resolution_scope").(*schema.Set)
		assert.Equal(t, 2, resScopeSet.Len())
		assert.True(t, resScopeSet.Contains("/orgs/default/projects/project1/vpcs/vpc-1"))
		assert.True(t, resScopeSet.Contains("/orgs/default/projects/project1/vpcs/vpc-2"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsZoneMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID).Return(nsxModel.DnsZone{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyDnsZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectDnsZoneData())
		d.SetId(projectDnsZoneID)

		err := resourceNsxtPolicyDnsZoneRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyDnsZoneUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsZoneMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID, gomock.Any()).Return(projectDnsZoneAPIResponse(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID).Return(projectDnsZoneAPIResponse(), nil),
		)

		res := resourceNsxtPolicyDnsZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectDnsZoneData())
		d.SetId(projectDnsZoneID)

		err := resourceNsxtPolicyDnsZoneUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyDnsZoneDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsZoneMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID).Return(nil)

		res := resourceNsxtPolicyDnsZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectDnsZoneData())
		d.SetId(projectDnsZoneID)

		err := resourceNsxtPolicyDnsZoneDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestUnitNsxt_policyDnsZoneFromSchema(t *testing.T) {
	res := resourceNsxtPolicyDnsZone()

	t.Run("basic fields and empty soa are converted", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":     "zone1",
			"description":      "a zone",
			"dns_domain_name":  "example.com",
			"ttl":              600,
			"resolution_scope": []interface{}{"/orgs/default/projects/proj1/vpcs/vpc1"},
		})
		obj := policyDnsZoneFromSchema(d)
		assert.Equal(t, "zone1", *obj.DisplayName)
		assert.Equal(t, "a zone", *obj.Description)
		assert.Equal(t, "example.com", *obj.DnsDomainName)
		assert.EqualValues(t, 600, *obj.Ttl)
		assert.ElementsMatch(t, []string{"/orgs/default/projects/proj1/vpcs/vpc1"}, obj.ResolutionScope)
		assert.Nil(t, obj.Soa, "all-zero-value soa block should not be sent")
	})

	t.Run("soa fields are converted when set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":    "zone1",
			"dns_domain_name": "example.com",
			"soa": []interface{}{
				map[string]interface{}{
					"primary_nameserver": "ns1.example.com.",
					"responsible_party":  "admin.example.com.",
					"serial_number":      5,
					"refresh_interval":   3600,
					"retry_interval":     600,
					"expire_time":        86400,
					"negative_cache_ttl": 60,
				},
			},
		})
		obj := policyDnsZoneFromSchema(d)
		require.NotNil(t, obj.Soa)
		assert.Equal(t, "ns1.example.com.", *obj.Soa.PrimaryNameserver)
		assert.Equal(t, "admin.example.com.", *obj.Soa.ResponsibleParty)
		assert.EqualValues(t, 5, *obj.Soa.SerialNumber)
		assert.EqualValues(t, 3600, *obj.Soa.RefreshInterval)
		assert.EqualValues(t, 600, *obj.Soa.RetryInterval)
		assert.EqualValues(t, 86400, *obj.Soa.ExpireTime)
		assert.EqualValues(t, 60, *obj.Soa.NegativeCacheTtl)
	})

	t.Run("soa with only negative_cache_ttl zero is still sent", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"display_name":    "zone1",
			"dns_domain_name": "example.com",
			"soa": []interface{}{
				map[string]interface{}{
					"primary_nameserver": "ns1.example.com.",
				},
			},
		})
		obj := policyDnsZoneFromSchema(d)
		require.NotNil(t, obj.Soa)
		assert.Equal(t, "ns1.example.com.", *obj.Soa.PrimaryNameserver)
	})
}

func TestUnitNsxt_policyDnsZoneToSchema(t *testing.T) {
	res := resourceNsxtPolicyDnsZone()

	t.Run("populates schema fields including soa", func(t *testing.T) {
		displayName := "zone1"
		description := "a zone"
		dnsDomain := "example.com"
		ttl := int64(600)
		primaryNS := "ns1.example.com."
		serial := int64(5)

		obj := nsxModel.DnsZone{
			DisplayName:     &displayName,
			Description:     &description,
			DnsDomainName:   &dnsDomain,
			Ttl:             &ttl,
			ResolutionScope: []string{"/orgs/default/projects/proj1/vpcs/vpc1"},
			Soa: &nsxModel.DnsZoneSoa{
				PrimaryNameserver: &primaryNS,
				SerialNumber:      &serial,
			},
		}

		d := res.TestResourceData()
		policyDnsZoneToSchema(d, obj)

		assert.Equal(t, displayName, d.Get("display_name"))
		assert.Equal(t, description, d.Get("description"))
		assert.Equal(t, dnsDomain, d.Get("dns_domain_name"))
		assert.Equal(t, 600, d.Get("ttl"))
		soa := d.Get("soa").([]interface{})
		require.Len(t, soa, 1)
		soaMap := soa[0].(map[string]interface{})
		assert.Equal(t, primaryNS, soaMap["primary_nameserver"])
		assert.Equal(t, 5, soaMap["serial_number"])
	})

	t.Run("nil soa leaves schema soa unset", func(t *testing.T) {
		displayName := "zone1"
		obj := nsxModel.DnsZone{DisplayName: &displayName}

		d := res.TestResourceData()
		policyDnsZoneToSchema(d, obj)

		assert.Empty(t, d.Get("soa").([]interface{}))
	})
}
