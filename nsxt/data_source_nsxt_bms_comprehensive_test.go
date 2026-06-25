// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccDataSourceNsxtPolicyBareMetalServers_basic tests nsxt_policy_baremetal_servers data source
func TestAccDataSourceNsxtPolicyBareMetalServers_basic(t *testing.T) {
	dataSourceName := "data.nsxt_policy_baremetal_servers.test"

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
				Config: testAccBMSDataSourceServersBasicTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					// Verify server results structure when available
					testAccBMSDataSourceServersEnhancedCheck(dataSourceName),
				),
			},
		},
	})
}

// TestAccDataSourceNsxtPolicyBareMetalServers_withFilters tests nsxt_policy_baremetal_servers with filters
func TestAccDataSourceNsxtPolicyBareMetalServers_withFilters(t *testing.T) {
	dataSourceName := "data.nsxt_policy_baremetal_servers.filtered"

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
				Config: testAccBMSDataSourceServersFilteredTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

// TestAccDataSourceNsxtPolicyBareMetalServer_basic tests nsxt_policy_baremetal_server data source
func TestAccDataSourceNsxtPolicyBareMetalServer_basic(t *testing.T) {
	serverID := getTestBMSServerID()

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
				Config: testAccBMSDataSourceServerTemplate(serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server.single", "id"),
					resource.TestCheckResourceAttr("data.nsxt_policy_baremetal_server.single", "external_id", serverID),
					// Verify individual server lookup attributes
					testAccBMSDataSourceSingleServerCheck("data.nsxt_policy_baremetal_server.single"),
				),
			},
		},
	})
}

// TestAccDataSourceNsxtPolicyBareMetalServerInterfaces_basic tests nsxt_policy_baremetal_server_interfaces data source
func TestAccDataSourceNsxtPolicyBareMetalServerInterfaces_basic(t *testing.T) {
	dataSourceName := "data.nsxt_policy_baremetal_server_interfaces.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_INTERFACE")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSDataSourceInterfacesTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

// TestAccDataSourceNsxtPolicyBareMetalServerInterface_basic tests nsxt_policy_baremetal_server_interface data source
func TestAccDataSourceNsxtPolicyBareMetalServerInterface_basic(t *testing.T) {
	interfaceID := getTestBMSInterfaceID()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_INTERFACE")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSDataSourceInterfaceTemplate(interfaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_interface.single", "id"),
				),
			},
		},
	})
}

// TestAccDataSourceNsxtPolicyBareMetalServerTags_validation tests nsxt_policy_baremetal_server_tags data source
func TestAccDataSourceNsxtPolicyBareMetalServerTags_validation(t *testing.T) {
	serverID := getTestBMSServerID()

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
				Config: testAccBMSDataSourceServerTagsTemplate(serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_tags.test", "id"),
				),
			},
		},
	})
}

// TestAccDataSourceNsxtPolicyBareMetalServerInterfaceTags_validation tests nsxt_policy_baremetal_server_interface_tags data source
func TestAccDataSourceNsxtPolicyBareMetalServerInterfaceTags_validation(t *testing.T) {
	interfaceID := getTestBMSInterfaceID()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_INTERFACE")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSDataSourceInterfaceTagsTemplate(interfaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_interface_tags.test", "id"),
				),
			},
		},
	})
}

// TestAccDataSourceNsxtPolicyBareMetalServer_validationErrors tests various error scenarios for data sources
func TestAccDataSourceNsxtPolicyBareMetalServer_validationErrors(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccOnlyLocalManager(t); testAccPreCheck(t); testAccNSXVersion(t, "9.0.0") },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccBMSDataSourceInvalidRegexTemplate(),
				ExpectError: regexp.MustCompile("invalid regex for display_name"),
			},
			{
				Config: testAccBMSDataSourceEmptyTagScopeTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.invalid_tag_scope", "id"),
				),
			},
			{
				Config:      testAccBMSDataSourceMissingIdentifierTemplate(),
				ExpectError: regexp.MustCompile("either external_id or display_name must be specified"),
			},
		},
	})
}

// TestAccDataSourceNsxtPolicyBareMetalServer_displayNameValidation tests display_name validation and error scenarios
func TestAccDataSourceNsxtPolicyBareMetalServer_displayNameValidation(t *testing.T) {
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
				// Test that non-existent display_name returns appropriate error
				Config:      testAccBMSDataSourceNonExistentDisplayNameTemplate(),
				ExpectError: regexp.MustCompile("No bare metal server found with display_name"),
			},
		},
	})
}

// TestAccDataSourceNsxtPolicyBareMetalServerInterface_displayNameValidation tests interface display_name validation
func TestAccDataSourceNsxtPolicyBareMetalServerInterface_displayNameValidation(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_INTERFACE")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test that non-existent interface display_name returns appropriate error
				Config:      testAccBMSInterfaceDataSourceNonExistentDisplayNameTemplate(),
				ExpectError: regexp.MustCompile("No bare metal server interface found with display_name"),
			},
		},
	})
}

// TestAccDataSourceNsxtPolicyBareMetalServerInterface_duplicateDisplayNameValidation tests interface duplicate display_name validation
func TestAccDataSourceNsxtPolicyBareMetalServerInterface_duplicateDisplayNameValidation(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_INTERFACE_NAME")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test that duplicate interface display_name returns appropriate error
				Config:      testAccBMSInterfaceDataSourceDuplicateDisplayNameTemplate(),
				ExpectError: regexp.MustCompile("Multiple bare metal server interfaces found with display_name.*Use external_id for exact match"),
			},
		},
	})
}

// Template functions

func testAccBMSDataSourceServersBasicTemplate() string {
	return `
data "nsxt_policy_baremetal_servers" "test" {}
`
}

func testAccBMSDataSourceServersFilteredTemplate() string {
	return `
data "nsxt_policy_baremetal_servers" "filtered" {
  display_name = "non-existent-server"
}
`
}

func testAccBMSDataSourceServerTemplate(serverID string) string {
	return fmt.Sprintf(`
# Get specific server by external ID from environment
data "nsxt_policy_baremetal_server" "single" {
  external_id = "%s"
}
`, serverID)
}

func testAccBMSDataSourceInterfacesTemplate() string {
	return `
data "nsxt_policy_baremetal_server_interfaces" "test" {}
`
}

func testAccBMSDataSourceInterfaceTemplate(interfaceID string) string {
	return fmt.Sprintf(`
# Get specific interface by external ID from environment
data "nsxt_policy_baremetal_server_interface" "single" {
  external_id = "%s"
}
`, interfaceID)
}

func testAccBMSDataSourceServerTagsTemplate(serverID string) string {
	return fmt.Sprintf(`
# Create tags for the test server
resource "nsxt_policy_baremetal_server_tags" "comprehensive_server_tags_test" {
  external_id = "%s"

  tag {
    scope = "test-ds"
    tag   = "data-source-test"
  }
}

# Read the tags using data source
data "nsxt_policy_baremetal_server_tags" "test" {
  external_id = "%s"
  depends_on = [nsxt_policy_baremetal_server_tags.comprehensive_server_tags_test]
}
`, serverID, serverID)
}

func testAccBMSDataSourceInterfaceTagsTemplate(interfaceID string) string {
	return fmt.Sprintf(`
# Create tags for the test interface
resource "nsxt_policy_baremetal_server_interface_tags" "comprehensive_interface_tags_test" {
  external_id = "%s"

  tag {
    scope = "test-ds"
    tag   = "data-source-test"
  }
}

# Read the tags using data source
data "nsxt_policy_baremetal_server_interface_tags" "test" {
  external_id = "%s"
  depends_on = [nsxt_policy_baremetal_server_interface_tags.comprehensive_interface_tags_test]
}
`, interfaceID, interfaceID)
}

func testAccBMSDataSourceInvalidRegexTemplate() string {
	return `
data "nsxt_policy_baremetal_servers" "invalid_regex" {
  display_name = "invalid[regex"
}
`
}

func testAccBMSDataSourceEmptyTagScopeTemplate() string {
	return `
data "nsxt_policy_baremetal_servers" "invalid_tag_scope" {
  tag_scope = ""
  tag = "production"
}
`
}

func testAccBMSDataSourceMissingIdentifierTemplate() string {
	return `
data "nsxt_policy_baremetal_server" "missing_identifier" {
  # Neither external_id nor display_name specified
}
`
}

func testAccBMSDataSourceNonExistentDisplayNameTemplate() string {
	return `
# Test that searching for a non-existent display_name returns proper error
data "nsxt_policy_baremetal_server" "nonexistent" {
  display_name = "absolutely-nonexistent-server-name-for-testing-12345"
}
`
}

func testAccBMSInterfaceDataSourceNonExistentDisplayNameTemplate() string {
	return `
# Test that searching for a non-existent interface display_name returns proper error  
data "nsxt_policy_baremetal_server_interface" "nonexistent" {
  display_name = "absolutely-nonexistent-interface-name-for-testing-12345"
}
`
}

func testAccBMSInterfaceDataSourceDuplicateDisplayNameTemplate() string {
	interfaceName := os.Getenv("NSXT_TEST_BMS_INTERFACE_NAME")
	return fmt.Sprintf(`
# Test that searching for a duplicate interface display_name returns proper error
# This test uses NSXT_TEST_BMS_INTERFACE environment variable which should be a common name like "eth0"

data "nsxt_policy_baremetal_server_interface" "duplicate" {
  display_name = "%s"
}
`, interfaceName)
}

// Enhanced validation functions for comprehensive data source tests
func testAccBMSDataSourceServersEnhancedCheck(dataSourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if rs, ok := s.RootModule().Resources[dataSourceName]; ok {
			if rs.Primary.Attributes["results.#"] != "0" {
				// If we have results, validate fields
				return resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.external_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.os_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.cpu_cores"),
				)(s)
			}
		}
		return nil
	}
}

func testAccBMSDataSourceSingleServerCheck(dataSourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if _, ok := s.RootModule().Resources[dataSourceName]; ok {
			return resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				resource.TestCheckResourceAttrSet(dataSourceName, "external_id"),
				resource.TestCheckResourceAttrSet(dataSourceName, "display_name"),
				resource.TestCheckResourceAttrSet(dataSourceName, "os_name"),
				resource.TestCheckResourceAttrSet(dataSourceName, "cpu_cores"),
			)(s)
		}
		return nil
	}
}
