/* Copyright © 2023 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNsxtComputeManager_basic(t *testing.T) {
	ComputeManagerName := getComputeManagerName()
	testResourceName := "data.nsxt_edge_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_COMPUTE_MANAGER")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXComputeManagerReadTemplate(ComputeManagerName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", ComputeManagerName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "description"),
					resource.TestCheckResourceAttrSet(testResourceName, "server"),
				),
			},
		},
	})
}

func testAccNSXComputeManagerReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_compute_manager" "test" {
  display_name = "%s"
}`, name)
}
