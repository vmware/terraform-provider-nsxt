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
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func boolPtr(b bool) *bool {
	return &b
}

func policyRuleTestResource() *schema.Resource {
	return &schema.Resource{Schema: map[string]*schema.Schema{
		"rule": getSecurityPolicyAndGatewayRulesSchema(false, false, false),
	}}
}

func TestUnitNsxt_getPolicyRulesFromSchema(t *testing.T) {
	res := policyRuleTestResource()

	t.Run("assigns auto-incrementing sequence numbers when unspecified", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"rule": []interface{}{
				map[string]interface{}{"display_name": "rule1"},
				map[string]interface{}{"display_name": "rule2"},
			},
		})
		rules := getPolicyRulesFromSchema(d)
		require.Len(t, rules, 2)
		assert.Equal(t, "rule1", *rules[0].DisplayName)
		assert.EqualValues(t, 1, *rules[0].SequenceNumber)
		assert.Equal(t, "rule2", *rules[1].DisplayName)
		assert.EqualValues(t, 2, *rules[1].SequenceNumber)
	})

	t.Run("out-of-order explicit sequence numbers are corrected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"rule": []interface{}{
				map[string]interface{}{"display_name": "rule1", "sequence_number": 5},
				map[string]interface{}{"display_name": "rule2", "sequence_number": 3},
			},
		})
		rules := getPolicyRulesFromSchema(d)
		require.Len(t, rules, 2)
		assert.EqualValues(t, 5, *rules[0].SequenceNumber)
		assert.EqualValues(t, 6, *rules[1].SequenceNumber)
	})

	t.Run("existing nsx_id is preserved as the rule id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"rule": []interface{}{
				map[string]interface{}{"display_name": "rule1", "nsx_id": "custom-id"},
			},
		})
		rules := getPolicyRulesFromSchema(d)
		require.Len(t, rules, 1)
		assert.Equal(t, "custom-id", *rules[0].Id)
	})

	t.Run("no rules returns empty list", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
		rules := getPolicyRulesFromSchema(d)
		assert.Empty(t, rules)
	})

	t.Run("ip_version NONE leaves IpProtocol nil", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"rule": []interface{}{
				map[string]interface{}{"display_name": "rule1", "ip_version": "NONE"},
			},
		})
		rules := getPolicyRulesFromSchema(d)
		require.Len(t, rules, 1)
		assert.Nil(t, rules[0].IpProtocol)
	})
}

func TestUnitNsxt_setPolicyRulesInSchema(t *testing.T) {
	res := policyRuleTestResource()

	t.Run("populates schema rule fields from model rules", func(t *testing.T) {
		displayName := "rule1"
		description := "a rule"
		path := "/infra/domains/default/security-policies/sp1/rules/rule1"
		notes := "some notes"
		logLabel := "label1"
		revision := int64(3)
		id := "rule1"
		ruleID := int64(42)
		tagScope := "scope1"
		tagTag := "tag1"

		rule := model.Rule{
			DisplayName:          &displayName,
			Description:          &description,
			Path:                 &path,
			Notes:                &notes,
			Logged:               boolPtr(true),
			Tag:                  &logLabel,
			Action:               strPtr(model.Rule_ACTION_ALLOW),
			DestinationsExcluded: boolPtr(true),
			SourcesExcluded:      boolPtr(false),
			IpProtocol:           strPtr(model.Rule_IP_PROTOCOL_IPV4),
			Direction:            strPtr(model.Rule_DIRECTION_IN),
			Disabled:             boolPtr(true),
			Revision:             &revision,
			SourceGroups:         []string{"/infra/domains/default/groups/src1"},
			DestinationGroups:    []string{"/infra/domains/default/groups/dst1"},
			Profiles:             []string{"ANY"},
			Services:             []string{"ANY"},
			Scope:                []string{"ANY"},
			SequenceNumber:       int64Ptr(1),
			Id:                   &id,
			RuleId:               &ruleID,
			Tags:                 []model.Tag{{Scope: &tagScope, Tag: &tagTag}},
		}

		d := res.TestResourceData()
		require.NoError(t, setPolicyRulesInSchema(d, []model.Rule{rule}))

		rules := d.Get("rule").([]interface{})
		require.Len(t, rules, 1)
		elem := rules[0].(map[string]interface{})
		assert.Equal(t, displayName, elem["display_name"])
		assert.Equal(t, description, elem["description"])
		assert.Equal(t, path, elem["path"])
		assert.Equal(t, notes, elem["notes"])
		assert.Equal(t, true, elem["logged"])
		assert.Equal(t, logLabel, elem["log_label"])
		assert.Equal(t, model.Rule_ACTION_ALLOW, elem["action"])
		assert.Equal(t, true, elem["destinations_excluded"])
		assert.Equal(t, false, elem["sources_excluded"])
		assert.Equal(t, model.Rule_IP_PROTOCOL_IPV4, elem["ip_version"])
		assert.Equal(t, model.Rule_DIRECTION_IN, elem["direction"])
		assert.Equal(t, true, elem["disabled"])
		assert.EqualValues(t, revision, elem["revision"])
		assert.EqualValues(t, ruleID, elem["rule_id"])
		assert.Equal(t, id, elem["nsx_id"])
		assert.Contains(t, elem["source_groups"].(*schema.Set).List(), "/infra/domains/default/groups/src1")
		assert.Contains(t, elem["destination_groups"].(*schema.Set).List(), "/infra/domains/default/groups/dst1")
	})

	t.Run("nil IpProtocol maps to NONE", func(t *testing.T) {
		displayName := "rule1"
		rule := model.Rule{DisplayName: &displayName}

		d := res.TestResourceData()
		require.NoError(t, setPolicyRulesInSchema(d, []model.Rule{rule}))

		rules := d.Get("rule").([]interface{})
		require.Len(t, rules, 1)
		assert.Equal(t, "NONE", rules[0].(map[string]interface{})["ip_version"])
	})

	t.Run("no rules clears the schema list", func(t *testing.T) {
		d := res.TestResourceData()
		require.NoError(t, setPolicyRulesInSchema(d, nil))
		assert.Empty(t, d.Get("rule").([]interface{}))
	})
}

func TestUnitNsxt_validatePolicyRuleSequence(t *testing.T) {
	res := policyRuleTestResource()

	t.Run("increasing explicit sequence numbers are valid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"rule": []interface{}{
				map[string]interface{}{"display_name": "rule1", "sequence_number": 1},
				map[string]interface{}{"display_name": "rule2", "sequence_number": 2},
			},
		})
		assert.NoError(t, validatePolicyRuleSequence(d))
	})

	t.Run("unspecified sequence numbers are valid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"rule": []interface{}{
				map[string]interface{}{"display_name": "rule1"},
				map[string]interface{}{"display_name": "rule2"},
			},
		})
		assert.NoError(t, validatePolicyRuleSequence(d))
	})

	t.Run("out-of-order explicit sequence numbers are rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
			"rule": []interface{}{
				map[string]interface{}{"display_name": "rule1", "sequence_number": 5},
				map[string]interface{}{"display_name": "rule2", "sequence_number": 5},
			},
		})
		err := validatePolicyRuleSequence(d)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "rule2")
	})
}
