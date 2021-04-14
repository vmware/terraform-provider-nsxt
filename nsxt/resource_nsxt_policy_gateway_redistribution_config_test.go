/* Copyright © 2021 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccNsxtPolicyGatewayRedistributionHelperName = getAccTestResourceName()

func TestAccResourceNsxtPolicyGatewayRedistributionConfig_basic(t *testing.T) {
	tier0ResourceName := "nsxt_policy_tier0_gateway.test"
	testResourceName := "nsxt_policy_gateway_redistribution_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyTier0CheckNoRedistribution(testAccNsxtPolicyGatewayRedistributionHelperName),
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewayRedistributionCreateTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(tier0ResourceName),
					resource.TestCheckResourceAttr(testResourceName, "bgp_enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "ospf_enabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.name", "test-rule-1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.types.#", "3"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "site_path"),
				),
			},
			{
				Config: testAccNsxtPolicyGatewayRedistributionUpdateTemplate(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(tier0ResourceName),
					resource.TestCheckResourceAttr(testResourceName, "bgp_enabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "ospf_enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.name", "test-rule-1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.types.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.name", "test-rule-2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.types.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "site_path"),
				),
			},
			{
				Config: testAccNsxtPolicyGatewayRedistributionUpdate2Template(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier0Exists(tier0ResourceName),
					resource.TestCheckResourceAttr(testResourceName, "bgp_enabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "ospf_enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "site_path"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyGatewayRedistributionConfig_importBasic(t *testing.T) {
	testResourceName := "nsxt_policy_gateway_redistribution_config.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyTier0CheckNoRedistribution(testAccNsxtPolicyGatewayRedistributionHelperName),
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewayRedistributionCreateTemplate(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNSXPolicyRedistributionConfigImporterGetID,
			},
		},
	})
}

func testAccNsxtPolicyTier0CheckNoRedistribution(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Tier0 Redistribution config resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy Tier0 Redistribution config resource ID not set in resources")
		}

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		localeServiceID := rs.Primary.Attributes["locale_service_id"]
		gwID := rs.Primary.Attributes["gateway_id"]
		localeService := policyTier0GetLocaleService(gwID, localeServiceID, connector, testAccIsGlobalManager())
		if localeService == nil {
			return fmt.Errorf("Error while retrieving Locale Service %s for Gateway %s", localeServiceID, gwID)
		}

		if localeService.RouteRedistributionConfig == nil {
			return fmt.Errorf("Error while retrieving policy Tier0 Redistribution config %s. HaVipConfigs is empty", resourceID)
		}

		return nil
	}

}

func testAccNSXPolicyRedistributionConfigImporterGetID(s *terraform.State) (string, error) {
	testResourceName := "nsxt_policy_gateway_redistribution_config.test"
	rs, ok := s.RootModule().Resources[testResourceName]
	if !ok {
		return "", fmt.Errorf("NSX Policy Redistribution config resource %s not found in resources", testResourceName)
	}
	resourceID := rs.Primary.ID
	if resourceID == "" {
		return "", fmt.Errorf("NSX Policy Redistribution config resource ID not set in resources ")
	}
	gwID := rs.Primary.Attributes["gateway_id"]
	if gwID == "" {
		return "", fmt.Errorf("NSX Policy Redistribution config Tier0 Gateway ID not set in resources ")
	}
	localeServiceID := rs.Primary.Attributes["locale_service_id"]
	if localeServiceID == "" {
		return "", fmt.Errorf("NSX Policy HA Vip config Tier0 Gateway locale service ID not set in resources ")
	}

	return fmt.Sprintf("%s/%s", gwID, localeServiceID), nil
}

func testAccNsxtPolicyGatewayRedistributionPrerequisites() string {
	return testAccNsxtPolicyGatewayFabricDeps(false) + fmt.Sprintf(`
resource "nsxt_policy_tier0_gateway" "test" {
  display_name      = "%s"
  %s
}`, testAccNsxtPolicyGatewayRedistributionHelperName, testAccNsxtPolicyTier0EdgeClusterTemplate())
}

func getAccTestSitePathConfig() string {
	if testAccIsGlobalManager() {
		return `site_path = data.nsxt_policy_site.test.path`
	}

	return ""
}

func testAccNsxtPolicyGatewayRedistributionCreateTemplate() string {
	return testAccNsxtPolicyGatewayRedistributionPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_gateway_redistribution_config" "test" {
  gateway_path = nsxt_policy_tier0_gateway.test.path
  %s

  bgp_enabled  = false
  ospf_enabled = true
  rule {
      name = "test-rule-1"
      types = ["TIER0_SEGMENT", "TIER0_EVPN_TEP_IP", "TIER1_CONNECTED"]
  }
}`, getAccTestSitePathConfig())
}

func testAccNsxtPolicyGatewayRedistributionUpdateTemplate() string {
	return testAccNsxtPolicyGatewayRedistributionPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_gateway_redistribution_config" "test" {
  gateway_path = nsxt_policy_tier0_gateway.test.path
  %s

  bgp_enabled  = true
  ospf_enabled = false
  rule {
      name = "test-rule-1"
      types = ["TIER1_CONNECTED"]
  }
  rule {
      name  = "test-rule-2"
      types = ["TIER1_LB_VIP"]
  }
}`, getAccTestSitePathConfig())
}

func testAccNsxtPolicyGatewayRedistributionUpdate2Template() string {
	return testAccNsxtPolicyGatewayRedistributionPrerequisites() + fmt.Sprintf(`
resource "nsxt_policy_gateway_redistribution_config" "test" {
  gateway_path = nsxt_policy_tier0_gateway.test.path
  %s

  bgp_enabled  = true
  ospf_enabled = false
}`, getAccTestSitePathConfig())
}
