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

	intrusionservicegatewaypolicies "github.com/vmware/terraform-provider-nsxt/api/infra/domains/intrusion_service_gateway_policies"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
	idsgwrulemocks "github.com/vmware/terraform-provider-nsxt/mocks/infra/domains/intrusion_service_gateway_policies"
)

var (
	idsGwRuleDSDomain     = "default"
	idsGwRuleDSPolicyID   = "ids-gw-pol-1"
	idsGwRuleDSID         = "ids-gw-rule-1"
	idsGwRuleDSName       = "ids-gw-rule"
	idsGwRuleDSPath       = "/infra/domains/default/intrusion-service-gateway-policies/ids-gw-pol-1/rules/ids-gw-rule-1"
	idsGwRuleDSPolicyPath = "/infra/domains/default/intrusion-service-gateway-policies/ids-gw-pol-1"
	idsGwRuleDSAction     = "DETECT"
	idsGwRuleDSDesc       = "gateway rule desc"
)

func idsGwRuleModel() nsxModel.IdsRule {
	logged := true
	disabled := false
	sequenceNumber := int64(1)
	ruleId := int64(67890)
	return nsxModel.IdsRule{
		Id:             &idsGwRuleDSID,
		DisplayName:    &idsGwRuleDSName,
		Description:    &idsGwRuleDSDesc,
		Path:           &idsGwRuleDSPath,
		Action:         &idsGwRuleDSAction,
		Logged:         &logged,
		Disabled:       &disabled,
		SequenceNumber: &sequenceNumber,
		RuleId:         &ruleId,
	}
}

func setupIntrusionServiceGwPolicyRuleDataSourceMock(t *testing.T, ctrl *gomock.Controller) (*idsgwrulemocks.MockRulesClient, func()) {
	t.Helper()
	mockSDK := idsgwrulemocks.NewMockRulesClient(ctrl)
	wrapper := &intrusionservicegatewaypolicies.IntrusionServiceGatewayRuleClientContext{
		Client:     mockSDK,
		ClientType: utl.Local,
	}
	orig := cliIntrusionServiceGatewayPolicyRulesClient
	cliIntrusionServiceGatewayPolicyRulesClient = func(_ utl.SessionContext, _ vapiProtocolClient.Connector) *intrusionservicegatewaypolicies.IntrusionServiceGatewayRuleClientContext {
		return wrapper
	}
	return mockSDK, func() { cliIntrusionServiceGatewayPolicyRulesClient = orig }
}

func TestMockDataSourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK, restore := setupIntrusionServiceGwPolicyRuleDataSourceMock(t, ctrl)
	defer restore()

	boolFalse := false

	t.Run("by id success", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwRuleDSDomain, idsGwRuleDSPolicyID, idsGwRuleDSID).Return(idsGwRuleModel(), nil)

		ds := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":          idsGwRuleDSID,
			"policy_path": idsGwRuleDSPolicyPath,
		})

		err := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsGwRuleDSID, d.Id())
		assert.Equal(t, idsGwRuleDSName, d.Get("display_name"))
		assert.Equal(t, idsGwRuleDSAction, d.Get("action"))
		assert.Equal(t, true, d.Get("logged"))
		assert.Equal(t, false, d.Get("disabled"))
		assert.Equal(t, 1, d.Get("sequence_number"))
		assert.Equal(t, 67890, d.Get("rule_id"))
	})

	t.Run("by id not found", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwRuleDSDomain, idsGwRuleDSPolicyID, idsGwRuleDSID).Return(nsxModel.IdsRule{}, vapiErrors.NotFound{})

		ds := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":          idsGwRuleDSID,
			"policy_path": idsGwRuleDSPolicyPath,
		})

		err := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "was not found")
	})

	t.Run("by id API error", func(t *testing.T) {
		mockSDK.EXPECT().Get(idsGwRuleDSDomain, idsGwRuleDSPolicyID, idsGwRuleDSID).Return(nsxModel.IdsRule{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id":          idsGwRuleDSID,
			"policy_path": idsGwRuleDSPolicyPath,
		})

		err := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error while reading")
	})

	t.Run("missing policy_path", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"id": idsGwRuleDSID,
		})

		err := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "policy_path must be specified")
	})

	t.Run("missing id and display_name", func(t *testing.T) {
		ds := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"policy_path": idsGwRuleDSPolicyPath,
		})

		err := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "must be specified")
	})

	t.Run("by display_name via list", func(t *testing.T) {
		rc := int64(1)
		mockSDK.EXPECT().List(idsGwRuleDSDomain, idsGwRuleDSPolicyID, nil, &boolFalse, nil, nil, nil, nil).Return(nsxModel.IdsRuleListResult{
			Results:     []nsxModel.IdsRule{idsGwRuleModel()},
			ResultCount: &rc,
		}, nil)

		ds := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"policy_path":  idsGwRuleDSPolicyPath,
			"display_name": idsGwRuleDSName,
		})

		err := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.NoError(t, err)
		assert.Equal(t, idsGwRuleDSID, d.Id())
		assert.Equal(t, idsGwRuleDSName, d.Get("display_name"))
		assert.Equal(t, idsGwRuleDSAction, d.Get("action"))
	})

	t.Run("list error", func(t *testing.T) {
		mockSDK.EXPECT().List(idsGwRuleDSDomain, idsGwRuleDSPolicyID, nil, &boolFalse, nil, nil, nil, nil).Return(nsxModel.IdsRuleListResult{}, vapiErrors.InternalServerError{})

		ds := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"policy_path":  idsGwRuleDSPolicyPath,
			"display_name": idsGwRuleDSName,
		})

		err := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
	})

	t.Run("multiple matches error", func(t *testing.T) {
		rc := int64(2)
		rule2 := idsGwRuleModel()
		id2 := "ids-gw-rule-2"
		rule2.Id = &id2

		mockSDK.EXPECT().List(idsGwRuleDSDomain, idsGwRuleDSPolicyID, nil, &boolFalse, nil, nil, nil, nil).Return(nsxModel.IdsRuleListResult{
			Results:     []nsxModel.IdsRule{idsGwRuleModel(), rule2},
			ResultCount: &rc,
		}, nil)

		ds := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRule()
		d := schema.TestResourceDataRaw(t, ds.Schema, map[string]interface{}{
			"policy_path":  idsGwRuleDSPolicyPath,
			"display_name": idsGwRuleDSName,
		})

		err := dataSourceNsxtPolicyIntrusionServiceGatewayPolicyRuleRead(d, newGoMockProviderClient())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Found multiple")
	})
}
