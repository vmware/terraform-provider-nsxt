// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var accTestPolicyLBGenericPersistenceProfileCreateAttributes = map[string]string{
	"display_name":                     getAccTestResourceName(),
	"description":                      "terraform created",
	"persistence_shared":               "true",
	"timeout":                          "2",
	"ha_persistence_mirroring_enabled": "true",
}

var accTestPolicyLBGenericPersistenceProfileUpdateAttributes = map[string]string{
	"display_name":                     getAccTestResourceName(),
	"description":                      "terraform updated",
	"persistence_shared":               "false",
	"timeout":                          "5",
	"ha_persistence_mirroring_enabled": "false",
}

func TestAccResourceNsxtPolicyLBGenericPersistenceProfile_basic(t *testing.T) {
	testResourceName := "nsxt_policy_lb_generic_persistence_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBGenericPersistenceProfileCheckDestroy(state, accTestPolicyLBGenericPersistenceProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBGenericPersistenceProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBPersistenceProfileExists(accTestPolicyLBGenericPersistenceProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBGenericPersistenceProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBGenericPersistenceProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "persistence_shared", accTestPolicyLBGenericPersistenceProfileCreateAttributes["persistence_shared"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestPolicyLBGenericPersistenceProfileCreateAttributes["timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_persistence_mirroring_enabled", accTestPolicyLBGenericPersistenceProfileCreateAttributes["ha_persistence_mirroring_enabled"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBGenericPersistenceProfileTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBPersistenceProfileExists(accTestPolicyLBGenericPersistenceProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBGenericPersistenceProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBGenericPersistenceProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "persistence_shared", accTestPolicyLBGenericPersistenceProfileUpdateAttributes["persistence_shared"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestPolicyLBGenericPersistenceProfileUpdateAttributes["timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_persistence_mirroring_enabled", accTestPolicyLBGenericPersistenceProfileUpdateAttributes["ha_persistence_mirroring_enabled"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBGenericPersistenceProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBPersistenceProfileExists(accTestPolicyLBGenericPersistenceProfileCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyLBGenericPersistenceProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_lb_generic_persistence_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBGenericPersistenceProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBGenericPersistenceProfileMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyLBGenericPersistenceProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_lb_generic_persistence_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyLBPersistenceProfileExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy LBGenericPersistenceProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyLBGenericPersistenceProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyLBGenericPersistenceProfileCreateAttributes
	} else {
		attrMap = accTestPolicyLBGenericPersistenceProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_lb_generic_persistence_profile" "test" {
  display_name = "%s"
  description  = "%s"

  persistence_shared = %s
  timeout            = %s
  
  ha_persistence_mirroring_enabled = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["persistence_shared"], attrMap["timeout"], attrMap["ha_persistence_mirroring_enabled"])
}

func testAccNsxtPolicyLBGenericPersistenceProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_lb_generic_persistence_profile" "test" {
  display_name = "%s"
}`, accTestPolicyLBGenericPersistenceProfileUpdateAttributes["display_name"])
}
