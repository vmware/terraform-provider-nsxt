/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var accTestGatewayConnectionCreateAttributes = map[string]string{
	"display_name":     getAccTestResourceName(),
	"description":      "terraform created",
	"aggregate_routes": "192.168.240.0/24",
}

var accTestGatewayConnectionUpdateAttributes = map[string]string{
	"display_name":     getAccTestResourceName(),
	"description":      "terraform updated",
	"aggregate_routes": "192.168.241.0/24",
}

func TestAccResourceNsxtGatewayConnection_basic(t *testing.T) {
	testResourceName := "nsxt_gateway_connection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtGatewayConnectionCheckDestroy(state, accTestGatewayConnectionUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtGatewayConnectionTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtGatewayConnectionExists(accTestGatewayConnectionCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestGatewayConnectionCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestGatewayConnectionCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "aggregate_routes.0", accTestGatewayConnectionCreateAttributes["aggregate_routes"]),
					resource.TestCheckResourceAttrSet(testResourceName, "tier0_path"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtGatewayConnectionTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtGatewayConnectionExists(accTestGatewayConnectionUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestGatewayConnectionUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestGatewayConnectionUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "aggregate_routes.0", accTestGatewayConnectionUpdateAttributes["aggregate_routes"]),
					resource.TestCheckResourceAttrSet(testResourceName, "tier0_path"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtGatewayConnectionMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtGatewayConnectionExists(accTestGatewayConnectionCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtGatewayConnection_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_gateway_connection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtGatewayConnectionCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtGatewayConnectionMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtGatewayConnectionExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy GatewayConnection resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy GatewayConnection resource ID not set in resources")
		}

		exists, err := resourceNsxtGatewayConnectionExists(resourceID, connector, testAccIsGlobalManager())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy GatewayConnection %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtGatewayConnectionCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_gateway_connection" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtGatewayConnectionExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy GatewayConnection %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtGatewayConnectionTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestGatewayConnectionCreateAttributes
	} else {
		attrMap = accTestGatewayConnectionUpdateAttributes
	}
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "terraformt0gw"
  edge_cluster_path = data.nsxt_policy_edge_cluster.EC.path
}

resource "nsxt_gateway_connection" "test" {
  display_name = "%s"
  description  = "%s"
  tier0_path = nsxt_policy_tier0_gateway.test.path
  aggregate_routes = ["%s"]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, getEdgeClusterName(), attrMap["display_name"], attrMap["description"], attrMap["aggregate_routes"])
}

func testAccNsxtGatewayConnectionMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_gateway_connection" "test" {
  display_name = "%s"

}`, accTestGatewayConnectionUpdateAttributes["display_name"])
}
