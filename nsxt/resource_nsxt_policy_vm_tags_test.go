/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceNsxtPolicyVMTags_basic(t *testing.T) {
	vmID := getTestVMID()
	testResourceName := "nsxt_policy_vm_tags.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccOnlyLocalManager(t); testAccPreCheck(t); testAccEnvDefined(t, "NSXT_TEST_VM_ID") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXPolicyVMTagsCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyVMTagsCreateTemplate(vmID),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXPolicyVMTagsCheckExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "instance_id", vmID),
				),
			},
			{
				Config: testAccNSXPolicyVMTagsUpdateTemplate(vmID),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXPolicyVMTagsCheckExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "instance_id", vmID),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyVMTags_withPorts(t *testing.T) {
	vmID := getTestVMID()
	testResourceName := "nsxt_policy_vm_tags.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_VM_ID")
			testAccEnvDefined(t, "NSXT_TEST_VM_SEGMENT_ID")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXPolicyVMTagsCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyVMPortTagsCreateTemplate(vmID),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXPolicyVMTagsCheckExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "port.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "port.0.segment_path"),
					resource.TestCheckResourceAttr(testResourceName, "port.0.tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "instance_id", vmID),
				),
			},
			{
				Config: testAccNSXPolicyVMPortTagsUpdateTemplate(vmID),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXPolicyVMTagsCheckExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "port.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "port.0.segment_path"),
					resource.TestCheckResourceAttr(testResourceName, "port.0.tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "instance_id", vmID),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyVMTags_import_basic(t *testing.T) {
	vmID := getTestVMID()
	testResourceName := "nsxt_policy_vm_tags.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccOnlyLocalManager(t); testAccPreCheck(t); testAccEnvDefined(t, "NSXT_TEST_VM_ID") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXPolicyVMTagsCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyVMTagsCreateTemplate(vmID),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id"},
			},
		},
	})
}

func testAccNSXPolicyVMTagsCheckExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX Policy VM Tags resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX Policy VM Tags resource ID not set in resources ")
		}

		_, err := findNsxtPolicyVMByID(connector, resourceID, testAccProvider.Meta())
		if err != nil {
			return fmt.Errorf("Failed to find VM %s", resourceID)
		}

		return nil
	}
}

func testAccNSXPolicyVMTagsCheckDestroy(state *terraform.State) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_vm_tags" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		vm, err := findNsxtPolicyVMByID(connector, resourceID, testAccProvider.Meta())
		if err != nil {
			return fmt.Errorf("Failed to find VM %s", resourceID)
		}

		if len(vm.Tags) > 0 {
			return fmt.Errorf("VM %s still has tags, although nsxt_policy_vm_tags was deleted", resourceID)
		}
	}
	return nil
}

func testAccNSXPolicyVMTagsCreateTemplate(instanceID string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_vm_tags" "test" {
  instance_id = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, instanceID)
}

func testAccNSXPolicyVMTagsUpdateTemplate(instanceID string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_vm_tags" "test" {
  instance_id = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}`, instanceID)
}

func testAccNSXPolicyVMPortTagsCreateTemplate(instanceID string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_vm_tags" "test" {
  instance_id = "%s"

  tag {
    scope = "color"
    tag   = "blue"
  }

  port {
    segment_path = "/infra/segments/%s"
    tag {
      scope = "color"
      tag   = "green"
    }
  }
}`, instanceID, getTestVMSegmentID())
}

func testAccNSXPolicyVMPortTagsUpdateTemplate(instanceID string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_vm_tags" "test" {
  instance_id = "%s"

  port {
    segment_path = "/infra/segments/%s"
    tag {
      scope = "color"
      tag   = "green"
    }
    tag {
      scope = "shape"
      tag   = "round"
    }
  }
}`, instanceID, getTestVMSegmentID())
}
