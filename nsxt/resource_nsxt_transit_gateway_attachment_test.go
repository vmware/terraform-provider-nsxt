/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var dependantEntityName = getAccTestResourceName()

var accTestTransitGatewayAttachmentCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
}

var accTestTransitGatewayAttachmentUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
}

func TestAccResourceNsxtTransitGatewayAttachment_basic(t *testing.T) {
	testResourceName := "nsxt_transit_gateway_attachment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtTransitGatewayAttachmentCheckDestroy(state, accTestTransitGatewayAttachmentUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtTransitGatewayAttachmentTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtTransitGatewayAttachmentExists(accTestTransitGatewayAttachmentCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestTransitGatewayAttachmentCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestTransitGatewayAttachmentCreateAttributes["description"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtTransitGatewayAttachmentTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtTransitGatewayAttachmentExists(accTestTransitGatewayAttachmentUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestTransitGatewayAttachmentUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestTransitGatewayAttachmentUpdateAttributes["description"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtTransitGatewayAttachment_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_transit_gateway_attachment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtTransitGatewayAttachmentCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtTransitGatewayAttachmentTemplate(true),
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

func testAccNsxtTransitGatewayAttachmentExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy TransitGatewayAttachment resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy TransitGatewayAttachment resource ID not set in resources")
		}

		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := resourceNsxtTransitGatewayAttachmentExists(testAccGetSessionContext(), parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy TransitGatewayAttachment %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtTransitGatewayAttachmentCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_transit_gateway_attachment" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := resourceNsxtTransitGatewayAttachmentExists(testAccGetSessionContext(), parentPath, resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy TransitGatewayAttachment %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtTransitGatewayAttachmentTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestTransitGatewayAttachmentCreateAttributes
	} else {
		attrMap = accTestTransitGatewayAttachmentUpdateAttributes
	}
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "%s"
  edge_cluster_path = data.nsxt_policy_edge_cluster.test.path
}

resource "nsxt_policy_gateway_connection" "test" {
  display_name      = "%s"
  tier0_path = nsxt_policy_tier0_gateway.test.path
  aggregate_routes = ["192.168.240.0/24"]
}

resource "nsxt_policy_project" "test" {
  display_name      = "%s"
  tier0_gateway_paths = [nsxt_policy_tier0_gateway.test.path]
  tgw_external_connections = [nsxt_gateway_connection.test.path]
}

resource "nsxt_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  display_name      = "%s"
  transit_subnets = ["192.168.7.0/24"]
}

resource "nsxt_transit_gateway_attachment" "test" {
  parent_path  = nsxt_transit_gateway.test.path
  connection_path = nsxt_gateway_connection.test.path
  display_name = "%s"
  description  = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, getEdgeClusterName(), dependantEntityName, dependantEntityName, dependantEntityName, dependantEntityName, attrMap["display_name"], attrMap["description"])
}
