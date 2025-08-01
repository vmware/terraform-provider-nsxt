// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// list the Tier1 gateways
func TestAccDataSourceNsxtPolicyTier1Gateways_basic(t *testing.T) {
	checkResourceName := "nsxt_policy_tier1_gateways_result"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			//	testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.1.0")
			//testAccEnvDefined(t, "NSXT_PROJECT_ID")

		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyTier1GatewaysReadTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput(checkResourceName, "1"),
				),
			},
		},
	})
}

func testAccNSXPolicyTier1GatewaysReadTemplate() string {
	return fmt.Sprintln(`
data "nsxt_policy_tier1_gateways" "all" {}

output "nsxt_policy_tier1_gateways_result"  {
  value = length(data.nsxt_policy_tier1_gateways.all.items)
}`)
}
