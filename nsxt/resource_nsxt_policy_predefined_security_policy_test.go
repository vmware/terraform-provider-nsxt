/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceNsxtPolicyPredefinedSecurityPolicy_basic(t *testing.T) {
	testResourceName := "nsxt_policy_predefined_security_policy.test"
	description1 := "test 1"
	description2 := "test 2"
	tags := `tag {
            scope = "color"
            tag   = "orange"
        }`

	// NOTE: These tests cannot be parallel, as they modify same default policy
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyPredefinedSecurityPolicyBasic(description1, tags),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "description", description1),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyPredefinedSecurityPolicyBasic(description2, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "description", description2),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyPredefinedSecurityPolicy_defaultRule(t *testing.T) {
	testResourceName := "nsxt_policy_predefined_security_policy.test"
	action1 := "DROP"
	action2 := "ALLOW"
	description1 := "test 1"
	description2 := "test 2"
	tags := `tag {
            scope = "color"
            tag   = "orange"
        }`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyPredefinedSecurityPolicyDefaultRule(description1, action1, action1, tags),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.0.action", action1),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.0.description", description1),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.0.log_label", action1),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.0.logged", "true"),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.0.tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "default_rule.0.revision"),
				),
			},
			{
				Config: testAccNsxtPolicyPredefinedSecurityPolicyDefaultRule(description2, action2, action2, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.0.action", action2),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.0.description", description2),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.0.logged", "true"),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.0.log_label", action2),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.0.tag.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "default_rule.0.revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyPredefinedSecurityPolicy_rules(t *testing.T) {
	testResourceName := "nsxt_policy_predefined_security_policy.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyPredefinedSecurityPolicyWithRules(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "rule2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", "ALLOW"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.source_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.log_label", "group2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.disabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.0.action", "ALLOW"),
				),
			},
			{
				Config: testAccNsxtPolicyPredefinedSecurityPolicyWithRulesUpdate1(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "rule1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", "DROP"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.source_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.log_label", "group1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.disabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.display_name", "rule2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.action", "ALLOW"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.source_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.log_label", "group2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.disabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.0.action", "ALLOW"),
				),
			},
			{
				Config: testAccNsxtPolicyPredefinedSecurityPolicyWithRulesUpdate2(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicySecurityPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "rule2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", "ALLOW"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.source_groups.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.log_label", "group2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.disabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "default_rule.0.action", "REJECT"),
				),
			},
			{
				Config: testAccNsxtPolicyPredefinedSecurityPolicyPrerequisites(),
			},
		},
	})
}

func testAccNsxtPolicyPredefinedSecurityPolicyPrerequisites() string {
	return `
resource "nsxt_policy_group" "group1" {
  display_name = "predefined-policy-test1"
}

resource "nsxt_policy_group" "group2" {
  display_name = "predefined-policy-test2"
}

data "nsxt_policy_security_policy" "test" {
  is_default = true
  category   = "Application"
}`
}

func testAccNsxtPolicyPredefinedSecurityPolicyBasic(description string, tags string) string {
	return testAccNsxtPolicyPredefinedSecurityPolicyPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_predefined_security_policy" "test" {
  path        = data.nsxt_policy_security_policy.test.path
  description = "%s"
  %s
}`, description, tags)
}

func testAccNsxtPolicyPredefinedSecurityPolicyDefaultRule(description string, action string, label string, tags string) string {
	return testAccNsxtPolicyPredefinedSecurityPolicyPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_predefined_security_policy" "test" {
  path = data.nsxt_policy_security_policy.test.path
  default_rule {
    description  = "%s"
    action       = "%s"
    log_label    = "%s"
    logged       = true
    %s
  }
}`, description, action, label, tags)
}

func testAccNsxtPolicyPredefinedSecurityPolicyWithRules() string {
	return testAccNsxtPolicyPredefinedSecurityPolicyPrerequisites() + `
resource "nsxt_policy_predefined_security_policy" "test" {
  path = data.nsxt_policy_security_policy.test.path

  rule {
      display_name  = "rule2"
      source_groups = [nsxt_policy_group.group2.path]
      log_label     = "group2"
      action        = "ALLOW"
      disabled      = true
  }
}`
}

func testAccNsxtPolicyPredefinedSecurityPolicyWithRulesUpdate1() string {
	return testAccNsxtPolicyPredefinedSecurityPolicyPrerequisites() + `
resource "nsxt_policy_predefined_security_policy" "test" {
  path = data.nsxt_policy_security_policy.test.path

  rule {
      display_name  = "rule1"
      source_groups = [nsxt_policy_group.group1.path]
      log_label     = "group1"
      action        = "DROP"
  }
  rule {
      display_name  = "rule2"
      source_groups = [nsxt_policy_group.group2.path]
      log_label     = "group2"
      action        = "ALLOW"
      disabled      = false
  }
}`
}

func testAccNsxtPolicyPredefinedSecurityPolicyWithRulesUpdate2() string {
	return testAccNsxtPolicyPredefinedSecurityPolicyPrerequisites() + `
resource "nsxt_policy_predefined_security_policy" "test" {
  path = data.nsxt_policy_security_policy.test.path

  rule {
      display_name  = "rule2"
      source_groups = [nsxt_policy_group.group1.path, nsxt_policy_group.group2.path]
      log_label     = "group2"
      action        = "ALLOW"
  }

  default_rule {
      action = "REJECT"
  }
}`
}
