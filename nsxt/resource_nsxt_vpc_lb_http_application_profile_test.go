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

var accTestVpcLBHttpApplicationProfileCreateAttributes = map[string]string{
	"display_name":           getAccTestResourceName(),
	"description":            "terraform created",
	"http_redirect_to":       "http://www.test-create.com",
	"http_redirect_to_https": "false",
	"idle_timeout":           "15",
	"request_body_size":      "256",
	"request_header_size":    "1024",
	"response_buffering":     "true",
	"response_header_size":   "2048",
	"response_timeout":       "20",
	"server_keep_alive":      "true",
	"x_forwarded_for":        "REPLACE",
}

var accTestVpcLBHttpApplicationProfileUpdateAttributes = map[string]string{
	"display_name":           getAccTestResourceName(),
	"description":            "terraform updated",
	"http_redirect_to":       "http://www.test-update.com",
	"http_redirect_to_https": "false",
	"idle_timeout":           "20",
	"request_body_size":      "512",
	"request_header_size":    "2048",
	"response_buffering":     "false",
	"response_header_size":   "4096",
	"response_timeout":       "10",
	"server_keep_alive":      "false",
	"x_forwarded_for":        "INSERT",
}

func TestAccResourceNsxtVpcLBHttpApplicationProfile_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_http_application_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBHttpApplicationProfileCheckDestroy(state, accTestVpcLBHttpApplicationProfileUpdateAttributes["display_name"])
		},
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBHttpApplicationProfileTemplate(true),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBApplicationProfileExists(accTestVpcLBHttpApplicationProfileCreateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcLBHttpApplicationProfileCreateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBHttpApplicationProfileCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "http_redirect_to", accTestVpcLBHttpApplicationProfileCreateAttributes["http_redirect_to"]),
					resource.TestCheckResourceAttr(testResourceName, "http_redirect_to_https", accTestVpcLBHttpApplicationProfileCreateAttributes["http_redirect_to_https"]),
					resource.TestCheckResourceAttr(testResourceName, "idle_timeout", accTestVpcLBHttpApplicationProfileCreateAttributes["idle_timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "request_body_size", accTestVpcLBHttpApplicationProfileCreateAttributes["request_body_size"]),
					resource.TestCheckResourceAttr(testResourceName, "request_header_size", accTestVpcLBHttpApplicationProfileCreateAttributes["request_header_size"]),
					resource.TestCheckResourceAttr(testResourceName, "response_buffering", accTestVpcLBHttpApplicationProfileCreateAttributes["response_buffering"]),
					resource.TestCheckResourceAttr(testResourceName, "response_header_size", accTestVpcLBHttpApplicationProfileCreateAttributes["response_header_size"]),
					resource.TestCheckResourceAttr(testResourceName, "response_timeout", accTestVpcLBHttpApplicationProfileCreateAttributes["response_timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "server_keep_alive", accTestVpcLBHttpApplicationProfileCreateAttributes["server_keep_alive"]),
					resource.TestCheckResourceAttr(testResourceName, "x_forwarded_for", accTestVpcLBHttpApplicationProfileCreateAttributes["x_forwarded_for"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBHttpApplicationProfileTemplate(false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBApplicationProfileExists(accTestVpcLBHttpApplicationProfileUpdateAttributes["display_name"], testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", accTestVpcLBHttpApplicationProfileUpdateAttributes["display_name"]),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBHttpApplicationProfileUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "http_redirect_to", accTestVpcLBHttpApplicationProfileUpdateAttributes["http_redirect_to"]),
					resource.TestCheckResourceAttr(testResourceName, "http_redirect_to_https", accTestVpcLBHttpApplicationProfileUpdateAttributes["http_redirect_to_https"]),
					resource.TestCheckResourceAttr(testResourceName, "idle_timeout", accTestVpcLBHttpApplicationProfileUpdateAttributes["idle_timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "request_body_size", accTestVpcLBHttpApplicationProfileUpdateAttributes["request_body_size"]),
					resource.TestCheckResourceAttr(testResourceName, "request_header_size", accTestVpcLBHttpApplicationProfileUpdateAttributes["request_header_size"]),
					resource.TestCheckResourceAttr(testResourceName, "response_buffering", accTestVpcLBHttpApplicationProfileUpdateAttributes["response_buffering"]),
					resource.TestCheckResourceAttr(testResourceName, "response_header_size", accTestVpcLBHttpApplicationProfileUpdateAttributes["response_header_size"]),
					resource.TestCheckResourceAttr(testResourceName, "response_timeout", accTestVpcLBHttpApplicationProfileUpdateAttributes["response_timeout"]),
					resource.TestCheckResourceAttr(testResourceName, "server_keep_alive", accTestVpcLBHttpApplicationProfileUpdateAttributes["server_keep_alive"]),
					resource.TestCheckResourceAttr(testResourceName, "x_forwarded_for", accTestVpcLBHttpApplicationProfileUpdateAttributes["x_forwarded_for"]),

					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBHttpApplicationProfileMinimalistic(),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBApplicationProfileExists(accTestVpcLBHttpApplicationProfileCreateAttributes["display_name"], testResourceName),
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

func TestAccResourceNsxtVpcLBHttpApplicationProfile_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_lb_http_application_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccOnlyVPC(t); testAccNSXVersion(t, "9.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBHttpApplicationProfileCheckDestroy(state, name)
		},
		Steps: withImportIdempotencyChecks(withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBHttpApplicationProfileMinimalistic(),
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

func testAccNsxtVpcLBApplicationProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy LB Application Profile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("Policy LB Application Profile resource ID not set in resources")
		}

		exists, err := resourceNsxtVpcLBAppProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Policy LB Application Profile %s does not exist", resourceID)
		}

		return nil
	}
}

func testAccNsxtVpcLBHttpApplicationProfileCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_vpc_lb_http_application_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcLBAppProfileExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("Policy LBHttpProfile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtVpcLBHttpApplicationProfileTemplate(createFlow bool) string {
	var attrMap map[string]string
	if createFlow {
		attrMap = accTestVpcLBHttpApplicationProfileCreateAttributes
	} else {
		attrMap = accTestVpcLBHttpApplicationProfileUpdateAttributes
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_http_application_profile" "test" {
%s
  display_name           = "%s"
  description            = "%s"
  http_redirect_to       = "%s"
  http_redirect_to_https = %s
  idle_timeout           = %s
  request_body_size      = %s
  request_header_size    = %s
  response_buffering     = %s
  response_header_size   = %s
  response_timeout       = %s
  server_keep_alive      = %s
  x_forwarded_for        = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), attrMap["display_name"], attrMap["description"], attrMap["http_redirect_to"], attrMap["http_redirect_to_https"], attrMap["idle_timeout"], attrMap["request_body_size"], attrMap["request_header_size"], attrMap["response_buffering"], attrMap["response_header_size"], attrMap["response_timeout"], attrMap["server_keep_alive"], attrMap["x_forwarded_for"])
}

func testAccNsxtVpcLBHttpApplicationProfileMinimalistic() string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_http_application_profile" "test" {
%s
  display_name = "%s"

}`, testAccNsxtPolicyMultitenancyContext(), accTestVpcLBHttpApplicationProfileUpdateAttributes["display_name"])
}
