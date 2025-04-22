/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var accTestPolicyEdgeHighAvailabilityProfileCreateAttributes = map[string]string{
	"display_name":                 getAccTestResourceName(),
	"description":                  "terraform created",
	"bfd_probe_interval":           "20000",
	"bfd_allowed_hops":             "140",
	"bfd_declare_dead_multiple":    "7",
	"standby_relocation_threshold": "2000",
}

var accTestPolicyEdgeHighAvailabilityProfileUpdateAttributes = map[string]string{
	"display_name":                 getAccTestResourceName(),
	"description":                  "terraform updated",
	"bfd_probe_interval":           "25000",
	"bfd_allowed_hops":             "122",
	"bfd_declare_dead_multiple":    "9",
	"standby_relocation_threshold": "5000",
}

func TestAccResourceNsxtPolicyEdgeHighAvailabilityProfile_basic(t *testing.T) {
	testResourceName := "nsxt_policy_edge_high_availability_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyEdgeHighAvailabilityProfileCheckDestroy(state, accTestPolicyEdgeHighAvailabilityProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyEdgeHighAvailabilityProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyEdgeHighAvailabilityProfileExists(accTestPolicyEdgeHighAvailabilityProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyEdgeHighAvailabilityProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyEdgeHighAvailabilityProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "bfd_probe_interval", accTestPolicyEdgeHighAvailabilityProfileCreateAttributes["bfd_probe_interval"]),
					resource.TestCheckResourceAttr(testResourceName, "bfd_allowed_hops", accTestPolicyEdgeHighAvailabilityProfileCreateAttributes["bfd_allowed_hops"]),
					resource.TestCheckResourceAttr(testResourceName, "bfd_declare_dead_multiple", accTestPolicyEdgeHighAvailabilityProfileCreateAttributes["bfd_declare_dead_multiple"]),
					resource.TestCheckResourceAttr(testResourceName, "standby_relocation_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "standby_relocation_config.0.standby_relocation_threshold", accTestPolicyEdgeHighAvailabilityProfileCreateAttributes["standby_relocation_threshold"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyEdgeHighAvailabilityProfileTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyEdgeHighAvailabilityProfileExists(accTestPolicyEdgeHighAvailabilityProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyEdgeHighAvailabilityProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyEdgeHighAvailabilityProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "bfd_probe_interval", accTestPolicyEdgeHighAvailabilityProfileUpdateAttributes["bfd_probe_interval"]),
					resource.TestCheckResourceAttr(testResourceName, "bfd_allowed_hops", accTestPolicyEdgeHighAvailabilityProfileUpdateAttributes["bfd_allowed_hops"]),
					resource.TestCheckResourceAttr(testResourceName, "bfd_declare_dead_multiple", accTestPolicyEdgeHighAvailabilityProfileUpdateAttributes["bfd_declare_dead_multiple"]),
					resource.TestCheckResourceAttr(testResourceName, "standby_relocation_config.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "standby_relocation_config.0.standby_relocation_threshold", accTestPolicyEdgeHighAvailabilityProfileUpdateAttributes["standby_relocation_threshold"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyEdgeHighAvailabilityProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyEdgeHighAvailabilityProfileExists(accTestPolicyEdgeHighAvailabilityProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyEdgeHighAvailabilityProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_edge_high_availability_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyLocalManager(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyEdgeHighAvailabilityProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyEdgeHighAvailabilityProfileMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNsxtPolicyEdgeHighAvailabilityProfileImporterGetID,
			},
		},
	})
}

func testAccNsxtPolicyEdgeHighAvailabilityProfileImporterGetID(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources["nsxt_policy_edge_high_availability_profile.test"]
	if !ok {
		return "", fmt.Errorf("PolicyEdgeHighAvailabilityProfile resource %s not found in resources", "nsxt_policy_edge_high_availability_profile.test")
	}
	resourceID := rs.Primary.ID
	if resourceID == "" {
		return "", fmt.Errorf("PolicyEdgeHighAvailabilityProfile resource ID not set in resources ")
	}
	path := rs.Primary.Attributes["path"]
	if path == "" {
		return "", fmt.Errorf("PolicyEdgeHighAvailabilityProfile path not set in resources ")
	}
	return path, nil
}

func testAccNsxtPolicyEdgeHighAvailabilityProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy PolicyEdgeHighAvailabilityProfile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy PolicyEdgeHighAvailabilityProfile resource ID not set in resources")
		}
		epID := rs.Primary.Attributes["enforcement_point"]
		sitePath := rs.Primary.Attributes["site_path"]
		siteID := getPolicyIDFromPath(sitePath)
		exists, err := resourceNsxtPolicyEdgeHighAvailabilityProfileExists(siteID, epID, resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy PolicyEdgeHighAvailabilityProfile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtPolicyEdgeHighAvailabilityProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_edge_high_availability_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		epID := rs.Primary.Attributes["enforcement_point"]
		sitePath := rs.Primary.Attributes["site_path"]
		siteID := getPolicyIDFromPath(sitePath)
		exists, err := resourceNsxtPolicyEdgeHighAvailabilityProfileExists(siteID, epID, resourceID, connector)
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy PolicyEdgeHighAvailabilityProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyEdgeHighAvailabilityProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyEdgeHighAvailabilityProfileCreateAttributes
	} else {
		attrMap = accTestPolicyEdgeHighAvailabilityProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_edge_high_availability_profile" "test" {
  display_name              = "%s"
  description               = "%s"
  bfd_probe_interval        = %s
  bfd_allowed_hops          = %s
  bfd_declare_dead_multiple = %s
  standby_relocation_config {
    standby_relocation_threshold = %s
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["bfd_probe_interval"], attrMap["bfd_allowed_hops"], attrMap["bfd_declare_dead_multiple"], attrMap["standby_relocation_threshold"])
}

func testAccNsxtPolicyEdgeHighAvailabilityProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_edge_high_availability_profile" "test" {
  display_name = "%s"
}`, accTestPolicyEdgeHighAvailabilityProfileUpdateAttributes["display_name"])
}
