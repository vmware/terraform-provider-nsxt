/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNsxtIPPool_basic(t *testing.T) {
	ipPoolName := getIPPoolName()
	if ipPoolName == "" {
		t.Skipf("No NSXT_TEST_IP_POOL set - skipping test")
	}
	testResourceName := "data.nsxt_ip_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccOnlyLocalManager(t); testAccTestDeprecated(t); testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXIPPoolReadTemplate(ipPoolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", ipPoolName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
				),
			},
		},
	})
}

func testAccNSXIPPoolReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_ip_pool" "test" {
  display_name = "%s"
}`, name)
}
