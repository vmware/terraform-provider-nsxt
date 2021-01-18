/* Copyright © 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceNsxtPolicyIntrusionServiceProfile_basic(t *testing.T) {
	name := getAccTestResourceName()
	updatedName := getAccTestResourceName()
	testResourceName := "nsxt_policy_intrusion_service_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIntrusionServiceProfileCheckDestroy(state, updatedName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceProfileCreate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServiceProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "severities.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.attack_types.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.attack_targets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.cvss.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.products_affected.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "overridden_signature.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "overridden_signature.0.signature_id"),
					resource.TestCheckResourceAttr(testResourceName, "overridden_signature.0.action", "REJECT"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyIntrusionServiceProfileUpdate(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServiceProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "severities.#", "3"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.attack_types.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.attack_targets.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.cvss.#", "4"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.products_affected.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "overridden_signature.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIntrusionServiceProfileMinimalistic(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServiceProfileExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttr(testResourceName, "severities.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.attack_types.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.attack_targets.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.cvss.#", "0"),
					resource.TestCheckResourceAttr(testResourceName, "criteria.0.products_affected.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "overridden_signature.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIntrusionServiceProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_intrusion_service_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIntrusionServiceProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceProfileMinimalistic(name),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyIntrusionServiceProfileExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyIntrusionServiceProfileExists(resourceID, connector, false)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Error while retrieving policy resource ID %s", resourceID)
		}
		return nil
	}
}

func testAccNsxtPolicyIntrusionServiceProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_intrusion_service_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyIntrusionServiceProfileExists(resourceID, connector, false)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy resource %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyIntrusionServiceProfileCreate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "%s"
  description  = "Acceptance Test"
  severities   = ["HIGH", "CRITICAL"]

  criteria {
    attack_types   = ["trojan-activity", "successful-admin"]
    attack_targets = ["SERVER"]
  }

  overridden_signature {
    signature_id = "2026323"
    action       = "REJECT"
    enabled      = true
  }

  tag {
    scope = "color"
    tag   = "orange"
  }
}`, name)
}

func testAccNsxtPolicyIntrusionServiceProfileUpdate(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "%s"
  description  = "Acceptance Test"
  severities   = ["HIGH", "CRITICAL", "LOW"]

  criteria {
    attack_types      = ["trojan-activity", "successful-admin"]
    cvss              = ["NONE", "MEDIUM", "HIGH", "CRITICAL"]
    products_affected = ["Linux"]
  }

  overridden_signature {
    signature_id = "2026323"
    action       = "REJECT"
    enabled      = true
  }

  overridden_signature {
    signature_id = "2026325"
    action       = "REJECT"
    enabled      = true
  }

  tag {
    scope = "color"
    tag   = "orange"
  }
}`, name)
}

func testAccNsxtPolicyIntrusionServiceProfileMinimalistic(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "%s"
  severities   = ["HIGH", "CRITICAL"]

  criteria {
    products_affected = ["Linux"]
  }

}`, name)
}
