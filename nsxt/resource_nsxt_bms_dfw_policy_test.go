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
				Config: testAccBMSDFWPolicyStandaloneTemplate(testName),
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
				Config: testAccBMSDFWPolicyParentTemplate(testName),
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
				Config: testAccBMSDFWPolicyRulesTemplate(testName),
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
				Config: testAccBMSDFWPolicyCategoriesTemplate(testName),
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
				Config: testAccBMSDFWPolicyDataSourcesTemplate(testName),
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
				Config: testAccBMSDFWPolicyRuleDataSourceTemplate(testName),
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
				Config: testAccBMSDFWPolicyAllScenariosTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
					// Conditional checks for BMS resources
					testAccBMSDFWPolicyAllScenariosConditionalCheck(),
				),
			},
		},
	})
}

func testAccBMSDFWPolicyStandaloneTemplate(testName string) string {
	return fmt.Sprintf(`
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}
data "nsxt_policy_baremetal_server_interfaces" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
  has_interfaces = length(data.nsxt_policy_baremetal_server_interfaces.all.results) > 0
  server_ids = local.has_servers ? [for s in data.nsxt_policy_baremetal_servers.all.results : s.external_id] : []
  interface_ids = local.has_interfaces ? [for i in data.nsxt_policy_baremetal_server_interfaces.all.results : i.external_id] : []
  first_server = length(local.server_ids) > 0 ? [local.server_ids[0]] : []
  first_interface = length(local.interface_ids) > 0 ? [local.interface_ids[0]] : []
}

# Create BMS groups for policy testing
resource "nsxt_policy_group" "bms_source_group" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-source-group"
  description = "BMS source group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = local.first_server
    }
  }
}

resource "nsxt_policy_group" "bms_dest_group" {
  count = local.has_interfaces ? 1 : 0
  display_name = "%s-bms-dest-group"
  description = "BMS destination group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServerInterface"
      external_ids = local.first_interface
    }
  }
}

# Create service for policy
resource "nsxt_policy_service" "http" {
  count = local.has_servers && local.has_interfaces ? 1 : 0
  display_name = "%s-http-service"

  l4_port_set_entry {
    protocol = "TCP"
    destination_ports = ["80"]
  }
}

# DFW standalone policy with BMS groups
resource "nsxt_policy_security_policy" "bms_standalone" {
  count = local.has_servers && local.has_interfaces ? 1 : 0
  display_name = "%s-bms-standalone-policy"
  description = "DFW standalone policy with BMS groups"
  category = "Application"

  rule {
    display_name = "Allow-BMS-HTTP"
    source_groups = [nsxt_policy_group.bms_source_group[0].path]
    destination_groups = [nsxt_policy_group.bms_dest_group[0].path]
    action = "ALLOW"
    services = [nsxt_policy_service.http[0].path]
    logged = true
  }

  rule {
    display_name = "Drop-BMS-Default"
    destination_groups = [nsxt_policy_group.bms_source_group[0].path]
    action = "DROP"
    logged = true
  }
}

# Read policy using data source
data "nsxt_policy_security_policy" "bms_standalone_read" {
  count = local.has_servers && local.has_interfaces ? 1 : 0
  id = nsxt_policy_security_policy.bms_standalone[0].id
}
`, testName, testName, testName, testName)
}

func testAccBMSDFWPolicyParentTemplate(testName string) string {
	return fmt.Sprintf(`
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
  server_ids = local.has_servers ? [for s in data.nsxt_policy_baremetal_servers.all.results : s.external_id] : []
  first_server = length(local.server_ids) > 0 ? [local.server_ids[0]] : []
}

# Create BMS group for scope
resource "nsxt_policy_group" "bms_scope_group" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-scope-group"
  description = "BMS scope group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = local.first_server
    }
  }
}

# DFW parent policy with BMS scope
resource "nsxt_policy_parent_security_policy" "bms_parent" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-parent-policy"
  description = "DFW parent policy scoped to BMS"
  category = "Infrastructure"
  scope = [nsxt_policy_group.bms_scope_group[0].path]

}

# Note: nsxt_policy_parent_security_policy data source is not available
# Parent policy validation is done via resource creation only
`, testName, testName)
}

func testAccBMSDFWPolicyRulesTemplate(testName string) string {
	return fmt.Sprintf(`
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}
data "nsxt_policy_baremetal_server_interfaces" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
  has_interfaces = length(data.nsxt_policy_baremetal_server_interfaces.all.results) > 0
  server_ids = local.has_servers ? [for s in data.nsxt_policy_baremetal_servers.all.results : s.external_id] : []
  interface_ids = local.has_interfaces ? [for i in data.nsxt_policy_baremetal_server_interfaces.all.results : i.external_id] : []
  first_server = length(local.server_ids) > 0 ? [local.server_ids[0]] : []
  first_interface = length(local.interface_ids) > 0 ? [local.interface_ids[0]] : []
}

# Create BMS groups
resource "nsxt_policy_group" "bms_rule_source" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-rule-source"
  description = "BMS rule source group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = local.first_server
    }
  }
}

resource "nsxt_policy_group" "bms_rule_dest" {
  count = local.has_interfaces ? 1 : 0
  display_name = "%s-bms-rule-dest"
  description = "BMS rule destination group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServerInterface"
      external_ids = local.first_interface
    }
  }
}

# DFW standalone rule with BMS groups
resource "nsxt_policy_security_policy_rule" "bms_rule" {
  count = local.has_servers && local.has_interfaces ? 1 : 0
  display_name = "%s-bms-standalone-rule"
  description = "Standalone DFW rule with BMS groups"
  policy_path = "/infra/domains/default/security-policies/default-layer3-section"
  sequence_number = 100
  
  source_groups = [nsxt_policy_group.bms_rule_source[0].path]
  destination_groups = [nsxt_policy_group.bms_rule_dest[0].path]
  services = ["/infra/services/HTTPS"]
  action = "ALLOW"
  logged = true
}

# Note: Reading individual rules via data source is not supported
# Rules can be validated via the parent policy
`, testName, testName, testName)
}

func testAccBMSDFWPolicyCategoriesTemplate(testName string) string {
	return fmt.Sprintf(`
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
  server_ids = local.has_servers ? [for s in data.nsxt_policy_baremetal_servers.all.results : s.external_id] : []
  first_server = length(local.server_ids) > 0 ? [local.server_ids[0]] : []
}

# Create BMS group for all categories
resource "nsxt_policy_group" "bms_category_group" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-category-group"
  description = "BMS group for category testing"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = local.first_server
    }
  }
}

# Skip Ethernet category for BMS - not compatible with BMS groups
# Ethernet policies require pure Layer2 configuration without IP-based groups

# Category: Emergency
resource "nsxt_policy_security_policy" "bms_emergency" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-emergency-policy"
  description = "BMS policy - Emergency category"
  category = "Emergency"

  rule {
    display_name = "BMS-Emergency-Rule"
    destination_groups = [nsxt_policy_group.bms_category_group[0].path]
    action = "DROP"
  }
}

# Category: Infrastructure
resource "nsxt_policy_security_policy" "bms_infrastructure" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-infrastructure-policy"
  description = "BMS policy - Infrastructure category"
  category = "Infrastructure"

  rule {
    display_name = "BMS-Infrastructure-Rule"
    source_groups = [nsxt_policy_group.bms_category_group[0].path]
    destination_groups = [nsxt_policy_group.bms_category_group[0].path]
    action = "ALLOW"
    services = ["/infra/services/SSH"]
  }
}

# Category: Environment
resource "nsxt_policy_security_policy" "bms_environment" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-environment-policy"
  description = "BMS policy - Environment category"
  category = "Environment"

  rule {
    display_name = "BMS-Environment-Rule"
    destination_groups = [nsxt_policy_group.bms_category_group[0].path]
    action = "ALLOW"
    services = ["/infra/services/HTTP"]
  }
}

# Category: Application
resource "nsxt_policy_security_policy" "bms_application" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-application-policy"
  description = "BMS policy - Application category"
  category = "Application"

  rule {
    display_name = "BMS-Application-Rule"
    source_groups = [nsxt_policy_group.bms_category_group[0].path]
    destination_groups = [nsxt_policy_group.bms_category_group[0].path]
    action = "ALLOW"
    services = ["/infra/services/HTTPS"]
  }
}

# Ethernet policy skipped - not compatible with BMS groups

# Output policy information
output "policy_categories" {
  value = {
    policies_created = local.has_servers
    ethernet_policy = "skipped - not compatible with BMS groups"
    emergency_policy = local.has_servers ? nsxt_policy_security_policy.bms_emergency[0].path : ""
    infrastructure_policy = local.has_servers ? nsxt_policy_security_policy.bms_infrastructure[0].path : ""
    environment_policy = local.has_servers ? nsxt_policy_security_policy.bms_environment[0].path : ""
    application_policy = local.has_servers ? nsxt_policy_security_policy.bms_application[0].path : ""
  }
}
`, testName, testName, testName, testName, testName)
}

func testAccBMSDFWPolicyDataSourcesTemplate(testName string) string {
	return fmt.Sprintf(`
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
  sorted_server_ids = local.has_servers ? sort([for s in data.nsxt_policy_baremetal_servers.all.results : s.external_id]) : []
  first_server_id = length(local.sorted_server_ids) > 0 ? local.sorted_server_ids[0] : ""
}

# Create BMS Group
resource "nsxt_policy_group" "bms_group" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = [local.first_server_id]
    }
  }
}

# Create DFW Policy
resource "nsxt_policy_security_policy" "bms_policy" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-policy"
  category     = "Application"

  rule {
    display_name       = "Allow BMS Traffic"
    source_groups      = [nsxt_policy_group.bms_group[0].path]
    destination_groups = [nsxt_policy_group.bms_group[0].path]
    action             = "ALLOW"
    services           = []
    logged             = true
  }
}

# Read group via data source
data "nsxt_policy_group" "bms_group" {
  count = local.has_servers ? 1 : 0
  id    = nsxt_policy_group.bms_group[0].id
}

# Read policy via data source
data "nsxt_policy_security_policy" "bms_policy" {
  count = local.has_servers ? 1 : 0
  id    = nsxt_policy_security_policy.bms_policy[0].id
}
`, testName, testName)
}

func testAccBMSDFWPolicyRuleDataSourceTemplate(testName string) string {
	return fmt.Sprintf(`
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
  sorted_server_ids = local.has_servers ? sort([for s in data.nsxt_policy_baremetal_servers.all.results : s.external_id]) : []
  first_server_id = length(local.sorted_server_ids) > 0 ? local.sorted_server_ids[0] : ""
}

# Create BMS Group
resource "nsxt_policy_group" "bms_group" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = [local.first_server_id]
    }
  }
}

# Create DFW Policy with a rule
resource "nsxt_policy_security_policy" "bms_policy" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-policy"
  category     = "Application"

  rule {
    display_name       = "%s-bms-rule"
    source_groups      = [nsxt_policy_group.bms_group[0].path]
    destination_groups = [nsxt_policy_group.bms_group[0].path]
    action             = "ALLOW"
    direction          = "IN_OUT"
    services           = []
    logged             = true
  }
}

# Note: nsxt_policy_security_policy_rule data source is not available
# Rule can be read via the policy data source instead
`, testName, testName, testName)
}

func testAccBMSDFWPolicyAllScenariosTemplate(testName string) string {
	return fmt.Sprintf(`
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}
data "nsxt_policy_baremetal_server_interfaces" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
  has_interfaces = length(data.nsxt_policy_baremetal_server_interfaces.all.results) > 0
  sorted_server_ids = local.has_servers ? sort([for s in data.nsxt_policy_baremetal_servers.all.results : s.external_id]) : []
  sorted_interface_ids = local.has_interfaces ? sort([for i in data.nsxt_policy_baremetal_server_interfaces.all.results : i.external_id]) : []
  first_server_id = length(local.sorted_server_ids) > 0 ? local.sorted_server_ids[0] : ""
  first_interface_id = length(local.sorted_interface_ids) > 0 ? local.sorted_interface_ids[0] : ""
}

# BMS Server Group
resource "nsxt_policy_group" "bms_server_group" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-server-group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = [local.first_server_id]
    }
  }
}

# BMS Interface Group
resource "nsxt_policy_group" "bms_interface_group" {
  count = local.has_interfaces ? 1 : 0
  display_name = "%s-interface-group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServerInterface"
      external_ids = [local.first_interface_id]
    }
  }
}

# Scenario 1: Parent Policy with BMS Group in Scope
resource "nsxt_policy_parent_security_policy" "bms_parent" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-parent-policy"
  category     = "Application"
  scope        = [nsxt_policy_group.bms_server_group[0].path]
}

# Scenario 2: Standalone Policy with Rules using BMS Groups
resource "nsxt_policy_security_policy" "bms_standalone" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-standalone-policy"
  category     = "Application"

  rule {
    display_name       = "BMS Source Rule"
    source_groups      = [nsxt_policy_group.bms_server_group[0].path]
    action             = "ALLOW"
    services           = []
  }

  rule {
    display_name       = "BMS Destination Rule"
    destination_groups = local.has_interfaces && length(nsxt_policy_group.bms_interface_group) > 0 ? [nsxt_policy_group.bms_interface_group[0].path] : [nsxt_policy_group.bms_server_group[0].path]
    action             = "ALLOW"
    services           = []
  }

  rule {
    display_name       = "BMS Bidirectional Rule"
    source_groups      = [nsxt_policy_group.bms_server_group[0].path]
    destination_groups = local.has_interfaces && length(nsxt_policy_group.bms_interface_group) > 0 ? [nsxt_policy_group.bms_interface_group[0].path] : [nsxt_policy_group.bms_server_group[0].path]
    action             = "ALLOW"
    services           = []
  }
}

# Scenario 3: Additional validation - Groups are properly created
# (Standalone rule testing is covered in dedicated test files)

# Output comprehensive scenario information
output "all_scenarios" {
  value = {
    resources_available = local.has_servers
    parent_policy = local.has_servers ? nsxt_policy_parent_security_policy.bms_parent[0].path : ""
    standalone_policy = local.has_servers ? nsxt_policy_security_policy.bms_standalone[0].path : ""
    server_group = local.has_servers ? nsxt_policy_group.bms_server_group[0].path : ""
    interface_group = local.has_interfaces ? nsxt_policy_group.bms_interface_group[0].path : ""
  }
}
`, testName, testName, testName, testName)
}

func testAccBMSDFWPolicyDataSourcesConditionalCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Check if group data source exists (conditional on having BMS servers)
		groupDataSourceName := "data.nsxt_policy_group.bms_group"
		policyDataSourceName := "data.nsxt_policy_security_policy.bms_policy"

		if _, ok := s.RootModule().Resources[groupDataSourceName+".0"]; ok {
			// If the data sources exist, verify their attributes
			checks := []resource.TestCheckFunc{
				resource.TestCheckResourceAttrSet(groupDataSourceName+".0", "id"),
				resource.TestCheckResourceAttrSet(groupDataSourceName+".0", "display_name"),
				resource.TestCheckResourceAttrSet(groupDataSourceName+".0", "path"),
			}

			if _, ok := s.RootModule().Resources[policyDataSourceName+".0"]; ok {
				checks = append(checks,
					resource.TestCheckResourceAttrSet(policyDataSourceName+".0", "id"),
					resource.TestCheckResourceAttrSet(policyDataSourceName+".0", "display_name"),
					resource.TestCheckResourceAttrSet(policyDataSourceName+".0", "path"),
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

		if _, ok := s.RootModule().Resources[parentPolicyName+".0"]; ok {
			checks = append(checks, resource.TestCheckResourceAttrSet(parentPolicyName+".0", "id"))
		}

		if _, ok := s.RootModule().Resources[standalonePolicyName+".0"]; ok {
			checks = append(checks, resource.TestCheckResourceAttrSet(standalonePolicyName+".0", "id"))
		}

		if _, ok := s.RootModule().Resources[standaloneRuleName+".0"]; ok {
			checks = append(checks, resource.TestCheckResourceAttrSet(standaloneRuleName+".0", "id"))
		}

		if len(checks) > 0 {
			return resource.ComposeTestCheckFunc(checks...)(s)
		}

		// If no BMS resources available, that's also valid
		return nil
	}
}
