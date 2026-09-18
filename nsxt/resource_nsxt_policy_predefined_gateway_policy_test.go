// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccResourceNsxtPolicyPredefinedGatewayPolicy_basic verifies that `rule`
// is rejected outside of VMC. Read on this resource reports every rule on the
// policy, not just the ones declared in `rule` here; on NSX Local Manager or
// Global Manager the policy's rules are typically also managed by another
// resource (nsxt_policy_gateway_policy, or nsxt_policy_parent_gateway_policy +
// nsxt_policy_gateway_policy_rule), so combining `rule` here with any of those
// causes a permanent configuration drift and can corrupt or delete rules
// owned elsewhere. There is no legitimate use of `rule` outside VMC.
func TestAccResourceNsxtPolicyPredefinedGatewayPolicy_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNsxtPolicyPredefinedGatewayPolicyRule(),
				ExpectError: regexp.MustCompile("rule is only supported for nsxt_policy_predefined_gateway_policy on VMC"),
			},
		},
	})
}

// TestAccResourceNsxtPolicyPredefinedGatewayPolicy_multitenancy verifies the
// same VMC-only restriction on `rule` holds under a multitenancy (Project)
// context. See TestAccResourceNsxtPolicyPredefinedGatewayPolicy_basic.
func TestAccResourceNsxtPolicyPredefinedGatewayPolicy_multitenancy(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyMultitenancy(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNsxtPolicyPredefinedGatewayPolicyRuleMultitenancy(),
				ExpectError: regexp.MustCompile("rule is only supported for nsxt_policy_predefined_gateway_policy on VMC"),
			},
		},
	})
}

func TestAccResourceNsxtPolicyPredefinedGatewayPolicy_defaultRule(t *testing.T) {
	testResourceName := "nsxt_policy_predefined_gateway_policy.test"
	testGatewayResourceName := "nsxt_policy_tier0_gateway.test"
	action1 := "REJECT"
	action2 := "ALLOW"
	description1 := "test 1"
	description2 := "test 2"
	tags := `tag {
            scope = "color"
            tag   = "orange"
        }`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyPredefinedGatewayPolicyDefaultRule(description1, action1, action1, tags),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayPolicyExists(testResourceName, defaultDomain),
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
				Config: testAccNsxtPolicyPredefinedGatewayPolicyDefaultRule(description2, action2, action2, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayPolicyExists(testResourceName, defaultDomain),
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
			{
				Config: testAccNsxtPolicyPredefinedGatewayPolicyPrerequisites(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testGatewayResourceName),
				),
			},
		},
	})
}

func testAccNsxtPolicyPredefinedGatewayPolicyPrerequisites() string {
	t0EdgeCluster := `edge_cluster_path = data.nsxt_policy_edge_cluster.test.path`
	if testAccIsGlobalManager() {
		t0EdgeCluster = fmt.Sprintf(`locale_service { %s }`, t0EdgeCluster)
	}

	return testAccNsxtPolicyEdgeClusterReadTemplate(getEdgeClusterName()) + fmt.Sprintf(`


resource "nsxt_policy_tier0_gateway" "test" {
  display_name = "predefined-gw-policy-test"
  %s
}

data "nsxt_policy_gateway_policy" "test" {
  category     = "Default"
  display_name = "Policy_Default_Infra-tier0-${nsxt_policy_tier0_gateway.test.nsx_id}"
}`, t0EdgeCluster)
}

func testAccNsxtPolicyPredefinedGatewayPolicyPrerequisitesMultitenancy() string {
	t1EdgeCluster := `edge_cluster_path = data.nsxt_policy_edge_cluster.test.path`
	context := testAccNsxtPolicyMultitenancyContext()
	return testAccNsxtPolicyEdgeClusterReadTemplate(getEdgeClusterName()) + fmt.Sprintf(`


resource "nsxt_policy_tier1_gateway" "test" {
%s
  display_name = "predefined-gw-policy-test"
  %s
}

data "nsxt_policy_gateway_policy" "test" {
%s
  category     = "Default"
  display_name = "Policy_Default_Infra-tier1-${nsxt_policy_tier1_gateway.test.nsx_id}"
}`, context, t1EdgeCluster, context)
}

func testAccNsxtPolicyPredefinedGatewayPolicyRule() string {
	return testAccNsxtPolicyPredefinedGatewayPolicyPrerequisites() + `

resource "nsxt_policy_predefined_gateway_policy" "test" {
  path = data.nsxt_policy_gateway_policy.test.path

  rule {
    display_name = "rule1"
    action       = "DROP"
    scope        = [nsxt_policy_tier0_gateway.test.path]
  }
}`
}

func testAccNsxtPolicyPredefinedGatewayPolicyRuleMultitenancy() string {
	return testAccNsxtPolicyPredefinedGatewayPolicyPrerequisitesMultitenancy() + fmt.Sprintf(`

resource "nsxt_policy_predefined_gateway_policy" "test" {
%s
  path = data.nsxt_policy_gateway_policy.test.path

  rule {
    display_name = "rule1"
    action       = "DROP"
    scope        = [nsxt_policy_tier1_gateway.test.path]
  }
}`, testAccNsxtPolicyMultitenancyContext())
}

func testAccNsxtPolicyPredefinedGatewayPolicyDefaultRule(description string, action string, label string, tags string) string {
	return testAccNsxtPolicyPredefinedGatewayPolicyPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_predefined_gateway_policy" "test" {
  path = data.nsxt_policy_gateway_policy.test.path
  default_rule {
    scope       = nsxt_policy_tier0_gateway.test.path
    description = "%s"
    action      = "%s"
    log_label   = "%s"
    logged      = true
    %s
  }
}`, description, action, label, tags)
}
