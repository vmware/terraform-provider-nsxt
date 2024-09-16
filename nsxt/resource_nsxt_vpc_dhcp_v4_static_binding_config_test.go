/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var accTestPolicyDhcpV4StaticBindingConfigCreateAttributes = map[string]string{
	"display_name":    getAccTestResourceName(),
	"description":     "terraform created",
	"resource_type":   "DhcpV4StaticBindingConfig",
	"gateway_address": "test-create",
	"host_name":       "test-create",
	"mac_address":     "test-create",
	"lease_time":      "2",
	"ip_address":      "test-create",
	"next_hop":        "test-create",
	"network":         "test-create",
	"code":            "2",
	"values":          "test-create",
}

var accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes = map[string]string{
	"display_name":    getAccTestResourceName(),
	"description":     "terraform updated",
	"resource_type":   "DhcpV6StaticBindingConfig",
	"gateway_address": "test-update",
	"host_name":       "test-update",
	"mac_address":     "test-update",
	"lease_time":      "5",
	"ip_address":      "test-update",
	"next_hop":        "test-update",
	"network":         "test-update",
	"code":            "5",
	"values":          "test-update",
}

func TestAccResourceNsxtPolicyDhcpV4StaticBindingConfig_basic(t *testing.T) {
	testResourceName := "nsxt_policy_dhcp_v4_static_binding_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDhcpV4StaticBindingConfigCheckDestroy(state, accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDhcpV4StaticBindingConfigTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDhcpV4StaticBindingConfigExists(accTestPolicyDhcpV4StaticBindingConfigCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDhcpV4StaticBindingConfigCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDhcpV4StaticBindingConfigCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "resource_type", accTestPolicyDhcpV4StaticBindingConfigCreateAttributes["resource_type"]),
					resource.TestCheckResourceAttr(testResourceName, "gateway_address", accTestPolicyDhcpV4StaticBindingConfigCreateAttributes["gateway_address"]),
					resource.TestCheckResourceAttr(testResourceName, "host_name", accTestPolicyDhcpV4StaticBindingConfigCreateAttributes["host_name"]),
					resource.TestCheckResourceAttr(testResourceName, "mac_address", accTestPolicyDhcpV4StaticBindingConfigCreateAttributes["mac_address"]),
					resource.TestCheckResourceAttr(testResourceName, "lease_time", accTestPolicyDhcpV4StaticBindingConfigCreateAttributes["lease_time"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_address", accTestPolicyDhcpV4StaticBindingConfigCreateAttributes["ip_address"]),
					resource.TestCheckResourceAttr(testResourceName, "options.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "options.0.option121.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "options.0.other.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "option121.0.static_route.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "static_route.0.next_hop", accTestPolicyDhcpV4StaticBindingConfigCreateAttributes["next_hop"]),
					resource.TestCheckResourceAttr(testResourceName, "static_route.0.network", accTestPolicyDhcpV4StaticBindingConfigCreateAttributes["network"]),
					resource.TestCheckResourceAttr(testResourceName, "other.0.code", accTestPolicyDhcpV4StaticBindingConfigCreateAttributes["code"]),
					resource.TestCheckResourceAttr(testResourceName, "other.0.values.0", accTestPolicyDhcpV4StaticBindingConfigCreateAttributes["values"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyDhcpV4StaticBindingConfigTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDhcpV4StaticBindingConfigExists(accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "resource_type", accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["resource_type"]),
					resource.TestCheckResourceAttr(testResourceName, "gateway_address", accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["gateway_address"]),
					resource.TestCheckResourceAttr(testResourceName, "host_name", accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["host_name"]),
					resource.TestCheckResourceAttr(testResourceName, "mac_address", accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["mac_address"]),
					resource.TestCheckResourceAttr(testResourceName, "lease_time", accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["lease_time"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_address", accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["ip_address"]),
					resource.TestCheckResourceAttr(testResourceName, "options.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "options.0.option121.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "options.0.other.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "option121.0.static_route.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "static_route.0.next_hop", accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["next_hop"]),
					resource.TestCheckResourceAttr(testResourceName, "static_route.0.network", accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["network"]),
					resource.TestCheckResourceAttr(testResourceName, "other.0.code", accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["code"]),
					resource.TestCheckResourceAttr(testResourceName, "other.0.values.0", accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["values"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyDhcpV4StaticBindingConfigMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyDhcpV4StaticBindingConfigExists(accTestPolicyDhcpV4StaticBindingConfigCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyDhcpV4StaticBindingConfig_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_dhcp_v4_static_binding_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyDhcpV4StaticBindingConfigCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyDhcpV4StaticBindingConfigMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyDhcpV4StaticBindingConfigExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy DhcpV4StaticBindingConfig resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy DhcpV4StaticBindingConfig resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyDhcpV4StaticBindingConfigExists(resourceID, connector, testAccIsGlobalManager())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy DhcpV4StaticBindingConfig %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyDhcpV4StaticBindingConfigCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_dhcp_v4_static_binding_config" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyDhcpV4StaticBindingConfigExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy DhcpV4StaticBindingConfig %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyDhcpV4StaticBindingConfigTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyDhcpV4StaticBindingConfigCreateAttributes
	} else {
		attrMap = accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_dhcp_v4_static_binding_config" "test" {
  display_name = "%s"
  description  = "%s"
  resource_type = %s
  gateway_address = %s
  host_name = %s
  mac_address = %s
  lease_time = %s
  ip_address = %s

  options {

    option121 {

      static_route {
        next_hop = %s
        network = %s
      }

    }


    other {
      code = %s
      values = [%s]
    }

  }


  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["resource_type"], attrMap["gateway_address"], attrMap["host_name"], attrMap["mac_address"], attrMap["lease_time"], attrMap["ip_address"], attrMap["next_hop"], attrMap["network"], attrMap["code"], attrMap["values"])
}

func testAccNsxtPolicyDhcpV4StaticBindingConfigMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_dhcp_v4_static_binding_config" "test" {
  display_name = "%s"
  resource_type = %s

}`, accTestPolicyDhcpV4StaticBindingConfigUpdateAttributes["display_name"], attrMap["resource_type"])
}
