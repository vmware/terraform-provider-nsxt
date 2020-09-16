/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNsxtPolicyGatewayPolicy_basic(t *testing.T) {
	name := "terraform_ds_test"
	category := "LocalGatewayRules"
	testResourceName := "data.nsxt_policy_gateway_policy.test"
	withCategory := fmt.Sprintf(`category = "%s"`, category)
	withDomain := `domain = "default"`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewayPolicyTemplate(name, category, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttr(testResourceName, "category", category),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyGatewayPolicyTemplate(name, category, withCategory),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttr(testResourceName, "category", category),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyGatewayPolicyTemplate(name, category, withDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttr(testResourceName, "category", category),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

// This test is applicable to earlier NSX versions, where a single default GW
// policy is auto-created. In later versions a policy is created per GW.
func TestAccDataSourceNsxtPolicyGatewayPolicy_default(t *testing.T) {
	testResourceName := "data.nsxt_policy_gateway_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersionLessThan(t, "3.1.0") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDefaultGatewayPolicyTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(testResourceName, "description"),
					resource.TestCheckResourceAttr(testResourceName, "category", "Default"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccNsxtPolicyGatewayPolicyTemplate(name string, category string, extra string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_gateway_policy" "test" {
  display_name = "%s"
  description  = "%s"
  category     = "%s"
}

data "nsxt_policy_gateway_policy" "test" {
  display_name = nsxt_policy_gateway_policy.test.display_name
  %s
}`, name, name, category, extra)
}

func testAccNsxtPolicyDefaultGatewayPolicyTemplate() string {
	return `
data "nsxt_policy_gateway_policy" "test" {
  category     = "Default"
}`
}
