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

var nsxtPolicyBlockSubnetPoolID = "tfpool1"

// TODO: remove extra test step config once IP Blocks don't need a delay to delete
func TestAccResourceNsxtPolicyIPPoolBlockSubnet_minimal(t *testing.T) {
	name := "blocksubnet1"
	testResourceName := "nsxt_policy_ip_pool_block_subnet.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXPolicyIPPoolBlockSubnetCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyIPPoolBlockSubnetCreateMinimalTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXPolicyIPPoolBlockSubnetCheckExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "size", "4"),
					resource.TestCheckResourceAttr(testResourceName, "auto_assign_gateway", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "block_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
				),
			},
			{
				Config: testAccNSXPolicyIPPoolBlockSubnetIPBlockTemplate(),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIPPoolBlockSubnet_basic(t *testing.T) {
	name := "tfpool1"
	updatedName := fmt.Sprintf("%s-updated", name)
	testResourceName := "nsxt_policy_ip_pool_block_subnet.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXPolicyIPPoolBlockSubnetCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyIPPoolBlockSubnetCreateTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXPolicyIPPoolBlockSubnetCheckExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "size", "4"),
					resource.TestCheckResourceAttr(testResourceName, "auto_assign_gateway", "false"),
					resource.TestCheckResourceAttrSet(testResourceName, "block_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
				),
			},
			{
				Config: testAccNSXPolicyIPPoolBlockSubnetUpdateTemplate(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXPolicyIPPoolBlockSubnetCheckExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "size", "4"),
					resource.TestCheckResourceAttr(testResourceName, "auto_assign_gateway", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "block_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
				),
			},
			{
				Config: testAccNSXPolicyIPPoolBlockSubnetIPBlockTemplate(),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIPPoolBlockSubnet_import_basic(t *testing.T) {
	name := "tfpool1"
	testResourceName := "nsxt_policy_ip_pool_block_subnet.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXPolicyIPPoolBlockSubnetCheckDestroy(state)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyIPPoolBlockSubnetCreateTemplate(name),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNSXPolicyIPPoolBlockSubnetImporterGetID,
			},
			{
				Config: testAccNSXPolicyIPPoolBlockSubnetIPBlockTemplate(),
			},
		},
	})
}

func testAccNSXPolicyIPPoolBlockSubnetImporterGetID(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["nsxt_policy_ip_pool_block_subnet.test"]
	if !ok {
		return "", fmt.Errorf("NSX Policy Block Subnet resource %s not found in resources", "nsxt_policy_ip_pool_block_subnet.test")
	}
	resourceID := rs.Primary.ID
	if resourceID == "" {
		return "", fmt.Errorf("NSX Policy Block Subnet resource ID not set in resources ")
	}
	poolPath := rs.Primary.Attributes["pool_path"]
	if poolPath == "" {
		return "", fmt.Errorf("NSX Policy Block Subnet pool_path not set in resources ")
	}
	poolID := getPolicyIDFromPath(poolPath)
	return fmt.Sprintf("%s/%s", poolID, resourceID), nil
}

func testAccNSXPolicyIPPoolBlockSubnetCheckExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
		client := ip_pools.NewDefaultIpSubnetsClient(connector)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX Policy Block Subnet resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX Policy Block Subnet resource ID not set in resources")
		}

		_, err := client.Get(nsxtPolicyBlockSubnetPoolID, resourceID)
		if err != nil {
			return fmt.Errorf("Failed to find Block Subnet %s", resourceID)
		}

		return nil
	}
}

func testAccNSXPolicyIPPoolBlockSubnetCheckDestroy(state *terraform.State) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	client := ip_pools.NewDefaultIpSubnetsClient(connector)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_ip_pool_block_subnet" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		_, err := client.Get(nsxtPolicyBlockSubnetPoolID, resourceID)
		if err == nil {
			return fmt.Errorf("Block Subnet still exists %s", resourceID)
		}

	}
	return nil
}

func testAccNSXPolicyIPPoolBlockSubnetIPBlockTemplate() string {
	return fmt.Sprintf(`
resource "nsxt_policy_ip_block" "block1" {
  display_name = "tfblock1"
  cidr         = "11.11.12.0/24"
}`)
}

func testAccNSXPolicyIPPoolBlockSubnetCreateMinimalTemplate(name string) string {
	return testAccNSXPolicyIPPoolBlockSubnetIPBlockTemplate() + fmt.Sprintf(`
resource "nsxt_policy_ip_pool" "pool1" {
  display_name = "%s"
  nsx_id       = "%s"
}

resource "nsxt_policy_ip_pool_block_subnet" "test" {
  display_name = "%s"
  size         = 4
  pool_path    = nsxt_policy_ip_pool.pool1.path
  block_path   = nsxt_policy_ip_block.block1.path
}`, nsxtPolicyBlockSubnetPoolID, nsxtPolicyBlockSubnetPoolID, name)
}

func testAccNSXPolicyIPPoolBlockSubnetCreateTemplate(name string) string {
	return testAccNSXPolicyIPPoolBlockSubnetIPBlockTemplate() + fmt.Sprintf(`
resource "nsxt_policy_ip_pool" "pool1" {
  display_name = "%s"
  nsx_id       = "%s"
}

resource "nsxt_policy_ip_pool_block_subnet" "test" {
  display_name        = "%s"
  description         = "Acceptance Test"
  size                = 4
  auto_assign_gateway = false
  pool_path           = nsxt_policy_ip_pool.pool1.path
  block_path          = nsxt_policy_ip_block.block1.path
  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}`, nsxtPolicyBlockSubnetPoolID, nsxtPolicyBlockSubnetPoolID, name)
}

func testAccNSXPolicyIPPoolBlockSubnetUpdateTemplate(name string) string {
	return testAccNSXPolicyIPPoolBlockSubnetIPBlockTemplate() + fmt.Sprintf(`
resource "nsxt_policy_ip_pool" "pool1" {
  display_name = "%s"
  nsx_id       = "%s"
}

resource "nsxt_policy_ip_pool_block_subnet" "test" {
  display_name        = "%s"
  description         = "Acceptance Test"
  size                = 4
  auto_assign_gateway = true
  pool_path           = nsxt_policy_ip_pool.pool1.path
  block_path          = nsxt_policy_ip_block.block1.path
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}`, nsxtPolicyBlockSubnetPoolID, nsxtPolicyBlockSubnetPoolID, name)
}
