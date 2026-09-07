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

	apipkg "github.com/vmware/terraform-provider-nsxt/api"
	apidomains "github.com/vmware/terraform-provider-nsxt/api/infra/domains"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	inframocks "github.com/vmware/terraform-provider-nsxt/mocks/infra"
	domainmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains"
)

var (
	predefinedGwPolicyID          = "default"
	predefinedGwPolicyDomain      = "default"
	predefinedGwPolicyPath        = "/infra/domains/default/gateway-policies/default"
	predefinedGwPolicyDescription = "Predefined Gateway Policy"
	predefinedGwPolicyRevision    = int64(1)
	predefinedGwPolicyCategory    = "Default"
)

func predefinedGwPolicyAPIResponse() nsxModel.GatewayPolicy {
	return nsxModel.GatewayPolicy{
		Id:          &predefinedGwPolicyID,
		Description: &predefinedGwPolicyDescription,
		Revision:    &predefinedGwPolicyRevision,
		Path:        &predefinedGwPolicyPath,
		Category:    &predefinedGwPolicyCategory,
	}
}

func minimalPredefinedGwPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"path": predefinedGwPolicyPath,
	}
}

func setupPredefinedGwPolicyMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockGatewayPoliciesClient, *inframocks.MockInfraClient, func()) {
	mockSDK := domainmocks.NewMockGatewayPoliciesClient(ctrl)
	mockWrapper := &apidomains.GatewayPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalGwPolicy := cliGatewayPoliciesClient
	cliGatewayPoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.GatewayPolicyClientContext {
		return mockWrapper
	}

	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)
	originalInfra := cliInfraClient
	cliInfraClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	return mockSDK, mockInfraSDK, func() {
		cliGatewayPoliciesClient = originalGwPolicy
		cliInfraClient = originalInfra
	}
}

func TestMockResourceNsxtPolicyPredefinedGatewayPolicyCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupPredefinedGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(predefinedGwPolicyAPIResponse(), nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(predefinedGwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())

		err := resourceNsxtPolicyPredefinedGatewayPolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, predefinedGwPolicyID, d.Id())
	})

	t.Run("Create fails when path is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"path": "",
		})

		err := resourceNsxtPolicyPredefinedGatewayPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyPredefinedGatewayPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupPredefinedGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(predefinedGwPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())
		d.SetId(predefinedGwPolicyID)

		err := resourceNsxtPolicyPredefinedGatewayPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, predefinedGwPolicyDescription, d.Get("description"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())

		err := resourceNsxtPolicyPredefinedGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails on API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(nsxModel.GatewayPolicy{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())
		d.SetId(predefinedGwPolicyID)

		err := resourceNsxtPolicyPredefinedGatewayPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyPredefinedGatewayPolicyUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupPredefinedGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(predefinedGwPolicyAPIResponse(), nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(predefinedGwPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())
		d.SetId(predefinedGwPolicyID)

		err := resourceNsxtPolicyPredefinedGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())

		err := resourceNsxtPolicyPredefinedGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyPredefinedGatewayPolicyDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupPredefinedGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Delete (revert) success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(predefinedGwPolicyAPIResponse(), nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
		)

		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())
		d.SetId(predefinedGwPolicyID)

		err := resourceNsxtPolicyPredefinedGatewayPolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())

		err := resourceNsxtPolicyPredefinedGatewayPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyPredefinedGatewayPolicyUpdateSystemPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupPredefinedGwPolicyMock(t, ctrl)
	defer restore()

	t.Run("Update fails when policy is a system policy", func(t *testing.T) {
		systemCategory := "SystemRules"
		systemPolicy := predefinedGwPolicyAPIResponse()
		systemPolicy.Category = &systemCategory
		mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(systemPolicy, nil)

		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())
		d.SetId(predefinedGwPolicyID)

		err := resourceNsxtPolicyPredefinedGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "System policy")
	})

	t.Run("Update fails when path has no domain segment", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedGatewayPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{"path": "not-a-valid-path"})
		d.SetId(predefinedGwPolicyID)

		err := resourceNsxtPolicyPredefinedGatewayPolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestUnitNsxt_revertGatewayPolicyDefaultRule(t *testing.T) {
	id := "rule-1"
	description := "custom description"
	action := nsxModel.Rule_ACTION_DROP
	tags := []nsxModel.Tag{{Scope: strPtr("s"), Tag: strPtr("t")}}
	rule := nsxModel.Rule{Id: &id, Description: &description, Action: &action, Tags: tags}

	reverted := revertGatewayPolicyDefaultRule(rule)
	assert.Equal(t, "", *reverted.Description)
	assert.Equal(t, nsxModel.Rule_ACTION_ALLOW, *reverted.Action)
	assert.Empty(t, reverted.Tags)
}

func TestUnitNsxt_createPolicyChildRule(t *testing.T) {
	ruleID := "rule-1"
	rule := nsxModel.Rule{Id: &ruleID}

	dataValue, err := createPolicyChildRule(ruleID, rule, false)
	require.NoError(t, err)
	assert.NotNil(t, dataValue)

	deleteValue, err := createPolicyChildRule(ruleID, rule, true)
	require.NoError(t, err)
	assert.NotNil(t, deleteValue)
}

func TestUnitNsxt_createChildDomainWithGatewayPolicy(t *testing.T) {
	policyID := "default"
	policy := nsxModel.GatewayPolicy{Id: &policyID}

	dataValue, err := createChildDomainWithGatewayPolicy("default", policyID, policy)
	require.NoError(t, err)
	assert.NotNil(t, dataValue)
}

func TestUnitNsxt_revertPolicyPredefinedGatewayPolicy(t *testing.T) {
	defaultRuleID := "default-rule"
	nonDefaultRuleID := "custom-rule"
	isDefault := true
	notDefault := false
	tags := []nsxModel.Tag{{Scope: strPtr("s"), Tag: strPtr("t")}}

	policy := nsxModel.GatewayPolicy{
		Rules: []nsxModel.Rule{
			{Id: &defaultRuleID, IsDefault: &isDefault},
			{Id: &nonDefaultRuleID, IsDefault: &notDefault},
		},
		Tags: tags,
	}

	reverted, err := revertPolicyPredefinedGatewayPolicy(policy, newGoMockProviderClient())
	require.NoError(t, err)
	assert.Equal(t, "", *reverted.Description)
	assert.Nil(t, reverted.Rules)
	assert.Len(t, reverted.Children, 2)
	assert.Empty(t, reverted.Tags)
}

func TestUnitNsxt_setPolicyDefaultRulesInSchema(t *testing.T) {
	res := resourceNsxtPolicyPredefinedGatewayPolicy()
	d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())

	ruleID := "rule-1"
	description := "a default rule"
	scope := []string{"ANY"}
	rules := []nsxModel.Rule{{Id: &ruleID, Description: &description, Scope: scope}}

	err := setPolicyDefaultRulesInSchema(d, rules)
	require.NoError(t, err)
	set := d.Get("default_rule").(*schema.Set)
	assert.Equal(t, 1, set.Len())
}

func TestUnitNsxt_updateGatewayPolicyDefaultRuleByScope(t *testing.T) {
	res := resourceNsxtPolicyPredefinedGatewayPolicy()

	t.Run("matches rule by scope and updates fields", func(t *testing.T) {
		data := minimalPredefinedGwPolicyData()
		data["default_rule"] = []interface{}{
			map[string]interface{}{
				"nsx_id":      "rule-1",
				"scope":       "/infra/domains/default/groups/g1",
				"description": "updated desc",
				"action":      nsxModel.Rule_ACTION_DROP,
				"logged":      true,
				"log_label":   "lbl",
				"tag":         []interface{}{},
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)

		ruleID := "rule-1"
		rule := nsxModel.Rule{Id: &ruleID, Scope: []string{"/infra/domains/default/groups/g1"}}

		updated := updateGatewayPolicyDefaultRuleByScope(rule, d, nil, false)
		require.NotNil(t, updated)
		assert.Equal(t, "updated desc", *updated.Description)
		assert.Equal(t, nsxModel.Rule_ACTION_DROP, *updated.Action)
	})

	t.Run("no match and not previously deleted returns nil", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())

		ruleID := "rule-unmatched"
		rule := nsxModel.Rule{Id: &ruleID, Scope: []string{"/infra/domains/default/groups/other"}}

		updated := updateGatewayPolicyDefaultRuleByScope(rule, d, nil, false)
		assert.Nil(t, updated)
	})
}

func TestUnitNsxt_nsxtPredefinedPolicyImporter(t *testing.T) {
	res := resourceNsxtPolicyPredefinedGatewayPolicy()

	t.Run("non-policy-path ID sets path and id directly", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId("plain-id")

		out, err := nsxtPredefinedPolicyImporter(d, newGoMockProviderClient())
		require.NoError(t, err)
		require.Len(t, out, 1)
		assert.Equal(t, "plain-id", d.Get("path"))
	})
}
