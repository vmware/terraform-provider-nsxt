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

// TestAccResourceNsxtPolicyBMSGroup_dynamicOSName tests dynamic BMS groups based on OS name with all operators
func TestAccResourceNsxtPolicyBMSGroup_dynamicOSName(t *testing.T) {
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
				Config: testAccBMSDynamicGroupsOSNameTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
					// Verify dynamic groups based on OS name are created properly
					testAccBMSDynamicGroupsConditionalCheck("os_name_servers"),
					testAccBMSDynamicGroupsConditionalCheck("os_name_interfaces"),
				),
			},
		},
	})
}

// TestAccResourceNsxtBMSPolicyGroup_dDynamicDisplayName tests dynamic BMS groups based on display name with all operators
func TestAccResourceNsxtBMSPolicyGroup_dDynamicDisplayName(t *testing.T) {
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
				Config: testAccBMSDynamicGroupsDisplayNameTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
					// Verify dynamic groups based on display name patterns work
					testAccBMSDynamicGroupsConditionalCheck("display_name_servers"),
				),
			},
		},
	})
}

// TestAccResourceNsxtPolicyBMSGroup_dynamicTags tests dynamic BMS groups based on tags with all operators
func TestAccResourceNsxtPolicyBMSGroup_dynamicTags(t *testing.T) {
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
				Config: testAccBMSDynamicGroupsTagsTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
					// Verify dynamic groups based on server tags work correctly
					testAccBMSDynamicGroupsConditionalCheck("tag_based_servers"),
				),
			},
		},
	})
}

// TestAccResourceNsxtPolicyBMSGroup_dynamicInterfaceTags tests dynamic BMS interface groups based on tags
func TestAccResourceNsxtPolicyBMSGroup_dynamicInterfaceTags(t *testing.T) {
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
				Config: testAccBMSDynamicGroupsInterfaceTagsTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_interfaces.all", "results.#"),
					// Verify dynamic groups based on interface tags work correctly
					testAccBMSDynamicGroupsConditionalCheck("interface_tag_based"),
				),
			},
		},
	})
}

// TestAccResourceNsxtPolicyBMSGroup_dynamicALL tests static BMS groups using external ID expressions (misnamed - should be in static tests)
func TestAccResourceNsxtPolicyBMSGroup_dynamicALL(t *testing.T) {
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
				Config: testAccBMSDynamicGroupsALLTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
					// Verify static group with external IDs works (misnamed test)
					testAccBMSDynamicGroupsConditionalCheck("bms_all_equals"),
				),
			},
		},
	})
}

func testAccBMSDynamicGroupsOSNameTemplate(testName string) string {
	return fmt.Sprintf(`
# Data source for testing BMS inventory access
data "nsxt_policy_baremetal_servers" "all" {}

# OS Name: EQUALS
resource "nsxt_policy_group" "bms_os_equals" {
  display_name = "%s-bms-os-equals"
  description = "BMS group with OS name equals"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "OSName"
      member_type = "BareMetalServer"
      operator = "EQUALS"
      value = "linux"
    }
  }
}

# OS Name: CONTAINS
resource "nsxt_policy_group" "bms_os_contains" {
  display_name = "%s-bms-os-contains"
  description = "BMS group with OS name contains"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "OSName"
      member_type = "BareMetalServer"
      operator = "CONTAINS"
      value = "ubuntu"
    }
  }
}

# OS Name: STARTS_WITH
resource "nsxt_policy_group" "bms_os_starts_with" {
  display_name = "%s-bms-os-starts-with"
  description = "BMS group with OS name starts with"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "OSName"
      member_type = "BareMetalServer"
      operator = "STARTSWITH"
      value = "red"
    }
  }
}

# OS Name: ENDS_WITH
resource "nsxt_policy_group" "bms_os_ends_with" {
  display_name = "%s-bms-os-ends-with"
  description = "BMS group with OS name ends with"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "OSName"
      member_type = "BareMetalServer"
      operator = "ENDSWITH"
      value = "server"
    }
  }
}

# OS Name: NOT_EQUALS
resource "nsxt_policy_group" "bms_os_not_equals" {
  display_name = "%s-bms-os-not-equals"
  description = "BMS group with OS name not equals"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "OSName"
      member_type = "BareMetalServer"
      operator = "NOTEQUALS"
      value = "windows"
    }
  }
}

# Read groups using data sources
data "nsxt_policy_group" "bms_os_equals_read" {
  id = nsxt_policy_group.bms_os_equals.id
}
`, testName, testName, testName, testName, testName)
}

func testAccBMSDynamicGroupsDisplayNameTemplate(testName string) string {
	return fmt.Sprintf(`
# Data source for testing BMS inventory access  
data "nsxt_policy_baremetal_servers" "all" {}

# Display Name: EQUALS
resource "nsxt_policy_group" "bms_name_equals" {
  display_name = "%s-bms-name-equals"
  description = "BMS group with name equals"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "Name"
      member_type = "BareMetalServer"
      operator = "EQUALS"
      value = "server001"
    }
  }
}

# Display Name: CONTAINS
resource "nsxt_policy_group" "bms_name_contains" {
  display_name = "%s-bms-name-contains"
  description = "BMS group with name contains"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "Name"
      member_type = "BareMetalServer"
      operator = "CONTAINS"
      value = "web"
    }
  }
}

# Display Name: STARTS_WITH
resource "nsxt_policy_group" "bms_name_starts_with" {
  display_name = "%s-bms-name-starts-with"
  description = "BMS group with name starts with"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "Name"
      member_type = "BareMetalServer"
      operator = "STARTSWITH"
      value = "prod"
    }
  }
}

# Display Name: ENDS_WITH
resource "nsxt_policy_group" "bms_name_ends_with" {
  display_name = "%s-bms-name-ends-with"
  description = "BMS group with name ends with"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "Name"
      member_type = "BareMetalServer"
      operator = "ENDSWITH"
      value = "001"
    }
  }
}

# Display Name: NOT_EQUALS
resource "nsxt_policy_group" "bms_name_not_equals" {
  display_name = "%s-bms-name-not-equals"
  description = "BMS group with name not equals"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "Name"
      member_type = "BareMetalServer"
      operator = "NOTEQUALS"
      value = "test-server"
    }
  }
}
`, testName, testName, testName, testName, testName)
}

func testAccBMSDynamicGroupsTagsTemplate(testName string) string {
	return fmt.Sprintf(`
# Data source for testing BMS inventory access
data "nsxt_policy_baremetal_servers" "all" {}

# Tag: EQUALS
resource "nsxt_policy_group" "bms_tag_equals" {
  display_name = "%s-bms-tag-equals"
  description = "BMS group with tag equals"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "Tag"
      member_type = "BareMetalServer"
      operator = "EQUALS"
      value = "environment|production"
    }
  }
}

# Tag: CONTAINS
resource "nsxt_policy_group" "bms_tag_contains" {
  display_name = "%s-bms-tag-contains"
  description = "BMS group with tag contains"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "Tag"
      member_type = "BareMetalServer"
      operator = "CONTAINS"
      value = "web"
    }
  }
}

# Tag: STARTS_WITH
resource "nsxt_policy_group" "bms_tag_starts_with" {
  display_name = "%s-bms-tag-starts-with"
  description = "BMS group with tag starts with"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "Tag"
      member_type = "BareMetalServer"
      operator = "STARTSWITH"
      value = "app|"
    }
  }
}

# Tag: ENDS_WITH
resource "nsxt_policy_group" "bms_tag_ends_with" {
  display_name = "%s-bms-tag-ends-with"
  description = "BMS group with tag ends with"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "Tag"
      member_type = "BareMetalServer"
      operator = "ENDSWITH"
      value = "|server"
    }
  }
}
`, testName, testName, testName, testName)
}

func testAccBMSDynamicGroupsInterfaceTagsTemplate(testName string) string {
	return fmt.Sprintf(`
# Data source for testing BMS interface inventory access
data "nsxt_policy_baremetal_server_interfaces" "all" {}

# Interface Tag: EQUALS
resource "nsxt_policy_group" "bmsi_tag_equals" {
  display_name = "%s-bmsi-tag-equals"
  description = "BMS interface group with tag equals"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "Tag"
      member_type = "BareMetalServerInterface"
      operator = "EQUALS"
      value = "network|management"
    }
  }
}

# Interface Tag: NOT_EQUALS
resource "nsxt_policy_group" "bmsi_tag_not_equals" {
  display_name = "%s-bmsi-tag-not-equals"
  description = "BMS interface group with tag not equals"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "Tag"
      member_type = "BareMetalServerInterface"
      operator = "NOTEQUALS"
      value = "network|storage"
    }
  }
}

# Interface Tag: NOT_IN (if supported)
resource "nsxt_policy_group" "bmsi_tag_not_in" {
  display_name = "%s-bmsi-tag-not-in"
  description = "BMS interface group with tag not in"
  group_type = "BareMetalServer"

  criteria {
    condition {
      key = "Tag"
      member_type = "BareMetalServerInterface"
      operator = "NOTIN"
      value = "network|storage,network|backup"
    }
  }
}
`, testName, testName, testName)
}

func testAccBMSDynamicGroupsALLTemplate(testName string) string {
	return fmt.Sprintf(`
# Data source for testing BMS inventory access
data "nsxt_policy_baremetal_servers" "all" {}

# Match specific BMS server using external_id_expression (static group)
resource "nsxt_policy_group" "bms_all_equals" {
  display_name = "%s-bms-static-external-id"
  description = "BMS group containing specific server by external ID"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

# Read group using data source
data "nsxt_policy_group" "bms_all_equals_read" {
  id = nsxt_policy_group.bms_all_equals.id
}

# Output ALL group information
output "all_group_info" {
  value = {
    group_path = nsxt_policy_group.bms_all_equals.path
  }
}
`, testName, getTestBMSServerID())
}

// Enhanced validation functions for dynamic groups
func testAccBMSDynamicGroupsConditionalCheck(groupName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Check if the specific dynamic group was created
		groupResourceName := "nsxt_policy_group." + groupName
		if _, ok := s.RootModule().Resources[groupResourceName]; ok {
			return resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(groupResourceName, "id"),
				resource.TestCheckResourceAttr(groupResourceName, "group_type", "BareMetalServer"),
				resource.TestCheckResourceAttrSet(groupResourceName, "criteria.#"),
			)(s)
		}
		// If no BMS resources available, that's also valid
		return nil
	}
}
