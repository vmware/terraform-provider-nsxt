// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/ServicesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/ServicesClient.go ServicesClient

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
