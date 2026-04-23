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

func setupSecurityPolicyDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockSecurityPoliciesClient, func()) {
	t.Helper()
	mockSDK := domainmocks.NewMockSecurityPoliciesClient(ctrl)
	mockWrapper := &apidomains.SecurityPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliSecurityPoliciesClient
	cliSecurityPoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.SecurityPolicyClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliSecurityPoliciesClient = orig }
}

func TestMockDataSourceNsxtPolicySecurityPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupSecurityPolicyDataSourceMock(t, ctrl)
	defer restore()

	boolFalse := false

	t.Run("by id success", func(t *testing.T) {
		mockSDK.EXPECT().Get(secPolicyDomain, secPolicyID).Return(secPolicyAPIResponse(), nil)

		ds := dataSourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     secPolicyID,
			"domain": secPolicyDomain,
		})

		err := dataSourceNsxtPolicySecurityPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, secPolicyID, d.Id())
		assert.Equal(t, secPolicyDisplayName, d.Get("display_name"))
		assert.Equal(t, secPolicyCategory, d.Get("category"))
	})

	t.Run("by id not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(secPolicyDomain, secPolicyID).Return(nsxModel.SecurityPolicy{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     secPolicyID,
			"domain": secPolicyDomain,
		})

		err := dataSourceNsxtPolicySecurityPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by id API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(secPolicyDomain, secPolicyID).Return(nsxModel.SecurityPolicy{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     secPolicyID,
			"domain": secPolicyDomain,
		})

		err := dataSourceNsxtPolicySecurityPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading Security Policy")
	})

	t.Run("missing id display_name and category", func(t *testing.T) {
		ds := dataSourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain": secPolicyDomain,
		})

		err := dataSourceNsxtPolicySecurityPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be specified")
	})

	t.Run("by display_name single match via list", func(t *testing.T) {
		mockSDK.EXPECT().List(secPolicyDomain, nil, nil, nil, nil, nil, &boolFalse, nil).Return(nsxModel.SecurityPolicyListResult{
			Results: []nsxModel.SecurityPolicy{secPolicyAPIResponse()},
		}, nil)

		ds := dataSourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": secPolicyDisplayName,
			"domain":       secPolicyDomain,
		})

		err := dataSourceNsxtPolicySecurityPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, secPolicyID, d.Id())
	})

	t.Run("list API error", func(t *testing.T) {
		mockSDK.EXPECT().List(secPolicyDomain, nil, nil, nil, nil, nil, &boolFalse, nil).Return(nsxModel.SecurityPolicyListResult{}, vapiErrors.ServiceUnavailable{})

		ds := dataSourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": secPolicyDisplayName,
			"domain":       secPolicyDomain,
		})

		err := dataSourceNsxtPolicySecurityPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading Security Policies")
	})

	t.Run("by display_name multiple exact matches", func(t *testing.T) {
		mockSDK.EXPECT().List(secPolicyDomain, nil, nil, nil, nil, nil, &boolFalse, nil).Return(nsxModel.SecurityPolicyListResult{
			Results: []nsxModel.SecurityPolicy{secPolicyAPIResponse(), secPolicyAPIResponse()},
		}, nil)

		ds := dataSourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": secPolicyDisplayName,
			"domain":       secPolicyDomain,
		})

		err := dataSourceNsxtPolicySecurityPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "multiple")
	})

	t.Run("by category single match without display_name", func(t *testing.T) {
		p := secPolicyAPIResponse()
		mockSDK.EXPECT().List(secPolicyDomain, nil, nil, nil, nil, nil, &boolFalse, nil).Return(nsxModel.SecurityPolicyListResult{
			Results: []nsxModel.SecurityPolicy{p},
		}, nil)

		ds := dataSourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":   secPolicyDomain,
			"category": secPolicyCategory,
		})

		err := dataSourceNsxtPolicySecurityPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, secPolicyID, d.Id())
	})

	t.Run("list returns nil client", func(t *testing.T) {
		orig := cliSecurityPoliciesClient
		cliSecurityPoliciesClient = func(utl.SessionContext, vapiProtocolClient.Connector) *apidomains.SecurityPolicyClientContext {
			return nil
		}
		defer func() { cliSecurityPoliciesClient = orig }()

		ds := dataSourceNsxtPolicySecurityPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"display_name": secPolicyDisplayName,
			"domain":       secPolicyDomain,
		})

		err := dataSourceNsxtPolicySecurityPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
