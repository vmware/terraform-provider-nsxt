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
	dnsARCID          = "dns-arc-id"
	dnsARCDisplayName = "test-dns-arc"
	dnsARCDescription = "Test DNS Auto Record Config"
	dnsARCRevision    = int64(1)
)

func dnsARCAPIResponse() nsxModel.ProjectDnsAutoRecordConfig {
	ipBlock := "/infra/ip-blocks/block-1"
	zonePath := "/orgs/default/projects/p1/dns-services/svc1/zones/zone1"
	pattern := "{vm_name}"
	ttl := int64(300)
	return nsxModel.ProjectDnsAutoRecordConfig{
		Id:            &dnsARCID,
		DisplayName:   &dnsARCDisplayName,
		Description:   &dnsARCDescription,
		Revision:      &dnsARCRevision,
		IpBlockPath:   &ipBlock,
		ZonePath:      &zonePath,
		NamingPattern: &pattern,
		Ttl:           &ttl,
	}
}

func minimalDnsARCData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":   dnsARCDisplayName,
		"description":    dnsARCDescription,
		"nsx_id":         dnsARCID,
		"ip_block_path":  "/infra/ip-blocks/block-1",
		"zone_path":      "/orgs/default/projects/p1/dns-services/svc1/zones/zone1",
		"naming_pattern": "{vm_name}",
		"ttl":            300,
	}
}

func setupDnsARCMock(t *testing.T, ctrl *gomock.Controller) (*projectmocks.MockDnsAutoRecordConfigsClient, func()) {
	mockSDK := projectmocks.NewMockDnsAutoRecordConfigsClient(ctrl)
	mockWrapper := &apiprojects.ProjectDnsAutoRecordConfigClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  "project1",
	}

	original := cliProjectDnsAutoRecordConfigsClient
	cliProjectDnsAutoRecordConfigsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiprojects.ProjectDnsAutoRecordConfigClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliProjectDnsAutoRecordConfigsClient = original }
}

func TestMockResourceNsxtPolicyDnsRecordAutoConfigCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsARCMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsARCID).Return(nsxModel.ProjectDnsAutoRecordConfig{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), dnsARCID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsARCID).Return(dnsARCAPIResponse(), nil),
		)

		res := resourceNsxtPolicyDnsRecordAutoConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsARCData())

		err := resourceNsxtPolicyDnsRecordAutoConfigCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dnsARCID, d.Id())
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "9.1.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtPolicyDnsRecordAutoConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsARCData())

		err := resourceNsxtPolicyDnsRecordAutoConfigCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})
}

func TestMockResourceNsxtPolicyDnsRecordAutoConfigRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsARCMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsARCID).Return(dnsARCAPIResponse(), nil)

		res := resourceNsxtPolicyDnsRecordAutoConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsARCData())
		d.SetId(dnsARCID)

		err := resourceNsxtPolicyDnsRecordAutoConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dnsARCDisplayName, d.Get("display_name"))
		assert.Equal(t, "{vm_name}", d.Get("naming_pattern"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsARCMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsARCID).Return(nsxModel.ProjectDnsAutoRecordConfig{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyDnsRecordAutoConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsARCData())
		d.SetId(dnsARCID)

		err := resourceNsxtPolicyDnsRecordAutoConfigRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyDnsRecordAutoConfigUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsARCMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), dnsARCID, gomock.Any()).Return(dnsARCAPIResponse(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsARCID).Return(dnsARCAPIResponse(), nil),
		)

		res := resourceNsxtPolicyDnsRecordAutoConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsARCData())
		d.SetId(dnsARCID)

		err := resourceNsxtPolicyDnsRecordAutoConfigUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyDnsRecordAutoConfigDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsARCMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), dnsARCID).Return(nil)

		res := resourceNsxtPolicyDnsRecordAutoConfig()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsARCData())
		d.SetId(dnsARCID)

		err := resourceNsxtPolicyDnsRecordAutoConfigDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
