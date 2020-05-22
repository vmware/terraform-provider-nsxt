/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/vmware/go-vmware-nsxt/common"
	"github.com/vmware/go-vmware-nsxt/manager"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var vmTagsResourceName = "test"
var vmTagsFullResourceName = "nsxt_vm_tags." + vmTagsResourceName

func TestAccResourceNsxtVMTags_basic(t *testing.T) {
	vmID := getTestVMID()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccEnvDefined(t, "NSXT_TEST_VM_ID") },
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
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "logical_port_tag.#", "1"),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "logical_port_tag.934310497.scope", "a"),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "logical_port_tag.934310497.tag", "b"),
				),
			},
			{
				Config: testAccNSXVMTagsUpdateTemplate(vmID),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXVMTagsCheckExists(),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "logical_port_tag.#", "1"),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "logical_port_tag.600822426.scope", "c"),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "logical_port_tag.600822426.tag", "d"),
				),
			},
		},
	})
}

func TestAccResourceNsxtVMTags_import_basic(t *testing.T) {
	vmID := getTestVMID()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccEnvDefined(t, "NSXT_TEST_VM_ID") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXVMTagsCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXVMTagsCreateTemplate(vmID),
			},
			{
				ResourceName:            vmTagsFullResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id"},
			},
		},
	})
}

func TestAccResourceNsxtVMTags_withSpecificLogicalPort(t *testing.T) {
	vmID := getTestVMID()
	logicalSwitchID := getTestLogicalSwitchID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccEnvDefined(t, "NSXT_TEST_VM_ID")
			testAccEnvDefined(t, "NSXT_TEST_VM_LOGICAL_SWITCH_ID")
			testAccEnvDefined(t, "NSXT_TEST_VM_NON_TAGGED_LOGICAL_SWITCH_ID")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXVMTagsCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXVMTagsCreateTemplateSpecificLogicalSwitch(vmID, logicalSwitchID),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXVMTagsCheckExists(),
					testAccNSXVMPortTagsCheckNotExist(
						[]common.Tag{
							{
								Scope: "a",
								Tag:   "b",
							},
						},
					),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "logical_port_tag.#", "1"),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "logical_port_tag.934310497.scope", "a"),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "logical_port_tag.934310497.tag", "b"),
				),
			},
			{
				Config: testAccNSXVMTagsUpdateTemplateSpecificLogicalSwitch(vmID, logicalSwitchID),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXVMTagsCheckExists(),
					testAccNSXVMPortTagsCheckNotExist(
						[]common.Tag{
							{
								Scope: "c",
								Tag:   "d",
							},
							{
								Scope: "e",
								Tag:   "f",
							},
						},
					),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "logical_port_tag.#", "2"),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "logical_port_tag.600822426.scope", "c"),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "logical_port_tag.600822426.tag", "d"),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "logical_port_tag.3616979359.scope", "e"),
					resource.TestCheckResourceAttr(vmTagsFullResourceName, "logical_port_tag.3616979359.tag", "f"),
				),
			},
		},
	})
}

func testAccNSXVMTagsCheckExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {

		nsxClient := testAccProvider.Meta().(nsxtClients).NsxtClient

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

func testAccNSXVMPortTagsCheckNotExist(checkedTags []common.Tag) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		nonTaggedLogicalSwitchID := getTestNonTaggedLogicalSwitchID()
		nsxClient := testAccProvider.Meta().(nsxtClients).NsxtClient

		rs, ok := state.RootModule().Resources[vmTagsFullResourceName]
		if !ok {
			return fmt.Errorf("NSX tags resource %s not found in resources", vmTagsFullResourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX vm tags resource ID not set in resources ")
		}

		ports, err := findPortsByExternalID(nsxClient, resourceID)
		if err != nil {
			return fmt.Errorf("Error during logical port retrieval: %v", err)
		}

		if len(ports) > 0 {
			var nonTaggedPort manager.LogicalPort
			for _, port := range ports {
				if port.LogicalSwitchId == nonTaggedLogicalSwitchID {
					nonTaggedPort = port
					break
				}
			}

			if checkIfTagInPort(nonTaggedPort, checkedTags) {
				return fmt.Errorf("The port attached to %s shouldn't be tagged with %s", nonTaggedLogicalSwitchID, checkedTags)
			}
		} else {
			return fmt.Errorf("VM %s is not attached to any networks", resourceID)
		}

		return nil
	}
}

func testAccNSXVMTagsCheckDestroy(state *terraform.State) error {
	nsxClient := testAccProvider.Meta().(nsxtClients).NsxtClient
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

  logical_port_tag {
    scope = "a"
    tag   = "b"
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

  logical_port_tag {
    scope = "c"
    tag   = "d"
  }
}`, vmTagsResourceName, instanceID)
}

func testAccNSXVMTagsCreateTemplateSpecificLogicalSwitch(instanceID string, logicalSwitchID string) string {
	return fmt.Sprintf(`
resource "nsxt_vm_tags" "%s" {
  instance_id = "%s"
  logical_switch_id = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  logical_port_tag {
    scope = "a"
    tag   = "b"
  }
}`, vmTagsResourceName, instanceID, logicalSwitchID)
}

func testAccNSXVMTagsUpdateTemplateSpecificLogicalSwitch(instanceID string, logicalSwitchID string) string {
	return fmt.Sprintf(`
resource "nsxt_vm_tags" "%s" {
  instance_id = "%s"
  logical_switch_id = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }

  logical_port_tag {
    scope = "c"
    tag   = "d"
  }

  logical_port_tag {
    scope = "e"
    tag   = "f"
  }
}`, vmTagsResourceName, instanceID, logicalSwitchID)
}

func checkIfTagInPort(nonTaggedPort manager.LogicalPort, checkedTags []common.Tag) bool {
	for _, tag := range nonTaggedPort.Tags {
		for _, checkedTag := range checkedTags {
			if tag.Scope == checkedTag.Scope && tag.Tag == checkedTag.Tag {
				return true
			}
		}
	}
	return false
}
