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

var EpRelatedResourceName = getAccTestResourceName()
var accTestPolicyTGWIPSecVpnLocalEndpointCreateAttributes = map[string]string{
	"display_name":  getAccTestResourceName(),
	"description":   "terraform created",
	"local_address": "20.20.0.10",
	"local_id":      "test-create",
}

var accTestPolicyTGWIPSecVpnLocalEndpointUpdateAttributes = map[string]string{
	"display_name":  getAccTestResourceName(),
	"description":   "terraform updated",
	"local_address": "20.20.0.20",
	"local_id":      "test-update",
}

func TestAccResourceNsxtPolicyTGWIPSecVpnLocalEndpoint_basic(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway_ipsec_vpn_local_endpoint.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTGWIPSecVpnLocalEndpointCheckDestroy(state, accTestPolicyTGWIPSecVpnLocalEndpointUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTGWIPSecVpnLocalEndpointTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTGWIPSecVpnLocalEndpointExists(accTestPolicyTGWIPSecVpnLocalEndpointCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyTGWIPSecVpnLocalEndpointCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTGWIPSecVpnLocalEndpointCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "local_address", accTestPolicyTGWIPSecVpnLocalEndpointCreateAttributes["local_address"]),
					resource.TestCheckResourceAttr(testResourceName, "local_id", accTestPolicyTGWIPSecVpnLocalEndpointCreateAttributes["local_id"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTGWIPSecVpnLocalEndpointTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTGWIPSecVpnLocalEndpointExists(accTestPolicyTGWIPSecVpnLocalEndpointUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyTGWIPSecVpnLocalEndpointUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTGWIPSecVpnLocalEndpointUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "local_address", accTestPolicyTGWIPSecVpnLocalEndpointUpdateAttributes["local_address"]),
					resource.TestCheckResourceAttr(testResourceName, "local_id", accTestPolicyTGWIPSecVpnLocalEndpointUpdateAttributes["local_id"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTGWIPSecVpnLocalEndpointMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTGWIPSecVpnLocalEndpointExists(accTestPolicyTGWIPSecVpnLocalEndpointCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyTGWIPSecVpnLocalEndpoint_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_transit_gateway_ipsec_vpn_local_endpoint.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTGWIPSecVpnLocalEndpointCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTGWIPSecVpnLocalEndpointMinimalistic(),
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

func testAccNsxtPolicyTGWIPSecVpnLocalEndpointExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy TransitGatewayIPSecVpnLocalEndpoint resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy TransitGatewayIPSecVpnLocalEndpoint resource ID not set in resources")
		}
		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := resourceNsxtPolicyTGWIPSecVpnLocalEndpointExists(testAccGetSessionContext(), parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy TransitGatewayIPSecVpnLocalEndpoint %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyTGWIPSecVpnLocalEndpointCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_transit_gateway_ipsec_vpn_local_endpoint" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := resourceNsxtPolicyTGWIPSecVpnLocalEndpointExists(testAccGetSessionContext(), parentPath, resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy TransitGatewayIPSecVpnLocalEndpoint %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyTGWIPSecVpnLocalEndpointTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyTGWIPSecVpnLocalEndpointCreateAttributes
	} else {
		attrMap = accTestPolicyTGWIPSecVpnLocalEndpointUpdateAttributes
	}
	return testAccNsxtPolicyTGWIPSecVpnLocalEndpointPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_ipsec_vpn_local_endpoint" "test" {
  parent_path   =  nsxt_policy_transit_gateway_ipsec_vpn_service.test.path
  display_name  = "%s"
  description   = "%s"
  local_address = "%s"
  local_id      = "%s"
    tag {
    scope = "scope1"
    tag   = "tag1"
  }
  depends_on = [nsxt_policy_transit_gateway_ipsec_vpn_service.test]
}`, attrMap["display_name"], attrMap["description"], attrMap["local_address"], attrMap["local_id"])
}

func testAccNsxtPolicyTGWIPSecVpnLocalEndpointMinimalistic() string {
	return testAccNsxtPolicyTGWIPSecVpnLocalEndpointPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_ipsec_vpn_local_endpoint" "test" {
  display_name = "%s"
  parent_path   =  nsxt_policy_transit_gateway_ipsec_vpn_service.test.path
  local_address = "fd21:3:3::1"
}`, accTestPolicyTGWIPSecVpnLocalEndpointUpdateAttributes["display_name"])
}

func testAccNsxtPolicyTGWIPSecVpnLocalEndpointPrerequisites() string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "test" {
  display_name = "%s"
}

resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "%s"
  edge_cluster_path = data.nsxt_policy_edge_cluster.test.path
  failover_mode            = "PREEMPTIVE"
  default_rule_logging     = false
  enable_firewall          = true
  ha_mode                  = "ACTIVE_STANDBY"
  internal_transit_subnets = ["102.64.0.0/16"]
  transit_subnets          = ["101.64.0.0/16"]
  vrf_transit_subnets      = ["100.64.0.0/16"]
  rd_admin_address         = "192.168.0.2"

  bgp_config {
    local_as_num    = "60000"
    multipath_relax = false

    route_aggregation {
      prefix = "12.10.10.0/24"
    }

    route_aggregation {
      prefix = "12.11.10.0/24"
    }
  }

  tag {
    scope = "color"
    tag   = "blue"
  }
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
  id = "default"
}
resource "nsxt_policy_transit_gateway_attachment" "test" {
  parent_path     = data.nsxt_policy_transit_gateway.test.path
  connection_path = nsxt_policy_gateway_connection.test.path
  display_name    = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
  depends_on = [data.nsxt_policy_transit_gateway.test]
}

resource "nsxt_policy_transit_gateway_ipsec_vpn_service" "test" {
  display_name = "%s"
  parent_path  = data.nsxt_policy_transit_gateway.test.path
  description  = "test resource"
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
  depends_on = [data.nsxt_policy_transit_gateway.test]
}`, getEdgeClusterName(), EpRelatedResourceName, EpRelatedResourceName, EpRelatedResourceName, EpRelatedResourceName, EpRelatedResourceName)
}
