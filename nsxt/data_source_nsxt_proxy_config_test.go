// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtProxyConfig_basic(t *testing.T) {
	testAccDataSourceNsxtPolicyProxyConfigBasic(t, false, func() {
		testAccPreCheck(t)
		testAccOnlyLocalManager(t)
		testAccNSXVersion(t, "3.1.0")
	})
}

func testAccDataSourceNsxtPolicyProxyConfigBasic(t *testing.T, withContext bool, preCheck func()) {
	testResourceName := "data.nsxt_proxy_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  preCheck,
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtProxyConfigMinimalistic(withContext) + testAccNsxtProxyConfigReadTemplate(withContext),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "enabled"),
				),
			},
		},
	})
}

func testAccNsxtProxyConfigReadTemplate(withContext bool) string {
	context := ""
	// Note: Proxy config is a singleton, no filter needed
	return fmt.Sprintf(`
data "nsxt_proxy_config" "test" {
%s
}`, context)
}

func testAccNsxtProxyConfigMinimalistic(withContext bool) string {
	context := ""
	// Note: display_name is not included as NSX doesn't support custom display names for proxy config singleton
	return fmt.Sprintf(`
resource "nsxt_proxy_config" "test" {
%s
  enabled      = false
  host         = "proxy.example.com"
}`, context)
}
