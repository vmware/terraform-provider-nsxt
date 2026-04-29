//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	secpolapi "github.com/vmware/terraform-provider-nsxt/api/infra/domains/security_policies"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	secpolmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains/security_policies"
)

var (
	spRuleID          = "sp-rule-1"
	spRuleDisplayName = "Test SP Rule"
	spRuleDescription = "Test security policy rule"
	spRuleRevision    = int64(1)
	spRuleDomain      = "default"
	spRulePolicyID    = "sp-001"
	spRulePolicyPath  = "/infra/domains/default/security-policies/sp-001"
	spRuleAction      = "ALLOW"
	spRuleDirection   = "IN_OUT"
	spRuleIPVersion   = "IPV4_IPV6"
)

func spRuleAPIResponse() nsxModel.Rule {
	resourceType := "Rule"
	return nsxModel.Rule{
		Id:           &spRuleID,
		DisplayName:  &spRuleDisplayName,
		Description:  &spRuleDescription,
		Revision:     &spRuleRevision,
		ResourceType: &resourceType,
		Action:       &spRuleAction,
		Direction:    &spRuleDirection,
	}
}

func minimalSPRuleData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": spRuleDisplayName,
		"description":  spRuleDescription,
		"policy_path":  spRulePolicyPath,
		"action":       spRuleAction,
		"direction":    spRuleDirection,
		"ip_version":   spRuleIPVersion,
		"logged":       false,
		"disabled":     false,
	}
}

func setupSPRuleMock(t *testing.T, ctrl *gomock.Controller) (*secpolmocks.MockRulesClient, func()) {
	mockSDK := secpolmocks.NewMockRulesClient(ctrl)
	mockWrapper := &secpolapi.RuleClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliSecurityPoliciesRulesClient
	cliSecurityPoliciesRulesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *secpolapi.RuleClientContext {
		return mockWrapper
	}

	return mockSDK, func() { cliSecurityPoliciesRulesClient = original }
}

func TestMockResourceNsxtPolicySecurityPolicyRuleCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupSPRuleMock(t, ctrl)
	defer restore()

	t.Run("Create success with auto-generated ID", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(spRuleDomain, spRulePolicyID, gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(spRuleDomain, spRulePolicyID, gomock.Any()).Return(spRuleAPIResponse(), nil),
		)

		res := resourceNsxtPolicySecurityPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSPRuleData())

		err := resourceNsxtPolicySecurityPolicyRuleCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicySecurityPolicyRuleRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupSPRuleMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(spRuleDomain, spRulePolicyID, spRuleID).Return(spRuleAPIResponse(), nil)

		res := resourceNsxtPolicySecurityPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSPRuleData())
		d.SetId(spRuleID)

		err := resourceNsxtPolicySecurityPolicyRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, spRuleDisplayName, d.Get("display_name"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(spRuleDomain, spRulePolicyID, spRuleID).Return(nsxModel.Rule{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicySecurityPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSPRuleData())
		d.SetId(spRuleID)

		err := resourceNsxtPolicySecurityPolicyRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Empty(t, d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySecurityPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSPRuleData())

		err := resourceNsxtPolicySecurityPolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySecurityPolicyRuleUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupSPRuleMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(spRuleDomain, spRulePolicyID, spRuleID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(spRuleDomain, spRulePolicyID, spRuleID).Return(spRuleAPIResponse(), nil),
		)

		res := resourceNsxtPolicySecurityPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSPRuleData())
		d.SetId(spRuleID)

		err := resourceNsxtPolicySecurityPolicyRuleUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicySecurityPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSPRuleData())

		err := resourceNsxtPolicySecurityPolicyRuleUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySecurityPolicyRuleDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupSPRuleMock(t, ctrl)
	defer restore()

	t.Run("Delete success", func(t *testing.T) {
		dataWithNsxID := minimalSPRuleData()
		dataWithNsxID["nsx_id"] = spRuleID
		mockSDK.EXPECT().Delete(spRuleDomain, spRulePolicyID, spRuleID).Return(nil)

		res := resourceNsxtPolicySecurityPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, dataWithNsxID)
		d.SetId(spRuleID)

		err := resourceNsxtPolicySecurityPolicyRuleDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when nsx_id is empty", func(t *testing.T) {
		res := resourceNsxtPolicySecurityPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalSPRuleData())
		d.SetId(spRuleID)

		err := resourceNsxtPolicySecurityPolicyRuleDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicySecurityPolicyRuleExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupSPRuleMock(t, ctrl)
	defer restore()

	m := newGoMockProviderClient()
	connector := getPolicyConnector(m)
	sctx := utl.SessionContext{ClientType: utl.Local}

	t.Run("exists when rule is present", func(t *testing.T) {
		mockSDK.EXPECT().Get(spRuleDomain, spRulePolicyID, spRuleID).Return(spRuleAPIResponse(), nil)
		ok, err := resourceNsxtPolicySecurityPolicyRuleExists(sctx, spRuleID, spRulePolicyPath, connector)
		require.NoError(t, err)
		assert.True(t, ok)
	})

	t.Run("not exists on NotFound", func(t *testing.T) {
		mockSDK.EXPECT().Get(spRuleDomain, spRulePolicyID, spRuleID).Return(nsxModel.Rule{}, vapiErrors.NotFound{})
		ok, err := resourceNsxtPolicySecurityPolicyRuleExists(sctx, spRuleID, spRulePolicyPath, connector)
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("error on unexpected API failure", func(t *testing.T) {
		mockSDK.EXPECT().Get(spRuleDomain, spRulePolicyID, spRuleID).Return(nsxModel.Rule{}, vapiErrors.ServiceUnavailable{})
		ok, err := resourceNsxtPolicySecurityPolicyRuleExists(sctx, spRuleID, spRulePolicyPath, connector)
		require.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("error when rules client is not supported (nil)", func(t *testing.T) {
		orig := cliSecurityPoliciesRulesClient
		cliSecurityPoliciesRulesClient = func(utl.SessionContext, vapiProtocolClient.Connector) *secpolapi.RuleClientContext {
			return nil
		}
		defer func() { cliSecurityPoliciesRulesClient = orig }()
		ok, err := resourceNsxtPolicySecurityPolicyRuleExists(sctx, spRuleID, spRulePolicyPath, connector)
		require.Error(t, err)
		assert.False(t, ok)
	})
}

func TestMockResourceNsxtPolicySecurityPolicyRuleImport(t *testing.T) {
	m := newGoMockProviderClient()
	res := resourceNsxtPolicySecurityPolicyRule()

	importPath := "/infra/domains/default/security-policies/sp-001/rules/sp-rule-1"
	ruleIdx := strings.Index(importPath, "rule")
	require.Greater(t, ruleIdx, 0)
	wantPolicyPath := importPath[:ruleIdx-1]

	t.Run("import success sets policy_path", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(importPath)

		rd, err := nsxtSecurityPolicyRuleImporter(d, m)
		require.NoError(t, err)
		require.Len(t, rd, 1)
		assert.Equal(t, wantPolicyPath, rd[0].Get("policy_path"))
	})

	t.Run("import fails when path has no rule segment", func(t *testing.T) {
		badPath := "/infra/domains/default/security-policies/sp-001"
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		d.SetId(badPath)

		_, err := nsxtSecurityPolicyRuleImporter(d, m)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid path")
	})
}
