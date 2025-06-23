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

var accTestPolicyLBFastTcpApplicationProfileCreateAttributes = map[string]string{
	"display_name":              getAccTestResourceName(),
	"description":               "terraform created",
	"idle_timeout":              "150",
	"close_timeout":             "10",
	"ha_flow_mirroring_enabled": "true",
}

var accTestPolicyLBFastTcpApplicationProfileUpdateAttributes = map[string]string{
	"display_name":              getAccTestResourceName(),
	"description":               "terraform updated",
	"idle_timeout":              "10",
	"close_timeout":             "15",
	"ha_flow_mirroring_enabled": "false",
}

func TestAccResourceNsxtPolicyLBFastTcpApplicationProfile_basic(t *testing.T) {
	testResourceName := "nsxt_policy_lb_fast_tcp_application_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBFastTcpApplicationProfileCheckDestroy(state, accTestPolicyLBFastTcpApplicationProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBFastTcpApplicationProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBApplicationProfileExists(accTestPolicyLBFastTcpApplicationProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBFastTcpApplicationProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBFastTcpApplicationProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "idle_timeout", accTestPolicyLBFastTcpApplicationProfileCreateAttributes["idle_timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "close_timeout", accTestPolicyLBFastTcpApplicationProfileCreateAttributes["close_timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_flow_mirroring_enabled", accTestPolicyLBFastTcpApplicationProfileCreateAttributes["ha_flow_mirroring_enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBFastTcpApplicationProfileTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBApplicationProfileExists(accTestPolicyLBFastTcpApplicationProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestPolicyLBFastTcpApplicationProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestPolicyLBFastTcpApplicationProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "idle_timeout", accTestPolicyLBFastTcpApplicationProfileUpdateAttributes["idle_timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "close_timeout", accTestPolicyLBFastTcpApplicationProfileUpdateAttributes["close_timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_flow_mirroring_enabled", accTestPolicyLBFastTcpApplicationProfileUpdateAttributes["ha_flow_mirroring_enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtPolicyLBFastTcpApplicationProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyLBApplicationProfileExists(accTestPolicyLBFastTcpApplicationProfileCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtPolicyLBFastTcpApplicationProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_policy_lb_fast_tcp_application_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyLocalManager(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyLBFastTcpApplicationProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyLBFastTcpApplicationProfileMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNsxtPolicyLBFastTcpApplicationProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_lb_fast_tcp_application_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyLBAppProfileExists(resourceID, connector, testAccIsGlobalManager())
		if err == nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy LB Fast Tcp Profile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyLBFastTcpApplicationProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestPolicyLBFastTcpApplicationProfileCreateAttributes
	} else {
		attrMap = accTestPolicyLBFastTcpApplicationProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_policy_lb_fast_tcp_application_profile" "test" {
  display_name  = "%s"
  description   = "%s"
  close_timeout = %s
  idle_timeout  = %s

  ha_flow_mirroring_enabled = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, attrMap["display_name"], attrMap["description"], attrMap["close_timeout"], attrMap["idle_timeout"], attrMap["ha_flow_mirroring_enabled"])
}

func testAccNsxtPolicyLBFastTcpApplicationProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_policy_lb_fast_tcp_application_profile" "test" {
  display_name = "%s"
}`, accTestPolicyLBFastTcpApplicationProfileUpdateAttributes["display_name"])
}
