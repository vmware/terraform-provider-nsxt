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
	orgrootmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsxt"
)

var (
	predefinedSecPolicyID          = "default"
	predefinedSecPolicyDomain      = "default"
	predefinedSecPolicyPath        = "/infra/domains/default/security-policies/default"
	predefinedSecPolicyDescription = "Predefined Security Policy"
	predefinedSecPolicyRevision    = int64(1)
	predefinedSecPolicyCategory    = "Default"
)

func predefinedSecPolicyAPIResponse() nsxModel.SecurityPolicy {
	return nsxModel.SecurityPolicy{
		Id:          &predefinedSecPolicyID,
		Description: &predefinedSecPolicyDescription,
		Revision:    &predefinedSecPolicyRevision,
		Path:        &predefinedSecPolicyPath,
		Category:    &predefinedSecPolicyCategory,
	}
}

func minimalPredefinedSecPolicyData() map[string]interface{} {
	return map[string]interface{}{
		"path": predefinedSecPolicyPath,
	}
}

func setupPredefinedSecPolicyMock(t *testing.T, ctrl *gomock.Controller) (*domainmocks.MockSecurityPoliciesClient, *inframocks.MockInfraClient, func()) {
	mockSDK := domainmocks.NewMockSecurityPoliciesClient(ctrl)
	mockWrapper := &apidomains.SecurityPolicyClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	originalSecPolicy := cliSecurityPoliciesClient
	cliSecurityPoliciesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apidomains.SecurityPolicyClientContext {
		return mockWrapper
	}

	mockInfraSDK := inframocks.NewMockInfraClient(ctrl)
	originalInfra := cliInfraClient
	cliInfraClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apipkg.InfraClientContext {
		return &apipkg.InfraClientContext{Client: mockInfraSDK, ClientType: utl.Local}
	}

	return mockSDK, mockInfraSDK, func() {
		cliSecurityPoliciesClient = originalSecPolicy
		cliInfraClient = originalInfra
	}
}

func TestMockResourceNsxtPolicyPredefinedSecurityPolicyCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupPredefinedSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Create success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())

		err := resourceNsxtPolicyPredefinedSecurityPolicyCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, predefinedSecPolicyID, d.Id())
	})

	t.Run("Create fails when path is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"path": "",
		})

		err := resourceNsxtPolicyPredefinedSecurityPolicyCreate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyPredefinedSecurityPolicyRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, _, restore := setupPredefinedSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil)

		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())
		d.SetId(predefinedSecPolicyID)

		err := resourceNsxtPolicyPredefinedSecurityPolicyRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, predefinedSecPolicyDescription, d.Get("description"))
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())

		err := resourceNsxtPolicyPredefinedSecurityPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("Read fails on API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(nsxModel.SecurityPolicy{}, vapiErrors.InternalServerError{})

		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())
		d.SetId(predefinedSecPolicyID)

		err := resourceNsxtPolicyPredefinedSecurityPolicyRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyPredefinedSecurityPolicyUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupPredefinedSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())
		d.SetId(predefinedSecPolicyID)

		err := resourceNsxtPolicyPredefinedSecurityPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())

		err := resourceNsxtPolicyPredefinedSecurityPolicyUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyPredefinedSecurityPolicyDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, mockInfra, restore := setupPredefinedSecPolicyMock(t, ctrl)
	defer restore()

	t.Run("Delete (revert) success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
		)

		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())
		d.SetId(predefinedSecPolicyID)

		err := resourceNsxtPolicyPredefinedSecurityPolicyDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())

		err := resourceNsxtPolicyPredefinedSecurityPolicyDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestUnitNsxt_updateSecurityPolicyDefaultRule(t *testing.T) {
	res := resourceNsxtPolicyPredefinedSecurityPolicy()

	t.Run("applies fields from the configured default_rule block", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"default_rule": []interface{}{
				map[string]interface{}{
					"description": "my default rule",
					"action":      "DROP",
					"log_label":   "mylabel",
					"logged":      true,
				},
			},
		})
		ruleID := "rule-1"
		rule := nsxModel.Rule{Id: &ruleID}

		updated := updateSecurityPolicyDefaultRule(rule, d)
		require.NotNil(t, updated)
		assert.Equal(t, "my default rule", *updated.Description)
		assert.Equal(t, "DROP", *updated.Action)
		assert.Equal(t, "mylabel", *updated.Tag)
		assert.True(t, *updated.Logged)
	})

	t.Run("no default_rule block and no change returns nil", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		ruleID := "rule-1"
		rule := nsxModel.Rule{Id: &ruleID}

		updated := updateSecurityPolicyDefaultRule(rule, d)
		assert.Nil(t, updated)
	})
}

func TestUnitNsxt_revertSecurityPolicyDefaultRule(t *testing.T) {
	t.Run("resets action, logged, tag and label to defaults", func(t *testing.T) {
		action := "DROP"
		logged := true
		tag := "mylabel"
		rule := nsxModel.Rule{
			Action: &action,
			Logged: &logged,
			Tag:    &tag,
			Tags:   []nsxModel.Tag{{Scope: strPtr("s"), Tag: strPtr("t")}},
		}

		reverted := revertSecurityPolicyDefaultRule(rule)
		assert.Equal(t, "ALLOW", *reverted.Action)
		assert.False(t, *reverted.Logged)
		assert.Equal(t, "", *reverted.Tag)
		assert.Empty(t, reverted.Tags)
	})
}

func TestUnitNsxt_revertPolicyPredefinedSecurityPolicy(t *testing.T) {
	t.Run("reverts default rules and deletes non-default rules", func(t *testing.T) {
		isDefault := true
		notDefault := false
		defaultRuleID := "default-rule"
		customRuleID := "custom-rule"
		action := "DROP"
		policy := nsxModel.SecurityPolicy{
			Description: strPtr("original"),
			Tags:        []nsxModel.Tag{{Scope: strPtr("s"), Tag: strPtr("t")}},
			Rules: []nsxModel.Rule{
				{Id: &defaultRuleID, IsDefault: &isDefault, Action: &action},
				{Id: &customRuleID, IsDefault: &notDefault},
			},
		}

		reverted, err := revertPolicyPredefinedSecurityPolicy(policy, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", *reverted.Description)
		assert.Nil(t, reverted.Rules)
		assert.Empty(t, reverted.Tags)
		require.Len(t, reverted.Children, 2)
	})

	t.Run("no rules is a no-op besides clearing description/tags", func(t *testing.T) {
		policy := nsxModel.SecurityPolicy{Description: strPtr("original")}

		reverted, err := revertPolicyPredefinedSecurityPolicy(policy, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", *reverted.Description)
		assert.Empty(t, reverted.Children)
	})
}

func TestMockNsxtSecurityPolicyInfraPatchVPC(t *testing.T) {
	t.Run("VPC context patches through the org root client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrgRootSDK := orgrootmocks.NewMockOrgRootClient(ctrl)
		original := cliOrgRootClient
		cliOrgRootClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apipkg.OrgRootClientContext {
			return &apipkg.OrgRootClientContext{Client: mockOrgRootSDK, ClientType: utl.VPC}
		}
		defer func() { cliOrgRootClient = original }()
		mockOrgRootSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil)

		policyID := "default"
		ctx := utl.SessionContext{ClientType: utl.VPC, ProjectID: "project1", VPCID: "vpc1"}
		policy := nsxModel.SecurityPolicy{Id: &policyID}

		err := securityPolicyInfraPatch(ctx, policy, "default", newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("VPC context propagates a Patch error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrgRootSDK := orgrootmocks.NewMockOrgRootClient(ctrl)
		original := cliOrgRootClient
		cliOrgRootClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *apipkg.OrgRootClientContext {
			return &apipkg.OrgRootClientContext{Client: mockOrgRootSDK, ClientType: utl.VPC}
		}
		defer func() { cliOrgRootClient = original }()
		mockOrgRootSDK.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(vapiErrors.InternalServerError{})

		policyID := "default"
		ctx := utl.SessionContext{ClientType: utl.VPC, ProjectID: "project1", VPCID: "vpc1"}
		policy := nsxModel.SecurityPolicy{Id: &policyID}

		err := securityPolicyInfraPatch(ctx, policy, "default", newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyPredefinedSecurityPolicyUpdateDefaultRuleChange(t *testing.T) {
	t.Run("Update patches a changed default_rule block", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK, mockInfra, restore := setupPredefinedSecPolicyMock(t, ctrl)
		defer restore()

		isDefault := true
		defaultRuleID := "default-rule-1"
		action := "ALLOW"
		respWithDefaultRule := predefinedSecPolicyAPIResponse()
		respWithDefaultRule.Rules = []nsxModel.Rule{
			{Id: &defaultRuleID, IsDefault: &isDefault, Action: &action},
		}

		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(respWithDefaultRule, nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil),
		)

		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		data := minimalPredefinedSecPolicyData()
		data["default_rule"] = []interface{}{
			map[string]interface{}{
				"description": "updated default rule",
				"action":      "DROP",
				"logged":      true,
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		d.SetId(predefinedSecPolicyID)

		err := resourceNsxtPolicyPredefinedSecurityPolicyUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})
}

func TestUnitNsxt_updatePolicyPredefinedSecurityPolicyEmptyDomain(t *testing.T) {
	t.Run("fails when domain cannot be extracted from path", func(t *testing.T) {
		res := resourceNsxtPolicyPredefinedSecurityPolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"path": "/not-a-domain-path",
		})

		err := updatePolicyPredefinedSecurityPolicy("some-id", d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "domain")
	})
}
