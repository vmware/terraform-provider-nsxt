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

	apidomains "github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	domainmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains"
)

func setupGatewayPolicyDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockGatewayPoliciesClient, func()) {
	t.Helper()
	mockSDK := domainmocks.NewMockGatewayPoliciesClient(ctrl)
	mockWrapper := &apidomains.GatewayPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliGatewayPoliciesClient
	cliGatewayPoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.GatewayPolicyClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliGatewayPoliciesClient = orig }
}

func TestMockDataSourceNsxtPolicyGatewayPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGatewayPolicyDataSourceMock(t, ctrl)
	defer restore()

	boolFalse := false

	t.Run("by id success", func(t *testing.T) {
		mockSDK.EXPECT().Get(gwPolicyDomain, gwPolicyID).Return(gwPolicyAPIResponse(), nil)

		ds := dataSourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     gwPolicyID,
			"domain": gwPolicyDomain,
		})

		err := dataSourceNsxtPolicyGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, gwPolicyID, d.Id())
		assert.Equal(t, gwPolicyDisplayName, d.Get("display_name"))
		assert.Equal(t, gwPolicyCategory, d.Get("category"))
	})

	t.Run("by id not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(gwPolicyDomain, gwPolicyID).Return(nsxModel.GatewayPolicy{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     gwPolicyID,
			"domain": gwPolicyDomain,
		})

		err := dataSourceNsxtPolicyGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by id API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(gwPolicyDomain, gwPolicyID).Return(nsxModel.GatewayPolicy{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     gwPolicyID,
			"domain": gwPolicyDomain,
		})

		err := dataSourceNsxtPolicyGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading Gateway Policy")
	})

	t.Run("missing id display_name and category", func(t *testing.T) {
		ds := dataSourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain": gwPolicyDomain,
		})

		err := dataSourceNsxtPolicyGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be specified")
	})

	t.Run("by display_name single match via list", func(t *testing.T) {
		mockSDK.EXPECT().List(gwPolicyDomain, nil, nil, nil, nil, nil, &boolFalse, nil).Return(nsxModel.GatewayPolicyListResult{
			Results: []nsxModel.GatewayPolicy{gwPolicyAPIResponse()},
		}, nil)

		ds := dataSourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": gwPolicyDisplayName,
			"domain":       gwPolicyDomain,
		})

		err := dataSourceNsxtPolicyGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, gwPolicyID, d.Id())
	})

	t.Run("list API error", func(t *testing.T) {
		mockSDK.EXPECT().List(gwPolicyDomain, nil, nil, nil, nil, nil, &boolFalse, nil).Return(nsxModel.GatewayPolicyListResult{}, vapiErrors.ServiceUnavailable{})

		ds := dataSourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": gwPolicyDisplayName,
			"domain":       gwPolicyDomain,
		})

		err := dataSourceNsxtPolicyGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading Gateway Policies")
	})

	t.Run("by display_name multiple exact matches", func(t *testing.T) {
		mockSDK.EXPECT().List(gwPolicyDomain, nil, nil, nil, nil, nil, &boolFalse, nil).Return(nsxModel.GatewayPolicyListResult{
			Results: []nsxModel.GatewayPolicy{gwPolicyAPIResponse(), gwPolicyAPIResponse()},
		}, nil)

		ds := dataSourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": gwPolicyDisplayName,
			"domain":       gwPolicyDomain,
		})

		err := dataSourceNsxtPolicyGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("by category single match without display_name", func(t *testing.T) {
		p := gwPolicyAPIResponse()
		mockSDK.EXPECT().List(gwPolicyDomain, nil, nil, nil, nil, nil, &boolFalse, nil).Return(nsxModel.GatewayPolicyListResult{
			Results: []nsxModel.GatewayPolicy{p},
		}, nil)

		ds := dataSourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":   gwPolicyDomain,
			"category": gwPolicyCategory,
		})

		err := dataSourceNsxtPolicyGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, gwPolicyID, d.Id())
	})

	t.Run("list returns nil client", func(t *testing.T) {
		orig := cliGatewayPoliciesClient
		cliGatewayPoliciesClient = func(utl.SessionContext, vapiProtocolClient.Connector) *apidomains.GatewayPolicyClientContext {
			return nil
		}
		defer func() { cliGatewayPoliciesClient = orig }()

		ds := dataSourceNsxtPolicyGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": gwPolicyDisplayName,
			"domain":       gwPolicyDomain,
		})

		err := dataSourceNsxtPolicyGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
