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

// TestAccBMSEndToEndWorkflow tests complete BMS workflow from inventory to policy
func TestAccBMSEndToEndWorkflow(t *testing.T) {
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
				// Step 1: Discovery and basic setup
				Config: testAccBMSEndToEndStep1Template(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_interfaces.all", "results.#"),
				),
			},
			{
				// Step 2: Add tagging
				Config: testAccBMSEndToEndStep2Template(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
				),
			},
			{
				// Step 3: Create groups and policies
				Config: testAccBMSEndToEndStep3Template(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
				),
			},
		},
	})
}

// TestAccBMSEndToEndMembersAssociations tests BMS members and associations
func TestAccBMSEndToEndMembersAssociations(t *testing.T) {
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
				Config: testAccBMSEndToEndMembersAssociationsTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
				),
			},
		},
	})
}

// TestAccBMSEndToEndDataSourcesCoverage tests comprehensive data source coverage
func TestAccBMSEndToEndDataSourcesCoverage(t *testing.T) {
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
				Config: testAccBMSEndToEndDataSourcesCoverageTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					// Inventory data sources
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_interfaces.all", "results.#"),
					// Conditional group data source check
					testAccBMSEndToEndDataSourcesConditionalCheck(),
				),
			},
		},
	})
}

func testAccBMSEndToEndStep1Template(testName string) string {
	return `
# Step 1: Discovery and inventory
data "nsxt_policy_baremetal_servers" "all" {}
data "nsxt_policy_baremetal_server_interfaces" "all" {}

# Basic filtering
data "nsxt_policy_baremetal_servers" "linux_servers" {
  os_name = "linux"
}

# Individual lookups
locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
  has_interfaces = length(data.nsxt_policy_baremetal_server_interfaces.all.results) > 0
  first_server_id = local.has_servers ? sort([for s in data.nsxt_policy_baremetal_servers.all.results : s.external_id])[0] : ""
  first_interface_id = local.has_interfaces ? sort([for i in data.nsxt_policy_baremetal_server_interfaces.all.results : i.external_id])[0] : ""
}

data "nsxt_policy_baremetal_server" "first" {
  count = local.has_servers ? 1 : 0
  external_id = local.first_server_id
}

data "nsxt_policy_baremetal_server_interface" "first" {
  count = local.has_interfaces ? 1 : 0
  external_id = local.first_interface_id
}

# Output discovery results
output "step1_discovery" {
  value = {
    total_servers = length(data.nsxt_policy_baremetal_servers.all.results)
    total_interfaces = length(data.nsxt_policy_baremetal_server_interfaces.all.results)
    linux_servers = length(data.nsxt_policy_baremetal_servers.linux_servers.results)
    first_server_id = local.first_server_id
    first_interface_id = local.first_interface_id
  }
}
`
}

func testAccBMSEndToEndStep2Template(testName string) string {
	return `
# Step 2: Discovery + Tagging
data "nsxt_policy_baremetal_servers" "all" {}
data "nsxt_policy_baremetal_server_interfaces" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
  has_interfaces = length(data.nsxt_policy_baremetal_server_interfaces.all.results) > 0
  first_server_id = local.has_servers ? sort([for s in data.nsxt_policy_baremetal_servers.all.results : s.external_id])[0] : ""
  first_interface_id = local.has_interfaces ? sort([for i in data.nsxt_policy_baremetal_server_interfaces.all.results : i.external_id])[0] : ""
}

# Add tags to server
resource "nsxt_policy_baremetal_server_tags" "web_server" {
  count = local.has_servers ? 1 : 0
  external_id = local.first_server_id

  tag {
    scope = "test-env"
    tag   = "production"
  }

  tag {
    scope = "test-app"
    tag   = "web-server"
  }

  tag {
    scope = "test-tier"
    tag   = "frontend"
  }
}

# Add tags to interface
resource "nsxt_policy_baremetal_server_interface_tags" "mgmt_interface" {
  count = local.has_interfaces ? 1 : 0
  external_id = local.first_interface_id

  tag {
    scope = "test-net"
    tag   = "management"
  }

  tag {
    scope = "test-vlan"
    tag   = "100"
  }
}

# Read tags back
data "nsxt_policy_baremetal_server_tags" "web_server_tags" {
  count = local.has_servers ? 1 : 0
  external_id = local.first_server_id
  depends_on = [nsxt_policy_baremetal_server_tags.web_server]
}

data "nsxt_policy_baremetal_server_interface_tags" "mgmt_interface_tags" {
  count = local.has_interfaces ? 1 : 0
  external_id = local.first_interface_id
  depends_on = [nsxt_policy_baremetal_server_interface_tags.mgmt_interface]
}

# Output tagging results
output "step2_tagging" {
  value = {
    server_tagged = local.has_servers
    interface_tagged = local.has_interfaces
    server_tags = local.has_servers ? nsxt_policy_baremetal_server_tags.web_server[0].tag : []
    interface_tags = local.has_interfaces ? nsxt_policy_baremetal_server_interface_tags.mgmt_interface[0].tag : []
  }
}
`
}

func testAccBMSEndToEndStep3Template(testName string) string {
	return fmt.Sprintf(`
# Step 3: Discovery + Tagging + Groups + Policies
data "nsxt_policy_baremetal_servers" "all" {}
data "nsxt_policy_baremetal_server_interfaces" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
  has_interfaces = length(data.nsxt_policy_baremetal_server_interfaces.all.results) > 0
  first_server_id = local.has_servers ? sort([for s in data.nsxt_policy_baremetal_servers.all.results : s.external_id])[0] : ""
  first_interface_id = local.has_interfaces ? sort([for i in data.nsxt_policy_baremetal_server_interfaces.all.results : i.external_id])[0] : ""
  can_create_policy = local.has_servers && local.has_interfaces
}

# Keep the tagging from Step 2
resource "nsxt_policy_baremetal_server_tags" "web_server" {
  count = local.has_servers ? 1 : 0
  external_id = local.first_server_id

  tag {
    scope = "test-env"
    tag   = "production"
  }

  tag {
    scope = "test-app"
    tag   = "web-server"
  }

  tag {
    scope = "test-tier"
    tag   = "frontend"
  }
}

resource "nsxt_policy_baremetal_server_interface_tags" "mgmt_interface" {
  count = local.has_interfaces ? 1 : 0
  external_id = local.first_interface_id

  tag {
    scope = "test-net"
    tag   = "management"
  }

  tag {
    scope = "test-vlan"
    tag   = "100"
  }
}

# Create static BMS group
resource "nsxt_policy_group" "web_servers" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-web-servers"
  description = "Web server BMS group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = [local.first_server_id]
    }
  }
}

# Create dynamic BMS group based on tags
resource "nsxt_policy_group" "production_servers" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-production-servers"
  description = "Production BMS servers (dynamic)"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "Tag"
      member_type = "BareMetalServer"
      operator = "EQUALS"
      value = "test-env-%s|production"
    }
  }
}

# Create BMS interface group
resource "nsxt_policy_group" "mgmt_interfaces" {
  count = local.has_interfaces ? 1 : 0
  display_name = "%s-mgmt-interfaces"
  description = "Management interfaces group"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "Tag"
      member_type = "BareMetalServerInterface"
      operator = "EQUALS"
      value = "test-net-%s|management"
    }
  }
}

# Create service for policy
resource "nsxt_policy_service" "web_service" {
  count = local.can_create_policy ? 1 : 0
  display_name = "%s-web-service"

  l4_port_set_entry {
    protocol = "TCP"
    destination_ports = ["80", "443"]
  }
}

# Create DFW policy with BMS groups
resource "nsxt_policy_security_policy" "web_policy" {
  count = local.can_create_policy ? 1 : 0
  display_name = "%s-web-policy"
  description = "Security policy for web servers"
  category = "Application"

  rule {
    display_name = "Allow-Web-Traffic"
    source_groups = [nsxt_policy_group.production_servers[0].path]
    destination_groups = [nsxt_policy_group.web_servers[0].path]
    action = "ALLOW"
    services = [nsxt_policy_service.web_service[0].path]
    logged = true
  }

  rule {
    display_name = "Allow-Management"
    destination_groups = [nsxt_policy_group.mgmt_interfaces[0].path]
    action = "ALLOW"
    services = ["/infra/services/SSH"]
    logged = true
  }

  rule {
    display_name = "Drop-Default"
    destination_groups = [nsxt_policy_group.web_servers[0].path]
    action = "DROP"
    logged = true
  }
}

# Read groups back using data sources
data "nsxt_policy_group" "web_servers_read" {
  count = local.has_servers ? 1 : 0
  id = nsxt_policy_group.web_servers[0].id
}

data "nsxt_policy_group" "production_servers_read" {
  count = local.has_servers ? 1 : 0
  id = nsxt_policy_group.production_servers[0].id
}

# Read policy back using data source
data "nsxt_policy_security_policy" "web_policy_read" {
  count = local.can_create_policy ? 1 : 0
  id = nsxt_policy_security_policy.web_policy[0].id
}

# Output complete workflow results
output "step3_complete_workflow" {
  value = {
    servers_available = local.has_servers
    interfaces_available = local.has_interfaces
    policy_created = local.can_create_policy
    web_servers_group = local.has_servers ? nsxt_policy_group.web_servers[0].path : ""
    production_servers_group = local.has_servers ? nsxt_policy_group.production_servers[0].path : ""
    mgmt_interfaces_group = local.has_interfaces ? nsxt_policy_group.mgmt_interfaces[0].path : ""
    web_policy = local.can_create_policy ? nsxt_policy_security_policy.web_policy[0].path : ""
  }
}
`, testName, testName, testName, testName, testName, testName, testName)
}

func testAccBMSEndToEndMembersAssociationsTemplate(testName string) string {
	return fmt.Sprintf(`
# Test BMS members and associations
data "nsxt_policy_baremetal_servers" "all" {}
data "nsxt_policy_baremetal_server_interfaces" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
  has_interfaces = length(data.nsxt_policy_baremetal_server_interfaces.all.results) > 0
  server_ids = local.has_servers ? [for s in data.nsxt_policy_baremetal_servers.all.results : s.external_id] : []
  interface_ids = local.has_interfaces ? [for i in data.nsxt_policy_baremetal_server_interfaces.all.results : i.external_id] : []
  selected_servers = length(local.server_ids) >= 2 ? slice(local.server_ids, 0, 2) : local.server_ids
  selected_interfaces = length(local.interface_ids) >= 2 ? slice(local.interface_ids, 0, 2) : local.interface_ids
}

# Create multiple BMS groups to test associations
resource "nsxt_policy_group" "bms_group1" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-group1"
  description = "First BMS group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = length(local.selected_servers) > 0 ? [local.selected_servers[0]] : []
    }
  }
}

resource "nsxt_policy_group" "bms_group2" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-group2"
  description = "Second BMS group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = length(local.selected_servers) > 1 ? [local.selected_servers[1]] : length(local.selected_servers) > 0 ? [local.selected_servers[0]] : []
    }
  }
}

resource "nsxt_policy_group" "bmsi_group1" {
  count = local.has_interfaces ? 1 : 0
  display_name = "%s-bmsi-group1"
  description = "First BMS interface group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServerInterface"
      external_ids = length(local.selected_interfaces) > 0 ? [local.selected_interfaces[0]] : []
    }
  }
}

resource "nsxt_policy_group" "bmsi_group2" {
  count = local.has_interfaces ? 1 : 0
  display_name = "%s-bmsi-group2"
  description = "Second BMS interface group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServerInterface"
      external_ids = length(local.selected_interfaces) > 1 ? [local.selected_interfaces[1]] : length(local.selected_interfaces) > 0 ? [local.selected_interfaces[0]] : []
    }
  }
}

# Create nested group that includes other BMS groups
resource "nsxt_policy_group" "nested_bms_group" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-nested-bms-group"
  description = "Nested group containing BMS groups"
  group_type = "BareMetalServer"

  criteria {
    path_expression {
      member_paths = [
        nsxt_policy_group.bms_group1[0].path,
        nsxt_policy_group.bms_group2[0].path
      ]
    }
  }
}

# Test group associations by creating policy using all groups
resource "nsxt_policy_security_policy" "associations_policy" {
  count = local.has_servers && local.has_interfaces ? 1 : 0
  display_name = "%s-associations-policy"
  description = "Policy testing BMS associations"
  category = "Application"

  rule {
    display_name = "BMS-Server-to-Interface"
    source_groups = [nsxt_policy_group.bms_group1[0].path]
    destination_groups = [nsxt_policy_group.bmsi_group1[0].path]
    action = "ALLOW"
    services = ["/infra/services/HTTP"]
  }

  rule {
    display_name = "Nested-Group-Rule"
    source_groups = [nsxt_policy_group.nested_bms_group[0].path]
    destination_groups = [nsxt_policy_group.bmsi_group2[0].path]
    action = "ALLOW"
    services = ["/infra/services/HTTPS"]
  }
}

# Output members and associations results
output "members_associations" {
  value = {
    servers_count = length(local.selected_servers)
    interfaces_count = length(local.selected_interfaces)
    bms_group1_members = length(local.selected_servers) > 0 ? [local.selected_servers[0]] : []
    bms_group2_members = length(local.selected_servers) > 1 ? [local.selected_servers[1]] : length(local.selected_servers) > 0 ? [local.selected_servers[0]] : []
    bmsi_group1_members = length(local.selected_interfaces) > 0 ? [local.selected_interfaces[0]] : []
    bmsi_group2_members = length(local.selected_interfaces) > 1 ? [local.selected_interfaces[1]] : length(local.selected_interfaces) > 0 ? [local.selected_interfaces[0]] : []
    nested_group_path = local.has_servers ? nsxt_policy_group.nested_bms_group[0].path : ""
    associations_policy = local.has_servers && local.has_interfaces ? nsxt_policy_security_policy.associations_policy[0].path : ""
  }
}
`, testName, testName, testName, testName, testName, testName)
}

func testAccBMSEndToEndDataSourcesCoverageTemplate(testName string) string {
	return fmt.Sprintf(`
# Comprehensive inventory coverage
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

# Individual server lookup
data "nsxt_policy_baremetal_server" "single" {
  count = local.has_servers ? 1 : 0
  external_id = local.first_server_id
}

# Individual interface lookup 
data "nsxt_policy_baremetal_server_interface" "single" {
  count = local.has_interfaces ? 1 : 0
  external_id = local.first_interface_id
}

# Create and tag a server
resource "nsxt_policy_baremetal_server_tags" "server_tags" {
  count = local.has_servers ? 1 : 0
  external_id = local.first_server_id

  tag {
    scope = "test-ds"
    tag   = "coverage-test"
  }
}

# Create and tag an interface
resource "nsxt_policy_baremetal_server_interface_tags" "interface_tags" {
  count = local.has_interfaces ? 1 : 0
  external_id = local.first_interface_id

  tag {
    scope = "test-ds"
    tag   = "coverage-test"
  }
}

# Read server tags via data source
data "nsxt_policy_baremetal_server_tags" "server_tags" {
  count = local.has_servers ? 1 : 0
  external_id = local.first_server_id
  depends_on = [nsxt_policy_baremetal_server_tags.server_tags]
}

# Read interface tags via data source
data "nsxt_policy_baremetal_server_interface_tags" "interface_tags" {
  count = local.has_interfaces ? 1 : 0
  external_id = local.first_interface_id
  depends_on = [nsxt_policy_baremetal_server_interface_tags.interface_tags]
}

# Create BMS groups
resource "nsxt_policy_group" "server_group" {
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

resource "nsxt_policy_group" "interface_group" {
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

# Read groups via data sources
data "nsxt_policy_group" "server_group" {
  count = local.has_servers ? 1 : 0
  id = nsxt_policy_group.server_group[0].id
}

data "nsxt_policy_group" "interface_group" {
  count = local.has_interfaces ? 1 : 0
  id = nsxt_policy_group.interface_group[0].id
}

# Create security policy
resource "nsxt_policy_security_policy" "coverage_policy" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-coverage-policy"
  category     = "Application"

  rule {
    display_name       = "Coverage Rule"
    source_groups      = [nsxt_policy_group.server_group[0].path]
    destination_groups = local.has_interfaces ? [nsxt_policy_group.interface_group[0].path] : [nsxt_policy_group.server_group[0].path]
    action             = "ALLOW"
    services           = []
  }
}

# Read policy via data source
data "nsxt_policy_security_policy" "coverage_policy" {
  count = local.has_servers ? 1 : 0
  id = nsxt_policy_security_policy.coverage_policy[0].id
}

# Test associations if available
data "nsxt_policy_baremetal_server_group_associations" "server_associations" {
  count = local.has_servers ? 1 : 0
  external_id = local.first_server_id
}

data "nsxt_policy_baremetal_server_interface_group_associations" "interface_associations" {
  count = local.has_interfaces ? 1 : 0
  external_id = local.first_interface_id
}

# Output comprehensive data source coverage
output "data_source_coverage" {
  value = {
    servers_found = length(data.nsxt_policy_baremetal_servers.all.results)
    interfaces_found = length(data.nsxt_policy_baremetal_server_interfaces.all.results)
    server_tags_read = local.has_servers ? length(data.nsxt_policy_baremetal_server_tags.server_tags[0].tag) : 0
    interface_tags_read = local.has_interfaces ? length(data.nsxt_policy_baremetal_server_interface_tags.interface_tags[0].tag) : 0
    server_group_path = local.has_servers ? data.nsxt_policy_group.server_group[0].path : ""
    interface_group_path = local.has_interfaces ? data.nsxt_policy_group.interface_group[0].path : ""
    policy_path = local.has_servers ? data.nsxt_policy_security_policy.coverage_policy[0].path : ""
  }
}
`, testName, testName, testName)
}

func testAccBMSEndToEndDataSourcesConditionalCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Check if server group data source exists (conditional on having BMS servers)
		serverGroupDataSourceName := "data.nsxt_policy_group.server_group"
		if _, ok := s.RootModule().Resources[serverGroupDataSourceName+".0"]; ok {
			// If the data source exists, verify its attributes
			return resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(serverGroupDataSourceName+".0", "id"),
				resource.TestCheckResourceAttrSet(serverGroupDataSourceName+".0", "display_name"),
				resource.TestCheckResourceAttrSet(serverGroupDataSourceName+".0", "path"),
			)(s)
		}
		// If no BMS servers available, that's also valid
		return nil
	}
}
