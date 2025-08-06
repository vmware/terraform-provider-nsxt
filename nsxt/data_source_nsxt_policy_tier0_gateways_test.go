// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// list the Tier0 gateways
func TestAccDataSourceNsxtPolicyTier0Gateways_basic(t *testing.T) {
	checkResourceName := "nsxt_policy_tier0_gateways_result"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyTier0GatewaysReadTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput(checkResourceName, "tier0-gateway-01"),
				),
			},
			{
				Config: testAccNSXPolicyTier0GatewaysReadTemplateWithRegex(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput(checkResourceName, "regex-tier0-gateway-01"),
				),
			},
		},
	})
}

func testAccNSXPolicyTier0GatewaysReadTemplate() string {
	return fmt.Sprintln(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "EDGECLUSTER1"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "tier0-gateway-01"
  nsx_id = "test-tier0-gateway-01"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

data "nsxt_policy_tier0_gateways" "all" {
depends_on = [nsxt_policy_tier0_gateway.test]
}

output "nsxt_policy_tier0_gateways_result"  {
  value = data.nsxt_policy_tier0_gateways.all.items["test-tier0-gateway-01"]
  depends_on = [data.nsxt_policy_tier0_gateways.all]
}`)
}

func testAccNSXPolicyTier0GatewaysReadTemplateWithRegex() string {
	return fmt.Sprintln(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "EDGECLUSTER1"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "regex-tier0-gateway-01"
  nsx_id = "testre-tier0-gateway-01"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

data "nsxt_policy_tier0_gateways" "all" {
display_name = ".*"
depends_on = [nsxt_policy_tier0_gateway.test]
}

output "nsxt_policy_tier0_gateways_result"  {
  value = data.nsxt_policy_tier0_gateways.all.items["testre-tier0-gateway-01"]
  depends_on = [data.nsxt_policy_tier0_gateways.all]
}`)
}
