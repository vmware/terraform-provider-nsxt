// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var accTestVpcIpAddressAllocationCreateAttributes = map[string]string{
	"display_name":    getAccTestResourceName(),
	"description":     "terraform created",
	"allocation_size": "2",
}

var accTestVpcIpAddressAllocationUpdateAttributes = map[string]string{
	"display_name":    getAccTestResourceName(),
	"description":     "terraform updated",
	"allocation_size": "2",
}

func TestAccResourceNsxtVpcIpAddressAllocation_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_ip_address_allocation.test"
	createName := getAccTestResourceName()
	updateName := getAccTestResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcIpAddressAllocationCheckDestroy(state, updateName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcIpAddressAllocationTemplate(true, createName, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcIpAddressAllocationExists(createName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcIpAddressAllocationCreateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "allocation_ips"),
					resource.TestCheckResourceAttr(testResourceName, "allocation_size", accTestVpcIpAddressAllocationCreateAttributes["allocation_size"]),
					resource.TestCheckResourceAttrSet(testResourceName, "ip_address_type"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcIpAddressAllocationTemplate(false, createName, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcIpAddressAllocationExists(updateName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcIpAddressAllocationUpdateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "allocation_ips"),
					resource.TestCheckResourceAttr(testResourceName, "allocation_size", accTestVpcIpAddressAllocationUpdateAttributes["allocation_size"]),
					resource.TestCheckResourceAttrSet(testResourceName, "ip_address_type"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcIpAddressAllocationMinimalistic(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcIpAddressAllocationExists(updateName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtVpcIpAddressAllocation_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_ip_address_allocation.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcIpAddressAllocationCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcIpAddressAllocationMinimalistic(name),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
		},
	})
}

func testAccNsxtVpcIpAddressAllocationExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy VpcIpAddressAllocation resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy VpcIpAddressAllocation resource ID not set in resources")
		}

		exists, err := resourceNsxtVpcIpAddressAllocationExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy VpcIpAddressAllocation %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcIpAddressAllocationCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_vpc_ip_address_allocation" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcIpAddressAllocationExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy VpcIpAddressAllocation %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcIpAddressAllocationTemplate(createFlow bool, createName, updateName string) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcIpAddressAllocationCreateAttributes
	} else {
		attrMap = accTestVpcIpAddressAllocationUpdateAttributes
	}
	displayName := updateName
	if createFlow {
		displayName = createName
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_ip_address_allocation" "test" {
  %s
  display_name    = "%s"
  description     = "%s"
  allocation_size = %s
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_vpc_ip_address_allocation" "test" {
  %s
  allocation_ips = nsxt_vpc_ip_address_allocation.test.allocation_ips
}`, testAccNsxtPolicyMultitenancyContext(), displayName, attrMap["description"], attrMap["allocation_size"], testAccNsxtPolicyMultitenancyContext())
}

func testAccNsxtVpcIpAddressAllocationMinimalistic(updateName string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_ip_address_allocation" "test" {
  %s
  display_name    = "%s"
  allocation_size = %s
}

data "nsxt_vpc_ip_address_allocation" "test" {
  %s
  allocation_ips = nsxt_vpc_ip_address_allocation.test.allocation_ips
}`, testAccNsxtPolicyMultitenancyContext(), updateName, accTestVpcIpAddressAllocationUpdateAttributes["allocation_size"], testAccNsxtPolicyMultitenancyContext())
}
