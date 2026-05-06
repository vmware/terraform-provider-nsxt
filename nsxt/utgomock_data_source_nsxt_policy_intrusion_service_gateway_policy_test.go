//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"

	apidomains "github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	domainmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains"
)

var (
	idsGwPolDomain   = "default"
	idsGwPolID       = "ids-gw-pol-1"
	idsGwPolName     = "ids-gateway-policy"
	idsGwPolPath     = "/infra/domains/default/gateway-policies/ids-gw-pol-1"
	idsGwPolCategory = "LOCALGATEWAY"
	idsGwPolDesc     = "desc"
)

func idsGatewayPolicyModel() nsxModel.IdsGatewayPolicy {
	resourceType := "IdsRule"
	ruleID := "test-gw-rule-id"
	ruleName := "test-gw-rule"
	ruleNotes := "Test Gateway rule notes"
	ruleAction := "DETECT"
	ruleDirection := "IN_OUT"
	ruleId := int64(12345)

	return nsxModel.IdsGatewayPolicy{
		Id:          &idsGwPolID,
		DisplayName: &idsGwPolName,
		Description: &idsGwPolDesc,
		Path:        &idsGwPolPath,
		Category:    &idsGwPolCategory,
		Rules: []nsxModel.IdsRule{
			{
				Id:           &ruleID,
				DisplayName:  &ruleName,
				Notes:        &ruleNotes,
				ResourceType: &resourceType,
				Action:       &ruleAction,
				Direction:    &ruleDirection,
				RuleId:       &ruleId,
			},
		},
	}
}

func setupIntrusionServiceGatewayPolicyDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockIntrusionServiceGatewayPoliciesClient, func()) {
	t.Helper()
	mockSDK := domainmocks.NewMockIntrusionServiceGatewayPoliciesClient(ctrl)
	wrapper := &apidomains.IntrusionServiceGatewayPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliIntrusionServiceGatewayPoliciesClient
	cliIntrusionServiceGatewayPoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.IntrusionServiceGatewayPolicyClientContext {
		return wrapper
	}
	return mockSDK, func() { cliIntrusionServiceGatewayPoliciesClient = orig }
}

func TestMockDataSourceNsxtPolicyIntrusionServiceGatewayPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIntrusionServiceGatewayPolicyDataSourceMock(t, ctrl)
	defer restore()

	boolFalse := false

	t.Run("by id success", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwPolDomain, idsGwPolID).Return(idsGatewayPolicyModel(), nil)

		ds := dataSourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     idsGwPolID,
			"domain": idsGwPolDomain,
		})

		err := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsGwPolID, d.Id())
		assert.Equal(t, idsGwPolName, d.Get("display_name"))
		assert.Equal(t, idsGwPolCategory, d.Get("category"))
		// Check that rules are present
		rules := d.Get("rule").([]interface{})
		assert.Len(t, rules, 1)
		rule := rules[0].(map[string]interface{})
		assert.Equal(t, 12345, rule["rule_id"])
		assert.Equal(t, "test-gw-rule", rule["display_name"])
		assert.Equal(t, "DETECT", rule["action"])
		_, hasNsxId := rule["nsx_id"]
		assert.True(t, hasNsxId, "nsx_id should be included in IDPS data source embedded rules for consistency with resources")
		_, hasScope := rule["scope"]
		assert.True(t, hasScope, "scope should be included in IDPS data source embedded rules")
		// Verify oversubscription is excluded for Gateway rules
		_, hasOversubscription := rule["oversubscription"]
		assert.False(t, hasOversubscription, "oversubscription should be excluded from Gateway IDPS rules")
	})

	t.Run("by id not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwPolDomain, idsGwPolID).Return(nsxModel.IdsGatewayPolicy{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     idsGwPolID,
			"domain": idsGwPolDomain,
		})

		err := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by id API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwPolDomain, idsGwPolID).Return(nsxModel.IdsGatewayPolicy{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     idsGwPolID,
			"domain": idsGwPolDomain,
		})

		err := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading")
	})

	t.Run("missing id name and category", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain": idsGwPolDomain,
		})

		err := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be specified")
	})

	t.Run("by display_name via list", func(t *testing.T) {
		rc := int64(1)
		mockSDK.EXPECT().List(idsGwPolDomain, nil, &boolFalse, nil, nil, nil, nil, nil).Return(nsxModel.IdsGatewayPolicyListResult{
			Results:     []nsxModel.IdsGatewayPolicy{idsGatewayPolicyModel()},
			ResultCount: &rc,
		}, nil)
		mockSDK.EXPECT().Get(idsGwPolDomain, idsGwPolID).Return(idsGatewayPolicyModel(), nil)

		ds := dataSourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":       idsGwPolDomain,
			"display_name": idsGwPolName,
		})

		err := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsGwPolID, d.Id())
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List(idsGwPolDomain, nil, &boolFalse, nil, nil, nil, nil, nil).Return(nsxModel.IdsGatewayPolicyListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIntrusionServiceGatewayPolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":       idsGwPolDomain,
			"display_name": idsGwPolName,
		})

		err := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
