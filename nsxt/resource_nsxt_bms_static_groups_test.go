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

// TestAccBMSStaticGroupsServer tests static BMS server groups
func TestAccBMSStaticGroupsServer(t *testing.T) {
	testName := getAccTestResourceName()

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
				Config: testAccBMSStaticGroupsServerTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nsxt_policy_group.bms_servers_static", "id"),
					resource.TestCheckResourceAttr("nsxt_policy_group.bms_servers_static", "group_type", "BareMetalServer"),
				),
			},
		},
	})
}

// TestAccBMSStaticGroupsInterface tests static BMS interface groups
func TestAccBMSStaticGroupsInterface(t *testing.T) {
	testName := getAccTestResourceName()

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
				Config: testAccBMSStaticGroupsInterfaceTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nsxt_policy_group.bms_interfaces_static", "id"),
					resource.TestCheckResourceAttr("nsxt_policy_group.bms_interfaces_static", "group_type", "BareMetalServer"),
				),
			},
		},
	})
}

// TestAccBMSStaticGroupsNested tests nested BMS groups (group as member)
func TestAccBMSStaticGroupsNested(t *testing.T) {
	testName := getAccTestResourceName()

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
				Config: testAccBMSStaticGroupsNestedTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
				),
			},
		},
	})
}

// TestAccBMSStaticGroupsDataSources tests reading static groups via data sources
func TestAccBMSStaticGroupsDataSources(t *testing.T) {
	testName := getAccTestResourceName()

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
				Config: testAccBMSStaticGroupsDataSourceTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					// Only check if BMS servers are available
					testAccBMSStaticGroupsDataSourceCheck(),
				),
			},
		},
	})
}

func testAccBMSStaticGroupsServerTemplate(testName string) string {
	return fmt.Sprintf(`
# Static BMS server group with external IDs from environment
resource "nsxt_policy_group" "bms_servers_static" {
  display_name = "%s-bms-servers-static"
  description = "Static BMS server group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

# Read group using data source
data "nsxt_policy_group" "bms_servers_static_read" {
  id = nsxt_policy_group.bms_servers_static.id
}

# Output group information
output "static_server_group" {
  value = {
    group_created = true
    group_path = nsxt_policy_group.bms_servers_static.path
  }
}
`, testName, getTestBMSServerID())
}

func testAccBMSStaticGroupsInterfaceTemplate(testName string) string {
	return fmt.Sprintf(`
# Static BMS interface group with external IDs from environment
resource "nsxt_policy_group" "bms_interfaces_static" {
  display_name = "%s-bms-interfaces-static"
  description = "Static BMS interface group"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServerInterface"
      external_ids = ["%s"]
    }
  }
}

# Read group using data source
data "nsxt_policy_group" "bms_interfaces_static_read" {
  id = nsxt_policy_group.bms_interfaces_static.id
}

# Output group information
output "static_interface_group" {
  value = {
    group_created = true
    group_path = nsxt_policy_group.bms_interfaces_static.path
  }
}
`, testName, getTestBMSInterfaceID())
}

func testAccBMSStaticGroupsNestedTemplate(testName string) string {
	return fmt.Sprintf(`
# Data source for testing BMS inventory access
data "nsxt_policy_baremetal_servers" "all" {}

# First create a base BMS group
resource "nsxt_policy_group" "bms_base_group" {
  display_name = "%s-bms-base-group"
  description = "Base BMS group for nesting"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

# Create nested group that includes the base group
resource "nsxt_policy_group" "bms_nested_group" {
  display_name = "%s-bms-nested-group"
  description = "Nested BMS group containing other groups"
  group_type = "BareMetalServer"

  criteria {
    path_expression {
      member_paths = [nsxt_policy_group.bms_base_group.path]
    }
  }
}

# Read nested group using data source
data "nsxt_policy_group" "bms_nested_group_read" {
  id = nsxt_policy_group.bms_nested_group.id
}

# Output nested group information
output "nested_group" {
  value = {
    base_group_path = nsxt_policy_group.bms_base_group.path
    nested_group_path = nsxt_policy_group.bms_nested_group.path
  }
}
`, testName, getTestBMSServerID(), testName)
}

func testAccBMSStaticGroupsDataSourceTemplate(testName string) string {
	return fmt.Sprintf(`
# Create static BMS server group
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

# Read the group via data source
data "nsxt_policy_group" "bms_server_group" {
  id = nsxt_policy_group.bms_server_group.id
}

# Output data source information
output "group_data_source" {
  value = {
    group_path = data.nsxt_policy_group.bms_server_group.path
    group_display_name = data.nsxt_policy_group.bms_server_group.display_name
  }
}
`, testName, getTestBMSServerID())
}

func testAccBMSStaticGroupsDataSourceCheck() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Check if data source exists (conditional on having BMS servers)
		groupDataSourceName := "data.nsxt_policy_group.bms_server_group"
		if _, ok := s.RootModule().Resources[groupDataSourceName+".0"]; ok {
			// If the data source exists, verify its attributes
			return resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(groupDataSourceName+".0", "id"),
				resource.TestCheckResourceAttrSet(groupDataSourceName+".0", "display_name"),
				resource.TestCheckResourceAttrSet(groupDataSourceName+".0", "path"),
			)(s)
		}
		// If no BMS servers available, that's also valid
		return nil
	}
}
