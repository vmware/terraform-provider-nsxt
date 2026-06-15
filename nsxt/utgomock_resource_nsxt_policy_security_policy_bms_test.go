//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

func TestMockResourceNsxtPolicySecurityPolicyBMSGroupsSchema(t *testing.T) {
	resource := resourceNsxtPolicySecurityPolicy()

	// Test schema structure
	assert.NotNil(t, resource.Schema)

	// Test rule schema
	ruleSchema := resource.Schema["rule"].Elem.(*schema.Resource).Schema
	assert.Contains(t, ruleSchema, "source_groups")
	assert.Contains(t, ruleSchema, "destination_groups")
	assert.Contains(t, ruleSchema, "scope")
}

func TestMockResourceNsxtPolicySecurityPolicyBMSRuleCreation(t *testing.T) {
	testCases := []struct {
		name         string
		sourceGroups []string
		destGroups   []string
		scope        []string
		action       string
		expectedRule model.Rule
	}{
		{
			name: "BMS to BMS communication rule",
			sourceGroups: []string{
				"/infra/domains/default/groups/production-bms",
			},
			destGroups: []string{
				"/infra/domains/default/groups/production-bms",
			},
			action: "ALLOW",
			expectedRule: model.Rule{
				Action:            &[]string{"ALLOW"}[0],
				SourceGroups:      []string{"/infra/domains/default/groups/production-bms"},
				DestinationGroups: []string{"/infra/domains/default/groups/production-bms"},
			},
		},
		{
			name: "Web traffic to BMS rule",
			destGroups: []string{
				"/infra/domains/default/groups/web-servers-bms",
			},
			scope: []string{
				"/infra/domains/default/groups/bms-data-interfaces",
			},
			action: "ALLOW",
			expectedRule: model.Rule{
				Action:            &[]string{"ALLOW"}[0],
				DestinationGroups: []string{"/infra/domains/default/groups/web-servers-bms"},
				Scope:             []string{"/infra/domains/default/groups/bms-data-interfaces"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rule := createSecurityPolicyRule(tc.sourceGroups, tc.destGroups, tc.scope, tc.action)

			assert.Equal(t, *tc.expectedRule.Action, *rule.Action)

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

func TestMockResourceNsxtPolicySecurityPolicyBMSGroupValidation(t *testing.T) {
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
			name:        "valid VM group path",
			groupPath:   "/infra/domains/default/groups/production-vms",
			expectValid: true,
			groupType:   "vm",
		},
		{
			name:        "invalid group path format",
			groupPath:   "invalid-path",
			expectValid: false,
			groupType:   "unknown",
		},
		{
			name:        "empty group path",
			groupPath:   "",
			expectValid: false,
			groupType:   "unknown",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := validateGroupPath(tc.groupPath)
			groupType := detectGroupType(tc.groupPath)

			assert.Equal(t, tc.expectValid, isValid)
			assert.Equal(t, tc.groupType, groupType)
		})
	}
}

func TestMockResourceNsxtPolicySecurityPolicyBMSMixedGroups(t *testing.T) {
	mixedGroupPaths := []string{
		"/infra/domains/default/groups/production-bms",
		"/infra/domains/default/groups/production-vms",
		"/infra/domains/default/groups/bms-interfaces",
	}

	testCases := []struct {
		name          string
		groups        []string
		expectedBMS   int
		expectedVM    int
		expectedMixed bool
	}{
		{
			name:          "mixed BMS and VM groups",
			groups:        mixedGroupPaths,
			expectedBMS:   2,
			expectedVM:    1,
			expectedMixed: true,
		},
		{
			name: "BMS only groups",
			groups: []string{
				"/infra/domains/default/groups/production-bms",
				"/infra/domains/default/groups/bms-interfaces",
			},
			expectedBMS:   2,
			expectedVM:    0,
			expectedMixed: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analysis := analyzeGroupMix(tc.groups)

			assert.Equal(t, tc.expectedBMS, analysis["bms_groups"])
			assert.Equal(t, tc.expectedVM, analysis["vm_groups"])
			assert.Equal(t, tc.expectedMixed, analysis["is_mixed"])
		})
	}
}

func TestMockResourceNsxtPolicySecurityPolicyBMSRuleValidation(t *testing.T) {
	testCases := []struct {
		name        string
		rule        map[string]interface{}
		expectValid bool
		errorMsg    string
	}{
		{
			name: "valid BMS rule",
			rule: map[string]interface{}{
				"display_name":       "allow_bms_internal",
				"action":             "ALLOW",
				"source_groups":      []string{"/infra/domains/default/groups/production-bms"},
				"destination_groups": []string{"/infra/domains/default/groups/production-bms"},
			},
			expectValid: true,
		},
		{
			name: "rule without action",
			rule: map[string]interface{}{
				"display_name":       "invalid_rule",
				"destination_groups": []string{"/infra/domains/default/groups/production-bms"},
			},
			expectValid: false,
			errorMsg:    "action required",
		},
		{
			name: "rule with invalid group path",
			rule: map[string]interface{}{
				"display_name":       "invalid_group_rule",
				"action":             "ALLOW",
				"destination_groups": []string{"invalid-path"},
			},
			expectValid: false,
			errorMsg:    "invalid group path",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid, errorMsg := validateSecurityPolicyRule(tc.rule)

			assert.Equal(t, tc.expectValid, isValid)
			if !tc.expectValid {
				assert.Contains(t, errorMsg, tc.errorMsg)
			}
		})
	}
}

// Helper functions to simulate actual implementation logic
func createSecurityPolicyRule(sourceGroups, destGroups, scope []string, action string) model.Rule {
	rule := model.Rule{
		Action: &action,
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

func validateGroupPath(groupPath string) bool {
	if groupPath == "" {
		return false
	}

	// Simple validation for group path format
	return len(groupPath) > 0 && groupPath[0] == '/' &&
		(contains(groupPath, "/groups/") || contains(groupPath, "/infra/"))
}

func detectGroupType(groupPath string) string {
	if !validateGroupPath(groupPath) {
		return "unknown"
	}

	if contains(groupPath, "bms") || contains(groupPath, "baremetal") {
		return "bms"
	} else if contains(groupPath, "vm") || contains(groupPath, "virtual") {
		return "vm"
	}

	return "unknown"
}

func analyzeGroupMix(groups []string) map[string]interface{} {
	result := make(map[string]interface{})

	bmsCount := 0
	vmCount := 0

	for _, group := range groups {
		groupType := detectGroupType(group)
		switch groupType {
		case "bms":
			bmsCount++
		case "vm":
			vmCount++
		}
	}

	result["bms_groups"] = bmsCount
	result["vm_groups"] = vmCount
	result["is_mixed"] = bmsCount > 0 && vmCount > 0

	return result
}

func validateSecurityPolicyRule(rule map[string]interface{}) (bool, string) {
	// Check required action
	action, hasAction := rule["action"]
	if !hasAction || action == "" {
		return false, "action required"
	}

	// Validate group paths if present
	if destGroups, ok := rule["destination_groups"]; ok {
		if groups, isSlice := destGroups.([]string); isSlice {
			for _, group := range groups {
				if !validateGroupPath(group) {
					return false, "invalid group path"
				}
			}
		}
	}

	return true, ""
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr ||
		(len(s) > len(substr) && contains(s[1:], substr))
}
