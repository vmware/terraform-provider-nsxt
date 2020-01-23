/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"testing"
)

func TestAccResourceNsxtLbHttpRequestRewriteRule_basic(t *testing.T) {
	name := "test"
	fullName := "nsxt_lb_http_request_rewrite_rule.test"
	matchStrategy := "ALL"
	updatedMatchStrategy := "ANY"
	matchType := "STARTS_WITH"
	updatedMatchType := "ENDS_WITH"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLbHTTPRequestRewriteRuleCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXLbHTTPRequestRewriteRuleCreateTemplate(matchStrategy, matchType, "true", "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLbHTTPRequestRewriteRuleExists(name, fullName),
					resource.TestCheckResourceAttr(fullName, "display_name", "test"),
					resource.TestCheckResourceAttr(fullName, "description", "test description"),
					resource.TestCheckResourceAttr(fullName, "match_strategy", matchStrategy),
					testLbRuleConditionCount(fullName, "body", 2),
					testLbRuleConditionAttr(fullName, "body", "468400354", "value", "VALUE1"),
					testLbRuleConditionAttr(fullName, "body", "468400354", "match_type", "EQUALS"),
					testLbRuleConditionAttr(fullName, "body", "468400354", "case_sensitive", "false"),
					testLbRuleConditionAttr(fullName, "body", "468400354", "inverse", "false"),
					testLbRuleConditionAttr(fullName, "body", "2026604645", "value", "VALUE2"),
					testLbRuleConditionAttr(fullName, "body", "2026604645", "match_type", matchType),
					testLbRuleConditionAttr(fullName, "body", "2026604645", "case_sensitive", "true"),
					testLbRuleConditionAttr(fullName, "body", "2026604645", "inverse", "false"),
					testLbRuleConditionCount(fullName, "header", 2),
					testLbRuleConditionAttr(fullName, "header", "3951969552", "name", "NAME1"),
					testLbRuleConditionAttr(fullName, "header", "3951969552", "value", "VALUE1"),
					testLbRuleConditionAttr(fullName, "header", "3951969552", "match_type", "EQUALS"),
					testLbRuleConditionAttr(fullName, "header", "3951969552", "case_sensitive", "false"),
					testLbRuleConditionAttr(fullName, "header", "3951969552", "inverse", "false"),
					testLbRuleConditionAttr(fullName, "header", "1621991098", "name", "NAME2"),
					testLbRuleConditionAttr(fullName, "header", "1621991098", "value", "VALUE2"),
					testLbRuleConditionAttr(fullName, "header", "1621991098", "match_type", matchType),
					testLbRuleConditionAttr(fullName, "header", "1621991098", "case_sensitive", "true"),
					testLbRuleConditionAttr(fullName, "header", "1621991098", "inverse", "false"),
					testLbRuleConditionCount(fullName, "cookie", 1),
					testLbRuleConditionAttr(fullName, "cookie", "1526572800", "name", "NAME1"),
					testLbRuleConditionAttr(fullName, "cookie", "1526572800", "value", "VALUE1"),
					testLbRuleConditionAttr(fullName, "cookie", "1526572800", "match_type", matchType),
					testLbRuleConditionAttr(fullName, "cookie", "1526572800", "case_sensitive", "true"),
					testLbRuleConditionAttr(fullName, "cookie", "1526572800", "inverse", "false"),
					testLbRuleConditionCount(fullName, "method", 1),
					testLbRuleConditionAttr(fullName, "method", "2121111353", "method", "HEAD"),
					testLbRuleConditionAttr(fullName, "method", "2121111353", "inverse", "false"),
					testLbRuleConditionCount(fullName, "version", 1),
					testLbRuleConditionAttr(fullName, "version", "675459547", "version", "HTTP_VERSION_1_1"),
					testLbRuleConditionAttr(fullName, "version", "675459547", "inverse", "false"),
					testLbRuleConditionCount(fullName, "uri", 1),
					testLbRuleConditionAttr(fullName, "uri", "2025832401", "uri", "/hello"),
					testLbRuleConditionAttr(fullName, "uri", "2025832401", "match_type", matchType),
					testLbRuleConditionAttr(fullName, "uri", "2025832401", "inverse", "false"),
					testLbRuleConditionCount(fullName, "uri_arguments", 1),
					testLbRuleConditionAttr(fullName, "uri_arguments", "3495637703", "uri_arguments", "hello"),
					testLbRuleConditionAttr(fullName, "uri_arguments", "3495637703", "match_type", matchType),
					testLbRuleConditionAttr(fullName, "uri_arguments", "3495637703", "inverse", "false"),
					testLbRuleConditionCount(fullName, "ip", 1),
					testLbRuleConditionAttr(fullName, "ip", "3204467959", "source_address", "1.1.1.1"),
					testLbRuleConditionAttr(fullName, "ip", "3204467959", "inverse", "false"),
					testLbRuleConditionCount(fullName, "tcp", 1),
					testLbRuleConditionAttr(fullName, "tcp", "494891186", "source_port", "7887"),
					testLbRuleConditionAttr(fullName, "tcp", "494891186", "inverse", "false"),

					resource.TestCheckResourceAttr(fullName, "header_rewrite_action.#", "1"),
					resource.TestCheckResourceAttr(fullName, "header_rewrite_action.1966732178.name", "NAME1"),
					resource.TestCheckResourceAttr(fullName, "header_rewrite_action.1966732178.value", "VALUE1"),
					resource.TestCheckResourceAttr(fullName, "uri_rewrite_action.#", "1"),
					resource.TestCheckResourceAttr(fullName, "uri_rewrite_action.2605625851.uri", "/hello"),
					resource.TestCheckResourceAttr(fullName, "uri_rewrite_action.2605625851.uri_arguments", "hello"),
					resource.TestCheckResourceAttr(fullName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNSXLbHTTPRequestRewriteRuleCreateTemplate(updatedMatchStrategy, updatedMatchType, "false", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLbHTTPRequestRewriteRuleExists(name, fullName),
					resource.TestCheckResourceAttr(fullName, "display_name", "test"),
					resource.TestCheckResourceAttr(fullName, "description", "test description"),
					resource.TestCheckResourceAttr(fullName, "match_strategy", updatedMatchStrategy),
					testLbRuleConditionCount(fullName, "body", 2),
					testLbRuleConditionAttr(fullName, "body", "468400354", "value", "VALUE1"),
					testLbRuleConditionAttr(fullName, "body", "468400354", "match_type", "EQUALS"),
					testLbRuleConditionAttr(fullName, "body", "468400354", "case_sensitive", "false"),
					testLbRuleConditionAttr(fullName, "body", "468400354", "inverse", "false"),
					testLbRuleConditionAttr(fullName, "body", "703991590", "value", "VALUE2"),
					testLbRuleConditionAttr(fullName, "body", "703991590", "match_type", updatedMatchType),
					testLbRuleConditionAttr(fullName, "body", "703991590", "case_sensitive", "false"),
					testLbRuleConditionAttr(fullName, "body", "703991590", "inverse", "true"),
					testLbRuleConditionCount(fullName, "header", 2),
					testLbRuleConditionAttr(fullName, "header", "3951969552", "name", "NAME1"),
					testLbRuleConditionAttr(fullName, "header", "3951969552", "value", "VALUE1"),
					testLbRuleConditionAttr(fullName, "header", "3951969552", "match_type", "EQUALS"),
					testLbRuleConditionAttr(fullName, "header", "3951969552", "case_sensitive", "false"),
					testLbRuleConditionAttr(fullName, "header", "3951969552", "inverse", "false"),
					testLbRuleConditionAttr(fullName, "header", "1831610353", "name", "NAME2"),
					testLbRuleConditionAttr(fullName, "header", "1831610353", "value", "VALUE2"),
					testLbRuleConditionAttr(fullName, "header", "1831610353", "match_type", updatedMatchType),
					testLbRuleConditionAttr(fullName, "header", "1831610353", "case_sensitive", "false"),
					testLbRuleConditionAttr(fullName, "header", "1831610353", "inverse", "true"),
					testLbRuleConditionCount(fullName, "cookie", 1),
					testLbRuleConditionAttr(fullName, "cookie", "1467752011", "name", "NAME1"),
					testLbRuleConditionAttr(fullName, "cookie", "1467752011", "value", "VALUE1"),
					testLbRuleConditionAttr(fullName, "cookie", "1467752011", "match_type", updatedMatchType),
					testLbRuleConditionAttr(fullName, "cookie", "1467752011", "case_sensitive", "false"),
					testLbRuleConditionAttr(fullName, "cookie", "1467752011", "inverse", "true"),
					testLbRuleConditionCount(fullName, "method", 1),
					testLbRuleConditionAttr(fullName, "method", "3814880847", "method", "HEAD"),
					testLbRuleConditionAttr(fullName, "method", "3814880847", "inverse", "true"),
					testLbRuleConditionCount(fullName, "version", 1),
					testLbRuleConditionAttr(fullName, "version", "1187949210", "version", "HTTP_VERSION_1_1"),
					testLbRuleConditionAttr(fullName, "version", "1187949210", "inverse", "true"),
					testLbRuleConditionCount(fullName, "uri", 1),
					testLbRuleConditionAttr(fullName, "uri", "3068503099", "uri", "/hello"),
					testLbRuleConditionAttr(fullName, "uri", "3068503099", "match_type", updatedMatchType),
					testLbRuleConditionAttr(fullName, "uri", "3068503099", "inverse", "true"),
					testLbRuleConditionCount(fullName, "uri_arguments", 1),
					testLbRuleConditionAttr(fullName, "uri_arguments", "4074102258", "uri_arguments", "hello"),
					testLbRuleConditionAttr(fullName, "uri_arguments", "4074102258", "match_type", updatedMatchType),
					testLbRuleConditionAttr(fullName, "uri_arguments", "4074102258", "inverse", "true"),
					testLbRuleConditionCount(fullName, "ip", 1),
					testLbRuleConditionAttr(fullName, "ip", "445373689", "source_address", "1.1.1.1"),
					testLbRuleConditionAttr(fullName, "ip", "445373689", "inverse", "true"),
					testLbRuleConditionCount(fullName, "tcp", 1),
					testLbRuleConditionAttr(fullName, "tcp", "3399348458", "source_port", "7887"),
					testLbRuleConditionAttr(fullName, "tcp", "3399348458", "inverse", "true"),

					resource.TestCheckResourceAttr(fullName, "header_rewrite_action.#", "1"),
					resource.TestCheckResourceAttr(fullName, "header_rewrite_action.1966732178.name", "NAME1"),
					resource.TestCheckResourceAttr(fullName, "header_rewrite_action.1966732178.value", "VALUE1"),
					resource.TestCheckResourceAttr(fullName, "uri_rewrite_action.#", "1"),
					resource.TestCheckResourceAttr(fullName, "uri_rewrite_action.2605625851.uri", "/hello"),
					resource.TestCheckResourceAttr(fullName, "uri_rewrite_action.2605625851.uri_arguments", "hello"),
					resource.TestCheckResourceAttr(fullName, "tag.#", "1"),
				),
			},
		},
	})
}

func testLbRuleConditionCount(resourceName string, conditionType string, count int) resource.TestCheckFunc {

	return resource.TestCheckResourceAttr(resourceName,
		fmt.Sprintf("%s_condition.#", conditionType),
		fmt.Sprintf("%d", count))
}

func testLbRuleConditionAttr(resourceName string, conditionType string, setHash string, attrName string, attrValue string) resource.TestCheckFunc {

	return resource.TestCheckResourceAttr(resourceName,
		fmt.Sprintf("%s_condition.%s.%s", conditionType, setHash, attrName),
		attrValue)
}

func TestAccResourceNsxtLbHttpRequestRewriteRule_importBasic(t *testing.T) {
	name := "test"
	resourceName := "nsxt_lb_http_request_rewrite_rule.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLbHTTPRequestRewriteRuleCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXLbHTTPRequestRewriteRuleCreateTemplateTrivial(),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNSXLbHTTPRequestRewriteRuleExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		nsxClient := testAccProvider.Meta().(nsxtClients).NsxtClient
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX LB rule resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX LB rule resource ID not set in resources ")
		}

		monitor, responseCode, err := nsxClient.ServicesApi.ReadLoadBalancerRule(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving LB rule with ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if LB rule %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		if displayName == monitor.DisplayName {
			return nil
		}
		return fmt.Errorf("NSX LB rule %s wasn't found", displayName)
	}
}

func testAccNSXLbHTTPRequestRewriteRuleCheckDestroy(state *terraform.State, displayName string) error {
	nsxClient := testAccProvider.Meta().(nsxtClients).NsxtClient
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_lb_http_request_rewrite_rule" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		monitor, responseCode, err := nsxClient.ServicesApi.ReadLoadBalancerRule(nsxClient.Context, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving LB rule with ID %s. Error: %v", resourceID, err)
		}

		if displayName == monitor.DisplayName {
			return fmt.Errorf("NSX LB rule %s still exists", displayName)
		}
	}
	return nil
}

func testAccNSXLbRuleNameValueConditionTemplate(conditionType string, matchType string, caseSensitive string, inverse string, num int) string {
	return fmt.Sprintf(`
    %s_condition {
      name           = "NAME%d"
      value          = "VALUE%d"
      match_type     = "%s"
      case_sensitive = %s
      inverse        = %s
    }
    `, conditionType, num, num, matchType, caseSensitive, inverse)
}

func testAccNSXLbRuleValueConditionTemplate(conditionType string, valueType string, matchType string, caseSensitive string, inverse string, num int) string {
	return fmt.Sprintf(`
    %s_condition {
      %s             = "VALUE%d"
      match_type     = "%s"
      case_sensitive = %s
      inverse        = %s
    }
    `, conditionType, valueType, num, matchType, caseSensitive, inverse)
}

func testAccNSXLbRuleSimpleConditionTemplate(conditionType string, value string, inverse string) string {
	return fmt.Sprintf(`
    %s_condition {
      %s      = "%s"
      inverse = %s
    }
    `, conditionType, conditionType, value, inverse)
}

func testAccNSXLbRuleSimpleMatchConditionTemplate(conditionType string, value string, matchType string, inverse string) string {
	return fmt.Sprintf(`
    %s_condition {
      %s         = "%s"
      match_type = "%s"
      inverse    = %s
    }
    `, conditionType, conditionType, value, matchType, inverse)
}

func testAccNSXLbHTTPRequestRewriteRuleCreateTemplate(matchStrategy string, matchType string, caseSensitive string, inverse string) string {
	return fmt.Sprintf(`
resource "nsxt_lb_http_request_rewrite_rule" "test" {
  display_name   = "test"
  description    = "test description"
  match_strategy = "%s"
`, matchStrategy) +
		testAccNSXLbRuleValueConditionTemplate("body", "value", "EQUALS", "false", "false", 1) +
		testAccNSXLbRuleValueConditionTemplate("body", "value", matchType, caseSensitive, inverse, 2) +
		testAccNSXLbRuleNameValueConditionTemplate("header", "EQUALS", "false", "false", 1) +
		testAccNSXLbRuleNameValueConditionTemplate("header", matchType, caseSensitive, inverse, 2) +
		testAccNSXLbRuleNameValueConditionTemplate("cookie", matchType, caseSensitive, inverse, 1) +
		testAccNSXLbRuleSimpleConditionTemplate("method", "HEAD", inverse) +
		testAccNSXLbRuleSimpleConditionTemplate("version", "HTTP_VERSION_1_1", inverse) +
		testAccNSXLbRuleSimpleMatchConditionTemplate("uri", "/hello", matchType, inverse) +
		testAccNSXLbRuleSimpleMatchConditionTemplate("uri_arguments", "hello", matchType, inverse) +
		fmt.Sprintf(`
  ip_condition {
    source_address = "1.1.1.1"
    inverse        = %s
  }

  tcp_condition {
    source_port = "7887"
    inverse     = %s
  }

  uri_rewrite_action {
    uri           = "/hello"
    uri_arguments = "hello"
  }

  header_rewrite_action {
    name  = "NAME1"
    value = "VALUE1"
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
`, inverse, inverse)
}

func testAccNSXLbHTTPRequestRewriteRuleCreateTemplateTrivial() string {
	return `
resource "nsxt_lb_http_request_rewrite_rule" "test" {
  description    = "test description"
  match_strategy = "ALL"

  method_condition {
    method = "POST"
  }

  header_rewrite_action {
    name  = "NAME1"
    value = "VALUE1"
  }
}
`
}
