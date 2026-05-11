// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyIntrusionServiceGatewayPolicy_basic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_gateway_policy.by_name"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyReadByName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					// Check that rules are present
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "test-gw-rule"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", "IN_OUT"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.rule_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.nsx_id"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.scope.#", "1"),
					// Verify oversubscription is excluded for Gateway IDPS rules
					resource.TestCheckNoResourceAttr(testResourceName, "rule.0.oversubscription"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIntrusionServiceGatewayPolicy_byCategory(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_gateway_policy.by_category"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyReadByCategory(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					// Validate embedded rules
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "test-gw-rule-category"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", "IN_OUT"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.rule_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.nsx_id"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.scope.#", "1"),
					// Verify oversubscription is excluded for Gateway IDPS rules
					resource.TestCheckNoResourceAttr(testResourceName, "rule.0.oversubscription"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIntrusionServiceGatewayPolicy_byID(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_gateway_policy.by_id"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyReadByID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					// Validate embedded rules
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "test-gw-rule-id"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", "IN_OUT"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.rule_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.nsx_id"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.scope.#", "1"),
					// Verify oversubscription is excluded for Gateway IDPS rules
					resource.TestCheckNoResourceAttr(testResourceName, "rule.0.oversubscription"),
				),
			},
		},
	})
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyReadByName(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_tier1_gateway" "test" {
  display_name = "test-t1-gw-%s"
  description  = "Test Tier1 gateway for IDPS policy tests"
}

resource "nsxt_policy_group" "gw_source" {
  display_name = "test-gw-source-%s"
  description  = "Gateway policy source group"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["172.16.10.0/24"]
    }
  }
}

resource "nsxt_policy_group" "gw_dest" {
  display_name = "test-gw-dest-%s"
  description  = "Gateway policy destination group"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["172.16.20.0/24"]
    }
  }
}

resource "nsxt_policy_intrusion_service_gateway_policy" "ids_gw_policy_for_ds" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "LocalGatewayRules"
  sequence_number = 3
  
  rule {
    display_name       = "test-gw-rule"
    notes              = "Test Gateway rule for data source"
    action             = "DETECT"
    direction          = "IN_OUT"
    source_groups      = [nsxt_policy_group.gw_source.path]
    destination_groups = [nsxt_policy_group.gw_dest.path]
    scope              = [nsxt_policy_tier1_gateway.test.path]
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.test.path]
  }
}

data "nsxt_policy_intrusion_service_gateway_policy" "by_name" {
  display_name = nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_ds.display_name
  depends_on   = [nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_ds]
}`, getAccTestResourceName(), name, name, name)
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyReadByCategory(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_tier1_gateway" "test" {
  display_name = "test-t1-gw-%s"
  description  = "Test Tier1 gateway for IDPS policy tests"
}

resource "nsxt_policy_group" "gw_source" {
  display_name = "test-gw-source-%s"
  description  = "Gateway policy source group"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["172.16.10.0/24"]
    }
  }
}

resource "nsxt_policy_group" "gw_dest" {
  display_name = "test-gw-dest-%s"
  description  = "Gateway policy destination group"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["172.16.20.0/24"]
    }
  }
}

resource "nsxt_policy_intrusion_service_gateway_policy" "ids_gw_policy_for_ds" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "LocalGatewayRules"
  sequence_number = 3
  
  rule {
    display_name       = "test-gw-rule-category"
    notes              = "Test Gateway rule for category-based lookup"
    action             = "DETECT"
    direction          = "IN_OUT"
    source_groups      = [nsxt_policy_group.gw_source.path]
    destination_groups = [nsxt_policy_group.gw_dest.path]
    scope              = [nsxt_policy_tier1_gateway.test.path]
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.test.path]
  }
}

data "nsxt_policy_intrusion_service_gateway_policy" "by_category" {
  display_name = nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_ds.display_name
  category     = "LocalGatewayRules"
  depends_on   = [nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_ds]
}`, getAccTestResourceName(), name, name, name)
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyReadByID(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_tier1_gateway" "test" {
  display_name = "test-t1-gw-%s"
  description  = "Test Tier1 gateway for IDPS policy tests"
}

resource "nsxt_policy_group" "gw_source" {
  display_name = "test-gw-source-%s"
  description  = "Gateway policy source group"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["172.16.10.0/24"]
    }
  }
}

resource "nsxt_policy_group" "gw_dest" {
  display_name = "test-gw-dest-%s"
  description  = "Gateway policy destination group"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["172.16.20.0/24"]
    }
  }
}

resource "nsxt_policy_intrusion_service_gateway_policy" "ids_gw_policy_for_ds" {
  display_name    = "%s"
  description     = "Acceptance Test"
  category        = "LocalGatewayRules"
  sequence_number = 3
  
  rule {
    display_name       = "test-gw-rule-id"
    notes              = "Test Gateway rule for ID-based lookup"
    action             = "DETECT"
    direction          = "IN_OUT"
    source_groups      = [nsxt_policy_group.gw_source.path]
    destination_groups = [nsxt_policy_group.gw_dest.path]
    scope              = [nsxt_policy_tier1_gateway.test.path]
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.test.path]
  }
}

data "nsxt_policy_intrusion_service_gateway_policy" "by_id" {
  id         = nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_ds.id
  depends_on = [nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_ds]
}`, getAccTestResourceName(), name, name, name)
}

func TestAccDataSourceNsxtPolicyIntrusionServiceGatewayPolicy_withDetailedRules(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_gateway_policy.with_detailed_rules"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyWithDetailedRules(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Comprehensive IDPS Gateway policy test"),
					resource.TestCheckResourceAttr(testResourceName, "category", "LocalGatewayRules"),
					resource.TestCheckResourceAttr(testResourceName, "domain", "default"),
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					// Check that rules are present with detailed configuration
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "2"),
					// First rule - comprehensive Gateway configuration
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "comprehensive-gw-rule"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action", "DETECT_PREVENT"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.direction", "OUT"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.logged", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.disabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.sources_excluded", "true"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destinations_excluded", "false"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ip_version", "IPV6"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.sequence_number"),
					// Validate Gateway rule field arrays
					resource.TestCheckResourceAttr(testResourceName, "rule.0.source_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.destination_groups.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.services.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.scope.#", "1"), // Single T1 gateway
					resource.TestCheckResourceAttr(testResourceName, "rule.0.ids_profiles.#", "1"),
					// Verify oversubscription is NOT present for Gateway rules
					resource.TestCheckNoResourceAttr(testResourceName, "rule.0.oversubscription"),
					// Second rule - simpler configuration
					resource.TestCheckResourceAttr(testResourceName, "rule.1.display_name", "simple-gw-rule"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.direction", "IN_OUT"),
					// Verify rule_id, nsx_id, and scope are present (from our latest alignment)
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.rule_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.1.rule_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.1.nsx_id"),
				),
			},
		},
	})
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyWithDetailedRules(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_tier1_gateway" "t1_gw" {
  display_name = "test-t1-gw-%s"
  description  = "Test primary Tier1 gateway"
}

resource "nsxt_policy_tier1_gateway" "t1_gw_secondary" {
  display_name = "test-t1-gw-secondary-%s"
  description  = "Test secondary Tier1 gateway"
}

resource "nsxt_policy_group" "external_clients" {
  display_name = "test-external-clients-%s"
  description  = "Test external clients group"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["203.0.113.0/24", "198.51.100.0/24"]
    }
  }
}

resource "nsxt_policy_group" "internal_servers" {
  display_name = "test-internal-servers-%s"
  description  = "Test internal servers group"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["10.10.10.0/24"]
    }
  }
}

resource "nsxt_policy_group" "web_tier" {
  display_name = "test-web-tier-%s"
  description  = "Test web tier group"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["172.16.100.0/24"]
    }
  }
}

data "nsxt_policy_service" "http" {
  display_name = "HTTP"
}

resource "nsxt_policy_intrusion_service_gateway_policy" "ids_gw_policy_detailed" {
  display_name    = "%s"
  description     = "Comprehensive IDPS Gateway policy test"
  category        = "LocalGatewayRules"
  sequence_number = 8
  locked          = false

  rule {
    display_name          = "comprehensive-gw-rule"
    description           = "Comprehensive Gateway rule testing all fields"
    notes                 = "Gateway rule with realistic configuration for north-south traffic"
    action                = "DETECT_PREVENT"
    direction             = "OUT"
    logged                = true
    disabled              = false
    sources_excluded      = true
    destinations_excluded = false
    ip_version            = "IPV6"
    sequence_number       = 15
    
    source_groups      = [nsxt_policy_group.external_clients.path]
    destination_groups = [nsxt_policy_group.internal_servers.path, nsxt_policy_group.web_tier.path]
    services          = [data.nsxt_policy_service.http.path]
    scope             = [nsxt_policy_tier1_gateway.t1_gw.path]
    ids_profiles      = [data.nsxt_policy_intrusion_service_profile.test.path]
    
    tag {
      scope = "environment"
      tag   = "production"
    }
    
    tag {
      scope = "traffic_type"
      tag   = "north_south"
    }
  }

  rule {
    display_name        = "simple-gw-rule"
    description         = "Simple Gateway rule for comparison"
    action              = "DETECT"
    direction           = "IN_OUT"
    logged              = false
    sequence_number     = 25
    scope              = [nsxt_policy_tier1_gateway.t1_gw.path]
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.test.path]
  }
}

data "nsxt_policy_intrusion_service_gateway_policy" "with_detailed_rules" {
  display_name = nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_detailed.display_name
  depends_on   = [nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_detailed]
}`, getAccTestResourceName(), getAccTestResourceName(), name, name, name, name)
}
