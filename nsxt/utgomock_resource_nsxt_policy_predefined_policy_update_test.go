//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vapiErrors "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"go.uber.org/mock/gomock"
)

// collectChildRules walks the H-API tree sent to Infra.Patch and returns the
// ChildRule entries it carries, keyed by rule ID with their marked_for_delete flag.
func collectChildRules(t *testing.T, infra nsxModel.Infra) map[string]bool {
	t.Helper()
	result := map[string]bool{}
	var walk func(v data.DataValue)
	walk = func(v data.DataValue) {
		switch val := v.(type) {
		case *data.StructValue:
			fields := val.Fields()
			if rt, ok := fields["resource_type"].(*data.StringValue); ok && rt.Value() == "ChildRule" {
				id := fields["id"].(*data.OptionalValue).Value().(*data.StringValue).Value()
				marked := false
				if opt := fields["marked_for_delete"].(*data.OptionalValue); opt.IsSet() {
					marked = opt.Value().(*data.BooleanValue).Value()
				}
				result[id] = marked
				return
			}
			for _, f := range fields {
				walk(f)
			}
		case *data.ListValue:
			for _, item := range val.List() {
				walk(item)
			}
		case *data.OptionalValue:
			if val.IsSet() {
				walk(val.Value())
			}
		}
	}
	for _, child := range infra.Children {
		walk(child)
	}
	return result
}

// resourceDataWithChange returns ResourceData carrying oldRaw as state and a
// real diff towards newRaw, so HasChange/GetChange behave as during an
// actual update (d.Set alone does not register as a change).
func resourceDataWithChange(t *testing.T, res *schema.Resource, oldRaw, newRaw map[string]interface{}, id string) *schema.ResourceData {
	t.Helper()
	old := schema.TestResourceDataRaw(t, res.Schema, oldRaw)
	old.SetId(id)
	state := old.State()

	diff, err := res.Diff(context.Background(), state, terraform.NewResourceConfigRaw(newRaw), nil)
	require.NoError(t, err)
	require.NotNil(t, diff)

	var captured *schema.ResourceData
	capture := &schema.Resource{
		Schema: res.Schema,
		Update: func(d *schema.ResourceData, _ interface{}) error {
			captured = d
			return nil
		},
	}
	_, diags := capture.Apply(context.Background(), state, diff, nil)
	require.False(t, diags.HasError(), "%v", diags)
	require.NotNil(t, captured)
	return captured
}

// predefinedPolicyDataWithRuleChange returns ResourceData whose state holds
// rules keep-rule and stale-rule, while config holds keep-rule and new-rule
// along with a changed description and tags.
func predefinedPolicyDataWithRuleChange(t *testing.T, res *schema.Resource, base map[string]interface{}, id string) *schema.ResourceData {
	rule := func(nsxID string) map[string]interface{} {
		return map[string]interface{}{
			"nsx_id":       nsxID,
			"display_name": nsxID,
			"scope":        []interface{}{"/infra/tier-1s/t1"},
			"action":       nsxModel.Rule_ACTION_ALLOW,
		}
	}
	oldRaw := map[string]interface{}{}
	newRaw := map[string]interface{}{}
	for k, v := range base {
		oldRaw[k] = v
		newRaw[k] = v
	}
	oldRaw["rule"] = []interface{}{rule("keep-rule"), rule("stale-rule")}
	newRaw["rule"] = []interface{}{rule("keep-rule"), rule("new-rule")}
	newRaw["description"] = "updated description"
	newRaw["tag"] = []interface{}{map[string]interface{}{"scope": "env", "tag": "prod"}}
	return resourceDataWithChange(t, res, oldRaw, newRaw, id)
}

func TestMockUpdatePolicyPredefinedSecurityPolicyNsxtRules(t *testing.T) {
	res := resourceNsxtPolicyPredefinedSecurityPolicy()

	t.Run("adds configured rules and deletes rules removed from config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockSDK, mockInfra, restore := setupPredefinedSecPolicyMock(t, ctrl)
		defer restore()

		var patched nsxModel.Infra
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).DoAndReturn(func(infra nsxModel.Infra, _ *bool) error {
				patched = infra
				return nil
			}),
		)

		d := predefinedPolicyDataWithRuleChange(t, res, minimalPredefinedSecPolicyData(), predefinedSecPolicyID)
		require.NoError(t, updatePolicyPredefinedSecurityPolicy(predefinedSecPolicyID, d, newGoMockProviderClient()))

		assert.Equal(t, map[string]bool{"keep-rule": false, "new-rule": false, "stale-rule": true}, collectChildRules(t, patched))
	})

	t.Run("Get error is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockSDK, _, restore := setupPredefinedSecPolicyMock(t, ctrl)
		defer restore()
		mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(nsxModel.SecurityPolicy{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())
		require.Error(t, updatePolicyPredefinedSecurityPolicy(predefinedSecPolicyID, d, newGoMockProviderClient()))
	})

	t.Run("Patch error is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockSDK, mockInfra, restore := setupPredefinedSecPolicyMock(t, ctrl)
		defer restore()
		mockSDK.EXPECT().Get(predefinedSecPolicyDomain, predefinedSecPolicyID).Return(predefinedSecPolicyAPIResponse(), nil)
		mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).Return(vapiErrors.InvalidRequest{})

		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedSecPolicyData())
		err := updatePolicyPredefinedSecurityPolicy(predefinedSecPolicyID, d, newGoMockProviderClient())
		require.Error(t, err)
	})
}

func TestMockUpdatePolicyPredefinedGatewayPolicyNsxtRules(t *testing.T) {
	res := resourceNsxtPolicyPredefinedGatewayPolicy()

	t.Run("adds configured rules and deletes rules removed from config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockSDK, mockInfra, restore := setupPredefinedGwPolicyMock(t, ctrl)
		defer restore()

		var patched nsxModel.Infra
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(predefinedGwPolicyAPIResponse(), nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).DoAndReturn(func(infra nsxModel.Infra, _ *bool) error {
				patched = infra
				return nil
			}),
		)

		d := predefinedPolicyDataWithRuleChange(t, res, minimalPredefinedGwPolicyData(), predefinedGwPolicyID)
		require.NoError(t, updatePolicyPredefinedGatewayPolicy(predefinedGwPolicyID, d, newGoMockVMCProviderClient()))

		assert.Equal(t, map[string]bool{"keep-rule": false, "new-rule": false, "stale-rule": true}, collectChildRules(t, patched))
	})

	t.Run("updates default rules matching every existing default scope", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockSDK, mockInfra, restore := setupPredefinedGwPolicyMock(t, ctrl)
		defer restore()

		isDefault := true
		defaultRuleID := "default-rule-t1"
		otherRuleID := "regular-rule"
		policy := predefinedGwPolicyAPIResponse()
		policy.Rules = []nsxModel.Rule{
			{Id: &defaultRuleID, IsDefault: &isDefault, Scope: []string{"/infra/tier-1s/t1"}},
			{Id: &otherRuleID, Scope: []string{"/infra/tier-1s/t1"}},
		}

		var patched nsxModel.Infra
		gomock.InOrder(
			mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(policy, nil),
			mockInfra.EXPECT().Patch(gomock.Any(), gomock.Any()).DoAndReturn(func(infra nsxModel.Infra, _ *bool) error {
				patched = infra
				return nil
			}),
		)

		data := minimalPredefinedGwPolicyData()
		data["default_rule"] = []interface{}{
			map[string]interface{}{
				"scope":  "/infra/tier-1s/t1",
				"action": nsxModel.Rule_ACTION_DROP,
			},
		}
		d := schema.TestResourceDataRaw(t, res.Schema, data)
		require.NoError(t, updatePolicyPredefinedGatewayPolicy(predefinedGwPolicyID, d, newGoMockProviderClient()))

		assert.Equal(t, map[string]bool{defaultRuleID: false}, collectChildRules(t, patched))
	})

	t.Run("Get error is returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockSDK, _, restore := setupPredefinedGwPolicyMock(t, ctrl)
		defer restore()
		mockSDK.EXPECT().Get(predefinedGwPolicyDomain, predefinedGwPolicyID).Return(nsxModel.GatewayPolicy{}, vapiErrors.NotFound{})

		d := schema.TestResourceDataRaw(t, res.Schema, minimalPredefinedGwPolicyData())
		require.Error(t, updatePolicyPredefinedGatewayPolicy(predefinedGwPolicyID, d, newGoMockProviderClient()))
	})
}
