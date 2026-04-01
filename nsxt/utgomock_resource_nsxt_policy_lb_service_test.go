//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/infra/LbServicesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt/infra/LbServicesClient.go LbServicesClient

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
	lbSvcMocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	"github.com/vmware/terraform-provider-nsxt/nsxt/util"
)

var (
	lbServiceID           = "lb-service-1"
	lbServiceDisplayName  = "lb-service-fooname"
	lbServiceDescription  = "lb service mock"
	lbServicePath         = "/infra/lb-services/lb-service-1"
	lbServiceRevision     = int64(1)
	lbServiceConnectivity = "/infra/tier-1s/t1"
	lbServiceSize         = model.LBService_SIZE_SMALL
)

func TestMockResourceNsxtPolicyLBServiceCreate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLbSDK := lbSvcMocks.NewMockLbServicesClient(ctrl)
	mockWrapper := &cliinfra.LBServiceClientContext{
		Client:     mockLbSDK,
		ClientType: utl.Local,
	}

	originalCli := cliLbServicesClient
	defer func() { cliLbServicesClient = originalCli }()
	cliLbServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.LBServiceClientContext {
		return mockWrapper
	}

	t.Run("Create success", func(t *testing.T) {
		mockLbSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)
		mockLbSDK.EXPECT().Get(gomock.Any()).Return(model.LBService{
			Id:               &lbServiceID,
			DisplayName:      &lbServiceDisplayName,
			Description:      &lbServiceDescription,
			Path:             &lbServicePath,
			Revision:         &lbServiceRevision,
			ConnectivityPath: &lbServiceConnectivity,
			Size:             &lbServiceSize,
		}, nil)

		res := resourceNsxtPolicyLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBServiceData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBServiceCreate(d, m)
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
		assert.Equal(t, d.Id(), d.Get("nsx_id"))
	})

	t.Run("Create fails when resource already exists", func(t *testing.T) {
		mockLbSDK.EXPECT().Get("existing-id").Return(model.LBService{Id: &lbServiceID}, nil)

		res := resourceNsxtPolicyLBService()
		data := minimalLBServiceData()
		data["nsx_id"] = "existing-id"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBServiceCreate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestMockResourceNsxtPolicyLBServiceRead(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLbSDK := lbSvcMocks.NewMockLbServicesClient(ctrl)
	mockWrapper := &cliinfra.LBServiceClientContext{
		Client:     mockLbSDK,
		ClientType: utl.Local,
	}

	originalCli := cliLbServicesClient
	defer func() { cliLbServicesClient = originalCli }()
	cliLbServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.LBServiceClientContext {
		return mockWrapper
	}

	t.Run("Read success", func(t *testing.T) {
		mockLbSDK.EXPECT().Get(lbServiceID).Return(model.LBService{
			Id:               &lbServiceID,
			DisplayName:      &lbServiceDisplayName,
			Description:      &lbServiceDescription,
			Path:             &lbServicePath,
			Revision:         &lbServiceRevision,
			ConnectivityPath: &lbServiceConnectivity,
			Size:             &lbServiceSize,
		}, nil)

		res := resourceNsxtPolicyLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(lbServiceID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBServiceRead(d, m)
		require.NoError(t, err)
		assert.Equal(t, lbServiceDisplayName, d.Get("display_name"))
		assert.Equal(t, lbServiceDescription, d.Get("description"))
		assert.Equal(t, lbServicePath, d.Get("path"))
		assert.Equal(t, int(lbServiceRevision), d.Get("revision"))
		assert.Equal(t, lbServiceID, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBServiceRead(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining LBService ID")
	})
}

func TestMockResourceNsxtPolicyLBServiceUpdate(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLbSDK := lbSvcMocks.NewMockLbServicesClient(ctrl)
	mockWrapper := &cliinfra.LBServiceClientContext{
		Client:     mockLbSDK,
		ClientType: utl.Local,
	}

	originalCli := cliLbServicesClient
	defer func() { cliLbServicesClient = originalCli }()
	cliLbServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.LBServiceClientContext {
		return mockWrapper
	}

	t.Run("Update success", func(t *testing.T) {
		mockLbSDK.EXPECT().Patch(lbServiceID, gomock.Any()).Return(nil)
		mockLbSDK.EXPECT().Get(lbServiceID).Return(model.LBService{
			Id:               &lbServiceID,
			DisplayName:      &lbServiceDisplayName,
			Description:      &lbServiceDescription,
			Path:             &lbServicePath,
			Revision:         &lbServiceRevision,
			ConnectivityPath: &lbServiceConnectivity,
			Size:             &lbServiceSize,
		}, nil)

		res := resourceNsxtPolicyLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBServiceData())
		d.SetId(lbServiceID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBServiceUpdate(d, m)
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalLBServiceData())

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBServiceUpdate(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining LBService ID")
	})
}

func TestMockResourceNsxtPolicyLBServiceDelete(t *testing.T) {
	util.NsxVersion = "9.1.0"
	defer func() { util.NsxVersion = "" }()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLbSDK := lbSvcMocks.NewMockLbServicesClient(ctrl)
	mockWrapper := &cliinfra.LBServiceClientContext{
		Client:     mockLbSDK,
		ClientType: utl.Local,
	}

	originalCli := cliLbServicesClient
	defer func() { cliLbServicesClient = originalCli }()
	cliLbServicesClient = func(sessionContext utl.SessionContext, connector client.Connector) *cliinfra.LBServiceClientContext {
		return mockWrapper
	}

	t.Run("Delete success", func(t *testing.T) {
		mockLbSDK.EXPECT().Delete(lbServiceID, gomock.Any()).Return(nil)

		res := resourceNsxtPolicyLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(lbServiceID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBServiceDelete(d, m)
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBServiceDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error obtaining LBService ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		mockLbSDK.EXPECT().Delete(lbServiceID, gomock.Any()).Return(errors.New("API error"))

		res := resourceNsxtPolicyLBService()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(lbServiceID)

		m := newGoMockProviderClient()
		err := resourceNsxtPolicyLBServiceDelete(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
	})
}

func minimalLBServiceData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":      lbServiceDisplayName,
		"description":       lbServiceDescription,
		"connectivity_path": lbServiceConnectivity,
	}
}
