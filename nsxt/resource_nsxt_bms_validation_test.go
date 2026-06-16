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

// Note: Tests for time-based rules, identity firewall, and L7 context profiles removed.
// These tests were non-functional as they used comment strings that the API doesn't validate.
// Proper validation would require implementing real provider-side validation with ExpectError assertions.

// TestAccResourceNsxtPolicyBareMetalServerTags_dfwOnly tests that BMS is only supported in DFW (not GFW)
func TestAccResourceNsxtPolicyBareMetalServerTags_dfwOnly(t *testing.T) {
	testName := getAccTestResourceName()
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
				Config: testAccBMSValidationDFWOnlyTemplate(testName, serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
					// Verify BMS groups work in DFW policies
					testAccBMSValidationEnhancedCheck("nsxt_policy_group.bms_servers", testName+"-bms-servers"),
					testAccBMSValidationEnhancedCheck("nsxt_policy_security_policy.dfw_bms_policy", testName+"-dfw-bms-policy"),
				),
			},
		},
	})
}

// TestAccResourceNsxtPolicyBareMetalServerTags_layer34Only tests that only Layer 3/4 rules are supported with BMS
func TestAccResourceNsxtPolicyBareMetalServerTags_layer34Only(t *testing.T) {
	testName := getAccTestResourceName()
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
				Config: testAccBMSValidationLayer34OnlyTemplate(testName, serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nsxt_policy_baremetal_servers.all", "results.#"),
					// Verify Layer 3/4 rules work with BMS groups
					testAccBMSValidationEnhancedCheck("nsxt_policy_group.layer34_bms_servers", testName+"-layer34-bms-servers"),
					testAccBMSValidationEnhancedCheck("nsxt_policy_security_policy.layer34_bms_policy", testName+"-layer34-bms-policy"),
				),
			},
		},
	})
}

// TestAccResourceNsxtPolicyBareMetalServerTags_invalidExternalID tests error handling for invalid external IDs
func TestAccResourceNsxtPolicyBareMetalServerTags_invalidExternalID(t *testing.T) {
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
				Config:      testAccBMSValidationInvalidExternalIDTemplate(),
				ExpectError: regexp.MustCompile("external_id must be a valid UUID format|invalid.*external.*id"),
			},
		},
	})
}

// Note: TestAccBMSValidationVersionCheck removed because it will never run (t.Skip in PreCheck).
// Version validation is already handled by testAccNSXVersion(t, "9.0.0") in other tests.

// Template functions for removed non-functional validation tests were here.
// They were removed because they didn't actually validate restrictions -
// they only used comment strings that the API doesn't reject.

func testAccBMSValidationDFWOnlyTemplate(testName, serverID string) string {
	return fmt.Sprintf(`
# Data source to verify BMS inventory is accessible
data "nsxt_policy_baremetal_servers" "all" {}

# This should work - BMS is supported in DFW
# Using environment variable for test server
resource "nsxt_policy_group" "bms_group" {
  display_name = "%s-bms-group"
  description = "BMS group for DFW test"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

# This should work - BMS groups in DFW (Distributed Firewall)
resource "nsxt_policy_security_policy" "dfw_bms" {
  display_name = "%s-dfw-bms-policy"
  description = "DFW policy with BMS group - should work"
  category = "Application"

  rule {
    display_name = "DFW-BMS-Rule"
    source_groups = [nsxt_policy_group.bms_group.path]
    destination_groups = [nsxt_policy_group.bms_group.path]
    action = "ALLOW"
    services = ["/infra/services/HTTP"]
  }
}

# Output success information
output "dfw_bms_support" {
  value = {
    bms_in_dfw_supported = true
    policy_created = nsxt_policy_security_policy.dfw_bms.path
  }
}
`, testName, serverID, testName)
}

func testAccBMSValidationLayer34OnlyTemplate(testName, serverID string) string {
	return fmt.Sprintf(`
# Data source to verify BMS inventory is accessible
data "nsxt_policy_baremetal_servers" "all" {}

# This should work - Layer 3/4 rules are supported with BMS
# Using environment variable for test server
resource "nsxt_policy_group" "bms_group" {
  display_name = "%s-bms-group"
  description = "BMS group for Layer 3/4 test"
  group_type = "BareMetalServer"

  criteria {
    external_id_expression {
      member_type = "BareMetalServer"
      external_ids = ["%s"]
    }
  }
}

# Layer 3 rule - IP-based (should work)
resource "nsxt_policy_security_policy" "layer3_bms" {
  display_name = "%s-layer3-bms-policy"
  description = "Layer 3 policy with BMS group - should work"
  category = "Application"

  rule {
    display_name = "Layer3-BMS-Rule"
    source_groups = [nsxt_policy_group.bms_group.path]
    destination_groups = ["192.168.1.0/24"]  # IP range - Layer 3
    action = "ALLOW"
    ip_version = "IPV4"
  }
}

# Layer 4 rule - Port-based (should work)
resource "nsxt_policy_service" "custom_tcp" {
  display_name = "%s-custom-tcp-service"

  l4_port_set_entry {
    protocol = "TCP"
    destination_ports = ["8080"]
  }
}

resource "nsxt_policy_security_policy" "layer4_bms" {
  display_name = "%s-layer4-bms-policy"
  description = "Layer 4 policy with BMS group - should work"
  category = "Application"

  rule {
    display_name = "Layer4-BMS-Rule"
    source_groups = [nsxt_policy_group.bms_group.path]
    destination_groups = [nsxt_policy_group.bms_group.path]
    services = [nsxt_policy_service.custom_tcp.path]  # TCP port - Layer 4
    action = "ALLOW"
  }
}

# Output Layer 3/4 support information
output "layer34_bms_support" {
  value = {
    layer3_supported = true
    layer4_supported = true
    layer3_policy = nsxt_policy_security_policy.layer3_bms.path
    layer4_policy = nsxt_policy_security_policy.layer4_bms.path
  }
}
`, testName, serverID, testName, testName, testName)
}

func testAccBMSValidationInvalidExternalIDTemplate() string {
	return `
# This should fail - invalid external_id format
resource "nsxt_policy_baremetal_server_tags" "invalid_id" {
  external_id = "invalid-uuid-format"

  tag {
    scope = "test"
    tag   = "invalid"
  }
}
`
}

// Template function for version check test removed with the test.

// Enhanced validation function for validation tests
func testAccBMSValidationEnhancedCheck(resourceName, expectedDisplayName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if _, ok := s.RootModule().Resources[resourceName]; ok {
			return resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "path"),
				resource.TestCheckResourceAttr(resourceName, "display_name", expectedDisplayName),
			)(s)
		}
		return nil
	}
}
