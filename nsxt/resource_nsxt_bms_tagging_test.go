// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccResourceNsxtPolicyBareMetalServerTags_crud tests complete CRUD operations for BMS server tags
func TestAccResourceNsxtPolicyBareMetalServerTags_crud(t *testing.T) {
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
				// Step 1: Create tags
				Config: testAccBMSTaggingServerCreateTemplate(serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nsxt_policy_baremetal_server_tags.server_tags_crud", "id"),
					resource.TestCheckResourceAttr("nsxt_policy_baremetal_server_tags.server_tags_crud", "tag.#", "2"),
					// Tags are a set, so we can't predict order. Just check they exist.
					resource.TestCheckTypeSetElemNestedAttrs("nsxt_policy_baremetal_server_tags.server_tags_crud", "tag.*", map[string]string{
						"scope": "test-env",
						"tag":   "development",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("nsxt_policy_baremetal_server_tags.server_tags_crud", "tag.*", map[string]string{
						"scope": "test-app",
						"tag":   "web-server",
					}),
				),
			},
			{
				// Step 2: Update tags
				Config: testAccBMSTaggingServerUpdateTemplate(serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nsxt_policy_baremetal_server_tags.server_tags_crud", "id"),
					resource.TestCheckResourceAttr("nsxt_policy_baremetal_server_tags.server_tags_crud", "tag.#", "3"),
					// Check all three tags exist in the set
					resource.TestCheckTypeSetElemNestedAttrs("nsxt_policy_baremetal_server_tags.server_tags_crud", "tag.*", map[string]string{
						"scope": "test-env",
						"tag":   "production",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("nsxt_policy_baremetal_server_tags.server_tags_crud", "tag.*", map[string]string{
						"scope": "test-app",
						"tag":   "database",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("nsxt_policy_baremetal_server_tags.server_tags_crud", "tag.*", map[string]string{
						"scope": "test-team",
						"tag":   "backend",
					}),
				),
			},
			{
				// Step 3: Read tags through data sources
				Config: testAccBMSTaggingServerReadTemplate(serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_tags.server_tags_crud", "id"),
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_tags.server_tags_crud", "external_id"),
					resource.TestCheckResourceAttr("data.nsxt_policy_baremetal_server_tags.server_tags_crud", "external_id", serverID),
					// Verify tags were properly read back (3 tags from template)
					testAccBMSTaggingEnhancedCheck("data.nsxt_policy_baremetal_server_tags.server_tags_crud", 3),
				),
			},
		},
	})
}

// TestAccResourceNsxtPolicyBareMetalServerInterfaceTags_crud tests complete CRUD operations for BMS interface tags
func TestAccResourceNsxtPolicyBareMetalServerInterfaceTags_crud(t *testing.T) {
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
				// Step 1: Create tags
				Config: testAccBMSTaggingInterfaceCreateTemplate(interfaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nsxt_policy_baremetal_server_interface_tags.interface_tags_crud", "id"),
					resource.TestCheckResourceAttr("nsxt_policy_baremetal_server_interface_tags.interface_tags_crud", "tag.#", "2"),
					// Check both tags exist in the set
					resource.TestCheckTypeSetElemNestedAttrs("nsxt_policy_baremetal_server_interface_tags.interface_tags_crud", "tag.*", map[string]string{
						"scope": "test-net",
						"tag":   "management",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("nsxt_policy_baremetal_server_interface_tags.interface_tags_crud", "tag.*", map[string]string{
						"scope": "test-vlan",
						"tag":   "100",
					}),
				),
			},
			{
				// Step 2: Update tags
				Config: testAccBMSTaggingInterfaceUpdateTemplate(interfaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("nsxt_policy_baremetal_server_interface_tags.interface_tags_crud", "id"),
					resource.TestCheckResourceAttr("nsxt_policy_baremetal_server_interface_tags.interface_tags_crud", "tag.#", "3"),
					// Check all three updated tags exist in the set
					resource.TestCheckTypeSetElemNestedAttrs("nsxt_policy_baremetal_server_interface_tags.interface_tags_crud", "tag.*", map[string]string{
						"scope": "test-net",
						"tag":   "data-plane",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("nsxt_policy_baremetal_server_interface_tags.interface_tags_crud", "tag.*", map[string]string{
						"scope": "test-vlan",
						"tag":   "200",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("nsxt_policy_baremetal_server_interface_tags.interface_tags_crud", "tag.*", map[string]string{
						"scope": "test-speed",
						"tag":   "10gb",
					}),
				),
			},
			{
				// Step 3: Read tags through data sources
				Config: testAccBMSTaggingInterfaceReadTemplate(interfaceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_interface_tags.interface_tags_crud", "id"),
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_server_interface_tags.interface_tags_crud", "external_id"),
					resource.TestCheckResourceAttr("data.nsxt_policy_baremetal_server_interface_tags.interface_tags_crud", "external_id", interfaceID),
					// Verify interface tags were properly read back (3 tags from template)
					testAccBMSTaggingEnhancedCheck("data.nsxt_policy_baremetal_server_interface_tags.interface_tags_crud", 3),
				),
			},
		},
	})
}

func testAccBMSTaggingServerCreateTemplate(serverID string) string {
	return fmt.Sprintf(`
# Create server tags
resource "nsxt_policy_baremetal_server_tags" "server_tags_crud" {
  external_id = "%s"

  tag {
    scope = "test-env"
    tag   = "development"
  }

  tag {
    scope = "test-app"
    tag   = "web-server"
  }
}
`, serverID)
}

func testAccBMSTaggingServerUpdateTemplate(serverID string) string {
	return fmt.Sprintf(`
# Update server tags
resource "nsxt_policy_baremetal_server_tags" "server_tags_crud" {
  external_id = "%s"

  tag {
    scope = "test-env"
    tag   = "production"  # Updated value
  }

  tag {
    scope = "test-app"
    tag   = "database"    # Updated value
  }

  tag {
    scope = "test-team"  # New tag
    tag   = "backend"
  }
}
`, serverID)
}

func testAccBMSTaggingServerReadTemplate(serverID string) string {
	return fmt.Sprintf(`
# Update server tags (same as previous step)
resource "nsxt_policy_baremetal_server_tags" "server_tags_crud" {
  external_id = "%s"

  tag {
    scope = "test-env"
    tag   = "production"
  }

  tag {
    scope = "test-app"
    tag   = "database"
  }

  tag {
    scope = "test-team"
    tag   = "backend"
  }
}

# Read tags through data source
data "nsxt_policy_baremetal_server_tags" "server_tags_crud" {
  external_id = "%s"
}

# Output tag information
output "server_tags_crud" {
  value = data.nsxt_policy_baremetal_server_tags.server_tags_crud.tag
}
`, serverID, serverID)
}

func testAccBMSTaggingInterfaceCreateTemplate(interfaceID string) string {
	return fmt.Sprintf(`
# Create interface tags
resource "nsxt_policy_baremetal_server_interface_tags" "interface_tags_crud" {
  external_id = "%s"

  tag {
    scope = "test-net"
    tag   = "management"
  }

  tag {
    scope = "test-vlan"
    tag   = "100"
  }
}
`, interfaceID)
}

func testAccBMSTaggingInterfaceUpdateTemplate(interfaceID string) string {
	return fmt.Sprintf(`
# Update interface tags
resource "nsxt_policy_baremetal_server_interface_tags" "interface_tags_crud" {
  external_id = "%s"

  tag {
    scope = "test-net"
    tag   = "data-plane"  # Updated value
  }

  tag {
    scope = "test-vlan"
    tag   = "200"         # Updated value
  }

  tag {
    scope = "test-speed"  # New tag
    tag   = "10gb"
  }
}
`, interfaceID)
}

func testAccBMSTaggingInterfaceReadTemplate(interfaceID string) string {
	return fmt.Sprintf(`
# Update interface tags (same as previous step)
resource "nsxt_policy_baremetal_server_interface_tags" "interface_tags_crud" {
  external_id = "%s"

  tag {
    scope = "test-net"
    tag   = "data-plane"
  }

  tag {
    scope = "test-vlan"
    tag   = "200"
  }

  tag {
    scope = "test-speed"
    tag   = "10gb"
  }
}

# Read tags through data source
data "nsxt_policy_baremetal_server_interface_tags" "interface_tags_crud" {
  external_id = "%s"
}

# Output tag information
output "interface_tags_crud" {
  value = data.nsxt_policy_baremetal_server_interface_tags.interface_tags_crud.tag
}
`, interfaceID, interfaceID)
}

// Enhanced validation function for tagging tests
func testAccBMSTaggingEnhancedCheck(resourceName string, expectedTagCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if _, ok := s.RootModule().Resources[resourceName]; ok {
			checks := []resource.TestCheckFunc{
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "external_id"),
				resource.TestCheckResourceAttr(resourceName, "tag.#", strconv.Itoa(expectedTagCount)),
			}

			// Only validate tag structure if we expect tags to exist
			if expectedTagCount > 0 {
				checks = append(checks,
					resource.TestCheckResourceAttrSet(resourceName, "tag.0.scope"),
					resource.TestCheckResourceAttrSet(resourceName, "tag.0.tag"),
				)
			}

			return resource.ComposeTestCheckFunc(checks...)(s)
		}
		return nil
	}
}

// Conditional check functions removed - using direct environment variable-based testing
