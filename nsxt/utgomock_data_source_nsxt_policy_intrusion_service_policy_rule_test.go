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

	idsruleapi "github.com/vmware/terraform-provider-nsxt/api/infra/domains/intrusion_service_policies"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	idsrulemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains/intrusion_service_policies"
)

var (
	idsDfwRuleDomain     = "default"
	idsDfwRulePolicyID   = "ids-dfw-pol-1"
	idsDfwRuleID         = "ids-dfw-rule-1"
	idsDfwRuleName       = "ids-dfw-rule"
	idsDfwRulePath       = "/infra/domains/default/intrusion-service-policies/ids-dfw-pol-1/rules/ids-dfw-rule-1"
	idsDfwRulePolicyPath = "/infra/domains/default/intrusion-service-policies/ids-dfw-pol-1"
	idsDfwRuleAction     = "DETECT"
	idsDfwRuleDesc       = "rule desc"
)

func idsDfwRuleModel() nsxModel.IdsRule {
	logged := true
	disabled := false
	sequenceNumber := int64(1)
	ruleId := int64(12345)
	return nsxModel.IdsRule{
		Id:             &idsDfwRuleID,
		DisplayName:    &idsDfwRuleName,
		Description:    &idsDfwRuleDesc,
		Path:           &idsDfwRulePath,
		Action:         &idsDfwRuleAction,
		Logged:         &logged,
		Disabled:       &disabled,
		SequenceNumber: &sequenceNumber,
		RuleId:         &ruleId,
	}
}

func setupIntrusionServiceDfwPolicyRuleDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*idsrulemocks.MockRulesClient, func()) {
	t.Helper()
	mockSDK := idsrulemocks.NewMockRulesClient(ctrl)
	wrapper := &idsruleapi.IntrusionServicePolicyRuleClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliIntrusionServicePolicyRulesClient
	cliIntrusionServicePolicyRulesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *idsruleapi.IntrusionServicePolicyRuleClientContext {
		return wrapper
	}
	return mockSDK, func() { cliIntrusionServicePolicyRulesClient = orig }
}

func TestMockDataSourceNsxtPolicyIntrusionServicePolicyRuleRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIntrusionServiceDfwPolicyRuleDataSourceMock(t, ctrl)
	defer restore()

	boolFalse := false

	t.Run("by id success", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsDfwRuleDomain, idsDfwRulePolicyID, idsDfwRuleID).Return(idsDfwRuleModel(), nil)

		ds := dataSourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":          idsDfwRuleID,
			"policy_path": idsDfwRulePolicyPath,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsDfwRuleID, d.Id())
		assert.Equal(t, idsDfwRuleName, d.Get("display_name"))
		assert.Equal(t, idsDfwRuleAction, d.Get("action"))
		assert.Equal(t, true, d.Get("logged"))
		assert.Equal(t, false, d.Get("disabled"))
		assert.Equal(t, 1, d.Get("sequence_number"))
		assert.Equal(t, 12345, d.Get("rule_id"))
	})

	t.Run("by id not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsDfwRuleDomain, idsDfwRulePolicyID, idsDfwRuleID).Return(nsxModel.IdsRule{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":          idsDfwRuleID,
			"policy_path": idsDfwRulePolicyPath,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by id API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsDfwRuleDomain, idsDfwRulePolicyID, idsDfwRuleID).Return(nsxModel.IdsRule{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":          idsDfwRuleID,
			"policy_path": idsDfwRulePolicyPath,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading")
	})

	t.Run("missing policy_path", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": idsDfwRuleID,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "policy_path must be specified")
	})

	t.Run("missing id and display_name", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"policy_path": idsDfwRulePolicyPath,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be specified")
	})

	t.Run("by display_name via list", func(t *testing.T) {
		rc := int64(1)
		mockSDK.EXPECT().List(idsDfwRuleDomain, idsDfwRulePolicyID, nil, &boolFalse, nil, nil, nil, nil).Return(nsxModel.IdsRuleListResult{
			Results:     []nsxModel.IdsRule{idsDfwRuleModel()},
			ResultCount: &rc,
		}, nil)

		ds := dataSourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"policy_path":  idsDfwRulePolicyPath,
			"display_name": idsDfwRuleName,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsDfwRuleID, d.Id())
		assert.Equal(t, idsDfwRuleName, d.Get("display_name"))
		assert.Equal(t, idsDfwRuleAction, d.Get("action"))
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List(idsDfwRuleDomain, idsDfwRulePolicyID, nil, &boolFalse, nil, nil, nil, nil).Return(nsxModel.IdsRuleListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"policy_path":  idsDfwRulePolicyPath,
			"display_name": idsDfwRuleName,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("multiple matches error", func(t *testing.T) {
		rc := int64(2)
		rule2 := idsDfwRuleModel()
		id2 := "ids-dfw-rule-2"
		rule2.Id = &id2

		mockSDK.EXPECT().List(idsDfwRuleDomain, idsDfwRulePolicyID, nil, &boolFalse, nil, nil, nil, nil).Return(nsxModel.IdsRuleListResult{
			Results:     []nsxModel.IdsRule{idsDfwRuleModel(), rule2},
			ResultCount: &rc,
		}, nil)

		ds := dataSourceNsxtPolicyIntrusionServicePolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"policy_path":  idsDfwRulePolicyPath,
			"display_name": idsDfwRuleName,
		})

		err := dataSourceNsxtPolicyIntrusionServicePolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Found multiple")
	})
}
