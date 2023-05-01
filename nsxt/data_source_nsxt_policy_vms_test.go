/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNsxtPolicyVMs_basic(t *testing.T) {
	testResourceName := "data.nsxt_policy_vms.test"
	checkDataSourceName := "data.nsxt_policy_vm.check"
	checkResourceName := "nsxt_policy_group.check"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccEnvDefined(t, "NSXT_TEST_VM_NAME")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVMsTemplate("bios_id"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrPair(checkResourceName, "display_name", checkDataSourceName, "bios_id"),
				),
			},
			{
				Config: testAccNsxtPolicyVMsTemplate("external_id"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrPair(checkResourceName, "display_name", checkDataSourceName, "external_id"),
				),
			},
			{
				Config: testAccNsxtPolicyVMsTemplate("instance_id"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrPair(checkResourceName, "display_name", checkDataSourceName, "instance_id"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyVMs_filter(t *testing.T) {
	testResourceName := "data.nsxt_policy_vms.test"
	checkResourceName := "nsxt_policy_group.check"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVMsTemplateFilter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(checkResourceName, "display_name"),
				),
			},
		},
	})
}

func testAccNsxtPolicyVMsTemplate(valueType string) string {
	return fmt.Sprintf(`
data "nsxt_policy_vms" "test" {
  value_type = "%s"
}

data "nsxt_policy_vm" "check" {
  display_name = "%s"
}

resource "nsxt_policy_group" "check" {
  display_name = data.nsxt_policy_vms.test.items["%s"]
}`, valueType, getTestVMName(), getTestVMName())
}

func testAccNsxtPolicyVMsTemplateFilter() string {
	return `
data "nsxt_policy_vms" "test" {
  state    = "running"
  guest_os = "ubuntu"
}

resource "nsxt_policy_group" "check" {
  display_name = length(data.nsxt_policy_vms.test.items)
}`
}
