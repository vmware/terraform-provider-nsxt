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

var accTestVpcLBFastTcpApplicationProfileCreateAttributes = map[string]string{
	"display_name":              getAccTestResourceName(),
	"description":               "terraform created",
	"idle_timeout":              "150",
	"close_timeout":             "10",
	"ha_flow_mirroring_enabled": "true",
}

var accTestVpcLBFastTcpApplicationProfileUpdateAttributes = map[string]string{
	"display_name":              getAccTestResourceName(),
	"description":               "terraform updated",
	"idle_timeout":              "10",
	"close_timeout":             "15",
	"ha_flow_mirroring_enabled": "false",
}

func TestAccResourceNsxtVpcLBFastTcpApplicationProfile_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_fast_tcp_application_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBFastTcpApplicationProfileCheckDestroy(state, accTestVpcLBFastTcpApplicationProfileUpdateAttributes["display_name"])
		},
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBFastTcpApplicationProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBApplicationProfileExists(accTestVpcLBFastTcpApplicationProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcLBFastTcpApplicationProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBFastTcpApplicationProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "idle_timeout", accTestVpcLBFastTcpApplicationProfileCreateAttributes["idle_timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "close_timeout", accTestVpcLBFastTcpApplicationProfileCreateAttributes["close_timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_flow_mirroring_enabled", accTestVpcLBFastTcpApplicationProfileCreateAttributes["ha_flow_mirroring_enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBFastTcpApplicationProfileTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBApplicationProfileExists(accTestVpcLBFastTcpApplicationProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcLBFastTcpApplicationProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBFastTcpApplicationProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "idle_timeout", accTestVpcLBFastTcpApplicationProfileUpdateAttributes["idle_timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "close_timeout", accTestVpcLBFastTcpApplicationProfileUpdateAttributes["close_timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "ha_flow_mirroring_enabled", accTestVpcLBFastTcpApplicationProfileUpdateAttributes["ha_flow_mirroring_enabled"]),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBFastTcpApplicationProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBApplicationProfileExists(accTestVpcLBFastTcpApplicationProfileCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtVpcLBFastTcpApplicationProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_lb_fast_tcp_application_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBFastTcpApplicationProfileCheckDestroy(state, name)
		},
		Steps: withImportIdempotencyChecks(withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBFastTcpApplicationProfileMinimalistic(),
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

func testAccNsxtVpcLBFastTcpApplicationProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_vpc_lb_fast_tcp_application_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcLBAppProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy LB Fast Tcp Profile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcLBFastTcpApplicationProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcLBFastTcpApplicationProfileCreateAttributes
	} else {
		attrMap = accTestVpcLBFastTcpApplicationProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_fast_tcp_application_profile" "test" {
%s
  display_name  = "%s"
  description   = "%s"
  close_timeout = %s
  idle_timeout  = %s

  ha_flow_mirroring_enabled = %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], attrMap["description"], attrMap["close_timeout"], attrMap["idle_timeout"], attrMap["ha_flow_mirroring_enabled"])
}

func testAccNsxtVpcLBFastTcpApplicationProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_fast_tcp_application_profile" "test" {
%s
  display_name = "%s"
}`, testAccNsxtPolicyMultitenancyContext(), accTestVpcLBFastTcpApplicationProfileUpdateAttributes["display_name"])
}
