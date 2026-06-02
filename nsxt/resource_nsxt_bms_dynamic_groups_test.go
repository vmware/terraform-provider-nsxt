// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccBMSDynamicGroupsOSName tests dynamic BMS groups based on OS name with all operators
func TestAccBMSDynamicGroupsOSName(t *testing.T) {
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
				Config: testAccBMSDynamicGroupsOSNameTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
				),
			},
		},
	})
}

// TestAccBMSDynamicGroupsDisplayName tests dynamic BMS groups based on display name with all operators
func TestAccBMSDynamicGroupsDisplayName(t *testing.T) {
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
				Config: testAccBMSDynamicGroupsDisplayNameTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
				),
			},
		},
	})
}

// TestAccBMSDynamicGroupsTags tests dynamic BMS groups based on tags with all operators
func TestAccBMSDynamicGroupsTags(t *testing.T) {
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
				Config: testAccBMSDynamicGroupsTagsTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
				),
			},
		},
	})
}

// TestAccBMSDynamicGroupsInterfaceTags tests dynamic BMS interface groups based on tags
func TestAccBMSDynamicGroupsInterfaceTags(t *testing.T) {
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
				Config: testAccBMSDynamicGroupsInterfaceTagsTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_interfaces.all", "results.#"),
				),
			},
		},
	})
}

// TestAccBMSDynamicGroupsALL tests dynamic BMS groups using ALL operator
func TestAccBMSDynamicGroupsALL(t *testing.T) {
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
				Config: testAccBMSDynamicGroupsALLTemplate(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
				),
			},
		},
	})
}

func testAccBMSDynamicGroupsOSNameTemplate(testName string) string {
	return fmt.Sprintf(`
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
}

# OS Name: EQUALS
resource "nsxt_policy_group" "bms_os_equals" {
  count = local.has_servers ? 1 : 0
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
  count = local.has_servers ? 1 : 0
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
  count = local.has_servers ? 1 : 0
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
  count = local.has_servers ? 1 : 0
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
  count = local.has_servers ? 1 : 0
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
  count = local.has_servers ? 1 : 0
  id = nsxt_policy_group.bms_os_equals[0].id
}
`, testName, testName, testName, testName, testName)
}

func testAccBMSDynamicGroupsDisplayNameTemplate(testName string) string {
	return fmt.Sprintf(`
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
}

# Display Name: EQUALS
resource "nsxt_policy_group" "bms_name_equals" {
  count = local.has_servers ? 1 : 0
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
  count = local.has_servers ? 1 : 0
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
  count = local.has_servers ? 1 : 0
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
  count = local.has_servers ? 1 : 0
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
  count = local.has_servers ? 1 : 0
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
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
}

# Tag: EQUALS
resource "nsxt_policy_group" "bms_tag_equals" {
  count = local.has_servers ? 1 : 0
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
  count = local.has_servers ? 1 : 0
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
  count = local.has_servers ? 1 : 0
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
  count = local.has_servers ? 1 : 0
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
# Discovery
data "nsxt_policy_baremetal_server_interfaces" "all" {}

locals {
  has_interfaces = length(data.nsxt_policy_baremetal_server_interfaces.all.results) > 0
}

# Interface Tag: EQUALS
resource "nsxt_policy_group" "bmsi_tag_equals" {
  count = local.has_interfaces ? 1 : 0
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
  count = local.has_interfaces ? 1 : 0
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
  count = local.has_interfaces ? 1 : 0
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
# Discovery
data "nsxt_policy_baremetal_servers" "all" {}

locals {
  has_servers = length(data.nsxt_policy_baremetal_servers.all.results) > 0
}

# Match all discovered BMS servers using external_id_expression
resource "nsxt_policy_group" "bms_all_equals" {
  count = local.has_servers ? 1 : 0
  display_name = "%s-bms-all-equals"
  description = "BMS group containing all discovered servers"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = data.nsxt_policy_baremetal_servers.all.results[*].external_id
    }
  }
}

# Read group using data source
data "nsxt_policy_group" "bms_all_equals_read" {
  count = local.has_servers ? 1 : 0
  id = nsxt_policy_group.bms_all_equals[0].id
}

# Output ALL group information
output "all_group_info" {
  value = {
    group_created = local.has_servers
    group_path = local.has_servers ? nsxt_policy_group.bms_all_equals[0].path : ""
  }
}
`, testName)
}
