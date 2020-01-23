/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/ip_pools"
	"testing"
)

var accTestPolicyIPAddressAllocationCreateAttributes = map[string]string{
	"display_name":  "terra-test",
	"description":   "terraform created",
	"allocation_ip": "12.12.12.11",
}

var accTestPolicyIPAddressAllocationUpdateAttributes = map[string]string{
	"display_name":  "terra-test-updated",
	"description":   "terraform updated",
	"allocation_ip": "12.12.12.12",
}

func TestAccResourceNsxtPolicyIPAddressAllocation_basic(t *testing.T) {
	testResourceName := "nsxt_policy_ip_address_allocation.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPAddressAllocationCheckDestroy(state, accTestPolicyIPAddressAllocationCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPAddressAllocationTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPAddressAllocationExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPAddressAllocationCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPAddressAllocationCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "allocation_ip", accTestPolicyIPAddressAllocationCreateAttributes["allocation_ip"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationDependenciesTemplate(),
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPAddressAllocationExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPAddressAllocationUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPAddressAllocationUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "allocation_ip", accTestPolicyIPAddressAllocationUpdateAttributes["allocation_ip"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationDependenciesTemplate(),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIPAddressAllocation_anyIPBasic(t *testing.T) {
	testResourceName := "nsxt_policy_ip_address_allocation.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPAddressAllocationCheckDestroy(state, accTestPolicyIPAddressAllocationCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPAddressAllocationAnyFreeIPTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPAddressAllocationExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPAddressAllocationCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPAddressAllocationCreateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "allocation_ip"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationDependenciesTemplate(),
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationAnyFreeIPTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPAddressAllocationExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPAddressAllocationUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPAddressAllocationUpdateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "allocation_ip"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationDependenciesTemplate(),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIPAddressAllocation_importBasic(t *testing.T) {
	name := "terra-test-import"
	testResourceName := "nsxt_policy_ip_address_allocation.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPAddressAllocationCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPAddressAllocationTemplate(true),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNSXPolicyIPAddressAllocationImporterGetID,
			},
			{
				Config: testAccNsxtPolicyIPAddressAllocationDependenciesTemplate(),
			},
		},
	})
}

func testAccNSXPolicyIPAddressAllocationImporterGetID(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["nsxt_policy_ip_address_allocation.test"]
	if !ok {
		return "", fmt.Errorf("NSX Policy IP Allocation resource %s not found in resources", "nsxt_policy_ip_address_allocation.test")
	}
	resourceID := rs.Primary.ID
	if resourceID == "" {
		return "", fmt.Errorf("NSX Policy IP Allocation resource ID not set in resources ")
	}
	poolPath := rs.Primary.Attributes["pool_path"]
	if poolPath == "" {
		return "", fmt.Errorf("NSX Policy IP Allocation pool_path not set in resources ")
	}
	poolID := getPolicyIDFromPath(poolPath)
	return fmt.Sprintf("%s/%s", poolID, resourceID), nil
}

func testAccNsxtPolicyIPAddressAllocationExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		nsxClient := ip_pools.NewDefaultIpAllocationsClient(connector)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy IPAddressAllocation resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy IPAddressAllocation resource ID not set in resources")
		}

		poolPath := rs.Primary.Attributes["pool_path"]
		if poolPath == "" {
			return fmt.Errorf("No pool_path found for IP Address Allocation with ID %s", resourceID)
		}
		poolID, err := resourceNsxtPolicyIPAddressParsePoolIDFromPath(poolPath, connector)
		if err != nil {
			return err
		}

		_, err = nsxClient.Get(poolID, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving policy IPAddressAllocation ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyIPAddressAllocationCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	nsxClient := ip_pools.NewDefaultIpAllocationsClient(connector)
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_ip_address_allocation" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		poolPath := rs.Primary.Attributes["pool_path"]
		if poolPath == "" {
			return fmt.Errorf("No pool_path found for IP Address Allocation with ID %s", resourceID)
		}
		poolID, err := resourceNsxtPolicyIPAddressParsePoolIDFromPath(poolPath, connector)
		if err != nil {
			return err
		}

		_, err = nsxClient.Get(poolID, resourceID)
		if err == nil {
			return fmt.Errorf("Policy IPAddressAllocation %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyIPAddressAllocationTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyIPAddressAllocationCreateAttributes
	} else {
		attrMap = accTestPolicyIPAddressAllocationUpdateAttributes
	}
	return testAccNsxtPolicyIPAddressAllocationDependenciesTemplate() + fmt.Sprintf(`

resource "nsxt_policy_ip_address_allocation" "test" {
  display_name  = "%s"
  description   = "%s"
  allocation_ip = "%s"
  pool_path     = nsxt_policy_ip_pool.test.path

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_ip_address_allocation.test.path
}`, attrMap["display_name"], attrMap["description"], attrMap["allocation_ip"])
}

func testAccNsxtPolicyIPAddressAllocationAnyFreeIPTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyIPAddressAllocationCreateAttributes
	} else {
		attrMap = accTestPolicyIPAddressAllocationUpdateAttributes
	}
	return testAccNsxtPolicyIPAddressAllocationDependenciesTemplate() + fmt.Sprintf(`

resource "nsxt_policy_ip_address_allocation" "test" {
  display_name  = "%s"
  description   = "%s"
  pool_path     = nsxt_policy_ip_pool.test.path

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_ip_address_allocation.test.path
}`, attrMap["display_name"], attrMap["description"])
}

func testAccNsxtPolicyIPAddressAllocationDependenciesTemplate() string {
	return fmt.Sprintf(`
resource "nsxt_policy_ip_pool" "test" {
  display_name = "tfpool1"
}

resource "nsxt_policy_ip_pool_static_subnet" "test" {
  display_name = "tfsnet1"
  pool_path    = nsxt_policy_ip_pool.test.path
  cidr         = "12.12.12.0/24"
  allocation_range {
    start = "12.12.12.10"
    end   = "12.12.12.20"
  }
}`)
}
