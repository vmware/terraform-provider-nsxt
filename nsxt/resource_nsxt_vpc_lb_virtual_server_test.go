// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var accTestVpcLBVirtualServerCreateAttributes = map[string]string{
	"description":                "terraform created",
	"access_log_enabled":         "true",
	"default_pool_member_ports":  "777-778",
	"enabled":                    "true",
	"max_concurrent_connections": "2",
	"max_new_connection_rate":    "2",
	"ports":                      "80-81",
	"log_significant_event_only": "true",
}

var accTestVpcLBVirtualServerUpdateAttributes = map[string]string{
	"description":                "terraform updated",
	"access_log_enabled":         "false",
	"default_pool_member_ports":  "555",
	"enabled":                    "false",
	"max_concurrent_connections": "5",
	"max_new_connection_rate":    "5",
	"ports":                      "443",
	"log_significant_event_only": "false",
}

func TestAccResourceNsxtVpcLBVirtualServer_basic(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_virtual_server.test"
	createName := getAccTestResourceName()
	updateName := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBVirtualServerCheckDestroy(state, updateName)
		},
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBVirtualServerTemplate(true, createName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", createName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBVirtualServerCreateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "access_log_enabled", accTestVpcLBVirtualServerCreateAttributes["access_log_enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "default_pool_member_ports.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "default_pool_member_ports.0", accTestVpcLBVirtualServerCreateAttributes["default_pool_member_ports"]),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestVpcLBVirtualServerCreateAttributes["enabled"]),
					resource.TestCheckResourceAttrPair(testResourceName, "ip_address", "nsxt_vpc_ip_address_allocation.vip1", "allocation_ips"),
					resource.TestCheckResourceAttr(testResourceName, "max_concurrent_connections", accTestVpcLBVirtualServerCreateAttributes["max_concurrent_connections"]),
					resource.TestCheckResourceAttr(testResourceName, "max_new_connection_rate", accTestVpcLBVirtualServerCreateAttributes["max_new_connection_rate"]),
					resource.TestCheckResourceAttr(testResourceName, "ports.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "ports.0", accTestVpcLBVirtualServerCreateAttributes["ports"]),
					resource.TestCheckResourceAttr(testResourceName, "log_significant_event_only", accTestVpcLBVirtualServerCreateAttributes["log_significant_event_only"]),
					resource.TestCheckResourceAttrSet(testResourceName, "application_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "persistence_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "sorry_pool_path"),
					resource.TestCheckResourceAttr(testResourceName, "context.0.project_id", os.Getenv("NSXT_VPC_PROJECT_ID")),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBVirtualServerTemplate(false, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", accTestVpcLBVirtualServerUpdateAttributes["description"]),
					resource.TestCheckResourceAttr(testResourceName, "access_log_enabled", accTestVpcLBVirtualServerUpdateAttributes["access_log_enabled"]),
					resource.TestCheckResourceAttr(testResourceName, "default_pool_member_ports.0", accTestVpcLBVirtualServerUpdateAttributes["default_pool_member_ports"]),
					resource.TestCheckResourceAttr(testResourceName, "enabled", accTestVpcLBVirtualServerUpdateAttributes["enabled"]),
					resource.TestCheckResourceAttrPair(testResourceName, "ip_address", "nsxt_vpc_ip_address_allocation.vip2", "allocation_ips"),
					resource.TestCheckResourceAttr(testResourceName, "max_concurrent_connections", accTestVpcLBVirtualServerUpdateAttributes["max_concurrent_connections"]),
					resource.TestCheckResourceAttr(testResourceName, "max_new_connection_rate", accTestVpcLBVirtualServerUpdateAttributes["max_new_connection_rate"]),
					resource.TestCheckResourceAttr(testResourceName, "ports.0", accTestVpcLBVirtualServerUpdateAttributes["ports"]),
					resource.TestCheckResourceAttr(testResourceName, "log_significant_event_only", accTestVpcLBVirtualServerUpdateAttributes["log_significant_event_only"]),
					resource.TestCheckResourceAttrSet(testResourceName, "application_profile_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "pool_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNsxtVpcLBVirtualServerMinimalistic(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updateName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttr(testResourceName, "pool_path", ""),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
				),
			},
		}),
	})
}

func TestAccResourceNsxtVpcLBVirtualServer_withRules(t *testing.T) {
	testResourceName := "nsxt_vpc_lb_virtual_server.test"
	name := getAccTestResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBVirtualServerCheckDestroy(state, name)
		},
		Steps: withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBVirtualServerWithRules(name),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.display_name", "select_pool_test"),
					resource.TestCheckResourceAttr(testResourceName, "rule.0.action.0.select_pool.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "rule.0.action.0.select_pool.0.pool_id"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.display_name", "http_redirect_test"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.action.0.http_redirect.0.redirect_status", "301"),
					resource.TestCheckResourceAttr(testResourceName, "rule.1.condition.0.http_request_body.0.body_value", "xyz"),
				),
			},
			{
				Config: testAccNsxtVpcLBVirtualServerHttp(name, false),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtVpcLBVirtualServerExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "rule.#", "0"),
				),
			},
		}),
	})
}

func TestAccResourceNsxtVpcLBVirtualServer_importBasic(t *testing.T) {
	name := getAccTestResourceName()
	testResourceName := "nsxt_vpc_lb_virtual_server.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccOnlyVPC(t)
			testAccNSXVersion(t, "9.0.0")
		},
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtVpcLBVirtualServerCheckDestroy(state, name)
		},
		Steps: withImportIdempotencyChecks(withIdempotencyChecks([]resource.TestStep{
			{
				Config: testAccNsxtVpcLBVirtualServerMinimalistic(name),
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

func testAccNsxtVpcLBVirtualServerExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("VPC LBVirtualServer resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("VPC LBVirtualServer resource ID not set in resources")
		}

		exists, err := resourceNsxtVpcLBVirtualServerExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("VPC LBVirtualServer %s does not exist", resourceID)
		}
		return nil
	}
}

func testAccNsxtVpcLBVirtualServerCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nsxt_vpc_lb_virtual_server" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		exists, err := resourceNsxtVpcLBVirtualServerExists(testAccGetSessionContext(), resourceID, connector)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("VPC LBVirtualServer %s still exists", displayName)
		}
	}
	return nil
}

// The virtual server IP address has to be reserved in the VPC, so allocate it.
func testAccNsxtVpcLBVirtualServerVips(name string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_ip_address_allocation" "vip1" {
  %[1]s
  display_name    = "%[2]s-vip1"
  allocation_size = 1
}

resource "nsxt_vpc_ip_address_allocation" "vip2" {
  %[1]s
  display_name    = "%[2]s-vip2"
  allocation_size = 1
}
`, testAccNsxtPolicyMultitenancyContext(), name)
}

// testAccNsxtVpcLBVirtualServerDeps creates the VPC LB objects the virtual server refers to.
func testAccNsxtVpcLBVirtualServerDeps(name string) string {
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_fast_tcp_application_profile" "dep" {
  %[1]s
  display_name = "%[2]s-app"
}

resource "nsxt_vpc_lb_source_ip_persistence_profile" "dep" {
  %[1]s
  display_name = "%[2]s-persistence"
}

resource "nsxt_vpc_lb_pool" "pool" {
  %[1]s
  display_name = "%[2]s-pool"
}

resource "nsxt_vpc_lb_pool" "sorry_pool" {
  %[1]s
  display_name = "%[2]s-sorry-pool"
}
`, testAccNsxtPolicyMultitenancyContext(), name) + testAccNsxtVpcLBVirtualServerVips(name)
}

func testAccNsxtVpcLBVirtualServerTemplate(createFlow bool, name string) string {
	var attrMap map[string]string
	var vip string
	if createFlow {
		attrMap = accTestVpcLBVirtualServerCreateAttributes
		vip = "vip1"
	} else {
		attrMap = accTestVpcLBVirtualServerUpdateAttributes
		vip = "vip2"
	}
	return testAccNsxtVpcLBVirtualServerDeps(name) + fmt.Sprintf(`
resource "nsxt_vpc_lb_virtual_server" "test" {
  %s
  display_name               = "%s"
  description                = "%s"
  access_log_enabled         = %s
  application_profile_path   = nsxt_vpc_lb_fast_tcp_application_profile.dep.path
  default_pool_member_ports  = ["%s"]
  enabled                    = %s
  ip_address                 = nsxt_vpc_ip_address_allocation.%s.allocation_ips
  max_concurrent_connections = %s
  max_new_connection_rate    = %s
  ports                      = ["%s"]
  log_significant_event_only = %s
  pool_path                  = nsxt_vpc_lb_pool.pool.path
  sorry_pool_path            = nsxt_vpc_lb_pool.sorry_pool.path
  persistence_profile_path   = nsxt_vpc_lb_source_ip_persistence_profile.dep.path

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, testAccNsxtPolicyMultitenancyContext(), name, attrMap["description"], attrMap["access_log_enabled"], attrMap["default_pool_member_ports"], attrMap["enabled"], vip, attrMap["max_concurrent_connections"], attrMap["max_new_connection_rate"], attrMap["ports"], attrMap["log_significant_event_only"])
}

func testAccNsxtVpcLBVirtualServerMinimalistic(name string) string {
	return testAccNsxtVpcLBVirtualServerDeps(name) + fmt.Sprintf(`
resource "nsxt_vpc_lb_virtual_server" "test" {
  %s
  display_name             = "%s"
  application_profile_path = nsxt_vpc_lb_fast_tcp_application_profile.dep.path
  ip_address               = nsxt_vpc_ip_address_allocation.vip1.allocation_ips
  ports                    = ["%s"]
}`, testAccNsxtPolicyMultitenancyContext(), name, accTestVpcLBVirtualServerUpdateAttributes["ports"])
}

func testAccNsxtVpcLBVirtualServerWithRules(name string) string {
	return testAccNsxtVpcLBVirtualServerHttp(name, true)
}

// testAccNsxtVpcLBVirtualServerHttp keeps the HTTP application profile in the config whether or
// not the rules are set, so the profile isn't deleted while the virtual server still refers to it.
func testAccNsxtVpcLBVirtualServerHttp(name string, withRules bool) string {
	rules := `
  rule {
    display_name = "select_pool_test"
    action {
      select_pool {
        pool_id = nsxt_vpc_lb_pool.pool.path
      }
    }
  }

  rule {
    display_name = "http_redirect_test"
    action {
      http_redirect {
        redirect_status = "301"
        redirect_url    = "dummy"
      }
    }
    condition {
      http_request_body {
        body_value = "xyz"
        match_type = "REGEX"
      }
    }
  }
`
	if !withRules {
		rules = ""
	}
	return fmt.Sprintf(`
resource "nsxt_vpc_lb_http_application_profile" "dep" {
  %[1]s
  display_name = "%[2]s-http-app"
}

resource "nsxt_vpc_lb_pool" "pool" {
  %[1]s
  display_name = "%[2]s-pool"
}

resource "nsxt_vpc_lb_virtual_server" "test" {
  %[1]s
  display_name             = "%[2]s"
  application_profile_path = nsxt_vpc_lb_http_application_profile.dep.path
  ip_address               = nsxt_vpc_ip_address_allocation.vip1.allocation_ips
  ports                    = ["80"]
  pool_path                = nsxt_vpc_lb_pool.pool.path
%[3]s
}`, testAccNsxtPolicyMultitenancyContext(), name, rules) + testAccNsxtVpcLBVirtualServerVips(name)
}
