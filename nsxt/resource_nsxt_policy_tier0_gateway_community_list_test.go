/* Copyright © 2020 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

        "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
        "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var accTestPolicyTier0GatewayCommunityListCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
        
}

var accTestPolicyTier0GatewayCommunityListUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
        
}

func TestAccResourceNsxtPolicyTier0GatewayCommunityList_basic(t *testing.T) {
	testResourceName := "nsxt_policy_tier0_gateway_community_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0GatewayCommunityListCheckDestroy(state, accTestPolicyTier0GatewayCommunityListUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0GatewayCommunityListTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0GatewayCommunityListExists(accTestPolicyTier0GatewayCommunityListCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyTier0GatewayCommunityListCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTier0GatewayCommunityListCreateAttributes["description"]),
                                        
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTier0GatewayCommunityListTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0GatewayCommunityListExists(accTestPolicyTier0GatewayCommunityListUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyTier0GatewayCommunityListUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTier0GatewayCommunityListUpdateAttributes["description"]),
                                        
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTier0GatewayCommunityListMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0GatewayCommunityListExists(accTestPolicyTier0GatewayCommunityListCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyTier0GatewayCommunityList_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_tier0_gateway_community_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier0GatewayCommunityListCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier0GatewayCommunityListMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyTier0GatewayCommunityListExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Tier0GatewayCommunityList resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy Tier0GatewayCommunityList resource ID not set in resources")
		}

                exists, err := resourceNsxtPolicyTier0GatewayCommunityListExists(resourceID, connector, testAccIsGlobalManager())
                if err != nil {
                    return err
                }
		if !exists {
                    return fmt.Errorf("Policy Tier0GatewayCommunityList %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyTier0GatewayCommunityListCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_tier0_gateway_community_list" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
                exists, err := resourceNsxtPolicyTier0GatewayCommunityListExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

                if exists {
                        return fmt.Errorf("Policy Tier0GatewayCommunityList %s still exists", displayName)
                }
	}
	return nil
}

func testAccNsxtPolicyTier0GatewayCommunityListTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyTier0GatewayCommunityListCreateAttributes
	} else {
		attrMap = accTestPolicyTier0GatewayCommunityListUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_tier0_gateway_community_list" "test" {
  display_name = "%s"
  description  = "%s"
  

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway_community_list.test.path
}`, attrMap["display_name"], attrMap["description"])
}

func testAccNsxtPolicyTier0GatewayCommunityListMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_tier0_gateway_community_list" "test" {
  display_name = "%s"

}

data "nsxt_policy_realization_info" "realization_info" {
  path = nsxt_policy_tier0_gateway_community_list.test.path
}`, accTestPolicyTier0GatewayCommunityListUpdateAttributes["display_name"])
}
