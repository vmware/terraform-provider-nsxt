// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var accTestProjectIpAddressAllocationCreateAttributes = map[string]string{
	"display_name":    getAccTestResourceName(),
	"description":     "terraform created",
	"allocation_size": "1",
}

var accTestProjectIpAddressAllocationUpdateAttributes = map[string]string{
	"display_name":    getAccTestResourceName(),
	"description":     "terraform updated",
	"allocation_size": "1",
}

func TestAccResourceNsxtPolicyProjectIpAddressAllocation_basic(t *testing.T) {
	testResourceName := "nsxt_policy_project_ip_address_allocation.test"
	createName := getAccTestResourceName()
	updateName := getAccTestResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectIpAddressAllocationCheckDestroy(state, updateName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyProjectIpAddressAllocationTemplate(true, createName, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectIpAddressAllocationExists(createName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestProjectIpAddressAllocationCreateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "allocation_ips"),
					resource.TestCheckResourceAttr(testResourceName, "allocation_size", accTestProjectIpAddressAllocationCreateAttributes["allocation_size"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyProjectIpAddressAllocationTemplate(false, createName, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectIpAddressAllocationExists(updateName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestProjectIpAddressAllocationUpdateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "allocation_ips"),
					resource.TestCheckResourceAttr(testResourceName, "allocation_size", accTestProjectIpAddressAllocationUpdateAttributes["allocation_size"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyProjectIpAddressAllocationMinimalistic(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectIpAddressAllocationExists(updateName, testResourceName),
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

func TestAccResourceNsxtPolicyProjectIpAddressAllocation_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_project_ip_address_allocation.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectIpAddressAllocationCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyProjectIpAddressAllocationMinimalistic(name),
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

func TestAccResourceNsxtPolicyProjectIpAddressAllocation_ipv6(t *testing.T) {
	testResourceName := "nsxt_policy_project_ip_address_allocation.test_ipv6"
	createName := getAccTestResourceName()
	updateName := getAccTestResourceName()
	fixtureBase := getAccTestResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectIpAddressAllocationCheckDestroy(state, updateName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyProjectIpAddressAllocationIpv6Template(true, fixtureBase, createName, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectIpAddressAllocationExists(createName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestProjectIpAddressAllocationCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_address_type", "IPV6"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_allocation_prefix_length", "64"),
					resource.TestCheckResourceAttrSet(testResourceName, "allocation_ips"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrPair("data.nsxt_policy_project_ip_address_allocation.test_ipv6", "path", testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyProjectIpAddressAllocationIpv6Template(false, fixtureBase, createName, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectIpAddressAllocationExists(updateName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestProjectIpAddressAllocationUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_address_type", "IPV6"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_allocation_prefix_length", "64"),
					resource.TestCheckResourceAttrSet(testResourceName, "allocation_ips"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrPair("data.nsxt_policy_project_ip_address_allocation.test_ipv6", "path", testResourceName, "path"),
				),
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

func testAccNsxtPolicyProjectIpAddressAllocationExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy ProjectIpAddressAllocation resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy ProjectIpAddressAllocation resource ID not set in resources")
		}

		sessionContext := testAccGetSessionContextFromParentPath(testAccProvider.Meta(), rs.Primary.Attributes["path"])
		if sessionContext.ProjectID == "" {
			sessionContext = testAccGetSessionContext()
		}
		exists, err := resourceNsxtPolicyProjectIpAddressAllocationExists(sessionContext, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy ProjectIpAddressAllocation %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyProjectIpAddressAllocationCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_project_ip_address_allocation" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		sessionContext := testAccGetSessionContextFromParentPath(testAccProvider.Meta(), rs.Primary.Attributes["path"])
		if sessionContext.ProjectID == "" {
			sessionContext = testAccGetSessionContext()
		}
		exists, err := resourceNsxtPolicyProjectIpAddressAllocationExists(sessionContext, resourceID, connector)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy ProjectIpAddressAllocation %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyProjectIpAddressAllocationTemplate(createFlow bool, createName, updateName string) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestProjectIpAddressAllocationCreateAttributes
	} else {
		attrMap = accTestProjectIpAddressAllocationUpdateAttributes
	}
	displayName := updateName
	if createFlow {
		displayName = createName
	}
	return fmt.Sprintf(`
data "nsxt_policy_project" "test" {
  id = "%s"
}

resource "nsxt_policy_project_ip_address_allocation" "test" {
  %s
  display_name    = "%s"
  description     = "%s"
  allocation_size = %s
  ip_block        = data.nsxt_policy_project.test.external_ipv4_blocks[0]
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_project_ip_address_allocation" "test" {
  %s
  allocation_ips = nsxt_policy_project_ip_address_allocation.test.allocation_ips
}`, os.Getenv("NSXT_VPC_PROJECT_ID"), testAccNsxtProjectContext(), displayName, attrMap["description"], attrMap["allocation_size"], testAccNsxtProjectContext())
}

func testAccNsxtPolicyProjectIpAddressAllocationMinimalistic(updateName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_project" "test" {
  id = "%s"
}

resource "nsxt_policy_project_ip_address_allocation" "test" {
  %s
  display_name    = "%s"
  allocation_size = %s
  ip_block        = data.nsxt_policy_project.test.external_ipv4_blocks[0]
}

data "nsxt_policy_project_ip_address_allocation" "test" {
  %s
  allocation_ips = nsxt_policy_project_ip_address_allocation.test.allocation_ips
}`, os.Getenv("NSXT_VPC_PROJECT_ID"), testAccNsxtProjectContext(), updateName, accTestProjectIpAddressAllocationUpdateAttributes["allocation_size"], testAccNsxtProjectContext())
}

func testAccNsxtPolicyProjectIpAddressAllocationIpv6Template(createFlow bool, fixtureBase, createName, updateName string) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestProjectIpAddressAllocationCreateAttributes
	} else {
		attrMap = accTestProjectIpAddressAllocationUpdateAttributes
	}
	displayName := updateName
	if createFlow {
		displayName = createName
	}
	return fmt.Sprintf(`
resource "nsxt_policy_ip_block" "ipv6_block" {
  display_name = "%s-v6-block"
  cidr         = "2001:db8:beef::/48"
  visibility   = "EXTERNAL"
}

resource "nsxt_policy_project" "ipv6_proj" {
  display_name = "%s-v6-proj"
  ipv6_blocks  = [nsxt_policy_ip_block.ipv6_block.path]
}

resource "nsxt_policy_project_ip_address_allocation" "test_ipv6" {
  context {
    project_id = nsxt_policy_project.ipv6_proj.id
  }
  display_name                  = "%s"
  description                   = "%s"
  ip_address_type               = "IPV6"
  ipv6_allocation_prefix_length = 64
  ip_block                      = nsxt_policy_ip_block.ipv6_block.path
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_project_ip_address_allocation" "test_ipv6" {
  context {
    project_id = nsxt_policy_project.ipv6_proj.id
  }
  allocation_ips = nsxt_policy_project_ip_address_allocation.test_ipv6.allocation_ips
}
`, fixtureBase, fixtureBase, displayName, attrMap["description"])
}
