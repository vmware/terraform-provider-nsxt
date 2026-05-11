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

var (
	idsParentDfwPolDomain   = "default"
	idsParentDfwPolID       = "ids-parent-dfw-pol-1"
	idsParentDfwPolName     = "ids-parent-dfw-policy"
	idsParentDfwPolPath     = "/infra/domains/default/intrusion-service-policies/ids-parent-dfw-pol-1"
	idsParentDfwPolCategory = "ThreatRules"
	idsParentDfwPolDesc     = "parent policy desc"
)

func idsParentDfwPolicyModel() nsxModel.IdsSecurityPolicy {
	stateful := true
	locked := false
	sequenceNumber := int64(10)
	return nsxModel.IdsSecurityPolicy{
		Id:             &idsParentDfwPolID,
		DisplayName:    &idsParentDfwPolName,
		Description:    &idsParentDfwPolDesc,
		Path:           &idsParentDfwPolPath,
		Category:       &idsParentDfwPolCategory,
		Stateful:       &stateful,
		Locked:         &locked,
		SequenceNumber: &sequenceNumber,
	}
}

func setupParentIntrusionServiceDfwPolicyDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockIntrusionServicePoliciesClient, func()) {
	t.Helper()
	mockSDK := domainmocks.NewMockIntrusionServicePoliciesClient(ctrl)
	wrapper := &apidomains.IdsSecurityPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliIntrusionServicePoliciesClient
	cliIntrusionServicePoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.IdsSecurityPolicyClientContext {
		return wrapper
	}
	return mockSDK, func() { cliIntrusionServicePoliciesClient = orig }
}

func TestMockDataSourceNsxtPolicyParentIntrusionServicePolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupParentIntrusionServiceDfwPolicyDataSourceMock(t, ctrl)
	defer restore()

	boolFalse := false

	t.Run("by id success", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsParentDfwPolDomain, idsParentDfwPolID).Return(idsParentDfwPolicyModel(), nil)

		ds := dataSourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     idsParentDfwPolID,
			"domain": idsParentDfwPolDomain,
		})

		err := dataSourceNsxtPolicyParentIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsParentDfwPolID, d.Id())
		assert.Equal(t, idsParentDfwPolName, d.Get("display_name"))
		assert.Equal(t, idsParentDfwPolCategory, d.Get("category"))
		assert.Equal(t, true, d.Get("stateful"))
		assert.Equal(t, false, d.Get("locked"))
		assert.Equal(t, 10, d.Get("sequence_number"))
	})

	t.Run("by id not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsParentDfwPolDomain, idsParentDfwPolID).Return(nsxModel.IdsSecurityPolicy{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     idsParentDfwPolID,
			"domain": idsParentDfwPolDomain,
		})

		err := dataSourceNsxtPolicyParentIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by id API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsParentDfwPolDomain, idsParentDfwPolID).Return(nsxModel.IdsSecurityPolicy{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     idsParentDfwPolID,
			"domain": idsParentDfwPolDomain,
		})

		err := dataSourceNsxtPolicyParentIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading")
	})

	t.Run("missing id name and category", func(t *testing.T) {
		ds := dataSourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain": idsParentDfwPolDomain,
		})

		err := dataSourceNsxtPolicyParentIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be specified")
	})

	t.Run("by display_name via list", func(t *testing.T) {
		rc := int64(1)
		mockSDK.EXPECT().List(idsParentDfwPolDomain, nil, &boolFalse, nil, nil, nil, nil, nil).Return(nsxModel.IdsSecurityPolicyListResult{
			Results:     []nsxModel.IdsSecurityPolicy{idsParentDfwPolicyModel()},
			ResultCount: &rc,
		}, nil)

		ds := dataSourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":       idsParentDfwPolDomain,
			"display_name": idsParentDfwPolName,
		})

		err := dataSourceNsxtPolicyParentIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsParentDfwPolID, d.Id())
		assert.Equal(t, idsParentDfwPolName, d.Get("display_name"))
		assert.Equal(t, idsParentDfwPolCategory, d.Get("category"))
	})

	t.Run("by category via list", func(t *testing.T) {
		rc := int64(1)
		mockSDK.EXPECT().List(idsParentDfwPolDomain, nil, &boolFalse, nil, nil, nil, nil, nil).Return(nsxModel.IdsSecurityPolicyListResult{
			Results:     []nsxModel.IdsSecurityPolicy{idsParentDfwPolicyModel()},
			ResultCount: &rc,
		}, nil)

		ds := dataSourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":   idsParentDfwPolDomain,
			"category": idsParentDfwPolCategory,
		})

		err := dataSourceNsxtPolicyParentIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsParentDfwPolID, d.Id())
		assert.Equal(t, idsParentDfwPolName, d.Get("display_name"))
		assert.Equal(t, idsParentDfwPolCategory, d.Get("category"))
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List(idsParentDfwPolDomain, nil, &boolFalse, nil, nil, nil, nil, nil).Return(nsxModel.IdsSecurityPolicyListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":       idsParentDfwPolDomain,
			"display_name": idsParentDfwPolName,
		})

		err := dataSourceNsxtPolicyParentIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("multiple matches error", func(t *testing.T) {
		rc := int64(2)
		policy2 := idsParentDfwPolicyModel()
		id2 := "ids-parent-dfw-pol-2"
		policy2.Id = &id2

		mockSDK.EXPECT().List(idsParentDfwPolDomain, nil, &boolFalse, nil, nil, nil, nil, nil).Return(nsxModel.IdsSecurityPolicyListResult{
			Results:     []nsxModel.IdsSecurityPolicy{idsParentDfwPolicyModel(), policy2},
			ResultCount: &rc,
		}, nil)

		ds := dataSourceNsxtPolicyParentIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":       idsParentDfwPolDomain,
			"display_name": idsParentDfwPolName,
		})

		err := dataSourceNsxtPolicyParentIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Found multiple")
	})
}
