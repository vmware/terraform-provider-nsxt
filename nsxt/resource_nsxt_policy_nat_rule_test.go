/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"testing"
)

var testAccResourcePolicyNATRuleName = "nsxt_policy_nat_rule.test"
var testAccResourcePolicyNATRuleSourceNet = "14.1.1.3"
var testAccResourcePolicyNATRuleDestNet = "15.1.1.3"
var testAccResourcePolicyNATRuleTransNet = "16.1.1.3"

func TestAccResourceNsxtPolicyNATRule_minimalT0(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-nat-rule-basic")
	action := model.PolicyNatRule_ACTION_REFLEXIVE

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyNATRuleCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyNATRuleTier0MinimalCreateTemplate(name, testAccResourcePolicyNATRuleSourceNet, testAccResourcePolicyNATRuleTransNet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyNATRuleExists(testAccResourcePolicyNATRuleName),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "display_name", name),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "source_networks.#", "1"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "translated_networks.#", "1"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "source_networks.0", testAccResourcePolicyNATRuleSourceNet),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "translated_networks.0", testAccResourcePolicyNATRuleTransNet),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "tag.#", "0"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "action", action),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "logging", "false"),
					resource.TestCheckResourceAttrSet(testAccResourcePolicyNATRuleName, "path"),
					resource.TestCheckResourceAttrSet(testAccResourcePolicyNATRuleName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyNATRule_basicT1(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-nat-rule-basic")
	updateName := name + "updated"
	snet := "22.1.1.2"
	dnet := "33.1.1.2"
	tnet := "44.1.1.2"
	action := model.PolicyNatRule_ACTION_DNAT

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyNATRuleCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyNATRuleTier1CreateTemplate(name, action, testAccResourcePolicyNATRuleSourceNet, testAccResourcePolicyNATRuleDestNet, testAccResourcePolicyNATRuleTransNet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyNATRuleExists(testAccResourcePolicyNATRuleName),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "display_name", name),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "destination_networks.#", "1"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "source_networks.#", "1"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "translated_networks.#", "1"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "destination_networks.0", testAccResourcePolicyNATRuleDestNet),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "source_networks.0", testAccResourcePolicyNATRuleSourceNet),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "translated_networks.0", testAccResourcePolicyNATRuleTransNet),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "action", action),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "logging", "false"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "firewall_match", model.PolicyNatRule_FIREWALL_MATCH_MATCH_EXTERNAL_ADDRESS),
					resource.TestCheckResourceAttrSet(testAccResourcePolicyNATRuleName, "path"),
					resource.TestCheckResourceAttrSet(testAccResourcePolicyNATRuleName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyNATRuleTier1CreateTemplate(updateName, action, snet, dnet, tnet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyNATRuleExists(testAccResourcePolicyNATRuleName),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "display_name", updateName),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "destination_networks.#", "1"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "source_networks.#", "1"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "translated_networks.#", "1"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "destination_networks.0", dnet),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "source_networks.0", snet),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "translated_networks.0", tnet),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "action", action),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "logging", "false"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "firewall_match", model.PolicyNatRule_FIREWALL_MATCH_MATCH_EXTERNAL_ADDRESS),
					resource.TestCheckResourceAttrSet(testAccResourcePolicyNATRuleName, "path"),
					resource.TestCheckResourceAttrSet(testAccResourcePolicyNATRuleName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyNATRuleTier1UpdateMultipleSourceNetworksTemplate(updateName, action, testAccResourcePolicyNATRuleSourceNet, snet, dnet, tnet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyNATRuleExists(testAccResourcePolicyNATRuleName),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "display_name", updateName),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "destination_networks.#", "1"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "source_networks.#", "2"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "translated_networks.#", "1"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "destination_networks.0", dnet),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "translated_networks.0", tnet),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "action", action),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "logging", "false"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "firewall_match", model.PolicyNatRule_FIREWALL_MATCH_MATCH_EXTERNAL_ADDRESS),
					resource.TestCheckResourceAttrSet(testAccResourcePolicyNATRuleName, "path"),
					resource.TestCheckResourceAttrSet(testAccResourcePolicyNATRuleName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyNATRule_basicT0(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-nat-rule-basic")
	updateName := name + "updated"
	snet := "22.1.1.2"
	tnet := "44.1.1.2"
	action := model.PolicyNatRule_ACTION_REFLEXIVE

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyNATRuleCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyNATRuleTier0CreateTemplate(name, action, testAccResourcePolicyNATRuleSourceNet, testAccResourcePolicyNATRuleTransNet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyNATRuleExists(testAccResourcePolicyNATRuleName),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "display_name", name),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "source_networks.#", "1"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "translated_networks.#", "1"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "source_networks.0", testAccResourcePolicyNATRuleSourceNet),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "translated_networks.0", testAccResourcePolicyNATRuleTransNet),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "action", action),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "logging", "false"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "firewall_match", model.PolicyNatRule_FIREWALL_MATCH_MATCH_EXTERNAL_ADDRESS),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "scope.#", "1"),
					resource.TestCheckResourceAttrSet(testAccResourcePolicyNATRuleName, "path"),
					resource.TestCheckResourceAttrSet(testAccResourcePolicyNATRuleName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyNATRuleTier0CreateTemplate(updateName, action, snet, tnet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyNATRuleExists(testAccResourcePolicyNATRuleName),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "display_name", updateName),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "source_networks.#", "1"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "translated_networks.#", "1"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "source_networks.0", snet),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "translated_networks.0", tnet),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "tag.#", "2"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "action", action),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "logging", "false"),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "firewall_match", model.PolicyNatRule_FIREWALL_MATCH_MATCH_EXTERNAL_ADDRESS),
					resource.TestCheckResourceAttr(testAccResourcePolicyNATRuleName, "scope.#", "1"),
					resource.TestCheckResourceAttrSet(testAccResourcePolicyNATRuleName, "path"),
					resource.TestCheckResourceAttrSet(testAccResourcePolicyNATRuleName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyNATRule_basicT1Import(t *testing.T) {
	name := fmt.Sprintf("test-nsx-policy-nat-rule-basic")
	action := model.PolicyNatRule_ACTION_DNAT

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyNATRuleCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyNATRuleTier1CreateTemplate(name, action, testAccResourcePolicyNATRuleSourceNet, testAccResourcePolicyNATRuleDestNet, testAccResourcePolicyNATRuleTransNet),
			},
			{
				ResourceName:      testAccResourcePolicyNATRuleName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNSXPolicyNATRuleImporterGetID,
			},
		},
	})
}

func testAccNSXPolicyNATRuleImporterGetID(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources[testAccResourcePolicyNATRuleName]
	if !ok {
		return "", fmt.Errorf("NSX Policy NAT Rule resource %s not found in resources", testAccResourcePolicyNATRuleName)
	}
	resourceID := rs.Primary.ID
	if resourceID == "" {
		return "", fmt.Errorf("NSX Policy NAT Rule resource ID not set in resources ")
	}
	gwPath := rs.Primary.Attributes["gateway_path"]
	if gwPath == "" {
		return "", fmt.Errorf("NSX Policy NAT Rule Gateway Policy Path not set in resources ")
	}
	_, gwID := parseGatewayPolicyPath(gwPath)
	return fmt.Sprintf("%s/%s", gwID, resourceID), nil
}

func testAccNsxtPolicyNATRuleExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy NAT Rule resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy NAT Rule resource ID not set in resources")
		}

		gwPath := rs.Primary.Attributes["gateway_path"]
		isT0, gwID := parseGatewayPolicyPath(gwPath)
		_, err := getNsxtPolicyNATRuleByID(connector, gwID, isT0, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving policy NAT Rule ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyNATRuleCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_nat_rule" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		gwPath := rs.Primary.Attributes["gateway_path"]
		isT0, gwID := parseGatewayPolicyPath(gwPath)
		_, err := getNsxtPolicyNATRuleByID(connector, gwID, isT0, resourceID)
		if err == nil {
			return fmt.Errorf("Policy NAT Rule %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyNATRuleTier0MinimalCreateTemplate(name, sourceNet, translatedNet string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_tier0_gateway" "t0test" {
  display_name              = "terraform-t0-gw"
  description               = "Acceptance Test"
  edge_cluster_path         = "${data.nsxt_policy_edge_cluster.EC.path}"
}

resource "nsxt_policy_nat_rule" "test" {
  display_name         = "%s"
  gateway_path         = "${nsxt_policy_tier0_gateway.t0test.path}"
  action               = "%s"
  source_networks      = ["%s"]
  translated_networks  = ["%s"]
}
`, getEdgeClusterName(), name, model.PolicyNatRule_ACTION_REFLEXIVE, sourceNet, translatedNet)
}

func testAccNsxtPolicyNATRuleTier1CreateTemplate(name string, action string, sourceNet string, destNet string, translatedNet string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_tier1_gateway" "t1test" {
  display_name              = "terraform-t1-gw"
  description               = "Acceptance Test"
  edge_cluster_path         = "${data.nsxt_policy_edge_cluster.EC.path}"
}

resource "nsxt_policy_nat_rule" "test" {
  display_name         = "%s"
  description          = "Acceptance Test"
  gateway_path         = "${nsxt_policy_tier1_gateway.t1test.path}"
  action               = "%s"
  source_networks      = ["%s"]
  destination_networks = ["%s"]
  translated_networks  = ["%s"]
  logging              = false
  firewall_match       = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}
`, getEdgeClusterName(), name, action, sourceNet, destNet, translatedNet, model.PolicyNatRule_FIREWALL_MATCH_MATCH_EXTERNAL_ADDRESS)
}

func testAccNsxtPolicyNATRuleTier1UpdateMultipleSourceNetworksTemplate(name string, action string, sourceNet1 string, sourceNet2 string, destNet string, translatedNet string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

resource "nsxt_policy_tier1_gateway" "t1test" {
  display_name              = "terraform-t1-gw"
  description               = "Acceptance Test"
  edge_cluster_path         = "${data.nsxt_policy_edge_cluster.EC.path}"
}

resource "nsxt_policy_nat_rule" "test" {
  display_name         = "%s"
  description          = "Acceptance Test"
  gateway_path         = "${nsxt_policy_tier1_gateway.t1test.path}"
  action               = "%s"
  source_networks      = ["%s", "%s"]
  destination_networks = ["%s"]
  translated_networks  = ["%s"]
  logging              = false
  firewall_match       = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}
`, getEdgeClusterName(), name, action, sourceNet1, sourceNet2, destNet, translatedNet, model.PolicyNatRule_FIREWALL_MATCH_MATCH_EXTERNAL_ADDRESS)
}

func testAccNsxtPolicyNATRuleTier0CreateTemplate(name string, action string, sourceNet string, translatedNet string) string {
	return fmt.Sprintf(`
data "nsxt_policy_edge_cluster" "EC" {
  display_name = "%s"
}

data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
}

resource "nsxt_policy_tier0_gateway" "t0test" {
  display_name              = "terraform-t0-gw"
  description               = "Acceptance Test"
  ha_mode                   = "ACTIVE_STANDBY"
  edge_cluster_path         = "${data.nsxt_policy_edge_cluster.EC.path}"
}

resource "nsxt_policy_vlan_segment" "test" {
  transport_zone_path = data.nsxt_policy_transport_zone.test.path
  display_name        = "interface_test"
  vlan_ids            = [12]
}

resource "nsxt_policy_tier0_gateway_interface" "test" {
  display_name = "t0gwinterface"
  type         = "SERVICE"
  gateway_path = nsxt_policy_tier0_gateway.t0test.path
  segment_path = nsxt_policy_vlan_segment.test.path
  subnets      = ["1.1.12.2/24"]
}

resource "nsxt_policy_nat_rule" "test" {
  display_name         = "%s"
  description          = "Acceptance Test"
  gateway_path         = "${nsxt_policy_tier0_gateway.t0test.path}"
  action               = "%s"
  source_networks      = ["%s"]
  translated_networks  = ["%s"]
  logging              = false
  firewall_match       = "%s"
  scope                = ["${nsxt_policy_tier0_gateway_interface.test.path}"]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }

  tag {
    scope = "scope2"
    tag   = "tag2"
  }
}

`, getEdgeClusterName(), getVlanTransportZoneName(), name, action, sourceNet, translatedNet, model.PolicyNatRule_FIREWALL_MATCH_MATCH_EXTERNAL_ADDRESS)
}
