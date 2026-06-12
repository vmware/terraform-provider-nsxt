// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccBMSDataSourceServers tests nsxt_policy_baremetal_servers data source
func TestAccBMSDataSourceServers(t *testing.T) {
	dataSourceName := "data.nsxt_policy_baremetal_servers.test"

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
				Config: testAccBMSDataSourceServersBasicTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

// TestAccBMSDataSourceServersWithFilters tests nsxt_policy_baremetal_servers with filters
func TestAccBMSDataSourceServersWithFilters(t *testing.T) {
	dataSourceName := "data.nsxt_policy_baremetal_servers.filtered"

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
				Config: testAccBMSDataSourceServersFilteredTemplate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

// TestAccBMSDataSourceServer tests nsxt_policy_baremetal_server data source
func TestAccBMSDataSourceServer(t *testing.T) {
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
				Config: testAccBMSDataSourceServerTemplate(serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server.single", "id"),
				),
			},
		},
	})
}

// TestAccBMSDataSourceInterfaces tests nsxt_policy_baremetal_server_interfaces data source
func TestAccBMSDataSourceInterfaces(t *testing.T) {
	dataSourceName := "data.nsxt_policy_baremetal_server_interfaces.test"

	resource.Test(t, resource.TestCase{
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

// TestAccBMSDataSourceInterface tests nsxt_policy_baremetal_server_interface data source
func TestAccBMSDataSourceInterface(t *testing.T) {
	interfaceID := getTestBMSInterfaceID()

	resource.Test(t, resource.TestCase{
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

// TestAccBMSDataSourceServerTags tests nsxt_policy_baremetal_server_tags data source
func TestAccBMSDataSourceServerTags(t *testing.T) {
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
				Config: testAccBMSDataSourceServerTagsTemplate(testName, serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_tags.test", "id"),
				),
			},
		},
	})
}

// TestAccBMSDataSourceInterfaceTags tests nsxt_policy_baremetal_server_interface_tags data source
func TestAccBMSDataSourceInterfaceTags(t *testing.T) {
	testName := getAccTestResourceName()
	interfaceID := getTestBMSInterfaceID()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccOnlyLocalManager(t)
			testAccPreCheck(t)
			testAccNSXVersion(t, "9.0.0")
			testAccEnvDefined(t, "NSXT_TEST_BMS_INTERFACE")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBMSDataSourceInterfaceTagsTemplate(testName, interfaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_interface_tags.test", "id"),
				),
			},
		},
	})
}

// TestAccBMSDataSourceValidationErrors tests various error scenarios for data sources
func TestAccBMSDataSourceValidationErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
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

func testAccBMSDataSourceServerTagsTemplate(testName, serverID string) string {
	return fmt.Sprintf(`
# Create tags for the test server
resource "nsxt_policy_baremetal_server_tags" "test" {
  external_id = "%s"

  tag {
    scope = "test-ds"
    tag   = "data-source-test"
  }
}

# Read the tags using data source
data "nsxt_policy_baremetal_server_tags" "test" {
  external_id = "%s"
  depends_on = [nsxt_policy_baremetal_server_tags.test]
}
`, serverID, serverID)
}

func testAccBMSDataSourceInterfaceTagsTemplate(testName, interfaceID string) string {
	return fmt.Sprintf(`
# Create tags for the test interface
resource "nsxt_policy_baremetal_server_interface_tags" "test" {
  external_id = "%s"

  tag {
    scope = "test-ds"
    tag   = "data-source-test"
  }
}

# Read the tags using data source
data "nsxt_policy_baremetal_server_interface_tags" "test" {
  external_id = "%s"
  depends_on = [nsxt_policy_baremetal_server_interface_tags.test]
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
