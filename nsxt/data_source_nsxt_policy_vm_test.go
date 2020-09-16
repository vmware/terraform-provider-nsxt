/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNsxtPolicyVM_basic(t *testing.T) {
	testResourceName := "data.nsxt_policy_vm.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_VM_ID")
			testAccEnvDefined(t, "NSXT_TEST_VM_NAME")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVMReadByNameTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "bios_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "external_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "instance_id"),
				),
			},
			{
				Config: testAccNsxtPolicyVMReadByIDTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "bios_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "external_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "instance_id"),
				),
			},
		},
	})
}

func testAccNsxtPolicyVMReadByNameTemplate() string {
	return fmt.Sprintf(`
data "nsxt_policy_vm" "test" {
  display_name = "%s"
}`, getTestVMName())
}

func testAccNsxtPolicyVMReadByIDTemplate() string {
	return fmt.Sprintf(`
data "nsxt_policy_vm" "test" {
  external_id = "%s"
}`, getTestVMID())
}
