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

	intrusionservicegatewaypolicies "github.com/vmware/terraform-provider-nsxt/api/infra/domains/intrusion_service_gateway_policies"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	isgwrulemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains/intrusion_service_gateway_policies"
)

var (
	idsGwRuleID          = "ids-gw-rule-1"
	idsGwRuleDisplayName = "Test IDS Gateway Policy Rule"
	idsGwRuleDescription = "Test IDS gateway policy rule"
	idsGwRuleRevision    = int64(1)
	idsGwRuleDomain      = "default"
	idsGwRulePolicyID    = "ids-gw-policy-001"
	idsGwRulePolicyPath  = "/infra/domains/default/ids-gateway-policies/ids-gw-policy-001"
	idsGwRuleAction      = "DETECT"
	idsGwRuleDirection   = "IN_OUT"
	idsGwRuleIPVersion   = "IPV4_IPV6"
	idsGwRuleScopePath   = "/infra/domains/default/gateway-policies/gw-policy-001"
	idsGwRuleProfilePath = "/infra/domains/default/ids-profiles/default"
)

func idsGwRuleAPIResponse() nsxModel.IdsRule {
	resourceType := "IdsRule"
	seq := int64(1)
	logged := false
	disabled := false
	srcEx := false
	dstEx := false
	return nsxModel.IdsRule{
		Id:                   &idsGwRuleID,
		DisplayName:          &idsGwRuleDisplayName,
		Description:          &idsGwRuleDescription,
		Revision:             &idsGwRuleRevision,
		ResourceType:         &resourceType,
		Action:               &idsGwRuleAction,
		Direction:            &idsGwRuleDirection,
		IpProtocol:           &idsGwRuleIPVersion,
		SequenceNumber:       &seq,
		Logged:               &logged,
		Disabled:             &disabled,
		SourcesExcluded:      &srcEx,
		DestinationsExcluded: &dstEx,
		Scope:                []string{idsGwRuleScopePath},
		IdsProfiles:          []string{idsGwRuleProfilePath},
	}
}

func minimalIdsGwRuleData() map[string]interface{} {
	return map[string]interface{}{
		"display_name":    idsGwRuleDisplayName,
		"description":     idsGwRuleDescription,
		"policy_path":     idsGwRulePolicyPath,
		"action":          idsGwRuleAction,
		"direction":       idsGwRuleDirection,
		"ip_version":      idsGwRuleIPVersion,
		"logged":          false,
		"disabled":        false,
		"sequence_number": 1,
		"scope":           []interface{}{idsGwRuleScopePath},
		"ids_profiles":    []interface{}{idsGwRuleProfilePath},
	}
}

func setupIdsGwRuleMock(t *testing.T, ctrl *gomock.Controller) (*isgwrulemocks.MockRulesClient, func()) {
	t.Helper()
	mockSDK := isgwrulemocks.NewMockRulesClient(ctrl)
	mockWrapper := &intrusionservicegatewaypolicies.IntrusionServiceGatewayRuleClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}

	original := cliIntrusionServiceGatewayPolicyRulesClient
	cliIntrusionServiceGatewayPolicyRulesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *intrusionservicegatewaypolicies.IntrusionServiceGatewayRuleClientContext {
		return mockWrapper
	}
	return mockSDK, func() { cliIntrusionServiceGatewayPolicyRulesClient = original }
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyRuleCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIdsGwRuleMock(t, ctrl)
	defer restore()

	t.Run("Create success with auto-generated ID", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(idsGwRuleDomain, idsGwRulePolicyID, gomock.Any(), gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(idsGwRuleDomain, idsGwRulePolicyID, gomock.Any()).Return(idsGwRuleAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleCreate(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.NotEmpty(t, d.Id())
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIdsGwRuleMock(t, ctrl)
	defer restore()

	t.Run("Read success", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID).Return(idsGwRuleAPIResponse(), nil)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		d.SetId(idsGwRuleID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsGwRuleDisplayName, d.Get("display_name"))
		assert.Equal(t, idsGwRuleAction, d.Get("action"))
	})

	t.Run("Read not found clears ID", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID).Return(nsxModel.IdsRule{}, vapiErrors.NotFound{})

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		d.SetId(idsGwRuleID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, "", d.Id())
	})

	t.Run("Read fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyRuleUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIdsGwRuleMock(t, ctrl)
	defer restore()

	t.Run("Update success", func(t *testing.T) {
		gomock.InOrder(
			mockSDK.EXPECT().Patch(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID, gomock.Any()).Return(nil),
			mockSDK.EXPECT().Get(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID).Return(idsGwRuleAPIResponse(), nil),
		)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		d.SetId(idsGwRuleID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleUpdate(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Update fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleUpdate(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockResourceNsxtPolicyIntrusionServiceGatewayPolicyRuleDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIdsGwRuleMock(t, ctrl)
	defer restore()

	t.Run("Delete success with nsx_id set", func(t *testing.T) {
		mockSDK.EXPECT().Delete(idsGwRuleDomain, idsGwRulePolicyID, idsGwRuleID).Return(nil)

		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())
		d.SetId(idsGwRuleID)
		d.Set("nsx_id", idsGwRuleID)

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleDelete(d, newGoMockProviderClient())
		require.NoError(t, err)
	})

	t.Run("Delete fails when ID is empty", func(t *testing.T) {
		res := resourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, res.Schema, minimalIdsGwRuleData())

		err := resourceNsxtPolicyIntrusionServiceGatewayPolicyRuleDelete(d, newGoMockProviderClient())
		require.Error(t, err)
	})
}
