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

func projectDnsZoneAPIResponse() nsxModel.ProjectDnsZone {
	domain := "example.com"
	ttl := int64(300)
	return nsxModel.ProjectDnsZone{
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
	mockWrapper := &dnssvcsapi.ProjectDnsZoneClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  "project1",
	}

	original := cliProjectDnsZonesClient
	cliProjectDnsZonesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *dnssvcsapi.ProjectDnsZoneClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliProjectDnsZonesClient = original }
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
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID).Return(nsxModel.ProjectDnsZone{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID).Return(projectDnsZoneAPIResponse(), nil),
		)

		res := resourceNsxtPolicyDnsZone()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalProjectDnsZoneData())

		err := resourceNsxtPolicyDnsZoneCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, projectDnsZoneID, d.Id())
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

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsZoneMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), projectDnsZoneID).Return(nsxModel.ProjectDnsZone{}, vapiErrors.NotFound{})

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
