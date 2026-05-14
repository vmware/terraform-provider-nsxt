// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_1s"
)

func TestAccResourceNsxtPolicyGatewaySecurityConfig_tier0Basic(t *testing.T) {
	testResourceName := "nsxt_policy_gateway_security_config.test"
	gwName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyGatewaySecurityConfigTier0CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewaySecurityConfigTier0Template(gwName, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewaySecurityConfigTier0Exists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "tier0_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "idps_enabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "idfw_enabled", "false"),
				),
			},
			{
				Config: testAccNsxtPolicyGatewaySecurityConfigTier0Template(gwName, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewaySecurityConfigTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "idps_enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "idfw_enabled", "true"),
				),
			},
			{
				Config: testAccNsxtPolicyGatewaySecurityConfigTier0Template(gwName, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewaySecurityConfigTier0Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "idps_enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "idfw_enabled", "false"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyGatewaySecurityConfig_tier1Basic(t *testing.T) {
	testResourceName := "nsxt_policy_gateway_security_config.test"
	gwName := getAccTestResourceName()

	// Note: IDPS and MALWAREPREVENTION on Tier-1 require an Edge Cluster (Service Router).
	// Tests use idfw_enabled and tls_enabled which work without an edge cluster.
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyGatewaySecurityConfigTier1CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewaySecurityConfigTier1Template(gwName, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewaySecurityConfigTier1Exists(testResourceName),
					resource.TestCheckResourceAttrSet(testResourceName, "tier1_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttr(testResourceName, "idfw_enabled", "true"),
					resource.TestCheckResourceAttr(testResourceName, "tls_enabled", "false"),
				),
			},
			{
				Config: testAccNsxtPolicyGatewaySecurityConfigTier1Template(gwName, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyGatewaySecurityConfigTier1Exists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "idfw_enabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "tls_enabled", "true"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyGatewaySecurityConfig_tier0Import(t *testing.T) {
	testResourceName := "nsxt_policy_gateway_security_config.test"
	gwName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyGatewaySecurityConfigTier0CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewaySecurityConfigTier0Template(gwName, true, false),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revision"},
			},
		},
	})
}

func TestAccResourceNsxtPolicyGatewaySecurityConfig_tier1Import(t *testing.T) {
	testResourceName := "nsxt_policy_gateway_security_config.test"
	gwName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccNsxtPolicyGatewaySecurityConfigTier1CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyGatewaySecurityConfigTier1Template(gwName, true, false),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"revision"},
			},
		},
	})
}

func testAccNsxtPolicyGatewaySecurityConfigTier0Exists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Gateway Security Config resource %s not found in resources", resourceName)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("Gateway Security Config ID not set")
		}

		// id is in format "tier0/<gateway-id>"
		parts := strings.SplitN(id, "/", 2)
		if len(parts) != 2 || parts[0] != "tier0" {
			return fmt.Errorf("unexpected resource ID format %q: expected 'tier0/<id>'", id)
		}
		gwID := parts[1]

		connector := getPolicyConnector(testAccProvider.Meta())
		client := tier_0s.NewSecurityConfigClient(connector)
		_, err := client.Get(gwID, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("error retrieving Tier0 Gateway Security Config for gateway %s: %v", gwID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyGatewaySecurityConfigTier1Exists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Gateway Security Config resource %s not found in resources", resourceName)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("Gateway Security Config ID not set")
		}

		// id is in format "tier1/<gateway-id>"
		parts := strings.SplitN(id, "/", 2)
		if len(parts) != 2 || parts[0] != "tier1" {
			return fmt.Errorf("unexpected resource ID format %q: expected 'tier1/<id>'", id)
		}
		gwID := parts[1]

		connector := getPolicyConnector(testAccProvider.Meta())
		client := tier_1s.NewSecurityConfigClient(connector)
		_, err := client.Get(gwID, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("error retrieving Tier1 Gateway Security Config for gateway %s: %v", gwID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyGatewaySecurityConfigTier0CheckDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_gateway_security_config" {
			continue
		}

		id := rs.Primary.ID
		parts := strings.SplitN(id, "/", 2)
		if len(parts) != 2 || parts[0] != "tier0" {
			continue
		}
		gwID := parts[1]

		connector := getPolicyConnector(testAccProvider.Meta())
		client := tier_0s.NewSecurityConfigClient(connector)
		config, err := client.Get(gwID, nil, nil, nil, nil, nil, nil)
		if err != nil {
			if isNotFoundError(err) {
				continue
			}
			continue
		}

		for _, f := range config.Features {
			if f.Enable != nil && *f.Enable {
				return fmt.Errorf("Tier0 Gateway Security Config %s still has enabled feature %v after destroy", gwID, f.Feature)
			}
		}
	}
	return nil
}

func testAccNsxtPolicyGatewaySecurityConfigTier1CheckDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_policy_gateway_security_config" {
			continue
		}

		id := rs.Primary.ID
		parts := strings.SplitN(id, "/", 2)
		if len(parts) != 2 || parts[0] != "tier1" {
			continue
		}
		gwID := parts[1]

		connector := getPolicyConnector(testAccProvider.Meta())
		client := tier_1s.NewSecurityConfigClient(connector)
		config, err := client.Get(gwID, nil, nil, nil, nil, nil, nil)
		if err != nil {
			if isNotFoundError(err) {
				continue
			}
			continue
		}

		for _, f := range config.Features {
			if f.Enable != nil && *f.Enable {
				return fmt.Errorf("Tier1 Gateway Security Config %s still has enabled feature %v after destroy", gwID, f.Feature)
			}
		}
	}
	return nil
}

func testAccNsxtPolicyGatewaySecurityConfigTier0Template(gwName string, idpsEnabled, idfwEnabled bool) string {
	return fmt.Sprintf(`
resource "nsxt_policy_tier0_gateway" "test" {
  display_name = "%s"
  ha_mode      = "ACTIVE_ACTIVE"
}

resource "nsxt_policy_gateway_security_config" "test" {
  tier0_id     = nsxt_policy_tier0_gateway.test.id
  idps_enabled = %t
  idfw_enabled = %t
}
`, gwName, idpsEnabled, idfwEnabled)
}

// testAccNsxtPolicyGatewaySecurityConfigTier1Template uses idfw_enabled and
// tls_enabled which work without an edge cluster (service router).
// IDPS and MALWAREPREVENTION on Tier-1 require an edge cluster for the service router.
func testAccNsxtPolicyGatewaySecurityConfigTier1Template(gwName string, idfwEnabled, tlsEnabled bool) string {
	return fmt.Sprintf(`
resource "nsxt_policy_tier1_gateway" "test" {
  display_name = "%s"
}

resource "nsxt_policy_gateway_security_config" "test" {
  tier1_id     = nsxt_policy_tier1_gateway.test.id
  idfw_enabled = %t
  tls_enabled  = %t
}
`, gwName, idfwEnabled, tlsEnabled)
}
