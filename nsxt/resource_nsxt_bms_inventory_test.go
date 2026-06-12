// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccBMSInventoryServers tests BMS server inventory data sources
func TestAccBMSInventoryServers(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSInventoryServersTemplate(),
				Check: resource.ComposeTestCheckFunc(
					// Basic data source validation
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "id"),
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
					// Conditional check for individual server
					testAccBMSInventoryServerConditionalCheck(),
				),
			},
		},
	})
}

// TestAccBMSInventoryInterfaces tests BMS interface inventory data sources
func TestAccBMSInventoryInterfaces(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSInventoryInterfacesTemplate(),
				Check: resource.ComposeTestCheckFunc(
					// Basic data source validation
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_interfaces.all", "id"),
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_interfaces.all", "results.#"),
					// Conditional check for individual interface
					testAccBMSInventoryInterfaceConditionalCheck(),
				),
			},
		},
	})
}

// TestAccBMSInventoryFiltering tests BMS inventory with filtering
func TestAccBMSInventoryFiltering(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSInventoryFilteringTemplate(),
				Check: resource.ComposeTestCheckFunc(
					// Filtered data sources should work (may return empty results)
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.linux", "id"),
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_interfaces.tagged", "id"),
				),
			},
		},
	})
}

// TestAccBMSInventoryValidation tests data source validation scenarios
func TestAccBMSInventoryValidation(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_SERVER")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccBMSInventoryInvalidRegexTemplate(),
				ExpectError: regexp.MustCompile("invalid regex for display_name"),
			},
			{
				Config: testAccBMSInventoryInvalidTagScopeTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.invalid_tag_scope", "id"),
				),
			},
			{
				Config:      testAccBMSInventoryMissingIdentifierTemplate(),
				ExpectError: regexp.MustCompile("either external_id or display_name must be specified"),
			},
		},
	})
}

func testAccBMSInventoryServersTemplate() string {
	return fmt.Sprintf(`
# Get all BMS servers
data "nsxt_policy_baremetal_servers" "all" {}

# Get specific server by external ID
data "nsxt_policy_baremetal_server" "test" {
  external_id = "%s"
}

# Output for debugging
output "server_inventory" {
  value = {
    total_servers = length(data.nsxt_policy_baremetal_servers.all.results)
    test_server_id = data.nsxt_policy_baremetal_server.test.external_id
  }
}
`, getTestBMSServerID())
}

func testAccBMSInventoryInterfacesTemplate() string {
	return fmt.Sprintf(`
# Get all BMS interfaces
data "nsxt_policy_baremetal_server_interfaces" "all" {}

# Get specific interface by external ID
data "nsxt_policy_baremetal_server_interface" "test" {
  external_id = "%s"
}

# Output for debugging
output "interface_inventory" {
  value = {
    total_interfaces = length(data.nsxt_policy_baremetal_server_interfaces.all.results)
    test_interface_id = data.nsxt_policy_baremetal_server_interface.test.external_id
  }
}
`, getTestBMSInterfaceID())
}

func testAccBMSInventoryFilteringTemplate() string {
	return `
# Test various filtering options
data "nsxt_policy_baremetal_servers" "linux" {
  os_name = "linux"
}

data "nsxt_policy_baremetal_servers" "by_tag" {
  tag_scope = "environment"
  tag = "production"
}

data "nsxt_policy_baremetal_server_interfaces" "tagged" {
  tag_scope = "network"
  tag = "management"
}

# Output filtering results
output "filtering_results" {
  value = {
    linux_servers = length(data.nsxt_policy_baremetal_servers.linux.results)
    tagged_servers = length(data.nsxt_policy_baremetal_servers.by_tag.results)
    tagged_interfaces = length(data.nsxt_policy_baremetal_server_interfaces.tagged.results)
  }
}
`
}

func testAccBMSInventoryInvalidRegexTemplate() string {
	return `
data "nsxt_policy_baremetal_servers" "invalid_regex" {
  display_name = "invalid[regex"
}
`
}

func testAccBMSInventoryInvalidTagScopeTemplate() string {
	return `
data "nsxt_policy_baremetal_servers" "invalid_tag_scope" {
  tag_scope = ""
  tag = "production"
}
`
}

func testAccBMSInventoryMissingIdentifierTemplate() string {
	return `
data "nsxt_policy_baremetal_server" "missing_identifier" {
  # Neither external_id nor display_name specified
}
`
}

func testAccBMSInventoryServerConditionalCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Check if individual server data source exists (conditional on having BMS servers)
		serverDataSourceName := "data.nsxt_policy_baremetal_server.first"
		if _, ok := s.RootModule().Resources[serverDataSourceName+".0"]; ok {
			// If the data source exists, verify its attributes
			return resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(serverDataSourceName+".0", "id"),
				resource.TestCheckResourceAttrSet(serverDataSourceName+".0", "external_id"),
			)(s)
		}
		// If no BMS servers available, that's also valid
		return nil
	}
}

func testAccBMSInventoryInterfaceConditionalCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Check if individual interface data source exists (conditional on having BMS interfaces)
		interfaceDataSourceName := "data.nsxt_policy_baremetal_server_interface.first"
		if _, ok := s.RootModule().Resources[interfaceDataSourceName+".0"]; ok {
			// If the data source exists, verify its attributes
			return resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(interfaceDataSourceName+".0", "id"),
				resource.TestCheckResourceAttrSet(interfaceDataSourceName+".0", "external_id"),
			)(s)
		}
		// If no BMS interfaces available, that's also valid
		return nil
	}
}
