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
	dnsRecordID          = "dns-record-id"
	dnsRecordDisplayName = "test-dns-record"
	dnsRecordDescription = "Test Project DNS Record"
	dnsRecordRevision    = int64(1)
)

func dnsRecordAPIResponse() nsxModel.ProjectDnsRecord {
	recordName := "www"
	recordType := nsxModel.ProjectDnsRecord_RECORD_TYPE_A
	zonePath := "/orgs/default/projects/p1/dns-services/svc1/zones/zone1"
	ttl := int64(300)
	fqdn := "www.example.com."
	return nsxModel.ProjectDnsRecord{
		Id:           &dnsRecordID,
		DisplayName:  &dnsRecordDisplayName,
		Description:  &dnsRecordDescription,
		Revision:     &dnsRecordRevision,
		RecordName:   &recordName,
		RecordType:   &recordType,
		RecordValues: []string{"192.168.1.10"},
		ZonePath:     &zonePath,
		Ttl:          &ttl,
		Fqdn:         &fqdn,
	}
}

func minimalDnsRecordData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":  dnsRecordDisplayName,
		"description":   dnsRecordDescription,
		"nsx_id":        dnsRecordID,
		"record_name":   "www",
		"record_type":   "A",
		"record_values": []interface{}{"192.168.1.10"},
		"zone_path":     "/orgs/default/projects/p1/dns-services/svc1/zones/zone1",
		"ttl":           300,
	}
}

func setupDnsRecordMock(t *testing.T, ctrl *gomock.Controller) (*projectmocks.MockDnsRecordsClient, func()) {
	mockSDK := projectmocks.NewMockDnsRecordsClient(ctrl)
	mockWrapper := &apiprojects.ProjectDnsRecordClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  "project1",
	}

	original := cliProjectDnsRecordsClient
	cliProjectDnsRecordsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiprojects.ProjectDnsRecordClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliProjectDnsRecordsClient = original }
}

func TestMockResourceNsxtPolicyDnsRecordCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		util.NsxVersion = "9.2.0"
		defer func() { util.NsxVersion = "" }()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsRecordMock(t, ctrl)
		defer restore()

		notFoundErr := vapiErrors.NotFound{}
		gomock.InOrder(
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsRecordID).Return(nsxModel.ProjectDnsRecord{}, notFoundErr),
			mockSDK.EXPECT().Patch(gomock.Any(), gomock.Any(), dnsRecordID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsRecordID).Return(dnsRecordAPIResponse(), nil),
		)

		res := resourceNsxtPolicyDnsRecord()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsRecordData())

		err := resourceNsxtPolicyDnsRecordCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dnsRecordID, d.Id())
	})

	t.Run("Create fails on old NSX version", func(t *testing.T) {
		util.NsxVersion = "9.1.0"
		defer func() { util.NsxVersion = "" }()

		res := resourceNsxtPolicyDnsRecord()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsRecordData())

		err := resourceNsxtPolicyDnsRecordCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "9.2.0")
	})
}

func TestMockResourceNsxtPolicyDnsRecordRead(t *testing.T) {
	t.Run("Read success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsRecordMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsRecordID).Return(dnsRecordAPIResponse(), nil)

		res := resourceNsxtPolicyDnsRecord()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsRecordData())
		d.SetId(dnsRecordID)

		err := resourceNsxtPolicyDnsRecordRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, dnsRecordDisplayName, d.Get("display_name"))
		assert.Equal(t, "www.example.com.", d.Get("fqdn"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsRecordMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsRecordID).Return(nsxModel.ProjectDnsRecord{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyDnsRecord()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsRecordData())
		d.SetId(dnsRecordID)

		err := resourceNsxtPolicyDnsRecordRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyDnsRecordUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsRecordMock(t, ctrl)
		defer restore()

		gomock.InOrder(
			mockSDK.EXPECT().Update(gomock.Any(), gomock.Any(), dnsRecordID, gomock.Any()).Return(dnsRecordAPIResponse(), nil),
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsRecordID).Return(dnsRecordAPIResponse(), nil),
		)

		res := resourceNsxtPolicyDnsRecord()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsRecordData())
		d.SetId(dnsRecordID)

		err := resourceNsxtPolicyDnsRecordUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestMockResourceNsxtPolicyDnsRecordDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupDnsRecordMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(gomock.Any(), gomock.Any(), dnsRecordID).Return(nil)

		res := resourceNsxtPolicyDnsRecord()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalDnsRecordData())
		d.SetId(dnsRecordID)

		err := resourceNsxtPolicyDnsRecordDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}
