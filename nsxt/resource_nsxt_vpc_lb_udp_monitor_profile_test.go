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

var accTestVpcLBUdpMonitorProfileCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
	"receive":      "test-create",
	"send":         "test-create",
	"fall_count":   "2",
	"interval":     "2",
	"monitor_port": "8080",
	"rise_count":   "2",
	"timeout":      "2",
}

var accTestVpcLBUdpMonitorProfileUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
	"receive":      "test-update",
	"send":         "test-update",
	"fall_count":   "5",
	"interval":     "5",
	"monitor_port": "8090",
	"rise_count":   "5",
	"timeout":      "5",
}

func TestAccResourceNsxtVpcLBUdpMonitorProfile_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_udp_monitor_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBUdpMonitorProfileCheckDestroy(state, accTestVpcLBUdpMonitorProfileUpdateAttributes["display_name"])
		},
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBUdpMonitorProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBUdpMonitorProfileExists(accTestVpcLBUdpMonitorProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcLBUdpMonitorProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBUdpMonitorProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "receive", accTestVpcLBUdpMonitorProfileCreateAttributes["receive"]),
					resource.TestCheckResourceAttr(testResourceName, "send", accTestVpcLBUdpMonitorProfileCreateAttributes["send"]),
					resource.TestCheckResourceAttr(testResourceName, "fall_count", accTestVpcLBUdpMonitorProfileCreateAttributes["fall_count"]),
					resource.TestCheckResourceAttr(testResourceName, "interval", accTestVpcLBUdpMonitorProfileCreateAttributes["interval"]),
					resource.TestCheckResourceAttr(testResourceName, "monitor_port", accTestVpcLBUdpMonitorProfileCreateAttributes["monitor_port"]),
					resource.TestCheckResourceAttr(testResourceName, "rise_count", accTestVpcLBUdpMonitorProfileCreateAttributes["rise_count"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestVpcLBUdpMonitorProfileCreateAttributes["timeout"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBUdpMonitorProfileTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBUdpMonitorProfileExists(accTestVpcLBUdpMonitorProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcLBUdpMonitorProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBUdpMonitorProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "receive", accTestVpcLBUdpMonitorProfileUpdateAttributes["receive"]),
					resource.TestCheckResourceAttr(testResourceName, "send", accTestVpcLBUdpMonitorProfileUpdateAttributes["send"]),
					resource.TestCheckResourceAttr(testResourceName, "fall_count", accTestVpcLBUdpMonitorProfileUpdateAttributes["fall_count"]),
					resource.TestCheckResourceAttr(testResourceName, "interval", accTestVpcLBUdpMonitorProfileUpdateAttributes["interval"]),
					resource.TestCheckResourceAttr(testResourceName, "monitor_port", accTestVpcLBUdpMonitorProfileUpdateAttributes["monitor_port"]),
					resource.TestCheckResourceAttr(testResourceName, "rise_count", accTestVpcLBUdpMonitorProfileUpdateAttributes["rise_count"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestVpcLBUdpMonitorProfileUpdateAttributes["timeout"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBUdpMonitorProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBUdpMonitorProfileExists(accTestVpcLBUdpMonitorProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		}),
	})
}

func TestAccResourceNsxtVpcLBUdpMonitorProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_lb_udp_monitor_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBUdpMonitorProfileCheckDestroy(state, name)
		},
		Steps: withImportIdempotencyChecks(withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBUdpMonitorProfileMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
		})),
	})
}

func testAccNsxtVpcLBUdpMonitorProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy LBUdpMonitorProfile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy LBUdpMonitorProfile resource ID not set in resources")
		}

		exists, err := resourceNsxtVpcLBMonitorProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy LBUdpMonitorProfile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcLBUdpMonitorProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_vpc_lb_udp_monitor_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcLBMonitorProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy LBUdpMonitorProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcLBUdpMonitorProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcLBUdpMonitorProfileCreateAttributes
	} else {
		attrMap = accTestVpcLBUdpMonitorProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_udp_monitor_profile" "test" {
%s
  display_name = "%s"
  description  = "%s"
  receive      = "%s"
  send         = "%s"
  fall_count   = %s
  interval     = %s
  monitor_port = %s
  rise_count   = %s
  timeout      = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], attrMap["description"], attrMap["receive"], attrMap["send"], attrMap["fall_count"], attrMap["interval"], attrMap["monitor_port"], attrMap["rise_count"], attrMap["timeout"])
}

func testAccNsxtVpcLBUdpMonitorProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_udp_monitor_profile" "test" {
%s
  display_name = "%s"
  receive = "%s"
  send = "%s"

}`, testAccNsxtPolicyMultitenancyContext(), accTestVpcLBUdpMonitorProfileUpdateAttributes["display_name"], accTestVpcLBUdpMonitorProfileUpdateAttributes["receive"], accTestVpcLBUdpMonitorProfileUpdateAttributes["send"])
}
