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

var accTestVpcLBPassiveMonitorProfileCreateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform created",
	"max_fails":    "2",
	"timeout":      "2",
}

var accTestVpcLBPassiveMonitorProfileUpdateAttributes = map[string]string{
	"display_name": getAccTestResourceName(),
	"description":  "terraform updated",
	"max_fails":    "5",
	"timeout":      "5",
}

func TestAccResourceNsxtVpcLBPassiveMonitorProfile_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_passive_monitor_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.2") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBPassiveMonitorProfileCheckDestroy(state, accTestVpcLBPassiveMonitorProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcLBPassiveMonitorProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBPassiveMonitorProfileExists(accTestVpcLBPassiveMonitorProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcLBPassiveMonitorProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBPassiveMonitorProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "max_fails", accTestVpcLBPassiveMonitorProfileCreateAttributes["max_fails"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestVpcLBPassiveMonitorProfileCreateAttributes["timeout"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBPassiveMonitorProfileTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBPassiveMonitorProfileExists(accTestVpcLBPassiveMonitorProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcLBPassiveMonitorProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBPassiveMonitorProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "max_fails", accTestVpcLBPassiveMonitorProfileUpdateAttributes["max_fails"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestVpcLBPassiveMonitorProfileUpdateAttributes["timeout"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBPassiveMonitorProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBPassiveMonitorProfileExists(accTestVpcLBPassiveMonitorProfileCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtVpcLBPassiveMonitorProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_lb_passive_monitor_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.2") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBPassiveMonitorProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcLBPassiveMonitorProfileMinimalistic(),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceNsxtPolicyImportIDRetriever(testResourceName),
			},
		},
	})
}

func testAccNsxtVpcLBPassiveMonitorProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy LBPassiveMonitorProfile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy LBPassiveMonitorProfile resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyLBMonitorProfileExistsWrapperCtx(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy LBPassiveMonitorProfile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcLBPassiveMonitorProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_vpc_lb_passive_monitor_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyLBMonitorProfileExistsWrapperCtx(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy LBPassiveMonitorProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcLBPassiveMonitorProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcLBPassiveMonitorProfileCreateAttributes
	} else {
		attrMap = accTestVpcLBPassiveMonitorProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_passive_monitor_profile" "test" {
%s
  display_name = "%s"
  description  = "%s"
  max_fails    = %s
  timeout      = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], attrMap["description"], attrMap["max_fails"], attrMap["timeout"])
}

func testAccNsxtVpcLBPassiveMonitorProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_passive_monitor_profile" "test" {
%s
  display_name = "%s"

}`, testAccNsxtPolicyMultitenancyContext(), accTestVpcLBPassiveMonitorProfileUpdateAttributes["display_name"])
}
