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

var accTestPolicyLBSourceIpPersistenceProfileCreateAttributes = map[string]string{
	"display_name":                     getAccTestResourceName(),
	"description":                      "terraform created",
	"persistence_shared":               "true",
	"purge":                            "NO_PURGE",
	"timeout":                          "2",
	"ha_persistence_mirroring_enabled": "true",
}

var accTestPolicyLBSourceIpPersistenceProfileUpdateAttributes = map[string]string{
	"display_name":                     getAccTestResourceName(),
	"description":                      "terraform updated",
	"persistence_shared":               "false",
	"purge":                            "FULL",
	"timeout":                          "5",
	"ha_persistence_mirroring_enabled": "false",
}

func TestAccResourceNsxtPolicyLBSourceIpPersistenceProfile_basic(t *testing.T) {
	testResourceName := "nsxt_policy_lb_source_ip_persistence_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBSourceIpPersistenceProfileCheckDestroy(state, accTestPolicyLBSourceIpPersistenceProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBSourceIpPersistenceProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBPersistenceProfileExists(accTestPolicyLBSourceIpPersistenceProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBSourceIpPersistenceProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBSourceIpPersistenceProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "persistence_shared", accTestPolicyLBSourceIpPersistenceProfileCreateAttributes["persistence_shared"]),
					resource.TestCheckResourceAttr(testResourceName, "purge", accTestPolicyLBSourceIpPersistenceProfileCreateAttributes["purge"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestPolicyLBSourceIpPersistenceProfileCreateAttributes["timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_persistence_mirroring_enabled", accTestPolicyLBSourceIpPersistenceProfileCreateAttributes["ha_persistence_mirroring_enabled"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBSourceIpPersistenceProfileTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBPersistenceProfileExists(accTestPolicyLBSourceIpPersistenceProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBSourceIpPersistenceProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBSourceIpPersistenceProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "persistence_shared", accTestPolicyLBSourceIpPersistenceProfileUpdateAttributes["persistence_shared"]),
					resource.TestCheckResourceAttr(testResourceName, "purge", accTestPolicyLBSourceIpPersistenceProfileUpdateAttributes["purge"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestPolicyLBSourceIpPersistenceProfileUpdateAttributes["timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_persistence_mirroring_enabled", accTestPolicyLBSourceIpPersistenceProfileUpdateAttributes["ha_persistence_mirroring_enabled"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBSourceIpPersistenceProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBPersistenceProfileExists(accTestPolicyLBSourceIpPersistenceProfileCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyLBSourceIpPersistenceProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_lb_source_ip_persistence_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBSourceIpPersistenceProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBSourceIpPersistenceProfileMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyLBSourceIpPersistenceProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_lb_source_ip_persistence_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyLBPersistenceProfileExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy LBSourceIpPersistenceProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyLBSourceIpPersistenceProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyLBSourceIpPersistenceProfileCreateAttributes
	} else {
		attrMap = accTestPolicyLBSourceIpPersistenceProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_lb_source_ip_persistence_profile" "test" {
  display_name = "%s"
  description  = "%s"

  persistence_shared = %s
  purge              = "%s"
  timeout            = %s

  ha_persistence_mirroring_enabled = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["persistence_shared"], attrMap["purge"], attrMap["timeout"], attrMap["ha_persistence_mirroring_enabled"])
}

func testAccNsxtPolicyLBSourceIpPersistenceProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_lb_source_ip_persistence_profile" "test" {
  display_name = "%s"
}`, accTestPolicyLBSourceIpPersistenceProfileUpdateAttributes["display_name"])
}
