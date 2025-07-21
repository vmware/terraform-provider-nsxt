// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyTransitGateway_basic(t *testing.T) {
	testAccDataSourceNsxtPolicyTransitGatewayBasic(t, true, func() {
		testAccPreCheck(t)
		testAccIsMultitenancy()
	})
}

func testAccDataSourceNsxtPolicyTransitGatewayBasic(t *testing.T, withContext bool, preCheck func()) {
	testResourceName := "data.nsxt_policy_transit_gateway.test"
	checkResourceName := "data.nsxt_policy_transit_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyTransitGatewayDefaultReadTemplate(withContext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(checkResourceName, "is_default", "true"),
				),
			},
		},
	})
}

func testAccNSXPolicyTransitGatewayDefaultReadTemplate(withContext bool) string {
	context := ""
	if withContext {
		context = testAccNsxtProjectContext()
	}
	return fmt.Sprintf(`
	
	data "nsxt_policy_transit_gateway" "test" {
	%s
	
		  is_default = "true"
		}`, context)
}
