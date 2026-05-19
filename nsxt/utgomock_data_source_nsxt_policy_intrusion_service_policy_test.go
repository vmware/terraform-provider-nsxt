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
	idsDfwPolDomain   = "default"
	idsDfwPolID       = "ids-dfw-pol-1"
	idsDfwPolName     = "ids-dfw-policy"
	idsDfwPolPath     = "/infra/domains/default/intrusion-service-policies/ids-dfw-pol-1"
	idsDfwPolCategory = "ThreatRules"
	idsDfwPolDesc     = "desc"
)

func idsDfwPolicyModel() nsxModel.IdsSecurityPolicy {
	stateful := true
	resourceType := "IdsRule"
	ruleID := "test-rule-id"
	ruleName := "test-rule"
	ruleNotes := "Test rule notes"
	ruleAction := "DETECT"
	ruleDirection := "IN_OUT"
	ruleId := int64(12345)

	return nsxModel.IdsSecurityPolicy{
		Id:          &idsDfwPolID,
		DisplayName: &idsDfwPolName,
		Description: &idsDfwPolDesc,
		Path:        &idsDfwPolPath,
		Category:    &idsDfwPolCategory,
		Stateful:    &stateful,
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

func setupIntrusionServiceDfwPolicyDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockIntrusionServicePoliciesClient, func()) {
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

func TestMockDataSourceNsxtPolicyIntrusionServicePolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIntrusionServiceDfwPolicyDataSourceMock(t, ctrl)
	defer restore()

	boolFalse := false

	t.Run("by id success", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsDfwPolDomain, idsDfwPolID).Return(idsDfwPolicyModel(), nil)

		ds := dataSourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     idsDfwPolID,
			"domain": idsDfwPolDomain,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsDfwPolID, d.Id())
		assert.Equal(t, idsDfwPolName, d.Get("display_name"))
		assert.Equal(t, idsDfwPolCategory, d.Get("category"))
		assert.Equal(t, true, d.Get("stateful"))
		// Check that rules are present
		rules := d.Get("rule").([]interface{})
		assert.Len(t, rules, 1)
		rule := rules[0].(map[string]interface{})
		assert.Equal(t, 12345, rule["rule_id"])
		assert.Equal(t, "test-rule", rule["display_name"])
		assert.Equal(t, "DETECT", rule["action"])
		_, hasNsxId := rule["nsx_id"]
		assert.True(t, hasNsxId, "nsx_id should be included in IDPS data source embedded rules for consistency with resources")
		// Verify scope is included in data source embedded rules (as per original implementation)
		_, hasScope := rule["scope"]
		assert.True(t, hasScope, "scope should be included in IDPS data source embedded rules")
	})

	t.Run("by id not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsDfwPolDomain, idsDfwPolID).Return(nsxModel.IdsSecurityPolicy{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     idsDfwPolID,
			"domain": idsDfwPolDomain,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by id API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsDfwPolDomain, idsDfwPolID).Return(nsxModel.IdsSecurityPolicy{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":     idsDfwPolID,
			"domain": idsDfwPolDomain,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading")
	})

	t.Run("missing id name and category", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain": idsDfwPolDomain,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be specified")
	})

	t.Run("by display_name via list", func(t *testing.T) {
		rc := int64(1)
		mockSDK.EXPECT().List(idsDfwPolDomain, nil, &boolFalse, nil, nil, nil, nil, nil).Return(nsxModel.IdsSecurityPolicyListResult{
			Results:     []nsxModel.IdsSecurityPolicy{idsDfwPolicyModel()},
			ResultCount: &rc,
		}, nil)
		mockSDK.EXPECT().Get(idsDfwPolDomain, idsDfwPolID).Return(idsDfwPolicyModel(), nil)

		ds := dataSourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":       idsDfwPolDomain,
			"display_name": idsDfwPolName,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsDfwPolID, d.Id())
		assert.Equal(t, idsDfwPolName, d.Get("display_name"))
		assert.Equal(t, idsDfwPolCategory, d.Get("category"))
	})

	t.Run("by category via list", func(t *testing.T) {
		rc := int64(1)
		mockSDK.EXPECT().List(idsDfwPolDomain, nil, &boolFalse, nil, nil, nil, nil, nil).Return(nsxModel.IdsSecurityPolicyListResult{
			Results:     []nsxModel.IdsSecurityPolicy{idsDfwPolicyModel()},
			ResultCount: &rc,
		}, nil)
		// Expect additional Get call to fetch complete policy with rules
		mockSDK.EXPECT().Get(idsDfwPolDomain, idsDfwPolID).Return(idsDfwPolicyModel(), nil)

		ds := dataSourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":   idsDfwPolDomain,
			"category": idsDfwPolCategory,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsDfwPolID, d.Id())
		assert.Equal(t, idsDfwPolName, d.Get("display_name"))
		assert.Equal(t, idsDfwPolCategory, d.Get("category"))
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List(idsDfwPolDomain, nil, &boolFalse, nil, nil, nil, nil, nil).Return(nsxModel.IdsSecurityPolicyListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":       idsDfwPolDomain,
			"display_name": idsDfwPolName,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("multiple matches error", func(t *testing.T) {
		rc := int64(2)
		policy2 := idsDfwPolicyModel()
		id2 := "ids-dfw-pol-2"
		policy2.Id = &id2

		mockSDK.EXPECT().List(idsDfwPolDomain, nil, &boolFalse, nil, nil, nil, nil, nil).Return(nsxModel.IdsSecurityPolicyListResult{
			Results:     []nsxModel.IdsSecurityPolicy{idsDfwPolicyModel(), policy2},
			ResultCount: &rc,
		}, nil)

		ds := dataSourceNsxtPolicyIntrusionServicePolicy()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"domain":       idsDfwPolDomain,
			"display_name": idsDfwPolName,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Found multiple")
	})
}
