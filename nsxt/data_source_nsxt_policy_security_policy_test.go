/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNsxtPolicySecurityPolicy_basic(t *testing.T) {
	name := "terraform_ds_test"
	category := "Application"
	testResourceName := "data.nsxt_policy_security_policy.test"
	withCategory := fmt.Sprintf(`category = "%s"`, category)
	withDomain := `domain = "default"`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicySecurityPolicyTemplate(name, category, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttr(testResourceName, "category", category),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicySecurityPolicyTemplate(name, category, withCategory),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttr(testResourceName, "category", category),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicySecurityPolicyTemplate(name, category, withDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttr(testResourceName, "category", category),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtEmptyTemplate(),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicySecurityPolicy_default(t *testing.T) {
	testResourceName := "data.nsxt_policy_security_policy.test"
	category1 := "Ethernet"
	category2 := "Application"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDefaultSecurityPolicyTemplate(category1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttr(testResourceName, "category", category1),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyDefaultSecurityPolicyTemplate(category2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttr(testResourceName, "category", category2),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccNsxtPolicySecurityPolicyTemplate(name string, category string, extra string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_security_policy" "test" {
  display_name = "%s"
  description  = "%s"
  category     = "%s"
}

data "nsxt_policy_security_policy" "test" {
  display_name = nsxt_policy_security_policy.test.display_name
  %s
}`, name, name, category, extra)
}

func testAccNsxtPolicyDefaultSecurityPolicyTemplate(category string) string {
	return fmt.Sprintf(`
data "nsxt_policy_security_policy" "test" {
  is_default = true
  category   = "%s"
}`, category)
}
