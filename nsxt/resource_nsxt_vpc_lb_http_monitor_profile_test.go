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

var accTestVpcLBHttpMonitorProfileCreateAttributes = map[string]string{
	"display_name":          getAccTestResourceName(),
	"description":           "terraform created",
	"request_body":          "test-create",
	"request_method":        "HEAD",
	"request_url":           "test-create",
	"request_version":       "HTTP_VERSION_1_1",
	"response_body":         "test-create",
	"response_status_codes": "200",
	"fall_count":            "2",
	"interval":              "2",
	"monitor_port":          "8080",
	"rise_count":            "2",
	"timeout":               "2",
	"header_name":           "test-create",
	"header_value":          "test-create",
}

var accTestVpcLBHttpMonitorProfileUpdateAttributes = map[string]string{
	"display_name":          getAccTestResourceName(),
	"description":           "terraform updated",
	"request_body":          "test-update",
	"request_method":        "GET",
	"request_url":           "test-update",
	"request_version":       "HTTP_VERSION_1_0",
	"response_body":         "test-update",
	"response_status_codes": "200",
	"fall_count":            "5",
	"interval":              "5",
	"monitor_port":          "8090",
	"rise_count":            "5",
	"timeout":               "5",
	"header_name":           "test-update",
	"header_value":          "test-update",
}

func TestAccResourceNsxtVpcLBHttpMonitorProfile_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_http_monitor_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBHttpMonitorProfileCheckDestroy(state, accTestVpcLBHttpMonitorProfileUpdateAttributes["display_name"])
		},
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBHttpMonitorProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBHttpMonitorProfileExists(accTestVpcLBHttpMonitorProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcLBHttpMonitorProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBHttpMonitorProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "request_body", accTestVpcLBHttpMonitorProfileCreateAttributes["request_body"]),
					resource.TestCheckResourceAttr(testResourceName, "request_header.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "request_method", accTestVpcLBHttpMonitorProfileCreateAttributes["request_method"]),
					resource.TestCheckResourceAttr(testResourceName, "request_url", accTestVpcLBHttpMonitorProfileCreateAttributes["request_url"]),
					resource.TestCheckResourceAttr(testResourceName, "request_version", accTestVpcLBHttpMonitorProfileCreateAttributes["request_version"]),
					resource.TestCheckResourceAttr(testResourceName, "response_body", accTestVpcLBHttpMonitorProfileCreateAttributes["response_body"]),
					resource.TestCheckResourceAttr(testResourceName, "response_status_codes.0", accTestVpcLBHttpMonitorProfileCreateAttributes["response_status_codes"]),
					resource.TestCheckResourceAttr(testResourceName, "fall_count", accTestVpcLBHttpMonitorProfileCreateAttributes["fall_count"]),
					resource.TestCheckResourceAttr(testResourceName, "interval", accTestVpcLBHttpMonitorProfileCreateAttributes["interval"]),
					resource.TestCheckResourceAttr(testResourceName, "monitor_port", accTestVpcLBHttpMonitorProfileCreateAttributes["monitor_port"]),
					resource.TestCheckResourceAttr(testResourceName, "rise_count", accTestVpcLBHttpMonitorProfileCreateAttributes["rise_count"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestVpcLBHttpMonitorProfileCreateAttributes["timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "request_header.0.name", accTestVpcLBHttpMonitorProfileCreateAttributes["header_name"]),
					resource.TestCheckResourceAttr(testResourceName, "request_header.0.value", accTestVpcLBHttpMonitorProfileCreateAttributes["header_value"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBHttpMonitorProfileTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBHttpMonitorProfileExists(accTestVpcLBHttpMonitorProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcLBHttpMonitorProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBHttpMonitorProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "request_body", accTestVpcLBHttpMonitorProfileUpdateAttributes["request_body"]),
					resource.TestCheckResourceAttr(testResourceName, "request_header.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "request_method", accTestVpcLBHttpMonitorProfileUpdateAttributes["request_method"]),
					resource.TestCheckResourceAttr(testResourceName, "request_url", accTestVpcLBHttpMonitorProfileUpdateAttributes["request_url"]),
					resource.TestCheckResourceAttr(testResourceName, "request_version", accTestVpcLBHttpMonitorProfileUpdateAttributes["request_version"]),
					resource.TestCheckResourceAttr(testResourceName, "response_body", accTestVpcLBHttpMonitorProfileUpdateAttributes["response_body"]),
					resource.TestCheckResourceAttr(testResourceName, "response_status_codes.0", accTestVpcLBHttpMonitorProfileUpdateAttributes["response_status_codes"]),
					resource.TestCheckResourceAttr(testResourceName, "fall_count", accTestVpcLBHttpMonitorProfileUpdateAttributes["fall_count"]),
					resource.TestCheckResourceAttr(testResourceName, "interval", accTestVpcLBHttpMonitorProfileUpdateAttributes["interval"]),
					resource.TestCheckResourceAttr(testResourceName, "monitor_port", accTestVpcLBHttpMonitorProfileUpdateAttributes["monitor_port"]),
					resource.TestCheckResourceAttr(testResourceName, "rise_count", accTestVpcLBHttpMonitorProfileUpdateAttributes["rise_count"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestVpcLBHttpMonitorProfileUpdateAttributes["timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "request_header.0.name", accTestVpcLBHttpMonitorProfileUpdateAttributes["header_name"]),
					resource.TestCheckResourceAttr(testResourceName, "request_header.0.value", accTestVpcLBHttpMonitorProfileUpdateAttributes["header_value"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBHttpMonitorProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBHttpMonitorProfileExists(accTestVpcLBHttpMonitorProfileCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtVpcLBHttpMonitorProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_lb_http_monitor_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBHttpMonitorProfileCheckDestroy(state, name)
		},
		Steps: withImportIdempotencyChecks(withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBHttpMonitorProfileMinimalistic(),
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

func testAccNsxtVpcLBHttpMonitorProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy LBHttpMonitorProfile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy LBHttpMonitorProfile resource ID not set in resources")
		}

		exists, err := resourceNsxtVpcLBMonitorProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy LBHttpMonitorProfile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcLBHttpMonitorProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_vpc_lb_http_monitor_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcLBMonitorProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy LBHttpMonitorProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcLBHttpMonitorProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcLBHttpMonitorProfileCreateAttributes
	} else {
		attrMap = accTestVpcLBHttpMonitorProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_http_monitor_profile" "test" {
%s
  display_name = "%s"
  description  = "%s"
  request_body = "%s"

  request_method        = "%s"
  request_url           = "%s"
  request_version       = "%s"
  response_body         = "%s"
  response_status_codes = [%s]
  fall_count            = %s
  interval              = %s
  monitor_port          = %s
  rise_count            = %s
  timeout               = %s

  request_header {
    name  = "%s"
    value = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], attrMap["description"], attrMap["request_body"], attrMap["request_method"], attrMap["request_url"], attrMap["request_version"], attrMap["response_body"], attrMap["response_status_codes"], attrMap["fall_count"], attrMap["interval"], attrMap["monitor_port"], attrMap["rise_count"], attrMap["timeout"], attrMap["header_name"], attrMap["header_value"])
}

func testAccNsxtVpcLBHttpMonitorProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_http_monitor_profile" "test" {
%s
  display_name = "%s"

}`, testAccNsxtPolicyMultitenancyContext(), accTestVpcLBHttpMonitorProfileUpdateAttributes["display_name"])
}
