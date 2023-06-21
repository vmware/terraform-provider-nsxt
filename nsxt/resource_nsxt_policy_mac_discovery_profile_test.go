/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var accTestPolicyMacDiscoveryProfileCreateAttributes = map[string]string{
	"display_name":                     getAccTestResourceName(),
	"description":                      "terraform created",
	"mac_change_enabled":               "true",
	"mac_learning_enabled":             "true",
	"mac_limit":                        "2",
	"mac_limit_policy":                 "ALLOW",
	"remote_overlay_mac_limit":         "2048",
	"unknown_unicast_flooding_enabled": "true",
}

var accTestPolicyMacDiscoveryProfileUpdateAttributes = map[string]string{
	"display_name":                     getAccTestResourceName(),
	"description":                      "terraform updated",
	"mac_change_enabled":               "false",
	"mac_learning_enabled":             "false",
	"mac_limit":                        "5",
	"mac_limit_policy":                 "DROP",
	"remote_overlay_mac_limit":         "4096",
	"unknown_unicast_flooding_enabled": "false",
}

func TestAccResourceNsxtPolicyMacDiscoveryProfile_basic(t *testing.T) {
	testResourceName := "nsxt_policy_mac_discovery_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyMacDiscoveryProfileCheckDestroy(state, accTestPolicyMacDiscoveryProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyMacDiscoveryProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyMacDiscoveryProfileExists(accTestPolicyMacDiscoveryProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyMacDiscoveryProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyMacDiscoveryProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "mac_change_enabled", accTestPolicyMacDiscoveryProfileCreateAttributes["mac_change_enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "mac_learning_enabled", accTestPolicyMacDiscoveryProfileCreateAttributes["mac_learning_enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "mac_limit", accTestPolicyMacDiscoveryProfileCreateAttributes["mac_limit"]),
					resource.TestCheckResourceAttr(testResourceName, "mac_limit_policy", accTestPolicyMacDiscoveryProfileCreateAttributes["mac_limit_policy"]),
					resource.TestCheckResourceAttr(testResourceName, "remote_overlay_mac_limit", accTestPolicyMacDiscoveryProfileCreateAttributes["remote_overlay_mac_limit"]),
					resource.TestCheckResourceAttr(testResourceName, "unknown_unicast_flooding_enabled", accTestPolicyMacDiscoveryProfileCreateAttributes["unknown_unicast_flooding_enabled"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyMacDiscoveryProfileTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyMacDiscoveryProfileExists(accTestPolicyMacDiscoveryProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyMacDiscoveryProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyMacDiscoveryProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "mac_change_enabled", accTestPolicyMacDiscoveryProfileUpdateAttributes["mac_change_enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "mac_learning_enabled", accTestPolicyMacDiscoveryProfileUpdateAttributes["mac_learning_enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "mac_limit", accTestPolicyMacDiscoveryProfileUpdateAttributes["mac_limit"]),
					resource.TestCheckResourceAttr(testResourceName, "mac_limit_policy", accTestPolicyMacDiscoveryProfileUpdateAttributes["mac_limit_policy"]),
					resource.TestCheckResourceAttr(testResourceName, "remote_overlay_mac_limit", accTestPolicyMacDiscoveryProfileUpdateAttributes["remote_overlay_mac_limit"]),
					resource.TestCheckResourceAttr(testResourceName, "unknown_unicast_flooding_enabled", accTestPolicyMacDiscoveryProfileUpdateAttributes["unknown_unicast_flooding_enabled"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyMacDiscoveryProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyMacDiscoveryProfileExists(accTestPolicyMacDiscoveryProfileCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyMacDiscoveryProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_mac_discovery_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyMacDiscoveryProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyMacDiscoveryProfileMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyMacDiscoveryProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy MacDiscoveryProfile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy MacDiscoveryProfile resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyMacDiscoveryProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy MacDiscoveryProfile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyMacDiscoveryProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_mac_discovery_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyMacDiscoveryProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy MacDiscoveryProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyMacDiscoveryProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyMacDiscoveryProfileCreateAttributes
	} else {
		attrMap = accTestPolicyMacDiscoveryProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_mac_discovery_profile" "test" {
  display_name = "%s"
  description  = "%s"
  mac_change_enabled = %s
  mac_learning_enabled = %s
  mac_limit = %s
  mac_limit_policy = "%s"
  remote_overlay_mac_limit = %s
  unknown_unicast_flooding_enabled = %s
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["mac_change_enabled"], attrMap["mac_learning_enabled"], attrMap["mac_limit"], attrMap["mac_limit_policy"], attrMap["remote_overlay_mac_limit"], attrMap["unknown_unicast_flooding_enabled"])
}

func testAccNsxtPolicyMacDiscoveryProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_mac_discovery_profile" "test" {
  display_name = "%s"
}`, accTestPolicyMacDiscoveryProfileUpdateAttributes["display_name"])
}
