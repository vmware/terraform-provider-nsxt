/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

var testAccResourcePolicyGatewayDNSForwarderName = "nsxt_policy_gateway_dns_forwarder.test"
var testAccPolicyDNSForwarderHelperNames = [2]string{
	getAccTestResourceName(),
	getAccTestResourceName(),
}

var accTestPolicyGatewayDNSForwarderCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
	"listener_ip":  "10.2.2.12",
	"enabled":      "true",
	"log_level":    model.PolicyDnsForwarder_LOG_LEVEL_FATAL,
	"cache_size":   "2048",
}

var accTestPolicyGatewayDNSForwarderUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
	"listener_ip":  "10.2.2.15",
	"enabled":      "false",
	"log_level":    model.PolicyDnsForwarder_LOG_LEVEL_DEBUG,
	"cache_size":   "4096",
}

func TestAccResourceNsxtPolicyGatewayDNSForwarder_tier0(t *testing.T) {
	testAccResourceNsxtPolicyGatewayDNSForwarder(t, true)
}

func TestAccResourceNsxtPolicyGatewayDNSForwarder_tier1(t *testing.T) {
	testAccResourceNsxtPolicyGatewayDNSForwarder(t, false)
}

func testAccResourceNsxtPolicyGatewayDNSForwarder(t *testing.T, isT0 bool) {
	resourceName := testAccResourcePolicyGatewayDNSForwarderName
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGatewayDNSForwarderCheckDestroy(state, accTestPolicyGatewayDNSForwarderUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewayDNSForwarderTemplate(isT0, true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayDNSForwarderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", accTestPolicyGatewayDNSForwarderCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(resourceName, "description", accTestPolicyGatewayDNSForwarderCreateAttributes["description"]),
					resource.TestCheckResourceAttr(resourceName, "listener_ip", accTestPolicyGatewayDNSForwarderCreateAttributes["listener_ip"]),
					resource.TestCheckResourceAttr(resourceName, "enabled", accTestPolicyGatewayDNSForwarderCreateAttributes["enabled"]),
					resource.TestCheckResourceAttr(resourceName, "log_level", accTestPolicyGatewayDNSForwarderCreateAttributes["log_level"]),
					resource.TestCheckResourceAttrSet(resourceName, "default_forwarder_zone_path"),
					resource.TestCheckResourceAttr(resourceName, "conditional_forwarder_zone_paths.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cache_size", accTestPolicyGatewayDNSForwarderCreateAttributes["cache_size"]),
					resource.TestCheckResourceAttr(resourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(resourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyGatewayDNSForwarderTemplate(isT0, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewayDNSForwarderExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", accTestPolicyGatewayDNSForwarderUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(resourceName, "description", accTestPolicyGatewayDNSForwarderUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(resourceName, "listener_ip", accTestPolicyGatewayDNSForwarderUpdateAttributes["listener_ip"]),
					resource.TestCheckResourceAttr(resourceName, "enabled", accTestPolicyGatewayDNSForwarderUpdateAttributes["enabled"]),
					resource.TestCheckResourceAttr(resourceName, "log_level", accTestPolicyGatewayDNSForwarderUpdateAttributes["log_level"]),
					resource.TestCheckResourceAttrSet(resourceName, "default_forwarder_zone_path"),
					resource.TestCheckResourceAttr(resourceName, "conditional_forwarder_zone_paths.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cache_size", accTestPolicyGatewayDNSForwarderUpdateAttributes["cache_size"]),
					resource.TestCheckResourceAttr(resourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "path"),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(resourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyGatewayDNSForwarder_importTier0(t *testing.T) {
	testAccResourceNsxtPolicyGatewayDNSForwarderImport(t, true)
}

func TestAccResourceNsxtPolicyGatewayDNSForwarder_importTier1(t *testing.T) {
	testAccResourceNsxtPolicyGatewayDNSForwarderImport(t, false)
}

func testAccResourceNsxtPolicyGatewayDNSForwarderImport(t *testing.T, isT0 bool) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyGatewayDNSForwarderCheckDestroy(state, accTestPolicyGatewayDNSForwarderCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewayDNSForwarderMinimalistic(isT0),
			},
			{
				ResourceName:      testAccResourcePolicyGatewayDNSForwarderName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNSXPolicyGatewayDNSForwarderImporterGetID,
			},
		},
	})
}

func testAccNSXPolicyGatewayDNSForwarderImporterGetID(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources[testAccResourcePolicyGatewayDNSForwarderName]
	if !ok {
		return "", fmt.Errorf("NSX Policy Gateway DNS Forwarder resource %s not found in resources", testAccResourcePolicyGatewayDNSForwarderName)
	}
	gwPath := rs.Primary.Attributes["gateway_path"]
	if gwPath == "" {
		return "", fmt.Errorf("NSX Policy Gateway DNS Forwarder Gateway Policy Path not set in resources ")
	}
	return gwPath, nil
}

func testAccNsxtPolicyGatewayDNSForwarderExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Gateway DNS Forwarder resource %s not found in resources", resourceName)
		}

		gwPath := rs.Primary.Attributes["gateway_path"]
		isT0, gwID := parseGatewayPolicyPath(gwPath)

		_, err := policyGatewayDNSForwarderGet(connector, gwID, isT0, testAccIsGlobalManager())
		if err != nil {
			return fmt.Errorf("Policy Gateway DNS Forwarder resource does not exist on %s", gwPath)
		}

		return nil
	}
}

func testAccNsxtPolicyGatewayDNSForwarderCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_gateway_dns_forwarder" {
			continue
		}

		gwPath := rs.Primary.Attributes["gateway_path"]
		isT0, gwID := parseGatewayPolicyPath(gwPath)

		_, err := policyGatewayDNSForwarderGet(connector, gwID, isT0, testAccIsGlobalManager())
		if err == nil {
			return fmt.Errorf("Policy Gateway DNS Forwarder %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyGatewayDNSForwarderPrerequisites(names [2]string, isT0 bool) string {
	return testAccNsxtPolicyEdgeClusterReadTemplate(getEdgeClusterName()) +
		testAccNsxtPolicyGatewayWithEdgeClusterTemplate("test", isT0, true) +
		fmt.Sprintf(`
resource "nsxt_policy_dns_forwarder_zone" "default" {
  display_name     = "%s"
  upstream_servers = ["1.1.1.1"]
}

resource "nsxt_policy_dns_forwarder_zone" "fqdn" {
  display_name     = "%s"
  upstream_servers = ["2.1.1.1"]
  dns_domain_names = ["conditional.domain.org"]
}`, names[0], names[1])
}

func testAccNsxtPolicyGatewayDNSForwarderTemplate(isT0 bool, createFlow bool) string {
	whyDoesGoNeedToBeSoComplicated := map[bool]int8{false: 1, true: 0}
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyGatewayDNSForwarderCreateAttributes
	} else {
		attrMap = accTestPolicyGatewayDNSForwarderUpdateAttributes
	}
	return testAccNsxtPolicyGatewayDNSForwarderPrerequisites(testAccPolicyDNSForwarderHelperNames, isT0) + fmt.Sprintf(`
resource "nsxt_policy_gateway_dns_forwarder" "test" {
  display_name = "%s"
  description  = "%s"
  gateway_path = nsxt_policy_tier%d_gateway.test.path
  listener_ip  = "%s"
  enabled      = %s
  log_level    = "%s"
  cache_size   = %s

  default_forwarder_zone_path      = nsxt_policy_dns_forwarder_zone.default.path
  conditional_forwarder_zone_paths = [nsxt_policy_dns_forwarder_zone.fqdn.path]

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
`, attrMap["display_name"], attrMap["description"], whyDoesGoNeedToBeSoComplicated[isT0], attrMap["listener_ip"], attrMap["enabled"], attrMap["log_level"], attrMap["cache_size"])
}

func testAccNsxtPolicyGatewayDNSForwarderMinimalistic(isT0 bool) string {
	whyDoesGoNeedToBeSoComplicated := map[bool]int8{false: 1, true: 0}
	return testAccNsxtPolicyGatewayDNSForwarderPrerequisites(testAccPolicyDNSForwarderHelperNames, isT0) + fmt.Sprintf(`
resource "nsxt_policy_gateway_dns_forwarder" "test" {
  display_name = "%s"
  gateway_path = nsxt_policy_tier%d_gateway.test.path
  listener_ip  = "78.2.1.12"

  default_forwarder_zone_path      = nsxt_policy_dns_forwarder_zone.default.path
}
`, accTestPolicyGatewayDNSForwarderCreateAttributes["display_name"], whyDoesGoNeedToBeSoComplicated[isT0])
}
