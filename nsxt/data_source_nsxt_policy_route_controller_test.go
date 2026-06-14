// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyRouteController_basic(t *testing.T) {
	testResourceName := "nsxt_policy_route_controller.test"
	testDSByName := "data.nsxt_policy_route_controller.by_name"
	testDSByID := "data.nsxt_policy_route_controller.by_id"
	displayName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNsxtPolicyRouteControllerPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyRouteControllerDataSourceTemplate(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
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

func testAccNsxtPolicyRouteControllerDataSourceTemplate(displayName string) string {
	return testAccNsxtPolicyRouteControllerVnaTemplate() + fmt.Sprintf(`
resource "nsxt_policy_route_controller" "test" {
  display_name                           = "%s"
  virtual_network_appliance_cluster_path = data.nsxt_policy_virtual_network_appliance_cluster.vna.path
}

data "nsxt_policy_route_controller" "by_name" {
  display_name = nsxt_policy_route_controller.test.display_name
}

data "nsxt_policy_route_controller" "by_id" {
  id = nsxt_policy_route_controller.test.id
}
`, displayName)
}
