/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceNsxtPolicyPredefinedGatewayPolicy_basic(t *testing.T) {
	testResourceName := "nsxt_policy_predefined_gateway_policy.test"
	testGatewayResourceName := "nsxt_policy_tier0_gateway.test"
	description1 := "test 1"
	description2 := "test 2"
	tags := `tag {
            scope = "color"
            tag   = "orange"
        }`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t); testAccNSXVersion(t, "3.1.0") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyPredefinedGatewayPolicyBasic(description1, tags),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "description", description1),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyPredefinedGatewayPolicyBasic(description2, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "description", description2),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
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
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t); testAccNSXVersion(t, "3.1.0") },
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
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "predefined-gw-policy-test"
  edge_cluster_path = data.nsxt_policy_edge_cluster.test.path
}

data "nsxt_policy_gateway_policy" "test" {
  category     = "Default"
  display_name = "Policy_Default_Infra-tier0-${nsxt_policy_tier0_gateway.test.nsx_id}"
}`, getEdgeClusterName())
}

func testAccNsxtPolicyPredefinedGatewayPolicyBasic(description string, tags string) string {
	return testAccNsxtPolicyPredefinedGatewayPolicyPrerequisites() + fmt.Sprintf(`

resource "nsxt_policy_predefined_gateway_policy" "test" {
  path        = data.nsxt_policy_gateway_policy.test.path
  description = "%s"
  %s
}`, description, tags)
}

func testAccNsxtPolicyPredefinedGatewayPolicyDefaultRule(description string, action string, label string, tags string) string {
	return testAccNsxtPolicyPredefinedGatewayPolicyPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_predefined_gateway_policy" "test" {
  path        = data.nsxt_policy_gateway_policy.test.path
  default_rule {
    scope        = nsxt_policy_tier0_gateway.test.path
    description  = "%s"
    action       = "%s"
    log_label    = "%s"
    logged       = true
    %s
  }
}`, description, action, label, tags)
}
