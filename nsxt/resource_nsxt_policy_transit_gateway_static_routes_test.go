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

var accTestPolicyTransitGatewayStaticRoutesCreateAttributes = map[string]string{
	"display_name":         getAccTestResourceName(),
	"description":          "terraform created",
	"enabled_on_secondary": "true",
	"network":              "2.2.2.0/24",
	"ip_address":           "3.1.1.1",
	"admin_distance":       "2",
}

var accTestPolicyTransitGatewayStaticRoutesUpdateAttributes = map[string]string{
	"display_name":         getAccTestResourceName(),
	"description":          "terraform updated",
	"enabled_on_secondary": "false",
	"network":              "3.3.3.0/24",
	"ip_address":           "4.1.1.1",
	"admin_distance":       "5",
}

func TestAccResourceNsxtPolicyTransitGatewayStaticRoutes_basic(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway_static_route.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.1.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayStaticRoutesCheckDestroy(state, accTestPolicyTransitGatewayStaticRoutesUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayStaticRoutesTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayStaticRoutesExists(accTestPolicyTransitGatewayStaticRoutesCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyTransitGatewayStaticRoutesCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTransitGatewayStaticRoutesCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "network", accTestPolicyTransitGatewayStaticRoutesCreateAttributes["network"]),
					resource.TestCheckResourceAttr(testResourceName, "next_hop.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "enabled_on_secondary", accTestPolicyTransitGatewayStaticRoutesCreateAttributes["enabled_on_secondary"]),
					resource.TestCheckResourceAttr(testResourceName, "next_hop.0.ip_address", accTestPolicyTransitGatewayStaticRoutesCreateAttributes["ip_address"]),
					resource.TestCheckResourceAttr(testResourceName, "next_hop.0.admin_distance", accTestPolicyTransitGatewayStaticRoutesCreateAttributes["admin_distance"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayStaticRoutesTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayStaticRoutesExists(accTestPolicyTransitGatewayStaticRoutesUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyTransitGatewayStaticRoutesUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTransitGatewayStaticRoutesUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "network", accTestPolicyTransitGatewayStaticRoutesUpdateAttributes["network"]),
					resource.TestCheckResourceAttr(testResourceName, "next_hop.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "enabled_on_secondary", accTestPolicyTransitGatewayStaticRoutesUpdateAttributes["enabled_on_secondary"]),
					resource.TestCheckResourceAttr(testResourceName, "next_hop.0.ip_address", accTestPolicyTransitGatewayStaticRoutesUpdateAttributes["ip_address"]),
					resource.TestCheckResourceAttr(testResourceName, "next_hop.0.admin_distance", accTestPolicyTransitGatewayStaticRoutesUpdateAttributes["admin_distance"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTransitGatewayStaticRoutesMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTransitGatewayStaticRoutesExists(accTestPolicyTransitGatewayStaticRoutesCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyTransitGatewayStaticRoutes_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_transit_gateway_static_route.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.1.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTransitGatewayStaticRoutesCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTransitGatewayStaticRoutesMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyTransitGatewayStaticRoutesExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy TransitGatewayStaticRoutes resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy ransitGatewayStaticRoutes resource ID not set in resources")
		}
		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := resourceNsxtPolicyTransitGatewayStaticRoutesExists(testAccGetSessionContext(), parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy ransitGatewayStaticRoutes %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyTransitGatewayStaticRoutesCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_transit_gateway_static_route" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := resourceNsxtPolicyTransitGatewayStaticRoutesExists(testAccGetSessionContext(), parentPath, resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy ransitGatewayStaticRoutes %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyTransitGatewayStaticRoutesTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyTransitGatewayStaticRoutesCreateAttributes
	} else {
		attrMap = accTestPolicyTransitGatewayStaticRoutesUpdateAttributes
	}
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "testEdgeCluster"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "testt0gw"
  edge_cluster_path = data.nsxt_policy_edge_cluster.test.path
}
resource "nsxt_policy_project" "proj" {
  display_name        = "demoproj1"
  description         = "Terraform provisioned Project"
  short_id            = "test"
  tier0_gateway_paths = [nsxt_policy_tier0_gateway.test.path]
}

resource "nsxt_policy_transit_gateway" "tgwStaticRoute" {
  context {
    project_id = resource.nsxt_policy_project.proj.id
  }

  display_name    = "tgwStaticRoute"
  description     = "Terraform provisioned TransitGateway"
  transit_subnets = ["10.203.4.0/24"]
}

resource "nsxt_policy_transit_gateway_static_route" "test" {
  display_name = "%s"
  description  = "%s"
  parent_path = "resource.nsxt_policy_transit_gateway.tgwStaticRoute.path"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"])
}

func testAccNsxtPolicyTransitGatewayStaticRoutesMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_static_route" "test" {
  display_name = "%s"

}`, accTestPolicyTransitGatewayStaticRoutesUpdateAttributes["display_name"])
}
