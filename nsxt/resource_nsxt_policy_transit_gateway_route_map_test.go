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

func TestAccResourceNsxtPolicyTransitGatewayRouteMap_basic(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway_route_map.test"
	prereqName := getAccTestResourceName()
	displayName := getAccTestResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayRouteMapCheckDestroy(state, displayName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayRouteMapCreateTemplate(prereqName, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayRouteMapExists(displayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "description", "terraform created"),
					resource.TestCheckResourceAttr(testResourceName, "entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "entry.0.action", "PERMIT"),
					resource.TestCheckResourceAttr(testResourceName, "entry.0.community_list_match.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "parent_path"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayRouteMapUpdateTemplate(prereqName, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayRouteMapExists(displayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(testResourceName, "description", "terraform updated"),
					resource.TestCheckResourceAttr(testResourceName, "entry.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "entry.0.action", "DENY"),
					resource.TestCheckResourceAttr(testResourceName, "entry.0.community_list_match.#", "2"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttrSet(testResourceName, "parent_path"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayRouteMapMinimalistic(prereqName, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayRouteMapExists(displayName, testResourceName),
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

func TestAccResourceNsxtPolicyTransitGatewayRouteMap_importBasic(t *testing.T) {
	prereqName := getAccTestResourceName()
	displayName := getAccTestResourceName()
	testResourceName := "nsxt_policy_transit_gateway_route_map.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.2.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayRouteMapCheckDestroy(state, displayName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayRouteMapMinimalistic(prereqName, displayName),
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

func testAccNsxtPolicyTransitGatewayRouteMapExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy TransitGatewayRouteMap resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy TransitGatewayRouteMap resource ID not set in resources")
		}

		parentPath := rs.Primary.Attributes["parent_path"]
		sessionContext := getSessionContextFromParentPath(testAccProvider.Meta(), parentPath)
		exists, err := resourceNsxtPolicyTransitGatewayRouteMapExists(sessionContext, parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy TransitGatewayRouteMap %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyTransitGatewayRouteMapCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_transit_gateway_route_map" {
			continue
		}
		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		sessionContext := getSessionContextFromParentPath(testAccProvider.Meta(), parentPath)
		exists, err := resourceNsxtPolicyTransitGatewayRouteMapExists(sessionContext, parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("Policy TransitGatewayRouteMap %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyTransitGatewayRouteMapPrerequisites(prereqName string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
}

resource "nsxt_policy_project" "test" {
  display_name = "%s"
  site_info {
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.test.path]
  }
}

resource "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  display_name    = "%s"
  transit_subnets = ["100.64.0.0/16"]
  centralized_config {
    ha_mode            = "ACTIVE_STANDBY"
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.test.path]
  }
}`, getEdgeClusterName(), prereqName, prereqName)
}

func testAccNsxtPolicyTransitGatewayRouteMapCreateTemplate(prereqName string, displayName string) string {
	return testAccNsxtPolicyTransitGatewayRouteMapPrerequisites(prereqName) + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_route_map" "test" {
  parent_path  = nsxt_policy_transit_gateway.test.path
  display_name = "%s"
  description  = "terraform created"

  entry {
    action = "PERMIT"
    community_list_match {
      criteria       = "11:22"
      match_operator = "MATCH_COMMUNITY_REGEX"
    }
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, displayName)
}

func testAccNsxtPolicyTransitGatewayRouteMapUpdateTemplate(prereqName string, displayName string) string {
	return testAccNsxtPolicyTransitGatewayRouteMapPrerequisites(prereqName) + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_route_map" "test" {
  parent_path  = nsxt_policy_transit_gateway.test.path
  display_name = "%s"
  description  = "terraform updated"

  entry {
    action = "DENY"
    community_list_match {
      criteria       = "11:22"
      match_operator = "MATCH_COMMUNITY_REGEX"
    }
    community_list_match {
      criteria       = "33:44"
      match_operator = "MATCH_COMMUNITY_REGEX"
    }
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, displayName)
}

func testAccNsxtPolicyTransitGatewayRouteMapMinimalistic(prereqName string, displayName string) string {
	return testAccNsxtPolicyTransitGatewayRouteMapPrerequisites(prereqName) + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_route_map" "test" {
  parent_path  = nsxt_policy_transit_gateway.test.path
  display_name = "%s"

  entry {
    action = "PERMIT"
    community_list_match {
      criteria       = "11:22"
      match_operator = "MATCH_COMMUNITY_REGEX"
    }
  }
}`, displayName)
}
