// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyTier1Gateways(t *testing.T) {
	testAccDataSourceNsxtPolicyTier1Gateways_basic(t, false, func() {
		testAccPreCheck(t)
	})
}
func TestAccDataSourceNsxtPolicyTier1Gateways_multitenancy(t *testing.T) {
	testAccDataSourceNsxtPolicyTier1Gateways_basic(t, true, func() {
		testAccPreCheck(t)
		testAccOnlyMultitenancy(t)
	})
}

// list the Tier1 gateways
func testAccDataSourceNsxtPolicyTier1Gateways_basic(t *testing.T, withContext bool, preCheck func()) {
	checkResourceName := "nsxt_policy_tier1_gateways_result"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyTier1GatewaysReadTemplate(withContext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput(checkResourceName, "tier1-gateway-01"),
				),
			},
			{
				Config: testAccNSXPolicyTier1GatewaysReadTemplateWithRegex(withContext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput(checkResourceName, "regex-tier1-gateway-01"),
				),
			},
		},
	})
}

func testAccNSXPolicyTier1GatewaysReadTemplate(withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "EDGECLUSTER1"
}

resource "nsxt_policy_tier1_gateway" "test" {
%s
  display_name      = "tier1-gateway-01"
  nsx_id = "test-tier1-gateway-01"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

data "nsxt_policy_tier1_gateways" "all" {
%s
depends_on = [nsxt_policy_tier1_gateway.test]
}

output "nsxt_policy_tier1_gateways_result"  {
  value = data.nsxt_policy_tier1_gateways.all.items["test-tier1-gateway-01"]
  depends_on = [data.nsxt_policy_tier1_gateways.all]
}`, context, context)
}

func testAccNSXPolicyTier1GatewaysReadTemplateWithRegex(withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtPolicyMultitenancyContext()
	}
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "EDGECLUSTER1"
}

resource "nsxt_policy_tier1_gateway" "test" {
%s
  display_name      = "regex-tier1-gateway-01"
  nsx_id = "test-tier1-gateway-01"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

data "nsxt_policy_tier1_gateways" "all" {
%s
display_name = ".*"
depends_on = [nsxt_policy_tier1_gateway.test]
}

output "nsxt_policy_tier1_gateways_result"  {
  value = data.nsxt_policy_tier1_gateways.all.items["test-tier1-gateway-01"]
  depends_on = [data.nsxt_policy_tier1_gateways.all]
}`, context, context)
}
