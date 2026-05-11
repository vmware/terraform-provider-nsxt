// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceNsxtPolicyIntrusionServiceGatewayPolicyRule_basic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_gateway_policy_rule.by_name"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyRuleReadByName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Gateway Rule"),
					resource.TestCheckResourceAttr(testResourceName, "action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "logged", "true"),
					resource.TestCheckResourceAttr(testResourceName, "disabled", "false"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "policy_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "sequence_number"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule_id"),
				),
			},
		},
	})
}

func TestAccDataSourceNsxtPolicyIntrusionServiceGatewayPolicyRule_byID(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_gateway_policy_rule.by_id"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyRuleReadByID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test Gateway Rule"),
					resource.TestCheckResourceAttr(testResourceName, "action", "DETECT"),
					resource.TestCheckResourceAttr(testResourceName, "logged", "true"),
					resource.TestCheckResourceAttr(testResourceName, "disabled", "false"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "policy_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "sequence_number"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule_id"),
				),
			},
		},
	})
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyRuleReadByName(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_tier1_gateway" "gw1" {
  display_name = "test-t1-gw-%s"
  description  = "Test Tier1 gateway for IDPS rule tests"
}

resource "nsxt_policy_group" "basic_source" {
  display_name = "test-basic-source-%s"
  description  = "Basic source group for simple test"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["10.10.10.0/24"]
    }
  }
}

resource "nsxt_policy_group" "basic_dest" {
  display_name = "test-basic-dest-%s"
  description  = "Basic destination group for simple test"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["10.10.20.0/24"]
    }
  }
}

resource "nsxt_policy_intrusion_service_gateway_policy" "ids_gw_policy_for_rule_ds" {
  display_name    = "%s-gw-policy"
  description     = "Acceptance Test Gateway Policy"
  category        = "LocalGatewayRules"
  sequence_number = 3
  
  rule {
    display_name    = "%s"
    description     = "Acceptance Test Gateway Rule"
    action          = "DETECT"
    logged          = true
    disabled        = false
    scope           = [nsxt_policy_tier1_gateway.gw1.path]
    sequence_number = 1
    ids_profiles    = [data.nsxt_policy_intrusion_service_profile.test.path]
  }
}

data "nsxt_policy_intrusion_service_gateway_policy_rule" "by_name" {
  display_name = "%s"
  policy_path  = nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_rule_ds.path
  depends_on   = [nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_rule_ds]
}`, getAccTestResourceName(), name, name, name, name, name)
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyRuleReadByID(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_tier1_gateway" "gw1" {
  display_name = "test-t1-gw-%s"
  description  = "Test Tier1 gateway for IDPS rule tests"
}

resource "nsxt_policy_group" "basic_source" {
  display_name = "test-basic-source-%s"
  description  = "Basic source group for simple test"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["10.10.10.0/24"]
    }
  }
}

resource "nsxt_policy_group" "basic_dest" {
  display_name = "test-basic-dest-%s"
  description  = "Basic destination group for simple test"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["10.10.20.0/24"]
    }
  }
}

resource "nsxt_policy_intrusion_service_gateway_policy" "ids_gw_policy_for_rule_ds" {
  display_name    = "%s-gw-policy"
  description     = "Acceptance Test Gateway Policy"
  category        = "LocalGatewayRules"
  sequence_number = 3
  
  rule {
    display_name       = "%s"
    description        = "Acceptance Test Gateway Rule"
    action             = "DETECT"
    logged             = true
    disabled           = false
    source_groups      = [nsxt_policy_group.basic_source.path]
    destination_groups = [nsxt_policy_group.basic_dest.path]
    services           = []
    scope              = [nsxt_policy_tier1_gateway.gw1.path]
    sequence_number    = 1
    ids_profiles       = [data.nsxt_policy_intrusion_service_profile.test.path]
  }
}

data "nsxt_policy_intrusion_service_gateway_policy_rule" "by_id" {
  display_name = "%s"
  policy_path  = nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_rule_ds.path
  depends_on   = [nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_rule_ds]
}`, getAccTestResourceName(), name, name, name, name, name)
}

func TestAccDataSourceNsxtPolicyIntrusionServiceGatewayPolicyRule_withDetailedConfiguration(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "data.nsxt_policy_intrusion_service_gateway_policy_rule.detailed_gw_rule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyRuleWithDetailedConfiguration(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Comprehensive standalone Gateway IDPS rule test"),
					resource.TestCheckResourceAttr(testResourceName, "action", "DETECT_PREVENT"),
					resource.TestCheckResourceAttr(testResourceName, "direction", "OUT"),
					resource.TestCheckResourceAttr(testResourceName, "logged", "true"),
					resource.TestCheckResourceAttr(testResourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(testResourceName, "sources_excluded", "false"),
					resource.TestCheckResourceAttr(testResourceName, "destinations_excluded", "true"),
					resource.TestCheckResourceAttr(testResourceName, "ip_version", "IPV6"),
					resource.TestCheckResourceAttrSet(testResourceName, "sequence_number"),
					resource.TestCheckResourceAttrSet(testResourceName, "policy_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					// Validate detailed Gateway rule configuration
					resource.TestCheckResourceAttr(testResourceName, "source_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "destination_groups.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "services.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "scope.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ids_profiles.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "2"),
					// Verify oversubscription is NOT present for Gateway rules
					resource.TestCheckNoResourceAttr(testResourceName, "oversubscription"),
					// Verify rule_id is present
					resource.TestCheckResourceAttrSet(testResourceName, "rule_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyRuleWithDetailedConfiguration(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_tier1_gateway" "primary_gw" {
  display_name = "test-primary-gw-%s"
  description  = "Test primary gateway for detailed IDPS rules"
}


resource "nsxt_policy_group" "internet_clients" {
  display_name = "test-internet-clients-%s"
  description  = "Test internet clients group"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["203.0.113.0/24", "198.51.100.0/24"]
    }
  }
}

resource "nsxt_policy_group" "public_web" {
  display_name = "test-public-web-%s"
  description  = "Test public web servers"
  
  criteria {
    ipaddress_expression {
      ip_addresses = ["192.168.50.10-192.168.50.20"]
    }
  }
}

data "nsxt_policy_service" "https_service" {
  display_name = "HTTPS"
}

resource "nsxt_policy_service" "api_service" {
  display_name = "test-api-service-%s"
  description  = "Test API service"

  l4_port_set_entry {
    display_name      = "api-port"
    description       = "API service on port 8443"
    protocol          = "TCP"
    destination_ports = ["8443"]
  }
}

resource "nsxt_policy_intrusion_service_gateway_policy" "ids_gw_policy_for_detailed_rule" {
  display_name    = "%s-gw-policy"
  description     = "Parent Gateway policy for detailed rule test"
  category        = "LocalGatewayRules"
  sequence_number = 9
  
  rule {
    display_name          = "%s"
    description           = "Comprehensive standalone Gateway IDPS rule test"
    notes                 = "Testing Gateway-specific IDPS rule fields with north-south traffic patterns"
    action                = "DETECT_PREVENT"
    direction             = "OUT"
    logged                = true
    log_label             = "gateway-security-audit"
    disabled              = false
    sources_excluded      = false
    destinations_excluded = true
    ip_version            = "IPV6"
    sequence_number       = 35
    
    source_groups      = [nsxt_policy_group.internet_clients.path]
    destination_groups = [nsxt_policy_group.public_web.path]
    services          = [data.nsxt_policy_service.https_service.path, nsxt_policy_service.api_service.path]
    scope             = [nsxt_policy_tier1_gateway.primary_gw.path]
    ids_profiles      = [data.nsxt_policy_intrusion_service_profile.test.path]
    
    tag {
      scope = "traffic_direction"
      tag   = "egress"
    }
    
    tag {
      scope = "policy_type"  
      tag   = "gateway_idps"
    }
  }
}

data "nsxt_policy_intrusion_service_gateway_policy_rule" "detailed_gw_rule" {
  display_name = "%s"
  policy_path  = nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_detailed_rule.path
  depends_on   = [nsxt_policy_intrusion_service_gateway_policy.ids_gw_policy_for_detailed_rule]
}`, getAccTestResourceName(), name, name, name, name, name, name)
}
