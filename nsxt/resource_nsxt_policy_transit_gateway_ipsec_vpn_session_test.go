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

var SessionRelatedResourceName = getAccTestResourceName()

var accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes = map[string]string{
	"description":                "terraform created",
	"enabled":                    "true",
	"vpn_type":                   "RouteBased",
	"authentication_mode":        "PSK",
	"compliance_suite":           "NONE",
	"ip_addresses":               "169.254.152.25",
	"prefix_length":              "24",
	"peer_address":               "18.18.18.22",
	"peer_id":                    "18.18.18.22",
	"psk":                        "secret1",
	"connection_initiation_mode": "RESPOND_ONLY",
}

var accTestPolicyTGWIPSecVpnSessionRouteBasedUpdateAttributes = map[string]string{
	"description":                "terraform updated",
	"enabled":                    "false",
	"vpn_type":                   "RouteBased",
	"authentication_mode":        "PSK",
	"compliance_suite":           "NONE",
	"ip_addresses":               "169.254.152.26",
	"prefix_length":              "24",
	"peer_address":               "18.18.18.21",
	"peer_id":                    "18.18.18.21",
	"psk":                        "secret2",
	"connection_initiation_mode": "ON_DEMAND",
}

var accTestPolicyTGWIPSecVpnSessionPolicyBasedCreateAttributes = map[string]string{
	"description":                "Terraform-provisioned IPsec Route-Based VPN",
	"enabled":                    "true",
	"vpn_type":                   "PolicyBased",
	"authentication_mode":        "PSK",
	"compliance_suite":           "NONE",
	"peer_address":               "18.18.18.21",
	"peer_id":                    "18.18.18.21",
	"psk":                        "VMware123!",
	"connection_initiation_mode": "RESPOND_ONLY",
	"sources":                    "192.170.10.0/24",
	"destinations":               "192.171.10.0/24",
	"action":                     "PROTECT",
}

var accTestPolicyTGWIPSecVpnSessionPolicyBasedUpdateAttributes = map[string]string{
	"description":                "terraform updated",
	"enabled":                    "true",
	"vpn_type":                   "PolicyBased",
	"authentication_mode":        "PSK",
	"compliance_suite":           "NONE",
	"peer_address":               "18.18.18.22",
	"peer_id":                    "18.18.18.22",
	"psk":                        "secret1",
	"connection_initiation_mode": "RESPOND_ONLY",
	"sources":                    "192.172.10.0/24",
	"destinations":               "192.173.10.0/24",
	"action":                     "PROTECT",
}

func TestAccResourceNsxtPolicyTGWIPSecVpnSessionRouteBased_basic(t *testing.T) {
	createDisplayName := getAccTestResourceName()
	updateDisplayName := getAccTestResourceName()
	testResourceName := "nsxt_policy_transit_gateway_ipsec_vpn_session.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTGWIPSecVpnSessionCheckDestroy(state, createDisplayName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTGWIPSecVpnSessionRouteBasedTemplate(true, createDisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTGWIPSecVpnSessionExists(createDisplayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createDisplayName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "ike_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "tunnel_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "dpd_profile_path"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "vpn_type", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["vpn_type"]),
					resource.TestCheckResourceAttr(testResourceName, "authentication_mode", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["authentication_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "compliance_suite", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["compliance_suite"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "prefix_length", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["prefix_length"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_address", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["peer_address"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_id", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["peer_id"]),
					resource.TestCheckResourceAttr(testResourceName, "psk", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["psk"]),
					resource.TestCheckResourceAttr(testResourceName, "connection_initiation_mode", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["connection_initiation_mode"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTGWIPSecVpnSessionRouteBasedTemplate(false, updateDisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTGWIPSecVpnSessionExists(updateDisplayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateDisplayName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTGWIPSecVpnSessionRouteBasedUpdateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "ike_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "tunnel_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "dpd_profile_path"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyTGWIPSecVpnSessionRouteBasedUpdateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "vpn_type", accTestPolicyTGWIPSecVpnSessionRouteBasedUpdateAttributes["vpn_type"]),
					resource.TestCheckResourceAttr(testResourceName, "authentication_mode", accTestPolicyTGWIPSecVpnSessionRouteBasedUpdateAttributes["authentication_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "compliance_suite", accTestPolicyTGWIPSecVpnSessionRouteBasedUpdateAttributes["compliance_suite"]),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestPolicyTGWIPSecVpnSessionRouteBasedUpdateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "prefix_length", accTestPolicyTGWIPSecVpnSessionRouteBasedUpdateAttributes["prefix_length"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_address", accTestPolicyTGWIPSecVpnSessionRouteBasedUpdateAttributes["peer_address"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_id", accTestPolicyTGWIPSecVpnSessionRouteBasedUpdateAttributes["peer_id"]),
					resource.TestCheckResourceAttr(testResourceName, "psk", accTestPolicyTGWIPSecVpnSessionRouteBasedUpdateAttributes["psk"]),
					resource.TestCheckResourceAttr(testResourceName, "connection_initiation_mode", accTestPolicyTGWIPSecVpnSessionRouteBasedUpdateAttributes["connection_initiation_mode"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTGWIPSecVpnSessionRouteBasedMinimalistic(createDisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTGWIPSecVpnSessionExists(createDisplayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createDisplayName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "ike_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "tunnel_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "dpd_profile_path"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "vpn_type", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["vpn_type"]),
					resource.TestCheckResourceAttr(testResourceName, "authentication_mode", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["authentication_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "compliance_suite", "NONE"),
					resource.TestCheckResourceAttr(testResourceName, "ip_addresses.0", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["ip_addresses"]),
					resource.TestCheckResourceAttr(testResourceName, "prefix_length", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["prefix_length"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_address", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["peer_address"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_id", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["peer_id"]),
					resource.TestCheckResourceAttr(testResourceName, "psk", accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes["psk"]),
					resource.TestCheckResourceAttr(testResourceName, "connection_initiation_mode", "INITIATOR"),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTGWIPSecVpnSessionPolicyBased_basic(t *testing.T) {
	createDisplayName := getAccTestResourceName()
	updateDisplayName := getAccTestResourceName()
	testResourceName := "nsxt_policy_transit_gateway_ipsec_vpn_session.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTGWIPSecVpnSessionCheckDestroy(state, createDisplayName)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTGWIPSecVpnSessionPolicyBasedTemplate(true, createDisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTGWIPSecVpnSessionExists(createDisplayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createDisplayName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTGWIPSecVpnSessionPolicyBasedCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyTGWIPSecVpnSessionPolicyBasedCreateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "vpn_type", accTestPolicyTGWIPSecVpnSessionPolicyBasedCreateAttributes["vpn_type"]),
					resource.TestCheckResourceAttr(testResourceName, "authentication_mode", accTestPolicyTGWIPSecVpnSessionPolicyBasedCreateAttributes["authentication_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "compliance_suite", accTestPolicyTGWIPSecVpnSessionPolicyBasedCreateAttributes["compliance_suite"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_address", accTestPolicyTGWIPSecVpnSessionPolicyBasedCreateAttributes["peer_address"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_id", accTestPolicyTGWIPSecVpnSessionPolicyBasedCreateAttributes["peer_id"]),
					resource.TestCheckResourceAttr(testResourceName, "psk", accTestPolicyTGWIPSecVpnSessionPolicyBasedCreateAttributes["psk"]),
					resource.TestCheckResourceAttr(testResourceName, "connection_initiation_mode", accTestPolicyTGWIPSecVpnSessionPolicyBasedCreateAttributes["connection_initiation_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.sources.0", accTestPolicyTGWIPSecVpnSessionPolicyBasedCreateAttributes["sources"]),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destinations.0", accTestPolicyTGWIPSecVpnSessionPolicyBasedCreateAttributes["destinations"]),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", accTestPolicyTGWIPSecVpnSessionPolicyBasedCreateAttributes["action"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyTGWIPSecVpnSessionPolicyBasedTemplate(false, updateDisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTGWIPSecVpnSessionExists(updateDisplayName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateDisplayName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyTGWIPSecVpnSessionPolicyBasedUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyTGWIPSecVpnSessionPolicyBasedUpdateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "vpn_type", accTestPolicyTGWIPSecVpnSessionPolicyBasedUpdateAttributes["vpn_type"]),
					resource.TestCheckResourceAttr(testResourceName, "authentication_mode", accTestPolicyTGWIPSecVpnSessionPolicyBasedUpdateAttributes["authentication_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "compliance_suite", accTestPolicyTGWIPSecVpnSessionPolicyBasedUpdateAttributes["compliance_suite"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_address", accTestPolicyTGWIPSecVpnSessionPolicyBasedUpdateAttributes["peer_address"]),
					resource.TestCheckResourceAttr(testResourceName, "peer_id", accTestPolicyTGWIPSecVpnSessionPolicyBasedUpdateAttributes["peer_id"]),
					resource.TestCheckResourceAttr(testResourceName, "psk", accTestPolicyTGWIPSecVpnSessionPolicyBasedUpdateAttributes["psk"]),
					resource.TestCheckResourceAttr(testResourceName, "connection_initiation_mode", accTestPolicyTGWIPSecVpnSessionPolicyBasedUpdateAttributes["connection_initiation_mode"]),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.sources.0", accTestPolicyTGWIPSecVpnSessionPolicyBasedUpdateAttributes["sources"]),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destinations.0", accTestPolicyTGWIPSecVpnSessionPolicyBasedUpdateAttributes["destinations"]),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", accTestPolicyTGWIPSecVpnSessionPolicyBasedUpdateAttributes["action"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func testAccNsxtPolicyTGWIPSecVpnSessionPolicyBasedTemplate(createFlow bool, displayName string) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyTGWIPSecVpnSessionPolicyBasedCreateAttributes
	} else {
		attrMap = accTestPolicyTGWIPSecVpnSessionPolicyBasedUpdateAttributes
	}
	return testAccNsxtPolicyTGWIPSecVpnSessionPreConditionTemplate() +
		fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_ipsec_vpn_session" "test" {
  display_name               = "%s"
  description                = "%s"
  enabled                    = "%s"
  parent_path                = nsxt_policy_transit_gateway_ipsec_vpn_service.test.path
  local_endpoint_path        = nsxt_policy_transit_gateway_ipsec_vpn_local_endpoint.test.path
  vpn_type                   = "%s"
  authentication_mode        = "%s"
  compliance_suite           = "%s"
  peer_address               = "%s"
  peer_id                    = "%s"
  psk                        = "%s"
  connection_initiation_mode = "%s"
  direction                  = "BOTH"
  max_segment_size           = null

  rule {
    sources      = ["%s"]
    destinations = ["%s"]
    action       = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, displayName, attrMap["description"], attrMap["enabled"], attrMap["vpn_type"],
			attrMap["authentication_mode"], attrMap["compliance_suite"], attrMap["peer_address"], attrMap["peer_id"],
			attrMap["psk"], attrMap["connection_initiation_mode"], attrMap["sources"], attrMap["destinations"], attrMap["action"])
}

func testAccNsxtPolicyTGWIPSecVpnSessionRouteBasedTemplate(createFlow bool, displayName string) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes
	} else {
		attrMap = accTestPolicyTGWIPSecVpnSessionRouteBasedUpdateAttributes
	}
	return testAccNsxtPolicyTGWIPSecVpnSessionPreConditionTemplate() +
		fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_ipsec_vpn_session" "test" {
  display_name               = "%s"
  description                = "%s"
  enabled                    = "%s"
  parent_path                = nsxt_policy_transit_gateway_ipsec_vpn_service.test.path
  local_endpoint_path        = nsxt_policy_transit_gateway_ipsec_vpn_local_endpoint.test.path
  vpn_type                   = "%s"
  authentication_mode        = "%s"
  compliance_suite           = "%s"
  ip_addresses               = ["%s"]
  prefix_length              = "%s"
  peer_address               = "%s"
  peer_id                    = "%s"
  psk                        = "%s"
  connection_initiation_mode = "%s"
  direction                  = "BOTH"
  max_segment_size           = null

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
  depends_on = [nsxt_policy_transit_gateway_ipsec_vpn_service.test, nsxt_policy_transit_gateway_ipsec_vpn_local_endpoint.test]
}`, displayName, attrMap["description"], attrMap["enabled"], attrMap["vpn_type"],
			attrMap["authentication_mode"], attrMap["compliance_suite"], attrMap["ip_addresses"], attrMap["prefix_length"], attrMap["peer_address"], attrMap["peer_id"], attrMap["psk"], attrMap["connection_initiation_mode"])
}

func testAccNsxtPolicyTGWIPSecVpnSessionPreConditionTemplate() string {
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
  depends_on = [data.nsxt_policy_edge_cluster.test]

  tag {
    scope = "color"
    tag   = "blue"
  }
}

resource "nsxt_policy_gateway_connection" "test" {
  display_name     = "%s"
  tier0_path       = nsxt_policy_tier0_gateway.test.path
  aggregate_routes = ["192.168.240.0/24"]
  depends_on = [nsxt_policy_tier0_gateway.test]
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
 depends_on = [nsxt_policy_gateway_connection.test, nsxt_policy_tier0_gateway.test, data.nsxt_policy_edge_cluster.test]
}

resource "nsxt_policy_project_ip_address_allocation" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  display_name = "ipsec_ipaddr_allocation"
  ip_block     = nsxt_policy_project.test.external_ipv4_blocks[0]
  depends_on = [nsxt_policy_project.test]
}

data "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }
  is_default = "true"
}

resource "nsxt_policy_transit_gateway_attachment" "test" {
  parent_path     = data.nsxt_policy_transit_gateway.test.path
  connection_path = nsxt_policy_gateway_connection.test.path
  display_name    = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
  depends_on = [data.nsxt_policy_transit_gateway.test, nsxt_policy_gateway_connection.test]
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
}

resource "nsxt_policy_transit_gateway_ipsec_vpn_local_endpoint" "test" {
  display_name  = "%s"
  parent_path   =  nsxt_policy_transit_gateway_ipsec_vpn_service.test.path
  description   = "IPSec VPN Local Endpoint"
  local_address = nsxt_policy_project_ip_address_allocation.test.allocation_ips
  local_id      = "10.110.0.0"
  depends_on = [nsxt_policy_transit_gateway_ipsec_vpn_service.test, nsxt_policy_project_ip_address_allocation.test]
}`, getEdgeClusterName(), SessionRelatedResourceName, SessionRelatedResourceName, SessionRelatedResourceName, SessionRelatedResourceName, SessionRelatedResourceName, SessionRelatedResourceName)
}

func TestAccResourceNsxtPolicyTGWIPSecVpnSession_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_transit_gateway_ipsec_vpn_session.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.1.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTGWIPSecVpnSessionCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTGWIPSecVpnSessionRouteBasedMinimalistic(name),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"psk"},
				ImportStateIdFunc:       testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
		},
	})
}

func testAccNsxtPolicyTGWIPSecVpnSessionExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy TGWIPSecVpnSession resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy TGWIPSecVpnSession resource ID not set in resources")
		}
		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := resourceNsxtPolicyTGWIPSecVpnSessionExists(getSessionContextFromParentPath(testAccProvider.Meta(), parentPath), parentPath, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy TGWIPSecVpnSession %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyTGWIPSecVpnSessionCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_transit_gateway_ipsec_vpn_session" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		parentPath := rs.Primary.Attributes["parent_path"]
		exists, err := resourceNsxtPolicyTGWIPSecVpnSessionExists(getSessionContextFromParentPath(testAccProvider.Meta(), parentPath), parentPath, resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy TGWIPSecVpnSession %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyTGWIPSecVpnSessionRouteBasedMinimalistic(displayName string) string {
	attrMap := accTestPolicyTGWIPSecVpnSessionRouteBasedCreateAttributes
	return testAccNsxtPolicyTGWIPSecVpnSessionPreConditionTemplate() +
		fmt.Sprintf(`
resource "nsxt_policy_transit_gateway_ipsec_vpn_session" "test" {
  display_name        = "%s"
  parent_path         = nsxt_policy_transit_gateway_ipsec_vpn_service.test.path
  local_endpoint_path = nsxt_policy_transit_gateway_ipsec_vpn_local_endpoint.test.path
  vpn_type            = "%s"
  peer_address        = "%s"
  peer_id             = "%s"
  ip_addresses        = ["%s"]
  prefix_length       = "%s"
  psk                 = "%s"
  depends_on = [nsxt_policy_transit_gateway_ipsec_vpn_service.test, nsxt_policy_transit_gateway_ipsec_vpn_local_endpoint.test]
}`, displayName, attrMap["vpn_type"], attrMap["peer_address"], attrMap["peer_id"], attrMap["ip_addresses"], attrMap["prefix_length"], attrMap["psk"])
}
