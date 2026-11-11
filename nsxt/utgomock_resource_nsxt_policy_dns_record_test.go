//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

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

func dnsRecordAPIResponse() nsxModel.DnsRecord {
	recordName := "www"
	recordType := nsxModel.DnsRecord_RECORD_TYPE_A
	zonePath := "/orgs/default/projects/p1/dns-services/svc1/zones/zone1"
	ttl := int64(300)
	fqdn := "www.example.com."
	return nsxModel.DnsRecord{
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
	mockWrapper := &apiprojects.DnsRecordClientContext{
		Client:     mockSDK,
		ClientType: utl.Multitenancy,
		ProjectID:  "project1",
	}

	original := cliProjectDnsRecordsClient
	cliProjectDnsRecordsClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apiprojects.DnsRecordClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliProjectDnsRecordsClient = original }
}

func TestUnitNsxt_PolicyDnsRecordTypeValidation(t *testing.T) {
	validate := resourceNsxtPolicyDnsRecord().Schema["record_type"].ValidateFunc
	for _, tc := range []struct {
		name  string
		value string
		valid bool
	}{
		{name: "A", value: "A", valid: true},
		{name: "AAAA", value: "AAAA", valid: true},
		{name: "CNAME", value: "CNAME", valid: true},
		{name: "PTR", value: "PTR", valid: true},
		{name: "NS", value: "NS", valid: true},
		{name: "TXT", value: "TXT", valid: true},
		{name: "unsupported", value: "INVALID"},
		{name: "lowercase", value: "txt"},
		{name: "empty"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, errs := validate(tc.value, "record_type")
			if tc.valid {
				require.Empty(t, errs)
			} else {
				require.NotEmpty(t, errs)
			}
		})
	}
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
			mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsRecordID).Return(nsxModel.DnsRecord{}, notFoundErr),
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

		mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsRecordID).Return(nsxModel.DnsRecord{}, vapiErrors.NotFound{})

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

func TestDnsRecordValueOrderDiff(t *testing.T) {
	for _, recordType := range []string{"A", "AAAA", "NS"} {
		t.Run(recordType, func(t *testing.T) {
			values := []interface{}{"192.168.1.10", "192.168.1.20", "192.168.1.30"}
			extra := "192.168.1.40"
			if recordType == "AAAA" {
				values = []interface{}{"2001:db8::10", "2001:db8::20", "2001:db8::30"}
				extra = "2001:db8::40"
			} else if recordType == "NS" {
				values = []interface{}{"ns1.example.com", "ns2.example.com", "ns3.example.com"}
				extra = "ns4.example.com"
			}
			for _, tc := range []struct {
				name       string
				values     []interface{}
				wantChange bool
			}{
				{"unchanged", values, false},
				{"reversed", []interface{}{values[2], values[1], values[0]}, false},
				{"rotated", []interface{}{values[1], values[2], values[0]}, false},
				{"added", append(append([]interface{}{}, values...), extra), true},
				{"duplicate added", append(append([]interface{}{}, values...), values[0]), false},
				{"removed", values[:2], true},
				{"replaced", []interface{}{values[0], values[1], values[1]}, true},
			} {
				t.Run(tc.name, func(t *testing.T) {
					res := resourceNsxtPolicyDnsRecord()
					stored := minimalDnsRecordData()
					stored["record_type"] = recordType
					stored["record_values"] = values
					d := schema.TestResourceDataRaw(t, res.Schema, stored)
					d.SetId(dnsRecordID)
					configured := minimalDnsRecordData()
					configured["record_type"] = recordType
					configured["record_values"] = tc.values
					diff, err := res.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(configured), nil)
					require.NoError(t, err)
					assert.Equal(t, tc.wantChange, diff != nil && !diff.Empty())
					actualSet, ok := d.Get("record_values").(*schema.Set)
					require.True(t, ok)
					assert.ElementsMatch(t, values, actualSet.List(), "comparison must not mutate existing state")
					assert.ElementsMatch(t, getStringListFromSchemaSet(d, "record_values"), policyDnsRecordFromSchema(d).RecordValues)
				})
			}
		})
	}
}

func TestDnsRecordRefreshLegacyListState(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSDK, restore := setupDnsRecordMock(t, ctrl)
	t.Cleanup(restore)
	legacy := resourceNsxtPolicyDnsRecord()
	legacyValues := *legacy.Schema["record_values"]
	legacyValues.Type = schema.TypeList
	legacy.Schema["record_values"] = &legacyValues
	configured := minimalDnsRecordData()
	configured["record_values"] = []interface{}{"192.168.1.10", "192.168.1.20"}
	old := schema.TestResourceDataRaw(t, legacy.Schema, configured)
	old.SetId(dnsRecordID)
	res := resourceNsxtPolicyDnsRecord()
	d := res.Data(old.State())
	response := dnsRecordAPIResponse()
	response.RecordValues = []string{"192.168.1.20", "192.168.1.10"}
	mockSDK.EXPECT().Get(gomock.Any(), gomock.Any(), dnsRecordID).Return(response, nil)
	require.NoError(t, res.Read(d, newGoMockProviderClient()))
	assert.Equal(t, dnsRecordID, d.Id())
	values, ok := d.Get("record_values").(*schema.Set)
	require.True(t, ok)
	assert.ElementsMatch(t, configured["record_values"], values.List())
	assert.ElementsMatch(t, response.RecordValues, policyDnsRecordFromSchema(d).RecordValues)
	diff, err := res.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(configured), nil)
	require.NoError(t, err)
	assert.True(t, diff == nil || diff.Empty(), "refreshed legacy state must have a clean plan")
}
