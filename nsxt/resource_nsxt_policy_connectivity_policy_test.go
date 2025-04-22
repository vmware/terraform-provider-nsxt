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

var accTestPolicyConnectivityPolicyCreateAttributes = map[string]string{
	"display_name":       getAccTestResourceName(),
	"description":        "terraform created",
	"connectivity_scope": "ISOLATED",
}

var accTestPolicyConnectivityPolicyUpdateAttributes = map[string]string{
	"display_name":       getAccTestResourceName(),
	"description":        "terraform updated",
	"connectivity_scope": "ISOLATED", // Updating connectivity scope will cause re-creation
}

func TestAccResourceNsxtPolicyConnectivityPolicy_basic(t *testing.T) {
	testResourceName := "nsxt_policy_connectivity_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.1.0")
			testAccOnlyVPC(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyConnectivityPolicyCheckDestroy(state, accTestPolicyConnectivityPolicyUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyConnectivityPolicyTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyConnectivityPolicyExists(accTestPolicyConnectivityPolicyCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyConnectivityPolicyCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyConnectivityPolicyCreateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "group_path"),
					resource.TestCheckResourceAttr(testResourceName, "connectivity_scope", accTestPolicyConnectivityPolicyCreateAttributes["connectivity_scope"]),
					resource.TestCheckResourceAttrSet(testResourceName, "internal_id"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyConnectivityPolicyTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyConnectivityPolicyExists(accTestPolicyConnectivityPolicyUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyConnectivityPolicyUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyConnectivityPolicyUpdateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "group_path"),
					resource.TestCheckResourceAttr(testResourceName, "connectivity_scope", accTestPolicyConnectivityPolicyUpdateAttributes["connectivity_scope"]),
					resource.TestCheckResourceAttrSet(testResourceName, "internal_id"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyConnectivityPolicyMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyConnectivityPolicyExists(accTestPolicyConnectivityPolicyCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyConnectivityPolicy_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_connectivity_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.1.0")
			testAccOnlyVPC(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyConnectivityPolicyCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyConnectivityPolicyMinimalistic(),
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

func testAccNsxtPolicyConnectivityPolicyExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy ConnectivityPolicy resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy ConnectivityPolicy resource ID not set in resources")
		}

		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := resourceNsxtPolicyConnectivityPolicyExists(testAccGetSessionProjectContext(), parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy ConnectivityPolicy %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyConnectivityPolicyCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_connectivity_policy" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]

		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := resourceNsxtPolicyConnectivityPolicyExists(testAccGetSessionProjectContext(), parentPath, resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy ConnectivityPolicy %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyConnectivityPolicyPrerequisites() string {
	return fmt.Sprintf(`
data "nsxt_policy_transit_gateway" "default" {
  %s
  id = "default"
}

resource "nsxt_policy_group" "test" {
  %s
  display_name = "%s"
}`, testAccNsxtProjectContext(), testAccNsxtProjectContext(), accTestPolicyConnectivityPolicyCreateAttributes["display_name"])
}

func testAccNsxtPolicyConnectivityPolicyTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyConnectivityPolicyCreateAttributes
	} else {
		attrMap = accTestPolicyConnectivityPolicyUpdateAttributes
	}
	return testAccNsxtPolicyConnectivityPolicyPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_connectivity_policy" "test" {
  parent_path  = data.nsxt_policy_transit_gateway.default.path
  display_name = "%s"
  description  = "%s"

  group_path         = nsxt_policy_group.test.path
  connectivity_scope = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["connectivity_scope"])
}

func testAccNsxtPolicyConnectivityPolicyMinimalistic() string {
	return testAccNsxtPolicyConnectivityPolicyPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_connectivity_policy" "test" {
  parent_path  = data.nsxt_policy_transit_gateway.default.path
  group_path   = nsxt_policy_group.test.path
  display_name = "%s"
}`, accTestPolicyConnectivityPolicyUpdateAttributes["display_name"])
}
