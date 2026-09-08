//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/ServicesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/ServicesClient.go ServicesClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	cliinfra "github.com/vmware/terraform-provider-nsxt/api/infra"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	svcmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	svcDisplayName = "service-fooname"
	svcDescription = "service mock description"
	svcPath        = "/infra/services/my-service"
	svcRevision    = int64(1)
	svcID          = "my-service"
)

func TestMockResourceNsxtPolicyServiceCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServicesSDK := svcmocks.NewMockServicesClient(ctrl)
	mockWrapper := &cliinfra.ServiceClientContext{
		Client:     mockServicesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliServicesClient
	defer func() { cliServicesClient = originalCli }()
	cliServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ServiceClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockServicesSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockServicesSDK.EXPECT().Get(gomock.Any()).Return(model.Service{
			DisplayName: &svcDisplayName,
			Description: &svcDescription,
			Path:        &svcPath,
			Revision:    &svcRevision,
		}, nil)

		res := resourceNsxtPolicyService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalServiceData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyServiceCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockServicesSDK.EXPECT().Get("existing-id").Return(model.Service{Id: &svcID}, nil)

		res := resourceNsxtPolicyService()
		data := minimalServiceData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyServiceCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyServiceRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServicesSDK := svcmocks.NewMockServicesClient(ctrl)
	mockWrapper := &cliinfra.ServiceClientContext{
		Client:     mockServicesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliServicesClient
	defer func() { cliServicesClient = originalCli }()
	cliServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ServiceClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockServicesSDK.EXPECT().Get(svcID).Return(model.Service{
			DisplayName: &svcDisplayName,
			Description: &svcDescription,
			Path:        &svcPath,
			Revision:    &svcRevision,
		}, nil)

		res := resourceNsxtPolicyService()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(svcID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyServiceRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, svcDisplayName, d.Get("display_name"))
		assert.Equal(t, svcDescription, d.Get("description"))
		assert.Equal(t, svcPath, d.Get("path"))
		assert.Equal(t, int(svcRevision), d.Get("revision"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyService()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyServiceRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining service id")
	})
}

func TestMockResourceNsxtPolicyServiceUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServicesSDK := svcmocks.NewMockServicesClient(ctrl)
	mockWrapper := &cliinfra.ServiceClientContext{
		Client:     mockServicesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliServicesClient
	defer func() { cliServicesClient = originalCli }()
	cliServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ServiceClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockServicesSDK.EXPECT().Update(svcID, gomock.Any()).Return(model.Service{
			DisplayName: &svcDisplayName,
			Description: &svcDescription,
			Path:        &svcPath,
			Revision:    &svcRevision,
		}, nil)
		mockServicesSDK.EXPECT().Get(svcID).Return(model.Service{
			DisplayName: &svcDisplayName,
			Description: &svcDescription,
			Path:        &svcPath,
			Revision:    &svcRevision,
		}, nil)

		res := resourceNsxtPolicyService()
		data := minimalServiceData()
		data["revision"] = 1
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(svcID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyServiceUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalServiceData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyServiceUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining service id")
	})
}

func TestMockResourceNsxtPolicyServiceDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServicesSDK := svcmocks.NewMockServicesClient(ctrl)
	mockWrapper := &cliinfra.ServiceClientContext{
		Client:     mockServicesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliServicesClient
	defer func() { cliServicesClient = originalCli }()
	cliServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ServiceClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockServicesSDK.EXPECT().Delete(svcID).Return(nil)

		res := resourceNsxtPolicyService()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(svcID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyServiceDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyService()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyServiceDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining service id")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockServicesSDK.EXPECT().Delete(svcID).Return(errors.New("API error"))

		res := resourceNsxtPolicyService()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(svcID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyServiceDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalServiceData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":         svcDisplayName,
		"description":          svcDescription,
		"icmp_entry":           []interface{}{},
		"l4_port_set_entry":    []interface{}{},
		"igmp_entry":           []interface{}{},
		"ether_type_entry":     []interface{}{},
		"ip_protocol_entry":    []interface{}{},
		"algorithm_entry":      []interface{}{},
		"nested_service_entry": []interface{}{},
	}
}

func allServiceEntryTypesData() map[string]interface{} {
	data := minimalServiceData()
	data["icmp_entry"] = []interface{}{
		map[string]interface{}{
			"display_name": "icmp-1",
			"description":  "icmp desc",
			"protocol":     "ICMPv4",
			"icmp_type":    "8",
			"icmp_code":    "0",
		},
	}
	data["l4_port_set_entry"] = []interface{}{
		map[string]interface{}{
			"display_name":      "l4-1",
			"description":       "l4 desc",
			"protocol":          "TCP",
			"destination_ports": []interface{}{"80"},
			"source_ports":      []interface{}{"1024-2048"},
		},
	}
	data["igmp_entry"] = []interface{}{
		map[string]interface{}{
			"display_name": "igmp-1",
			"description":  "igmp desc",
		},
	}
	data["ether_type_entry"] = []interface{}{
		map[string]interface{}{
			"display_name": "ether-1",
			"description":  "ether desc",
			"ether_type":   2048,
		},
	}
	data["ip_protocol_entry"] = []interface{}{
		map[string]interface{}{
			"display_name": "ipprot-1",
			"description":  "ipprot desc",
			"protocol":     6,
		},
	}
	data["algorithm_entry"] = []interface{}{
		map[string]interface{}{
			"display_name":     "alg-1",
			"description":      "alg desc",
			"destination_port": "21",
			"source_ports":     []interface{}{"1024-2048"},
			"algorithm":        "FTP",
		},
	}
	data["nested_service_entry"] = []interface{}{
		map[string]interface{}{
			"display_name":        "nested-1",
			"description":         "nested desc",
			"nested_service_path": "/infra/services/other-service",
		},
	}
	return data
}

func TestUnitNsxt_getServiceEntriesFromSchema(t *testing.T) {
	t.Run("converts one entry of every type", func(t *testing.T) {
		res := resourceNsxtPolicyService()
		d := schema.TestResourceDataRaw(t, res.Schema, allServiceEntryTypesData())

		entries, err := getServiceEntriesFromSchema(d)
		require.NoError(t, err)
		assert.Len(t, entries, 7)
	})

	t.Run("fails on invalid icmp_type", func(t *testing.T) {
		res := resourceNsxtPolicyService()
		data := minimalServiceData()
		data["icmp_entry"] = []interface{}{
			map[string]interface{}{"protocol": "ICMPv4", "icmp_type": "not-a-number", "icmp_code": "", "display_name": "", "description": ""},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		_, err := getServiceEntriesFromSchema(d)
		require.Error(t, err)
	})

	t.Run("fails on invalid icmp_code", func(t *testing.T) {
		res := resourceNsxtPolicyService()
		data := minimalServiceData()
		data["icmp_entry"] = []interface{}{
			map[string]interface{}{"protocol": "ICMPv4", "icmp_type": "", "icmp_code": "not-a-number", "display_name": "", "description": ""},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		_, err := getServiceEntriesFromSchema(d)
		require.Error(t, err)
	})

	t.Run("works against a nested map instead of ResourceData", func(t *testing.T) {
		nestedMap := map[string]interface{}{
			"icmp_entry": schema.NewSet(schema.HashResource(getIcmpEntrySchema().Elem.(*schema.Resource)), nil),
		}
		entries, err := getServiceEntriesFromSchema(nestedMap)
		require.NoError(t, err)
		assert.Empty(t, entries)
	})
}

func TestUnitNsxt_setServiceEntriesInSchema(t *testing.T) {
	t.Run("round-trips every entry type back into schema", func(t *testing.T) {
		res := resourceNsxtPolicyService()
		d := schema.TestResourceDataRaw(t, res.Schema, allServiceEntryTypesData())

		entries, err := getServiceEntriesFromSchema(d)
		require.NoError(t, err)

		emptyData := schema.TestResourceDataRaw(t, res.Schema, minimalServiceData())
		err = setServiceEntriesInSchema(emptyData, entries, true)
		require.NoError(t, err)

		assert.Len(t, emptyData.Get("icmp_entry").(*schema.Set).List(), 1)
		assert.Len(t, emptyData.Get("l4_port_set_entry").(*schema.Set).List(), 1)
		assert.Len(t, emptyData.Get("igmp_entry").(*schema.Set).List(), 1)
		assert.Len(t, emptyData.Get("ether_type_entry").(*schema.Set).List(), 1)
		assert.Len(t, emptyData.Get("ip_protocol_entry").(*schema.Set).List(), 1)
		assert.Len(t, emptyData.Get("algorithm_entry").(*schema.Set).List(), 1)
		assert.Len(t, emptyData.Get("nested_service_entry").(*schema.Set).List(), 1)
	})

	t.Run("skips nested_service_entry when nestedSupported is false", func(t *testing.T) {
		res := resourceNsxtPolicyService()
		d := schema.TestResourceDataRaw(t, res.Schema, allServiceEntryTypesData())

		entries, err := getServiceEntriesFromSchema(d)
		require.NoError(t, err)

		emptyData := schema.TestResourceDataRaw(t, res.Schema, minimalServiceData())
		err = setServiceEntriesInSchema(emptyData, entries, false)
		require.NoError(t, err)
		assert.Empty(t, emptyData.Get("nested_service_entry").(*schema.Set).List())
	})
}

func TestUnitNsxt_filterServiceEntryDisplayName(t *testing.T) {
	t.Run("returns empty string when display name is nil", func(t *testing.T) {
		assert.Equal(t, "", filterServiceEntryDisplayName(nil, nil))
	})

	t.Run("returns empty string when display name matches the generated id", func(t *testing.T) {
		name := "auto-id-123"
		assert.Equal(t, "", filterServiceEntryDisplayName(&name, &name))
	})

	t.Run("returns the display name when it differs from the id", func(t *testing.T) {
		name, id := "custom-name", "auto-id-123"
		assert.Equal(t, "custom-name", filterServiceEntryDisplayName(&name, &id))
	})
}

func TestMockResourceNsxtPolicyServiceExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServicesSDK := svcmocks.NewMockServicesClient(ctrl)
	mockWrapper := &cliinfra.ServiceClientContext{
		Client:     mockServicesSDK,
		ClientType: utl.Local,
	}

	originalCli := cliServicesClient
	defer func() { cliServicesClient = originalCli }()
	cliServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.ServiceClientContext {
		return mockWrapper
	}

	t.Run("returns true when the service exists", func(t *testing.T) {
		mockServicesSDK.EXPECT().Get(svcID).Return(model.Service{Id: &svcID}, nil)

		exists, err := resourceNsxtPolicyServiceExists(utl.SessionContext{}, svcID, nil)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("returns false when the service is not found", func(t *testing.T) {
		mockServicesSDK.EXPECT().Get("missing-id").Return(model.Service{}, vapiErrors.NotFound{})

		exists, err := resourceNsxtPolicyServiceExists(utl.SessionContext{}, "missing-id", nil)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("propagates other API errors", func(t *testing.T) {
		mockServicesSDK.EXPECT().Get("err-id").Return(model.Service{}, errors.New("boom"))

		_, err := resourceNsxtPolicyServiceExists(utl.SessionContext{}, "err-id", nil)
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyServiceCreateEntryConversionError(t *testing.T) {
	res := resourceNsxtPolicyService()
	data := minimalServiceData()
	data["icmp_entry"] = []interface{}{
		map[string]interface{}{"protocol": "ICMPv4", "icmp_type": "not-a-number", "icmp_code": "", "display_name": "", "description": ""},
	}
	d := schema.TestResourceDataRaw(t, res.Schema, data)

	m := newGoMockProviderClient()
	err := resourceNsxtPolicyServiceCreate(d, m)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Service entries conversion")
}
