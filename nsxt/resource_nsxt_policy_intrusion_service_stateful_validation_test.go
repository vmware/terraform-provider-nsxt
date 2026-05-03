// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test that verifies stateful field is computed-only for DFW IDPS policies
func TestAccResourceNsxtPolicyIntrusionServicePolicy_statefulComputedOnly(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_intrusion_service_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServicePolicyStatefulTest(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServicePolicyExists(testResourceName, "default"),
					// Verify stateful is always true (computed)
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					// Verify other basic attributes
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Test stateful field behavior"),
				),
			},
		},
	})
}

// Test that verifies stateful field is computed-only for Gateway IDPS policies
func TestAccResourceNsxtPolicyIntrusionServiceGatewayPolicy_statefulComputedOnly(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_intrusion_service_gateway_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyIntrusionServiceGatewayPolicyStatefulTest(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServiceGatewayPolicyExists(testResourceName, "default"),
					// Verify stateful is always true (computed)
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					// Verify other basic attributes
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Test stateful field behavior for gateway"),
				),
			},
		},
	})
}

// Test that verifies stateful field is computed-only for Parent DFW IDPS policies
func TestAccResourceNsxtPolicyParentIntrusionServicePolicy_statefulComputedOnly(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_parent_intrusion_service_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentIntrusionServicePolicyStatefulTest(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServicePolicyExists(testResourceName, "default"),
					// Verify stateful is always true (computed)
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					// Verify other basic attributes
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Test stateful field behavior for parent"),
				),
			},
		},
	})
}

// Test that verifies stateful field is computed-only for Parent Gateway IDPS policies
func TestAccResourceNsxtPolicyParentIntrusionServiceGatewayPolicy_statefulComputedOnly(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_parent_intrusion_service_gateway_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "4.2.0")
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyParentIntrusionServiceGatewayPolicyStatefulTest(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyIntrusionServiceGatewayPolicyExists(testResourceName, "default"),
					// Verify stateful is always true (computed)
					resource.TestCheckResourceAttr(testResourceName, "stateful", "true"),
					// Verify other basic attributes
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Test stateful field behavior for parent gateway"),
				),
			},
		},
	})
}

// Test configurations without stateful field (should be computed as true)

func testAccNsxtPolicyIntrusionServicePolicyStatefulTest(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_intrusion_service_policy" "test" {
  display_name    = "%s"
  description     = "Test stateful field behavior"
  category        = "ThreatRules"
  locked          = false
  sequence_number = 15
  # stateful field is intentionally omitted - should be computed as true

  rule {
    display_name = "test-rule"
    action       = "DETECT"
    ids_profiles = [data.nsxt_policy_intrusion_service_profile.test.path]
  }
}`, name)
}

func testAccNsxtPolicyIntrusionServiceGatewayPolicyStatefulTest(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_intrusion_service_profile" "test" {
  display_name = "DefaultIDSProfile"
}

resource "nsxt_policy_tier1_gateway" "test" {
  display_name = "terraform-test-ids-gw-t1"
  description  = "T1 gateway for IDPS Gateway Policy stateful test"
}

resource "nsxt_policy_intrusion_service_gateway_policy" "test" {
  display_name    = "%s"
  description     = "Test stateful field behavior for gateway"
  category        = "LocalGatewayRules"
  locked          = false
  sequence_number = 16
  # stateful field is intentionally omitted - should be computed as true

  rule {
    display_name = "test-gw-rule"
    action       = "DETECT"
    scope        = [nsxt_policy_tier1_gateway.test.path]
    ids_profiles = [data.nsxt_policy_intrusion_service_profile.test.path]
  }
}`, name)
}

func testAccNsxtPolicyParentIntrusionServicePolicyStatefulTest(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_parent_intrusion_service_policy" "test" {
  display_name    = "%s"
  description     = "Test stateful field behavior for parent"
  category        = "ThreatRules"
  locked          = false
  sequence_number = 17
  # stateful field is intentionally omitted - should be computed as true
}`, name)
}

func testAccNsxtPolicyParentIntrusionServiceGatewayPolicyStatefulTest(name string) string {
	return fmt.Sprintf(`
resource "nsxt_policy_parent_intrusion_service_gateway_policy" "test" {
  display_name    = "%s"
  description     = "Test stateful field behavior for parent gateway"
  category        = "SharedPreRules"
  locked          = false
  sequence_number = 18
  # stateful field is intentionally omitted - should be computed as true
}`, name)
}
