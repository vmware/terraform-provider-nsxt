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
	checkResourceName := "data.nsxt_policy_tier0_gateways.all"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.1.0")
			testAccEnvDefined(t, "NSXT_PROJECT_ID")

		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyTier0GatewaysReadTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(checkResourceName, "items.#", "1"),
				),
			},
		},
	})
}

func testAccNSXPolicyTier0GatewaysReadTemplate() string {
	return fmt.Sprintln(`
	data "nsxt_policy_tier0_gateways" "all" {}`)
}
