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

var RelatedResourceName = getAccTestResourceName()

var accTestPolicyTGWIPSecVpnServicesCreateAttributes = map[string]string{
	"display_name":  getAccTestResourceName(),
	"description":   "terraform created",
	"enabled":       "true",
	"ha_sync":       "true",
	"ike_log_level": "INFO",
	"sources":       "192.168.10.0/24",
	"destinations":  "192.169.10.0/24",
	"action":        "BYPASS",
}

var accTestPolicyTGWIPSecVpnServicesUpdateAttributes = map[string]string{

	"description":   "terraform updated",
	"enabled":       "false",
	"ha_sync":       "false",
	"ike_log_level": "INFO",
	"sources":       "192.170.10.0/24",
	"destinations":  "192.171.10.0/24",
	"action":        "BYPASS",
}

func TestAccResourceNsxtPolicyTGWIPSecVpnServices_basic(t *testing.T) {
	testResourceName := "nsxt_policy_transit_gateway_ipsec_vpn_service.test"
	createDisplayName := getAccTestResourceName()
	updateDisplayName := getAccTestResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTGWIPSecVpnServicesCheckDestroy(state, updateDisplayName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTGWIPSecVpnServicesTemplate(true, createDisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTGWIPSecVpnServicesExists(createDisplayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createDisplayName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTGWIPSecVpnServicesCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyTGWIPSecVpnServicesCreateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_sync", accTestPolicyTGWIPSecVpnServicesCreateAttributes["ha_sync"]),
					resource.TestCheckResourceAttr(testResourceName, "ike_log_level", accTestPolicyTGWIPSecVpnServicesCreateAttributes["ike_log_level"]),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.0.sources.0", accTestPolicyTGWIPSecVpnServicesCreateAttributes["sources"]),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.0.destinations.0", accTestPolicyTGWIPSecVpnServicesCreateAttributes["destinations"]),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.0.action", accTestPolicyTGWIPSecVpnServicesCreateAttributes["action"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTGWIPSecVpnServicesTemplate(false, updateDisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTGWIPSecVpnServicesExists(updateDisplayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateDisplayName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTGWIPSecVpnServicesUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyTGWIPSecVpnServicesUpdateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_sync", accTestPolicyTGWIPSecVpnServicesUpdateAttributes["ha_sync"]),
					resource.TestCheckResourceAttr(testResourceName, "ike_log_level", accTestPolicyTGWIPSecVpnServicesUpdateAttributes["ike_log_level"]),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.0.sources.0", accTestPolicyTGWIPSecVpnServicesUpdateAttributes["sources"]),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.0.destinations.0", accTestPolicyTGWIPSecVpnServicesUpdateAttributes["destinations"]),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.0.action", accTestPolicyTGWIPSecVpnServicesUpdateAttributes["action"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTGWIPSecVpnServicesMinimalistic(updateDisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTGWIPSecVpnServicesExists(createDisplayName, testResourceName),
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

func TestAccResourceNsxtPolicyTGWIPSecVpnServices_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_transit_gateway_ipsec_vpn_service.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTGWIPSecVpnServicesCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTGWIPSecVpnServicesMinimalistic(name),
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

func testAccNsxtPolicyTGWIPSecVpnServicesExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy TransitGatewayIPSecVpnServices resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy TransitGatewayIPSecVpnServices resource ID not set in resources")
		}
		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := resourceNsxtPolicyTGWIPSecVpnServicesExists(testAccGetSessionContext(), parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy TransitGatewayIPSecVpnServices %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyTGWIPSecVpnServicesCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_transit_gateway_ipsec_vpn_service" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := resourceNsxtPolicyTGWIPSecVpnServicesExists(testAccGetSessionContext(), parentPath, resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy TransitGatewayIPSecVpnServices %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyTGWIPSecVpnServicesTemplate(createFlow bool, displayName string) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyTGWIPSecVpnServicesCreateAttributes
	} else {
		attrMap = accTestPolicyTGWIPSecVpnServicesUpdateAttributes
	}
	return testAccNsxtPolicyTGWIPSecVpnServicesPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_ipsec_vpn_service" "test" {
  display_name = "%s"
  parent_path  = data.nsxt_policy_transit_gateway.test.path
  description  = "%s"
  enabled       = "%s"
  ha_sync       = "%s"
  ike_log_level = "%s"
  bypass_rule {
    sources      = ["%s"]
    destinations = ["%s"]
    action       = "%s"
  }
  tag {
    scope = "scope1"
    tag   = "tag1"
  }
  depends_on = [data.nsxt_policy_transit_gateway.test]
}`, displayName, attrMap["description"], attrMap["enabled"], attrMap["ha_sync"], attrMap["ike_log_level"], attrMap["sources"], attrMap["destinations"], attrMap["action"])
}

func testAccNsxtPolicyTGWIPSecVpnServicesMinimalistic(displayName string) string {
	return testAccNsxtPolicyTGWIPSecVpnServicesPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_ipsec_vpn_service" "test" {
  display_name = "%s"
  parent_path      = data.nsxt_policy_transit_gateway.test.path

}`, displayName)
}

func testAccNsxtPolicyTGWIPSecVpnServicesPrerequisites() string {
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

resource "nsxt_policy_ip_block" "extblk" {
  display_name = "ipsec-testing-extblk"
  cidr         = "10.110.0.0/22"
  visibility   = "EXTERNAL"
}
  
resource "nsxt_policy_project" "test" {
  display_name             = "%s"
  tier0_gateway_paths      = [nsxt_policy_tier0_gateway.test.path]
  tgw_external_connections = [nsxt_policy_gateway_connection.test.path]
  site_info {
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.test.path]
  }
 external_ipv4_blocks = [nsxt_policy_ip_block.extblk.path] 
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
}`, getEdgeClusterName(), RelatedResourceName, RelatedResourceName, RelatedResourceName, RelatedResourceName)
}
