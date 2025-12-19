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

var dependantResourceName = getAccTestResourceName()

var accTestPolicyTransitGatewayStaticRouteCreateAttributes = map[string]string{
	"display_name":         getAccTestResourceName(),
	"description":          "terraform created",
	"enabled_on_secondary": "true",
	"network":              "2.2.2.0/24",
	"admin_distance":       "2",
}

var accTestPolicyTransitGatewayStaticRouteUpdateAttributes = map[string]string{
	"display_name":         getAccTestResourceName(),
	"description":          "terraform updated",
	"enabled_on_secondary": "false",
	"network":              "3.3.3.0/24",
	"admin_distance":       "5",
}

func TestAccResourceNsxtPolicyTransitGatewayStaticRoute_basic(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway_static_route.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayStaticRouteCheckDestroy(state, accTestPolicyTransitGatewayStaticRouteUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayStaticRouteTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayStaticRouteExists(accTestPolicyTransitGatewayStaticRouteCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyTransitGatewayStaticRouteCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTransitGatewayStaticRouteCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "network", accTestPolicyTransitGatewayStaticRouteCreateAttributes["network"]),
					resource.TestCheckResourceAttr(testResourceName, "next_hop.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "enabled_on_secondary", accTestPolicyTransitGatewayStaticRouteCreateAttributes["enabled_on_secondary"]),
					resource.TestCheckResourceAttr(testResourceName, "next_hop.0.admin_distance", accTestPolicyTransitGatewayStaticRouteCreateAttributes["admin_distance"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayStaticRouteTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayStaticRouteExists(accTestPolicyTransitGatewayStaticRouteUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyTransitGatewayStaticRouteUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTransitGatewayStaticRouteUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "network", accTestPolicyTransitGatewayStaticRouteUpdateAttributes["network"]),
					resource.TestCheckResourceAttr(testResourceName, "next_hop.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "enabled_on_secondary", accTestPolicyTransitGatewayStaticRouteUpdateAttributes["enabled_on_secondary"]),
					resource.TestCheckResourceAttr(testResourceName, "next_hop.0.admin_distance", accTestPolicyTransitGatewayStaticRouteUpdateAttributes["admin_distance"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayStaticRouteMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayStaticRouteExists(accTestPolicyTransitGatewayStaticRouteCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyTransitGatewayStaticRoute_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_transit_gateway_static_route.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayStaticRouteCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayStaticRouteMinimalistic(),
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

func testAccNsxtPolicyTransitGatewayStaticRouteExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy TransitGatewayStaticRoute resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy ransitGatewayStaticRoute resource ID not set in resources")
		}
		parentPath := rs.Primary.Attributes["parent_path"]
		sessionContext := getSessionContextFromParentPath(testAccProvider.Meta(), parentPath)
		exists, err := resourceNsxtPolicyTransitGatewayStaticRouteExists(sessionContext, parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy transitGatewayStaticRoute %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyTransitGatewayStaticRouteCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_transit_gateway_static_route" {
			continue
		}
		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		sessionContext := getSessionContextFromParentPath(testAccProvider.Meta(), parentPath)
		exists, err := resourceNsxtPolicyTransitGatewayStaticRouteExists(sessionContext, parentPath, resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy ransitGatewayStaticRoute %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyTransitGatewayStaticRouteTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyTransitGatewayStaticRouteCreateAttributes
	} else {
		attrMap = accTestPolicyTransitGatewayStaticRouteUpdateAttributes
	}
	return testAccNsxtPolicyTransitGatewayStaticRoutePrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_static_route" "test" {
  display_name          = "%s"
  description           = "%s"
  parent_path           = data.nsxt_policy_transit_gateway.test.path
  enabled_on_secondary  = %s
  network               = "%s"

  next_hop {
    admin_distance = %s
    scope          = [nsxt_policy_transit_gateway_attachment.test.path]
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
  depends_on = [data.nsxt_policy_transit_gateway.test]
}`, attrMap["display_name"], attrMap["description"], attrMap["enabled_on_secondary"], attrMap["network"], attrMap["admin_distance"])
}

func testAccNsxtPolicyTransitGatewayStaticRouteMinimalistic() string {
	return testAccNsxtPolicyTransitGatewayStaticRoutePrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_static_route" "test" {
  display_name     = "%s"
  parent_path      = data.nsxt_policy_transit_gateway.test.path
  description      = ""
  network          = "%s"
  next_hop {
    admin_distance = %s
    scope          = [nsxt_policy_transit_gateway_attachment.test.path]
  }
  depends_on = [data.nsxt_policy_transit_gateway.test]
}`, accTestPolicyTransitGatewayStaticRouteUpdateAttributes["display_name"], accTestPolicyTransitGatewayStaticRouteUpdateAttributes["network"], accTestPolicyTransitGatewayStaticRouteUpdateAttributes["admin_distance"])
}

func testAccNsxtPolicyTransitGatewayStaticRoutePrerequisites() string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "%s"
  edge_cluster_path = data.nsxt_policy_edge_cluster.test.path
}

resource "nsxt_policy_gateway_connection" "test" {
  display_name     = "%s"
  tier0_path       = nsxt_policy_tier0_gateway.test.path
  aggregate_routes = ["192.168.240.0/24"]
}

resource "nsxt_policy_project" "test" {
  display_name             = "%s"
  tier0_gateway_paths      = [nsxt_policy_tier0_gateway.test.path]
  tgw_external_connections = [nsxt_policy_gateway_connection.test.path]
  site_info {
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.test.path]
  }
}

data "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  is_default = true
}

resource "nsxt_policy_transit_gateway_attachment" "test" {
  parent_path     = data.nsxt_policy_transit_gateway.test.path
  connection_path = nsxt_policy_gateway_connection.test.path
  display_name    = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, getEdgeClusterName(), dependantResourceName, dependantResourceName, dependantResourceName, dependantResourceName)
}
