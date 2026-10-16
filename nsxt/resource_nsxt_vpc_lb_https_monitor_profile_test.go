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

var accTestVpcLBHttpsMonitorProfileCreateAttributes = map[string]string{
	"display_name":            getAccTestResourceName(),
	"description":             "terraform created",
	"request_body":            "test-create",
	"request_method":          "POST",
	"request_url":             "test-create",
	"request_version":         "HTTP_VERSION_1_1",
	"response_body":           "test-create",
	"response_status_codes":   "200",
	"fall_count":              "2",
	"interval":                "2",
	"monitor_port":            "8080",
	"rise_count":              "2",
	"timeout":                 "2",
	"header_name":             "test-create",
	"header_value":            "test-create",
	"certificate_chain_depth": "2",
	"server_auth":             "IGNORE",
}

var accTestVpcLBHttpsMonitorProfileUpdateAttributes = map[string]string{
	"display_name":            getAccTestResourceName(),
	"description":             "terraform updated",
	"request_body":            "test-update",
	"request_method":          "OPTIONS",
	"request_url":             "test-update",
	"request_version":         "HTTP_VERSION_1_0",
	"response_body":           "test-update",
	"response_status_codes":   "200",
	"fall_count":              "5",
	"interval":                "5",
	"monitor_port":            "8090",
	"rise_count":              "5",
	"timeout":                 "5",
	"header_name":             "test-update",
	"header_value":            "test-update",
	"certificate_chain_depth": "5",
	"server_auth":             "IGNORE",
}

func TestAccResourceNsxtVpcLBHttpsMonitorProfile_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_https_monitor_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.2") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBHttpsMonitorProfileCheckDestroy(state, accTestVpcLBHttpsMonitorProfileUpdateAttributes["display_name"])
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcLBHttpsMonitorProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBHttpsMonitorProfileExists(accTestVpcLBHttpsMonitorProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcLBHttpsMonitorProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBHttpsMonitorProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "request_body", accTestVpcLBHttpsMonitorProfileCreateAttributes["request_body"]),
					resource.TestCheckResourceAttr(testResourceName, "request_header.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "request_method", accTestVpcLBHttpsMonitorProfileCreateAttributes["request_method"]),
					resource.TestCheckResourceAttr(testResourceName, "request_url", accTestVpcLBHttpsMonitorProfileCreateAttributes["request_url"]),
					resource.TestCheckResourceAttr(testResourceName, "request_version", accTestVpcLBHttpsMonitorProfileCreateAttributes["request_version"]),
					resource.TestCheckResourceAttr(testResourceName, "response_body", accTestVpcLBHttpsMonitorProfileCreateAttributes["response_body"]),
					resource.TestCheckResourceAttr(testResourceName, "response_status_codes.0", accTestVpcLBHttpsMonitorProfileCreateAttributes["response_status_codes"]),
					resource.TestCheckResourceAttr(testResourceName, "server_ssl.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "fall_count", accTestVpcLBHttpsMonitorProfileCreateAttributes["fall_count"]),
					resource.TestCheckResourceAttr(testResourceName, "interval", accTestVpcLBHttpsMonitorProfileCreateAttributes["interval"]),
					resource.TestCheckResourceAttr(testResourceName, "monitor_port", accTestVpcLBHttpsMonitorProfileCreateAttributes["monitor_port"]),
					resource.TestCheckResourceAttr(testResourceName, "rise_count", accTestVpcLBHttpsMonitorProfileCreateAttributes["rise_count"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestVpcLBHttpsMonitorProfileCreateAttributes["timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "request_header.0.name", accTestVpcLBHttpsMonitorProfileCreateAttributes["header_name"]),
					resource.TestCheckResourceAttr(testResourceName, "request_header.0.value", accTestVpcLBHttpsMonitorProfileCreateAttributes["header_value"]),
					resource.TestCheckResourceAttr(testResourceName, "server_ssl.0.certificate_chain_depth", accTestVpcLBHttpsMonitorProfileCreateAttributes["certificate_chain_depth"]),
					resource.TestCheckResourceAttr(testResourceName, "server_ssl.0.server_auth", accTestVpcLBHttpsMonitorProfileCreateAttributes["server_auth"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBHttpsMonitorProfileTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBHttpsMonitorProfileExists(accTestVpcLBHttpsMonitorProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcLBHttpsMonitorProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBHttpsMonitorProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "request_body", accTestVpcLBHttpsMonitorProfileUpdateAttributes["request_body"]),
					resource.TestCheckResourceAttr(testResourceName, "request_header.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "request_method", accTestVpcLBHttpsMonitorProfileUpdateAttributes["request_method"]),
					resource.TestCheckResourceAttr(testResourceName, "request_url", accTestVpcLBHttpsMonitorProfileUpdateAttributes["request_url"]),
					resource.TestCheckResourceAttr(testResourceName, "request_version", accTestVpcLBHttpsMonitorProfileUpdateAttributes["request_version"]),
					resource.TestCheckResourceAttr(testResourceName, "response_body", accTestVpcLBHttpsMonitorProfileUpdateAttributes["response_body"]),
					resource.TestCheckResourceAttr(testResourceName, "response_status_codes.0", accTestVpcLBHttpsMonitorProfileUpdateAttributes["response_status_codes"]),
					resource.TestCheckResourceAttr(testResourceName, "server_ssl.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "fall_count", accTestVpcLBHttpsMonitorProfileUpdateAttributes["fall_count"]),
					resource.TestCheckResourceAttr(testResourceName, "interval", accTestVpcLBHttpsMonitorProfileUpdateAttributes["interval"]),
					resource.TestCheckResourceAttr(testResourceName, "monitor_port", accTestVpcLBHttpsMonitorProfileUpdateAttributes["monitor_port"]),
					resource.TestCheckResourceAttr(testResourceName, "rise_count", accTestVpcLBHttpsMonitorProfileUpdateAttributes["rise_count"]),
					resource.TestCheckResourceAttr(testResourceName, "timeout", accTestVpcLBHttpsMonitorProfileUpdateAttributes["timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "request_header.0.name", accTestVpcLBHttpsMonitorProfileUpdateAttributes["header_name"]),
					resource.TestCheckResourceAttr(testResourceName, "request_header.0.value", accTestVpcLBHttpsMonitorProfileUpdateAttributes["header_value"]),
					resource.TestCheckResourceAttr(testResourceName, "server_ssl.0.certificate_chain_depth", accTestVpcLBHttpsMonitorProfileUpdateAttributes["certificate_chain_depth"]),
					resource.TestCheckResourceAttr(testResourceName, "server_ssl.0.server_auth", accTestVpcLBHttpsMonitorProfileUpdateAttributes["server_auth"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBHttpsMonitorProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBHttpsMonitorProfileExists(accTestVpcLBHttpsMonitorProfileCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtVpcLBHttpsMonitorProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_lb_https_monitor_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.2") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBHttpsMonitorProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtVpcLBHttpsMonitorProfileMinimalistic(),
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

func testAccNsxtVpcLBHttpsMonitorProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy LBHttpsMonitorProfile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy LBHttpsMonitorProfile resource ID not set in resources")
		}

		exists, err := resourceNsxtPolicyLBMonitorProfileExistsWrapperCtx(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy LBHttpsMonitorProfile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcLBHttpsMonitorProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_vpc_lb_https_monitor_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtPolicyLBMonitorProfileExistsWrapperCtx(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy LBHttpsMonitorProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcLBHttpsMonitorProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcLBHttpsMonitorProfileCreateAttributes
	} else {
		attrMap = accTestVpcLBHttpsMonitorProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_https_monitor_profile" "test" {
%s
  display_name = "%s"
  description  = "%s"
  request_body = "%s"


  request_method        = "%s"
  request_url           = "%s"
  request_version       = "%s"
  response_body         = "%s"
  response_status_codes = [%s]

  request_header {
    name  = "%s"
    value = "%s"
  }

  fall_count   = %s
  interval     = %s
  monitor_port = %s
  rise_count   = %s
  timeout      = %s

  server_ssl {
    certificate_chain_depth = %s
    server_auth             = "%s"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], attrMap["description"], attrMap["request_body"], attrMap["request_method"], attrMap["request_url"], attrMap["request_version"], attrMap["response_body"], attrMap["response_status_codes"], attrMap["header_name"], attrMap["header_value"], attrMap["fall_count"], attrMap["interval"], attrMap["monitor_port"], attrMap["rise_count"], attrMap["timeout"], attrMap["certificate_chain_depth"], attrMap["server_auth"])
}

func testAccNsxtVpcLBHttpsMonitorProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_https_monitor_profile" "test" {
%s
  display_name = "%s"

}`, testAccNsxtPolicyMultitenancyContext(), accTestVpcLBHttpsMonitorProfileUpdateAttributes["display_name"])
}
