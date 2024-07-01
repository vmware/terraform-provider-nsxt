/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var accTestPolicyVpcSubnetCreateAttributes = map[string]string{
	"display_name":           getAccTestResourceName(),
	"description":            "terraform created",
	"ipv4_subnet_size":       "2",
	"ip_addresses":           "test-create",
	"access_mode":            "Private",
	"enabled":                "true",
	"dhcp_relay_config_path": "test-create",
	"enable_dhcp":            "true",
	"dns_server_ips":         "test-create",
	"ipv4_pool_size":         "2",
}

var accTestPolicyVpcSubnetUpdateAttributes = map[string]string{
	"display_name":           getAccTestResourceName(),
	"description":            "terraform updated",
	"ipv4_subnet_size":       "5",
	"ip_addresses":           "test-update",
	"access_mode":            "Public",
	"enabled":                "false",
	"dhcp_relay_config_path": "test-update",
	"enable_dhcp":            "false",
	"dns_server_ips":         "test-update",
	"ipv4_pool_size":         "5",
}

func TestAccResourceNsxtPolicyVpcSubnet_basic(t *testing.T) {
	testResourceName := "nsxt_policy_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyVpcSubnetCheckDestroy(state, accTestPolicyVpcSubnetUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVpcSubnetTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVpcSubnetExists(accTestPolicyVpcSubnetCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyVpcSubnetCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyVpcSubnetCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ipv4_subnet_size", accTestPolicyVpcSubnetCreateAttributes["ipv4_subnet_size"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestPolicyVpcSubnetCreateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "access_mode", accTestPolicyVpcSubnetCreateAttributes["access_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.0.static_ip_allocation.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "static_ip_allocation.0.enabled", accTestPolicyVpcSubnetCreateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_relay_config_path", accTestPolicyVpcSubnetCreateAttributes["dhcp_relay_config_path"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dns_client_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.enable_dhcp", accTestPolicyVpcSubnetCreateAttributes["enable_dhcp"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.static_pool_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dns_client_config.0.dns_server_ips.0", accTestPolicyVpcSubnetCreateAttributes["dns_server_ips"]),
					resource.TestCheckResourceAttr(testResourceName, "static_pool_config.0.ipv4_pool_size", accTestPolicyVpcSubnetCreateAttributes["ipv4_pool_size"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyVpcSubnetTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVpcSubnetExists(accTestPolicyVpcSubnetUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyVpcSubnetUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyVpcSubnetUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ipv4_subnet_size", accTestPolicyVpcSubnetUpdateAttributes["ipv4_subnet_size"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestPolicyVpcSubnetUpdateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "access_mode", accTestPolicyVpcSubnetUpdateAttributes["access_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "advanced_config.0.static_ip_allocation.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "static_ip_allocation.0.enabled", accTestPolicyVpcSubnetUpdateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dhcp_relay_config_path", accTestPolicyVpcSubnetUpdateAttributes["dhcp_relay_config_path"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.dns_client_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.enable_dhcp", accTestPolicyVpcSubnetUpdateAttributes["enable_dhcp"]),
					resource.TestCheckResourceAttr(testResourceName, "dhcp_config.0.static_pool_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "dns_client_config.0.dns_server_ips.0", accTestPolicyVpcSubnetUpdateAttributes["dns_server_ips"]),
					resource.TestCheckResourceAttr(testResourceName, "static_pool_config.0.ipv4_pool_size", accTestPolicyVpcSubnetUpdateAttributes["ipv4_pool_size"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyVpcSubnetMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyVpcSubnetExists(accTestPolicyVpcSubnetCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyVpcSubnet_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyVpcSubnetCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyVpcSubnetMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyVpcSubnetExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy VpcSubnet resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy VpcSubnet resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyVpcSubnetExists(resourceID, connector, testAccIsGlobalManager())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy VpcSubnet %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyVpcSubnetCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_subnet" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyVpcSubnetExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy VpcSubnet %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyVpcSubnetTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyVpcSubnetCreateAttributes
	} else {
		attrMap = accTestPolicyVpcSubnetUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_subnet" "test" {
  display_name = "%s"
  description  = "%s"

  advanced_config {

    static_ip_allocation {
      enabled = %s
    }

  }

  ipv4_subnet_size = %s
  ip_addresses = [%s]
  access_mode = %s

  dhcp_config {
    dhcp_relay_config_path = %s

    dns_client_config {
      dns_server_ips = [%s]
    }

    enable_dhcp = %s

    static_pool_config {
      ipv4_pool_size = %s
    }

  }


  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["ipv4_subnet_size"], attrMap["ip_addresses"], attrMap["access_mode"], attrMap["enabled"], attrMap["dhcp_relay_config_path"], attrMap["enable_dhcp"], attrMap["dns_server_ips"], attrMap["ipv4_pool_size"])
}

func testAccNsxtPolicyVpcSubnetMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_subnet" "test" {
  display_name = "%s"

}`, accTestPolicyVpcSubnetUpdateAttributes["display_name"])
}
