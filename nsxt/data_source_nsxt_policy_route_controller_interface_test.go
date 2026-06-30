// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyRouteControllerInterface_basic(t *testing.T) {
	testDSByName := "data.nsxt_policy_route_controller_interface.by_name"
	testDSByID := "data.nsxt_policy_route_controller_interface.by_id"
	displayName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRCInterfacePreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRCInterfaceDataSourceTemplate(displayName),
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

func testAccNsxtPolicyRCInterfaceDataSourceTemplate(displayName string) string {
	return testAccNsxtPolicyRCInterfaceRouteControllerTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller_interface" "test" {
  display_name        = "%s"
  parent_path         = nsxt_policy_route_controller.rc.path
  floating_ip_subnets = ["192.168.202.10/24"]

  interface_address {
    subnets                        = ["192.168.202.1/24"]
    portgroup_id                   = "%s"
    virtual_network_appliance_path = data.nsxt_policy_virtual_network_appliance.vna.path
  }
}

data "nsxt_policy_route_controller_interface" "by_name" {
  display_name = nsxt_policy_route_controller_interface.test.display_name
  parent_path  = nsxt_policy_route_controller.rc.path
}

data "nsxt_policy_route_controller_interface" "by_id" {
  id          = nsxt_policy_route_controller_interface.test.nsx_id
  parent_path = nsxt_policy_route_controller.rc.path
}
`, displayName, getPortgroupID())
}
