/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceNsxtLogicalTier0Router_basic(t *testing.T) {
	routerName := getTier0RouterName()
	testResourceName := "data.nsxt_logical_tier0_router.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXLogicalTier0RouterReadTemplate(routerName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", routerName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "edge_cluster_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "high_availability_mode"),
				),
			},
		},
	})
}

func testAccNSXLogicalTier0RouterReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_logical_tier0_router" "test" {
  display_name = "%s"
}`, name)
}
