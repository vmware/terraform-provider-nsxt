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
)

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
