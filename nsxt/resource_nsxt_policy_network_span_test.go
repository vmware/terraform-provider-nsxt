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

var accTestPolicyNetworkSpanCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
	"exclusive":    "true",
}

var accTestPolicyNetworkSpanUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
	"exclusive":    "false",
}

func TestAccResourceNsxtPolicyNetworkSpan_basic(t *testing.T) {
	testResourceName := "nsxt_policy_network_span.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyNetworkSpanCheckDestroy(state, accTestPolicyNetworkSpanUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyNetworkSpanTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyNetworkSpanExists(accTestPolicyNetworkSpanCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyNetworkSpanCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyNetworkSpanCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "exclusive", accTestPolicyNetworkSpanCreateAttributes["exclusive"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyNetworkSpanTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyNetworkSpanExists(accTestPolicyNetworkSpanUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyNetworkSpanUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyNetworkSpanUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "exclusive", accTestPolicyNetworkSpanUpdateAttributes["exclusive"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyNetworkSpanMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyNetworkSpanExists(accTestPolicyNetworkSpanCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyNetworkSpan_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_network_span.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyNetworkSpanCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyNetworkSpanMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyNetworkSpanExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy NetworkSpan resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy NetworkSpan resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyNetworkSpanExists(resourceID, connector, testAccIsGlobalManager())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy NetworkSpan %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyNetworkSpanCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_network_span" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyNetworkSpanExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy NetworkSpan %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyNetworkSpanTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyNetworkSpanCreateAttributes
	} else {
		attrMap = accTestPolicyNetworkSpanUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_network_span" "test" {
  display_name = "%s"
  description  = "%s"
  exclusive    = %s
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["exclusive"])
}

func testAccNsxtPolicyNetworkSpanMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_network_span" "test" {
  display_name = "%s"

}`, accTestPolicyNetworkSpanUpdateAttributes["display_name"])
}
