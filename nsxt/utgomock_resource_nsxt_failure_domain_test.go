//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/nsx/FailureDomainsClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/FailureDomainsClient.go FailureDomainsClient

package nsxt

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"

	apinsx "github.com/vmware/terraform-provider-nsxt/api/nsx"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	nsxmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx"
)

var (
	fdID          = "failure-domain-uuid-1"
	fdDisplayName = "test-failure-domain"
	fdDescription = "Test failure domain"
	fdRevision    = int64(1)
)

func failureDomainAPIResponse() nsxModel.FailureDomain {
	activePrefer := true
	return nsxModel.FailureDomain{
		Id:                          &fdID,
		DisplayName:                 &fdDisplayName,
		Description:                 &fdDescription,
		Revision:                    &fdRevision,
		PreferredActiveEdgeServices: &activePrefer,
	}
}

func failureDomainAPIResponseNoPreference() nsxModel.FailureDomain {
	resp := failureDomainAPIResponse()
	resp.PreferredActiveEdgeServices = nil
	return resp
}

func setupFailureDomainMock(t *testing.T, ctrl *gomock.Controller) (*nsxmocks.MockFailureDomainsClient, func()) {
	mockSDK := nsxmocks.NewMockFailureDomainsClient(ctrl)
	mockWrapper := &apinsx.FailureDomainClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalCli := cliFailureDomainsClient
	cliFailureDomainsClient = func(_ utl.SessionContext, _ client.Connector) *apinsx.FailureDomainClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliFailureDomainsClient = originalCli }
}

func minimalFailureDomainData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":            fdDisplayName,
		"description":             fdDescription,
		"preferred_edge_services": "active",
	}
}

func TestMockResourceNsxtFailureDomainCreate(t *testing.T) {
	t.Run("Create success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupFailureDomainMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Create(gomock.Any()).Return(failureDomainAPIResponse(), nil)
		mockSDK.EXPECT().Get(fdID).Return(failureDomainAPIResponse(), nil)

		res := resourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFailureDomainData())

		err := resourceNsxtFailureDomainCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, fdID, d.Id())
		assert.Equal(t, fdDisplayName, d.Get("display_name"))
		assert.Equal(t, "active", d.Get("preferred_edge_services"))
	})

	t.Run("Create success with no_preference", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupFailureDomainMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Create(gomock.Any()).Return(failureDomainAPIResponseNoPreference(), nil)
		mockSDK.EXPECT().Get(fdID).Return(failureDomainAPIResponseNoPreference(), nil)

		res := resourceNsxtFailureDomain()
		data := minimalFailureDomainData()
		data["preferred_edge_services"] = "no_preference"
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		err := resourceNsxtFailureDomainCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "no_preference", d.Get("preferred_edge_services"))
	})

	t.Run("Create fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupFailureDomainMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Create(gomock.Any()).Return(nsxModel.FailureDomain{}, errors.New("create API error"))

		res := resourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFailureDomainData())

		err := resourceNsxtFailureDomainCreate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "create API error")
	})
}

func TestMockResourceNsxtFailureDomainRead(t *testing.T) {
	t.Run("Read success with active preference", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupFailureDomainMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(fdID).Return(failureDomainAPIResponse(), nil)

		res := resourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(fdID)

		err := resourceNsxtFailureDomainRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, fdDisplayName, d.Get("display_name"))
		assert.Equal(t, fdDescription, d.Get("description"))
		assert.Equal(t, int(fdRevision), d.Get("revision"))
		assert.Equal(t, "active", d.Get("preferred_edge_services"))
	})

	t.Run("Read success with standby preference", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupFailureDomainMock(t, ctrl)
		defer restore()

		standby := false
		resp := failureDomainAPIResponse()
		resp.PreferredActiveEdgeServices = &standby
		mockSDK.EXPECT().Get(fdID).Return(resp, nil)

		res := resourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(fdID)

		err := resourceNsxtFailureDomainRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "standby", d.Get("preferred_edge_services"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtFailureDomainRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FailureDomain ID")
	})

	t.Run("Read fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupFailureDomainMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Get(fdID).Return(nsxModel.FailureDomain{}, errors.New("read API error"))

		res := resourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(fdID)

		err := resourceNsxtFailureDomainRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "read API error")
	})
}

func TestMockResourceNsxtFailureDomainUpdate(t *testing.T) {
	t.Run("Update success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupFailureDomainMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Update(fdID, gomock.Any()).Return(failureDomainAPIResponse(), nil)
		mockSDK.EXPECT().Get(fdID).Return(failureDomainAPIResponse(), nil)

		res := resourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFailureDomainData())
		d.SetId(fdID)

		err := resourceNsxtFailureDomainUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, fdDisplayName, d.Get("display_name"))
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFailureDomainData())

		err := resourceNsxtFailureDomainUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FailureDomain ID")
	})

	t.Run("Update fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupFailureDomainMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Update(fdID, gomock.Any()).Return(nsxModel.FailureDomain{}, errors.New("update API error"))

		res := resourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalFailureDomainData())
		d.SetId(fdID)

		err := resourceNsxtFailureDomainUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "update API error")
	})
}

func TestMockResourceNsxtFailureDomainDelete(t *testing.T) {
	t.Run("Delete success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupFailureDomainMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(fdID).Return(nil)

		res := resourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(fdID)

		err := resourceNsxtFailureDomainDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

		err := resourceNsxtFailureDomainDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FailureDomain ID")
	})

	t.Run("Delete fails when API returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, restore := setupFailureDomainMock(t, ctrl)
		defer restore()

		mockSDK.EXPECT().Delete(fdID).Return(errors.New("delete API error"))

		res := resourceNsxtFailureDomain()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(fdID)

		err := resourceNsxtFailureDomainDelete(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete API error")
	})
}
