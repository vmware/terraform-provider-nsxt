// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtProxyConfig_basic(t *testing.T) {
	testResourceName := "data.nsxt_proxy_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "3.1.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtProxyConfigMinimalistic() + testAccNsxtProxyConfigReadTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "enabled"),
				),
			},
		},
	})
}

func testAccNsxtProxyConfigReadTemplate() string {
	// Note: Proxy config is a singleton, no filter needed
	return `
data "nsxt_proxy_config" "test" {
}`
}

func testAccNsxtProxyConfigMinimalistic() string {
	// Note: display_name is not included as NSX doesn't support custom display names for proxy config singleton
	return `
resource "nsxt_proxy_config" "test" {
  enabled      = false
  host         = "proxy.example.com"
}`
}
