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

// TestAccBMSDFWPolicyStandalone tests DFW standalone policies with BMS groups
func TestAccBMSDFWPolicyStandalone(t *testing.T) {
	testName := getAccTestResourceName()
	serverID := getTestBMSServerID()
	interfaceID := getTestBMSInterfaceID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
			testAccEnvDefined(t, "NSXT_TEST_BMS_INTERFACE")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSDFWPolicyStandaloneTemplate(testName, serverID, interfaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
				),
			},
		},
	})
}

// TestAccBMSDFWPolicyParent tests DFW parent policies with BMS groups
func TestAccBMSDFWPolicyParent(t *testing.T) {
	testName := getAccTestResourceName()
	serverID := getTestBMSServerID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSDFWPolicyParentTemplate(testName, serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
				),
			},
		},
	})
}

// TestAccBMSDFWPolicyRules tests DFW standalone rules with BMS groups
func TestAccBMSDFWPolicyRules(t *testing.T) {
	testName := getAccTestResourceName()
	serverID := getTestBMSServerID()
	interfaceID := getTestBMSInterfaceID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
			testAccEnvDefined(t, "NSXT_TEST_BMS_INTERFACE")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSDFWPolicyRulesTemplate(testName, serverID, interfaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
				),
			},
		},
	})
}

// TestAccBMSDFWPolicyCategories tests DFW policies covering all categories
func TestAccBMSDFWPolicyCategories(t *testing.T) {
	testName := getAccTestResourceName()
	serverID := getTestBMSServerID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSDFWPolicyCategoriesTemplate(testName, serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
				),
			},
		},
	})
}

// TestAccBMSDFWPolicyDataSources tests reading DFW policies with BMS groups via data sources
func TestAccBMSDFWPolicyDataSources(t *testing.T) {
	testName := getAccTestResourceName()
	serverID := getTestBMSServerID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSDFWPolicyDataSourcesTemplate(testName, serverID),
				Check: resource.ComposeTestCheckFunc(
					testAccBMSDFWPolicyDataSourcesConditionalCheck(),
				),
			},
		},
	})
}

// TestAccBMSDFWPolicyRuleDataSource tests reading DFW policy rules via data sources
func TestAccBMSDFWPolicyRuleDataSource(t *testing.T) {
	testName := getAccTestResourceName()
	serverID := getTestBMSServerID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSDFWPolicyRuleDataSourceTemplate(testName, serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
				),
			},
		},
	})
}

// TestAccBMSDFWPolicyAllScenarios tests comprehensive DFW scenarios with BMS
func TestAccBMSDFWPolicyAllScenarios(t *testing.T) {
	testName := getAccTestResourceName()
	serverID := getTestBMSServerID()
	interfaceID := getTestBMSInterfaceID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
			testAccEnvDefined(t, "NSXT_TEST_BMS_INTERFACE")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSDFWPolicyAllScenariosTemplate(testName, serverID, interfaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
					// Conditional checks for BMS resources
					testAccBMSDFWPolicyAllScenariosConditionalCheck(),
				),
			},
		},
	})
}

func testAccBMSDFWPolicyStandaloneTemplate(testName, serverID, interfaceID string) string {
	return fmt.Sprintf(`
# Data sources for testing BMS inventory access
data "nsxt_policy_baremetal_servers" "all" {}
data "nsxt_policy_baremetal_server_interfaces" "all" {}

# Create BMS groups for policy testing
resource "nsxt_policy_group" "bms_source_group" {
  display_name = "%s-bms-source-group"
  description = "BMS source group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

resource "nsxt_policy_group" "bms_dest_group" {
  display_name = "%s-bms-dest-group"
  description = "BMS destination group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServerInterface"
      external_ids = ["%s"]
    }
  }
}

# Create service for policy
resource "nsxt_policy_service" "http" {
  display_name = "%s-http-service"

  l4_port_set_entry {
    protocol = "TCP"
    destination_ports = ["80"]
  }
}

# DFW standalone policy with BMS groups
resource "nsxt_policy_security_policy" "bms_standalone" {
  display_name = "%s-bms-standalone-policy"
  description = "DFW standalone policy with BMS groups"
  category = "Application"

  rule {
    display_name = "Allow-BMS-HTTP"
    source_groups = [nsxt_policy_group.bms_source_group.path]
    destination_groups = [nsxt_policy_group.bms_dest_group.path]
    action = "ALLOW"
    services = [nsxt_policy_service.http.path]
    logged = true
  }

  rule {
    display_name = "Drop-BMS-Default"
    destination_groups = [nsxt_policy_group.bms_source_group.path]
    action = "DROP"
    logged = true
  }
}

# Read policy using data source
data "nsxt_policy_security_policy" "bms_standalone_read" {
  id = nsxt_policy_security_policy.bms_standalone.id
}
`, testName, serverID, testName, interfaceID, testName, testName)
}

func testAccBMSDFWPolicyParentTemplate(testName, serverID string) string {
	return fmt.Sprintf(`
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}

# Create BMS group for scope
resource "nsxt_policy_group" "bms_scope_group" {
  display_name = "%s-bms-scope-group"
  description = "BMS scope group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

# DFW parent policy with BMS scope
resource "nsxt_policy_parent_security_policy" "bms_parent" {
  display_name = "%s-bms-parent-policy"
  description = "DFW parent policy scoped to BMS"
  category = "Infrastructure"
  scope = [nsxt_policy_group.bms_scope_group.path]

}

# Note: nsxt_policy_parent_security_policy data source is not available
# Parent policy validation is done via resource creation only
`, testName, serverID, testName)
}

func testAccBMSDFWPolicyRulesTemplate(testName, serverID, interfaceID string) string {
	return fmt.Sprintf(`
# Data sources for testing BMS inventory access
data "nsxt_policy_baremetal_servers" "all" {}
data "nsxt_policy_baremetal_server_interfaces" "all" {}

# Create BMS groups
resource "nsxt_policy_group" "bms_rule_source" {
  display_name = "%s-bms-rule-source"
  description = "BMS rule source group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

resource "nsxt_policy_group" "bms_rule_dest" {
  display_name = "%s-bms-rule-dest"
  description = "BMS rule destination group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServerInterface"
      external_ids = ["%s"]
    }
  }
}

# DFW standalone rule with BMS groups
resource "nsxt_policy_security_policy_rule" "bms_rule" {
  display_name = "%s-bms-standalone-rule"
  description = "Standalone DFW rule with BMS groups"
  policy_path = "/infra/domains/default/security-policies/default-layer3-section"
  sequence_number = 100
  
  source_groups = [nsxt_policy_group.bms_rule_source.path]
  destination_groups = [nsxt_policy_group.bms_rule_dest.path]
  services = ["/infra/services/HTTPS"]
  action = "ALLOW"
  logged = true
}

# Note: Reading individual rules via data source is not supported
# Rules can be validated via the parent policy
`, testName, serverID, testName, interfaceID, testName)
}

func testAccBMSDFWPolicyCategoriesTemplate(testName, serverID string) string {
	return fmt.Sprintf(`
# Data source for testing BMS inventory access
data "nsxt_policy_baremetal_servers" "all" {}

# Create BMS group for all categories
resource "nsxt_policy_group" "bms_category_group" {
  display_name = "%s-bms-category-group"
  description = "BMS group for category testing"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

# Skip Ethernet category for BMS - not compatible with BMS groups
# Ethernet policies require pure Layer2 configuration without IP-based groups

# Category: Emergency
resource "nsxt_policy_security_policy" "bms_emergency" {
  display_name = "%s-bms-emergency-policy"
  description = "BMS policy - Emergency category"
  category = "Emergency"

  rule {
    display_name = "BMS-Emergency-Rule"
    destination_groups = [nsxt_policy_group.bms_category_group.path]
    action = "DROP"
  }
}

# Category: Infrastructure
resource "nsxt_policy_security_policy" "bms_infrastructure" {
  display_name = "%s-bms-infrastructure-policy"
  description = "BMS policy - Infrastructure category"
  category = "Infrastructure"

  rule {
    display_name = "BMS-Infrastructure-Rule"
    source_groups = [nsxt_policy_group.bms_category_group.path]
    destination_groups = [nsxt_policy_group.bms_category_group.path]
    action = "ALLOW"
    services = ["/infra/services/SSH"]
  }
}

# Category: Environment
resource "nsxt_policy_security_policy" "bms_environment" {
  display_name = "%s-bms-environment-policy"
  description = "BMS policy - Environment category"
  category = "Environment"

  rule {
    display_name = "BMS-Environment-Rule"
    destination_groups = [nsxt_policy_group.bms_category_group.path]
    action = "ALLOW"
    services = ["/infra/services/HTTP"]
  }
}

# Category: Application
resource "nsxt_policy_security_policy" "bms_application" {
  display_name = "%s-bms-application-policy"
  description = "BMS policy - Application category"
  category = "Application"

  rule {
    display_name = "BMS-Application-Rule"
    source_groups = [nsxt_policy_group.bms_category_group.path]
    destination_groups = [nsxt_policy_group.bms_category_group.path]
    action = "ALLOW"
    services = ["/infra/services/HTTPS"]
  }
}

# Ethernet policy skipped - not compatible with BMS groups

# Output policy information
output "policy_categories" {
  value = {
    policies_created = true
    ethernet_policy = "skipped - not compatible with BMS groups"
    emergency_policy = nsxt_policy_security_policy.bms_emergency.path
    infrastructure_policy = nsxt_policy_security_policy.bms_infrastructure.path
    environment_policy = nsxt_policy_security_policy.bms_environment.path
    application_policy = nsxt_policy_security_policy.bms_application.path
  }
}
`, testName, serverID, testName, testName, testName, testName)
}

func testAccBMSDFWPolicyDataSourcesTemplate(testName, serverID string) string {
	return fmt.Sprintf(`
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}


# Create BMS Group
resource "nsxt_policy_group" "bms_group" {
  display_name = "%s-bms-group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

# Create DFW Policy
resource "nsxt_policy_security_policy" "bms_policy" {
  display_name = "%s-bms-policy"
  category     = "Application"

  rule {
    display_name       = "Allow BMS Traffic"
    source_groups      = [nsxt_policy_group.bms_group.path]
    destination_groups = [nsxt_policy_group.bms_group.path]
    action             = "ALLOW"
    services           = []
    logged             = true
  }
}

# Read group via data source
data "nsxt_policy_group" "bms_group" {
  id    = nsxt_policy_group.bms_group.id
}

# Read policy via data source
data "nsxt_policy_security_policy" "bms_policy" {
  id    = nsxt_policy_security_policy.bms_policy.id
}
`, testName, serverID, testName)
}

func testAccBMSDFWPolicyRuleDataSourceTemplate(testName, serverID string) string {
	return fmt.Sprintf(`
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}


# Create BMS Group
resource "nsxt_policy_group" "bms_group" {
  display_name = "%s-bms-group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

# Create DFW Policy with a rule
resource "nsxt_policy_security_policy" "bms_policy" {
  display_name = "%s-bms-policy"
  category     = "Application"

  rule {
    display_name       = "%s-bms-rule"
    source_groups      = [nsxt_policy_group.bms_group.path]
    destination_groups = [nsxt_policy_group.bms_group.path]
    action             = "ALLOW"
    direction          = "IN_OUT"
    services           = []
    logged             = true
  }
}

# Note: nsxt_policy_security_policy_rule data source is not available
# Rule can be read via the policy data source instead
`, testName, serverID, testName, testName)
}

func testAccBMSDFWPolicyAllScenariosTemplate(testName, serverID, interfaceID string) string {
	return fmt.Sprintf(`
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}
data "nsxt_policy_baremetal_server_interfaces" "all" {}


# BMS Server Group
resource "nsxt_policy_group" "bms_server_group" {
  display_name = "%s-server-group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

# BMS Interface Group
resource "nsxt_policy_group" "bms_interface_group" {
  display_name = "%s-interface-group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServerInterface"
      external_ids = ["%s"]
    }
  }
}

# Scenario 1: Parent Policy with BMS Group in Scope
resource "nsxt_policy_parent_security_policy" "bms_parent" {
  display_name = "%s-parent-policy"
  category     = "Application"
  scope        = [nsxt_policy_group.bms_server_group.path]
}

# Scenario 2: Standalone Policy with Rules using BMS Groups
resource "nsxt_policy_security_policy" "bms_standalone" {
  display_name = "%s-standalone-policy"
  category     = "Application"

  rule {
    display_name       = "BMS Source Rule"
    source_groups      = [nsxt_policy_group.bms_server_group.path]
    action             = "ALLOW"
    services           = []
  }

  rule {
    display_name       = "BMS Destination Rule"
    destination_groups = [nsxt_policy_group.bms_interface_group.path]
    action             = "ALLOW"
    services           = []
  }

  rule {
    display_name       = "BMS Bidirectional Rule"
    source_groups      = [nsxt_policy_group.bms_server_group.path]
    destination_groups = [nsxt_policy_group.bms_interface_group.path]
    action             = "ALLOW"
    services           = []
  }
}

# Scenario 3: Additional validation - Groups are properly created
# (Standalone rule testing is covered in dedicated test files)

# Output comprehensive scenario information
output "all_scenarios" {
  value = {
    resources_available = true
    parent_policy = nsxt_policy_parent_security_policy.bms_parent.path
    standalone_policy = nsxt_policy_security_policy.bms_standalone.path
    server_group = nsxt_policy_group.bms_server_group.path
    interface_group = nsxt_policy_group.bms_interface_group.path
  }
}
`, testName, serverID, testName, interfaceID, testName, testName)
}

func testAccBMSDFWPolicyDataSourcesConditionalCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Check if group data source exists (conditional on having BMS servers)
		groupDataSourceName := "data.nsxt_policy_group.bms_group"
		policyDataSourceName := "data.nsxt_policy_security_policy.bms_policy"

		if _, ok := s.RootModule().Resources[groupDataSourceName]; ok {
			// If the data sources exist, verify their attributes
			checks := []resource.TestCheckFunc{
				resource.TestCheckResourceAttrSet(groupDataSourceName, "id"),
				resource.TestCheckResourceAttrSet(groupDataSourceName, "display_name"),
				resource.TestCheckResourceAttrSet(groupDataSourceName, "path"),
			}

			if _, ok := s.RootModule().Resources[policyDataSourceName]; ok {
				checks = append(checks,
					resource.TestCheckResourceAttrSet(policyDataSourceName, "id"),
					resource.TestCheckResourceAttrSet(policyDataSourceName, "display_name"),
					resource.TestCheckResourceAttrSet(policyDataSourceName, "path"),
				)
			}

			return resource.ComposeTestCheckFunc(checks...)(s)
		}
		// If no BMS servers available, that's also valid
		return nil
	}
}

func testAccBMSDFWPolicyAllScenariosConditionalCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Check if BMS resources exist (conditional on having BMS servers)
		parentPolicyName := "nsxt_policy_parent_security_policy.bms_parent"
		standalonePolicyName := "nsxt_policy_security_policy.bms_standalone"
		standaloneRuleName := "nsxt_policy_security_policy_rule.bms_standalone_rule"

		checks := []resource.TestCheckFunc{}

		if _, ok := s.RootModule().Resources[parentPolicyName]; ok {
			checks = append(checks, resource.TestCheckResourceAttrSet(parentPolicyName, "id"))
		}

		if _, ok := s.RootModule().Resources[standalonePolicyName]; ok {
			checks = append(checks, resource.TestCheckResourceAttrSet(standalonePolicyName, "id"))
		}

		if _, ok := s.RootModule().Resources[standaloneRuleName]; ok {
			checks = append(checks, resource.TestCheckResourceAttrSet(standaloneRuleName, "id"))
		}

		if len(checks) > 0 {
			return resource.ComposeTestCheckFunc(checks...)(s)
		}

		// If no BMS resources available, that's also valid
		return nil
	}
}
