/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TODO: Rewrite this test based on GW Policy data source when available,
// use 3.1.0 as baseline for the test and remove testAccOnlyLocalManager
func TestAccResourceNsxtPolicyFixedGatewayPolicy_basic(t *testing.T) {
	path := "/infra/domains/default/gateway-policies/Policy_Default_Infra"
	testResourceName := "nsxt_policy_fixed_gateway_policy.test"
	testGatewayResourceName := "nsxt_policy_tier0_gateway.test"
	description1 := "test 1"
	description2 := "test 2"
	tags := `tag {
            scope = "color"
            tag   = "orange"
        }`

	// NOTE: These tests cannot be parallel, as they modify same default policy
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyFixedGatewayPolicyBasic(path, description1, tags),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "description", description1),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyFixedGatewayPolicyBasic(path, description2, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayPolicyExists(testResourceName, defaultDomain),
					resource.TestCheckResourceAttr(testResourceName, "description", description2),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyFixedGatewayPolicyPrerequisites(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testGatewayResourceName),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyFixedGatewayPolicy_defaultRule(t *testing.T) {
	path := "/infra/domains/default/gateway-policies/Policy_Default_Infra"
	testResourceName := "nsxt_policy_fixed_gateway_policy.test"
	testGatewayResourceName := "nsxt_policy_tier0_gateway.test"
	action1 := "REJECT"
	action2 := "ALLOW"
	description1 := "test 1"
	description2 := "test 2"
	tags := `tag {
            scope = "color"
            tag   = "orange"
        }`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyFixedGatewayPolicyDefaultRule(path, description1, action1, action1, tags),
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
				Config: testAccNsxtPolicyFixedGatewayPolicyDefaultRule(path, description2, action2, action2, ""),
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
				Config: testAccNsxtPolicyFixedGatewayPolicyPrerequisites(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(testGatewayResourceName),
				),
			},
		},
	})
}

func testAccNsxtPolicyFixedGatewayPolicyPrerequisites() string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "fixed-gw-policy-test"
  edge_cluster_path = data.nsxt_policy_edge_cluster.test.path
}`, getEdgeClusterName())
}

func testAccNsxtPolicyFixedGatewayPolicyBasic(path string, description string, tags string) string {
	return testAccNsxtPolicyFixedGatewayPolicyPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_fixed_gateway_policy" "test" {
  path        = "%s"
  description = "%s"
  %s
}`, path, description, tags)
}

func testAccNsxtPolicyFixedGatewayPolicyDefaultRule(path string, description string, action string, label string, tags string) string {
	return testAccNsxtPolicyFixedGatewayPolicyPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_fixed_gateway_policy" "test" {
  path        = "%s"
  default_rule {
    scope        = nsxt_policy_tier0_gateway.test.path
    description  = "%s"
    action       = "%s"
    log_label    = "%s"
    logged       = true
    %s
  }
}`, path, description, action, label, tags)
}
