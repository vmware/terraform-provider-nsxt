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

	gwpolapi "github.com/vmware/terraform-provider-nsxt/api/infra/domains/gateway_policies"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	gwpolmocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains/gateway_policies"
)

var (
	gwRuleID          = "gw-rule-1"
	gwRuleDisplayName = "Test GW Policy Rule"
	gwRuleDescription = "Test gateway policy rule"
	gwRuleRevision    = int64(1)
	gwRuleDomain      = "default"
	gwRulePolicyID    = "gw-policy-001"
	gwRulePolicyPath  = "/infra/domains/default/gateway-policies/gw-policy-001"
	gwRuleAction      = "ALLOW"
	gwRuleDirection   = "IN_OUT"
	gwRuleIPVersion   = "IPV4_IPV6"
)

func gwRuleAPIResponse() nsxModel.Rule {
	resourceType := "Rule"
	return nsxModel.Rule{
		Id:           &gwRuleID,
		DisplayName:  &gwRuleDisplayName,
		Description:  &gwRuleDescription,
		Revision:     &gwRuleRevision,
		ResourceType: &resourceType,
		Action:       &gwRuleAction,
		Direction:    &gwRuleDirection,
	}
}

func minimalGwRuleData() map[string]interface{} {
	return map[string]interface{}{
		"display_name": gwRuleDisplayName,
		"description":  gwRuleDescription,
		"policy_path":  gwRulePolicyPath,
		"action":       gwRuleAction,
		"direction":    gwRuleDirection,
		"ip_version":   gwRuleIPVersion,
		"logged":       false,
		"disabled":     false,
	}
}

func setupGwRuleMock(t *testing.T, ctrl *gomock.Controller) (*gwpolmocks.MockRulesClient, func()) {
	mockSDK := gwpolmocks.NewMockRulesClient(ctrl)
	mockWrapper := &gwpolapi.RuleClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliGatewayPoliciesRulesClient
	cliGatewayPoliciesRulesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *gwpolapi.RuleClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliGatewayPoliciesRulesClient = original }
}

func TestMockResourceNsxtPolicyGatewayPolicyRuleCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGwRuleMock(t, ctrl)
	defer restore()

	t.Run("Create success with auto-generated ID", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(gwRuleDomain, gwRulePolicyID, gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gwRuleDomain, gwRulePolicyID, gomock.Any()).Return(gwRuleAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGatewayRulePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwRuleData())

		err := resourceNsxtPolicyGatewayPolicyRuleCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyGatewayPolicyRuleRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGwRuleMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(gwRuleDomain, gwRulePolicyID, gwRuleID).Return(gwRuleAPIResponse(), nil)

		res := resourceNsxtPolicyGatewayRulePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwRuleData())
		d.SetId(gwRuleID)

		err := resourceNsxtPolicyGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, gwRuleDisplayName, d.Get("display_name"))
		assert.Equal(t, gwRuleAction, d.Get("action"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(gwRuleDomain, gwRulePolicyID, gwRuleID).Return(nsxModel.Rule{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyGatewayRulePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwRuleData())
		d.SetId(gwRuleID)

		err := resourceNsxtPolicyGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayRulePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwRuleData())

		err := resourceNsxtPolicyGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayPolicyRuleUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGwRuleMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(gwRuleDomain, gwRulePolicyID, gwRuleID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(gwRuleDomain, gwRulePolicyID, gwRuleID).Return(gwRuleAPIResponse(), nil),
		)

		res := resourceNsxtPolicyGatewayRulePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwRuleData())
		d.SetId(gwRuleID)

		err := resourceNsxtPolicyGatewayPolicyRuleUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayRulePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwRuleData())

		err := resourceNsxtPolicyGatewayPolicyRuleUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyGatewayPolicyRuleDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupGwRuleMock(t, ctrl)
	defer restore()

	t.Run("Delete success with nsx_id set", func(t *testing.T) {
		mockSDK.EXPECT().Delete(gwRuleDomain, gwRulePolicyID, gwRuleID).Return(nil)

		res := resourceNsxtPolicyGatewayRulePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwRuleData())
		d.SetId(gwRuleID)
		d.Set("nsx_id", gwRuleID)

		err := resourceNsxtPolicyGatewayPolicyRuleDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when nsx_id is empty", func(t *testing.T) {
		res := resourceNsxtPolicyGatewayRulePolicy()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalGwRuleData())
		d.SetId(gwRuleID)

		err := resourceNsxtPolicyGatewayPolicyRuleDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
