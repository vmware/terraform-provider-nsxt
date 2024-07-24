/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var accTestTransitGatewayCreateAttributes = map[string]string{
	"display_name":    getAccTestResourceName(),
	"description":     "terraform created",
	"transit_subnets": "192.168.5.0/24",
}

var accTestTransitGatewayUpdateAttributes = map[string]string{
	"display_name":    getAccTestResourceName(),
	"description":     "terraform updated",
	"transit_subnets": "192.168.7.0/24",
}

func TestAccResourceNsxtTransitGateway_basic(t *testing.T) {
	testResourceName := "nsxt_transit_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtTransitGatewayCheckDestroy(state, accTestTransitGatewayUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtTransitGatewayTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtTransitGatewayExists(accTestTransitGatewayCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestTransitGatewayCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestTransitGatewayCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "transit_subnets.0", accTestTransitGatewayCreateAttributes["transit_subnets"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtTransitGatewayTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtTransitGatewayExists(accTestTransitGatewayUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestTransitGatewayUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestTransitGatewayUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "transit_subnets.0", accTestTransitGatewayUpdateAttributes["transit_subnets"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtTransitGatewayMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtTransitGatewayExists(accTestTransitGatewayCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtTransitGateway_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_transit_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtTransitGatewayCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtTransitGatewayMinimalistic(),
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

func testAccNsxtTransitGatewayExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy TransitGateway resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy TransitGateway resource ID not set in resources")
		}

		exists, err := resourceNsxtTransitGatewayExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy TransitGateway %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtTransitGatewayCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_transit_gateway" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtTransitGatewayExists(testAccGetSessionContext(), resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy TransitGateway %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtTransitGatewayTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestTransitGatewayCreateAttributes
	} else {
		attrMap = accTestTransitGatewayUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_transit_gateway" "test" {
%s
  display_name    = "%s"
  description     = "%s"
  transit_subnets = ["%s"]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtProjectContext(), attrMap["display_name"], attrMap["description"], attrMap["transit_subnets"])
}

func testAccNsxtTransitGatewayMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_transit_gateway" "test" {
%s
  display_name    = "%s"
  transit_subnets = ["%s"]
}`, testAccNsxtProjectContext(), accTestTransitGatewayUpdateAttributes["display_name"], accTestTransitGatewayUpdateAttributes["transit_subnets"])
}
