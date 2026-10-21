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

var accTestVpcLBIcmpMonitorProfileCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
	"data_length":  "2",
	"fall_count":   "2",
	"interval":     "2",
	"rise_count":   "2",
	"timeout":      "2",
}

var accTestVpcLBIcmpMonitorProfileUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
	"data_length":  "5",
	"fall_count":   "5",
	"interval":     "5",
	"rise_count":   "5",
	"timeout":      "5",
}

func TestAccResourceNsxtVpcLBIcmpMonitorProfile_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_icmp_monitor_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBIcmpMonitorProfileCheckDestroy(state, accTestVpcLBIcmpMonitorProfileUpdateAttributes["display_name"])
		},
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBIcmpMonitorProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBIcmpMonitorProfileExists(accTestVpcLBIcmpMonitorProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcLBIcmpMonitorProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBIcmpMonitorProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "data_length", accTestVpcLBIcmpMonitorProfileCreateAttributes["data_length"]),
					resource.TestCheckResourceAttr(testResourceName, "fall_count", accTestVpcLBIcmpMonitorProfileCreateAttributes["fall_count"]),
					resource.TestCheckResourceAttr(testResourceName, "interval", accTestVpcLBIcmpMonitorProfileCreateAttributes["interval"]),
					resource.TestCheckResourceAttr(testResourceName, "rise_count", accTestVpcLBIcmpMonitorProfileCreateAttributes["rise_count"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestVpcLBIcmpMonitorProfileCreateAttributes["timeout"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBIcmpMonitorProfileTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBIcmpMonitorProfileExists(accTestVpcLBIcmpMonitorProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcLBIcmpMonitorProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBIcmpMonitorProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "data_length", accTestVpcLBIcmpMonitorProfileUpdateAttributes["data_length"]),
					resource.TestCheckResourceAttr(testResourceName, "fall_count", accTestVpcLBIcmpMonitorProfileUpdateAttributes["fall_count"]),
					resource.TestCheckResourceAttr(testResourceName, "interval", accTestVpcLBIcmpMonitorProfileUpdateAttributes["interval"]),
					resource.TestCheckResourceAttr(testResourceName, "rise_count", accTestVpcLBIcmpMonitorProfileUpdateAttributes["rise_count"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestVpcLBIcmpMonitorProfileUpdateAttributes["timeout"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBIcmpMonitorProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBIcmpMonitorProfileExists(accTestVpcLBIcmpMonitorProfileCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtVpcLBIcmpMonitorProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_lb_icmp_monitor_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBIcmpMonitorProfileCheckDestroy(state, name)
		},
		Steps: withImportIdempotencyChecks(withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBIcmpMonitorProfileMinimalistic(),
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

func testAccNsxtVpcLBIcmpMonitorProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy LBIcmpMonitorProfile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy LBIcmpMonitorProfile resource ID not set in resources")
		}

		exists, err := resourceNsxtVpcLBMonitorProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy LBIcmpMonitorProfile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcLBIcmpMonitorProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_vpc_lb_icmp_monitor_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcLBMonitorProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy LBIcmpMonitorProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcLBIcmpMonitorProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcLBIcmpMonitorProfileCreateAttributes
	} else {
		attrMap = accTestVpcLBIcmpMonitorProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_icmp_monitor_profile" "test" {
%s
  display_name = "%s"
  description  = "%s"
  data_length  = %s
  fall_count   = %s
  interval     = %s
  rise_count   = %s
  timeout      = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], attrMap["description"], attrMap["data_length"], attrMap["fall_count"], attrMap["interval"], attrMap["rise_count"], attrMap["timeout"])
}

func testAccNsxtVpcLBIcmpMonitorProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_icmp_monitor_profile" "test" {
%s
  display_name = "%s"

}`, testAccNsxtPolicyMultitenancyContext(), accTestVpcLBIcmpMonitorProfileUpdateAttributes["display_name"])
}
