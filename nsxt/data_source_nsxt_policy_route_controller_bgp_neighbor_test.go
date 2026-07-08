// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyRouteControllerBgpNeighbor_basic(t *testing.T) {
	testDSByName := "data.nsxt_policy_route_controller_bgp_neighbor.by_name"
	testDSByID := "data.nsxt_policy_route_controller_bgp_neighbor.by_id"
	displayName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRCBgpNeighborPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRCBgpNeighborDataSourceTemplate(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testDSByName, "display_name", displayName),
					resource.TestCheckResourceAttrSet(testDSByName, "id"),
					resource.TestCheckResourceAttrSet(testDSByName, "path"),
					resource.TestCheckResourceAttr(testDSByID, "display_name", displayName),
					resource.TestCheckResourceAttrSet(testDSByID, "id"),
					resource.TestCheckResourceAttrSet(testDSByID, "path"),
				),
			},
		},
	})
}

func testAccNsxtPolicyRCBgpNeighborDataSourceTemplate(displayName string) string {
	return testAccNsxtPolicyRCBgpNeighborRouteControllerTemplate() + testAccNsxtPolicyRCBgpNeighborInterfaceTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller_bgp_neighbor" "test" {
  display_name     = "%s"
  parent_path      = "${nsxt_policy_route_controller.rc.path}/bgp"
  neighbor_address = "192.168.100.1"
  remote_as_num    = "65099"
  source_addresses = ["192.168.200.1"]

  depends_on = [nsxt_policy_route_controller_interface.bgp_src]
}

data "nsxt_policy_route_controller_bgp_neighbor" "by_name" {
  display_name = nsxt_policy_route_controller_bgp_neighbor.test.display_name
  parent_path  = "${nsxt_policy_route_controller.rc.path}/bgp"
}

data "nsxt_policy_route_controller_bgp_neighbor" "by_id" {
  id          = nsxt_policy_route_controller_bgp_neighbor.test.nsx_id
  parent_path = "${nsxt_policy_route_controller.rc.path}/bgp"
}
`, displayName)
}
