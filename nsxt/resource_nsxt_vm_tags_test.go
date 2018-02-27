/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt"
	"testing"
)

var vmTagsResourceName = "test"
var vmTagsFullResourceName = "nsxt_vm_tags." + vmTagsResourceName

func TestAccResourceNsxtVMTags(t *testing.T) {
	vmID := getTestVMID()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccEnvDefined(t, "NSX_TEST_VM_ID") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXVMTagsCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXVMTagsCreateTemplate(vmID),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXVMTagsCheckExists(),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNSXVMTagsUpdateTemplate(vmID),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXVMTagsCheckExists(),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "tag.#", "2"),
				),
			},
		},
	})
}

func testAccNSXVMTagsCheckExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)

		rs, ok := state.RootModule().Resources[vmTagsFullResourceName]
		if !ok {
			return fmt.Errorf("NSX vm tags resource %s not found in resources", vmTagsFullResourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX vm tags resource ID not set in resources ")
		}

		_, err := findVMByExternalID(nsxClient, resourceID)
		if err != nil {
			return fmt.Errorf("Failed to find VM %s", resourceID)
		}

		return nil
	}
}

func testAccNSXVMTagsCheckDestroy(state *terraform.State) error {
	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_vm_tags" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		vm, err := findVMByExternalID(nsxClient, resourceID)
		if err != nil {
			return fmt.Errorf("Failed to find VM %s", resourceID)
		}

		if len(vm.Tags) > 0 {
			return fmt.Errorf("VM %s still has tags, although nsxt_vm_tags was deleted", resourceID)
		}
	}
	return nil
}

func testAccNSXVMTagsCreateTemplate(instanceID string) string {
	return fmt.Sprintf(`
resource "nsxt_vm_tags" "%s" {
  instance_id = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, vmTagsResourceName, instanceID)
}

func testAccNSXVMTagsUpdateTemplate(instanceID string) string {
	return fmt.Sprintf(`
resource "nsxt_vm_tags" "%s" {
  instance_id = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}`, vmTagsResourceName, instanceID)
}
