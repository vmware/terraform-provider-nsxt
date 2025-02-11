// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var accTestPolicyLBFastUdpApplicationProfileCreateAttributes = map[string]string{
	"display_name":           getAccTestResourceName(),
	"description":            "terraform created",
	"idle_timeout":           "150",
	"flow_mirroring_enabled": "true",
}

var accTestPolicyLBFastUdpApplicationProfileUpdateAttributes = map[string]string{
	"display_name":           getAccTestResourceName(),
	"description":            "terraform updated",
	"idle_timeout":           "10",
	"flow_mirroring_enabled": "false",
}

func TestAccResourceNsxtPolicyLBFastUdpApplicationProfile_basic(t *testing.T) {
	testResourceName := "nsxt_policy_lb_fast_udp_application_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBFastUdpApplicationProfileCheckDestroy(state, accTestPolicyLBFastUdpApplicationProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBFastUdpApplicationProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBApplicationProfileExists(accTestPolicyLBFastUdpApplicationProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBFastUdpApplicationProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBFastUdpApplicationProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "idle_timeout", accTestPolicyLBFastUdpApplicationProfileCreateAttributes["idle_timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "flow_mirroring_enabled", accTestPolicyLBFastUdpApplicationProfileCreateAttributes["flow_mirroring_enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBFastUdpApplicationProfileTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBApplicationProfileExists(accTestPolicyLBFastUdpApplicationProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBFastUdpApplicationProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBFastUdpApplicationProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "idle_timeout", accTestPolicyLBFastUdpApplicationProfileUpdateAttributes["idle_timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "flow_mirroring_enabled", accTestPolicyLBFastUdpApplicationProfileUpdateAttributes["flow_mirroring_enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBFastUdpApplicationProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBApplicationProfileExists(accTestPolicyLBFastUdpApplicationProfileCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyLBFastUdpApplicationProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_lb_fast_udp_application_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBFastUdpApplicationProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBFastUdpApplicationProfileMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyLBFastUdpApplicationProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_lb_fast_udp_application_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyLBAppProfileExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy LB Fast Udp Profile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyLBFastUdpApplicationProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyLBFastUdpApplicationProfileCreateAttributes
	} else {
		attrMap = accTestPolicyLBFastUdpApplicationProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_lb_fast_udp_application_profile" "test" {
  display_name  = "%s"
  description   = "%s"
  idle_timeout  = %s

  flow_mirroring_enabled = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["idle_timeout"], attrMap["flow_mirroring_enabled"])
}

func testAccNsxtPolicyLBFastUdpApplicationProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_lb_fast_udp_application_profile" "test" {
  display_name = "%s"
}`, accTestPolicyLBFastUdpApplicationProfileUpdateAttributes["display_name"])
}
