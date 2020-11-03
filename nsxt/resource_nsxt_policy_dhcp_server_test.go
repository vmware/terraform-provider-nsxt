/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var accTestPolicyDhcpServerCreateAttributes = map[string]string{
	"display_name": "terra-test",
	"description":  "terraform created",
	"lease_time":   "200",
}

var accTestPolicyDhcpServerUpdateAttributes = map[string]string{
	"display_name": "terra-test-updated",
	"description":  "terraform updated",
	"lease_time":   "500",
}

func TestAccResourceNsxtPolicyDhcpServer_basic(t *testing.T) {
	testResourceName := "nsxt_policy_dhcp_server.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDhcpServerCheckDestroy(state, accTestPolicyDhcpServerCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDhcpServerCreateTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDhcpServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDhcpServerCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDhcpServerCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "lease_time", accTestPolicyDhcpServerCreateAttributes["lease_time"]),

					resource.TestCheckResourceAttrSet(testResourceName, "edge_cluster_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "preferred_edge_paths.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "server_addresses.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyDhcpServerUpdateTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDhcpServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDhcpServerUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDhcpServerUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "lease_time", accTestPolicyDhcpServerUpdateAttributes["lease_time"]),

					resource.TestCheckResourceAttrSet(testResourceName, "edge_cluster_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "preferred_edge_paths.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "server_addresses.#", "2"),
				),
			},
			{
				Config: testAccNsxtPolicyDhcpServerMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDhcpServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "preferred_edge_paths.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "server_addresses.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyDhcpServer_importBasic(t *testing.T) {
	name := "terra-test-import"
	testResourceName := "nsxt_policy_dhcp_server.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccOnlyLocalManager(t); testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDhcpServerCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDhcpServerMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyDhcpServerExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy DhcpServer resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy DhcpServer resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyDhcpServerExists(resourceID, connector, testAccIsGlobalManager())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Error while retrieving policy DhcpServer ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyDhcpServerCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_dhcp_server" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyDhcpServerExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy DhcpServer %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyDhcpServerCreateTemplate() string {
	attrMap := accTestPolicyDhcpServerCreateAttributes

	return testAccNsxtPolicyGatewayFabricDeps(false) + fmt.Sprintf(`

resource "nsxt_policy_dhcp_server" "test" {
  display_name = "%s"
  description  = "%s"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
  lease_time = %s
  server_addresses = ["110.64.0.1/16"]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
`, attrMap["display_name"], attrMap["description"], attrMap["lease_time"])
}

func testAccNsxtPolicyDhcpServerUpdateTemplate() string {
	attrMap := accTestPolicyDhcpServerUpdateAttributes

	return testAccNsxtPolicyGatewayFabricDeps(false) + fmt.Sprintf(`
resource "nsxt_policy_dhcp_server" "test" {
  display_name = "%s"
  description  = "%s"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
  lease_time = %s
  server_addresses = ["2001::1234:abcd:ffff:c0a8:101/64", "110.64.0.1/16"]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
`, attrMap["display_name"], attrMap["description"], attrMap["lease_time"])
}

func testAccNsxtPolicyDhcpServerMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_dhcp_server" "test" {
  display_name = "%s"
  server_addresses = ["110.64.0.1/16"]
}
`, accTestPolicyDhcpServerUpdateAttributes["display_name"])
}
