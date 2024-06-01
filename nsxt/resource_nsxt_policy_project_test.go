/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/vmware/terraform-provider-nsxt/nsxt/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var shortID = getAccTestRandomString(6)

var accTestPolicyProjectCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
	"short_id":     shortID,
}

var accTestPolicyProjectUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
	"short_id":     shortID,
}

func getExpectedSiteInfoCount(t *testing.T) string {
	if util.NsxVersion == "" {
		connector, err := testAccGetPolicyConnector()
		if err != nil {
			t.Errorf("Failed to get policy connector")
			return "0"
		}

		err = initNSXVersion(connector)
		if err != nil {
			t.Errorf("Failed to retrieve NSX version")
			return "0"
		}
	}
	if util.NsxVersionHigherOrEqual("4.1.1") {
		return "1"
	}
	return "0"
}

func TestAccResourceNsxtPolicyProject_basic(t *testing.T) {
	testResourceName := "nsxt_policy_project.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectCheckDestroy(state, accTestPolicyProjectUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyProjectTemplate(true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectExists(accTestPolicyProjectCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyProjectCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyProjectCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "short_id", accTestPolicyProjectCreateAttributes["short_id"]),
					//TODO: add site info validation
					resource.TestCheckResourceAttr(testResourceName, "site_info.#", getExpectedSiteInfoCount(t)),
					resource.TestCheckResourceAttr(testResourceName, "tier0_gateway_paths.#", "1"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyProjectTemplate(false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectExists(accTestPolicyProjectUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyProjectUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyProjectUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "short_id", accTestPolicyProjectCreateAttributes["short_id"]),
					resource.TestCheckResourceAttr(testResourceName, "site_info.#", getExpectedSiteInfoCount(t)),
					resource.TestCheckResourceAttr(testResourceName, "tier0_gateway_paths.#", "1"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyProjectTemplate(false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectExists(accTestPolicyProjectUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyProjectUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyProjectUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "short_id", accTestPolicyProjectCreateAttributes["short_id"]),
					resource.TestCheckResourceAttr(testResourceName, "site_info.#", getExpectedSiteInfoCount(t)),
					resource.TestCheckResourceAttr(testResourceName, "tier0_gateway_paths.#", "0"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyProjectMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyProjectExists(accTestPolicyProjectCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyProject_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_project.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyProjectCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyProjectMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyProjectExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Project resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy Project resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyProjectExists(resourceID, connector, testAccIsGlobalManager())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy Project %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyProjectCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_project" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyProjectExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy Project %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyProjectTemplate(createFlow, includeT0GW bool) string {
	var attrMap map[string]string
	shortIDSpec := ""
	t0GW := ""
	if createFlow {
		attrMap = accTestPolicyProjectCreateAttributes
		shortIDSpec = fmt.Sprintf("short_id     = \"%s\"", attrMap["short_id"])
	} else {
		attrMap = accTestPolicyProjectUpdateAttributes
	}
	if includeT0GW {
		t0GW = "tier0_gateway_paths = [data.nsxt_policy_tier0_gateway.test.path]"
	}
	return fmt.Sprintf(`
data "nsxt_policy_tier0_gateway" "test" {
  display_name = "%s"
}
resource "nsxt_policy_project" "test" {
  display_name = "%s"
  description  = "%s"
  %s
  %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, getTier0RouterName(), attrMap["display_name"], attrMap["description"], shortIDSpec, t0GW)
}

func testAccNsxtPolicyProjectMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_project" "test" {
  display_name = "%s"

}`, accTestPolicyProjectUpdateAttributes["display_name"])
}
