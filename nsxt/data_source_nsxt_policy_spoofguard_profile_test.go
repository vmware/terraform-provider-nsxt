/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNsxtPolicySpoofGuardProfile_basic(t *testing.T) {
	// Use existing system defined profile
	name := "default-spoofguard-profile"
	testResourceName := "data.nsxt_policy_spoofguard_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicySpoofGuardProfileReadTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
		},
	})
}

func testAccNsxtPolicySpoofGuardProfileReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_spoofguard_profile" "test" {
  display_name = "%s"
}`, name)
}
