// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyGatewaySecurityConfig_tier0(t *testing.T) {
	testDataSourceName := "data.nsxt_policy_gateway_security_config.test"
	testResourceName := "nsxt_policy_gateway_security_config.test"
	gwName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewaySecurityConfigTier0ReadTemplate(gwName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSourceName, "tier0_id"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
					resource.TestCheckResourceAttr(testDataSourceName, "idps_enabled", "true"),
					resource.TestCheckResourceAttrPair(testDataSourceName, "tier0_id", testResourceName, "tier0_id"),
					resource.TestCheckResourceAttrPair(testDataSourceName, "idps_enabled", testResourceName, "idps_enabled"),
					resource.TestCheckResourceAttrPair(testDataSourceName, "path", testResourceName, "path"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyGatewaySecurityConfig_tier1(t *testing.T) {
	testDataSourceName := "data.nsxt_policy_gateway_security_config.test"
	testResourceName := "nsxt_policy_gateway_security_config.test"
	gwName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewaySecurityConfigTier1ReadTemplate(gwName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSourceName, "tier1_id"),
					resource.TestCheckResourceAttrSet(testDataSourceName, "path"),
					resource.TestCheckResourceAttr(testDataSourceName, "idfw_enabled", "true"),
					resource.TestCheckResourceAttrPair(testDataSourceName, "tier1_id", testResourceName, "tier1_id"),
					resource.TestCheckResourceAttrPair(testDataSourceName, "idfw_enabled", testResourceName, "idfw_enabled"),
					resource.TestCheckResourceAttrPair(testDataSourceName, "path", testResourceName, "path"),
				),
			},
		},
	})
}

func testAccNsxtPolicyGatewaySecurityConfigTier0ReadTemplate(gwName string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_tier0_gateway" "test" {
  display_name = "%s"
  ha_mode      = "ACTIVE_ACTIVE"
}

resource "nsxt_policy_gateway_security_config" "test" {
  tier0_id     = nsxt_policy_tier0_gateway.test.id
  idps_enabled = true
}

data "nsxt_policy_gateway_security_config" "test" {
  tier0_id = nsxt_policy_gateway_security_config.test.tier0_id
}
`, gwName)
}

// testAccNsxtPolicyGatewaySecurityConfigTier1ReadTemplate uses idfw_enabled
// which works without an edge cluster (service router) on Tier-1.
func testAccNsxtPolicyGatewaySecurityConfigTier1ReadTemplate(gwName string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_tier1_gateway" "test" {
  display_name = "%s"
}

resource "nsxt_policy_gateway_security_config" "test" {
  tier1_id     = nsxt_policy_tier1_gateway.test.id
  idfw_enabled = true
}

data "nsxt_policy_gateway_security_config" "test" {
  tier1_id = nsxt_policy_gateway_security_config.test.tier1_id
}
`, gwName)
}
