/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceNsxtPolicyContextProfile_basic(t *testing.T) {
	// Use existing system defined profile
	name := "AMQP"
	testResourceName := "data.nsxt_policy_context_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyContextProfileReadTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttrSet(testResourceName, "description"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyEmptyTemplate(),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyContextProfile_prefix(t *testing.T) {
	// Use existing system defined profile
	name := "DIAMETER"
	namePrefix := string(name[0:5])
	testResourceName := "data.nsxt_policy_context_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyContextProfileReadTemplate(namePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttrSet(testResourceName, "description"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyEmptyTemplate(),
			},
		},
	})
}

func testAccNsxtPolicyContextProfileReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_context_profile" "test" {
  display_name = "%s"
}`, name)
}
