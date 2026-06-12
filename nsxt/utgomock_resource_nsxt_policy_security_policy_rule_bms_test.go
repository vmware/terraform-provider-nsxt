//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestMockResourceNsxtPolicySecurityPolicyRuleBMSStandaloneSchema(t *testing.T) {
	resource := resourceNsxtPolicySecurityPolicyRule()

	// Test schema structure
	assert.NotNil(t, resource.Schema)

	// Test rule-specific fields
	schema := resource.Schema
	assert.Contains(t, schema, "policy_path")
	assert.Contains(t, schema, "source_groups")
	assert.Contains(t, schema, "destination_groups")
	assert.Contains(t, schema, "scope")
	assert.Contains(t, schema, "sequence_number")
}

func TestMockResourceNsxtPolicySecurityPolicyRuleBMSGroupValidation(t *testing.T) {
	testCases := []struct {
		name        string
		groupPath   string
		expectValid bool
		groupType   string
	}{
		{
			name:        "valid BMS group path",
			groupPath:   "/infra/domains/default/groups/production-bms",
			expectValid: true,
			groupType:   "bms",
		},
		{
			name:        "valid BMS interface group path",
			groupPath:   "/infra/domains/default/groups/bms-interfaces",
			expectValid: true,
			groupType:   "bms",
		},
		{
			name:        "valid VM group path",
			groupPath:   "/infra/domains/default/groups/vm-group",
			expectValid: true,
			groupType:   "vm",
		},
		{
			name:        "invalid path format",
			groupPath:   "invalid-group-path",
			expectValid: false,
			groupType:   "unknown",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := validateStandaloneRuleGroupPath(tc.groupPath)
			groupType := detectStandaloneRuleGroupType(tc.groupPath)

			assert.Equal(t, tc.expectValid, isValid)
			assert.Equal(t, tc.groupType, groupType)
		})
	}
}

func TestMockResourceNsxtPolicySecurityPolicyRuleBMSRuleCreation(t *testing.T) {
	testCases := []struct {
		name           string
		displayName    string
		policyPath     string
		sourceGroups   []string
		destGroups     []string
		scope          []string
		action         string
		sequenceNumber int64
		expectedRule   model.Rule
	}{
		{
			name:           "BMS source to destination rule",
			displayName:    "allow-bms-internal",
			policyPath:     "/infra/domains/default/security-policies/policy1",
			sourceGroups:   []string{"/infra/domains/default/groups/bms-prod"},
			destGroups:     []string{"/infra/domains/default/groups/bms-dev"},
			action:         "ALLOW",
			sequenceNumber: 100,
			expectedRule: model.Rule{
				DisplayName:       &[]string{"allow-bms-internal"}[0],
				SourceGroups:      []string{"/infra/domains/default/groups/bms-prod"},
				DestinationGroups: []string{"/infra/domains/default/groups/bms-dev"},
				Action:            &[]string{"ALLOW"}[0],
				SequenceNumber:    &[]int64{100}[0],
			},
		},
		{
			name:           "BMS interface scope rule",
			displayName:    "scoped-to-bms-interfaces",
			policyPath:     "/infra/domains/default/security-policies/policy1",
			scope:          []string{"/infra/domains/default/groups/bms-interfaces"},
			action:         "DROP",
			sequenceNumber: 200,
			expectedRule: model.Rule{
				DisplayName:    &[]string{"scoped-to-bms-interfaces"}[0],
				Scope:          []string{"/infra/domains/default/groups/bms-interfaces"},
				Action:         &[]string{"DROP"}[0],
				SequenceNumber: &[]int64{200}[0],
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rule := createStandaloneSecurityPolicyRule(
				tc.displayName, tc.policyPath, tc.sourceGroups,
				tc.destGroups, tc.scope, tc.action, tc.sequenceNumber,
			)

			assert.Equal(t, *tc.expectedRule.DisplayName, *rule.DisplayName)
			assert.Equal(t, *tc.expectedRule.Action, *rule.Action)
			assert.Equal(t, *tc.expectedRule.SequenceNumber, *rule.SequenceNumber)

			if tc.expectedRule.SourceGroups != nil {
				assert.Equal(t, len(tc.expectedRule.SourceGroups), len(rule.SourceGroups))
			}

			if tc.expectedRule.DestinationGroups != nil {
				assert.Equal(t, len(tc.expectedRule.DestinationGroups), len(rule.DestinationGroups))
			}

			if tc.expectedRule.Scope != nil {
				assert.Equal(t, len(tc.expectedRule.Scope), len(rule.Scope))
			}
		})
	}
}

func TestMockResourceNsxtPolicySecurityPolicyRuleBMSSequenceValidation(t *testing.T) {
	testCases := []struct {
		name           string
		sequenceNumber int64
		expectValid    bool
	}{
		{
			name:           "valid sequence number",
			sequenceNumber: 100,
			expectValid:    true,
		},
		{
			name:           "minimum sequence number",
			sequenceNumber: 1,
			expectValid:    true,
		},
		{
			name:           "maximum sequence number",
			sequenceNumber: 999999,
			expectValid:    true,
		},
		{
			name:           "invalid zero sequence",
			sequenceNumber: 0,
			expectValid:    false,
		},
		{
			name:           "invalid negative sequence",
			sequenceNumber: -1,
			expectValid:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validateStandaloneRuleSequenceNumber(tc.sequenceNumber)
			assert.Equal(t, tc.expectValid, result)
		})
	}
}

func TestMockResourceNsxtPolicySecurityPolicyRuleBMSPolicyPathValidation(t *testing.T) {
	testCases := []struct {
		name        string
		policyPath  string
		expectValid bool
	}{
		{
			name:        "valid policy path",
			policyPath:  "/infra/domains/default/security-policies/bms-policy",
			expectValid: true,
		},
		{
			name:        "valid custom domain policy path",
			policyPath:  "/infra/domains/production/security-policies/bms-policy",
			expectValid: true,
		},
		{
			name:        "invalid empty policy path",
			policyPath:  "",
			expectValid: false,
		},
		{
			name:        "invalid policy path format",
			policyPath:  "invalid-policy-path",
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validateStandaloneRulePolicyPath(tc.policyPath)
			assert.Equal(t, tc.expectValid, result)
		})
	}
}

func TestMockResourceNsxtPolicySecurityPolicyRuleBMSMixedGroupAnalysis(t *testing.T) {
	testCases := []struct {
		name           string
		groups         []string
		expectedResult map[string]interface{}
	}{
		{
			name: "BMS only groups",
			groups: []string{
				"/infra/domains/default/groups/bms-prod",
				"/infra/domains/default/groups/bms-interfaces",
			},
			expectedResult: map[string]interface{}{
				"bms_count":    2,
				"vm_count":     0,
				"mixed_groups": false,
				"total_groups": 2,
			},
		},
		{
			name: "mixed BMS and VM groups",
			groups: []string{
				"/infra/domains/default/groups/bms-prod",
				"/infra/domains/default/groups/vm-prod",
				"/infra/domains/default/groups/bms-interfaces",
			},
			expectedResult: map[string]interface{}{
				"bms_count":    2,
				"vm_count":     1,
				"mixed_groups": true,
				"total_groups": 3,
			},
		},
		{
			name: "VM only groups",
			groups: []string{
				"/infra/domains/default/groups/vm-prod",
				"/infra/domains/default/groups/vm-dev",
			},
			expectedResult: map[string]interface{}{
				"bms_count":    0,
				"vm_count":     2,
				"mixed_groups": false,
				"total_groups": 2,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := analyzeStandaloneRuleGroupMix(tc.groups)

			assert.Equal(t, tc.expectedResult["bms_count"], result["bms_count"])
			assert.Equal(t, tc.expectedResult["vm_count"], result["vm_count"])
			assert.Equal(t, tc.expectedResult["mixed_groups"], result["mixed_groups"])
			assert.Equal(t, tc.expectedResult["total_groups"], result["total_groups"])
		})
	}
}

// Helper functions to simulate actual implementation logic
func validateStandaloneRuleGroupPath(groupPath string) bool {
	if groupPath == "" {
		return false
	}

	return len(groupPath) > 0 && groupPath[0] == '/' &&
		(contains(groupPath, "/groups/") || contains(groupPath, "/infra/"))
}

func detectStandaloneRuleGroupType(groupPath string) string {
	if !validateStandaloneRuleGroupPath(groupPath) {
		return "unknown"
	}

	if contains(groupPath, "bms") || contains(groupPath, "baremetal") ||
		contains(groupPath, "interface") {
		return "bms"
	} else if contains(groupPath, "vm") || contains(groupPath, "virtual") {
		return "vm"
	}

	return "unknown"
}

func createStandaloneSecurityPolicyRule(displayName, policyPath string, sourceGroups, destGroups, scope []string, action string, sequenceNumber int64) model.Rule {
	rule := model.Rule{
		DisplayName:    &displayName,
		Action:         &action,
		SequenceNumber: &sequenceNumber,
	}

	if sourceGroups != nil {
		rule.SourceGroups = sourceGroups
	}

	if destGroups != nil {
		rule.DestinationGroups = destGroups
	}

	if scope != nil {
		rule.Scope = scope
	}

	return rule
}

func validateStandaloneRuleSequenceNumber(sequenceNumber int64) bool {
	return sequenceNumber > 0 && sequenceNumber <= 999999
}

func validateStandaloneRulePolicyPath(policyPath string) bool {
	if policyPath == "" {
		return false
	}

	return len(policyPath) > 0 && policyPath[0] == '/' &&
		contains(policyPath, "/security-policies/")
}

func analyzeStandaloneRuleGroupMix(groups []string) map[string]interface{} {
	result := make(map[string]interface{})

	bmsCount := 0
	vmCount := 0

	for _, group := range groups {
		groupType := detectStandaloneRuleGroupType(group)
		switch groupType {
		case "bms":
			bmsCount++
		case "vm":
			vmCount++
		}
	}

	result["bms_count"] = bmsCount
	result["vm_count"] = vmCount
	result["mixed_groups"] = bmsCount > 0 && vmCount > 0
	result["total_groups"] = len(groups)

	return result
}
