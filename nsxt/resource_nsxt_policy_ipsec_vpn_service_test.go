/* Copyright © 2022 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var accTestPolicyIPSecVpnServiceCreateAttributes = map[string]string{
	"display_name":  getAccTestResourceName(),
	"description":   "terraform created",
	"enabled":       "true",
	"ha_sync":       "true",
	"ike_log_level": "INFO",
	"sources":       "192.168.10.0/24",
	"destinations":  "192.169.10.0/24",
	"action":        "BYPASS",
}

var accTestPolicyIPSecVpnServiceUpdateAttributes = map[string]string{
	"display_name":  getAccTestResourceName(),
	"description":   "terraform updated",
	"enabled":       "false",
	"ha_sync":       "true",
	"ike_log_level": "DEBUG",
	"sources":       "192.170.10.0/24",
	"destinations":  "192.171.10.0/24",
	"action":        "BYPASS",
}

func TestAccResourceNsxtPolicyIPSecVpnService_basic(t *testing.T) {
	testResourceName := "nsxt_policy_ipsec_vpn_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPSecVpnServiceCheckDestroy(state, accTestPolicyIPSecVpnServiceUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPSecVpnServiceTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPSecVpnServiceExists(accTestPolicyIPSecVpnServiceCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPSecVpnServiceCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPSecVpnServiceCreateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_path"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyIPSecVpnServiceCreateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_sync", accTestPolicyIPSecVpnServiceCreateAttributes["ha_sync"]),
					resource.TestCheckResourceAttr(testResourceName, "ike_log_level", accTestPolicyIPSecVpnServiceCreateAttributes["ike_log_level"]),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.0.sources.0", accTestPolicyIPSecVpnServiceCreateAttributes["sources"]),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.0.destinations.0", accTestPolicyIPSecVpnServiceCreateAttributes["destinations"]),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.0.action", accTestPolicyIPSecVpnServiceCreateAttributes["action"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIPSecVpnServiceTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPSecVpnServiceExists(accTestPolicyIPSecVpnServiceUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyIPSecVpnServiceUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyIPSecVpnServiceUpdateAttributes["description"]),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_path"),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestPolicyIPSecVpnServiceUpdateAttributes["enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_sync", accTestPolicyIPSecVpnServiceUpdateAttributes["ha_sync"]),
					resource.TestCheckResourceAttr(testResourceName, "ike_log_level", accTestPolicyIPSecVpnServiceUpdateAttributes["ike_log_level"]),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.0.sources.0", accTestPolicyIPSecVpnServiceUpdateAttributes["sources"]),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.0.destinations.0", accTestPolicyIPSecVpnServiceUpdateAttributes["destinations"]),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.0.action", accTestPolicyIPSecVpnServiceUpdateAttributes["action"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyIPSecVpnServiceMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIPSecVpnServiceExists(accTestPolicyIPSecVpnServiceCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_path"),
					resource.TestCheckResourceAttr(testResourceName, "bypass_rule.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyIPSecVpnService_Import(t *testing.T) {
	resourceName := "nsxt_policy_ipsec_vpn_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyIPSecVpnServiceCheckDestroy(state, accTestPolicyIPSecVpnServiceCreateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIPSecVpnServiceTemplate(true),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNsxtPolicyIPSecVpnServiceImporterGetID,
			},
		},
	})
}

func testAccNsxtPolicyIPSecVpnServiceImporterGetID(s *terraform.State) (string, error) {
	resourceName := "nsxt_policy_ipsec_vpn_service.test"
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return "", fmt.Errorf("Policy IPSecVpnService resource %s not found in resources", resourceName)
	}
	resourceID := rs.Primary.ID
	if resourceID == "" {
		return "", fmt.Errorf("Policy IPSecVpnService resource ID not set in resources")
	}
	localeServicePath := rs.Primary.Attributes["locale_service_path"]
	return fmt.Sprintf("%s/ipsec-vpn-services/%s", localeServicePath, resourceID), nil
}

func testAccNsxtPolicyIPSecVpnServiceExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy IPSecVpnService resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy IPSecVpnService resource ID not set in resources")
		}

		localeServicePath := rs.Primary.Attributes["locale_service_path"]
		isT0, gwID, localeServiceID, err := parseLocaleServicePolicyPath(localeServicePath)

		if err != nil {
			return nil
		}
		_, err1 := getNsxtPolicyIPSecVpnServiceByID(connector, gwID, isT0, localeServiceID, resourceID, testAccIsGlobalManager())
		if err1 != nil {
			return fmt.Errorf("Policy IPSecVpnService %s does not exist", displayName)
		}

		return nil
	}
}

func testAccNsxtPolicyIPSecVpnServiceCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_ipsec_vpn_service" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		localeServicePath := rs.Primary.Attributes["locale_service_path"]
		isT0, gwID, localeServiceID, err := parseLocaleServicePolicyPath(localeServicePath)

		if err != nil {
			return nil
		}

		_, err1 := getNsxtPolicyIPSecVpnServiceByID(connector, gwID, isT0, localeServiceID, resourceID, testAccIsGlobalManager())
		if err1 == nil {
			return fmt.Errorf("Policy IPSecVpnService %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyIPSecVpnServiceTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyIPSecVpnServiceCreateAttributes
	} else {
		attrMap = accTestPolicyIPSecVpnServiceUpdateAttributes
	}
	return testAccNsxtPolicyEdgeClusterReadTemplate(getEdgeClusterName()) +
		testAccNsxtPolicyTier0WithEdgeClusterForVPN("test") + fmt.Sprintf(`
resource "nsxt_policy_ipsec_vpn_service" "test" {
	display_name                   = "%s"
	description                    = "%s"
	locale_service_path   		   = one(nsxt_policy_tier0_gateway.test.locale_service).path
	enabled                        = "%s"
	ha_sync                        = "%s"
	ike_log_level                  = "%s"
	bypass_rule {
		sources                    = ["%s"]
		destinations               = ["%s"]
		action                     = "%s"
	}

	tag {
	scope = "scope1"
	tag   = "tag1"
	}
   }`, attrMap["display_name"], attrMap["description"], attrMap["enabled"], attrMap["ha_sync"], attrMap["ike_log_level"], attrMap["sources"], attrMap["destinations"], attrMap["action"])
}

func testAccNsxtPolicyIPSecVpnServiceMinimalistic() string {
	return testAccNsxtPolicyEdgeClusterReadTemplate(getEdgeClusterName()) +
		testAccNsxtPolicyTier0WithEdgeClusterForVPN("test") + fmt.Sprintf(`
   resource "nsxt_policy_ipsec_vpn_service" "test" {
	 display_name          = "%s"
	 locale_service_path   = one(nsxt_policy_tier0_gateway.test.locale_service).path
   }`, accTestPolicyIPSecVpnServiceUpdateAttributes["display_name"])
}
